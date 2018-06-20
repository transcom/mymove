package notifications

import (
	"fmt"
	// "net/url"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/uuid"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/models"
)

// MoveCanceled has notification content for approved moves
type MoveCanceled struct {
	db      *pop.Connection
	logger  *zap.Logger
	moveID  uuid.UUID
	session *auth.Session // TODO - remove this when we move permissions up to handlers and out of models
}

// NewMoveCanceled returns a new move approval notification
func NewMoveCanceled(db *pop.Connection,
	logger *zap.Logger,
	session *auth.Session,
	moveID uuid.UUID) *MoveCanceled {

	return &MoveCanceled{
		db:      db,
		logger:  logger,
		moveID:  moveID,
		session: session,
	}
}

func (m MoveCanceled) emails() ([]emailContent, error) {
	var emails []emailContent

	move, err := models.FetchMove(m.db, m.session, m.moveID)
	if err != nil {
		return emails, err
	}

	orders, err := models.FetchOrder(m.db, m.session, move.OrdersID)
	if err != nil {
		return emails, err
	}

	serviceMember, err := models.FetchServiceMember(m.db, m.session, orders.ServiceMemberID)
	if err != nil {
		return emails, err
	}

	if serviceMember.PersonalEmail == nil {
		return emails, fmt.Errorf("no email found for service member")
	}

	// Set up various text segments. Copy comes from here:
	// https://docs.google.com/document/d/1bgE0Q_-_c93uruMP8dcNSHugXo8Pidz6YFojWBKn1Gg/edit#heading=h.h3ys1ur2qhpn
	// TODO: we will want some sort of templating system
	// TODO: there is currently (6/20) no cancel email text in the doc above, so what's here is placeholder/suggested

	introText := `Your move has been canceled.`
	if move.PersonallyProcuredMoves != nil {
		introText = fmt.Sprintf("%s", introText)
	}

	// TODO: Add the PPPO contact info
	closingText := `If you have any questions, contact your origin PPPO.`

	smEmail := emailContent{
		recipientEmail: *serviceMember.PersonalEmail,
		subject:        "MOVE.MIL: Your move has been canceled.",
		htmlBody:       fmt.Sprintf("%s<br/>%s", introText, closingText),
		textBody:       fmt.Sprintf("%s\n%s", introText, closingText),
	}

	m.logger.Info("Sent move cancellation email to service member",
		zap.String("sevice member email address", *serviceMember.PersonalEmail))

	// TODO: Send email to trusted contacts when that's supported
	return append(emails, smEmail), nil
}
