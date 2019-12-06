package paymentrequest

import (
	"github.com/gobuffalo/pop"

	"github.com/transcom/mymove/pkg/services"

	"github.com/transcom/mymove/pkg/models"
)

type paymentRequestListFetcher struct {
	db *pop.Connection
}

func NewPaymentRequestListFetcher(db *pop.Connection) services.PaymentRequestListFetcher {
	return &paymentRequestListFetcher{db}
}

func (f *paymentRequestListFetcher) FetchPaymentRequestList() ([]models.PaymentRequest, error) {
	paymentRequests := []models.PaymentRequest{}

	err := f.db.All(&paymentRequests)
	if err != nil {
		return paymentRequests, err
	}

	return paymentRequests, err
}