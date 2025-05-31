package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"mangadex-cli/internal/config"
	"mangadex-cli/internal/db"

	"github.com/spf13/cobra"
)

var (
	cfgFile string
	cfg     *config.Config
	database *db.DB
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "mangadex-cli",
	Short: "MangaDex CLI notification service",
	Long: `A CLI tool that notifies you via email when your favorite manga on MangaDex have new chapters.
	
The tool allows you to subscribe to manga series and receive email notifications
when new chapters are available. You can manage your subscriptions, 
configure email settings, and test the notification service through this CLI.`,
	// This is called before any subcommand
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// Skip initialization for these commands as they don't require full setup
		if cmd.CommandPath() == "mangadex-cli help" {
			return nil
		}

		var err error
		
		// Create config directory if it doesn't exist
		configDir := filepath.Dir(cfgFile)
		if _, err := os.Stat(configDir); os.IsNotExist(err) {
			if err := os.MkdirAll(configDir, 0755); err != nil {
				return fmt.Errorf("failed to create config directory: %w", err)
			}
		}
		
		// Load configuration
		cfg, err = config.Load(cfgFile)
		if err != nil {
			// For first run, create default config
			if os.IsNotExist(err) {
				cfg = config.Default()
				if err := cfg.Save(cfgFile); err != nil {
					return fmt.Errorf("failed to create default config: %w", err)
				}
				fmt.Printf("Default configuration created at %s\n", cfgFile)
			} else {
				return fmt.Errorf("failed to load config: %w", err)
			}
		}
		
		// Initialize database
		database, err = db.NewDB(cfg.DatabasePath)
		if err != nil {
			return fmt.Errorf("failed to initialize database: %w", err)
		}
		
		return nil
	},
	
	// Cleanup on exit
	PersistentPostRunE: func(cmd *cobra.Command, args []string) error {
		if database != nil {
			if err := database.Close(); err != nil {
				return fmt.Errorf("failed to close database: %w", err)
			}
		}
		return nil
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	// Set default config file location to user's home directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not determine home directory: %v\n", err)
		os.Exit(1)
	}
	
	defaultConfigPath := filepath.Join(homeDir, ".mangadex-cli", "config.yaml")
	
	// Define persistent flags
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", defaultConfigPath, "config file path")
	
	// Add subcommands
	rootCmd.AddCommand(authCmd)
	rootCmd.AddCommand(configCmd)
	rootCmd.AddCommand(subscriptionCmd)
	rootCmd.AddCommand(serviceCmd)
	rootCmd.AddCommand(checkCmd)
}