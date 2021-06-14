package paymentrequest

import (
	"testing"
	"time"

	"github.com/transcom/mymove/pkg/etag"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/query"
	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/unit"
)

func (suite *PaymentRequestServiceSuite) TestUpdateProcessedPaymentRequestStatus() {
	builder := query.NewQueryBuilder(suite.DB())

	suite.T().Run("If we get a payment request pointer with a valid status it should update and return no error", func(t *testing.T) {
		// Payment request being updated with the UpdateProcessedPaymentRequestStatus can only be:
		// "SENT_TO_GEX", "RECEIVED_BY_GEX", "PAID", or "EDI_ERROR"
		approvedPRStatuses := [4]models.PaymentRequestStatus{"SENT_TO_GEX", "RECEIVED_BY_GEX", "PAID", "EDI_ERROR"}
		for _, approvedPRStatus := range approvedPRStatuses {
			paymentRequest := testdatagen.MakeDefaultPaymentRequest(suite.DB())

			paymentRequest.Status = approvedPRStatus

			updater := NewPaymentRequestStatusUpdater(builder)

			_, err := updater.UpdateProcessedPaymentRequestStatus(&paymentRequest, etag.GenerateEtag(paymentRequest.UpdatedAt))
			suite.NoError(err)
		}
	})

	suite.T().Run("Should return a PreconditionFailedError with a stale etag", func(t *testing.T) {
		paymentRequest := testdatagen.MakeDefaultPaymentRequest(suite.DB())
		paymentRequest.Status = models.PaymentRequestStatusSentToGex

		updater := NewPaymentRequestStatusUpdater(builder)

		_, err := updater.UpdateProcessedPaymentRequestStatus(&paymentRequest, etag.GenerateEtag(time.Now()))
		suite.Error(err)
		suite.IsType(services.PreconditionFailedError{}, err)
	})

}

func (suite *PaymentRequestServiceSuite) TestUpdateReviewedPaymentRequestStatus() {
	builder := query.NewQueryBuilder(suite.DB())

	suite.T().Run("If we get a payment request pointer with a status it should update and return no error", func(t *testing.T) {
		paymentRequest := testdatagen.MakeDefaultPaymentRequest(suite.DB())

		paymentRequest.Status = models.PaymentRequestStatusReviewed

		updater := NewPaymentRequestStatusUpdater(builder)

		_, err := updater.UpdateReviewedPaymentRequestStatus(&paymentRequest, etag.GenerateEtag(paymentRequest.UpdatedAt))
		suite.NoError(err)
	})

	suite.T().Run("Should return a ConflictError if the payment request has any service items that have not been reviewed", func(t *testing.T) {
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

		_, err := updater.UpdateReviewedPaymentRequestStatus(&paymentRequest, etag.GenerateEtag(paymentRequest.UpdatedAt))
		suite.Error(err)
		suite.IsType(services.ConflictError{}, err)
	})

	suite.T().Run("Should update and return no error if the payment request has service items that have all been reviewed", func(t *testing.T) {
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

		_, err := updater.UpdateReviewedPaymentRequestStatus(&paymentRequest, etag.GenerateEtag(paymentRequest.UpdatedAt))
		suite.NoError(err)
	})

	suite.T().Run("Should return a PreconditionFailedError with a stale etag", func(t *testing.T) {
		paymentRequest := testdatagen.MakeDefaultPaymentRequest(suite.DB())
		paymentRequest.Status = models.PaymentRequestStatusReviewed

		updater := NewPaymentRequestStatusUpdater(builder)

		_, err := updater.UpdateReviewedPaymentRequestStatus(&paymentRequest, etag.GenerateEtag(time.Now()))
		suite.Error(err)
		suite.IsType(services.PreconditionFailedError{}, err)
	})

	suite.T().Run("Should return an InvalidInput error with a wrong status", func(t *testing.T) {
		// Payment request being updated with the UpdateReviewedPaymentRequestStatus can only be:
		// REVIEWED or REVIEWED_AND_ALL_SERVICE_ITEMS_REJECTED
		nonApprovedPRStatuses := [5]models.PaymentRequestStatus{"SENT_TO_GEX", "RECEIVED_BY_GEX", "PAID", "EDI_ERROR", "PENDING"}
		for _, nonApprovedPRStatus := range nonApprovedPRStatuses {
			paymentRequest := testdatagen.MakeDefaultPaymentRequest(suite.DB())

			paymentRequest.Status = nonApprovedPRStatus

			updater := NewPaymentRequestStatusUpdater(builder)

			_, err := updater.UpdateReviewedPaymentRequestStatus(&paymentRequest, etag.GenerateEtag(time.Now()))
			suite.Error(err)
			suite.IsType(services.InvalidInputError{}, err)
		}
	})
}
