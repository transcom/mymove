package internalapi

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/gofrs/uuid"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/apperror"

	"github.com/transcom/mymove/pkg/appcontext"
	weightticketops "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/ppm"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/handlers/internalapi/internal/payloads"
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
			returnPayload := payloads.CreateWeightTicket(weightTicket)
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
			returnPayload := payloads.UpdateWeightTicket(*updateWeightTicket)
			return weightticketops.NewUpdateWeightTicketOK().WithPayload(returnPayload), nil
		})
}
