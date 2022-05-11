package notifications

import (
	"bytes"
	"fmt"
	html "html/template"
	text "text/template"

	"github.com/dustin/go-humanize"
	"github.com/gofrs/uuid"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/assets"
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
	moveID       uuid.UUID
	htmlTemplate *html.Template
	textTemplate *text.Template
}

// NewMoveSubmitted returns a new move submitted notification
func NewMoveSubmitted(moveID uuid.UUID) *MoveSubmitted {

	return &MoveSubmitted{
		moveID:       moveID,
		htmlTemplate: moveSubmittedHTMLTemplate,
		textTemplate: moveSubmittedTextTemplate,
	}
}

func (m MoveSubmitted) emails(appCtx appcontext.AppContext) ([]emailContent, error) {
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

	originDSTransportInfo, err := models.FetchDLContactInfo(appCtx.DB(), serviceMember.DutyLocationID)
	if err != nil {
		return emails, err
	}

	var originDutyLocation, originDutyLocationPhoneLine *string
	if originDSTransportInfo != nil {
		originDutyLocation = &originDSTransportInfo.Name
		originDutyLocationPhoneLine = &originDSTransportInfo.PhoneLine

	}

	totalEntitlement, err := models.GetEntitlement(*serviceMember.Rank, orders.HasDependents)
	if err != nil {
		return emails, err
	}

	if serviceMember.PersonalEmail == nil {
		return emails, fmt.Errorf("no email found for service member")
	}

	htmlBody, textBody, err := m.renderTemplates(appCtx, moveSubmittedEmailData{
		Link:                        "https://my.move.mil/",
		PpmLink:                     "https://office.move.mil/downloads/ppm_info_sheet.pdf",
		OriginDutyLocation:          originDutyLocation,
		DestinationDutyLocation:     orders.NewDutyLocation.Name,
		OriginDutyLocationPhoneLine: originDutyLocationPhoneLine,
		Locator:                     move.Locator,
		WeightAllowance:             humanize.Comma(int64(totalEntitlement)),
	})

	if err != nil {
		appCtx.Logger().Error("error rendering template", zap.Error(err))
	}

	smEmail := emailContent{
		recipientEmail: *serviceMember.PersonalEmail,
		subject:        "Thank you for submitting your move details",
		htmlBody:       htmlBody,
		textBody:       textBody,
	}

	appCtx.Logger().Info("Generated move submitted email",
		zap.String("moveLocator", move.Locator))

	// TODO: Send email to trusted contacts when that's supported
	return append(emails, smEmail), nil
}

func (m MoveSubmitted) renderTemplates(appCtx appcontext.AppContext, data moveSubmittedEmailData) (string, string, error) {
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

type moveSubmittedEmailData struct {
	Link                        string
	PpmLink                     string
	OriginDutyLocation          *string
	DestinationDutyLocation     string
	OriginDutyLocationPhoneLine *string
	Locator                     string
	WeightAllowance             string
}

// RenderHTML renders the html for the email
func (m MoveSubmitted) RenderHTML(appCtx appcontext.AppContext, data moveSubmittedEmailData) (string, error) {
	var htmlBuffer bytes.Buffer
	if err := m.htmlTemplate.Execute(&htmlBuffer, data); err != nil {
		appCtx.Logger().Error("cant render html template ", zap.Error(err))
	}
	return htmlBuffer.String(), nil
}

// RenderText renders the text for the email
func (m MoveSubmitted) RenderText(appCtx appcontext.AppContext, data moveSubmittedEmailData) (string, error) {
	var textBuffer bytes.Buffer
	if err := m.textTemplate.Execute(&textBuffer, data); err != nil {
		appCtx.Logger().Error("cant render text template ", zap.Error(err))
		return "", err
	}
	return textBuffer.String(), nil
}
