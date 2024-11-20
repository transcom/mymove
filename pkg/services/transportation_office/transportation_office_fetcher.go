package transportationoffice

import (
	"database/sql"

	"github.com/gofrs/uuid"
	"github.com/pkg/errors"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

type transportationOfficesFetcher struct {
}

func NewTransportationOfficesFetcher() services.TransportationOfficesFetcher {
	return &transportationOfficesFetcher{}
}

func (o transportationOfficesFetcher) GetTransportationOffice(appCtx appcontext.AppContext, transportationOfficeID uuid.UUID, includeOnlyPPMCloseoutOffices bool) (*models.TransportationOffice, error) {
	var transportationOffice models.TransportationOffice
	err := appCtx.DB().EagerPreload("Address", "Address.Country").
		Where("provides_ppm_closeout = ?", includeOnlyPPMCloseoutOffices).
		Find(&transportationOffice, transportationOfficeID)

	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, apperror.NewNotFoundError(transportationOfficeID, "while looking for TransportationOffice")
		default:
			return nil, apperror.NewQueryError("GetTransportationOffice by transportationOfficeID", err, "")
		}
	}

	return &transportationOffice, nil
}

func (o transportationOfficesFetcher) GetTransportationOffices(appCtx appcontext.AppContext, search string, forPpm bool) (*models.TransportationOffices, error) {
	officeList, err := FindTransportationOffice(appCtx, search, forPpm)

	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return &officeList, apperror.NewNotFoundError(uuid.Nil, "Search string: "+search)
		default:
			return &officeList, err
		}
	}

	return &officeList, nil
}

func FindTransportationOffice(appCtx appcontext.AppContext, search string, forPpm bool) (models.TransportationOffices, error) {
	var officeList []models.TransportationOffice

	// The % operator filters out strings that are below this similarity threshold
	err := appCtx.DB().Q().RawQuery("SET pg_trgm.similarity_threshold = 0.03").Exec()
	if err != nil {
		return officeList, err
	}
	providesPPMCloseout := `and provides_ppm_closeout is true`

	sqlQuery := `
		with names as (select office.id as transportation_office_id, office.name, similarity(office.name, $1) as sim
        from transportation_offices as office
        where name % $1 `
	if forPpm {
		sqlQuery += providesPPMCloseout
	}
	sqlQuery += `
		order by sim desc
        limit 5)
		select office.*
        from names n inner join transportation_offices office on n.transportation_office_id = office.id
        group by office.id
        order by max(n.sim) desc, office.name
        limit 5`
	query := appCtx.DB().Q().RawQuery(sqlQuery, search)
	if err := query.All(&officeList); err != nil {
		if errors.Cause(err).Error() != models.RecordNotFoundErrorString {
			return officeList, err
		}
	}
	for i := range officeList {
		err := appCtx.DB().Load(&officeList[i], "Address", "Address.Country")
		if err != nil {
			return officeList, err
		}
	}
	return officeList, nil
}

func (o transportationOfficesFetcher) GetAllGBLOCs(appCtx appcontext.AppContext) (*models.GBLOCs, error) {
	gblocsList, err := ListDistinctGBLOCs(appCtx)

	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return &gblocsList, apperror.NewNotFoundError(uuid.Nil, "No GBLOCS found")
		default:
			return &gblocsList, err
		}
	}

	return &gblocsList, nil
}

func ListDistinctGBLOCs(appCtx appcontext.AppContext) (models.GBLOCs, error) {
	var gblocList models.GBLOCs

	err := appCtx.DB().RawQuery("SELECT DISTINCT gbloc FROM transportation_offices ORDER BY gbloc ASC").All(&gblocList)
	if err != nil {
		return gblocList, err
	}

	return gblocList, err
}

func (o transportationOfficesFetcher) GetCounselingOffices(appCtx appcontext.AppContext, dutyLocationID uuid.UUID) (*models.TransportationOffices, error) {
	officeList, err := findCounselingOffice(appCtx, dutyLocationID)

	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return &officeList, apperror.NewNotFoundError(uuid.Nil, "dutyLocationID not found")
		default:
			return &officeList, err
		}
	}

	return &officeList, nil
}

// return all the transportation offices in the GBLOC of the given duty location where provides_services_counseling = true
func findCounselingOffice(appCtx appcontext.AppContext, dutyLocationID uuid.UUID) (models.TransportationOffices, error) {
	var officeList []models.TransportationOffice

	sqlQuery := `
	with counseling_offices as (
                SELECT transportation_offices.id, transportation_offices.name, transportation_offices.address_id as counseling_address, substring(addresses.postal_code, 1,3 ) as pickup_zip
                        FROM postal_code_to_gblocs
                        JOIN addresses on postal_code_to_gblocs.postal_code = addresses.postal_code
                        JOIN duty_locations on addresses.id = duty_locations.address_id
                        JOIN transportation_offices on postal_code_to_gblocs.gbloc = transportation_offices.gbloc
                        WHERE duty_locations.provides_services_counseling = true and duty_locations.id = $1
                )
        SELECT counseling_offices.id, counseling_offices.name
                FROM counseling_offices
                JOIN duty_locations duty_locations2 on counseling_offices.id = duty_locations2.transportation_office_id
                JOIN addresses on counseling_offices.counseling_address = addresses.id
                JOIN re_us_post_regions on addresses.postal_code = re_us_post_regions.uspr_zip_id
                LEFT JOIN zip3_distances ON (
		                (re_us_post_regions.zip3 = zip3_distances.to_zip3
		            AND counseling_offices.pickup_zip = zip3_distances.from_zip3)
		                OR
		                (re_us_post_regions.zip3 = zip3_distances.from_zip3
		            AND counseling_offices.pickup_zip = zip3_distances.to_zip3)
		        )
                WHERE duty_locations2.provides_services_counseling = true
        group by counseling_offices.id, counseling_offices.name, zip3_distances.distance_miles
                ORDER BY coalesce(zip3_distances.distance_miles,0), counseling_offices.name asc`

	query := appCtx.DB().Q().RawQuery(sqlQuery, dutyLocationID)
	if err := query.All(&officeList); err != nil {
		if errors.Cause(err).Error() != models.RecordNotFoundErrorString {
			return officeList, err
		}
	}

	return officeList, nil
}
