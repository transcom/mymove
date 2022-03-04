package notifications

import (
	"bytes"
	"fmt"
	html "html/template"
	text "text/template"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/assets"
	"github.com/transcom/mymove/pkg/unit"

	"go.uber.org/zap"
)

var (
	movePaymentReminderRawText  = string(assets.MustAsset("pkg/notifications/templates/move_payment_reminder_template.txt"))
	paymentReminderTextTemplate = text.Must(text.New("text_template").Parse(movePaymentReminderRawText))
	movePaymentReminderRawHTML  = string(assets.MustAsset("pkg/notifications/templates/move_payment_reminder_template.html"))
	paymentReminderHTMLTemplate = html.Must(html.New("text_template").Parse(movePaymentReminderRawHTML))
)

// PaymentReminder has notification content for approved moves
type PaymentReminder struct {
	emailAfter    string
	noEmailBefore string
	htmlTemplate  *html.Template
	textTemplate  *text.Template
}

// NewPaymentReminder returns a new payment reminder notification
func NewPaymentReminder() (*PaymentReminder, error) {

	return &PaymentReminder{
		emailAfter:    "10 DAYS",
		noEmailBefore: "2019-06-01",
		htmlTemplate:  paymentReminderHTMLTemplate,
		textTemplate:  paymentReminderTextTemplate,
	}, nil
}

// PaymentReminderEmailInfos is a slice of PaymentReminderEmailInfo
type PaymentReminderEmailInfos []PaymentReminderEmailInfo

// PaymentReminderEmailInfo contains payment reminder data for rendering a template
type PaymentReminderEmailInfo struct {
	ServiceMemberID      uuid.UUID   `db:"id"`
	Email                *string     `db:"personal_email"`
	NewDutyStationName   string      `db:"new_duty_station_name"`
	WeightEstimate       *unit.Pound `db:"weight_estimate"`
	IncentiveEstimateMin *unit.Cents `db:"incentive_estimate_min"`
	IncentiveEstimateMax *unit.Cents `db:"incentive_estimate_max"`
	IncentiveTxt         string
	TOName               *string `db:"transportation_office_name"`
	TOPhone              *string `db:"transportation_office_phone"`
	MoveDate             string  `db:"move_date"`
	Locator              string  `db:"locator"`
}

// GetEmailInfo fetches payment email information
func (m PaymentReminder) GetEmailInfo(appCtx appcontext.AppContext) (PaymentReminderEmailInfos, error) {
	query := `SELECT sm.id as id, sm.personal_email AS personal_email,
	COALESCE(ppm.weight_estimate, 0) AS weight_estimate,
	COALESCE(ppm.incentive_estimate_min, 0) AS incentive_estimate_min,
	COALESCE(ppm.incentive_estimate_max, 0) AS incentive_estimate_max,
	ppm.original_move_date as move_date,
	dln.name AS new_duty_station_name,
	tos.name AS transportation_office_name,
	opl.number AS transportation_office_phone,
	m.locator
FROM personally_procured_moves ppm
	JOIN moves m ON ppm.move_id = m.id
	JOIN orders o ON m.orders_id = o.id
	JOIN service_members sm ON o.service_member_id = sm.id
	JOIN duty_locations dln ON o.new_duty_location_id = dln.id
	JOIN transportation_offices tos ON tos.id = dln.transportation_office_id
	LEFT JOIN office_phone_lines opl on opl.transportation_office_id = tos.id and opl.id =
	(
		SELECT opl2.id FROM office_phone_lines opl2
		WHERE opl2.is_dsn_number IS false
		AND tos.id = opl2.transportation_office_id
		LIMIT 1
	)
	LEFT JOIN notifications n ON sm.id = n.service_member_id
	WHERE ppm.original_move_date <= now() - ($1)::INTERVAL
	AND ppm.original_move_date >= $2
	AND ppm.status = 'APPROVED'
	AND (notification_type != 'MOVE_PAYMENT_REMINDER_EMAIL' OR n.service_member_id IS NULL)
	AND m.status = 'APPROVED'
	AND m.show IS true
	ORDER BY m.locator;`

	paymentReminderEmailInfos := PaymentReminderEmailInfos{}
	err := appCtx.DB().RawQuery(query, m.emailAfter, m.noEmailBefore).All(&paymentReminderEmailInfos)

	return paymentReminderEmailInfos, err
}

// NotificationSendingContext expects a `notification` with an `emails` method,
// so we implement `email` to satisfy that interface
func (m PaymentReminder) emails(appCtx appcontext.AppContext) ([]emailContent, error) {
	paymentReminderEmailInfos, err := m.GetEmailInfo(appCtx)
	if err != nil {
		appCtx.Logger().Error("error retrieving email info")
		return []emailContent{}, err
	}
	if len(paymentReminderEmailInfos) == 0 {
		appCtx.Logger().Info("no emails to be sent")
		return []emailContent{}, nil
	}
	return m.formatEmails(appCtx, paymentReminderEmailInfos)
}

// formatEmails formats email data using both html and text template
func (m PaymentReminder) formatEmails(appCtx appcontext.AppContext, PaymentReminderEmailInfos PaymentReminderEmailInfos) ([]emailContent, error) {
	var emails []emailContent
	for _, PaymentReminderEmailInfo := range PaymentReminderEmailInfos {
		incentiveTxt := ""
		if PaymentReminderEmailInfo.WeightEstimate.Int() > 0 && PaymentReminderEmailInfo.IncentiveEstimateMin.Int() > 0 && PaymentReminderEmailInfo.IncentiveEstimateMax.Int() > 0 {
			incentiveTxt = fmt.Sprintf("You expected to move about %d lbs, which gives you an estimated incentive of %s-%s.", PaymentReminderEmailInfo.WeightEstimate.Int(), PaymentReminderEmailInfo.IncentiveEstimateMin.ToDollarString(), PaymentReminderEmailInfo.IncentiveEstimateMax.ToDollarString())
		}
		var toPhone *string
		if PaymentReminderEmailInfo.TOPhone != nil {
			toPhone = PaymentReminderEmailInfo.TOPhone
		}

		var toName *string
		if PaymentReminderEmailInfo.TOPhone != nil {
			toName = PaymentReminderEmailInfo.TOName
		}

		htmlBody, textBody, err := m.renderTemplates(appCtx, PaymentReminderEmailData{
			DestinationDutyStation: PaymentReminderEmailInfo.NewDutyStationName,
			WeightEstimate:         fmt.Sprintf("%d", PaymentReminderEmailInfo.WeightEstimate.Int()),
			IncentiveEstimateMin:   PaymentReminderEmailInfo.IncentiveEstimateMin.ToDollarString(),
			IncentiveEstimateMax:   PaymentReminderEmailInfo.IncentiveEstimateMax.ToDollarString(),
			IncentiveTxt:           incentiveTxt,
			TOName:                 toName,
			TOPhone:                toPhone,
			Locator:                PaymentReminderEmailInfo.Locator,
		})
		if err != nil {
			appCtx.Logger().Error("error rendering template", zap.Error(err))
			continue
		}
		if PaymentReminderEmailInfo.Email == nil {
			appCtx.Logger().Info("no email found for service member",
				zap.String("service member uuid", PaymentReminderEmailInfo.ServiceMemberID.String()))
			continue
		}
		smEmail := emailContent{
			recipientEmail: *PaymentReminderEmailInfo.Email,
			subject:        fmt.Sprintf("[MilMove] Reminder: request payment for your move to %s (move %s)", PaymentReminderEmailInfo.NewDutyStationName, PaymentReminderEmailInfo.Locator),
			htmlBody:       htmlBody,
			textBody:       textBody,
			onSuccess:      m.OnSuccess(appCtx, PaymentReminderEmailInfo),
		}
		appCtx.Logger().Info("generated payment reminder email to service member",
			zap.String("service member uuid", PaymentReminderEmailInfo.ServiceMemberID.String()),
			zap.String("moveLocator", PaymentReminderEmailInfo.Locator),
		)
		emails = append(emails, smEmail)
	}
	return emails, nil
}

func (m PaymentReminder) renderTemplates(appCtx appcontext.AppContext, data PaymentReminderEmailData) (string, string, error) {
	htmlBody, err := m.RenderHTML(appCtx, data)
	if err != nil {
		return "", "", fmt.Errorf("error rendering html template using %#v", data)
	}
	textBody, err := m.RenderText(appCtx, data)
	if err != nil {
		return "", "", fmt.Errorf("error rendering text template using %#v", data)
	}
	return htmlBody, textBody, nil
}

// OnSuccess callback passed to be invoked by NewNotificationSender when an email successfully sent
// saves the svs the email info along with the SES mail id to the notifications table
func (m PaymentReminder) OnSuccess(appCtx appcontext.AppContext, PaymentReminderEmailInfo PaymentReminderEmailInfo) func(string) error {
	return func(msgID string) error {
		n := models.Notification{
			ServiceMemberID:  PaymentReminderEmailInfo.ServiceMemberID,
			SESMessageID:     msgID,
			NotificationType: models.MovePaymentReminderEmail,
		}
		err := appCtx.DB().Create(&n)
		if err != nil {
			dataString := fmt.Sprintf("%#v", n)
			appCtx.Logger().Error("adding notification to notifications table", zap.String("notification", dataString))
			return err
		}
		return nil
	}
}

// PaymentReminderEmailData is used to render an email template
type PaymentReminderEmailData struct {
	DestinationDutyStation string
	WeightEstimate         string
	IncentiveEstimateMin   string
	IncentiveEstimateMax   string
	IncentiveTxt           string
	TOName                 *string
	TOPhone                *string
	Locator                string
}

// RenderHTML renders the html for the email
func (m PaymentReminder) RenderHTML(appCtx appcontext.AppContext, data PaymentReminderEmailData) (string, error) {
	var htmlBuffer bytes.Buffer
	if err := m.htmlTemplate.Execute(&htmlBuffer, data); err != nil {
		appCtx.Logger().Error("cant render html template ", zap.Error(err))
	}
	return htmlBuffer.String(), nil
}

// RenderText renders the text for the email
func (m PaymentReminder) RenderText(appCtx appcontext.AppContext, data PaymentReminderEmailData) (string, error) {
	var textBuffer bytes.Buffer
	if err := m.textTemplate.Execute(&textBuffer, data); err != nil {
		appCtx.Logger().Error("cant render text template ", zap.Error(err))
		return "", err
	}
	return textBuffer.String(), nil
}
