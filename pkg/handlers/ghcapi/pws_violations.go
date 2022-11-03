package ghcapi

import (
	"github.com/go-openapi/runtime/middleware"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/appcontext"
	pwsviolationsop "github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations/pws_violations"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/handlers/ghcapi/internal/payloads"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

// PWSViolationsHandler is a struct that describes getting PWS violations for evaluation reports
type GetPWSViolationsHandler struct {
	handlers.HandlerConfig
	services.PWSViolationsFetcher
}

// Handle handles the handling of getting PWS violations for evaluation reports
func (h GetPWSViolationsHandler) Handle(params pwsviolationsop.GetPWSViolationsParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {

			pwsViolations, err := h.GetPWSViolations(appCtx)
			if err != nil {
				if err == models.ErrFetchNotFound {
					appCtx.Logger().Error("Error fetching PWS violations: ", zap.Error(err))
					return pwsviolationsop.NewGetPWSViolationsNotFound(), err
				}
				appCtx.Logger().Error("Error fetching PWS violations: ", zap.Error(err))
				return pwsviolationsop.NewGetPWSViolationsInternalServerError(), err
			}

			returnPayload := payloads.PWSViolations(*pwsViolations)
			return pwsviolationsop.NewGetPWSViolationsOK().WithPayload(returnPayload), nil
		})
}
