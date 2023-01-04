package services

import (
	"github.com/gofrs/uuid"
	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
)

//go:generate mockery --name TransportationOfficesFetcher --disable-version-string
type TransportationOfficesFetcher interface {
	GetTransportationOffices(appCtx appcontext.AppContext, search string) (*models.TransportationOffices, error)
	GetTransportationOffice(appCtx appcontext.AppContext, transportationOfficeID uuid.UUID, includeOnlyPPMCloseoutOffices bool) (*models.TransportationOffice, error)
}
