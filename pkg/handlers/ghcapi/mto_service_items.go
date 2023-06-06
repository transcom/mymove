package ghcapi

import (
	"fmt"

	"github.com/go-openapi/runtime/middleware"
	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	mtoserviceitemop "github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations/mto_service_item"
	"github.com/transcom/mymove/pkg/gen/ghcmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/handlers/ghcapi/internal/payloads"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/models/roles"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/audit"
	"github.com/transcom/mymove/pkg/services/event"
	"github.com/transcom/mymove/pkg/services/query"
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

// UpdateMTOServiceItemStatusHandler struct that describes updating service item status
type UpdateMTOServiceItemStatusHandler struct {
	handlers.HandlerConfig
	services.MTOServiceItemUpdater
	services.Fetcher
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

			payload := payloads.MTOServiceItemModel(updatedMTOServiceItem)
			return mtoserviceitemop.NewUpdateMTOServiceItemStatusOK().WithPayload(payload), nil
		})
}

// ListMTOServiceItemsHandler struct that describes listing service items for the move task order
type ListMTOServiceItemsHandler struct {
	handlers.HandlerConfig
	services.ListFetcher
	services.Fetcher
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
					zap.Error(fmt.Errorf("Move Task Order ID: %s", moveTaskOrder.ID)),
					zap.Error(err))

				return mtoserviceitemop.NewListMTOServiceItemsNotFound(), err
			}

			queryFilters = []services.QueryFilter{
				query.NewQueryFilter("move_id", "=", moveTaskOrderID.String()),
			}
			queryAssociations := query.NewQueryAssociationsPreload([]services.QueryAssociation{
				query.NewQueryAssociation("ReService"),
				query.NewQueryAssociation("Dimensions"),
				query.NewQueryAssociation("SITDestinationOriginalAddress"),
				query.NewQueryAssociation("SITDestinationFinalAddress"),
				query.NewQueryAssociation("SITAddressUpdates.OldAddress"),
				query.NewQueryAssociation("SITAddressUpdates.NewAddress"),
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
					serviceItem.ReService.Code == models.ReServiceCodeDDFSIT {
					loadErr := appCtx.DB().Load(&serviceItems[i], "CustomerContacts")
					if loadErr != nil {
						return mtoserviceitemop.NewListMTOServiceItemsInternalServerError(), loadErr
					}
				}
			}

			returnPayload := payloads.MTOServiceItemModels(serviceItems)
			return mtoserviceitemop.NewListMTOServiceItemsOK().WithPayload(returnPayload), nil
		})
}

// CreateSITAddressUpdateHandler creates a SIT Address Update in the approved state
type CreateSITAddressUpdateHandler struct {
	handlers.HandlerConfig
	services.ApprovedSITAddressUpdateRequestCreator
}

// Handle creates the approved SIT Address Update
func (h CreateSITAddressUpdateHandler) Handle(params mtoserviceitemop.CreateSITAddressUpdateParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {

			payload := params.Body
			serviceItemID := params.MtoServiceItemID

			handleError := func(err error) (middleware.Responder, error) {
				appCtx.Logger().Error("ghcapi.CreateSITAddressUpdate error", zap.Error(err))
				switch e := err.(type) {
				case apperror.NotFoundError:
					payload := ghcmessages.Error{
						Message: handlers.FmtString(err.Error()),
					}
					return mtoserviceitemop.NewCreateSITAddressUpdateNotFound().WithPayload(&payload), err
				case apperror.InvalidInputError:
					payload := payloadForValidationError(
						"Validation errors",
						"CreateSITAddressUpdate",
						h.GetTraceIDFromRequest(params.HTTPRequest),
						e.ValidationErrors)
					return mtoserviceitemop.NewCreateSITAddressUpdateUnprocessableEntity().WithPayload(payload), err
				case apperror.QueryError:
					if e.Unwrap() != nil {
						// If you can unwrap, log the internal error (usually a pq error) for better debugging
						appCtx.Logger().Error("ghcapi.CreateSITAddressUpdate query error", zap.Error(e.Unwrap()))
					}
					return mtoserviceitemop.NewCreateSITAddressUpdateInternalServerError(), err
				case apperror.ForbiddenError:
					return mtoserviceitemop.NewCreateSITAddressUpdateForbidden().
						WithPayload(&ghcmessages.Error{Message: handlers.FmtString(err.Error())}), err
				default:
					return mtoserviceitemop.NewCreateSITAddressUpdateInternalServerError(), err
				}
			}

			if !appCtx.Session().IsOfficeUser() || !appCtx.Session().Roles.HasRole(roles.RoleTypeTOO) {
				return handleError(apperror.NewForbiddenError("is not a TOO"))
			}

			sitAddressUpdate := payloads.ApprovedSITAddressUpdateFromCreate(payload, serviceItemID)
			createdSITAddressUpdate, err := h.ApprovedSITAddressUpdateRequestCreator.CreateApprovedSITAddressUpdate(appCtx, sitAddressUpdate)
			if err != nil {
				return handleError(err)
			}

			returnPayload := payloads.MTOServiceItemModel(&createdSITAddressUpdate.MTOServiceItem)
			return mtoserviceitemop.NewCreateSITAddressUpdateOK().WithPayload(returnPayload), nil
		})
}
