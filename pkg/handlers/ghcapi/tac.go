package ghcapi

import (
	"strings"

	"github.com/go-openapi/runtime/middleware"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	tacop "github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations/tac"
	"github.com/transcom/mymove/pkg/gen/ghcmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
)

// TacValidationHandler validates a TAC value
type TacValidationHandler struct {
	handlers.HandlerConfig
}

// Handle accepts the TAC value and returns a payload showing if it is valid
func (h TacValidationHandler) Handle(params tacop.TacValidationParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {
			if appCtx.Session() == nil {
				sessionErr := apperror.NewSessionError(
					"user is not authorized",
				)
				appCtx.Logger().Error(sessionErr.Error())
				return tacop.NewTacValidationUnauthorized(), sessionErr
			}

			if !appCtx.Session().IsOfficeApp() || !appCtx.Session().IsOfficeUser() {
				sessionOfficeErr := apperror.NewSessionError(
					"user is not authenticated with TOO office role",
				)
				appCtx.Logger().Error(sessionOfficeErr.Error())
				return tacop.NewTacValidationForbidden(), sessionOfficeErr
			}

			db := appCtx.DB()
			isValid, err := db.Where("tac = $1", strings.ToUpper(params.Tac)).
				Exists(&models.TransportationAccountingCode{})
			if err != nil {
				appCtx.Logger().
					Error("Error looking for transportation accounting code", zap.Error(err))
				return tacop.NewTacValidationInternalServerError(), err
			}

			tacValidationPayload := &ghcmessages.TacValid{
				IsValid: &isValid,
			}

			return tacop.NewTacValidationOK().WithPayload(tacValidationPayload), nil
		})
}
