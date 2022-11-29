package ghcapi

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
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

		movingExpense := payloads.MovingExpenseModelFromUpdate(payload)

		movingExpense.ID = uuid.FromStringOrNil(params.MovingExpenseID.String())

		updatedMovingExpense, _ := h.movingExpenseUpdater.UpdateMovingExpense(appCtx, *movingExpense, params.IfMatch)

		returnPayload := payloads.MovingExpense(h.FileStorer(), updatedMovingExpense)

		return movingexpenseops.NewUpdateMovingExpenseOK().WithPayload(returnPayload), nil
	})
}
