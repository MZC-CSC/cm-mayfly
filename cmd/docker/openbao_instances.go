package docker

// Additional OpenBao instances in `infra info`.
//
// The stack used to carry exactly one OpenBao — the one backing cb-tumblebug's
// credential store — and internal/openbao assumes it throughout: the address is
// a constant, the token lives in .env under VAULT_TOKEN, and the unseal key
// sits under conf/docker/data/openbao/secrets. printOpenbaoInfo reports that
// instance, including whether OpenBao and .env still agree.
//
// cm-honeybee changed that. It keeps source-side CSP credentials and SSH
// secrets in an OpenBao of its own (service openbao-honeybee) and manages that
// instance itself — initializing, unsealing and enabling KV without help from
// mayfly. There is no token of ours to compare and no unseal key of ours to
// read, so the consistency verdict does not apply. What remains worth showing
// is whether the instance is up, initialized and unsealed: while it is sealed
// every source registration fails, and without a line in `infra info` that is
// only visible in the container logs.

import (
	"encoding/json"
	"fmt"
	"strings"
)

// openbaoInstancePrefix selects the compose services treated as OpenBao
// instances. Same rule the service categorizer (categorizeService) and the
// --clean-db data guard use, so naming an instance openbao-* is all it takes to
// be picked up everywhere.
const openbaoInstancePrefix = "openbao-"

// notAnOpenbaoInstance lists openbao-* services that are not OpenBao servers.
// openbao-unseal is the sidecar watching the primary instance.
var notAnOpenbaoInstance = map[string]bool{
	"openbao-unseal": true,
}

// openbaoInstanceStatus is what can be said about an instance we hold no token
// for.
type openbaoInstanceStatus struct {
	Service     string
	Reachable   bool
	Initialized bool
	Sealed      bool
	Detail      string // why it is unreachable, when it is
}

// baoStatusJSON is the subset of `bao status -format=json` that is read.
type baoStatusJSON struct {
	Initialized bool `json:"initialized"`
	Sealed      bool `json:"sealed"`
}

// openbaoInstanceServices returns the openbao-* services declared in the
// compose file, minus the ones that are not instances.
func openbaoInstanceServices() []string {
	var out []string
	for _, name := range getServicesFromCompose() {
		if !strings.HasPrefix(name, openbaoInstancePrefix) || notAnOpenbaoInstance[name] {
			continue
		}
		out = append(out, name)
	}
	return out
}

// queryOpenbaoInstance asks one instance for its seal status.
//
// It goes through composeOutput — `docker compose -f <file> exec` with
// COMPOSE_PROJECT_NAME set — rather than `docker exec <container>`. The compose
// file pins container_name on every service, so those names are global to the
// host: `docker exec openbao-honeybee …` would reach whatever container carries
// that name, including one from another project. Resolving the service within
// our project keeps the query inside the stack the user asked about, and it
// follows a -p override.
//
// Querying through the container also avoids publishing a host port per
// instance. The primary OpenBao is reachable over localhost only because it
// publishes 8200, and a status line does not justify exposing another secret
// store.
func queryOpenbaoInstance(service string) openbaoInstanceStatus {
	st := openbaoInstanceStatus{Service: service}

	out, err := composeOutput("exec", "-T", service, "bao", "status", "-format=json")

	// `bao status` exits non-zero when sealed (1) or uninitialized (2) while
	// still printing the JSON, so a non-nil error is not by itself a failure —
	// the output decides. Only when no JSON comes back is the instance treated
	// as unreachable.
	trimmed := strings.TrimSpace(string(out))
	if i := strings.Index(trimmed, "{"); i >= 0 {
		var s baoStatusJSON
		if json.Unmarshal([]byte(trimmed[i:]), &s) == nil {
			st.Reachable = true
			st.Initialized = s.Initialized
			st.Sealed = s.Sealed
			return st
		}
	}

	st.Detail = "not running or not reachable"
	if err != nil && trimmed == "" {
		st.Detail = "not running"
	}
	return st
}

// printOpenbaoInstancesInfo appends one section per additional instance. It
// prints nothing when the stack carries none, so callers can call it
// unconditionally.
func printOpenbaoInstancesInfo() {
	for _, svc := range openbaoInstanceServices() {
		st := queryOpenbaoInstance(svc)
		fmt.Println()
		fmt.Printf("[OpenBao: %s]\n", svc)
		if !st.Reachable {
			fmt.Printf("  API        : reachable=false (%s)\n", st.Detail)
			continue
		}
		fmt.Printf("  API        : reachable=true initialized=%v sealed=%v\n",
			st.Initialized, st.Sealed)
		switch {
		case !st.Initialized:
			fmt.Println("  state      : uninitialized — the owning service initializes it on startup")
		case st.Sealed:
			fmt.Println("  state      : sealed — secrets stay unavailable until it is unsealed")
		default:
			fmt.Println("  state      : ready")
		}
	}
}
