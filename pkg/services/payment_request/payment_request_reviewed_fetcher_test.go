package paymentrequest

import (
	"testing"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *PaymentRequestServiceSuite) TestFetchReviewedPaymentRequest() {
	reviewedPaymentRequestFetcher := NewPaymentRequestReviewedFetcher(suite.DB())

	testdatagen.MakePaymentRequest(suite.DB(), testdatagen.Assertions{
		PaymentRequest: models.PaymentRequest{
			Status: models.PaymentRequestStatusReviewed,
		},
	})
	testdatagen.MakePaymentRequest(suite.DB(), testdatagen.Assertions{
		PaymentRequest: models.PaymentRequest{
			Status: models.PaymentRequestStatusPending,
		},
	})
	testdatagen.MakePaymentRequest(suite.DB(), testdatagen.Assertions{
		PaymentRequest: models.PaymentRequest{
			Status: models.PaymentRequestStatusReviewedAllRejected,
		},
	})

	suite.T().Run("check for reviewed payment requests", func(t *testing.T) {
		result, err := reviewedPaymentRequestFetcher.FetchReviewedPaymentRequest()
		suite.NoError(err)
		suite.Equal(1, len(result))
	})

}
