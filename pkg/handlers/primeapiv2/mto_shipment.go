package primeapiv2

import (
	"fmt"

	"github.com/go-openapi/runtime/middleware"
	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	mtoshipmentops "github.com/transcom/mymove/pkg/gen/primev2api/primev2operations/mto_shipment"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/handlers/primeapi"
	"github.com/transcom/mymove/pkg/handlers/primeapiv2/payloads"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

// CreateMTOShipmentHandler is the handler to create MTO shipments
type CreateMTOShipmentHandler struct {
	handlers.HandlerConfig
	services.ShipmentCreator
	mtoAvailabilityChecker services.MoveTaskOrderChecker
}

// Handle creates the mto shipment
func (h CreateMTOShipmentHandler) Handle(params mtoshipmentops.CreateMTOShipmentParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {
			payload := params.Body
			if payload == nil {
				err := apperror.NewBadDataError("the MTO Shipment request body cannot be empty")
				appCtx.Logger().Error(err.Error())
				return mtoshipmentops.NewCreateMTOShipmentBadRequest().WithPayload(payloads.ClientError(handlers.BadRequestErrMessage,
					err.Error(), h.GetTraceIDFromRequest(params.HTTPRequest))), err
			}

			for _, mtoServiceItem := range params.Body.MtoServiceItems() {
				// restrict creation to a list
				if _, ok := CreateableServiceItemMap[mtoServiceItem.ModelType()]; !ok {
					// throw error if modelType() not on the list
					mapKeys := primeapi.GetMapKeys(primeapi.CreateableServiceItemMap)
					detailErr := fmt.Sprintf("MTOServiceItem modelType() not allowed: %s ", mtoServiceItem.ModelType())
					verrs := validate.NewErrors()
					verrs.Add("modelType", fmt.Sprintf("allowed modelType() %v", mapKeys))

					appCtx.Logger().Error("primeapiv2.CreateMTOShipmentHandler error", zap.Error(verrs))
					return mtoshipmentops.NewCreateMTOShipmentUnprocessableEntity().WithPayload(payloads.ValidationError(
						detailErr, h.GetTraceIDFromRequest(params.HTTPRequest), verrs)), verrs
				}
			}

			mtoShipment := payloads.MTOShipmentModelFromCreate(payload)
			mtoShipment.Status = models.MTOShipmentStatusSubmitted
			mtoServiceItemsList, verrs := payloads.MTOServiceItemModelListFromCreate(payload)

			if verrs != nil && verrs.HasAny() {
				appCtx.Logger().Error("Error validating mto service item list: ", zap.Error(verrs))

				return mtoshipmentops.NewCreateMTOShipmentUnprocessableEntity().WithPayload(payloads.ValidationError(
					"The MTO service item list is invalid.", h.GetTraceIDFromRequest(params.HTTPRequest), nil)), verrs
			}

			mtoShipment.MTOServiceItems = mtoServiceItemsList

			moveTaskOrderID := uuid.FromStringOrNil(payload.MoveTaskOrderID.String())
			mtoAvailableToPrime, err := h.mtoAvailabilityChecker.MTOAvailableToPrime(appCtx, moveTaskOrderID)

			if mtoAvailableToPrime {
				mtoShipment, err = h.ShipmentCreator.CreateShipment(appCtx, mtoShipment)
			} else if err == nil {
				appCtx.Logger().Error("primeapiv2.CreateMTOShipmentHandler error - MTO is not available to Prime")
				return mtoshipmentops.NewCreateMTOShipmentNotFound().WithPayload(payloads.ClientError(
					handlers.NotFoundMessage, fmt.Sprintf("id: %s not found for moveTaskOrder", moveTaskOrderID), h.GetTraceIDFromRequest(params.HTTPRequest))), err
			}

			// Could be the error from MTOAvailableToPrime or CreateMTOShipment:
			if err != nil {
				appCtx.Logger().Error("primeapiv2.CreateMTOShipmentHandler error", zap.Error(err))
				switch e := err.(type) {
				case apperror.NotFoundError:
					return mtoshipmentops.NewCreateMTOShipmentNotFound().WithPayload(
						payloads.ClientError(handlers.NotFoundMessage, err.Error(), h.GetTraceIDFromRequest(params.HTTPRequest))), err
				case apperror.InvalidInputError:
					return mtoshipmentops.NewCreateMTOShipmentUnprocessableEntity().WithPayload(
						payloads.ValidationError(err.Error(), h.GetTraceIDFromRequest(params.HTTPRequest), e.ValidationErrors)), err
				case apperror.QueryError:
					if e.Unwrap() != nil {
						// If you can unwrap, log the internal error (usually a pq error) for better debugging
						appCtx.Logger().Error("primeapiv2.CreateMTOShipmentHandler query error", zap.Error(e.Unwrap()))
					}
					return mtoshipmentops.NewCreateMTOShipmentInternalServerError().WithPayload(
						payloads.InternalServerError(nil, h.GetTraceIDFromRequest(params.HTTPRequest))), err
				default:
					return mtoshipmentops.NewCreateMTOShipmentInternalServerError().WithPayload(
						payloads.InternalServerError(nil, h.GetTraceIDFromRequest(params.HTTPRequest))), err
				}
			}
			returnPayload := payloads.MTOShipment(mtoShipment)
			return mtoshipmentops.NewCreateMTOShipmentOK().WithPayload(returnPayload), nil
		})
}
