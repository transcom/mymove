package country

import (
	"fmt"
	"strings"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

type countrySearcher struct {
}

func NewCountrySearcher() services.CountrySearcher {
	return &countrySearcher{}
}

func (o countrySearcher) SearchCountries(appCtx appcontext.AppContext, searchQuery *string) (models.Countries, error) {
	var countries models.Countries

	// If searchQuery is nil, return ALL
	if searchQuery == nil {
		err := appCtx.DB().Order("country_name asc").All(&countries)
		if err != nil {
			return countries, err
		}

		return countries, nil
	}

	// If searchQuery len is 2 chars, search against country(code)
	if len(*searchQuery) == 2 {
		err := appCtx.DB().Where("country = ?", strings.ToUpper(*searchQuery)).All(&countries)
		if err != nil {
			return countries, err
		}

		return countries, nil
	}

	// If searchQuery len is greater than 2 chars search for partial match against country name
	partialSearch := fmt.Sprintf("%%%s%%", *searchQuery)
	err := appCtx.DB().Where("UPPER(country_name) ILIKE ?", strings.ToUpper(partialSearch)).All(&countries)
	if err != nil {
		return countries, err
	}

	return countries, nil
}
