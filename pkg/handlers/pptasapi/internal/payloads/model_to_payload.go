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

// ListReport payload
func ListReport(appCtx appcontext.AppContext, report *models.Report) *pptasmessages.ListReport {
	if report == nil {
		return nil
	}

	payload := &pptasmessages.ListReport{
		FirstName:              *report.FirstName,
		LastName:               *report.LastName,
		MiddleInitial:          *report.MiddleInitial,
		Affiliation:            (*pptasmessages.Affiliation)(report.Affiliation),
		PayGrade:               (*string)(report.PayGrade),
		Edipi:                  *report.Edipi,
		PhonePrimary:           *report.PhonePrimary,
		PhoneSecondary:         report.PhoneSecondary,
		EmailPrimary:           *report.EmailPrimary,
		EmailSecondary:         report.EmailSecondary,
		OrdersType:             string(report.OrdersType),
		OrdersNumber:           *report.OrdersNumber,
		OrdersDate:             strfmt.DateTime(*report.OrdersDate),
		OriginAddress:          Address(report.OriginAddress),
		DestinationAddress:     Address(report.DestinationAddress),
		OriginGbloc:            report.OriginGBLOC,
		DestinationGbloc:       report.DestinationGBLOC,
		DepCD:                  &report.DepCD,
		MoveDate:               (*strfmt.Date)(report.MoveDate),
		Tac:                    report.TAC,
		FiscalYear:             report.FiscalYear,
		ShipmentNum:            int64(report.ShipmentNum),
		TransmitCD:             report.TransmitCd,
		Dd2278IssueDate:        strfmt.Date(*report.DD2278IssueDate),
		ShipmentID:             strfmt.UUID(report.ShipmentId.String()),
		Scac:                   report.SCAC,
		Loa:                    report.LOA,
		SitInDate:              (*strfmt.Date)(report.SitInDate),
		SitOutDate:             (*strfmt.Date)(report.SitOutDate),
		SitType:                report.SitType,
		PaidDate:               (*strfmt.Date)(report.PaidDate),
		LinehaulTotal:          report.LinehaulTotal,
		LinehaulFuelTotal:      report.LinehaulFuelTotal,
		OriginPrice:            report.OriginPrice,
		DestinationPrice:       report.DestinationPrice,
		PackingPrice:           report.PackingPrice,
		UnpackingPrice:         report.UnpackingPrice,
		SitOriginFirstDayTotal: report.SITOriginFirstDayTotal,
		SitOriginAddlDaysTotal: report.SITOriginAddlDaysTotal,
		SitDestFirstDayTotal:   report.SITDestFirstDayTotal,
		SitDestAddlDaysTotal:   report.SITDestAddlDaysTotal,
		SitPickupTotal:         report.SITPickupTotal,
		SitDeliveryTotal:       report.SITDeliveryTotal,
		SitOriginFuelSurcharge: report.SITOriginFuelSurcharge,
		SitDestFuelSurcharge:   report.SITDestFuelSurcharge,
		CratingTotal:           report.CratingTotal,
		UncratingTotal:         report.UncratingTotal,
		CratingDimensions:      report.CratingDimensions,
		ShuttleTotal:           report.ShuttleTotal,
		MoveManagementFeeTotal: report.MoveManagementFeeTotal,
		CounselingFeeTotal:     report.CounselingFeeTotal,
		InvoicePaidAmt:         report.InvoicePaidAmt,
		PpmLinehaul:            report.PpmLinehaul,
		PpmFuelRateAdjTotal:    report.PpmFuelRateAdjTotal,
		PpmOriginPrice:         report.PpmOriginPrice,
		PpmDestPrice:           report.PpmDestPrice,
		PpmPacking:             report.PpmPacking,
		PpmUnpacking:           report.PpmUnpacking,
		PpmStorage:             report.PpmStorage,
		PpmTotal:               report.PpmTotal,
		FinancialReviewFlag:    report.FinancialReviewFlag,
	}

	if report.OrderNumber != nil {
		payload.OrderNumber = report.OrderNumber
	}

	if report.PickupDate != nil {
		payload.PickupDate = strfmt.Date(*report.PickupDate)
	}

	if report.PBPAndE != nil {
		payload.PbpAnde = models.Float64Pointer(report.PBPAndE.Float64())
	}

	if report.TravelAdvance != nil {
		payload.TravelAdvance = models.Float64Pointer(report.TravelAdvance.Float64())
	}

	if report.NetWeight != nil {
		payload.NetWeight = models.Int64Pointer(report.NetWeight.Int64())
	}

	if report.EntitlementWeight != nil {
		payload.EntitlementWeight = models.Int64Pointer(report.EntitlementWeight.Int64())
	}

	if report.Miles != nil {
		payload.Miles = int64(*report.Miles)
	}

	if report.Address != nil {
		payload.Address = Address(report.Address)
	}

	if report.ShipmentType != nil {
		payload.ShipmentType = *report.ShipmentType
	}

	if report.TravelType != nil {
		payload.TravelType = *report.TravelType
	}

	if report.Appro != nil {
		payload.Appro = report.Appro
	}

	if report.Subhead != nil {
		payload.Subhead = report.Subhead
	}

	if report.ObjClass != nil {
		payload.ObjClass = report.ObjClass
	}

	if report.BCN != nil {
		payload.Bcn = report.BCN
	}

	if report.SubAllotCD != nil {
		payload.SubAllotCD = report.SubAllotCD
	}

	if report.AAA != nil {
		payload.Aaa = report.AAA
	}

	if report.TravelType != nil {
		payload.TravelType = *report.TravelType
	}

	if report.PAA != nil {
		payload.Paa = report.PAA
	}

	if report.CostCD != nil {
		payload.CostCD = report.CostCD
	}

	if report.DDCD != nil {
		payload.Ddcd = report.DDCD
	}

	if report.TravelClassCode != nil {
		payload.TravelClassCode = *report.TravelClassCode
	}

	if report.WeightEstimate != nil {
		payload.WeightEstimate = models.Float64Pointer(report.WeightEstimate.Float64())
	}

	if report.ActualOriginNetWeight != nil {
		payload.ActualOriginNetWeight = models.Float64Pointer(report.ActualOriginNetWeight.Float64())
	}

	if report.DestinationReweighNetWeight != nil {
		payload.DestinationReweighNetWeight = models.Float64Pointer(report.DestinationReweighNetWeight.Float64())
	}

	if report.DeliveryDate != nil {
		payload.DeliveryDate = strfmt.Date(*report.DeliveryDate)
	}

	if report.CounseledDate != nil {
		payload.CounseledDate = strfmt.Date(*report.CounseledDate)
	}

	return payload
}

// ListReports payload
func ListReports(appCtx appcontext.AppContext, reports *models.Reports) []*pptasmessages.ListReport {
	payload := make(pptasmessages.ListReports, len(*reports))

	for i, report := range *reports {
		copyOfReport := report // Make copy to avoid implicit memory aliasing of items from a range statement.
		payload[i] = ListReport(appCtx, &copyOfReport)
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
