package handlers

import (
	"github.com/wneessen/go-mail"
	"log"
)

func SendEmail(to string, subject string, body string) error {
	message := mail.NewMsg()
	if err := message.From("no-reply@manganotifier.com"); err != nil {
		log.Fatalf("failed to set From address: %s", err)
	}
	if err := message.To(to); err != nil {
		log.Fatalf("failed to set To address: %s", err)
	}

	message.Subject(subject)
	//Can customize body of mail
	message.SetBodyString(mail.TypeTextPlain, body)

	return nil
}
