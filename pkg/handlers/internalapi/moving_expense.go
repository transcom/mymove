package internalapi

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/gofrs/uuid"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	movingexpenseops "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/ppm"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/handlers/internalapi/internal/payloads"
	"github.com/transcom/mymove/pkg/services"
)

// CreateMovingExpenseHandler

type CreateMovingExpenseHandler struct {
	handlers.HandlerConfig
	movingExpenseCreator services.MovingExpenseCreator
}

// Handle creates a moving expense
func (h CreateMovingExpenseHandler) Handle(params movingexpenseops.CreateMovingExpenseParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest, func(appCtx appcontext.AppContext) (middleware.Responder, error) {
		if appCtx.Session() == nil {
			noSessionErr := apperror.NewSessionError("No user session")
			return movingexpenseops.NewCreateMovingExpenseUnauthorized(), noSessionErr
		}
		if !appCtx.Session().IsMilApp() && appCtx.Session().ServiceMemberID == uuid.Nil {
			noServiceMemberIDErr := apperror.NewSessionError("No service member ID")
			return movingexpenseops.NewCreateMovingExpenseForbidden(), noServiceMemberIDErr
		}

		// No need for payload_to_model for Create
		ppmShipmentID, err := uuid.FromString(params.PpmShipmentID.String())
		if err != nil {
			appCtx.Logger().Error("missing PPM Shipment ID", zap.Error(err))
			return movingexpenseops.NewCreateMovingExpenseBadRequest(), nil
		}

		movingExpense, err := h.movingExpenseCreator.CreateMovingExpense(appCtx, ppmShipmentID)

		if err != nil {
			appCtx.Logger().Error("internalapi.CreateMovingExpenseHandler", zap.Error(err))
			// Can get a status error
			// Can get an DB error - check if moving expense doc creates
			// Can get an error for whether the PPM exist
			switch err.(type) {
			case apperror.ForbiddenError:
				return movingexpenseops.NewCreateMovingExpenseForbidden(), err
			case apperror.NotFoundError:
				return movingexpenseops.NewCreateMovingExpenseNotFound(), err
			default:
				return movingexpenseops.NewCreateMovingExpenseInternalServerError(), err
			}
		}
		// Need to add to payload
		returnPayload := payloads.MovingExpense(h.FileStorer(), movingExpense)
		return movingexpenseops.NewCreateMovingExpenseOK().WithPayload(returnPayload), nil
	})
}

// UpdateMovingExpenseHandler
type UpdateMovingExpenseHandler struct {
	handlers.HandlerConfig
	movingExpenseUpdater services.MovingExpenseUpdater
}

func (h UpdateMovingExpenseHandler) Handle(params movingexpenseops.UpdateMovingExpenseParams) middleware.Responder {
	// track every request with middleware:
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {

			if appCtx.Session() == nil {
				noSessionErr := apperror.NewSessionError("No user session")
				return movingexpenseops.NewUpdateMovingExpenseUnauthorized(), noSessionErr
			}

			if !appCtx.Session().IsMilApp() && appCtx.Session().ServiceMemberID == uuid.Nil {
				noServiceMemberIDErr := apperror.NewSessionError("No service member ID")
				return movingexpenseops.NewUpdateMovingExpenseForbidden(), noServiceMemberIDErr
			}

			payload := params.UpdateMovingExpense
			if payload == nil {
				noBodyErr := apperror.NewBadDataError("Invalid moving expense: params UpdateMovingExpense is nil")
				appCtx.Logger().Error(noBodyErr.Error())
				return movingexpenseops.NewUpdateMovingExpenseBadRequest().WithPayload(payloads.ClientError(handlers.BadRequestErrMessage,
					"The moving expense request payload cannot be empty.", h.GetTraceIDFromRequest(params.HTTPRequest))), noBodyErr
			}

			movingExpense := payloads.MovingExpenseModelFromUpdate(payload)
			movingExpense.ID = uuid.FromStringOrNil(params.PpmShipmentID.String())

			updateMovingExpense, err := h.movingExpenseUpdater.UpdateMovingExpense(appCtx, *movingExpense, params.IfMatch)

			if err != nil {
				appCtx.Logger().Error("internalapi.UpdateMovingExpenseHandler", zap.Error(err))
				switch e := err.(type) {
				case apperror.InvalidInputError:
					return movingexpenseops.NewUpdateWeightTicketUnprocessableEntity().WithPayload(payloads.ValidationError(handlers.ValidationErrMessage, h.GetTraceIDFromRequest(params.HTTPRequest), e.ValidationErrors)), err
				case apperror.PreconditionFailedError:
					return movingexpenseops.NewUpdateMovingExpensePreconditionFailed().WithPayload(payloads.ClientError(handlers.PreconditionErrMessage, err.Error(), h.GetTraceIDFromRequest(params.HTTPRequest))), err
				case apperror.ForbiddenError:
					return movingexpenseops.NewUpdateMovingExpenseForbidden().WithPayload(payloads.ClientError(handlers.PreconditionErrMessage, err.Error(), h.GetTraceIDFromRequest(params.HTTPRequest))), err
				case apperror.NotFoundError:
					return movingexpenseops.NewUpdateMovingExpenseNotFound().WithPayload(payloads.ClientError(handlers.NotFoundMessage, err.Error(), h.GetTraceIDFromRequest(params.HTTPRequest))), err
				case apperror.QueryError:
					if e.Unwrap() != nil {
						// If you can unwrap, log the internal error (usually a pq error) for better debugging
						appCtx.
							Logger().
							Error(
								"internalapi.UpdateMovingExpenseHandler error",
								zap.Error(e.Unwrap()),
							)
					}
					return movingexpenseops.
						NewUpdateMovingExpenseInternalServerError().
						WithPayload(
							payloads.InternalServerError(
								nil,
								h.GetTraceIDFromRequest(params.HTTPRequest),
							),
						), err
				default:
					return movingexpenseops.
						NewUpdateMovingExpenseInternalServerError().
						WithPayload(
							payloads.InternalServerError(
								nil,
								h.GetTraceIDFromRequest(params.HTTPRequest),
							),
						), err
				}

			}
			returnPayload := payloads.MovingExpense(h.FileStorer(), updateMovingExpense)
			return movingexpenseops.NewUpdateMovingExpenseOK().WithPayload(returnPayload), nil
		})
}
