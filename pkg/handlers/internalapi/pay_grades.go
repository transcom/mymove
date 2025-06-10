package internalapi

import (
	"github.com/go-openapi/runtime/middleware"

	"github.com/transcom/mymove/pkg/appcontext"
	ordersop "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/orders"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/handlers/internalapi/internal/payloads"
	"github.com/transcom/mymove/pkg/models"
)

type GetPayGradesHandler struct {
	handlers.HandlerConfig
}

// Handle retrieves pay grades for a given affiliation
func (h GetPayGradesHandler) Handle(params ordersop.GetPayGradesParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {
			payGrades, err := models.GetPayGradesForAffiliation(appCtx.DB(), params.Affiliation)
			if err != nil {
				return handlers.ResponseForError(appCtx.Logger(), err), err
			}

			if len(payGrades) < 1 {
				return ordersop.NewGetPayGradesNotFound(), nil
			}

			payload := payloads.PayGrades(payGrades)

			return ordersop.NewGetPayGradesOK().WithPayload(payload), nil
		})
}
