package notifications

import (
	"fmt"
	"time"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/unit"
)

func (suite *NotificationSuite) createPaymentReminderMoves(assertions []testdatagen.Assertions) []models.PersonallyProcuredMove {
	ppms := make([]models.PersonallyProcuredMove, 0)
	estimateMin := unit.Cents(1000)
	estimateMax := unit.Cents(2000)

	for _, assertion := range assertions {
		assertion.PersonallyProcuredMove.Status = models.PPMStatusAPPROVED
		assertion.PersonallyProcuredMove.IncentiveEstimateMin = &estimateMin
		assertion.PersonallyProcuredMove.IncentiveEstimateMax = &estimateMax

		ppm := testdatagen.MakePPM(suite.DB(), assertion)
		ppms = append(ppms, ppm)
	}
	return ppms
}

func offsetDate(dayOffset int) time.Time {
	currentDatetime := time.Now()
	return currentDatetime.AddDate(0, 0, dayOffset)
}

// cutoff date for sending payment reminders (don't send if older than this...)
func cutoffDate() time.Time {
	cutoffDate, _ := time.Parse("2019-01-01", "2019-10-01")
	return cutoffDate
}

func (suite *NotificationSuite) TestPaymentReminderFetchSomeFound() {
	db := suite.DB()
	date10DaysAgo := offsetDate(-10)
	date9DaysAgo := offsetDate(-9)

	moves := []testdatagen.Assertions{
		{PersonallyProcuredMove: models.PersonallyProcuredMove{ReviewedDate: &date10DaysAgo}},
		{PersonallyProcuredMove: models.PersonallyProcuredMove{ReviewedDate: &date9DaysAgo}},
	}

	ppms := suite.createPaymentReminderMoves(moves)

	PaymentReminder, err := NewPaymentReminder(db, suite.logger)
	suite.NoError(err)
	emailInfo, err := PaymentReminder.GetEmailInfo()
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
	date9DaysAgo := offsetDate(-9)
	dateTooOld := cutoffDate()

	moves := []testdatagen.Assertions{
		{PersonallyProcuredMove: models.PersonallyProcuredMove{ReviewedDate: &date9DaysAgo}},
		{PersonallyProcuredMove: models.PersonallyProcuredMove{ReviewedDate: &dateTooOld}},
	}

	suite.createPaymentReminderMoves(moves)

	PaymentReminder, err := NewPaymentReminder(db, suite.logger)
	suite.NoError(err)
	emailInfo, err := PaymentReminder.GetEmailInfo()

	suite.NoError(err)
	suite.Len(emailInfo, 0)
}

func (suite *NotificationSuite) TestPaymentReminderFetchAlreadySentEmail() {
	db := suite.DB()

	date10DaysAgo := offsetDate(-10)
	dateTooOld := cutoffDate()

	moves := []testdatagen.Assertions{
		{PersonallyProcuredMove: models.PersonallyProcuredMove{ReviewedDate: &date10DaysAgo}},
		{PersonallyProcuredMove: models.PersonallyProcuredMove{ReviewedDate: &dateTooOld}},
	}
	suite.createPaymentReminderMoves(moves)

	suite.createPPMMoves(moves)
	PaymentReminder, err := NewPaymentReminder(db, suite.logger)
	suite.NoError(err)
	emailInfoBeforeSending, err := PaymentReminder.GetEmailInfo()
	suite.NoError(err)
	suite.Len(emailInfoBeforeSending, 2)

	// simulate successfully sending an email and then check that this email does not get sent again.
	err = PaymentReminder.OnSuccess(emailInfoBeforeSending[0])("SES_MOVE_ID")
	suite.NoError(err)
	emailInfoAfterSending, err := PaymentReminder.GetEmailInfo()
	suite.NoError(err)
	suite.Len(emailInfoAfterSending, 1)
}

func (suite *NotificationSuite) TestPaymentReminderOnSuccess() {
	db := suite.DB()
	sm := testdatagen.MakeDefaultServiceMember(db)
	ei := PaymentReminderEmailInfo{
		ServiceMemberID: sm.ID,
	}

	PaymentReminder, err := NewPaymentReminder(db, suite.logger)
	suite.NoError(err)
	err = PaymentReminder.OnSuccess(ei)("SESID")
	suite.NoError(err)

	n := models.Notification{}
	err = db.First(&n)
	suite.NoError(err)
	suite.Equal(sm.ID, n.ServiceMemberID)
	suite.Equal(models.MovePaymentReminderEmail, n.NotificationType)
	suite.Equal("SESID", n.SESMessageID)
}

func (suite *NotificationSuite) TestPaymentReminderHTMLTemplateRender() {
	pr, err := NewPaymentReminder(suite.DB(), suite.logger)
	suite.NoError(err)
	s := PaymentReminderEmailData{
		DestinationDutyStation: "DestDutyStation",
		WeightEstimate:         "1500",
		IncentiveEstimateMin:   "500",
		IncentiveEstimateMax:   "1000",
		TOName:                 "TEST PPPO",
		TOPhone:                "555-555-5555",
	}
	expectedHTMLContent := `<p>We hope your move to DestDutyStation went well.</p>

<p>It’s been a couple of weeks, so we want to make sure you get paid for that move. You expected to move about 1500 lbs, which gives you an estimated incentive of 500-1000.</p>

<p>To get your incentive, you need to request payment.</p>

<p>Log in to MilMove and request payment</p>

<p>We want to pay you for your PPM, but we can’t do that until you document expenses and request payment.</p>

<p>To do that</p>

<p>Log in to MilMove</p>
<ul>
  <li>Click Request Payment</li>
  <li>Follow the instructions.</li>
  <li>What documents do you need?</li>
</ul>

<p>To request payment, you should have copies of:</p>
<p>Weight tickets from certified scales, documenting empty and full weights for all vehicles and trailers you used for your move</p>
<p>Receipts for reimbursable expenses (see our moving tips PDF for more info)</p>

<p>MilMove will ask you to upload copies of your documents as you complete your payment request.</p>

<p>What if you’re missing documents?</p>

<p>If you’re missing receipts, you can still request payment. You might not get reimbursement or a tax credit for those expenses.</p>

<p>If you’re missing certified weight tickets, your PPPO will have to help. Call the TEST PPPO at 555-555-5555 to have them walk you through it.</p>

<p>Log in to MilMove to complete your request and get paid.</p>

<p>Request payment within 45 days of your move date or you might not be able to get paid.</p>

<p>If you have any questions or concerns, you can talk to a human! Call your local PPPO at TEST PPPO at 555-555-5555.</p>
`

	htmlContent, err := pr.RenderHTML(s)

	suite.NoError(err)
	suite.Equal(expectedHTMLContent, htmlContent)

}

func (suite *NotificationSuite) TestPaymentReminderTextTemplateRender() {
	pr, err := NewPaymentReminder(suite.DB(), suite.logger)
	suite.NoError(err)

	s := PaymentReminderEmailData{
		DestinationDutyStation: "DestDutyStation",
		WeightEstimate:         "1500",
		IncentiveEstimateMin:   "500",
		IncentiveEstimateMax:   "1000",
		TOName:                 "TEST PPPO",
		TOPhone:                "555-555-5555",
	}
	expectedTextContent := `We hope your move to DestDutyStation went well.

It’s been a couple of weeks, so we want to make sure you get paid for that move. You expected to move about 1500 lbs, which gives you an estimated incentive of 500-1000.

To get your incentive, you need to request payment.

Log in to MilMove and request payment

We want to pay you for your PPM, but we can’t do that until you document expenses and request payment.

To do that

Log in to MilMove
  * Click Request Payment
  * Follow the instructions.
  * What documents do you need?

To request payment, you should have copies of:
Weight tickets from certified scales, documenting empty and full weights for all vehicles and trailers you used for your move
Receipts for reimbursable expenses (see our moving tips PDF for more info)

MilMove will ask you to upload copies of your documents as you complete your payment request.

What if you’re missing documents?

If you’re missing receipts, you can still request payment. You might not get reimbursement or a tax credit for those expenses.

If you’re missing certified weight tickets, your PPPO will have to help. Call the TEST PPPO at 555-555-5555 to have them walk you through it.

Log in to MilMove to complete your request and get paid.

Request payment within 45 days of your move date or you might not be able to get paid.

If you have any questions or concerns, you can talk to a human! Call your local PPPO at TEST PPPO at 555-555-5555.
`

	textContent, err := pr.RenderText(s)

	suite.NoError(err)
	suite.Equal(expectedTextContent, textContent)
}

func (suite *NotificationSuite) TestFormatPaymentRequestedEmails() {
	pr, err := NewPaymentReminder(suite.DB(), suite.logger)
	suite.NoError(err)
	email1 := "email1"
	weightEst1 := unit.Pound(100)
	estimateMin1 := unit.Cents(1000)
	estimateMax1 := unit.Cents(1100)

	email2 := "email2"
	weightEst2 := unit.Pound(200)
	estimateMin2 := unit.Cents(2000)
	estimateMax2 := unit.Cents(2200)

	weightEst3 := unit.Pound(0)
	estimateMin3 := unit.Cents(0)
	estimateMax3 := unit.Cents(0)

	emailInfos := PaymentReminderEmailInfos{
		{
			Email:                &email1,
			NewDutyStationName:   "nd1",
			WeightEstimate:       &weightEst1,
			IncentiveEstimateMin: &estimateMin1,
			IncentiveEstimateMax: &estimateMax1,
			TOName:               "to1",
			TOPhone:              "111-111-1111",
		},
		{
			Email:                &email2,
			NewDutyStationName:   "nd2",
			WeightEstimate:       &weightEst2,
			IncentiveEstimateMin: &estimateMin2,
			IncentiveEstimateMax: &estimateMax2,
			TOName:               "to2",
			TOPhone:              "222-222-2222",
		},
		{
			// nil emails should be skipped
			Email:                nil,
			NewDutyStationName:   "nd0",
			WeightEstimate:       &weightEst3,
			IncentiveEstimateMin: &estimateMin3,
			IncentiveEstimateMax: &estimateMax3,
			TOName:               "to0",
			TOPhone:              "000-000-0000",
		},
	}

	formattedEmails, err := pr.formatEmails(emailInfos)

	suite.NoError(err)
	for i, actualEmailContent := range formattedEmails {
		emailInfo := emailInfos[i]

		data := PaymentReminderEmailData{
			DestinationDutyStation: emailInfo.NewDutyStationName,
			WeightEstimate:         fmt.Sprintf("%d", emailInfo.WeightEstimate),
			IncentiveEstimateMin:   emailInfo.IncentiveEstimateMin.ToDollarString(),
			IncentiveEstimateMax:   emailInfo.IncentiveEstimateMax.ToDollarString(),
			TOName:                 emailInfo.TOName,
			TOPhone:                emailInfo.TOPhone,
		}

		htmlBody, err := pr.RenderHTML(data)
		suite.NoError(err)
		textBody, err := pr.RenderText(data)
		suite.NoError(err)
		expectedEmailContent := emailContent{
			recipientEmail: *emailInfo.Email,
			subject:        fmt.Sprintf("[MilMove] Reminder: request payment for your move to %s", emailInfo.NewDutyStationName),
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
