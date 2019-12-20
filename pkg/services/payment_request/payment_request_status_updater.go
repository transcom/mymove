package paymentrequest

import (
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

type paymentRequestStatusQueryBuilder interface {
	UpdateOne(model interface{}) error
}

type paymentRequestStatusUpdater struct {
	builder paymentRequestStatusQueryBuilder
}

func NewPaymentRequestStatusUpdater(builder paymentRequestStatusQueryBuilder) services.PaymentRequestStatusUpdater {
	return &paymentRequestStatusUpdater{builder}
}

func (p *paymentRequestStatusUpdater) UpdatePaymentRequestStatus(paymentRequest *models.PaymentRequest) error {
	err := p.builder.UpdateOne(paymentRequest)
	return err
}
