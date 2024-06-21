package payloads

import (
	"github.com/go-openapi/strfmt"
	"github.com/gofrs/uuid"

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
func ListReport(report *models.Report) *pptasmessages.ListReport {
	if report == nil {
		return nil
	}

	payload := &pptasmessages.ListReport{
		// ID:        *report.ID,
		FirstName:          *report.FirstName,
		LastName:           *report.LastName,
		MiddleInitial:      *report.MiddleInitial,
		Affiliation:        (*pptasmessages.Affiliation)(report.Affiliation),
		PayGrade:           (*string)(report.PayGrade),
		Edipi:              *report.Edipi,
		PhonePrimary:       *report.PhonePrimary,
		PhoneSecondary:     report.PhoneSecondary,
		EmailPrimary:       *report.EmailPrimary,
		EmailSecondary:     report.EmailSecondary,
		OrdersType:         string(report.OrdersType),
		OrdersNumber:       *report.OrdersNumber,
		OrdersDate:         strfmt.DateTime(*report.OrdersDate),
		Address:            nil,
		OriginAddress:      nil,
		DestinationAddress: nil,
		OriginGbloc:        report.OriginGBLOC,
		DestinationGbloc:   report.DestinationGBLOC,
		DepCD:              report.DepCD,
		TravelAdvance:      nil, // report.TravelAdvance,
		MoveDate:           (*strfmt.Date)(report.MoveDate),
		TAC:                report.TAC,
		FiscalYear:         nil,
		Appro:              nil, // report.Appro,
		Subhead:            nil, // report.Subhead,
		ObjClass:           nil, // report.ObjClass,
		BCN:                nil, // report.BCN,
		SubAllotCD:         nil, // report.SubAllotCD,
		AAA:                nil, // report.AAA,
		TypeCD:             nil, // report.TypeCD,
		PAA:                nil, // report.PAA,
		CostCD:             nil, // report.CostCD,
		DDCD:               nil, // report.DDCD,
		ShipmentNum:        int64(report.ShipmentNum),
		WeightEstimate:     report.WeightEstimate.Float64(),
		TransmitCD:         nil, // report.TransmitCd,
		DD2278IssueDate:    strfmt.Date(*report.DD2278IssueDate),
		Miles:              0,   // int64(*report.Miles),
		WeightAuthorized:   0.0, //report.WeightAuthorized.Float64(),
		// ShipmentID:         strfmt.UUID(report.ShipmentId),
		SCAC:                        report.SCAC,
		OrderNumber:                 *report.OrderNumber,
		LOA:                         nil, // report.LOA,
		ShipmentType:                "",  // *report.ShipmentType,
		EntitlementWeight:           0,   // report.EntitlementWeight.Int64(),
		NetWeight:                   0,   // report.NetWeight.Int64(),
		PBPAndE:                     0.0, // report.PBPAndE.Float64(),
		PickupDate:                  strfmt.Date(*report.PickupDate),
		SitInDate:                   (*strfmt.Date)(report.SitInDate),
		SitOutDate:                  (*strfmt.Date)(report.SitOutDate),
		SitType:                     report.SitType,
		Rate:                        nil, // report.Rate,
		PaidDate:                    (*strfmt.Date)(report.PaidDate),
		LinehaulTotal:               nil, // report.LinehaulTotal,
		SitTotal:                    nil, // report.SitTotal,
		AccessorialTotal:            nil, // report.AccessorialTotal,
		FuelTotal:                   nil, // report.FuelTotal,
		OtherTotal:                  nil, // report.OtherTotal,
		InvoicePaidAmt:              0.0, // report.InvoicePaidAmt.Float64(),
		TravelType:                  *report.TravelType,
		TravelClassCode:             *report.TravelClassCode,
		DeliveryDate:                strfmt.Date(*report.DeliveryDate),
		ActualOriginNetWeight:       *report.ActualOriginNetWeight,
		DestinationReweighNetWeight: *report.DestinationReweighNetWeight,
		CounseledDate:               strfmt.Date(*report.CounseledDate),
	}

	return payload
}

// ListReports payload
func ListReports(reports *models.Reports) []*pptasmessages.ListReport {
	payload := make(pptasmessages.ListReports, len(*reports))

	for i, report := range *reports {
		copyOfReport := report // Make copy to avoid implicit memory aliasing of items from a range statement.
		payload[i] = ListReport(&copyOfReport)
	}
	return payload
}
