package ghcapi

import (
	"fmt"
	"github.com/go-openapi/runtime/middleware"
	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"
	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	movingexpenseops "github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations/ppm"
	"github.com/transcom/mymove/pkg/gen/ghcmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/services"
	"go.uber.org/zap"
)

// UpdateMovingExpenseHandler
type UpdateMovingExpenseHandler struct {
	handlers.HandlerConfig
	movingExpenseUpdater services.MovingExpenseUpdater
}

func (h UpdateMovingExpenseHandler) Handle(params movingexpenseops.UpdateMovingExpenseParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest, func(appCtx appcontext.AppContext) (middleware.Responder, error) {
		if appCtx.Session() == nil {
			noSessionErr := apperror.NewSessionError("No user session")
			return movingexpenseops.NewUpdateMovingExpenseUnauthorized(), noSessionErr
		}
	})

	func(appCtx appcontext.AppContext) (middleware.Responder, error) {
		payload := params.UpdateMovingExpense
		if payload == nil {
			appCtx.Logger().Error("Invalid Moving Expense: params Body is nil")
			emptyBodyError := apperror.NewBadDataError("The request body cannot be empty.")
			payload := payloadForValidationError(
				"Empty body error",
				emptyBodyError.Error(),
				h.GetTraceIDFromRequest(params.HTTPRequest),
				validate.NewErrors(),
			)

			return movingexpenseops.NewUpdateMovingExpenseUnprocessableEntity().WithPayload(payload), emptyBodyError
		}

		ppmshipmentID := uuid.FromStringOrNil(params.PpmShipmentID.String())
		//oldPPMShipment, err := mtoshipment.FindShipment(appCtx, ppmshipmentID)
		// Can't find original moving expense
		if err != nil {
			appCtx.Logger().Error("ghcapi.UpdateShipmentHandler", zap.Error(err))
			switch err.(type) {
			case apperror.NotFoundError:
				return movingexpenseops.NewUpdateMovingExpenseNotFound(), err
			default:
				msg := fmt.Sprintf("%v | Instance: %v", handlers.FmtString(err.Error()), h.GetTraceIDFromRequest(params.HTTPRequest))

				return movingexpenseops.NewUpdateMovingExpenseInternalServerError().WithPayload(
					&ghcmessages.Error{Message: &msg},
				), err
			}
		}

		movingExpense := payloads.MovingExpenseModelFromUpdate(payload)
		movingExpense.ID = ppmshipmentID

		handleError := func(err error) (middleware.Responder, error) {
			appCtx.Logger().Error("ghcapi.UpdateMovingExpenseHandler", zap.Error(err))

			switch e := err.(type) {
			case apperror.NotFoundError:
				return movingexpenseops.NewUpdateMovingExpenseNotFound(), err
			case apperror.ForbiddenError:
				msg := fmt.Sprintf("%v | Instance: %v", handlers.FmtString(err.Error()), h.GetTraceIDFromRequest(params.HTTPRequest))
				return movingexpenseops.NewUpdateMovingExpenseForbidden().WithPayload(
					&ghcmessages.Error{Message: &msg},
				), err
			case apperror.InvalidInputError:
				return movingexpenseops.NewUpdateMovingExpenseUnprocessableEntity().WithPayload(
					payloadForValidationError(
						handlers.ValidationErrMessage,
						err.Error(),
						h.GetTraceIDFromRequest(params.HTTPRequest),
						e.ValidationErrors,
					),
				), err
			case apperror.PreconditionFailedError:
				msg := fmt.Sprintf("%v | Instance: %v", handlers.FmtString(err.Error()), h.GetTraceIDFromRequest(params.HTTPRequest))
				return movingexpenseops.NewUpdateMovingExpensePreconditionFailed().WithPayload(
					&ghcmessages.Error{Message: &msg},
				), err
			default:
				msg := fmt.Sprintf("%v | Instance: %v", handlers.FmtString(err.Error()), h.GetTraceIDFromRequest(params.HTTPRequest))

				return movingexpenseops.NewUpdateMovingExpenseInternalServerError().WithPayload(
					&ghcmessages.Error{Message: &msg},
				), err
			}
		}

		updatedMovingExpense, err := h.movingExpenseUpdater.UpdateMovingExpense(appCtx, movingExpense, params.IfMatch)
		if err != nil {
			return handleError(err)
		}
		returnPayload := payloads.MovingExpense(updatedMovingExpense)
		return movingexpenseops.NewUpdateMovingExpenseOK().WithPayload(returnPayload), nil
	})
}

