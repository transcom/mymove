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

func FindLocationsByZipCity(appCtx appcontext.AppContext, search string) (models.VLocations, error) {
	var locationList []models.VLocation
	search = strings.ReplaceAll(search, ",", "") //remove any commas so they are not used in the search
	searchSlice := strings.Split(search, " ")
	city := ""
	state := ""
	postalCode := ""

	var postalCodeRegex = regexp.MustCompile(`^[0-9]+$`)

	if len(searchSlice) == 1 {
		// check if this is a zip only search
		if postalCodeRegex.MatchString(search) {
			postalCode = search
		} else {
			city = search
		}
	} else if postalCode == "" && len(searchSlice) == 2 {
		city = strings.TrimSpace(searchSlice[0])
		state = strings.TrimSpace(searchSlice[1])
	} else if len(searchSlice) == 3 {
		if postalCodeRegex.MatchString(searchSlice[2]) {
			postalCode = strings.TrimSpace(searchSlice[2])
		}
		city = strings.TrimSpace(searchSlice[0])
		state = strings.TrimSpace(searchSlice[1])
	}

	// user may have typed a comma as part of the city name we need to remove that do to the query
	if city != "" {
		city = strings.ReplaceAll(city, ",", "")
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

	sqlQuery := ""

	if isAlaskaEnabled {
		sqlQuery = fmt.Sprintf(`
			select vl.city_name, vl.state, vl.usprc_county_nm, vl.uspr_zip_id, vl.uprc_id
				from v_locations vl where vl.uspr_zip_id like '%[1]s%%' and
				vl.city_name like upper('%[2]s%%') and vl.state like upper('%[3]s%%') limit 30`, postalCode, city, state)
	} else {
		sqlQuery = fmt.Sprintf(`
		select vl.city_name, vl.state, vl.usprc_county_nm, vl.uspr_zip_id, vl.uprc_id
			from v_locations vl where vl.uspr_zip_id like '%[1]s%%' and
			vl.city_name like upper('%[2]s%%') and vl.state like upper('%[3]s%%') and vl.state NOT in ('AK','HI') limit 30`, postalCode, city, state)
	}

	query := appCtx.DB().Q().RawQuery(sqlQuery)
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