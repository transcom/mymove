package moveorder

import (
	"database/sql"

	"github.com/gobuffalo/pop"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

type moveOrderFetcher struct {
	db *pop.Connection
}

func (f moveOrderFetcher) ListMoveOrders() ([]models.MoveOrder, error) {
	var moveOrders []models.MoveOrder
	err := f.db.All(&moveOrders)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return []models.MoveOrder{}, services.NotFoundError{}
		default:
			return []models.MoveOrder{}, err
		}
	}

	// Attempting to load these associations using Eager() returns an error, so this loop
	// loads them one at a time. This is creating a N + 1 query for each association, which is
	// bad. But that's also what the current implementation of Eager does, so this is no worse
	// that what we had.
	for i := range moveOrders {
		f.db.Load(&moveOrders[i], "Customer")
		f.db.Load(&moveOrders[i], "ConfirmationNumber")
		f.db.Load(&moveOrders[i], "DestinationDutyStation")
		f.db.Load(moveOrders[i].DestinationDutyStation, "Address")
		f.db.Load(&moveOrders[i], "OriginDutyStation")
		f.db.Load(moveOrders[i].OriginDutyStation, "Address")
		f.db.Load(&moveOrders[i], "Entitlement")
	}

	return moveOrders, nil
}

// NewMoveOrderFetcher creates a new struct with the service dependencies
func NewMoveOrderFetcher(db *pop.Connection) services.MoveOrderFetcher {
	return &moveOrderFetcher{db}
}

// FetchMoveOrder retrieves a MoveOrder for a given UUID
func (f moveOrderFetcher) FetchMoveOrder(moveOrderID uuid.UUID) (*models.MoveOrder, error) {
	moveOrder := &models.MoveOrder{}
	err := f.db.Find(moveOrder, moveOrderID)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return &models.MoveOrder{}, services.NewNotFoundError(moveOrderID, "")
		default:
			return &models.MoveOrder{}, err
		}
	}

	f.db.Load(moveOrder, "Customer")
	f.db.Load(moveOrder, "ConfirmationNumber")
	f.db.Load(moveOrder, "DestinationDutyStation")
	f.db.Load(moveOrder.DestinationDutyStation, "Address")
	f.db.Load(moveOrder, "OriginDutyStation")
	f.db.Load(moveOrder.OriginDutyStation, "Address")
	f.db.Load(moveOrder, "Entitlement")

	return moveOrder, nil
}
