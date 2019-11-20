package paymentrequest

import (
	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/validate"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

type paymentRequestCreator struct {
	db *pop.Connection
}

func NewPaymentRequestCreator(db *pop.Connection) services.PaymentRequestCreator {
	return &paymentRequestCreator{db}
}

func (p *paymentRequestCreator) CreatePaymentRequest(pr *models.PaymentRequest) (paymentRequest *models.PaymentRequest, verrs *validate.Errors, err error) {

	verrs, err = p.db.ValidateAndCreate(&paymentRequest)
	if err != nil || verrs.HasAny() {
		return nil, verrs, err
	}

	return paymentRequest, nil, nil
}
