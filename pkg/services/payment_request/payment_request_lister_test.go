package paymentrequest

import (
	"testing"
)

func (suite *PaymentRequestServiceSuite) TestListPaymentRequests() {

	lister := NewPaymentRequestLister(suite.DB())

	// Happy path
	suite.T().Run("Payment request mock data is returned", func(t *testing.T) {
		paymentRequests, _, _ := lister.ListPaymentRequests()
		suite.Equal(2, len(*paymentRequests))
	})
}
