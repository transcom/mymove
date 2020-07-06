package services

import (
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
)

// CustomerFetcher is the service object interface for FetchCustomer
//go:generate mockery -name CustomerFetcher
type CustomerFetcher interface {
	FetchCustomer(customerID uuid.UUID) (*models.ServiceMember, error)
}
