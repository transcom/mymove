package paymentrequest

import (
	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/models"
)

func (suite *PaymentRequestServiceSuite) TestFetchReviewedPaymentRequest() {
	reviewedPaymentRequestFetcher := NewPaymentRequestReviewedFetcher()

	setupTestData := func() {
		factory.BuildPaymentRequest(suite.DB(), []factory.Customization{
			{
				Model: models.PaymentRequest{
					Status: models.PaymentRequestStatusReviewed,
				},
			},
		}, nil)
		factory.BuildPaymentRequest(suite.DB(), []factory.Customization{
			{
				Model: models.PaymentRequest{
					Status: models.PaymentRequestStatusPending,
				},
			},
		}, nil)
		factory.BuildPaymentRequest(suite.DB(), []factory.Customization{
			{
				Model: models.PaymentRequest{
					Status: models.PaymentRequestStatusReviewedAllRejected,
				},
			},
		}, nil)
	}

	suite.Run("check for reviewed payment requests", func() {
		setupTestData()
		result, err := reviewedPaymentRequestFetcher.FetchReviewedPaymentRequest(suite.AppContextForTest())
		suite.NoError(err)
		suite.Equal(1, len(result))
	})

}
