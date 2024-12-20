package primeapi

import (
	"fmt"

	"github.com/go-openapi/runtime/middleware"
	"github.com/gofrs/uuid"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	mtoshipmentops "github.com/transcom/mymove/pkg/gen/primeapi/primeoperations/mto_shipment"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/handlers/primeapi/payloads"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	mtoshipment "github.com/transcom/mymove/pkg/services/mto_shipment"
)

// UpdateShipmentDestinationAddressHandler is the handler to create address update request for non-SIT
type UpdateShipmentDestinationAddressHandler struct {
	handlers.HandlerConfig
	services.ShipmentAddressUpdateRequester
}

// Handle creates the address update request for non-SIT
func (h UpdateShipmentDestinationAddressHandler) Handle(params mtoshipmentops.UpdateShipmentDestinationAddressParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {
			payload := params.Body
			shipmentID := uuid.FromStringOrNil(params.MtoShipmentID.String())

			addressUpdate := payloads.ShipmentAddressUpdateModel(payload, shipmentID)

			eTag := params.IfMatch

			response, err := h.ShipmentAddressUpdateRequester.RequestShipmentDeliveryAddressUpdate(appCtx, shipmentID, addressUpdate.NewAddress, addressUpdate.ContractorRemarks, eTag)

			if err != nil {
				appCtx.Logger().Error("primeapi.UpdateShipmentDestinationAddressHandler error", zap.Error(err))

				switch e := err.(type) {

				// NotFoundError -> Not Found response
				case apperror.NotFoundError:
					return mtoshipmentops.NewUpdateShipmentDestinationAddressNotFound().WithPayload(payloads.ClientError(handlers.NotFoundMessage, err.Error(), h.GetTraceIDFromRequest(params.HTTPRequest))), err
				// ConflictError -> Request conflict reponse
				case apperror.ConflictError:
					return mtoshipmentops.NewUpdateShipmentDestinationAddressConflict().WithPayload(payloads.ClientError(handlers.ConflictErrMessage, err.Error(), h.GetTraceIDFromRequest(params.HTTPRequest))), err
				// PreconditionError -> precondition failed reponse
				case apperror.PreconditionFailedError:
					return mtoshipmentops.NewUpdateShipmentDestinationAddressPreconditionFailed().WithPayload(payloads.ClientError(handlers.PreconditionErrMessage, err.Error(), h.GetTraceIDFromRequest(params.HTTPRequest))), err
				// InvalidInputError -> Unprocessable Entity reponse
				case apperror.InvalidInputError:
					return mtoshipmentops.NewUpdateShipmentDestinationAddressUnprocessableEntity().WithPayload(payloads.ValidationError(err.Error(), h.GetTraceIDFromRequest(params.HTTPRequest), e.ValidationErrors)), err
				// Unknown -> Internal Server Error
				default:
					return mtoshipmentops.NewUpdateShipmentDestinationAddressInternalServerError().WithPayload(payloads.InternalServerError(nil, h.GetTraceIDFromRequest(params.HTTPRequest))), err
				}
			}
			returnPayload := payloads.ShipmentAddressUpdate(response)
			return mtoshipmentops.NewUpdateShipmentDestinationAddressCreated().WithPayload(returnPayload), nil
		})
}

// DeleteMTOShipmentHandler is the handler to soft delete MTO shipments
type DeleteMTOShipmentHandler struct {
	handlers.HandlerConfig
	services.ShipmentDeleter
}

// Handle handler that updates a mto shipment
func (h DeleteMTOShipmentHandler) Handle(params mtoshipmentops.DeleteMTOShipmentParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {
			shipmentID := uuid.FromStringOrNil(params.MtoShipmentID.String())
			_, err := h.DeleteShipment(appCtx, shipmentID)
			if err != nil {
				appCtx.Logger().Error("primeapi.DeleteMTOShipmentHandler", zap.Error(err))

				switch err.(type) {
				case apperror.NotFoundError:
					return mtoshipmentops.NewDeleteMTOShipmentNotFound().WithPayload(
						payloads.ClientError(handlers.NotFoundMessage, err.Error(), h.GetTraceIDFromRequest(params.HTTPRequest))), err
				case apperror.ConflictError:
					return mtoshipmentops.NewDeleteMTOShipmentConflict(), err
				case apperror.ForbiddenError:
					return mtoshipmentops.NewDeleteMTOShipmentForbidden().WithPayload(
						payloads.ClientError(handlers.ForbiddenErrMessage, err.Error(), h.GetTraceIDFromRequest(params.HTTPRequest))), err
				case apperror.UnprocessableEntityError:
					return mtoshipmentops.NewDeleteMTOShipmentUnprocessableEntity(), err
				default:
					return mtoshipmentops.NewDeleteMTOShipmentInternalServerError().WithPayload(
						payloads.InternalServerError(nil, h.GetTraceIDFromRequest(params.HTTPRequest))), err
				}
			}

			return mtoshipmentops.NewDeleteMTOShipmentNoContent(), nil
		})
}

// UpdateMTOShipmentStatusHandler is the handler to update MTO Shipments' status
type UpdateMTOShipmentStatusHandler struct {
	handlers.HandlerConfig
	checker services.MTOShipmentUpdater
	updater services.MTOShipmentStatusUpdater
}

// Handle handler that updates a mto shipment's status
func (h UpdateMTOShipmentStatusHandler) Handle(params mtoshipmentops.UpdateMTOShipmentStatusParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {

			shipmentID := uuid.FromStringOrNil(params.MtoShipmentID.String())

			availableToPrime, err := h.checker.MTOShipmentsMTOAvailableToPrime(appCtx, shipmentID)
			if err != nil {
				appCtx.Logger().Error("primeapi.UpdateMTOShipmentHandler error - MTO is not available to prime", zap.Error(err))
				switch e := err.(type) {
				case apperror.NotFoundError:
					return mtoshipmentops.NewUpdateMTOShipmentStatusNotFound().WithPayload(
						payloads.ClientError(handlers.NotFoundMessage, e.Error(), h.GetTraceIDFromRequest(params.HTTPRequest))), err
				default:
					return mtoshipmentops.NewUpdateMTOShipmentStatusInternalServerError().WithPayload(
						payloads.InternalServerError(nil, h.GetTraceIDFromRequest(params.HTTPRequest))), err
				}
			}
			if !availableToPrime {
				return mtoshipmentops.NewUpdateMTOShipmentStatusInternalServerError().WithPayload(
					payloads.InternalServerError(nil, h.GetTraceIDFromRequest(params.HTTPRequest))), err
			}

			status := models.MTOShipmentStatus(params.Body.Status)
			eTag := params.IfMatch

			featureFlagValues, err := handlers.GetAllDomesticMHFlags(appCtx, h.HandlerConfig.FeatureFlagFetcher())
			if err != nil {
				appCtx.Logger().Panic(fmt.Sprintf("Error fetching mobile home feature flags: %s", err))
			}
			shipment, err := h.updater.UpdateMTOShipmentStatus(appCtx, shipmentID, status, nil, nil, eTag, featureFlagValues)
			if err != nil {
				appCtx.Logger().Error("UpdateMTOShipmentStatusStatus error: ", zap.Error(err))

				switch e := err.(type) {
				case apperror.NotFoundError:
					return mtoshipmentops.NewUpdateMTOShipmentStatusNotFound().WithPayload(
						payloads.ClientError(handlers.NotFoundMessage, err.Error(), h.GetTraceIDFromRequest(params.HTTPRequest))), err
				case apperror.InvalidInputError:
					return mtoshipmentops.NewUpdateMTOShipmentStatusUnprocessableEntity().WithPayload(
						payloads.ValidationError("The input provided did not pass validation.", h.GetTraceIDFromRequest(params.HTTPRequest), e.ValidationErrors)), err
				case apperror.PreconditionFailedError:
					return mtoshipmentops.NewUpdateMTOShipmentStatusPreconditionFailed().WithPayload(
						payloads.ClientError(handlers.PreconditionErrMessage, err.Error(), h.GetTraceIDFromRequest(params.HTTPRequest))), err
				case mtoshipment.ConflictStatusError:
					return mtoshipmentops.NewUpdateMTOShipmentStatusConflict().WithPayload(
						payloads.ClientError(handlers.ConflictErrMessage, err.Error(), h.GetTraceIDFromRequest(params.HTTPRequest))), err
				default:
					return mtoshipmentops.NewUpdateMTOShipmentStatusInternalServerError().WithPayload(
						payloads.InternalServerError(handlers.FmtString(err.Error()), h.GetTraceIDFromRequest(params.HTTPRequest))), err
				}
			}

			return mtoshipmentops.NewUpdateMTOShipmentStatusOK().WithPayload(payloads.MTOShipment(shipment)), nil
		})
}
