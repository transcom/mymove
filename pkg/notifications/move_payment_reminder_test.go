package notifications

import (
	"time"

	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/models"
)

func offsetDate(dayOffset int) time.Time {
	currentDatetime := time.Now()
	return currentDatetime.AddDate(0, 0, dayOffset)
}

func (suite *NotificationSuite) CreatePPMShipmentDateTooOld() models.PPMShipment {
	cutoffDate, _ := time.Parse("2006-01-02", "2019-05-31")
	ppm := factory.BuildPPMShipment(suite.DB(), []factory.Customization{
		{
			Model: models.PPMShipment{
				ExpectedDepartureDate: cutoffDate,
			},
		},
	}, []factory.Trait{
		factory.GetTraitPPMShipmentReadyForPaymentRequest,
	})
	return ppm
}

func (suite *NotificationSuite) GetPPMShipment(offset int) models.PPMShipment {
	expectedDate := offsetDate(offset)

	ppm := factory.BuildPPMShipment(suite.DB(), []factory.Customization{
		{
			Model: models.PPMShipment{
				ExpectedDepartureDate: expectedDate,
			},
		},
	}, []factory.Trait{
		factory.GetTraitPPMShipmentReadyForPaymentRequest,
	})

	return ppm
}

func (suite *NotificationSuite) CreatePPMShipment(offset int) {
	expectedDate := offsetDate(offset)

	factory.BuildPPMShipment(suite.DB(), []factory.Customization{
		{
			Model: models.PPMShipment{
				ExpectedDepartureDate: expectedDate,
			},
		},
	}, []factory.Trait{
		factory.GetTraitPPMShipmentReadyForPaymentRequest,
	})
}

func (suite *NotificationSuite) TestPaymentReminderFetchSomeFound() {

	ppms := []models.PPMShipment{suite.GetPPMShipment(-9), suite.GetPPMShipment(-14), suite.GetPPMShipment(-15)}

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
	suite.CreatePPMShipment(-10)
	suite.CreatePPMShipment(-9)
	suite.CreatePPMShipmentDateTooOld()

	PaymentReminder, err := NewPaymentReminder()
	suite.NoError(err)
	emailInfo, err := PaymentReminder.GetEmailInfo(suite.AppContextForTest())

	suite.NoError(err)
	suite.Len(emailInfo, 0)
}

func (suite *NotificationSuite) TestPaymentReminderFetchAlreadySentEmail() {
	suite.CreatePPMShipmentDateTooOld()
	suite.CreatePPMShipment(-14)

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
	serviceMember := factory.BuildServiceMember(suite.DB(), nil, nil)
	serviceMenmberID := PaymentReminderEmailInfo{
		ServiceMemberID: serviceMember.ID,
	}

	PaymentReminder, err := NewPaymentReminder()
	suite.NoError(err)
	err = PaymentReminder.OnSuccess(suite.AppContextForTest(), serviceMenmberID)("SESID")
	suite.NoError(err)

	notification := models.Notification{}
	err = suite.DB().First(&notification)
	suite.NoError(err)
	suite.Equal(serviceMember.ID, notification.ServiceMemberID)
	suite.Equal(models.MovePaymentReminderEmail, notification.NotificationType)
	suite.Equal("SESID", notification.SESMessageID)
}

func (suite *NotificationSuite) TestPaymentReminderHTMLTemplateRender() {
	pr, err := NewPaymentReminder()
	suite.NoError(err)

	paymentReminderData := PaymentReminderEmailData{
		OriginDutyLocation:      "OriginDutyLocation",
		DestinationDutyLocation: "DestDutyLocation",
		Locator:                 "abc123",
		OneSourceLink:           OneSourceTransportationOfficeLink,
		MyMoveLink:              MyMoveLink,
	}
	expectedHTMLContent := `<p><strong>***</strong> DO NOT REPLY directly to this email <strong>***</strong></p>

<p>This is a reminder that your PPM with the <strong>assigned move code ` + paymentReminderData.Locator + `</strong> from
<strong>` + paymentReminderData.OriginDutyLocation + `</strong> to <strong>` + paymentReminderData.DestinationDutyLocation + `</strong> is awaiting action in MilMove.</p>

<p>Next steps:</p>

<p>To get your payment, you need to login to MilMove, document expenses, and request payment.</p>

<p>To do that:</p>

<p>
<ul>
  <li>Log into <a href=` + MyMoveLink + `>MilMove</a></li>
  <li>Click on "Upload PPM Documents"</li>
  <li>Follow the instructions</li>
</ul>
</p>

<p>To request payment, you should have copies of:</p>

<ul>
<li>       Weight tickets from certified scales, documenting empty and full weights for all vehicles and trailers you used for your move.</li>
<li>       Receipts for reimbursable expenses.</li>
</ul>

<p>MilMove will ask you to upload copies of your documents as you complete your payment request.

<p>If you are missing reciepts, you may still be able to request payment, but you will need assistance from your transportation office.</p>

<p>Payment request must be submitted within 45 days of your move date.</p>

<p>If you have any questions, contact a government transportation office. You can see a listing of
transportation offices on Military OneSource here: &lt;<a href="` + paymentReminderData.OneSourceLink + `">` + paymentReminderData.OneSourceLink + `</a>&gt;</p>

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
		MyMoveLink:              MyMoveLink,
	}
	expectedTextContent := `*** DO NOT REPLY directly to this email ***

This is a reminder that your PPM with the assigned move code ` + paymentReminderData.Locator + ` from ` + paymentReminderData.OriginDutyLocation +
		`
to ` + paymentReminderData.DestinationDutyLocation + ` is awaiting action in MilMove.

Next steps:

To get your payment, you need to login to MilMove, document expenses, and request payment.

To do that:

  * Log into MilMove<` + MyMoveLink + `>
  * Click on "Upload PPM Documents"
  * Follow the instructions

To request payment, you should have copies of:

*       Weight tickets from certified scales, documenting empty and full weights for all vehicles and trailers you used for your move.
*       Receipts for reimbursable expenses.

MilMove will ask you to upload copies of your documents as you complete your payment request.

If you are missing reciepts, you may still be able to request payment, but you will need assistance from your transportation office.

Payment request must be submitted within 45 days of your move date.

If you have any questions, contact a government transportation office. You can see a listing of transportation offices on Military OneSource here: <` + paymentReminderData.OneSourceLink + `>

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
			MyMoveLink:              MyMoveLink,
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

func (suite *NotificationSuite) TestFormatPaymentRequestedEmailsForRetireeSeparation() {
	pr, err := NewPaymentReminder()
	suite.NoError(err)

	email1 := "email1"
	streetOne1 := "100 Street Rd"
	streetTwo1 := "STE 1"
	streetThree1 := "Floor 1"
	city1 := "Alpha City"
	state1 := "Alabama"
	postalCode1 := "11111"
	expectedDestination1 := "100 Street Rd STE 1 Floor 1, Alpha City, Alabama 11111"

	email2 := "email2"
	expectedDestination2 := "nd2" // no street address, fall back to duty station

	email3 := "email3"
	streetOne3 := "300 Highway Ln"
	city3 := "Charlie City"
	state3 := "California"
	postalCode3 := "33333"
	expectedDestination3 := "300 Highway Ln, Charlie City, California 33333"

	email4 := "email4"
	streetOne4 := "400 Unused Blvd"
	city4 := "Delta City"
	state4 := "Delaware"
	postalCode4 := "44444"
	expectedDestination4 := "nd4" // Permanent Change of Station, ignore street address

	email5 := "email5"
	streetOne5 := "500 Parkway Dr"
	expectedDestination5 := "500 Parkway Dr" // Tolerate other nil address fields

	expectedDestinations := []string{
		expectedDestination1,
		expectedDestination2,
		expectedDestination3,
		expectedDestination4,
		expectedDestination5,
	}

	emailInfos := PaymentReminderEmailInfos{
		{
			Email:                 &email1,
			NewDutyLocationName:   "nd1",
			Locator:               "abc123",
			OrdersType:            internalmessages.OrdersTypeRETIREMENT,
			DestinationStreet1:    &streetOne1,
			DestinationStreet2:    &streetTwo1,
			DestinationStreet3:    &streetThree1,
			DestinationCity:       &city1,
			DestinationState:      &state1,
			DestinationPostalCode: &postalCode1,
		},
		{
			Email:               &email2,
			NewDutyLocationName: "nd2",
			Locator:             "abc456",
			OrdersType:          internalmessages.OrdersTypeRETIREMENT,
		},
		{
			Email:                 &email3,
			NewDutyLocationName:   "nd3",
			Locator:               "def123",
			OrdersType:            internalmessages.OrdersTypeSEPARATION,
			DestinationStreet1:    &streetOne3,
			DestinationCity:       &city3,
			DestinationState:      &state3,
			DestinationPostalCode: &postalCode3,
		},
		{
			Email:                 &email4,
			NewDutyLocationName:   "nd4",
			Locator:               "def456",
			OrdersType:            internalmessages.OrdersTypePERMANENTCHANGEOFSTATION,
			DestinationStreet1:    &streetOne4,
			DestinationCity:       &city4,
			DestinationState:      &state4,
			DestinationPostalCode: &postalCode4,
		},
		{
			Email:                 &email5,
			NewDutyLocationName:   "nd5",
			Locator:               "ghi123",
			OrdersType:            internalmessages.OrdersTypeRETIREMENT,
			DestinationStreet1:    &streetOne5,
			DestinationCity:       nil,
			DestinationState:      nil,
			DestinationPostalCode: nil,
		},
	}
	formattedEmails, err := pr.formatEmails(suite.AppContextForTest(), emailInfos)
	suite.NoError(err)

	for i, actualEmailContent := range formattedEmails {
		emailInfo := emailInfos[i]
		expectedDestination := expectedDestinations[i]

		data := PaymentReminderEmailData{
			DestinationDutyLocation: expectedDestination,
			Locator:                 emailInfo.Locator,
			OneSourceLink:           OneSourceTransportationOfficeLink,
			MyMoveLink:              MyMoveLink,
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
	suite.Len(formattedEmails, 5)
}
