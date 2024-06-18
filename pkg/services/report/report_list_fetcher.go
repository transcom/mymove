package report

import (
	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/unit"
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
	query := appCtx.DB().EagerPreload(
		"MoveTaskOrder", "MoveTaskOrder.MTOShipments", "MoveTaskOrder.Orders", "MoveTaskOrder.Orders.SAC", "MoveTaskOrder.Orders.Moves",
		"MoveTaskOrder.Orders.Entitlement", "MoveTaskOrder.Orders.ServiceMember", // "MoveTaskOrder.Orders.TransportationAccountingCode" // "MoveTaskOrder.Orders.DutyLocations",
	).
		InnerJoin("moves", "payment_requests.move_id = moves.id").
		InnerJoin("orders", "orders.id = moves.orders_id").
		InnerJoin("entitlements", "entitlements.id = orders.entitlement_id").
		InnerJoin("service_members", "orders.service_member_id = service_members.id").
		InnerJoin("transportation_accounting_codes", "orders.tac = transportation_accounting_codes.tac").
		// InnerJoin("lines_of_accounting", "transportation_accounting_codes.loa_id = lines_of_accounting.id").
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
			newReport.LastName = paymentRequest.MoveTaskOrder.Orders.ServiceMember.LastName
			newReport.MiddleInitial = paymentRequest.MoveTaskOrder.Orders.ServiceMember.MiddleName
			newReport.Affiliation = paymentRequest.MoveTaskOrder.Orders.ServiceMember.Affiliation
			newReport.PayGrade = paymentRequest.MoveTaskOrder.Orders.Grade
			newReport.Edipi = paymentRequest.MoveTaskOrder.Orders.ServiceMember.Edipi
			newReport.PhonePrimary = paymentRequest.MoveTaskOrder.Orders.ServiceMember.Telephone
			newReport.PhoneSecondary = paymentRequest.MoveTaskOrder.Orders.ServiceMember.SecondaryTelephone
			newReport.EmailPrimary = paymentRequest.MoveTaskOrder.Orders.ServiceMember.PersonalEmail
			// newReport.EmailSecondary =
			newReport.OrdersType = paymentRequest.MoveTaskOrder.Orders.OrdersType
			newReport.OrdersNumber = paymentRequest.MoveTaskOrder.Orders.OrdersNumber
			newReport.OrdersDate = &paymentRequest.MoveTaskOrder.Orders.IssueDate
			newReport.Address = paymentRequest.MoveTaskOrder.Orders.ServiceMember.ResidentialAddress

			// Which shipment are we using for address
			// newReport.OriginAddress = &paymentRequest.MoveTaskOrder.MTOShipments[0].DeliveryAddressUpdate.OriginalAddress
			// newReport.DestinationAddress = paymentRequest.MoveTaskOrder.MTOShipments[0].DestinationAddress

			// I don't know what to preload here to get data in the query
			// newReport.OriginGBLOC = &paymentRequest.MoveTaskOrder.Orders.OriginDutyLocation.TransportationOffice.Gbloc
			// newReport.DestinationGBLOC = &paymentRequest.MoveTaskOrder.Orders.NewDutyLocation.TransportationOffice.Gbloc

			// still don't know what this is
			// newReport.DepCD =

			// need to determine which shipment(s) have a ppm and get the travel advances and add them up
			// TravelAdvance = paymentRequest.MoveTaskOrder.MTOShipments[0].

			// Get the MoveDate from HHG/PPM
			// newReport.MoveDate = paymentRequest.MoveTaskOrder.

			// newReport.TAC = paymentRequest.MoveTaskOrder.Orders.TAC

			// "get with navy on what they need for fiscal year"
			// newReport.FiscalYear = paymentRequest.

			// need to figure out how to preload the TAC and LOA tables to get some of these. Some are still unknown where there location lives
			// newReport.Appro
			// newReport.Subhead
			// newReport.ObjClass
			// newReport.BCN
			// newReport.SubAllotCD
			// newReport.AAA
			// newReport.TypeCD
			// newReport.PAA
			// newReport.CostCD
			// newReport.DDCD

			newReport.ShipmentNum = len(paymentRequest.MoveTaskOrder.MTOShipments)
			newReport.WeightEstimate = calculateTotalWeightEstimate(paymentRequest.MoveTaskOrder.MTOShipments)

			// newReport.TransmitCD
			newReport.DD2278IssueDate = paymentRequest.MoveTaskOrder.ServiceCounselingCompletedAt
			// newReport.Miles
			newReport.WeightAuthorized = (*unit.Pound)(paymentRequest.MoveTaskOrder.Orders.Entitlement.DBAuthorizedWeight) // entitlement table isn't working
			newReport.ShipmentId = paymentRequest.MoveTaskOrderID
			// newReport.SCAC
			newReport.OrderNumber = paymentRequest.MoveTaskOrder.Orders.SAC
			// newReport.LOA
			// newReport.ShipmentType =
			// newReport.EntitlementWeight = (*unit.Pound)(paymentRequest.MoveTaskOrder.Orders.Entitlement // entitlement table isn't working
			// newReport.NetWeight =

			reports = append(reports, newReport)
		}
	} else {
		return nil
	}

	return reports
}

func calculateTotalWeightEstimate(shipments models.MTOShipments) *unit.Pound {
	var weightEstimate unit.Pound
	for _, shipment := range shipments {
		if shipment.PPMShipment != nil {
			weightEstimate += *shipment.PPMShipment.EstimatedWeight
		}

		if shipment.PrimeEstimatedWeight != nil {
			weightEstimate += *shipment.PrimeEstimatedWeight
		}
	}

	return &weightEstimate
}
