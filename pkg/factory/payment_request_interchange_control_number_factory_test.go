package factory

import (
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
)

func (suite *FactorySuite) TestBuildPaymentRequestToInterchangeControlNumber() {
	suite.Run("Successful creation of payment request control number ", func() {
		// Under test:      BuildPaymentRequestToInterchangeControlNumber
		// Mocked:          None
		// Set up:          Create a payment request control number with no customizations or traits
		// Expected outcome:paymentRequest2InterchangeControlNumber should be created with default values

		// SETUP
		// CALL FUNCTION UNDER TEST
		pr2cin := BuildPaymentRequestToInterchangeControlNumber(suite.DB(), nil, nil)

		// VALIDATE RESULTS
		suite.Equal(models.EDIType858, pr2cin.EDIType)
		suite.NotNil(pr2cin.InterchangeControlNumber)
		suite.NotNil(pr2cin.PaymentRequestID)
	})

	suite.Run("Successful creation of customized paymentRequestControlNumber", func() {
		// Under test:      BuildPaymentRequestToInterchangeControlNumber
		// Mocked:          None
		// Set up:          Create a payment request control number with custom fields
		// Expected outcome:paymentRequest2InterchangeControlNumber should be created with custom values

		// SETUP
		customPr2cin := models.PaymentRequestToInterchangeControlNumber{
			InterchangeControlNumber: 994,
			EDIType:                  models.EDIType997,
		}

		customPaymentRequest := BuildPaymentRequest(suite.DB(), nil, nil)

		// CALL FUNCTION UNDER TEST
		pr2cin := BuildPaymentRequestToInterchangeControlNumber(suite.DB(), []Customization{
			{
				Model: customPr2cin,
			},
			{
				Model:    customPaymentRequest,
				LinkOnly: true,
			},
		}, nil)

		// VALIDATE RESULTS
		suite.Equal(customPr2cin.InterchangeControlNumber, pr2cin.InterchangeControlNumber)
		suite.Equal(customPr2cin.EDIType, pr2cin.EDIType)
		suite.Equal(customPaymentRequest.ID, pr2cin.PaymentRequestID)
	})

	suite.Run("Successful return of linkOnly paymentRequestControlNumber", func() {
		// Under test:       BuildPaymentRequestToInterchangeControlNumber
		// Set up:           Pass in a linkOnly paymentRequestControlNumber
		// Expected outcome: No new PaymentRequestControlNumber should be created.

		// Check num PaymentRequestToInterchangeControlNumber records
		precount, err := suite.DB().Count(&models.PaymentRequestToInterchangeControlNumber{})
		suite.NoError(err)

		id := uuid.Must(uuid.NewV4())
		pr2cin := BuildPaymentRequestToInterchangeControlNumber(suite.DB(), []Customization{
			{
				Model: models.PaymentRequestToInterchangeControlNumber{
					ID: id,
				},
				LinkOnly: true,
			},
		}, nil)
		count, err := suite.DB().Count(&models.PaymentRequestToInterchangeControlNumber{})
		suite.Equal(precount, count)
		suite.NoError(err)
		suite.Equal(id, pr2cin.ID)
	})

	suite.Run("Successful return of stubbed paymentRequestControlNumber", func() {
		// Under test:       BuildPaymentRequestToInterchangeControlNumber
		// Set up:           Pass in nil db
		// Expected outcome: No new PaymentRequestControlNumber should be created.

		// Check num PaymentRequestToInterchangeControlNumber records
		precount, err := suite.DB().Count(&models.PaymentRequestToInterchangeControlNumber{})
		suite.NoError(err)

		customPr2cin := models.PaymentRequestToInterchangeControlNumber{
			EDIType:                  models.EDIType997,
			InterchangeControlNumber: 48412,
		}
		// Nil passed in as db
		pr2cin := BuildPaymentRequestToInterchangeControlNumber(nil, []Customization{
			{
				Model: customPr2cin,
			},
		}, nil)

		count, err := suite.DB().Count(&models.PaymentRequestToInterchangeControlNumber{})
		suite.Equal(precount, count)
		suite.NoError(err)

		suite.Equal(customPr2cin.InterchangeControlNumber, pr2cin.InterchangeControlNumber)
		suite.Equal(customPr2cin.EDIType, pr2cin.EDIType)
	})
}
