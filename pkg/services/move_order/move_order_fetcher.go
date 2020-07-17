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

func (f moveOrderFetcher) ListMoveOrders() ([]models.Order, error) {
	var moveOrders []models.Order
	err := f.db.Eager(
		"ServiceMember",
		"NewDutyStation",
		"OriginDutyStation",
		"Entitlement",
	).All(&moveOrders)

	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return []models.Order{}, services.NotFoundError{}
		default:
			return []models.Order{}, err
		}
	}

	for i := range moveOrders {
		// Due to a bug in pop (https://github.com/gobuffalo/pop/issues/578), we
		// cannot eager load the address as "OriginDutyStation.Address" because
		// OriginDutyStation is a pointer.
		if moveOrders[i].OriginDutyStation != nil {
			f.db.Load(moveOrders[i].OriginDutyStation, "Address")
		}
	}

	return moveOrders, nil
}

// NewMoveOrderFetcher creates a new struct with the service dependencies
func NewMoveOrderFetcher(db *pop.Connection) services.MoveOrderFetcher {
	return &moveOrderFetcher{db}
}

// FetchMoveOrder retrieves a MoveOrder for a given UUID
func (f moveOrderFetcher) FetchMoveOrder(moveOrderID uuid.UUID) (*models.Order, error) {
	moveOrder := &models.Order{}
	err := f.db.Eager(
		"ServiceMember",
		"NewDutyStation.Address",
		"OriginDutyStation",
		"Entitlement",
	).Find(moveOrder, moveOrderID)

	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return &models.Order{}, services.NewNotFoundError(moveOrderID, "")
		default:
			return &models.Order{}, err
		}
	}

	// Due to a bug in pop (https://github.com/gobuffalo/pop/issues/578), we
	// cannot eager load the address as "OriginDutyStation.Address" because
	// OriginDutyStation is a pointer.
	if moveOrder.OriginDutyStation != nil {
		f.db.Load(moveOrder.OriginDutyStation, "Address")
	}

	return moveOrder, nil
}
