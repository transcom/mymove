package paymentrequest

import (
	"testing"
	"time"

	"github.com/transcom/mymove/pkg/etag"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/query"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *PaymentRequestServiceSuite) TestUpdatePaymentRequestStatus() {
	builder := query.NewQueryBuilder(suite.DB())

	suite.T().Run("If we get a payment request pointer with a status it should update and return no error", func(t *testing.T) {
		paymentRequest := testdatagen.MakeDefaultPaymentRequest(suite.DB())
		paymentRequest.Status = models.PaymentRequestStatusReviewed

		updater := NewPaymentRequestStatusUpdater(builder)

		_, err := updater.UpdatePaymentRequestStatus(&paymentRequest, etag.GenerateEtag(paymentRequest.UpdatedAt))
		suite.NoError(err)
	})

	suite.T().Run("Should return a PreconditionFailedError with a stale etag", func(t *testing.T) {
		paymentRequest := testdatagen.MakeDefaultPaymentRequest(suite.DB())
		paymentRequest.Status = models.PaymentRequestStatusReviewed

		updater := NewPaymentRequestStatusUpdater(builder)

		_, err := updater.UpdatePaymentRequestStatus(&paymentRequest, etag.GenerateEtag(time.Now()))
		suite.Error(err)
		suite.IsType(services.PreconditionFailedError{}, err)
	})

}
