package report

import (
	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

type reportListFetcher struct {
}

func NewReportListFetcher() services.ReportListFetcher {
	return &reportListFetcher{}
}

// Fetch Moves with an approved Payment Request for Navy service members and ignore TIO and GBLOC rules
func (f *reportListFetcher) FetchMovesForReports(appCtx appcontext.AppContext, params *services.MoveFetcherParams) (models.Moves, error) {
	var moves models.Moves

	// Raw query may not be viable without a new model created
	// rawQueryTest, qerr := query.GetSQLQueryByName("report_builder")
	// if qerr != nil {
	// 	return moves, apperror.NewQueryError("AuditHistory", qerr, "")
	// }
	// queryy := appCtx.DB().RawQuery(rawQueryTest)
	// errr := queryy.All(&moves)

	// if errr != nil {
	// 	return nil, errr
	// }

	approvedStatuses := []string{models.PaymentRequestStatusReviewed.String(), models.PaymentRequestStatusSentToGex.String(), models.PaymentRequestStatusReceivedByGex.String()}
	query := appCtx.DB().EagerPreload(
		"PaymentRequests",
		"PaymentRequests.PaymentServiceItems",
		"PaymentRequests.PaymentServiceItems.MTOServiceItem",
		"PaymentRequests.PaymentServiceItems.MTOServiceItem.ReService",
		"PaymentRequests.PaymentServiceItems.MTOServiceItem.ReService.Name",
		"PaymentRequests.PaymentServiceItems.PaymentServiceItemParams.ServiceItemParamKey",
		"MTOServiceItems.ReService",
		"MTOServiceItems.Dimensions",
		"MTOServiceItems.ServiceRequestDocuments.ServiceRequestDocumentUploads",
		"MTOShipments.DestinationAddress",
		"MTOShipments.PickupAddress",
		"MTOShipments.SecondaryDeliveryAddress",
		"MTOShipments.SecondaryPickupAddress",
		"MTOShipments.MTOAgents",
		"MTOShipments.Reweigh",
		"Orders.ServiceMember",
		"Orders.ServiceMember.BackupContacts",
		"Orders.Entitlement",
		"Orders.Entitlement.WeightAllotted",
		"Orders.NewDutyLocation.Address",
		"Orders.NewDutyLocation.TransportationOffice.Gbloc",
		"Orders.OriginDutyLocation.Address",
		"Orders.TAC",
		"LockedByOfficeUser",
	).
		InnerJoin("payment_requests", "moves.id = payment_requests.move_id").
		InnerJoin("payment_service_items", "payment_service_items.payment_request_id = payment_requests.id").
		InnerJoin("mto_service_items", "payment_service_items.mto_service_item_id = mto_service_items.id").
		InnerJoin("re_services", "re_services.id = mto_service_items.re_service_id").
		InnerJoin("orders", "orders.id = moves.orders_id").
		InnerJoin("entitlements", "entitlements.id = orders.entitlement_id").
		InnerJoin("service_members", "orders.service_member_id = service_members.id").
		InnerJoin("transportation_accounting_codes", "orders.tac = transportation_accounting_codes.tac").
		InnerJoin("lines_of_accounting", "transportation_accounting_codes.loa_id = lines_of_accounting.id").
		Where("payment_requests.status in (?)", approvedStatuses).
		Where("service_members.affiliation = ?", models.AffiliationNAVY)

	err := query.All(&moves)

	if err != nil {
		return nil, err
	}

	return moves, nil
}
