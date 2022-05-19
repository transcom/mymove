package primeapi

import (
	"fmt"
	"reflect"

	"github.com/go-openapi/runtime/middleware"
	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	mtoserviceitemops "github.com/transcom/mymove/pkg/gen/primeapi/primeoperations/mto_service_item"
	"github.com/transcom/mymove/pkg/gen/primemessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/handlers/primeapi/payloads"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

// CreateableServiceItemMap is a map of MTOServiceItemModelTypes and their allowed statuses
// THIS WILL NEED TO BE UPDATED AS WE CONTINUE TO ADD MORE SERVICE ITEMS.
// We will eventually remove this when all service items are added.
var CreateableServiceItemMap = map[primemessages.MTOServiceItemModelType]bool{
	primemessages.MTOServiceItemModelTypeMTOServiceItemOriginSIT:       true,
	primemessages.MTOServiceItemModelTypeMTOServiceItemDestSIT:         true,
	primemessages.MTOServiceItemModelTypeMTOServiceItemShuttle:         true,
	primemessages.MTOServiceItemModelTypeMTOServiceItemDomesticCrating: true,
}

// CreateMTOServiceItemHandler is the handler to create MTO service items
type CreateMTOServiceItemHandler struct {
	handlers.HandlerConfig
	mtoServiceItemCreator  services.MTOServiceItemCreator
	mtoAvailabilityChecker services.MoveTaskOrderChecker
}

// Handle handler that creates a mto service item
func (h CreateMTOServiceItemHandler) Handle(params mtoserviceitemops.CreateMTOServiceItemParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {

			// restrict creation to a list
			if _, ok := CreateableServiceItemMap[params.Body.ModelType()]; !ok {
				// throw error if modelType() not on the list
				mapKeys := GetMapKeys(CreateableServiceItemMap)
				detailErr := fmt.Sprintf("MTOServiceItem modelType() not allowed: %s ", params.Body.ModelType())
				verrs := validate.NewErrors()
				verrs.Add("modelType", fmt.Sprintf("allowed modelType() %v", mapKeys))

				appCtx.Logger().Error("primeapi.CreateMTOServiceItemHandler error", zap.Error(verrs))
				return mtoserviceitemops.NewCreateMTOServiceItemUnprocessableEntity().WithPayload(payloads.ValidationError(
					detailErr, h.GetTraceIDFromRequest(params.HTTPRequest), verrs)), verrs
			}

			// validation errors passed back if any
			mtoServiceItem, verrs := payloads.MTOServiceItemModel(params.Body)

			if verrs != nil && verrs.HasAny() {
				return mtoserviceitemops.NewCreateMTOServiceItemUnprocessableEntity().WithPayload(payloads.ValidationError(
					"Invalid input found in service item", h.GetTraceIDFromRequest(params.HTTPRequest), verrs)), verrs
			} else if mtoServiceItem == nil {
				return mtoserviceitemops.NewCreateMTOServiceItemUnprocessableEntity().WithPayload(
					payloads.ValidationError("Unable to process service item", h.GetTraceIDFromRequest(params.HTTPRequest), nil)), verrs
			}

			moveTaskOrderID := uuid.FromStringOrNil(mtoServiceItem.MoveTaskOrderID.String())
			mtoAvailableToPrime, err := h.mtoAvailabilityChecker.MTOAvailableToPrime(appCtx, moveTaskOrderID)
			var mtoServiceItems *models.MTOServiceItems

			if mtoAvailableToPrime {
				mtoServiceItem.Status = models.MTOServiceItemStatusSubmitted
				mtoServiceItems, verrs, err = h.mtoServiceItemCreator.CreateMTOServiceItem(appCtx, mtoServiceItem)
			} else if err == nil {
				primeErr := apperror.NewNotFoundError(moveTaskOrderID, "primeapi.CreateMTOServiceItemHandler error - MTO is not available to Prime")
				appCtx.Logger().Error(primeErr.Error())
				return mtoserviceitemops.NewCreateMTOServiceItemNotFound().WithPayload(payloads.ClientError(
					handlers.NotFoundMessage, fmt.Sprintf("id: %s not found for moveTaskOrder", moveTaskOrderID), h.GetTraceIDFromRequest(params.HTTPRequest))), primeErr
			}

			if verrs != nil && verrs.HasAny() {
				return mtoserviceitemops.NewCreateMTOServiceItemUnprocessableEntity().WithPayload(payloads.ValidationError(
					verrs.Error(), h.GetTraceIDFromRequest(params.HTTPRequest), verrs)), verrs
			}

			// Could be the error from MTOAvailableToPrime or CreateMTOServiceItem:
			if err != nil {
				appCtx.Logger().Error("primeapi.CreateMTOServiceItemHandler error", zap.Error(err))
				switch e := err.(type) {
				case apperror.NotFoundError:
					return mtoserviceitemops.NewCreateMTOServiceItemNotFound().WithPayload(payloads.ClientError(handlers.NotFoundMessage, err.Error(), h.GetTraceIDFromRequest(params.HTTPRequest))), err
				case apperror.InvalidInputError:
					return mtoserviceitemops.NewCreateMTOServiceItemUnprocessableEntity().WithPayload(payloads.ValidationError(e.Error(), h.GetTraceIDFromRequest(params.HTTPRequest), e.ValidationErrors)), err
				case apperror.ConflictError:
					return mtoserviceitemops.NewCreateMTOServiceItemConflict().WithPayload(payloads.ClientError(handlers.ConflictErrMessage, err.Error(), h.GetTraceIDFromRequest(params.HTTPRequest))), err
				case apperror.QueryError:
					if e.Unwrap() != nil {
						// If you can unwrap, log the internal error (usually a pq error) for better debugging
						appCtx.Logger().Error("primeapi.CreateMTOServiceItemHandler query error", zap.Error(e.Unwrap()))
					}
					return mtoserviceitemops.NewCreateMTOServiceItemInternalServerError().WithPayload(payloads.InternalServerError(nil, h.GetTraceIDFromRequest(params.HTTPRequest))), err
				default:
					return mtoserviceitemops.NewCreateMTOServiceItemInternalServerError().WithPayload(payloads.InternalServerError(nil, h.GetTraceIDFromRequest(params.HTTPRequest))), err
				}
			}

			mtoServiceItemsPayload := *payloads.MTOServiceItems(mtoServiceItems)
			return mtoserviceitemops.NewCreateMTOServiceItemOK().WithPayload(mtoServiceItemsPayload), nil
		})
}

// UpdateMTOServiceItemHandler is the handler to update MTO service items
type UpdateMTOServiceItemHandler struct {
	handlers.HandlerConfig
	services.MTOServiceItemUpdater
}

// Handle handler that updates an MTOServiceItem. Only a limited number of service items and fields may be updated.
func (h UpdateMTOServiceItemHandler) Handle(params mtoserviceitemops.UpdateMTOServiceItemParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {

			mtoServiceItem, verrs := payloads.MTOServiceItemModelFromUpdate(params.MtoServiceItemID, params.Body)
			if verrs != nil && verrs.HasAny() {
				return mtoserviceitemops.NewUpdateMTOServiceItemUnprocessableEntity().WithPayload(payloads.ValidationError(
					verrs.Error(), h.GetTraceIDFromRequest(params.HTTPRequest), verrs)), verrs
			}

			eTag := params.IfMatch
			updatedMTOServiceItem, err := h.MTOServiceItemUpdater.UpdateMTOServiceItemPrime(appCtx, mtoServiceItem, eTag)

			if err != nil {
				appCtx.Logger().Error("primeapi.UpdateMTOServiceItemHandler error", zap.Error(err))
				switch e := err.(type) {
				case apperror.NotFoundError:
					return mtoserviceitemops.NewUpdateMTOServiceItemNotFound().WithPayload(
						payloads.ClientError(handlers.NotFoundMessage, err.Error(), h.GetTraceIDFromRequest(params.HTTPRequest))), err
				case apperror.InvalidInputError:
					return mtoserviceitemops.NewUpdateMTOServiceItemUnprocessableEntity().WithPayload(
						payloads.ValidationError(e.Error(), h.GetTraceIDFromRequest(params.HTTPRequest), e.ValidationErrors)), err
				case apperror.ConflictError:
					return mtoserviceitemops.NewUpdateMTOServiceItemConflict().WithPayload(
						payloads.ClientError(handlers.ConflictErrMessage, err.Error(), h.GetTraceIDFromRequest(params.HTTPRequest))), err
				case apperror.PreconditionFailedError:
					return mtoserviceitemops.NewUpdateMTOServiceItemPreconditionFailed().WithPayload(
						payloads.ClientError(handlers.PreconditionErrMessage, err.Error(), h.GetTraceIDFromRequest(params.HTTPRequest))), err
				case apperror.QueryError:
					if e.Unwrap() != nil {
						// If you can unwrap, log the internal error (usually a pq error) for better debugging
						appCtx.Logger().Error("primeapi.UpdateMTOServiceItemHandler query error", zap.Error(e.Unwrap()))
					}
					return mtoserviceitemops.NewUpdateMTOServiceItemInternalServerError().WithPayload(
						payloads.InternalServerError(nil, h.GetTraceIDFromRequest(params.HTTPRequest))), err
				default:
					return mtoserviceitemops.NewUpdateMTOServiceItemInternalServerError().WithPayload(
						payloads.InternalServerError(nil, h.GetTraceIDFromRequest(params.HTTPRequest))), err
				}
			}

			return mtoserviceitemops.NewUpdateMTOServiceItemOK().WithPayload(payloads.MTOServiceItem(updatedMTOServiceItem)), nil
		})
}

// GetMapKeys is a helper function that returns the keys that are MTOServiceItemModelTypes from the map
func GetMapKeys(m map[primemessages.MTOServiceItemModelType]bool) []reflect.Value {
	return reflect.ValueOf(m).MapKeys()
}
