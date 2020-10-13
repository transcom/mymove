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

// AllowedServiceItemMap is a map of MTOServiceItemModelTypes and their allowed statuses
// THIS WILL NEED TO BE UPDATED AS WE CONTINUE TO ADD MORE SERVICE ITEMS.
// We will eventually remove this when all service items are added.
var AllowedServiceItemMap = map[primemessages.MTOServiceItemModelType]bool{
	primemessages.MTOServiceItemModelTypeMTOServiceItemDOFSIT:          true,
	primemessages.MTOServiceItemModelTypeMTOServiceItemDDFSIT:          true,
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
	if _, ok := AllowedServiceItemMap[params.Body.ModelType()]; !ok {
		// throw error if modelType() not on the list
		mapKeys := GetMapKeys(AllowedServiceItemMap)
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

// GetMapKeys is a helper function that returns the keys that are MTOServiceItemModelTypes from the map
func GetMapKeys(m map[primemessages.MTOServiceItemModelType]bool) []reflect.Value {
	return reflect.ValueOf(m).MapKeys()
}
