package publicapi

import (
	"strings"

	"github.com/go-openapi/runtime/middleware"
	"github.com/gofrs/uuid"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/gen/apimessages"
	accesscodeop "github.com/transcom/mymove/pkg/gen/restapi/apioperations/accesscode"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

func payloadForAccessCodeModel(accessCode models.AccessCode) *apimessages.AccessCode {
	payload := &apimessages.AccessCode{
		ID:        handlers.FmtUUID(accessCode.ID),
		Code:      handlers.FmtStringPtr(&accessCode.Code),
		MoveType:  handlers.FmtString(accessCode.MoveType.String()),
		CreatedAt: handlers.FmtDateTime(accessCode.CreatedAt),
	}

	if accessCode.ServiceMemberID != nil {
		payload.ServiceMemberID = *handlers.FmtUUID(*accessCode.ServiceMemberID)
	}

	if accessCode.ClaimedAt != nil {
		payload.ClaimedAt = handlers.FmtDateTime(*accessCode.ClaimedAt)
	}

	return payload
}

// FetchAccessCodeHandler fetches an access code associated with a service member
type FetchAccessCodeHandler struct {
	handlers.HandlerContext
	accessCodeFetcher services.AccessCodeFetcher
}

// Handle fetches the access code for a service member
func (h FetchAccessCodeHandler) Handle(params accesscodeop.FetchAccessCodeParams) middleware.Responder {
	session, logger := h.SessionAndLoggerFromRequest(params.HTTPRequest)

	if session == nil {
		return accesscodeop.NewFetchAccessCodeUnauthorized()
	}

	// Fetch access code
	accessCode, err := h.accessCodeFetcher.FetchAccessCode(session.ServiceMemberID)
	var fetchAccessCodePayload *apimessages.AccessCode

	if err != nil {
		logger.Error("Error retrieving access_code for service member", zap.Error(err))
		fetchAccessCodePayload = &apimessages.AccessCode{}
		return accesscodeop.NewFetchAccessCodeOK().WithPayload(fetchAccessCodePayload)
	}

	fetchAccessCodePayload = payloadForAccessCodeModel(*accessCode)

	return accesscodeop.NewFetchAccessCodeOK().WithPayload(fetchAccessCodePayload)
}

// ValidateAccessCodeHandler validates an access code to allow access to the MilMove platform as a service member
type ValidateAccessCodeHandler struct {
	handlers.HandlerContext
	accessCodeValidator services.AccessCodeValidator
}

// Handle accepts the code - validates the access code
func (h ValidateAccessCodeHandler) Handle(params accesscodeop.ValidateAccessCodeParams) middleware.Responder {
	session, logger := h.SessionAndLoggerFromRequest(params.HTTPRequest)

	if session == nil {
		return accesscodeop.NewValidateAccessCodeUnauthorized()
	}

	splitParams := strings.Split(*params.Code, "-")
	moveType, code := splitParams[0], splitParams[1]

	accessCode, valid, _ := h.accessCodeValidator.ValidateAccessCode(code, models.SelectedMoveType(moveType))
	var validateAccessCodePayload *apimessages.AccessCode

	if !valid {
		logger.Warn("Access code not valid")
		validateAccessCodePayload = &apimessages.AccessCode{}
		return accesscodeop.NewValidateAccessCodeOK().WithPayload(validateAccessCodePayload)
	}

	validateAccessCodePayload = payloadForAccessCodeModel(*accessCode)

	return accesscodeop.NewValidateAccessCodeOK().WithPayload(validateAccessCodePayload)
}

// ClaimAccessCodeHandler updates an access code to mark it as claimed
type ClaimAccessCodeHandler struct {
	handlers.HandlerContext
	accessCodeClaimer services.AccessCodeClaimer
}

// Handle accepts the code - updates the access code
func (h ClaimAccessCodeHandler) Handle(params accesscodeop.ClaimAccessCodeParams) middleware.Responder {
	session, logger := h.SessionAndLoggerFromRequest(params.HTTPRequest)

	if session == nil || session.ServiceMemberID == uuid.Nil {
		return accesscodeop.NewClaimAccessCodeUnauthorized()
	}

	accessCode, verrs, err := h.accessCodeClaimer.ClaimAccessCode(*params.AccessCode.Code, session.ServiceMemberID)

	if err != nil || verrs.HasAny() {
		return handlers.ResponseForVErrors(logger, verrs, err)
	}

	accessCodePayload := payloadForAccessCodeModel(*accessCode)

	return accesscodeop.NewClaimAccessCodeOK().WithPayload(accessCodePayload)
}
