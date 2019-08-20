package notifications

import (
	"context"
	"fmt"

	"github.com/gobuffalo/pop"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/models"
)

// MoveCanceled has notification content for approved moves
type MoveCanceled struct {
	db      *pop.Connection
	logger  Logger
	moveID  uuid.UUID
	session *auth.Session // TODO - remove this when we move permissions up to handlers and out of models
}

// NewMoveCanceled returns a new move approval notification
func NewMoveCanceled(db *pop.Connection, logger Logger, session *auth.Session, moveID uuid.UUID) *MoveCanceled {

	return &MoveCanceled{
		db:      db,
		logger:  logger,
		moveID:  moveID,
		session: session,
	}
}

func (m MoveCanceled) emails(ctx context.Context) ([]emailContent, error) {
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

	dsTransportInfo, err := models.FetchDSContactInfo(m.db, serviceMember.DutyStationID)
	if err != nil {
		return emails, err
	}

	if orders.NewDutyStation.Name == "" {
		return emails, fmt.Errorf("missing new duty station for service member")
	}

	// Set up various text segments. Copy comes from here:
	// https://docs.google.com/document/d/1gIQZprWzJJE_sAAyg5NViPwy9ckL5RK37gFq1fEfipU
	// TODO: we will want some sort of templating system

	introText := `Update on your move`
	nextSteps := fmt.Sprintf("Upon review, the office has determined that MilMove can’t handle your move from %s to %s, and canceled it in the system. You’re still moving, but we’ll need to manage your move with a different system.",
		dsTransportInfo.Name, orders.NewDutyStation.Name)
	closingText := fmt.Sprintf("Please call the PPPO at %s at %s and they’ll help you figure out what to do next.",
		dsTransportInfo.Name, dsTransportInfo.PhoneLine)

	smEmail := emailContent{
		recipientEmail: *serviceMember.PersonalEmail,
		subject:        fmt.Sprintf("[MilMove] %s", introText),
		htmlBody:       fmt.Sprintf("%s<br/><br/><br/>%s", nextSteps, closingText),
		textBody:       fmt.Sprintf("%s\n\n\n%s", nextSteps, closingText),
	}

	// TODO: Send email to trusted contacts when that's supported
	return append(emails, smEmail), nil
}
