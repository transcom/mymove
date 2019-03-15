package notifications

import (
	"context"
	"fmt"

	"github.com/gobuffalo/pop"
	"github.com/gofrs/uuid"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/models"
)

// MoveSubmitted has notification content for submitted moves
type MoveSubmitted struct {
	db      *pop.Connection
	logger  Logger
	moveID  uuid.UUID
	session *auth.Session // TODO - remove this when we move permissions up to handlers and out of models
}

// NewMoveSubmitted returns a new move submitted notification
func NewMoveSubmitted(db *pop.Connection, logger Logger, session *auth.Session, moveID uuid.UUID) *MoveSubmitted {

	return &MoveSubmitted{
		db:      db,
		logger:  logger,
		moveID:  moveID,
		session: session,
	}
}

func (m MoveSubmitted) emails(ctx context.Context) ([]emailContent, error) {
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

	submittedText := "Your move has been submitted to your local transportation office for review. "
	processText := "This process can take up to 3 business days. "
	pppoText := "If you have questions or need expedited processing contact your local transportation office."

	if serviceMember.DutyStationID != nil {
		originDSTransportInfo, err := models.FetchDSContactInfo(m.db, serviceMember.DutyStationID)
		if err != nil {
			return emails, err
		}
		destinationDutyStation, err := models.FetchDutyStation(context.Background(), m.db, orders.NewDutyStationID)
		if err != nil {
			return emails, err
		}

		pppoText = fmt.Sprintf(
			"If you have questions or need expedited processing contact your local PPPO %s at %s.",
			originDSTransportInfo.Name,
			originDSTransportInfo.PhoneLine,
		)
		submittedText = fmt.Sprintf(
			"Your move from %s to %s has been submitted to your local transportation office for review. ",
			originDSTransportInfo.Name,
			destinationDutyStation.Name,
		)
	}

	text := submittedText + processText + pppoText
	smEmail := emailContent{
		recipientEmail: *serviceMember.PersonalEmail,
		subject:        "MOVE.MIL: Your move has been submitted.",
		htmlBody:       text,
		textBody:       text,
	}

	m.logger.Info("Generated move submitted email to service member",
		zap.String("service member email address", *serviceMember.PersonalEmail))

	// TODO: Send email to trusted contacts when that's supported
	return append(emails, smEmail), nil
}
