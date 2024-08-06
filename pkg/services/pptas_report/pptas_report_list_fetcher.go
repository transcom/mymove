package report

import (
	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/gen/pptasmessages"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/unit"
)

type pptasReportListFetcher struct {
	estimator   services.PPMEstimator
	moveFetcher services.MoveFetcher
	tacFetcher  services.TransportationAccountingCodeFetcher
	loaFetcher  services.LineOfAccountingFetcher
}

func NewPPTASReportListFetcher(estimator services.PPMEstimator, moveFetcher services.MoveFetcher, tacFetcher services.TransportationAccountingCodeFetcher, loaFetcher services.LineOfAccountingFetcher) services.PPTASReportListFetcher {
	return &pptasReportListFetcher{
		estimator:   estimator,
		moveFetcher: moveFetcher,
		tacFetcher:  tacFetcher,
		loaFetcher:  loaFetcher,
	}
}

// Builds a list of reports for PPTAS
func (f *pptasReportListFetcher) BuildPPTASReportsFromMoves(appCtx appcontext.AppContext, params *services.MoveTaskOrderFetcherParams) (models.PPTASReports, error) {
	var fullreport models.PPTASReports
	moves, err := f.moveFetcher.FetchMovesForReports(appCtx, params)

	if err != nil {
		return nil, err
	}

	for _, move := range moves {
		var report models.PPTASReport
		orders := move.Orders
		var middleInitial string
		if orders.ServiceMember.MiddleName != nil && *orders.ServiceMember.MiddleName != "" {
			middleInitial = string([]rune(*orders.ServiceMember.MiddleName)[0])
		}

		scac := "HSFR"
		transmitCd := "T"

		// handle orders and service member information here
		report.FirstName = orders.ServiceMember.FirstName
		report.LastName = orders.ServiceMember.LastName
		report.MiddleInitial = &middleInitial
		report.Affiliation = orders.ServiceMember.Affiliation
		report.PayGrade = orders.Grade
		report.Edipi = orders.ServiceMember.Edipi
		report.PhonePrimary = orders.ServiceMember.Telephone
		report.PhoneSecondary = orders.ServiceMember.SecondaryTelephone
		report.EmailPrimary = orders.ServiceMember.PersonalEmail
		report.OrdersType = orders.OrdersType
		report.TravelClassCode = (*string)(&orders.OrdersType)
		report.OrdersNumber = orders.OrdersNumber
		report.OrdersDate = &orders.IssueDate
		report.TAC = orders.TAC
		report.ShipmentNum = len(move.MTOShipments)
		report.SCAC = &scac
		report.OriginGBLOC = orders.OriginDutyLocationGBLOC
		report.DestinationGBLOC = &orders.NewDutyLocation.TransportationOffice.Gbloc
		report.DepCD = orders.HasDependents
		report.TransmitCd = &transmitCd
		report.CounseledDate = move.ServiceCounselingCompletedAt

		if orders.Grade != nil && orders.Entitlement != nil {
			orders.Entitlement.SetWeightAllotment(string(*orders.Grade))
		}

		weightAllotment := orders.Entitlement.WeightAllotment()

		var totalWeight unit.Pound
		if orders.Entitlement.DBAuthorizedWeight != nil && weightAllotment != nil {
			if orders.Entitlement.DependentsAuthorized != nil {
				totalWeight = unit.Pound(weightAllotment.TotalWeightSelfPlusDependents)

				report.WeightAuthorized = (*unit.Pound)(orders.Entitlement.DBAuthorizedWeight)
			} else {
				totalWeight = unit.Pound(weightAllotment.TotalWeightSelf)
			}
		}

		report.EntitlementWeight = &totalWeight

		if orders.ServiceMember.BackupContacts != nil {
			report.EmailSecondary = &orders.ServiceMember.BackupContacts[0].Email
		}

		// if orders.OrdersTypeDetail != nil {
		// 	pptasShipment.ShipmentType = string(*orders.OrdersTypeDetail)
		// 	pptasShipment.TravelType = string(*orders.OrdersTypeDetail)
		// }

		if orders.SAC != nil {
			report.OrderNumber = orders.SAC
		}

		populateShipmentFields(&report, appCtx, move, orders)
		fullreport = append(fullreport, report)
	}

	return fullreport, nil
}

// iterate through mtoshipments and build out PPTASShipment objects for pptas report.
func populateShipmentFields(report *models.PPTASReport, appCtx appcontext.AppContext, move models.Move, orders models.Order) {
	var pptasShipments []*pptasmessages.PPTASShipment
	for _, shipment := range move.MTOShipments {
		var pptasShipment pptasmessages.PPTASShipment

		// pptasShipment.ShipmentID = strfmt.UUID(shipment.ID.String())
		pptasShipment.ShipmentType = string(shipment.ShipmentType)
		pptasShipments = append(pptasShipments, &pptasShipment)
	}

	/**
	for _, shipment := range move.MTOShipments {
		var pptasShipment *pptasmessages.PPTASShipment

		// pptasShipment.ShipmentID = strfmt.UUID(shipment.ID)
		// pptasShipment.ShipmentType = string(shipment.ShipmentType)

		// progear := unit.Pound(0)
		// sitTotal := unit.Pound(0)
		// originActualWeight := unit.Pound(0)
		// travelAdvance := unit.Cents(0)

		// var moveDate *time.Time
		// if shipment.PPMShipment != nil {
		// 	moveDate = &shipment.PPMShipment.ExpectedDepartureDate
		// } else if shipment.ActualPickupDate != nil {
		// 	moveDate = shipment.ActualPickupDate
		// }

		// if moveDate != nil {
		// 	pptasShipment.DeliveryDate = strfmt.Date(*moveDate)
		// }

		// if shipment.ActualPickupDate != nil {
		// 	pptasShipment.PickupDate = strfmt.Date(*shipment.ActualPickupDate)
		// }

		var linehaulTotal, managementTotal, fuelPrice, domesticOriginTotal, domesticDestTotal, domesticPacking,
			domesticUnpacking, domesticCrating, domesticUncrating, counselingTotal, sitPickuptotal, sitOriginFuelSurcharge,
			sitOriginShuttle, sitOriginAddlDays, sitOriginFirstDay, sitDeliveryTotal, sitDestFuelSurcharge, sitDestShuttle,
			sitDestAddlDays, sitDestFirstDay float64

		var allCrates []*pptasmessages.Crate
		var paymentRequests []models.PaymentRequest
		prQuerry := appCtx.DB().EagerPreload(
			"PaymentServiceItems.MTOServiceItem.Dimensions",
			"PaymentServiceItems.MTOServiceItem.ReService").
			Where("payment_requests.move_id = ?", move.ID).
			All(&paymentRequests)
		if prQuerry != nil {
			return nil, apperror.NewQueryError("failed to query payment request", err, ".")
		}

		for _, pr := range paymentRequests {
			if pr.Status == models.PaymentRequestStatusReviewed || pr.Status == models.PaymentRequestStatusSentToGex {
				paymentRequests = append(paymentRequests, pr)
			}
		}

		if len(paymentRequests) < 1 && move.MTOShipments[0].PPMShipment == nil {
			continue
		}

		// tacErr := inputReportTAC(pptasShipment, orders, appCtx, f.tacFetcher, f.loaFetcher)
		// if tacErr != nil {
		// 	return nil, tacErr
		// }

		// this adds up all the different payment service items across all payment requests for a move
		for _, pr := range paymentRequests {
			for _, serviceItem := range pr.PaymentServiceItems {
				totalPrice := serviceItem.PriceCents.Float64()

				switch serviceItem.MTOServiceItem.ReService.Name {
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
				case "Domestic crating":
					crate := buildServiceItemCrate(serviceItem.MTOServiceItem)
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
					if pptasShipment.SitType == nil || *pptasShipment.SitType == "" {
						pptasShipment.SitType = models.StringPointer("Origin")
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
					if pptasShipment.SitType == models.StringPointer("Origin") || pptasShipment.SitType == nil {
						sitType := "Destination"
						pptasShipment.SitType = &sitType
					}
					sitDestFirstDay += totalPrice
				case "Counseling":
					counselingTotal += totalPrice
				default:
					continue
				}
			}
		}

		print(shipment.ActualDeliveryDate)

		// shuttleTotal := sitOriginShuttle + sitDestShuttle
		// pptasShipment.LinehaulTotal = &linehaulTotal
		// pptasShipment.LinehaulFuelTotal = &fuelPrice
		// pptasShipment.OriginPrice = &domesticOriginTotal
		// pptasShipment.DestinationPrice = &domesticDestTotal
		// pptasShipment.PackingPrice = &domesticPacking
		// pptasShipment.UnpackingPrice = &domesticUnpacking
		// pptasShipment.CratingTotal = &domesticCrating
		// pptasShipment.UncratingTotal = &domesticUncrating
		// pptasShipment.ShuttleTotal = &shuttleTotal
		// pptasShipment.MoveManagementFeeTotal = &managementTotal
		// pptasShipment.CounselingFeeTotal = &counselingTotal
		// pptasShipment.CratingDimensions = allCrates

		// // calculate total invoice cost
		// invoicePaidAmt := shuttleTotal + linehaulTotal + fuelPrice + domesticOriginTotal + domesticDestTotal + domesticPacking + domesticUnpacking +
		// 	sitOriginFirstDay + sitOriginAddlDays + sitDestFirstDay + sitDestAddlDays + sitPickuptotal + sitDeliveryTotal + sitOriginFuelSurcharge +
		// 	sitDestFuelSurcharge + domesticCrating + domesticUncrating
		// pptasShipment.InvoicePaidAmt = &invoicePaidAmt

		// var ppmLinehaul, ppmFuel, ppmOriginPrice, ppmDestPrice, ppmPacking, ppmUnpacking float64

		// sharing this for loop for all MTOShipment calculations
		// for _, shipment := range move.MTOShipments {
		// 	if pptasShipment.OriginAddress == nil {
		// 		pptasShipment.OriginAddress = shipment.PickupAddress
		// 	}
		// 	if pptasShipment.DestinationAddress == nil {
		// 		pptasShipment.DestinationAddress = shipment.DestinationAddress
		// 	}

		// 	// calculate total progear for entire move
		// 	if shipment.PPMShipment != nil {
		// 		// query the ppmshipment for all it's child needs for the price breakdown
		// 		var ppmShipment models.PPMShipment
		// 		ppmQ := appCtx.DB().Q().EagerPreload("PickupAddress", "DestinationAddress", "WeightTickets", "Shipment").
		// 			InnerJoin("mto_shipments", "mto_shipments.id = ppm_shipments.shipment_id").
		// 			Where("ppm_shipments.id = ?", shipment.PPMShipment.ID).
		// 			Where("ppm_shipments.status = ?", models.PPMShipmentStatusCloseoutComplete).
		// 			First(&ppmShipment)

		// 			// if the ppm isn't in closeout complete status skip to the next shipment
		// 		if ppmQ != nil && ppmQ.Error() == models.RecordNotFoundErrorString {
		// 			continue
		// 		}

		// 		if ppmQ != nil {
		// 			return nil, apperror.NewQueryError("failed to query ppm ", ppmQ, ".")
		// 		}

		// 		var shipmentTotalProgear float64
		// 		if ppmShipment.ProGearWeight != nil {
		// 			shipmentTotalProgear += ppmShipment.ProGearWeight.Float64()
		// 		}

		// 		if ppmShipment.SpouseProGearWeight != nil {
		// 			shipmentTotalProgear += ppmShipment.SpouseProGearWeight.Float64()
		// 		}

		// 		progear += unit.Pound(shipmentTotalProgear)

		// 		// need to determine which shipment(s) have a ppm and get the travel advances and add them up
		// 		if ppmShipment.AdvanceAmountReceived != nil {
		// 			travelAdvance += *ppmShipment.AdvanceAmountReceived
		// 		}

		// 		// add SIT estimated weights
		// 		if *ppmShipment.SITExpected {
		// 			sitTotal += *ppmShipment.SITEstimatedWeight

		// 			// SIT Fields
		// 			pptasShipment.SitInDate = ppmShipment.SITEstimatedEntryDate
		// 			pptasShipment.SitOutDate = ppmShipment.SITEstimatedDepartureDate
		// 		}

		// 		// do the ppm cost breakdown here
		// 		linehaul, fuel, origin, dest, packing, unpacking, err := f.estimator.PriceBreakdown(appCtx, &ppmShipment)
		// 		if err != nil {
		// 			return nil, apperror.NewUnprocessableEntityError("ppm price breakdown")
		// 		}

		// 		ppmLinehaul += linehaul.Float64()
		// 		ppmFuel += fuel.Float64()
		// 		ppmOriginPrice += origin.Float64()
		// 		ppmDestPrice += dest.Float64()
		// 		ppmPacking += packing.Float64()
		// 		ppmUnpacking += unpacking.Float64()
		// 	}

		// 	if shipment.PrimeActualWeight != nil {
		// 		originActualWeight += *shipment.PrimeActualWeight
		// 	}
		// }

		// if pptasShipment.SitInDate != nil || pptasShipment.SitOutDate != nil {
		// 	pptasShipment.SitOriginFirstDayTotal = &sitOriginFirstDay
		// 	pptasShipment.SitOriginAddlDaysTotal = &sitOriginAddlDays
		// 	pptasShipment.SitDestFirstDayTotal = &sitDestFirstDay
		// 	pptasShipment.SitDestAddlDaysTotal = &sitDestAddlDays
		// 	pptasShipment.SitPickupTotal = &sitPickuptotal
		// 	pptasShipment.SitDeliveryTotal = &sitDeliveryTotal
		// 	pptasShipment.SitOriginFuelSurcharge = &sitOriginFuelSurcharge
		// 	pptasShipment.SitDestFuelSurcharge = &sitDestFuelSurcharge
		// }

		// pptasShipment.PpmLinehaul = &ppmLinehaul
		// pptasShipment.PpmFuelRateAdjTotal = &ppmFuel
		// pptasShipment.PpmOriginPrice = &ppmOriginPrice
		// pptasShipment.PpmDestPrice = &ppmDestPrice
		// pptasShipment.PpmPacking = &ppmPacking
		// pptasShipment.PpmUnpacking = &ppmUnpacking
		// ppmTotal := ppmLinehaul + ppmFuel + ppmOriginPrice + ppmDestPrice + ppmPacking + ppmUnpacking
		// pptasShipment.PpmTotal = &ppmTotal

		// pptasShipment.ActualOriginNetWeight = &originActualWeight
		// pptasShipment.PbpAnde = progear
		// reweigh := shipment.Reweigh
		// if reweigh != nil {
		// 	pptasShipment.DestinationReweighNetWeight = reweigh.Weight
		// } else {
		// 	pptasShipment.DestinationReweighNetWeight = nil
		// }

		addressLoad := appCtx.DB().Load(&orders.ServiceMember, "ResidentialAddress")
		if addressLoad != nil {
			return nil, apperror.NewQueryError("failed to load residential address", addressLoad, ".")
		}

		// netWeight := models.GetTotalNetWeightFromMTO(move)

		// pptasShipment.Address = orders.ServiceMember.ResidentialAddress
		// pptasShipment.TravelAdvance = travelAdvance
		// pptasShipment.MoveDate = moveDate
		// pptasShipment.WeightEstimate = calculateTotalWeightEstimate(move.MTOShipments)
		// pptasShipment.Dd2278IssueDate = strfmt.Date(*move.ServiceCounselingCompletedAt)
		// pptasShipment.Miles = int64(*shipment.Distance)
		// pptasShipment.NetWeight = netWeight
		// if len(paymentRequests) > 0 && paymentRequests[0].ReviewedAt != nil {
		// 	pptasShipment.PaidDate = paymentRequests[0].ReviewedAt
		// }

		financialFlag := move.FinancialReviewFlag
		pptasShipment.FinancialReviewFlag = &financialFlag

		pptasShipments = append(pptasShipments, pptasShipment)
	}
	**/

	report.Shipments = pptasShipments
}

func buildServiceItemCrate(serviceItem models.MTOServiceItem) pptasmessages.Crate {
	var newServiceItemCrate pptasmessages.Crate
	var newCrateDimensions pptasmessages.MTOServiceItemDimension
	var newItemDimensions pptasmessages.MTOServiceItemDimension

	for dimensionIndex := range serviceItem.Dimensions {
		if serviceItem.Dimensions[dimensionIndex].Type == "ITEM" {
			newItemDimensions.Type = pptasmessages.DimensionTypeITEM
			newItemDimensions.Height = int32(serviceItem.Dimensions[dimensionIndex].Height)
			newItemDimensions.Length = int32(serviceItem.Dimensions[dimensionIndex].Length)
			newItemDimensions.Width = int32(serviceItem.Dimensions[dimensionIndex].Width)
			newServiceItemCrate.ItemDimensions = &newItemDimensions
		}
		if serviceItem.Dimensions[dimensionIndex].Type == "CRATE" {
			newCrateDimensions.Type = pptasmessages.DimensionTypeCRATE
			newCrateDimensions.Height = int32(serviceItem.Dimensions[dimensionIndex].Height)
			newCrateDimensions.Length = int32(serviceItem.Dimensions[dimensionIndex].Length)
			newCrateDimensions.Width = int32(serviceItem.Dimensions[dimensionIndex].Width)
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

// inputs all TAC related fields and builds full line of accounting string
func inputReportTAC(pptasShipment *pptasmessages.PPTASShipment, orders models.Order, appCtx appcontext.AppContext, tacFetcher services.TransportationAccountingCodeFetcher, loa services.LineOfAccountingFetcher) error {
	tac, err := tacFetcher.FetchOrderTransportationAccountingCodes(*orders.ServiceMember.Affiliation, orders.IssueDate, *orders.TAC, appCtx)
	if err != nil {
		return err
	} else if len(tac) < 1 {
		return apperror.NewNotFoundError(orders.ID, "No valid TAC found")
	}

	longLoa := loa.BuildFullLineOfAccountingString(tac[0].LineOfAccounting)

	pptasShipment.Loa = &longLoa
	pptasShipment.FiscalYear = tac[0].TacFyTxt
	pptasShipment.Appro = tac[0].LineOfAccounting.LoaBafID
	pptasShipment.Subhead = tac[0].LineOfAccounting.LoaObjClsID
	pptasShipment.ObjClass = tac[0].LineOfAccounting.LoaAlltSnID
	pptasShipment.Bcn = tac[0].LineOfAccounting.LoaSbaltmtRcpntID
	pptasShipment.SubAllotCD = tac[0].LineOfAccounting.LoaInstlAcntgActID
	pptasShipment.Aaa = tac[0].LineOfAccounting.LoaTrnsnID
	pptasShipment.TypeCD = tac[0].LineOfAccounting.LoaJbOrdNm
	pptasShipment.Paa = tac[0].LineOfAccounting.LoaDocID
	pptasShipment.CostCD = tac[0].LineOfAccounting.LoaPgmElmntID
	pptasShipment.Ddcd = tac[0].LineOfAccounting.LoaDptID

	return nil
}
