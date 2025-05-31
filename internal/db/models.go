package db

import (
	"strings"
	"time"
)

// User represents a user who receives notifications
type User struct {
	ID        int       `gorm:"primaryKey" json:"id"`
	Email     string    `gorm:"uniqueIndex" json:"email"`
	Name      string    `json:"name"`
	Active    bool      `gorm:"default:true" json:"active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Subscription represents a manga subscription for a user
type Subscription struct {
	ID             int       `gorm:"primaryKey" json:"id"`
	UserID         int       `json:"user_id"`
	MangaID        string    `json:"manga_id"`
	MangaTitle     string    `json:"manga_title"`
	Languages      string    `json:"languages"` // Comma-separated language codes
	LastCheckTime  time.Time `json:"last_check_time"`
	LastChapterTime time.Time `json:"last_chapter_time"`
	Active         bool      `gorm:"default:true" json:"active"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// GetLanguages returns the list of languages for this subscription
func (s *Subscription) GetLanguages() []string {
	if s.Languages == "" {
		return []string{"en"} // Default to English
	}
	
	// Split the comma-separated language string and trim spaces
	languages := strings.Split(s.Languages, ",")
	for i, lang := range languages {
		languages[i] = strings.TrimSpace(lang)
	}
	
	return languages
}