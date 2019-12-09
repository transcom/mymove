package paymentrequest

import (
	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/validate"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/testdatagen"
)

type paymentRequestLister struct {
	db *pop.Connection
}

func NewPaymentRequestLister(db *pop.Connection) services.PaymentRequestLister {
	return &paymentRequestLister{db}
}

func (p *paymentRequestLister) ListPaymentRequests() (*[]models.PaymentRequest, *validate.Errors, error) {
	// A small collection of mock payment requests. This is temporary and will be replaced with real data eventually.
	uuid1, _ := uuid.NewV4()
	uuid2, _ := uuid.NewV4()
	mockPaymentRequests := []models.PaymentRequest{
		{
			ID:              uuid1,
			IsFinal:         false,
			RejectionReason: "",
			CreatedAt:       testdatagen.PeakRateCycleStart,
			UpdatedAt:       testdatagen.PeakRateCycleStart,
		},
		{
			ID:              uuid2,
			IsFinal:         false,
			RejectionReason: "",
			CreatedAt:       testdatagen.DateInsidePerformancePeriod,
			UpdatedAt:       testdatagen.DateInsidePerformancePeriod,
		},
	}

	return &mockPaymentRequests, nil, nil
}
