package paymentrequest

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/gobuffalo/pop/v6"
	"github.com/gofrs/uuid"
	"github.com/lib/pq"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/models/roles"
	"github.com/transcom/mymove/pkg/services"
	officeuser "github.com/transcom/mymove/pkg/services/office_user"
)

type paymentRequestListFetcher struct {
}

// NewPaymentRequestListFetcher returns a new payment request list fetcher
func NewPaymentRequestListFetcher() services.PaymentRequestListFetcher {
	return &paymentRequestListFetcher{}
}

// QueryOption defines the type for the functional arguments passed to ListOrders
type QueryOption func(*pop.Query)

type paymentRequestRow struct {
	PaymentRequest   json.RawMessage `db:"payment_request"`
	Move             json.RawMessage `db:"move"`
	Orders           json.RawMessage `db:"orders"`
	OriginToOffice   json.RawMessage `db:"origin_to_office"`
	TIOUser          json.RawMessage `db:"tio_user"`
	CounselingOffice json.RawMessage `db:"counseling_office"`
	TotalCount       int             `db:"total_count"`
}

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

	privileges, err := roles.FetchPrivilegesForUser(appCtx.DB(), appCtx.Session().UserID)
	if err != nil {
		appCtx.Logger().Error("Error retreiving user privileges", zap.Error(err))
	}
	hasSafetyPrivilege := privileges.HasPrivilege(roles.PrivilegeTypeSafety)

	var rows []paymentRequestRow
	err = appCtx.DB().
		RawQuery(
			`SELECT * FROM get_payment_request_queue($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16)`,
			gbloc,
			params.Branch,
			params.Locator,
			params.Edipi,
			params.Emplid,
			params.CustomerName,
			params.OriginDutyLocation,
			pq.Array(params.Status),
			params.SubmittedAt,
			params.AssignedTo,
			params.CounselingOffice,
			hasSafetyPrivilege,
			params.Page,
			params.PerPage,
			params.Sort,
			params.Order,
		).
		All(&rows)
	if err != nil {
		return nil, 0, err
	}

	moves := make(models.PaymentRequests, len(rows))
	var total int
	for i, r := range rows {
		total = r.TotalCount
		var pr models.PaymentRequest
		if err := json.Unmarshal(r.PaymentRequest, &pr); err != nil {
			return nil, 0, fmt.Errorf("error unmarshaling payment request JSON: %w", err)
		}
		if err := json.Unmarshal(r.Move, &pr.MoveTaskOrder); err != nil {
			return nil, 0, fmt.Errorf("error unmarshaling move task order JSON: %w", err)
		}
		if err := json.Unmarshal(r.Orders, &pr.MoveTaskOrder.Orders); err != nil {
			return nil, 0, fmt.Errorf("error unmarshaling orders JSON: %w", err)
		}
		if err := json.Unmarshal(r.OriginToOffice, &pr.MoveTaskOrder.Orders.OriginDutyLocation.TransportationOffice); err != nil {
			return nil, 0, fmt.Errorf("error unmarshaling origin duty location transportation office JSON: %w", err)
		}
		if err := json.Unmarshal(r.TIOUser, &pr.MoveTaskOrder.TIOPaymentRequestAssignedUser); err != nil {
			return nil, 0, fmt.Errorf("error unmarshaling assigned TIO user JSON: %w", err)
		}
		if err := json.Unmarshal(r.CounselingOffice, &pr.MoveTaskOrder.CounselingOffice); err != nil {
			return nil, 0, fmt.Errorf("error unmarshaling move's counseling office JSON: %w", err)
		}

		moves[i] = pr
	}

	return &moves, total, nil
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

// When a queue is sorted by a non-unique value (ex: status, branch) the order within each vlaue is inconsistent at different page sizes
// Adding a secondary sort ensures a consistent order within the primary sort column
func locatorFilter(locator *string) QueryOption {
	return func(query *pop.Query) {
		if locator != nil {
			query.Where("moves.locator = ?", strings.ToUpper(*locator))
		}
	}
}
