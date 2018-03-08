package handlers

import (
	"fmt"
	"time"

	"github.com/go-openapi/runtime/middleware"
	"github.com/markbates/pop"
	"go.uber.org/zap"

	"github.com/gorilla/context"
	"github.com/satori/go.uuid"
	certop "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/certification"
	"github.com/transcom/mymove/pkg/models"
)

// CreateSignedCertificationHandler creates a new issue via POST /issue
type CreateSignedCertificationHandler struct {
	db     *pop.Connection
	logger *zap.Logger
}

// NewCreateSignedCertificationHandler returns a new CreateSignedCertificationHandler
func NewCreateSignedCertificationHandler(db *pop.Connection, logger *zap.Logger) CreateSignedCertificationHandler {
	return CreateSignedCertificationHandler{
		db:     db,
		logger: logger,
	}
}

// Handle creates a new SignedCertification from a request payload
func (h CreateSignedCertificationHandler) Handle(params certop.CreateSignedCertificationParams) middleware.Responder {
	var response middleware.Responder
	userID, ok := context.Get(params.HTTPRequest, "user_id").(string)
	if ok {
		response = certop.NewCreateSignedCertificationUnauthorized()
		return response
	}

	userUUID, err := uuid.FromString(userID)
	if err != nil {
		response = certop.NewCreateSignedCertificationUnauthorized()
		return response
	}

	user, err := models.GetUserByID(h.db, userUUID)
	if err != nil {
		response = certop.NewCreateSignedCertificationUnauthorized()
		return response
	}
	fmt.Println(user)

	var userCantUseThisMove bool
	if userCantUseThisMove {
		response = certop.NewCreateSignedCertificationForbidden()
		return response
	}

	newSignedCertification := models.SignedCertification{
		CertificationText: *params.CreateSignedCertificationPayload.CertificationText,
		Signature:         *params.CreateSignedCertificationPayload.Signature,
		Date:              (time.Time)(*params.CreateSignedCertificationPayload.Date),
		SubmittingUserID:  user.ID,
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
