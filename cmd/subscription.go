package cmd

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"mangadex-cli/internal/api"
	"mangadex-cli/internal/db"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

var (
	mangaTitle    string
	mangaID       string
	userEmail     string
	subscriptionID int
	languages     string
)

// subscriptionCmd represents the subscription command
var subscriptionCmd = &cobra.Command{
	Use:   "subscription",
	Short: "Manage manga subscriptions",
	Long:  `Commands for adding, removing, and listing manga subscriptions.`,
}

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a manga subscription",
	Long: `Add a manga to your subscription list.
You can specify a manga by title (search) or by its MangaDex ID.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Initialize API client
		client := api.NewMangaDexClient(cfg.MangaDexAPIURL)
		
		// Set auth token if available
		if cfg.AuthToken != "" {
			client.SessionToken = cfg.AuthToken
			client.RefreshToken = cfg.RefreshToken
			client.TokenExpiry = cfg.TokenExpiry
		}
		
		// Validate required parameters
		if mangaTitle == "" && mangaID == "" {
			return fmt.Errorf("either manga title or manga ID must be provided")
		}
		
		if userEmail == "" {
			return fmt.Errorf("user email must be provided")
		}
		
		// Get user or create if doesn't exist
		user, err := database.GetUserByEmail(userEmail)
		if err != nil {
			// Create new user if not found
			user = &db.User{
				Email:     userEmail,
				Active:    true,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}
			if err := database.AddUser(user); err != nil {
				return fmt.Errorf("failed to create user: %w", err)
			}
		}
		
		// Parse language preferences
		var languageList []string
		if languages != "" {
			languageList = strings.Split(languages, ",")
			for i, lang := range languageList {
				languageList[i] = strings.TrimSpace(lang)
			}
		} else {
			// Default to English
			languageList = []string{"en"}
		}
		
		var manga *api.Manga
		// Search by title or get by ID
		if mangaID != "" {
			// Get manga by ID
			manga, err = client.GetManga(mangaID)
			if err != nil {
				return fmt.Errorf("failed to get manga with ID %s: %w", mangaID, err)
			}
		} else {
			// Search manga by title
			results, err := client.SearchManga(mangaTitle)
			if err != nil {
				return fmt.Errorf("failed to search for manga: %w", err)
			}
			
			if len(results) == 0 {
				return fmt.Errorf("no manga found matching '%s'", mangaTitle)
			}
			
			// Display results for selection
			fmt.Println("Multiple manga found. Please select one:")
			table := tablewriter.NewWriter(os.Stdout)
			table.SetHeader([]string{"#", "ID", "Title"})
			
			for i, m := range results {
				title := m.GetTitle()
				row := []string{
					strconv.Itoa(i + 1),
					m.ID,
					title,
				}
				table.Append(row)
			}
			table.Render()
			
			// Get user selection
			var selection int
			fmt.Print("Enter selection number: ")
			_, err = fmt.Scanln(&selection)
			if err != nil || selection < 1 || selection > len(results) {
				return fmt.Errorf("invalid selection")
			}
			
			manga = results[selection-1]
		}
		
		// Create subscription
		subscription := &db.Subscription{
			UserID:         user.ID,
			MangaID:        manga.ID,
			MangaTitle:     manga.GetTitle(),
			Languages:      strings.Join(languageList, ","),
			LastCheckTime:  time.Now(),
			LastChapterTime: time.Now(),
			Active:         true,
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		}
		
		// Add subscription to database
		if err := database.AddSubscription(subscription); err != nil {
			return fmt.Errorf("failed to add subscription: %w", err)
		}
		
		fmt.Printf("Successfully subscribed to \"%s\" for %s\n", manga.GetTitle(), userEmail)
		return nil
	},
}

// removeCmd represents the remove command
var removeCmd = &cobra.Command{
	Use:   "remove",
	Short: "Remove a manga subscription",
	Long:  `Remove a manga from your subscription list.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Validate required parameters
		if subscriptionID == 0 {
			return fmt.Errorf("subscription ID must be provided")
		}
		
		// Get subscription to confirm
		subscription, err := database.GetSubscription(subscriptionID)
		if err != nil {
			return fmt.Errorf("failed to find subscription with ID %d: %w", subscriptionID, err)
		}
		
		// Confirm deletion
		fmt.Printf("Are you sure you want to remove the subscription for \"%s\"? (y/n): ", subscription.MangaTitle)
		var confirm string
		fmt.Scanln(&confirm)
		
		if strings.ToLower(confirm) != "y" && strings.ToLower(confirm) != "yes" {
			fmt.Println("Subscription removal cancelled")
			return nil
		}
		
		// Delete subscription
		if err := database.DeleteSubscription(subscriptionID); err != nil {
			return fmt.Errorf("failed to remove subscription: %w", err)
		}
		
		fmt.Println("Subscription removed successfully")
		return nil
	},
}

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List manga subscriptions",
	Long:  `List all manga subscriptions.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Get all subscriptions or filter by email
		var subscriptions []db.Subscription
		var err error
		
		if userEmail != "" {
			// Get user by email
			user, err := database.GetUserByEmail(userEmail)
			if err != nil {
				return fmt.Errorf("failed to find user with email %s: %w", userEmail, err)
			}
			
			// Get subscriptions for user
			subscriptions, err = database.GetSubscriptionsByUserID(user.ID)
			if err != nil {
				return fmt.Errorf("failed to get subscriptions: %w", err)
			}
		} else {
			// Get all subscriptions
			subscriptions, err = database.ListSubscriptions()
			if err != nil {
				return fmt.Errorf("failed to list subscriptions: %w", err)
			}
		}
		
		if len(subscriptions) == 0 {
			if userEmail != "" {
				fmt.Printf("No subscriptions found for %s\n", userEmail)
			} else {
				fmt.Println("No subscriptions found")
			}
			return nil
		}
		
		// Display subscriptions
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"ID", "Manga Title", "User Email", "Languages", "Last Check", "Status"})
		
		for _, sub := range subscriptions {
			// Get user email
			user, err := database.GetUser(sub.UserID)
			if err != nil {
				// Use placeholder if user not found
				user = &db.User{Email: "unknown"}
			}
			
			status := "Active"
			if !sub.Active {
				status = "Inactive"
			}
			
			row := []string{
				strconv.Itoa(sub.ID),
				sub.MangaTitle,
				user.Email,
				sub.Languages,
				sub.LastCheckTime.Format("2006-01-02 15:04"),
				status,
			}
			table.Append(row)
		}
		table.Render()
		
		return nil
	},
}

func init() {
	subscriptionCmd.AddCommand(addCmd)
	subscriptionCmd.AddCommand(removeCmd)
	subscriptionCmd.AddCommand(listCmd)
	
	// Add flags for add command
	addCmd.Flags().StringVarP(&mangaTitle, "title", "t", "", "Manga title to search for")
	addCmd.Flags().StringVarP(&mangaID, "id", "i", "", "MangaDex manga ID")
	addCmd.Flags().StringVarP(&userEmail, "email", "e", "", "User email address")
	addCmd.Flags().StringVarP(&languages, "languages", "l", "en", "Comma-separated language codes (e.g., 'en,es,fr')")
	
	// Add flags for remove command
	removeCmd.Flags().IntVarP(&subscriptionID, "id", "i", 0, "Subscription ID to remove")
	removeCmd.MarkFlagRequired("id")
	
	// Add flags for list command
	listCmd.Flags().StringVarP(&userEmail, "email", "e", "", "Filter subscriptions by user email")
}
