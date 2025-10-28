package mailer

import (
	"log"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type SendGridMailer struct {
	fromEmail string
	apiKey    string
	client    *sendgrid.Client
}

// creating a constructor to initialize the struct
func NewSendgrid(apiKey, fromEmail string) *SendGridMailer {
	client := sendgrid.NewSendClient(apiKey)

	return &SendGridMailer{
		fromEmail: fromEmail,
		apiKey:    apiKey,
		client:    client,
	}
}

func (m *SendGridMailer) Send(templateFile, username, email string, data any, isSandBox bool) error {
	log.Printf("sandbox is %v", isSandBox)
	from := mail.NewEmail(FromName, m.fromEmail)
	to := mail.NewEmail(username, email)

	//template parsing and building
	subject, body, err := parseTemplate(templateFile, data)
	if err != nil {
		return err
	}

	message := mail.NewSingleEmail(from, subject.String(), to, "", body.String())

	message.SetMailSettings(
		&mail.MailSettings{
			SandboxMode: &mail.Setting{
				Enable: &isSandBox,
			},
		},
	)

	// for i := range maxRetries {
	// 	response, err := m.client.Send(message)
	// 	if err != nil {
	// 		log.Printf("Failed to send email to %v, attempt %d of %d", email, i+1, maxRetries)
	// 		log.Printf("Error: %v", err.Error())

	// 		// exponential backoff which delays for period before send a mail again in case an error occurs.
	// 		time.Sleep(time.Second * time.Duration(i+1))
	// 		continue
	// 	}

	// 	log.Printf("Email - %v sent with status code %v", email, response.StatusCode)
	// 	return nil
	// }

	// return fmt.Errorf("failed to send email after %d attemps", maxRetries)

	err = handleRetries(email, func() error {
		response, err := m.client.Send(message)
		if err != nil {
			return err
		}
		log.Printf("Email - %v sent with status code %v", email, response.StatusCode)
		return nil
	})

	return err
}
