package ghcapi

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	movingexpenseops "github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations/ppm"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/handlers/ghcapi/internal/payloads"
	"github.com/transcom/mymove/pkg/services"
)

// UpdateMovingExpenseHandler
type UpdateMovingExpenseHandler struct {
	handlers.HandlerConfig
	movingExpenseUpdater services.MovingExpenseUpdater
}

func (h UpdateMovingExpenseHandler) Handle(params movingexpenseops.UpdateMovingExpenseParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest, func(appCtx appcontext.AppContext) (middleware.Responder, error) {
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

		movingExpense := payloads.MovingExpenseModelFromUpdate(payload)

		movingExpense.ID = uuid.FromStringOrNil(params.MovingExpenseID.String())

		updatedMovingExpense, _ := h.movingExpenseUpdater.UpdateMovingExpense(appCtx, *movingExpense, params.IfMatch)

		returnPayload := payloads.MovingExpense(h.FileStorer(), updatedMovingExpense)

		return movingexpenseops.NewUpdateMovingExpenseOK().WithPayload(returnPayload), nil
	})
}
