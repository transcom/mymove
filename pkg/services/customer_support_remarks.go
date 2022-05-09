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
