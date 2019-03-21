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
	moveID  uuid.UUID
	session *auth.Session // TODO - remove this when we move permissions up to handlers and out of models
}

// NewMoveApproved returns a new move approval notification
func NewMoveApproved(db *pop.Connection,
	logger Logger,
	session *auth.Session,
	moveID uuid.UUID) *MoveApproved {

	return &MoveApproved{
		db:      db,
		logger:  logger,
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

	if serviceMember.PersonalEmail == nil {
		return emails, fmt.Errorf("no email found for service member")
	}

	// Set up various text segments. Copy comes from here:
	// https://docs.google.com/document/d/1bgE0Q_-_c93uruMP8dcNSHugXo8Pidz6YFojWBKn1Gg/edit#heading=h.h3ys1ur2qhpn
	// TODO: we will want some sort of templating system

	ppmInfoSheetURL := url.URL{
		Scheme: "https",
		Host:   m.session.Hostname,
		Path:   "downloads/ppm_info_sheet.pdf",
	}

	introText := `Your move has been approved and you are ready to move!`
	if move.PersonallyProcuredMoves != nil {
		introText = fmt.Sprintf("%s %s %s",
			introText,
			`Please review the PPM info sheet for more detailed instructions: `,
			ppmInfoSheetURL.String())
	}

	nextStepsText := `Next steps:`

	ppmText := ""
	if move.PersonallyProcuredMoves != nil {
		ppmText = `For your “Do-it-Yourself” shipment, you can begin your move whenever you are ready. Be sure to save your weight tickets and any receipts associated with your move for when you request payment later on in the process.`
	}

	// TODO: Add the PPPO contact info
	closingText := `If you have any questions, contact your origin PPPO.`

	smEmail := emailContent{
		recipientEmail: *serviceMember.PersonalEmail,
		subject:        "MOVE.MIL: Your move has been approved.",
		htmlBody:       fmt.Sprintf("%s<br/>%s<br/>%s<br/>%s", introText, nextStepsText, ppmText, closingText),
		textBody:       fmt.Sprintf("%s\n%s\n%s\n%s", introText, nextStepsText, ppmText, closingText),
	}

	// TODO: Send email to trusted contacts when that's supported
	return append(emails, smEmail), nil
}
