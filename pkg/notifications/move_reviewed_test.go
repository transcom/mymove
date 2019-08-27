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
	ppms := suite.createPPMMoves(moves)

	moveReviewed := NewMoveReviewed(db, suite.logger, onDate)
	emailInfo, err := moveReviewed.GetEmailInfo(onDate)

	suite.NoError(err)
	suite.NotNil(emailInfo)
	suite.Len(emailInfo, 1)
	suite.Equal(ppms[0].Move.Orders.NewDutyStation.Name, emailInfo[0].NewDutyStationName)
	suite.NotNil(emailInfo[0].Email)
	suite.Equal(*ppms[0].Move.Orders.ServiceMember.PersonalEmail, *emailInfo[0].Email)
	suite.Equal(ppms[0].Move.Orders.ServiceMember.DutyStation.Name, emailInfo[0].DutyStationName)
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

	moveReviewed := NewMoveReviewed(db, suite.logger, startDate)
	emailInfo, err := moveReviewed.GetEmailInfo(startDate)

	suite.NoError(err)
	suite.Len(emailInfo, 0)
}

func (suite *NotificationSuite) TestHTMLTemplateRender() {
	startDate := time.Date(2019, 1, 7, 0, 0, 0, 0, time.UTC)
	onDate := startDate.AddDate(0, 0, -6)
	mr := NewMoveReviewed(suite.DB(), suite.logger, onDate)
	s := moveReviewedEmailData{
		Link:                   "www.survey",
		OriginDutyStation:      "OriginDutyStation",
		DestinationDutyStation: "DestDutyStation",
		Email:                  "email",
	}
	expectedHTMLContent := `<em>Good news:</em> Your move from OriginDutyStation to DestDutyStation has been processed for payment.

Can we ask a quick favor? <a href="www.survey"> Tell us about your experience</a> with requesting and receiving payment.

We’ll use your feedback to make MilMove better for your fellow service members.

Thank you for your thoughts, and <em>congratulations on your move.</em>`

	htmlContent := mr.RenderHTML(s)

	suite.Equal(expectedHTMLContent, htmlContent)

}

func (suite *NotificationSuite) TestTextTemplateRender() {
	startDate := time.Date(2019, 1, 7, 0, 0, 0, 0, time.UTC)
	onDate := startDate.AddDate(0, 0, -6)
	mr := NewMoveReviewed(suite.DB(), suite.logger, onDate)
	s := moveReviewedEmailData{
		Link:                   "www.survey",
		OriginDutyStation:      "OriginDutyStation",
		DestinationDutyStation: "DestDutyStation",
		Email:                  "email",
	}
	expectedTextContent := `Good news: Your move from OriginDutyStation to DestDutyStation has been processed for payment.

Can we ask a quick favor? Tell us about your experience with requesting and receiving payment at www.survey.

We’ll use your feedback to make MilMove better for your fellow service members.

Thank you for your thoughts, and congratulations on your move.`

	textContent := mr.RenderText(s)

	suite.Equal(expectedTextContent, textContent)
}

func (suite *NotificationSuite) TestFormatEmails() {
	startDate := time.Date(2019, 1, 7, 0, 0, 0, 0, time.UTC)
	onDate := startDate.AddDate(0, 0, -6)
	mr := NewMoveReviewed(suite.DB(), suite.logger, onDate)
	email1 := "email1"
	email2 := "email2"
	emailInfos := EmailInfos{
		{
			Email:              &email1,
			DutyStationName:    "d1",
			NewDutyStationName: "nd2",
		},
		{
			Email:              &email2,
			DutyStationName:    "d2",
			NewDutyStationName: "nd2",
		},
		{
			// nil emails should be skipped
			Email:              nil,
			DutyStationName:    "d2",
			NewDutyStationName: "nd2",
		},
	}

	formattedEmails, err := mr.formatEmails(emailInfos)

	suite.NoError(err)
	for i, actualEmailContent := range formattedEmails {
		emailInfo := emailInfos[i]
		data := moveReviewedEmailData{
			Link:                   surveyLink,
			OriginDutyStation:      emailInfo.DutyStationName,
			DestinationDutyStation: emailInfo.NewDutyStationName,
			Email:                  *emailInfo.Email,
		}
		expectedEmailContent := emailContent{
			recipientEmail: *emailInfo.Email,
			subject:        "[MilMove] Let us know how we did",
			htmlBody:       mr.RenderHTML(data),
			textBody:       mr.RenderText(data),
		}
		if emailInfo.Email != nil {
			suite.Equal(expectedEmailContent, actualEmailContent)
		}
	}
	// only expect the two moves with non-nil email addresses to get added to formattedEmails
	suite.Len(formattedEmails, 2)
}
