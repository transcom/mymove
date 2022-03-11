package internalapi

import (
	"time"

	"go.uber.org/zap"

	"github.com/go-openapi/runtime/middleware"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
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
	return h.AuditableAppContextFromRequest(params.HTTPRequest,
		func(appCtx appcontext.AppContext) middleware.Responder {
			moveID, _ := uuid.FromString(params.MoveID.String())
			payload := params.CreateSignedCertificationPayload

			var ppmID *uuid.UUID
			if payload.PersonallyProcuredMoveID != nil {
				ppmID, err := uuid.FromString((*payload.PersonallyProcuredMoveID).String())
				if err == nil {
					_, err = models.FetchPersonallyProcuredMove(appCtx.DB(), appCtx.Session(), ppmID)
					if err != nil {
						return handlers.ResponseForError(appCtx.Logger(), err)
					}
				}
			}

			var ptrCertType *models.SignedCertificationType
			if payload.CertificationType != nil {
				certType := models.SignedCertificationType(*payload.CertificationType)
				ptrCertType = &certType
			}

			move, err := models.FetchMove(appCtx.DB(), appCtx.Session(), moveID)
			if err != nil {
				return handlers.ResponseForError(appCtx.Logger(), err)
			}

			newSignedCertification, verrs, err := move.CreateSignedCertification(appCtx.DB(),
				appCtx.Session().UserID,
				*payload.CertificationText,
				*payload.Signature,
				(time.Time)(*payload.Date),
				ppmID,
				ptrCertType)
			if verrs.HasAny() || err != nil {
				return handlers.ResponseForVErrors(appCtx.Logger(), verrs, err)
			}
			signedCertificationPayload := payloadForSignedCertificationModel(*newSignedCertification)
			stringCertType := ""
			if signedCertificationPayload.CertificationType != nil {
				stringCertType = string(*signedCertificationPayload.CertificationType)
			}

			appCtx.Logger().Info("signedCertification created",
				zap.String("id", signedCertificationPayload.ID.String()),
				zap.String("moveId", signedCertificationPayload.MoveID.String()),
				zap.String("createdAt", signedCertificationPayload.CreatedAt.String()),
				zap.String("certification_type", stringCertType),
				zap.String("certification_text", *signedCertificationPayload.CertificationText),
			)

			return certop.NewCreateSignedCertificationCreated().WithPayload(signedCertificationPayload)
		})
}

// IndexSignedCertificationsHandler gets all signed certifications associated with a move
type IndexSignedCertificationsHandler struct {
	handlers.HandlerContext
}

// Handle gets a list of SignedCertifications for a move
func (h IndexSignedCertificationsHandler) Handle(params certop.IndexSignedCertificationParams) middleware.Responder {
	appCtx := h.AppContextFromRequest(params.HTTPRequest)
	moveID, _ := uuid.FromString(params.MoveID.String())

	_, err := models.FetchMove(appCtx.DB(), appCtx.Session(), moveID)
	if err != nil {
		return handlers.ResponseForError(appCtx.Logger(), err)
	}

	signedCertifications, err := models.FetchSignedCertifications(appCtx.DB(), appCtx.Session(), moveID)
	var signedCertificationsPayload internalmessages.SignedCertifications
	for _, sc := range signedCertifications {
		signedCertificationsPayload = append(signedCertificationsPayload, payloadForSignedCertificationModel(*sc))
	}
	if err != nil {
		return handlers.ResponseForError(appCtx.Logger(), err)
	}
	return certop.NewIndexSignedCertificationOK().WithPayload(signedCertificationsPayload)
}
