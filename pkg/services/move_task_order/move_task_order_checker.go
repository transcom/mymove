package movetaskorder

import (
	"database/sql"

	"github.com/gobuffalo/pop"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

type moveTaskOrderChecker struct {
	db *pop.Connection
}

// NewMoveTaskOrderChecker creates a new struct with the service dependencies
func NewMoveTaskOrderChecker(db *pop.Connection) services.MoveTaskOrderChecker {
	return &moveTaskOrderChecker{db}
}

//IsAvailableToPrime retrieves a MoveTaskOrder for a given UUID
func (f moveTaskOrderChecker) IsAvailableToPrime(moveTaskOrderID uuid.UUID) error {
	mto := &models.MoveTaskOrder{}
	err := f.db.RawQuery("SELECT * from move_task_orders WHERE id = $1", moveTaskOrderID).First(mto)

	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return services.NewNotFoundError(moveTaskOrderID, "")
		default:
			return err
		}
	}

	if mto.AvailableToPrimeAt == nil {
		return services.NewInvalidInputError(mto.ID, nil, nil, "MTO is not available to prime")
	}

	return nil
}
