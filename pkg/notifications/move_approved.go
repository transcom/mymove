package notifications

import (
	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/uuid"
	"github.com/transcom/mymove/pkg/models"
)

// MoveApproved has notification content for approved moves
type MoveApproved struct {
	db     *pop.Connection
	moveID uuid.UUID
	reqApp string
	user   models.User
}

func (m MoveApproved) emails() ([]emailContent, error) {
	var emails []emailContent

	move, err := models.FetchMove(m.db, m.user, m.reqApp, m.moveID)
	if err != nil {
		return emails, err
	}

	orders, err := models.FetchOrder(m.db, m.user, m.reqApp, move.OrdersID)
	if err != nil {
		return emails, err
	}

	serviceMember, err := models.FetchServiceMember(m.db, m.user, m.reqApp, orders.ServiceMemberID)
	if err != nil {
		return emails, err
	}

	smEmail := emailContent{
		recipientEmail: *serviceMember.PersonalEmail,
		subject:        "Move Approved",
		htmlBody:       "Congrats!<br>Your move has been approved!",
		textBody:       "Congrats! Your move has been approved!",
	}

	// TODO: Send email to trusted contacts when that's supported
	return append(emails, smEmail), nil
}
