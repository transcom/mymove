package paymentrequest

import (
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/query"
	"github.com/transcom/mymove/pkg/testingsuite"
)

func (suite *PaymentRequestSyncadaFileFetcherSuite) TestFetchPaymentRequestSyncadaFile() {
	builder := query.NewQueryBuilder()
	fetcher := NewPaymentRequestSyncadaFileFetcher(builder)

	suite.Run("Fetch Syncada files", func() {
		// Set up test data
		paymentRequestEdiFile := BuildPaymentRequestEdiRecord("858.rec1", "someStringedi", "1234-7654-1")
		err := suite.DB().Create(&paymentRequestEdiFile)
		suite.NoError(err)

		result, err := fetcher.FetchPaymentRequestSyncadaFile(suite.AppContextForTest(), []services.QueryFilter{})
		suite.NoError(err)
		suite.NotNil(result)
		// Add more assertions here to verify the content of the result
	})
}

type PaymentRequestSyncadaFileFetcherSuite struct {
	testingsuite.PopTestSuite
}

func BuildPaymentRequestEdiRecord(fileName string, ediString string, prNumber string) models.PaymentRequestEdiFile {
	var paymentRequestEdiFile models.PaymentRequestEdiFile
	paymentRequestEdiFile.ID = uuid.Must(uuid.NewV4())
	paymentRequestEdiFile.EdiString = ediString
	paymentRequestEdiFile.Filename = fileName
	paymentRequestEdiFile.PaymentRequestNumber = prNumber

	return paymentRequestEdiFile
}
