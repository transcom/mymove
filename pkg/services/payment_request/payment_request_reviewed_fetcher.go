package paymentrequest

import (
	"fmt"

	"github.com/gobuffalo/pop/v5"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

type paymentRequestReviewedFetcher struct {
	db *pop.Connection
}

// NewPaymentRequestReviewedFetcher returns a new payment request fetcher
func NewPaymentRequestReviewedFetcher(db *pop.Connection) services.PaymentRequestReviewedFetcher {
	return &paymentRequestReviewedFetcher{db}
}

//FetchReviewedPaymentRequest finds all payment request with status 'reviewed'
func (p *paymentRequestReviewedFetcher) FetchReviewedPaymentRequest() (models.PaymentRequests, error) {
	var reviewedPaymentRequests models.PaymentRequests
	err := p.db.Q().
		Where("status = ?", models.PaymentRequestStatusReviewed).
		All(&reviewedPaymentRequests)
	if err != nil {
		return reviewedPaymentRequests, services.NewQueryError("PaymentRequests", err, fmt.Sprintf("Could not find reviewed payment requests: %s", err))
	}
	return reviewedPaymentRequests, err
}

const limitOfPRsToProcess int = 100
const lockTimeout string = "1s"

//FetchAndLockReviewedPaymentRequest finds all payment request with status 'reviewed'
func (p *paymentRequestReviewedFetcher) FetchAndLockReviewedPaymentRequest() (models.PaymentRequests, error) {
	var reviewedPaymentRequests models.PaymentRequests
	query := `
		SET LOCAL lock_timeout = $1;
		SELECT * FROM payment_requests
		WHERE status = $2 FOR UPDATE
		LIMIT $3;
	`
	err := p.db.RawQuery(query, lockTimeout, models.PaymentRequestStatusReviewed, limitOfPRsToProcess).All(&reviewedPaymentRequests)

	if err != nil {
		return reviewedPaymentRequests, services.NewQueryError("PaymentRequests", err, fmt.Sprintf("Could not find reviewed payment requests: %s", err))
	}
	return reviewedPaymentRequests, err
}
