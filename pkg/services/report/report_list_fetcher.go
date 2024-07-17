package report

import (
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/gen/pptasmessages"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/unit"
)

type reportListFetcher struct {
	estimator services.PPMEstimator
}

func NewReportListFetcher(estimator services.PPMEstimator) services.ReportListFetcher {
	return &reportListFetcher{
		estimator: estimator,
	}
}

func (f *reportListFetcher) BuildReportFromMoves(appCtx appcontext.AppContext, params *services.MoveTaskOrderFetcherParams) (models.Reports, error) {
	var fullreport models.Reports
	moves, err := FetchMovesForReports(appCtx, params)

	if err != nil {
		return nil, err
	}

	for _, move := range moves {
		var report models.Report

		orders := move.Orders
		var paymentRequests []models.PaymentRequest
		for _, pr := range move.PaymentRequests {
			if pr.Status == models.PaymentRequestStatusReviewed || pr.Status == models.PaymentRequestStatusReceivedByGex || pr.Status == models.PaymentRequestStatusSentToGex {
				paymentRequests = append(paymentRequests, pr)
			}
		}

		tac := FetchTACForMmove(appCtx, orders)

		var middleInitial string
		if *orders.ServiceMember.MiddleName != "" {
			middleInitial = string([]rune(*orders.ServiceMember.MiddleName)[0])
		}

		progear := unit.Pound(0)
		sitTotal := unit.Pound(0)
		originActualWeight := unit.Pound(0)
		travelAdvance := unit.Cents(0)
		scac := "HSFR"
		transmitCd := "T"

		var moveDate *time.Time
		if move.MTOShipments[0].PPMShipment != nil {
			moveDate = &move.MTOShipments[0].PPMShipment.ExpectedDepartureDate
		} else if move.MTOShipments[0].ActualPickupDate != nil {
			moveDate = move.MTOShipments[0].ActualPickupDate
		}

		if moveDate != nil {
			report.DeliveryDate = moveDate
		}

		if !reflect.ValueOf(orders.OrdersTypeDetail).IsNil() {
			report.ShipmentType = (*string)(orders.OrdersTypeDetail)
			report.TravelType = (*string)(orders.OrdersTypeDetail)
		}

		if !reflect.ValueOf(move.MTOShipments[0].ActualPickupDate).IsNil() {
			report.PickupDate = move.MTOShipments[0].ActualPickupDate
		}

		if orders.Grade != nil && orders.Entitlement != nil {
			orders.Entitlement.SetWeightAllotment(string(*orders.Grade))
		}

		weightAllotment := orders.Entitlement.WeightAllotment()

		var totalWeight unit.Pound
		if orders.Entitlement.DBAuthorizedWeight != nil {
			if orders.Entitlement.DependentsAuthorized != nil {
				totalWeight = unit.Pound(weightAllotment.TotalWeightSelfPlusDependents)

				report.WeightAuthorized = (*unit.Pound)(orders.Entitlement.DBAuthorizedWeight)
			} else {
				totalWeight = unit.Pound(weightAllotment.TotalWeightSelf)
			}
		}

		report.EntitlementWeight = &totalWeight

		var longLoa string
		if len(tac) > 0 {
			longLoa = buildFullLineOfAccountingString(tac[0].LineOfAccounting)

			report.LOA = &longLoa
			report.FiscalYear = tac[0].TacFyTxt
			report.Appro = tac[0].LineOfAccounting.LoaBafID
			report.Subhead = tac[0].LineOfAccounting.LoaObjClsID
			report.ObjClass = tac[0].LineOfAccounting.LoaAlltSnID
			report.BCN = tac[0].LineOfAccounting.LoaSbaltmtRcpntID
			report.SubAllotCD = tac[0].LineOfAccounting.LoaInstlAcntgActID
			report.AAA = tac[0].LineOfAccounting.LoaTrnsnID
			report.TypeCD = tac[0].LineOfAccounting.LoaJbOrdNm
			report.PAA = tac[0].LineOfAccounting.LoaDocID
			report.CostCD = tac[0].LineOfAccounting.LoaPgmElmntID
			report.DDCD = tac[0].LineOfAccounting.LoaDptID
		}

		var linehaulTotal float64
		var managementTotal float64
		var fuelPrice float64
		var domesticOriginTotal float64
		var domesticDestTotal float64
		var domesticPacking float64
		var domesticUnpacking float64
		var domesticCrating float64
		var domesticUncrating float64
		var counselingTotal float64
		var sitPickuptotal float64
		var sitOriginFuelSurcharge float64
		var sitOriginShuttle float64
		var sitOriginAddlDays float64
		var sitOriginFirstDay float64
		var sitDeliveryTotal float64
		var sitDestFuelSurcharge float64
		var sitDestShuttle float64
		var sitDestAddlDays float64
		var sitDestFirstDay float64

		var allCrates []*pptasmessages.Crate

		// this adds up all the different payment service items across all payment requests for a move
		for _, pr := range paymentRequests {
			for _, serviceItem := range pr.PaymentServiceItems {
				var mtoServiceItem models.MTOServiceItem
				msiErr := appCtx.DB().Q().EagerPreload("ReService", "Dimensions").
					InnerJoin("re_services", "re_services.id = mto_service_items.re_service_id").
					Where("mto_service_items.id = ?", serviceItem.MTOServiceItemID).
					First(&mtoServiceItem)
				if msiErr != nil {
					return nil, apperror.NewQueryError("failed to query service items", msiErr, ".")
				}

				totalPrice := serviceItem.PriceCents.Float64()
				sitType := ""

				switch mtoServiceItem.ReService.Name {
				case "Domestic linehaul":
				case "Domestic shorthaul":
					linehaulTotal += totalPrice
				case "Move management":
					managementTotal += totalPrice
				case "Fuel surcharge":
					fuelPrice += totalPrice
				case "Domestic origin price":
					domesticOriginTotal += totalPrice
				case "Domestic destination price":
					domesticDestTotal += totalPrice
				case "Domestic packing":
					domesticPacking += totalPrice
				case "Domestic unpacking":
					domesticUnpacking += totalPrice
				case "Domestic uncrating":
					domesticUncrating += totalPrice
				// case "Domestic crating - standalone":
				case "Domestic crating":
					crate := buildServiceItemCrate(mtoServiceItem)
					allCrates = append(allCrates, &crate)
					domesticCrating += totalPrice
				case "Domestic origin SIT pickup":
					sitPickuptotal += totalPrice
				case "Domestic origin SIT fuel surcharge":
					sitOriginFuelSurcharge += totalPrice
				case "Domestic origin shuttle service":
					sitOriginShuttle += totalPrice
				case "Domestic origin add'l SIT":
					sitOriginAddlDays += totalPrice
				case "Domestic origin 1st day SIT":
					if sitType == "" {
						sitType = "Origin"
						report.SitType = &sitType
					}
					sitOriginFirstDay += totalPrice
				case "Domestic destination SIT fuel surcharge":
					sitDestFuelSurcharge += totalPrice
				case "Domestic destination SIT delivery":
					sitDeliveryTotal += totalPrice
				case "Domestic destination shuttle service":
					sitDestShuttle += totalPrice
				case "Domestic destination add'l SIT":
					sitDestAddlDays += totalPrice
				case "Domestic destination 1st day SIT":
					if sitType == "Origin" || sitType == "" {
						sitType := "Destination"
						report.SitType = &sitType
					}
					sitDestFirstDay += totalPrice
				case "Counseling":
					counselingTotal += totalPrice
				}
			}
		}

		shuttleTotal := sitOriginShuttle + sitDestShuttle
		report.LinehaulTotal = &linehaulTotal
		report.LinehaulFuelTotal = &fuelPrice
		report.OriginPrice = &domesticOriginTotal
		report.DestinationPrice = &domesticDestTotal
		report.PackingPrice = &domesticPacking
		report.UnpackingPrice = &domesticUnpacking
		report.CratingTotal = &domesticCrating
		report.UncratingTotal = &domesticUncrating
		report.ShuttleTotal = &shuttleTotal
		report.MoveManagementFeeTotal = &managementTotal
		report.CounselingFeeTotal = &counselingTotal
		report.CratingDimensions = allCrates

		// calculate total invoice cost
		invoicePaidAmt := shuttleTotal + linehaulTotal + fuelPrice + domesticOriginTotal + domesticDestTotal + domesticPacking + domesticUnpacking +
			sitOriginFirstDay + sitOriginAddlDays + sitDestFirstDay + sitDestAddlDays + sitPickuptotal + sitDeliveryTotal + sitOriginFuelSurcharge +
			sitDestFuelSurcharge + domesticCrating + domesticUncrating
		report.InvoicePaidAmt = &invoicePaidAmt

		var ppmLinehaul float64
		var ppmFuel float64
		var ppmOriginPrice float64
		var ppmDestPrice float64
		var ppmPacking float64
		var ppmUnpacking float64

		// sharing this for loop for all MTOShipment calculations
		for _, shipment := range move.MTOShipments {
			// calculate total progear for entire move
			if shipment.PPMShipment != nil {
				var shipmentTotalProgear float64
				if !reflect.ValueOf(shipment.PPMShipment.ProGearWeight).IsNil() {
					shipmentTotalProgear += shipment.PPMShipment.ProGearWeight.Float64()
				}

				if !reflect.ValueOf(shipment.PPMShipment.SpouseProGearWeight).IsNil() {
					shipmentTotalProgear += shipment.PPMShipment.SpouseProGearWeight.Float64()
				}

				progear += unit.Pound(shipmentTotalProgear)

				// need to determine which shipment(s) have a ppm and get the travel advances and add them up
				if shipment.PPMShipment.AdvanceAmountReceived != nil {
					travelAdvance += *shipment.PPMShipment.AdvanceAmountReceived
				}

				// add SIT estimated weights
				if *shipment.PPMShipment.SITExpected {
					sitTotal += *shipment.PPMShipment.SITEstimatedWeight

					// SIT Fields
					report.SitInDate = shipment.PPMShipment.SITEstimatedEntryDate
					report.SitOutDate = shipment.PPMShipment.SITEstimatedDepartureDate
				}

				// query the ppmshipment for all it's child needs for the price breakdown
				var ppmShipment models.PPMShipment
				ppmQ := appCtx.DB().Q().EagerPreload("PickupAddress", "DestinationAddress", "WeightTickets", "Shipment").
					InnerJoin("mto_shipments", "mto_shipments.id = ppm_shipments.shipment_id").
					Where("ppm_shipments.id = ?", shipment.PPMShipment.ID).
					First(&ppmShipment)
				if ppmQ != nil {
					return nil, apperror.NewQueryError("failed to query ppm ", ppmQ, ".")
				}

				// do the ppm cost breakdown here
				linehaul, fuel, origin, dest, packing, unpacking, err := f.estimator.PriceBreakdown(appCtx, &ppmShipment)
				if err != nil {
					return nil, apperror.NewUnprocessableEntityError("ppm price breakdown")
				}

				ppmLinehaul += linehaul.Float64()
				ppmFuel += fuel.Float64()
				ppmOriginPrice += origin.Float64()
				ppmDestPrice += dest.Float64()
				ppmPacking += packing.Float64()
				ppmUnpacking += unpacking.Float64()

			}

			if shipment.PrimeActualWeight != nil {
				originActualWeight += *shipment.PrimeActualWeight
			}
		}

		if report.SitInDate != nil || report.SitOutDate != nil {
			report.SITOriginFirstDayTotal = &sitOriginFirstDay
			report.SITOriginAddlDaysTotal = &sitOriginAddlDays
			report.SITDestFirstDayTotal = &sitDestFirstDay
			report.SITDestAddlDaysTotal = &sitDestAddlDays
			report.SITPickupTotal = &sitPickuptotal
			report.SITDeliveryTotal = &sitDeliveryTotal
			report.SITOriginFuelSurcharge = &sitOriginFuelSurcharge
			report.SITDestFuelSurcharge = &sitDestFuelSurcharge
		}

		report.PpmLinehaul = &ppmLinehaul
		report.PpmFuelRateAdjTotal = &ppmFuel
		report.PpmOriginPrice = &ppmOriginPrice
		report.PpmDestPrice = &ppmDestPrice
		report.PpmPacking = &ppmPacking
		report.PpmUnpacking = &ppmUnpacking
		ppmTotal := ppmLinehaul + ppmFuel + ppmOriginPrice + ppmDestPrice + ppmPacking + ppmUnpacking
		report.PpmTotal = &ppmTotal

		report.ActualOriginNetWeight = &originActualWeight
		report.PBPAndE = &progear
		reweigh := move.MTOShipments[0].Reweigh
		if reweigh != nil {
			report.DestinationReweighNetWeight = reweigh.Weight
		} else {
			report.DestinationReweighNetWeight = nil
		}

		if orders.SAC != nil {
			report.OrderNumber = orders.SAC
		}

		addressLoad := appCtx.DB().Load(&orders.ServiceMember, "ResidentialAddress")
		if addressLoad != nil {
			return nil, apperror.NewQueryError("failed to load residential address", addressLoad, ".")
		}

		netWeight := models.GetTotalNetWeightForMove(move)

		report.FirstName = orders.ServiceMember.FirstName
		report.LastName = orders.ServiceMember.LastName
		report.MiddleInitial = &middleInitial
		report.Affiliation = orders.ServiceMember.Affiliation
		report.PayGrade = orders.Grade
		report.Edipi = orders.ServiceMember.Edipi
		report.PhonePrimary = orders.ServiceMember.Telephone
		report.PhoneSecondary = orders.ServiceMember.SecondaryTelephone
		report.EmailPrimary = orders.ServiceMember.PersonalEmail
		report.EmailSecondary = &orders.ServiceMember.BackupContacts[0].Email
		report.OrdersType = orders.OrdersType
		report.TravelClassCode = (*string)(&orders.OrdersType)
		report.OrdersNumber = orders.OrdersNumber
		report.OrdersDate = &orders.IssueDate
		report.Address = orders.ServiceMember.ResidentialAddress
		report.OriginAddress = move.MTOShipments[0].PickupAddress
		report.DestinationAddress = move.MTOShipments[0].DestinationAddress
		report.OriginGBLOC = orders.OriginDutyLocationGBLOC
		report.DestinationGBLOC = &orders.NewDutyLocation.TransportationOffice.Gbloc
		report.DepCD = orders.HasDependents
		report.TravelAdvance = &travelAdvance
		report.MoveDate = moveDate
		report.TAC = orders.TAC
		report.ShipmentNum = len(move.MTOShipments)
		report.WeightEstimate = calculateTotalWeightEstimate(move.MTOShipments)
		report.TransmitCd = &transmitCd
		report.DD2278IssueDate = move.ServiceCounselingCompletedAt
		report.Miles = move.MTOShipments[0].Distance
		report.ShipmentId = move.ID
		report.SCAC = &scac
		report.NetWeight = &netWeight
		report.PaidDate = paymentRequests[0].ReviewedAt
		report.CounseledDate = move.ServiceCounselingCompletedAt

		fullreport = append(fullreport, report)
	}

	return fullreport, nil
}

// Fetch Moves with an approved Payment Request for Navy service members and ignore TIO and GBLOC rules
func FetchMovesForReports(appCtx appcontext.AppContext, params *services.MoveTaskOrderFetcherParams) (models.Moves, error) {
	var moves models.Moves

	approvedStatuses := []string{models.PaymentRequestStatusReviewed.String(), models.PaymentRequestStatusSentToGex.String(), models.PaymentRequestStatusReceivedByGex.String()}
	query := appCtx.DB().EagerPreload(
		"PaymentRequests",
		"PaymentRequests.PaymentServiceItems",
		"PaymentRequests.PaymentServiceItems.PriceCents",
		"PaymentRequests.PaymentServiceItems.PaymentServiceItemParams.ServiceItemParamKey",
		"MTOShipments.DestinationAddress",
		"MTOShipments.PickupAddress",
		"MTOShipments.SecondaryDeliveryAddress",
		"MTOShipments.SecondaryPickupAddress",
		"MTOShipments.MTOAgents",
		"MTOShipments.Reweigh",
		"MTOShipments.PPMShipment",
		"Orders.ServiceMember",
		"Orders.ServiceMember.ResidentialAddress",
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
		InnerJoin("orders", "orders.id = moves.orders_id").
		InnerJoin("entitlements", "entitlements.id = orders.entitlement_id").
		InnerJoin("service_members", "orders.service_member_id = service_members.id").
		LeftJoin("personally_procured_moves", "personally_procured_moves.move_id = moves.id").
		InnerJoin("mto_shipments", "mto_shipments.move_id = moves.id").
		InnerJoin("addresses", "addresses.id in (mto_shipments.pickup_address_id, mto_shipments.destination_address_id, service_members.residential_address_id)").
		Where("payment_requests.status in (?)", approvedStatuses).
		Where("service_members.affiliation = ?", models.AffiliationNAVY).
		GroupBy("moves.id")

	if params.Since != nil {
		query.Where("payment_requests.updated_at >= ?", params.Since)
	}

	err := query.All(&moves)

	if err != nil {
		return nil, err
	}

	return moves, nil
}

func FetchTACForMmove(appCtx appcontext.AppContext, orders models.Order) []models.TransportationAccountingCode {
	var tac []models.TransportationAccountingCode

	tacQueryError := appCtx.DB().Q().
		EagerPreload(
			"LineOfAccounting",
			"LineOfAccounting.LoaTrsySfxTx",
			"LineOfAccounting.LoaObjClsID",
			"LineOfAccounting.LoaAlltSnID",
			"LineOfAccounting.LoaSbaltmtRcpntID",
			"LineOfAccounting.LoaInstlAcntgActID",
			"LineOfAccounting.LoaTrnsnID",
			"LineOfAccounting.LoaJbOrdNm",
			"LineOfAccounting.LoaDocID",
			"LineOfAccounting.LoaPgmElmntID",
			"LineOfAccounting.LoaDptID",
		).
		Join("lines_of_accounting loa", "loa.loa_sys_id = transportation_accounting_codes.loa_sys_id").
		Where("transportation_accounting_codes.tac = ?", orders.TAC).
		Where("? BETWEEN transportation_accounting_codes.trnsprtn_acnt_bgn_dt AND transportation_accounting_codes.trnsprtn_acnt_end_dt", orders.IssueDate).
		Where("? BETWEEN loa.loa_bgn_dt AND loa.loa_end_dt", orders.IssueDate).
		Where("loa.loa_hs_gds_cd != ?", models.LineOfAccountingHouseholdGoodsCodeNTS).
		All(&tac)

	if tacQueryError != nil {
		return nil
	}

	return tac
}

func buildFullLineOfAccountingString(loa *models.LineOfAccounting) string {
	emptyString := ""
	var fiscalYear string
	if fmt.Sprint(*loa.LoaBgFyTx) != "" && fmt.Sprint(*loa.LoaEndFyTx) != "" {
		fiscalYear = fmt.Sprint(*loa.LoaBgFyTx) + fmt.Sprint(*loa.LoaEndFyTx)
	} else {
		fiscalYear = ""
	}

	if loa.LoaDptID == nil {
		loa.LoaDptID = &emptyString
	}
	if loa.LoaTnsfrDptNm == nil {
		loa.LoaTnsfrDptNm = &emptyString
	}
	if loa.LoaBafID == nil {
		loa.LoaBafID = &emptyString
	}
	if loa.LoaTrsySfxTx == nil {
		loa.LoaTrsySfxTx = &emptyString
	}
	if loa.LoaMajClmNm == nil {
		loa.LoaMajClmNm = &emptyString
	}
	if loa.LoaOpAgncyID == nil {
		loa.LoaOpAgncyID = &emptyString
	}
	if loa.LoaAlltSnID == nil {
		loa.LoaAlltSnID = &emptyString
	}
	if loa.LoaUic == nil {
		loa.LoaUic = &emptyString
	}
	if loa.LoaPgmElmntID == nil {
		loa.LoaPgmElmntID = &emptyString
	}
	if loa.LoaTskBdgtSblnTx == nil {
		loa.LoaTskBdgtSblnTx = &emptyString
	}
	if loa.LoaDfAgncyAlctnRcpntID == nil {
		loa.LoaDfAgncyAlctnRcpntID = &emptyString
	}
	if loa.LoaJbOrdNm == nil {
		loa.LoaJbOrdNm = &emptyString
	}
	if loa.LoaSbaltmtRcpntID == nil {
		loa.LoaSbaltmtRcpntID = &emptyString
	}
	if loa.LoaWkCntrRcpntNm == nil {
		loa.LoaWkCntrRcpntNm = &emptyString
	}
	if loa.LoaMajRmbsmtSrcID == nil {
		loa.LoaMajRmbsmtSrcID = &emptyString
	}
	if loa.LoaDtlRmbsmtSrcID == nil {
		loa.LoaDtlRmbsmtSrcID = &emptyString
	}
	if loa.LoaCustNm == nil {
		loa.LoaCustNm = &emptyString
	}
	if loa.LoaObjClsID == nil {
		loa.LoaObjClsID = &emptyString
	}
	if loa.LoaSrvSrcID == nil {
		loa.LoaSrvSrcID = &emptyString
	}
	if loa.LoaSpclIntrID == nil {
		loa.LoaSpclIntrID = &emptyString
	}
	if loa.LoaBdgtAcntClsNm == nil {
		loa.LoaBdgtAcntClsNm = &emptyString
	}
	if loa.LoaDocID == nil {
		loa.LoaDocID = &emptyString
	}
	if loa.LoaClsRefID == nil {
		loa.LoaClsRefID = &emptyString
	}
	if loa.LoaInstlAcntgActID == nil {
		loa.LoaInstlAcntgActID = &emptyString
	}
	if loa.LoaLclInstlID == nil {
		loa.LoaLclInstlID = &emptyString
	}
	if loa.LoaTrnsnID == nil {
		loa.LoaTrnsnID = &emptyString
	}
	if loa.LoaFmsTrnsactnID == nil {
		loa.LoaFmsTrnsactnID = &emptyString
	}

	LineOfAccountingDfasElementOrder := []string{
		*loa.LoaDptID,               // "LoaDptID"
		*loa.LoaTnsfrDptNm,          // "LoaTnsfrDptNm",
		fiscalYear,                  // "LoaEndFyTx",
		*loa.LoaBafID,               // "LoaBafID",
		*loa.LoaTrsySfxTx,           // "LoaTrsySfxTx",
		*loa.LoaMajClmNm,            // "LoaMajClmNm",
		*loa.LoaOpAgncyID,           // "LoaOpAgncyID",
		*loa.LoaAlltSnID,            // "LoaAlltSnID",
		*loa.LoaUic,                 // "LoaUic",
		*loa.LoaPgmElmntID,          // "LoaPgmElmntID",
		*loa.LoaTskBdgtSblnTx,       // "LoaTskBdgtSblnTx",
		*loa.LoaDfAgncyAlctnRcpntID, // "LoaDfAgncyAlctnRcpntID",
		*loa.LoaJbOrdNm,             // "LoaJbOrdNm",
		*loa.LoaSbaltmtRcpntID,      // "LoaSbaltmtRcpntID",
		*loa.LoaWkCntrRcpntNm,       // "LoaWkCntrRcpntNm",
		*loa.LoaMajRmbsmtSrcID,      // "LoaMajRmbsmtSrcID",
		*loa.LoaDtlRmbsmtSrcID,      // "LoaDtlRmbsmtSrcID",
		*loa.LoaCustNm,              // "LoaCustNm",
		*loa.LoaObjClsID,            // "LoaObjClsID",
		*loa.LoaSrvSrcID,            // "LoaSrvSrcID",
		*loa.LoaSpclIntrID,          // "LoaSpcLIntrID",
		*loa.LoaBdgtAcntClsNm,       // "LoaBdgtAcntCLsNm",
		*loa.LoaDocID,               // "LoaDocID",
		*loa.LoaClsRefID,            // "LoaCLsRefID",
		*loa.LoaInstlAcntgActID,     // "LoaInstLAcntgActID",
		*loa.LoaLclInstlID,          // "LoaLcLInstLID",
		*loa.LoaTrnsnID,             // "LoaTrnsnID",
		*loa.LoaFmsTrnsactnID,       // "LoaFmsTrnsactnID",
	}

	longLoa := strings.Join(LineOfAccountingDfasElementOrder, "*")
	longLoa = strings.ReplaceAll(longLoa, " *", "*")

	return longLoa
}

func buildServiceItemCrate(serviceItem models.MTOServiceItem) pptasmessages.Crate {
	var newServiceItemCrate pptasmessages.Crate
	var newCrateDimensions pptasmessages.CrateCrateDimensions
	var newItemDimensions pptasmessages.CrateItemDimensions

	for dimensionIndex := range serviceItem.Dimensions {
		if serviceItem.Dimensions[dimensionIndex].Type == "ITEM" {
			newItemDimensions.Height = serviceItem.Dimensions[dimensionIndex].Height.ToInches()
			newItemDimensions.Length = serviceItem.Dimensions[dimensionIndex].Length.ToInches()
			newItemDimensions.Width = serviceItem.Dimensions[dimensionIndex].Width.ToInches()
			newServiceItemCrate.ItemDimensions = &newItemDimensions
		}
		if serviceItem.Dimensions[dimensionIndex].Type == "CRATE" {
			newCrateDimensions.Height = serviceItem.Dimensions[dimensionIndex].Height.ToInches()
			newCrateDimensions.Length = serviceItem.Dimensions[dimensionIndex].Length.ToInches()
			newCrateDimensions.Width = serviceItem.Dimensions[dimensionIndex].Width.ToInches()
			newServiceItemCrate.CrateDimensions = &newCrateDimensions
		}
	}

	newServiceItemCrate.Description = *serviceItem.Description

	return newServiceItemCrate
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
