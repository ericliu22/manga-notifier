package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

// MangaDexClient handles API communication with MangaDex
type MangaDexClient struct {
	BaseURL      string
	SessionToken string
	RefreshToken string
	TokenExpiry  time.Time
	httpClient   *http.Client
}

// NewMangaDexClient creates a new MangaDex API client
func NewMangaDexClient(baseURL string) *MangaDexClient {
	return &MangaDexClient{
		BaseURL:    baseURL,
		httpClient: &http.Client{Timeout: 10 * time.Second},
	}
}

// Login authenticates with MangaDex API
func (client *MangaDexClient) Login(username, password string) error {
	// Prepare login data
	loginData := map[string]string{
		"username": username,
		"password": password,
	}
	
	jsonData, err := json.Marshal(loginData)
	if err != nil {
		return fmt.Errorf("error preparing login request: %w", err)
	}
	
	// Make auth request
	resp, err := client.httpClient.Post(
		fmt.Sprintf("%s/auth/login", client.BaseURL),
		"application/json",
		bytes.NewBuffer(jsonData),
	)
	
	if err != nil {
		return fmt.Errorf("login request failed: %w", err)
	}
	defer resp.Body.Close()
	
	// Read response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read login response: %w", err)
	}
	
	// Check status code
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("login failed with status code %d: %s", resp.StatusCode, string(body))
	}
	
	// Parse response
	var loginResp struct {
		Result   string `json:"result"`
		Token    struct {
			Session string `json:"session"`
			Refresh string `json:"refresh"`
		} `json:"token"`
	}
	
	if err := json.Unmarshal(body, &loginResp); err != nil {
		return fmt.Errorf("failed to parse login response: %w", err)
	}
	
	// Store tokens
	client.SessionToken = loginResp.Token.Session
	client.RefreshToken = loginResp.Token.Refresh
	client.TokenExpiry = time.Now().Add(15 * time.Minute) // Token expires after 15 minutes
	
	return nil
}

// RefreshToken refreshes the authentication token
func (client *MangaDexClient) RefreshClientToken() error {
	// Prepare refresh token request
	refreshData := map[string]string{
		"token": client.RefreshToken,
	}
	
	jsonData, err := json.Marshal(refreshData)
	if err != nil {
		return fmt.Errorf("error preparing refresh token request: %w", err)
	}
	
	// Make refresh request
	resp, err := client.httpClient.Post(
		fmt.Sprintf("%s/auth/refresh", client.BaseURL),
		"application/json",
		bytes.NewBuffer(jsonData),
	)
	
	if err != nil {
		return fmt.Errorf("refresh token request failed: %w", err)
	}
	defer resp.Body.Close()
	
	// Read response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read refresh token response: %w", err)
	}
	
	// Check status code
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("refresh token failed with status code %d: %s", resp.StatusCode, string(body))
	}
	
	// Parse response
	var refreshResp struct {
		Result string `json:"result"`
		Token  struct {
			Session string `json:"session"`
			Refresh string `json:"refresh"`
		} `json:"token"`
	}
	
	if err := json.Unmarshal(body, &refreshResp); err != nil {
		return fmt.Errorf("failed to parse refresh token response: %w", err)
	}
	
	// Store new tokens
	client.SessionToken = refreshResp.Token.Session
	client.RefreshToken = refreshResp.Token.Refresh
	client.TokenExpiry = time.Now().Add(15 * time.Minute) // Token expires after 15 minutes
	
	return nil
}

// ensureValidToken makes sure the session token is valid
func (client *MangaDexClient) ensureValidToken() error {
	if client.SessionToken == "" {
		return fmt.Errorf("not authenticated")
	}
	
	if time.Now().After(client.TokenExpiry) {
		return client.RefreshClientToken()
	}
	
	return nil
}

// makeRequest makes an authenticated request to the MangaDex API
func (client *MangaDexClient) makeRequest(method, endpoint string, queryParams map[string]string) ([]byte, error) {
	// Build URL with query parameters
	reqURL, err := url.Parse(fmt.Sprintf("%s%s", client.BaseURL, endpoint))
	if err != nil {
		return nil, fmt.Errorf("failed to parse URL: %w", err)
	}
	
	// Add query parameters if provided
	if queryParams != nil {
		q := reqURL.Query()
		for key, value := range queryParams {
			q.Add(key, value)
		}
		reqURL.RawQuery = q.Encode()
	}
	
	// Create request
	req, err := http.NewRequest(method, reqURL.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	
	// Add authentication header if we have a token
	if client.SessionToken != "" {
		if err := client.ensureValidToken(); err != nil {
			return nil, fmt.Errorf("failed to ensure valid token: %w", err)
		}
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", client.SessionToken))
	}
	
	// Make request
	resp, err := client.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()
	
	// Read response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}
	
	// Check for error status codes
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("request failed with status code %d: %s", resp.StatusCode, string(body))
	}
	
	return body, nil
}

// GetManga gets details for a specific manga by ID
func (client *MangaDexClient) GetManga(id string) (*Manga, error) {
	body, err := client.makeRequest(http.MethodGet, fmt.Sprintf("/manga/%s", id), nil)
	if err != nil {
		return nil, err
	}
	
	var response struct {
		Result string `json:"result"`
		Data   struct {
			ID         string               `json:"id"`
			Attributes MangaAttributesDTO   `json:"attributes"`
			Relationships []RelationshipDTO `json:"relationships"`
		} `json:"data"`
	}
	
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to parse manga response: %w", err)
	}
	
	// Convert API response to our Manga model
	manga := &Manga{
		ID:          response.Data.ID,
		Title:       response.Data.Attributes.Title,
		Description: response.Data.Attributes.Description,
		Status:      response.Data.Attributes.Status,
		CreatedAt:   response.Data.Attributes.CreatedAt,
		UpdatedAt:   response.Data.Attributes.UpdatedAt,
	}
	
	// Extract cover art URL if available
	for _, rel := range response.Data.Relationships {
		if rel.Type == "cover_art" {
			manga.CoverArtID = rel.ID
			break
		}
	}
	
	if manga.CoverArtID != "" {
		manga.CoverArtURL = fmt.Sprintf("%s/covers/%s/%s.256.jpg", client.BaseURL, manga.ID, manga.CoverArtID)
	}
	
	return manga, nil
}

// SearchManga searches for manga by title
func (client *MangaDexClient) SearchManga(title string) ([]*Manga, error) {
	params := map[string]string{
		"title": title,
		"limit": "5",
		"order[relevance]": "desc",
	}
	
	body, err := client.makeRequest(http.MethodGet, "/manga", params)
	if err != nil {
		return nil, err
	}
	
	var response struct {
		Result string `json:"result"`
		Data   []struct {
			ID         string               `json:"id"`
			Attributes MangaAttributesDTO   `json:"attributes"`
			Relationships []RelationshipDTO `json:"relationships"`
		} `json:"data"`
	}
	
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to parse manga search response: %w", err)
	}
	
	// Convert API response to our Manga model
	mangas := make([]*Manga, 0, len(response.Data))
	for _, data := range response.Data {
		manga := &Manga{
			ID:          data.ID,
			Title:       data.Attributes.Title,
			Description: data.Attributes.Description,
			Status:      data.Attributes.Status,
			CreatedAt:   data.Attributes.CreatedAt,
			UpdatedAt:   data.Attributes.UpdatedAt,
		}
		
		// Extract cover art URL if available
		for _, rel := range data.Relationships {
			if rel.Type == "cover_art" {
				manga.CoverArtID = rel.ID
				break
			}
		}
		
		if manga.CoverArtID != "" {
			manga.CoverArtURL = fmt.Sprintf("%s/covers/%s/%s.256.jpg", client.BaseURL, manga.ID, manga.CoverArtID)
		}
		
		mangas = append(mangas, manga)
	}
	
	return mangas, nil
}

// GetMangaChapters gets chapters for a manga, optionally since a specific time
func (client *MangaDexClient) GetMangaChapters(mangaID string, since time.Time) ([]Chapter, error) {
	params := map[string]string{
		"manga":              mangaID,
		"order[publishAt]":   "desc",
		"limit":              "50",
		"contentRating[]":    "safe",
	}
	
	// Add "createdAt" filter if "since" is not zero time
	if !since.IsZero() {
		params["createdAtSince"] = since.Format(time.RFC3339)
	}
	
	body, err := client.makeRequest(http.MethodGet, "/chapter", params)
	if err != nil {
		return nil, err
	}
	
	var response struct {
		Result string `json:"result"`
		Data   []struct {
			ID         string                `json:"id"`
			Attributes ChapterAttributesDTO  `json:"attributes"`
			Relationships []RelationshipDTO  `json:"relationships"`
		} `json:"data"`
	}
	
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to parse chapter response: %w", err)
	}
	
	// Convert API response to our Chapter model
	chapters := make([]Chapter, 0, len(response.Data))
	for _, data := range response.Data {
		chapter := Chapter{
			ID:                data.ID,
			Title:             data.Attributes.Title,
			Volume:            data.Attributes.Volume,
			Chapter:           data.Attributes.Chapter,
			TranslatedLanguage: data.Attributes.TranslatedLanguage,
			PublishAt:         data.Attributes.PublishAt,
			CreatedAt:         data.Attributes.CreatedAt,
			UpdatedAt:         data.Attributes.UpdatedAt,
		}
		
		// Extract scanlation groups
		groups := make([]string, 0)
		for _, rel := range data.Relationships {
			if rel.Type == "scanlation_group" {
				groups = append(groups, rel.ID)
			}
		}
		chapter.Groups = groups
		
		chapters = append(chapters, chapter)
	}
	
	return chapters, nil
}

// GetChapterDetails gets detailed information for a chapter
func (client *MangaDexClient) GetChapterDetails(chapterID string) (*Chapter, error) {
	body, err := client.makeRequest(http.MethodGet, fmt.Sprintf("/chapter/%s", chapterID), nil)
	if err != nil {
		return nil, err
	}
	
	var response struct {
		Result string `json:"result"`
		Data   struct {
			ID         string               `json:"id"`
			Attributes ChapterAttributesDTO `json:"attributes"`
			Relationships []RelationshipDTO `json:"relationships"`
		} `json:"data"`
	}
	
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to parse chapter details response: %w", err)
	}
	
	chapter := &Chapter{
		ID:                response.Data.ID,
		Title:             response.Data.Attributes.Title,
		Volume:            response.Data.Attributes.Volume,
		Chapter:           response.Data.Attributes.Chapter,
		TranslatedLanguage: response.Data.Attributes.TranslatedLanguage,
		PublishAt:         response.Data.Attributes.PublishAt,
		CreatedAt:         response.Data.Attributes.CreatedAt,
		UpdatedAt:         response.Data.Attributes.UpdatedAt,
	}
	
	// Extract scanlation groups
	groups := make([]string, 0)
	for _, rel := range response.Data.Relationships {
		if rel.Type == "scanlation_group" {
			groups = append(groups, rel.ID)
		}
	}
	chapter.Groups = groups
	
	return chapter, nil
}
