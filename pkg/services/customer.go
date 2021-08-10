package services

import (
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appconfig"
	"github.com/transcom/mymove/pkg/models"
)

// CustomerFetcher is the service object interface for FetchCustomer
//go:generate mockery --name CustomerFetcher --disable-version-string
type CustomerFetcher interface {
	FetchCustomer(appCfg appconfig.AppConfig, customerID uuid.UUID) (*models.ServiceMember, error)
}

// CustomerUpdater is the service object interface for updating fields of a ServiceMember
//go:generate mockery --name CustomerUpdater --disable-version-string
type CustomerUpdater interface {
	UpdateCustomer(appCfg appconfig.AppConfig, eTag string, customer models.ServiceMember) (*models.ServiceMember, error)
}
