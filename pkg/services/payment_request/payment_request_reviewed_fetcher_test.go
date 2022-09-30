package paymentrequest

import (
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *PaymentRequestServiceSuite) TestFetchReviewedPaymentRequest() {
	reviewedPaymentRequestFetcher := NewPaymentRequestReviewedFetcher()

	setupTestData := func() {
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
	}

	suite.Run("check for reviewed payment requests", func() {
		setupTestData()
		result, err := reviewedPaymentRequestFetcher.FetchReviewedPaymentRequest(suite.AppContextForTest())
		suite.NoError(err)
		suite.Equal(1, len(result))
	})

}
