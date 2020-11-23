package notifications

import (
	"bytes"
	"context"
	"fmt"
	html "html/template"
	text "text/template"

	"github.com/gobuffalo/pop/v5"
	"github.com/gofrs/uuid"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/assets"
	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/models"
)

var (
	moveSubmittedRawTextTemplate = string(assets.MustAsset("pkg/notifications/templates/move_submitted_template.txt"))
	moveSubmittedTextTemplate    = text.Must(text.New("text_template").Parse(moveSubmittedRawTextTemplate))
	moveSubmittedRawHTMLTemplate = string(assets.MustAsset("pkg/notifications/templates/move_submitted_template.html"))
	moveSubmittedHTMLTemplate    = html.Must(html.New("text_template").Parse(moveSubmittedRawHTMLTemplate))
)

// MoveSubmitted has notification content for submitted moves
type MoveSubmitted struct {
	db           *pop.Connection
	logger       Logger
	moveID       uuid.UUID
	session      *auth.Session // TODO - remove this when we move permissions up to handlers and out of models
	htmlTemplate *html.Template
	textTemplate *text.Template
}

// NewMoveSubmitted returns a new move submitted notification
func NewMoveSubmitted(db *pop.Connection, logger Logger, session *auth.Session, moveID uuid.UUID) *MoveSubmitted {

	return &MoveSubmitted{
		db:           db,
		logger:       logger,
		moveID:       moveID,
		session:      session,
		htmlTemplate: moveSubmittedHTMLTemplate,
		textTemplate: moveSubmittedTextTemplate,
	}
}

func (m MoveSubmitted) emails(ctx context.Context) ([]emailContent, error) {
	var emails []emailContent

	move, err := models.FetchMove(m.db, m.session, m.moveID)
	if err != nil {
		return emails, err
	}

	orders, err := models.FetchOrderForUser(m.db, m.session, move.OrdersID)
	if err != nil {
		return emails, err
	}

	serviceMember, err := models.FetchServiceMemberForUser(ctx, m.db, m.session, orders.ServiceMemberID)
	if err != nil {
		return emails, err
	}

	originDSTransportInfo, err := models.FetchDSContactInfo(m.db, serviceMember.DutyStationID)
	if err != nil {
		return emails, err
	}

	totalEntitlement, err := models.GetEntitlement(*serviceMember.Rank, orders.HasDependents, orders.SpouseHasProGear)
	if err != nil {
		return emails, err
	}

	if serviceMember.PersonalEmail == nil {
		return emails, fmt.Errorf("no email found for service member")
	}

	htmlBody, textBody, err := m.renderTemplates(moveSubmittedEmailData{
		Link:                       "https://my.move.mil/",
		PpmLink:                    "https://office.move.mil/downloads/ppm_info_sheet.pdf",
		OriginDutyStation:          originDSTransportInfo.Name,
		DestinationDutyStation:     orders.NewDutyStation.Name,
		OriginDutyStationPhoneLine: originDSTransportInfo.PhoneLine,
		Locator:                    move.Locator,
		WeightAllowance:            totalEntitlement,
	})

	if err != nil {
		m.logger.Error("error rendering template", zap.Error(err))
	}

	smEmail := emailContent{
		recipientEmail: *serviceMember.PersonalEmail,
		subject:        "Thank you for submitting your move details",
		htmlBody:       htmlBody,
		textBody:       textBody,
	}

	m.logger.Info("Generated move submitted email",
		zap.String("moveLocator", move.Locator))

	// TODO: Send email to trusted contacts when that's supported
	return append(emails, smEmail), nil
}

func (m MoveSubmitted) renderTemplates(data moveSubmittedEmailData) (string, string, error) {
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

type moveSubmittedEmailData struct {
	Link                       string
	PpmLink                    string
	OriginDutyStation          string
	DestinationDutyStation     string
	OriginDutyStationPhoneLine string
	Locator                    string
	WeightAllowance            int
}

// RenderHTML renders the html for the email
func (m MoveSubmitted) RenderHTML(data moveSubmittedEmailData) (string, error) {
	var htmlBuffer bytes.Buffer
	if err := m.htmlTemplate.Execute(&htmlBuffer, data); err != nil {
		m.logger.Error("cant render html template ", zap.Error(err))
	}
	return htmlBuffer.String(), nil
}

// RenderText renders the text for the email
func (m MoveSubmitted) RenderText(data moveSubmittedEmailData) (string, error) {
	var textBuffer bytes.Buffer
	if err := m.textTemplate.Execute(&textBuffer, data); err != nil {
		m.logger.Error("cant render text template ", zap.Error(err))
		return "", err
	}
	return textBuffer.String(), nil
}
