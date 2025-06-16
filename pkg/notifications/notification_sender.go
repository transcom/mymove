package notifications

import (
	"bytes"
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ses"
	"github.com/aws/aws-sdk-go-v2/service/ses/types"
	"github.com/go-gomail/gomail"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/cli"
	"github.com/transcom/mymove/pkg/uploader"
)

// Notification is an interface for creating emails
type Notification interface {
	emails(appCtx appcontext.AppContext) ([]emailContent, error)
}

type RawEmailSender interface {
	SendRawEmail(ctx context.Context, params *ses.SendRawEmailInput, optFns ...func(*ses.Options)) (*ses.SendRawEmailOutput, error)
}

type emailContent struct {
	attachments     []string
	recipientEmail  string
	recipientEmails []string
	subject         string
	htmlBody        string
	textBody        string
	onSuccess       func(string) error
}

// NotificationSender is an interface for sending notifications
//
//go:generate mockery --name NotificationSender
type NotificationSender interface {
	SendNotification(appCtx appcontext.AppContext, notification Notification) error
}

// NotificationSendingContext provides context to a notification sender
type NotificationSendingContext struct {
	svc           RawEmailSender
	domain        string
	sysAdminEmail string
}

// NewNotificationSender returns a new NotificationSendingContext
func NewNotificationSender(svc RawEmailSender, domain string, sysAdminEmail string) NotificationSendingContext {
	return NotificationSendingContext{
		svc:           svc,
		domain:        domain,
		sysAdminEmail: sysAdminEmail,
	}
}

// SendNotification sends a one or more notifications for all supported mediums
func (n NotificationSendingContext) SendNotification(appCtx appcontext.AppContext, notification Notification) error {
	emails, err := notification.emails(appCtx)
	if err != nil {
		return err
	}

	return sendEmails(appCtx, emails, n.svc, n.domain)
}

// InitEmail initializes the email backend
func InitEmail(v *viper.Viper, logger *zap.Logger) (NotificationSender, error) {
	if v.GetString(cli.EmailBackendFlag) == "ses" {
		// Setup Amazon SES (email) service TODO: This might be able
		// to be combined with the AWS Session that we're using for S3
		// down below.

		awsSESRegion := v.GetString(cli.AWSSESRegionFlag)
		awsSESDomain := v.GetString(cli.AWSSESDomainFlag)
		sysAdminEmail := v.GetString(cli.SysAdminEmail)
		logger.Info("Using ses email backend",
			zap.String("region", awsSESRegion),
			zap.String("domain", awsSESDomain))
		cfg, err := config.LoadDefaultConfig(context.Background(),
			config.WithRegion(awsSESRegion),
		)
		if err != nil {
			logger.Fatal("error loading ses aws config", zap.Error(err))
		}

		sesService := ses.NewFromConfig(cfg)
		input := &ses.GetAccountSendingEnabledInput{}
		result, err := sesService.GetAccountSendingEnabled(context.Background(), input)
		if err != nil || result == nil || !result.Enabled {
			logger.Error("email sending not enabled", zap.Error(err))
			return NewNotificationSender(nil, awsSESDomain, sysAdminEmail), err
		}
		return NewNotificationSender(sesService, awsSESDomain, sysAdminEmail), nil
	}

	domain := "milmovelocal"
	logger.Info("Using local email backend", zap.String("domain", domain))
	return NewStubNotificationSender(domain), nil
}

// GetSysAdminEmail returns the System Administrators' email address that has been set in the NotificationSender
func GetSysAdminEmail(sender NotificationSender) (email string) {
	if senderContext, ok := sender.(NotificationSendingContext); ok {
		email = senderContext.sysAdminEmail
	}
	return email
}

func sendEmails(appCtx appcontext.AppContext, emails []emailContent, svc RawEmailSender, domain string) error {
	for i, email := range emails {
		rawMessage, err := formatRawEmailMessage(email, domain)
		if err != nil {
			return err
		}
		var destinations []string
		if len(email.recipientEmails) > 0 {
			destinations = email.recipientEmails
		} else {
			destinations = []string{email.recipientEmail}
		}

		input := ses.SendRawEmailInput{
			Destinations: destinations,
			RawMessage:   &types.RawMessage{Data: rawMessage},
			Source:       aws.String(senderEmail(domain)),
		}

		// Returns the message ID. Should we store that somewhere?
		sendRawEmailOutput, err := svc.SendRawEmail(context.Background(), &input)
		if err != nil {
			return errors.Wrap(err, "Failed to send email using SES")
		}
		if email.onSuccess != nil && sendRawEmailOutput.MessageId != nil {
			err := email.onSuccess(*sendRawEmailOutput.MessageId)
			if err != nil {
				appCtx.Logger().Error("email.onSuccess error", zap.Error(err))
			}
		}
		// rate limited if exceed > 80 emails / second. delay to prevent hitting the limit
		if i > 0 {
			time.Sleep(20 * time.Millisecond)
		}
	}

	return nil
}

func formatRawEmailMessage(email emailContent, domain string) ([]byte, error) {
	m := gomail.NewMessage()
	m.SetHeader("From", senderEmail(domain))
	if len(email.recipientEmails) > 0 {
		m.SetHeader("To", email.recipientEmails...)
	} else {
		m.SetHeader("To", email.recipientEmail)
	}
	m.SetHeader("Subject", email.subject)
	m.SetBody(uploader.FileTypeText, email.textBody)
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
