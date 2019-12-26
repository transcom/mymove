package services

import (
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/uploader"

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
type PaymentRequestFetcher interface {
	FetchPaymentRequest(paymentRequestID uuid.UUID) (*models.PaymentRequest, error)
}

// PaymentRequestUploadCreator is the exported interface for creating a payment request upload
type PaymentRequestUploadCreator interface {
	CreateUpload(file uploader.File, paymentRequestID uuid.UUID, userID uuid.UUID) (*models.Upload, error)
}
