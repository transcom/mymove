package notifications

import (
	"bytes"
	"context"
	html "html/template"
	text "text/template"
	"time"

	"github.com/gobuffalo/pop"
	"go.uber.org/zap"
)

const link = "https://www.surveymonkey.com/r/MilMovePt3-08191"

var htmlTemplate = html.Must(html.New("email_survey_html").Parse(`<em>Good news:</em> Your move from {{.OriginDutyStation}} to {{.DestinationDutyStation}} has been processed for payment.

Can we ask a quick favor? <a href="{{.Link}}"> Tell us about your experience</a> with requesting and receiving payment.

We’ll use your feedback to make MilMove better for your fellow service members.

Thank you for your thoughts, and <em>congratulations on your move.</em>`))

var textTemplate = text.Must(text.New("email_survey_text").Parse(`Good news: Your move from {{.OriginDutyStation}} to {{.DestinationDutyStation}} has been processed for payment.

Can we ask a quick favor? Tell us about your experience with requesting and receiving payment at {{.Link}}.

We’ll use your feedback to make MilMove better for your fellow service members.

Thank you for your thoughts, and congratulations on your move.`))

// MoveReviewed has notification content for completed/reviewed moves
type MoveReviewed struct {
	db     *pop.Connection
	logger Logger
	date   time.Time
}

// NewMoveReviewed returns a new move submitted notification
func NewMoveReviewed(db *pop.Connection, logger Logger, date time.Time) *MoveReviewed {

	return &MoveReviewed{
		db:     db,
		logger: logger,
		date:   date,
	}
}

type EmailInfos []EmailInfo

type EmailInfo struct {
	Email              *string `db:"personal_email"`
	DutyStationName    string  `db:"name"`
	NewDutyStationName string  `db:"name"`
}

func (m MoveReviewed) GetEmailInfo(date time.Time) (*EmailInfos, error) {
	dateString := date.Format("2006-01-02")
	query := `SELECT sm.personal_email, dsn.name, dso.name
	FROM personally_procured_moves
	         JOIN moves m ON personally_procured_moves.move_id = m.id
	         JOIN orders o ON m.orders_id = o.id
	         JOIN service_members sm ON o.service_member_id = sm.id
	         JOIN duty_stations dso ON sm.duty_station_id = dso.id
	         JOIN duty_stations dsn ON o.new_duty_station_id = dsn.id
	WHERE CAST(reviewed_date as date) = $1;`

	emailInfo := &EmailInfos{}
	err := m.db.RawQuery(query, dateString).All(emailInfo)
	return emailInfo, err
}

// Notifications expects emails to be implemented so we do
func (m MoveReviewed) emails(ctx context.Context) ([]emailContent, error) {
	emailInfos, err := m.GetEmailInfo(m.date)
	if emailInfos == nil || err == nil {
		m.logger.Info("no emails to be sent for", zap.String("date", m.date.String()))
		return nil, nil
	}
	return m.FormatEmails(emailInfos)
}

//FormatEmails formats email data using both html and text template
func (m MoveReviewed) FormatEmails(emailInfos *EmailInfos) ([]emailContent, error) {
	var emails []emailContent
	for _, emailInfo := range *emailInfos {
		var email string
		if emailInfo.Email == nil {
			m.logger.Error("no email found for service member")
			continue
		}
		email = *emailInfo.Email
		data := emailData{
			Link:                   link,
			OriginDutyStation:      emailInfo.DutyStationName,
			DestinationDutyStation: emailInfo.NewDutyStationName,
			Email:                  email,
		}
		smEmail := emailContent{
			recipientEmail: email,
			subject:        "[MilMove] Let us know how we did",
			htmlBody:       m.RenderHTML(data),
			textBody:       m.RenderText(data),
		}
		m.logger.Info("Generated move reviewed email to service member",
			zap.String("service member email address", email))
		// TODO: Send email to trusted contacts when that's supported
		emails = append(emails, smEmail)
	}
	return emails, nil
}

type emailData struct {
	Link                   string
	OriginDutyStation      string
	DestinationDutyStation string
	Email                  string
}

// RenderHTML renders the html for the email
func (m MoveReviewed) RenderHTML(data emailData) string {
	var htmlBuffer bytes.Buffer
	if err := htmlTemplate.Execute(&htmlBuffer, data); err != nil {
		m.logger.Error("cant render html template for: ",
			zap.String("service member email address", data.Email))
	}
	return htmlBuffer.String()
}

// RenderText renders the text for the email
func (m MoveReviewed) RenderText(data emailData) string {
	var textBuffer bytes.Buffer
	if err := textTemplate.Execute(&textBuffer, data); err != nil {
		m.logger.Error("cant render text template for: ",
			zap.String("service member email address", data.Email))
	}
	return textBuffer.String()
}
