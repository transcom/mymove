package paymentrequest

import (
	"testing"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *PaymentRequestServiceSuite) TestPaymentRequestList() {
	paymentRequestListFetcher := NewPaymentRequestListFetcher(suite.DB())

	suite.T().Run("Successful fetch of payment requests", func(t *testing.T) {
		paymentRequests := []models.PaymentRequest{}
		_ := suite.DB().All(&paymentRequests)
		numPaymentRequests := len(paymentRequests)

		numToMake := 3
		for i := 0; i < numToMake; i++ {
			testdatagen.MakeDefaultPaymentRequest(suite.DB())
		}

		allPaymentRequests, err := paymentRequestListFetcher.FetchPaymentRequestList()
		expectedLength := numPaymentRequests + numToMake
		suite.NoError(err)
		suite.Equal(expectedLength, len(allPaymentRequests))
	})
}
