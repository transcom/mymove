package primeapi

import (
	"fmt"
	"reflect"

	"github.com/gofrs/uuid"

	"github.com/go-openapi/runtime/middleware"
	"github.com/gobuffalo/validate"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/gen/primemessages"

	mtoserviceitemops "github.com/transcom/mymove/pkg/gen/primeapi/primeoperations/mto_service_item"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/handlers/primeapi/internal/payloads"
	"github.com/transcom/mymove/pkg/services"
)

// THIS WILL NEED TO BE UPDATED AS WE CONTINUE TO ADD MORE SERVICE ITEMS.
// We will eventually remove this when all service items are added.
var allowedServiceItemMap = map[primemessages.MTOServiceItemModelType]bool{
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
	if _, ok := allowedServiceItemMap[params.Body.ModelType()]; !ok {
		// throw error if modelType() not on the list
		mapKeys := getMapKeys(allowedServiceItemMap)
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
	mtoAvailableToPrime, mtoErr := h.mtoAvailabilityChecker.MTOAvailableToPrime(moveTaskOrderID)
	if mtoErr != nil {
		logger.Error("primeapi.CreateMTOServiceItemHandler error - checking MTO availability to Prime", zap.Error(mtoErr))
		switch e := mtoErr.(type) {
		case services.NotFoundError:
			return mtoserviceitemops.NewCreateMTOServiceItemNotFound().WithPayload(
				payloads.ClientError(handlers.NotFoundMessage, e.Error(), h.GetTraceID()))
		default:
			return mtoserviceitemops.NewCreateMTOServiceItemInternalServerError().WithPayload(payloads.InternalServerError(nil, h.GetTraceID()))
		}
	} else if !mtoAvailableToPrime {
		logger.Error("primeapi.CreateMTOServiceItemHandler error - MTO is not available to Prime", zap.Error(mtoErr))
		return mtoserviceitemops.NewCreateMTOServiceItemNotFound().WithPayload(payloads.ClientError(
			handlers.NotFoundMessage, fmt.Sprintf("id: %s not found for moveTaskOrder", moveTaskOrderID), h.GetTraceID()))
	}

	mtoServiceItem, verrs, err := h.mtoServiceItemCreator.CreateMTOServiceItem(mtoServiceItem)
	if verrs != nil && verrs.HasAny() {
		return mtoserviceitemops.NewCreateMTOServiceItemUnprocessableEntity().WithPayload(payloads.ValidationError(
			verrs.Error(), h.GetTraceID(), verrs))
	}

	if err != nil {
		logger.Error("primeapi.CreateMTOServiceItemHandler error", zap.Error(err))
		switch e := err.(type) {
		case services.NotFoundError:
			return mtoserviceitemops.NewCreateMTOServiceItemNotFound().WithPayload(payloads.ClientError(handlers.NotFoundMessage, err.Error(), h.GetTraceID()))
		case services.InvalidInputError:
			return mtoserviceitemops.NewCreateMTOServiceItemUnprocessableEntity().WithPayload(payloads.ValidationError(e.Error(), h.GetTraceID(), e.ValidationErrors))
		case services.ConflictError:
			return mtoserviceitemops.NewCreateMTOServiceItemConflict().WithPayload(payloads.ClientError(handlers.ConflictErrMessage, e.Error(), h.GetTraceID()))
		case services.QueryError:
			if e.Unwrap() != nil {
				// If you can unwrap, log the internal error (usually a pq error) for better debugging
				logger.Error("primeapi.CreateMTOServiceItemHandler error", zap.Error(e.Unwrap()))
			}
			return mtoserviceitemops.NewCreateMTOServiceItemInternalServerError().WithPayload(payloads.InternalServerError(nil, h.GetTraceID()))
		default:
			return mtoserviceitemops.NewCreateMTOServiceItemInternalServerError().WithPayload(payloads.InternalServerError(nil, h.GetTraceID()))
		}
	}
	mtoServiceItemPayload := payloads.MTOServiceItem(mtoServiceItem)
	return mtoserviceitemops.NewCreateMTOServiceItemOK().WithPayload(mtoServiceItemPayload)
}

// helper to get the keys from a map
func getMapKeys(m map[primemessages.MTOServiceItemModelType]bool) []reflect.Value {
	return reflect.ValueOf(m).MapKeys()
}
