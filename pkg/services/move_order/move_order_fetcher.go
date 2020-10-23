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

func (f moveOrderFetcher) ListMoveOrders(officeUserID uuid.UUID, options ...func(query *pop.Query)) ([]models.Order, error) {
	// Now that we've joined orders and move_orders, we only want to return orders that
	// have an associated move.
	var moveOrders []models.Order
	var transportationOffice models.TransportationOffice
	// select the GBLOC associated with the transportation office of the session's current office user
	err := f.db.Q().
		Join("office_users", "transportation_offices.id = office_users.transportation_office_id").
		Where("office_users.id = ?", officeUserID).First(&transportationOffice)

	if err != nil {
		return []models.Order{}, err
	}

	gbloc := transportationOffice.Gbloc

	query := f.db.Q().Eager(
		"ServiceMember",
		"NewDutyStation.Address",
		"OriginDutyStation",
		"Entitlement",
		"Moves.MTOShipments",
		"Moves.MTOServiceItems",
	).InnerJoin("moves", "orders.id = moves.orders_id").
		InnerJoin("mto_shipments", "moves.id = mto_shipments.move_id").
		InnerJoin("duty_stations", "orders.origin_duty_station_id = duty_stations.id").
		InnerJoin("transportation_offices", "duty_stations.transportation_office_id = transportation_offices.id").
		Where("transportation_offices.gbloc = ?", gbloc).
		// TODO: Let's include the status in filters that are passed into this service once we build that feature for the TXO queue (instead of it being hardcoded like it is below right now).
		Where("moves.status NOT IN ('DRAFT', 'CANCELLED')").
		GroupBy("orders.id")

	for _, option := range options {
		if option != nil {
			option(query)
		}
	}

	err = query.All(&moveOrders)
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
			f.db.Load(moveOrders[i].OriginDutyStation, "Address", "TransportationOffice")
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
	// Now that we've joined orders and move_orders, we only want to return orders that
	// have an associated move_task_order.
	moveOrder := &models.Order{}
	err := f.db.Q().Eager(
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
