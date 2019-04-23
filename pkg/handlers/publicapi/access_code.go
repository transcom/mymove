package publicapi

import (
	"strings"

	"github.com/go-openapi/runtime/middleware"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/gen/apimessages"
	accesscodeop "github.com/transcom/mymove/pkg/gen/restapi/apioperations/accesscode"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

// ValidateAccessCodeHandler validates an access code to allow access to the MilMove platform as a service member
type ValidateAccessCodeHandler struct {
	handlers.HandlerContext
	accessCodeValidator services.AccessCodeValidator
}

func payloadForAccessCodeModel(accessCode models.AccessCode) *apimessages.AccessCodePayload {
	return &apimessages.AccessCodePayload{
		ID:        handlers.FmtUUID(accessCode.ID),
		Code:      handlers.FmtStringPtr(&accessCode.Code),
		MoveID:    *handlers.FmtUUID(accessCode.MoveID),
		MoveType:  handlers.FmtString(accessCode.MoveType.String()),
		CreatedAt: handlers.FmtDateTime(accessCode.CreatedAt),
		UpdatedAt: handlers.FmtDateTime(accessCode.UpdatedAt),
	}
}

// Handle accepts the code - validates the access code
func (h ValidateAccessCodeHandler) Handle(params accesscodeop.ValidateAccessCodeParams) middleware.Responder {
	session := auth.SessionFromRequestContext(params.HTTPRequest)

	if session == nil {
		return accesscodeop.NewValidateAccessCodeUnauthorized()
	}

	splitParams := strings.Split(*params.Code, "-")
	moveType, code := splitParams[0], splitParams[1]

	accessCode, valid, err := h.accessCodeValidator.ValidateAccessCode(code, models.SelectedMoveType(moveType))

	if !valid && err != nil {
		h.Logger().Warn("Access code not valid", zap.Error(err))
		return handlers.ResponseForError(h.Logger(), err)
	}

	accessCodePayload := payloadForAccessCodeModel(*accessCode)

	validateAccessCodePayload := &apimessages.ValidateAccessCodePayload{
		Valid:      &valid,
		AccessCode: accessCodePayload,
	}

	return accesscodeop.NewValidateAccessCodeOK().WithPayload(validateAccessCodePayload)
}
