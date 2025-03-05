package pwdless

import (
	"os"
	"time"

	"github.com/dhax/go-base/email"
)

// ContentLoginToken defines content for login token email template.
type ContentLoginToken struct {
	Email  string
	Name   string
	URL    string
	Token  string
	Expiry time.Time
}

// LoginTokenEmail creates and sends a login token email with provided template content.
func LoginTokenEmail(name, address string, content ContentLoginToken) email.Message {
	return email.Message{
		From:     email.NewEmail(os.Getenv("EMAIL_FROM_NAME"), os.Getenv("EMAIL_FROM_ADDRESS")),
		To:       email.NewEmail(name, address),
		Subject:  "Login Token",
		Template: "loginToken",
		Content:  content,
	}
}
