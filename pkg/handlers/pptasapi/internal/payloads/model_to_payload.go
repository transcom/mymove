package payloads

import (
	"fmt"
	"strings"
	"time"

	"github.com/go-openapi/strfmt"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/etag"
	"github.com/transcom/mymove/pkg/gen/pptasmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/unit"
)

// InternalServerError describes errors in a standard structure to be returned in the payload.
// If detail is nil, string defaults to "An internal server error has occurred."
func InternalServerError(detail *string, traceID uuid.UUID) *pptasmessages.ClientError {
	instanceToUse := strfmt.UUID(traceID.String())
	payload := pptasmessages.ClientError{
		Title:    handlers.FmtString(handlers.InternalServerErrMessage),
		Detail:   handlers.FmtString(handlers.InternalServerErrDetail),
		Instance: &instanceToUse,
	}
	if detail != nil {
		payload.Detail = detail
	}
	return &payload
}

// ListReport payload
func ListReport(appCtx appcontext.AppContext, move *models.Move) *pptasmessages.ListReport {
	var emptyReport *pptasmessages.ListReport

	if move == nil {
		return nil
	}

	Orders := move.Orders
	var PaymentRequest []models.PaymentRequest
	// get payment requests that have been approved
	for _, pr := range move.PaymentRequests {
		if pr.Status == models.PaymentRequestStatusReviewed || pr.Status == models.PaymentRequestStatusSentToGex || pr.Status == models.PaymentRequestStatusReceivedByGex {
			PaymentRequest = append(PaymentRequest, pr)
		}
	}

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
		Where("transportation_accounting_codes.tac = ?", Orders.TAC).
		Where("? BETWEEN transportation_accounting_codes.trnsprtn_acnt_bgn_dt AND transportation_accounting_codes.trnsprtn_acnt_end_dt", Orders.IssueDate).
		Where("? BETWEEN loa.loa_bgn_dt AND loa.loa_end_dt", Orders.IssueDate).
		Where("loa.loa_hs_gds_cd != ?", models.LineOfAccountingHouseholdGoodsCodeNTS).
		All(&tac)

	if tacQueryError != nil {
		return emptyReport
	}

	var middleInitial string
	if *Orders.ServiceMember.MiddleName != "" {
		middleInitial = string([]rune(*Orders.ServiceMember.MiddleName)[0])
	}
	progear := unit.Pound(0)
	sitTotal := unit.Pound(0)
	originActualWeight := unit.Pound(0)
	travelAdvance := unit.Cents(0)
	scac := "HSFR"
	transmitCd := "T"
	longLoa := buildFullLineOfAccountingString(tac[0].LineOfAccounting)

	var moveDate *time.Time
	if move.MTOShipments[0].PPMShipment != nil {
		moveDate = &move.MTOShipments[0].PPMShipment.ExpectedDepartureDate
	} else {
		moveDate = move.MTOShipments[0].ActualPickupDate
	}

	// get weight entitlements
	if Orders.Grade != nil && Orders.Entitlement != nil {
		Orders.Entitlement.SetWeightAllotment(string(*Orders.Grade))
	}

	weightAllotment := Orders.Entitlement.WeightAllotment()

	var totalWeight int64
	if *Orders.Entitlement.DependentsAuthorized {
		totalWeight = int64(weightAllotment.TotalWeightSelfPlusDependents)
	} else {
		totalWeight = int64(weightAllotment.TotalWeightSelf)
	}

	payload := &pptasmessages.ListReport{
		FirstName:          *Orders.ServiceMember.FirstName,
		LastName:           *Orders.ServiceMember.LastName,
		MiddleInitial:      middleInitial,
		Affiliation:        (*pptasmessages.Affiliation)(Orders.ServiceMember.Affiliation),
		PayGrade:           (*string)(Orders.Grade),
		Edipi:              *Orders.ServiceMember.Edipi,
		PhonePrimary:       *Orders.ServiceMember.Telephone,
		PhoneSecondary:     Orders.ServiceMember.SecondaryTelephone,
		EmailPrimary:       *Orders.ServiceMember.PersonalEmail,
		EmailSecondary:     &Orders.ServiceMember.BackupContacts[0].Email,
		OrdersType:         string(Orders.OrdersType),
		OrdersNumber:       *Orders.OrdersNumber,
		OrdersDate:         strfmt.DateTime(Orders.IssueDate),
		Address:            Address(Orders.ServiceMember.ResidentialAddress),
		OriginAddress:      Address(move.MTOShipments[0].PickupAddress),
		DestinationAddress: Address(move.MTOShipments[0].DestinationAddress),
		OriginGbloc:        Orders.OriginDutyLocationGBLOC,
		DestinationGbloc:   &Orders.NewDutyLocation.TransportationOffice.Gbloc,
		DepCD:              &Orders.HasDependents, // has dependants?
		TravelAdvance:      models.Float64Pointer(travelAdvance.Float64()),
		MoveDate:           (*strfmt.Date)(moveDate),
		Tac:                Orders.TAC,
		FiscalYear:         tac[0].TacFyTxt,
		Appro:              tac[0].LineOfAccounting.LoaBafID,
		Subhead:            tac[0].LineOfAccounting.LoaObjClsID,
		ObjClass:           tac[0].LineOfAccounting.LoaAlltSnID,
		Bcn:                tac[0].LineOfAccounting.LoaSbaltmtRcpntID,
		SubAllotCD:         tac[0].LineOfAccounting.LoaInstlAcntgActID,
		Aaa:                tac[0].LineOfAccounting.LoaTrnsnID,
		TypeCD:             tac[0].LineOfAccounting.LoaJbOrdNm,
		Paa:                tac[0].LineOfAccounting.LoaDocID,
		CostCD:             tac[0].LineOfAccounting.LoaPgmElmntID,
		Ddcd:               tac[0].LineOfAccounting.LoaDptID,
		ShipmentNum:        int64(len(move.MTOShipments)),
		WeightEstimate:     calculateTotalWeightEstimate(move.MTOShipments).Float64(),
		TransmitCD:         &transmitCd, // report.TransmitCd,
		Dd2278IssueDate:    strfmt.Date(*move.ServiceCounselingCompletedAt),
		Miles:              int64(*move.MTOShipments[0].Distance),
		WeightAuthorized:   float64(*Orders.Entitlement.DBAuthorizedWeight),
		ShipmentID:         strfmt.UUID(move.ID.String()),
		Scac:               &scac,
		Loa:                &longLoa,
		ShipmentType:       string(*Orders.OrdersTypeDetail),
		EntitlementWeight:  totalWeight,
		NetWeight:          int64(models.GetTotalNetWeightForMove(*move)), // this only calculates PPM is that correct?
		PickupDate:         strfmt.Date(*move.MTOShipments[0].ActualPickupDate),
		PaidDate:           (*strfmt.Date)(PaymentRequest[0].ReviewedAt),

		// PpmLineHaul:
		// PpmFuelRateAdjTotal:
		// PpmOriginPrice:
		// PpmDestPrice:
		// PpmPacking:
		// PpmUnpacking:
		// PpmStorage:
		// PpmTotal:
		TravelType:      string(*Orders.OrdersTypeDetail),
		TravelClassCode: string(Orders.OrdersType),
		DeliveryDate:    strfmt.Date(*moveDate),
		CounseledDate:   strfmt.Date(*move.ServiceCounselingCompletedAt),
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
	for _, pr := range PaymentRequest {
		for _, serviceItem := range pr.PaymentServiceItems {
			var mtoServiceItem models.MTOServiceItem
			msiErr := appCtx.DB().Q().EagerPreload("ReService", "Dimensions").
				InnerJoin("re_services", "re_services.id = mto_service_items.re_service_id").
				Where("mto_service_items.id = ?", serviceItem.MTOServiceItemID).
				First(&mtoServiceItem)
			if msiErr != nil {
				return emptyReport
			}

			totalPrice := serviceItem.PriceCents.Float64()

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
				sitType := "Origin"
				payload.SitType = &sitType
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
				sitType := "Destination"
				payload.SitType = &sitType
				sitDestFirstDay += totalPrice
			case "Counseling":
				counselingTotal += totalPrice
			}

		}
	}

	shuttleTotal := sitOriginShuttle + sitDestShuttle
	payload.LinehaulTotal = &linehaulTotal
	payload.LinehaulFuelTotal = &fuelPrice
	payload.OriginPrice = &domesticOriginTotal
	payload.DestinationPrice = &domesticDestTotal
	payload.PackingPrice = &domesticPacking
	payload.UnpackingPrice = &domesticUnpacking
	payload.SitOriginFirstDayTotal = &sitOriginFirstDay
	payload.SitOriginAddlDaysTotal = &sitOriginAddlDays
	payload.SitDestFirstDayTotal = &sitDestFirstDay
	payload.SitDestAddlDaysTotal = &sitDestAddlDays
	payload.SitPickupTotal = &sitPickuptotal
	payload.SitDeliveryTotal = &sitDeliveryTotal
	payload.SitOriginFuelSurcharge = &sitOriginFuelSurcharge
	payload.SitDestFuelSurcharge = &sitDestFuelSurcharge
	payload.CratingTotal = &domesticCrating
	payload.UncratingTotal = &domesticUncrating
	payload.ShuttleTotal = &shuttleTotal
	payload.MoveManagementFeeTotal = &managementTotal
	payload.CounselingFeeTotal = &counselingTotal
	payload.CratingDimensions = allCrates

	// calculate total invoice cost
	invoicePaidAmt := shuttleTotal + linehaulTotal + fuelPrice + domesticOriginTotal + domesticDestTotal + domesticPacking + domesticUnpacking +
		sitOriginFirstDay + sitOriginAddlDays + sitDestFirstDay + sitDestAddlDays + sitPickuptotal + sitDeliveryTotal + sitOriginFuelSurcharge +
		sitDestFuelSurcharge + domesticCrating + domesticUncrating
	payload.InvoicePaidAmt = &invoicePaidAmt

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
				payload.SitInDate = (*strfmt.Date)(shipment.PPMShipment.SITEstimatedEntryDate)
				payload.SitOutDate = (*strfmt.Date)(shipment.PPMShipment.SITEstimatedDepartureDate)
				// newreport.SitType = shipment.PPMShipment.
			}
		}

		if shipment.PrimeActualWeight != nil {
			originActualWeight += *shipment.PrimeActualWeight
		}
	}

	payload.ActualOriginNetWeight = float64(originActualWeight)
	payload.PbpAnde = progear.Float64()
	reweigh := move.MTOShipments[0].Reweigh
	if reweigh != nil {
		payload.DestinationReweighNetWeight = models.Float64Pointer(reweigh.Weight.Float64())
	} else {
		payload.DestinationReweighNetWeight = nil
	}

	if Orders.SAC != nil {
		payload.OrdersNumber = *Orders.SAC
	} else {
		payload.OrderNumber = ""
	}

	return payload
}

// ListReports payload
func ListReports(appCtx appcontext.AppContext, moves *models.Moves) []*pptasmessages.ListReport {
	payload := make(pptasmessages.ListReports, len(*moves))

	for i, move := range *moves {
		copyOfMove := move // Make copy to avoid implicit memory aliasing of items from a range statement.
		payload[i] = ListReport(appCtx, &copyOfMove)
	}
	return payload
}

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
		Country:        address.Country,
		County:         &address.County,
		ETag:           etag.GenerateEtag(address.UpdatedAt),
	}
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
