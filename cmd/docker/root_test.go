package docker

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// The DB *user* keys must be required, not just the passwords: each postgres
// service's healthcheck is `pg_isready -U ${*_DB_USER}` and fails when the user
// is blank. Omitting them (the original bug) let a blank TUMBLEBUG_DB_USER slip
// past the pre-run check and deadlock the cb-tumblebug chain on a fresh install.
func TestRequiredEnvKeysIncludeDBUsers(t *testing.T) {
	want := []string{"TUMBLEBUG_DB_USER", "ANT_DB_USER", "BUTTERFLY_DB_USER", "AIRFLOW_DB_USER"}
	have := make(map[string]bool, len(requiredEnvKeys))
	for _, k := range requiredEnvKeys {
		have[k] = true
	}
	for _, k := range want {
		if !have[k] {
			t.Errorf("requiredEnvKeys is missing %q — a blank DB user would break its postgres healthcheck", k)
		}
	}
}

// fullValidEnv returns a .env body that sets every required key to a non-empty
// value, so tests can then blank exactly one key to assert it is flagged.
func fullValidEnv() string {
	var b strings.Builder
	b.WriteString("# fixture\n")
	for _, k := range requiredEnvKeys {
		b.WriteString(k + "=x\n")
	}
	return b.String()
}

// withEnvFixture writes body to <tmp>/.env and points DockerFilePath at
// <tmp>/docker-compose.yaml for the duration of the test.
func withEnvFixture(t *testing.T, body string) {
	t.Helper()
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, ".env"), []byte(body), 0o600); err != nil {
		t.Fatal(err)
	}
	prev := DockerFilePath
	DockerFilePath = filepath.Join(dir, "docker-compose.yaml")
	t.Cleanup(func() { DockerFilePath = prev })
}

func TestValidateDockerEnvFile(t *testing.T) {
	t.Run("all set → ok", func(t *testing.T) {
		withEnvFixture(t, fullValidEnv())
		if err := validateDockerEnvFile(); err != nil {
			t.Errorf("fully-populated .env should validate, got: %v", err)
		}
	})

	t.Run("blank TUMBLEBUG_DB_USER → flagged", func(t *testing.T) {
		body := strings.Replace(fullValidEnv(), "TUMBLEBUG_DB_USER=x\n", "TUMBLEBUG_DB_USER=\n", 1)
		withEnvFixture(t, body)
		err := validateDockerEnvFile()
		if err == nil {
			t.Fatal("a blank TUMBLEBUG_DB_USER must be reported, got nil")
		}
		if !strings.Contains(err.Error(), "TUMBLEBUG_DB_USER") {
			t.Errorf("error should name TUMBLEBUG_DB_USER, got: %v", err)
		}
	})

	t.Run("blank password still flagged (regression)", func(t *testing.T) {
		body := strings.Replace(fullValidEnv(), "SPIDER_PASSWORD=x\n", "SPIDER_PASSWORD=\n", 1)
		withEnvFixture(t, body)
		if err := validateDockerEnvFile(); err == nil || !strings.Contains(err.Error(), "SPIDER_PASSWORD") {
			t.Errorf("blank SPIDER_PASSWORD must still be flagged, got: %v", err)
		}
	})
}
