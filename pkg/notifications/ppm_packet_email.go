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
	ppmPacketEmailRawText      = string(assets.MustAsset("notifications/templates/ppm_packet_email_template.txt"))
	ppmPacketEmailTextTemplate = text.Must(text.New("text_template").Parse(ppmPacketEmailRawText))
	ppmPacketEmailRawHTML      = string(assets.MustAsset("notifications/templates/ppm_packet_email_template.html"))
	ppmPacketEmailHTMLTemplate = html.Must(html.New("html_template").Parse(ppmPacketEmailRawHTML))
)

// PpmPacketEmail has notification content for approved moves
type PpmPacketEmail struct {
	ppmShipmentID uuid.UUID
	htmlTemplate  *html.Template
	textTemplate  *text.Template
}

// ppmPacketEmailData is used to render an email template
// Uses ZIPs only if no city/state data is provided
type PpmPacketEmailData struct {
	OriginZIP                         *string
	OriginCity                        *string
	OriginState                       *string
	DestinationZIP                    *string
	DestinationCity                   *string
	DestinationState                  *string
	SubmitLocation                    string
	ServiceBranch                     string
	Locator                           string
	IsActualExpenseReimbursement      *bool
	OneSourceTransportationOfficeLink string
	WashingtonHQServicesLink          string
	MyMoveLink                        string
	SmartVoucherLink                  string
}

// Used to get logging data from GetEmailData
type LoggerData struct {
	ServiceMember models.ServiceMember
	PPMShipmentID uuid.UUID
	MoveLocator   string
}

// NewPpmPacketEmail returns a new payment reminder notification 14 days after actual move in date
func NewPpmPacketEmail(ppmShipmentID uuid.UUID) *PpmPacketEmail {

	return &PpmPacketEmail{
		ppmShipmentID: ppmShipmentID,
		htmlTemplate:  ppmPacketEmailHTMLTemplate,
		textTemplate:  ppmPacketEmailTextTemplate,
	}
}

// NotificationSendingContext expects a `notification` with an `emails` method,
// so we implement `email` to satisfy that interface
func (p PpmPacketEmail) emails(appCtx appcontext.AppContext) ([]emailContent, error) {
	var emails []emailContent

	appCtx.Logger().Info("ppm SHIPMENT UUID",
		zap.String("uuid", p.ppmShipmentID.String()),
	)

	emailData, loggerData, err := p.GetEmailData(appCtx)
	if err != nil {
		return nil, err
	}

	appCtx.Logger().Info("generated PPM Closeout Packet email",
		zap.String("service member uuid", loggerData.ServiceMember.ID.String()),
		zap.String("PPM Shipment ID", loggerData.PPMShipmentID.String()),
		zap.String("Move Locator", loggerData.MoveLocator),
	)

	var htmlBody, textBody string
	htmlBody, textBody, err = p.renderTemplates(appCtx, emailData)

	if err != nil {
		appCtx.Logger().Error("error rendering template", zap.Error(err))
	}

	ppmEmail := emailContent{
		recipientEmail: *loggerData.ServiceMember.PersonalEmail,
		subject:        "Your Personally Procured Move (PPM) closeout has been processed and is now available for your review.",
		htmlBody:       htmlBody,
		textBody:       textBody,
	}

	return append(emails, ppmEmail), nil
}

func (p PpmPacketEmail) GetEmailData(appCtx appcontext.AppContext) (PpmPacketEmailData, LoggerData, error) {
	var ppmShipment models.PPMShipment
	err := appCtx.DB().Find(&ppmShipment, p.ppmShipmentID)
	if err != nil {
		return PpmPacketEmailData{}, LoggerData{}, err
	} else if ppmShipment.PickupAddressID == nil || ppmShipment.DestinationAddressID == nil {
		return PpmPacketEmailData{}, LoggerData{}, fmt.Errorf("no pickup or destination address found for this shipment")
	}

	var mtoShipment models.MTOShipment
	err = appCtx.DB().Find(&mtoShipment, ppmShipment.ShipmentID)
	if err != nil {
		return PpmPacketEmailData{}, LoggerData{}, err
	}

	var move models.Move
	err = appCtx.DB().Find(&move, mtoShipment.MoveTaskOrderID)
	if err != nil {
		return PpmPacketEmailData{}, LoggerData{}, err
	}

	serviceMember, err := models.GetCustomerFromShipment(appCtx.DB(), ppmShipment.ShipmentID)
	if err != nil {
		return PpmPacketEmailData{}, LoggerData{}, err
	}

	if serviceMember.PersonalEmail == nil {
		return PpmPacketEmailData{}, LoggerData{}, fmt.Errorf("no email found for service member")
	}

	var submitLocation string
	if *serviceMember.Affiliation == models.AffiliationARMY {
		submitLocation = `the Defense Finance and Accounting Service (DFAS)`
	} else {
		submitLocation = `your local finance office`
	}

	var affiliationDisplayValue = map[models.ServiceMemberAffiliation]string{
		models.AffiliationARMY:       "Army",
		models.AffiliationNAVY:       "Marine Corps, Navy, and Coast Guard",
		models.AffiliationMARINES:    "Marine Corps, Navy, and Coast Guard",
		models.AffiliationAIRFORCE:   "Air Force and Space Force",
		models.AffiliationSPACEFORCE: "Air Force and Space Force",
		models.AffiliationCOASTGUARD: "Marine Corps, Navy, and Coast Guard",
	}

	// If address IDs are available for this PPM shipment, then do another query to get the city/state for origin and destination.
	// Note: This is a conditional put in because this work was done before address_ids were added to the ppm_shipments table.
	if ppmShipment.PickupAddressID != nil && ppmShipment.DestinationAddressID != nil {
		var pickupAddress, destinationAddress models.Address
		err = appCtx.DB().Find(&pickupAddress, ppmShipment.PickupAddressID)
		if err != nil {
			return PpmPacketEmailData{}, LoggerData{}, err
		}
		err = appCtx.DB().Find(&destinationAddress, ppmShipment.DestinationAddressID)
		if err != nil {
			return PpmPacketEmailData{}, LoggerData{}, err
		}

		return PpmPacketEmailData{
				OriginCity:                        &pickupAddress.City,
				OriginState:                       &pickupAddress.State,
				OriginZIP:                         &pickupAddress.PostalCode,
				DestinationCity:                   &destinationAddress.City,
				DestinationState:                  &destinationAddress.State,
				DestinationZIP:                    &destinationAddress.PostalCode,
				SubmitLocation:                    submitLocation,
				ServiceBranch:                     affiliationDisplayValue[*serviceMember.Affiliation],
				Locator:                           move.Locator,
				IsActualExpenseReimbursement:      ppmShipment.IsActualExpenseReimbursement,
				OneSourceTransportationOfficeLink: OneSourceTransportationOfficeLink,
				WashingtonHQServicesLink:          WashingtonHQServicesLink,
				MyMoveLink:                        MyMoveLink,
				SmartVoucherLink:                  SmartVoucherLink,
			},
			LoggerData{
				ServiceMember: *serviceMember,
				PPMShipmentID: ppmShipment.ID,
				MoveLocator:   move.Locator,
			}, nil
	}

	// Fallback to using ZIPs if the above if-block for city,state doesn't happen
	return PpmPacketEmailData{
			OriginZIP:                         &ppmShipment.PickupAddress.PostalCode,
			DestinationZIP:                    &ppmShipment.DestinationAddress.PostalCode,
			SubmitLocation:                    submitLocation,
			ServiceBranch:                     affiliationDisplayValue[*serviceMember.Affiliation],
			Locator:                           move.Locator,
			IsActualExpenseReimbursement:      ppmShipment.IsActualExpenseReimbursement,
			OneSourceTransportationOfficeLink: OneSourceTransportationOfficeLink,
			WashingtonHQServicesLink:          WashingtonHQServicesLink,
			MyMoveLink:                        MyMoveLink,
			SmartVoucherLink:                  SmartVoucherLink,
		},
		LoggerData{
			ServiceMember: *serviceMember,
			PPMShipmentID: ppmShipment.ID,
			MoveLocator:   move.Locator,
		}, nil

}

func (p PpmPacketEmail) renderTemplates(appCtx appcontext.AppContext, data PpmPacketEmailData) (string, string, error) {
	htmlBody, err := p.RenderHTML(appCtx, data)
	if err != nil {
		return "", "", fmt.Errorf("error rendering html template using %#v", data)
	}
	textBody, err := p.RenderText(appCtx, data)
	if err != nil {
		return "", "", fmt.Errorf("error rendering text template using %#v", data)
	}
	return htmlBody, textBody, nil
}

// RenderHTML renders the html for the email
func (p PpmPacketEmail) RenderHTML(appCtx appcontext.AppContext, data PpmPacketEmailData) (string, error) {
	var htmlBuffer bytes.Buffer
	if err := p.htmlTemplate.Execute(&htmlBuffer, data); err != nil {
		appCtx.Logger().Error("cant render html template ", zap.Error(err))
	}
	return htmlBuffer.String(), nil
}

// RenderText renders the text for the email
func (p PpmPacketEmail) RenderText(appCtx appcontext.AppContext, data PpmPacketEmailData) (string, error) {
	var textBuffer bytes.Buffer
	if err := p.textTemplate.Execute(&textBuffer, data); err != nil {
		appCtx.Logger().Error("cant render text template ", zap.Error(err))
		return "", err
	}
	return textBuffer.String(), nil
}
