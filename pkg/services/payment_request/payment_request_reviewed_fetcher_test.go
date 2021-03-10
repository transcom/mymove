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

func (suite *PaymentRequestServiceSuite) TestFetchAndLockReviewedPaymentRequest() {
	reviewedPaymentRequestFetcher := NewPaymentRequestReviewedFetcher(suite.DB())

	_ = suite.createPaymentRequest(100)

	suite.T().Run("successfully fetch given reviewed payment requests", func(t *testing.T) {
		result, err := reviewedPaymentRequestFetcher.FetchAndLockReviewedPaymentRequest()
		suite.NoError(err)
		suite.Equal(100, len(result))
	})

	// suite.T().Run("throw an error if a locked payment request is updated", func(t *testing.T) {
	// 	errTransaction := suite.DB().Transaction(func(tx *pop.Connection) error {
	// 		_, err := reviewedPaymentRequestFetcher.FetchAndLockReviewedPaymentRequest()
	// 		suite.NoError(err)
	// 		return err
	// 	})
	// 	suite.DB().Transaction(func(tx *pop.Connection) error {
	// 		err := suite.DB().RawQuery(`UPDATE payment_requests SET status = $1 WHERE id = $2;`, models.PaymentRequestStatusPaid, prs[99].ID).Exec()

	// 		suite.NoError(err)
	// 		return err
	// 	})
	// 	suite.Error(errTransaction)
	// })
	_ = suite.createPaymentRequest(101)
	var paymentRequests models.PaymentRequests
	suite.DB().All(&paymentRequests)
	suite.T().Run("retrieve only the number of payment requests set by limitOfPRsToProcess", func(t *testing.T) {
		result, err := reviewedPaymentRequestFetcher.FetchAndLockReviewedPaymentRequest()
		suite.NoError(err)
		suite.Equal(limitOfPRsToProcess, len(result))
		suite.Equal(201, len(paymentRequests))
	})
}
