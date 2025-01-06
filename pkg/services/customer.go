package services

import (
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
)

// CustomerFetcher is the service object interface for FetchCustomer
//
//go:generate mockery --name CustomerFetcher
type CustomerFetcher interface {
	FetchCustomer(appCtx appcontext.AppContext, customerID uuid.UUID) (*models.ServiceMember, error)
}

// CustomerUpdater is the service object interface for updating fields of a ServiceMember
//
//go:generate mockery --name CustomerUpdater
type CustomerUpdater interface {
	UpdateCustomer(appCtx appcontext.AppContext, eTag string, customer models.ServiceMember) (*models.ServiceMember, error)
}

//go:generate mockery --name CustomerSearcher
type CustomerSearcher interface {
	SearchCustomers(appCtx appcontext.AppContext, params *SearchCustomersParams) (models.ServiceMemberSearchResults, int, error)
}

type SearchCustomersParams struct {
	Edipi         *string
	Emplid        *string
	Branch        *string
	CustomerName  *string
	PersonalEmail *string
	Telephone     *string
	Page          int64
	PerPage       int64
	Sort          *string
	Order         *string
}
