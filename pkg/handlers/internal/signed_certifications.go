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
		ID:                fmtUUID(cert.ID),
		CreatedAt:         fmtDateTime(cert.CreatedAt),
		UpdatedAt:         fmtDateTime(cert.UpdatedAt),
		CertificationText: fmtString(cert.CertificationText),
		Signature:         fmtString(cert.Signature),
		Date:              fmtDate(cert.Date),
	}
}

// CreateSignedCertificationHandler creates a new issue via POST /issue
type CreateSignedCertificationHandler utils.HandlerContext

// Handle creates a new SignedCertification from a request payload
func (h CreateSignedCertificationHandler) Handle(params certop.CreateSignedCertificationParams) middleware.Responder {
	session := auth.SessionFromRequestContext(params.HTTPRequest)
	// User should always be populated by middleware
	moveID, _ := uuid.FromString(params.MoveID.String())

	move, err := models.FetchMove(h.Db, session, moveID)
	if err != nil {
		return responseForError(h.Logger, err)
	}

	payload := params.CreateSignedCertificationPayload
	_, verrs, err := move.CreateSignedCertification(h.Db, session.UserID, *payload.CertificationText, *payload.Signature, (time.Time)(*payload.Date))
	if verrs.HasAny() || err != nil {
		return responseForVErrors(h.Logger, verrs, err)
	}

	return certop.NewCreateSignedCertificationCreated()
}

// IndexSignedCertificationsHandler creates a new issue via POST /issue
type IndexSignedCertificationsHandler utils.HandlerContext

// Handle returns a SignedCertification for a given moveID
func (h IndexSignedCertificationsHandler) Handle(params certop.IndexSignedCertificationsParams) middleware.Responder {
	session := auth.SessionFromRequestContext(params.HTTPRequest)
	// #nosec Format of UUID is checked by swagger
	moveID, _ := uuid.FromString(params.MoveID.String())

	move, err := models.FetchMove(h.Db, session, moveID)
	if err != nil {
		return responseForError(h.Logger, err)
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
