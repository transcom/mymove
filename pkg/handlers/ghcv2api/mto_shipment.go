package ghcv2api

import (
	"fmt"

	"github.com/go-openapi/runtime/middleware"
	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	mtoshipmentops "github.com/transcom/mymove/pkg/gen/ghcv2api/ghcv2operations/mto_shipment"
	"github.com/transcom/mymove/pkg/gen/ghcv2messages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/handlers/ghcv2api/internal/payloads"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/event"
	mtoshipment "github.com/transcom/mymove/pkg/services/mto_shipment"
)

// UpdateShipmentHandler updates shipments
type UpdateShipmentHandler struct {
	handlers.HandlerConfig
	services.ShipmentUpdater
	services.ShipmentSITStatus
}

// Handle updates shipments
func (h UpdateShipmentHandler) Handle(params mtoshipmentops.UpdateMTOShipmentParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {

			payload := params.Body
			if payload == nil {
				appCtx.Logger().Error("Invalid mto shipment: params Body is nil")
				emptyBodyError := apperror.NewBadDataError("The MTO Shipment request body cannot be empty.")
				payload := payloadForValidationError(
					"Empty body error",
					emptyBodyError.Error(),
					h.GetTraceIDFromRequest(params.HTTPRequest),
					validate.NewErrors(),
				)

				return mtoshipmentops.NewUpdateMTOShipmentUnprocessableEntity().WithPayload(payload), emptyBodyError
			}

			shipmentID := uuid.FromStringOrNil(params.ShipmentID.String())
			oldShipment, err := mtoshipment.FindShipment(appCtx, shipmentID)

			if err != nil {
				appCtx.Logger().Error("ghcapi.UpdateShipmentHandler", zap.Error(err))
				switch err.(type) {
				case apperror.NotFoundError:
					return mtoshipmentops.NewUpdateMTOShipmentNotFound(), err
				default:
					msg := fmt.Sprintf("%v | Instance: %v", handlers.FmtString(err.Error()), h.GetTraceIDFromRequest(params.HTTPRequest))

					return mtoshipmentops.NewUpdateMTOShipmentInternalServerError().WithPayload(
						&ghcv2messages.Error{Message: &msg},
					), err
				}
			}

			mtoShipment := payloads.MTOShipmentModelFromUpdate(payload)
			mtoShipment.ID = shipmentID
			mtoShipment.ShipmentType = oldShipment.ShipmentType

			//MTOShipmentModelFromUpdate defaults UsesExternalVendor to false if it's nil in the payload
			if payload.UsesExternalVendor == nil {
				mtoShipment.UsesExternalVendor = oldShipment.UsesExternalVendor
			}
			// booleans not passed will update to false
			mtoShipment.Diversion = oldShipment.Diversion

			handleError := func(err error) (middleware.Responder, error) {
				appCtx.Logger().Error("ghcapi.UpdateShipmentHandler", zap.Error(err))

				switch e := err.(type) {
				case apperror.NotFoundError:
					return mtoshipmentops.NewUpdateMTOShipmentNotFound(), err
				case apperror.ForbiddenError:
					msg := fmt.Sprintf("%v | Instance: %v", handlers.FmtString(err.Error()), h.GetTraceIDFromRequest(params.HTTPRequest))
					return mtoshipmentops.NewUpdateMTOShipmentForbidden().WithPayload(
						&ghcv2messages.Error{Message: &msg},
					), err
				case apperror.InvalidInputError:
					return mtoshipmentops.NewUpdateMTOShipmentUnprocessableEntity().WithPayload(
						payloadForValidationError(
							handlers.ValidationErrMessage,
							err.Error(),
							h.GetTraceIDFromRequest(params.HTTPRequest),
							e.ValidationErrors,
						),
					), err
				case apperror.PreconditionFailedError:
					msg := fmt.Sprintf("%v | Instance: %v", handlers.FmtString(err.Error()), h.GetTraceIDFromRequest(params.HTTPRequest))
					return mtoshipmentops.NewUpdateMTOShipmentPreconditionFailed().WithPayload(
						&ghcv2messages.Error{Message: &msg},
					), err
				case apperror.QueryError:
					if e.Unwrap() != nil {
						// If you can unwrap, log the internal error (usually a pq error) for better debugging
						appCtx.Logger().Error("ghcapi.UpdateShipmentHandler error", zap.Error(e.Unwrap()))
					}

					msg := fmt.Sprintf("%v | Instance: %v", handlers.FmtString(err.Error()), h.GetTraceIDFromRequest(params.HTTPRequest))

					return mtoshipmentops.NewUpdateMTOShipmentInternalServerError().WithPayload(
						&ghcv2messages.Error{Message: &msg},
					), err
				default:
					msg := fmt.Sprintf("%v | Instance: %v", handlers.FmtString(err.Error()), h.GetTraceIDFromRequest(params.HTTPRequest))

					return mtoshipmentops.NewUpdateMTOShipmentInternalServerError().WithPayload(
						&ghcv2messages.Error{Message: &msg},
					), err
				}
			}
			updatedMtoShipment, err := h.ShipmentUpdater.UpdateShipment(appCtx, mtoShipment, params.IfMatch)
			if err != nil {
				return handleError(err)
			}

			_, err = event.TriggerEvent(event.Event{
				EndpointKey: event.GhcUpdateMTOShipmentEndpointKey,
				// Endpoint that is being handled
				EventKey:        event.MTOShipmentUpdateEventKey,    // Event that you want to trigger
				UpdatedObjectID: updatedMtoShipment.ID,              // ID of the updated logical object
				MtoID:           updatedMtoShipment.MoveTaskOrderID, // ID of the associated Move
				AppContext:      appCtx,
				TraceID:         h.GetTraceIDFromRequest(params.HTTPRequest),
			})
			// If the event trigger fails, just log the error.
			if err != nil {
				appCtx.Logger().Error("ghcapi.UpdateMTOShipment could not generate the event")
			}

			shipmentSITStatus, err := h.CalculateShipmentSITStatus(appCtx, *updatedMtoShipment)
			if err != nil {
				return handleError(err)
			}
			sitStatusPayload := payloads.SITStatus(shipmentSITStatus, h.FileStorer())

			returnPayload := payloads.MTOShipment(h.FileStorer(), updatedMtoShipment, sitStatusPayload)
			return mtoshipmentops.NewUpdateMTOShipmentOK().WithPayload(returnPayload), nil
		})
}

func payloadForValidationError(title string, detail string, instance uuid.UUID, validationErrors *validate.Errors) *ghcv2messages.ValidationError {
	payload := &ghcv2messages.ValidationError{
		ClientError: *payloadForClientError(title, detail, instance),
	}

	if validationErrors != nil {
		payload.InvalidFields = handlers.NewValidationErrorsResponse(validationErrors).Errors
	}

	return payload
}

func payloadForClientError(title string, detail string, instance uuid.UUID) *ghcv2messages.ClientError {
	return &ghcv2messages.ClientError{
		Title:    handlers.FmtString(title),
		Detail:   handlers.FmtString(detail),
		Instance: handlers.FmtUUID(instance),
	}
}
