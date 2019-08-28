package notifications

import (
	"bytes"
	"context"
	html "html/template"
	text "text/template"
	"time"

	"github.com/gobuffalo/pop"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/models"
)

const surveyLink = "https://www.surveymonkey.com/r/MilMovePt3-08191"

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
	PersonallyProcuredMove models.PersonallyProcuredMove `db:"ppm"`
	Email                  *string                       `db:"personal_email"`
	DutyStationName        string                        `db:"duty_station_name"`
	NewDutyStationName     string                        `db:"new_duty_station_name"`
}

func (m MoveReviewed) GetEmailInfo(date time.Time) (EmailInfos, error) {
	dateString := date.Format("2006-01-02")
	query := `SELECT p.* as ppm, sm.personal_email, dsn.name as new_duty_station_name, dso.name as duty_station_name
	FROM personally_procured_moves p
	         JOIN moves m ON p.move_id = m.id
	         JOIN orders o ON m.orders_id = o.id
	         JOIN service_members sm ON o.service_member_id = sm.id
	         JOIN duty_stations dso ON sm.duty_station_id = dso.id
	         JOIN duty_stations dsn ON o.new_duty_station_id = dsn.id
	WHERE CAST(p.reviewed_date as date) = $1
		AND p.survey_email_sent = false;`

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
		return []emailContent{}, nil
	}
	if emailInfos == nil {
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
		emailInfo.PersonallyProcuredMove.SurveyEmailSent = true
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
func (m MoveReviewed) RenderHTML(data moveReviewedEmailData) string {
	var htmlBuffer bytes.Buffer
	if err := htmlTemplate.Execute(&htmlBuffer, data); err != nil {
		m.logger.Error("cant render html template for: ",
			zap.String("service member email address", data.Email))
	}
	return htmlBuffer.String()
}

// RenderText renders the text for the email
func (m MoveReviewed) RenderText(data moveReviewedEmailData) string {
	var textBuffer bytes.Buffer
	if err := textTemplate.Execute(&textBuffer, data); err != nil {
		m.logger.Error("cant render text template for: ",
			zap.String("service member email address", data.Email))
	}
	return textBuffer.String()
}
