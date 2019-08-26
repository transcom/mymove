package notifications

import (
	"context"
	"fmt"
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gofrs/uuid"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/models"
)

// MoveReviewed has notification content for completed/reviewed moves
type MoveReviewed struct {
	db      *pop.Connection
	logger  Logger
	moveID  uuid.UUID
	session *auth.Session // TODO - remove this when we move permissions up to handlers and out of models
}

// NewMoveReviewed returns a new move submitted notification
func NewMoveReviewed(db *pop.Connection, logger Logger) *MoveReviewed {

	return &MoveReviewed{
		db:     db,
		logger: logger,
	}
}

type EmailInfos []EmailInfo

type EmailInfo struct {
	Email              string `db:"personal_email"`
	DutyStationName    string `db:"name"`
	NewDutyStationName string `db:"name"`
}

func (m MoveReviewed) GetEmailInfo(begRange time.Time, endRange time.Time) (*EmailInfos, error) {
	//TODO query below needs to be updated to use review date not approve date
	query := `SELECT sm.personal_email, dsn.name, dso.name
	FROM personally_procured_moves
	         JOIN moves m ON personally_procured_moves.move_id = m.id
	         JOIN orders o ON m.orders_id = o.id
	         JOIN service_members sm ON o.service_member_id = sm.id
	         JOIN duty_stations dso ON sm.duty_station_id = dso.id
	         JOIN duty_stations dsn ON o.new_duty_station_id = dsn.id
	WHERE approve_date BETWEEN $1 AND $2;`

	emailInfo := &EmailInfos{}
	err := m.db.RawQuery(query, begRange, endRange).All(emailInfo)
	return emailInfo, err
}

func (m MoveReviewed) emails(ctx context.Context) ([]emailContent, error) {
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

	link := "https://www.surveymonkey.com/r/MilMovePt3-08191"
	startText := "Good news: Your move has been processed for payment."
	startTextHTML := "<strong>Good news:</strong> Your move has been processed for payment."
	surveyText := fmt.Sprintf("Can we ask a quick favor? Tell us about your experience with requesting and receiving payment at %s.", link)
	surveyTextHTML := fmt.Sprintf("Can we ask a quick favor? <a href=\"%s\">Tell us about your experience</a> with requesting and receiving payment.", link)
	feedbackText := "We’ll use your feedback to make MilMove better for your fellow service members."
	closingText := "Thank you for your thoughts, and congratulations on your move."
	closingTextHTML := "Thank you for your thoughts, and <strong>congratulations on your move.</strong>"

	if serviceMember.DutyStationID != nil {
		originDSTransportInfo, err := models.FetchDSContactInfo(m.db, serviceMember.DutyStationID)
		if err != nil {
			return emails, err
		}
		destinationDutyStation, err := models.FetchDutyStation(context.Background(), m.db, orders.NewDutyStationID)
		if err != nil {
			return emails, err
		}

		startText = fmt.Sprintf(
			"Good news: Your move from %s to %s has been processed for payment. ",
			originDSTransportInfo.Name,
			destinationDutyStation.Name,
		)
		startTextHTML = fmt.Sprintf(
			"<strong>Good news:</strong> Your move from %s to %s has been processed for payment. ",
			originDSTransportInfo.Name,
			destinationDutyStation.Name,
		)
	}

	smEmail := emailContent{
		recipientEmail: *serviceMember.PersonalEmail,
		subject:        "[MilMove] Let us know how we did",
		htmlBody:       fmt.Sprintf("%s<br/><br/>%s<br/><br/>%s<br/><br/><br/>%s", startTextHTML, surveyTextHTML, feedbackText, closingTextHTML),
		textBody:       fmt.Sprintf("%s\n\n%s\n\n%s\n\n\n%s", startText, surveyText, feedbackText, closingText),
	}

	m.logger.Info("Generated move reviewed email to service member",
		zap.String("service member email address", *serviceMember.PersonalEmail))

	// TODO: Send email to trusted contacts when that's supported
	return append(emails, smEmail), nil
}
