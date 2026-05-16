package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"hotel_site/internal/service"
)

// Repository is a MySQL-backed implementation for content storage and retrieval.
type Repository struct {
	// db is the shared SQL connection pool used by all repository methods.
	db *sql.DB
}

// New creates a Repository bound to the provided SQL connection pool.
func New(db *sql.DB) *Repository {
	return &Repository{db: db}
}

// EnsureSchema creates required tables and additive columns if they do not yet exist.
// The operation is idempotent and safe to call on every startup.
func (r *Repository) EnsureSchema(ctx context.Context) error {
	// Base CREATE TABLE statements for all domain entities used by the app.
	statements := []string{
		`CREATE TABLE IF NOT EXISTS site_content (
			id BIGINT PRIMARY KEY AUTO_INCREMENT,
			brand_name VARCHAR(255) NOT NULL DEFAULT '',
			tagline VARCHAR(255) NOT NULL DEFAULT '',
			hero_title VARCHAR(255) NOT NULL,
			hero_subtitle TEXT NOT NULL,
			hero_image_path VARCHAR(255) NOT NULL DEFAULT '',
			season_note TEXT NOT NULL,
			booking_url VARCHAR(255) NOT NULL,
			address VARCHAR(255) NOT NULL,
			phone VARCHAR(64) NOT NULL,
			email VARCHAR(255) NOT NULL,
			check_in VARCHAR(32) NOT NULL,
			check_out VARCHAR(32) NOT NULL,
			about_title VARCHAR(255) NOT NULL,
			about_body TEXT NOT NULL,
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS rooms (
			id BIGINT PRIMARY KEY AUTO_INCREMENT,
			name VARCHAR(255) NOT NULL,
			summary TEXT NOT NULL,
			capacity INT NOT NULL,
			size_sqm INT NOT NULL,
			room_count INT NOT NULL,
			feature_note VARCHAR(255) NOT NULL,
			image_path VARCHAR(255) NOT NULL DEFAULT ''
		)`,
		`CREATE TABLE IF NOT EXISTS gallery_items (
			id BIGINT PRIMARY KEY AUTO_INCREMENT,
			title VARCHAR(255) NOT NULL,
			subtitle TEXT NOT NULL,
			tone VARCHAR(64) NOT NULL,
			image_path VARCHAR(255) NOT NULL DEFAULT ''
		)`,
		`CREATE TABLE IF NOT EXISTS scrape_snapshots (
			id BIGINT PRIMARY KEY AUTO_INCREMENT,
			source_url VARCHAR(255) NOT NULL,
			page_title VARCHAR(255) NOT NULL,
			description TEXT NOT NULL,
			markdown MEDIUMTEXT NOT NULL,
			html MEDIUMTEXT NOT NULL,
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
		)`,
	}

	for _, stmt := range statements {
		if _, err := r.db.ExecContext(ctx, stmt); err != nil {
			return fmt.Errorf("ensure schema: %w", err)
		}
	}

	// alterSpec describes an additive migration for a single missing column.
	type alterSpec struct {
		table      string
		column     string
		definition string
	}

	// Backward-compatible additive migrations for previously created databases.
	alters := []alterSpec{
		{table: "site_content", column: "brand_name", definition: "VARCHAR(255) NOT NULL DEFAULT ''"},
		{table: "site_content", column: "tagline", definition: "VARCHAR(255) NOT NULL DEFAULT ''"},
		{table: "site_content", column: "hero_image_path", definition: "VARCHAR(255) NOT NULL DEFAULT ''"},
		{table: "rooms", column: "image_path", definition: "VARCHAR(255) NOT NULL DEFAULT ''"},
		{table: "gallery_items", column: "image_path", definition: "VARCHAR(255) NOT NULL DEFAULT ''"},
	}

	for _, spec := range alters {
		// Skip ALTER when target column already exists.
		exists, err := r.columnExists(ctx, spec.table, spec.column)
		if err != nil {
			return fmt.Errorf("ensure alters: %w", err)
		}
		if exists {
			continue
		}

		stmt := fmt.Sprintf("ALTER TABLE %s ADD COLUMN %s %s", spec.table, spec.column, spec.definition)
		if _, err := r.db.ExecContext(ctx, stmt); err != nil {
			return fmt.Errorf("ensure alters: %w", err)
		}
	}

	return nil
}

// columnExists checks information_schema to detect whether a table column already exists.
func (r *Repository) columnExists(ctx context.Context, tableName, columnName string) (bool, error) {
	var count int
	err := r.db.QueryRowContext(ctx, `
		SELECT COUNT(*)
		FROM information_schema.columns
		WHERE table_schema = DATABASE() AND table_name = ? AND column_name = ?`,
		tableName, columnName,
	).Scan(&count)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// SeedDefaults inserts starter records for site, rooms, and gallery when tables are empty.
// Existing user-managed data is preserved and never overwritten.
func (r *Repository) SeedDefaults(ctx context.Context) error {
	if err := r.seedSite(ctx); err != nil {
		return err
	}
	if err := r.seedRooms(ctx); err != nil {
		return err
	}
	if err := r.seedGallery(ctx); err != nil {
		return err
	}

	return nil
}

// seedSite inserts exactly one starter site_content row when the table has no data.
func (r *Repository) seedSite(ctx context.Context) error {
	var count int
	if err := r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM site_content`).Scan(&count); err != nil {
		return err
	}
	if count > 0 {
		return nil
	}

	site := service.DefaultSiteContent()
	_, err := r.db.ExecContext(ctx, `INSERT INTO site_content (
		brand_name, tagline, hero_title, hero_subtitle, hero_image_path, season_note, booking_url, address, phone, email,
		check_in, check_out, about_title, about_body
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		site.BrandName, site.Tagline, site.HeroTitle, site.HeroSubtitle, site.HeroImagePath, site.SeasonNote, site.BookingURL, site.Address,
		site.Phone, site.Email, site.CheckIn, site.CheckOut, site.AboutTitle, site.AboutBody,
	)
	return err
}

// seedRooms inserts starter room categories when the rooms table is empty.
func (r *Repository) seedRooms(ctx context.Context) error {
	var count int
	if err := r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM rooms`).Scan(&count); err != nil {
		return err
	}
	if count > 0 {
		return nil
	}

	for _, room := range service.DefaultRooms() {
		if _, err := r.db.ExecContext(ctx, `INSERT INTO rooms (
			name, summary, capacity, size_sqm, room_count, feature_note, image_path
		) VALUES (?, ?, ?, ?, ?, ?, ?)`,
			room.Name, room.Summary, room.Capacity, room.SizeSQM, room.RoomCount, room.FeatureNote, room.ImagePath,
		); err != nil {
			return err
		}
	}

	return nil
}

// seedGallery inserts starter gallery cards when gallery_items is empty.
func (r *Repository) seedGallery(ctx context.Context) error {
	var count int
	if err := r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM gallery_items`).Scan(&count); err != nil {
		return err
	}
	if count > 0 {
		return nil
	}

	for _, item := range service.DefaultGallery() {
		if _, err := r.db.ExecContext(ctx, `INSERT INTO gallery_items (title, subtitle, tone, image_path) VALUES (?, ?, ?, ?)`,
			item.Title, item.Subtitle, item.Tone, item.ImagePath,
		); err != nil {
			return err
		}
	}

	return nil
}

// HomePageData loads all data blocks needed by homepage render/API:
// site singleton row, ordered rooms, ordered gallery, and optional latest snapshot.
func (r *Repository) HomePageData(ctx context.Context) (service.HomePageData, error) {
	data := service.HomePageData{}

	row := r.db.QueryRowContext(ctx, `SELECT id, brand_name, tagline, hero_title, hero_subtitle, hero_image_path, season_note, booking_url, address, phone, email, check_in, check_out, about_title, about_body, created_at, updated_at
		FROM site_content ORDER BY id ASC LIMIT 1`)
	if err := row.Scan(
		&data.Site.ID, &data.Site.BrandName, &data.Site.Tagline, &data.Site.HeroTitle, &data.Site.HeroSubtitle, &data.Site.HeroImagePath, &data.Site.SeasonNote,
		&data.Site.BookingURL, &data.Site.Address, &data.Site.Phone, &data.Site.Email,
		&data.Site.CheckIn, &data.Site.CheckOut, &data.Site.AboutTitle, &data.Site.AboutBody,
		&data.Site.CreatedAt, &data.Site.UpdatedAt,
	); err != nil {
		return service.HomePageData{}, err
	}

	roomsRows, err := r.db.QueryContext(ctx, `SELECT id, name, summary, capacity, size_sqm, room_count, feature_note, image_path FROM rooms ORDER BY id ASC`)
	if err != nil {
		return service.HomePageData{}, err
	}
	defer roomsRows.Close()

	for roomsRows.Next() {
		var room service.Room
		if err := roomsRows.Scan(&room.ID, &room.Name, &room.Summary, &room.Capacity, &room.SizeSQM, &room.RoomCount, &room.FeatureNote, &room.ImagePath); err != nil {
			return service.HomePageData{}, err
		}
		data.Rooms = append(data.Rooms, room)
	}
	if err := roomsRows.Err(); err != nil {
		return service.HomePageData{}, err
	}

	galleryRows, err := r.db.QueryContext(ctx, `SELECT id, title, subtitle, tone, image_path FROM gallery_items ORDER BY id ASC`)
	if err != nil {
		return service.HomePageData{}, err
	}
	defer galleryRows.Close()

	for galleryRows.Next() {
		var item service.GalleryItem
		if err := galleryRows.Scan(&item.ID, &item.Title, &item.Subtitle, &item.Tone, &item.ImagePath); err != nil {
			return service.HomePageData{}, err
		}
		data.Gallery = append(data.Gallery, item)
	}
	if err := galleryRows.Err(); err != nil {
		return service.HomePageData{}, err
	}

	snapshot, err := r.LatestSnapshot(ctx)
	if err == nil {
		data.Snapshot = &snapshot
	} else if !errors.Is(err, sql.ErrNoRows) {
		return service.HomePageData{}, err
	}

	return data, nil
}

// UpdateSiteContent updates mutable fields for an existing site_content row by ID.
func (r *Repository) UpdateSiteContent(ctx context.Context, site service.SiteContent) error {
	_, err := r.db.ExecContext(ctx, `UPDATE site_content SET
		brand_name = ?, tagline = ?, hero_title = ?, hero_subtitle = ?, hero_image_path = ?, season_note = ?,
		booking_url = ?, address = ?, phone = ?, email = ?, check_in = ?, check_out = ?, about_title = ?, about_body = ?
		WHERE id = ?`,
		site.BrandName, site.Tagline, site.HeroTitle, site.HeroSubtitle, site.HeroImagePath, site.SeasonNote,
		site.BookingURL, site.Address, site.Phone, site.Email, site.CheckIn, site.CheckOut, site.AboutTitle, site.AboutBody, site.ID,
	)
	return err
}

// CreateRoom inserts a new room category row.
func (r *Repository) CreateRoom(ctx context.Context, room service.Room) error {
	_, err := r.db.ExecContext(ctx, `INSERT INTO rooms (
		name, summary, capacity, size_sqm, room_count, feature_note, image_path
	) VALUES (?, ?, ?, ?, ?, ?, ?)`,
		room.Name, room.Summary, room.Capacity, room.SizeSQM, room.RoomCount, room.FeatureNote, room.ImagePath,
	)
	return err
}

// UpdateRoom updates an existing room category by ID.
func (r *Repository) UpdateRoom(ctx context.Context, room service.Room) error {
	_, err := r.db.ExecContext(ctx, `UPDATE rooms SET
		name = ?, summary = ?, capacity = ?, size_sqm = ?, room_count = ?, feature_note = ?, image_path = ?
		WHERE id = ?`,
		room.Name, room.Summary, room.Capacity, room.SizeSQM, room.RoomCount, room.FeatureNote, room.ImagePath, room.ID,
	)
	return err
}

// DeleteRoom removes a room category by primary key.
func (r *Repository) DeleteRoom(ctx context.Context, id int64) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM rooms WHERE id = ?`, id)
	return err
}

// CreateGalleryItem inserts a new gallery entry.
func (r *Repository) CreateGalleryItem(ctx context.Context, item service.GalleryItem) error {
	_, err := r.db.ExecContext(ctx, `INSERT INTO gallery_items (title, subtitle, tone, image_path) VALUES (?, ?, ?, ?)`,
		item.Title, item.Subtitle, item.Tone, item.ImagePath,
	)
	return err
}

// UpdateGalleryItem updates an existing gallery entry by ID.
func (r *Repository) UpdateGalleryItem(ctx context.Context, item service.GalleryItem) error {
	_, err := r.db.ExecContext(ctx, `UPDATE gallery_items SET title = ?, subtitle = ?, tone = ?, image_path = ? WHERE id = ?`,
		item.Title, item.Subtitle, item.Tone, item.ImagePath, item.ID,
	)
	return err
}

// DeleteGalleryItem removes a gallery entry by primary key.
func (r *Repository) DeleteGalleryItem(ctx context.Context, id int64) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM gallery_items WHERE id = ?`, id)
	return err
}

// SaveSnapshot persists a captured source-page snapshot for later reference.
func (r *Repository) SaveSnapshot(ctx context.Context, snapshot service.Snapshot) error {
	_, err := r.db.ExecContext(ctx, `INSERT INTO scrape_snapshots (
		source_url, page_title, description, markdown, html
	) VALUES (?, ?, ?, ?, ?)`,
		snapshot.SourceURL, snapshot.PageTitle, snapshot.Description, snapshot.Markdown, snapshot.HTML,
	)

	return err
}

// LatestSnapshot returns the most recent scrape snapshot by descending ID.
func (r *Repository) LatestSnapshot(ctx context.Context) (service.Snapshot, error) {
	var snapshot service.Snapshot
	err := r.db.QueryRowContext(ctx, `SELECT id, source_url, page_title, description, markdown, html, created_at
		FROM scrape_snapshots ORDER BY id DESC LIMIT 1`).Scan(
		&snapshot.ID, &snapshot.SourceURL, &snapshot.PageTitle, &snapshot.Description,
		&snapshot.Markdown, &snapshot.HTML, &snapshot.CreatedAt,
	)
	if err != nil {
		return service.Snapshot{}, err
	}

	return snapshot, nil
}

// OpenMySQL initializes SQL driver connection pool, applies pool limits,
// and performs a startup ping with timeout to fail fast on bad connectivity.
func OpenMySQL(ctx context.Context, dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	db.SetConnMaxLifetime(5 * time.Minute)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := db.PingContext(pingCtx); err != nil {
		return nil, err
	}

	return db, nil
}
