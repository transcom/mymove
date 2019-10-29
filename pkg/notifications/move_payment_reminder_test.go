package notifications

import (
	"time"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/unit"
)

func (suite *NotificationSuite) createPPMMoves1(assertions []testdatagen.Assertions) []models.PersonallyProcuredMove {
	ppms := make([]models.PersonallyProcuredMove, 0)
	for _, assertion := range assertions {
		ppm := testdatagen.MakePPM(suite.DB(), assertion)
		ppms = append(ppms, ppm)
	}
	return ppms
}

func (suite *NotificationSuite) TestPaymentReminderFetchSomeFound() {
	db := suite.DB()
	currentDatetime := time.Now()
	paymentReminderDate := currentDatetime.AddDate(0, 0, -10)
	noPaymentReminderDate := currentDatetime.AddDate(0, 0, -9)
	estimateMin := unit.Cents(1000)
	estimateMax := unit.Cents(2000)

	moves := []testdatagen.Assertions{
		{PersonallyProcuredMove: models.PersonallyProcuredMove{Status: models.PPMStatusAPPROVED, ReviewedDate: &paymentReminderDate, IncentiveEstimateMin: &estimateMin, IncentiveEstimateMax: &estimateMax}},
		{PersonallyProcuredMove: models.PersonallyProcuredMove{Status: models.PPMStatusAPPROVED, ReviewedDate: &noPaymentReminderDate, IncentiveEstimateMin: &estimateMin, IncentiveEstimateMax: &estimateMax}},
	}
	ppms := suite.createPPMMoves1(moves)

	PaymentReminder, err := NewPaymentReminder(db, suite.logger, paymentReminderDate)
	suite.NoError(err)
	emailInfo, err := PaymentReminder.GetEmailInfo(paymentReminderDate)
	suite.NoError(err)

	suite.NotNil(emailInfo)
	suite.Len(emailInfo, 1, "Wrong number of rows returned")
	suite.Equal(ppms[0].Move.Orders.NewDutyStation.Name, emailInfo[0].NewDutyStationName)
	suite.NotNil(emailInfo[0].Email)
	suite.Equal(*ppms[0].Move.Orders.ServiceMember.PersonalEmail, *emailInfo[0].Email)
	suite.Equal(ppms[0].WeightEstimate, emailInfo[0].WeightEstimate)
	suite.Equal(ppms[0].IncentiveEstimateMin, emailInfo[0].IncentiveEstimateMin)
	suite.Equal(ppms[0].IncentiveEstimateMax, emailInfo[0].IncentiveEstimateMax)
	suite.Equal(ppms[0].Move.Orders.ServiceMember.DutyStation.TransportationOffice.Name, emailInfo[0].TOName)
	suite.Equal(ppms[0].Move.Orders.ServiceMember.DutyStation.TransportationOffice.PhoneLines[0].Number, emailInfo[0].TOPhone)
}

func (suite *NotificationSuite) TestPaymentReminderFetchNoneFound() {
	db := suite.DB()
	currentDatetime := time.Now()
	paymentReminderDateTooNew, _ := time.Parse("2019-01-01", "2019-10-01") // cutoff date for sending payment reminders (don't send if older than this...
	paymentReminderDateTooOld := currentDatetime.AddDate(0, 0, -9)
	estimateMin := unit.Cents(1000)
	estimateMax := unit.Cents(2000)
	startDate := time.Date(2019, 1, 7, 0, 0, 0, 0, time.UTC)
	moves := []testdatagen.Assertions{
		{PersonallyProcuredMove: models.PersonallyProcuredMove{Status: models.PPMStatusAPPROVED, ReviewedDate: &paymentReminderDateTooNew, IncentiveEstimateMin: &estimateMin, IncentiveEstimateMax: &estimateMax}},
		{PersonallyProcuredMove: models.PersonallyProcuredMove{Status: models.PPMStatusAPPROVED, ReviewedDate: &paymentReminderDateTooOld, IncentiveEstimateMin: &estimateMin, IncentiveEstimateMax: &estimateMax}},
	}
	suite.createPPMMoves(moves)

	PaymentReminder, err := NewPaymentReminder(db, suite.logger, startDate)
	suite.NoError(err)
	emailInfo, err := PaymentReminder.GetEmailInfo(startDate)

	suite.NoError(err)
	suite.Len(emailInfo, 0)
}

// func (suite *NotificationSuite) TestPaymentReminderFetchAlreadySentEmail() {
// 	db := suite.DB()
// 	startDate := time.Date(2019, 1, 7, 0, 0, 0, 0, time.UTC)
// 	moves := []testdatagen.Assertions{
// 		{PersonallyProcuredMove: models.PersonallyProcuredMove{Status: models.PPMStatusAPPROVED, ReviewedDate: &startDate}},
// 		{PersonallyProcuredMove: models.PersonallyProcuredMove{Status: models.PPMStatusAPPROVED, ReviewedDate: &startDate}},
// 	}
// 	suite.createPPMMoves(moves)
// 	PaymentReminder, err := NewPaymentReminder(db, suite.logger, startDate)
// 	suite.NoError(err)
// 	emailInfoBeforeSending, err := PaymentReminder.GetEmailInfo(startDate)
// 	suite.NoError(err)
// 	suite.Len(emailInfoBeforeSending, 2)

// 	// simulate successfully sending an email and then check that this email does not get sent again.
// 	err = PaymentReminder.OnSuccess(emailInfoBeforeSending[0])("SES_MOVE_ID")
// 	suite.NoError(err)
// 	emailInfoAfterSending, err := PaymentReminder.GetEmailInfo(startDate)
// 	suite.NoError(err)
// 	suite.Len(emailInfoAfterSending, 1)
// }

// func (suite *NotificationSuite) TestPaymentReminderOnSuccess() {
// 	db := suite.DB()
// 	sm := testdatagen.MakeDefaultServiceMember(db)
// 	ei := EmailInfo{
// 		ServiceMemberID: sm.ID,
// 	}
// 	startDate := time.Date(2019, 1, 7, 0, 0, 0, 0, time.UTC)
// 	PaymentReminder, err := NewPaymentReminder(db, suite.logger, startDate)
// 	suite.NoError(err)
// 	err = PaymentReminder.OnSuccess(ei)("SESID")
// 	suite.NoError(err)

// 	n := models.Notification{}
// 	err = db.First(&n)
// 	suite.NoError(err)
// 	suite.Equal(sm.ID, n.ServiceMemberID)
// 	suite.Equal(models.MovePaymentReminderEmail, n.NotificationType)
// 	suite.Equal("SESID", n.SESMessageID)
// }

// func (suite *NotificationSuite) TestHTMLTemplateRender() {
// 	startDate := time.Date(2019, 1, 7, 0, 0, 0, 0, time.UTC)
// 	onDate := startDate.AddDate(0, 0, -6)
// 	mr, err := NewPaymentReminder(suite.DB(), suite.logger, onDate)
// 	suite.NoError(err)
// 	s := PaymentReminderEmailData{
// 		Link:                   "www.survey",
// 		OriginDutyStation:      "OriginDutyStation",
// 		DestinationDutyStation: "DestDutyStation",
// 	}
// 	expectedHTMLContent := `<p><strong>Good news:</strong> Your move from OriginDutyStation to DestDutyStation has been processed for payment.</p>

// <p>Can we ask a quick favor? <a href="www.survey"> Tell us about your experience</a> with requesting and receiving payment.</p>

// <p>We'll use your feedback to make MilMove better for your fellow service members.</p>

// <p>Thank you for your thoughts, and <strong>congratulations on your move.</strong></p>`

// 	htmlContent, err := mr.RenderHTML(s)

// 	suite.NoError(err)
// 	suite.Equal(expectedHTMLContent, htmlContent)

// }

// func (suite *NotificationSuite) TestTextTemplateRender() {
// 	startDate := time.Date(2019, 1, 7, 0, 0, 0, 0, time.UTC)
// 	onDate := startDate.AddDate(0, 0, -6)
// 	mr, err := NewPaymentReminder(suite.DB(), suite.logger, onDate)
// 	suite.NoError(err)
// 	s := PaymentReminderEmailData{
// 		Link:                   "www.survey",
// 		OriginDutyStation:      "OriginDutyStation",
// 		DestinationDutyStation: "DestDutyStation",
// 	}
// 	expectedTextContent := `Good news: Your move from OriginDutyStation to DestDutyStation has been processed for payment.

// Can we ask a quick favor? Tell us about your experience with requesting and receiving payment at www.survey.

// We'll use your feedback to make MilMove better for your fellow service members.

// Thank you for your thoughts, and congratulations on your move.`

// 	textContent, err := mr.RenderText(s)

// 	suite.NoError(err)
// 	suite.Equal(expectedTextContent, textContent)
// }

// func (suite *NotificationSuite) TestFormatEmails() {
// 	startDate := time.Date(2019, 1, 7, 0, 0, 0, 0, time.UTC)
// 	onDate := startDate.AddDate(0, 0, -6)
// 	mr, err := NewPaymentReminder(suite.DB(), suite.logger, onDate)
// 	suite.NoError(err)
// 	email1 := "email1"
// 	email2 := "email2"
// 	emailInfos := EmailInfos{
// 		{
// 			Email:              &email1,
// 			DutyStationName:    "d1",
// 			NewDutyStationName: "nd2",
// 		},
// 		{
// 			Email:              &email2,
// 			DutyStationName:    "d2",
// 			NewDutyStationName: "nd2",
// 		},
// 		{
// 			// nil emails should be skipped
// 			Email:              nil,
// 			DutyStationName:    "d2",
// 			NewDutyStationName: "nd2",
// 		},
// 	}

// 	formattedEmails, err := mr.formatEmails(emailInfos)

// 	suite.NoError(err)
// 	for i, actualEmailContent := range formattedEmails {
// 		emailInfo := emailInfos[i]
// 		data := PaymentReminderEmailData{
// 			Link:                   surveyLink,
// 			OriginDutyStation:      emailInfo.DutyStationName,
// 			DestinationDutyStation: emailInfo.NewDutyStationName,
// 		}
// 		htmlBody, err := mr.RenderHTML(data)
// 		suite.NoError(err)
// 		textBody, err := mr.RenderText(data)
// 		suite.NoError(err)
// 		expectedEmailContent := emailContent{
// 			recipientEmail: *emailInfo.Email,
// 			subject:        "[MilMove] Let us know how we did",
// 			htmlBody:       htmlBody,
// 			textBody:       textBody,
// 		}
// 		if emailInfo.Email != nil {
// 			suite.Equal(expectedEmailContent.recipientEmail, actualEmailContent.recipientEmail)
// 			suite.Equal(expectedEmailContent.subject, actualEmailContent.subject)
// 			suite.Equal(expectedEmailContent.htmlBody, actualEmailContent.htmlBody)
// 			suite.Equal(expectedEmailContent.textBody, actualEmailContent.textBody)
// 		}
// 	}
// 	// only expect the two moves with non-nil email addresses to get added to formattedEmails
// 	suite.Len(formattedEmails, 2)
// }
