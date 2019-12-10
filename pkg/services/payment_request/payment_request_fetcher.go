package paymentrequest

import (
	"github.com/go-openapi/swag"
	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/validate"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/testdatagen"
)

type paymentRequestFetcher struct {
	db *pop.Connection
}

func NewPaymentRequestFetcher(db *pop.Connection) services.PaymentRequestFetcher {
	return &paymentRequestFetcher{db}
}

func (p *paymentRequestFetcher) FetchPaymentRequest(paymentRequestID uuid.UUID) (*models.PaymentRequest, *validate.Errors, error) {
	// A mock payment request. This is temporary and will be replaced with real data eventually.
	mockPaymentRequest := models.PaymentRequest{
		ID:              paymentRequestID,
		IsFinal:         false,
		RejectionReason: swag.String(""),
		CreatedAt:       testdatagen.PeakRateCycleStart,
		UpdatedAt:       testdatagen.PeakRateCycleStart,
	}

	return &mockPaymentRequest, nil, nil
}
