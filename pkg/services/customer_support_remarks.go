package services

import (
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	customersupportremarksop "github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations/customer_support_remarks"
	"github.com/transcom/mymove/pkg/models"
)

// CustomerSupportRemarksFetcher is the exported interface for fetching office remarks for a move.
//
//go:generate mockery --name CustomerSupportRemarksFetcher
type CustomerSupportRemarksFetcher interface {
	ListCustomerSupportRemarks(appCtx appcontext.AppContext, moveCode string) (*models.CustomerSupportRemarks, error)
}

//go:generate mockery --name CustomerSupportRemarksCreator
type CustomerSupportRemarksCreator interface {
	CreateCustomerSupportRemark(appCtx appcontext.AppContext, customerSupportRemark *models.CustomerSupportRemark, moveCode string) (*models.CustomerSupportRemark, error)
}

//go:generate mockery --name CustomerSupportRemarkUpdater
type CustomerSupportRemarkUpdater interface {
	UpdateCustomerSupportRemark(appCtx appcontext.AppContext, params customersupportremarksop.UpdateCustomerSupportRemarkForMoveParams) (*models.CustomerSupportRemark, error)
}

//go:generate mockery --name CustomerSupportRemarkDeleter
type CustomerSupportRemarkDeleter interface {
	DeleteCustomerSupportRemark(appCtx appcontext.AppContext, customerSupportRemarkID uuid.UUID) error
}
