package transportationoffice

import (
	"database/sql"

	"github.com/gofrs/uuid"
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

func (o transportationOfficesFetcher) GetTransportationOffices(appCtx appcontext.AppContext, search string) (*models.TransportationOffices, error) {
	officeList, err := models.FindTransportationOffice(appCtx.DB(), search)

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
