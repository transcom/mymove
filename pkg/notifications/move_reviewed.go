package notifications

import (
	"errors"
	"fmt"
	"time"

	"github.com/gobuffalo/pop"
	"go.uber.org/zap"
)

// MoveReviewed has notification content for completed/reviewed moves
type MoveReviewed struct {
	db     *pop.Connection
	logger Logger
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

func (m MoveReviewed) GetEmailInfo(date time.Time) (*EmailInfos, error) {
	dateString := date.Format("2006-01-02")
	query := `SELECT sm.personal_email, dsn.name, dso.name
	FROM personally_procured_moves
	         JOIN moves m ON personally_procured_moves.move_id = m.id
	         JOIN orders o ON m.orders_id = o.id
	         JOIN service_members sm ON o.service_member_id = sm.id
	         JOIN duty_stations dso ON sm.duty_station_id = dso.id
	         JOIN duty_stations dsn ON o.new_duty_station_id = dsn.id
	WHERE CAST(reviewed_date as date) = $1;`

	emailInfo := &EmailInfos{}
	err := m.db.RawQuery(query, dateString).All(emailInfo)
	return emailInfo, err
}
func (m MoveReviewed) emails(date time.Time) ([]emailContent, error) {
	var emails []emailContent
	// now := time.Now()
	// then := now.AddDate(0, 0, offsetDays)
	emailInfos, err := m.GetEmailInfo(date)

	if err != nil {
		return nil, err
	}

	if emailInfos == nil {
		return nil, errors.New("TODO")
	}
	for _, emailInfo := range *emailInfos {
		// TODO email should not be nil but can be in db
		// if emailInfo.Email == nil {
		// 	return emails, fmt.Errorf("no email found for service member")
		// }

		link := "https://www.surveymonkey.com/r/MilMovePt3-08191"
		startText := fmt.Sprintf(
			"Good news: Your move from %s to %s has been processed for payment. ",
			emailInfo.DutyStationName,
			emailInfo.NewDutyStationName,
		)
		startTextHTML := fmt.Sprintf(
			"<strong>Good news:</strong> Your move from %s to %s has been processed for payment. ",
			emailInfo.DutyStationName,
			emailInfo.NewDutyStationName,
		)
		surveyText := fmt.Sprintf("Can we ask a quick favor? Tell us about your experience with requesting and receiving payment at %s.", link)
		surveyTextHTML := fmt.Sprintf("Can we ask a quick favor? <a href=\"%s\">Tell us about your experience</a> with requesting and receiving payment.", link)
		feedbackText := "We’ll use your feedback to make MilMove better for your fellow service members."
		closingText := "Thank you for your thoughts, and congratulations on your move."
		closingTextHTML := "Thank you for your thoughts, and <strong>congratulations on your move.</strong>"

		smEmail := emailContent{
			recipientEmail: emailInfo.Email,
			subject:        "[MilMove] Let us know how we did",
			htmlBody:       fmt.Sprintf("%s<br/><br/>%s<br/><br/>%s<br/><br/><br/>%s", startTextHTML, surveyTextHTML, feedbackText, closingTextHTML),
			textBody:       fmt.Sprintf("%s\n\n%s\n\n%s\n\n\n%s", startText, surveyText, feedbackText, closingText),
		}

		m.logger.Info("Generated move reviewed email to service member",
			zap.String("service member email address", emailInfo.Email))

		// TODO: Send email to trusted contacts when that's supported
		emails = append(emails, smEmail)
	}
	return emails, nil
}
