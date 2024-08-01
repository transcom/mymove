package payloads

import (
	"github.com/go-openapi/strfmt"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/etag"
	"github.com/transcom/mymove/pkg/gen/pptasmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
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

// PPTASReport payload
func PPTASReport(appCtx appcontext.AppContext, pptasReport *models.PPTASReport) *pptasmessages.PPTASReport {
	if pptasReport == nil {
		return nil
	}

	payload := &pptasmessages.PPTASReport{
		MiddleInitial:          pptasReport.MiddleInitial,
		PhoneSecondary:         pptasReport.PhoneSecondary,
		EmailSecondary:         pptasReport.EmailSecondary,
		OrdersType:             string(pptasReport.OrdersType),
		PayGrade:               (*string)(pptasReport.PayGrade),
		OriginAddress:          Address(pptasReport.OriginAddress),
		DestinationAddress:     Address(pptasReport.DestinationAddress),
		OriginGbloc:            pptasReport.OriginGBLOC,
		DestinationGbloc:       pptasReport.DestinationGBLOC,
		DepCD:                  &pptasReport.DepCD,
		Affiliation:            (*pptasmessages.Affiliation)(pptasReport.Affiliation),
		Tac:                    pptasReport.TAC,
		FiscalYear:             pptasReport.FiscalYear,
		ShipmentNum:            int64(pptasReport.ShipmentNum),
		TransmitCD:             pptasReport.TransmitCd,
		Dd2278IssueDate:        strfmt.Date(*pptasReport.DD2278IssueDate),
		ShipmentID:             strfmt.UUID(pptasReport.ShipmentId.String()),
		Scac:                   pptasReport.SCAC,
		Loa:                    pptasReport.LOA,
		SitType:                pptasReport.SitType,
		LinehaulTotal:          pptasReport.LinehaulTotal,
		LinehaulFuelTotal:      pptasReport.LinehaulFuelTotal,
		OriginPrice:            pptasReport.OriginPrice,
		DestinationPrice:       pptasReport.DestinationPrice,
		PackingPrice:           pptasReport.PackingPrice,
		UnpackingPrice:         pptasReport.UnpackingPrice,
		PaidDate:               (*strfmt.Date)(pptasReport.PaidDate),
		MoveDate:               (*strfmt.Date)(pptasReport.MoveDate),
		SitInDate:              (*strfmt.Date)(pptasReport.SitInDate),
		SitOutDate:             (*strfmt.Date)(pptasReport.SitOutDate),
		SitOriginFirstDayTotal: pptasReport.SITOriginFirstDayTotal,
		SitOriginAddlDaysTotal: pptasReport.SITOriginAddlDaysTotal,
		SitDestFirstDayTotal:   pptasReport.SITDestFirstDayTotal,
		SitDestAddlDaysTotal:   pptasReport.SITDestAddlDaysTotal,
		SitPickupTotal:         pptasReport.SITPickupTotal,
		SitDeliveryTotal:       pptasReport.SITDeliveryTotal,
		SitOriginFuelSurcharge: pptasReport.SITOriginFuelSurcharge,
		SitDestFuelSurcharge:   pptasReport.SITDestFuelSurcharge,
		CratingTotal:           pptasReport.CratingTotal,
		UncratingTotal:         pptasReport.UncratingTotal,
		CratingDimensions:      pptasReport.CratingDimensions,
		ShuttleTotal:           pptasReport.ShuttleTotal,
		MoveManagementFeeTotal: pptasReport.MoveManagementFeeTotal,
		CounselingFeeTotal:     pptasReport.CounselingFeeTotal,
		InvoicePaidAmt:         pptasReport.InvoicePaidAmt,
		FinancialReviewFlag:    pptasReport.FinancialReviewFlag,
	}

	if pptasReport.FirstName != nil {
		payload.FirstName = *pptasReport.FirstName
	}

	if pptasReport.LastName != nil {
		payload.LastName = *pptasReport.LastName
	}

	if pptasReport.OrdersDate != nil {
		payload.OrdersDate = strfmt.DateTime(*pptasReport.OrdersDate)
	}

	if pptasReport.Edipi != nil {
		payload.Edipi = *pptasReport.Edipi
	}

	if pptasReport.PhonePrimary != nil {
		payload.PhonePrimary = *pptasReport.PhonePrimary
	}

	if pptasReport.EmailPrimary != nil {
		payload.EmailPrimary = *pptasReport.EmailPrimary
	}

	if pptasReport.OrdersNumber != nil {
		payload.OrdersNumber = *pptasReport.OrdersNumber
	}

	if pptasReport.OrderNumber != nil {
		payload.OrderNumber = pptasReport.OrderNumber
	}

	if pptasReport.PickupDate != nil {
		payload.PickupDate = strfmt.Date(*pptasReport.PickupDate)
	}

	if pptasReport.PBPAndE != nil {
		payload.PbpAnde = models.Float64Pointer(pptasReport.PBPAndE.Float64())
	}

	if pptasReport.TravelAdvance != nil {
		payload.TravelAdvance = models.Float64Pointer(pptasReport.TravelAdvance.Float64())
	}

	if pptasReport.NetWeight != nil {
		payload.NetWeight = models.Int64Pointer(pptasReport.NetWeight.Int64())
	}

	if pptasReport.EntitlementWeight != nil {
		payload.EntitlementWeight = models.Int64Pointer(pptasReport.EntitlementWeight.Int64())
	}

	if pptasReport.Miles != nil {
		payload.Miles = int64(*pptasReport.Miles)
	}

	if pptasReport.Address != nil {
		payload.Address = Address(pptasReport.Address)
	}

	if pptasReport.ShipmentType != nil {
		payload.ShipmentType = *pptasReport.ShipmentType
	}

	if pptasReport.TravelType != nil {
		payload.TravelType = *pptasReport.TravelType
	}

	if pptasReport.Appro != nil {
		payload.Appro = pptasReport.Appro
	}

	if pptasReport.Subhead != nil {
		payload.Subhead = pptasReport.Subhead
	}

	if pptasReport.ObjClass != nil {
		payload.ObjClass = pptasReport.ObjClass
	}

	if pptasReport.BCN != nil {
		payload.Bcn = pptasReport.BCN
	}

	if pptasReport.SubAllotCD != nil {
		payload.SubAllotCD = pptasReport.SubAllotCD
	}

	if pptasReport.AAA != nil {
		payload.Aaa = pptasReport.AAA
	}

	if pptasReport.TravelType != nil {
		payload.TravelType = *pptasReport.TravelType
	}

	if pptasReport.PAA != nil {
		payload.Paa = pptasReport.PAA
	}

	if pptasReport.CostCD != nil {
		payload.CostCD = pptasReport.CostCD
	}

	if pptasReport.DDCD != nil {
		payload.Ddcd = pptasReport.DDCD
	}

	if pptasReport.TravelClassCode != nil {
		payload.TravelClassCode = *pptasReport.TravelClassCode
	}

	if pptasReport.WeightEstimate != nil {
		payload.WeightEstimate = models.Float64Pointer(pptasReport.WeightEstimate.Float64())
	}

	if pptasReport.ActualOriginNetWeight != nil {
		payload.ActualOriginNetWeight = models.Float64Pointer(pptasReport.ActualOriginNetWeight.Float64())
	}

	if pptasReport.DestinationReweighNetWeight != nil {
		payload.DestinationReweighNetWeight = models.Float64Pointer(pptasReport.DestinationReweighNetWeight.Float64())
	}

	if pptasReport.DeliveryDate != nil {
		payload.DeliveryDate = strfmt.Date(*pptasReport.DeliveryDate)
	}

	if pptasReport.CounseledDate != nil {
		payload.CounseledDate = strfmt.Date(*pptasReport.CounseledDate)
	}

	emptyCost := float64(0)
	if pptasReport.PpmLinehaul != nil && *pptasReport.PpmLinehaul != emptyCost {
		payload.PpmLinehaul = pptasReport.PpmLinehaul
	}

	if pptasReport.PpmFuelRateAdjTotal != nil && *pptasReport.PpmFuelRateAdjTotal != emptyCost {
		payload.PpmFuelRateAdjTotal = pptasReport.PpmFuelRateAdjTotal
	}

	if pptasReport.PpmOriginPrice != nil && *pptasReport.PpmOriginPrice != emptyCost {
		payload.PpmOriginPrice = pptasReport.PpmOriginPrice
	}

	if pptasReport.PpmDestPrice != nil && *pptasReport.PpmDestPrice != emptyCost {
		payload.PpmDestPrice = pptasReport.PpmDestPrice
	}

	if pptasReport.PpmPacking != nil && *pptasReport.PpmPacking != emptyCost {
		payload.PpmPacking = pptasReport.PpmPacking
	}

	if pptasReport.PpmUnpacking != nil && *pptasReport.PpmUnpacking != emptyCost {
		payload.PpmUnpacking = pptasReport.PpmUnpacking
	}

	if pptasReport.PpmTotal != nil && *pptasReport.PpmTotal != emptyCost {
		payload.PpmTotal = pptasReport.PpmTotal
	}

	return payload
}

// PPTASReports payload
func PPTASReports(appCtx appcontext.AppContext, pptasReports *models.PPTASReports) []*pptasmessages.PPTASReport {
	payload := make(pptasmessages.PPTASReports, len(*pptasReports))

	for i, pptasReport := range *pptasReports {
		copyOfPPTASReport := pptasReport // Make copy to avoid implicit memory aliasing of items from a range statement.
		payload[i] = PPTASReport(appCtx, &copyOfPPTASReport)
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
