package models

import (
	"time"
)

type Subscription struct {
	ID                  string    `json:"id"`                    // UUID of the subscription
	UserID              string    `json:"user_id"`               // UUID of the user
	MangaID             string    `json:"manga_id"`              // UUID of the manga
	LastNotifiedChapter int       `json:"last_notified_chapter"` // Last chapter notified
	SubscribedAt        time.Time `json:"subscribed_at"`         // Subscription timestamp
}
