package paymentrequest

import (
	"github.com/gobuffalo/pop/v5"
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

	// fetch the payment request first with proof of service docs
	// will error if payment request not found
	err := p.db.Eager("PaymentServiceItems", "ProofOfServiceDocs").
		Find(&paymentRequest, paymentRequestID)
	if err != nil {
		return models.PaymentRequest{}, err
	}

	// then fetch the uploads separately to omit soft deleted items
	// empty records are expected
	for index, posd := range paymentRequest.ProofOfServiceDocs {
		var primeUploads models.PrimeUploads
		err = p.db.Q().
			Where("prime_uploads.proof_of_service_docs_id = ? AND prime_uploads.deleted_at IS NULL AND u.deleted_at IS NULL", posd.ID).
			Eager("Upload").
			Join("uploads as u", "u.id = prime_uploads.upload_id").
			All(&primeUploads)
		if err != nil {
			return models.PaymentRequest{}, err
		}

		paymentRequest.ProofOfServiceDocs[index].PrimeUploads = primeUploads
	}

	return paymentRequest, err
}
