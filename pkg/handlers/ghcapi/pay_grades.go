package ghcapi

import (
	"github.com/go-openapi/runtime/middleware"

	"github.com/transcom/mymove/pkg/appcontext"
	ordersop "github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations/orders"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/handlers/ghcapi/internal/payloads"
	"github.com/transcom/mymove/pkg/models"
)

type GetPayGradesHandler struct {
	handlers.HandlerConfig
}

// Handle retrieves orders in the system belonging to the logged in user given order ID
func (h GetPayGradesHandler) Handle(params ordersop.GetPayGradesParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {
			payGrades, err := models.GetPayGradesForAffiliation(appCtx.DB(), params.Affiliation)
			if err != nil {
				return handlers.ResponseForError(appCtx.Logger(), err), err
			}

			appCtx.DB()

			if len(payGrades) < 1 {
				return ordersop.NewGetPayGradesNotFound(), nil
			}

			payload := payloads.PayGrades(payGrades)

			return ordersop.NewGetPayGradesOK().WithPayload(payload), nil
		})
}
