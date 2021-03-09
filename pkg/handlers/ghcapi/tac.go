package ghcapi

import (
	"github.com/go-openapi/runtime/middleware"
	"go.uber.org/zap"

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
	session, logger := h.SessionAndLoggerFromRequest(params.HTTPRequest)

	if session == nil {
		return tacop.NewTacValidationUnauthorized()
	}

	if !session.IsOfficeApp() || !session.IsOfficeUser() {
		return tacop.NewTacValidationForbidden()
	}

	db := h.DB()
	isValid, err := db.Where("tac = $1", params.Tac).Exists(&models.TransportationAccountingCode{})

	if err != nil {
		logger.Error("Error looking for transportation accounting code", zap.Error(err))
		return tacop.NewTacValidationInternalServerError()
	}

	tacValidationPayload := &ghcmessages.TacValid{
		IsValid: &isValid,
	}

	return tacop.NewTacValidationOK().WithPayload(tacValidationPayload)
}
