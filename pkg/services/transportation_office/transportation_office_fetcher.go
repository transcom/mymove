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
	q := appCtx.DB().EagerPreload("Address", "Address.Country")
	if includeOnlyPPMCloseoutOffices {
		q.Where("provides_ppm_closeout = ?", includeOnlyPPMCloseoutOffices)
	}
	err := q.Find(&transportationOffice, transportationOfficeID)

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

func (o transportationOfficesFetcher) GetTransportationOffices(appCtx appcontext.AppContext, search string, forPpm bool, forAdminOfficeUserReqFilter bool) (*models.TransportationOffices, error) {
	officeList, err := FindTransportationOffice(appCtx, search, forPpm, forAdminOfficeUserReqFilter)

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

func FindTransportationOffice(appCtx appcontext.AppContext, search string, forPpm bool, forAdminOfficeUserReqFilter bool) (models.TransportationOffices, error) {
	var officeList []models.TransportationOffice

	// Changing return limit for Admin Requested Office Users Transportation Office Filter implementation
	var limit = 5
	if forAdminOfficeUserReqFilter {
		limit = 50
	}

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
        limit $2)
		select office.*
        from names n inner join transportation_offices office on n.transportation_office_id = office.id
        group by office.id
        order by max(n.sim) desc, office.name
        limit $2`
	query := appCtx.DB().Q().RawQuery(sqlQuery, search, limit)
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

// return all the transportation offices in the GBLOC of the given duty location where provides_services_counseling = true
// serviceMemberID is only provided when this function is called by the office handler
func (o transportationOfficesFetcher) GetCounselingOffices(appCtx appcontext.AppContext, dutyLocationID uuid.UUID, serviceMemberID uuid.UUID) (*models.TransportationOffices, error) {
	officeList, err := models.GetCounselingOffices(appCtx.DB(), dutyLocationID, serviceMemberID)

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
