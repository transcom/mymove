package handlers

import (
	"time"

	"github.com/go-openapi/runtime/middleware"
	"github.com/gobuffalo/uuid"

	"github.com/transcom/mymove/pkg/auth"
	certop "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/certification"
	"github.com/transcom/mymove/pkg/models"
)

// CreateSignedCertificationHandler creates a new issue via POST /issue
type CreateSignedCertificationHandler HandlerContext

// Handle creates a new SignedCertification from a request payload
func (h CreateSignedCertificationHandler) Handle(params certop.CreateSignedCertificationParams) middleware.Responder {
	// User should always be populated by middleware
	user, _ := auth.GetUser(params.HTTPRequest.Context())
	moveID, _ := uuid.FromString(params.MoveID.String())

	move, err := models.FetchMove(h.db, user, moveID)
	if err != nil {
		return responseForError(h.logger, err)
	}

	payload := params.CreateSignedCertificationPayload
	_, verrs, err := move.CreateSignedCertification(h.db, user, *payload.CertificationText, *payload.Signature, (time.Time)(*payload.Date))
	if verrs.HasAny() || err != nil {
		return responseForVErrors(h.logger, verrs, err)
	}

	return certop.NewCreateSignedCertificationCreated()
}
