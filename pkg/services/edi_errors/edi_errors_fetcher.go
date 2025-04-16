package edi_errors

import (
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

type ediErrorFetcher struct{}

// NewEDIErrorFetcher returns an instance that implements the EDIErrorFetcher interface
func NewEDIErrorFetcher() services.EDIErrorFetcher {
	return &ediErrorFetcher{}
}

// FetchEdiErrors returns all edi_errors related to payment requests with status EDI_ERROR
func (f *ediErrorFetcher) FetchEdiErrors(appCtx appcontext.AppContext) (models.EdiErrors, error) {
	var ediErrorPaymentRequests models.PaymentRequests

	err := appCtx.DB().Q().
		Where("status = ?", models.PaymentRequestStatusEDIError).
		All(&ediErrorPaymentRequests)

	if err != nil {
		return models.EdiErrors{}, apperror.NewQueryError("payment_requests", err, "Could not find payment requests with EDI_ERROR status")
	}

	var ediErrorPaymentRequestIds []uuid.UUID
	for _, pr := range ediErrorPaymentRequests {
		ediErrorPaymentRequestIds = append(ediErrorPaymentRequestIds, pr.ID)
	}

	if len(ediErrorPaymentRequestIds) == 0 {
		return models.EdiErrors{}, nil
	}

	var ediErrors models.EdiErrors
	err = appCtx.DB().Q().
		Where("payment_request_id IN (?)", ediErrorPaymentRequestIds).
		Eager("PaymentRequest").
		All(&ediErrors)

	if err != nil {
		return models.EdiErrors{}, apperror.NewQueryError("edi_errors", err, "Could not find EDI error details for payment requests in EDI_ERROR status")
	}

	return ediErrors, nil
}
