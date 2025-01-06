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
	moveIssuedToPrimeRawTextTemplate = string(assets.MustAsset("notifications/templates/move_issued_to_prime_template.txt"))
	moveIssuedToPrimeTextTemplate    = text.Must(text.New("text_template").Parse(moveIssuedToPrimeRawTextTemplate))
	moveIssuedToPrimeRawHTMLTemplate = string(assets.MustAsset("notifications/templates/move_issued_to_prime_template.html"))
	moveIssuedToPrimeHTMLTemplate    = html.Must(html.New("text_template").Parse(moveIssuedToPrimeRawHTMLTemplate))
)

// MoveIssuedToPrime has notification content for submitted moves
type MoveIssuedToPrime struct {
	moveID       uuid.UUID
	htmlTemplate *html.Template
	textTemplate *text.Template
}

// NewMoveIssuedToPrime returns a new move submitted notification
func NewMoveIssuedToPrime(moveID uuid.UUID) *MoveIssuedToPrime {

	return &MoveIssuedToPrime{
		moveID:       moveID,
		htmlTemplate: moveIssuedToPrimeHTMLTemplate,
		textTemplate: moveIssuedToPrimeTextTemplate,
	}
}

func (m MoveIssuedToPrime) emails(appCtx appcontext.AppContext) ([]emailContent, error) {
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

	var originDutyLocation *string
	if originDSTransportInfo != nil {
		originDutyLocation = &originDSTransportInfo.Name
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
	var providesGovernmentCounseling bool
	if orders.OriginDutyLocation != nil {
		providesGovernmentCounseling = orders.OriginDutyLocation.ProvidesServicesCounseling
	}

	if serviceMember.PersonalEmail == nil {
		return emails, fmt.Errorf("no email found for service member")
	}

	htmlBody, textBody, err := m.renderTemplates(appCtx, moveIssuedToPrimeEmailData{
		MilitaryOneSourceLink:        OneSourceTransportationOfficeLink,
		OriginDutyLocation:           originDutyLocation,
		DestinationLocation:          destinationAddress,
		ProvidesGovernmentCounseling: providesGovernmentCounseling,
		Locator:                      move.Locator,
	})

	if err != nil {
		appCtx.Logger().Error("error rendering template", zap.Error(err))
	}

	smEmail := emailContent{
		recipientEmail: *serviceMember.PersonalEmail,
		subject:        "Your personal property move has been ordered.",
		htmlBody:       htmlBody,
		textBody:       textBody,
	}

	appCtx.Logger().Info("Generated move issued to prime email",
		zap.String("moveLocator", move.Locator))

	// TODO: Send email to trusted contacts when that's supported
	return append(emails, smEmail), nil
}

func (m MoveIssuedToPrime) renderTemplates(appCtx appcontext.AppContext, data moveIssuedToPrimeEmailData) (string, string, error) {
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

type moveIssuedToPrimeEmailData struct {
	MilitaryOneSourceLink        string
	OriginDutyLocation           *string
	DestinationLocation          string
	ProvidesGovernmentCounseling bool
	Locator                      string
}

// RenderHTML renders the html for the email
func (m MoveIssuedToPrime) RenderHTML(appCtx appcontext.AppContext, data moveIssuedToPrimeEmailData) (string, error) {
	var htmlBuffer bytes.Buffer
	if err := m.htmlTemplate.Execute(&htmlBuffer, data); err != nil {
		appCtx.Logger().Error("cant render html template ", zap.Error(err))
	}
	return htmlBuffer.String(), nil
}

// RenderText renders the text for the email
func (m MoveIssuedToPrime) RenderText(appCtx appcontext.AppContext, data moveIssuedToPrimeEmailData) (string, error) {
	var textBuffer bytes.Buffer
	if err := m.textTemplate.Execute(&textBuffer, data); err != nil {
		appCtx.Logger().Error("cant render text template ", zap.Error(err))
		return "", err
	}
	return textBuffer.String(), nil
}
