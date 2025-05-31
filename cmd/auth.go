package cmd

import (
	"fmt"
	"syscall"
	"time"

	"mangadex-cli/internal/api"

	"github.com/spf13/cobra"
	"golang.org/x/term"
)

var (
	username string
	password string
)

// authCmd represents the auth command
var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "Manage MangaDex authentication",
	Long:  `Commands for logging in and out of MangaDex.`,
}

// loginCmd represents the login command
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login to MangaDex",
	Long: `Login to your MangaDex account to access subscription functionality.
The login process uses OAuth 2.0 to authenticate with MangaDex.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// If username is not provided via flag, prompt for it
		if username == "" {
			fmt.Print("MangaDex username: ")
			fmt.Scanln(&username)
		}

		// If password is not provided via flag, prompt for it securely
		if password == "" {
			fmt.Print("Password: ")
			bytePassword, err := term.ReadPassword(int(syscall.Stdin))
			if err != nil {
				return fmt.Errorf("failed to read password: %w", err)
			}
			fmt.Println() // Add newline after password input
			password = string(bytePassword)
		}

		// Initialize API client
		client := api.NewMangaDexClient(cfg.MangaDexAPIURL)

		// Perform login
		if err := client.Login(username, password); err != nil {
			return fmt.Errorf("login failed: %w", err)
		}

		// Store auth tokens in config
		cfg.AuthToken = client.SessionToken
		cfg.RefreshToken = client.RefreshToken
		cfg.TokenExpiry = client.TokenExpiry

		// Save config
		if err := cfg.Save(cfgFile); err != nil {
			return fmt.Errorf("failed to save authentication tokens: %w", err)
		}

		fmt.Println("Login successful! Authentication tokens saved.")
		return nil
	},
}

// logoutCmd represents the logout command
var logoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Logout from MangaDex",
	Long:  `Logout from your MangaDex account and remove authentication tokens.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Check if already logged in
		if cfg.AuthToken == "" && cfg.RefreshToken == "" {
			fmt.Println("You are not logged in.")
			return nil
		}

		// Clear auth tokens
		cfg.AuthToken = ""
		cfg.RefreshToken = ""
		cfg.TokenExpiry = time.Now()

		// Save config
		if err := cfg.Save(cfgFile); err != nil {
			return fmt.Errorf("failed to remove authentication tokens: %w", err)
		}

		fmt.Println("Logout successful. Authentication tokens removed.")
		return nil
	},
}

func init() {
	authCmd.AddCommand(loginCmd)
	authCmd.AddCommand(logoutCmd)

	// Define flags for login command
	loginCmd.Flags().StringVarP(&username, "username", "u", "", "MangaDex username")
	loginCmd.Flags().StringVarP(&password, "password", "p", "", "MangaDex password (not recommended, use interactive prompt instead)")
}
