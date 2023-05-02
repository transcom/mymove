package factory

import (
	"github.com/benbjohnson/clock"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/unit"
)

func (suite *FactorySuite) TestBuildPaymentServiceItemParam() {
	suite.Run("Successful creation of default PaymentServiceItemParam", func() {
		// Under test:      BuildPaymentServiceItemParam
		// Mocked:          None
		// Set up:          Create a payment service item param with no customizations or traits
		// Expected outcome:paymentServiceItemParam should be created with default values

		// SETUP
		defaultPaymentServiceItemParam := models.PaymentServiceItemParam{
			Value: "123",
		}

		// CALL FUNCTION UNDER TEST
		paymentServiceItemParam := BuildPaymentServiceItemParam(suite.DB(), nil, nil)

		// VALIDATE RESULTS
		suite.Equal(defaultPaymentServiceItemParam.Value, paymentServiceItemParam.Value)
		suite.NotNil(paymentServiceItemParam.PaymentServiceItemID)
		suite.NotNil(paymentServiceItemParam.ServiceItemParamKeyID)
	})

	suite.Run("Successful creation of customized PaymentServiceItemParam", func() {
		// Under test:      BuildPaymentServiceItemParam
		// Set up:          Create a payment service item param and pass custom fields
		// Expected outcome:paymentServiceItemParam should be created with custom fields

		// SETUP
		customPaymentServiceItemParam := models.PaymentServiceItemParam{
			Value: "456",
		}

		// CALL FUNCTION UNDER TEST
		paymentServiceItemParam := BuildPaymentServiceItemParam(suite.DB(), []Customization{
			{
				Model: customPaymentServiceItemParam,
			},
		}, nil)

		// VALIDATE RESULTS
		suite.Equal(customPaymentServiceItemParam.Value, paymentServiceItemParam.Value)
	})

	suite.Run("Successful creation of payment service item param with customized payment service item and service item param key", func() {
		// Under test:      BuildPaymentServiceItemParam
		// Set up:          Create a payment service item param and pass custom fields
		// Expected outcome:paymentServiceItemParam should be created with custom payment service item and service item param key

		// SETUP
		customPaymentServiceItemParam := models.PaymentServiceItemParam{
			Value: "456",
		}

		customPaymentServiceItem := models.PaymentServiceItem{
			Status: models.PaymentServiceItemStatusApproved,
		}
		customServiceItemParam := models.ServiceItemParamKey{
			Key:  models.ServiceItemParamNameRequestedPickupDate,
			Type: models.ServiceItemParamTypeDate,
		}

		// CALL FUNCTION UNDER TEST
		paymentServiceItemParam := BuildPaymentServiceItemParam(suite.DB(), []Customization{
			{
				Model: customPaymentServiceItemParam,
			},
			{
				Model: customPaymentServiceItem,
			},
			{
				Model: customServiceItemParam,
			},
		}, nil)

		// VALIDATE RESULTS
		suite.Equal(customPaymentServiceItemParam.Value, paymentServiceItemParam.Value)

		suite.Equal(customPaymentServiceItem.Status, paymentServiceItemParam.PaymentServiceItem.Status)
		suite.Equal(customServiceItemParam.Key, paymentServiceItemParam.ServiceItemParamKey.Key)
		suite.Equal(customServiceItemParam.Type, paymentServiceItemParam.ServiceItemParamKey.Type)
	})

	suite.Run("Successful return of linkOnly paymentServiceItemParam", func() {
		// Under test:       BuildPaymentServiceItemParam
		// Set up:           Pass in a linkOnly paymentServiceItemParam
		// Expected outcome: No new paymentServiceItemParam should be created.

		// Check num PaymentServiceItemParam records
		precount, err := suite.DB().Count(&models.PaymentServiceItemParam{})
		suite.NoError(err)

		id := uuid.Must(uuid.NewV4())
		paymentServiceItemParam := BuildPaymentServiceItemParam(suite.DB(), []Customization{
			{
				Model: models.PaymentServiceItemParam{
					ID: id,
				},
				LinkOnly: true,
			},
		}, nil)
		count, err := suite.DB().Count(&models.PaymentServiceItemParam{})
		suite.Equal(precount, count)
		suite.NoError(err)
		suite.Equal(id, paymentServiceItemParam.ID)
	})

	suite.Run("Successful creation of customized paymentServiceItem with params", func() {
		// Under test:      BuildPaymentServiceItemWithParams
		//
		// Set up: Create a Payment Service Item with custom serviceItemParamKey, paymentServiceItem,
		// & list of params to create
		//
		// Expected outcome:paymentServiceItem should be created with multiple paymentServiceItemParams

		// SETUP
		const testDateFormat = "20060102"
		priceCents := unit.Cents(800)
		currentTime := clock.NewMock().Now()
		serviceCode := models.ReServiceCodeCS

		customPaymentServiceItem := models.PaymentServiceItem{
			Status:     models.PaymentServiceItemStatusApproved,
			PriceCents: &priceCents,
		}

		paymentServiceItemParams := []CreatePaymentServiceItemParams{
			{
				Key:     models.ServiceItemParamNameContractCode,
				KeyType: models.ServiceItemParamTypeString,
				Value:   "Test_value",
			},
			{
				Key:     models.ServiceItemParamNameReferenceDate,
				KeyType: models.ServiceItemParamTypeDate,
				Value:   currentTime.Format(testDateFormat),
			},
			{
				Key:     models.ServiceItemParamNameWeightBilled,
				KeyType: models.ServiceItemParamTypeInteger,
				Value:   "4242",
			},
			{
				Key:     models.ServiceItemParamNameDistanceZip,
				KeyType: models.ServiceItemParamTypeInteger,
				Value:   "24246",
			},
		}

		// CALL FUNCTION UNDER TEST
		paymentServiceItem := BuildPaymentServiceItemWithParams(suite.DB(), serviceCode, paymentServiceItemParams, []Customization{
			{
				Model: customPaymentServiceItem,
			},
		}, nil)

		// VALIDATE RESULTS
		suite.Equal(customPaymentServiceItem.Status, paymentServiceItem.Status)
		suite.Equal(customPaymentServiceItem.PriceCents, paymentServiceItem.PriceCents)

		suite.Equal(4, len(paymentServiceItem.PaymentServiceItemParams))
	})
}
