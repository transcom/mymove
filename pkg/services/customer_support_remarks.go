package services

import (
	"github.com/transcom/mymove/pkg/appcontext"
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
