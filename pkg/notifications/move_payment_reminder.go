package notifications

import (
	"bytes"
	"fmt"
	html "html/template"
	text "text/template"

	"github.com/gofrs/uuid"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/assets"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/unit"
)

var (
	movePaymentReminderRawText  = string(assets.MustAsset("notifications/templates/move_payment_reminder_template.txt"))
	paymentReminderTextTemplate = text.Must(text.New("text_template").Parse(movePaymentReminderRawText))
	movePaymentReminderRawHTML  = string(assets.MustAsset("notifications/templates/move_payment_reminder_template.html"))
	paymentReminderHTMLTemplate = html.Must(html.New("text_template").Parse(movePaymentReminderRawHTML))
)

// PaymentReminder has notification content for approved moves
type PaymentReminder struct {
	emailAfter    string
	noEmailBefore string
	htmlTemplate  *html.Template
	textTemplate  *text.Template
}

// NewPaymentReminder returns a new payment reminder notification 14 days after actual move in date
func NewPaymentReminder() (*PaymentReminder, error) {

	return &PaymentReminder{
		emailAfter:    "14 DAYS",
		noEmailBefore: "2019-06-01",
		htmlTemplate:  paymentReminderHTMLTemplate,
		textTemplate:  paymentReminderTextTemplate,
	}, nil
}

// PaymentReminderEmailInfos is a slice of PaymentReminderEmailInfo
type PaymentReminderEmailInfos []PaymentReminderEmailInfo

// PaymentReminderEmailInfo contains payment reminder data for rendering a template
type PaymentReminderEmailInfo struct {
	ServiceMemberID        uuid.UUID                   `db:"id"`
	Email                  *string                     `db:"personal_email"`
	NewDutyLocationName    string                      `db:"new_duty_location_name"`
	OriginDutyLocationName string                      `db:"origin_duty_location_name"`
	MoveDate               string                      `db:"move_date"`
	Locator                string                      `db:"locator"`
	WeightEstimate         *unit.Pound                 `db:"weight_estimate"`
	IncentiveEstimate      *unit.Cents                 `db:"incentive_estimate"`
	DestinationStreet1     *string                     `db:"destination_street_address_1"`
	DestinationStreet2     *string                     `db:"destination_street_address_2"`
	DestinationStreet3     *string                     `db:"destination_street_address_3"`
	DestinationCity        *string                     `db:"destination_city"`
	DestinationState       *string                     `db:"destination_state"`
	DestinationPostalCode  *string                     `db:"destination_postal_code"`
	OrdersType             internalmessages.OrdersType `db:"orders_type"`
}

// GetEmailInfo fetches payment email information
// left joins on duty locations to allow for those fields to be null
func (m PaymentReminder) GetEmailInfo(appCtx appcontext.AppContext) (PaymentReminderEmailInfos, error) {
	query := `SELECT DISTINCT sm.id as id, sm.personal_email AS personal_email,
	COALESCE(ps.estimated_weight, 0) AS weight_estimate,
	COALESCE(ps.estimated_incentive, 0) AS incentive_estimate,
	ps.expected_departure_date  as move_date,
	dln.name AS new_duty_location_name,
	dln2.name AS origin_duty_location_name,
	m.locator,
	da.street_address_1 AS destination_street_address_1,
	da.street_address_2 AS destination_street_address_2,
	da.street_address_3 AS destination_street_address_3,
	da.city AS destination_city,
	da.state AS destination_state,
	da.postal_code AS destination_postal_code,
	o.orders_type
FROM ppm_shipments ps
	JOIN mto_shipments ms on ms.id = ps.shipment_id
	JOIN moves m ON ms.move_id  = m.id
	JOIN orders o ON m.orders_id = o.id
	JOIN service_members sm ON o.service_member_id = sm.id
	JOIN duty_locations dln ON o.new_duty_location_id = dln.id
	JOIN duty_locations dln2 ON o.origin_duty_location_id = dln2.id
	JOIN addresses da ON ps.destination_postal_address_id = da.id
	WHERE ps.status = 'WAITING_ON_CUSTOMER'::public."ppm_shipment_status"
	AND ms.status = 'APPROVED'::public."mto_shipment_status"
	AND ps.expected_departure_date <= now() - ($1)::interval
	AND ps.expected_departure_date  >= $2
	AND NOT EXISTS (
        SELECT 1 FROM notifications n
        WHERE sm.id = n.service_member_id
		AND n.notification_type  = 'MOVE_PAYMENT_REMINDER_EMAIL'
    )`

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

// TODO: rename to DestinationLocation
// formatEmails formats email data using both html and text template
func (m PaymentReminder) formatEmails(appCtx appcontext.AppContext, PaymentReminderEmailInfos PaymentReminderEmailInfos) ([]emailContent, error) {
	var emails []emailContent
	for _, PaymentReminderEmailInfo := range PaymentReminderEmailInfos {
		htmlBody, textBody, err := m.renderTemplates(appCtx, PaymentReminderEmailData{
			OriginDutyLocation:      PaymentReminderEmailInfo.OriginDutyLocationName,
			DestinationDutyLocation: getDestinationLocation(appCtx, PaymentReminderEmailInfo),
			Locator:                 PaymentReminderEmailInfo.Locator,
			OneSourceLink:           OneSourceTransportationOfficeLink,
			MyMoveLink:              MyMoveLink,
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
			subject:        "Complete your Personally Procured Move (PPM)",
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

func getDestinationLocation(appCtx appcontext.AppContext, PaymentReminderEmailInfo PaymentReminderEmailInfo) string {
	destinationLocation := PaymentReminderEmailInfo.NewDutyLocationName
	ordersType := PaymentReminderEmailInfo.OrdersType
	street1 := PaymentReminderEmailInfo.DestinationStreet1
	if street1 != nil {
		appCtx.Logger().Error("Street1 is: " + *street1)
	} else {
		appCtx.Logger().Error("Street1 is nil")
	}
	isSeparateeOrRetireeOrder := ordersType == internalmessages.OrdersTypeRETIREMENT || ordersType == internalmessages.OrdersTypeSEPARATION
	if isSeparateeOrRetireeOrder {
		appCtx.Logger().Debug("isSeparateeOrRetireeOrder: true")
	} else {
		appCtx.Logger().Debug("isSeparateeOrRetireeOrder: false")
	}
	if isSeparateeOrRetireeOrder && street1 != nil {
		appCtx.Logger().Debug("In address section")
		street2, street3, city, state, postalCode := "", "", "", "", ""
		if PaymentReminderEmailInfo.DestinationStreet2 != nil {
			street2 = " " + *PaymentReminderEmailInfo.DestinationStreet2
		}
		if PaymentReminderEmailInfo.DestinationStreet3 != nil {
			street3 = " " + *PaymentReminderEmailInfo.DestinationStreet3
		}
		if PaymentReminderEmailInfo.DestinationCity != nil {
			city = ", " + *PaymentReminderEmailInfo.DestinationCity
		}
		if PaymentReminderEmailInfo.DestinationState != nil {
			state = ", " + *PaymentReminderEmailInfo.DestinationState
		}
		if PaymentReminderEmailInfo.DestinationPostalCode != nil {
			postalCode = " " + *PaymentReminderEmailInfo.DestinationPostalCode
		}
		destinationLocation = fmt.Sprintf("%s%s%s%s%s%s", *street1, street2, street3, city, state, postalCode)
		appCtx.Logger().Debug("New location: " + destinationLocation)
	}
	return destinationLocation
}

// OnSuccess callback passed to be invoked by NewNotificationSender when an email successfully sent
// saves the svs the email info along with the SES mail id to the notifications table
func (m PaymentReminder) OnSuccess(appCtx appcontext.AppContext, PaymentReminderEmailInfo PaymentReminderEmailInfo) func(string) error {
	return func(msgID string) error {
		notification := models.Notification{
			ServiceMemberID:  PaymentReminderEmailInfo.ServiceMemberID,
			SESMessageID:     msgID,
			NotificationType: models.MovePaymentReminderEmail,
		}
		err := appCtx.DB().Create(&notification)
		if err != nil {
			dataString := fmt.Sprintf("%#v", notification)
			appCtx.Logger().Error("adding notification to notifications table", zap.String("notification", dataString))
			return err
		}

		return nil
	}
}

// PaymentReminderEmailData is used to render an email template
type PaymentReminderEmailData struct {
	OriginDutyLocation      string
	DestinationDutyLocation string
	Locator                 string
	OneSourceLink           string
	MyMoveLink              string
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
