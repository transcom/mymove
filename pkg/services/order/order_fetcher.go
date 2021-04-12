package order

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/go-openapi/swag"
	"github.com/gobuffalo/pop/v5"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

type orderFetcher struct {
	db *pop.Connection
}

// QueryOption defines the type for the functional arguments used for private functions in OrderFetcher
type QueryOption func(*pop.Query)

func (f orderFetcher) ListOrders(officeUserID uuid.UUID, params *services.ListOrderParams) ([]models.Move, int, error) {
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
	locatorQuery := locatorFilter(params.Locator)
	dodIDQuery := dodIDFilter(params.DodID)
	lastNameQuery := lastNameFilter(params.LastName)
	dutyStationQuery := destinationDutyStationFilter(params.DestinationDutyStation)
	moveStatusQuery := moveStatusFilter(params.Status)
	sortOrderQuery := sortOrder(params.Sort, params.Order)
	// Adding to an array so we can iterate over them and apply the filters after the query structure is set below
	options := [8]QueryOption{branchQuery, locatorQuery, dodIDQuery, lastNameQuery, dutyStationQuery, moveStatusQuery, gblocQuery, sortOrderQuery}

	query := f.db.Q().EagerPreload(
		"Orders.ServiceMember",
		"Orders.NewDutyStation.Address",
		"Orders.OriginDutyStation.Address",
		// See note further below about having to do this in a separate Load call due to a Pop issue.
		// "Orders.OriginDutyStation.TransportationOffice",
		"Orders.Entitlement",
		"MTOShipments",
		"MTOServiceItems",
	).InnerJoin("orders", "orders.id = moves.orders_id").
		InnerJoin("service_members", "orders.service_member_id = service_members.id").
		InnerJoin("mto_shipments", "moves.id = mto_shipments.move_id").
		InnerJoin("duty_stations as origin_ds", "orders.origin_duty_station_id = origin_ds.id").
		InnerJoin("transportation_offices as origin_to", "origin_ds.transportation_office_id = origin_to.id").
		LeftJoin("duty_stations as dest_ds", "dest_ds.id = orders.new_duty_station_id").
		Where("show = ?", swag.Bool(true)).
		Where("moves.selected_move_type NOT IN (?)", models.SelectedMoveTypeUB, models.SelectedMoveTypePOV)

	for _, option := range options {
		if option != nil {
			option(query) // mutates
		}
	}

	// Pass zeros into paginate in this case. Which will give us 1 page and 20 per page respectively
	if params.Page == nil {
		params.Page = swag.Int64(0)
	}
	if params.PerPage == nil {
		params.PerPage = swag.Int64(0)
	}

	var groupByColumms []string
	groupByColumms = append(groupByColumms, "service_members.id", "orders.id", "origin_ds.id")
	if params.Sort != nil && *params.Sort == "destinationDutyStation" {
		groupByColumms = append(groupByColumms, "dest_ds.name")
	}

	err = query.GroupBy("moves.id", groupByColumms...).Paginate(int(*params.Page), int(*params.PerPage)).All(&moves)
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
		// There appears to be a bug in Pop for EagerPreload when you have two or more eager paths with 3+ levels
		// where the first 2 levels match.  For example:
		//   "Orders.OriginDutyStation.Address" and "Orders.OriginDutyStation.TransportationOffice"
		// In those cases, only the last relationship is loaded in the results.  So, we can only do one of the paths
		// in the EagerPreload above and request the second one explicitly with a separate Load call.
		//
		// Note that we also had a problem before with Eager as well.  Here's what we found with it:
		//   Due to a bug in pop (https://github.com/gobuffalo/pop/issues/578), we
		//   cannot eager load the address as "OriginDutyStation.Address" because
		//   OriginDutyStation is a pointer.
		if moves[i].Orders.OriginDutyStation != nil {
			loadErr := f.db.Load(moves[i].Orders.OriginDutyStation, "TransportationOffice")
			if loadErr != nil {
				return []models.Move{}, 0, err
			}
		}
	}

	return moves, count, nil
}

// NewOrderFetcher creates a new struct with the service dependencies
func NewOrderFetcher(db *pop.Connection) services.OrderFetcher {
	return &orderFetcher{db}
}

// FetchOrder retrieves an Order for a given UUID
func (f orderFetcher) FetchOrder(orderID uuid.UUID) (*models.Order, error) {
	// Now that we've joined orders and move_orders, we only want to return orders that
	// have an associated move_task_order.
	order := &models.Order{}
	err := f.db.Q().Eager(
		"ServiceMember.BackupContacts",
		"ServiceMember.ResidentialAddress",
		"NewDutyStation.Address",
		"OriginDutyStation",
		"Entitlement",
		"Moves",
	).Find(order, orderID)

	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return &models.Order{}, services.NewNotFoundError(orderID, "")
		default:
			return &models.Order{}, err
		}
	}

	// Due to a bug in pop (https://github.com/gobuffalo/pop/issues/578), we
	// cannot eager load the address as "OriginDutyStation.Address" because
	// OriginDutyStation is a pointer.
	if order.OriginDutyStation != nil {
		err = f.db.Load(order.OriginDutyStation, "Address")
		if err != nil {
			return order, err
		}
	}

	return order, nil
}

// These are a bunch of private functions that are used to cobble our list Orders filters together.
func branchFilter(branch *string) QueryOption {
	return func(query *pop.Query) {
		if branch == nil {
			query.Where("service_members.affiliation != ?", models.AffiliationMARINES)
		}
		if branch != nil {
			query.Where("service_members.affiliation = ?", *branch)
		}
	}
}

func lastNameFilter(lastName *string) QueryOption {
	return func(query *pop.Query) {
		if lastName != nil {
			nameSearch := fmt.Sprintf("%s%%", *lastName)
			query.Where("service_members.last_name ILIKE ?", nameSearch)
		}
	}
}

func dodIDFilter(dodID *string) QueryOption {
	return func(query *pop.Query) {
		if dodID != nil {
			query.Where("service_members.edipi = ?", dodID)
		}
	}
}

func locatorFilter(locator *string) QueryOption {
	return func(query *pop.Query) {
		if locator != nil {
			query.Where("moves.locator = ?", *locator)
		}
	}
}
func destinationDutyStationFilter(destinationDutyStation *string) QueryOption {
	return func(query *pop.Query) {
		if destinationDutyStation != nil {
			nameSearch := fmt.Sprintf("%s%%", *destinationDutyStation)
			query.Where("dest_ds.name ILIKE ?", nameSearch)
		}
	}
}

func moveStatusFilter(statuses []string) QueryOption {
	return func(query *pop.Query) {
		// If we have statuses let's use them
		if len(statuses) > 0 {
			var translatedStatuses []string
			for _, status := range statuses {
				if strings.EqualFold(status, string(models.MoveStatusSUBMITTED)) {
					translatedStatuses = append(translatedStatuses, string(models.MoveStatusSUBMITTED), string(models.MoveStatusServiceCounselingCompleted))
				} else {
					translatedStatuses = append(translatedStatuses, status)
				}
			}
			query.Where("moves.status IN (?)", translatedStatuses)
		}
		// The TOO should never see moves that are in the following statuses: Draft, Canceled, Needs Service Counseling
		if len(statuses) <= 0 {
			query.Where("moves.status NOT IN (?)", models.MoveStatusDRAFT, models.MoveStatusCANCELED, models.MoveStatusNeedsServiceCounseling)
		}
	}
}

func gblocFilter(gbloc string) QueryOption {
	return func(query *pop.Query) {
		query.Where("origin_to.gbloc = ?", gbloc)
	}
}

func sortOrder(sort *string, order *string) QueryOption {
	parameters := map[string]string{
		"lastName":               "service_members.last_name",
		"dodID":                  "service_members.edipi",
		"branch":                 "service_members.affiliation",
		"locator":                "moves.locator",
		"status":                 "moves.status",
		"submittedAt":            "moves.submitted_at",
		"destinationDutyStation": "dest_ds.name",
	}

	return func(query *pop.Query) {
		// If we have a sort and order defined let's use it. Otherwise we'll use our default status desc sort order.
		if sort != nil && order != nil {
			if sortTerm, ok := parameters[*sort]; ok {
				if sortTerm == "lastName" {
					query.Order(fmt.Sprintf("service_members.last_name %s, service_members.first_name %s", *order, *order))
				} else {
					query.Order(fmt.Sprintf("%s %s", sortTerm, *order))
				}
			} else {
				query.Order("moves.status desc")
			}
		} else {
			query.Order("moves.status desc")
		}
	}
}
