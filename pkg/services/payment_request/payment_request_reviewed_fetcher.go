package paymentrequest

import (
	"fmt"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

type paymentRequestReviewedFetcher struct {
}

// NewPaymentRequestReviewedFetcher returns a new payment request fetcher
func NewPaymentRequestReviewedFetcher() services.PaymentRequestReviewedFetcher {
	return &paymentRequestReviewedFetcher{}
}

// FetchReviewedPaymentRequest finds all payment request with status 'reviewed'
func (p *paymentRequestReviewedFetcher) FetchReviewedPaymentRequest(appCtx appcontext.AppContext) (models.PaymentRequests, error) {
	var reviewedPaymentRequests models.PaymentRequests
	err := appCtx.DB().Q().
		Where("status = ?", models.PaymentRequestStatusReviewed).
		All(&reviewedPaymentRequests)
	if err != nil {
		return reviewedPaymentRequests, apperror.NewQueryError("PaymentRequests", err, fmt.Sprintf("Could not find reviewed payment requests: %s", err))
	}
	return reviewedPaymentRequests, err
}
