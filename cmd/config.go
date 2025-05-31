package cmd

import (
	"fmt"
	"strings"

	"mangadex-cli/internal/email"

	"github.com/spf13/cobra"
)

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage configuration settings",
	Long:  `Commands for viewing and updating configuration settings.`,
}

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get [setting]",
	Short: "Get configuration values",
	Long: `Get the current value of a configuration setting.
If no setting is specified, all configuration values are shown.`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		// If no setting specified, show all (except sensitive data)
		if len(args) == 0 {
			fmt.Println("Current Configuration:")
			fmt.Printf("Database Path: %s\n", cfg.DatabasePath)
			fmt.Printf("MangaDex API URL: %s\n", cfg.MangaDexAPIURL)
			fmt.Printf("Update Check Interval: %d seconds\n", cfg.UpdateCheckInterval)
			
			// Show auth status but not the actual tokens
			if cfg.AuthToken != "" {
				fmt.Println("Authentication: Logged in")
			} else {
				fmt.Println("Authentication: Not logged in")
			}
			
			// Show email config without password
			fmt.Printf("SMTP Server: %s\n", cfg.SMTPSettings.Server)
			fmt.Printf("SMTP Port: %d\n", cfg.SMTPSettings.Port)
			fmt.Printf("SMTP Username: %s\n", cfg.SMTPSettings.Username)
			fmt.Printf("SMTP Use TLS: %t\n", cfg.SMTPSettings.UseTLS)
			fmt.Printf("SMTP From Email: %s\n", cfg.SMTPSettings.FromEmail)
			fmt.Printf("SMTP From Name: %s\n", cfg.SMTPSettings.FromName)
			
			return nil
		}
		
		// Get specific setting
		setting := strings.ToLower(args[0])
		switch setting {
		case "databasepath":
			fmt.Printf("Database Path: %s\n", cfg.DatabasePath)
		case "mangadexapiurl":
			fmt.Printf("MangaDex API URL: %s\n", cfg.MangaDexAPIURL)
		case "updatecheckinterval":
			fmt.Printf("Update Check Interval: %d seconds\n", cfg.UpdateCheckInterval)
		case "smtpserver":
			fmt.Printf("SMTP Server: %s\n", cfg.SMTPSettings.Server)
		case "smtpport":
			fmt.Printf("SMTP Port: %d\n", cfg.SMTPSettings.Port)
		case "smtpusername":
			fmt.Printf("SMTP Username: %s\n", cfg.SMTPSettings.Username)
		case "smtpusetls":
			fmt.Printf("SMTP Use TLS: %t\n", cfg.SMTPSettings.UseTLS)
		case "smtpfromemail":
			fmt.Printf("SMTP From Email: %s\n", cfg.SMTPSettings.FromEmail)
		case "smtpfromname":
			fmt.Printf("SMTP From Name: %s\n", cfg.SMTPSettings.FromName)
		default:
			return fmt.Errorf("unknown setting: %s", setting)
		}
		
		return nil
	},
}

// setCmd represents the set command
var setCmd = &cobra.Command{
	Use:   "set [setting] [value]",
	Short: "Update configuration settings",
	Long:  `Set a new value for a configuration setting.`,
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		setting := strings.ToLower(args[0])
		value := args[1]
		
		switch setting {
		case "databasepath":
			cfg.DatabasePath = value
			fmt.Printf("Database Path set to: %s\n", value)
		case "mangadexapiurl":
			cfg.MangaDexAPIURL = value
			fmt.Printf("MangaDex API URL set to: %s\n", value)
		case "updatecheckinterval":
			var interval int
			if _, err := fmt.Sscanf(value, "%d", &interval); err != nil {
				return fmt.Errorf("invalid interval value, must be a number: %w", err)
			}
			cfg.UpdateCheckInterval = interval
			fmt.Printf("Update Check Interval set to: %d seconds\n", interval)
		case "smtpserver":
			cfg.SMTPSettings.Server = value
			fmt.Printf("SMTP Server set to: %s\n", value)
		case "smtpport":
			var port int
			if _, err := fmt.Sscanf(value, "%d", &port); err != nil {
				return fmt.Errorf("invalid port value, must be a number: %w", err)
			}
			cfg.SMTPSettings.Port = port
			fmt.Printf("SMTP Port set to: %d\n", port)
		case "smtpusername":
			cfg.SMTPSettings.Username = value
			fmt.Printf("SMTP Username set to: %s\n", value)
		case "smtppassword":
			cfg.SMTPSettings.Password = value
			fmt.Println("SMTP Password updated")
		case "smtpusetls":
			var useTLS bool
			if strings.ToLower(value) == "true" {
				useTLS = true
			} else if strings.ToLower(value) == "false" {
				useTLS = false
			} else {
				return fmt.Errorf("invalid TLS setting, must be true or false")
			}
			cfg.SMTPSettings.UseTLS = useTLS
			fmt.Printf("SMTP Use TLS set to: %t\n", useTLS)
		case "smtpfromemail":
			cfg.SMTPSettings.FromEmail = value
			fmt.Printf("SMTP From Email set to: %s\n", value)
		case "smtpfromname":
			cfg.SMTPSettings.FromName = value
			fmt.Printf("SMTP From Name set to: %s\n", value)
		default:
			return fmt.Errorf("unknown setting: %s", setting)
		}
		
		// Save config
		if err := cfg.Save(cfgFile); err != nil {
			return fmt.Errorf("failed to save configuration: %w", err)
		}
		
		return nil
	},
}

// testCmd represents the test command
var testCmd = &cobra.Command{
	Use:   "test email [recipient]",
	Short: "Test configuration",
	Long:  `Test configuration components like email sending.`,
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		component := strings.ToLower(args[0])
		value := args[1]
		
		switch component {
		case "email":
			// Test email configuration
			emailService := email.NewEmailService(cfg.SMTPSettings)
			if err := emailService.Connect(); err != nil {
				return fmt.Errorf("failed to connect to email server: %w", err)
			}
			defer emailService.Disconnect()
			
			if err := emailService.SendTestEmail(value); err != nil {
				return fmt.Errorf("failed to send test email: %w", err)
			}
			
			fmt.Printf("Test email sent to %s\n", value)
			return nil
			
		default:
			return fmt.Errorf("unknown test component: %s", component)
		}
	},
}

func init() {
	configCmd.AddCommand(getCmd)
	configCmd.AddCommand(setCmd)
	configCmd.AddCommand(testCmd)
}