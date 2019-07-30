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
	PPMTransitioner
	MoveDocumentTransitioner
}

// PPMTransitioner is an interface for moveDocument implementation
//go:generate mockery -name PPMTransitioner
type PPMTransitioner interface {
	UpdateMoveDocumentStatus(params movedocop.UpdateMoveDocumentParams, moveDocument *models.MoveDocument, session *auth.Session) (*models.MoveDocument, *validate.Errors, error)
}

// MoveDocumentTransitioner is an interface for moveDocument implementation
//go:generate mockery -name MoveDocumentTransitioner
type MoveDocumentTransitioner interface {
	Commit(params movedocop.UpdateMoveDocumentParams, moveDocument *models.MoveDocument, session *auth.Session) (*models.MoveDocument, *validate.Errors, error)
}
