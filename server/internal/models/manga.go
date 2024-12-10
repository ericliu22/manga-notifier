package models

import (
	"time"
)

type Manga struct {
	ID             string    `json:"id"`
	Email          string    `json:"email"`
	latest_chapter int       `json:"latest_chapter"`
	CreatedAt      time.Time `json:"created_at"`
}
