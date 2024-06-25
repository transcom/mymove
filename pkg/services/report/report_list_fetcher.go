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

// Fetch Moves with an approved Payment Request for Navy service members and ignore TIO and GBLOC rules
func (f *reportListFetcher) FetchMovesForReports(appCtx appcontext.AppContext, params *services.MoveFetcherParams) (models.Moves, error) {
	var moves models.Moves

	approvedStatuses := []string{models.PaymentRequestStatusReviewed.String(), models.PaymentRequestStatusSentToGex.String(), models.PaymentRequestStatusReceivedByGex.String()}
	query := appCtx.DB().EagerPreload(
		"PaymentRequests",
		"PaymentRequests.PaymentServiceItems",
		"PaymentRequests.PaymentServiceItems.PaymentServiceItemParams.ServiceItemParamKey",
		"MTOServiceItems.ReService",
		"MTOServiceItems.Dimensions",
		"MTOServiceItems.ServiceRequestDocuments.ServiceRequestDocumentUploads",
		"MTOShipments.DestinationAddress",
		"MTOShipments.PickupAddress",
		"MTOShipments.SecondaryDeliveryAddress",
		"MTOShipments.SecondaryPickupAddress",
		"MTOShipments.MTOAgents",
		"Orders.ServiceMember",
		"Orders.ServiceMember.BackupContacts",
		"Orders.Entitlement",
		"Orders.Entitlement.WeightAllotted",
		"Orders.NewDutyLocation.Address",
		"Orders.OriginDutyLocation.Address",
		"LockedByOfficeUser",
		"CloseoutOfficeID",
	).
		InnerJoin("payment_requests", "moves.id = payment_requests.move_id").
		InnerJoin("orders", "orders.id = moves.orders_id").
		InnerJoin("entitlements", "entitlements.id = orders.entitlement_id").
		InnerJoin("service_members", "orders.service_member_id = service_members.id").
		LeftJoin("transportation_offices", "moves.closeout_office_id = transportation_offices.id").
		Where("payment_requests.status in (?)", approvedStatuses).
		Where("service_members.affiliation = ?", models.AffiliationNAVY)

	err := query.All(&moves)
	if err != nil {
		return nil, err
	}

	return moves, nil
}

func (f *reportListFetcher) BuildReportListFromMoves(appCtx appcontext.AppContext, moves models.Moves) models.Reports {
	var reports models.Reports

	for _, move := range moves {
		var newReport models.Report
		move := move

		Orders := move.Orders

		progear := unit.Pound(0)
		sitTotal := unit.Pound(0)
		travelAdvance := unit.Cents(0)

		// sharing this for loop for all MTOShipment calculations
		for _, shipment := range move.MTOShipments {
			// calculate total progear for entire move
			if shipment.PPMShipment != nil {
				shipmentTotalProgear := shipment.PPMShipment.ProGearWeight.Float64() + shipment.PPMShipment.SpouseProGearWeight.Float64()
				progear += unit.Pound(shipmentTotalProgear)

				// need to determine which shipment(s) have a ppm and get the travel advances and add them up
				if shipment.PPMShipment.AdvanceAmountReceived != nil {
					travelAdvance += *shipment.PPMShipment.AdvanceAmountReceived
				}

				// add SIT estimated weights
				if *shipment.PPMShipment.SITExpected {
					sitTotal += *shipment.PPMShipment.SITEstimatedWeight

					// SIT Fields
					newReport.SitInDate = shipment.PPMShipment.SITEstimatedEntryDate
					newReport.SitOutDate = shipment.PPMShipment.SITEstimatedDepartureDate
					// newreport.SitType = // Example data is destination.. ??
				}
			}
		}

		newReport.ID = move.ID
		newReport.FirstName = Orders.ServiceMember.FirstName
		newReport.LastName = Orders.ServiceMember.LastName
		newReport.MiddleInitial = Orders.ServiceMember.MiddleName
		newReport.Affiliation = Orders.ServiceMember.Affiliation
		newReport.PayGrade = Orders.Grade
		newReport.Edipi = Orders.ServiceMember.Edipi
		newReport.PhonePrimary = Orders.ServiceMember.Telephone
		newReport.PhoneSecondary = Orders.ServiceMember.SecondaryTelephone
		newReport.EmailPrimary = Orders.ServiceMember.PersonalEmail
		// newReport.EmailSecondary = // Are we using Backup contact email for secondary email?
		newReport.OrdersType = Orders.OrdersType
		newReport.OrdersNumber = Orders.OrdersNumber
		newReport.OrdersDate = &Orders.IssueDate
		newReport.Address = Orders.ServiceMember.ResidentialAddress

		newReport.OriginAddress = move.MTOShipments[0].PickupAddress
		newReport.DestinationAddress = move.MTOShipments[0].DestinationAddress

		// I don't know what to preload here to get data in the query
		// newReport.OriginGBLOC = &Orders.OriginDutyLocation.TransportationOffice.Gbloc
		// newReport.DestinationGBLOC = &Orders.NewDutyLocation.TransportationOffice.Gbloc

		// still don't know what this is
		// newReport.DepCD =

		newReport.TravelAdvance = &travelAdvance

		if move.MTOShipments[0].PPMShipment != nil {
			// use departure date for PPM
			newReport.MoveDate = &move.MTOShipments[0].PPMShipment.ExpectedDepartureDate
		} else {
			// use requested pickup date for HHG
			newReport.MoveDate = move.MTOShipments[0].ActualPickupDate
		}

		// newReport.TAC = MoveTaskOrder.Orders.TAC

		// "get with navy on what they need for fiscal year"
		// newReport.FiscalYear = paymentRequest.

		// need to figure out how to preload the TAC and LOA tables to get some of these. Some are still unknown where there location lives
		// LOA/TAC Fields
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

		newReport.ShipmentNum = len(move.MTOShipments)
		newReport.WeightEstimate = calculateTotalWeightEstimate(move.MTOShipments)

		// newReport.TransmitCD
		newReport.DD2278IssueDate = move.ServiceCounselingCompletedAt
		// newReport.Miles
		newReport.WeightAuthorized = (*unit.Pound)(Orders.Entitlement.WeightAllowance())
		newReport.ShipmentId = move.ID
		// newReport.SCAC =
		if Orders.SAC != nil {
			newReport.OrderNumber = Orders.SAC
		} else {
			emptySAC := ""
			newReport.OrderNumber = &emptySAC
		}
		// newReport.LOA
		// newReport.ShipmentType =
		newReport.EntitlementWeight = (*unit.Pound)(Orders.Entitlement.DBAuthorizedWeight)
		// newReport.NetWeight =
		newReport.PBPAndE = &progear
		newReport.PickupDate = move.MTOShipments[0].ActualPickupDate
		// newReport.Rate =
		// newReport.PaidDate =
		// newReport.LinehaulTotal =
		// newReport.AccessorialTotal =
		// newReport.FuelTotal =
		// newReport.OtherTotal =
		// newReport.InvoicePaidAmt =
		newReport.TravelType = (*string)(Orders.OrdersTypeDetail)
		newReport.TravelClassCode = (*string)(&Orders.OrdersType)
		newReport.DeliveryDate = move.MTOShipments[0].ActualDeliveryDate
		// newReport.ActualOriginNetWeight =
		// newReport.DestinationReweighNetWeight = MoveTaskOrder.MTOShipments[0].
		newReport.CounseledDate = move.ServiceCounselingCompletedAt

		reports = append(reports, newReport)
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
