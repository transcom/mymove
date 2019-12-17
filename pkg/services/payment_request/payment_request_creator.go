package paymentrequest

import (
	"fmt"

	"github.com/gobuffalo/pop"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

type paymentRequestCreator struct {
	db *pop.Connection
}

func NewPaymentRequestCreator(db *pop.Connection) services.PaymentRequestCreator {
	return &paymentRequestCreator{db}
}

func (p *paymentRequestCreator) CreatePaymentRequest(paymentRequest *models.PaymentRequest) (*models.PaymentRequest, error) {
	verrs, err := p.db.ValidateAndCreate(paymentRequest)
	if err != nil {
		return nil, fmt.Errorf("failure creating payment request: %w", err)
	}
	if verrs.HasAny() {
		return nil, fmt.Errorf("validation error saving PaymentRequest: %w", verrs)
	}

	return paymentRequest, err
}
