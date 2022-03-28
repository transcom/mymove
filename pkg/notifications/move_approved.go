package notifications

import (
	"bytes"
	"fmt"
	html "html/template"
	"net/url"
	text "text/template"

	"github.com/gofrs/uuid"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/assets"
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
	host         string
	moveID       uuid.UUID
	htmlTemplate *html.Template
	textTemplate *text.Template
}

// NewMoveApproved returns a new move approval notification
func NewMoveApproved(
	host string,
	moveID uuid.UUID) *MoveApproved {

	return &MoveApproved{
		host:         host,
		moveID:       moveID,
		htmlTemplate: moveApprovedHTMLTemplate,
		textTemplate: moveApprovedTextTemplate,
	}
}

func (m MoveApproved) emails(appCtx appcontext.AppContext) ([]emailContent, error) {
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

	var originDutyLocation, originDutyLocationPhoneLine *string
	dsTransportInfo, err := models.FetchDLContactInfo(appCtx.DB(), serviceMember.DutyLocationID)
	if err != nil {
		return emails, err
	}

	if dsTransportInfo != nil {
		originDutyLocation = &dsTransportInfo.Name
		originDutyLocationPhoneLine = &dsTransportInfo.PhoneLine
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

	htmlBody, textBody, err := m.renderTemplates(appCtx, moveApprovedEmailData{
		Link:                        ppmInfoSheetURL.String(),
		OriginDutyLocation:          originDutyLocation,
		DestinationDutyLocation:     orders.NewDutyLocation.Name,
		OriginDutyLocationPhoneLine: originDutyLocationPhoneLine,
		Locator:                     move.Locator,
	})

	if err != nil {
		appCtx.Logger().Error("error rendering template", zap.Error(err))
	}

	smEmail := emailContent{
		subject:        fmt.Sprintf("[MilMove] Your Move is approved (move: %s)", move.Locator),
		recipientEmail: *serviceMember.PersonalEmail,
		htmlBody:       htmlBody,
		textBody:       textBody,
	}

	appCtx.Logger().Info("generated move approved email to service member",
		zap.String("service member uuid", serviceMember.ID.String()),
		zap.String("moveLocator", move.Locator),
	)

	// TODO: Send email to trusted contacts when that's supported
	return append(emails, smEmail), nil
}

func (m MoveApproved) renderTemplates(appCtx appcontext.AppContext, data moveApprovedEmailData) (string, string, error) {
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

// moveApprovedEmailData has content for email template
type moveApprovedEmailData struct {
	Link                        string
	OriginDutyLocation          *string
	DestinationDutyLocation     string
	OriginDutyLocationPhoneLine *string
	Locator                     string
}

// RenderHTML renders the html for the email
func (m MoveApproved) RenderHTML(appCtx appcontext.AppContext, data moveApprovedEmailData) (string, error) {
	var htmlBuffer bytes.Buffer
	if err := m.htmlTemplate.Execute(&htmlBuffer, data); err != nil {
		appCtx.Logger().Error("cant render html template ", zap.Error(err))
	}
	return htmlBuffer.String(), nil
}

// RenderText renders the text for the email
func (m MoveApproved) RenderText(appCtx appcontext.AppContext, data moveApprovedEmailData) (string, error) {
	var textBuffer bytes.Buffer
	if err := m.textTemplate.Execute(&textBuffer, data); err != nil {
		appCtx.Logger().Error("cant render text template ", zap.Error(err))
		return "", err
	}
	return textBuffer.String(), nil
}
