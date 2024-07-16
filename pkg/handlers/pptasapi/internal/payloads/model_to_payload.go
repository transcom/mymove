package payloads

import (
	"reflect"

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
		OrdersDate:             strfmt.DateTime(*report.OrdersDate),
		Address:                Address(report.Address),
		OriginAddress:          Address(report.OriginAddress),
		DestinationAddress:     Address(report.DestinationAddress),
		OriginGbloc:            report.OriginGBLOC,
		DestinationGbloc:       report.DestinationGBLOC,
		DepCD:                  &report.DepCD,
		TravelAdvance:          models.Float64Pointer(report.TravelAdvance.Float64()),
		MoveDate:               (*strfmt.Date)(report.MoveDate),
		Tac:                    report.TAC,
		FiscalYear:             report.FiscalYear,
		ShipmentNum:            int64(report.ShipmentNum),
		TransmitCD:             report.TransmitCd,
		Dd2278IssueDate:        strfmt.Date(*report.DD2278IssueDate),
		Miles:                  int64(*report.Miles),
		ShipmentID:             strfmt.UUID(report.ShipmentId.String()),
		Scac:                   report.SCAC,
		Loa:                    report.LOA,
		ShipmentType:           *report.ShipmentType,
		EntitlementWeight:      report.EntitlementWeight.Int64(),
		NetWeight:              report.NetWeight.Int64(),
		PbpAnde:                report.PBPAndE.Float64(),
		PickupDate:             strfmt.Date(*report.PickupDate),
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
		TravelType:             *report.TravelType,
		DeliveryDate:           strfmt.Date(*report.DeliveryDate),
		ActualOriginNetWeight:  report.ActualOriginNetWeight.Float64(),
		CounseledDate:          strfmt.Date(*report.CounseledDate),
	}

	if report.OrderNumber != "" {
		payload.OrderNumber = report.OrderNumber
	}

	if !reflect.ValueOf(report.Appro).IsNil() {
		payload.Appro = report.Appro
	}

	if !reflect.ValueOf(report.Subhead).IsNil() {
		payload.Subhead = report.Subhead
	}

	if !reflect.ValueOf(report.ObjClass).IsNil() {
		payload.ObjClass = report.ObjClass
	}

	if !reflect.ValueOf(report.BCN).IsNil() {
		payload.Bcn = report.BCN
	}

	if !reflect.ValueOf(report.SubAllotCD).IsNil() {
		payload.SubAllotCD = report.SubAllotCD
	}

	if !reflect.ValueOf(report.AAA).IsNil() {
		payload.Aaa = report.AAA
	}

	if !reflect.ValueOf(report.TravelType).IsNil() {
		payload.TravelType = *report.TravelType
	}

	if !reflect.ValueOf(report.PAA).IsNil() {
		payload.Paa = report.PAA
	}

	if !reflect.ValueOf(report.CostCD).IsNil() {
		payload.CostCD = report.CostCD
	}

	if !reflect.ValueOf(report.DDCD).IsNil() {
		payload.Ddcd = report.DDCD
	}

	if !reflect.ValueOf(report.TravelClassCode).IsNil() {
		payload.TravelClassCode = *report.TravelClassCode
	}

	if !reflect.ValueOf(report.WeightEstimate).IsNil() {
		payload.WeightEstimate = report.WeightEstimate.Float64()
	}

	if !reflect.ValueOf(report.DestinationReweighNetWeight).IsNil() {
		payload.DestinationReweighNetWeight = models.Float64Pointer(report.DestinationReweighNetWeight.Float64())
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
