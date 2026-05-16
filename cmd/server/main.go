package main

import (
	"context"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"

	"hotel_site/internal/config"
	httpapp "hotel_site/internal/http"
	"hotel_site/internal/repository"
)

// main composes the whole application lifecycle:
// 1) load runtime configuration,
// 2) open and validate database connectivity,
// 3) ensure schema and default data exist,
// 4) build HTTP handlers,
// 5) start the HTTP server and block forever.
func main() {
	// Load process-level settings (bind address and database DSN).
	cfg := config.Load()
	// Use a root background context for startup operations.
	ctx := context.Background()

	// Create a MySQL connection pool and verify connectivity with a startup ping.
	db, err := repository.OpenMySQL(ctx, cfg.MySQLDSN)
	if err != nil {
		log.Fatalf("open mysql: %v", err)
	}
	// Release DB resources on process exit.
	defer db.Close()

	// Create concrete repository implementation bound to this DB pool.
	repo := repository.New(db)
	// Ensure all required tables and expected additive columns exist.
	if err := repo.EnsureSchema(ctx); err != nil {
		log.Fatalf("ensure schema: %v", err)
	}
	// Seed initial data only when tables are empty.
	if err := repo.SeedDefaults(ctx); err != nil {
		log.Fatalf("seed defaults: %v", err)
	}

	// Build HTTP application (templates, routes, middleware).
	app, err := httpapp.New(repo)
	if err != nil {
		log.Fatalf("build app: %v", err)
	}

	// Start serving requests; ListenAndServe returns only on fatal server stop.
	log.Printf("listening on %s", cfg.AppAddr)
	if err := http.ListenAndServe(cfg.AppAddr, app.Routes()); err != nil {
		log.Fatalf("listen: %v", err)
	}
}
