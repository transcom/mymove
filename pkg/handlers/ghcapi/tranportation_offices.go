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

			// B-21022: forPpm param is set true. This is used by PPM closeout widget. Need to ensure certain offices are included/excluded
			// if location has ppm closedout enabled.
			transportationOffices, err := h.TransportationOfficesFetcher.GetTransportationOffices(appCtx, params.Search, true)

			if err != nil {
				appCtx.Logger().Error("Error searching for Transportation Offices: ", zap.Error(err))
				return transportationofficeop.NewGetTransportationOfficesInternalServerError(), err
			}

			returnPayload := payloads.TransportationOffices(*transportationOffices)
			return transportationofficeop.NewGetTransportationOfficesOK().WithPayload(returnPayload), nil
		})
}

type GetTransportationOfficesOpenHandler struct {
	handlers.HandlerConfig
	services.TransportationOfficesFetcher
}

func (h GetTransportationOfficesOpenHandler) Handle(params transportationofficeop.GetTransportationOfficesOpenParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {

			transportationOffices, err := h.TransportationOfficesFetcher.GetTransportationOffices(appCtx, params.Search, false)
			if err != nil {
				appCtx.Logger().Error("Error searching for Transportation Offices: ", zap.Error(err))
				return transportationofficeop.NewGetTransportationOfficesOpenInternalServerError(), err
			}

			returnPayload := payloads.TransportationOffices(*transportationOffices)
			return transportationofficeop.NewGetTransportationOfficesOpenOK().WithPayload(returnPayload), nil
		})
}

type GetTransportationOfficesGBLOCsHandler struct {
	handlers.HandlerConfig
	services.TransportationOfficesFetcher
}

func (h GetTransportationOfficesGBLOCsHandler) Handle(params transportationofficeop.GetTransportationOfficesGBLOCsParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {

			transportationOffices, err := h.TransportationOfficesFetcher.GetAllGBLOCs(appCtx)
			if err != nil {
				appCtx.Logger().Error("Error listing distinct GBLOCs: ", zap.Error(err))
				return transportationofficeop.NewGetTransportationOfficesGBLOCsInternalServerError(), err
			}

			returnPayload := payloads.GBLOCs(*transportationOffices)
			return transportationofficeop.NewGetTransportationOfficesGBLOCsOK().WithPayload(returnPayload), nil
		})
}
