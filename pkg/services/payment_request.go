package services

import (
	"github.com/gobuffalo/validate"

	"github.com/transcom/mymove/pkg/models"
)

// PaymentRequestCreator is the exported interface for creating a payment request
//go:generate mockery -name PaymentRequestCreator
type PaymentRequestCreator interface {
	CreatePaymentRequest(paymentRequest *models.PaymentRequest) (*models.PaymentRequest, error)
}

// PaymentRequestListFetcher is the exported interface for fetching a list of payment requests
//go:generate mockery -name PaymentRequestListFetcher
type PaymentRequestListFetcher interface {
	FetchPaymentRequestList() (*models.PaymentRequests, error)
}

// PaymentRequestFetcher is the exported interface for fetching a payment request
//go:generate mockery -name PaymentRequestFetcher
type PaymentRequestFetcher interface {
	FetchPaymentRequest(filters []QueryFilter) (models.PaymentRequest, error)
}

// PaymentRequestStatusUpdater is the exported interface for updating the status of a payment request
//go:generate mockery -name PaymentRequestStatusUpdater
type PaymentRequestStatusUpdater interface {
	UpdatePaymentRequestStatus(*models.PaymentRequest) (*validate.Errors, error)
}
