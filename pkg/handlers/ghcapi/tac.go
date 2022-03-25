package ghcapi

import (
	"fmt"
	"strings"

	"github.com/go-openapi/runtime/middleware"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/appcontext"
	tacop "github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations/tac"
	"github.com/transcom/mymove/pkg/gen/ghcmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
)

// TacValidationHandler validates a TAC value
type TacValidationHandler struct {
	handlers.HandlerContext
}

// Handle accepts the TAC value and returns a payload showing if it is valid
func (h TacValidationHandler) Handle(params tacop.TacValidationParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {
			rollbackErr := fmt.Errorf("TAC validation error")

			if appCtx.Session() == nil {
				return tacop.NewTacValidationUnauthorized(), rollbackErr
			}

			if !appCtx.Session().IsOfficeApp() || !appCtx.Session().IsOfficeUser() {
				return tacop.NewTacValidationForbidden(), rollbackErr
			}

			db := appCtx.DB()
			isValid, err := db.Where("tac = $1", strings.ToUpper(params.Tac)).Exists(&models.TransportationAccountingCode{})

			if err != nil {
				appCtx.Logger().Error("Error looking for transportation accounting code", zap.Error(err))
				return tacop.NewTacValidationInternalServerError(), rollbackErr
			}

			tacValidationPayload := &ghcmessages.TacValid{
				IsValid: &isValid,
			}

			return tacop.NewTacValidationOK().WithPayload(tacValidationPayload), nil
		})
}
