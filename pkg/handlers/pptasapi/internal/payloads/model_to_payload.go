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
		MiddleInitial:       pptasReport.MiddleInitial,
		PhoneSecondary:      pptasReport.PhoneSecondary,
		EmailSecondary:      pptasReport.EmailSecondary,
		OrdersType:          string(pptasReport.OrdersType),
		PayGrade:            (*string)(pptasReport.PayGrade),
		OriginGbloc:         pptasReport.OriginGBLOC,
		DestinationGbloc:    pptasReport.DestinationGBLOC,
		DepCD:               &pptasReport.DepCD,
		Affiliation:         (*pptasmessages.Affiliation)(pptasReport.Affiliation),
		Tac:                 pptasReport.TAC,
		ShipmentNum:         int64(pptasReport.ShipmentNum),
		TransmitCD:          pptasReport.TransmitCd,
		Scac:                pptasReport.SCAC,
		FinancialReviewFlag: pptasReport.FinancialReviewFlag,
	}

	if len(pptasReport.Shipments) > 0 {
		payload.Shipments = pptasReport.Shipments
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

	if pptasReport.Address != nil {
		payload.Address = Address(pptasReport.Address)
	}

	if pptasReport.EntitlementWeight != nil {
		payload.EntitlementWeight = models.Int64Pointer(pptasReport.EntitlementWeight.Int64())
	}

	if pptasReport.WeightAuthorized != nil {
		payload.WeightAuthorized = models.Float64Pointer(pptasReport.WeightAuthorized.Float64())
	}

	if pptasReport.TravelType != nil {
		payload.TravelType = *pptasReport.TravelType
	}

	if pptasReport.TravelClassCode != nil {
		payload.TravelClassCode = *pptasReport.TravelClassCode
	}

	if pptasReport.CounseledDate != nil {
		payload.CounseledDate = strfmt.Date(*pptasReport.CounseledDate)
	}

	if pptasReport.FinancialReviewRemarks != nil {
		payload.FinancialReviewRemarks = pptasReport.FinancialReviewRemarks
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
