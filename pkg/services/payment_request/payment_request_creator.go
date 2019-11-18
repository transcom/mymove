package payment_request

import (
	"github.com/gobuffalo/pop"
	"github.com/transcom/mymove/pkg/services"
)

type createPaymentRequest struct {
	db *pop.Connection
}

func NewCreatePaymentRequest(db *pop.Connection) services.PaymentRequestCreator {
	return &createPaymentRequest{db}
}

func