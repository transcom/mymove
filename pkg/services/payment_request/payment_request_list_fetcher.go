package paymentrequest

import (
	"fmt"
	"strings"
	"time"

	"github.com/go-openapi/swag"

	"github.com/gobuffalo/pop/v5"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/services"
	officeuser "github.com/transcom/mymove/pkg/services/office_user"

	"github.com/transcom/mymove/pkg/models"
)

type paymentRequestListFetcher struct {
}

var parameters = map[string]string{
	"lastName":           "service_members.last_name",
	"dodID":              "service_members.edipi",
	"submittedAt":        "payment_requests.created_at",
	"branch":             "service_members.affiliation",
	"locator":            "moves.locator",
	"status":             "payment_requests.status",
	"age":                "payment_requests.created_at",
	"originDutyLocation": "duty_locations.name",
}

// NewPaymentRequestListFetcher returns a new payment request list fetcher
func NewPaymentRequestListFetcher() services.PaymentRequestListFetcher {
	return &paymentRequestListFetcher{}
}

// QueryOption defines the type for the functional arguments passed to ListOrders
type QueryOption func(*pop.Query)

// FetchPaymentRequestList returns a list of payment requests
func (f *paymentRequestListFetcher) FetchPaymentRequestList(appCtx appcontext.AppContext, officeUserID uuid.UUID, params *services.FetchPaymentRequestListParams) (*models.PaymentRequests, int, error) {

	gblocFetcher := officeuser.NewOfficeUserGblocFetcher()
	gbloc, gblocErr := gblocFetcher.FetchGblocForOfficeUser(appCtx, officeUserID)
	if gblocErr != nil {
		return &models.PaymentRequests{}, 0, gblocErr
	}

	paymentRequests := models.PaymentRequests{}
	query := appCtx.DB().Q().EagerPreload(
		"MoveTaskOrder.Orders.OriginDutyLocation.TransportationOffice",
		// See note further below about having to do this in a separate Load call due to a Pop issue.
		// "MoveTaskOrder.Orders.ServiceMember",
	).
		InnerJoin("moves", "payment_requests.move_id = moves.id").
		InnerJoin("orders", "orders.id = moves.orders_id").
		InnerJoin("service_members", "orders.service_member_id = service_members.id").
		InnerJoin("duty_locations", "duty_locations.id = orders.origin_duty_location_id").
		// Need to use left join because some duty locations do not have transportation offices
		LeftJoin("transportation_offices", "duty_locations.transportation_office_id = transportation_offices.id").
		// If a customer puts in an invalid ZIP for their pickup address, it won't show up in this view,
		// and we don't want it to get hidden from services counselors.
		LeftJoin("move_to_gbloc", "move_to_gbloc.move_id = moves.id").
		Where("moves.show = ?", swag.Bool(true))

	branchQuery := branchFilter(params.Branch)
	// If the user is associated with the USMC GBLOC we want to show them ALL the USMC moves, so let's override here.
	// We also only want to do the gbloc filtering thing if we aren't a USMC user, which we cover with the else.
	var gblocQuery QueryOption
	if gbloc == "USMC" {
		branchQuery = branchFilter(swag.String(string(models.AffiliationMARINES)))
		gblocQuery = nil
	} else {
		gblocQuery = shipmentGBLOCFilter(&gbloc)
	}
	locatorQuery := locatorFilter(params.Locator)
	dodIDQuery := dodIDFilter(params.DodID)
	lastNameQuery := lastNameFilter(params.LastName)
	dutyLocationQuery := destinationDutyLocationFilter(params.DestinationDutyStation)
	statusQuery := paymentRequestsStatusFilter(params.Status)
	submittedAtQuery := submittedAtFilter(params.SubmittedAt)
	originDutyLocationQuery := dutyLocationFilter(params.OriginDutyLocation)
	orderQuery := sortOrder(params.Sort, params.Order)

	options := [10]QueryOption{branchQuery, locatorQuery, dodIDQuery, lastNameQuery, dutyLocationQuery, statusQuery, originDutyLocationQuery, submittedAtQuery, gblocQuery, orderQuery}

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

	err := query.GroupBy("payment_requests.id, service_members.id, moves.id, duty_locations.id, duty_locations.name").Paginate(int(*params.Page), int(*params.PerPage)).All(&paymentRequests)
	if err != nil {
		return nil, 0, err
	}
	// Get the count
	count := query.Paginator.TotalEntriesSize

	for i := range paymentRequests {
		// There appears to be a bug in Pop for EagerPreload when you have two or more eager paths with 3+ levels
		// where the first 2 levels match.  For example:
		//   "MoveTaskOrder.Orders.OriginDutyLocation.TransportationOffice" and "MoveTaskOrder.Orders.ServiceMember"
		// In those cases, only the last relationship is loaded in the results.  So, we can only do one of the paths
		// in the EagerPreload above and request the second one explicitly with a separate Load call.
		//
		// Note that we also had a problem before with Eager as well.  Here's what we found with it:
		//   Due to a bug in pop (https://github.com/gobuffalo/pop/issues/578), we
		//   cannot eager load the address as "OriginDutyLocation.Address" because
		//   OriginDutyLocation is a pointer.
		loadErr := appCtx.DB().Load(&paymentRequests[i].MoveTaskOrder.Orders, "ServiceMember")
		if loadErr != nil {
			return nil, 0, err
		}
		loadErr = appCtx.DB().Load(&paymentRequests[i].MoveTaskOrder, "ShipmentGBLOC")
		if loadErr != nil {
			return nil, 0, err
		}
	}

	return &paymentRequests, count, nil
}

// FetchPaymentRequestListByMove returns a payment request by move locator id
func (f *paymentRequestListFetcher) FetchPaymentRequestListByMove(appCtx appcontext.AppContext, officeUserID uuid.UUID, locator string) (*models.PaymentRequests, error) {
	gblocFetcher := officeuser.NewOfficeUserGblocFetcher()
	gbloc, gblocErr := gblocFetcher.FetchGblocForOfficeUser(appCtx, officeUserID)
	if gblocErr != nil {
		return &models.PaymentRequests{}, gblocErr
	}

	paymentRequests := models.PaymentRequests{}

	// Replaced EagerPreload due to nullable fka on Contractor
	query := appCtx.DB().Q().Eager(
		"PaymentServiceItems.PaymentServiceItemParams.ServiceItemParamKey",
		"PaymentServiceItems.MTOServiceItem.ReService",
		"PaymentServiceItems.MTOServiceItem.MTOShipment",
		"MoveTaskOrder.Contractor",
		"MoveTaskOrder.Orders").
		InnerJoin("moves", "payment_requests.move_id = moves.id").
		InnerJoin("orders", "orders.id = moves.orders_id").
		InnerJoin("service_members", "orders.service_member_id = service_members.id").
		InnerJoin("contractors", "contractors.id = moves.contractor_id").
		InnerJoin("duty_locations", "duty_locations.id = orders.origin_duty_location_id").
		// Need to use left join because some duty locations do not have transportation offices
		LeftJoin("transportation_offices", "duty_locations.transportation_office_id = transportation_offices.id").
		// If a customer puts in an invalid ZIP for their pickup address, it won't show up in this view,
		// and we don't want it to get hidden from services counselors.
		LeftJoin("move_to_gbloc", "move_to_gbloc.move_id = moves.id").
		Where("moves.show = ?", swag.Bool(true))

	var branchQuery QueryOption
	// If the user is associated with the USMC GBLOC we want to show them ALL the USMC moves, so let's override here.
	// We also only want to do the gbloc filtering thing if we aren't a USMC user, which we cover with the else.
	var gblocQuery QueryOption
	if gbloc == "USMC" {
		branchQuery = branchFilter(swag.String(string(models.AffiliationMARINES)))
	} else {
		gblocQuery = shipmentGBLOCFilter(&gbloc)
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
		return nil, err
	}

	return &paymentRequests, nil
}

func orderName(query *pop.Query, order *string) *pop.Query {
	query.Order(fmt.Sprintf("service_members.last_name %s, service_members.first_name %s", *order, *order))
	return query
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
				query.Order(fmt.Sprintf("%s %s", sortTerm, reverseOrder(order)))
			} else {
				query.Order(fmt.Sprintf("%s %s", sortTerm, *order))
			}
		} else {
			query.Order("payment_requests.created_at asc")
		}
	}
}

func branchFilter(branch *string) QueryOption {
	return func(query *pop.Query) {
		// When no branch filter is selected we want to filter out Marine Corps payment requests
		if branch == nil {
			query.Where("service_members.affiliation != ?", models.AffiliationMARINES)
		} else {
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

func dutyLocationFilter(dutyLocation *string) QueryOption {
	return func(query *pop.Query) {
		if dutyLocation != nil {
			locationSearch := fmt.Sprintf("%s%%", *dutyLocation)
			query.Where("duty_locations.name ILIKE ?", locationSearch)
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
func destinationDutyLocationFilter(destinationDutyLocation *string) QueryOption {
	return func(query *pop.Query) {
		if destinationDutyLocation != nil {
			nameSearch := fmt.Sprintf("%s%%", *destinationDutyLocation)
			query.InnerJoin("duty_locations as destination_duty_location", "orders.new_duty_location = destination_duty_location.id").Where("destination_duty_location.name ILIKE ?", nameSearch)
		}
	}
}

func submittedAtFilter(submittedAt *time.Time) QueryOption {
	return func(query *pop.Query) {
		if submittedAt != nil {
			// Between is inclusive, so the end date is set to 1 milsecond prior to the next day
			submittedAtEnd := submittedAt.AddDate(0, 0, 1).Add(-1 * time.Millisecond)
			query.Where("payment_requests.created_at between ? and ?", submittedAt.Format(time.RFC3339), submittedAtEnd.Format(time.RFC3339))
		}
	}
}

func shipmentGBLOCFilter(gbloc *string) QueryOption {
	return func(query *pop.Query) {
		if gbloc != nil {
			query.Where("move_to_gbloc.gbloc = ?", *gbloc)
		}
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
				} else if strings.EqualFold(status, "Rejected") {
					translatedStatuses = append(translatedStatuses,
						models.PaymentRequestStatusReviewedAllRejected.String())
				} else if strings.EqualFold(status, "Paid") {
					translatedStatuses = append(translatedStatuses, models.PaymentRequestStatusPaid.String())
				}
			}
			query.Where("payment_requests.status in (?)", translatedStatuses)
		}
	}

}
