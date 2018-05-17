package notifications

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/aws/aws-sdk-go/service/ses/sesiface"
)

type notification interface {
	emails() ([]emailContent, error)
}

type emailContent struct {
	recipientEmail string
	subject        string
	htmlBody       string
	textBody       string
}

const sesRegion = "us-west-2"
const senderEmail = "noreply@dp3.us"
const charset = "UTF-8"

// SendNotification sends a one or more notifications for all supported mediums
func SendNotification(notification notification, svc sesiface.SESAPI) error {
	emails, err := notification.emails()
	if err != nil {
		return err
	}

	if svc == nil {
		session, err := session.NewSession(&aws.Config{
			Region: aws.String(sesRegion),
		})
		if err != nil {
			return err
		}
		svc = ses.New(session)
	}

	return sendEmails(emails, svc)
}

func sendEmails(emails []emailContent, svc sesiface.SESAPI) error {
	for _, email := range emails {
		input := &ses.SendEmailInput{
			Destination: &ses.Destination{
				ToAddresses: []*string{
					aws.String(email.recipientEmail),
				},
			},
			Message: &ses.Message{
				Body: &ses.Body{
					Html: &ses.Content{
						Charset: aws.String(charset),
						Data:    aws.String(email.htmlBody),
					},
					Text: &ses.Content{
						Charset: aws.String(charset),
						Data:    aws.String(email.textBody),
					},
				},
				Subject: &ses.Content{
					Charset: aws.String(charset),
					Data:    aws.String(email.subject),
				},
			},
			Source: aws.String(senderEmail),
		}

		// Returns the message ID. Should we store that somewhere?
		_, err := svc.SendEmail(input)

		if err != nil {
			return err
		}
	}

	return nil
}
