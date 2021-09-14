package notifications

import (
	"bytes"
	"fmt"
	html "html/template"
	"strings"
	text "text/template"

	"github.com/gobuffalo/pop/v5"
	"github.com/gofrs/uuid"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/assets"
	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/models"
)

var (
	reweighRequestedRawTextTemplate = string(assets.MustAsset("pkg/notifications/templates/reweigh_requested_template.txt"))
	reweighRequestedTextTemplate    = text.Must(text.New("text_template").Parse(reweighRequestedRawTextTemplate))
	reweighRequestedRawHTMLTemplate = string(assets.MustAsset("pkg/notifications/templates/reweigh_requested_template.html"))
	reweighRequestedHTMLTemplate    = html.Must(html.New("text_template").Parse(reweighRequestedRawHTMLTemplate))
)

// ReweighRequested has notification content for submitted moves
type ReweighRequested struct {
	db           *pop.Connection
	logger       Logger
	moveID       uuid.UUID
	shipment     models.MTOShipment
	session      *auth.Session // TODO - remove this when we move permissions up to handlers and out of models
	htmlTemplate *html.Template
	textTemplate *text.Template
}

// NewReweighRequested returns a new move submitted notification
func NewReweighRequested(db *pop.Connection, logger Logger, session *auth.Session, moveID uuid.UUID, shipment models.MTOShipment) *ReweighRequested {

	return &ReweighRequested{
		db:           db,
		logger:       logger,
		moveID:       moveID,
		shipment:     shipment,
		session:      session,
		htmlTemplate: reweighRequestedHTMLTemplate,
		textTemplate: reweighRequestedTextTemplate,
	}
}

func (m ReweighRequested) emails() ([]emailContent, error) {
	var emails []emailContent

	var serviceMember models.ServiceMember
	err := m.db.Q().
		InnerJoin("orders", "orders.service_member_id = service_members.id").
		InnerJoin("moves", "moves.orders_id = orders.id").
		Where("moves.id = ?", m.shipment.MoveTaskOrderID).
		First(&serviceMember)
	if err != nil {
		return emails, fmt.Errorf("error fetching service member email for shipment ID: %s with error %w", m.shipment.MoveTaskOrderID, err)
	}

	htmlBody, textBody, err := m.renderTemplates(reweighRequestedEmailData{})

	if err != nil {
		m.logger.Error("error rendering template", zap.Error(err))
	}

	shipmentType := strings.Split(string(m.shipment.ShipmentType), "_")[0]

	smEmail := emailContent{
		recipientEmail: *serviceMember.PersonalEmail,
		subject:        fmt.Sprintf("FYI: Your %v should be reweighed before it is delivered", shipmentType),
		htmlBody:       htmlBody,
		textBody:       textBody,
	}

	m.logger.Info("Generated reweigh requested email",
		zap.String("moveLocator", m.shipment.MoveTaskOrder.Locator))

	// TODO: Send email to trusted contacts when that's supported
	return append(emails, smEmail), nil
}

func (m ReweighRequested) renderTemplates(data reweighRequestedEmailData) (string, string, error) {
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

type reweighRequestedEmailData struct{}

// RenderHTML renders the html for the email
func (m ReweighRequested) RenderHTML(data reweighRequestedEmailData) (string, error) {
	var htmlBuffer bytes.Buffer
	if err := m.htmlTemplate.Execute(&htmlBuffer, data); err != nil {
		m.logger.Error("cant render html template ", zap.Error(err))
	}
	return htmlBuffer.String(), nil
}

// RenderText renders the text for the email
func (m ReweighRequested) RenderText(data reweighRequestedEmailData) (string, error) {
	var textBuffer bytes.Buffer
	if err := m.textTemplate.Execute(&textBuffer, data); err != nil {
		m.logger.Error("cant render text template ", zap.Error(err))
		return "", err
	}
	return textBuffer.String(), nil
}
