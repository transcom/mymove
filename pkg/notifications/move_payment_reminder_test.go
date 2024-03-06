package notifications

import (
	"time"

	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/unit"
)

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

func (suite *NotificationSuite) TestPaymentReminderHTMLTemplateRender() {
	pr, err := NewPaymentReminder()
	suite.NoError(err)

	paymentReminderData := PaymentReminderEmailData{
		OriginDutyLocation:      "OriginDutyLocation",
		DestinationDutyLocation: "DestDutyLocation",
		Locator:                 "abc123",
		OneSourceLink:           OneSourceTransportationOfficeLink,
	}
	expectedHTMLContent := `<p>*** DO NOT REPLY directly to this email ***</p>

<p>This is a reminder that your PPM with the <strong>assigned move code ` + paymentReminderData.Locator + `</strong> from <strong>` + paymentReminderData.OriginDutyLocation +
		`</strong>
to <strong>` + paymentReminderData.DestinationDutyLocation + `</strong> is awaiting action in MilMove.</p>

<p>To get your payment, you need to login to MilMove, document expenses, and request payment.</p>

<p>To do that:</p>

<ul>
  <li> Log into MilMove</li>
  <li> Click on "Upload PPM Documents"</li>
  <li> Follow the instructions</li>
</ul>

To request payment, you should have copies of:</p>

<ul>
<li>       Weight tickets from certified scales, documenting empty and full weights for all vehicles and
trailers you used for your move.</li>
<li>       Receipts for reimbursable expenses (see our moving tips PDF for more info <a href="` + paymentReminderData.OneSourceLink + `">` +
		paymentReminderData.OneSourceLink + `)</a></li>
</ul>

<p>MilMove will ask you to upload copies of your documents as you complete your payment request.

<p>If you are missing reciepts, you can still request payment but may not get reimbursement or a tax credit
for those expenses.</p>

<p>Payment request must be submitted within 45 days of your move date.</p>

<p>If you have any questions, contact a government transportation office. You can see a listing of</p>
transportation offices on Military OneSource here: &lt;<a href="` + paymentReminderData.OneSourceLink + `">` + paymentReminderData.OneSourceLink + `</a>&gt;

<p>Thank you,</p>

<p>USTRANSCOM MilMove Team</p>

<p>The information contained in this email may contain Privacy Act information and is therefore protected
under the Privacy Act of 1974. Failure to protect Privacy Act information could result in a $5,000 fine.</p>`

	htmlContent, err := pr.RenderHTML(suite.AppContextForTest(), paymentReminderData)

	suite.NoError(err)
	suite.Equal(expectedHTMLContent, htmlContent)

}

func (suite *NotificationSuite) TestPaymentReminderTextTemplateRender() {
	pr, err := NewPaymentReminder()
	suite.NoError(err)
	paymentReminderData := PaymentReminderEmailData{
		OriginDutyLocation:      "OriginDutyLocation",
		DestinationDutyLocation: "DestDutyLocation",
		Locator:                 "abc123",
		OneSourceLink:           OneSourceTransportationOfficeLink,
	}
	expectedTextContent := `*** DO NOT REPLY directly to this email ***

This is a reminder that your PPM with the assigned move code ` + paymentReminderData.Locator + ` from ` + paymentReminderData.OriginDutyLocation +
		`
to ` + paymentReminderData.DestinationDutyLocation + ` is awaiting action in MilMove.

To get your payment, you need to login to MilMove, document expenses, and request payment.

To do that:

  * Log into MilMove
  * Click on "Upload PPM Documents"
  * Follow the instructions

To request payment, you should have copies of:

*       Weight tickets from certified scales, documenting empty and full weights for all vehicles and
trailers you used for your move.
*       Receipts for reimbursable expenses (see our moving tips PDF for more info ` + paymentReminderData.OneSourceLink + `)

MilMove will ask you to upload copies of your documents as you complete your payment request.

If you are missing reciepts, you can still request payment but may not get reimbursement or a tax credit
for those expenses.

Payment request must be submitted within 45 days of your move date.

If you have any questions, contact a government transportation office. You can see a listing of
transportation offices on Military OneSource here: <` + paymentReminderData.OneSourceLink + `>

Thank you,

USTRANSCOM MilMove Team

The information contained in this email may contain Privacy Act information and is therefore protected
under the Privacy Act of 1974. Failure to protect Privacy Act information could result in a $5,000 fine.`

	textContent, err := pr.RenderText(suite.AppContextForTest(), paymentReminderData)

	suite.NoError(err)
	suite.Equal(expectedTextContent, textContent)
}

func (suite *NotificationSuite) TestFormatPaymentRequestedEmails() {
	pr, err := NewPaymentReminder()
	suite.NoError(err)
	email1 := "email1"

	email2 := "email2"

	email3 := "email3"

	emailInfos := PaymentReminderEmailInfos{
		{
			Email:               &email1,
			NewDutyLocationName: "nd1",
			Locator:             "abc123",
		},
		{
			Email:               &email2,
			NewDutyLocationName: "nd2",
			Locator:             "abc456",
		},
		{
			Email:               &email3,
			NewDutyLocationName: "nd3",
			Locator:             "def123",
		},
		{
			// nil emails should be skipped
			Email:               nil,
			NewDutyLocationName: "nd0",
			Locator:             "def456",
		},
	}
	formattedEmails, err := pr.formatEmails(suite.AppContextForTest(), emailInfos)

	suite.NoError(err)
	for i, actualEmailContent := range formattedEmails {
		emailInfo := emailInfos[i]

		data := PaymentReminderEmailData{
			DestinationDutyLocation: emailInfo.NewDutyLocationName,
			Locator:                 emailInfo.Locator,
			OneSourceLink:           OneSourceTransportationOfficeLink,
		}
		htmlBody, err := pr.RenderHTML(suite.AppContextForTest(), data)
		suite.NoError(err)
		textBody, err := pr.RenderText(suite.AppContextForTest(), data)
		suite.NoError(err)
		expectedEmailContent := emailContent{
			recipientEmail: *emailInfo.Email,
			subject:        "Complete your Personally Procured Move (PPM)",
			htmlBody:       htmlBody,
			textBody:       textBody,
		}
		if emailInfo.Email != nil {
			suite.Equal(expectedEmailContent.recipientEmail, actualEmailContent.recipientEmail)
			suite.Equal(expectedEmailContent.subject, actualEmailContent.subject)
			suite.Equal(expectedEmailContent.textBody, actualEmailContent.textBody)
		}
	}
	// only expect the three moves with non-nil email addresses to get added to formattedEmails
	suite.Len(formattedEmails, 3)
}
