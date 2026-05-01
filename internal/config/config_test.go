package config

import "testing"

// TestLoadBuildsDSNFromParts verifies that Load() composes a valid DSN
// from MYSQL_* parts when MYSQL_DSN is intentionally left empty.
func TestLoadBuildsDSNFromParts(t *testing.T) {
	t.Setenv("MYSQL_DSN", "")
	t.Setenv("MYSQL_USER", "hotel")
	t.Setenv("MYSQL_PASSWORD", "secret")
	t.Setenv("MYSQL_HOST", "db")
	t.Setenv("MYSQL_PORT", "3307")
	t.Setenv("MYSQL_DATABASE", "resort")

	cfg := Load()

	want := "hotel:secret@tcp(db:3307)/resort?parseTime=true&charset=utf8mb4"
	if cfg.MySQLDSN != want {
		t.Fatalf("unexpected dsn: got %q want %q", cfg.MySQLDSN, want)
	}
}

// TestLoadPrefersExplicitDSN verifies that a full MYSQL_DSN value has priority
// over part-based MYSQL_* values.
func TestLoadPrefersExplicitDSN(t *testing.T) {
	t.Setenv("MYSQL_DSN", "custom-dsn")
	t.Setenv("MYSQL_HOST", "ignored")

	cfg := Load()

	if cfg.MySQLDSN != "custom-dsn" {
		t.Fatalf("expected explicit dsn to win, got %q", cfg.MySQLDSN)
	}
}
