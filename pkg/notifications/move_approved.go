package notifications

import (
	"bytes"
	"context"
	"fmt"
	html "html/template"
	"net/url"
	text "text/template"

	"github.com/gobuffalo/pop/v5"
	"github.com/gofrs/uuid"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/assets"
	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/models"
)

var (
	moveApprovedRawTextTemplate = string(assets.MustAsset("pkg/notifications/templates/move_approved_template.txt"))
	moveApprovedTextTemplate    = text.Must(text.New("text_template").Parse(moveApprovedRawTextTemplate))
	moveApprovedRawHTMLTemplate = string(assets.MustAsset("pkg/notifications/templates/move_approved_template.html"))
	moveApprovedHTMLTemplate    = html.Must(html.New("text_template").Parse(moveApprovedRawHTMLTemplate))
)

// MoveApproved has notification content for approved moves
type MoveApproved struct {
	db           *pop.Connection
	logger       Logger
	host         string
	moveID       uuid.UUID
	session      *auth.Session // TODO - remove this when we move permissions up to handlers and out of models
	htmlTemplate *html.Template
	textTemplate *text.Template
}

// NewMoveApproved returns a new move approval notification
func NewMoveApproved(db *pop.Connection,
	logger Logger,
	session *auth.Session,
	host string,
	moveID uuid.UUID) *MoveApproved {

	return &MoveApproved{
		db:           db,
		logger:       logger,
		host:         host,
		moveID:       moveID,
		session:      session,
		htmlTemplate: moveApprovedHTMLTemplate,
		textTemplate: moveApprovedTextTemplate,
	}
}

func (m MoveApproved) emails(ctx context.Context) ([]emailContent, error) {
	var emails []emailContent

	move, err := models.FetchMove(m.db, m.session, m.moveID, nil)
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

	dsTransportInfo, err := models.FetchDSContactInfo(m.db, serviceMember.DutyStationID)
	if err != nil {
		return emails, err
	}

	if serviceMember.PersonalEmail == nil {
		return emails, fmt.Errorf("no email found for service member")
	}

	// Set up various text segments. Copy comes from here:
	// https://docs.google.com/document/d/1Hm8qTm_biHmdT5LCyHY8QJxqXqlDPGJMnhDk_0LE5Gc/
	// TODO: we will want some sort of templating system
	ppmInfoSheetURL := url.URL{
		Scheme: "https",
		Host:   m.host,
		Path:   "downloads/ppm_info_sheet.pdf",
	}

	htmlBody, textBody, err := m.renderTemplates(moveApprovedEmailData{
		Link:                       ppmInfoSheetURL.String(),
		OriginDutyStation:          dsTransportInfo.Name,
		DestinationDutyStation:     orders.NewDutyStation.Name,
		OriginDutyStationPhoneLine: dsTransportInfo.PhoneLine,
		Locator:                    move.Locator,
	})

	if err != nil {
		m.logger.Error("error rendering template", zap.Error(err))
	}

	smEmail := emailContent{
		subject:        fmt.Sprintf("[MilMove] Your Move is approved (move: %s)", move.Locator),
		recipientEmail: *serviceMember.PersonalEmail,
		htmlBody:       htmlBody,
		textBody:       textBody,
	}

	m.logger.Info("generated move approved email to service member",
		zap.String("service member uuid", serviceMember.ID.String()),
		zap.String("moveLocator", move.Locator),
	)

	// TODO: Send email to trusted contacts when that's supported
	return append(emails, smEmail), nil
}

func (m MoveApproved) renderTemplates(data moveApprovedEmailData) (string, string, error) {
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

// moveApprovedEmailData has content for email template
type moveApprovedEmailData struct {
	Link                       string
	OriginDutyStation          string
	DestinationDutyStation     string
	OriginDutyStationPhoneLine string
	Locator                    string
}

// RenderHTML renders the html for the email
func (m MoveApproved) RenderHTML(data moveApprovedEmailData) (string, error) {
	var htmlBuffer bytes.Buffer
	if err := m.htmlTemplate.Execute(&htmlBuffer, data); err != nil {
		m.logger.Error("cant render html template ", zap.Error(err))
	}
	return htmlBuffer.String(), nil
}

// RenderText renders the text for the email
func (m MoveApproved) RenderText(data moveApprovedEmailData) (string, error) {
	var textBuffer bytes.Buffer
	if err := m.textTemplate.Execute(&textBuffer, data); err != nil {
		m.logger.Error("cant render text template ", zap.Error(err))
		return "", err
	}
	return textBuffer.String(), nil
}
