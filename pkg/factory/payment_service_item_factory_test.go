package factory

import (
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/unit"
)

func (suite *FactorySuite) TestBuildPaymentServiceItem() {
	suite.Run("Successful creation of payment request ", func() {
		// Under test:      BuildPaymentServiceItem
		// Mocked:          None
		// Set up:          Create a payment request with no customizations or traits
		// Expected outcome:paymentServiceItem should be created with default values

		// SETUP
		// CALL FUNCTION UNDER TEST
		paymentServiceItem := BuildPaymentServiceItem(suite.DB(), nil, nil)

		// VALIDATE RESULTS
		defaultCents := unit.Cents(888)
		suite.False(paymentServiceItem.PaymentRequestID.IsNil())
		suite.False(paymentServiceItem.PaymentRequest.ID.IsNil())
		suite.False(paymentServiceItem.MTOServiceItemID.IsNil())
		suite.False(paymentServiceItem.MTOServiceItem.ID.IsNil())
		suite.Equal(models.PaymentServiceItemStatusRequested,
			paymentServiceItem.Status)
		suite.NotNil(paymentServiceItem.PriceCents)
		suite.Equal(defaultCents, *paymentServiceItem.PriceCents)
	})

	suite.Run("Successful creation of customized PaymentServiceItem", func() {
		// Under test:      BuildPaymentServiceItem
		// Mocked:          None
		// Set up:          Create a payment request with and pass custom fields
		// Expected outcome:paymentServiceItem should be created with custom values

		// SETUP
		customPaymentRequest := models.PaymentRequest{
			PaymentRequestNumber: "123",
		}
		customMTOServiceItem := models.MTOServiceItem{
			Reason: models.StringPointer("custom reason"),
		}
		customCents := unit.Cents(789)
		customPaymentServiceItem := models.PaymentServiceItem{
			PriceCents: &customCents,
		}
		customs := []Customization{
			{
				Model: customPaymentRequest,
			},
			{
				Model: customMTOServiceItem,
			},
			{
				Model: customPaymentServiceItem,
			},
		}
		// CALL FUNCTION UNDER TEST
		paymentServiceItem := BuildPaymentServiceItem(suite.DB(), customs, nil)

		// VALIDATE RESULTS
		suite.False(paymentServiceItem.PaymentRequestID.IsNil())
		suite.False(paymentServiceItem.PaymentRequest.ID.IsNil())
		suite.Equal(customPaymentRequest.PaymentRequestNumber,
			paymentServiceItem.PaymentRequest.PaymentRequestNumber)

		suite.False(paymentServiceItem.MTOServiceItemID.IsNil())
		suite.False(paymentServiceItem.MTOServiceItem.ID.IsNil())
		suite.NotNil(paymentServiceItem.MTOServiceItem.Reason)
		suite.Equal(*customMTOServiceItem.Reason,
			*paymentServiceItem.MTOServiceItem.Reason)

		suite.Equal(models.PaymentServiceItemStatusRequested,
			paymentServiceItem.Status)
		suite.NotNil(paymentServiceItem.PriceCents)
		suite.Equal(customCents, *paymentServiceItem.PriceCents)
	})

	suite.Run("Successful return of linkOnly PaymentServiceItem", func() {
		// Under test:       BuildPaymentServiceItem
		// Set up:           Pass in a linkOnly paymentServiceItem
		// Expected outcome: No new PaymentServiceItem should be created.

		// Check num PaymentServiceItem records
		precount, err := suite.DB().Count(&models.PaymentServiceItem{})
		suite.NoError(err)

		id := uuid.Must(uuid.NewV4())
		paymentServiceItem := BuildPaymentServiceItem(suite.DB(), []Customization{
			{
				Model: models.PaymentServiceItem{
					ID: id,
				},
				LinkOnly: true,
			},
		}, nil)
		count, err := suite.DB().Count(&models.PaymentServiceItem{})
		suite.Equal(precount, count)
		suite.NoError(err)
		suite.Equal(id, paymentServiceItem.ID)
	})

	suite.Run("Successful return of stubbed PaymentServiceItem", func() {
		// Under test:       BuildPaymentServiceItem
		// Set up:           Pass in nil db
		// Expected outcome: No new PaymentServiceItem should be created.

		// Check num PaymentServiceItem records
		precount, err := suite.DB().Count(&models.PaymentServiceItem{})
		suite.NoError(err)

		customPaymentServiceItem := models.PaymentServiceItem{
			Status: models.PaymentServiceItemStatusEDIError,
		}
		// Nil passed in as db
		paymentServiceItem := BuildPaymentServiceItem(nil, []Customization{
			{
				Model: customPaymentServiceItem,
			},
		}, nil)

		count, err := suite.DB().Count(&models.PaymentServiceItem{})
		suite.Equal(precount, count)
		suite.NoError(err)

		suite.Equal(customPaymentServiceItem.Status, paymentServiceItem.Status)
	})

}
