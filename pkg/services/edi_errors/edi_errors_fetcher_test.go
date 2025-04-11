package edi_errors

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/testingsuite"
)

type EdiErrorsSuite struct {
	*testingsuite.PopTestSuite
}

func TestEdiErrorsSuite(t *testing.T) {
	suite.Run(t, &EdiErrorsSuite{
		PopTestSuite: testingsuite.NewPopTestSuite(testingsuite.CurrentPackage(), testingsuite.WithPerTestTransaction()),
	})
}

func (suite *EdiErrorsSuite) TestFetchEdiErrors() {
	appCtx := suite.AppContextForTest()

	paymentRequest := testdatagen.MakePaymentRequest(appCtx.DB(), testdatagen.Assertions{
		PaymentRequest: models.PaymentRequest{
			Status: models.PaymentRequestStatusEDIError,
		},
	})

	ediCode810 := "FailureForEDI810"
	ediDescription810 := "DescriptionForEDI810"
	ediError810 := testdatagen.MakeEdiError(appCtx.DB(), testdatagen.Assertions{
		EdiError: models.EdiError{
			PaymentRequestID: paymentRequest.ID,
			Code:             &ediCode810,
			EDIType:          models.EDIType810,
			Description:      &ediDescription810,
		},
	})
	suite.NotNil(ediError810)

	paymentRequest2 := testdatagen.MakePaymentRequest(appCtx.DB(), testdatagen.Assertions{
		PaymentRequest: models.PaymentRequest{
			Status: models.PaymentRequestStatusEDIError,
		},
	})

	ediCode858 := "FailureForEDI858"
	ediDescription858 := "DescriptionForEDI858"
	ediError858 := testdatagen.MakeEdiError(appCtx.DB(), testdatagen.Assertions{
		EdiError: models.EdiError{
			PaymentRequestID: paymentRequest2.ID,
			Code:             &ediCode858,
			EDIType:          models.EDIType858,
			Description:      &ediDescription858,
		},
	})
	suite.NotNil(ediError858)

	fetcher := NewEDIErrorFetcher()
	results, err := fetcher.FetchEdiErrors(appCtx)

	suite.NoError(err)
	suite.NotEmpty(results)
	suite.Equal(2, len(results))
	suite.Equal(ediCode810, *results[0].Code)
	suite.Equal(ediCode858, *results[1].Code)
	suite.Equal(models.EDIType810, results[0].EDIType)
	suite.Equal(models.EDIType858, results[1].EDIType)
	suite.Equal(ediDescription810, *results[0].Description)
	suite.Equal(ediDescription858, *results[1].Description)
}

func (suite *EdiErrorsSuite) TestFetchEdiErrorsNoPaymentRequestsFound() {
	// no payment requests in EDI_ERROR is fine, we don't want to error out so we return an empty slice
	appCtx := suite.AppContextForTest()

	fetcher := NewEDIErrorFetcher()
	results, err := fetcher.FetchEdiErrors(appCtx)

	suite.NoError(err)
	suite.Empty(results)
}
