// Module path used for imports inside this repository.
module github.com/pavom/resort-clone

// Minimum Go toolchain version expected by this project.
go 1.23.0

// Direct runtime dependency: MySQL driver for database/sql.
require github.com/go-sql-driver/mysql v1.8.1

// Indirect dependency required by transitive package graph.
require filippo.io/edwards25519 v1.1.0 // indirect
