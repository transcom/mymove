package ghcapi

import (
	"fmt"

	"github.com/go-openapi/runtime/middleware"
	"github.com/gofrs/uuid"
	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	ppmsitops "github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations/ppm"
	"github.com/transcom/mymove/pkg/gen/ghcmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/handlers/ghcapi/internal/payloads"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"go.uber.org/zap"
)

// GetPPMSITEstimatedCostHandler is the handler that calculates SIT Estimated Cost for the PPM Shipment
type GetPPMSITEstimatedCostHandler struct {
	handlers.HandlerConfig
	services.PPMEstimator
	services.PPMShipmentFetcher
}

// Handle calculates and returns SIT Estimated Cost for the PPM Shipment
func (h GetPPMSITEstimatedCostHandler) Handle(params ppmsitops.GetPPMSITEstimatedCostParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {

			handleError := func(err error) (middleware.Responder, error) {
				appCtx.Logger().Error("GetPPMSITEstimatedCost error", zap.Error(err))
				payload := &ghcmessages.Error{Message: handlers.FmtString(err.Error())}
				switch err.(type) {
				case apperror.NotFoundError:
					return ppmsitops.NewGetPPMSITEstimatedCostNotFound().WithPayload(payload), err
				case apperror.ForbiddenError:
					return ppmsitops.NewGetPPMSITEstimatedCostForbidden().WithPayload(payload), err
				case apperror.QueryError:
					return ppmsitops.NewGetPPMSITEstimatedCostInternalServerError().WithPayload(payload), err
				default:
					return ppmsitops.NewGetPPMSITEstimatedCostInternalServerError().WithPayload(payload), err
				}
			}

			errInstance := fmt.Sprintf("Instance: %s", h.GetTraceIDFromRequest(params.HTTPRequest))
			errPayload := &ghcmessages.Error{Message: &errInstance}

			if !appCtx.Session().IsOfficeApp() {
				return ppmsitops.NewGetPPMSITEstimatedCostForbidden().WithPayload(errPayload), apperror.NewSessionError("Request should come from the office app.")
			}

			ppmShipmentID := uuid.FromStringOrNil(params.PpmShipmentID.String())
			ppmShipment, err := h.GetPPMShipment(appCtx, ppmShipmentID, nil, nil)

			if err != nil {
				return handleError(err)
			}

			sitLocationOrigin := models.SITLocationTypeOrigin
			sitLocationDestination := models.SITLocationTypeDestination

			if params.SitLocation == (string)(sitLocationOrigin) {
				ppmShipment.SITLocation = &sitLocationOrigin
			} else {
				ppmShipment.SITLocation = &sitLocationDestination
			}

			calculatedCost, err := h.PPMEstimator.CalculatePPMSITEstimatedCost(appCtx, ppmShipment)

			if err != nil {
				return handleError(err)
			}

			returnPayload := payloads.PPMSITEstimatedCost(calculatedCost)

			return ppmsitops.NewGetPPMSITEstimatedCostOK().WithPayload(returnPayload), nil
		})
}

// UpdatePPMSITHandler is the handler that updates SIT related data
type UpdatePPMSITHandler struct {
	handlers.HandlerConfig
	services.PPMShipmentUpdater
	services.PPMShipmentFetcher
}

// Handle updates SIT related data for PPM
func (h UpdatePPMSITHandler) Handle(params ppmsitops.UpdatePPMSITParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {

			handleError := func(err error) (middleware.Responder, error) {
				appCtx.Logger().Error("UpdatePPMSIT error", zap.Error(err))
				payload := &ghcmessages.Error{Message: handlers.FmtString(err.Error())}
				switch err.(type) {
				case apperror.NotFoundError:
					return ppmsitops.NewUpdatePPMSITNotFound().WithPayload(payload), err
				case apperror.ForbiddenError:
					return ppmsitops.NewUpdatePPMSITForbidden().WithPayload(payload), err
				case apperror.QueryError:
					return ppmsitops.NewUpdatePPMSITInternalServerError().WithPayload(payload), err
				default:
					return ppmsitops.NewUpdatePPMSITInternalServerError().WithPayload(payload), err
				}
			}

			errInstance := fmt.Sprintf("Instance: %s", h.GetTraceIDFromRequest(params.HTTPRequest))
			errPayload := &ghcmessages.Error{Message: &errInstance}

			if !appCtx.Session().IsOfficeApp() {
				return ppmsitops.NewUpdatePPMSITForbidden().WithPayload(errPayload), apperror.NewSessionError("Request should come from the office app.")
			}

			ppmShipmentID := uuid.FromStringOrNil(params.PpmShipmentID.String())
			ppmShipment, err := h.GetPPMShipment(appCtx, ppmShipmentID, nil, nil)

			if err != nil {
				return handleError(err)
			}

			payload := params.Body

			if payload == nil {
				invalidShipmentError := apperror.NewBadDataError("Invalid ppm shipment: params Body is nil")
				appCtx.Logger().Error(invalidShipmentError.Error())
				return ppmsitops.NewUpdatePPMSITBadRequest(), invalidShipmentError
			}

			ppmShipment.SITLocation = (*models.SITLocationType)(payload.SitLocation)
			updatedPPMShipment, err := h.PPMShipmentUpdater.UpdatePPMShipmentSITEstimatedCost(appCtx, ppmShipment)

			if err != nil {
				return handleError(err)
			}

			returnPayload := payloads.PPMShipment(h.FileStorer(), updatedPPMShipment)

			return ppmsitops.NewUpdatePPMSITOK().WithPayload(returnPayload), nil
		})
}
