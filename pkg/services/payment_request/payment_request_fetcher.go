package paymentrequest

import (
	"github.com/gobuffalo/pop"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

type paymentRequestFetcher struct {
	db *pop.Connection
}

// NewPaymentRequestFetcher returns a new payment request fetcher
func NewPaymentRequestFetcher(db *pop.Connection) services.PaymentRequestFetcher {
	return &paymentRequestFetcher{db}
}

//FetchPaymentRequest finds the payment request by id
func (p *paymentRequestFetcher) FetchPaymentRequest(paymentRequestID uuid.UUID) (models.PaymentRequest, error) {
	var paymentRequest models.PaymentRequest
	err := p.db.Eager("PaymentServiceItems", "ProofOfServiceDocs.PrimeUploads.Upload").Find(&paymentRequest, paymentRequestID)
	if err != nil {
		return models.PaymentRequest{}, err
	}

	return paymentRequest, err
}
