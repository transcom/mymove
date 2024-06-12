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

func (f *reportListFetcher) FetchReportList(appCtx appcontext.AppContext, params *services.FetchPaymentRequestListParams) (models.Reports, error) {
	paymentRequests, err := f.FetchPaymentRequestListForReports(appCtx, params)
	if err != nil {
		return nil, err
	}

	reports := f.BuildReportListFromPaymentRequests(appCtx, paymentRequests)
	return reports, nil
}

// Fetch Payment Requests for Navy service members and ignore TIO and GBLOC rules
func (f *reportListFetcher) FetchPaymentRequestListForReports(appCtx appcontext.AppContext, params *services.FetchPaymentRequestListParams) (*models.PaymentRequests, error) {
	paymentRequests := models.PaymentRequests{}

	approvedStatuses := []string{models.PaymentRequestStatusReviewed.String(), models.PaymentRequestStatusSentToGex.String(), models.PaymentRequestStatusReceivedByGex.String()}
	query := appCtx.DB().Q().EagerPreload(
		"MoveTaskOrder", "MoveTaskOrder.Orders", "MoveTaskOrder.Orders.ServiceMember",
	).
		InnerJoin("moves", "payment_requests.move_id = moves.id").
		InnerJoin("orders", "orders.id = moves.orders_id").
		InnerJoin("service_members", "orders.service_member_id = service_members.id").
		Where("moves.show = ?", models.BoolPointer(true)).
		Where("service_members.affiliation = ?", models.AffiliationNAVY).
		Where("payment_requests.status in (?)", approvedStatuses)

	err := query.GroupBy("payment_requests.id, service_members.id, moves.id").All(&paymentRequests)

	if err != nil {
		return nil, err
	}

	return &paymentRequests, nil
}

func (f *reportListFetcher) BuildReportListFromPaymentRequests(appCtx appcontext.AppContext, paymentRequests *models.PaymentRequests) models.Reports {
	var reports models.Reports

	if paymentRequests != nil {
		for _, paymentRequest := range *paymentRequests {
			var newReport models.Report
			newReport.ID = paymentRequest.ID
			newReport.FirstName = paymentRequest.MoveTaskOrder.Orders.ServiceMember.FirstName
			newReport.Edipi = paymentRequest.MoveTaskOrder.Orders.ServiceMember.Edipi
			newReport.Address1 = paymentRequest.MoveTaskOrder.Orders.ServiceMember.ResidentialAddress
			newReport.Address2 = paymentRequest.MoveTaskOrder.Orders.ServiceMember.BackupMailingAddress
			reports = append(reports, newReport)
		}
	} else {
		return nil
	}

	return reports
}
