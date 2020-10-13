package services

import (
	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/gen/internalmessages"

	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/models"
)

// MoveDocumentUpdater is an interface for moveDocument implementation
//go:generate mockery -name MoveDocumentUpdater
type MoveDocumentUpdater interface {
	Update(moveDocumentPayload *internalmessages.MoveDocumentPayload, moveID uuid.UUID, session *auth.Session) (*models.MoveDocument, *validate.Errors, error)
	MoveDocumentStatusUpdater
}

// MoveDocumentStatusUpdater is an interface for moveDocument implementation
//go:generate mockery -name MoveDocumentStatusUpdater
type MoveDocumentStatusUpdater interface {
	UpdateMoveDocumentStatus(moveDocumentPayload *internalmessages.MoveDocumentPayload, moveDocument *models.MoveDocument, session *auth.Session) (*models.MoveDocument, *validate.Errors, error)
}
