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

// CountrySearcher is the exported interface for searching a country
//
//go:generate mockery --name CountrySearcher
type CountrySearcher interface {
	SearchCountries(appCtx appcontext.AppContext, queryFilter *string) (models.Countries, error)
}
