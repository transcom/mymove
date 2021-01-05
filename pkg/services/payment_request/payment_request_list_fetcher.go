package paymentrequest

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/go-openapi/swag"

	"github.com/gobuffalo/pop/v5"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/services"
	officeuser "github.com/transcom/mymove/pkg/services/office_user"

	"github.com/transcom/mymove/pkg/models"
)

type paymentRequestListFetcher struct {
	db *pop.Connection
}

var parameters = map[string]string{
	"lastName":    "service_members.last_name",
	"dodID":       "service_members.edipi",
	"submittedAt": "payment_requests.created_at",
	"branch":      "service_members.affiliation",
	"locator":     "moves.locator",
	"status":      "payment_requests.status",
	"age":         "payment_requests.created_at",
}

// NewPaymentRequestListFetcher returns a new payment request list fetcher
func NewPaymentRequestListFetcher(db *pop.Connection) services.PaymentRequestListFetcher {
	return &paymentRequestListFetcher{db}
}

// QueryOption defines the type for the functional arguments passed to ListMoveOrders
type QueryOption func(*pop.Query)

func (f *paymentRequestListFetcher) FetchPaymentRequestList(officeUserID uuid.UUID, params *services.FetchPaymentRequestListParams) (*models.PaymentRequests, int, error) {

	gblocFetcher := officeuser.NewOfficeUserGblocFetcher(f.db)
	gbloc, gblocErr := gblocFetcher.FetchGblocForOfficeUser(officeUserID)
	if gblocErr != nil {
		return &models.PaymentRequests{}, 0, gblocErr
	}

	paymentRequests := models.PaymentRequests{}
	query := f.db.Q().Eager(
		"MoveTaskOrder.Orders.OriginDutyStation",
		"MoveTaskOrder.Orders.ServiceMember",
	).
		InnerJoin("moves", "payment_requests.move_id = moves.id").
		InnerJoin("orders", "orders.id = moves.orders_id").
		InnerJoin("service_members", "orders.service_member_id = service_members.id").
		InnerJoin("duty_stations", "duty_stations.id = orders.origin_duty_station_id").
		InnerJoin("transportation_offices", "transportation_offices.id = duty_stations.transportation_office_id").
		Where("moves.show = ?", swag.Bool(true))

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
	statusQuery := paymentRequestsStatusFilter(params.Status)
	submittedAtQuery := submittedAtFilter(params.SubmittedAt)
	orderQuery := sortOrder(params.Sort, params.Order)

	options := [9]QueryOption{branchQuery, locatorQuery, dodIDQuery, lastNameQuery, dutyStationQuery, statusQuery, submittedAtQuery, gblocQuery, orderQuery}

	for _, option := range options {
		if option != nil {
			option(query)
		}
	}

	if params.Page == nil {
		params.Page = swag.Int64(1)
	}

	if params.PerPage == nil {
		params.PerPage = swag.Int64(20)
	}

	err := query.GroupBy("payment_requests.id, service_members.id, moves.id").Paginate(int(*params.Page), int(*params.PerPage)).All(&paymentRequests)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, 0, services.NotFoundError{}
		default:
			return nil, 0, err
		}
	}

	for i := range paymentRequests {
		// Due to a bug in pop (https://github.com/gobuffalo/pop/issues/578), we
		// cannot eager load the address as "OriginDutyStation.Address" because
		// OriginDutyStation is a pointer.
		if originDutyStation := paymentRequests[i].MoveTaskOrder.Orders.OriginDutyStation; originDutyStation != nil {
			f.db.Load(originDutyStation, "TransportationOffice")
		}
	}

	// Get the count
	count := query.Paginator.TotalEntriesSize

	return &paymentRequests, count, nil
}

func (f *paymentRequestListFetcher) FetchPaymentRequestListByMove(officeUserID uuid.UUID, locator string) (*models.PaymentRequests, error) {
	gblocFetcher := officeuser.NewOfficeUserGblocFetcher(f.db)
	gbloc, gblocErr := gblocFetcher.FetchGblocForOfficeUser(officeUserID)
	if gblocErr != nil {
		return &models.PaymentRequests{}, gblocErr
	}

	paymentRequests := models.PaymentRequests{}

	query := f.db.Q().EagerPreload("PaymentServiceItems").
		InnerJoin("moves", "payment_requests.move_id = moves.id").
		InnerJoin("orders", "orders.id = moves.orders_id").
		InnerJoin("service_members", "orders.service_member_id = service_members.id").
		InnerJoin("duty_stations", "duty_stations.id = orders.origin_duty_station_id").
		InnerJoin("transportation_offices", "transportation_offices.id = duty_stations.transportation_office_id").
		Where("moves.show = ?", swag.Bool(true))

	var branchQuery QueryOption
	// If the user is associated with the USMC GBLOC we want to show them ALL the USMC moves, so let's override here.
	// We also only want to do the gbloc filtering thing if we aren't a USMC user, which we cover with the else.
	var gblocQuery QueryOption
	if gbloc == "USMC" {
		branchQuery = branchFilter(swag.String(string(models.AffiliationMARINES)))
	} else {
		gblocQuery = gblocFilter(gbloc)
	}
	locatorQuery := locatorFilter(&locator)

	options := [3]QueryOption{branchQuery, gblocQuery, locatorQuery}

	for _, option := range options {
		if option != nil {
			option(query)
		}
	}

	err := query.All(&paymentRequests)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, services.NotFoundError{}
		default:
			return nil, err
		}
	}

	return &paymentRequests, nil
}

func orderName(query *pop.Query, order *string) *pop.Query {
	return query.Order(fmt.Sprintf("service_members.last_name %s, service_members.first_name %s", *order, *order))
}

func reverseOrder(order *string) string {
	if *order == "asc" {
		return "desc"
	}
	return "asc"
}

func sortOrder(sort *string, order *string) QueryOption {
	return func(query *pop.Query) {
		if sort != nil && order != nil {
			sortTerm := parameters[*sort]
			if *sort == "lastName" {
				orderName(query, order)
			} else if *sort == "age" {
				query = query.Order(fmt.Sprintf("%s %s", sortTerm, reverseOrder(order)))
			} else {
				query = query.Order(fmt.Sprintf("%s %s", sortTerm, *order))
			}
		} else {
			query = query.Order("payment_requests.created_at asc")
		}
	}
}

func branchFilter(branch *string) QueryOption {
	return func(query *pop.Query) {
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

func locatorFilter(locator *string) QueryOption {
	return func(query *pop.Query) {
		if locator != nil {
			query = query.Where("moves.locator = ?", *locator)
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

func submittedAtFilter(submittedAt *string) QueryOption {
	return func(query *pop.Query) {
		if submittedAt != nil {
			query = query.Where("CAST(payment_requests.created_at AS DATE) = ?", *submittedAt)
		}
	}
}

func gblocFilter(gbloc string) QueryOption {
	return func(query *pop.Query) {
		query = query.Where("transportation_offices.gbloc = ?", gbloc)
	}
}

func paymentRequestsStatusFilter(statuses []string) QueryOption {
	return func(query *pop.Query) {
		var translatedStatuses []string
		if len(statuses) > 0 {
			for _, status := range statuses {
				if strings.EqualFold(status, "Payment requested") {
					translatedStatuses = append(translatedStatuses, models.PaymentRequestStatusPending.String())

				} else if strings.EqualFold(status, "Reviewed") {
					translatedStatuses = append(translatedStatuses,
						models.PaymentRequestStatusReviewed.String(),
						models.PaymentRequestStatusSentToGex.String(),
						models.PaymentRequestStatusReceivedByGex.String())
				} else if strings.EqualFold(status, "Paid") {
					translatedStatuses = append(translatedStatuses, models.PaymentRequestStatusPaid.String())
				}
			}
			query = query.Where("payment_requests.status in (?)", translatedStatuses)
		}
	}

}
