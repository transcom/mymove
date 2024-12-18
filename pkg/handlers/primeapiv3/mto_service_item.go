package primeapiv3

import (
	"fmt"

	"github.com/go-openapi/runtime/middleware"
	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	mtoserviceitemops "github.com/transcom/mymove/pkg/gen/primeapi/primeoperations/mto_service_item"
	"github.com/transcom/mymove/pkg/gen/primev3messages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/handlers/primeapi"
	primeapipayloads "github.com/transcom/mymove/pkg/handlers/primeapi/payloads"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

// CreateableServiceItemMap is a map of MTOServiceItemModelTypes and their allowed statuses
// THIS WILL NEED TO BE UPDATED AS WE CONTINUE TO ADD MORE SERVICE ITEMS.
// We will eventually remove this when all service items are added.
var CreateableServiceItemMap = map[primev3messages.MTOServiceItemModelType]bool{
	primev3messages.MTOServiceItemModelTypeMTOServiceItemOriginSIT:            true,
	primev3messages.MTOServiceItemModelTypeMTOServiceItemDestSIT:              true,
	primev3messages.MTOServiceItemModelTypeMTOServiceItemShuttle:              true,
	primev3messages.MTOServiceItemModelTypeMTOServiceItemDomesticCrating:      true,
	primev3messages.MTOServiceItemModelTypeMTOServiceItemInternationalCrating: true,
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

			// ** Create service item can not be done for PPM shipment **/
			shipment, err := models.FetchShipmentByID(appCtx.DB(), uuid.FromStringOrNil(params.Body.MtoShipmentID().String()))
			if err != nil {
				appCtx.Logger().Error("primeapi.CreateMTOServiceItemHandler.v3 Error Fetch Shipment", zap.Error(err))
				switch err {
				case models.ErrFetchNotFound:
					return mtoserviceitemops.NewCreateMTOServiceItemNotFound().WithPayload(primeapipayloads.ClientError(handlers.NotFoundMessage, "Fetch Shipment", h.GetTraceIDFromRequest(params.HTTPRequest))), err
				default:
					return mtoserviceitemops.NewCreateMTOServiceItemInternalServerError().WithPayload(primeapipayloads.InternalServerError(nil, h.GetTraceIDFromRequest(params.HTTPRequest))), err
				}
			}

			if shipment.ShipmentType == models.MTOShipmentTypePPM {
				verrs := validate.NewErrors()
				verrs.Add("mtoShipmentID", params.Body.MtoShipmentID().String())
				appCtx.Logger().Error("primeapi.CreateMTOServiceItemHandler.v3 - Create Service Item is not allowed for PPM shipments", zap.Error(verrs))
				return mtoserviceitemops.NewCreateMTOServiceItemUnprocessableEntity().WithPayload(primeapipayloads.ValidationError(
					"Create Service Item is not allowed for PPM shipments", h.GetTraceIDFromRequest(params.HTTPRequest), verrs)), verrs
			}

			/** Feature Flag - Alaska **/
			isAlaskaEnabled := false
			featureFlagName := "enable_alaska"
			flag, err := h.FeatureFlagFetcher().GetBooleanFlag(params.HTTPRequest.Context(), appCtx.Logger(), "", featureFlagName, map[string]string{})
			if err != nil {
				appCtx.Logger().Error("Error fetching feature flag", zap.String("featureFlagKey", featureFlagName), zap.Error(err))
			} else {
				isAlaskaEnabled = flag.Match
			}

			/** Turn on/off international crating/uncrating service items **/
			if !isAlaskaEnabled {
				delete(CreateableServiceItemMap, primev3messages.MTOServiceItemModelTypeMTOServiceItemInternationalCrating)
			}

			// restrict creation to a list
			if _, ok := primeapi.CreateableServiceItemMap[params.Body.ModelType()]; !ok {
				// throw error if modelType() not on the list
				mapKeys := primeapi.GetMapKeys(primeapi.CreateableServiceItemMap)
				detailErr := fmt.Sprintf("MTOServiceItem modelType() not allowed: %s ", params.Body.ModelType())
				verrs := validate.NewErrors()
				verrs.Add("modelType", fmt.Sprintf("allowed modelType() %v", mapKeys))

				appCtx.Logger().Error("primeapi.CreateMTOServiceItemHandler.v3 error", zap.Error(verrs))
				return mtoserviceitemops.NewCreateMTOServiceItemUnprocessableEntity().WithPayload(primeapipayloads.ValidationError(
					detailErr, h.GetTraceIDFromRequest(params.HTTPRequest), verrs)), verrs
			}

			// validation errors passed back if any
			mtoServiceItem, verrs := primeapipayloads.MTOServiceItemModel(params.Body)

			if verrs != nil && verrs.HasAny() {
				return mtoserviceitemops.NewCreateMTOServiceItemUnprocessableEntity().WithPayload(primeapipayloads.ValidationError(
					"Invalid input found in service item", h.GetTraceIDFromRequest(params.HTTPRequest), verrs)), verrs
			} else if mtoServiceItem == nil {
				return mtoserviceitemops.NewCreateMTOServiceItemUnprocessableEntity().WithPayload(
					primeapipayloads.ValidationError("Unable to process service item", h.GetTraceIDFromRequest(params.HTTPRequest), nil)), verrs
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
				return mtoserviceitemops.NewCreateMTOServiceItemNotFound().WithPayload(primeapipayloads.ClientError(
					handlers.NotFoundMessage, fmt.Sprintf("id: %s not found for moveTaskOrder", moveTaskOrderID), h.GetTraceIDFromRequest(params.HTTPRequest))), primeErr
			}

			if verrs != nil && verrs.HasAny() {
				return mtoserviceitemops.NewCreateMTOServiceItemUnprocessableEntity().WithPayload(primeapipayloads.ValidationError(
					verrs.Error(), h.GetTraceIDFromRequest(params.HTTPRequest), verrs)), verrs
			}

			// Could be the error from MTOAvailableToPrime or CreateMTOServiceItem:
			if err != nil {
				appCtx.Logger().Error("primeapi.CreateMTOServiceItemHandler error", zap.Error(err))
				switch e := err.(type) {
				case apperror.NotFoundError:
					return mtoserviceitemops.NewCreateMTOServiceItemNotFound().WithPayload(primeapipayloads.ClientError(handlers.NotFoundMessage, err.Error(), h.GetTraceIDFromRequest(params.HTTPRequest))), err
				case apperror.InvalidInputError:
					return mtoserviceitemops.NewCreateMTOServiceItemUnprocessableEntity().WithPayload(primeapipayloads.ValidationError(e.Error(), h.GetTraceIDFromRequest(params.HTTPRequest), e.ValidationErrors)), err
				case apperror.ConflictError:
					return mtoserviceitemops.NewCreateMTOServiceItemConflict().WithPayload(primeapipayloads.ClientError(handlers.ConflictErrMessage, err.Error(), h.GetTraceIDFromRequest(params.HTTPRequest))), err
				case apperror.QueryError:
					if e.Unwrap() != nil {
						// If you can unwrap, log the internal error (usually a pq error) for better debugging
						appCtx.Logger().Error("primeapi.CreateMTOServiceItemHandler query error", zap.Error(e.Unwrap()))
					}
					return mtoserviceitemops.NewCreateMTOServiceItemInternalServerError().WithPayload(primeapipayloads.InternalServerError(nil, h.GetTraceIDFromRequest(params.HTTPRequest))), err
				default:
					return mtoserviceitemops.NewCreateMTOServiceItemInternalServerError().WithPayload(primeapipayloads.InternalServerError(nil, h.GetTraceIDFromRequest(params.HTTPRequest))), err
				}
			}

			mtoServiceItemsPayload := *primeapipayloads.MTOServiceItems(mtoServiceItems)
			return mtoserviceitemops.NewCreateMTOServiceItemOK().WithPayload(mtoServiceItemsPayload), nil
		})
}
