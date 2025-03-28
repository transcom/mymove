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
	"github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations/ppm"
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
func (h GetPPMSITEstimatedCostHandler) Handle(params ppm.GetPPMSITEstimatedCostParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {

			handleError := func(err error) (middleware.Responder, error) {
				appCtx.Logger().Error("GetPPMSITEstimatedCost error", zap.Error(err))
				payload := &ghcmessages.Error{Message: handlers.FmtString(err.Error())}
				switch err.(type) {
				case apperror.NotFoundError:
					return ppm.NewGetPPMSITEstimatedCostNotFound().WithPayload(payload), err
				case apperror.ForbiddenError:
					return ppm.NewGetPPMSITEstimatedCostForbidden().WithPayload(payload), err
				case apperror.QueryError:
					return ppm.NewGetPPMSITEstimatedCostInternalServerError().WithPayload(payload), err
				default:
					return ppm.NewGetPPMSITEstimatedCostInternalServerError().WithPayload(payload), err
				}
			}

			errInstance := fmt.Sprintf("Instance: %s", h.GetTraceIDFromRequest(params.HTTPRequest))
			errPayload := &ghcmessages.Error{Message: &errInstance}

			if !appCtx.Session().IsOfficeApp() {
				return ppm.NewGetPPMSITEstimatedCostForbidden().WithPayload(errPayload), apperror.NewSessionError("Request should come from the office app.")
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

			return ppm.NewGetPPMSITEstimatedCostOK().WithPayload(returnPayload), nil
		})
}

// UpdatePPMSITHandler is the handler that updates SIT related data
type UpdatePPMSITHandler struct {
	handlers.HandlerConfig
	services.PPMShipmentUpdater
	services.PPMShipmentFetcher
}

// Handle updates SIT related data for PPM
func (h UpdatePPMSITHandler) Handle(params ppm.UpdatePPMSITParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {

			handleError := func(err error) (middleware.Responder, error) {
				appCtx.Logger().Error("UpdatePPMSIT error", zap.Error(err))
				payload := &ghcmessages.Error{Message: handlers.FmtString(err.Error())}
				switch err.(type) {
				case apperror.NotFoundError:
					return ppm.NewUpdatePPMSITNotFound().WithPayload(payload), err
				case apperror.ForbiddenError:
					return ppm.NewUpdatePPMSITForbidden().WithPayload(payload), err
				case apperror.QueryError:
					return ppm.NewUpdatePPMSITInternalServerError().WithPayload(payload), err
				default:
					return ppm.NewUpdatePPMSITInternalServerError().WithPayload(payload), err
				}
			}

			errInstance := fmt.Sprintf("Instance: %s", h.GetTraceIDFromRequest(params.HTTPRequest))
			errPayload := &ghcmessages.Error{Message: &errInstance}

			if !appCtx.Session().IsOfficeApp() {
				return ppm.NewUpdatePPMSITForbidden().WithPayload(errPayload), apperror.NewSessionError("Request should come from the office app.")
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
				return ppm.NewUpdatePPMSITBadRequest(), invalidShipmentError
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

			return ppm.NewUpdatePPMSITOK().WithPayload(returnPayload), nil
		})
}

// SendPPMToCustomerHandler is the handler that sends PPM to customer
type SendPPMToCustomerHandler struct {
	handlers.HandlerConfig
	services.PPMShipmentFetcher
	services.MoveTaskOrderUpdater
}

// Handle send PPM to customer status change
func (h SendPPMToCustomerHandler) Handle(params ppm.SendPPMToCustomerParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {

			handleError := func(err error) (middleware.Responder, error) {
				appCtx.Logger().Error("UpdatePPMSendToCustomer error", zap.Error(err))
				payload := &ghcmessages.Error{Message: handlers.FmtString(err.Error())}
				switch e := err.(type) {
				case apperror.NotFoundError:
					return ppm.NewSendPPMToCustomerNotFound().WithPayload(payload), err
				case apperror.ForbiddenError:
					return ppm.NewSendPPMToCustomerForbidden().WithPayload(payload), err
				case apperror.PreconditionFailedError:
					return ppm.NewSendPPMToCustomerPreconditionFailed().WithPayload(payload), err
				case apperror.InvalidInputError:
					payload := payloadForValidationError("Validation errors", "SendShipmentToCustomer", h.GetTraceIDFromRequest(params.HTTPRequest), e.ValidationErrors)
					return ppm.NewSendPPMToCustomerUnprocessableEntity().WithPayload(payload), err
				case apperror.QueryError:
					return ppm.NewSendPPMToCustomerInternalServerError().WithPayload(payload), err
				default:
					return ppm.NewSendPPMToCustomerInternalServerError().WithPayload(payload), err
				}
			}

			errInstance := fmt.Sprintf("Instance: %s", h.GetTraceIDFromRequest(params.HTTPRequest))
			errPayload := &ghcmessages.Error{Message: &errInstance}

			if !appCtx.Session().IsOfficeApp() {
				return ppm.NewSendPPMToCustomerForbidden().WithPayload(errPayload), apperror.NewSessionError("Request should come from the office app.")
			}

			ppmShipmentID, err := uuid.FromString(params.PpmShipmentID.String())
			if err != nil || ppmShipmentID.IsNil() {
				appCtx.Logger().Error("error with PPM Shipment ID", zap.Error(err))

				return ppm.NewSendPPMToCustomerBadRequest().WithPayload(errPayload), err
			}

			ppmEagerAssociations := []string{"PickupAddress",
				"DestinationAddress",
				"SecondaryPickupAddress",
				"SecondaryDestinationAddress",
				"Shipment.MoveTaskOrder"}

			ppmShipment, err := h.GetPPMShipment(appCtx, ppmShipmentID, ppmEagerAssociations, nil)
			if err != nil {
				return handleError(err)
			}

			ppmShipment, err = h.UpdateStatusServiceCounselingSendPPMToCustomer(appCtx, *ppmShipment, params.IfMatch, &ppmShipment.Shipment.MoveTaskOrder)
			if err != nil {
				return handleError(err)
			}

			returnPayload := payloads.PPMShipment(h.FileStorer(), ppmShipment)

			return ppm.NewSendPPMToCustomerOK().WithPayload(returnPayload), nil
		})
}

// SubmitPPMShipmentDocumentationHandler is the handler to save a PPMShipment signature and route the PPM shipment to the office
type SubmitPPMShipmentDocumentationHandler struct {
	handlers.HandlerConfig
	services.PPMShipmentNewSubmitter
}

// Handle saves a new customer signature for PPMShipment documentation submission and routes PPM shipment to the
// service counselor.
func (h SubmitPPMShipmentDocumentationHandler) Handle(params ppm.SubmitPPMShipmentDocumentationParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {
			handleError := func(err error) (middleware.Responder, error) {
				appCtx.Logger().Error("GetPPMSITEstimatedCost error", zap.Error(err))
				payload := &ghcmessages.Error{Message: handlers.FmtString(err.Error())}
				switch err.(type) {
				case apperror.NotFoundError:
					return ppm.NewGetPPMSITEstimatedCostNotFound().WithPayload(payload), err
				case apperror.ForbiddenError:
					return ppm.NewGetPPMSITEstimatedCostForbidden().WithPayload(payload), err
				case apperror.QueryError:
					return ppm.NewGetPPMSITEstimatedCostInternalServerError().WithPayload(payload), err
				default:
					return ppm.NewGetPPMSITEstimatedCostInternalServerError().WithPayload(payload), err
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
				appCtx.Logger().Error("internalapi.SubmitPPMShipmentDocumentationHandler", zap.Error(err))
				return handleError(err)
			}

			returnPayload := payloads.PPMShipment(h.FileStorer(), ppmShipment)

			return ppm.NewSubmitPPMShipmentDocumentationOK().WithPayload(returnPayload), nil
		})
}
