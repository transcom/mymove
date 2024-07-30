package internalapi

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	transportationofficeop "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/transportation_offices"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/handlers/internalapi/internal/payloads"
	"github.com/transcom/mymove/pkg/models"
)

// ShowCounselingOfficesHandler returns the counseling offices for a duty location ID
type ShowCounselingOfficesHandler struct {
	handlers.HandlerConfig
}

// Handle retrieves the counseling office in the system for a given duty location ID
/*func (h ShowCounselingOfficesHandler) Handle(params transportationofficeop.ShowCounselingOfficesParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {
			//dutyLocationID, _ := uuid.FromString(params.DutyLocationID.String())
			dutyLocationID, _ := uuid.FromString("1234")
			counselingOffices, err := models.GetCounselingOffices(appCtx.DB(), dutyLocationID)
			if err != nil {
				return handlers.ResponseForError(appCtx.Logger(), err), err
			}
			counselingOfficePayload := payloads.TransportationOffices(counselingOffices)

			return transportationofficeop.NewShowCounselingOfficesOK().WithPayload(counselingOfficePayload), nil
		})
}*/
func (h ShowCounselingOfficesHandler) Handle(params transportationofficeop.ShowDutyLocationTransportationOfficeParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {
			//dutyLocationID, _ := uuid.FromString(params.DutyLocationID.String())
			dutyLocationID, _ := uuid.FromString("1234")
			counselingOffices, err := models.GetCounselingOffices(appCtx.DB(), dutyLocationID)
			if err != nil {
				return handlers.ResponseForError(appCtx.Logger(), err), err
			}
			counselingOfficePayload := payloads.TransportationOffices(counselingOffices)

			return transportationofficeop.NewGetTransportationOfficesOK().WithPayload(counselingOfficePayload), nil
		})
}
