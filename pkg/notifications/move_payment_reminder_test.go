package notifications

import (
	"fmt"
	"time"

	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/unit"
)

// func (suite *NotificationSuite) createPaymentReminderMoves(assertions []models.PPMShipment) []models.PPMShipment {
// 	//ppms := make([]models.PPMShipment, 0)
// 	var ppms []models.PPMShipment
// 	estimateIncentive := unit.Cents(1000)

// 	for _, assertion := range assertions {
// 		assertion.EstimatedIncentive = &estimateIncentive
// 		assertion.Shipment.Status = models.MTOShipmentStatusApproved
// 		assertion.ExpectedDepartureDate = offsetDate(-14)

// 		//ppm2 := testdatagen.MakePPM(suite.DB(), assertion)
// 		ppm := factory.BuildPPMShipment(suite.AppContextForTest().DB(), nil, nil)
// 		ppms = append(ppms, ppm)
// 		ppms = append(ppms, assertions...)
// 	}
// 	return ppms
// }

func offsetDate(dayOffset int) time.Time {
	currentDatetime := time.Now()
	return currentDatetime.AddDate(0, 0, dayOffset)
}

// cutoff date for sending payment reminders (don't send if older than this...)
func cutoffDate() time.Time {
	cutoffDate, _ := time.Parse("2006-01-02", "2019-05-31")

	return cutoffDate
}

func (suite *NotificationSuite) TestPaymentReminderFetchSomeFound() {
	date14DaysAgo := offsetDate(-14)
	date9DaysAgo := offsetDate(-9)

	weightEstimate := unit.Pound(300)

	ppms := []models.PPMShipment{
		factory.BuildPPMShipment(suite.DB(), []factory.Customization{
			{
				Model: models.PPMShipment{
					ExpectedDepartureDate: date14DaysAgo,
					Status:                models.PPMShipmentStatusWaitingOnCustomer,
					EstimatedWeight:       &weightEstimate,
				},
			},
			{
				Model: models.Move{
					Status: models.MoveStatusAPPROVED,
				},
			},
			{
				Model: models.MTOShipment{
					Status: models.MTOShipmentStatusApproved,
				},
			},
		}, nil),
		factory.BuildPPMShipment(suite.DB(), []factory.Customization{
			{
				Model: models.PPMShipment{
					ExpectedDepartureDate: offsetDate(-15),
					Status:                models.PPMShipmentStatusWaitingOnCustomer,
					EstimatedWeight:       &weightEstimate,
				},
			},
			{
				Model: models.Move{
					Status: models.MoveStatusAPPROVED,
				},
			},
			{
				Model: models.MTOShipment{
					Status: models.MTOShipmentStatusApproved,
				},
			},
		}, nil),
		factory.BuildPPMShipment(suite.DB(), []factory.Customization{
			{
				Model: models.PPMShipment{
					ExpectedDepartureDate: date9DaysAgo,
					Status:                models.PPMShipmentStatusWaitingOnCustomer,
					EstimatedWeight:       &weightEstimate,
				},
			},
			{
				Model: models.Move{
					Status: models.MoveStatusAPPROVED,
				},
			},
			{
				Model: models.MTOShipment{
					Status: models.MTOShipmentStatusApproved,
				},
			},
		}, nil),
	}

	PaymentReminder, err := NewPaymentReminder()
	suite.NoError(err)
	emailInfo, err := PaymentReminder.GetEmailInfo(suite.AppContextForTest())
	suite.NoError(err)
	suite.NotNil(emailInfo)
	suite.Len(emailInfo, 2, "Wrong number of rows returned")

	for i := 0; i < len(ppms); i++ {
		for j := 0; j < len(emailInfo); j++ {
			if ppms[i].Shipment.MoveTaskOrder.Locator == emailInfo[j].Locator {
				suite.Equal(ppms[i].Shipment.MoveTaskOrder.Orders.NewDutyLocation.Name, emailInfo[j].NewDutyLocationName)
				suite.Equal(ppms[i].Shipment.MoveTaskOrder.Orders.NewDutyLocation.Name, emailInfo[j].NewDutyLocationName)
				suite.NotNil(emailInfo[j].Email)
				suite.Equal(*ppms[i].Shipment.MoveTaskOrder.Orders.ServiceMember.PersonalEmail, *emailInfo[j].Email)
				suite.Equal(ppms[i].EstimatedWeight, emailInfo[j].WeightEstimate)
				suite.Equal(ppms[i].EstimatedIncentive, emailInfo[j].IncentiveEstimate)
				suite.Equal(ppms[i].Shipment.MoveTaskOrder.Orders.OriginDutyLocation.TransportationOffice.Name, *emailInfo[j].TOName)
				suite.Equal(ppms[i].Shipment.MoveTaskOrder.Orders.OriginDutyLocation.TransportationOffice.PhoneLines[0].Number, *emailInfo[j].TOPhone)
				suite.Equal(ppms[i].Shipment.MoveTaskOrder.Locator, emailInfo[j].Locator)
			}
		}
	}
}

func (suite *NotificationSuite) TestPaymentReminderFetchNoneFound() {
	date10DaysAgo := offsetDate(-10)
	date9DaysAgo := offsetDate(-9)
	dateTooOld := cutoffDate()
	weightEstimate := unit.Pound(100)

	factory.BuildPPMShipment(suite.DB(), []factory.Customization{
		{
			Model: models.PPMShipment{
				ExpectedDepartureDate: date9DaysAgo,
				Status:                models.PPMShipmentStatusWaitingOnCustomer,
				EstimatedWeight:       &weightEstimate,
			},
		},
		{
			Model: models.Move{
				Status: models.MoveStatusAPPROVED,
			},
		},
		{
			Model: models.MTOShipment{
				Status: models.MTOShipmentStatusApproved,
			},
		},
	}, nil)

	factory.BuildPPMShipment(suite.DB(), []factory.Customization{
		{
			Model: models.PPMShipment{
				ExpectedDepartureDate: dateTooOld,
				Status:                models.PPMShipmentStatusWaitingOnCustomer,
				EstimatedWeight:       &weightEstimate,
			},
		},
		{
			Model: models.Move{
				Status: models.MoveStatusAPPROVED,
			},
		},
		{
			Model: models.MTOShipment{
				Status: models.MTOShipmentStatusApproved,
			},
		},
	}, nil)

	factory.BuildPPMShipment(suite.DB(), []factory.Customization{
		{
			Model: models.PPMShipment{
				ExpectedDepartureDate: date10DaysAgo,
				Status:                models.PPMShipmentStatusWaitingOnCustomer,
				EstimatedWeight:       &weightEstimate,
			},
		},
		{
			Model: models.Move{
				Status: models.MoveStatusAPPROVED,
			},
		},
		{
			Model: models.MTOShipment{
				Status: models.MTOShipmentStatusApproved,
			},
		},
	}, nil)

	factory.BuildPPMShipment(suite.DB(), []factory.Customization{
		{
			Model: models.PPMShipment{
				ExpectedDepartureDate: date9DaysAgo,
				Status:                models.PPMShipmentStatusWaitingOnCustomer,
				EstimatedWeight:       &weightEstimate,
			},
		},
		{
			Model: models.Move{
				Show: models.BoolPointer(false),
			},
		},
		{
			Model: models.MTOShipment{
				Status: models.MTOShipmentStatusApproved,
			},
		},
	}, nil)

	PaymentReminder, err := NewPaymentReminder()
	suite.NoError(err)
	emailInfo, err := PaymentReminder.GetEmailInfo(suite.AppContextForTest())

	suite.NoError(err)
	suite.Len(emailInfo, 0)
}

func (suite *NotificationSuite) TestPaymentReminderFetchAlreadySentEmail() {
	date14DaysAgo := offsetDate(-14)
	dateTooOld := cutoffDate()
	weightEstimate := unit.Pound(200)

	factory.BuildPPMShipment(suite.DB(), []factory.Customization{
		{
			Model: models.PPMShipment{
				ExpectedDepartureDate: date14DaysAgo,
				Status:                models.PPMShipmentStatusWaitingOnCustomer,
				EstimatedWeight:       &weightEstimate,
			},
		},
		{
			Model: models.Move{
				Status: models.MoveStatusAPPROVED,
			},
		},
		{
			Model: models.MTOShipment{
				Status: models.MTOShipmentStatusApproved,
			},
		},
	}, nil)

	factory.BuildPPMShipment(suite.DB(), []factory.Customization{
		{
			Model: models.PPMShipment{
				ExpectedDepartureDate: dateTooOld,
				Status:                models.PPMShipmentStatusWaitingOnCustomer,
				EstimatedWeight:       &weightEstimate,
			},
		},
		{
			Model: models.Move{
				Status: models.MoveStatusAPPROVED,
			},
		},
		{
			Model: models.MTOShipment{
				Status: models.MTOShipmentStatusApproved,
			},
		},
	}, nil)

	PaymentReminder, err := NewPaymentReminder()
	suite.NoError(err)
	emailInfoBeforeSending, err := PaymentReminder.GetEmailInfo(suite.AppContextForTest())
	suite.NoError(err)
	suite.Len(emailInfoBeforeSending, 1)

	err = PaymentReminder.OnSuccess(suite.AppContextForTest(), emailInfoBeforeSending[0])("SESID")
	suite.NoError(err)
	emailInfoAfterSending, err := PaymentReminder.GetEmailInfo(suite.AppContextForTest())
	suite.NoError(err)
	suite.Len(emailInfoAfterSending, 0)
}

func (suite *NotificationSuite) TestPaymentReminderOnSuccess() {
	sm := factory.BuildServiceMember(suite.DB(), nil, nil)
	ei := PaymentReminderEmailInfo{
		ServiceMemberID: sm.ID,
	}

	PaymentReminder, err := NewPaymentReminder()
	suite.NoError(err)
	err = PaymentReminder.OnSuccess(suite.AppContextForTest(), ei)("SESID")
	suite.NoError(err)

	n := models.Notification{}
	err = suite.DB().First(&n)
	suite.NoError(err)
	suite.Equal(sm.ID, n.ServiceMemberID)
	suite.Equal(models.MovePaymentReminderEmail, n.NotificationType)
	suite.Equal("SESID", n.SESMessageID)
}

// func (suite *NotificationSuite) TestPaymentReminderHTMLTemplateRender() {
// 	test := "https://my.move.mil/sign-in"
// 	pr, err := NewPaymentReminder()
// 	suite.NoError(err)

// 	name := "TEST PPPO"
// 	phone := "555-555-5555"
// 	s := PaymentReminderEmailData{
// 		DestinationDutyLocation: "DestDutyLocation",
// 		WeightEstimate:          "1500",
// 		IncentiveEstimate:       "500",
// 		IncentiveTxt:            "You expected to move about 1500 lbs, which gives you an estimated incentive of $500-$1000.",
// 		TOName:                  &name,
// 		TOPhone:                 &phone,
// 		Locator:                 "abc123",
// 	}
// 	expectedHTMLContent := `<p>We hope your move to DestDutyLocation went well.</p>

// <p>It’s been a couple of weeks, so we want to make sure you get paid for that move. You expected to move about 1500 lbs, which gives you an estimated incentive of $500-$1000.</p>

// <p>To get your incentive, you need to request payment.</p>

// <p>Log in to MilMove and request payment</p>

// <p>We want to pay you for your PPM, but we can’t do that until you document expenses and request payment.</p>

// <p>To do that</p>

// <ul>
//   <li><a href="` + test + `">Log in to MilMove</a></li>
//   <li>Click Request Payment</li>
//   <li>Follow the instructions.</li>
// </ul>

// <p>What documents do you need?</p>

// <p>To request payment, you should have copies of:</p>
// <ul>
//   <li>Weight tickets from certified scales, documenting empty and full weights for all vehicles and trailers you used for your move</li>
//   <li>Receipts for reimbursable expenses (see our moving tips PDF for more info)</li>
// </ul>

// <p>MilMove will ask you to upload copies of your documents as you complete your payment request.</p>

// <p>What if you’re missing documents?</p>

// <p>If you’re missing receipts, you can still request payment. You might not get reimbursement or a tax credit for those expenses.</p>

// <p>If you’re missing certified weight tickets, your PPPO will have to help. Call TEST PPPO at 555-555-5555 to have them walk you through it. Reference your move locator code: abc123.</p>

// <p>Log in to MilMove to complete your request and get paid.</p>

// <p>Request payment within 45 days of your move date or you might not be able to get paid.</p>

// <p>If you have any questions or concerns, you can talk to a human! Call your local PPPO at TEST PPPO at 555-555-5555. Reference your move locator code: abc123.</p>
// `

// 	htmlContent, err := pr.RenderHTML(suite.AppContextForTest(), s)

// 	suite.NoError(err)
// 	suite.Equal(expectedHTMLContent, htmlContent)

// }

func (suite *NotificationSuite) TestPaymentReminderHTMLTemplateRenderNoOriginDutyLocation() {
	pr, err := NewPaymentReminder()
	suite.NoError(err)

	s := PaymentReminderEmailData{
		DestinationDutyLocation: "DestDutyLocation",
		WeightEstimate:          "1500",
		IncentiveEstimate:       "500",
		IncentiveTxt:            "You expected to move about 1500 lbs, which gives you an estimated incentive of $500-$1000.",
		TOName:                  nil,
		TOPhone:                 nil,
		Locator:                 "abc123",
	}
	expectedHTMLContent := `<p>We hope your move to DestDutyLocation went well.</p>

<p>It’s been a couple of weeks, so we want to make sure you get paid for that move. You expected to move about 1500 lbs, which gives you an estimated incentive of $500-$1000.</p>

<p>To get your incentive, you need to request payment.</p>

<p>Log in to MilMove and request payment</p>

<p>We want to pay you for your PPM, but we can’t do that until you document expenses and request payment.</p>

<p>To do that</p>

<ul>
  <li><a href="https://my.move.mil">Log in to MilMove</a></li>
  <li>Click Request Payment</li>
  <li>Follow the instructions.</li>
</ul>

<p>What documents do you need?</p>

<p>To request payment, you should have copies of:</p>
<ul>
  <li>Weight tickets from certified scales, documenting empty and full weights for all vehicles and trailers you used for your move</li>
  <li>Receipts for reimbursable expenses (see our moving tips PDF for more info)</li>
</ul>

<p>MilMove will ask you to upload copies of your documents as you complete your payment request.</p>

<p>What if you’re missing documents?</p>

<p>If you’re missing receipts, you can still request payment. You might not get reimbursement or a tax credit for those expenses.</p>

<p>If you are missing weight tickets, someone from the government will have to help. Consult Military OneSource's <a href="https://www.militaryonesource.mil/moving-housing/moving/planning-your-move/customer-service-contacts-for-military-pcs/">directory of PCS-related contacts</a> to find your best contact and reference your move code abc123.</p>

<p>Log in to MilMove to complete your request and get paid.</p>

<p>Request payment within 45 days of your move date or you might not be able to get paid.</p>


`

	htmlContent, err := pr.RenderHTML(suite.AppContextForTest(), s)

	suite.NoError(err)
	suite.Equal(expectedHTMLContent, htmlContent)

}

func (suite *NotificationSuite) TestPaymentReminderTextTemplateRender() {
	pr, err := NewPaymentReminder()
	suite.NoError(err)

	name := "TEST PPPO"
	phone := "555-555-5555"
	s := PaymentReminderEmailData{
		DestinationDutyLocation: "DestDutyLocation",
		WeightEstimate:          "1500",
		IncentiveEstimate:       "500",
		IncentiveTxt:            "You expected to move about 1500 lbs, which gives you an estimated incentive of $500-$1000.",
		TOName:                  &name,
		TOPhone:                 &phone,
		Locator:                 "abc123",
	}
	expectedTextContent := `We hope your move to DestDutyLocation went well.

It’s been a couple of weeks, so we want to make sure you get paid for that move. You expected to move about 1500 lbs, which gives you an estimated incentive of $500-$1000.

To get your incentive, you need to request payment.

Log in to MilMove and request payment

We want to pay you for your PPM, but we can’t do that until you document expenses and request payment.

To do that

  * Log in to MilMove
  * Click Request Payment
  * Follow the instructions.

What documents do you need?

To request payment, you should have copies of:
  * Weight tickets from certified scales, documenting empty and full weights for all vehicles and trailers you used for your move
  * Receipts for reimbursable expenses (see our moving tips PDF for more info)

MilMove will ask you to upload copies of your documents as you complete your payment request.

What if you’re missing documents?

If you’re missing receipts, you can still request payment. You might not get reimbursement or a tax credit for those expenses.

If you’re missing certified weight tickets, your PPPO will have to help. Call TEST PPPO at 555-555-5555 to have them walk you through it. Reference your move locator code: abc123.

Log in to MilMove to complete your request and get paid.

Request payment within 45 days of your move date or you might not be able to get paid.

If you have any questions or concerns, you can talk to a human! Call your local PPPO at TEST PPPO at 555-555-5555. Reference your move locator code: abc123.
`

	textContent, err := pr.RenderText(suite.AppContextForTest(), s)

	suite.NoError(err)
	suite.Equal(expectedTextContent, textContent)
}

func (suite *NotificationSuite) TestFormatPaymentRequestedEmails() {
	pr, err := NewPaymentReminder()
	suite.NoError(err)
	email1 := "email1"
	weightEst1 := unit.Pound(100)
	estimate1 := unit.Cents(1000)
	phone1 := "111-111-1111"

	email2 := "email2"
	weightEst2 := unit.Pound(200)
	estimate2 := unit.Cents(2000)
	phone2 := ""

	email3 := "email3"
	weightEst3 := unit.Pound(0)
	estimate3 := unit.Cents(0)

	phone := "000-000-0000"

	name1 := "to1"
	name2 := "to2"
	name3 := "to3"
	name4 := "to4"

	emailInfos := PaymentReminderEmailInfos{
		{
			Email:               &email1,
			NewDutyLocationName: "nd1",
			WeightEstimate:      &weightEst1,
			IncentiveEstimate:   &estimate1,
			IncentiveTxt:        fmt.Sprintf("You expected to move about %d lbs, which gives you an estimated incentive of %s.", weightEst1.Int(), estimate1.ToDollarString()),
			TOName:              &name1,
			TOPhone:             &phone1,
			Locator:             "abc123",
		},
		{
			Email:               &email2,
			NewDutyLocationName: "nd2",
			WeightEstimate:      &weightEst2,
			IncentiveEstimate:   &estimate2,
			IncentiveTxt:        fmt.Sprintf("You expected to move about %d lbs, which gives you an estimated incentive of %s.", weightEst2.Int(), estimate2.ToDollarString()),
			TOName:              &name2,
			TOPhone:             &phone2,
			Locator:             "abc456",
		},
		{
			Email:               &email3,
			NewDutyLocationName: "nd3",
			WeightEstimate:      &weightEst3,
			IncentiveEstimate:   &estimate3,
			IncentiveTxt:        "",
			TOName:              &name3,
			TOPhone:             &phone,
			Locator:             "def123",
		},
		{
			// nil emails should be skipped
			Email:               nil,
			NewDutyLocationName: "nd0",
			WeightEstimate:      &weightEst3,
			IncentiveEstimate:   &estimate3,
			IncentiveTxt:        "",
			TOName:              &name4,
			TOPhone:             &phone,
			Locator:             "def456",
		},
	}
	formattedEmails, err := pr.formatEmails(suite.AppContextForTest(), emailInfos)

	suite.NoError(err)
	for i, actualEmailContent := range formattedEmails {
		emailInfo := emailInfos[i]

		data := PaymentReminderEmailData{
			DestinationDutyLocation: emailInfo.NewDutyLocationName,
			WeightEstimate:          fmt.Sprintf("%d", emailInfo.WeightEstimate.Int()),
			IncentiveEstimate:       emailInfo.IncentiveEstimate.ToDollarString(),
			IncentiveTxt:            emailInfo.IncentiveTxt,
			TOName:                  emailInfo.TOName,
			TOPhone:                 emailInfo.TOPhone,
			Locator:                 emailInfo.Locator,
		}
		htmlBody, err := pr.RenderHTML(suite.AppContextForTest(), data)
		suite.NoError(err)
		textBody, err := pr.RenderText(suite.AppContextForTest(), data)
		suite.NoError(err)
		expectedEmailContent := emailContent{
			recipientEmail: *emailInfo.Email,
			subject:        fmt.Sprintf("[MilMove] Reminder: request payment for your move to %s (move %s)", emailInfo.NewDutyLocationName, emailInfo.Locator),
			htmlBody:       htmlBody,
			textBody:       textBody,
		}
		if emailInfo.Email != nil {
			suite.Equal(expectedEmailContent.recipientEmail, actualEmailContent.recipientEmail)
			suite.Equal(expectedEmailContent.subject, actualEmailContent.subject)
			suite.Equal(expectedEmailContent.htmlBody, actualEmailContent.htmlBody, "htmlBody diffferent: %s", emailInfo.TOName)
			suite.Equal(expectedEmailContent.textBody, actualEmailContent.textBody)
		}
	}
	// only expect the three moves with non-nil email addresses to get added to formattedEmails
	suite.Len(formattedEmails, 3)
}
