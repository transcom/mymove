package services

import (
	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
)

type AddressCreator interface {
	CreateAddress(appCtx appcontext.AppContext, address *models.Address) (*models.Address, error)
}

type AddressUpdater interface {
	UpdateAddress(appCtx appcontext.AppContext, address *models.Address, eTag string) (*models.Address, error)
}

//go:generate mockery --name VLocation
type VLocation interface {
	GetLocationsByZipCityState(appCtx appcontext.AppContext, search string, exclusionStateFilters []string, exactMatch ...bool) (*models.VLocations, error)
}

//go:generate mockery --name VIntlLocation
type VIntlLocation interface {
	GetOconusLocations(appCtx appcontext.AppContext, country string, search string, exactMatch bool) (*models.VIntlLocations, error)
}
