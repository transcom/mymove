package paymentrequest

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/gobuffalo/pop/v6"
	"github.com/gofrs/uuid"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	officeuser "github.com/transcom/mymove/pkg/services/office_user"
)

type paymentRequestListFetcher struct {
}

var parameters = map[string]string{
	"customerName":       "(service_members.last_name || ' ' || service_members.first_name)",
	"dodID":              "service_members.edipi",
	"emplid":             "service_members.emplid",
	"submittedAt":        "payment_requests.created_at",
	"branch":             "service_members.affiliation",
	"locator":            "moves.locator",
	"status":             "payment_requests.status",
	"age":                "payment_requests.created_at",
	"originDutyLocation": "duty_locations.name",
	"assignedTo":         "assigned_user.last_name,assigned_user.first_name",
}

// NewPaymentRequestListFetcher returns a new payment request list fetcher
func NewPaymentRequestListFetcher() services.PaymentRequestListFetcher {
	return &paymentRequestListFetcher{}
}

// QueryOption defines the type for the functional arguments passed to ListOrders
type QueryOption func(*pop.Query)

// FetchPaymentRequestList returns a list of payment requests
func (f *paymentRequestListFetcher) FetchPaymentRequestList(appCtx appcontext.AppContext, officeUserID uuid.UUID, params *services.FetchPaymentRequestListParams) (*models.PaymentRequests, int, error) {
	var gbloc string
	if params.ViewAsGBLOC != nil {
		gbloc = *params.ViewAsGBLOC
	} else {
		var gblocErr error
		gblocFetcher := officeuser.NewOfficeUserGblocFetcher()
		gbloc, gblocErr = gblocFetcher.FetchGblocForOfficeUser(appCtx, officeUserID)
		if gblocErr != nil {
			return &models.PaymentRequests{}, 0, gblocErr
		}
	}

	privileges, err := models.FetchPrivilegesForUser(appCtx.DB(), appCtx.Session().UserID)
	if err != nil {
		appCtx.Logger().Error("Error retreiving user privileges", zap.Error(err))
	}

	paymentRequests := models.PaymentRequests{}
	query := appCtx.DB().Q().EagerPreload(
		"MoveTaskOrder.Orders.OriginDutyLocation.TransportationOffice",
		"MoveTaskOrder.Orders.OriginDutyLocation.Address",
		"MoveTaskOrder.TIOAssignedUser",
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
		LeftJoin("office_users as assigned_user", "moves.tio_assigned_id = assigned_user.id").
		Where("moves.show = ?", models.BoolPointer(true))

	if !privileges.HasPrivilege(models.PrivilegeTypeSafety) {
		query.Where("orders.orders_type != (?)", "SAFETY")
	}

	branchQuery := branchFilter(params.Branch)
	// If the user is associated with the USMC GBLOC we want to show them ALL the USMC moves, so let's override here.
	// We also only want to do the gbloc filtering thing if we aren't a USMC user, which we cover with the else.
	var gblocQuery QueryOption
	if gbloc == "USMC" {
		branchQuery = branchFilter(models.StringPointer(string(models.AffiliationMARINES)))
		gblocQuery = nil
	} else {
		gblocQuery = shipmentGBLOCFilter(&gbloc)
	}
	locatorQuery := locatorFilter(params.Locator)
	dodIDQuery := dodIDFilter(params.DodID)
	emplidQuery := emplidFilter(params.Emplid)
	customerNameQuery := customerNameFilter(params.CustomerName)
	dutyLocationQuery := destinationDutyLocationFilter(params.DestinationDutyLocation)
	statusQuery := paymentRequestsStatusFilter(params.Status)
	submittedAtQuery := submittedAtFilter(params.SubmittedAt)
	originDutyLocationQuery := dutyLocationFilter(params.OriginDutyLocation)
	orderQuery := sortOrder(params.Sort, params.Order)
	tioAssignedUserQuery := tioAssignedUserFilter(params.TIOAssignedUser)

	options := [12]QueryOption{branchQuery, locatorQuery, dodIDQuery, customerNameQuery, dutyLocationQuery, statusQuery, originDutyLocationQuery, submittedAtQuery, gblocQuery, orderQuery, emplidQuery, tioAssignedUserQuery}

	for _, option := range options {
		if option != nil {
			option(query)
		}
	}

	if params.Page == nil {
		params.Page = models.Int64Pointer(1)
	}

	if params.PerPage == nil {
		params.PerPage = models.Int64Pointer(20)
	}

	err = query.GroupBy("payment_requests.id, service_members.id, moves.id, duty_locations.id, duty_locations.name, assigned_user.last_name, assigned_user.first_name").Paginate(int(*params.Page), int(*params.PerPage)).All(&paymentRequests)
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
func (f *paymentRequestListFetcher) FetchPaymentRequestListByMove(appCtx appcontext.AppContext, locator string) (*models.PaymentRequests, error) {
	paymentRequests := models.PaymentRequests{}

	// Replaced EagerPreload due to nullable fka on Contractor
	query := appCtx.DB().Q().Eager(
		"PaymentServiceItems.PaymentServiceItemParams.ServiceItemParamKey",
		"PaymentServiceItems.MTOServiceItem.ReService",
		"PaymentServiceItems.MTOServiceItem.MTOShipment",
		"ProofOfServiceDocs.PrimeUploads.Upload",
		"MoveTaskOrder.Contractor",
		"MoveTaskOrder.Orders.ServiceMember",
		"MoveTaskOrder.Orders.NewDutyLocation.Address",
		"MoveTaskOrder.LockedByOfficeUser").
		InnerJoin("moves", "payment_requests.move_id = moves.id").
		InnerJoin("orders", "orders.id = moves.orders_id").
		InnerJoin("service_members", "orders.service_member_id = service_members.id").
		InnerJoin("contractors", "contractors.id = moves.contractor_id").
		InnerJoin("duty_locations", "duty_locations.id = orders.origin_duty_location_id").
		// Need to use left join because some duty locations do not have transportation offices
		LeftJoin("transportation_offices", "duty_locations.transportation_office_id = transportation_offices.id").
		LeftJoin("office_users", "office_users.id = moves.locked_by").
		// If a customer puts in an invalid ZIP for their pickup address, it won't show up in this view,
		// and we don't want it to get hidden from services counselors.
		Where("moves.show = ?", models.BoolPointer(true))

	locatorQuery := locatorFilter(&locator)

	options := [1]QueryOption{locatorQuery}

	for _, option := range options {
		if option != nil {
			option(query)
		}
	}

	err := query.All(&paymentRequests)
	if err != nil {
		return nil, err
	}

	for i := range paymentRequests {
		for j := range paymentRequests[i].ProofOfServiceDocs {
			paymentRequests[i].ProofOfServiceDocs[j].PrimeUploads = paymentRequests[i].ProofOfServiceDocs[j].PrimeUploads.FilterDeleted()
		}

		mostRecentEdiErrorForPaymentRequest, errFetchingEdiError := fetchEDIErrorsForPaymentRequest(appCtx, &paymentRequests[i])
		if errFetchingEdiError != nil {
			return nil, errFetchingEdiError
		}

		// We process TPPS Paid Invoice Reports to get payment information for each payment service item
		// As well as the total amount paid for the overall payment request, and the date it was paid
		// This report tells us how much TPPS paid HS, then we store and display it
		tppsReportEntryList, errFetchingTPPSInformation := fetchTPPSPaidInvoiceReportDataPaymentRequest(appCtx, &paymentRequests[i])
		if errFetchingTPPSInformation != nil {
			return nil, errFetchingTPPSInformation
		}

		paymentRequests[i].EdiErrors = append(paymentRequests[i].EdiErrors, mostRecentEdiErrorForPaymentRequest)
		paymentRequests[i].TPPSPaidInvoiceReports = tppsReportEntryList

	}

	return &paymentRequests, nil
}

// fetchEDIErrorsForPaymentRequest returns the edi_error with the most recent created_at date for a payment request
func fetchEDIErrorsForPaymentRequest(appCtx appcontext.AppContext, pr *models.PaymentRequest) (models.EdiError, error) {

	// find any associated errors in the edi_errors table from processing the EDI 858, 824, or 997
	var ediError []models.EdiError
	ediErrorInfo := models.EdiError{}

	// regardless of PR status, find any associated edi_errors
	// 997s could have edi_errors logged but not have a status of EDI_ERROR
	err := appCtx.DB().Q().
		Where("edi_errors.payment_request_id = $1", pr.ID).
		Order("created_at DESC").
		All(&ediError)

	if err != nil {
		return ediErrorInfo, err
	} else if len(ediError) == 0 {
		return ediErrorInfo, nil
	}
	if len(ediError) > 0 {
		// since we ordered by created_at desc, the first result will be the most recent error we want to grab
		ediErrorInfo = ediError[0]
	}

	return ediErrorInfo, nil
}

// fetchTPPSPaidInvoiceReportDataPaymentRequest returns entries in the tpps_paid_invoice_reports
// for a payment request by matching the payment request number to the TPPS invoice number
func fetchTPPSPaidInvoiceReportDataPaymentRequest(appCtx appcontext.AppContext, pr *models.PaymentRequest) (models.TPPSPaidInvoiceReportEntrys, error) {

	var tppsPaidInvoiceReport []models.TPPSPaidInvoiceReportEntry
	tppsPaidInformation := models.TPPSPaidInvoiceReportEntrys{}

	err := appCtx.DB().Q().
		Where("tpps_paid_invoice_reports.invoice_number = $1", pr.PaymentRequestNumber).
		All(&tppsPaidInvoiceReport)

	if err != nil {
		return tppsPaidInformation, err
	} else if len(tppsPaidInvoiceReport) == 0 {
		return tppsPaidInformation, nil
	}
	if len(tppsPaidInvoiceReport) > 0 {
		tppsPaidInformation = tppsPaidInvoiceReport
		return tppsPaidInformation, nil
	}

	return tppsPaidInformation, nil
}

func orderName(query *pop.Query, order *string) *pop.Query {
	query.Order(fmt.Sprintf("service_members.last_name %s, service_members.first_name %s", *order, *order))
	return query
}

func orderAssignedName(query *pop.Query, order *string) *pop.Query {
	query.Order(fmt.Sprintf("assigned_user.last_name %s, assigned_user.first_name %s", *order, *order))
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
			} else if *sort == "assignedTo" {
				orderAssignedName(query, order)
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
				if strings.EqualFold(status, "Pending") || strings.EqualFold(status, "Payment Requested") {
					translatedStatuses = append(translatedStatuses, models.PaymentRequestStatusPending.String())

				} else if strings.EqualFold(status, "Reviewed") {
					translatedStatuses = append(translatedStatuses,
						models.PaymentRequestStatusReviewed.String(),
						models.PaymentRequestStatusSentToGex.String(),
						models.PaymentRequestStatusTppsReceived.String())
				} else if strings.EqualFold(status, "Rejected") || strings.EqualFold(status, "REVIEWED_AND_ALL_SERVICE_ITEMS_REJECTED") {
					translatedStatuses = append(translatedStatuses,
						models.PaymentRequestStatusReviewedAllRejected.String())
				} else if strings.EqualFold(status, "Paid") {
					translatedStatuses = append(translatedStatuses, models.PaymentRequestStatusPaid.String())
				} else if strings.EqualFold(status, "Deprecated") {
					translatedStatuses = append(translatedStatuses, models.PaymentRequestStatusDeprecated.String())
				} else if strings.EqualFold(status, "Error") || strings.EqualFold(status, "EDI_ERROR") {
					translatedStatuses = append(translatedStatuses, models.PaymentRequestStatusEDIError.String())
				}
			}
			query.Where("payment_requests.status in (?)", translatedStatuses)
		}
	}

}

func tioAssignedUserFilter(tioAssigned *string) QueryOption {
	return func(query *pop.Query) {
		if tioAssigned != nil {
			nameSearch := fmt.Sprintf("%s%%", *tioAssigned)
			query.Where("assigned_user.last_name ILIKE ?", nameSearch)
		}
	}
}
