package ghcapi

import (
	"fmt"
	"time"

	"github.com/go-openapi/runtime/middleware"
	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/etag"
	mtoserviceitemop "github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations/mto_service_item"
	"github.com/transcom/mymove/pkg/gen/ghcmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/handlers/ghcapi/internal/payloads"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/audit"
	"github.com/transcom/mymove/pkg/services/event"
	mtoshipment "github.com/transcom/mymove/pkg/services/mto_shipment"
	"github.com/transcom/mymove/pkg/services/query"
	"github.com/transcom/mymove/pkg/unit"
)

func payloadForClientError(title string, detail string, instance uuid.UUID) *ghcmessages.ClientError {
	return &ghcmessages.ClientError{
		Title:    handlers.FmtString(title),
		Detail:   handlers.FmtString(detail),
		Instance: handlers.FmtUUID(instance),
	}
}

func payloadForValidationError(title string, detail string, instance uuid.UUID, validationErrors *validate.Errors) *ghcmessages.ValidationError {
	payload := &ghcmessages.ValidationError{
		ClientError: *payloadForClientError(title, detail, instance),
	}

	if validationErrors != nil {
		payload.InvalidFields = handlers.NewValidationErrorsResponse(validationErrors).Errors
	}

	return payload
}

// GetMTOServiceItem returns an MTO Service item stored in the mto_service_items table
// requires a uuid to find the service item
type GetMTOServiceItemHandler struct {
	handlers.HandlerConfig
	mtoServiceItemFetcher services.MTOServiceItemFetcher
}

func (h GetMTOServiceItemHandler) Handle(params mtoserviceitemop.GetMTOServiceItemParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {
			// handling error responses based on error values
			handleError := func(err error) (middleware.Responder, error) {
				appCtx.Logger().Error("GetServiceItem error", zap.Error(err))
				payload := &ghcmessages.Error{Message: handlers.FmtString(err.Error())}
				switch err.(type) {
				case apperror.NotFoundError:
					return mtoserviceitemop.NewGetMTOServiceItemNotFound().WithPayload(payload), err
				case apperror.ForbiddenError:
					return mtoserviceitemop.NewGetMTOServiceItemForbidden().WithPayload(payload), err
				case apperror.QueryError:
					return mtoserviceitemop.NewGetMTOServiceItemInternalServerError(), err
				default:
					return mtoserviceitemop.NewGetMTOServiceItemInternalServerError(), err
				}
			}

			mtoServiceItemID, err := uuid.FromString(params.MtoServiceItemID)
			// return parsing errors
			if err != nil {
				parsingError := fmt.Errorf("UUID parsing failed for mtoServiceItem: %w", err).Error()
				appCtx.Logger().Error(parsingError)
				payload := payloadForValidationError(
					"UUID(s) parsing error",
					parsingError,
					h.GetTraceIDFromRequest(params.HTTPRequest),
					validate.NewErrors())

				return mtoserviceitemop.NewUpdateMTOServiceItemStatusUnprocessableEntity().WithPayload(payload), err
			}
			serviceItem, err := h.mtoServiceItemFetcher.GetServiceItem(appCtx, mtoServiceItemID)
			if err != nil {
				return handleError(err)
			}
			payload := payloads.MTOServiceItemSingleModel(serviceItem)
			return mtoserviceitemop.NewGetMTOServiceItemOK().WithPayload(payload), nil
		})
}

type UpdateServiceItemSitEntryDateHandler struct {
	handlers.HandlerConfig
	sitEntryDateUpdater services.SitEntryDateUpdater
	services.ShipmentSITStatus
	services.MTOShipmentFetcher
	services.ShipmentUpdater
}

func (h UpdateServiceItemSitEntryDateHandler) Handle(params mtoserviceitemop.UpdateServiceItemSitEntryDateParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {

			mtoServiceItemID, err := uuid.FromString(params.MtoServiceItemID)
			// return parsing errors
			if err != nil {
				parsingError := fmt.Errorf("UUID parsing failed for mtoServiceItem: %w", err).Error()
				appCtx.Logger().Error(parsingError)
				payload := payloadForValidationError(
					"UUID(s) parsing error",
					parsingError,
					h.GetTraceIDFromRequest(params.HTTPRequest),
					validate.NewErrors())

				return mtoserviceitemop.NewUpdateServiceItemSitEntryDateUnprocessableEntity().WithPayload(payload), err
			}
			sitEntryDateModel := models.SITEntryDateUpdate{
				ID:           mtoServiceItemID,
				SITEntryDate: (*time.Time)(params.Body.SitEntryDate),
			}
			serviceItem, err := h.sitEntryDateUpdater.UpdateSitEntryDate(appCtx, &sitEntryDateModel)
			if err != nil {
				databaseError := fmt.Errorf("UpdateSitEntryDate failed for service item: %w", err).Error()
				appCtx.Logger().Error(databaseError)
				payload := payloadForValidationError(
					"Database error",
					databaseError,
					h.GetTraceIDFromRequest(params.HTTPRequest),
					validate.NewErrors())

				return mtoserviceitemop.NewUpdateServiceItemSitEntryDateUnprocessableEntity().WithPayload(payload), err
			}

			// on service item sit entry date update, update the shipment SIT auth end date
			mtoshipmentID := *serviceItem.MTOShipmentID
			if mtoshipmentID != uuid.Nil {
				eagerAssociations := []string{"MTOServiceItems",
					"MTOServiceItems.SITDepartureDate",
					"MTOServiceItems.SITEntryDate",
					"MTOServiceItems.ReService",
					"SITDurationUpdates",
				}
				shipment, err := mtoshipment.FindShipment(appCtx, mtoshipmentID, eagerAssociations...)
				if shipment != nil {
					_, shipmentWithSITInfo, err := h.CalculateShipmentSITStatus(appCtx, *shipment)
					if err != nil {
						appCtx.Logger().Error(fmt.Sprintf("Could not calculate the shipment SIT status for shipment ID: %s: %s", shipment.ID, err))
					}

					existingETag := etag.GenerateEtag(shipment.UpdatedAt)

					shipment, err = h.UpdateShipment(appCtx, &shipmentWithSITInfo, existingETag, "ghc")
					if err != nil {
						appCtx.Logger().Error(fmt.Sprintf("Could not update the shipment SIT auth end date for shipment ID: %s: %s", shipment.ID, err))
					}

				}
				if err != nil {
					appCtx.Logger().Error(fmt.Sprintf("Could not find a shipment for the service item with ID: %s: %s", mtoServiceItemID, err))
				}
			}

			payload := payloads.MTOServiceItemSingleModel(serviceItem)

			return mtoserviceitemop.NewUpdateServiceItemSitEntryDateOK().WithPayload(payload), nil
		})
}

// UpdateMTOServiceItemStatusHandler struct that describes updating service item status
type UpdateMTOServiceItemStatusHandler struct {
	handlers.HandlerConfig
	services.MTOServiceItemUpdater
	services.Fetcher
	services.ShipmentSITStatus
	services.MTOShipmentFetcher
	services.ShipmentUpdater
}

// Handle handler that handles the handling for updating service item status
func (h UpdateMTOServiceItemStatusHandler) Handle(params mtoserviceitemop.UpdateMTOServiceItemStatusParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {
			existingMTOServiceItem := models.MTOServiceItem{}

			mtoServiceItemID, err := uuid.FromString(params.MtoServiceItemID)
			// return parsing errors
			if err != nil {
				parsingError := fmt.Errorf("UUID parsing failed for mtoServiceItem: %w", err).Error()
				appCtx.Logger().Error(parsingError)
				payload := payloadForValidationError(
					"UUID(s) parsing error",
					parsingError,
					h.GetTraceIDFromRequest(params.HTTPRequest),
					validate.NewErrors())

				return mtoserviceitemop.NewUpdateMTOServiceItemStatusUnprocessableEntity().WithPayload(payload), err
			}

			// Fetch the existing service item
			filter := []services.QueryFilter{query.NewQueryFilter("id", "=", mtoServiceItemID)}
			err = h.Fetcher.FetchRecord(appCtx, &existingMTOServiceItem, filter)

			if err != nil {
				appCtx.Logger().Error(fmt.Sprintf(
					"Error finding MTOServiceItem for status update with ID: %s",
					mtoServiceItemID),
					zap.Error(err))
				return mtoserviceitemop.NewUpdateMTOServiceItemStatusNotFound(), err
			}

			// Capture update attempt in audit log
			_, err = audit.Capture(appCtx, &existingMTOServiceItem, nil, params.HTTPRequest)
			if err != nil {
				appCtx.Logger().Error("Auditing service error for service item update.", zap.Error(err))
				return mtoserviceitemop.NewUpdateMTOServiceItemStatusInternalServerError(), err
			}

			updatedMTOServiceItem, err := h.MTOServiceItemUpdater.ApproveOrRejectServiceItem(
				appCtx,
				mtoServiceItemID,
				models.MTOServiceItemStatus(params.Body.Status),
				params.Body.RejectionReason, params.IfMatch)

			if err != nil {
				switch err.(type) {
				case apperror.NotFoundError:
					return mtoserviceitemop.NewUpdateMTOServiceItemStatusNotFound().
							WithPayload(&ghcmessages.Error{Message: handlers.FmtString(err.Error())}),
						err
				case apperror.PreconditionFailedError:
					return mtoserviceitemop.NewUpdateMTOServiceItemStatusPreconditionFailed().
							WithPayload(&ghcmessages.Error{Message: handlers.FmtString(err.Error())}),
						err
				case apperror.InvalidInputError:
					payload := payloadForValidationError(
						"Unable to complete request",
						err.Error(),
						h.GetTraceIDFromRequest(params.HTTPRequest),
						validate.NewErrors())
					return mtoserviceitemop.NewUpdateMTOServiceItemStatusUnprocessableEntity().WithPayload(payload), err
				default:
					appCtx.Logger().Error(fmt.Sprintf("Error saving payment request status for ID: %s: %s", mtoServiceItemID, err))
					return mtoserviceitemop.NewUpdateMTOServiceItemStatusInternalServerError(), err
				}
			}

			// on service item update, update the shipment SIT auth end date
			mtoshipmentID := existingMTOServiceItem.MTOShipment.ID
			if mtoshipmentID != uuid.Nil {
				eagerAssociations := []string{"MTOServiceItems",
					"MTOServiceItems.SITDepartureDate",
					"MTOServiceItems.SITEntryDate",
					"MTOServiceItems.ReService",
					"SITDurationUpdates",
				}
				shipment, err := mtoshipment.FindShipment(appCtx, mtoshipmentID, eagerAssociations...)
				if shipment != nil {
					_, shipmentWithSITInfo, err := h.CalculateShipmentSITStatus(appCtx, *shipment)
					if err != nil {
						appCtx.Logger().Error(fmt.Sprintf("Could not calculate the shipment SIT status for shipment ID: %s: %s", shipment.ID, err))
					}

					existingETag := etag.GenerateEtag(shipment.UpdatedAt)

					shipment, err = h.UpdateShipment(appCtx, &shipmentWithSITInfo, existingETag, "ghc")
					if err != nil {
						appCtx.Logger().Error(fmt.Sprintf("Could not update the shipment SIT auth end date for shipment ID: %s: %s", shipment.ID, err))
					}

				}
				if err != nil {
					appCtx.Logger().Error(fmt.Sprintf("Could not find a shipment for the service item with ID: %s: %s", mtoServiceItemID, err))
				}
			}

			// trigger webhook event for Prime
			_, err = event.TriggerEvent(event.Event{
				EventKey:        event.MTOServiceItemUpdateEventKey,
				MtoID:           existingMTOServiceItem.MoveTaskOrder.ID,
				UpdatedObjectID: existingMTOServiceItem.ID,
				EndpointKey:     event.GhcUpdateMTOServiceItemStatusEndpointKey,
				AppContext:      appCtx,
				TraceID:         h.GetTraceIDFromRequest(params.HTTPRequest),
			})

			if err != nil {
				appCtx.Logger().Error("ghcapi.UpdateMTOServiceItemStatusHandler could not generate the event")
			}

			payload := payloads.MTOServiceItemModel(updatedMTOServiceItem, h.FileStorer())
			return mtoserviceitemop.NewUpdateMTOServiceItemStatusOK().WithPayload(payload), nil
		})
}

// ListMTOServiceItemsHandler struct that describes listing service items for the move task order
type ListMTOServiceItemsHandler struct {
	handlers.HandlerConfig
	services.ListFetcher
	services.Fetcher
	counselingPricer     services.CounselingServicesPricer
	moveManagementPricer services.ManagementServicesPricer
}

// Handle handler that lists mto service items for the move task order
func (h ListMTOServiceItemsHandler) Handle(params mtoserviceitemop.ListMTOServiceItemsParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {
			moveTaskOrderID, err := uuid.FromString(params.MoveTaskOrderID.String())
			// return any parsing error
			if err != nil {
				parsingError := fmt.Errorf("UUID Parsing for %s: %w", "MoveTaskOrderID", err).Error()
				appCtx.Logger().Error(parsingError)
				payload := payloadForValidationError(
					"UUID(s) parsing error",
					parsingError,
					h.GetTraceIDFromRequest(params.HTTPRequest),
					validate.NewErrors())

				return mtoserviceitemop.NewListMTOServiceItemsUnprocessableEntity().WithPayload(payload), err
			}

			// check if move task order exists first
			queryFilters := []services.QueryFilter{
				query.NewQueryFilter("id", "=", moveTaskOrderID.String()),
			}

			moveTaskOrder := &models.Move{}
			err = h.Fetcher.FetchRecord(appCtx, moveTaskOrder, queryFilters)
			if err != nil {
				appCtx.Logger().Error(
					"Error fetching move task order: ",
					zap.Error(fmt.Errorf("move Task Order ID: %s", moveTaskOrder.ID)),
					zap.Error(err))

				return mtoserviceitemop.NewListMTOServiceItemsNotFound(), err
			}

			queryFilters = []services.QueryFilter{
				query.NewQueryFilter("move_id", "=", moveTaskOrderID.String()),
			}
			queryAssociations := query.NewQueryAssociationsPreload([]services.QueryAssociation{
				query.NewQueryAssociation("ServiceRequestDocuments.ServiceRequestDocumentUploads.Upload"),
				query.NewQueryAssociation("ReService"),
				query.NewQueryAssociation("Dimensions"),
				query.NewQueryAssociation("SITDestinationOriginalAddress"),
				query.NewQueryAssociation("SITDestinationFinalAddress"),
				query.NewQueryAssociation("SITOriginHHGOriginalAddress"),
				query.NewQueryAssociation("SITOriginHHGActualAddress"),
			})

			var serviceItems models.MTOServiceItems
			err = h.ListFetcher.FetchRecordList(appCtx, &serviceItems, queryFilters, queryAssociations, nil, nil)
			// return any errors
			if err != nil {
				appCtx.Logger().Error("Error fetching mto service items: ", zap.Error(err))

				return mtoserviceitemop.NewListMTOServiceItemsInternalServerError(), err
			}

			// Due to a Pop bug we are unable to use EagerPreload to fetch customer contacts, so we need to load them here.
			for i, serviceItem := range serviceItems {
				if serviceItem.ReService.Code == models.ReServiceCodeDDASIT ||
					serviceItem.ReService.Code == models.ReServiceCodeDDDSIT ||
					serviceItem.ReService.Code == models.ReServiceCodeDDFSIT ||
					serviceItem.ReService.Code == models.ReServiceCodeDDSFSC {
					loadErr := appCtx.DB().Load(&serviceItems[i], "CustomerContacts")
					if loadErr != nil {
						return mtoserviceitemop.NewListMTOServiceItemsInternalServerError(), loadErr
					}
				}
			}

			// Get the move management and counseling service items for this move if applicable
			var indices []int
			for i, mtoServiceItem := range serviceItems {
				if mtoServiceItem.MTOShipmentID == nil {
					indices = append(indices, i)
				}
			}

			if len(indices) > 0 {
				if err != nil {
					return mtoserviceitemop.NewListMTOServiceItemsInternalServerError(), err
				}

				for _, index := range indices {
					var price unit.Cents
					var displayParams services.PricingDisplayParams
					var err error
					if serviceItems[index].ReService.Code == "CS" {
						price, displayParams, err = h.counselingPricer.Price(appCtx, serviceItems[index].LockedPriceCents)
					} else if serviceItems[index].ReService.Code == "MS" {
						price, displayParams, err = h.moveManagementPricer.Price(appCtx, serviceItems[index].LockedPriceCents)
					}

					for _, param := range displayParams {
						appCtx.Logger().Debug("key: " + param.Key.String())
						appCtx.Logger().Debug("value: " + param.Value)
					}
					if err != nil {
						return mtoserviceitemop.NewListMTOServiceItemsInternalServerError(), err
					}

					serviceItems[index].PricingEstimate = &price
				}
			}

			returnPayload := payloads.MTOServiceItemModels(serviceItems, h.FileStorer())
			return mtoserviceitemop.NewListMTOServiceItemsOK().WithPayload(returnPayload), nil
		})
}
