package paymentrequest

import (
	"testing"

	"github.com/transcom/mymove/pkg/testdatagen"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
)

func (suite *PaymentRequestServiceSuite) TestFetchPaymentRequest() {
	suite.T().Run("If a payment request is fetched, it should be returned", func(t *testing.T) {

		fetcher := NewPaymentRequestFetcher(suite.DB())

		pr := testdatagen.MakePaymentRequest(suite.DB(), testdatagen.Assertions{})
		paymentRequest, err := fetcher.FetchPaymentRequest(pr.ID)

		suite.NoError(err)
		suite.Equal(pr.ID, paymentRequest.ID)
	})

	suite.T().Run("if there is an error, we get it with zero payment request", func(t *testing.T) {
		fetcher := NewPaymentRequestFetcher(suite.DB())

		paymentRequest, err := fetcher.FetchPaymentRequest(uuid.Nil)

		suite.Error(err)
		suite.Equal(err.Error(), "sql: no rows in result set")
		suite.Equal(models.PaymentRequest{}, paymentRequest)
	})
}
