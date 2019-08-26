package notifications

import (
	"log"
	"time"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *NotificationSuite) createPPMMoves(assertions []testdatagen.Assertions) []models.PersonallyProcuredMove {
	ppms := make([]models.PersonallyProcuredMove, 0)
	for _, assertion := range assertions {
		ppm := testdatagen.MakePPM(suite.DB(), assertion)
		ppms = append(ppms, ppm)
	}
	return ppms
}

func (suite *NotificationSuite) TestMoveReviewedFetchInRange() {
	db := suite.DB()
	endRange := time.Date(2019, 1, 7, 0, 0, 0, 0, time.UTC)
	begRange := endRange.AddDate(0, 0, -7)
	inRange := begRange.AddDate(0, 0, 1)
	outOfRange := endRange.AddDate(0, 0, 1)
	moves := []testdatagen.Assertions{
		{PersonallyProcuredMove: models.PersonallyProcuredMove{Status: models.PPMStatusAPPROVED, ApproveDate: &inRange}},
		{PersonallyProcuredMove: models.PersonallyProcuredMove{Status: models.PPMStatusAPPROVED, ApproveDate: &outOfRange}},
	}
	ppms := suite.createPPMMoves(moves)
	for _, v := range ppms {
		log.Println(*v.Move.Orders.ServiceMember.PersonalEmail)
	}
	moveReviewed := NewMoveReviewed(db, suite.logger)

	emailInfo, err := moveReviewed.GetEmailInfo(begRange, endRange)

	suite.NoError(err)
	suite.Len(*emailInfo, 1)
}

func (suite *NotificationSuite) TestMoveReviewedFetchNoneInRange() {
	db := suite.DB()
	endRange := time.Date(2019, 1, 7, 0, 0, 0, 0, time.UTC)
	begRange := endRange.AddDate(0, 0, -7)
	outOfRange1 := begRange.AddDate(0, 0, -1)
	outOfRange2 := endRange.AddDate(0, 0, 1)
	moves := []testdatagen.Assertions{
		{PersonallyProcuredMove: models.PersonallyProcuredMove{Status: models.PPMStatusAPPROVED, ApproveDate: &outOfRange1}},
		{PersonallyProcuredMove: models.PersonallyProcuredMove{Status: models.PPMStatusAPPROVED, ApproveDate: &outOfRange2}},
	}
	ppms := suite.createPPMMoves(moves)
	for _, v := range ppms {
		log.Println(*v.Move.Orders.ServiceMember.PersonalEmail)
	}
	moveReviewed := NewMoveReviewed(db, suite.logger)

	emailInfo, err := moveReviewed.GetEmailInfo(begRange, endRange)

	suite.NoError(err)
	suite.Len(*emailInfo, 0)
}

func (suite *NotificationSuite) TestMoveReviewedFetchAllInRange() {
	db := suite.DB()
	endRange := time.Date(2019, 1, 7, 0, 0, 0, 0, time.UTC)
	begRange := endRange.AddDate(0, 0, -7)
	inRange1 := begRange.AddDate(0, 0, 1)
	inRange2 := endRange.AddDate(0, 0, -1)
	moves := []testdatagen.Assertions{
		{PersonallyProcuredMove: models.PersonallyProcuredMove{Status: models.PPMStatusAPPROVED, ApproveDate: &inRange1}},
		{PersonallyProcuredMove: models.PersonallyProcuredMove{Status: models.PPMStatusAPPROVED, ApproveDate: &inRange2}},
	}
	ppms := suite.createPPMMoves(moves)
	for _, v := range ppms {
		log.Println(*v.Move.Orders.ServiceMember.PersonalEmail)
	}
	moveReviewed := NewMoveReviewed(db, suite.logger)

	emailInfo, err := moveReviewed.GetEmailInfo(begRange, endRange)

	suite.NoError(err)
	suite.Len(*emailInfo, 2)
}
