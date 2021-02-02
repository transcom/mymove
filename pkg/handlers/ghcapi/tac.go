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

	isValid := false
	invalidTACs := [4]string{
		"2LGT",
		"4EVR",
		"5ALV",
		"14DA$",
	}

	for _, code := range invalidTACs {
		if code == params.Tac {
			isValid = false
			break
		}
		isValid = true
	}

	tacValidationPayload := &ghcmessages.TacValid{
		IsValid: isValid,
	}

	return tacop.NewTacValidationOK().WithPayload(tacValidationPayload)
}
