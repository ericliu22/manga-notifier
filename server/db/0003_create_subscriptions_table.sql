-- Subscriptions table
CREATE TABLE subscriptions (
	id UUID PRIMARY KEY DEFAULT uuid_generate_v4(), -- UUID as primary key
	user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
	manga_id UUID NOT NULL REFERENCES manga(id) ON DELETE CASCADE,
	last_notified_chapter INT DEFAULT 0,
	subscribed_at TIMESTAMP DEFAULT NOW(),
	UNIQUE(user_id, manga_id)
);
