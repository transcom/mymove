package services

import (
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/gen/ghcmessages"
	"github.com/transcom/mymove/pkg/models"
)

// CustomerSupportRemarksFetcher is the exported interface for fetching office remarks for a move.
//go:generate mockery --name CustomerSupportRemarksFetcher --disable-version-string
type CustomerSupportRemarksFetcher interface {
	ListCustomerSupportRemarks(appCtx appcontext.AppContext, moveCode string) (*models.CustomerSupportRemarks, error)
}

//go:generate mockery --name CustomerSupportRemarksCreator --disable-version-string
type CustomerSupportRemarksCreator interface {
	CreateCustomerSupportRemark(appCtx appcontext.AppContext, customerSupportRemark *models.CustomerSupportRemark, moveCode string) (*models.CustomerSupportRemark, error)
}

//go:generate mockery --name CustomerSupportRemarkUpdater --disable-version-string
type CustomerSupportRemarkUpdater interface {
	UpdateCustomerSupportRemark(appCtx appcontext.AppContext, payload ghcmessages.UpdateCustomerSupportRemarkPayload) (*models.CustomerSupportRemark, error)
}

//go:generate mockery --name CustomerSupportRemarkDeleter --disable-version-string
type CustomerSupportRemarkDeleter interface {
	DeleteCustomerSupportRemark(appCtx appcontext.AppContext, customerSupportRemarkID uuid.UUID) (uuid.UUID, error)
}
