package services

import (
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
)

// CustomerFetcher is the service object interface for FetchCustomer
//go:generate mockery --name CustomerFetcher
type CustomerFetcher interface {
	FetchCustomer(customerID uuid.UUID) (*models.ServiceMember, error)
}

// CustomerUpdater is the service object interface for updating fields of a ServiceMember
//go:generate mockery --name CustomerUpdater
type CustomerUpdater interface {
	UpdateCustomer(eTag string, customer models.ServiceMember) (*models.ServiceMember, error)
}
