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

//MTOAvailableToPrime retrieves a MoveTaskOrder for a given UUID and checks if it is available to prime
func (f moveTaskOrderChecker) MTOAvailableToPrime(moveTaskOrderID uuid.UUID) (bool, error) {
	mto := &models.MoveTaskOrder{}
	err := f.db.RawQuery("SELECT * from move_task_orders WHERE id = $1", moveTaskOrderID).First(mto)

	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return false, services.NewNotFoundError(moveTaskOrderID, "")
		default:
			return false, err
		}
	}

	if mto.AvailableToPrimeAt == nil {
		return false, services.NewNotFoundError(mto.ID, "MTO with that ID not available")
	}

	return true, nil
}
