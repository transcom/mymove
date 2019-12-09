package paymentrequest

import (
	"testing"

	"github.com/transcom/mymove/pkg/testdatagen"

	"github.com/transcom/mymove/pkg/models"
)

func (suite *PaymentRequestServiceSuite) TestFetchPaymentRequestList() {
	paymentRequestListFetcher := NewPaymentRequestListFetcher(suite.DB())

	suite.T().Run("Successful fetch of payment requests", func(t *testing.T) {
		expectedPaymentRequests := []models.PaymentRequest{}
		numToMake := 3
		for i := 0; i < numToMake; i++ {
			testdatagen.MakeDefaultPaymentRequest(suite.DB())
		}

		suite.DB().All(&expectedPaymentRequests)

		allPaymentRequests, err := paymentRequestListFetcher.FetchPaymentRequestList()
		suite.NoError(err)
		suite.Equal(len(expectedPaymentRequests), len(*allPaymentRequests))
	})
}
