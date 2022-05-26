package internalapi

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/gofrs/uuid"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"

	mtoshipmentops "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/mto_shipment"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/handlers/internalapi/internal/payloads"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

//
// CREATE
//

// CreateMTOShipmentHandler is the handler to create MTO shipments
type CreateMTOShipmentHandler struct {
	handlers.HandlerConfig
	shipmentCreator services.ShipmentCreator
}

// Handle creates the mto shipment
func (h CreateMTOShipmentHandler) Handle(params mtoshipmentops.CreateMTOShipmentParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {

			if appCtx.Session() == nil || (!appCtx.Session().IsMilApp() && appCtx.Session().ServiceMemberID == uuid.Nil) {
				noSessionErr := apperror.NewSessionError("No service member ID")
				return mtoshipmentops.NewCreateMTOShipmentUnauthorized(), noSessionErr
			}

			payload := params.Body
			if payload == nil {
				noBodyErr := apperror.NewBadDataError("Invalid mto shipment: params Body is nil")
				appCtx.Logger().Error(noBodyErr.Error())
				return mtoshipmentops.NewCreateMTOShipmentBadRequest().WithPayload(payloads.ClientError(handlers.BadRequestErrMessage,
					"The MTO Shipment request body cannot be empty.", h.GetTraceIDFromRequest(params.HTTPRequest))), noBodyErr
			}
			mtoShipment := payloads.MTOShipmentModelFromCreate(payload)
			var err error

			mtoShipment, err = h.shipmentCreator.CreateShipment(appCtx, mtoShipment)

			if err != nil {
				appCtx.Logger().Error("internalapi.CreateMTOShipmentHandler", zap.Error(err))
				switch e := err.(type) {
				case apperror.NotFoundError:
					return mtoshipmentops.
						NewCreateMTOShipmentNotFound().
						WithPayload(
							payloads.ClientError(
								handlers.NotFoundMessage,
								err.Error(),
								h.GetTraceIDFromRequest(params.HTTPRequest),
							),
						), err
				case apperror.InvalidInputError:
					return mtoshipmentops.
						NewCreateMTOShipmentUnprocessableEntity().
						WithPayload(
							payloads.ValidationError(
								handlers.ValidationErrMessage,
								h.GetTraceIDFromRequest(params.HTTPRequest),
								e.ValidationErrors,
							),
						), err
				case apperror.QueryError:
					if e.Unwrap() != nil {
						// If you can unwrap, log the internal error (usually a pq error) for better debugging
						appCtx.Logger().Error("internalapi.CreateMTOServiceItemHandler error", zap.Error(e.Unwrap()))
					}
					return mtoshipmentops.
						NewCreateMTOShipmentInternalServerError().
						WithPayload(
							payloads.InternalServerError(
								nil,
								h.GetTraceIDFromRequest(params.HTTPRequest),
							),
						), err
				default:
					return mtoshipmentops.NewCreateMTOShipmentInternalServerError().
						WithPayload(
							payloads.InternalServerError(
								nil,
								h.GetTraceIDFromRequest(params.HTTPRequest),
							),
						), err
				}
			}

			returnPayload := payloads.MTOShipment(mtoShipment)
			return mtoshipmentops.NewCreateMTOShipmentOK().WithPayload(returnPayload), nil
		})
}

//
// UPDATE
//

// UpdateMTOShipmentHandler is the handler to update MTO shipments
type UpdateMTOShipmentHandler struct {
	handlers.HandlerConfig
	shipmentUpdater services.ShipmentUpdater
}

// Handle updates the mto shipment
func (h UpdateMTOShipmentHandler) Handle(params mtoshipmentops.UpdateMTOShipmentParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {

			if appCtx.Session() == nil {
				noSessionErr := apperror.NewSessionError("No user session")
				return mtoshipmentops.NewUpdateMTOShipmentUnauthorized(), noSessionErr
			}

			if !appCtx.Session().IsMilApp() && appCtx.Session().ServiceMemberID == uuid.Nil {
				noServiceMemberIDErr := apperror.NewSessionError("No service member ID")
				return mtoshipmentops.NewUpdateMTOShipmentForbidden(), noServiceMemberIDErr
			}

			payload := params.Body
			if payload == nil {
				noBodyErr := apperror.NewBadDataError("Invalid mto shipment: params Body is nil")
				appCtx.Logger().Error(noBodyErr.Error())
				return mtoshipmentops.NewUpdateMTOShipmentBadRequest().WithPayload(payloads.ClientError(handlers.BadRequestErrMessage,
					"The MTO Shipment request body cannot be empty.", h.GetTraceIDFromRequest(params.HTTPRequest))), noBodyErr
			}

			mtoShipment := payloads.MTOShipmentModelFromUpdate(payload)
			mtoShipment.ID = uuid.FromStringOrNil(params.MtoShipmentID.String())

			status := mtoShipment.Status
			if status != "" && status != models.MTOShipmentStatusDraft && status != models.MTOShipmentStatusSubmitted {
				invalidShipmentStatusErr := apperror.NewBadDataError("Invalid mto shipment status: shipment in service member app can only have draft or submitted status")
				appCtx.Logger().Error(invalidShipmentStatusErr.Error())

				return mtoshipmentops.NewUpdateMTOShipmentBadRequest().WithPayload(
					payloads.ClientError(handlers.BadRequestErrMessage,
						"When present, the MTO Shipment status must be one of: DRAFT or SUBMITTED.",
						h.GetTraceIDFromRequest(params.HTTPRequest))), invalidShipmentStatusErr
			}

			updatedMTOShipment, err := h.shipmentUpdater.UpdateShipment(appCtx, mtoShipment, params.IfMatch)

			if err != nil {
				appCtx.Logger().Error("internalapi.UpdateMTOShipmentHandler", zap.Error(err))
				switch e := err.(type) {
				case apperror.NotFoundError:
					return mtoshipmentops.
						NewUpdateMTOShipmentNotFound().
						WithPayload(
							payloads.ClientError(
								handlers.NotFoundMessage,
								err.Error(),
								h.GetTraceIDFromRequest(params.HTTPRequest),
							),
						), err
				case apperror.InvalidInputError:
					return mtoshipmentops.
						NewUpdateMTOShipmentUnprocessableEntity().
						WithPayload(payloads.
							ValidationError(
								handlers.ValidationErrMessage,
								h.GetTraceIDFromRequest(
									params.HTTPRequest,
								), e.ValidationErrors,
							),
						), err
				case apperror.PreconditionFailedError:
					return mtoshipmentops.
						NewUpdateMTOShipmentPreconditionFailed().
						WithPayload(
							payloads.ClientError(
								handlers.PreconditionErrMessage,
								err.Error(),
								h.GetTraceIDFromRequest(params.HTTPRequest),
							),
						), err
				case apperror.QueryError:
					if e.Unwrap() != nil {
						// If you can unwrap, log the internal error (usually a pq error) for better debugging
						appCtx.
							Logger().
							Error(
								"internalapi.UpdateMTOServiceItemHandler error",
								zap.Error(e.Unwrap()),
							)
					}
					return mtoshipmentops.
						NewUpdateMTOShipmentInternalServerError().
						WithPayload(
							payloads.InternalServerError(
								nil,
								h.GetTraceIDFromRequest(params.HTTPRequest),
							),
						), err
				default:
					return mtoshipmentops.
						NewUpdateMTOShipmentInternalServerError().
						WithPayload(
							payloads.InternalServerError(
								nil,
								h.GetTraceIDFromRequest(params.HTTPRequest),
							),
						), err
				}
			}

			returnPayload := payloads.MTOShipment(updatedMTOShipment)

			return mtoshipmentops.NewUpdateMTOShipmentOK().WithPayload(returnPayload), nil
		})
}

//
// GET ALL
//

// ListMTOShipmentsHandler returns a list of MTO Shipments
type ListMTOShipmentsHandler struct {
	handlers.HandlerConfig
	services.MTOShipmentFetcher
}

// Handle listing mto shipments for the move task order
func (h ListMTOShipmentsHandler) Handle(params mtoshipmentops.ListMTOShipmentsParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {

			if appCtx.Session() == nil || (!appCtx.Session().IsMilApp() && appCtx.Session().ServiceMemberID == uuid.Nil) {
				noSessionErr := apperror.NewSessionError("No session or service member ID")
				return mtoshipmentops.NewListMTOShipmentsUnauthorized(), noSessionErr
			}

			moveID, err := uuid.FromString(params.MoveTaskOrderID.String())
			// return any parsing error
			if err != nil {
				appCtx.Logger().Error("Invalid request: move ID not valid")
				return mtoshipmentops.NewListMTOShipmentsBadRequest().WithPayload(payloads.ClientError(handlers.BadRequestErrMessage,
					"The MTO Shipment request body cannot be empty.", h.GetTraceIDFromRequest(params.HTTPRequest))), err
			}

			// Search for shipments
			shipments, err := h.MTOShipmentFetcher.ListMTOShipments(appCtx, moveID)
			if err != nil {
				appCtx.Logger().Error("internalapi.ListMTOShipmentsHandler", zap.Error(err))
				return mtoshipmentops.NewListMTOShipmentsInternalServerError(), err
			}

			payload := payloads.MTOShipments((*models.MTOShipments)(&shipments))
			return mtoshipmentops.NewListMTOShipmentsOK().WithPayload(*payload), nil
		})
}

//
// DELETE
//

// DeleteShipmentHandler soft deletes a shipment
type DeleteShipmentHandler struct {
	handlers.HandlerConfig
	services.ShipmentDeleter
}

// Handle soft deletes a shipment
func (h DeleteShipmentHandler) Handle(params mtoshipmentops.DeleteShipmentParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {
			shipmentID := uuid.FromStringOrNil(params.MtoShipmentID.String())

			sm, err := models.GetCustomerFromShipment(appCtx.DB(), shipmentID)
			if err != nil {
				return mtoshipmentops.NewDeleteShipmentNotFound(), err
			}

			if appCtx.Session().ServiceMemberID != sm.ID {
				return mtoshipmentops.NewDeleteShipmentForbidden(), err
			}

			_, err = h.DeleteShipment(appCtx, shipmentID)
			if err != nil {
				appCtx.Logger().Error("internalapi.DeleteShipmentHandler", zap.Error(err))

				switch err.(type) {
				case apperror.NotFoundError:
					return mtoshipmentops.NewDeleteShipmentNotFound(), err
				case apperror.ConflictError:
					return mtoshipmentops.NewDeleteShipmentConflict(), err
				case apperror.ForbiddenError:
					return mtoshipmentops.NewDeleteShipmentForbidden(), err
				case apperror.UnprocessableEntityError:
					return mtoshipmentops.NewDeleteShipmentUnprocessableEntity(), err
				default:
					return mtoshipmentops.NewDeleteShipmentInternalServerError(), err
				}
			}

			return mtoshipmentops.NewDeleteShipmentNoContent(), nil
		})
}
