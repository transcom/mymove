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

	moveReviewed, err := NewMoveReviewed(db, suite.logger, onDate)
	suite.NoError(err)
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

	moveReviewed, err := NewMoveReviewed(db, suite.logger, startDate)
	suite.NoError(err)
	emailInfo, err := moveReviewed.GetEmailInfo(startDate)

	suite.NoError(err)
	suite.Len(emailInfo, 0)
}

func (suite *NotificationSuite) TestHTMLTemplateRender() {
	startDate := time.Date(2019, 1, 7, 0, 0, 0, 0, time.UTC)
	onDate := startDate.AddDate(0, 0, -6)
	mr, err := NewMoveReviewed(suite.DB(), suite.logger, onDate)
	suite.NoError(err)
	s := moveReviewedEmailData{
		Link:                   "www.survey",
		OriginDutyStation:      "OriginDutyStation",
		DestinationDutyStation: "DestDutyStation",
		Email:                  "email",
	}
	expectedHTMLContent := `<p><strong>Good news:</strong> Your move from OriginDutyStation to DestDutyStation has been processed for payment.</p>

<p>Can we ask a quick favor? <a href="www.survey"> Tell us about your experience</a> with requesting and receiving payment.</p>

<p>We'll use your feedback to make MilMove better for your fellow service members.</p>

<p>Thank you for your thoughts, and <strong>congratulations on your move.</strong></p>`

	htmlContent, err := mr.RenderHTML(s)

	suite.NoError(err)
	suite.Equal(expectedHTMLContent, htmlContent)

}

func (suite *NotificationSuite) TestTextTemplateRender() {
	startDate := time.Date(2019, 1, 7, 0, 0, 0, 0, time.UTC)
	onDate := startDate.AddDate(0, 0, -6)
	mr, err := NewMoveReviewed(suite.DB(), suite.logger, onDate)
	suite.NoError(err)
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

	textContent, err := mr.RenderText(s)

	suite.NoError(err)
	suite.Equal(expectedTextContent, textContent)
}

func (suite *NotificationSuite) TestFormatEmails() {
	startDate := time.Date(2019, 1, 7, 0, 0, 0, 0, time.UTC)
	onDate := startDate.AddDate(0, 0, -6)
	mr, err := NewMoveReviewed(suite.DB(), suite.logger, onDate)
	suite.NoError(err)
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
		htmlBody, err := mr.RenderHTML(data)
		suite.NoError(err)
		textBody, err := mr.RenderText(data)
		suite.NoError(err)
		expectedEmailContent := emailContent{
			recipientEmail: *emailInfo.Email,
			subject:        "[MilMove] Let us know how we did",
			htmlBody:       htmlBody,
			textBody:       textBody,
		}
		if emailInfo.Email != nil {
			suite.Equal(expectedEmailContent, actualEmailContent)
		}
	}
	// only expect the two moves with non-nil email addresses to get added to formattedEmails
	suite.Len(formattedEmails, 2)
}
