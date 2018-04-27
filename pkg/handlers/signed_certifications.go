package handlers

import (
	"time"

	"github.com/go-openapi/runtime/middleware"
	"github.com/gobuffalo/uuid"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/auth"
	certop "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/certification"
	"github.com/transcom/mymove/pkg/models"
)

// CreateSignedCertificationHandler creates a new issue via POST /issue
type CreateSignedCertificationHandler HandlerContext

// Handle creates a new SignedCertification from a request payload
func (h CreateSignedCertificationHandler) Handle(params certop.CreateSignedCertificationParams) middleware.Responder {
	var response middleware.Responder
	// User should always be populated by middleware
	user, _ := auth.GetUser(params.HTTPRequest.Context())

	moveID, err := uuid.FromString(params.MoveID.String())
	if err != nil {
		response = certop.NewCreateSignedCertificationUnauthorized()
		return response
	}

	moveResult, err := models.GetMoveForUser(h.db, user.ID, moveID)
	if err != nil {
		h.logger.Error("DB Error checking on move validity", zap.Error(err))
		return certop.NewCreateSignedCertificationInternalServerError()
	}
	if !moveResult.IsValid() {
		switch errCode := moveResult.ErrorCode(); errCode {
		case models.FetchErrorNotFound: // this won't work yet...
			response = certop.NewCreateSignedCertificationNotFound()
		case models.FetchErrorForbidden:
			response = certop.NewCreateSignedCertificationForbidden()
		default:
			h.logger.Fatal("An error type has occurred that is unaccounted for in this case statement.")
		}
		return response
	}

	move := moveResult.Move()

	newSignedCertification := models.SignedCertification{
		CertificationText: *params.CreateSignedCertificationPayload.CertificationText,
		Signature:         *params.CreateSignedCertificationPayload.Signature,
		Date:              (time.Time)(*params.CreateSignedCertificationPayload.Date),
		SubmittingUserID:  user.ID,
		MoveID:            move.ID,
	}
	if verrs, err := h.db.ValidateAndCreate(&newSignedCertification); verrs.HasAny() || err != nil {
		if verrs.HasAny() {
			h.logger.Error("DB Validation", zap.Error(verrs))
		} else {
			h.logger.Error("DB Insertion", zap.Error(err))
		}
		response = certop.NewCreateSignedCertificationBadRequest()
	} else {
		response = certop.NewCreateSignedCertificationCreated()

	}
	return response
}
