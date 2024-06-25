package report

import (
	"fmt"

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
	query := appCtx.DB().Q().EagerPreload(
		"MoveTaskOrder", "MoveTaskOrder.MTOShipments", "MoveTaskOrder.Orders", "MoveTaskOrder.Orders.SAC",
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
			paymentRequest := paymentRequest

			// pop doesn't support associations 3+ deep so there's data we need to load here in order to get the necessary data
			loadErr := appCtx.DB().Load(&paymentRequest.MoveTaskOrder.Orders, "ServiceMember")
			if loadErr != nil {
				fmt.Printf("Failed to load Entitlement table, %s", loadErr.Error())
			}

			loadErr = appCtx.DB().Load(&paymentRequest.MoveTaskOrder.Orders, "Entitlement")
			if loadErr != nil {
				fmt.Printf("Failed to load Entitlement table, %s", loadErr.Error())
			}

			// this doesn't do what I thought it would...
			// loadErr = appCtx.DB().Load(&paymentRequest.MoveTaskOrder.Orders, "DutyLocations")
			// if loadErr != nil {
			// 	fmt.Printf("Failed to load Entitlement table, %s", loadErr.Error())
			// }

			MoveTaskOrder := paymentRequest.MoveTaskOrder
			Orders := MoveTaskOrder.Orders

			progear := unit.Pound(0)
			sitTotal := unit.Pound(0)
			travelAdvance := unit.Cents(0)

			// sharing this for loop for all MTOShipment calculations
			for _, shipment := range MoveTaskOrder.MTOShipments {
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

			newReport.ID = paymentRequest.ID
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

			// Which shipment are we using for address
			// originAddress := models.Address(MoveTaskOrder.MTOShipments[0].DeliveryAddressUpdate.OriginalAddress)
			// newReport.OriginAddress = &originAddress
			// destinationAddress := models.Address(*MoveTaskOrder.MTOShipments[0].DestinationAddress)
			// newReport.DestinationAddress = &destinationAddress

			// I don't know what to preload here to get data in the query
			// newReport.OriginGBLOC = &Orders.OriginDutyLocation.TransportationOffice.Gbloc
			// newReport.DestinationGBLOC = &Orders.NewDutyLocation.TransportationOffice.Gbloc

			// still don't know what this is
			// newReport.DepCD =

			newReport.TravelAdvance = &travelAdvance

			if MoveTaskOrder.MTOShipments[0].PPMShipment != nil {
				// use departure date for PPM
				newReport.MoveDate = &MoveTaskOrder.MTOShipments[0].PPMShipment.ExpectedDepartureDate
			} else {
				// use requested pickup date for HHG
				newReport.MoveDate = MoveTaskOrder.MTOShipments[0].ActualPickupDate
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

			newReport.ShipmentNum = len(paymentRequest.MoveTaskOrder.MTOShipments)
			newReport.WeightEstimate = calculateTotalWeightEstimate(paymentRequest.MoveTaskOrder.MTOShipments)

			// newReport.TransmitCD
			newReport.DD2278IssueDate = MoveTaskOrder.ServiceCounselingCompletedAt
			// newReport.Miles
			newReport.WeightAuthorized = (*unit.Pound)(Orders.Entitlement.WeightAllowance())
			newReport.ShipmentId = paymentRequest.MoveTaskOrderID
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
			newReport.PickupDate = MoveTaskOrder.MTOShipments[0].ActualPickupDate
			// newReport.Rate =
			// newReport.PaidDate =
			// newReport.LinehaulTotal =
			// newReport.AccessorialTotal =
			// newReport.FuelTotal =
			// newReport.OtherTotal =
			// newReport.InvoicePaidAmt =
			newReport.TravelType = (*string)(Orders.OrdersTypeDetail)
			newReport.TravelClassCode = (*string)(&Orders.OrdersType)
			newReport.DeliveryDate = MoveTaskOrder.MTOShipments[0].ActualDeliveryDate
			// newReport.ActualOriginNetWeight =
			// newReport.DestinationReweighNetWeight = MoveTaskOrder.MTOShipments[0].
			newReport.CounseledDate = MoveTaskOrder.ServiceCounselingCompletedAt

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
