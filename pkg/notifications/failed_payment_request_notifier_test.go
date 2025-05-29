package notifications

import (
	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *NotificationSuite) TestPaymentRequestFailedEmails() {
	paymentRequest, _ := testdatagen.MakePaymentRequest(suite.DB(), testdatagen.Assertions{
		PaymentRequest: models.PaymentRequest{
			Status: models.PaymentRequestStatusReviewed,
		},
	})
	err := suite.DB().Create(&models.ApplicationParameters{
		ParameterName:  stringPointer("src_email"),
		ParameterValue: stringPointer("src@example.com"),
	})
	suite.NoError(err)

	ediError := models.EdiError{
		PaymentRequestID: paymentRequest.ID,
		Code:             stringPointer("123"),
		Description:      stringPointer("Test error"),
	}
	err = suite.DB().Create(&ediError)
	suite.NoError(err)
	notification := NewPaymentRequestFailed(paymentRequest)
	emails, err := notification.emails(suite.AppContextWithSessionForTest(&auth.Session{}))

	suite.NoError(err)
	suite.Equal(1, len(emails))
	suite.Equal("Payment Request Failed", emails[0].subject)
	suite.Contains(emails[0].htmlBody, "123")
	suite.Contains(emails[0].textBody, "Test error")

}

func (suite *NotificationSuite) TestPaymentRequestFailedEmailsNoEmails() {
	paymentRequest, _ := testdatagen.MakePaymentRequest(suite.DB(), testdatagen.Assertions{
		PaymentRequest: models.PaymentRequest{
			Status: models.PaymentRequestStatusReviewed,
		},
	})
	ediError := models.EdiError{
		PaymentRequestID: paymentRequest.ID,
		Code:             stringPointer("123"),
		Description:      stringPointer("Test error"),
	}
	err := suite.DB().Create(&ediError)
	suite.NoError(err)
	notification := NewPaymentRequestFailed(paymentRequest)
	_, err = notification.emails(suite.AppContextWithSessionForTest(&auth.Session{}))

	suite.Error(err)

}

func (suite *NotificationSuite) TestPaymentRequestFailedEmailsNoEDIError() {
	paymentRequest, _ := testdatagen.MakePaymentRequest(suite.DB(), testdatagen.Assertions{
		PaymentRequest: models.PaymentRequest{
			Status: models.PaymentRequestStatusReviewed,
		},
	})
	err := suite.DB().Create(&models.ApplicationParameters{
		ParameterName:  stringPointer("src_email"),
		ParameterValue: stringPointer("source@test.com"),
	})
	suite.NoError(err)

	notification := PaymentRequestFailed{paymentRequest: paymentRequest}
	emails, err := notification.emails(suite.AppContextWithSessionForTest(&auth.Session{}))

	suite.NoError(err)
	suite.Equal(0, len(emails))
}

func stringPointer(s string) *string {
	return &s
}
