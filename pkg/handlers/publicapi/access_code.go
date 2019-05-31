package publicapi

import (
	"strings"

	"github.com/go-openapi/runtime/middleware"
	"github.com/gofrs/uuid"

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
	payload := &apimessages.AccessCodePayload{
		ID:        handlers.FmtUUID(accessCode.ID),
		Code:      handlers.FmtStringPtr(&accessCode.Code),
		MoveType:  handlers.FmtString(accessCode.MoveType.String()),
		CreatedAt: handlers.FmtDateTime(accessCode.CreatedAt),
	}

	if accessCode.ServiceMemberID != nil {
		payload.ServiceMemberID = *handlers.FmtUUID(*accessCode.ServiceMemberID)
	}

	if accessCode.ClaimedAt != nil {
		payload.ClaimedAt = *handlers.FmtDateTime(*accessCode.ClaimedAt)
	}

	return payload
}

// Handle accepts the code - validates the access code
func (h ValidateAccessCodeHandler) Handle(params accesscodeop.ValidateAccessCodeParams) middleware.Responder {
	session := auth.SessionFromRequestContext(params.HTTPRequest)

	if session == nil {
		return accesscodeop.NewValidateAccessCodeUnauthorized()
	}

	splitParams := strings.Split(*params.Code, "-")
	moveType, code := splitParams[0], splitParams[1]

	accessCode, valid, _ := h.accessCodeValidator.ValidateAccessCode(code, models.SelectedMoveType(moveType))
	var validateAccessCodePayload *apimessages.ValidateAccessCodePayload

	if !valid {
		h.Logger().Warn("Access code not valid")
		validateAccessCodePayload = &apimessages.ValidateAccessCodePayload{
			Valid: &valid,
		}
	}

	accessCodePayload := payloadForAccessCodeModel(*accessCode)
	validateAccessCodePayload = &apimessages.ValidateAccessCodePayload{
		Valid:      &valid,
		AccessCode: accessCodePayload,
	}

	return accesscodeop.NewValidateAccessCodeOK().WithPayload(validateAccessCodePayload)
}

// ClaimAccessCodeHandler updates an access code to mark it as claimed
type ClaimAccessCodeHandler struct {
	handlers.HandlerContext
	accessCodeClaimer services.AccessCodeClaimer
}

// Handle accepts the code - updates the access code
func (h ClaimAccessCodeHandler) Handle(params accesscodeop.ClaimAccessCodeParams) middleware.Responder {
	session := auth.SessionFromRequestContext(params.HTTPRequest)

	if session == nil || session.ServiceMemberID == uuid.Nil {
		return accesscodeop.NewClaimAccessCodeUnauthorized()
	}

	accessCode, verrs, err := h.accessCodeClaimer.ClaimAccessCode(*params.AccessCodePayload.Code, session.ServiceMemberID)

	if err != nil || verrs.HasAny() {
		return handlers.ResponseForVErrors(h.Logger(), verrs, err)
	}

	accessCodePayload := payloadForAccessCodeModel(*accessCode)

	return accesscodeop.NewClaimAccessCodeOK().WithPayload(accessCodePayload)
}
