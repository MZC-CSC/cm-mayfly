package docker

import (
	"os"
	"path/filepath"
	"testing"
)

// --clean-db keeps the shared OpenBao and drops everything else, companion
// instances included.
//
// A companion (openbao-<service>) cannot survive its owner. Its secrets are
// keyed by the owner's rows and its unseal key lives in the owner's database,
// so a preserved companion is an initialized store nobody can open — not
// initialized again because it already is, not unsealed because the key went
// with the owner. Dropping it is what lets the owner start clean.
func TestHostDataTargetsKeepsSharedOpenbaoAndDropsCompanions(t *testing.T) {
	prevDB, prevAll := cleanDBFlag, cleanAllFlag
	defer func() { cleanDBFlag, cleanAllFlag = prevDB, prevAll }()
	cleanDBFlag, cleanAllFlag = true, false

	makeDataDirs(t,
		writeTestCompose(t, "cb-spider", "openbao", "openbao-honeybee", "cm-honeybee"),
		"cb-spider", "openbao", "openbao-honeybee", "cm-honeybee")

	targets, err := hostDataTargets(nil)
	wiped := basenames(t, targets, err)

	if wiped["openbao"] {
		t.Error("--clean-db wiped the shared openbao; its credentials must be preserved")
	}
	for _, gone := range []string{"openbao-honeybee", "cm-honeybee", "cb-spider"} {
		if !wiped[gone] {
			t.Errorf("--clean-db kept %s; only the shared openbao is preserved", gone)
		}
	}
}

// --clean-all is the full reset and takes the shared instance too.
func TestHostDataTargetsCleanAllTakesEverything(t *testing.T) {
	prevDB, prevAll := cleanDBFlag, cleanAllFlag
	defer func() { cleanDBFlag, cleanAllFlag = prevDB, prevAll }()
	cleanDBFlag, cleanAllFlag = true, true

	makeDataDirs(t, writeTestCompose(t, "openbao", "openbao-honeybee"),
		"openbao", "openbao-honeybee")

	targets, err := hostDataTargets(nil)
	wiped := basenames(t, targets, err)
	for _, name := range []string{"openbao", "openbao-honeybee"} {
		if !wiped[name] {
			t.Errorf("--clean-all kept %s; a full reset removes every instance", name)
		}
	}
}

// A service-scoped --clean-db takes the service's own OpenBao with it. Removing
// cm-honeybee's data while leaving openbao-honeybee behind is the state that
// cannot recover, so the two go together.
func TestHostDataTargetsServiceScopedTakesCompanion(t *testing.T) {
	prevDB, prevAll := cleanDBFlag, cleanAllFlag
	defer func() { cleanDBFlag, cleanAllFlag = prevDB, prevAll }()
	cleanDBFlag, cleanAllFlag = true, false

	makeDataDirs(t, writeTestCompose(t, "cm-honeybee", "openbao-honeybee", "cb-spider"),
		"cm-honeybee", "openbao-honeybee", "cb-spider")

	targets, err := hostDataTargets([]string{"cm-honeybee"})
	wiped := basenames(t, targets, err)
	if !wiped["cm-honeybee"] || !wiped["openbao-honeybee"] {
		t.Errorf("hostDataTargets = %v, want cm-honeybee and its openbao-honeybee", keys(wiped))
	}
	if wiped["cb-spider"] {
		t.Error("a service-scoped wipe touched cb-spider")
	}
}

// Most services have no companion; naming one must not invent a target for a
// directory that does not exist.
func TestHostDataTargetsServiceScopedWithoutCompanion(t *testing.T) {
	prevDB, prevAll := cleanDBFlag, cleanAllFlag
	defer func() { cleanDBFlag, cleanAllFlag = prevDB, prevAll }()
	cleanDBFlag, cleanAllFlag = true, false

	makeDataDirs(t, writeTestCompose(t, "cb-spider"), "cb-spider")

	targets, err := hostDataTargets([]string{"cb-spider"})
	if err != nil {
		t.Fatalf("hostDataTargets: %v", err)
	}
	if len(targets) != 1 || filepath.Base(targets[0]) != "cb-spider" {
		t.Fatalf("hostDataTargets = %v, want only the cb-spider directory", targets)
	}
}

// The instance list drives what `infra info` reports: the openbao-* services,
// minus the unseal sidecar, which is not an OpenBao server. The shared instance
// is reported separately, with its consistency verdict.
func TestOpenbaoInstanceServices(t *testing.T) {
	writeTestCompose(t, "cb-spider", "openbao", "openbao-unseal", "openbao-honeybee", "cm-honeybee")

	got := openbaoInstanceServices()
	if len(got) != 1 || got[0] != "openbao-honeybee" {
		t.Fatalf("openbaoInstanceServices = %v, want [openbao-honeybee]", got)
	}
}

// A stack with no companion yields an empty list, so `infra info` prints no
// extra section.
func TestOpenbaoInstanceServicesNone(t *testing.T) {
	writeTestCompose(t, "cb-spider", "openbao", "openbao-unseal")

	if got := openbaoInstanceServices(); len(got) != 0 {
		t.Fatalf("openbaoInstanceServices = %v, want none", got)
	}
}

// makeDataDirs creates conf/docker/data/<name>/ next to the test compose file.
func makeDataDirs(t *testing.T, composePath string, names ...string) string {
	t.Helper()
	dataRoot := filepath.Join(filepath.Dir(composePath), hostDataDirName)
	for _, n := range names {
		if err := os.MkdirAll(filepath.Join(dataRoot, n), 0o750); err != nil {
			t.Fatal(err)
		}
	}
	return dataRoot
}

func basenames(t *testing.T, targets []string, err error) map[string]bool {
	t.Helper()
	if err != nil {
		t.Fatalf("hostDataTargets: %v", err)
	}
	out := map[string]bool{}
	for _, p := range targets {
		out[filepath.Base(p)] = true
	}
	return out
}

func keys(m map[string]bool) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	return out
}
