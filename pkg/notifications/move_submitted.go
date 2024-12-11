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
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/models"
)

var (
	moveSubmittedRawTextTemplate = string(assets.MustAsset("notifications/templates/move_submitted_template.txt"))
	moveSubmittedTextTemplate    = text.Must(text.New("text_template").Parse(moveSubmittedRawTextTemplate))
	moveSubmittedRawHTMLTemplate = string(assets.MustAsset("notifications/templates/move_submitted_template.html"))
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

	destinationAddress := orders.NewDutyLocation.Name
	isSeparateeRetiree := orders.OrdersType == internalmessages.OrdersTypeRETIREMENT || orders.OrdersType == internalmessages.OrdersTypeSEPARATION
	if isSeparateeRetiree && len(move.MTOShipments) > 0 && move.MTOShipments[0].DestinationAddress != nil {
		destAddr := *move.MTOShipments[0].DestinationAddress
		optionalStreetAddress2 := ""
		if destAddr.StreetAddress2 != nil {
			optionalStreetAddress2 = " " + *destAddr.StreetAddress2
		}
		optionalStreetAddress3 := ""
		if destAddr.StreetAddress3 != nil {
			optionalStreetAddress3 = " " + *destAddr.StreetAddress3
		}
		destinationAddress = fmt.Sprintf("%s%s%s, %s, %s %s", destAddr.StreetAddress1, optionalStreetAddress2, optionalStreetAddress3, destAddr.City, destAddr.State, destAddr.PostalCode)
	}

	originDutyLocation := orders.OriginDutyLocation
	providesGovernmentCounseling := false
	if originDutyLocation != nil {
		providesGovernmentCounseling = originDutyLocation.ProvidesServicesCounseling
	}

	serviceMember, err := models.FetchServiceMemberForUser(appCtx.DB(), appCtx.Session(), orders.ServiceMemberID)
	if err != nil {
		return emails, err
	}

	originDSTransportInfo, err := models.FetchDLContactInfo(appCtx.DB(), orders.OriginDutyLocationID)
	if err != nil {
		return emails, err
	}

	var originDutyLocationName, originDutyLocationPhoneLine *string
	if originDSTransportInfo != nil {
		originDutyLocationName = &originDSTransportInfo.Name
		originDutyLocationPhoneLine = &originDSTransportInfo.PhoneLine
	} else if originDutyLocation != nil {
		originDutyLocationName = &originDutyLocation.Name
	}

	totalEntitlement := models.GetWeightAllotment(*orders.Grade)
	unaccompaniedBaggageAllowance, err := models.GetUBWeightAllowance(appCtx, originDutyLocation.Address.IsOconus, orders.NewDutyLocation.Address.IsOconus, orders.ServiceMember.Affiliation, orders.Grade, &orders.OrdersType, orders.Entitlement.DependentsAuthorized, orders.Entitlement.AccompaniedTour, orders.Entitlement.DependentsUnderTwelve, orders.Entitlement.DependentsTwelveAndOver)
	if err == nil {
		totalEntitlement.UnaccompaniedBaggageAllowance = unaccompaniedBaggageAllowance
	}

	weight := totalEntitlement.TotalWeightSelf
	if orders.HasDependents {
		weight = totalEntitlement.TotalWeightSelfPlusDependents
	}

	if serviceMember.PersonalEmail == nil {
		return emails, fmt.Errorf("no email found for service member")
	}

	htmlBody, textBody, err := m.renderTemplates(appCtx, moveSubmittedEmailData{
		OriginDutyLocation:                originDutyLocationName,
		DestinationLocation:               destinationAddress,
		OriginDutyLocationPhoneLine:       originDutyLocationPhoneLine,
		Locator:                           move.Locator,
		WeightAllowance:                   humanize.Comma(int64(weight)),
		ProvidesGovernmentCounseling:      providesGovernmentCounseling,
		OneSourceTransportationOfficeLink: OneSourceTransportationOfficeLink,
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
	OriginDutyLocation                *string
	DestinationLocation               string
	OriginDutyLocationPhoneLine       *string
	Locator                           string
	WeightAllowance                   string
	ProvidesGovernmentCounseling      bool
	OneSourceTransportationOfficeLink string
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
