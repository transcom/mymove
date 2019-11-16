package services

import (
	"github.com/gobuffalo/validate"

	"github.com/transcom/mymove/pkg/models"
)

// PaymentRequestCreator is the exported interface for fetching multiple transportation offices
//go:generate mockery -name PaymentRequestCreator
type PaymentRequestCreator interface {
	CreatePaymentRequest(paymentRequest *models.PaymentRequest, moveTaskOrderIDFilter []QueryFilter) (*models.PaymentRequest, *validate.Errors, error)
}