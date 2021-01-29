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
	moveCanceledRawTextTemplate = string(assets.MustAsset("pkg/notifications/templates/move_canceled_template.txt"))
	moveCanceledTextTemplate    = text.Must(text.New("text_template").Parse(moveCanceledRawTextTemplate))
	moveCanceledRawHTMLTemplate = string(assets.MustAsset("pkg/notifications/templates/move_canceled_template.html"))
	moveCanceledHTMLTemplate    = html.Must(html.New("text_template").Parse(moveCanceledRawHTMLTemplate))
)

// MoveCanceled has notification content for approved moves
type MoveCanceled struct {
	db           *pop.Connection
	logger       Logger
	moveID       uuid.UUID
	session      *auth.Session // TODO - remove this when we move permissions up to handlers and out of models
	htmlTemplate *html.Template
	textTemplate *text.Template
}

// NewMoveCanceled returns a new move approval notification
func NewMoveCanceled(db *pop.Connection, logger Logger, session *auth.Session, moveID uuid.UUID) *MoveCanceled {

	return &MoveCanceled{
		db:           db,
		logger:       logger,
		moveID:       moveID,
		session:      session,
		htmlTemplate: moveCanceledHTMLTemplate,
		textTemplate: moveCanceledTextTemplate,
	}
}

func (m MoveCanceled) emails(ctx context.Context) ([]emailContent, error) {
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

	if serviceMember.PersonalEmail == nil {
		return emails, fmt.Errorf("no email found for service member")
	}

	dsTransportInfo, err := models.FetchDSContactInfo(m.db, serviceMember.DutyStationID)
	if err != nil {
		return emails, err
	}

	if orders.NewDutyStation.Name == "" {
		return emails, fmt.Errorf("missing new duty station for service member")
	}

	// Set up various text segments. Copy comes from here:
	// https://docs.google.com/document/d/1gIQZprWzJJE_sAAyg5NViPwy9ckL5RK37gFq1fEfipU
	// TODO: we will want some sort of templating system

	htmlBody, textBody, err := m.renderTemplates(moveCanceledEmailData{
		OriginDutyStation:          dsTransportInfo.Name,
		DestinationDutyStation:     orders.NewDutyStation.Name,
		OriginDutyStationPhoneLine: dsTransportInfo.PhoneLine,
	})

	if err != nil {
		m.logger.Error("error rendering template", zap.Error(err))
	}

	smEmail := emailContent{
		recipientEmail: *serviceMember.PersonalEmail,
		subject:        "[MilMove] Update on your move",
		htmlBody:       htmlBody,
		textBody:       textBody,
	}

	// TODO: Send email to trusted contacts when that's supported
	return append(emails, smEmail), nil
}

func (m MoveCanceled) renderTemplates(data moveCanceledEmailData) (string, string, error) {
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

type moveCanceledEmailData struct {
	OriginDutyStation          string
	DestinationDutyStation     string
	OriginDutyStationPhoneLine string
}

// RenderHTML renders the html for the email
func (m MoveCanceled) RenderHTML(data moveCanceledEmailData) (string, error) {
	var htmlBuffer bytes.Buffer
	if err := m.htmlTemplate.Execute(&htmlBuffer, data); err != nil {
		m.logger.Error("cant render html template ", zap.Error(err))
	}
	return htmlBuffer.String(), nil
}

// RenderText renders the text for the email
func (m MoveCanceled) RenderText(data moveCanceledEmailData) (string, error) {
	var textBuffer bytes.Buffer
	if err := m.textTemplate.Execute(&textBuffer, data); err != nil {
		m.logger.Error("cant render text template ", zap.Error(err))
		return "", err
	}
	return textBuffer.String(), nil
}
