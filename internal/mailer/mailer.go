package mailer

import (
	"bytes"
	"embed"
	"fmt"
	"html/template"
	"log"
	"time"
)

const (
	FromName            = "GopherSocial"
	maxRetries          = 3
	UserWelcomeTemplate = "user_invitation.tmpl"
)

//go:embed "template"
var FS embed.FS

type Client interface {
	Send(templateFile, username, email string, data any, isSandBox bool) error
}

func parseTemplate(templateFile string, data any) (*bytes.Buffer, *bytes.Buffer, error) {
	tmpl, error := template.ParseFS(FS, "template/"+templateFile)
	if error != nil {
		return nil, nil, error
	}
	subject := new(bytes.Buffer)
	err := tmpl.ExecuteTemplate(subject, "subject", data)
	if err != nil {
		return nil, nil, error
	}
	body := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(body, "body", data)
	if err != nil {
		return nil, nil, error
	}

	return subject, body, nil
}

func handleRetries(email string, fn func() error) error {
	for i := range maxRetries {
		if err := fn(); err != nil {
			log.Printf("Failed to send email to %v, attempt %d of %d", email, i+1, maxRetries)
			log.Printf("Error: %v", err.Error())

			// exponential backoff which delays for period before send a mail again in case an error occurs.
			time.Sleep(time.Second * time.Duration(i+1))
			continue
		}
		return nil
	}

	return fmt.Errorf("failed to send email after %d attemps", maxRetries)
}
