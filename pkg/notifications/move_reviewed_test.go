package notifications

import (
	"fmt"
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
	startDate := time.Date(2019, 1, 7, 0, 0, 0, 0, time.UTC)
	onDate := startDate.AddDate(0, 0, -6)
	offDate := startDate.AddDate(0, 0, -7)
	moves := []testdatagen.Assertions{
		{PersonallyProcuredMove: models.PersonallyProcuredMove{Status: models.PPMStatusAPPROVED, ReviewedDate: &onDate}},
		{PersonallyProcuredMove: models.PersonallyProcuredMove{Status: models.PPMStatusAPPROVED, ReviewedDate: &offDate}},
	}
	ppms := suite.createPPMMoves(moves)

	moveReviewed, err := NewMoveReviewed(onDate)
	suite.NoError(err)
	emailInfo, err := moveReviewed.GetEmailInfo(suite.AppContextForTest(), onDate)

	suite.NoError(err)
	suite.NotNil(emailInfo)
	suite.Len(emailInfo, 1)
	suite.Equal(ppms[0].Move.Orders.NewDutyLocation.Name, emailInfo[0].NewDutyStationName)
	suite.NotNil(emailInfo[0].Email)
	suite.Equal(*ppms[0].Move.Orders.ServiceMember.PersonalEmail, *emailInfo[0].Email)
	suite.Equal(ppms[0].Move.Orders.ServiceMember.DutyStation.Name, emailInfo[0].DutyStationName)
}

func (suite *NotificationSuite) TestMoveReviewedFetchNoneFound() {
	startDate := time.Date(2019, 1, 7, 0, 0, 0, 0, time.UTC)
	offDate := startDate.AddDate(0, 0, -7)
	moves := []testdatagen.Assertions{
		{PersonallyProcuredMove: models.PersonallyProcuredMove{Status: models.PPMStatusAPPROVED, ReviewedDate: &offDate}},
		{PersonallyProcuredMove: models.PersonallyProcuredMove{Status: models.PPMStatusAPPROVED, ReviewedDate: &offDate}},
	}
	suite.createPPMMoves(moves)

	moveReviewed, err := NewMoveReviewed(startDate)
	suite.NoError(err)
	emailInfo, err := moveReviewed.GetEmailInfo(suite.AppContextForTest(), startDate)

	suite.NoError(err)
	suite.Len(emailInfo, 0)
}

func (suite *NotificationSuite) TestMoveReviewedFetchAlreadySentEmail() {
	startDate := time.Date(2019, 1, 7, 0, 0, 0, 0, time.UTC)
	moves := []testdatagen.Assertions{
		{PersonallyProcuredMove: models.PersonallyProcuredMove{Status: models.PPMStatusAPPROVED, ReviewedDate: &startDate}},
		{PersonallyProcuredMove: models.PersonallyProcuredMove{Status: models.PPMStatusAPPROVED, ReviewedDate: &startDate}},
	}
	suite.createPPMMoves(moves)
	moveReviewed, err := NewMoveReviewed(startDate)
	suite.NoError(err)
	emailInfoBeforeSending, err := moveReviewed.GetEmailInfo(suite.AppContextForTest(), startDate)
	suite.NoError(err)
	suite.Len(emailInfoBeforeSending, 2)

	// simulate successfully sending an email and then check that this email does not get sent again.
	err = moveReviewed.OnSuccess(suite.AppContextForTest(), emailInfoBeforeSending[0])("SES_MOVE_ID")
	suite.NoError(err)
	emailInfoAfterSending, err := moveReviewed.GetEmailInfo(suite.AppContextForTest(), startDate)
	suite.NoError(err)
	suite.Len(emailInfoAfterSending, 1)
}

func (suite *NotificationSuite) TestMoveReviewedOnSuccess() {
	db := suite.DB()
	sm := testdatagen.MakeDefaultServiceMember(db)
	ei := EmailInfo{
		ServiceMemberID: sm.ID,
	}
	startDate := time.Date(2019, 1, 7, 0, 0, 0, 0, time.UTC)
	moveReviewed, err := NewMoveReviewed(startDate)
	suite.NoError(err)
	err = moveReviewed.OnSuccess(suite.AppContextForTest(), ei)("SESID")
	suite.NoError(err)

	n := models.Notification{}
	err = db.First(&n)
	suite.NoError(err)
	suite.Equal(sm.ID, n.ServiceMemberID)
	suite.Equal(models.MoveReviewedEmail, n.NotificationType)
	suite.Equal("SESID", n.SESMessageID)
}

func (suite *NotificationSuite) TestHTMLTemplateRender() {
	startDate := time.Date(2019, 1, 7, 0, 0, 0, 0, time.UTC)
	onDate := startDate.AddDate(0, 0, -6)
	mr, err := NewMoveReviewed(onDate)
	suite.NoError(err)
	s := moveReviewedEmailData{
		Link:                   "www.survey",
		OriginDutyStation:      "OriginDutyLocation",
		DestinationDutyStation: "DestDutyStation",
	}
	expectedHTMLContent := `<p><strong>Good news:</strong> Your move from OriginDutyLocation to DestDutyStation has been processed for payment.</p>

<p>Can we ask a quick favor? <a href="www.survey"> Tell us about your experience</a> with requesting and receiving payment.</p>

<p>We'll use your feedback to make MilMove better for your fellow service members.</p>

<p>Thank you for your thoughts, and <strong>congratulations on your move.</strong></p>`

	htmlContent, err := mr.RenderHTML(suite.AppContextForTest(), s)

	suite.NoError(err)
	suite.Equal(expectedHTMLContent, htmlContent)

}

func (suite *NotificationSuite) TestTextTemplateRender() {
	startDate := time.Date(2019, 1, 7, 0, 0, 0, 0, time.UTC)
	onDate := startDate.AddDate(0, 0, -6)
	mr, err := NewMoveReviewed(onDate)
	suite.NoError(err)
	s := moveReviewedEmailData{
		Link:                   "www.survey",
		OriginDutyStation:      "OriginDutyLocation",
		DestinationDutyStation: "DestDutyStation",
	}
	expectedTextContent := `Good news: Your move from OriginDutyLocation to DestDutyStation has been processed for payment.

Can we ask a quick favor? Tell us about your experience with requesting and receiving payment at www.survey.

We'll use your feedback to make MilMove better for your fellow service members.

Thank you for your thoughts, and congratulations on your move.`

	textContent, err := mr.RenderText(suite.AppContextForTest(), s)

	suite.NoError(err)
	suite.Equal(expectedTextContent, textContent)
}

func (suite *NotificationSuite) TestFormatEmails() {
	startDate := time.Date(2019, 1, 7, 0, 0, 0, 0, time.UTC)
	onDate := startDate.AddDate(0, 0, -6)
	mr, err := NewMoveReviewed(onDate)
	suite.NoError(err)
	email1 := "email1"
	email2 := "email2"
	emailInfos := EmailInfos{
		{
			Email:              &email1,
			DutyStationName:    "d1",
			NewDutyStationName: "nd2",
			Locator:            "abc123",
		},
		{
			Email:              &email2,
			DutyStationName:    "d2",
			NewDutyStationName: "nd2",
			Locator:            "abc456",
		},
		{
			// nil emails should be skipped
			Email:              nil,
			DutyStationName:    "d2",
			NewDutyStationName: "nd2",
			Locator:            "abc788",
		},
	}

	formattedEmails, err := mr.formatEmails(suite.AppContextForTest(), emailInfos)

	suite.NoError(err)
	for i, actualEmailContent := range formattedEmails {
		emailInfo := emailInfos[i]
		data := moveReviewedEmailData{
			Link:                   surveyLink,
			OriginDutyStation:      emailInfo.DutyStationName,
			DestinationDutyStation: emailInfo.NewDutyStationName,
		}
		htmlBody, err := mr.RenderHTML(suite.AppContextForTest(), data)
		suite.NoError(err)
		textBody, err := mr.RenderText(suite.AppContextForTest(), data)
		suite.NoError(err)
		expectedEmailContent := emailContent{
			recipientEmail: *emailInfo.Email,
			subject:        fmt.Sprintf("[MilMove] Tell us how we did with your move (%s)", emailInfo.Locator),
			htmlBody:       htmlBody,
			textBody:       textBody,
		}
		if emailInfo.Email != nil {
			suite.Equal(expectedEmailContent.recipientEmail, actualEmailContent.recipientEmail)
			suite.Equal(expectedEmailContent.subject, actualEmailContent.subject)
			suite.Equal(expectedEmailContent.htmlBody, actualEmailContent.htmlBody)
			suite.Equal(expectedEmailContent.textBody, actualEmailContent.textBody)
		}
	}
	// only expect the two moves with non-nil email addresses to get added to formattedEmails
	suite.Len(formattedEmails, 2)
}
