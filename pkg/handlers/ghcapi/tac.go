package ghcapi

import (
	"github.com/go-openapi/runtime/middleware"

	tacop "github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations/tac"
	"github.com/transcom/mymove/pkg/gen/ghcmessages"
	"github.com/transcom/mymove/pkg/handlers"
)

// TacValidationHandler validates a TAC value
type TacValidationHandler struct {
	handlers.HandlerContext
}

// Handle accepts the TAC value and returns a payload showing if it is valid
func (h TacValidationHandler) Handle(params tacop.TacValidationParams) middleware.Responder {
	session, _ := h.SessionAndLoggerFromRequest(params.HTTPRequest)

	if session == nil {
		return tacop.NewTacValidationUnauthorized()
	}

	tacValidationPayload := &ghcmessages.TacValid{
		IsValid: true,
	}

	return tacop.NewTacValidationOK().WithPayload(tacValidationPayload)
}
