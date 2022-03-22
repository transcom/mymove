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
	moveID       uuid.UUID
	htmlTemplate *html.Template
	textTemplate *text.Template
}

// NewMoveCanceled returns a new move approval notification
func NewMoveCanceled(moveID uuid.UUID) *MoveCanceled {

	return &MoveCanceled{
		moveID:       moveID,
		htmlTemplate: moveCanceledHTMLTemplate,
		textTemplate: moveCanceledTextTemplate,
	}
}

func (m MoveCanceled) emails(appCtx appcontext.AppContext) ([]emailContent, error) {
	var emails []emailContent

	move, err := models.FetchMove(appCtx.DB(), appCtx.Session(), m.moveID)
	if err != nil {
		return emails, err
	}

	orders, err := models.FetchOrderForUser(appCtx.DB(), appCtx.Session(), move.OrdersID)
	if err != nil {
		return emails, err
	}

	serviceMember, err := models.FetchServiceMemberForUser(appCtx.DB(), appCtx.Session(), orders.ServiceMemberID)
	if err != nil {
		return emails, err
	}

	if serviceMember.PersonalEmail == nil {
		return emails, fmt.Errorf("no email found for service member")
	}

	dsTransportInfo, err := models.FetchDLContactInfo(appCtx.DB(), serviceMember.DutyLocationID)
	if err != nil {
		return emails, err
	}

	var originDutyLocation, originDutyLocationPhoneLine *string
	if dsTransportInfo != nil {
		originDutyLocation = &dsTransportInfo.Name
		originDutyLocationPhoneLine = &dsTransportInfo.PhoneLine
	}

	if orders.NewDutyLocation.Name == "" {
		return emails, fmt.Errorf("missing new duty station for service member")
	}

	// Set up various text segments. Copy comes from here:
	// https://docs.google.com/document/d/1gIQZprWzJJE_sAAyg5NViPwy9ckL5RK37gFq1fEfipU
	// TODO: we will want some sort of templating system

	htmlBody, textBody, err := m.renderTemplates(appCtx, moveCanceledEmailData{
		OriginDutyLocation:          originDutyLocation,
		DestinationDutyLocation:     orders.NewDutyLocation.Name,
		OriginDutyLocationPhoneLine: originDutyLocationPhoneLine,
		Locator:                     move.Locator,
	})

	if err != nil {
		appCtx.Logger().Error("error rendering template", zap.Error(err))
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

func (m MoveCanceled) renderTemplates(appCtx appcontext.AppContext, data moveCanceledEmailData) (string, string, error) {
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

type moveCanceledEmailData struct {
	OriginDutyLocation          *string
	DestinationDutyLocation     string
	OriginDutyLocationPhoneLine *string
	Locator                     string
}

// RenderHTML renders the html for the email
func (m MoveCanceled) RenderHTML(appCtx appcontext.AppContext, data moveCanceledEmailData) (string, error) {
	var htmlBuffer bytes.Buffer
	if err := m.htmlTemplate.Execute(&htmlBuffer, data); err != nil {
		appCtx.Logger().Error("cant render html template ", zap.Error(err))
	}
	return htmlBuffer.String(), nil
}

// RenderText renders the text for the email
func (m MoveCanceled) RenderText(appCtx appcontext.AppContext, data moveCanceledEmailData) (string, error) {
	var textBuffer bytes.Buffer
	if err := m.textTemplate.Execute(&textBuffer, data); err != nil {
		appCtx.Logger().Error("cant render text template ", zap.Error(err))
		return "", err
	}
	return textBuffer.String(), nil
}
