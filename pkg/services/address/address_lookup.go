package address

import (
	"context"
	"database/sql"
	"fmt"
	"regexp"
	"strings"

	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/cli"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/featureflag"
)

type vLocation struct {
}

func NewVLocation() services.VLocation {
	return &vLocation{}
}

func (o vLocation) GetLocationsByZipCityState(appCtx appcontext.AppContext, search string) (*models.VLocations, error) {
	locationList, err := FindLocationsByZipCity(appCtx, search)

	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return &locationList, apperror.NewNotFoundError(uuid.Nil, "Search string: "+search)
		default:
			return &locationList, err
		}
	}

	return &locationList, nil
}

// Returns a VLocation array containing all results for the search
// This method expects a comma to be entered after the city name has been entered and is used
// to determine when the state and postal code need to be parsed from the search string
// If there is only one result and no comma and the search string is all numbers we then search
// using the entered postal code rather than city name
func FindLocationsByZipCity(appCtx appcontext.AppContext, search string) (models.VLocations, error) {
	var locationList []models.VLocation
	searchSlice := strings.Split(search, ",")
	city := ""
	state := ""
	postalCode := ""
	var postalCodeRegex = regexp.MustCompile(`^[0-9]+$`)

	if len(searchSlice) > 1 {
		city = searchSlice[0]
		searchSlice = strings.Split(searchSlice[1], " ")
		state = searchSlice[1]

		if len(searchSlice) > 2 {
			postalCode = searchSlice[2]
		}
	} else {
		if postalCodeRegex.MatchString(search) {
			postalCode = strings.TrimSpace(search)
		} else {
			city = search
		}
	}

	/** Feature Flag - Alaska - Determines if AK be included/excluded **/
	isAlaskaEnabled := false
	featureFlagName := "enable_alaska"
	config := cli.GetFliptFetcherConfig(viper.GetViper())
	flagFetcher, err := featureflag.NewFeatureFlagFetcher(config)
	if err != nil {
		appCtx.Logger().Error("Error initializing FeatureFlagFetcher", zap.String("featureFlagKey", featureFlagName), zap.Error(err))
	}

	flag, err := flagFetcher.GetBooleanFlagForUser(context.TODO(), appCtx, featureFlagName, map[string]string{})
	if err != nil {
		appCtx.Logger().Error("Error fetching feature flag", zap.String("featureFlagKey", featureFlagName), zap.Error(err))
	} else {
		isAlaskaEnabled = flag.Match
	}

	isHawaiiEnabled := false
	featureFlagName = "enable_hawaii"
	config = cli.GetFliptFetcherConfig(viper.GetViper())
	flagFetcher, err = featureflag.NewFeatureFlagFetcher(config)
	if err != nil {
		appCtx.Logger().Error("Error initializing FeatureFlagFetcher", zap.String("featureFlagKey", featureFlagName), zap.Error(err))
	}

	flag, err = flagFetcher.GetBooleanFlagForUser(context.TODO(), appCtx, featureFlagName, map[string]string{})
	if err != nil {
		appCtx.Logger().Error("Error fetching feature flag", zap.String("featureFlagKey", featureFlagName), zap.Error(err))
	} else {
		isHawaiiEnabled = flag.Match
	}

	sqlQuery := `SELECT vl.city_name, vl.state, vl.usprc_county_nm, vl.uspr_zip_id, vl.uprc_id
	FROM v_locations vl where vl.uspr_zip_id like ? AND
	vl.city_name like upper(?) AND vl.state like upper(?) `

	if !isAlaskaEnabled {
		sqlQuery += ` AND vl.state NOT in ('AK')`
	}

	if !isHawaiiEnabled {
		sqlQuery += ` AND vl.state NOT in ('HI')`
	}

	sqlQuery += ` limit 30`

	query := appCtx.DB().RawQuery(sqlQuery, fmt.Sprintf("%s%%", postalCode), fmt.Sprintf("%s%%", city), fmt.Sprintf("%s%%", state))
	if err := query.All(&locationList); err != nil {
		if errors.Cause(err).Error() != models.RecordNotFoundErrorString {
			return locationList, err
		}
	}
	for i := range locationList {
		err := appCtx.DB().Load(&locationList[i], "State")
		if err != nil {
			return locationList, err
		}
	}
	return locationList, nil
}
