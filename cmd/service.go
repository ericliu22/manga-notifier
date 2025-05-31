package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"mangadex-cli/internal/api"
	"mangadex-cli/internal/email"
	"mangadex-cli/internal/scheduler"

	"github.com/spf13/cobra"
)

var (
	daemonize  bool
	foreground bool
)

// serviceCmd represents the service command
var serviceCmd = &cobra.Command{
	Use:   "service",
	Short: "Manage notification service",
	Long:  `Commands for starting, stopping, and checking the notification service.`,
}

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start the notification service",
	Long: `Start the notification service that checks for manga updates.
The service can run in the foreground or as a background process.`,
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
		
		// Initialize scheduler
		sched := scheduler.NewCronScheduler(
			database,
			client,
			emailService,
			cfg.UpdateCheckInterval,
		)
		
		// Start the scheduler
		if err := sched.Start(); err != nil {
			return fmt.Errorf("failed to start scheduler: %w", err)
		}
		
		fmt.Printf("Notification service started. Checking for updates every %d seconds\n", cfg.UpdateCheckInterval)
		
		// If running in the foreground, wait for interruption
		if foreground {
			fmt.Println("Press Ctrl+C to stop the service")
			
			// Set up signal handling
			sigChan := make(chan os.Signal, 1)
			signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
			
			// Wait for signal
			<-sigChan
			
			// Stop the scheduler
			if err := sched.Stop(); err != nil {
				return fmt.Errorf("failed to stop scheduler: %w", err)
			}
			
			fmt.Println("Notification service stopped")
		}
		
		return nil
	},
}

// statusCmd represents the status command
var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Check notification service status",
	Long:  `Check if the notification service is running.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// In a real implementation, we would check if the service is running
		// For now, we'll just display a message since we don't have a persistent service
		fmt.Println("Service status checking is not fully implemented yet.")
		fmt.Println("To run the service in the foreground, use: mangadex-cli service start --foreground")
		return nil
	},
}

func init() {
	serviceCmd.AddCommand(startCmd)
	serviceCmd.AddCommand(statusCmd)
	
	// Add flags for start command
	startCmd.Flags().BoolVarP(&daemonize, "daemon", "d", false, "Run as a daemon (background process)")
	startCmd.Flags().BoolVarP(&foreground, "foreground", "f", false, "Run in the foreground")
}
