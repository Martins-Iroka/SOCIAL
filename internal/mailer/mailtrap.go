package mailer

import (
	"errors"
	"log"

	gomail "gopkg.in/mail.v2"
)

type mailtrapClient struct {
	fromEmail string
	apiKey    string
}

func NewMailTrapClient(apiKey, fromEmail string) (mailtrapClient, error) {
	if apiKey == "" {
		return mailtrapClient{}, errors.New("api key is required")
	}

	return mailtrapClient{
		fromEmail: fromEmail,
		apiKey:    apiKey,
	}, nil
}

func (m mailtrapClient) Send(templateFile, username, email string, data any, isSandbox bool) error {
	//Template parsing and building
	//template parsing and building
	subject, body, err := parseTemplate(templateFile, data)
	if err != nil {
		return err
	}

	message := gomail.NewMessage()
	message.SetHeader("From", "hello@demomailtrap.co")
	message.SetHeader("To", "martdev17@gmail.com")
	message.SetHeader("Subject", subject.String())

	message.AddAlternative("text/html", body.String())

	dialer := gomail.NewDialer("live.smtp.mailtrap.io", 587, "api", m.apiKey)
	err = handleRetries(email, func() error {
		if err := dialer.DialAndSend(message); err != nil {
			return err
		}
		log.Printf("Email - %v sent with status code %v", email, 200)
		return nil
	})

	return err
}
