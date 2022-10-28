package notifications

import (
	"bytes"
	"fmt"
	html "html/template"
	"strings"
	text "text/template"

	"github.com/gofrs/uuid"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/assets"
	"github.com/transcom/mymove/pkg/models"
)

var (
	reweighRequestedRawTextTemplate = string(assets.MustAsset("notifications/templates/reweigh_requested_template.txt"))
	reweighRequestedTextTemplate    = text.Must(text.New("text_template").Parse(reweighRequestedRawTextTemplate))
	reweighRequestedRawHTMLTemplate = string(assets.MustAsset("notifications/templates/reweigh_requested_template.html"))
	reweighRequestedHTMLTemplate    = html.Must(html.New("text_template").Parse(reweighRequestedRawHTMLTemplate))
)

// ReweighRequested has notification content for submitted moves
type ReweighRequested struct {
	moveID       uuid.UUID
	shipment     models.MTOShipment
	htmlTemplate *html.Template
	textTemplate *text.Template
}

// NewReweighRequested returns a new move submitted notification
func NewReweighRequested(moveID uuid.UUID, shipment models.MTOShipment) *ReweighRequested {

	return &ReweighRequested{
		moveID:       moveID,
		shipment:     shipment,
		htmlTemplate: reweighRequestedHTMLTemplate,
		textTemplate: reweighRequestedTextTemplate,
	}
}

func (m ReweighRequested) emails(appCtx appcontext.AppContext) ([]emailContent, error) {
	var emails []emailContent

	serviceMember, err := models.GetCustomerFromShipment(appCtx.DB(), m.shipment.ID)
	if err != nil {
		appCtx.Logger().Error("error retrieving service member associated with this shipment", zap.Error(err))
	}
	if len(*serviceMember.PersonalEmail) == 0 {
		return emails, fmt.Errorf("no email found for service member")
	}

	htmlBody, textBody, err := m.renderTemplates(appCtx, reweighRequestedEmailData{})

	if err != nil {
		appCtx.Logger().Error("error rendering template", zap.Error(err))
	}

	shipmentType := strings.Split(string(m.shipment.ShipmentType), "_")[0]

	smEmail := emailContent{
		recipientEmail: *serviceMember.PersonalEmail,
		subject:        fmt.Sprintf("FYI: Your %v should be reweighed before it is delivered", shipmentType),
		htmlBody:       htmlBody,
		textBody:       textBody,
	}

	appCtx.Logger().Info("Generated reweigh requested email",
		zap.String("moveLocator", m.shipment.MoveTaskOrder.Locator))

	// TODO: Send email to trusted contacts when that's supported
	return append(emails, smEmail), nil
}

func (m ReweighRequested) renderTemplates(appCtx appcontext.AppContext, data reweighRequestedEmailData) (string, string, error) {
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

type reweighRequestedEmailData struct{}

// RenderHTML renders the html for the email
func (m ReweighRequested) RenderHTML(appCtx appcontext.AppContext, data reweighRequestedEmailData) (string, error) {
	var htmlBuffer bytes.Buffer
	if err := m.htmlTemplate.Execute(&htmlBuffer, data); err != nil {
		appCtx.Logger().Error("cant render html template ", zap.Error(err))
	}
	return htmlBuffer.String(), nil
}

// RenderText renders the text for the email
func (m ReweighRequested) RenderText(appCtx appcontext.AppContext, data reweighRequestedEmailData) (string, error) {
	var textBuffer bytes.Buffer
	if err := m.textTemplate.Execute(&textBuffer, data); err != nil {
		appCtx.Logger().Error("cant render text template ", zap.Error(err))
		return "", err
	}
	return textBuffer.String(), nil
}
