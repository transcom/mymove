package order

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/go-openapi/swag"
	"github.com/gobuffalo/pop/v5"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/apperror"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

type orderFetcher struct {
}

// QueryOption defines the type for the functional arguments used for private functions in OrderFetcher
type QueryOption func(*pop.Query)

func (f orderFetcher) ListOrders(appCtx appcontext.AppContext, officeUserID uuid.UUID, params *services.ListOrderParams) ([]models.Move, int, error) {
	var moves []models.Move
	var transportationOffice models.TransportationOffice
	// select the GBLOC associated with the transportation office of the session's current office user
	err := appCtx.DB().Q().
		Join("office_users", "transportation_offices.id = office_users.transportation_office_id").
		Where("office_users.id = ?", officeUserID).First(&transportationOffice)

	if err != nil {
		return []models.Move{}, 0, err
	}

	officeUserGbloc := transportationOffice.Gbloc

	// Alright let's build our query based on the filters we got from the handler. These use the FilterOption type above.
	// Essentially these are private functions that return query objects that we can mash together to form a complete
	// query from modular parts.

	// The services counselor queue does not base exclude marine results.
	// Only the TIO and TOO queues should.
	needsCounseling := false
	if len(params.Status) > 0 {
		for _, status := range params.Status {
			if status == string(models.MoveStatusNeedsServiceCounseling) || status == string(models.MoveStatusServiceCounselingCompleted) {
				needsCounseling = true
			}
		}
	}

	branchQuery := branchFilter(params.Branch, needsCounseling)

	// If the user is associated with the USMC GBLOC we want to show them ALL the USMC moves, so let's override here.
	// We also only want to do the gbloc filtering thing if we aren't a USMC user, which we cover with the else.
	// var gblocQuery QueryOption
	var gblocToFilterBy *string
	if officeUserGbloc == "USMC" && !needsCounseling {
		branchQuery = branchFilter(swag.String(string(models.AffiliationMARINES)), needsCounseling)
		gblocToFilterBy = params.OriginGBLOC
	} else {
		gblocToFilterBy = &officeUserGbloc
	}

	// We need to use two different GBLOC filter queries because:
	//  - The Services Counselor queue filters based on the GBLOC of the origin duty station's
	//    transportation office
	//  - The TOO queue uses the GBLOC we get from examining the postal code of the first shipment's
	//    pickup address. However, if that shipment happens to be an NTS-Release, we instead drop
	//    back to the GBLOC of the origin duty station's transportation office since an NTS-Release
	//    does not populate the pickup address field.
	var gblocQuery QueryOption
	if needsCounseling {
		gblocQuery = gblocFilterForSC(gblocToFilterBy)
	} else {
		gblocQuery = gblocFilterForTOO(gblocToFilterBy)
	}
	locatorQuery := locatorFilter(params.Locator)
	dodIDQuery := dodIDFilter(params.DodID)
	lastNameQuery := lastNameFilter(params.LastName)
	dutyStationQuery := destinationDutyStationFilter(params.DestinationDutyLocation)
	originDutyLocationQuery := originDutyLocationFilter(params.OriginDutyLocation)
	moveStatusQuery := moveStatusFilter(params.Status)
	submittedAtQuery := submittedAtFilter(params.SubmittedAt)
	requestedMoveDateQuery := requestedMoveDateFilter(params.RequestedMoveDate)
	sortOrderQuery := sortOrder(params.Sort, params.Order)
	// Adding to an array so we can iterate over them and apply the filters after the query structure is set below
	options := [11]QueryOption{branchQuery, locatorQuery, dodIDQuery, lastNameQuery, dutyStationQuery, originDutyLocationQuery, moveStatusQuery, gblocQuery, submittedAtQuery, requestedMoveDateQuery, sortOrderQuery}

	query := appCtx.DB().Q().EagerPreload(
		"Orders.ServiceMember",
		"Orders.NewDutyLocation.Address",
		"Orders.OriginDutyLocation.Address",
		// See note further below about having to do this in a separate Load call due to a Pop issue.
		// "Orders.OriginDutyLocation.TransportationOffice",
		"Orders.Entitlement",
		"MTOShipments",
		"MTOServiceItems",
		"ShipmentGBLOC",
		"OriginDutyLocationGBLOC",
	).InnerJoin("orders", "orders.id = moves.orders_id").
		InnerJoin("service_members", "orders.service_member_id = service_members.id").
		InnerJoin("mto_shipments", "moves.id = mto_shipments.move_id").
		InnerJoin("duty_locations as origin_dl", "orders.origin_duty_location_id = origin_dl.id").
		// Need to use left join because some duty locations do not have transportation offices
		LeftJoin("transportation_offices as origin_to", "origin_dl.transportation_office_id = origin_to.id").
		// If a customer puts in an invalid ZIP for their pickup address, it won't show up in this view,
		// and we don't want it to get hidden from services counselors.
		LeftJoin("move_to_gbloc", "move_to_gbloc.move_id = moves.id").
		InnerJoin("origin_duty_location_to_gbloc as o_gbloc", "o_gbloc.move_id = moves.id").
		LeftJoin("duty_locations as dest_dl", "dest_dl.id = orders.new_duty_location_id").
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
	groupByColumms = append(groupByColumms, "service_members.id", "orders.id", "origin_dl.id")

	if params.Sort != nil && *params.Sort == "destinationDutyStation" {
		groupByColumms = append(groupByColumms, "dest_dl.name")
	}

	if params.Sort != nil && *params.Sort == "originDutyLocation" {
		groupByColumms = append(groupByColumms, "origin_dl.name")
	}

	if params.Sort != nil && *params.Sort == "originGBLOC" {
		groupByColumms = append(groupByColumms, "origin_to.id")
	}

	err = query.GroupBy("moves.id", groupByColumms...).Paginate(int(*params.Page), int(*params.PerPage)).All(&moves)
	if err != nil {
		return []models.Move{}, 0, err
	}
	// Get the count
	count := query.Paginator.TotalEntriesSize

	for i := range moves {
		// There appears to be a bug in Pop for EagerPreload when you have two or more eager paths with 3+ levels
		// where the first 2 levels match.  For example:
		//   "Orders.OriginDutyLocation.Address" and "Orders.OriginDutyLocation.TransportationOffice"
		// In those cases, only the last relationship is loaded in the results.  So, we can only do one of the paths
		// in the EagerPreload above and request the second one explicitly with a separate Load call.
		//
		// Note that we also had a problem before with Eager as well.  Here's what we found with it:
		//   Due to a bug in pop (https://github.com/gobuffalo/pop/issues/578), we
		//   cannot eager load the address as "OriginDutyLocation.Address" because
		//   OriginDutyLocation is a pointer.
		if moves[i].Orders.OriginDutyLocation != nil {
			loadErr := appCtx.DB().Load(moves[i].Orders.OriginDutyLocation, "TransportationOffice")
			if loadErr != nil {
				return []models.Move{}, 0, err
			}
		}

		err := appCtx.DB().Load(&moves[i].Orders.ServiceMember, "BackupContacts")
		if err != nil {
			return []models.Move{}, 0, err
		}
	}

	return moves, count, nil
}

// NewOrderFetcher creates a new struct with the service dependencies
func NewOrderFetcher() services.OrderFetcher {
	return &orderFetcher{}
}

// FetchOrder retrieves an Order for a given UUID
func (f orderFetcher) FetchOrder(appCtx appcontext.AppContext, orderID uuid.UUID) (*models.Order, error) {
	order := &models.Order{}
	err := appCtx.DB().Q().Eager(
		"ServiceMember.BackupContacts",
		"ServiceMember.ResidentialAddress",
		"NewDutyLocation.Address",
		"OriginDutyLocation",
		"Entitlement",
		"Moves",
	).Find(order, orderID)

	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return &models.Order{}, apperror.NewNotFoundError(orderID, "")
		default:
			return &models.Order{}, apperror.NewQueryError("Order", err, "")
		}
	}

	// Due to a bug in pop (https://github.com/gobuffalo/pop/issues/578), we
	// cannot eager load the address as "OriginDutyLocation.Address" because
	// OriginDutyLocation is a pointer.
	if order.OriginDutyLocation != nil {
		err = appCtx.DB().Load(order.OriginDutyLocation, "Address")
		if err != nil {
			return order, err
		}
	}

	return order, nil
}

// These are a bunch of private functions that are used to cobble our list Orders filters together.
func branchFilter(branch *string, needsCounseling bool) QueryOption {
	return func(query *pop.Query) {
		if branch == nil && !needsCounseling {
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
			query.Where("dest_dl.name ILIKE ?", nameSearch)
		}
	}
}

func originDutyLocationFilter(originDutyLocation *string) QueryOption {
	return func(query *pop.Query) {
		if originDutyLocation != nil {
			nameSearch := fmt.Sprintf("%s%%", *originDutyLocation)
			query.Where("origin_dl.name ILIKE ?", nameSearch)
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

func submittedAtFilter(submittedAt *time.Time) QueryOption {
	return func(query *pop.Query) {
		if submittedAt != nil {
			// Between is inclusive, so the end date is set to 1 milsecond prior to the next day
			submittedAtEnd := submittedAt.AddDate(0, 0, 1).Add(-1 * time.Millisecond)
			query.Where("moves.submitted_at between ? and ?", submittedAt.Format(time.RFC3339), submittedAtEnd.Format(time.RFC3339))
		}
	}
}

func requestedMoveDateFilter(requestedMoveDate *string) QueryOption {
	return func(query *pop.Query) {
		if requestedMoveDate != nil {
			query.Where("mto_shipments.requested_pickup_date = ?", *requestedMoveDate)
		}
	}
}

func gblocFilterForSC(gbloc *string) QueryOption {
	// The SC should only see moves where the origin duty station's GBLOC matches the given GBLOC.
	return func(query *pop.Query) {
		if gbloc != nil {
			query.Where("o_gbloc.gbloc = ?", *gbloc)
		}
	}
}

func gblocFilterForTOO(gbloc *string) QueryOption {
	// The TOO should only see moves where the GBLOC for the first shipment's pickup address matches the given GBLOC
	// unless we're dealing with an NTS-Release shipment. For NTS-Release shipments, we drop back to looking at the
	// origin duty station's GBLOC since an NTS-Release does not populate the pickup address.
	return func(query *pop.Query) {
		if gbloc != nil {
			// Note: extra parens necessary to keep precedence correct when AND'ing all filters together.
			query.Where("((mto_shipments.shipment_type != ? AND move_to_gbloc.gbloc = ?) OR (mto_shipments.shipment_type = ? AND o_gbloc.gbloc = ?))",
				models.MTOShipmentTypeHHGOutOfNTSDom, *gbloc, models.MTOShipmentTypeHHGOutOfNTSDom, *gbloc)
		}
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
		"destinationDutyStation": "dest_dl.name",
		"originDutyLocation":     "origin_dl.name",
		"requestedMoveDate":      "min(mto_shipments.requested_pickup_date)",
		"originGBLOC":            "origin_to.gbloc",
	}

	return func(query *pop.Query) {
		// If we have a sort and order defined let's use it. Otherwise we'll use our default status desc sort order.
		if sort != nil && order != nil {
			if sortTerm, ok := parameters[*sort]; ok {
				if *sort == "lastName" {
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
