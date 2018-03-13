package handlers

import (
	"time"

	"github.com/go-openapi/runtime/middleware"
	"github.com/satori/go.uuid"
	"go.uber.org/zap"

	certop "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/certifications"
	"github.com/transcom/mymove/pkg/models"
)

// CreateSignedCertificationHandler creates a new issue via POST /issue
type CreateSignedCertificationHandler HandlerContext

func userCanModifyMove(move models.Move, user models.User) bool {
	// TODO: Handle case where more than one user is authorized to modify move
	return move.UserID == user.ID
}

// Handle creates a new SignedCertification from a request payload
func (h CreateSignedCertificationHandler) Handle(params certop.CreateSignedCertificationParams) middleware.Responder {
	var response middleware.Responder
	user, err := models.GetUserFromRequest(h.db, params.HTTPRequest)
	if err != nil {
		response = certop.NewCreateSignedCertificationUnauthorized()
		return response
	}

	moveID, err := uuid.FromString(params.MoveID.String())
	if err != nil {
		response = certop.NewCreateSignedCertificationUnauthorized()
		return response
	}

	move, err := models.GetMoveByID(h.db, moveID)
	if err != nil {
		// TODO: Think about returning a 404 not found instead
		response = certop.NewCreateSignedCertificationForbidden()
		return response
	}

	if !userCanModifyMove(move, user) {
		response = certop.NewCreateSignedCertificationForbidden()
		return response
	}

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
