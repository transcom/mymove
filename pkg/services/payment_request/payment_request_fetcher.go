package paymentrequest

import (
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

type paymentRequestQueryBuilder interface {
	FetchOne(model interface{}, filters []services.QueryFilter) error
}

type paymentRequestFetcher struct {
	builder paymentRequestQueryBuilder
}

func NewPaymentRequestFetcher(builder paymentRequestQueryBuilder) services.PaymentRequestFetcher {
	return &paymentRequestFetcher{builder}
}

func (p *paymentRequestFetcher) FetchPaymentRequest(filters []services.QueryFilter) (models.PaymentRequest, error) {
	var paymentRequest models.PaymentRequest
	err := p.builder.FetchOne(&paymentRequest, filters)
	return paymentRequest, err
}
