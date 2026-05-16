package config

import (
	"fmt"
	"os"
)

// Config is the normalized runtime configuration used by the application.
// Values come from environment variables with sensible local defaults.
type Config struct {
	// AppAddr is the bind address for net/http (for example ":8080").
	AppAddr string
	// MySQLDSN is the final DSN passed into the MySQL driver.
	MySQLDSN string
}

// Load reads environment variables and builds the final Config object.
// If MYSQL_DSN is not provided, it composes a DSN from MYSQL_* parts.
func Load() Config {
	// Prefer a full DSN when explicitly provided by the caller.
	mysqlDSN := os.Getenv("MYSQL_DSN")
	if mysqlDSN == "" {
		// Build DSN from separate env parts to simplify local/dev setup.
		mysqlDSN = fmt.Sprintf(
			"%s:%s@tcp(%s:%s)/%s?parseTime=true&charset=utf8mb4",
			getEnv("MYSQL_USER", "root"),
			getEnv("MYSQL_PASSWORD", "password"),
			getEnv("MYSQL_HOST", "127.0.0.1"),
			getEnv("MYSQL_PORT", "3306"),
			getEnv("MYSQL_DATABASE", "resort_clone"),
		)
	}

	// Collect all configuration values in one immutable struct.
	cfg := Config{
		AppAddr:  getEnv("APP_ADDR", ":8080"),
		MySQLDSN: mysqlDSN,
	}

	return cfg
}

// getEnv returns environment variable value by key.
// It falls back to the provided default when the key is missing or empty.
func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}

	return fallback
}
