package address

import (
	"database/sql"
	"fmt"

	"github.com/gofrs/uuid"
	"github.com/pkg/errors"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

type usPostRegionCity struct {
}

func NewUsPostRegionCity() services.UsPostRegionCity {
	return &usPostRegionCity{}
}

func (o usPostRegionCity) GetLocationsByZipCity(appCtx appcontext.AppContext, search string) (*models.UsPostRegionCities, error) {
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

func FindLocationsByZipCity(appCtx appcontext.AppContext, search string) (models.UsPostRegionCities, error) {
	var locationList []models.UsPostRegionCity
	sqlQuery := fmt.Sprintf(`
		select uprc.u_s_post_region_city_nm, uprc.state, uprc.usprc_county_nm, uprc.uspr_zip_id
			from us_post_region_cities uprc where uprc.uspr_zip_id like '%[1]s%%' or
			uprc.u_s_post_region_city_nm like upper('%[1]s%%') limit 30`, search)
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
