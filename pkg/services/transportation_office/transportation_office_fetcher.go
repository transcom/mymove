package transportationoffice

import (
	"database/sql"
	"time"

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

func (o transportationOfficesFetcher) GetTransportationOffices(appCtx appcontext.AppContext) (*models.TransportationOffices, error) {

	// Mock response values, to be replaced with real values in later ticket
	mockTransportationOffices := models.TransportationOffices{
		{
			ID:   uuid.FromStringOrNil("a2119fd0-bd06-4055-a94a-a266385780dc"),
			Name: "Transportation Office 1",
			Address: models.Address{
				ID:             uuid.FromStringOrNil("a2119fd0-bd06-4055-a94a-a266385780dc"),
				StreetAddress1: "123 Main St",
				City:           "Anytown",
				State:          "CA",
				PostalCode:     "90210",
			},
			PhoneLines: models.OfficePhoneLines{
				{
					ID:     uuid.FromStringOrNil("b3119fd0-bd06-4055-a94a-a266385780ab"),
					Number: "555-555-5555",
					Type:   "voice",
				},
			},
			Gbloc:     "LKNQ",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			ID:   uuid.FromStringOrNil("21d3faf3-812c-455c-b1d4-339b43081f40"),
			Name: "Transportation Office 2",
			Address: models.Address{
				ID:             uuid.FromStringOrNil("a2119fd0-bd06-4055-a94a-a266385780dc"),
				StreetAddress1: "123 ABC St",
				City:           "Anycity",
				State:          "CO",
				PostalCode:     "80487",
			},
			PhoneLines: models.OfficePhoneLines{
				{
					ID:     uuid.FromStringOrNil("b3119fd0-bd06-4055-a94a-a266385780ab"),
					Number: "555-555-5556",
					Type:   "voice",
				},
			},
			Gbloc:     "KKFA",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}
	return &mockTransportationOffices, nil
}
