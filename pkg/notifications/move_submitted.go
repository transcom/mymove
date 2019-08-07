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
	processText := "This can take up to 3 business days. The office will email you once your move has been approved."
	pppoText := "If you have questions or need expedited processing contact your local transportation office."
	closingText := "You can check the status of your move at any time at https://my.move.mil/"
	if serviceMember.DutyStationID != nil {
		originDSTransportInfo, err := models.FetchDSContactInfo(m.db, serviceMember.DutyStationID)
		if err != nil {
			return emails, err
		}
		destinationDutyStation, err := models.FetchDutyStation(context.Background(), m.db, orders.NewDutyStationID)
		if err != nil {
			return emails, err
		}

		submittedText = fmt.Sprintf(
			"Your move from %s to %s has been submitted to your local transportation office for review.",
			originDSTransportInfo.Name,
			destinationDutyStation.Name,
		)

		pppoText = fmt.Sprintf(
			"In the meantime, if you have questions or need expedited processing, call the %s PPPO at %s.",
			originDSTransportInfo.Name,
			originDSTransportInfo.PhoneLine,
		)
	}

	smEmail := emailContent{
		recipientEmail: *serviceMember.PersonalEmail,
		subject:        "[MilMove] You’ve submitted your move details",
		htmlBody:       fmt.Sprintf("%s<br/><br/>%s<br/><br/>%s<br/><br/><br/>%s", submittedText, processText, pppoText, closingText),
		textBody:       fmt.Sprintf("%s\n\n%s\n\n%s\n\n\n%s", submittedText, processText, pppoText, closingText),
	}

	m.logger.Info("Generated move submitted email to service member",
		zap.String("service member email address", *serviceMember.PersonalEmail))

	// TODO: Send email to trusted contacts when that's supported
	return append(emails, smEmail), nil
}
