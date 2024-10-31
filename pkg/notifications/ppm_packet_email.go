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
	ppmPacketEmailHTMLTemplate = html.Must(html.New("text_template").Parse(ppmPacketEmailRawHTML))
)

func formatMessageConfirmation(locator, pickupCity, pickupState, destinationCity, destinationState string) string {
	return fmt.Sprintf("This is a confirmation that your Personally Procured Move (PPM) with the assigned move code %s from %s, %s to %s, %s has been processed in MilMove.",
		locator,
		pickupCity,
		pickupState,
		destinationCity,
		destinationState,
	)
}

var templateShouldBeHTML = true

func formatMessageOpeningBranchCondition(branch, submitLocation string) string {
	switch branch {
	case GetAffiliationDisplayValues()[models.AffiliationMARINES],
		GetAffiliationDisplayValues()[models.AffiliationCOASTGUARD],
		GetAffiliationDisplayValues()[models.AffiliationNAVY]:
		if !templateShouldBeHTML {
			return "For Marine Corps, Navy, and Coast Guard personnel:\n\nYou can now log into MilMove " + MyMoveLink +
				"and view your payment packet; however, you do not need to forward your payment packet to finance as your closeout " +
				"location is associated with your finance office and they will handle this step for you.\n\nNote: Not all claimed " +
				"expenses may have been accepted during PPM closeout if they did not meet the definition of a valid expense.\n\n" +
				"Please be advised, your local finance office may require a DD Form 1351-2 to process payment. You can obtain a copy " +
				"of this form by utilizing the search feature at " + WashingtonHQServicesLink + "."
		}

		return "For Marine Corps, Navy, and Coast Guard personnel:<br><br>You can now log into MilMove " + MyMoveLink +
			"and view your payment packet; however, you do not need to forward your payment packet to finance as your closeout " +
			"location is associated with your finance office and they will handle this step for you.<br><br>Note: Not all claimed " +
			"expenses may have been accepted during PPM closeout if they did not meet the definition of a valid expense.<br><br>" +
			"Please be advised, your local finance office may require a DD Form 1351-2 to process payment. You can obtain a copy " +
			"of this form by utilizing the search feature at " + WashingtonHQServicesLink + "."
	case GetAffiliationDisplayValues()[models.AffiliationARMY]:
		if !templateShouldBeHTML {
			return "For Army personnel (FURTHER ACTION REQUIRED):\n\nLog in to " +
				"SmartVoucher at " + SmartVoucherLink + " using your CAC or myPay username and password.\n\nThis will allow you to edit " +
				"your voucher, and complete and sign DD Form 1351-2."
		}

		return "For Army personnel (FURTHER ACTION REQUIRED):<br><br>Log in to " +
			"SmartVoucher at " + SmartVoucherLink + " using your CAC or myPay username and password. This will allow you to edit " +
			"your voucher, and complete and sign DD Form 1351-2."
	case GetAffiliationDisplayValues()[models.AffiliationAIRFORCE], GetAffiliationDisplayValues()[models.AffiliationSPACEFORCE]:
		if !templateShouldBeHTML {
			return "For Air Force and Space Force personnel (FURTHER ACTION REQUIRED):\n\nYou can now log into MilMove <" + MyMoveLink +
				"> and download your payment packet to submit to " + submitLocation + ". You must complete this step to receive final " +
				"settlement of your PPM.\n\nNote: The Transportation Office does not determine claimable expenses. Claimable expenses " +
				"will be determined by finance.\n\nPlease be advised, your local finance office may require a DD Form 1351-2 " +
				"to process payment. You can obtain a copy of this form by utilizing the search feature at " + WashingtonHQServicesLink + "."
		}

		return "For Air Force and Space Force personnel (FURTHER ACTION REQUIRED):<br><br>You can now log into MilMove " + MyMoveLink +
			" and download your payment packet to submit to " + submitLocation + ". You must complete this step to receive final " +
			"settlement of your PPM.<br><br>Note: The Transportation Office does not determine claimable expenses. Claimable expenses " +
			"will be determined by finance.<br><br>Please be advised, your local finance office may require a DD Form 1351-2 " +
			"to process payment. You can obtain a copy of this form by utilizing the search feature at " + WashingtonHQServicesLink + "."
	}

	return ""
}

func formatMessageClosing() string {
	if !templateShouldBeHTML {
		return "If you have any questions, contact a government transportation office. You can see a listing of transportation " +
			"offices on Military One Source here: " + OneSourceTransportationOfficeLink + "\n\n" +
			"Thank you,\n\nUSTRANSCOM MilMove Team\n\nThe information contained in this email may contain Privacy Act information " +
			"and is therefore protected under the Privacy Act of 1974.  Failure to protect Privacy Act information could result in a $5,000 fine."
	}

	return "If you have any questions, contact a government transportation office. You can see a listing of transportation " +
		"offices on Military One Source here: " + OneSourceTransportationOfficeLink + "<br><br>" +
		"Thank you,<br><br>USTRANSCOM MilMove Team<br><br>The information contained in this email may contain Privacy Act information " +
		"and is therefore protected under the Privacy Act of 1974.  Failure to protect Privacy Act information could result in a $5,000 fine."
}

const (
	MessageDoNotReply                 = "*** DO NOT REPLY directly to this email ***"
	MessageActualExpenseReimbursement = "Please Note: Your PPM has been designated as Actual Expense Reimbursement. " +
		"This is the standard entitlement for Civilian employees. For uniformed Service Members, your PPM may have been " +
		"designated as Actual Expense Reimbursement due to failure to receive authorization prior to movement or failure " +
		"to obtain certified weight tickets. Actual Expense Reimbursement means reimbursement for expenses not to exceed " +
		"the Government Constructed Cost (GCC)."
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
	IsActualExpenseReimbursement      string
	OneSourceTransportationOfficeLink string
	WashingtonHQServicesLink          string
	MyMoveLink                        string
	SmartVoucherLink                  string
	MessageOpening                    string
	MessageConfirmation               string
	MessageBranchCondition            string
	MessageClosing                    string
	MessageAER                        string
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

func (p PpmPacketEmail) ConvertBoolToString(b *bool) string {
	if b != nil && *b {
		return "true"
	}
	return "false"
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
				ServiceBranch:                     GetAffiliationDisplayValues()[*serviceMember.Affiliation],
				Locator:                           move.Locator,
				IsActualExpenseReimbursement:      p.ConvertBoolToString(ppmShipment.IsActualExpenseReimbursement),
				OneSourceTransportationOfficeLink: OneSourceTransportationOfficeLink,
				WashingtonHQServicesLink:          WashingtonHQServicesLink,
				MyMoveLink:                        MyMoveLink,
				SmartVoucherLink:                  SmartVoucherLink,
				MessageOpening:                    MessageDoNotReply,
				MessageConfirmation:               formatMessageConfirmation(move.Locator, pickupAddress.City, pickupAddress.State, destinationAddress.City, destinationAddress.State),
				MessageBranchCondition:            formatMessageOpeningBranchCondition(GetAffiliationDisplayValues()[*serviceMember.Affiliation], submitLocation),
				MessageClosing:                    formatMessageClosing(),
				MessageAER:                        MessageActualExpenseReimbursement,
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
			ServiceBranch:                     GetAffiliationDisplayValues()[*serviceMember.Affiliation],
			Locator:                           move.Locator,
			IsActualExpenseReimbursement:      p.ConvertBoolToString(ppmShipment.IsActualExpenseReimbursement),
			OneSourceTransportationOfficeLink: OneSourceTransportationOfficeLink,
			WashingtonHQServicesLink:          WashingtonHQServicesLink,
			MyMoveLink:                        MyMoveLink,
			SmartVoucherLink:                  SmartVoucherLink,
			MessageOpening:                    MessageDoNotReply,
			MessageConfirmation:               formatMessageConfirmation(move.Locator, ppmShipment.PickupAddress.City, ppmShipment.PickupAddress.State, ppmShipment.DestinationAddress.City, ppmShipment.DestinationAddress.State),
			MessageBranchCondition:            formatMessageOpeningBranchCondition(GetAffiliationDisplayValues()[*serviceMember.Affiliation], submitLocation),
			MessageClosing:                    formatMessageClosing(),
			MessageAER:                        MessageActualExpenseReimbursement,
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
	templateShouldBeHTML = true
	if err := p.htmlTemplate.Execute(&htmlBuffer, data); err != nil {
		appCtx.Logger().Error("cant render html template ", zap.Error(err))
	}
	return htmlBuffer.String(), nil
}

// RenderText renders the text for the email
func (p PpmPacketEmail) RenderText(appCtx appcontext.AppContext, data PpmPacketEmailData) (string, error) {
	var textBuffer bytes.Buffer
	templateShouldBeHTML = false
	if err := p.textTemplate.Execute(&textBuffer, data); err != nil {
		appCtx.Logger().Error("cant render text template ", zap.Error(err))
		return "", err
	}
	return textBuffer.String(), nil
}
