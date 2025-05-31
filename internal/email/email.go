package email

import (
	"fmt"
	"mangadex-cli/internal/api"
	"mangadex-cli/internal/config"
	"strings"
	"time"
	
	"github.com/go-gomail/gomail"
)

// EmailService handles sending email notifications
type EmailService struct {
	Config config.SMTPConfig
}

// NewEmailService creates a new email service
func NewEmailService(config config.SMTPConfig) *EmailService {
	return &EmailService{
		Config: config,
	}
}

// Connect tests the connection to the email server
func (e *EmailService) Connect() error {
	dialer := e.createDialer()
	sender, err := dialer.Dial()
	if err != nil {
		return fmt.Errorf("failed to connect to email server: %w", err)
	}
	defer sender.Close()
	
	return nil
}

// Disconnect is a no-op as we don't maintain persistent connections
func (e *EmailService) Disconnect() error {
	// No persistent connection to close
	return nil
}

// createDialer creates a new gomail dialer with the configured settings
func (e *EmailService) createDialer() *gomail.Dialer {
	dialer := gomail.NewDialer(e.Config.Server, e.Config.Port, e.Config.Username, e.Config.Password)
	dialer.SSL = e.Config.UseTLS
	
	return dialer
}

// SendTestEmail sends a test email to verify configuration
func (e *EmailService) SendTestEmail(recipient string) error {
	// Create message
	m := gomail.NewMessage()
	m.SetHeader("From", e.createFromHeader())
	m.SetHeader("To", recipient)
	m.SetHeader("Subject", "MangaDex CLI Notification - Test Email")
	
	// Email body
	body := `
	<html>
		<head>
			<style>
				body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
				.container { max-width: 600px; margin: 0 auto; padding: 20px; }
				.header { background-color: #4a86e8; color: white; padding: 10px; text-align: center; }
				.footer { font-size: 12px; color: #777; margin-top: 30px; text-align: center; }
			</style>
		</head>
		<body>
			<div class="container">
				<div class="header">
					<h1>MangaDex CLI Notification</h1>
				</div>
				<div class="content">
					<h2>Test Email</h2>
					<p>This is a test email sent from your MangaDex CLI Notification Service.</p>
					<p>If you're receiving this message, your email configuration is working correctly!</p>
					<p>Time: %s</p>
				</div>
				<div class="footer">
					<p>This email was sent from the MangaDex CLI Notification Service.</p>
				</div>
			</div>
		</body>
	</html>
	`
	
	formattedBody := fmt.Sprintf(body, time.Now().Format(time.RFC1123))
	
	m.SetBody("text/html", formattedBody)
	m.AddAlternative("text/plain", fmt.Sprintf(
		"MangaDex CLI Notification - Test Email\n\n"+
			"This is a test email sent from your MangaDex CLI Notification Service.\n"+
			"If you're receiving this message, your email configuration is working correctly!\n\n"+
			"Time: %s\n\n"+
			"This email was sent from the MangaDex CLI Notification Service.",
		time.Now().Format(time.RFC1123)))
	
	// Send the email
	dialer := e.createDialer()
	if err := dialer.DialAndSend(m); err != nil {
		return fmt.Errorf("failed to send test email: %w", err)
	}
	
	return nil
}

// SendNotification sends a manga update notification
func (e *EmailService) SendNotification(recipient string, manga *api.Manga, chapters []api.Chapter) error {
	// Create message
	m := gomail.NewMessage()
	m.SetHeader("From", e.createFromHeader())
	m.SetHeader("To", recipient)
	
	// Subject line
	var subject string
	if len(chapters) == 1 {
		subject = fmt.Sprintf("New Chapter: %s - Chapter %s", manga.GetTitle(), chapters[0].Chapter)
	} else {
		subject = fmt.Sprintf("%d New Chapters for %s", len(chapters), manga.GetTitle())
	}
	m.SetHeader("Subject", subject)
	
	// Generate HTML and text content
	html, text := e.renderNotificationTemplate(manga, chapters)
	
	m.SetBody("text/html", html)
	m.AddAlternative("text/plain", text)
	
	// Send the email
	dialer := e.createDialer()
	if err := dialer.DialAndSend(m); err != nil {
		return fmt.Errorf("failed to send notification: %w", err)
	}
	
	return nil
}

// createFromHeader creates the From header with proper formatting
func (e *EmailService) createFromHeader() string {
	if e.Config.FromName != "" {
		return fmt.Sprintf("%s <%s>", e.Config.FromName, e.Config.FromEmail)
	}
	return e.Config.FromEmail
}

// renderNotificationTemplate generates the email content for a notification
func (e *EmailService) renderNotificationTemplate(manga *api.Manga, chapters []api.Chapter) (string, string) {
	// Sort chapters in ascending order
	// This is a simple implementation; for production, we would want to sort numerically
	// but chapter numbers can be complex ("10.5", "Extra", etc.)
	
	// HTML Template
	htmlTemplate := `
	<html>
		<head>
			<style>
				body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
				.container { max-width: 600px; margin: 0 auto; padding: 20px; }
				.header { background-color: #4a86e8; color: white; padding: 10px; text-align: center; }
				.manga-info { display: flex; margin: 20px 0; }
				.manga-cover { width: 120px; height: auto; margin-right: 20px; }
				.manga-details { flex: 1; }
				.chapter-list { margin: 20px 0; }
				.chapter { padding: 10px; border-bottom: 1px solid #eee; }
				.chapter:last-child { border-bottom: none; }
				.chapter-number { font-weight: bold; }
				.footer { font-size: 12px; color: #777; margin-top: 30px; text-align: center; }
				.read-button { display: inline-block; background-color: #4a86e8; color: white; padding: 8px 15px; text-decoration: none; border-radius: 3px; }
			</style>
		</head>
		<body>
			<div class="container">
				<div class="header">
					<h1>MangaDex Update</h1>
				</div>
				<div class="content">
					<h2>New Chapters for %s</h2>
					
					<div class="manga-info">
						%s
						<div class="manga-details">
							<p><strong>Status:</strong> %s</p>
							<p>%s</p>
						</div>
					</div>
					
					<div class="chapter-list">
						<h3>New Chapters:</h3>
						%s
					</div>
					
					<p><a href="https://mangadex.org/title/%s" class="read-button">View on MangaDex</a></p>
				</div>
				<div class="footer">
					<p>This email was sent from the MangaDex CLI Notification Service.</p>
					<p>Time: %s</p>
				</div>
			</div>
		</body>
	</html>
	`
	
	// Generate cover image HTML if available
	var coverHTML string
	if manga.CoverArtURL != "" {
		coverHTML = fmt.Sprintf(`<img src="%s" class="manga-cover" alt="%s Cover">`, manga.CoverArtURL, manga.GetTitle())
	}
	
	// Generate chapter list HTML
	var chapterListHTML strings.Builder
	for _, chapter := range chapters {
		chapterHTML := `<div class="chapter">
			<span class="chapter-number">Chapter %s</span>
			%s
			<p>Language: %s</p>
			<p>Published: %s</p>
			<p><a href="https://mangadex.org/chapter/%s">Read Chapter</a></p>
		</div>`
		
		title := ""
		if chapter.Title != "" {
			title = fmt.Sprintf(`<p>%s</p>`, chapter.Title)
		}
		
		chapterListHTML.WriteString(fmt.Sprintf(
			chapterHTML,
			chapter.Chapter,
			title,
			chapter.TranslatedLanguage,
			chapter.PublishAt.Format("January 2, 2006"),
			chapter.ID,
		))
	}
	
	// Format the HTML template
	html := fmt.Sprintf(
		htmlTemplate,
		manga.GetTitle(),
		coverHTML,
		manga.Status,
		manga.GetDescription(),
		chapterListHTML.String(),
		manga.ID,
		time.Now().Format(time.RFC1123),
	)
	
	// Text version
	textTemplate := `
MangaDex Update - %s

New Chapters for %s

Status: %s

%s

New Chapters:
%s

View on MangaDex: https://mangadex.org/title/%s

This email was sent from the MangaDex CLI Notification Service.
Time: %s
`
	
	// Generate chapter list text
	var chapterListText strings.Builder
	for _, chapter := range chapters {
		titleText := ""
		if chapter.Title != "" {
			titleText = fmt.Sprintf(" - %s", chapter.Title)
		}
		
		chapterListText.WriteString(fmt.Sprintf(
			"- Chapter %s%s | Language: %s | Published: %s | https://mangadex.org/chapter/%s\n",
			chapter.Chapter,
			titleText,
			chapter.TranslatedLanguage,
			chapter.PublishAt.Format("January 2, 2006"),
			chapter.ID,
		))
	}
	
	// Format the text template
	text := fmt.Sprintf(
		textTemplate,
		manga.GetTitle(),
		manga.GetTitle(),
		manga.Status,
		manga.GetDescription(),
		chapterListText.String(),
		manga.ID,
		time.Now().Format(time.RFC1123),
	)
	
	return html, text
}