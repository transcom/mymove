package services

import (
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
)

// MovingExpenseCreator creates a MovingExpense that is associated with a PPMShipment
//
//go:generate mockery --name MovingExpenseCreator
type MovingExpenseCreator interface {
	CreateMovingExpense(appCtx appcontext.AppContext, ppmShipmentID uuid.UUID) (*models.MovingExpense, error)
}

// MovingExpenseUpdater updates a MovingExpense
//
//go:generate mockery --name MovingExpenseUpdater
type MovingExpenseUpdater interface {
	UpdateMovingExpense(appCtx appcontext.AppContext, movingExpense models.MovingExpense, eTag string) (*models.MovingExpense, error)
}

// MovingExpenseDeleter deletes a MovingExpense
//
//go:generate mockery --name MovingExpenseDeleter
type MovingExpenseDeleter interface {
	DeleteMovingExpense(appCtx appcontext.AppContext, ppmID uuid.UUID, movingExpenseID uuid.UUID) error
}
