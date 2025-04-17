package services

import (
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
)

// EDIErrorFetcher is the exported interface for fetching edi_errors for payment requests in 'EDI_ERROR' status
//
//go:generate mockery --name EDIErrorFetcher
type EDIErrorFetcher interface {
	FetchEdiErrors(appCtx appcontext.AppContext, pagination Pagination) (models.EdiErrors, int, error)
	FetchEdiErrorByID(appCtx appcontext.AppContext, id uuid.UUID) (models.EdiError, error)
}
