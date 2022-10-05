package paymentrequest

import (
	"time"

	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/etag"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services/query"
	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/unit"
)

func (suite *PaymentRequestServiceSuite) TestUpdatePaymentRequestStatus() {
	builder := query.NewQueryBuilder()

	suite.Run("If we get a payment request pointer with a status it should update and return no error", func() {
		paymentRequest := testdatagen.MakeDefaultPaymentRequest(suite.DB())
		paymentRequest.Status = models.PaymentRequestStatusReviewed

		updater := NewPaymentRequestStatusUpdater(builder)

		_, err := updater.UpdatePaymentRequestStatus(suite.AppContextForTest(), &paymentRequest, etag.GenerateEtag(paymentRequest.UpdatedAt))
		suite.NoError(err)
	})

	suite.Run("Should return a ConflictError if the payment request has any service items that have not been reviewed", func() {
		paymentRequest := testdatagen.MakeDefaultPaymentRequest(suite.DB())

		psiCost := unit.Cents(10000)
		testdatagen.MakePaymentServiceItem(suite.DB(), testdatagen.Assertions{
			PaymentServiceItem: models.PaymentServiceItem{
				PriceCents: &psiCost,
				Status:     models.PaymentServiceItemStatusRequested,
			},
			PaymentRequest: paymentRequest,
		})
		testdatagen.MakePaymentServiceItem(suite.DB(), testdatagen.Assertions{
			PaymentServiceItem: models.PaymentServiceItem{
				PriceCents: &psiCost,
				Status:     models.PaymentServiceItemStatusApproved,
			},
			PaymentRequest: paymentRequest,
		})

		paymentRequest.Status = models.PaymentRequestStatusReviewed
		updater := NewPaymentRequestStatusUpdater(builder)

		_, err := updater.UpdatePaymentRequestStatus(suite.AppContextForTest(), &paymentRequest, etag.GenerateEtag(paymentRequest.UpdatedAt))
		suite.Error(err)
		suite.IsType(apperror.ConflictError{}, err)
	})

	suite.Run("Should update and return no error if the payment request has service items that have all been reviewed", func() {
		paymentRequest := testdatagen.MakeDefaultPaymentRequest(suite.DB())

		psiCost := unit.Cents(10000)
		testdatagen.MakePaymentServiceItem(suite.DB(), testdatagen.Assertions{
			PaymentServiceItem: models.PaymentServiceItem{
				PriceCents: &psiCost,
				Status:     models.PaymentServiceItemStatusApproved,
			},
			PaymentRequest: paymentRequest,
		})
		testdatagen.MakePaymentServiceItem(suite.DB(), testdatagen.Assertions{
			PaymentServiceItem: models.PaymentServiceItem{
				PriceCents: &psiCost,
				Status:     models.PaymentServiceItemStatusDenied,
			},
			PaymentRequest: paymentRequest,
		})

		paymentRequest.Status = models.PaymentRequestStatusReviewed
		updater := NewPaymentRequestStatusUpdater(builder)

		_, err := updater.UpdatePaymentRequestStatus(suite.AppContextForTest(), &paymentRequest, etag.GenerateEtag(paymentRequest.UpdatedAt))
		suite.NoError(err)
	})

	suite.Run("Should return a PreconditionFailedError with a stale etag", func() {
		paymentRequest := testdatagen.MakeDefaultPaymentRequest(suite.DB())
		paymentRequest.Status = models.PaymentRequestStatusReviewed

		updater := NewPaymentRequestStatusUpdater(builder)

		_, err := updater.UpdatePaymentRequestStatus(suite.AppContextForTest(), &paymentRequest, etag.GenerateEtag(time.Now()))
		suite.Error(err)
		suite.IsType(apperror.PreconditionFailedError{}, err)
	})

}
