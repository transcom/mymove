package order

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/gobuffalo/pop/v6"
	"github.com/gofrs/uuid"
	"github.com/jinzhu/copier"
	"github.com/lib/pq"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/db/utilities"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/models/roles"
	"github.com/transcom/mymove/pkg/services"
	officeuser "github.com/transcom/mymove/pkg/services/office_user"
)

// Since timestamps in a postgres DB are stored with at the microsecond precision, we want to ensure that we are checking all timestamps up until that point to prevent moves from not showing up
// If we only checked values to the second mark, moves towards the end of the day (post 23:59:59 but before 00:00:00) would be lost and not properly show up in the associated filter
const RFC3339Micro = "2006-01-02T15:04:05.999999Z07:00"

type orderFetcher struct {
	waf services.WeightAllotmentFetcher
}

// QueryOption defines the type for the functional arguments used for private functions in OrderFetcher
type QueryOption func(*pop.Query)

func (f orderFetcher) ListOrders(appCtx appcontext.AppContext, officeUserID uuid.UUID, role roles.RoleType, params *services.ListOrderParams) ([]models.Move, int, error) {
	var moves []models.Move

	var officeUserGbloc string
	if params.ViewAsGBLOC != nil {
		officeUserGbloc = *params.ViewAsGBLOC
	} else {
		var gblocErr error
		gblocFetcher := officeuser.NewOfficeUserGblocFetcher()
		officeUserGbloc, gblocErr = gblocFetcher.FetchGblocForOfficeUser(appCtx, officeUserID)
		if gblocErr != nil {
			return []models.Move{}, 0, gblocErr
		}
	}

	privileges, err := roles.FetchPrivilegesForUser(appCtx.DB(), appCtx.Session().UserID)
	if err != nil {
		appCtx.Logger().Error("Error retreiving user privileges", zap.Error(err))
	}

	// Alright let's build our query based on the filters we got from the handler. These use the FilterOption type above.
	// Essentially these are private functions that return query objects that we can mash together to form a complete
	// query from modular parts.

	// The services counselor queue does not base exclude marine results.
	// Only the TIO and TOO queues should.
	needsCounseling := false
	if len(params.Status) > 0 {
		for _, status := range params.Status {
			if status == string(models.MoveStatusNeedsServiceCounseling) {
				needsCounseling = true
			}
		}
	}

	ppmCloseoutGblocs := officeUserGbloc == "NAVY" || officeUserGbloc == "TVCB" || officeUserGbloc == "USCG"

	// Services Counselors in closeout GBLOCs should only see closeout moves
	if needsCounseling && ppmCloseoutGblocs && params.NeedsPPMCloseout != nil && !*params.NeedsPPMCloseout {
		return []models.Move{}, 0, nil
	}

	branchQuery := branchFilter(params.Branch, needsCounseling, ppmCloseoutGblocs)

	// If the user is associated with the USMC GBLOC we want to show them ALL the USMC moves, so let's override here.
	// We also only want to do the gbloc filtering thing if we aren't a USMC user, which we cover with the else.
	// var gblocQuery QueryOption
	var gblocToFilterBy *string
	if officeUserGbloc == "USMC" && !needsCounseling {
		branchQuery = branchFilter(models.StringPointer(string(models.AffiliationMARINES)), needsCounseling, ppmCloseoutGblocs)
		gblocToFilterBy = params.OriginGBLOC
	} else {
		gblocToFilterBy = &officeUserGbloc
	}

	// We need to use three different GBLOC filter queries because:
	//  - The Services Counselor queue filters based on the GBLOC of the origin duty location's
	//    transportation office
	//  - There is a separate queue for the GBLOCs: NAVY, TVCB and USCG. These GBLOCs are used by SC doing PPM Closeout
	//  - The TOO queue uses the GBLOC we get from examining the postal code of the first shipment's
	//    pickup address. However, if that shipment happens to be an NTS-Release, we instead drop
	//    back to the GBLOC of the origin duty location's transportation office since an NTS-Release
	//    does not populate the pickup address field.
	var gblocQuery QueryOption
	if ppmCloseoutGblocs {
		gblocQuery = gblocFilterForPPMCloseoutForNavyMarineAndCG(gblocToFilterBy)
	} else if needsCounseling {
		gblocQuery = gblocFilterForSC(gblocToFilterBy)
	} else if params.NeedsPPMCloseout != nil && *params.NeedsPPMCloseout {
		gblocQuery = gblocFilterForSCinArmyAirForce(gblocToFilterBy)
	} else {
		gblocQuery = gblocFilterForTOO(gblocToFilterBy)
	}
	locatorQuery := locatorFilter(params.Locator)
	dodIDQuery := dodIDFilter(params.Edipi)
	emplidQuery := emplidFilter(params.Emplid)
	customerNameQuery := customerNameFilter(params.CustomerName)
	originDutyLocationQuery := originDutyLocationFilter(params.OriginDutyLocation)
	destinationDutyLocationQuery := destinationDutyLocationFilter(params.DestinationDutyLocation)
	moveStatusQuery := moveStatusFilter(params.Status)
	submittedAtQuery := submittedAtFilter(params.SubmittedAt)
	appearedInTOOAtQuery := appearedInTOOAtFilter(params.AppearedInTOOAt)
	requestedMoveDateQuery := requestedMoveDateFilter(params.RequestedMoveDate)
	closeoutInitiatedQuery := closeoutInitiatedFilter(params.CloseoutInitiated)
	closeoutLocationQuery := closeoutLocationFilter(params.CloseoutLocation, ppmCloseoutGblocs)
	ppmTypeQuery := ppmTypeFilter(params.PPMType)
	ppmStatusQuery := ppmStatusFilter(params.PPMStatus)
	assignedToQuery := assignedUserFilter(params.AssignedTo)
	sortOrderQuery := sortOrder(params.Sort, params.Order, ppmCloseoutGblocs)
	secondarySortOrderQuery := secondarySortOrder(params.Sort)
	counselingQuery := counselingOfficeFilter(params.CounselingOffice)
	tooFilterOutDestinationRequestsQuery := tooQueueOriginRequestsFilter(role)
	// Adding to an array so we can iterate over them and apply the filters after the query structure is set below
	options := [21]QueryOption{branchQuery, locatorQuery, dodIDQuery, emplidQuery, customerNameQuery, originDutyLocationQuery, destinationDutyLocationQuery, moveStatusQuery, gblocQuery, submittedAtQuery, appearedInTOOAtQuery, requestedMoveDateQuery, ppmTypeQuery, closeoutInitiatedQuery, closeoutLocationQuery, ppmStatusQuery, sortOrderQuery, secondarySortOrderQuery, assignedToQuery, counselingQuery, tooFilterOutDestinationRequestsQuery}

	var query *pop.Query
	if ppmCloseoutGblocs {
		query = appCtx.DB().Q().Scope(utilities.ExcludeDeletedScope(models.MTOShipment{})).EagerPreload(
			"Orders.ServiceMember",
			"Orders.NewDutyLocation.Address",
			"Orders.OriginDutyLocation.Address",
			"Orders.Entitlement",
			"Orders.OrdersType",
			"MTOShipments.PPMShipment",
			"LockedByOfficeUser",
			"SCCounselingAssignedUser",
			"SCCloseoutAssignedUser",
			"CounselingOffice",
		).InnerJoin("orders", "orders.id = moves.orders_id").
			InnerJoin("service_members", "orders.service_member_id = service_members.id").
			InnerJoin("mto_shipments", "moves.id = mto_shipments.move_id").
			InnerJoin("ppm_shipments", "ppm_shipments.shipment_id = mto_shipments.id").
			InnerJoin("duty_locations as origin_dl", "orders.origin_duty_location_id = origin_dl.id").
			LeftJoin("duty_locations as dest_dl", "dest_dl.id = orders.new_duty_location_id").
			LeftJoin("office_users", "office_users.id = moves.locked_by").
			LeftJoin("office_users as assigned_user", "moves.sc_closeout_assigned_id  = assigned_user.id").
			LeftJoin("transportation_offices", "moves.counseling_transportation_office_id = transportation_offices.id").
			Where("show = ?", models.BoolPointer(true))

		if !privileges.HasPrivilege(roles.PrivilegeTypeSafety) {
			query.Where("orders.orders_type != (?)", "SAFETY")
		}
	} else {
		query = appCtx.DB().Q().Scope(utilities.ExcludeDeletedScope(models.MTOShipment{})).EagerPreload(
			"Orders.ServiceMember",
			"Orders.NewDutyLocation.Address",
			"Orders.OriginDutyLocation.Address",
			// See note further below about having to do this in a separate Load call due to a Pop issue.
			// "Orders.OriginDutyLocation.TransportationOffice",
			"Orders.Entitlement",
			"Orders.OrdersType",
			"MTOShipments",
			"MTOShipments.SITDurationUpdates",
			"MTOShipments.DeliveryAddressUpdate",
			"MTOServiceItems",
			"MTOServiceItems.ReService",
			"ShipmentGBLOC",
			"MTOShipments.PPMShipment",
			"CloseoutOffice",
			"LockedByOfficeUser",
			"CounselingOffice",
			"SCCounselingAssignedUser",
			"SCCloseoutAssignedUser",
			"TOOAssignedUser",
		).InnerJoin("orders", "orders.id = moves.orders_id").
			InnerJoin("service_members", "orders.service_member_id = service_members.id").
			InnerJoin("mto_shipments", "moves.id = mto_shipments.move_id").
			InnerJoin("duty_locations as origin_dl", "orders.origin_duty_location_id = origin_dl.id").
			// Need to use left join because some duty locations do not have transportation offices
			LeftJoin("transportation_offices as origin_to", "origin_dl.transportation_office_id = origin_to.id").
			// If a customer puts in an invalid ZIP for their pickup address, it won't show up in this view,
			// and we don't want it to get hidden from services counselors.
			LeftJoin("move_to_gbloc", "move_to_gbloc.move_id = moves.id").
			LeftJoin("duty_locations as dest_dl", "dest_dl.id = orders.new_duty_location_id").
			LeftJoin("office_users", "office_users.id = moves.locked_by").
			LeftJoin("transportation_offices", "moves.counseling_transportation_office_id = transportation_offices.id").
			Where("show = ?", models.BoolPointer(true))

		if !privileges.HasPrivilege(roles.PrivilegeTypeSafety) {
			query.Where("orders.orders_type != (?)", "SAFETY")
		}
		if role == roles.RoleTypeTOO {
			query.LeftJoin("office_users as assigned_user", "moves.too_assigned_id  = assigned_user.id")
		}

		if params.NeedsPPMCloseout != nil {
			if *params.NeedsPPMCloseout {
				query.InnerJoin("ppm_shipments", "ppm_shipments.shipment_id = mto_shipments.id").
					LeftJoin("transportation_offices as closeout_to", "closeout_to.id = moves.closeout_office_id").
					LeftJoin("office_users as assigned_user", "moves.sc_closeout_assigned_id  = assigned_user.id").
					Where("ppm_shipments.status IN (?)", models.PPMShipmentStatusNeedsCloseout).
					Where("service_members.affiliation NOT IN (?)", models.AffiliationNAVY, models.AffiliationMARINES, models.AffiliationCOASTGUARD)
			} else {
				query.LeftJoin("ppm_shipments", "ppm_shipments.shipment_id = mto_shipments.id").
					LeftJoin("office_users as assigned_user", "moves.sc_counseling_assigned_id  = assigned_user.id").
					Where("(ppm_shipments.status IS NULL OR ppm_shipments.status NOT IN (?))", models.PPMShipmentStatusWaitingOnCustomer, models.PPMShipmentStatusNeedsCloseout, models.PPMShipmentStatusCloseoutComplete)
			}
		} else {
			if appCtx.Session().ActiveRole.RoleType == roles.RoleTypeTOO {
				query.Where("(moves.ppm_type IS NULL OR (moves.ppm_type = 'PARTIAL' or (moves.ppm_type = 'FULL' and origin_dl.provides_services_counseling = 'false')))")
			}
			// TODO  not sure we'll need this once we're in a situation where closeout param is always passed
			query.LeftJoin("ppm_shipments", "ppm_shipments.shipment_id = mto_shipments.id")
		}
	}

	for _, option := range options {
		if option != nil {
			option(query) // mutates
		}
	}

	// Pass zeros into paginate in this case. Which will give us 1 page and 20 per page respectively
	if params.Page == nil {
		params.Page = models.Int64Pointer(0)
	}
	if params.PerPage == nil {
		params.PerPage = models.Int64Pointer(0)
	}

	var groupByColumms []string
	groupByColumms = append(groupByColumms, "service_members.id", "orders.id", "origin_dl.id")

	if params.Sort != nil && *params.Sort == "originDutyLocation" {
		groupByColumms = append(groupByColumms, "origin_dl.name")
	}

	if params.Sort != nil && *params.Sort == "destinationDutyLocation" {
		groupByColumms = append(groupByColumms, "dest_dl.name")
	}

	if params.Sort != nil && *params.Sort == "originGBLOC" {
		groupByColumms = append(groupByColumms, "origin_to.id")
	}

	if params.Sort != nil && *params.Sort == "closeoutLocation" && !ppmCloseoutGblocs {
		groupByColumms = append(groupByColumms, "closeout_to.id")
	}

	if params.Sort != nil && *params.Sort == "ppmStatus" {
		groupByColumms = append(groupByColumms, "ppm_shipments.id")
	}

	if params.Sort != nil && *params.Sort == "counselingOffice" {
		groupByColumms = append(groupByColumms, "transportation_offices.id")
	}
	if params.Sort != nil && *params.Sort == "assignedTo" {
		groupByColumms = append(groupByColumms, "assigned_user.last_name", "assigned_user.first_name")
	}

	err = query.GroupBy("moves.id", groupByColumms...).Paginate(int(*params.Page), int(*params.PerPage)).All(&moves)
	if err != nil {
		return []models.Move{}, 0, err
	}
	// Get the count
	count := query.Paginator.TotalEntriesSize

	// Services Counselors in PPM Closeout GBLOCs should see their closeout GBLOC in the CloseoutOffice field for every
	// move.
	// We send that field back as a Transportation Office, and the transportation office's name gets displayed.
	// There are transportation offices corresponding to each of the closeout GBLOCs, but their names don't match what
	// we want displayed. So our options are to either fake a transportation office here that has the closeout GBLOC for
	// its name, or use a real transportation office, and have the frontend render it differently if it detects that
	// it's a closeout office. We went with the former approach to keep the logic contained on the backend.
	overwriteCloseoutOfficeWithGBLOC := ppmCloseoutGblocs && params.NeedsPPMCloseout != nil && *params.NeedsPPMCloseout
	closeoutOffice := models.TransportationOffice{
		Name: officeUserGbloc,
	}

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

		// Overwrite each move's closeout office if we are in a PPM closeout GBLOC
		if overwriteCloseoutOfficeWithGBLOC {
			moves[i].CloseoutOffice = &closeoutOffice
		}
	}

	return moves, count, nil
}

// this is a custom/temporary struct used in the below service object to convert moves consumed from the db into Go structs
type MoveWithCount struct {
	models.Move
	OrdersRaw                     json.RawMessage              `json:"orders" db:"orders"`
	Orders                        *models.Order                `json:"-"`
	MTOShipmentsRaw               json.RawMessage              `json:"mto_shipments" db:"mto_shipments"`
	MTOShipments                  *models.MTOShipments         `json:"-"`
	CounselingOfficeRaw           json.RawMessage              `json:"counseling_transportation_office" db:"counseling_transportation_office"`
	CounselingOffice              *models.TransportationOffice `json:"-"`
	TOOAssignedUserRaw            json.RawMessage              `json:"too_assigned" db:"too_assigned"`
	TOOAssignedUser               *models.OfficeUser           `json:"-"`
	TOODestinationAssignedUserRaw json.RawMessage              `json:"too_destination_assigned" db:"too_destination_assigned"`
	TOODestinationAssignedUser    *models.OfficeUser           `json:"-"`
	MTOServiceItemsRaw            json.RawMessage              `json:"mto_service_items" db:"mto_service_items"`
	MTOServiceItems               *models.MTOServiceItems      `json:"-"`
	TotalCount                    int64                        `json:"total_count" db:"total_count"`
}

type JSONB []byte

func (j *JSONB) UnmarshalJSON(data []byte) error {
	*j = data
	return nil
}

func (f orderFetcher) ListOriginRequestsOrders(appCtx appcontext.AppContext, officeUserID uuid.UUID, params *services.ListOrderParams) ([]models.Move, int, error) {
	var moves []models.Move
	var movesWithCount []MoveWithCount

	var officeUserGbloc string
	hasSafetyPrivilege := false
	if params.ViewAsGBLOC != nil {
		officeUserGbloc = *params.ViewAsGBLOC
	} else {
		var gblocErr error
		gblocFetcher := officeuser.NewOfficeUserGblocFetcher()
		officeUserGbloc, gblocErr = gblocFetcher.FetchGblocForOfficeUser(appCtx, officeUserID)
		if gblocErr != nil {
			return []models.Move{}, 0, gblocErr
		}
	}

	privileges, privErr := roles.FetchPrivilegesForUser(appCtx.DB(), appCtx.Session().UserID)
	if privErr == nil && privileges.HasPrivilege(roles.PrivilegeTypeSafety) {
		hasSafetyPrivilege = true
	} else if privErr != nil {
		appCtx.Logger().Error("Error retrieving user privileges", zap.Error(privErr))
	}

	err := appCtx.DB().RawQuery("SELECT * FROM get_origin_queue($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17)",
		officeUserGbloc,
		params.CustomerName,
		params.Edipi,
		params.Emplid,
		pq.Array(params.Status),
		params.Locator,
		params.RequestedMoveDate,
		params.AppearedInTOOAt,
		params.Branch,
		strings.Join(params.OriginDutyLocation, " "),
		params.CounselingOffice,
		params.AssignedTo,
		hasSafetyPrivilege,
		params.Page,
		params.PerPage,
		params.Sort,
		params.Order).
		All(&movesWithCount)

	if err != nil {
		return []models.Move{}, 0, err
	}

	var count int64
	if len(movesWithCount) > 0 {
		count = movesWithCount[0].TotalCount
	} else {
		count = 0
	}

	for i := range movesWithCount {
		var order models.Order
		if err := json.Unmarshal(movesWithCount[i].OrdersRaw, &order); err != nil {
			return nil, 0, fmt.Errorf("error unmarshaling orders JSON: %w", err)
		}
		movesWithCount[i].OrdersRaw = nil
		movesWithCount[i].Orders = &order

		var shipments models.MTOShipments
		if err := json.Unmarshal(movesWithCount[i].MTOShipmentsRaw, &shipments); err != nil {
			return nil, 0, fmt.Errorf("error unmarshaling shipments JSON: %w", err)
		}
		movesWithCount[i].MTOShipmentsRaw = nil
		movesWithCount[i].MTOShipments = &shipments

		var serviceItems models.MTOServiceItems
		if err := json.Unmarshal(movesWithCount[i].MTOServiceItemsRaw, &serviceItems); err != nil {
			return nil, 0, fmt.Errorf("error unmarshaling service items JSON: %w", err)
		}
		movesWithCount[i].MTOServiceItemsRaw = nil
		movesWithCount[i].MTOServiceItems = &serviceItems

		var counselingTransportationOffice models.TransportationOffice
		if err := json.Unmarshal(movesWithCount[i].CounselingOfficeRaw, &counselingTransportationOffice); err != nil {
			return nil, 0, fmt.Errorf("error unmarshaling counseling_transportation_office JSON: %w", err)
		}
		movesWithCount[i].CounselingOfficeRaw = nil
		movesWithCount[i].CounselingOffice = &counselingTransportationOffice

		var tooAssigned models.OfficeUser
		if err := json.Unmarshal(movesWithCount[i].TOOAssignedUserRaw, &tooAssigned); err != nil {
			return nil, 0, fmt.Errorf("error unmarshaling too_assigned JSON: %w", err)
		}
		movesWithCount[i].TOOAssignedUserRaw = nil
		movesWithCount[i].TOOAssignedUser = &tooAssigned
	}

	for _, moveWithCount := range movesWithCount {
		var move models.Move
		if err := copier.Copy(&move, &moveWithCount); err != nil {
			return nil, 0, fmt.Errorf("error copying movesWithCount into Moves struct: %w", err)
		}
		moves = append(moves, move)
	}

	return moves, int(count), nil
}

func (f orderFetcher) ListDestinationRequestsOrders(appCtx appcontext.AppContext, officeUserID uuid.UUID, role roles.RoleType, params *services.ListOrderParams) ([]models.Move, int, error) {
	var moves []models.Move
	var movesWithCount []MoveWithCount

	// getting the office user's GBLOC
	var officeUserGbloc string
	hasSafetyPrivilege := false
	if params.ViewAsGBLOC != nil {
		officeUserGbloc = *params.ViewAsGBLOC
	} else {
		var gblocErr error
		gblocFetcher := officeuser.NewOfficeUserGblocFetcher()
		officeUserGbloc, gblocErr = gblocFetcher.FetchGblocForOfficeUser(appCtx, officeUserID)
		if gblocErr != nil {
			return []models.Move{}, 0, gblocErr
		}
	}
	privileges, privErr := roles.FetchPrivilegesForUser(appCtx.DB(), appCtx.Session().UserID)
	if privErr == nil && privileges.HasPrivilege(roles.PrivilegeTypeSafety) {
		hasSafetyPrivilege = true
	} else if privErr != nil {
		appCtx.Logger().Error("Error retrieving user privileges", zap.Error(privErr))
	}
	// calling the database function with all passed in parameters
	err := appCtx.DB().RawQuery("SELECT * FROM get_destination_queue($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17)",
		officeUserGbloc,
		params.CustomerName,
		params.Edipi,
		params.Emplid,
		pq.Array(params.Status),
		params.Locator,
		params.RequestedMoveDate,
		params.AppearedInTOOAt,
		params.Branch,
		params.DestinationDutyLocation,
		params.CounselingOffice,
		params.AssignedTo,
		hasSafetyPrivilege,
		params.Page,
		params.PerPage,
		params.Sort,
		params.Order).
		All(&movesWithCount)

	if err != nil {
		return []models.Move{}, 0, err
	}

	// each row is sent back with the total count from the db func, so we will take the value from the first one
	var count int64
	if len(movesWithCount) > 0 {
		count = movesWithCount[0].TotalCount
	} else {
		count = 0
	}

	// we have to manually loop through each move and populate the nested objects that the queue uses/needs
	for i := range movesWithCount {
		// populating Move.Orders struct
		var order models.Order
		if err := json.Unmarshal(movesWithCount[i].OrdersRaw, &order); err != nil {
			return nil, 0, fmt.Errorf("error unmarshaling orders JSON: %w", err)
		}
		movesWithCount[i].OrdersRaw = nil
		movesWithCount[i].Orders = &order

		// populating Move.MTOShipments array
		var shipments models.MTOShipments
		if err := json.Unmarshal(movesWithCount[i].MTOShipmentsRaw, &shipments); err != nil {
			return nil, 0, fmt.Errorf("error unmarshaling shipments JSON: %w", err)
		}
		movesWithCount[i].MTOShipmentsRaw = nil
		movesWithCount[i].MTOShipments = &shipments

		// populating Moves.CounselingOffice struct
		var counselingTransportationOffice models.TransportationOffice
		if err := json.Unmarshal(movesWithCount[i].CounselingOfficeRaw, &counselingTransportationOffice); err != nil {
			return nil, 0, fmt.Errorf("error unmarshaling counseling_transportation_office JSON: %w", err)
		}
		movesWithCount[i].CounselingOfficeRaw = nil
		movesWithCount[i].CounselingOffice = &counselingTransportationOffice

		// populating Moves.TOOAssigned struct
		var tooAssigned models.OfficeUser
		if err := json.Unmarshal(movesWithCount[i].TOODestinationAssignedUserRaw, &tooAssigned); err != nil {
			return nil, 0, fmt.Errorf("error unmarshaling too_assigned JSON: %w", err)
		}
		movesWithCount[i].TOODestinationAssignedUserRaw = nil
		movesWithCount[i].TOODestinationAssignedUser = &tooAssigned

		// populating Moves.MTOServiceItems struct
		var serviceItems models.MTOServiceItems
		if movesWithCount[i].MTOServiceItemsRaw != nil {
			if err := json.Unmarshal(movesWithCount[i].MTOServiceItemsRaw, &serviceItems); err != nil {
				return nil, 0, fmt.Errorf("error unmarshaling mto_service_items JSON: %w", err)
			}
		}
		movesWithCount[i].MTOServiceItemsRaw = nil
		movesWithCount[i].MTOServiceItems = &serviceItems
	}

	// the handler consumes a Move object and NOT the MoveWithCount struct used in this func
	// so we have to copy our custom struct into the Move struct
	for _, moveWithCount := range movesWithCount {
		var move models.Move
		if err := copier.Copy(&move, &moveWithCount); err != nil {
			return nil, 0, fmt.Errorf("error copying movesWithCount into Moves struct: %w", err)
		}
		moves = append(moves, move)
	}

	return moves, int(count), nil
}

func (f orderFetcher) ListAllOrderLocations(appCtx appcontext.AppContext, officeUserID uuid.UUID, params *services.ListOrderParams) ([]models.Move, error) {
	var moves []models.Move
	var err error
	var officeUserGbloc string

	if params.ViewAsGBLOC != nil {
		officeUserGbloc = *params.ViewAsGBLOC
	} else {
		gblocFetcher := officeuser.NewOfficeUserGblocFetcher()
		officeUserGbloc, err = gblocFetcher.FetchGblocForOfficeUser(appCtx, officeUserID)
	}

	if err != nil {
		return []models.Move{}, err
	}

	privileges, err := roles.FetchPrivilegesForUser(appCtx.DB(), appCtx.Session().UserID)
	if err != nil {
		appCtx.Logger().Error("Error retreiving user privileges", zap.Error(err))
	}

	needsCounseling := false
	if len(params.Status) > 0 {
		for _, status := range params.Status {
			if status == string(models.MoveStatusNeedsServiceCounseling) {
				needsCounseling = true
			}
		}
	}

	ppmCloseoutGblocs := officeUserGbloc == "NAVY" || officeUserGbloc == "TVCB" || officeUserGbloc == "USCG"

	// Services Counselors in closeout GBLOCs should only see closeout moves
	if needsCounseling && ppmCloseoutGblocs && params.NeedsPPMCloseout != nil && !*params.NeedsPPMCloseout {
		return []models.Move{}, nil
	}

	branchQuery := branchFilter(params.Branch, needsCounseling, ppmCloseoutGblocs)

	// If the user is associated with the USMC GBLOC we want to show them ALL the USMC moves, so let's override here.
	// We also only want to do the gbloc filtering thing if we aren't a USMC user, which we cover with the else.
	// var gblocQuery QueryOption
	if officeUserGbloc == "USMC" && !needsCounseling {
		branchQuery = branchFilter(models.StringPointer(string(models.AffiliationMARINES)), needsCounseling, ppmCloseoutGblocs)
	}

	// We need to use three different GBLOC filter queries because:
	//  - The Services Counselor queue filters based on the GBLOC of the origin duty location's
	//    transportation office
	//  - There is a separate queue for the GBLOCs: NAVY, TVCB and USCG. These GBLOCs are used by SC doing PPM Closeout
	//  - The TOO queue uses the GBLOC we get from examining the postal code of the first shipment's
	//    pickup address. However, if that shipment happens to be an NTS-Release, we instead drop
	//    back to the GBLOC of the origin duty location's transportation office since an NTS-Release
	//    does not populate the pickup address field.
	var gblocQuery QueryOption
	if ppmCloseoutGblocs {
		gblocQuery = gblocFilterForPPMCloseoutForNavyMarineAndCG(&officeUserGbloc)
	} else if needsCounseling {
		gblocQuery = gblocFilterForSC(&officeUserGbloc)
	} else if params.NeedsPPMCloseout != nil && *params.NeedsPPMCloseout {
		gblocQuery = gblocFilterForSCinArmyAirForce(&officeUserGbloc)
	} else {
		gblocQuery = gblocFilterForTOO(&officeUserGbloc)
	}
	moveStatusQuery := moveStatusFilter(params.Status)
	// Adding to an array so we can iterate over them and apply the filters after the query structure is set below
	options := [15]QueryOption{branchQuery, moveStatusQuery, gblocQuery}

	var query *pop.Query
	if ppmCloseoutGblocs {
		query = appCtx.DB().Q().Scope(utilities.ExcludeDeletedScope(models.MTOShipment{})).EagerPreload(
			"Orders.ServiceMember",
			"Orders.NewDutyLocation.Address",
			"Orders.OriginDutyLocation.Address",
			"Orders.Entitlement",
			"Orders.OrdersType",
			"MTOShipments.PPMShipment",
			"LockedByOfficeUser",
		).InnerJoin("orders", "orders.id = moves.orders_id").
			InnerJoin("service_members", "orders.service_member_id = service_members.id").
			InnerJoin("mto_shipments", "moves.id = mto_shipments.move_id").
			InnerJoin("ppm_shipments", "ppm_shipments.shipment_id = mto_shipments.id").
			InnerJoin("duty_locations as origin_dl", "orders.origin_duty_location_id = origin_dl.id").
			LeftJoin("duty_locations as dest_dl", "dest_dl.id = orders.new_duty_location_id").
			LeftJoin("transportation_offices as closeout_to", "closeout_to.id = moves.closeout_office_id").
			LeftJoin("office_users", "office_users.id = moves.locked_by").
			Where("show = ?", models.BoolPointer(true))

		if !privileges.HasPrivilege(roles.PrivilegeTypeSafety) {
			query.Where("orders.orders_type != (?)", roles.PrivilegeSearchTypeSafety)
		}
	} else {
		query = appCtx.DB().Q().Scope(utilities.ExcludeDeletedScope(models.MTOShipment{})).EagerPreload(
			"Orders.ServiceMember",
			"Orders.NewDutyLocation.Address",
			"Orders.OriginDutyLocation.Address",
			"Orders.Entitlement",
			"Orders.OrdersType",
			"MTOShipments",
			"MTOServiceItems",
			"ShipmentGBLOC",
			"MTOShipments.PPMShipment",
			"CloseoutOffice",
			"LockedByOfficeUser",
		).InnerJoin("orders", "orders.id = moves.orders_id").
			InnerJoin("service_members", "orders.service_member_id = service_members.id").
			InnerJoin("mto_shipments", "moves.id = mto_shipments.move_id").
			InnerJoin("duty_locations as origin_dl", "orders.origin_duty_location_id = origin_dl.id").
			// Need to use left join because some duty locations do not have transportation offices
			LeftJoin("transportation_offices as origin_to", "origin_dl.transportation_office_id = origin_to.id").
			// If a customer puts in an invalid ZIP for their pickup address, it won't show up in this view,
			// and we don't want it to get hidden from services counselors.
			LeftJoin("move_to_gbloc", "move_to_gbloc.move_id = moves.id").
			LeftJoin("duty_locations as dest_dl", "dest_dl.id = orders.new_duty_location_id").
			LeftJoin("transportation_offices as closeout_to", "closeout_to.id = moves.closeout_office_id").
			LeftJoin("office_users", "office_users.id = moves.locked_by").
			Where("show = ?", models.BoolPointer(true))

		if !privileges.HasPrivilege(roles.PrivilegeTypeSafety) {
			query.Where("orders.orders_type != (?)", roles.PrivilegeSearchTypeSafety)
		}

		if params.NeedsPPMCloseout != nil {
			if *params.NeedsPPMCloseout {
				query.InnerJoin("ppm_shipments", "ppm_shipments.shipment_id = mto_shipments.id").
					Where("ppm_shipments.status IN (?)", models.PPMShipmentStatusNeedsCloseout, models.PPMShipmentStatusCloseoutComplete).
					Where("service_members.affiliation NOT IN (?)", models.AffiliationNAVY, models.AffiliationMARINES, models.AffiliationCOASTGUARD)
			} else {
				query.LeftJoin("ppm_shipments", "ppm_shipments.shipment_id = mto_shipments.id").
					Where("(ppm_shipments.status IS NULL OR ppm_shipments.status NOT IN (?))", models.PPMShipmentStatusWaitingOnCustomer, models.PPMShipmentStatusNeedsCloseout, models.PPMShipmentStatusCloseoutComplete)
			}
		} else {
			if appCtx.Session().ActiveRole.RoleType == roles.RoleTypeTOO {
				query.Where("(moves.ppm_type IS NULL OR (moves.ppm_type = (?) or (moves.ppm_type = (?) and origin_dl.provides_services_counseling = 'false')))", models.MovePPMTypePARTIAL, models.MovePPMTypeFULL)
			}
			query.LeftJoin("ppm_shipments", "ppm_shipments.shipment_id = mto_shipments.id")
		}
	}

	for _, option := range options {
		if option != nil {
			option(query) // mutates
		}
	}

	var groupByColumms []string
	groupByColumms = append(groupByColumms, "service_members.id", "orders.id", "origin_dl.id")

	err = query.GroupBy("moves.id", groupByColumms...).All(&moves)
	if err != nil {
		return []models.Move{}, err
	}

	return moves, nil
}

// NewOrderFetcher creates a new struct with the service dependencies
func NewOrderFetcher(weightAllotmentFetcher services.WeightAllotmentFetcher) services.OrderFetcher {
	return &orderFetcher{waf: weightAllotmentFetcher}
}

// FetchOrder retrieves an Order for a given UUID
func (f orderFetcher) FetchOrder(appCtx appcontext.AppContext, orderID uuid.UUID) (*models.Order, error) {
	order := &models.Order{}
	err := appCtx.DB().Q().Eager(
		"ServiceMember.BackupContacts",
		"ServiceMember.ResidentialAddress",
		"NewDutyLocation.Address.Country",
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

	// Construct weight allotted if grade is present
	if order.Grade != nil {
		allotment, err := f.waf.GetWeightAllotment(appCtx, string(*order.Grade), order.OrdersType)
		if err != nil {
			return nil, err
		}
		order.Entitlement.WeightAllotted = &allotment
	}

	// Due to a bug in pop (https://github.com/gobuffalo/pop/issues/578), we
	// cannot eager load the address as "OriginDutyLocation.Address" because
	// OriginDutyLocation is a pointer.
	if order.OriginDutyLocation != nil {
		err = appCtx.DB().Load(order.OriginDutyLocation, "Address", "Address.Country")
		if err != nil {
			return order, err
		}
	}

	return order, nil
}

// These are a bunch of private functions that are used to cobble our list Orders filters together.
func branchFilter(branch *string, needsCounseling bool, ppmCloseoutGblocs bool) QueryOption {
	return func(query *pop.Query) {
		if branch == nil && !needsCounseling && !ppmCloseoutGblocs {
			query.Where("service_members.affiliation != ?", models.AffiliationMARINES)
		}
		if branch != nil {
			query.Where("service_members.affiliation ILIKE ?", *branch)
		}
	}
}

func customerNameFilter(name *string) QueryOption {
	return func(query *pop.Query) {
		if name == nil {
			return
		}

		// Remove "," that user may enter between names (displayed on frontend column)
		nameQueryParam := *name
		removeCharsRegex := regexp.MustCompile("[,]+")
		nameQueryParam = removeCharsRegex.ReplaceAllString(nameQueryParam, "")
		nameQueryParam = fmt.Sprintf("%%%s%%", nameQueryParam)

		// Search for partial within both (last first) and (first last) in one go
		query.Where("(service_members.last_name || ' ' || service_members.first_name || service_members.first_name || ' ' || service_members.last_name) ILIKE ?", nameQueryParam)
	}
}

func dodIDFilter(dodID *string) QueryOption {
	return func(query *pop.Query) {
		if dodID != nil {
			query.Where("service_members.edipi = ?", dodID)
		}
	}
}

func emplidFilter(emplid *string) QueryOption {
	return func(query *pop.Query) {
		if emplid != nil {
			query.Where("service_members.emplid = ?", emplid)
		}
	}
}

func locatorFilter(locator *string) QueryOption {
	return func(query *pop.Query) {
		if locator != nil {
			query.Where("moves.locator = ?", strings.ToUpper(*locator))
		}
	}
}

func originDutyLocationFilter(originDutyLocation []string) QueryOption {
	return func(query *pop.Query) {
		if len(originDutyLocation) > 0 {
			query.Where("origin_dl.name ILIKE ?", "%"+strings.Join(originDutyLocation, " ")+"%")
		}
	}
}

func destinationDutyLocationFilter(destinationDutyLocation *string) QueryOption {
	return func(query *pop.Query) {
		if destinationDutyLocation != nil {
			nameSearch := fmt.Sprintf("%s%%", *destinationDutyLocation)
			query.Where("dest_dl.name ILIKE ?", nameSearch)
		}
	}
}

func counselingOfficeFilter(office *string) QueryOption {
	return func(query *pop.Query) {
		if office != nil {
			query.Where("transportation_offices.name ILIKE ?", "%"+*office+"%")
		}
	}
}

func moveStatusFilter(statuses []string) QueryOption {
	return func(query *pop.Query) {
		// If we have statuses let's use them
		if len(statuses) > 0 {
			query.Where("moves.status IN (?)", statuses)
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

func appearedInTOOAtFilter(appearedInTOOAt *time.Time) QueryOption {
	return func(query *pop.Query) {
		if appearedInTOOAt != nil {
			start := appearedInTOOAt.Format(RFC3339Micro)
			// Between is inclusive, so the end date is set to 1 microsecond prior to the next day
			appearedInTOOAtEnd := appearedInTOOAt.AddDate(0, 0, 1).Add(-1 * time.Microsecond)
			end := appearedInTOOAtEnd.Format(RFC3339Micro)
			query.Where("(moves.submitted_at between ? AND ? OR moves.service_counseling_completed_at between ? AND ? OR moves.approvals_requested_at between ? AND ?)", start, end, start, end, start, end)
		}
	}
}

func requestedMoveDateFilter(requestedMoveDate *string) QueryOption {
	return func(query *pop.Query) {
		if requestedMoveDate != nil {
			query.Where("(mto_shipments.requested_pickup_date = ? OR ppm_shipments.expected_departure_date = ? OR (mto_shipments.shipment_type = 'HHG_OUTOF_NTS' AND mto_shipments.requested_delivery_date = ?))", *requestedMoveDate, *requestedMoveDate, *requestedMoveDate)
		}
	}
}

func closeoutInitiatedFilter(closeoutInitiated *time.Time) QueryOption {
	return func(query *pop.Query) {
		if closeoutInitiated != nil {
			// Between is inclusive, so the end date is set to 1 microsecond prior to the next day
			closeoutInitiatedEnd := closeoutInitiated.AddDate(0, 0, 1).Add(-1 * time.Microsecond)
			query.Having("MAX(ppm_shipments.submitted_at) between ? and ?", closeoutInitiated.Format(RFC3339Micro), closeoutInitiatedEnd.Format(RFC3339Micro))
		}
	}
}

func ppmTypeFilter(ppmType *string) QueryOption {
	return func(query *pop.Query) {
		if ppmType != nil {
			query.Where("moves.ppm_type = ?", *ppmType)
		}
	}
}

func ppmStatusFilter(ppmStatus *string) QueryOption {
	return func(query *pop.Query) {
		if ppmStatus != nil {
			query.Where("ppm_shipments.status = ?", *ppmStatus)
		}
	}
}

func assignedUserFilter(user *string) QueryOption {
	return func(query *pop.Query) {
		if user != nil {
			nameSearch := fmt.Sprintf("%s%%", *user)
			query.Where("assigned_user.last_name ILIKE ?", nameSearch)
		}
	}
}

func closeoutLocationFilter(closeoutLocation *string, ppmCloseoutGblocs bool) QueryOption {
	return func(query *pop.Query) {
		// If the office user is in a closeout gbloc, every single result they're seeing will have
		// the same closeout location, which will be identical to their gbloc, so there's no reason
		// to do this search.
		if closeoutLocation != nil && !ppmCloseoutGblocs {
			nameSearch := fmt.Sprintf("%s%%", *closeoutLocation)
			query.Where("closeout_to.name ILIKE ?", nameSearch)
		}
	}
}

func gblocFilterForSC(gbloc *string) QueryOption {
	// The SC should only see moves where the origin duty location's GBLOC matches the given GBLOC.
	return func(query *pop.Query) {
		if gbloc != nil {
			query.Where("orders.gbloc = ?", *gbloc)
		}
	}
}

func gblocFilterForTOO(gbloc *string) QueryOption {
	// The TOO should only see moves where the GBLOC for the first shipment's pickup address matches the given GBLOC
	// unless we're dealing with an NTS-Release shipment. For NTS-Release shipments, we drop back to looking at the
	// origin duty location's GBLOC since an NTS-Release does not populate the pickup address.
	return func(query *pop.Query) {
		if gbloc != nil {
			// Note: extra parens necessary to keep precedence correct when AND'ing all filters together.
			query.Where("((mto_shipments.shipment_type != ? AND move_to_gbloc.gbloc = ?) OR (mto_shipments.shipment_type = ? AND orders.gbloc = ?))",
				models.MTOShipmentTypeHHGOutOfNTS, *gbloc, models.MTOShipmentTypeHHGOutOfNTS, *gbloc)
		}
	}
}

func gblocFilterForSCinArmyAirForce(gbloc *string) QueryOption {
	// A services counselor in a transportation office that provides Services Counseling should see all moves with PPMs that have selected a closeout office that matches the GBLOC of their transportation office that is in waiting for customer, needs payment approval, or payment approved statuses. The Army and Air Force SCs should see moves in the PPM closeout Tab when the postal code or origin duty station is in a different GBLOC.
	return func(query *pop.Query) {
		if gbloc != nil {
			query.Where("mto_shipments.shipment_type = ? AND closeout_to.gbloc = ?", models.MTOShipmentTypePPM, *gbloc)
		}
	}
}

func gblocFilterForPPMCloseoutForNavyMarineAndCG(gbloc *string) QueryOption {
	// For PPM Closeout the SC should see moves that have ppm shipments
	// And the GBLOC should map to the service member's affiliation
	navyGbloc := "NAVY"
	tvcbGbloc := "TVCB"
	uscgGbloc := "USCG"
	return func(query *pop.Query) {
		if gbloc != nil {
			if *gbloc == navyGbloc {
				query.Where("mto_shipments.shipment_type = ? AND service_members.affiliation = ? AND ppm_shipments.status = ?", models.MTOShipmentTypePPM, models.AffiliationNAVY, models.PPMShipmentStatusNeedsCloseout)
			} else if *gbloc == tvcbGbloc {
				query.Where("mto_shipments.shipment_type = ? AND service_members.affiliation = ? AND ppm_shipments.status = ?", models.MTOShipmentTypePPM, models.AffiliationMARINES, models.PPMShipmentStatusNeedsCloseout)
			} else if *gbloc == uscgGbloc {
				query.Where("mto_shipments.shipment_type = ? AND service_members.affiliation = ? AND ppm_shipments.status = ?", models.MTOShipmentTypePPM, models.AffiliationCOASTGUARD, models.PPMShipmentStatusNeedsCloseout)
			}
		}
	}
}

func sortOrder(sort *string, order *string, ppmCloseoutGblocs bool) QueryOption {
	parameters := map[string]string{
		"customerName":            "(service_members.last_name || ' ' || service_members.first_name)",
		"edipi":                   "service_members.edipi",
		"emplid":                  "service_members.emplid",
		"branch":                  "service_members.affiliation",
		"locator":                 "moves.locator",
		"status":                  "moves.status",
		"submittedAt":             "moves.submitted_at",
		"appearedInTooAt":         "GREATEST(moves.submitted_at, moves.service_counseling_completed_at, moves.approvals_requested_at)",
		"originDutyLocation":      "origin_dl.name",
		"destinationDutyLocation": "dest_dl.name",
		"requestedMoveDate":       "LEAST(COALESCE(MIN(mto_shipments.requested_pickup_date), 'infinity'), COALESCE(MIN(ppm_shipments.expected_departure_date), 'infinity'), COALESCE(MIN(mto_shipments.requested_delivery_date), 'infinity'))",
		"originGBLOC":             "origin_to.gbloc",
		"ppmType":                 "moves.ppm_type",
		"ppmStatus":               "ppm_shipments.status",
		"closeoutLocation":        "closeout_to.name",
		"closeoutInitiated":       "MAX(ppm_shipments.submitted_at)",
		"counselingOffice":        "transportation_offices.name",
		"assignedTo":              "assigned_user.last_name,assigned_user.first_name",
	}

	return func(query *pop.Query) {
		// If we have a sort and order defined let's use it. Otherwise we'll use our default status desc sort order.
		if sort != nil && order != nil {
			// If an office user is in a closeout GBLOC, every move they see will have the same closeout location
			// so we can skip the sorting.
			if *sort == "closeoutLocation" && ppmCloseoutGblocs {
				return
			}
			if sortTerm, ok := parameters[*sort]; ok {
				if *sort == "customerName" {
					query.Order(fmt.Sprintf("service_members.last_name %s, service_members.first_name %s", *order, *order))
				} else if *sort == "assignedTo" {
					query.Order(fmt.Sprintf("assigned_user.last_name %s, assigned_user.first_name %s", *order, *order))
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

// When a queue is sorted by a non-unique value (ex: status, branch) the order within each value is inconsistent at different page sizes
// Adding a secondary sort ensures a consistent order within the primary sort column
func secondarySortOrder(sort *string) QueryOption {
	return func(query *pop.Query) {
		if sort == nil || *sort != "locator" {
			query.Order("moves.locator asc")
		}
	}
}

// We want to filter out any moves that have ONLY destination type requests to them, such as destination SIT, shuttle, out of the
// task order queue. If the moves have origin SIT, excess weight risks, or sit extensions with origin SIT service items, they
// should still appear in the task order queue, which is what this query looks for
func tooQueueOriginRequestsFilter(role roles.RoleType) QueryOption {
	return func(query *pop.Query) {
		if role == roles.RoleTypeTOO {
			baseQuery := `
			(mto_shipments.status IN ('SUBMITTED','APPROVALS_REQUESTED')
				-- keep moves in the TOO queue if they have an unacknowledged excess weight risk
        		OR (mto_shipments.status = 'APPROVED'
        			AND
                		(
                        	moves.excess_weight_qualified_at IS NOT NULL
                        	AND moves.excess_weight_acknowledged_at IS NULL
		                )
        		        OR (
                	        moves.excess_unaccompanied_baggage_weight_qualified_at IS NOT NULL
                    	    AND moves.excess_unaccompanied_baggage_weight_acknowledged_at IS NULL
                		)
        		)
    		)
			AND
			-- check for moves with destination requests and NOT origin requests, then return the inverse for the TOO queue with the NOT wrapped around the query
			NOT (
					-- check for moves with destination requests
					(
						-- moves with submitted destination SIT or shuttle submitted service items
						EXISTS (
							SELECT 1
							FROM mto_service_items msi
							JOIN re_services rs ON msi.re_service_id = rs.id
							WHERE msi.mto_shipment_id = mto_shipments.id
							AND msi.status = 'SUBMITTED'
							AND rs.code IN ('DDFSIT', 'DDASIT', 'DDDSIT', 'DDSHUT', 'DDSFSC',
											'IDFSIT', 'IDASIT', 'IDDSIT', 'IDSHUT', 'IDSFSC')
						)
						-- requested shipment address update
						OR EXISTS (
							SELECT 1
							FROM shipment_address_updates sau
							WHERE sau.shipment_id = mto_shipments.id
							AND sau.status = 'REQUESTED'
						)
						-- Moves with SIT extensions and ONLY destination SIT service items we filter out of TOO queue
						OR (
							EXISTS (
								SELECT 1
								FROM sit_extensions se
								JOIN mto_service_items msi ON se.mto_shipment_id = msi.mto_shipment_id
								JOIN re_services rs ON msi.re_service_id = rs.id
								WHERE se.mto_shipment_id = mto_shipments.id
								AND se.status = 'PENDING'
								AND rs.code IN ('DDFSIT', 'DDASIT', 'DDDSIT', 'DDSHUT', 'DDSFSC',
												'IDFSIT', 'IDASIT', 'IDDSIT', 'IDSHUT', 'IDSFSC')
							)
							-- make sure there are NO origin SIT service items (otherwise, it should be in both queues)
							AND NOT EXISTS (
								SELECT 1
								FROM mto_service_items msi
								JOIN re_services rs ON msi.re_service_id = rs.id
								WHERE msi.mto_shipment_id = mto_shipments.id
								AND msi.status = 'SUBMITTED'
								AND rs.code IN ('ICRT', 'IUBPK', 'IOFSIT', 'IOASIT', 'IOPSIT', 'IOSHUT',
												'IHUPK', 'IUCRT', 'DCRT', 'MS', 'CS', 'DOFSIT', 'DOASIT',
												'DOPSIT', 'DOSFSC', 'IOSFSC', 'DUPK', 'DUCRT', 'DOSHUT',
												'FSC', 'DMHF', 'DBTF', 'DBHF', 'IBTF', 'IBHF', 'DCRTSA',
												'DLH', 'DOP', 'DPK', 'DSH', 'DNPK', 'INPK', 'UBP',
												'ISLH', 'POEFSC', 'PODFSC', 'IHPK')
							)
						)
					)
					-- check for moves with origin requests or conditions where move should appear in TOO queue
					AND NOT (
						-- keep moves in the TOO queue with origin submitted service items
						EXISTS (
							SELECT 1
							FROM mto_service_items msi
							JOIN re_services rs ON msi.re_service_id = rs.id
							WHERE msi.mto_shipment_id = mto_shipments.id
							AND msi.status = 'SUBMITTED'
							AND rs.code IN ('ICRT', 'IUBPK', 'IOFSIT', 'IOASIT', 'IOPSIT', 'IOSHUT',
											'IHUPK', 'IUCRT', 'DCRT', 'MS', 'CS', 'DOFSIT', 'DOASIT',
											'DOPSIT', 'DOSFSC', 'IOSFSC', 'DUPK', 'DUCRT', 'DOSHUT',
											'FSC', 'DMHF', 'DBTF', 'DBHF', 'IBTF', 'IBHF', 'DCRTSA',
											'DLH', 'DOP', 'DPK', 'DSH', 'DNPK', 'INPK', 'UBP',
											'ISLH', 'POEFSC', 'PODFSC', 'IHPK')
						)
					)
				)

            `
			query.Where(baseQuery)
		}
	}
}
