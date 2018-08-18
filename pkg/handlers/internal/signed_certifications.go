package internal

import (
	"time"

	"github.com/go-openapi/runtime/middleware"
	"github.com/gobuffalo/uuid"

	"github.com/transcom/mymove/pkg/auth"
	certop "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/certification"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/handlers/utils"
	"github.com/transcom/mymove/pkg/models"
)

func payloadForSignedCertificationModel(cert models.SignedCertification) *internalmessages.SignedCertificationPayload {
	return &internalmessages.SignedCertificationPayload{
		ID:                utils.FmtUUID(cert.ID),
		CreatedAt:         utils.FmtDateTime(cert.CreatedAt),
		UpdatedAt:         utils.FmtDateTime(cert.UpdatedAt),
		CertificationText: utils.FmtString(cert.CertificationText),
		Signature:         utils.FmtString(cert.Signature),
		Date:              utils.FmtDate(cert.Date),
	}
}

// CreateSignedCertificationHandler creates a new issue via POST /issue
type CreateSignedCertificationHandler HandlerContext

// Handle creates a new SignedCertification from a request payload
func (h CreateSignedCertificationHandler) Handle(params certop.CreateSignedCertificationParams) middleware.Responder {
	session := auth.SessionFromRequestContext(params.HTTPRequest)
	// User should always be populated by middleware
	moveID, _ := uuid.FromString(params.MoveID.String())

	move, err := models.FetchMove(h.db, session, moveID)
	if err != nil {
		return utils.ResponseForError(h.logger, err)
	}

	payload := params.CreateSignedCertificationPayload
	_, verrs, err := move.CreateSignedCertification(h.db, session.UserID, *payload.CertificationText, *payload.Signature, (time.Time)(*payload.Date))
	if verrs.HasAny() || err != nil {
		return utils.ResponseForVErrors(h.logger, verrs, err)
	}

	return certop.NewCreateSignedCertificationCreated()
}

// IndexSignedCertificationsHandler creates a new issue via POST /issue
type IndexSignedCertificationsHandler HandlerContext

// Handle returns a SignedCertification for a given moveID
func (h IndexSignedCertificationsHandler) Handle(params certop.IndexSignedCertificationsParams) middleware.Responder {
	session := auth.SessionFromRequestContext(params.HTTPRequest)
	// #nosec Format of UUID is checked by swagger
	moveID, _ := uuid.FromString(params.MoveID.String())

	move, err := models.FetchMove(h.db, session, moveID)
	if err != nil {
		return utils.ResponseForError(h.logger, err)
	}

	certs := move.SignedCertifications

	limit := len(certs)
	if params.Limit != nil && limit > int(*params.Limit) {
		limit = int(*params.Limit)
	}

	payload := make(internalmessages.IndexSignedCertificationsPayload, limit)
	for i := 0; i < limit; i++ {
		cert := certs[i]
		payload[i] = payloadForSignedCertificationModel(cert)
	}

	return certop.NewIndexSignedCertificationsOK().WithPayload(payload)
}
