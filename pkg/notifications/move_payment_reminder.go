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
	// "github.com/transcom/mymove/pkg/handlers"

	"github.com/transcom/mymove/pkg/assets"

	"github.com/gobuffalo/pop"
	"go.uber.org/zap"
)

var (
	movePaymentReminderRawText  = string(assets.MustAsset("pkg/notifications/templates/move_payment_reminder_template.txt"))
	paymentReminderTextTemplate = text.Must(text.New("text_template").Parse(movePaymentReminderRawText))
	movePaymentReminderRawHTML  = string(assets.MustAsset("pkg/notifications/templates/move_payment_reminder_template.html"))
	paymentReminderHTMLTemplate = html.Must(html.New("text_template").Parse(movePaymentReminderRawHTML))
)

// PaymentReminder has notification content for completed/reviewed moves
type PaymentReminder struct {
	db           *pop.Connection
	logger       Logger
	date         time.Time
	htmlTemplate *html.Template
	textTemplate *text.Template
}

// NewPaymentReminder returns a new move submitted notification
func NewPaymentReminder(db *pop.Connection, logger Logger, date time.Time) (*PaymentReminder, error) {

	return &PaymentReminder{
		db:           db,
		logger:       logger,
		date:         date,
		htmlTemplate: paymentReminderHTMLTemplate,
		textTemplate: paymentReminderTextTemplate,
	}, nil
}

type PaymentReminderEmailInfos []PaymentReminderEmailInfo

type PaymentReminderEmailInfo struct {
	ServiceMemberID    uuid.UUID `db:"id"`
	Email              *string   `db:"personal_email"`
	DutyStationName    string    `db:"duty_station_name"`
	NewDutyStationName string    `db:"new_duty_station_name"`
}

func (m PaymentReminder) GetEmailInfo(date time.Time) (PaymentReminderEmailInfos, error) {
	// 	dateString := date.Format("2006-01-02")
	query := `SELECT 'e7edaddf-f4f9-401f-940b-b6c3be84195d' as id, 'lindsay+test1@truss.works' as personal_email, 'abc' as duty_station_name, '123' as new_duty_station_name`
	// 	query := `SELECT sm.id, sm.personal_email, dsn.name AS new_duty_station_name, dso.name AS duty_station_name
	// FROM personally_procured_moves p
	//          JOIN moves m ON p.move_id = m.id
	//          JOIN orders o ON m.orders_id = o.id
	//          JOIN service_members sm ON o.service_member_id = sm.id
	//          JOIN duty_stations dso ON sm.duty_station_id = dso.id
	//          JOIN duty_stations dsn ON o.new_duty_station_id = dsn.id
	//          LEFT JOIN notifications n ON sm.id = n.service_member_id
	// WHERE CAST(reviewed_date AS date) = $1
	// --  send email if haven't sent them a MOVE_REVIEWED_EMAIL yet OR we haven't sent them any emails at all
	//     AND (notification_type != 'MOVE_REVIEWED_EMAIL' OR n.service_member_id IS NULL);`

	paymentReminderEmailInfos := PaymentReminderEmailInfos{}
	err := m.db.RawQuery(query).All(&paymentReminderEmailInfos)

	return paymentReminderEmailInfos, err
	// return PaymentReminderEmailInfos, nil
}

// NotificationSendingContext expects a `notification` with an `emails` method,
// so we implement `email` to satisfy that interface
func (m PaymentReminder) emails(ctx context.Context) ([]emailContent, error) {
	paymentReminderEmailInfos, err := m.GetEmailInfo(m.date)
	if err != nil {
		m.logger.Error("error retrieving email info", zap.String("date", m.date.String()))
		return []emailContent{}, err
	}
	if len(PaymentReminderEmailInfos) == 0 {
		m.logger.Info("no emails to be sent", zap.String("date", m.date.String()))
		return []emailContent{}, nil
	}
	return m.formatEmails(paymentReminderEmailInfos)
}

// formatEmails formats email data using both html and text template
func (m PaymentReminder) formatEmails(PaymentReminderEmailInfos PaymentReminderEmailInfos) ([]emailContent, error) {
	var emails []emailContent
	for _, PaymentReminderemailInfo := range PaymentReminderEmailInfos {
		htmlBody, textBody, err := m.renderTemplates(paymentReminderEmailData{
			OriginDutyStation:      PaymentReminderemailInfo.DutyStationName,
			DestinationDutyStation: PaymentReminderemailInfo.NewDutyStationName,
		})
		if err != nil {
			m.logger.Error("error rendering template", zap.Error(err))
			continue
		}
		if PaymentReminderemailInfo.Email == nil {
			m.logger.Info("no email found for service member",
				zap.String("service member uuid", PaymentReminderemailInfo.ServiceMemberID.String()))
			continue
		}
		smEmail := emailContent{
			recipientEmail: *PaymentReminderemailInfo.Email,
			subject:        "[MilMove] Let us know how we did",
			htmlBody:       htmlBody,
			textBody:       textBody,
			onSuccess:      m.OnSuccess(PaymentReminderemailInfo),
		}
		m.logger.Info("generated move reviewed email to service member",
			zap.String("service member uuid", PaymentReminderemailInfo.ServiceMemberID.String()))
		emails = append(emails, smEmail)
	}
	return emails, nil
}

func (m PaymentReminder) renderTemplates(data paymentReminderEmailData) (string, string, error) {
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
func (m PaymentReminder) OnSuccess(PaymentReminderemailInfo PaymentReminderEmailInfo) func(string) error {
	return func(msgID string) error {
		n := models.Notification{
			ServiceMemberID:  PaymentReminderemailInfo.ServiceMemberID,
			SESMessageID:     msgID,
			NotificationType: models.MovePaymentReminderEmail,
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

type paymentReminderEmailData struct {
	Link                   string
	OriginDutyStation      string
	DestinationDutyStation string
}

// RenderHTML renders the html for the email
func (m PaymentReminder) RenderHTML(data paymentReminderEmailData) (string, error) {
	var htmlBuffer bytes.Buffer
	if err := m.htmlTemplate.Execute(&htmlBuffer, data); err != nil {
		m.logger.Error("cant render html template ", zap.Error(err))
	}
	return htmlBuffer.String(), nil
}

// RenderText renders the text for the email
func (m PaymentReminder) RenderText(data paymentReminderEmailData) (string, error) {
	var textBuffer bytes.Buffer
	if err := m.textTemplate.Execute(&textBuffer, data); err != nil {
		m.logger.Error("cant render text template ", zap.Error(err))
		return "", err
	}
	return textBuffer.String(), nil
}
