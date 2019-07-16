package notifications

import (
	"context"
	"fmt"
	"net/url"

	"github.com/gobuffalo/pop"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/models"
)

// MoveApproved has notification content for approved moves
type MoveApproved struct {
	db      *pop.Connection
	logger  Logger
	host    string
	moveID  uuid.UUID
	session *auth.Session // TODO - remove this when we move permissions up to handlers and out of models
}

// NewMoveApproved returns a new move approval notification
func NewMoveApproved(db *pop.Connection,
	logger Logger,
	session *auth.Session,
	host string,
	moveID uuid.UUID) *MoveApproved {

	return &MoveApproved{
		db:      db,
		logger:  logger,
		host:    host,
		moveID:  moveID,
		session: session,
	}
}

func (m MoveApproved) emails(ctx context.Context) ([]emailContent, error) {
	var emails []emailContent

	move, err := models.FetchMove(m.db, m.session, m.moveID)
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

	dsTransportInfo, err := models.FetchDSContactInfo(m.db, serviceMember.DutyStationID)
	if err != nil {
		return emails, err
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

	introText := `You're all set to move!`
	dutyStationText := fmt.Sprintf("The local transportation office <strong>approved your move</strong> from <strong>%s</strong> to <strong>%s</strong>.", dsTransportInfo.Name, orders.NewDutyStation.Name)
	ppmInfoSheetInstructions := fmt.Sprintf("Please <a href=\"%s\">review the Personally Procured Move (PPM) info sheet</a> for detailed instructions.", ppmInfoSheetURL.String())

	if move.PersonallyProcuredMoves != nil {
		introText = fmt.Sprintf("<strong>%s</strong><br /><br /> %s <br /><br />%s<br /><br />",
			introText,
			dutyStationText,
			ppmInfoSheetInstructions,
		)
	}

	nextStepsText := `<strong>Next steps</strong>`

	ppmText := ""
	if move.PersonallyProcuredMoves != nil {
		ppmText = `Because you’ve chosen a do-it-yourself move, you can start whenever you are ready.<br /><br >
		Be sure to <strong>save your weight tickets and any receipts</strong> associated with your move. You’ll need them to request payment later in the process.`
	}

	// TODO: Add the PPPO contact info
	closingText := fmt.Sprintf("If you have any questions, call the <strong>%s</strong> PPPO at %s.<br /><br />You can <a href=\"%s\">check the status of your move</a> anytime at https://my.move.mil", dsTransportInfo.Name, dsTransportInfo.PhoneLine, "https://my.move.mil")

	smEmail := emailContent{
		recipientEmail: *serviceMember.PersonalEmail,
		subject:        "[MilMove] Your move is approved",
		htmlBody:       fmt.Sprintf("%s<br/>%s<br/>%s<br/><br />%s", introText, nextStepsText, ppmText, closingText),
		textBody:       fmt.Sprintf("%s\n%s\n%s\n%s", introText, nextStepsText, ppmText, closingText),
	}

	// TODO: Send email to trusted contacts when that's supported
	return append(emails, smEmail), nil
}
