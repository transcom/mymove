package report

import (
	"time"

	"github.com/go-openapi/strfmt"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/etag"
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

func (f *pptasReportListFetcher) GetMovesForReportBuilder(appCtx appcontext.AppContext, params *services.MoveTaskOrderFetcherParams) (models.Moves, error) {
	moves, err := f.moveFetcher.FetchMovesForPPTASReports(appCtx, params)

	if err != nil {
		return nil, err
	}

	return moves, err
}

// Builds a list of reports for PPTAS
func (f *pptasReportListFetcher) BuildPPTASReportsFromMoves(appCtx appcontext.AppContext, moves models.Moves) (models.PPTASReports, error) {
	var fullreport models.PPTASReports

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

		financialFlag := move.FinancialReviewFlag
		report.FinancialReviewFlag = &financialFlag

		financialRemarks := move.FinancialReviewRemarks
		report.FinancialReviewRemarks = financialRemarks

		addressLoad := appCtx.DB().Load(&orders.ServiceMember, "ResidentialAddress")
		if addressLoad != nil {
			return nil, apperror.NewQueryError("failed to load residential address", addressLoad, ".")
		}
		report.Address = orders.ServiceMember.ResidentialAddress

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

				report.WeightAuthorized = (*unit.Pound)(orders.Entitlement.DBAuthorizedWeight)
			}
		}

		report.EntitlementWeight = &totalWeight

		if orders.ServiceMember.BackupContacts != nil {
			report.EmailSecondary = &orders.ServiceMember.BackupContacts[0].Email
		}

		if orders.OrdersTypeDetail != nil {
			report.TravelType = (*string)(orders.OrdersTypeDetail)
		}

		err := populateShipmentFields(&report, appCtx, move, orders, f.tacFetcher, f.loaFetcher, f.estimator)
		if err != nil {
			return nil, err
		}

		fullreport = append(fullreport, report)
	}

	return fullreport, nil
}

// iterate through mtoshipments and build out PPTASShipment objects for pptas report.
func populateShipmentFields(
	report *models.PPTASReport, appCtx appcontext.AppContext, move models.Move,
	orders models.Order, tacFetcher services.TransportationAccountingCodeFetcher,
	loaFetcher services.LineOfAccountingFetcher, estimator services.PPMEstimator) error {
	var pptasShipments []*pptasmessages.PPTASShipment
	for _, shipment := range move.MTOShipments {
		var pptasShipment pptasmessages.PPTASShipment

		pptasShipment.ShipmentID = strfmt.UUID(shipment.ID.String())
		pptasShipment.ShipmentType = string(shipment.ShipmentType)

		var moveDate *time.Time
		if shipment.ActualPickupDate != nil {
			moveDate = shipment.ActualPickupDate
			pptasShipment.MoveDate = (*strfmt.Date)(moveDate)
		}

		if moveDate != nil && shipment.ActualDeliveryDate != nil {
			pptasShipment.DeliveryDate = strfmt.Date(*shipment.ActualDeliveryDate)
		}

		if shipment.ActualPickupDate != nil {
			pptasShipment.PickupDate = strfmt.Date(*shipment.ActualPickupDate)
		}

		pptasShipment.MoveDate = (*strfmt.Date)(moveDate)

		if move.ServiceCounselingCompletedAt != nil {
			pptasShipment.Dd2278IssueDate = strfmt.Date(*move.ServiceCounselingCompletedAt)
		} else if move.PrimeCounselingCompletedAt != nil {
			pptasShipment.Dd2278IssueDate = strfmt.Date(*move.PrimeCounselingCompletedAt)
		}

		// location fields
		if pptasShipment.OriginAddress == nil {
			pptasShipment.OriginAddress = Address(shipment.PickupAddress)
		}
		if pptasShipment.DestinationAddress == nil {
			pptasShipment.DestinationAddress = Address(shipment.DestinationAddress)
		}

		// populate TGET data
		tacErr := inputReportTAC(report, &pptasShipment, orders, appCtx, tacFetcher, loaFetcher)
		if tacErr != nil {
			return tacErr
		}

		// populate payment request data
		err := populatePaymentRequestFields(&pptasShipment, appCtx, shipment)
		if err != nil {
			return err
		}

		// populate ppm data
		err = populatePPMFields(appCtx, &pptasShipment, shipment, estimator)
		if err != nil {
			return err
		}

		var originActualWeight float64
		if pptasShipment.ActualOriginNetWeight == nil && shipment.PrimeActualWeight != nil {
			originActualWeight = shipment.PrimeActualWeight.Float64()
			pptasShipment.ActualOriginNetWeight = &originActualWeight
		}

		if shipment.Reweigh != nil && shipment.Reweigh.Weight != nil {
			reweigh := shipment.Reweigh.Weight.Float64()
			pptasShipment.DestinationReweighNetWeight = &reweigh
		}

		netWeight := models.GetTotalNetWeightForMTOShipment(shipment).Int64()
		pptasShipment.NetWeight = &netWeight

		var weightEstimate float64
		if shipment.PPMShipment != nil && shipment.PPMShipment.EstimatedWeight != nil {
			weightEstimate = shipment.PPMShipment.EstimatedWeight.Float64()
		}

		if shipment.PrimeEstimatedWeight != nil {
			weightEstimate = shipment.PrimeEstimatedWeight.Float64()
		}
		pptasShipment.WeightEstimate = &weightEstimate

		if shipment.Distance != nil {
			pptasShipment.Miles = int64(*shipment.Distance)
		}

		pptasShipments = append(pptasShipments, &pptasShipment)
	}

	report.Shipments = pptasShipments

	return nil
}

func populatePaymentRequestFields(pptasShipment *pptasmessages.PPTASShipment, appCtx appcontext.AppContext, shipment models.MTOShipment) error {
	var paymentRequests []models.PaymentRequest
	approvedStatuses := []string{
		models.PaymentRequestStatusReviewed.String(),
		models.PaymentRequestStatusSentToGex.String(),
		models.PaymentRequestStatusPaid.String(),
		models.PaymentRequestStatusEDIError.String(),
		models.PaymentRequestStatusTppsReceived.String(),
	}

	prQErr := appCtx.DB().EagerPreload(
		"PaymentServiceItems.MTOServiceItem.ReService").
		InnerJoin("payment_service_items", "payment_requests.id = payment_service_items.payment_request_id").
		InnerJoin("mto_service_items", "mto_service_items.id = payment_service_items.mto_service_item_id").
		Where("mto_service_items.mto_shipment_id = ?", shipment.ID).
		Where("payment_requests.status in (?)", approvedStatuses).
		GroupBy("payment_requests.id").
		All(&paymentRequests)
	if prQErr != nil {
		return apperror.NewQueryError("failed to query payment request", prQErr, ".")
	}

	if len(paymentRequests) < 1 {
		return nil
	}

	var linehaulTotal, managementTotal, fuelPrice, domesticOriginTotal, domesticDestTotal, domesticPacking,
		domesticUnpacking, domesticCrating, domesticUncrating, counselingTotal, sitPickuptotal, sitOriginFuelSurcharge,
		sitOriginShuttle, sitOriginAddlDays, sitOriginFirstDay, sitDeliveryTotal, sitDestFuelSurcharge, sitDestShuttle,
		sitDestAddlDays, sitDestFirstDay float64

	var allCrates []*pptasmessages.Crate

	// assign the service item cost to the corresponding variable
	for _, pr := range paymentRequests {
		for _, serviceItem := range pr.PaymentServiceItems {
			mtoServiceItem := serviceItem.MTOServiceItem

			err := appCtx.DB().Load(&mtoServiceItem, "Dimensions")
			if err != nil {
				return err
			}

			var totalPrice float64
			if serviceItem.PriceCents != nil {
				totalPrice = serviceItem.PriceCents.Float64()
			}

			if serviceItem.MTOServiceItem.SITEntryDate != nil {
				pptasShipment.SitInDate = (*strfmt.Date)(serviceItem.MTOServiceItem.SITEntryDate)
			}

			if serviceItem.MTOServiceItem.SITDepartureDate != nil {
				pptasShipment.SitOutDate = (*strfmt.Date)(serviceItem.MTOServiceItem.SITDepartureDate)
			}

			switch serviceItem.MTOServiceItem.ReService.Name {
			case "Domestic linehaul":
				linehaulTotal += totalPrice
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

		// Paid date is the earliest payment request date
		if pr.PaidAt != nil && pptasShipment.PaidDate == nil {
			paidDate := strfmt.Date(*pr.PaidAt)
			pptasShipment.PaidDate = &paidDate
		} else if pr.PaidAt != nil && !pr.PaidAt.After(time.Time(*pptasShipment.PaidDate)) {
			paidDate := strfmt.Date(*pr.PaidAt)
			pptasShipment.PaidDate = &paidDate
		}
	}

	shuttleTotal := sitOriginShuttle + sitDestShuttle
	pptasShipment.LinehaulTotal = &linehaulTotal
	pptasShipment.LinehaulFuelTotal = &fuelPrice
	pptasShipment.OriginPrice = &domesticOriginTotal
	pptasShipment.DestinationPrice = &domesticDestTotal
	pptasShipment.PackingPrice = &domesticPacking
	pptasShipment.UnpackingPrice = &domesticUnpacking
	pptasShipment.CratingTotal = &domesticCrating
	pptasShipment.UncratingTotal = &domesticUncrating
	pptasShipment.ShuttleTotal = &shuttleTotal
	pptasShipment.MoveManagementFeeTotal = &managementTotal
	pptasShipment.CounselingFeeTotal = &counselingTotal
	pptasShipment.CratingDimensions = allCrates

	// calculate total invoice cost
	invoicePaidAmt := shuttleTotal + linehaulTotal + fuelPrice + domesticOriginTotal + domesticDestTotal + domesticPacking + domesticUnpacking +
		sitOriginFirstDay + sitOriginAddlDays + sitDestFirstDay + sitDestAddlDays + sitPickuptotal + sitDeliveryTotal + sitOriginFuelSurcharge +
		sitDestFuelSurcharge + domesticCrating + domesticUncrating
	pptasShipment.InvoicePaidAmt = &invoicePaidAmt

	if pptasShipment.SitInDate != nil || pptasShipment.SitOutDate != nil {
		pptasShipment.SitOriginFirstDayTotal = &sitOriginFirstDay
		pptasShipment.SitOriginAddlDaysTotal = &sitOriginAddlDays
		pptasShipment.SitDestFirstDayTotal = &sitDestFirstDay
		pptasShipment.SitDestAddlDaysTotal = &sitDestAddlDays
		pptasShipment.SitPickupTotal = &sitPickuptotal
		pptasShipment.SitDeliveryTotal = &sitDeliveryTotal
		pptasShipment.SitOriginFuelSurcharge = &sitOriginFuelSurcharge
		pptasShipment.SitDestFuelSurcharge = &sitDestFuelSurcharge
	}

	return nil
}

// populates ppm related fields (progear, ppm costs, SIT)
func populatePPMFields(appCtx appcontext.AppContext, pptasShipment *pptasmessages.PPTASShipment, shipment models.MTOShipment, estimator services.PPMEstimator) error {
	var travelAdvance float64

	var ppmLinehaul, ppmFuel, ppmOriginPrice, ppmDestPrice, ppmPacking, ppmUnpacking, ppmStorage float64
	if shipment.PPMShipment != nil && (shipment.PPMShipment.Status == models.PPMShipmentStatusCloseoutComplete || shipment.PPMShipment.Status == models.PPMShipmentStatusComplete) {
		// query the ppmshipment for all it's child needs for the price breakdown
		var ppmShipment models.PPMShipment
		ppmQ := appCtx.DB().Q().EagerPreload("PickupAddress", "DestinationAddress", "WeightTickets", "Shipment").
			InnerJoin("mto_shipments", "mto_shipments.id = ppm_shipments.shipment_id").
			Where("ppm_shipments.id = ?", shipment.PPMShipment.ID).
			Where("ppm_shipments.status = ?", models.PPMShipmentStatusCloseoutComplete).
			First(&ppmShipment)

		// if the ppm isn't in closeout complete status skip to the next shipment
		if ppmQ != nil && ppmQ.Error() == models.RecordNotFoundErrorString {
			return ppmQ
		}

		if ppmQ != nil {
			return apperror.NewQueryError("failed to query ppm ", ppmQ, ".")
		}

		if pptasShipment.OriginAddress == nil {
			pptasShipment.OriginAddress = Address(ppmShipment.PickupAddress)
		}
		if pptasShipment.DestinationAddress == nil {
			pptasShipment.DestinationAddress = Address(ppmShipment.DestinationAddress)
		}

		moveDate := &shipment.PPMShipment.ExpectedDepartureDate
		pptasShipment.MoveDate = (*strfmt.Date)(moveDate)

		pptasShipment.DeliveryDate = strfmt.Date(*ppmShipment.ActualMoveDate)

		ppmNetWeight := calculatePPMNetWeight(ppmShipment)
		pptasShipment.ActualOriginNetWeight = models.Float64Pointer(ppmNetWeight)

		var shipmentTotalProgear float64
		if ppmShipment.ProGearWeight != nil {
			shipmentTotalProgear += ppmShipment.ProGearWeight.Float64()
		}

		if ppmShipment.SpouseProGearWeight != nil {
			shipmentTotalProgear += ppmShipment.SpouseProGearWeight.Float64()
		}

		pptasShipment.PbpAnde = &shipmentTotalProgear

		// need to determine which shipment(s) have a ppm and get the travel advances and add them up
		if ppmShipment.AdvanceAmountReceived != nil {
			travelAdvance = ppmShipment.AdvanceAmountReceived.Float64()
			pptasShipment.TravelAdvance = &travelAdvance
		}

		// add SIT fields
		if ppmShipment.SITExpected != nil && *ppmShipment.SITExpected {
			pptasShipment.SitInDate = (*strfmt.Date)(ppmShipment.SITEstimatedEntryDate)
			pptasShipment.SitOutDate = (*strfmt.Date)(ppmShipment.SITEstimatedDepartureDate)
		}

		// do the ppm cost breakdown here
		linehaul, fuel, origin, dest, packing, unpacking, storage, err := estimator.PriceBreakdown(appCtx, &ppmShipment)
		if err != nil {
			return apperror.NewUnprocessableEntityError("ppm price breakdown")
		}

		ppmLinehaul += linehaul.Float64()
		ppmFuel += fuel.Float64()
		ppmOriginPrice += origin.Float64()
		ppmDestPrice += dest.Float64()
		ppmPacking += packing.Float64()
		ppmUnpacking += unpacking.Float64()
		ppmStorage += storage.Float64()
		ppmTotal := ppmLinehaul + ppmFuel + ppmOriginPrice + ppmDestPrice + ppmPacking + ppmUnpacking + ppmStorage

		pptasShipment.PpmLinehaul = &ppmLinehaul
		pptasShipment.PpmFuelRateAdjTotal = &ppmFuel
		pptasShipment.PpmOriginPrice = &ppmOriginPrice
		pptasShipment.PpmDestPrice = &ppmDestPrice
		pptasShipment.PpmPacking = &ppmPacking
		pptasShipment.PpmUnpacking = &ppmUnpacking
		pptasShipment.PpmStorage = &ppmStorage
		pptasShipment.PpmTotal = &ppmTotal
	}

	return nil
}

// calculate the ppm net weight by taking the difference in full weight and empty weight in the weight tickets
func calculatePPMNetWeight(ppmShipment models.PPMShipment) float64 {
	totalNetWeight := unit.Pound(0)

	for _, weightTicket := range ppmShipment.WeightTickets {
		totalNetWeight += *weightTicket.FullWeight - *weightTicket.EmptyWeight
	}

	return totalNetWeight.Float64()
}

// #nosec G115: it is unrealistic that an imperial measurement will exceed int32 limits
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

// inputs all TAC related fields and builds full line of accounting string
func inputReportTAC(report *models.PPTASReport, pptasShipment *pptasmessages.PPTASShipment, orders models.Order, appCtx appcontext.AppContext, tacFetcher services.TransportationAccountingCodeFetcher, loa services.LineOfAccountingFetcher) error {
	tac, err := tacFetcher.FetchOrderTransportationAccountingCodes(models.DepartmentIndicator(*orders.DepartmentIndicator), orders.IssueDate, *orders.TAC, appCtx)
	if err != nil {
		return err
	} else if len(tac) < 1 {
		return nil
	}

	if tac[0].LineOfAccounting != nil {
		longLoa := loa.BuildFullLineOfAccountingString(*tac[0].LineOfAccounting)

		pptasShipment.Loa = &longLoa
		pptasShipment.FiscalYear = tac[0].TacFyTxt
		pptasShipment.Appro = tac[0].LineOfAccounting.LoaBafID
		pptasShipment.Subhead = tac[0].LineOfAccounting.LoaTrsySfxTx
		pptasShipment.ObjClass = tac[0].LineOfAccounting.LoaObjClsID
		pptasShipment.Bcn = tac[0].LineOfAccounting.LoaAlltSnID
		pptasShipment.SubAllotCD = tac[0].LineOfAccounting.LoaSbaltmtRcpntID
		pptasShipment.Aaa = tac[0].LineOfAccounting.LoaTrnsnID
		pptasShipment.TypeCD = tac[0].LineOfAccounting.LoaJbOrdNm
		pptasShipment.Paa = tac[0].LineOfAccounting.LoaInstlAcntgActID
		pptasShipment.CostCD = tac[0].LineOfAccounting.LoaPgmElmntID
		pptasShipment.Ddcd = tac[0].LineOfAccounting.LoaDptID

		if report.OrderNumber == nil {
			report.OrderNumber = tac[0].LineOfAccounting.LoaDocID
		}
	}

	return nil
}

// Country payload
func Country(country *models.Country) *string {
	if country == nil {
		return nil
	}
	return &country.Country
}

// converts models.Address into payload address
func Address(address *models.Address) *pptasmessages.Address {
	if address == nil {
		return nil
	}
	return &pptasmessages.Address{
		ID:             strfmt.UUID(address.ID.String()),
		StreetAddress1: &address.StreetAddress1,
		StreetAddress2: address.StreetAddress2,
		StreetAddress3: address.StreetAddress3,
		City:           &address.City,
		State:          &address.State,
		PostalCode:     &address.PostalCode,
		Country:        Country(address.Country),
		County:         address.County,
		ETag:           etag.GenerateEtag(address.UpdatedAt),
	}
}
