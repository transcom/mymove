package ghcapi

import (
	"fmt"

	"github.com/go-openapi/runtime/middleware"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	ppmdocumentops "github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations/ppm"
	"github.com/transcom/mymove/pkg/gen/ghcmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/handlers/ghcapi/internal/payloads"
)

// GetPPMCloseoutDataHandler is the handler that fetches all of the documents for a PPM shipment for the office api
type GetPPMCloseoutDataHandler struct {
	handlers.HandlerConfig
}

// Handle retrieves all documents for a PPM shipment
func (h GetPPMCloseoutDataHandler) Handle(params ppmdocumentops.GetPPMCloseoutDataParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {
			errInstance := fmt.Sprintf("Instance: %s", h.GetTraceIDFromRequest(params.HTTPRequest))

			errPayload := &ghcmessages.Error{Message: &errInstance}

			if !appCtx.Session().IsOfficeApp() {
				return ppmdocumentops.NewGetPPMCloseoutDataForbidden().WithPayload(errPayload), apperror.NewSessionError("Request should come from the office app.")
			}

			ppmShipmentID := uuid.FromStringOrNil(params.PpmShipmentID.String())

			ppmShipmentIDString := ppmShipmentID.String()

			// if err != nil {
			// 	appCtx.Logger().Error("ghcapi.GetPPMCloseoutDataHandler error", zap.Error(err))

			// 	switch e := err.(type) {
			// 	case apperror.QueryError:
			// 		if e.Unwrap() != nil {
			// 			// If you can unwrap, log the error (usually a pq error) for better debugging
			// 			appCtx.Logger().Error(
			// 				"ghcapi.GetPPMCloseoutDataHandler error",
			// 				zap.Error(e.Unwrap()),
			// 			)
			// 		}

			// 		return ppmdocumentops.NewGetPPMCloseoutDataInternalServerError().WithPayload(errPayload), nil
			// 	default:
			// 		return ppmdocumentops.NewGetPPMCloseoutDataInternalServerError().WithPayload(errPayload), nil
			// 	}
			// }

			returnPayload := payloads.PPMCloseoutData(&ppmShipmentIDString)

			return ppmdocumentops.NewGetPPMCloseoutDataOK().WithPayload(returnPayload), nil
		})
}
