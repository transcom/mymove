package internalapi

import (
	"reflect"
	"time"

	"github.com/go-openapi/runtime/middleware"
	"github.com/gofrs/uuid"
	beeline "github.com/honeycombio/beeline-go"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/auth"
	certop "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/certification"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
)

func payloadForSignedCertificationModel(cert models.SignedCertification) *internalmessages.SignedCertificationPayload {
	return &internalmessages.SignedCertificationPayload{
		ID:                handlers.FmtUUID(cert.ID),
		CreatedAt:         handlers.FmtDateTime(cert.CreatedAt),
		UpdatedAt:         handlers.FmtDateTime(cert.UpdatedAt),
		CertificationText: handlers.FmtString(cert.CertificationText),
		Signature:         handlers.FmtString(cert.Signature),
		Date:              handlers.FmtDate(cert.Date),
	}
}

// CreateSignedCertificationHandler creates a new issue via POST /issue
type CreateSignedCertificationHandler struct {
	handlers.HandlerContext
}

// Handle creates a new SignedCertification from a request payload
func (h CreateSignedCertificationHandler) Handle(params certop.CreateSignedCertificationParams) middleware.Responder {
	ctx, span := beeline.StartSpan(params.HTTPRequest.Context(), reflect.TypeOf(h).Name())
	defer span.Send()

	session := auth.SessionFromRequestContext(params.HTTPRequest)
	// User should always be populated by middleware
	moveID, _ := uuid.FromString(params.MoveID.String())

	move, err := models.FetchMove(h.DB(), session, moveID)
	if err != nil {
		return h.RespondAndTraceError(ctx, err, "error fetching move", zap.String("move_id", moveID.String()))
	}

	payload := params.CreateSignedCertificationPayload
	_, verrs, err := move.CreateSignedCertification(h.DB(), session.UserID, *payload.CertificationText, *payload.Signature, (time.Time)(*payload.Date))
	if verrs.HasAny() || err != nil {
		return handlers.ResponseForVErrors(h.Logger(), verrs, err)
	}

	return certop.NewCreateSignedCertificationCreated()
}

// IndexSignedCertificationsHandler creates a new issue via POST /issue
type IndexSignedCertificationsHandler struct {
	handlers.HandlerContext
}
