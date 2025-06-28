package edi_errors

import (
	"testing"

	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/suite"

	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services/pagination"
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
	suite.Run("returns list of edi errors", func() {
		appCtx := suite.AppContextForTest()

		paymentRequest, err := testdatagen.MakePaymentRequest(appCtx.DB(), testdatagen.Assertions{
			PaymentRequest: models.PaymentRequest{
				Status: models.PaymentRequestStatusEDIError,
			},
		})
		suite.NoError(err)

		ediCode824 := "FailureForEDI824"
		ediDescription824 := "DescriptionForEDI824"
		ediError824 := testdatagen.MakeEdiError(appCtx.DB(), testdatagen.Assertions{
			EdiError: models.EdiError{
				PaymentRequestID: paymentRequest.ID,
				Code:             &ediCode824,
				EDIType:          models.EDIType824,
				Description:      &ediDescription824,
			},
		})
		suite.NotNil(ediError824)

		paymentRequest2, err := testdatagen.MakePaymentRequest(appCtx.DB(), testdatagen.Assertions{
			PaymentRequest: models.PaymentRequest{
				Status: models.PaymentRequestStatusEDIError,
			},
		})
		suite.NoError(err)

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
		pagination := pagination.NewPagination(models.Int64Pointer(pagination.DefaultPage()), models.Int64Pointer(pagination.DefaultPerPage()))
		results, _, err := fetcher.FetchEdiErrors(appCtx, pagination)

		suite.NoError(err)
		suite.NotEmpty(results)
		suite.Equal(2, len(results))
		suite.Equal(ediCode858, *results[0].Code)
		suite.Equal(ediCode824, *results[1].Code)
		suite.Equal(models.EDIType858, results[0].EDIType)
		suite.Equal(models.EDIType824, results[1].EDIType)
		suite.Equal(ediDescription858, *results[0].Description)
		suite.Equal(ediDescription824, *results[1].Description)
	})

	suite.Run("does not return anything if no payment requests found in EDI_ERROR status", func() {
		// no payment requests in EDI_ERROR is a valid scenario, we don't want to error out so we return an empty slice
		appCtx := suite.AppContextForTest()

		fetcher := NewEDIErrorFetcher()
		pagination := pagination.NewPagination(models.Int64Pointer(pagination.DefaultPage()), models.Int64Pointer(pagination.DefaultPerPage()))
		results, _, err := fetcher.FetchEdiErrors(appCtx, pagination)

		suite.NoError(err)
		suite.Empty(results)
	})
}

func (suite *EdiErrorsSuite) TestFetchEdiErrorByID() {
	suite.Run("returns edi error by ID", func() {

		appCtx := suite.AppContextForTest()

		// fetch a single edi error
		paymentRequest, err := testdatagen.MakePaymentRequest(appCtx.DB(), testdatagen.Assertions{
			PaymentRequest: models.PaymentRequest{
				Status: models.PaymentRequestStatusEDIError,
			},
		})
		suite.NoError(err)

		ediCode824 := "FailureForEDI824"
		ediDescription824 := "DescriptionForEDI824"
		ediError824 := testdatagen.MakeEdiError(appCtx.DB(), testdatagen.Assertions{
			EdiError: models.EdiError{
				PaymentRequestID: paymentRequest.ID,
				Code:             &ediCode824,
				EDIType:          models.EDIType824,
				Description:      &ediDescription824,
			},
		})
		suite.NotNil(ediError824)

		fetcher := NewEDIErrorFetcher()
		result, err := fetcher.FetchEdiErrorByID(appCtx, ediError824.ID)

		suite.NoError(err)
		suite.NotEmpty(result)
		suite.Equal(ediError824.ID, result.ID)
		suite.Equal(ediCode824, *result.Code)
		suite.Equal(models.EDIType824, result.EDIType)
		suite.Equal(ediDescription824, *result.Description)
		suite.Equal(paymentRequest.ID, result.PaymentRequestID)
		suite.Equal(paymentRequest.PaymentRequestNumber, result.PaymentRequest.PaymentRequestNumber)

		// fetch a second edi error
		paymentRequest2, err := testdatagen.MakePaymentRequest(appCtx.DB(), testdatagen.Assertions{
			PaymentRequest: models.PaymentRequest{
				Status: models.PaymentRequestStatusEDIError,
			},
		})
		suite.NoError(err)

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

		result, err = fetcher.FetchEdiErrorByID(appCtx, ediError858.ID)

		suite.NoError(err)
		suite.NotEmpty(result)
		suite.Equal(ediError858.ID, result.ID)
		suite.Equal(ediCode858, *result.Code)
		suite.Equal(models.EDIType858, result.EDIType)
		suite.Equal(ediDescription858, *result.Description)
		suite.Equal(paymentRequest2.ID, result.PaymentRequestID)
		suite.Equal(paymentRequest2.PaymentRequestNumber, result.PaymentRequest.PaymentRequestNumber)
	})

	suite.Run("returns not found error when EDIError does not exist", func() {
		appCtx := suite.AppContextForTest()
		fetcher := NewEDIErrorFetcher()

		nonexistentID := uuid.Must(uuid.NewV4())
		result, err := fetcher.FetchEdiErrorByID(appCtx, nonexistentID)

		suite.Error(err)
		suite.IsType(apperror.NotFoundError{}, err)
		suite.Contains(err.Error(), "EDIError not found")
		suite.Equal(models.EdiError{}, result)
	})
}
