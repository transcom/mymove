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
func (f *ediErrorFetcher) FetchEdiErrors(appCtx appcontext.AppContext, pagination services.Pagination) (models.EdiErrors, int, error) {
	var ediErrors models.EdiErrors

	query := appCtx.DB().Q().
		Join("payment_requests", "payment_requests.id = edi_errors.payment_request_id").
		Where("payment_requests.status = ?", models.PaymentRequestStatusEDIError).
		Eager("PaymentRequest").
		Order("edi_errors.created_at DESC")

	paginator := query.Paginate(pagination.Page(), pagination.PerPage())
	err := paginator.All(&ediErrors)
	if err != nil {
		return nil, 0, apperror.NewQueryError("edi_errors", err, "Could not fetch paginated EDI errors")
	}

	count := paginator.Paginator.TotalEntriesSize
	return ediErrors, count, nil
}

// FetchEdiErrorByID returns a single edi_error the edi_error ID for a payment_request with status EDI_ERROR
func (f *ediErrorFetcher) FetchEdiErrorByID(appCtx appcontext.AppContext, id uuid.UUID) (models.EdiError, error) {
	var ediError models.EdiError

	err := appCtx.DB().Q().
		Eager("PaymentRequest").
		Where("id = ?", id).
		First(&ediError)

	if err != nil {
		return models.EdiError{}, apperror.NewNotFoundError(id, "EDIError not found")
	}

	return ediError, nil
}
