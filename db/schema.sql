-- schema.sql defines the baseline relational structure required by the app.
-- Tables are designed for a public-facing resort website with optional source snapshots.

-- site_content stores singleton-style homepage and contact content.
-- The application currently reads the first row (lowest id) as active site content.
CREATE TABLE IF NOT EXISTS site_content (
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
);

-- rooms stores catalog cards for room categories displayed on homepage.
CREATE TABLE IF NOT EXISTS rooms (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(255) NOT NULL,
    summary TEXT NOT NULL,
    capacity INT NOT NULL,
    size_sqm INT NOT NULL,
    room_count INT NOT NULL,
    feature_note VARCHAR(255) NOT NULL,
    image_path VARCHAR(255) NOT NULL DEFAULT ''
);

-- gallery_items stores visual cards for the gallery section.
-- tone controls fallback styling when image_path is empty.
CREATE TABLE IF NOT EXISTS gallery_items (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    title VARCHAR(255) NOT NULL,
    subtitle TEXT NOT NULL,
    tone VARCHAR(64) NOT NULL,
    image_path VARCHAR(255) NOT NULL DEFAULT ''
);

-- scrape_snapshots stores source-page captures used for informational display.
-- It is append-only in current app flow, with latest record shown on homepage.
CREATE TABLE IF NOT EXISTS scrape_snapshots (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    source_url VARCHAR(255) NOT NULL,
    page_title VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    markdown MEDIUMTEXT NOT NULL,
    html MEDIUMTEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
