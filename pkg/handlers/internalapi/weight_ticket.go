package internalapi

import (
	"database/sql"

	"github.com/go-openapi/runtime/middleware"
	"github.com/gofrs/uuid"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/db/utilities"
	weightticketops "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/ppm"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/handlers/internalapi/internal/payloads"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

// CreateWeightTicketHandler
type CreateWeightTicketHandler struct {
	handlers.HandlerConfig
	weightTicketCreator services.WeightTicketCreator
}

// Handle creates a weight ticket
// Depending on the SO, may need to change the document params to weight ticket params
func (h CreateWeightTicketHandler) Handle(params weightticketops.CreateWeightTicketParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {
			if appCtx.Session() == nil {
				noSessionErr := apperror.NewSessionError("No user session")
				return weightticketops.NewCreateWeightTicketUnauthorized(), noSessionErr
			}

			if !appCtx.Session().IsMilApp() && appCtx.Session().ServiceMemberID == uuid.Nil {
				noServiceMemberIDErr := apperror.NewSessionError("No service member ID")
				return weightticketops.NewCreateWeightTicketForbidden(), noServiceMemberIDErr
			}

			// NO NEED FOR payload_to_model, will need for Update
			ppmShipmentID, err := uuid.FromString(params.PpmShipmentID.String())
			if err != nil {
				appCtx.Logger().Error("missing PPM Shipment ID", zap.Error(err))
				return weightticketops.NewCreateWeightTicketBadRequest(), nil
			}

			weightTicket, err := h.weightTicketCreator.CreateWeightTicket(appCtx, ppmShipmentID)

			if err != nil {
				appCtx.Logger().Error("internalapi.CreateWeightTicketHandler", zap.Error(err))
				// Can get a status error
				// Can get an DB error - does the weight ticket, doc create?
				// Can get an error for whether the PPM exist
				switch err.(type) {
				case apperror.InvalidInputError:
					return weightticketops.NewCreateWeightTicketUnprocessableEntity(), err
				case apperror.ForbiddenError:
					return weightticketops.NewCreateWeightTicketForbidden(), err
				case apperror.NotFoundError:
					return weightticketops.NewCreateWeightTicketNotFound(), err
				default:
					return weightticketops.NewCreateWeightTicketInternalServerError(), err
				}
			}
			returnPayload := payloads.WeightTicket(h.FileStorer(), weightTicket)
			return weightticketops.NewCreateWeightTicketOK().WithPayload(returnPayload), nil
		})
}

// UpdateWeightTicketHandler
type UpdateWeightTicketHandler struct {
	handlers.HandlerConfig
	weightTicketUpdater services.WeightTicketUpdater
}

func (h UpdateWeightTicketHandler) Handle(params weightticketops.UpdateWeightTicketParams) middleware.Responder {
	// track every request with middleware:
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {

			if appCtx.Session() == nil {
				noSessionErr := apperror.NewSessionError("No user session")
				return weightticketops.NewUpdateWeightTicketUnauthorized(), noSessionErr
			}

			if !appCtx.Session().IsMilApp() && appCtx.Session().ServiceMemberID == uuid.Nil {
				noServiceMemberIDErr := apperror.NewSessionError("No service member ID")
				return weightticketops.NewUpdateWeightTicketForbidden(), noServiceMemberIDErr
			}

			payload := params.UpdateWeightTicketPayload
			if payload == nil {
				noBodyErr := apperror.NewBadDataError("Invalid weight ticket: params UpdateWeightTicketPayload is nil")
				appCtx.Logger().Error(noBodyErr.Error())
				return weightticketops.NewUpdateWeightTicketBadRequest().WithPayload(payloads.ClientError(handlers.BadRequestErrMessage,
					"The Weight Ticket request payload cannot be empty.", h.GetTraceIDFromRequest(params.HTTPRequest))), noBodyErr
			}

			weightTicket := payloads.WeightTicketModelFromUpdate(payload)
			weightTicket.ID = uuid.FromStringOrNil(params.WeightTicketID.String())

			updateWeightTicket, err := h.weightTicketUpdater.UpdateWeightTicket(appCtx, *weightTicket, params.IfMatch)

			if err != nil {
				appCtx.Logger().Error("internalapi.UpdateWeightTicketHandler", zap.Error(err))
				switch e := err.(type) {
				case apperror.InvalidInputError:
					return weightticketops.NewUpdateWeightTicketUnprocessableEntity().WithPayload(payloads.ValidationError(handlers.ValidationErrMessage, h.GetTraceIDFromRequest(params.HTTPRequest), e.ValidationErrors)), err
				case apperror.PreconditionFailedError:
					return weightticketops.NewUpdateWeightTicketPreconditionFailed().WithPayload(payloads.ClientError(handlers.PreconditionErrMessage, err.Error(), h.GetTraceIDFromRequest(params.HTTPRequest))), err
				case apperror.NotFoundError:
					return weightticketops.NewUpdateWeightTicketNotFound().WithPayload(payloads.ClientError(handlers.NotFoundMessage, err.Error(), h.GetTraceIDFromRequest(params.HTTPRequest))), err
				case apperror.QueryError:
					if e.Unwrap() != nil {
						// If you can unwrap, log the internal error (usually a pq error) for better debugging
						appCtx.
							Logger().
							Error(
								"internalapi.UpdateWeightTicketHandler error",
								zap.Error(e.Unwrap()),
							)
					}
					return weightticketops.
						NewUpdateWeightTicketInternalServerError().
						WithPayload(
							payloads.InternalServerError(
								nil,
								h.GetTraceIDFromRequest(params.HTTPRequest),
							),
						), err
				default:
					return weightticketops.
						NewUpdateWeightTicketInternalServerError().
						WithPayload(
							payloads.InternalServerError(
								nil,
								h.GetTraceIDFromRequest(params.HTTPRequest),
							),
						), err
				}

			}
			returnPayload := payloads.WeightTicket(h.FileStorer(), updateWeightTicket)
			return weightticketops.NewUpdateWeightTicketOK().WithPayload(returnPayload), nil
		})
}

// DeleteWeightTicketHandler
type DeleteWeightTicketHandler struct {
	handlers.HandlerConfig
	weightTicketDeleter services.WeightTicketDeleter
}

// Handle deletes a weight ticket
func (h DeleteWeightTicketHandler) Handle(params weightticketops.DeleteWeightTicketParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {
			if appCtx.Session() == nil {
				noSessionErr := apperror.NewSessionError("No user session")
				appCtx.Logger().Error("internalapi.DeleteWeightTicketHandler", zap.Error(noSessionErr))
				return weightticketops.NewDeleteWeightTicketUnauthorized(), noSessionErr
			}

			if !appCtx.Session().IsMilApp() || appCtx.Session().ServiceMemberID == uuid.Nil {
				noServiceMemberIDErr := apperror.NewSessionError("No service member ID")
				appCtx.Logger().Error("internalapi.DeleteWeightTicketHandler", zap.Error(noServiceMemberIDErr))
				return weightticketops.NewDeleteWeightTicketForbidden(), noServiceMemberIDErr
			}

			// Make sure the service member is not modifying another service member's PPM
			ppmID := uuid.FromStringOrNil(params.PpmShipmentID.String())
			var ppmShipment models.PPMShipment
			err := appCtx.DB().Scope(utilities.ExcludeDeletedScope()).
				EagerPreload(
					"Shipment.MoveTaskOrder.Orders",
				).
				Find(&ppmShipment, ppmID)
			if err != nil {
				if err == sql.ErrNoRows {
					return weightticketops.NewDeleteWeightTicketNotFound(), err
				}
				return weightticketops.NewDeleteWeightTicketInternalServerError(), err
			}
			if ppmShipment.Shipment.MoveTaskOrder.Orders.ServiceMemberID != appCtx.Session().ServiceMemberID {
				wrongServiceMemberIDErr := apperror.NewSessionError("Attempted delete by wrong service member")
				appCtx.Logger().Error("internalapi.DeleteWeightTicketHandler", zap.Error(wrongServiceMemberIDErr))
				return weightticketops.NewDeleteWeightTicketForbidden(), wrongServiceMemberIDErr
			}

			weightTicketID := uuid.FromStringOrNil(params.WeightTicketID.String())

			err = h.weightTicketDeleter.DeleteWeightTicket(appCtx, weightTicketID)
			if err != nil {
				appCtx.Logger().Error("internalapi.DeleteWeightTicketHandler", zap.Error(err))

				switch err.(type) {
				case apperror.NotFoundError:
					return weightticketops.NewDeleteWeightTicketNotFound(), err
				case apperror.ConflictError:
					return weightticketops.NewDeleteWeightTicketConflict(), err
				case apperror.ForbiddenError:
					return weightticketops.NewDeleteWeightTicketForbidden(), err
				case apperror.UnprocessableEntityError:
					return weightticketops.NewDeleteWeightTicketUnprocessableEntity(), err
				default:
					return weightticketops.NewDeleteWeightTicketInternalServerError(), err
				}
			}

			return weightticketops.NewDeleteWeightTicketNoContent(), nil
		})
}
