-- Manga table
CREATE TABLE manga (
	id UUID PRIMARY KEY DEFAULT uuid_generate_v4(), -- UUID as primary key
	name VARCHAR(255) NOT NULL UNIQUE,
	latest_chapter INT NOT NULL DEFAULT 0,
	created_at TIMESTAMP DEFAULT NOW()
);

