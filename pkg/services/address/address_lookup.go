package address

import (
	"database/sql"

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

	// The % operator filters out strings that are below this similarity threshold
	err := appCtx.DB().Q().RawQuery("SET pg_trgm.similarity_threshold = 0.03").Exec()
	if err != nil {
		return locationList, err
	}

	// sqlQuery := `
	// 	with locations as (select office.id as transportation_office_id, office.name, similarity(office.name, $1) as sim
	//     from transportation_offices as office
	//     where name % $1 `

	// sqlQuery += `
	// 	order by sim desc
	//     limit 5)
	// 	select office.*
	//     from names n inner join transportation_offices office on n.transportation_office_id = office.id
	//     group by office.id
	//     order by max(n.sim) desc, office.name
	//     limit 5`
	sqlQuery := `
		select uprc.u_s_post_region_city_nm, uprc.state, uprc.usprc_county_nm, uprc.uspr_zip_id
			from us_post_region_cities uprc where position(upper($1) in uprc.uspr_zip_id) > 0 or
			position(upper($1) in uprc.u_s_post_region_city_nm) > 0
			limit 5`
	query := appCtx.DB().Q().RawQuery(sqlQuery, search)
	if err := query.All(&locationList); err != nil {
		if errors.Cause(err).Error() != models.RecordNotFoundErrorString {
			return locationList, err
		}
	}
	for i := range locationList {
		err := appCtx.DB().Load(&locationList[i], "Address")
		if err != nil {
			return locationList, err
		}
	}
	return locationList, nil
}
