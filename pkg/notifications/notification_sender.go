package notifications

import (
	"bytes"
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/aws/aws-sdk-go/service/ses/sesiface"
	"github.com/go-gomail/gomail"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

type notification interface {
	emails(ctx context.Context) ([]emailContent, error)
}

type emailContent struct {
	attachments    []string
	recipientEmail string
	subject        string
	htmlBody       string
	textBody       string
}

// NotificationSender is an interface for sending notifications
type NotificationSender interface {
	SendNotification(ctx context.Context, notification notification) error
}

// NotificationSendingContext provides context to a notification sender
type NotificationSendingContext struct {
	svc    sesiface.SESAPI
	domain string
	logger Logger
}

// NewNotificationSender returns a new NotificationSendingContext
func NewNotificationSender(svc sesiface.SESAPI, domain string, logger Logger) NotificationSendingContext {
	return NotificationSendingContext{
		svc:    svc,
		domain: domain,
		logger: logger,
	}
}

// SendNotification sends a one or more notifications for all supported mediums
func (n NotificationSendingContext) SendNotification(ctx context.Context, notification notification) error {
	emails, err := notification.emails(ctx)
	if err != nil {
		return err
	}

	return sendEmails(emails, n.svc, n.domain, n.logger)
}

func sendEmails(emails []emailContent, svc sesiface.SESAPI, domain string, logger Logger) error {
	for _, email := range emails {
		rawMessage, err := formatRawEmailMessage(email, domain)
		if err != nil {
			return err
		}

		input := ses.SendRawEmailInput{
			Destinations: []*string{aws.String(email.recipientEmail)},
			RawMessage:   &ses.RawMessage{Data: rawMessage},
			Source:       aws.String(senderEmail(domain)),
		}

		// Returns the message ID. Should we store that somewhere?
		_, err = svc.SendRawEmail(&input)
		if err != nil {
			return errors.Wrap(err, "Failed to send email using SES")
		}

		logger.Info("Sent email to service member",
			zap.String("service member email address", email.recipientEmail))
	}

	return nil
}

func formatRawEmailMessage(email emailContent, domain string) ([]byte, error) {
	m := gomail.NewMessage()
	m.SetHeader("From", senderEmail(domain))
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

func senderEmail(domain string) string {
	return "noreply@" + domain
}
