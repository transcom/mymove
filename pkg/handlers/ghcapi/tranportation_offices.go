package ghcapi

import (
	"github.com/go-openapi/runtime/middleware"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/appcontext"
	transportationofficeop "github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations/transportation_office"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/handlers/ghcapi/internal/payloads"
	"github.com/transcom/mymove/pkg/services"
)

type GetTransportationOfficesHandler struct {
	handlers.HandlerConfig
	services.TransportationOfficesFetcher
}

func (h GetTransportationOfficesHandler) Handle(params transportationofficeop.GetTransportationOfficesParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {

			transportationOffices, err := h.TransportationOfficesFetcher.GetTransportationOffices(appCtx)
			if err != nil {
				appCtx.Logger().Error("Error searching for Transportation Offices: ", zap.Error(err))
				return transportationofficeop.NewGetTransportationOfficesInternalServerError(), err
			}

			returnPayload := payloads.TransportationOffices(*transportationOffices)
			return transportationofficeop.NewGetTransportationOfficesOK().WithPayload(returnPayload), nil
		})
}
