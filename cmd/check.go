package cmd

import (
	"fmt"
	"time"

	"mangadex-cli/internal/api"
	"mangadex-cli/internal/email"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"os"
)

// checkCmd represents the check command
var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "Manually check for manga updates",
	Long: `Manually check for updates to your manga subscriptions.
This command performs the same check that the service would do on schedule.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Initialize API client
		client := api.NewMangaDexClient(cfg.MangaDexAPIURL)
		
		// Set auth token if available
		if cfg.AuthToken != "" {
			client.SessionToken = cfg.AuthToken
			client.RefreshToken = cfg.RefreshToken
			client.TokenExpiry = cfg.TokenExpiry 
		}
		
		// Initialize email service
		emailService := email.NewEmailService(cfg.SMTPSettings)
		
		// Get active subscriptions
		subscriptions, err := database.ListActiveSubscriptions()
		if err != nil {
			return fmt.Errorf("failed to get subscriptions: %w", err)
		}
		
		if len(subscriptions) == 0 {
			fmt.Println("No active subscriptions found")
			return nil
		}
		
		fmt.Printf("Checking updates for %d subscriptions...\n", len(subscriptions))
		
		// Track new chapters
		type UpdateInfo struct {
			MangaTitle string
			ChapterIDs []string
			Chapters   []api.Chapter
		}
		
		updates := make(map[int]*UpdateInfo) // UserID -> UpdateInfo
		
		// Check each subscription for updates
		for _, sub := range subscriptions {
			fmt.Printf("Checking \"%s\"...\n", sub.MangaTitle)
			
			// Get new chapters since last check
			chapters, err := client.GetMangaChapters(sub.MangaID, sub.LastCheckTime)
			if err != nil {
				fmt.Printf("Error checking \"%s\": %v\n", sub.MangaTitle, err)
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
			if err := database.UpdateSubscription(&sub); err != nil {
				fmt.Printf("Error updating subscription check time: %v\n", err)
			}
			
			// If no new chapters, continue
			if len(filteredChapters) == 0 {
				fmt.Printf("No new chapters for \"%s\"\n", sub.MangaTitle)
				continue
			}
			
			// Get manga details
			/*
			manga, err := client.GetManga(sub.MangaID)
			if err != nil {
				fmt.Printf("Error getting manga details for \"%s\": %v\n", sub.MangaTitle, err)
				continue
			}
			*/
			
			// Group updates by user
			if _, ok := updates[sub.UserID]; !ok {
				updates[sub.UserID] = &UpdateInfo{
					MangaTitle: sub.MangaTitle,
					ChapterIDs: make([]string, 0),
					Chapters:   make([]api.Chapter, 0),
				}
			}
			
			// Add chapters to user updates
			updateInfo := updates[sub.UserID]
			for _, chapter := range filteredChapters {
				updateInfo.ChapterIDs = append(updateInfo.ChapterIDs, chapter.ID)
				updateInfo.Chapters = append(updateInfo.Chapters, chapter)
			}
			
			fmt.Printf("Found %d new chapter(s) for \"%s\"\n", len(filteredChapters), sub.MangaTitle)
		}
		
		// Send email notifications for updates
		for userID, updateInfo := range updates {
			// Get user
			user, err := database.GetUser(userID)
			if err != nil {
				fmt.Printf("Error getting user with ID %d: %v\n", userID, err)
				continue
			}
			
			// Get manga
			manga, err := client.GetManga(updateInfo.MangaTitle)
			if err != nil {
				fmt.Printf("Error getting manga details: %v\n", err)
				continue
			}
			
			// Connect to email server
			if err := emailService.Connect(); err != nil {
				fmt.Printf("Error connecting to email server: %v\n", err)
				continue
			}
			
			// Send notification
			if err := emailService.SendNotification(user.Email, manga, updateInfo.Chapters); err != nil {
				fmt.Printf("Error sending notification to %s: %v\n", user.Email, err)
			} else {
				fmt.Printf("Notification sent to %s about %d new chapter(s) for \"%s\"\n", 
					user.Email, len(updateInfo.Chapters), updateInfo.MangaTitle)
			}
			
			// Disconnect from email server
			emailService.Disconnect()
		}
		
		// Display summary
		if len(updates) == 0 {
			fmt.Println("No updates found for any subscriptions")
		} else {
			table := tablewriter.NewWriter(os.Stdout)
			table.SetHeader([]string{"User Email", "Manga", "New Chapters"})
			
			for userID, updateInfo := range updates {
				user, _ := database.GetUser(userID)
				row := []string{
					user.Email,
					updateInfo.MangaTitle,
					fmt.Sprintf("%d", len(updateInfo.Chapters)),
				}
				table.Append(row)
			}
			
			fmt.Println("\nUpdate Summary:")
			table.Render()
		}
		
		return nil
	},
}

func init() {
	// No specific flags for check command
}
