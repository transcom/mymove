package notifications

import (
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

func (suite *NotificationSuite) TestMoveReviewedFetchSomeFound() {
	db := suite.DB()
	startDate := time.Date(2019, 1, 7, 0, 0, 0, 0, time.UTC)
	onDate := startDate.AddDate(0, 0, -6)
	offDate := startDate.AddDate(0, 0, -7)
	moves := []testdatagen.Assertions{
		{PersonallyProcuredMove: models.PersonallyProcuredMove{Status: models.PPMStatusAPPROVED, ReviewedDate: &onDate}},
		{PersonallyProcuredMove: models.PersonallyProcuredMove{Status: models.PPMStatusAPPROVED, ReviewedDate: &offDate}},
	}
	suite.createPPMMoves(moves)
	moveReviewed := NewMoveReviewed(db, suite.logger)

	emailInfo, err := moveReviewed.GetEmailInfo(onDate)

	suite.NoError(err)
	suite.Len(*emailInfo, 1)
}

func (suite *NotificationSuite) TestMoveReviewedFetchNoneFound() {
	db := suite.DB()
	startDate := time.Date(2019, 1, 7, 0, 0, 0, 0, time.UTC)
	offDate := startDate.AddDate(0, 0, -7)
	moves := []testdatagen.Assertions{
		{PersonallyProcuredMove: models.PersonallyProcuredMove{Status: models.PPMStatusAPPROVED, ReviewedDate: &offDate}},
		{PersonallyProcuredMove: models.PersonallyProcuredMove{Status: models.PPMStatusAPPROVED, ReviewedDate: &offDate}},
	}
	suite.createPPMMoves(moves)

	moveReviewed := NewMoveReviewed(db, suite.logger)

	emailInfo, err := moveReviewed.GetEmailInfo(startDate)

	suite.NoError(err)
	suite.Len(*emailInfo, 0)
}
