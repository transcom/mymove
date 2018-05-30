package notifications

import (
	"bytes"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/aws/aws-sdk-go/service/ses/sesiface"
	"github.com/go-gomail/gomail"
	"github.com/pkg/errors"
)

type notification interface {
	emails() ([]emailContent, error)
}

type emailContent struct {
	attachments    []string
	recipientEmail string
	subject        string
	htmlBody       string
	textBody       string
}

// SendNotification sends a one or more notifications for all supported mediums
func SendNotification(notification notification, svc sesiface.SESAPI) error {
	emails, err := notification.emails()
	if err != nil {
		return err
	}

	return sendEmails(emails, svc)
}

func sendEmails(emails []emailContent, svc sesiface.SESAPI) error {
	for _, email := range emails {
		rawMessage, err := formatRawEmailMessage(email)
		if err != nil {
			return err
		}

		input := ses.SendRawEmailInput{
			Destinations: []*string{aws.String(email.recipientEmail)},
			RawMessage:   &ses.RawMessage{Data: rawMessage},
			Source:       aws.String(senderEmail()),
		}

		// Returns the message ID. Should we store that somewhere?
		fmt.Printf("input is: %v\n", input)
		fmt.Printf("svc is: %v\n", svc)
		//fmt.Printf("function is: %v\n", svc.SendRawEmail)

		_, err = svc.SendRawEmail(&input)
		if err != nil {
			return errors.Wrap(err, "Failed to send email using SES")
		}
	}

	return nil
}

func formatRawEmailMessage(email emailContent) ([]byte, error) {
	m := gomail.NewMessage()
	m.SetHeader("From", senderEmail())
	m.SetHeader("To", email.recipientEmail)
	m.SetHeader("Subject", email.subject)
	m.SetBody("text/plain", email.textBody)
	m.AddAlternative("text/html", email.htmlBody)
	for _, attachment := range email.attachments {
		m.Attach(attachment)
	}

	buf := new(bytes.Buffer)
	_, err := m.WriteTo(buf)
	if err != nil {
		return buf.Bytes(), errors.Wrap(err, "Failed to generate raw email notification message")
	}

	return buf.Bytes(), nil
}

// SenderEmail returns the "noreply" sender address for this environment
func senderEmail() string {
	return "noreply@" + os.Getenv("AWS_SES_DOMAIN")
}
