package address

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

func (o countrySearcher) SearchCountries(appCtx appcontext.AppContext, queryFilter *string) (models.Countries, error) {
	var countries models.Countries

	query := ""

	if queryFilter != nil {
		query = strings.ToUpper(strings.TrimSpace(*queryFilter))
	}

	// If no query is provided return all.
	if len(query) == 0 {
		err := appCtx.DB().Order("country_name asc").All(&countries)
		if err != nil {
			return countries, err
		}

		return countries, nil
	}

	startsWithFilter := fmt.Sprintf("%s%%", query)

	// If query len is 2 chars: match on country code or starts with match on country name.
	if len(query) == 2 {
		sql := `SELECT * FROM re_countries where country = ?
                union
                SELECT * FROM re_countries where country_name ILIKE ?`
		err := appCtx.DB().RawQuery(sql, query, startsWithFilter).All(&countries)
		if err != nil {
			return countries, err
		}

		return countries, nil
	}

	// If query len is not 2 chars do starts with match on country name.
	err := appCtx.DB().Where("country_name ILIKE ?", startsWithFilter).All(&countries)
	if err != nil {
		return countries, err
	}

	return countries, nil
}
