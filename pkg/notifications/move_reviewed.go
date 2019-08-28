package notifications

import (
	"bytes"
	"context"
	"fmt"
	html "html/template"
	text "text/template"
	"time"

	"github.com/transcom/mymove/pkg/assets"

	"github.com/gobuffalo/pop"
	"go.uber.org/zap"
)

const surveyLink = "https://www.surveymonkey.com/r/MilMovePt3-08191"

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

	htmlTemplate, err := initHMTLTemplate(logger)
	if err != nil {
		return &MoveReviewed{}, err
	}

	textTemplate, err := initTextTemplate(logger)
	if err != nil {
		return &MoveReviewed{}, err
	}

	return &MoveReviewed{
		db:           db,
		logger:       logger,
		date:         date,
		htmlTemplate: htmlTemplate,
		textTemplate: textTemplate,
	}, nil
}

func initTextTemplate(logger Logger) (*text.Template, error) {
	tt, err := assets.Asset("pkg/notifications/templates/move_reviewed_template.txt")
	if err != nil {
		logger.Error("text template pathing error")
		return nil, err
	}
	templateString := string(tt)
	textTemplate, err := text.New("text_template").Parse(templateString)
	if err != nil {
		logger.Error("unable to parse text template", zap.String("template:", templateString))
		return nil, err
	}
	return textTemplate, nil
}

func initHMTLTemplate(logger Logger) (*html.Template, error) {
	ht, err := assets.Asset("pkg/notifications/templates/move_reviewed_template.html")
	if err != nil {
		logger.Error("html template pathing error")
		return nil, err
	}
	htmlTemplateString := string(ht)
	htmlTemplate, err := html.New("html_template").Parse(htmlTemplateString)
	if err != nil {
		logger.Error("unable to parse html template", zap.String("template:", htmlTemplateString))
		return nil, err
	}
	return htmlTemplate, err
}

type EmailInfos []EmailInfo

type EmailInfo struct {
	Email              *string `db:"personal_email"`
	DutyStationName    string  `db:"duty_station_name"`
	NewDutyStationName string  `db:"new_duty_station_name"`
}

func (m MoveReviewed) GetEmailInfo(date time.Time) (EmailInfos, error) {
	dateString := date.Format("2006-01-02")
	query := `SELECT sm.personal_email, dsn.name as new_duty_station_name, dso.name as duty_station_name
	FROM personally_procured_moves p
	         JOIN moves m ON p.move_id = m.id
	         JOIN orders o ON m.orders_id = o.id
	         JOIN service_members sm ON o.service_member_id = sm.id
	         JOIN duty_stations dso ON sm.duty_station_id = dso.id
	         JOIN duty_stations dsn ON o.new_duty_station_id = dsn.id
	WHERE CAST(reviewed_date as date) = $1;`

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

// TODO figure out best way to fix linter complaint that []emailContent private but method is public
// formatEmails formats email data using both html and text template
func (m MoveReviewed) formatEmails(emailInfos EmailInfos) ([]emailContent, error) {
	var emails []emailContent
	for _, emailInfo := range emailInfos {
		var email string
		if emailInfo.Email == nil {
			m.logger.Info("no email found for service member")
			continue
		}
		email = *emailInfo.Email
		data := moveReviewedEmailData{
			Link:                   surveyLink,
			OriginDutyStation:      emailInfo.DutyStationName,
			DestinationDutyStation: emailInfo.NewDutyStationName,
			Email:                  email,
		}
		htmlBody, err := m.RenderHTML(data)
		if err != nil {
			dataString := fmt.Sprintf("%#v", data)
			m.logger.Error("error rendering html template using", zap.String("data", dataString))
			continue
		}
		textBody, err := m.RenderText(data)
		if err != nil {
			dataString := fmt.Sprintf("%#v", data)
			m.logger.Error("error rendering text template using", zap.String("data", dataString))
			continue
		}
		smEmail := emailContent{
			recipientEmail: email,
			subject:        "[MilMove] Let us know how we did",
			htmlBody:       htmlBody,
			textBody:       textBody,
		}
		m.logger.Info("Generated move reviewed email to service member",
			zap.String("service member email address", email))
		// TODO: Send email to trusted contacts when that's supported
		emails = append(emails, smEmail)
	}
	return emails, nil
}

type moveReviewedEmailData struct {
	Link                   string
	OriginDutyStation      string
	DestinationDutyStation string
	Email                  string
}

// RenderHTML renders the html for the email
func (m MoveReviewed) RenderHTML(data moveReviewedEmailData) (string, error) {
	var htmlBuffer bytes.Buffer
	if err := m.htmlTemplate.Execute(&htmlBuffer, data); err != nil {
		m.logger.Error("cant render html template for: ",
			zap.String("service member email address", data.Email))
	}
	return htmlBuffer.String(), nil
}

// RenderText renders the text for the email
func (m MoveReviewed) RenderText(data moveReviewedEmailData) (string, error) {
	var textBuffer bytes.Buffer
	if err := m.textTemplate.Execute(&textBuffer, data); err != nil {
		m.logger.Error("cant render text template for: ",
			zap.String("service member email address", data.Email))
		return "", err
	}
	return textBuffer.String(), nil
}
