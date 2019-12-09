package services

import (
	"github.com/gobuffalo/validate"

	"github.com/transcom/mymove/pkg/models"
)

// PaymentRequestCreator is the exported interface for creating a payment request
//go:generate mockery -name PaymentRequestCreator
type PaymentRequestCreator interface {
	CreatePaymentRequest(paymentRequest *models.PaymentRequest) (*models.PaymentRequest, *validate.Errors, error)
}

// PaymentRequestLister is the exported interface for fetching a collection of payment requests
//go:generate mockery -name PaymentRequestLister
type PaymentRequestLister interface {
	ListPaymentRequests() (*[]models.PaymentRequest, *validate.Errors, error)
}
