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

	// middleInitial := move.Orders.ServiceMember.MiddleName[0]

	// payload := &pptasmessages.ListReport{
	// 	ID:            strfmt.UUID(pr.MoveTaskOrderID.String()),
	// 	LastName:      *pr.MoveTaskOrder.Orders.ServiceMember.LastName,
	// 	FirstName:     *pr.MoveTaskOrder.Orders.ServiceMember.FirstName,
	// 	MiddleInitial: "w",
	// 	Affiliation:   (*pptasmessages.Affiliation)(pr.MoveTaskOrder.Orders.ServiceMember.Affiliation),
	// 	Grade:         (*string)(pr.MoveTaskOrder.Orders.Grade.Pointer()),
	// 	Edipi:         *pr.MoveTaskOrder.Orders.ServiceMember.Edipi,
	// }

	payload := &pptasmessages.ListReport{
		// ID: *pr.
		FirstName: *report.FirstName,
		Edipi:     *report.Edipi,
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
