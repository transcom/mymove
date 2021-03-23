package primeapi

import (
	"fmt"
	"reflect"

	"github.com/transcom/mymove/pkg/handlers/primeapi/payloads"

	"github.com/transcom/mymove/pkg/models"

	"github.com/gofrs/uuid"

	"github.com/go-openapi/runtime/middleware"
	"github.com/gobuffalo/validate/v3"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/gen/primemessages"

	mtoserviceitemops "github.com/transcom/mymove/pkg/gen/primeapi/primeoperations/mto_service_item"
	"github.com/transcom/mymove/pkg/handlers"
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

// CreateMTOServiceItemHandler is the handler to update MTO shipments
type CreateMTOServiceItemHandler struct {
	handlers.HandlerContext
	mtoServiceItemCreator  services.MTOServiceItemCreator
	mtoAvailabilityChecker services.MoveTaskOrderChecker
}

// Handle handler that updates a mto shipment
func (h CreateMTOServiceItemHandler) Handle(params mtoserviceitemops.CreateMTOServiceItemParams) middleware.Responder {
	logger := h.LoggerFromRequest(params.HTTPRequest)

	// restrict creation to a list
	if _, ok := CreateableServiceItemMap[params.Body.ModelType()]; !ok {
		// throw error if modelType() not on the list
		mapKeys := GetMapKeys(CreateableServiceItemMap)
		detailErr := fmt.Sprintf("MTOServiceItem modelType() not allowed: %s ", params.Body.ModelType())
		verrs := validate.NewErrors()
		verrs.Add("modelType", fmt.Sprintf("allowed modelType() %v", mapKeys))

		logger.Error("primeapi.CreateMTOServiceItemHandler error", zap.Error(verrs))
		return mtoserviceitemops.NewCreateMTOServiceItemUnprocessableEntity().WithPayload(payloads.ValidationError(
			detailErr, h.GetTraceID(), verrs))
	}

	// validation errors passed back if any
	mtoServiceItem, verrs := payloads.MTOServiceItemModel(params.Body)

	if verrs != nil && verrs.HasAny() {
		return mtoserviceitemops.NewCreateMTOServiceItemUnprocessableEntity().WithPayload(payloads.ValidationError(
			verrs.Error(), h.GetTraceID(), verrs))
	} else if mtoServiceItem == nil {
		return mtoserviceitemops.NewCreateMTOServiceItemUnprocessableEntity().WithPayload(
			payloads.ValidationError("Unable to process service item", h.GetTraceID(), nil))
	}

	moveTaskOrderID := uuid.FromStringOrNil(mtoServiceItem.MoveTaskOrderID.String())
	mtoAvailableToPrime, err := h.mtoAvailabilityChecker.MTOAvailableToPrime(moveTaskOrderID)
	var mtoServiceItems *models.MTOServiceItems

	if mtoAvailableToPrime {
		mtoServiceItem.Status = models.MTOServiceItemStatusSubmitted
		mtoServiceItems, verrs, err = h.mtoServiceItemCreator.CreateMTOServiceItem(mtoServiceItem)
	} else if err == nil {
		logger.Error("primeapi.CreateMTOServiceItemHandler error - MTO is not available to Prime")
		return mtoserviceitemops.NewCreateMTOServiceItemNotFound().WithPayload(payloads.ClientError(
			handlers.NotFoundMessage, fmt.Sprintf("id: %s not found for moveTaskOrder", moveTaskOrderID), h.GetTraceID()))
	}

	if verrs != nil && verrs.HasAny() {
		return mtoserviceitemops.NewCreateMTOServiceItemUnprocessableEntity().WithPayload(payloads.ValidationError(
			verrs.Error(), h.GetTraceID(), verrs))
	}

	// Could be the error from MTOAvailableToPrime or CreateMTOServiceItem:
	if err != nil {
		logger.Error("primeapi.CreateMTOServiceItemHandler error", zap.Error(err))
		switch e := err.(type) {
		case services.NotFoundError:
			return mtoserviceitemops.NewCreateMTOServiceItemNotFound().WithPayload(payloads.ClientError(handlers.NotFoundMessage, err.Error(), h.GetTraceID()))
		case services.InvalidInputError:
			return mtoserviceitemops.NewCreateMTOServiceItemUnprocessableEntity().WithPayload(payloads.ValidationError(e.Error(), h.GetTraceID(), e.ValidationErrors))
		case services.ConflictError:
			return mtoserviceitemops.NewCreateMTOServiceItemConflict().WithPayload(payloads.ClientError(handlers.ConflictErrMessage, err.Error(), h.GetTraceID()))
		case services.QueryError:
			if e.Unwrap() != nil {
				// If you can unwrap, log the internal error (usually a pq error) for better debugging
				logger.Error("primeapi.CreateMTOServiceItemHandler query error", zap.Error(e.Unwrap()))
			}
			return mtoserviceitemops.NewCreateMTOServiceItemInternalServerError().WithPayload(payloads.InternalServerError(nil, h.GetTraceID()))
		default:
			return mtoserviceitemops.NewCreateMTOServiceItemInternalServerError().WithPayload(payloads.InternalServerError(nil, h.GetTraceID()))
		}
	}

	mtoServiceItemsPayload := *payloads.MTOServiceItems(mtoServiceItems)
	return mtoserviceitemops.NewCreateMTOServiceItemOK().WithPayload(mtoServiceItemsPayload)
}

// UpdateMTOServiceItemHandler is the handler to update MTO shipments
type UpdateMTOServiceItemHandler struct {
	handlers.HandlerContext
	services.MTOServiceItemUpdater
}

// Handle handler that updates an MTOServiceItem. Only a limited number of service items and fields may be updated.
func (h UpdateMTOServiceItemHandler) Handle(params mtoserviceitemops.UpdateMTOServiceItemParams) middleware.Responder {
	logger := h.LoggerFromRequest(params.HTTPRequest)

	mtoServiceItem, verrs := payloads.MTOServiceItemModelFromUpdate(params.MtoServiceItemID, params.Body)
	if verrs != nil && verrs.HasAny() {
		return mtoserviceitemops.NewUpdateMTOServiceItemUnprocessableEntity().WithPayload(payloads.ValidationError(
			verrs.Error(), h.GetTraceID(), verrs))
	}

	eTag := params.IfMatch
	updatedMTOServiceItem, err := h.MTOServiceItemUpdater.UpdateMTOServiceItemPrime(h.DB(), mtoServiceItem, eTag)

	if err != nil {
		logger.Error("primeapi.UpdateMTOServiceItemHandler error", zap.Error(err))
		switch e := err.(type) {
		case services.NotFoundError:
			return mtoserviceitemops.NewUpdateMTOServiceItemNotFound().WithPayload(payloads.ClientError(handlers.NotFoundMessage, err.Error(), h.GetTraceID()))
		case services.InvalidInputError:
			return mtoserviceitemops.NewUpdateMTOServiceItemUnprocessableEntity().WithPayload(payloads.ValidationError(e.Error(), h.GetTraceID(), e.ValidationErrors))
		case services.ConflictError:
			return mtoserviceitemops.NewUpdateMTOServiceItemConflict().WithPayload(payloads.ClientError(handlers.ConflictErrMessage, err.Error(), h.GetTraceID()))
		case services.PreconditionFailedError:
			return mtoserviceitemops.NewUpdateMTOServiceItemPreconditionFailed().WithPayload(payloads.ClientError(handlers.PreconditionErrMessage, err.Error(), h.GetTraceID()))
		case services.QueryError:
			if e.Unwrap() != nil {
				// If you can unwrap, log the internal error (usually a pq error) for better debugging
				logger.Error("primeapi.UpdateMTOServiceItemHandler query error", zap.Error(e.Unwrap()))
			}
			return mtoserviceitemops.NewUpdateMTOServiceItemInternalServerError().WithPayload(payloads.InternalServerError(nil, h.GetTraceID()))
		default:
			return mtoserviceitemops.NewUpdateMTOServiceItemInternalServerError().WithPayload(payloads.InternalServerError(nil, h.GetTraceID()))
		}
	}

	return mtoserviceitemops.NewUpdateMTOServiceItemOK().WithPayload(payloads.MTOServiceItem(updatedMTOServiceItem))
}

// GetMapKeys is a helper function that returns the keys that are MTOServiceItemModelTypes from the map
func GetMapKeys(m map[primemessages.MTOServiceItemModelType]bool) []reflect.Value {
	return reflect.ValueOf(m).MapKeys()
}
