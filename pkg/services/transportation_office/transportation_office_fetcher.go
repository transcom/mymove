package transportationoffice

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

type oconusGblocDepartmentIndicator struct {
	Gbloc               string  `db:"gbloc" rw:"r"`
	RateAreaName        string  `db:"rate_area_name" rw:"r"`
	DepartmentIndicator *string `db:"department_indicator" rw:"r"`
}

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

	duty_location, err := models.FetchDutyLocation(appCtx.DB(), dutyLocationID)
	if err != nil {
		return officeList, err
	}

	var sqlQuery string

	// ********************************
	// Find for oconus duty location
	// ********************************
	if *duty_location.Address.IsOconus {
		gblocDepartmentIndicator, err := findOconusGblocDepartmentIndicator(appCtx, duty_location)
		if err != nil {
			return officeList, err
		}

		sqlQuery = `
		with counseling_offices as (
			SELECT transportation_offices.id, transportation_offices.name, transportation_offices.address_id as counseling_address, substring(a.postal_code, 1,3 ) as pickup_zip
			FROM duty_locations
			JOIN addresses a on duty_locations.address_id = a.id
			JOIN v_locations v on (a.us_post_region_cities_id = v.uprc_id or v.uprc_id is null)
					and a.country_id = v.country_id
			JOIN re_oconus_rate_areas r on r.us_post_region_cities_id = v.uprc_id
			JOIN gbloc_aors on gbloc_aors.oconus_rate_area_id = r.id
			JOIN jppso_regions j on gbloc_aors.jppso_regions_id = j.id
			JOIN transportation_offices on j.code = transportation_offices.gbloc
			join addresses a2 on a2.id = transportation_offices.address_id
			WHERE duty_locations.provides_services_counseling = true and duty_locations.id = $1 and j.code = $2
			)
		SELECT counseling_offices.id, counseling_offices.name
			FROM counseling_offices
			JOIN addresses cnsl_address on counseling_offices.counseling_address = cnsl_address.id
			LEFT JOIN zip3_distances ON (
				(substring(cnsl_address.postal_code,1 ,3) = zip3_distances.to_zip3
				AND counseling_offices.pickup_zip = zip3_distances.from_zip3)
				OR
				(substring(cnsl_address.postal_code,1 ,3) = zip3_distances.from_zip3
				AND counseling_offices.pickup_zip = zip3_distances.to_zip3)
			)
			group by counseling_offices.id, counseling_offices.name, zip3_distances.distance_miles
			ORDER BY coalesce(zip3_distances.distance_miles,0) asc`

		query := appCtx.DB().Q().RawQuery(sqlQuery, dutyLocationID, gblocDepartmentIndicator.Gbloc)
		if err := query.All(&officeList); err != nil {
			if errors.Cause(err).Error() != models.RecordNotFoundErrorString {
				return officeList, err
			}
		}
		return officeList, nil
	}

	// ********************************
	// Find for conus duty location
	// ********************************
	sqlQuery = `
	with counseling_offices as (
		SELECT transportation_offices.id, transportation_offices.name, transportation_offices.address_id as counseling_address, substring(addresses.postal_code, 1,3 ) as origin_zip, substring(a2.postal_code, 1,3 ) as dest_zip
			FROM postal_code_to_gblocs
			JOIN addresses on postal_code_to_gblocs.postal_code = addresses.postal_code
			JOIN duty_locations on addresses.id = duty_locations.address_id
			JOIN transportation_offices on postal_code_to_gblocs.gbloc = transportation_offices.gbloc
			join addresses a2 on a2.id = transportation_offices.address_id
			WHERE duty_locations.provides_services_counseling = true and duty_locations.id = $1
		)
	SELECT counseling_offices.id, counseling_offices.name
		FROM counseling_offices
		JOIN duty_locations duty_locations2 on counseling_offices.id = duty_locations2.transportation_office_id
		JOIN addresses on counseling_offices.counseling_address = addresses.id
		LEFT JOIN zip3_distances ON (
	    	(substring(addresses.postal_code,1 ,3) = zip3_distances.to_zip3
	        AND counseling_offices.origin_zip = zip3_distances.from_zip3)
	    	OR
	    	(substring(addresses.postal_code,1 ,3) = zip3_distances.from_zip3
	        AND counseling_offices.origin_zip = zip3_distances.to_zip3)
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

func findOconusGblocDepartmentIndicator(appCtx appcontext.AppContext, dutyLocation models.DutyLocation) (*oconusGblocDepartmentIndicator, error) {
	serviceMember, err := models.FetchServiceMember(appCtx.DB(), appCtx.Session().ServiceMemberID)
	if err != nil {
		return nil, err
	}

	var oconusGblocDepartmentIndicator []oconusGblocDepartmentIndicator

	sqlQuery := `
    select j.code gbloc, r.name rate_area_name, g.department_indicator
    from addresses a,
    v_locations v,
    re_oconus_rate_areas o,
    re_rate_areas r,
    jppso_regions j,
    gbloc_aors g
    where a.id = $1
    and a.us_post_region_cities_id = v.uprc_id
    and v.uprc_id = o.us_post_region_cities_id
    and o.rate_area_id = r.id
    and o.id = g.oconus_rate_area_id
    and j.id = g.jppso_regions_id`

	query := appCtx.DB().Q().RawQuery(sqlQuery, dutyLocation.Address.ID)
	err = query.All(&oconusGblocDepartmentIndicator)
	if err != nil {
		return nil, err
	}

	// Determine departmentIndicator based on service member's affiliation
	var departmentIndicator *string = nil
	if serviceMember.Affiliation != nil && (*serviceMember.Affiliation == models.AffiliationAIRFORCE || *serviceMember.Affiliation == models.AffiliationSPACEFORCE) {
		departmentIndicator = models.StringPointer(models.DepartmentIndicatorAIRANDSPACEFORCE.String())
	} else if serviceMember.Affiliation != nil && (*serviceMember.Affiliation == models.AffiliationNAVY || *serviceMember.Affiliation == models.AffiliationMARINES) {
		departmentIndicator = models.StringPointer(models.DepartmentIndicatorNAVYANDMARINES.String())
	} else if serviceMember.Affiliation != nil && (*serviceMember.Affiliation == models.AffiliationARMY) {
		departmentIndicator = models.StringPointer(models.DepartmentIndicatorARMY.String())
	} else if serviceMember.Affiliation != nil && (*serviceMember.Affiliation == models.AffiliationCOASTGUARD) {
		departmentIndicator = models.StringPointer(models.DepartmentIndicatorCOASTGUARD.String())
	}

	// Is there a matching GBLOC for duty location address specifically for user's affiliation and duty location Zone(I-V)?
	// sample oconusGblocAffiliationInfo[]:
	// Gbloc     RateAreaName	       DepartmentIndicator
	// JEAT   	 Alaska (Zone) II      NULL                       (default)
	// MBFL	     Alaska (Zone) II      AIR_AND_SPACE_FORCE
	for _, info := range oconusGblocDepartmentIndicator {
		if (info.DepartmentIndicator != nil && departmentIndicator != nil) && (*info.DepartmentIndicator == *departmentIndicator) {
			appCtx.Logger().Debug(fmt.Sprintf("Found specific department match -- serviceMember.Affiliation: %s, DutyLocaton: %s, GBLOC: %s, departmentIndicator: %s, RateAreaName: %s, dutyLocation.Address.ID: %s",
				serviceMember.Affiliation, dutyLocation.Name, info.Gbloc, *departmentIndicator, info.RateAreaName, dutyLocation.Address.ID))
			return &info, nil
		}
	}

	// These were no departmentIndicator specific GBLOC for duty location -- return default GBLOC
	for _, info := range oconusGblocDepartmentIndicator {
		if info.DepartmentIndicator == nil {
			appCtx.Logger().Debug(fmt.Sprintf("Did not find specific department match, return default -- serviceMember.Affiliation: %s, DutyLocaton: %s, GBLOC: %s, departmentIndicator: %s, RateAreaName: %s, dutyLocation.Address.ID: %s",
				serviceMember.Affiliation, dutyLocation.Name, info.Gbloc, "NIL/NULL", info.RateAreaName, dutyLocation.Address.ID))
			return &info, nil
		}
	}

	// There is no default and department specific oconusGblocDepartmentIndicator. There is nothing in system to support. This should never happen.
	return nil, apperror.NewImplementationError(fmt.Sprintf("Error: Cannot determine GBLOC -- serviceMember.Affiliation: %s, DutyLocaton: %s, departmentIndicator: %s, dutyLocation.Address.ID: %s",
		serviceMember.Affiliation, dutyLocation.Name, *departmentIndicator, dutyLocation.Address.ID))
}
