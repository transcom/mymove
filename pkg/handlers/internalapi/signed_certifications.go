package internalapi

import (
	"time"

	"github.com/go-openapi/runtime/middleware"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/auth"
	certop "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/certification"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
)

func payloadForSignedCertificationModel(cert models.SignedCertification) *internalmessages.SignedCertificationPayload {
	return &internalmessages.SignedCertificationPayload{
		CertificationText:        handlers.FmtString(cert.CertificationText),
		CertificationType:        internalmessages.SignedCertificationType(cert.CertificationType),
		CreatedAt:                handlers.FmtDateTime(cert.CreatedAt),
		Date:                     handlers.FmtDate(cert.Date),
		ID:                       handlers.FmtUUID(cert.ID),
		PersonallyProcuredMoveID: handlers.FmtUUIDPtr(cert.PersonallyProcuredMoveID),
		ShipmentID:               handlers.FmtUUIDPtr(cert.ShipmentID),
		Signature:                handlers.FmtString(cert.Signature),
		UpdatedAt:                handlers.FmtDateTime(cert.UpdatedAt),
	}
}

// CreateSignedCertificationHandler creates a new issue via POST /issue
type CreateSignedCertificationHandler struct {
	handlers.HandlerContext
}

// Handle creates a new SignedCertification from a request payload
func (h CreateSignedCertificationHandler) Handle(params certop.CreateSignedCertificationParams) middleware.Responder {
	session := auth.SessionFromRequestContext(params.HTTPRequest)
	// User should always be populated by middleware
	moveID, _ := uuid.FromString(params.MoveID.String())
	payload := params.CreateSignedCertificationPayload

	//TODO Has to be another way.
	var ppmID *uuid.UUID
	tmpPpmID, err := uuid.FromString(payload.PersonallyProcuredMoveID.String())
	if err == nil {
		ppmID = &tmpPpmID
		_, err = models.FetchPersonallyProcuredMove(h.DB(), session, *ppmID)
		if err != nil {
			return handlers.ResponseForError(h.Logger(), err)
		}
	}
	var shipmentID *uuid.UUID
	tmpShipmentID, err := uuid.FromString(payload.ShipmentID.String())
	if err == nil {
		shipmentID = &tmpShipmentID
		_, err = models.FetchShipment(h.DB(), session, *shipmentID)
		if err != nil {
			return handlers.ResponseForError(h.Logger(), err)
		}
	}
	certType := models.SignedCertificationType(payload.CertificationType)

	move, err := models.FetchMove(h.DB(), session, moveID)
	if err != nil {
		return handlers.ResponseForError(h.Logger(), err)
	}

	_, verrs, err := move.CreateSignedCertification(h.DB(),
		session.UserID,
		*payload.CertificationText,
		*payload.Signature, (time.Time)(*payload.Date),
		ppmID,
		shipmentID,
		certType)
	if verrs.HasAny() || err != nil {
		return handlers.ResponseForVErrors(h.Logger(), verrs, err)
	}

	return certop.NewCreateSignedCertificationCreated()
}

// IndexSignedCertificationsHandler creates a new issue via POST /issue
type IndexSignedCertificationsHandler struct {
	handlers.HandlerContext
}
