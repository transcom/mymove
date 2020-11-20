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

// QueryOption defines the type for the functional arguments used for private functions in MoveOrderFetcher
type QueryOption func(*pop.Query)

func (f moveOrderFetcher) ListMoveOrders(officeUserID uuid.UUID, params *services.ListMoveOrderParams) ([]models.Move, int, error) {
	// Now that we've joined orders and move_orders, we only want to return orders that
	// have an associated move.
	var moves []models.Move
	var transportationOffice models.TransportationOffice
	// select the GBLOC associated with the transportation office of the session's current office user
	err := f.db.Q().
		Join("office_users", "transportation_offices.id = office_users.transportation_office_id").
		Where("office_users.id = ?", officeUserID).First(&transportationOffice)

	if err != nil {
		return []models.Move{}, 0, err
	}

	gbloc := transportationOffice.Gbloc

	// Alright let's build our query based on the filters we got from the handler. These use the FilterOption type above.
	// Essentially these are private functions that return query objects that we can mash together to form a complete
	// query from modular parts.

	branchQuery := branchFilter(params.Branch)
	// If the user is associated with the USMC GBLOC we want to show them ALL the USMC moves, so let's override here.
	// We also only want to do the gbloc filtering thing if we aren't a USMC user, which we cover with the else.
	var gblocQuery QueryOption
	if gbloc == "USMC" {
		branchQuery = branchFilter(swag.String(string(models.AffiliationMARINES)))
	} else {
		gblocQuery = gblocFilter(gbloc)
	}
	moveIDQuery := moveIDFilter(params.MoveID)
	dodIDQuery := dodIDFilter(params.DodID)
	lastNameQuery := lastNameFilter(params.LastName)
	dutyStationQuery := destinationDutyStationFilter(params.DestinationDutyStation)
	moveStatusQuery := moveStatusFilter(params.Status)
	sortOrderQuery := sortOrder(params.Sort, params.Order)
	// Adding to an array so we can iterate over them and apply the filters after the query structure is set below
	options := [8]QueryOption{branchQuery, moveIDQuery, dodIDQuery, lastNameQuery, dutyStationQuery, moveStatusQuery, gblocQuery, sortOrderQuery}

	query := f.db.Q().Eager(
		"Orders.ServiceMember",
		"Orders.NewDutyStation.Address",
		"Orders.OriginDutyStation",
		"Orders.Entitlement",
		"MTOShipments",
		"MTOServiceItems",
	).InnerJoin("orders", "orders.id = moves.orders_id").
		InnerJoin("service_members", "orders.service_member_id = service_members.id").
		InnerJoin("mto_shipments", "moves.id = mto_shipments.move_id").
		InnerJoin("duty_stations", "orders.origin_duty_station_id = duty_stations.id").
		InnerJoin("transportation_offices", "duty_stations.transportation_office_id = transportation_offices.id").
		Where("show = ?", swag.Bool(true))

	for _, option := range options {
		if option != nil {
			option(query)
		}
	}

	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return []models.Move{}, 0, services.NotFoundError{}
		default:
			return []models.Move{}, 0, err
		}
	}
	// Pass zeros into paginate in this case. Which will give us 1 page and 20 per page respectively
	if params.Page == nil {
		params.Page = swag.Int64(0)
	}
	if params.PerPage == nil {
		params.PerPage = swag.Int64(0)
	}

	err = query.GroupBy("moves.id, service_members.id, orders.id, duty_stations.id").Paginate(int(*params.Page), int(*params.PerPage)).All(&moves)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return []models.Move{}, 0, services.NotFoundError{}
		default:
			return []models.Move{}, 0, err
		}
	}
	// Get the count
	count := query.Paginator.TotalEntriesSize

	for i := range moves {
		// Due to a bug in pop (https://github.com/gobuffalo/pop/issues/578), we
		// cannot eager load the address as "OriginDutyStation.Address" because
		// OriginDutyStation is a pointer.
		if moves[i].Orders.OriginDutyStation != nil {
			f.db.Load(moves[i].Orders.OriginDutyStation, "Address", "TransportationOffice")
		}
	}

	return moves, count, nil
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
func branchFilter(branch *string) QueryOption {
	return func(query *pop.Query) {
		if branch == nil {
			query = query.Where("service_members.affiliation != ?", models.AffiliationMARINES)
		}
		if branch != nil {
			query = query.Where("service_members.affiliation = ?", *branch)
		}
	}
}

func lastNameFilter(lastName *string) QueryOption {
	return func(query *pop.Query) {
		if lastName != nil {
			nameSearch := fmt.Sprintf("%s%%", *lastName)
			query = query.Where("service_members.last_name ILIKE ?", nameSearch)
		}
	}
}

func dodIDFilter(dodID *string) QueryOption {
	return func(query *pop.Query) {
		if dodID != nil {
			query = query.Where("service_members.edipi = ?", dodID)
		}
	}
}

func moveIDFilter(moveID *string) QueryOption {
	return func(query *pop.Query) {
		if moveID != nil {
			query = query.Where("moves.locator = ?", *moveID)
		}
	}
}
func destinationDutyStationFilter(destinationDutyStation *string) QueryOption {
	return func(query *pop.Query) {
		if destinationDutyStation != nil {
			nameSearch := fmt.Sprintf("%s%%", *destinationDutyStation)
			query = query.InnerJoin("duty_stations as destination_duty_station", "orders.new_duty_station_id = destination_duty_station.id").Where("destination_duty_station.name ILIKE ?", nameSearch)
		}
	}
}

func moveStatusFilter(statuses []string) QueryOption {
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

func gblocFilter(gbloc string) QueryOption {
	return func(query *pop.Query) {
		query = query.Where("transportation_offices.gbloc = ?", gbloc)
	}
}

func sortOrder(sort *string, order *string) QueryOption {
	parameters := map[string]string{
		"lastName":               "service_members.last_name",
		"dodID":                  "service_members.edipi",
		"branch":                 "service_members.affiliation",
		"moveID":                 "moves.locator",
		"status":                 "moves.status",
		"destinationDutyStation": "duty_stations.name",
	}

	return func(query *pop.Query) {
		// If we have a sort and order defined let's use it. Otherwise we'll use our default status desc sort order.
		if sort != nil && order != nil {
			sortTerm := parameters[*sort]
			if sortTerm == "service_members.last_name" {
				query = query.Order(fmt.Sprintf("service_members.last_name %s, service_members.first_name %s", *order, *order))
			} else {
				query = query.Order(fmt.Sprintf("%s %s", sortTerm, *order))
			}
		} else {
			query = query.Order("moves.status desc")
		}
	}
}
