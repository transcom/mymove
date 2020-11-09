package moveorder

import (
	"database/sql"
	"fmt"

	"github.com/go-openapi/swag"
	"github.com/gobuffalo/pop/v5"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

type moveOrderFetcher struct {
	db *pop.Connection
}

// FilterOption defines the type for the functional arguments used for private functions in MoveOrderFetcher
type FilterOption func(*pop.Query)

func (f moveOrderFetcher) ListMoveOrders(officeUserID uuid.UUID, params *services.ListMoveOrderParams) ([]models.Order, int, error) {
	// Now that we've joined orders and move_orders, we only want to return orders that
	// have an associated move.
	var moveOrders []models.Order
	var transportationOffice models.TransportationOffice
	// select the GBLOC associated with the transportation office of the session's current office user
	err := f.db.Q().
		Join("office_users", "transportation_offices.id = office_users.transportation_office_id").
		Where("office_users.id = ?", officeUserID).First(&transportationOffice)

	if err != nil {
		return []models.Order{}, 0, err
	}

	gbloc := transportationOffice.Gbloc

	// Alright let's build our query based on the filters we got from the handler. These use the FilterOption type above.
	// Essentially these are private functions that return query objects that we can mash together to form a complete
	// query from modular parts.
	branchQuery := branchFilter(params.Branch)
	moveIDQuery := moveIDFilter(params.MoveID)
	dodIDQuery := dodIDFilter(params.DodID)
	lastNameQuery := lastNameFilter(params.LastName)
	dutyStationQuery := destinationDutyStationFilter(params.DestinationDutyStation)
	moveStatusQuery := moveStatusFilter(params.Status)
	// Adding to an array so we can iterate over them and apply the filters after the query structure is set below
	options := [6]FilterOption{branchQuery, moveIDQuery, dodIDQuery, lastNameQuery, dutyStationQuery, moveStatusQuery}

	query := f.db.Q().Eager(
		"ServiceMember",
		"NewDutyStation.Address",
		"OriginDutyStation",
		"Entitlement",
		"Moves.MTOShipments",
		"Moves.MTOServiceItems",
	).InnerJoin("moves", "orders.id = moves.orders_id").
		InnerJoin("service_members", "orders.service_member_id = service_members.id").
		InnerJoin("mto_shipments", "moves.id = mto_shipments.move_id").
		InnerJoin("duty_stations", "orders.origin_duty_station_id = duty_stations.id").
		InnerJoin("transportation_offices", "duty_stations.transportation_office_id = transportation_offices.id").
		Where("transportation_offices.gbloc = ?", gbloc).Order("status desc")

	for _, option := range options {
		if option != nil {
			option(query)
		}
	}

	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return []models.Order{}, 0, services.NotFoundError{}
		default:
			return []models.Order{}, 0, err
		}
	}
	// Pass zeros into paginate in this case. Which will give us 1 page and 20 per page respectively
	if params.Page == nil {
		params.Page = swag.Int64(0)
	}
	if params.PerPage == nil {
		params.Page = swag.Int64(0)
	}

	err = query.GroupBy("orders.id").Paginate(int(*params.Page), int(*params.PerPage)).All(&moveOrders)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return []models.Order{}, 0, services.NotFoundError{}
		default:
			return []models.Order{}, 0, err
		}
	}
	// Get the count
	count := query.Paginator.TotalEntriesSize

	for i := range moveOrders {
		// Due to a bug in pop (https://github.com/gobuffalo/pop/issues/578), we
		// cannot eager load the address as "OriginDutyStation.Address" because
		// OriginDutyStation is a pointer.
		if moveOrders[i].OriginDutyStation != nil {
			f.db.Load(moveOrders[i].OriginDutyStation, "Address", "TransportationOffice")
		}
	}

	return moveOrders, count, nil
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

// These are a bunch of private functions that are used to cobble our list MoveOrders filters together.
func branchFilter(branch *string) FilterOption {
	return func(query *pop.Query) {
		if branch != nil {
			query = query.Where("service_members.affiliation = ?", *branch)
		}
	}
}

func lastNameFilter(lastName *string) FilterOption {
	return func(query *pop.Query) {
		if lastName != nil {
			nameSearch := fmt.Sprintf("%s%%", *lastName)
			query = query.Where("service_members.last_name ILIKE ?", nameSearch)
		}
	}
}

func dodIDFilter(dodID *string) FilterOption {
	return func(query *pop.Query) {
		if dodID != nil {
			query = query.Where("service_members.edipi = ?", dodID)
		}
	}
}

func moveIDFilter(moveID *string) FilterOption {
	return func(query *pop.Query) {
		if moveID != nil {
			query = query.Where("moves.locator = ?", *moveID)
		}
	}
}
func destinationDutyStationFilter(destinationDutyStation *string) FilterOption {
	return func(query *pop.Query) {
		if destinationDutyStation != nil {
			nameSearch := fmt.Sprintf("%s%%", *destinationDutyStation)
			query = query.InnerJoin("duty_stations as destination_duty_station", "orders.new_duty_station_id = destination_duty_station.id").Where("destination_duty_station.name ILIKE ?", nameSearch)
		}
	}
}

func moveStatusFilter(statuses []string) FilterOption {
	return func(query *pop.Query) {
		// If we have statuses let's use them
		if len(statuses) > 0 {
			query = query.Where("moves.status IN (?)", statuses)
		}
		// If we don't have statuses let's just filter out cancelled and draft moves (they should not be in the queue)
		if len(statuses) <= 0 {
			query = query.Where("moves.status NOT IN (?)", models.MoveStatusDRAFT, models.MoveStatusCANCELED)
		}
	}
}
