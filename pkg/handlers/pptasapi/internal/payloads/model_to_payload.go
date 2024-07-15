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
	payload := &pptasmessages.ListReport{
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
		TravelClassCode:    *report.TravelClassCode,
		OrdersNumber:       *report.OrdersNumber,
		OrdersDate:         strfmt.DateTime(*report.OrdersDate),
		Address:            Address(report.Address),
		OriginAddress:      Address(report.OriginAddress),
		DestinationAddress: Address(report.DestinationAddress),
		OriginGbloc:        report.OriginGBLOC,
		DestinationGbloc:   report.DestinationGBLOC,
		DepCD:              &report.DepCD,
		TravelAdvance:      models.Float64Pointer(report.TravelAdvance.Float64()),
		MoveDate:           (*strfmt.Date)(report.MoveDate),
		Tac:                report.TAC,
		ShipmentNum:        int64(report.ShipmentNum),
		WeightEstimate:     report.WeightEstimate.Float64(),
		TransmitCD:         report.TransmitCd,
		Dd2278IssueDate:    strfmt.Date(*report.DD2278IssueDate),
		Miles:              int64(*report.Miles),
		ShipmentID:         strfmt.UUID(report.ShipmentId.String()),
		Scac:               report.SCAC,
		NetWeight:          report.NetWeight.Int64(),
		PaidDate:           (*strfmt.Date)(report.PaidDate),
		CounseledDate:      strfmt.Date(*report.CounseledDate),
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
