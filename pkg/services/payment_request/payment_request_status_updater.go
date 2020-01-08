package paymentrequest

import (
	"github.com/gobuffalo/validate"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

type paymentRequestStatusQueryBuilder interface {
	UpdateOne(model interface{}) (*validate.Errors, error)
}

type paymentRequestStatusUpdater struct {
	builder paymentRequestStatusQueryBuilder
}

func NewPaymentRequestStatusUpdater(builder paymentRequestStatusQueryBuilder) services.PaymentRequestStatusUpdater {
	return &paymentRequestStatusUpdater{builder}
}

func (p *paymentRequestStatusUpdater) UpdatePaymentRequestStatus(paymentRequest *models.PaymentRequest) (*validate.Errors, error) {
	verrs, err := p.builder.UpdateOne(paymentRequest)
	return verrs, err
}
