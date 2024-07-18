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
	err := appCtx.DB().EagerPreload("Address").
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

func (o transportationOfficesFetcher) GetTransportationOffices(appCtx appcontext.AppContext, search string) (*models.TransportationOffices, error) {
	officeList, err := FindTransportationOffice(appCtx, search)

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

func FindTransportationOffice(appCtx appcontext.AppContext, search string) (models.TransportationOffices, error) {
	var officeList []models.TransportationOffice

	// The % operator filters out strings that are below this similarity threshold
	err := appCtx.DB().Q().RawQuery("SET pg_trgm.similarity_threshold = 0.03").Exec()
	if err != nil {
		return officeList, err
	}

	sqlQuery := `
	with names as (select office.id as transportation_office_id, office.name, similarity(office.name, $1) as sim
        from transportation_offices as office
        where name % $1
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
		err := appCtx.DB().Load(&officeList[i], "Address")
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
