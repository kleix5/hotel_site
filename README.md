# Resort Clone Starter

<!-- Project purpose: quick overview of what this starter includes. -->
This project is a fresh Go + MySQL starter that recreates the structure of the reference hotel site at `https://xn--d1abkwbsq8g.su` without copying its exact markup.

## Stack

<!-- Core technology choices and why these parts exist. -->
- Go HTTP server
- MySQL for page content, room inventory, gallery items, and source snapshots
- Server-rendered HTML templates with custom CSS

## Run locally

<!-- Local setup flow for running app directly without containers. -->
1. Install Go 1.23+ and MySQL 8+.
2. Copy `.env.example` to `.env` or export the variables manually.
3. Create the database:

```sql
CREATE DATABASE resort_clone CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

4. Download Go dependencies:

```bash
go mod tidy
```

5. Start the app:

```bash
go run ./cmd/server
```

6. Open `http://localhost:8080`.

## Run with Docker Compose

<!-- Container-based setup flow for teams that prefer reproducible local environments. -->
1. Copy `.env.example` to `.env`.
2. Start the stack:

```bash
docker compose up --build
```

3. Open `http://localhost:8080`.

The compose setup starts:

- `app`: the Go web server
- `db`: MySQL 8 with a persistent named volume

Important: for Docker Compose, the app connects to MySQL using the service name `db`. You do not need to set `MYSQL_DSN` in `.env`.

To stop it:

```bash
docker compose down
```

To stop it and remove the MySQL volume too:

```bash
docker compose down -v
```

## Separate admin direction

<!-- Product-direction note explaining why richer fields exist in the schema. -->
This app is now focused on the public website only. The richer content fields and image paths remain in the data model so a separate admin application can manage them later over the same database or through an internal API.

## Notes

<!-- Operational notes and caveats for first-time adopters. -->
- The homepage content is seeded automatically the first time the app boots.
- Source snapshots are stored in `scrape_snapshots` for analysis and do not auto-overwrite page copy.
- Replace the starter text, legal details, and visuals before treating this as production-ready hospitality content.
- Git is safe to initialize now; `.env` and local volume data are ignored.
