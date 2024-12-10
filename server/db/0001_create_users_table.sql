CREATE EXTENSION IF NOT EXISTS "uuid-ossp"; -- Required for UUID generation

-- Users table
CREATE TABLE users (
	id UUID PRIMARY KEY DEFAULT uuid_generate_v4(), -- UUID as primary key
	email VARCHAR(255) NOT NULL UNIQUE,
	created_at TIMESTAMP DEFAULT NOW()
);
