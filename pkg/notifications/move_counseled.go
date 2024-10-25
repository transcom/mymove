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
)

var (
	moveCounseledRawTextTemplate = string(assets.MustAsset("notifications/templates/move_counseled_template.txt"))
	moveCounseledTextTemplate    = text.Must(text.New("text_template").Parse(moveCounseledRawTextTemplate))
	moveCounseledRawHTMLTemplate = string(assets.MustAsset("notifications/templates/move_counseled_template.html"))
	moveCounseledHTMLTemplate    = html.Must(html.New("text_template").Parse(moveCounseledRawHTMLTemplate))
)

// MoveCounseled has notification content for counseled moves (before TOO approval)
type MoveCounseled struct {
	moveID       uuid.UUID
	htmlTemplate *html.Template
	textTemplate *text.Template
}

// NewMoveCounseled returns a new move counseled notification (before TOO approval)
func NewMoveCounseled(moveID uuid.UUID) *MoveCounseled {

	return &MoveCounseled{
		moveID:       moveID,
		htmlTemplate: moveCounseledHTMLTemplate,
		textTemplate: moveCounseledTextTemplate,
	}
}

func (m MoveCounseled) emails(appCtx appcontext.AppContext) ([]emailContent, error) {
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

	originDSTransportInfo, err := models.FetchDLContactInfo(appCtx.DB(), orders.OriginDutyLocationID)
	if err != nil {
		return emails, err
	}

	var originDutyLocationName *string
	if originDSTransportInfo != nil {
		originDutyLocationName = &originDSTransportInfo.Name
	}

	destinationAddress := orders.NewDutyLocation.Name
	isSeparateeOrRetireeOrder := orders.OrdersType == internalmessages.OrdersTypeRETIREMENT || orders.OrdersType == internalmessages.OrdersTypeSEPARATION
	if isSeparateeOrRetireeOrder && len(move.MTOShipments) > 0 && move.MTOShipments[0].DestinationAddress != nil {
		mtoShipDestinationAddress, streetAddr2, streetAddr3 := *move.MTOShipments[0].DestinationAddress, "", ""
		if mtoShipDestinationAddress.StreetAddress2 != nil {
			streetAddr2 = " " + *mtoShipDestinationAddress.StreetAddress2
		}
		if mtoShipDestinationAddress.StreetAddress3 != nil {
			streetAddr3 = " " + *mtoShipDestinationAddress.StreetAddress3
		}
		destinationAddress = fmt.Sprintf("%s%s%s, %s, %s %s", mtoShipDestinationAddress.StreetAddress1, streetAddr2, streetAddr3, mtoShipDestinationAddress.City, mtoShipDestinationAddress.State, mtoShipDestinationAddress.PostalCode)
	}

	if serviceMember.PersonalEmail == nil {
		return emails, fmt.Errorf("no email found for service member")
	}

	htmlBody, textBody, err := m.renderTemplates(appCtx, MoveCounseledEmailData{
		OriginDutyLocation:         originDutyLocationName,
		DestinationLocation:        destinationAddress,
		Locator:                    move.Locator,
		MyMoveLink:                 MyMoveLink,
		ActualExpenseReimbursement: true,
	})

	if err != nil {
		appCtx.Logger().Error("error rendering template", zap.Error(err))
	}

	smEmail := emailContent{
		recipientEmail: *serviceMember.PersonalEmail,
		subject:        "Your counselor has approved your move details",
		htmlBody:       htmlBody,
		textBody:       textBody,
	}

	appCtx.Logger().Info("Generated move counseled email",
		zap.String("moveLocator", move.Locator))

	// TODO: Send email to trusted contacts when that's supported
	return append(emails, smEmail), nil
}

func (m MoveCounseled) renderTemplates(appCtx appcontext.AppContext, data MoveCounseledEmailData) (string, string, error) {
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

type MoveCounseledEmailData struct {
	OriginDutyLocation         *string
	DestinationLocation        string
	Locator                    string
	MyMoveLink                 string
	ActualExpenseReimbursement bool
}

// RenderHTML renders the html for the email
func (m MoveCounseled) RenderHTML(appCtx appcontext.AppContext, data MoveCounseledEmailData) (string, error) {
	var htmlBuffer bytes.Buffer
	if err := m.htmlTemplate.Execute(&htmlBuffer, data); err != nil {
		appCtx.Logger().Error("cant render html template ", zap.Error(err))
	}
	return htmlBuffer.String(), nil
}

// RenderText renders the text for the email
func (m MoveCounseled) RenderText(appCtx appcontext.AppContext, data MoveCounseledEmailData) (string, error) {
	var textBuffer bytes.Buffer
	if err := m.textTemplate.Execute(&textBuffer, data); err != nil {
		appCtx.Logger().Error("cant render text template ", zap.Error(err))
		return "", err
	}
	return textBuffer.String(), nil
}
