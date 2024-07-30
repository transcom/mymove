package factory

import (
	"time"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
)

func (suite *FactorySuite) TestBuildPaymentRequest() {
	suite.Run("Successful creation of payment request ", func() {
		// Under test:      BuildPaymentRequest
		// Mocked:          None
		// Set up:          Create a payment request with no customizations or traits
		// Expected outcome:paymentRequest should be created with default values

		// SETUP
		// CALL FUNCTION UNDER TEST
		paymentRequest := BuildPaymentRequest(suite.DB(), nil, nil)

		// VALIDATE RESULTS
		suite.False(paymentRequest.MoveTaskOrderID.IsNil())
		suite.False(paymentRequest.MoveTaskOrder.ID.IsNil())
		suite.False(paymentRequest.IsFinal)
		suite.Nil(paymentRequest.RejectionReason)
		suite.Equal(models.PaymentRequestStatusPending, paymentRequest.Status)
		suite.Equal(1, paymentRequest.SequenceNumber)
		suite.NotNil(paymentRequest.MoveTaskOrder.ReferenceID)
		suite.Equal(*paymentRequest.MoveTaskOrder.ReferenceID+"-1", paymentRequest.PaymentRequestNumber)
	})

	suite.Run("Successful creation of customized PaymentRequest", func() {
		// Under test:      BuildPaymentRequest
		// Mocked:          None
		// Set up:          Create a payment request with and pass custom fields
		// Expected outcome:paymentRequest should be created with custom values

		// SETUP
		customMove := models.Move{
			Locator: "ABC123",
			Show:    models.BoolPointer(true),
		}
		customPaymentRequest := models.PaymentRequest{
			Status:               models.PaymentRequestStatusTppsReceived,
			IsFinal:              true,
			RejectionReason:      models.StringPointer("custom reason"),
			PaymentRequestNumber: "abc-2",
			SequenceNumber:       2,
			CreatedAt:            time.Now().Add(time.Hour * -24),
		}
		customs := []Customization{
			{
				Model: customMove,
			},
			{
				Model: customPaymentRequest,
			},
		}
		// CALL FUNCTION UNDER TEST
		paymentRequest := BuildPaymentRequest(suite.DB(), customs, nil)

		// VALIDATE RESULTS
		suite.False(paymentRequest.MoveTaskOrderID.IsNil())
		suite.False(paymentRequest.MoveTaskOrder.ID.IsNil())
		suite.Equal(customMove.Locator, paymentRequest.MoveTaskOrder.Locator)
		suite.Equal(customMove.Show, paymentRequest.MoveTaskOrder.Show)

		suite.Equal(customPaymentRequest.Status, paymentRequest.Status)
		suite.Equal(customPaymentRequest.IsFinal, paymentRequest.IsFinal)
		suite.NotNil(paymentRequest.RejectionReason)
		suite.Equal(*customPaymentRequest.RejectionReason, *paymentRequest.RejectionReason)
		suite.Equal(customPaymentRequest.SequenceNumber, paymentRequest.SequenceNumber)
		suite.Equal(customPaymentRequest.PaymentRequestNumber, paymentRequest.PaymentRequestNumber)
		suite.Equal(customPaymentRequest.CreatedAt, paymentRequest.CreatedAt)
	})

	suite.Run("Successful return of linkOnly PaymentRequest", func() {
		// Under test:       BuildPaymentRequest
		// Set up:           Pass in a linkOnly paymentRequest
		// Expected outcome: No new PaymentRequest should be created.

		// Check num PaymentRequest records
		precount, err := suite.DB().Count(&models.PaymentRequest{})
		suite.NoError(err)

		id := uuid.Must(uuid.NewV4())
		paymentRequest := BuildPaymentRequest(suite.DB(), []Customization{
			{
				Model: models.PaymentRequest{
					ID: id,
				},
				LinkOnly: true,
			},
		}, nil)
		count, err := suite.DB().Count(&models.PaymentRequest{})
		suite.Equal(precount, count)
		suite.NoError(err)
		suite.Equal(id, paymentRequest.ID)
	})

	suite.Run("Successful return of stubbed PaymentRequest", func() {
		// Under test:       BuildPaymentRequest
		// Set up:           Pass in nil db
		// Expected outcome: No new PaymentRequest should be created.

		// Check num PaymentRequest records
		precount, err := suite.DB().Count(&models.PaymentRequest{})
		suite.NoError(err)

		customPaymentRequest := models.PaymentRequest{
			Status: models.PaymentRequestStatusDeprecated,
		}
		// Nil passed in as db
		paymentRequest := BuildPaymentRequest(nil, []Customization{
			{
				Model: customPaymentRequest,
			},
		}, nil)

		count, err := suite.DB().Count(&models.PaymentRequest{})
		suite.Equal(precount, count)
		suite.NoError(err)

		suite.Equal(customPaymentRequest.Status, paymentRequest.Status)
	})

}
