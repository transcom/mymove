package services

import (
	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
)

// CountrySearcher is the exported interface for searching a country
//
//go:generate mockery --name CountrySearcher
type CountrySearcher interface {
	SearchCountries(appCtx appcontext.AppContext, searchQuery *string) (models.Countries, error)
}
