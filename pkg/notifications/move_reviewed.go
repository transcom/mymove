package notifications

import (
	"bytes"
	"context"
	"fmt"
	html "html/template"
	text "text/template"
	"time"

	"github.com/transcom/mymove/pkg/models"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/assets"

	"github.com/gobuffalo/pop"
	"go.uber.org/zap"
)

const surveyLink = "https://www.surveymonkey.com/r/MilMovePt3-08191"

var (
	moveReviewedRawTextTemplate = string(assets.MustAsset("pkg/notifications/templates/move_reviewed_template.txt"))
	textTemplate                = text.Must(text.New("text_template").Parse(moveReviewedRawTextTemplate))
	moveReviewedRawHTMLTemplate = string(assets.MustAsset("pkg/notifications/templates/move_reviewed_template.html"))
	HTMLTemplate                = html.Must(html.New("text_template").Parse(moveReviewedRawHTMLTemplate))
)

// MoveReviewed has notification content for completed/reviewed moves
type MoveReviewed struct {
	db           *pop.Connection
	logger       Logger
	date         time.Time
	htmlTemplate *html.Template
	textTemplate *text.Template
}

// NewMoveReviewed returns a new move submitted notification
func NewMoveReviewed(db *pop.Connection, logger Logger, date time.Time) (*MoveReviewed, error) {

	return &MoveReviewed{
		db:           db,
		logger:       logger,
		date:         date,
		htmlTemplate: HTMLTemplate,
		textTemplate: textTemplate,
	}, nil
}

type EmailInfos []EmailInfo

type EmailInfo struct {
	ServiceMemberID    uuid.UUID `db:"id"`
	Email              *string   `db:"personal_email"`
	DutyStationName    string    `db:"duty_station_name"`
	NewDutyStationName string    `db:"new_duty_station_name"`
}

func (m MoveReviewed) GetEmailInfo(date time.Time) (EmailInfos, error) {
	dateString := date.Format("2006-01-02")
	query := `SELECT sm.id, sm.personal_email, dsn.name AS new_duty_station_name, dso.name AS duty_station_name
FROM personally_procured_moves p
         JOIN moves m ON p.move_id = m.id
         JOIN orders o ON m.orders_id = o.id
         JOIN service_members sm ON o.service_member_id = sm.id
         JOIN duty_stations dso ON sm.duty_station_id = dso.id
         JOIN duty_stations dsn ON o.new_duty_station_id = dsn.id
         LEFT JOIN notifications n ON sm.id = n.service_member_id
WHERE CAST(reviewed_date AS date) = $1
--  send email if haven't sent them a MOVE_REVIEWED_EMAIL yet OR we haven't sent them any emails at all
    AND (notification_type != 'MOVE_REVIEWED_EMAIL' OR n.service_member_id IS NULL);`

	emailInfos := EmailInfos{}
	err := m.db.RawQuery(query, dateString).All(&emailInfos)
	return emailInfos, err
}

// NotificationSendingContext expects a `notification` with an `emails` method,
// so we implement `email` to satisfy that interface
func (m MoveReviewed) emails(ctx context.Context) ([]emailContent, error) {
	emailInfos, err := m.GetEmailInfo(m.date)
	if err != nil {
		m.logger.Error("error retrieving email info for", zap.String("date", m.date.String()))
		return []emailContent{}, err
	}
	if len(emailInfos) == 0 {
		m.logger.Info("no emails to be sent for", zap.String("date", m.date.String()))
		return []emailContent{}, nil
	}
	return m.formatEmails(emailInfos)
}

// formatEmails formats email data using both html and text template
func (m MoveReviewed) formatEmails(emailInfos EmailInfos) ([]emailContent, error) {
	var emails []emailContent
	for _, emailInfo := range emailInfos {
		htmlBody, textBody, err := m.renderTemplates(moveReviewedEmailData{
			Link:                   surveyLink,
			OriginDutyStation:      emailInfo.DutyStationName,
			DestinationDutyStation: emailInfo.NewDutyStationName,
		})
		if err != nil {
			m.logger.Error("error rendering template", zap.Error(err))
			continue
		}
		if emailInfo.Email == nil {
			m.logger.Info("no email found for service member")
			continue
		}
		smEmail := emailContent{
			recipientEmail: *emailInfo.Email,
			subject:        "[MilMove] Let us know how we did",
			htmlBody:       htmlBody,
			textBody:       textBody,
			onSuccess:      m.OnSuccess(emailInfo),
		}
		m.logger.Info("Generated move reviewed email to service member",
			zap.String("service member email address", *emailInfo.Email))
		emails = append(emails, smEmail)
	}
	return emails, nil
}

func (m MoveReviewed) renderTemplates(data moveReviewedEmailData) (string, string, error) {
	htmlBody, err := m.RenderHTML(data)
	if err != nil {
		return "", "", fmt.Errorf("error rendering html template using %#v", data)
	}
	textBody, err := m.RenderText(data)
	if err != nil {
		return "", "", fmt.Errorf("error rendering text template using %#v", data)
	}
	return htmlBody, textBody, nil
}

// OnSuccess callback passed to be invoked by NewNotificationSender when an email successfully sent
// saves the svs the email info along with the SES mail id to the notifications table
func (m MoveReviewed) OnSuccess(emailInfo EmailInfo) func(string) error {
	return func(msgID string) error {
		n := models.Notification{
			ServiceMemberID:  emailInfo.ServiceMemberID,
			SESMessageID:     msgID,
			NotificationType: models.MoveReviewedEmail,
		}
		err := m.db.Create(&n)
		if err != nil {
			dataString := fmt.Sprintf("%#v", n)
			m.logger.Error("adding notification to notifications table", zap.String("notification", dataString))
			return err
		}
		return nil
	}
}

type moveReviewedEmailData struct {
	Link                   string
	OriginDutyStation      string
	DestinationDutyStation string
}

// RenderHTML renders the html for the email
func (m MoveReviewed) RenderHTML(data moveReviewedEmailData) (string, error) {
	var htmlBuffer bytes.Buffer
	if err := m.htmlTemplate.Execute(&htmlBuffer, data); err != nil {
		m.logger.Error("cant render html template ", zap.Error(err))
	}
	return htmlBuffer.String(), nil
}

// RenderText renders the text for the email
func (m MoveReviewed) RenderText(data moveReviewedEmailData) (string, error) {
	var textBuffer bytes.Buffer
	if err := m.textTemplate.Execute(&textBuffer, data); err != nil {
		m.logger.Error("cant render text template ", zap.Error(err))
		return "", err
	}
	return textBuffer.String(), nil
}
