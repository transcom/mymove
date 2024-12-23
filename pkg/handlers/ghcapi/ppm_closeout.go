package ghcapi

import (
	"fmt"

	"github.com/go-openapi/runtime/middleware"
	"github.com/gofrs/uuid"
	"go.uber.org/zap"

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
	services.PPMCloseoutFetcher
}

// Handle retrieves all calcuations for a PPM closeout
func (h GetPPMCloseoutHandler) Handle(params ppmcloseoutops.GetPPMCloseoutParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {

			handleError := func(err error) (middleware.Responder, error) {
				appCtx.Logger().Error("GetShipment error", zap.Error(err))
				payload := &ghcmessages.Error{Message: handlers.FmtString(err.Error())}
				switch err.(type) {
				case apperror.NotFoundError:
					return ppmcloseoutops.NewGetPPMCloseoutNotFound().WithPayload(payload), err
				case apperror.PPMNotReadyForCloseoutError:
					return ppmcloseoutops.NewGetPPMCloseoutNotFound().WithPayload(payload), err
				case apperror.ForbiddenError:
					return ppmcloseoutops.NewGetPPMCloseoutForbidden().WithPayload(payload), err
				case apperror.QueryError:
					return ppmcloseoutops.NewGetPPMCloseoutInternalServerError().WithPayload(payload), err
				default:
					return ppmcloseoutops.NewGetPPMCloseoutInternalServerError().WithPayload(payload), err
				}
			}
			errInstance := fmt.Sprintf("Instance: %s", h.GetTraceIDFromRequest(params.HTTPRequest))

			errPayload := &ghcmessages.Error{Message: &errInstance}

			if !appCtx.Session().IsOfficeApp() {
				return ppmcloseoutops.NewGetPPMCloseoutForbidden().WithPayload(errPayload), apperror.NewSessionError("Request should come from the office app.")
			}
			ppmShipmentID := uuid.FromStringOrNil(params.PpmShipmentID.String())

			ppmCloseout, err := h.PPMCloseoutFetcher.GetPPMCloseout(appCtx, ppmShipmentID)
			if err != nil {
				return handleError(err)
			}

			returnPayload := payloads.PPMCloseout(ppmCloseout)

			return ppmcloseoutops.NewGetPPMCloseoutOK().WithPayload(returnPayload), nil
		})
}

type GetPPMActualWeightHandler struct {
	handlers.HandlerConfig
	services.PPMCloseoutFetcher
	ppmShipmentFetcher services.PPMShipmentFetcher
}

// Handle retrieves actual weight for a PPM pending closeout
func (h GetPPMActualWeightHandler) Handle(params ppmcloseoutops.GetPPMActualWeightParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {

			handleError := func(err error) (middleware.Responder, error) {
				appCtx.Logger().Error("GetShipment error", zap.Error(err))
				payload := &ghcmessages.Error{Message: handlers.FmtString(err.Error())}
				switch err.(type) {
				case apperror.NotFoundError:
					return ppmcloseoutops.NewGetPPMActualWeightNotFound().WithPayload(payload), err
				case apperror.PPMNotReadyForCloseoutError:
					return ppmcloseoutops.NewGetPPMActualWeightNotFound().WithPayload(payload), err
				case apperror.ForbiddenError:
					return ppmcloseoutops.NewGetPPMActualWeightForbidden().WithPayload(payload), err
				case apperror.QueryError:
					return ppmcloseoutops.NewGetPPMActualWeightInternalServerError().WithPayload(payload), err
				default:
					return ppmcloseoutops.NewGetPPMActualWeightInternalServerError().WithPayload(payload), err
				}
			}
			errInstance := fmt.Sprintf("Instance: %s", h.GetTraceIDFromRequest(params.HTTPRequest))

			errPayload := &ghcmessages.Error{Message: &errInstance}

			if !appCtx.Session().IsOfficeApp() {
				return ppmcloseoutops.NewGetPPMActualWeightForbidden().WithPayload(errPayload), apperror.NewSessionError("Request should come from the office app.")
			}
			ppmShipmentID := uuid.FromStringOrNil(params.PpmShipmentID.String())

			eagerAssociations := []string{
				"WeightTickets",
			}
			ppmShipment, err := h.ppmShipmentFetcher.GetPPMShipment(appCtx, ppmShipmentID, eagerAssociations, nil)
			if err != nil {
				return handleError(err)
			}

			ppmActualWeight, err := h.PPMCloseoutFetcher.GetActualWeight(ppmShipment)
			if err != nil {
				return handleError(err)
			}

			returnPayload := payloads.PPMActualWeight(&ppmActualWeight)

			return ppmcloseoutops.NewGetPPMActualWeightOK().WithPayload(returnPayload), nil
		})
}
