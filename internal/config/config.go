package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// SMTPConfig stores email server settings
type SMTPConfig struct {
	Server    string `json:"server"`
	Port      int    `json:"port"`
	Username  string `json:"username"`
	Password  string `json:"password"`
	UseTLS    bool   `json:"use_tls"`
	FromEmail string `json:"from_email"`
	FromName  string `json:"from_name"`
}

// Config stores the application configuration
type Config struct {
	DatabasePath       string     `json:"database_path"`
	SMTPSettings       SMTPConfig `json:"smtp_settings"`
	UpdateCheckInterval int        `json:"update_check_interval"` // in seconds
	MangaDexAPIURL     string     `json:"mangadex_api_url"`
	AuthToken          string     `json:"auth_token"`
	RefreshToken       string     `json:"refresh_token"`
	TokenExpiry        time.Time     `json:"token_expiry"` // ISO8601/RFC3339 format
}

// Load reads the configuration from a file
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var config Config
	err = json.Unmarshal(data, &config)
	if err != nil {
		return nil, fmt.Errorf("error parsing config file: %w", err)
	}

	return &config, nil
}

// Save writes the configuration to a file
func (c *Config) Save(path string) error {
	// Create directory if it doesn't exist
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return fmt.Errorf("error serializing config: %w", err)
	}

	err = os.WriteFile(path, data, 0600) // Restricted permissions for security
	if err != nil {
		return fmt.Errorf("error writing config file: %w", err)
	}

	return nil
}

// Default returns a default configuration
func Default() *Config {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		homeDir = "."
	}

	return &Config{
		DatabasePath: filepath.Join(homeDir, ".mangadex-cli", "mangadex.db"),
		SMTPSettings: SMTPConfig{
			Server:    "smtp.example.com",
			Port:      587,
			Username:  "",
			Password:  "",
			UseTLS:    true,
			FromEmail: "manga-updates@example.com",
			FromName:  "MangaDex Notifier",
		},
		UpdateCheckInterval: 3600, // 1 hour
		MangaDexAPIURL:      "https://api.mangadex.org",
		AuthToken:           "",
		RefreshToken:        "",
		TokenExpiry:         time.Now(),
	}
}
