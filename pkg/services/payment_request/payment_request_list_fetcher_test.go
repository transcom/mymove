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

		for i := 0; i < 3; i++ {
			testdatagen.MakeDefaultPaymentRequest(suite.DB())
		}

		allPaymentRequests, err := paymentRequestListFetcher.FetchPaymentRequestList()
		expectedLength := numPaymentRequests + 3
		suite.NoError(err)
		suite.Equal(expectedLength, len(allPaymentRequests))
	})
}
