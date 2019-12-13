package paymentrequest

import (
	"fmt"

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

func (f *paymentRequestListFetcher) FetchPaymentRequestList() (*models.PaymentRequests, error) {
	paymentRequests := models.PaymentRequests{}

	err := f.db.All(&paymentRequests)
	if err != nil {
		return nil, fmt.Errorf("failure fetching payment requests: %w", err)
	}

	return &paymentRequests, err
}