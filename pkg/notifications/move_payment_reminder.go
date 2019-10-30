package notifications

import (
	"bytes"
	"context"
	"fmt"
	html "html/template"
	text "text/template"

	"github.com/transcom/mymove/pkg/models"

	"github.com/gofrs/uuid"
	// "github.com/transcom/mymove/pkg/handlers"

	"github.com/transcom/mymove/pkg/assets"
	"github.com/transcom/mymove/pkg/unit"

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
	htmlTemplate *html.Template
	textTemplate *text.Template
}

// NewPaymentReminder returns a new move submitted notification
func NewPaymentReminder(db *pop.Connection, logger Logger) (*PaymentReminder, error) {

	return &PaymentReminder{
		db:           db,
		logger:       logger,
		htmlTemplate: paymentReminderHTMLTemplate,
		textTemplate: paymentReminderTextTemplate,
	}, nil
}

type PaymentReminderEmailInfos []PaymentReminderEmailInfo

type PaymentReminderEmailInfo struct {
	ServiceMemberID      uuid.UUID   `db:"id"`
	Email                *string     `db:"personal_email"`
	DutyStationName      string      `db:"duty_station_name"`
	NewDutyStationName   string      `db:"new_duty_station_name"`
	WeightEstimate       *unit.Pound `db:"weight_estimate"`
	IncentiveEstimateMin *unit.Cents `db:"incentive_estimate_min"`
	IncentiveEstimateMax *unit.Cents `db:"incentive_estimate_max"`
	TOName               string      `db:"transportation_office_name"`
	TOPhone              string      `db:"transportation_office_phone"`
	ReviewedDate         string      `db:"reviewed_date"`
}

func (m PaymentReminder) GetEmailInfo() (PaymentReminderEmailInfos, error) {
	// 	query := `SELECT 'e7edaddf-f4f9-401f-940b-b6c3be84195d' as id, 'lindsay+test1@truss.works' as personal_email, 'Yuma AFB' as duty_station_name, 'Fort Gordon' as new_duty_station_name,
	// 8000 as weight_estimate, 500 as incentive_estimate_min, 1000 as incentive_estimate_max, 'blah PPO' as transportation_office_name, '555-555-1212' as transportation_office_phone`
	query := `SELECT sm.id as id, sm.personal_email as personal_email,
                     ppm.weight_estimate, ppm.incentive_estimate_min, ppm.incentive_estimate_max,
                     ppm.reviewed_date as reviewed_date,
                     dsn.name AS new_duty_station_name,
                     toff.name AS transportation_office_name,
                     opl.number AS transportation_office_phone
	FROM personally_procured_moves ppm
	         JOIN moves m ON ppm.move_id = m.id
	         JOIN orders o ON m.orders_id = o.id
	         JOIN service_members sm ON o.service_member_id = sm.id
	         JOIN duty_stations dsn ON o.new_duty_station_id = dsn.id
	         JOIN transportation_offices toff ON toff.id = dsn.transportation_office_id
	         JOIN office_phone_lines opl on toff.id = opl.transportation_office_id
	         LEFT JOIN notifications n ON sm.id = n.service_member_id
	where ppm.reviewed_date <= now() - INTERVAL '10 DAYS'
	      AND ppm.reviewed_date >= '2019-10-01'
          AND (notification_type != 'MOVE_PAYMENT_REMINDER_EMAIL' OR n.service_member_id IS NULL);
`

	paymentReminderEmailInfos := PaymentReminderEmailInfos{}
	err := m.db.RawQuery(query).All(&paymentReminderEmailInfos)

	return paymentReminderEmailInfos, err
}

// NotificationSendingContext expects a `notification` with an `emails` method,
// so we implement `email` to satisfy that interface
func (m PaymentReminder) emails(ctx context.Context) ([]emailContent, error) {
	paymentReminderEmailInfos, err := m.GetEmailInfo()
	if err != nil {
		m.logger.Error("error retrieving email info")
		return []emailContent{}, err
	}
	if len(paymentReminderEmailInfos) == 0 {
		m.logger.Info("no emails to be sent")
		return []emailContent{}, nil
	}
	return m.formatEmails(paymentReminderEmailInfos)
}

// formatEmails formats email data using both html and text template
func (m PaymentReminder) formatEmails(PaymentReminderEmailInfos PaymentReminderEmailInfos) ([]emailContent, error) {
	var emails []emailContent
	for _, PaymentReminderemailInfo := range PaymentReminderEmailInfos {
		htmlBody, textBody, err := m.renderTemplates(PaymentReminderEmailData{
			DestinationDutyStation: PaymentReminderemailInfo.NewDutyStationName,
			WeightEstimate:         fmt.Sprintf("%d", PaymentReminderemailInfo.WeightEstimate),
			IncentiveEstimateMin:   PaymentReminderemailInfo.IncentiveEstimateMin.ToDollarString(),
			IncentiveEstimateMax:   PaymentReminderemailInfo.IncentiveEstimateMax.ToDollarString(),
			TOName:                 PaymentReminderemailInfo.TOName,
			TOPhone:                PaymentReminderemailInfo.TOPhone,
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

func (m PaymentReminder) renderTemplates(data PaymentReminderEmailData) (string, string, error) {
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

type PaymentReminderEmailData struct {
	DestinationDutyStation string
	WeightEstimate         string
	IncentiveEstimateMin   string
	IncentiveEstimateMax   string
	TOName                 string
	TOPhone                string
}

// RenderHTML renders the html for the email
func (m PaymentReminder) RenderHTML(data PaymentReminderEmailData) (string, error) {
	var htmlBuffer bytes.Buffer
	if err := m.htmlTemplate.Execute(&htmlBuffer, data); err != nil {
		m.logger.Error("cant render html template ", zap.Error(err))
	}
	return htmlBuffer.String(), nil
}

// RenderText renders the text for the email
func (m PaymentReminder) RenderText(data PaymentReminderEmailData) (string, error) {
	var textBuffer bytes.Buffer
	if err := m.textTemplate.Execute(&textBuffer, data); err != nil {
		m.logger.Error("cant render text template ", zap.Error(err))
		return "", err
	}
	return textBuffer.String(), nil
}
