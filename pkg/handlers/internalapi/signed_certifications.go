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
	var ptrCertificationType *internalmessages.SignedCertificationType
	if cert.CertificationType != nil {
		certificationType := internalmessages.SignedCertificationType(*cert.CertificationType)
		ptrCertificationType = &certificationType
	}

	return &internalmessages.SignedCertificationPayload{
		CertificationText:        handlers.FmtString(cert.CertificationText),
		CertificationType:        ptrCertificationType,
		CreatedAt:                handlers.FmtDateTime(cert.CreatedAt),
		Date:                     handlers.FmtDateTime(cert.Date),
		ID:                       handlers.FmtUUID(cert.ID),
		MoveID:                   handlers.FmtUUID(cert.MoveID),
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

	var ppmID *uuid.UUID
	if payload.PersonallyProcuredMoveID != nil {
		ppmID, err := uuid.FromString((*payload.PersonallyProcuredMoveID).String())
		if err == nil {
			_, err = models.FetchPersonallyProcuredMove(h.DB(), session, ppmID)
			if err != nil {
				return handlers.ResponseForError(h.Logger(), err)
			}
		}
	}
	var shipmentID *uuid.UUID
	if payload.ShipmentID != nil {
		shipmentID, err := uuid.FromString((*payload.ShipmentID).String())
		if err == nil {
			_, err = models.FetchShipment(h.DB(), session, shipmentID)
			if err != nil {
				return handlers.ResponseForError(h.Logger(), err)
			}
		}
	}

	var ptrCertType *models.SignedCertificationType
	if payload.CertificationType != nil {
		certType := models.SignedCertificationType(*payload.CertificationType)
		ptrCertType = &certType
	}

	move, err := models.FetchMove(h.DB(), session, moveID)
	if err != nil {
		return handlers.ResponseForError(h.Logger(), err)
	}

	newSignedCertification, verrs, err := move.CreateSignedCertification(h.DB(),
		session.UserID,
		*payload.CertificationText,
		*payload.Signature,
		(time.Time)(*payload.Date),
		ppmID,
		shipmentID,
		ptrCertType)
	if verrs.HasAny() || err != nil {
		return handlers.ResponseForVErrors(h.Logger(), verrs, err)
	}
	signedCertificationPayload := payloadForSignedCertificationModel(*newSignedCertification)

	return certop.NewCreateSignedCertificationCreated().WithPayload(signedCertificationPayload)
}

// IndexSignedCertificationsHandler gets all signed certifications associated with a move
type IndexSignedCertificationsHandler struct {
	handlers.HandlerContext
}

// Handle gets a list of SignedCertifications for a move
func (h IndexSignedCertificationsHandler) Handle(params certop.IndexSignedCertificationParams) middleware.Responder {
	session := auth.SessionFromRequestContext(params.HTTPRequest)
	moveID, _ := uuid.FromString(params.MoveID.String())

	_, err := models.FetchMove(h.DB(), session, moveID)
	if err != nil {
		return handlers.ResponseForError(h.Logger(), err)
	}

	signedCertifications, err := models.FetchSignedCertifications(h.DB(), session, moveID)
	var signedCertificationsPayload internalmessages.SignedCertifications
	for _, sc := range signedCertifications {
		signedCertificationsPayload = append(signedCertificationsPayload, payloadForSignedCertificationModel(*sc))
	}
	if err != nil {
		return handlers.ResponseForError(h.Logger(), err)
	}
	return certop.NewIndexSignedCertificationOK().WithPayload(signedCertificationsPayload)
}
