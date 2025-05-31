package api

import (
	"time"
)

// Manga represents a manga series
type Manga struct {
	ID          string             `json:"id"`
	Title       map[string]string  `json:"title"`
	Description map[string]string  `json:"description"`
	CoverArtID  string             `json:"cover_art_id"`
	CoverArtURL string             `json:"cover_art_url"`
	Tags        []string           `json:"tags"`
	Status      string             `json:"status"`
	CreatedAt   time.Time          `json:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at"`
}

// Chapter represents a manga chapter
type Chapter struct {
	ID                string    `json:"id"`
	Title             string    `json:"title"`
	Volume            string    `json:"volume"`
	Chapter           string    `json:"chapter"`
	TranslatedLanguage string    `json:"translated_language"`
	Groups            []string  `json:"groups"`
	PublishAt         time.Time `json:"publish_at"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

// MangaAttributesDTO represents the attributes in the MangaDex API manga response
type MangaAttributesDTO struct {
	Title       map[string]string     `json:"title"`
	Description map[string]string     `json:"description"`
	Status      string                `json:"status"`
	CreatedAt   time.Time             `json:"createdAt"`
	UpdatedAt   time.Time             `json:"updatedAt"`
}

// ChapterAttributesDTO represents the attributes in the MangaDex API chapter response
type ChapterAttributesDTO struct {
	Title             string    `json:"title"`
	Volume            string    `json:"volume"`
	Chapter           string    `json:"chapter"`
	TranslatedLanguage string    `json:"translatedLanguage"`
	PublishAt         time.Time `json:"publishAt"`
	CreatedAt         time.Time `json:"createdAt"`
	UpdatedAt         time.Time `json:"updatedAt"`
}

// RelationshipDTO represents a relationship in MangaDex API responses
type RelationshipDTO struct {
	ID   string `json:"id"`
	Type string `json:"type"`
}

// GetTitle returns the title in the preferred language, falling back to English or the first available
func (m *Manga) GetTitle() string {
	if m.Title == nil {
		return "Unknown"
	}
	
	// Try English first
	if title, ok := m.Title["en"]; ok {
		return title
	}
	
	// Then try Japanese
	if title, ok := m.Title["ja"]; ok {
		return title
	}
	
	// Otherwise return the first available title
	for _, title := range m.Title {
		return title
	}
	
	return "Unknown"
}

// GetDescription returns the description in the preferred language
func (m *Manga) GetDescription() string {
	if m.Description == nil {
		return ""
	}
	
	// Try English first
	if desc, ok := m.Description["en"]; ok {
		return desc
	}
	
	// Then try Japanese
	if desc, ok := m.Description["ja"]; ok {
		return desc
	}
	
	// Otherwise return the first available description
	for _, desc := range m.Description {
		return desc
	}
	
	return ""
}