package notifications

import (
	"fmt"
	"log"
	"time"

	"github.com/gobuffalo/pop"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

type EmailInfo struct {
	Email              string `db:"personal_email"`
	DutyStationName    string `db:"name"`
	NewDutyStationName string `db:"name"`
}
type EmailInfos []EmailInfo

func getEmailInfo(err error, db *pop.Connection, begRange time.Time, endRange time.Time) (*EmailInfos, error) {
	query := `SELECT sm.personal_email, dsn.name, dso.name
	FROM personally_procured_moves
	         JOIN moves m ON personally_procured_moves.move_id = m.id
	         JOIN orders o ON m.orders_id = o.id
	         JOIN service_members sm ON o.service_member_id = sm.id
	         JOIN duty_stations dso ON sm.duty_station_id = dso.id
	         JOIN duty_stations dsn ON o.new_duty_station_id = dsn.id
	WHERE approve_date BETWEEN $1 AND $2;`

	emailInfo := &EmailInfos{}
	pop.Debug = true
	err = db.RawQuery(query, begRange, endRange).All(emailInfo)
	return emailInfo, err
}

func (suite *NotificationSuite) TestMoveReviewedFetch() {
	db := suite.DB()
	//query := models.DB.LeftJoin("roles", "roles.id=user_roles.role_id").
	//	LeftJoin("users u", "u.id=user_roles.user_id").
	//	Where(`roles.name like ?`, name).Paginate(page, perpage)
	ppm1 := testdatagen.MakeDefaultPPM(db)
	d1 := time.Date(2019, 1, 1, 0, 0, 0, 0, time.UTC)
	err := ppm1.Submit(d1)
	suite.Nil(err)
	//TODO Update to use correct date that does not yet exist
	err = ppm1.Approve(d1)
	suite.Nil(err)
	verrs, err := models.SavePersonallyProcuredMove(db, &ppm1)
	suite.NoVerrs(verrs)
	suite.Nil(err)
	//TODO add a couple moves, create helper to setup test. Will need to be very careful with timezones
	//TODO maybe even create a table sent email to, or check this via AWS, since don't want to spam people?
	endRange := time.Date(2019, 1, 7, 0, 0, 0, 0, time.UTC)
	fmt.Println(endRange.Format(time.RFC3339))
	begRange := endRange.AddDate(0, 0, -7)
	log.Print(endRange)
	emailInfo, err := getEmailInfo(err, db, begRange, endRange)
	suite.NoError(err)
	log.Fatal(emailInfo)
	suite.Len(emailInfo, 2)
}
