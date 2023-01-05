package internalapi

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/gofrs/uuid"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	transportationofficeop "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/transportation_offices"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/handlers/internalapi/internal/payloads"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

// ShowDutyLocationTransportationOfficeHandler returns the transportation office for a duty location ID
type ShowDutyLocationTransportationOfficeHandler struct {
	handlers.HandlerConfig
}

// Handle retrieves the transportation office in the system for a given duty location ID
func (h ShowDutyLocationTransportationOfficeHandler) Handle(params transportationofficeop.ShowDutyLocationTransportationOfficeParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {
			dutyLocationID, _ := uuid.FromString(params.DutyLocationID.String())
			transportationOffice, err := models.FetchDutyLocationTransportationOffice(appCtx.DB(), dutyLocationID)
			if err != nil {
				return handlers.ResponseForError(appCtx.Logger(), err), err
			}
			transportationOfficePayload := payloads.TransportationOffice(transportationOffice)

			return transportationofficeop.NewShowDutyLocationTransportationOfficeOK().WithPayload(transportationOfficePayload), nil
		})
}

type GetTransportationOfficesHandler struct {
	handlers.HandlerConfig
	services.TransportationOfficesFetcher
}

func (h GetTransportationOfficesHandler) Handle(params transportationofficeop.GetTransportationOfficesParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {
			if appCtx.Session() == nil {
				noSessionErr := apperror.NewSessionError("No user session")
				return transportationofficeop.NewGetTransportationOfficesUnauthorized(), noSessionErr
			}
			if !appCtx.Session().IsMilApp() && appCtx.Session().ServiceMemberID == uuid.Nil {
				noServiceMemberIDErr := apperror.NewSessionError("No service member ID")
				return transportationofficeop.NewGetTransportationOfficesForbidden(), noServiceMemberIDErr
			}

			transportationOffices, err := h.TransportationOfficesFetcher.GetTransportationOffices(appCtx, params.Search)
			if err != nil {
				appCtx.Logger().Error("Error searching for Transportation Offices: ", zap.Error(err))
				return transportationofficeop.NewGetTransportationOfficesInternalServerError(), err
			}

			returnPayload := payloads.TransportationOffices(*transportationOffices)
			return transportationofficeop.NewGetTransportationOfficesOK().WithPayload(returnPayload), nil
		})
}
