package services

import (
	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/models"
)

// MoveDocumentUpdater is an interface for moveDocument implementation
//go:generate mockery --name MoveDocumentUpdater --disable-version-string
type MoveDocumentUpdater interface {
	Update(appCtx appcontext.AppContext, moveDocumentPayload *internalmessages.MoveDocumentPayload, moveID uuid.UUID) (*models.MoveDocument, *validate.Errors, error)
	MoveDocumentStatusUpdater
}

// MoveDocumentStatusUpdater is an interface for moveDocument implementation
//go:generate mockery --name MoveDocumentStatusUpdater --disable-version-string
type MoveDocumentStatusUpdater interface {
	UpdateMoveDocumentStatus(appCtx appcontext.AppContext, moveDocumentPayload *internalmessages.MoveDocumentPayload, moveDocument *models.MoveDocument) (*models.MoveDocument, *validate.Errors, error)
}
