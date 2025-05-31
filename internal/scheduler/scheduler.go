package scheduler

import (
	"fmt"
	"log"
	"time"

	"mangadex-cli/internal/api"
	"mangadex-cli/internal/db"
	"mangadex-cli/internal/email"

	"github.com/robfig/cron/v3"
)


type UpdateInfo struct {
	MangaID    string
	MangaTitle string
	ChapterIDs []string
	Chapters   []api.Chapter
}
// CronScheduler handles periodic checking for manga updates
type CronScheduler struct {
	db           *db.DB
	apiClient    *api.MangaDexClient
	emailService *email.EmailService
	cron         *cron.Cron
	interval     int // seconds
	running      bool
}

// NewCronScheduler creates a new scheduler
func NewCronScheduler(
	database *db.DB,
	client *api.MangaDexClient,
	emailService *email.EmailService,
	checkInterval int,
) *CronScheduler {
	return &CronScheduler{
		db:           database,
		apiClient:    client,
		emailService: emailService,
		interval:     checkInterval,
		running:      false,
	}
}

// Start begins the update checking scheduler
func (s *CronScheduler) Start() error {
	if s.running {
		return fmt.Errorf("scheduler is already running")
	}

	// Create new cron scheduler
	s.cron = cron.New(cron.WithSeconds())

	// Schedule update checks
	schedule := fmt.Sprintf("@every %ds", s.interval)
	_, err := s.cron.AddFunc(schedule, func() {
		if err := s.CheckForUpdates(); err != nil {
			log.Printf("Error checking for updates: %v", err)
		}
	})

	if err != nil {
		return fmt.Errorf("failed to schedule update checks: %w", err)
	}

	// Start the cron scheduler
	s.cron.Start()
	s.running = true

	return nil
}

// Stop stops the update checking scheduler
func (s *CronScheduler) Stop() error {
	if !s.running || s.cron == nil {
		return nil
	}

	// Stop the cron scheduler
	s.cron.Stop()
	s.running = false

	return nil
}

// CheckForUpdates checks all active subscriptions for new chapters
func (s *CronScheduler) CheckForUpdates() error {
	log.Printf("Running scheduled update check at %s", time.Now().Format(time.RFC3339))

	// Get active subscriptions
	subscriptions, err := s.db.ListActiveSubscriptions()
	if err != nil {
		return fmt.Errorf("failed to get subscriptions: %w", err)
	}

	if len(subscriptions) == 0 {
		log.Println("No active subscriptions found")
		return nil
	}

	log.Printf("Checking updates for %d subscriptions...", len(subscriptions))

	// Track new chapters by user

	updates := make(map[int]map[string]*UpdateInfo) // UserID -> MangaID -> UpdateInfo

	// Check each subscription for updates
	for _, sub := range subscriptions {
		log.Printf("Checking \"%s\"...", sub.MangaTitle)

		// Get new chapters since last check
		chapters, err := s.apiClient.GetMangaChapters(sub.MangaID, sub.LastCheckTime)
		if err != nil {
			log.Printf("Error checking \"%s\": %v", sub.MangaTitle, err)
			continue
		}

		// Filter chapters by language
		languages := sub.GetLanguages()
		filteredChapters := make([]api.Chapter, 0)

		for _, chapter := range chapters {
			for _, lang := range languages {
				if chapter.TranslatedLanguage == lang {
					filteredChapters = append(filteredChapters, chapter)
					break
				}
			}
		}

		// Update last check time
		sub.LastCheckTime = time.Now()
		if err := s.db.UpdateSubscription(&sub); err != nil {
			log.Printf("Error updating subscription check time: %v", err)
		}

		// If no new chapters, continue
		if len(filteredChapters) == 0 {
			log.Printf("No new chapters for \"%s\"", sub.MangaTitle)
			continue
		}

		// Group updates by user and manga
		if _, ok := updates[sub.UserID]; !ok {
			updates[sub.UserID] = make(map[string]*UpdateInfo)
		}

		if _, ok := updates[sub.UserID][sub.MangaID]; !ok {
			updates[sub.UserID][sub.MangaID] = &UpdateInfo{
				MangaID:    sub.MangaID,
				MangaTitle: sub.MangaTitle,
				ChapterIDs: make([]string, 0),
				Chapters:   make([]api.Chapter, 0),
			}
		}

		// Add chapters to user updates
		updateInfo := updates[sub.UserID][sub.MangaID]
		for _, chapter := range filteredChapters {
			updateInfo.ChapterIDs = append(updateInfo.ChapterIDs, chapter.ID)
			updateInfo.Chapters = append(updateInfo.Chapters, chapter)
		}

		log.Printf("Found %d new chapter(s) for \"%s\"", len(filteredChapters), sub.MangaTitle)
	}

	// Process notifications for each user
	for userID, mangaUpdates := range updates {
		if err := s.ProcessUserNotifications(userID, mangaUpdates); err != nil {
			log.Printf("Error processing notifications for user %d: %v", userID, err)
		}
	}

	return nil
}

// ProcessUserNotifications sends notifications for a specific user's manga updates
func (s *CronScheduler) ProcessUserNotifications(userID int, mangaUpdates map[string]*UpdateInfo) error {
	// Get user
	user, err := s.db.GetUser(userID)
	if err != nil {
		return fmt.Errorf("failed to get user with ID %d: %w", userID, err)
	}

	// Connect to email server
	if err := s.emailService.Connect(); err != nil {
		return fmt.Errorf("failed to connect to email server: %w", err)
	}
	defer s.emailService.Disconnect()

	// Send notification for each manga with updates
	for _, updateInfo := range mangaUpdates {
		// Get manga details
		manga, err := s.apiClient.GetManga(updateInfo.MangaID)
		if err != nil {
			log.Printf("Error getting manga details for \"%s\": %v", updateInfo.MangaTitle, err)
			continue
		}

		// Send notification
		if err := s.emailService.SendNotification(user.Email, manga, updateInfo.Chapters); err != nil {
			log.Printf("Error sending notification to %s: %v", user.Email, err)
		} else {
			log.Printf("Notification sent to %s about %d new chapter(s) for \"%s\"",
				user.Email, len(updateInfo.Chapters), updateInfo.MangaTitle)
		}
	}

	return nil
}
