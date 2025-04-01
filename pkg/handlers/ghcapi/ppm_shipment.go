package ghcapi

import (
	"fmt"
	"time"

	"github.com/go-openapi/runtime/middleware"
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	ppmsitops "github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations/ppm"
	"github.com/transcom/mymove/pkg/gen/ghcmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/handlers/ghcapi/internal/payloads"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
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
			ppmEagerAssociations := []string{"PickupAddress",
				"DestinationAddress",
				"SecondaryPickupAddress",
				"SecondaryDestinationAddress",
				"Shipment",
			}
			ppmShipment, err := h.GetPPMShipment(appCtx, ppmShipmentID, ppmEagerAssociations, nil)

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

			ppmShipment.SITEstimatedEntryDate = (*time.Time)(&params.SitEntryDate)
			ppmShipment.SITEstimatedDepartureDate = (*time.Time)(&params.SitDepartureDate)
			ppmShipment.SITEstimatedWeight = handlers.PoundPtrFromInt64Ptr(&params.WeightStored)
			sitExpected := true
			ppmShipment.SITExpected = &sitExpected

			calculatedCostDetails, err := h.PPMEstimator.CalculatePPMSITEstimatedCostBreakdown(appCtx, ppmShipment)

			if err != nil {
				return handleError(err)
			}

			returnPayload := payloads.PPMSITEstimatedCost(calculatedCostDetails)

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
			ppmEagerAssociations := []string{"PickupAddress",
				"DestinationAddress",
				"SecondaryPickupAddress",
				"SecondaryDestinationAddress",
				"Shipment",
			}
			ppmShipment, err := h.GetPPMShipment(appCtx, ppmShipmentID, ppmEagerAssociations, nil)

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

			// We set sitExpected to true because this is a storage moving expense therefore SIT has to be true
			// The case where this could be false at this point is when the Customer created the shipment they answered No to SIT Expected question,
			// but later decided they needed SIT and submitted a moving expense for storage or if the Service Counselor adds one.
			sitExpected := true
			ppmShipment.SITExpected = &sitExpected
			updatedPPMShipment, err := h.PPMShipmentUpdater.UpdatePPMShipmentSITEstimatedCost(appCtx, ppmShipment)

			if err != nil {
				return handleError(err)
			}

			returnPayload := payloads.PPMShipment(h.FileStorer(), updatedPPMShipment)

			return ppmsitops.NewUpdatePPMSITOK().WithPayload(returnPayload), nil
		})
}

// SubmitPPMShipmentDocumentationHandler is the handler to save a PPMShipment signature and route the PPM shipment to the office
type SubmitPPMShipmentDocumentationHandler struct {
	handlers.HandlerConfig
	services.PPMShipmentNewSubmitter
}

// Handle saves a new customer signature for PPMShipment documentation submission and routes PPM shipment to the
// service counselor.
func (h SubmitPPMShipmentDocumentationHandler) Handle(params ppmsitops.SubmitPPMShipmentDocumentationParams) middleware.Responder {
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

			ppmShipmentID, err := uuid.FromString(params.PpmShipmentID.String())
			if err != nil {
				appCtx.Logger().Error("error with PPM Shipment ID", zap.Error(err))
				return handleError(err)
			} else if ppmShipmentID.IsNil() {
				appCtx.Logger().Error("nil PPM Shipment ID")
				return nil, errors.New("nil PPM shipment ID")
			}

			ppmShipment, err := h.PPMShipmentNewSubmitter.SubmitNewCustomerCloseOut(appCtx, ppmShipmentID, models.SignedCertification{})
			if err != nil {
				appCtx.Logger().Error("ghcapi.SubmitPPMShipmentDocumentationHandler", zap.Error(err))
				return handleError(err)
			}

			returnPayload := payloads.PPMShipment(h.FileStorer(), ppmShipment)

			return ppmsitops.NewSubmitPPMShipmentDocumentationOK().WithPayload(returnPayload), nil
		})
}
