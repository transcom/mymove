package services

import (
	"github.com/gobuffalo/validate"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/auth"
	movedocop "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/move_docs"
	"github.com/transcom/mymove/pkg/models"
)

// MoveDocumentUpdater is an interface for moveDocument implementation
//go:generate mockery -name MoveDocumentUpdater
type MoveDocumentUpdater interface {
	Update(moveDocument movedocop.UpdateMoveDocumentParams, moveId uuid.UUID, session *auth.Session) (*models.MoveDocument, *validate.Errors, error)
	MoveDocumentStatusUpdater
}

// MoveDocumentStatusUpdater is an interface for moveDocument implementation
//go:generate mockery -name MoveDocumentStatusUpdater
type MoveDocumentStatusUpdater interface {
	UpdateMoveDocumentStatus(params movedocop.UpdateMoveDocumentParams, moveDocument *models.MoveDocument, session *auth.Session) (*models.MoveDocument, *validate.Errors, error)
}

