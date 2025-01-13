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
	mtoserviceitem "github.com/transcom/mymove/pkg/services/mto_service_item"
	mtoshipment "github.com/transcom/mymove/pkg/services/mto_shipment"
)

// CreateableServiceItemMap is a map of MTOServiceItemModelTypes and their allowed statuses
// THIS WILL NEED TO BE UPDATED AS WE CONTINUE TO ADD MORE SERVICE ITEMS.
// We will eventually remove this when all service items are added.
var CreateableServiceItemMap = map[primemessages.MTOServiceItemModelType]bool{
	primemessages.MTOServiceItemModelTypeMTOServiceItemOriginSIT:            true,
	primemessages.MTOServiceItemModelTypeMTOServiceItemDestSIT:              true,
	primemessages.MTOServiceItemModelTypeMTOServiceItemShuttle:              true,
	primemessages.MTOServiceItemModelTypeMTOServiceItemDomesticCrating:      true,
	primemessages.MTOServiceItemModelTypeMTOServiceItemInternationalCrating: true,
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
				appCtx.Logger().Error("primeapi.CreateMTOServiceItemHandler Error Fetch Shipment", zap.Error(err))
				switch err {
				case models.ErrFetchNotFound:
					return mtoserviceitemops.NewCreateMTOServiceItemNotFound().WithPayload(payloads.ClientError(handlers.NotFoundMessage, "Fetch Shipment", h.GetTraceIDFromRequest(params.HTTPRequest))), err
				default:
					return mtoserviceitemops.NewCreateMTOServiceItemInternalServerError().WithPayload(payloads.InternalServerError(nil, h.GetTraceIDFromRequest(params.HTTPRequest))), err
				}
			}

			if shipment.ShipmentType == models.MTOShipmentTypePPM {
				verrs := validate.NewErrors()
				verrs.Add("mtoShipmentID", params.Body.MtoShipmentID().String())
				appCtx.Logger().Error("primeapi.CreateMTOServiceItemHandler - Create Service Item is not allowed for PPM shipments", zap.Error(verrs))
				return mtoserviceitemops.NewCreateMTOServiceItemUnprocessableEntity().WithPayload(payloads.ValidationError(
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
				delete(CreateableServiceItemMap, primemessages.MTOServiceItemModelTypeMTOServiceItemInternationalCrating)
			}

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
				case apperror.UnprocessableEntityError:
					payload := payloads.ValidationError(
						err.Error(),
						h.GetTraceIDFromRequest(params.HTTPRequest),
						validate.NewErrors())
					return mtoserviceitemops.NewCreateMTOServiceItemUnprocessableEntity().WithPayload(payload), err
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

			// We need to get the shipment in case we need to update the Required Delivery Date for an Origin SIT
			eagerAssociations := []string{"MoveTaskOrder",
				"PickupAddress",
				"DestinationAddress",
				"SecondaryPickupAddress",
				"SecondaryDeliveryAddress",
				"TertiaryPickupAddress",
				"TertiaryDeliveryAddress",
				"MTOServiceItems.ReService",
				"StorageFacility.Address",
				"PPMShipment"}
			serviceItem, err := mtoserviceitem.NewMTOServiceItemFetcher().GetServiceItem(appCtx, mtoServiceItem.ID)

			if err != nil {
				appCtx.Logger().Error("primeapi.UpdateMTOServiceItemHandler error", zap.Error(err))
				return mtoserviceitemops.NewUpdateMTOServiceItemNotFound().WithPayload(
					payloads.ClientError(handlers.NotFoundMessage, err.Error(), h.GetTraceIDFromRequest(params.HTTPRequest))), err
			}

			shipmentID := *serviceItem.MTOShipmentID
			eTag := params.IfMatch
			shipment, err := mtoshipment.NewMTOShipmentFetcher().GetShipment(appCtx, shipmentID, eagerAssociations...)

			if err != nil {
				appCtx.Logger().Error("primeapi.UpdateMTOServiceItemHandler error", zap.Error(err))
				return mtoserviceitemops.NewUpdateMTOServiceItemUnprocessableEntity().WithPayload(payloads.ValidationError(
					verrs.Error(), h.GetTraceIDFromRequest(params.HTTPRequest), verrs)), verrs
			}

			updatedMTOServiceItem, err := h.MTOServiceItemUpdater.UpdateMTOServiceItemPrime(appCtx, mtoServiceItem, h.HandlerConfig.HHGPlanner(), *shipment, eTag)

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
				case apperror.UnprocessableEntityError:
					verrs := validate.NewErrors()
					return mtoserviceitemops.NewUpdateMTOServiceItemUnprocessableEntity().WithPayload(payloads.ValidationError(
						err.Error(), h.GetTraceIDFromRequest(params.HTTPRequest), verrs)), verrs
				case *apperror.UnsupportedPortCodeError:
					verrs := validate.NewErrors()
					return mtoserviceitemops.NewUpdateMTOServiceItemUnprocessableEntity().WithPayload(payloads.ValidationError(
						err.Error(), h.GetTraceIDFromRequest(params.HTTPRequest), verrs)), verrs
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
