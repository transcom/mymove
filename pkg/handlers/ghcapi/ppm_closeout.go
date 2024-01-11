package ghcapi

import (
	"fmt"

	"github.com/go-openapi/runtime/middleware"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	ppmcloseoutops "github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations/ppm"
	"github.com/transcom/mymove/pkg/gen/ghcmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/handlers/ghcapi/internal/payloads"
	"github.com/transcom/mymove/pkg/services"
)

// GetPPMCloseoutHandler is the handler that fetches all of the calculations for a PPM closeout for the office api
type GetPPMCloseoutHandler struct {
	handlers.HandlerConfig
	ppmCloseoutFetcher services.PPMCloseoutFetcher
}

// Handle retrieves all calcuations for a PPM closeout
func (h GetPPMCloseoutHandler) Handle(params ppmcloseoutops.GetPPMCloseoutParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {
			// TODO - uncomment and edit once closeout return starts returning values besides ID.
			// handleError := func(err error) (middleware.Responder, error) {
			// 	appCtx.Logger().Error("ListMTOShipmentsHandler error", zap.Error(err))
			// 	payload := &ghcmessages.Error{Message: handlers.FmtString(err.Error())}
			// 	switch err.(type) {
			// 	case apperror.NotFoundError:
			// 		return mtoshipmentops.NewListMTOShipmentsNotFound().WithPayload(payload), err
			// 	case apperror.ForbiddenError:
			// 		return mtoshipmentops.NewListMTOShipmentsForbidden().WithPayload(payload), err
			// 	case apperror.QueryError:
			// 		return mtoshipmentops.NewListMTOShipmentsInternalServerError(), err
			// 	default:
			// 		return mtoshipmentops.NewListMTOShipmentsInternalServerError(), err
			// 	}
			// }

			errInstance := fmt.Sprintf("Instance: %s", h.GetTraceIDFromRequest(params.HTTPRequest))

			errPayload := &ghcmessages.Error{Message: &errInstance}

			if !appCtx.Session().IsOfficeApp() {
				return ppmcloseoutops.NewGetPPMCloseoutForbidden().WithPayload(errPayload), apperror.NewSessionError("Request should come from the office app.")
			}

			ppmShipmentID := uuid.FromStringOrNil(params.PpmShipmentID.String())

			// TODO - uncomment and edit once closeout return starts returning values besides ID.
			// ppmCloseout, err := h.ppmCloseoutFetcher.GetPPMCloseout()(appCtx, ppmShipmentID)
			// if err != nil {
			// 	return handleError(err)
			// }

			ppmShipmentIDString := ppmShipmentID.String()

			returnPayload := payloads.PPMCloseout(&ppmShipmentIDString)

			return ppmcloseoutops.NewGetPPMCloseoutOK().WithPayload(returnPayload), nil
		})
}
