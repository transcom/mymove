package notifications

import (
	"bytes"
	"fmt"
	html "html/template"
	text "text/template"

	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/assets"
	"github.com/transcom/mymove/pkg/gen/primemessages"
)

var (
	PrimeCounselingCompleteRawText      = string(assets.MustAsset("notifications/templates/prime_counseling_complete_template.txt"))
	PrimeCounselingCompleteTextTemplate = text.Must(text.New("text_template").Parse(PrimeCounselingCompleteRawText))
	PrimeCounselingCompleteRawHTML      = string(assets.MustAsset("notifications/templates/prime_counseling_complete_template.html"))
	PrimeCounselingCompleteHTMLTemplate = html.Must(html.New("text_template").Parse(PrimeCounselingCompleteRawHTML))
)

// PrimeCounselingComplete has notification content for moves that have had their counseling completed by the Prime
type PrimeCounselingComplete struct {
	moveTaskOrder primemessages.MoveTaskOrder
	htmlTemplate  *html.Template
	textTemplate  *text.Template
}

// PrimeCounselingCompleteData is used to render an email template
type PrimeCounselingCompleteData struct {
	CustomerEmail                     string
	OriginDutyLocation                string
	DestinationDutyLocation           string
	Locator                           string
	OneSourceTransportationOfficeLink string
	MyMoveLink                        string
}

// NewPrimeCounselingComplete returns a new payment reminder notification 14 days after actual move in date
func NewPrimeCounselingComplete(moveTaskOrder primemessages.MoveTaskOrder) *PrimeCounselingComplete {

	return &PrimeCounselingComplete{
		moveTaskOrder: moveTaskOrder,
		htmlTemplate:  PrimeCounselingCompleteHTMLTemplate,
		textTemplate:  PrimeCounselingCompleteTextTemplate,
	}
}

// NotificationSendingContext expects a `notification` with an `emails` method,
// so we implement `email` to satisfy that interface
func (p PrimeCounselingComplete) emails(appCtx appcontext.AppContext) ([]emailContent, error) {
	var emails []emailContent

	appCtx.Logger().Info("MTO (Move Task Order) Locator",
		zap.String("uuid", p.moveTaskOrder.MoveCode),
	)

	emailData, err := p.GetEmailData(p.moveTaskOrder, appCtx)
	if err != nil {
		return nil, err
	}
	var htmlBody, textBody string
	htmlBody, textBody, err = p.renderTemplates(appCtx, emailData)

	if err != nil {
		appCtx.Logger().Error("error rendering template", zap.Error(err))
	}

	primeCounselingEmail := emailContent{
		recipientEmail: emailData.CustomerEmail,
		subject:        "Your counselor has approved your move details",
		htmlBody:       htmlBody,
		textBody:       textBody,
	}

	return append(emails, primeCounselingEmail), nil
}

func (p PrimeCounselingComplete) GetEmailData(m primemessages.MoveTaskOrder, appCtx appcontext.AppContext) (PrimeCounselingCompleteData, error) {
	if m.Order.Customer.Email == "" {
		return PrimeCounselingCompleteData{}, fmt.Errorf("no email found for service member")
	}

	appCtx.Logger().Info("generated Prime Counseling Completed email",
		zap.String("service member uuid", string(m.Order.Customer.ID)),
		zap.String("service member email", string(m.Order.Customer.Email)),
		zap.String("Move Locator", string(m.MoveCode)),
		zap.String("Origin Duty Location Name", string(m.Order.OriginDutyLocation.Name)),
		zap.String("Destination Duty Location Name", string(m.Order.DestinationDutyLocation.Name)),
	)

	return PrimeCounselingCompleteData{
		CustomerEmail:                     m.Order.Customer.Email,
		OriginDutyLocation:                m.Order.OriginDutyLocation.Name,
		DestinationDutyLocation:           m.Order.DestinationDutyLocation.Name,
		Locator:                           m.MoveCode,
		OneSourceTransportationOfficeLink: OneSourceTransportationOfficeLink,
		MyMoveLink:                        MyMoveLink,
	}, nil
}

func (p PrimeCounselingComplete) renderTemplates(appCtx appcontext.AppContext, data PrimeCounselingCompleteData) (string, string, error) {
	htmlBody, err := p.RenderHTML(appCtx, data)
	if err != nil {
		return "", "", fmt.Errorf("error rendering html template using %#v", data)
	}
	textBody, err := p.RenderText(appCtx, data)
	if err != nil {
		return "", "", fmt.Errorf("error rendering text template using %#v", data)
	}
	return htmlBody, textBody, nil
}

// RenderHTML renders the html for the email
func (p PrimeCounselingComplete) RenderHTML(appCtx appcontext.AppContext, data PrimeCounselingCompleteData) (string, error) {
	var htmlBuffer bytes.Buffer
	if err := p.htmlTemplate.Execute(&htmlBuffer, data); err != nil {
		appCtx.Logger().Error("cant render html template ", zap.Error(err))
	}
	return htmlBuffer.String(), nil
}

// RenderText renders the text for the email
func (p PrimeCounselingComplete) RenderText(appCtx appcontext.AppContext, data PrimeCounselingCompleteData) (string, error) {
	var textBuffer bytes.Buffer
	if err := p.textTemplate.Execute(&textBuffer, data); err != nil {
		appCtx.Logger().Error("cant render text template ", zap.Error(err))
		return "", err
	}
	return textBuffer.String(), nil
}
