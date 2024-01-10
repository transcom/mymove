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

// PpmPacketEmail has notification content for approved moves
type PpmPacketEmail struct {
	ppmShipmentID uuid.UUID
	htmlTemplate  *html.Template
	textTemplate  *text.Template
}

// ppmPacketEmailData is used to render an email template
type ppmPacketEmailData struct {
	OriginZIP      string
	DestinationZIP string
	SubmitLocation string
	NavalBranch    bool
	ServiceBranch  string
	Locator        string
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
func (m PpmPacketEmail) emails(appCtx appcontext.AppContext) ([]emailContent, error) {
	var emails []emailContent

	appCtx.Logger().Info("ppm SHIPMENT UUID",
		zap.String("uuid", m.ppmShipmentID.String()),
	)

	var ppmShipment models.PPMShipment
	err := appCtx.DB().Find(&ppmShipment, m.ppmShipmentID)
	if err != nil {
		return emails, err
	} else if ppmShipment.PickupPostalCode == "" || ppmShipment.DestinationPostalCode == "" {
		return emails, fmt.Errorf("no pickup or destination postal code found for this shipment")
	}

	var mtoShipment models.MTOShipment
	err = appCtx.DB().Find(&mtoShipment, ppmShipment.ShipmentID)
	if err != nil {
		return emails, err
	}

	var move models.Move
	err = appCtx.DB().Find(&move, mtoShipment.MoveTaskOrderID)
	if err != nil {
		return emails, err
	}

	serviceMember, err := models.GetCustomerFromShipment(appCtx.DB(), ppmShipment.ShipmentID)
	if err != nil {
		return emails, err
	}

	if serviceMember.PersonalEmail == nil {
		return emails, fmt.Errorf("no email found for service member")
	}

	appCtx.Logger().Info("generated PPM Closeout Packet email for service member",
		zap.String("service member uuid", serviceMember.ID.String()),
		zap.String("PPM Shipment ID", ppmShipment.ID.String()),
		zap.String("Move Locator", move.Locator),
	)

	var submitLocation string
	if *serviceMember.Affiliation == models.AffiliationARMY {
		submitLocation = `the Defense Finance and Accounting Service (DFAS)`
	} else {
		submitLocation = `your local finance office`
	}

	var navalBranch bool
	if *serviceMember.Affiliation == models.AffiliationNAVY || *serviceMember.Affiliation == models.AffiliationMARINES ||
		*serviceMember.Affiliation == models.AffiliationCOASTGUARD {
		navalBranch = true
	} else {
		navalBranch = false
	}

	var affiliationDisplayValue = map[models.ServiceMemberAffiliation]string{
		models.AffiliationARMY:       "Army",
		models.AffiliationNAVY:       "Marine Corps, Navy and Coast Guard",
		models.AffiliationMARINES:    "Marine Corps, Navy and Coast Guard",
		models.AffiliationAIRFORCE:   "Air Force and Space Force",
		models.AffiliationSPACEFORCE: "Air Force and Space Force",
		models.AffiliationCOASTGUARD: "Marine Corps, Navy and Coast Guard",
	}

	htmlBody, textBody, err := m.renderTemplates(appCtx, ppmPacketEmailData{
		OriginZIP:      ppmShipment.PickupPostalCode,
		DestinationZIP: ppmShipment.DestinationPostalCode,
		SubmitLocation: submitLocation,
		NavalBranch:    navalBranch,
		ServiceBranch:  affiliationDisplayValue[*serviceMember.Affiliation],
		Locator:        move.Locator,
	})

	if err != nil {
		appCtx.Logger().Error("error rendering template", zap.Error(err))
	}

	ppmEmail := emailContent{
		recipientEmail: *serviceMember.PersonalEmail,
		subject:        "Your Personally Procured Move (PPM) closeout has been processed and is now available for your review.",
		htmlBody:       htmlBody,
		textBody:       textBody,
	}

	return append(emails, ppmEmail), nil
}

func (m PpmPacketEmail) renderTemplates(appCtx appcontext.AppContext, data ppmPacketEmailData) (string, string, error) {
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

// RenderHTML renders the html for the email
func (m PpmPacketEmail) RenderHTML(appCtx appcontext.AppContext, data ppmPacketEmailData) (string, error) {
	var htmlBuffer bytes.Buffer
	if err := m.htmlTemplate.Execute(&htmlBuffer, data); err != nil {
		appCtx.Logger().Error("cant render html template ", zap.Error(err))
	}
	return htmlBuffer.String(), nil
}

// RenderText renders the text for the email
func (m PpmPacketEmail) RenderText(appCtx appcontext.AppContext, data ppmPacketEmailData) (string, error) {
	var textBuffer bytes.Buffer
	if err := m.textTemplate.Execute(&textBuffer, data); err != nil {
		appCtx.Logger().Error("cant render text template ", zap.Error(err))
		return "", err
	}
	return textBuffer.String(), nil
}
