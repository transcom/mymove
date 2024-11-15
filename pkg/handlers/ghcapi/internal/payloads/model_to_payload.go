package payloads

import (
	"encoding/json"
	"errors"
	"math"
	"strings"
	"time"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/gofrs/uuid"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/etag"
	"github.com/transcom/mymove/pkg/gen/ghcmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/models/roles"
	"github.com/transcom/mymove/pkg/services"
	mtoshipment "github.com/transcom/mymove/pkg/services/mto_shipment"
	"github.com/transcom/mymove/pkg/storage"
	"github.com/transcom/mymove/pkg/unit"
)

// Contractor payload
func Contractor(contractor *models.Contractor) *ghcmessages.Contractor {
	if contractor == nil {
		return nil
	}

	payload := &ghcmessages.Contractor{
		ID:             strfmt.UUID(contractor.ID.String()),
		ContractNumber: contractor.ContractNumber,
		Name:           contractor.Name,
		Type:           contractor.Type,
	}

	return payload
}

func OfficeUser(officeUser *models.OfficeUser) *ghcmessages.LockedOfficeUser {
	if officeUser != nil {
		payload := ghcmessages.LockedOfficeUser{
			FirstName:              officeUser.FirstName,
			LastName:               officeUser.LastName,
			TransportationOfficeID: *handlers.FmtUUID(officeUser.TransportationOfficeID),
			TransportationOffice:   TransportationOffice(&officeUser.TransportationOffice),
		}
		return &payload
	}
	return nil
}

func AssignedOfficeUser(officeUser *models.OfficeUser) *ghcmessages.AssignedOfficeUser {
	if officeUser != nil {
		payload := ghcmessages.AssignedOfficeUser{
			OfficeUserID: strfmt.UUID(officeUser.ID.String()),
			FirstName:    officeUser.FirstName,
			LastName:     officeUser.LastName,
		}
		return &payload
	}
	return nil
}

// Move payload
func Move(move *models.Move, storer storage.FileStorer) (*ghcmessages.Move, error) {
	if move == nil {
		return nil, nil
	}
	// Adds shipmentGBLOC to be used for TOO/TIO's origin GBLOC
	var gbloc ghcmessages.GBLOC
	if len(move.ShipmentGBLOC) > 0 && move.ShipmentGBLOC[0].GBLOC != nil {
		gbloc = ghcmessages.GBLOC(*move.ShipmentGBLOC[0].GBLOC)
	} else if move.Orders.OriginDutyLocationGBLOC != nil {
		gbloc = ghcmessages.GBLOC(*move.Orders.OriginDutyLocationGBLOC)
	}

	var additionalDocumentsPayload *ghcmessages.Document
	var err error
	if move.AdditionalDocuments != nil {
		additionalDocumentsPayload, err = PayloadForDocumentModel(storer, *move.AdditionalDocuments)
	}
	if err != nil {
		return nil, err
	}

	payload := &ghcmessages.Move{
		ID:                           strfmt.UUID(move.ID.String()),
		AvailableToPrimeAt:           handlers.FmtDateTimePtr(move.AvailableToPrimeAt),
		ApprovedAt:                   handlers.FmtDateTimePtr(move.ApprovedAt),
		ContractorID:                 handlers.FmtUUIDPtr(move.ContractorID),
		Contractor:                   Contractor(move.Contractor),
		Locator:                      move.Locator,
		OrdersID:                     strfmt.UUID(move.OrdersID.String()),
		Orders:                       Order(&move.Orders),
		ReferenceID:                  handlers.FmtStringPtr(move.ReferenceID),
		Status:                       ghcmessages.MoveStatus(move.Status),
		ExcessWeightQualifiedAt:      handlers.FmtDateTimePtr(move.ExcessWeightQualifiedAt),
		BillableWeightsReviewedAt:    handlers.FmtDateTimePtr(move.BillableWeightsReviewedAt),
		CreatedAt:                    strfmt.DateTime(move.CreatedAt),
		SubmittedAt:                  handlers.FmtDateTimePtr(move.SubmittedAt),
		ApprovalsRequestedAt:         handlers.FmtDateTimePtr(move.ApprovalsRequestedAt),
		UpdatedAt:                    strfmt.DateTime(move.UpdatedAt),
		ETag:                         etag.GenerateEtag(move.UpdatedAt),
		ServiceCounselingCompletedAt: handlers.FmtDateTimePtr(move.ServiceCounselingCompletedAt),
		ExcessWeightAcknowledgedAt:   handlers.FmtDateTimePtr(move.ExcessWeightAcknowledgedAt),
		TioRemarks:                   handlers.FmtStringPtr(move.TIORemarks),
		FinancialReviewFlag:          move.FinancialReviewFlag,
		FinancialReviewRemarks:       move.FinancialReviewRemarks,
		CloseoutOfficeID:             handlers.FmtUUIDPtr(move.CloseoutOfficeID),
		CloseoutOffice:               TransportationOffice(move.CloseoutOffice),
		ShipmentGBLOC:                gbloc,
		LockedByOfficeUserID:         handlers.FmtUUIDPtr(move.LockedByOfficeUserID),
		LockedByOfficeUser:           OfficeUser(move.LockedByOfficeUser),
		LockExpiresAt:                handlers.FmtDateTimePtr(move.LockExpiresAt),
		AdditionalDocuments:          additionalDocumentsPayload,
		SCAssignedUser:               AssignedOfficeUser(move.SCAssignedUser),
		TOOAssignedUser:              AssignedOfficeUser(move.TOOAssignedUser),
		TIOAssignedUser:              AssignedOfficeUser(move.TIOAssignedUser),
	}

	return payload, nil
}

// ListMove payload
func ListMove(move *models.Move) *ghcmessages.ListPrimeMove {
	if move == nil {
		return nil
	}
	payload := &ghcmessages.ListPrimeMove{
		ID:                 strfmt.UUID(move.ID.String()),
		MoveCode:           move.Locator,
		CreatedAt:          strfmt.DateTime(move.CreatedAt),
		AvailableToPrimeAt: handlers.FmtDateTimePtr(move.AvailableToPrimeAt),
		ApprovedAt:         handlers.FmtDateTimePtr(move.ApprovedAt),
		OrderID:            strfmt.UUID(move.OrdersID.String()),
		ReferenceID:        *move.ReferenceID,
		UpdatedAt:          strfmt.DateTime(move.UpdatedAt),
		ETag:               etag.GenerateEtag(move.UpdatedAt),
		OrderType:          string(move.Orders.OrdersType),
	}

	if move.PPMType != nil {
		payload.PpmType = *move.PPMType
	}

	return payload
}

// ListMoves payload
func ListMoves(moves *models.Moves) []*ghcmessages.ListPrimeMove {
	listMoves := make(ghcmessages.ListPrimeMoves, len(*moves))

	for i, move := range *moves {
		// Create a local copy of the loop variable
		moveCopy := move
		listMoves[i] = ListMove(&moveCopy)
	}
	return listMoves
}

// CustomerSupportRemark payload
func CustomerSupportRemark(customerSupportRemark *models.CustomerSupportRemark) *ghcmessages.CustomerSupportRemark {
	if customerSupportRemark == nil {
		return nil
	}
	id := strfmt.UUID(customerSupportRemark.ID.String())
	moveID := strfmt.UUID(customerSupportRemark.MoveID.String())
	officeUserID := strfmt.UUID(customerSupportRemark.OfficeUserID.String())

	payload := &ghcmessages.CustomerSupportRemark{
		Content:             &customerSupportRemark.Content,
		ID:                  &id,
		CreatedAt:           strfmt.DateTime(customerSupportRemark.CreatedAt),
		UpdatedAt:           strfmt.DateTime(customerSupportRemark.UpdatedAt),
		MoveID:              &moveID,
		OfficeUserEmail:     customerSupportRemark.OfficeUser.Email,
		OfficeUserFirstName: customerSupportRemark.OfficeUser.FirstName,
		OfficeUserID:        &officeUserID,
		OfficeUserLastName:  customerSupportRemark.OfficeUser.LastName,
	}
	return payload
}

// CustomerSupportRemarks payload
func CustomerSupportRemarks(customerSupportRemarks models.CustomerSupportRemarks) ghcmessages.CustomerSupportRemarks {
	payload := make(ghcmessages.CustomerSupportRemarks, len(customerSupportRemarks))
	for i, v := range customerSupportRemarks {
		customerSupportRemark := v
		payload[i] = CustomerSupportRemark(&customerSupportRemark)
	}
	return payload
}

// EvaluationReportList payload
func EvaluationReportList(evaluationReports models.EvaluationReports) ghcmessages.EvaluationReportList {
	payload := make(ghcmessages.EvaluationReportList, len(evaluationReports))
	for i, v := range evaluationReports {
		evaluationReport := v
		payload[i] = EvaluationReport(&evaluationReport)
	}
	return payload
}

func ReportViolations(reportViolations models.ReportViolations) ghcmessages.ReportViolations {
	payload := make(ghcmessages.ReportViolations, len(reportViolations))
	for i, v := range reportViolations {
		reportViolation := v
		payload[i] = ReportViolation(&reportViolation)
	}
	return payload
}

func GsrAppeals(gsrAppeals models.GsrAppeals) ghcmessages.GSRAppeals {
	payload := make(ghcmessages.GSRAppeals, len(gsrAppeals))
	for i, v := range gsrAppeals {
		gsrAppeal := v
		payload[i] = GsrAppeal(&gsrAppeal)
	}
	return payload
}

func EvaluationReportOfficeUser(officeUser models.OfficeUser) ghcmessages.EvaluationReportOfficeUser {
	payload := ghcmessages.EvaluationReportOfficeUser{
		Email:     officeUser.Email,
		FirstName: officeUser.FirstName,
		ID:        strfmt.UUID(officeUser.ID.String()),
		LastName:  officeUser.LastName,
		Phone:     officeUser.Telephone,
	}
	return payload
}

// EvaluationReport payload
func EvaluationReport(evaluationReport *models.EvaluationReport) *ghcmessages.EvaluationReport {
	if evaluationReport == nil {
		return nil
	}
	id := *handlers.FmtUUID(evaluationReport.ID)
	moveID := *handlers.FmtUUID(evaluationReport.MoveID)
	shipmentID := handlers.FmtUUIDPtr(evaluationReport.ShipmentID)

	var inspectionType *ghcmessages.EvaluationReportInspectionType
	if evaluationReport.InspectionType != nil {
		tempInspectionType := ghcmessages.EvaluationReportInspectionType(*evaluationReport.InspectionType)
		inspectionType = &tempInspectionType
	}
	var location *ghcmessages.EvaluationReportLocation
	if evaluationReport.Location != nil {
		tempLocation := ghcmessages.EvaluationReportLocation(*evaluationReport.Location)
		location = &tempLocation
	}
	reportType := ghcmessages.EvaluationReportType(evaluationReport.Type)

	evaluationReportOfficeUserPayload := EvaluationReportOfficeUser(evaluationReport.OfficeUser)

	var timeDepart *string
	if evaluationReport.TimeDepart != nil {
		td := evaluationReport.TimeDepart.Format(timeHHMMFormat)
		timeDepart = &td
	}

	var evalStart *string
	if evaluationReport.EvalStart != nil {
		es := evaluationReport.EvalStart.Format(timeHHMMFormat)
		evalStart = &es
	}

	var evalEnd *string
	if evaluationReport.EvalEnd != nil {
		ee := evaluationReport.EvalEnd.Format(timeHHMMFormat)
		evalEnd = &ee
	}

	payload := &ghcmessages.EvaluationReport{
		CreatedAt:                          strfmt.DateTime(evaluationReport.CreatedAt),
		ID:                                 id,
		InspectionDate:                     handlers.FmtDatePtr(evaluationReport.InspectionDate),
		InspectionType:                     inspectionType,
		Location:                           location,
		LocationDescription:                evaluationReport.LocationDescription,
		MoveID:                             moveID,
		ObservedShipmentPhysicalPickupDate: handlers.FmtDatePtr(evaluationReport.ObservedShipmentPhysicalPickupDate),
		ObservedShipmentDeliveryDate:       handlers.FmtDatePtr(evaluationReport.ObservedShipmentDeliveryDate),
		Remarks:                            evaluationReport.Remarks,
		ShipmentID:                         shipmentID,
		SubmittedAt:                        handlers.FmtDateTimePtr(evaluationReport.SubmittedAt),
		TimeDepart:                         timeDepart,
		EvalStart:                          evalStart,
		EvalEnd:                            evalEnd,
		Type:                               reportType,
		ViolationsObserved:                 evaluationReport.ViolationsObserved,
		MoveReferenceID:                    evaluationReport.Move.ReferenceID,
		OfficeUser:                         &evaluationReportOfficeUserPayload,
		SeriousIncident:                    evaluationReport.SeriousIncident,
		SeriousIncidentDesc:                evaluationReport.SeriousIncidentDesc,
		ObservedClaimsResponseDate:         handlers.FmtDatePtr(evaluationReport.ObservedClaimsResponseDate),
		ObservedPickupDate:                 handlers.FmtDatePtr(evaluationReport.ObservedPickupDate),
		ObservedPickupSpreadStartDate:      handlers.FmtDatePtr(evaluationReport.ObservedPickupSpreadStartDate),
		ObservedPickupSpreadEndDate:        handlers.FmtDatePtr(evaluationReport.ObservedPickupSpreadEndDate),
		ObservedDeliveryDate:               handlers.FmtDatePtr(evaluationReport.ObservedDeliveryDate),
		ETag:                               etag.GenerateEtag(evaluationReport.UpdatedAt),
		UpdatedAt:                          strfmt.DateTime(evaluationReport.UpdatedAt),
		ReportViolations:                   ReportViolations(evaluationReport.ReportViolations),
		GsrAppeals:                         GsrAppeals(evaluationReport.GsrAppeals),
	}
	return payload
}

// PWSViolationItem payload
func PWSViolationItem(violation *models.PWSViolation) *ghcmessages.PWSViolation {
	if violation == nil {
		return nil
	}

	payload := &ghcmessages.PWSViolation{
		ID:                   strfmt.UUID(violation.ID.String()),
		DisplayOrder:         int64(violation.DisplayOrder),
		ParagraphNumber:      violation.ParagraphNumber,
		Title:                violation.Title,
		Category:             string(violation.Category),
		SubCategory:          violation.SubCategory,
		RequirementSummary:   violation.RequirementSummary,
		RequirementStatement: violation.RequirementStatement,
		IsKpi:                violation.IsKpi,
		AdditionalDataElem:   violation.AdditionalDataElem,
	}

	return payload
}

// PWSViolations payload
func PWSViolations(violations models.PWSViolations) ghcmessages.PWSViolations {
	payload := make(ghcmessages.PWSViolations, len(violations))

	for i, v := range violations {
		violation := v
		payload[i] = PWSViolationItem(&violation)
	}
	return payload
}

func ReportViolation(reportViolation *models.ReportViolation) *ghcmessages.ReportViolation {
	if reportViolation == nil {
		return nil
	}
	id := *handlers.FmtUUID(reportViolation.ID)
	violationID := *handlers.FmtUUID(reportViolation.ViolationID)
	reportID := *handlers.FmtUUID(reportViolation.ReportID)

	payload := &ghcmessages.ReportViolation{
		ID:          id,
		ViolationID: violationID,
		ReportID:    reportID,
		Violation:   PWSViolationItem(&reportViolation.Violation),
		GsrAppeals:  GsrAppeals(reportViolation.GsrAppeals),
	}
	return payload
}

func GsrAppeal(gsrAppeal *models.GsrAppeal) *ghcmessages.GSRAppeal {
	if gsrAppeal == nil {
		return nil
	}
	id := *handlers.FmtUUID(gsrAppeal.ID)
	reportID := *handlers.FmtUUID(gsrAppeal.EvaluationReportID)
	officeUserID := *handlers.FmtUUID(gsrAppeal.OfficeUserID)
	officeUser := EvaluationReportOfficeUser(*gsrAppeal.OfficeUser)
	isSeriousIncident := false
	if gsrAppeal.IsSeriousIncidentAppeal != nil {
		isSeriousIncident = *gsrAppeal.IsSeriousIncidentAppeal
	}

	payload := &ghcmessages.GSRAppeal{
		ID:                id,
		ReportID:          reportID,
		OfficeUserID:      officeUserID,
		OfficeUser:        &officeUser,
		IsSeriousIncident: isSeriousIncident,
		AppealStatus:      ghcmessages.GSRAppealStatusType(gsrAppeal.AppealStatus),
		Remarks:           gsrAppeal.Remarks,
		CreatedAt:         strfmt.DateTime(gsrAppeal.CreatedAt),
	}

	if gsrAppeal.ReportViolationID != nil {
		payload.ViolationID = *handlers.FmtUUID(*gsrAppeal.ReportViolationID)
	}
	return payload
}

// TransportationOffice payload
func TransportationOffice(office *models.TransportationOffice) *ghcmessages.TransportationOffice {
	if office == nil || office.ID == uuid.Nil {
		return nil
	}

	phoneLines := []string{}
	for _, phoneLine := range office.PhoneLines {
		if phoneLine.Type == "voice" {
			phoneLines = append(phoneLines, phoneLine.Number)
		}
	}

	payload := &ghcmessages.TransportationOffice{
		ID:         handlers.FmtUUID(office.ID),
		CreatedAt:  handlers.FmtDateTime(office.CreatedAt),
		UpdatedAt:  handlers.FmtDateTime(office.UpdatedAt),
		Name:       models.StringPointer(office.Name),
		Gbloc:      office.Gbloc,
		Address:    Address(&office.Address),
		PhoneLines: phoneLines,
	}
	return payload
}

func TransportationOffices(transportationOffices models.TransportationOffices) ghcmessages.TransportationOffices {
	payload := make(ghcmessages.TransportationOffices, len(transportationOffices))

	for i, to := range transportationOffices {
		transportationOffice := to
		payload[i] = TransportationOffice(&transportationOffice)
	}
	return payload
}

func GBLOCs(gblocs []string) ghcmessages.GBLOCs {
	payload := make(ghcmessages.GBLOCs, len(gblocs))

	for i, gbloc := range gblocs {
		payload[i] = string(gbloc)
	}
	return payload
}

// MoveHistory payload
func MoveHistory(logger *zap.Logger, moveHistory *models.MoveHistory) *ghcmessages.MoveHistory {
	payload := &ghcmessages.MoveHistory{
		HistoryRecords: moveHistoryRecords(logger, moveHistory.AuditHistories),
		ID:             strfmt.UUID(moveHistory.ID.String()),
		Locator:        moveHistory.Locator,
		ReferenceID:    moveHistory.ReferenceID,
	}

	return payload
}

// MoveAuditHistory payload
func MoveAuditHistory(logger *zap.Logger, auditHistory models.AuditHistory) *ghcmessages.MoveAuditHistory {

	payload := &ghcmessages.MoveAuditHistory{
		Action:               auditHistory.Action,
		ActionTstampClk:      strfmt.DateTime(auditHistory.ActionTstampClk),
		ActionTstampStm:      strfmt.DateTime(auditHistory.ActionTstampStm),
		ActionTstampTx:       strfmt.DateTime(auditHistory.ActionTstampTx),
		ChangedValues:        removeEscapeJSONtoObject(logger, auditHistory.ChangedData),
		OldValues:            removeEscapeJSONtoObject(logger, auditHistory.OldData),
		EventName:            auditHistory.EventName,
		ID:                   strfmt.UUID(auditHistory.ID.String()),
		ObjectID:             handlers.FmtUUIDPtr(auditHistory.ObjectID),
		RelID:                auditHistory.RelID,
		SessionUserID:        handlers.FmtUUIDPtr(auditHistory.SessionUserID),
		SessionUserFirstName: auditHistory.SessionUserFirstName,
		SessionUserLastName:  auditHistory.SessionUserLastName,
		SessionUserEmail:     auditHistory.SessionUserEmail,
		SessionUserTelephone: auditHistory.SessionUserTelephone,
		Context:              removeEscapeJSONtoArray(logger, auditHistory.Context),
		ContextID:            auditHistory.ContextID,
		StatementOnly:        auditHistory.StatementOnly,
		TableName:            auditHistory.AuditedTable,
		SchemaName:           auditHistory.SchemaName,
		TransactionID:        auditHistory.TransactionID,
	}

	return payload
}

func removeEscapeJSONtoObject(logger *zap.Logger, data *string) map[string]interface{} {
	var result map[string]interface{}
	if data == nil || *data == "" {
		return result
	}
	var byteData = []byte(*data)

	err := json.Unmarshal(byteData, &result)

	if err != nil {
		logger.Error("error unmarshalling the escaped json to object", zap.Error(err))
	}

	return result

}

func removeEscapeJSONtoArray(logger *zap.Logger, data *string) []map[string]string {
	var result []map[string]string
	if data == nil || *data == "" {
		return result
	}
	var byteData = []byte(*data)

	err := json.Unmarshal(byteData, &result)

	if err != nil {
		logger.Error("error unmarshalling the escaped json to array", zap.Error(err))
	}

	return result
}

func moveHistoryRecords(logger *zap.Logger, auditHistories models.AuditHistories) ghcmessages.MoveAuditHistories {
	payload := make(ghcmessages.MoveAuditHistories, len(auditHistories))

	for i, a := range auditHistories {
		payload[i] = MoveAuditHistory(logger, a)
	}
	return payload
}

// MoveTaskOrder payload
func MoveTaskOrder(moveTaskOrder *models.Move) *ghcmessages.MoveTaskOrder {
	if moveTaskOrder == nil {
		return nil
	}

	payload := &ghcmessages.MoveTaskOrder{
		ID:                 strfmt.UUID(moveTaskOrder.ID.String()),
		CreatedAt:          strfmt.DateTime(moveTaskOrder.CreatedAt),
		AvailableToPrimeAt: handlers.FmtDateTimePtr(moveTaskOrder.AvailableToPrimeAt),
		ApprovedAt:         handlers.FmtDateTimePtr(moveTaskOrder.ApprovedAt),
		OrderID:            strfmt.UUID(moveTaskOrder.OrdersID.String()),
		ReferenceID:        *moveTaskOrder.ReferenceID,
		UpdatedAt:          strfmt.DateTime(moveTaskOrder.UpdatedAt),
		ETag:               etag.GenerateEtag(moveTaskOrder.UpdatedAt),
		Locator:            moveTaskOrder.Locator,
	}
	return payload
}

// Customer payload
func Customer(customer *models.ServiceMember) *ghcmessages.Customer {
	if customer == nil {
		return nil
	}

	payload := ghcmessages.Customer{
		Agency:             swag.StringValue((*string)(customer.Affiliation)),
		CurrentAddress:     Address(customer.ResidentialAddress),
		Edipi:              swag.StringValue(customer.Edipi),
		Email:              customer.PersonalEmail,
		FirstName:          swag.StringValue(customer.FirstName),
		ID:                 strfmt.UUID(customer.ID.String()),
		LastName:           swag.StringValue(customer.LastName),
		Phone:              customer.Telephone,
		Suffix:             customer.Suffix,
		MiddleName:         customer.MiddleName,
		UserID:             strfmt.UUID(customer.UserID.String()),
		ETag:               etag.GenerateEtag(customer.UpdatedAt),
		BackupContact:      BackupContact(customer.BackupContacts),
		BackupAddress:      Address(customer.BackupMailingAddress),
		SecondaryTelephone: customer.SecondaryTelephone,
		PhoneIsPreferred:   swag.BoolValue(customer.PhoneIsPreferred),
		EmailIsPreferred:   swag.BoolValue(customer.EmailIsPreferred),
		CacValidated:       &customer.CacValidated,
		Emplid:             customer.Emplid,
	}
	return &payload
}

func CreatedCustomer(sm *models.ServiceMember, oktaUser *models.CreatedOktaUser, backupContact *models.BackupContact) *ghcmessages.CreatedCustomer {
	if sm == nil || oktaUser == nil || backupContact == nil {
		return nil
	}

	bc := &ghcmessages.BackupContact{
		Name:  &backupContact.Name,
		Email: &backupContact.Email,
		Phone: backupContact.Phone,
	}

	payload := ghcmessages.CreatedCustomer{
		ID:                 strfmt.UUID(sm.ID.String()),
		UserID:             strfmt.UUID(sm.UserID.String()),
		OktaID:             oktaUser.ID,
		OktaEmail:          oktaUser.Profile.Email,
		Affiliation:        swag.StringValue((*string)(sm.Affiliation)),
		Edipi:              sm.Edipi,
		FirstName:          swag.StringValue(sm.FirstName),
		MiddleName:         sm.MiddleName,
		LastName:           swag.StringValue(sm.LastName),
		Suffix:             sm.Suffix,
		ResidentialAddress: Address(sm.ResidentialAddress),
		BackupAddress:      Address(sm.BackupMailingAddress),
		PersonalEmail:      *sm.PersonalEmail,
		Telephone:          sm.Telephone,
		SecondaryTelephone: sm.SecondaryTelephone,
		PhoneIsPreferred:   swag.BoolValue(sm.PhoneIsPreferred),
		EmailIsPreferred:   swag.BoolValue(sm.EmailIsPreferred),
		BackupContact:      bc,
		CacValidated:       swag.BoolValue(&sm.CacValidated),
	}
	return &payload
}

// Order payload
func Order(order *models.Order) *ghcmessages.Order {
	if order == nil {
		return nil
	}
	if order.ID == uuid.Nil {
		return nil
	}

	destinationDutyLocation := DutyLocation(&order.NewDutyLocation)
	originDutyLocation := DutyLocation(order.OriginDutyLocation)
	if order.Grade != nil && order.Entitlement != nil {
		order.Entitlement.SetWeightAllotment(string(*order.Grade))
	}
	entitlements := Entitlement(order.Entitlement)

	var deptIndicator ghcmessages.DeptIndicator
	if order.DepartmentIndicator != nil {
		deptIndicator = ghcmessages.DeptIndicator(*order.DepartmentIndicator)
	}

	var ordersTypeDetail ghcmessages.OrdersTypeDetail
	if order.OrdersTypeDetail != nil {
		ordersTypeDetail = ghcmessages.OrdersTypeDetail(*order.OrdersTypeDetail)
	}

	var grade ghcmessages.Grade
	if order.Grade != nil {
		grade = ghcmessages.Grade(*order.Grade)
	}
	//
	var affiliation ghcmessages.Affiliation
	if order.ServiceMember.Affiliation != nil {
		affiliation = ghcmessages.Affiliation(*order.ServiceMember.Affiliation)
	}

	var moveCode string
	var moveTaskOrderID strfmt.UUID
	if len(order.Moves) > 0 {
		moveCode = order.Moves[0].Locator
		moveTaskOrderID = strfmt.UUID(order.Moves[0].ID.String())
	}

	payload := ghcmessages.Order{
		DestinationDutyLocation:        destinationDutyLocation,
		DestinationDutyLocationGBLOC:   ghcmessages.GBLOC(swag.StringValue(order.DestinationGBLOC)),
		Entitlement:                    entitlements,
		Grade:                          &grade,
		OrderNumber:                    order.OrdersNumber,
		OrderTypeDetail:                &ordersTypeDetail,
		ID:                             strfmt.UUID(order.ID.String()),
		OriginDutyLocation:             originDutyLocation,
		ETag:                           etag.GenerateEtag(order.UpdatedAt),
		Agency:                         &affiliation,
		CustomerID:                     strfmt.UUID(order.ServiceMemberID.String()),
		Customer:                       Customer(&order.ServiceMember),
		FirstName:                      swag.StringValue(order.ServiceMember.FirstName),
		LastName:                       swag.StringValue(order.ServiceMember.LastName),
		ReportByDate:                   strfmt.Date(order.ReportByDate),
		DateIssued:                     strfmt.Date(order.IssueDate),
		OrderType:                      ghcmessages.OrdersType(order.OrdersType),
		DepartmentIndicator:            &deptIndicator,
		Tac:                            handlers.FmtStringPtr(order.TAC),
		Sac:                            handlers.FmtStringPtr(order.SAC),
		NtsTac:                         handlers.FmtStringPtr(order.NtsTAC),
		NtsSac:                         handlers.FmtStringPtr(order.NtsSAC),
		SupplyAndServicesCostEstimate:  order.SupplyAndServicesCostEstimate,
		PackingAndShippingInstructions: order.PackingAndShippingInstructions,
		MethodOfPayment:                order.MethodOfPayment,
		Naics:                          order.NAICS,
		UploadedOrderID:                strfmt.UUID(order.UploadedOrdersID.String()),
		UploadedAmendedOrderID:         handlers.FmtUUIDPtr(order.UploadedAmendedOrdersID),
		AmendedOrdersAcknowledgedAt:    handlers.FmtDateTimePtr(order.AmendedOrdersAcknowledgedAt),
		MoveCode:                       moveCode,
		MoveTaskOrderID:                moveTaskOrderID,
		OriginDutyLocationGBLOC:        ghcmessages.GBLOC(swag.StringValue(order.OriginDutyLocationGBLOC)),
	}

	return &payload
}

// Entitlement payload
func Entitlement(entitlement *models.Entitlement) *ghcmessages.Entitlements {
	if entitlement == nil {
		return nil
	}
	var proGearWeight, proGearWeightSpouse, totalWeight int64
	proGearWeight = int64(entitlement.ProGearWeight)
	proGearWeightSpouse = int64(entitlement.ProGearWeightSpouse)

	if weightAllotment := entitlement.WeightAllotment(); weightAllotment != nil {
		if *entitlement.DependentsAuthorized {
			totalWeight = int64(weightAllotment.TotalWeightSelfPlusDependents)
		} else {
			totalWeight = int64(weightAllotment.TotalWeightSelf)
		}
	}
	var authorizedWeight *int64
	if entitlement.AuthorizedWeight() != nil {
		aw := int64(*entitlement.AuthorizedWeight())
		authorizedWeight = &aw
	}
	var sit *int64
	if entitlement.StorageInTransit != nil {
		sitValue := int64(*entitlement.StorageInTransit)
		sit = &sitValue
	}
	var totalDependents int64
	if entitlement.TotalDependents != nil {
		totalDependents = int64(*entitlement.TotalDependents)
	}
	requiredMedicalEquipmentWeight := int64(entitlement.RequiredMedicalEquipmentWeight)
	gunSafe := entitlement.GunSafe
	return &ghcmessages.Entitlements{
		ID:                             strfmt.UUID(entitlement.ID.String()),
		AuthorizedWeight:               authorizedWeight,
		DependentsAuthorized:           entitlement.DependentsAuthorized,
		NonTemporaryStorage:            entitlement.NonTemporaryStorage,
		PrivatelyOwnedVehicle:          entitlement.PrivatelyOwnedVehicle,
		ProGearWeight:                  proGearWeight,
		ProGearWeightSpouse:            proGearWeightSpouse,
		StorageInTransit:               sit,
		TotalDependents:                totalDependents,
		TotalWeight:                    totalWeight,
		RequiredMedicalEquipmentWeight: requiredMedicalEquipmentWeight,
		OrganizationalClothingAndIndividualEquipment: entitlement.OrganizationalClothingAndIndividualEquipment,
		GunSafe: gunSafe,
		ETag:    etag.GenerateEtag(entitlement.UpdatedAt),
	}
}

// DutyLocation payload
func DutyLocation(dutyLocation *models.DutyLocation) *ghcmessages.DutyLocation {
	if dutyLocation == nil {
		return nil
	}
	address := Address(&dutyLocation.Address)
	payload := ghcmessages.DutyLocation{
		Address:   address,
		AddressID: address.ID,
		ID:        strfmt.UUID(dutyLocation.ID.String()),
		Name:      dutyLocation.Name,
		ETag:      etag.GenerateEtag(dutyLocation.UpdatedAt),
	}
	return &payload
}

// Country payload
func Country(country *models.Country) *string {
	if country == nil {
		return nil
	}
	return &country.Country
}

// Address payload
func Address(address *models.Address) *ghcmessages.Address {
	if address == nil {
		return nil
	}
	return &ghcmessages.Address{
		ID:             strfmt.UUID(address.ID.String()),
		StreetAddress1: &address.StreetAddress1,
		StreetAddress2: address.StreetAddress2,
		StreetAddress3: address.StreetAddress3,
		City:           &address.City,
		State:          &address.State,
		PostalCode:     &address.PostalCode,
		Country:        Country(address.Country),
		County:         &address.County,
		ETag:           etag.GenerateEtag(address.UpdatedAt),
		IsOconus:       address.IsOconus,
	}
}

// PPM destination Address payload
func PPMDestinationAddress(address *models.Address) *ghcmessages.Address {
	payload := Address(address)

	if payload == nil {
		return nil
	}

	// Street address 1 is optional per business rule but not nullable on the database level.
	// Check if streetAddress 1 is using place holder value to represent 'NULL'.
	// If so return empty string.
	if strings.EqualFold(*payload.StreetAddress1, models.STREET_ADDRESS_1_NOT_PROVIDED) {
		payload.StreetAddress1 = models.StringPointer("")
	}
	return payload
}

// StorageFacility payload
func StorageFacility(storageFacility *models.StorageFacility) *ghcmessages.StorageFacility {
	if storageFacility == nil {
		return nil
	}

	payload := ghcmessages.StorageFacility{
		ID:           strfmt.UUID(storageFacility.ID.String()),
		FacilityName: storageFacility.FacilityName,
		Address:      Address(&storageFacility.Address),
		LotNumber:    storageFacility.LotNumber,
		Phone:        storageFacility.Phone,
		Email:        storageFacility.Email,
		ETag:         etag.GenerateEtag(storageFacility.UpdatedAt),
	}

	return &payload
}

// BackupContact payload
func BackupContact(contacts models.BackupContacts) *ghcmessages.BackupContact {
	if len(contacts) == 0 {
		return nil
	}
	var name, email, phone string

	if len(contacts) != 0 {
		contact := contacts[0]
		name = contact.Name
		email = contact.Email
		phone = ""
		contactPhone := contact.Phone
		if contactPhone != nil {
			phone = *contactPhone
		}
	}

	return &ghcmessages.BackupContact{
		Name:  &name,
		Email: &email,
		Phone: &phone,
	}
}

// SITDurationUpdate payload
func SITDurationUpdate(sitDurationUpdate *models.SITDurationUpdate) *ghcmessages.SITExtension {
	if sitDurationUpdate == nil {
		return nil
	}
	payload := &ghcmessages.SITExtension{
		ID:                strfmt.UUID(sitDurationUpdate.ID.String()),
		ETag:              etag.GenerateEtag(sitDurationUpdate.UpdatedAt),
		MtoShipmentID:     strfmt.UUID(sitDurationUpdate.MTOShipmentID.String()),
		RequestReason:     string(sitDurationUpdate.RequestReason),
		RequestedDays:     int64(sitDurationUpdate.RequestedDays),
		Status:            string(sitDurationUpdate.Status),
		CreatedAt:         strfmt.DateTime(sitDurationUpdate.CreatedAt),
		UpdatedAt:         strfmt.DateTime(sitDurationUpdate.UpdatedAt),
		ApprovedDays:      handlers.FmtIntPtrToInt64(sitDurationUpdate.ApprovedDays),
		ContractorRemarks: handlers.FmtStringPtr(sitDurationUpdate.ContractorRemarks),
		DecisionDate:      handlers.FmtDateTimePtr(sitDurationUpdate.DecisionDate),
		OfficeRemarks:     handlers.FmtStringPtr(sitDurationUpdate.OfficeRemarks),
	}

	return payload
}

// SITDurationUpdates payload
func SITDurationUpdates(sitDurationUpdates *models.SITDurationUpdates) *ghcmessages.SITExtensions {
	payload := make(ghcmessages.SITExtensions, len(*sitDurationUpdates))

	if len(*sitDurationUpdates) > 0 {
		for i, m := range *sitDurationUpdates {
			copyOfSITDurationUpdate := m // Make copy to avoid implicit memory aliasing of items from a range statement.
			payload[i] = SITDurationUpdate(&copyOfSITDurationUpdate)
		}
		// Reversing the SIT duration updates as they are saved in the order
		// they are created and we want to always display them in the reverse
		// order.
		for i, j := 0, len(payload)-1; i < j; i, j = i+1, j-1 {
			payload[i], payload[j] = payload[j], payload[i]
		}
	}
	return &payload
}

func currentSIT(currentSIT *services.CurrentSIT) *ghcmessages.SITStatusCurrentSIT {
	if currentSIT == nil {
		return nil
	}
	return &ghcmessages.SITStatusCurrentSIT{
		ServiceItemID:        *handlers.FmtUUID(currentSIT.ServiceItemID), // TODO: Refactor out service item ID dependence in GHC API. This should be based on SIT groupings / summaries
		Location:             currentSIT.Location,
		DaysInSIT:            handlers.FmtIntPtrToInt64(&currentSIT.DaysInSIT),
		SitEntryDate:         handlers.FmtDate(currentSIT.SITEntryDate),
		SitDepartureDate:     handlers.FmtDatePtr(currentSIT.SITDepartureDate),
		SitAuthorizedEndDate: handlers.FmtDate(currentSIT.SITAuthorizedEndDate),
		SitCustomerContacted: handlers.FmtDatePtr(currentSIT.SITCustomerContacted),
		SitRequestedDelivery: handlers.FmtDatePtr(currentSIT.SITRequestedDelivery),
	}
}

// SITStatus payload
func SITStatus(shipmentSITStatuses *services.SITStatus, storer storage.FileStorer) *ghcmessages.SITStatus {
	if shipmentSITStatuses == nil {
		return nil
	}

	payload := &ghcmessages.SITStatus{
		PastSITServiceItemGroupings: SITServiceItemGroupings(shipmentSITStatuses.PastSITs, storer),
		TotalSITDaysUsed:            handlers.FmtIntPtrToInt64(&shipmentSITStatuses.TotalSITDaysUsed),
		TotalDaysRemaining:          handlers.FmtIntPtrToInt64(&shipmentSITStatuses.TotalDaysRemaining),
		CalculatedTotalDaysInSIT:    handlers.FmtIntPtrToInt64(&shipmentSITStatuses.CalculatedTotalDaysInSIT),
		CurrentSIT:                  currentSIT(shipmentSITStatuses.CurrentSIT),
	}

	return payload
}

// SITStatuses payload
func SITStatuses(shipmentSITStatuses map[string]services.SITStatus, storer storage.FileStorer) map[string]*ghcmessages.SITStatus {
	sitStatuses := map[string]*ghcmessages.SITStatus{}
	if len(shipmentSITStatuses) == 0 {
		return sitStatuses
	}

	for _, sitStatus := range shipmentSITStatuses {
		copyOfSITStatus := sitStatus
		sitStatuses[sitStatus.ShipmentID.String()] = SITStatus(&copyOfSITStatus, storer)
	}

	return sitStatuses
}

// PPMShipment payload
func PPMShipment(_ storage.FileStorer, ppmShipment *models.PPMShipment) *ghcmessages.PPMShipment {
	if ppmShipment == nil || ppmShipment.ID.IsNil() {
		return nil
	}

	payloadPPMShipment := &ghcmessages.PPMShipment{
		ID:                             *handlers.FmtUUID(ppmShipment.ID),
		ShipmentID:                     *handlers.FmtUUID(ppmShipment.ShipmentID),
		CreatedAt:                      strfmt.DateTime(ppmShipment.CreatedAt),
		UpdatedAt:                      strfmt.DateTime(ppmShipment.UpdatedAt),
		Status:                         ghcmessages.PPMShipmentStatus(ppmShipment.Status),
		ExpectedDepartureDate:          handlers.FmtDate(ppmShipment.ExpectedDepartureDate),
		ActualMoveDate:                 handlers.FmtDatePtr(ppmShipment.ActualMoveDate),
		SubmittedAt:                    handlers.FmtDateTimePtr(ppmShipment.SubmittedAt),
		ReviewedAt:                     handlers.FmtDateTimePtr(ppmShipment.ReviewedAt),
		ApprovedAt:                     handlers.FmtDateTimePtr(ppmShipment.ApprovedAt),
		PickupAddress:                  Address(ppmShipment.PickupAddress),
		DestinationAddress:             PPMDestinationAddress(ppmShipment.DestinationAddress),
		ActualPickupPostalCode:         ppmShipment.ActualPickupPostalCode,
		ActualDestinationPostalCode:    ppmShipment.ActualDestinationPostalCode,
		SitExpected:                    ppmShipment.SITExpected,
		HasSecondaryPickupAddress:      ppmShipment.HasSecondaryPickupAddress,
		HasSecondaryDestinationAddress: ppmShipment.HasSecondaryDestinationAddress,
		HasTertiaryPickupAddress:       ppmShipment.HasTertiaryPickupAddress,
		HasTertiaryDestinationAddress:  ppmShipment.HasTertiaryDestinationAddress,
		EstimatedWeight:                handlers.FmtPoundPtr(ppmShipment.EstimatedWeight),
		HasProGear:                     ppmShipment.HasProGear,
		ProGearWeight:                  handlers.FmtPoundPtr(ppmShipment.ProGearWeight),
		SpouseProGearWeight:            handlers.FmtPoundPtr(ppmShipment.SpouseProGearWeight),
		EstimatedIncentive:             handlers.FmtCost(ppmShipment.EstimatedIncentive),
		HasRequestedAdvance:            ppmShipment.HasRequestedAdvance,
		AdvanceAmountRequested:         handlers.FmtCost(ppmShipment.AdvanceAmountRequested),
		HasReceivedAdvance:             ppmShipment.HasReceivedAdvance,
		AdvanceAmountReceived:          handlers.FmtCost(ppmShipment.AdvanceAmountReceived),
		SitEstimatedWeight:             handlers.FmtPoundPtr(ppmShipment.SITEstimatedWeight),
		SitEstimatedEntryDate:          handlers.FmtDatePtr(ppmShipment.SITEstimatedEntryDate),
		SitEstimatedDepartureDate:      handlers.FmtDatePtr(ppmShipment.SITEstimatedDepartureDate),
		SitEstimatedCost:               handlers.FmtCost(ppmShipment.SITEstimatedCost),
		IsActualExpenseReimbursement:   ppmShipment.IsActualExpenseReimbursement,
		ETag:                           etag.GenerateEtag(ppmShipment.UpdatedAt),
	}

	if ppmShipment.SITLocation != nil {
		sitLocation := ghcmessages.SITLocationType(*ppmShipment.SITLocation)
		payloadPPMShipment.SitLocation = &sitLocation
	}

	if ppmShipment.AdvanceStatus != nil {
		advanceStatus := ghcmessages.PPMAdvanceStatus(*ppmShipment.AdvanceStatus)
		payloadPPMShipment.AdvanceStatus = &advanceStatus
	}

	if ppmShipment.W2Address != nil {
		payloadPPMShipment.W2Address = Address(ppmShipment.W2Address)
	}

	if ppmShipment.SecondaryPickupAddress != nil {
		payloadPPMShipment.SecondaryPickupAddress = Address(ppmShipment.SecondaryPickupAddress)
	}

	if ppmShipment.SecondaryDestinationAddress != nil {
		payloadPPMShipment.SecondaryDestinationAddress = Address(ppmShipment.SecondaryDestinationAddress)
	}

	if ppmShipment.TertiaryPickupAddress != nil {
		payloadPPMShipment.TertiaryPickupAddress = Address(ppmShipment.TertiaryPickupAddress)
	}

	if ppmShipment.TertiaryDestinationAddress != nil {
		payloadPPMShipment.TertiaryDestinationAddress = Address(ppmShipment.TertiaryDestinationAddress)
	}

	if ppmShipment.IsActualExpenseReimbursement != nil {
		payloadPPMShipment.IsActualExpenseReimbursement = ppmShipment.IsActualExpenseReimbursement
	}

	return payloadPPMShipment
}

// BoatShipment payload
func BoatShipment(storer storage.FileStorer, boatShipment *models.BoatShipment) *ghcmessages.BoatShipment {
	if boatShipment == nil || boatShipment.ID.IsNil() {
		return nil
	}

	payloadBoatShipment := &ghcmessages.BoatShipment{
		ID:             *handlers.FmtUUID(boatShipment.ID),
		ShipmentID:     *handlers.FmtUUID(boatShipment.ShipmentID),
		CreatedAt:      strfmt.DateTime(boatShipment.CreatedAt),
		UpdatedAt:      strfmt.DateTime(boatShipment.UpdatedAt),
		Type:           models.StringPointer(string(boatShipment.Type)),
		Year:           handlers.FmtIntPtrToInt64(boatShipment.Year),
		Make:           boatShipment.Make,
		Model:          boatShipment.Model,
		LengthInInches: handlers.FmtIntPtrToInt64(boatShipment.LengthInInches),
		WidthInInches:  handlers.FmtIntPtrToInt64(boatShipment.WidthInInches),
		HeightInInches: handlers.FmtIntPtrToInt64(boatShipment.HeightInInches),
		HasTrailer:     boatShipment.HasTrailer,
		IsRoadworthy:   boatShipment.IsRoadworthy,
		ETag:           etag.GenerateEtag(boatShipment.UpdatedAt),
	}

	return payloadBoatShipment
}

// MobileHomeShipment payload
func MobileHomeShipment(storer storage.FileStorer, mobileHomeShipment *models.MobileHome) *ghcmessages.MobileHome {
	if mobileHomeShipment == nil || mobileHomeShipment.ID.IsNil() {
		return nil
	}

	payloadMobileHomeShipment := &ghcmessages.MobileHome{
		ID:             *handlers.FmtUUID(mobileHomeShipment.ID),
		ShipmentID:     *handlers.FmtUUID(mobileHomeShipment.ShipmentID),
		Make:           *mobileHomeShipment.Make,
		Model:          *mobileHomeShipment.Model,
		Year:           *handlers.FmtIntPtrToInt64(mobileHomeShipment.Year),
		LengthInInches: *handlers.FmtIntPtrToInt64(mobileHomeShipment.LengthInInches),
		HeightInInches: *handlers.FmtIntPtrToInt64(mobileHomeShipment.HeightInInches),
		WidthInInches:  *handlers.FmtIntPtrToInt64(mobileHomeShipment.WidthInInches),
		CreatedAt:      strfmt.DateTime(mobileHomeShipment.CreatedAt),
		UpdatedAt:      strfmt.DateTime(mobileHomeShipment.UpdatedAt),
		ETag:           etag.GenerateEtag(mobileHomeShipment.UpdatedAt),
	}

	return payloadMobileHomeShipment
}

// ProGearWeightTickets sets up a ProGearWeightTicket slice for the api using model data.
func ProGearWeightTickets(storer storage.FileStorer, proGearWeightTickets models.ProgearWeightTickets) []*ghcmessages.ProGearWeightTicket {
	payload := make([]*ghcmessages.ProGearWeightTicket, len(proGearWeightTickets))
	for i, proGearWeightTicket := range proGearWeightTickets {
		copyOfProGearWeightTicket := proGearWeightTicket
		proGearWeightTicketPayload := ProGearWeightTicket(storer, &copyOfProGearWeightTicket)
		payload[i] = proGearWeightTicketPayload
	}
	return payload
}

// ProGearWeightTicket payload
func ProGearWeightTicket(storer storage.FileStorer, progear *models.ProgearWeightTicket) *ghcmessages.ProGearWeightTicket {
	ppmShipmentID := strfmt.UUID(progear.PPMShipmentID.String())

	document, err := PayloadForDocumentModel(storer, progear.Document)
	if err != nil {
		return nil
	}

	payload := &ghcmessages.ProGearWeightTicket{
		ID:               strfmt.UUID(progear.ID.String()),
		PpmShipmentID:    ppmShipmentID,
		CreatedAt:        *handlers.FmtDateTime(progear.CreatedAt),
		UpdatedAt:        *handlers.FmtDateTime(progear.UpdatedAt),
		DocumentID:       *handlers.FmtUUID(progear.DocumentID),
		Document:         document,
		Weight:           handlers.FmtPoundPtr(progear.Weight),
		BelongsToSelf:    progear.BelongsToSelf,
		HasWeightTickets: progear.HasWeightTickets,
		Description:      progear.Description,
		ETag:             etag.GenerateEtag(progear.UpdatedAt),
	}

	if progear.Status != nil {
		status := ghcmessages.OmittablePPMDocumentStatus(*progear.Status)
		payload.Status = &status
	}

	if progear.Reason != nil {
		reason := ghcmessages.PPMDocumentStatusReason(*progear.Reason)
		payload.Reason = &reason
	}

	return payload
}

// MovingExpense payload
func MovingExpense(storer storage.FileStorer, movingExpense *models.MovingExpense) *ghcmessages.MovingExpense {

	document, err := PayloadForDocumentModel(storer, movingExpense.Document)
	if err != nil {
		return nil
	}

	payload := &ghcmessages.MovingExpense{
		ID:               *handlers.FmtUUID(movingExpense.ID),
		PpmShipmentID:    *handlers.FmtUUID(movingExpense.PPMShipmentID),
		DocumentID:       *handlers.FmtUUID(movingExpense.DocumentID),
		Document:         document,
		CreatedAt:        strfmt.DateTime(movingExpense.CreatedAt),
		UpdatedAt:        strfmt.DateTime(movingExpense.UpdatedAt),
		Description:      movingExpense.Description,
		PaidWithGtcc:     movingExpense.PaidWithGTCC,
		Amount:           handlers.FmtCost(movingExpense.Amount),
		MissingReceipt:   movingExpense.MissingReceipt,
		ETag:             etag.GenerateEtag(movingExpense.UpdatedAt),
		SitEstimatedCost: handlers.FmtCost(movingExpense.SITEstimatedCost),
	}
	if movingExpense.MovingExpenseType != nil {
		movingExpenseType := ghcmessages.OmittableMovingExpenseType(*movingExpense.MovingExpenseType)
		payload.MovingExpenseType = &movingExpenseType
	}

	if movingExpense.Status != nil {
		status := ghcmessages.OmittablePPMDocumentStatus(*movingExpense.Status)
		payload.Status = &status
	}

	if movingExpense.Reason != nil {
		reason := ghcmessages.PPMDocumentStatusReason(*movingExpense.Reason)
		payload.Reason = &reason
	}

	if movingExpense.SITStartDate != nil {
		payload.SitStartDate = handlers.FmtDatePtr(movingExpense.SITStartDate)
	}

	if movingExpense.SITEndDate != nil {
		payload.SitEndDate = handlers.FmtDatePtr(movingExpense.SITEndDate)
	}

	if movingExpense.WeightStored != nil {
		payload.WeightStored = handlers.FmtPoundPtr(movingExpense.WeightStored)
	}

	if movingExpense.SITLocation != nil {
		sitLocation := ghcmessages.SITLocationType(*movingExpense.SITLocation)
		payload.SitLocation = &sitLocation
	}

	if movingExpense.SITReimburseableAmount != nil {
		payload.SitReimburseableAmount = handlers.FmtCost(movingExpense.SITReimburseableAmount)
	}

	return payload
}

func MovingExpenses(storer storage.FileStorer, movingExpenses models.MovingExpenses) []*ghcmessages.MovingExpense {
	payload := make([]*ghcmessages.MovingExpense, len(movingExpenses))
	for i, movingExpense := range movingExpenses {
		copyOfMovingExpense := movingExpense
		payload[i] = MovingExpense(storer, &copyOfMovingExpense)
	}
	return payload
}

func WeightTickets(storer storage.FileStorer, weightTickets models.WeightTickets) []*ghcmessages.WeightTicket {
	payload := make([]*ghcmessages.WeightTicket, len(weightTickets))
	for i, weightTicket := range weightTickets {
		copyOfWeightTicket := weightTicket
		weightTicketPayload := WeightTicket(storer, &copyOfWeightTicket)
		payload[i] = weightTicketPayload
	}
	return payload
}

// WeightTicket payload
func WeightTicket(storer storage.FileStorer, weightTicket *models.WeightTicket) *ghcmessages.WeightTicket {
	ppmShipment := strfmt.UUID(weightTicket.PPMShipmentID.String())

	emptyDocument, err := PayloadForDocumentModel(storer, weightTicket.EmptyDocument)
	if err != nil {
		return nil
	}

	fullDocument, err := PayloadForDocumentModel(storer, weightTicket.FullDocument)
	if err != nil {
		return nil
	}

	proofOfTrailerOwnershipDocument, err := PayloadForDocumentModel(storer, weightTicket.ProofOfTrailerOwnershipDocument)
	if err != nil {
		return nil
	}

	payload := &ghcmessages.WeightTicket{
		ID:                                strfmt.UUID(weightTicket.ID.String()),
		PpmShipmentID:                     ppmShipment,
		CreatedAt:                         *handlers.FmtDateTime(weightTicket.CreatedAt),
		UpdatedAt:                         *handlers.FmtDateTime(weightTicket.UpdatedAt),
		VehicleDescription:                weightTicket.VehicleDescription,
		EmptyWeight:                       handlers.FmtPoundPtr(weightTicket.EmptyWeight),
		MissingEmptyWeightTicket:          weightTicket.MissingEmptyWeightTicket,
		EmptyDocumentID:                   *handlers.FmtUUID(weightTicket.EmptyDocumentID),
		EmptyDocument:                     emptyDocument,
		FullWeight:                        handlers.FmtPoundPtr(weightTicket.FullWeight),
		MissingFullWeightTicket:           weightTicket.MissingFullWeightTicket,
		FullDocumentID:                    *handlers.FmtUUID(weightTicket.FullDocumentID),
		FullDocument:                      fullDocument,
		OwnsTrailer:                       weightTicket.OwnsTrailer,
		TrailerMeetsCriteria:              weightTicket.TrailerMeetsCriteria,
		ProofOfTrailerOwnershipDocumentID: *handlers.FmtUUID(weightTicket.ProofOfTrailerOwnershipDocumentID),
		ProofOfTrailerOwnershipDocument:   proofOfTrailerOwnershipDocument,
		AdjustedNetWeight:                 handlers.FmtPoundPtr(weightTicket.AdjustedNetWeight),
		AllowableWeight:                   handlers.FmtPoundPtr(weightTicket.AllowableWeight),
		NetWeightRemarks:                  weightTicket.NetWeightRemarks,
		ETag:                              etag.GenerateEtag(weightTicket.UpdatedAt),
	}

	if weightTicket.Status != nil {
		status := ghcmessages.OmittablePPMDocumentStatus(*weightTicket.Status)
		payload.Status = &status
	}

	if weightTicket.Reason != nil {
		reason := ghcmessages.PPMDocumentStatusReason(*weightTicket.Reason)
		payload.Reason = &reason
	}

	return payload
}

// PPMDocuments payload
func PPMDocuments(storer storage.FileStorer, ppmDocuments *models.PPMDocuments) *ghcmessages.PPMDocuments {

	if ppmDocuments == nil {
		return nil
	}

	payload := &ghcmessages.PPMDocuments{
		WeightTickets:        WeightTickets(storer, ppmDocuments.WeightTickets),
		MovingExpenses:       MovingExpenses(storer, ppmDocuments.MovingExpenses),
		ProGearWeightTickets: ProGearWeightTickets(storer, ppmDocuments.ProgearWeightTickets),
	}

	return payload
}

// PPMCloseout payload
func PPMCloseout(ppmCloseout *models.PPMCloseout) *ghcmessages.PPMCloseout {
	if ppmCloseout == nil {
		return nil
	}
	payload := &ghcmessages.PPMCloseout{
		ID:                    strfmt.UUID(ppmCloseout.ID.String()),
		PlannedMoveDate:       handlers.FmtDatePtr(ppmCloseout.PlannedMoveDate),
		ActualMoveDate:        handlers.FmtDatePtr(ppmCloseout.ActualMoveDate),
		Miles:                 handlers.FmtIntPtrToInt64(ppmCloseout.Miles),
		EstimatedWeight:       handlers.FmtPoundPtr(ppmCloseout.EstimatedWeight),
		ActualWeight:          handlers.FmtPoundPtr(ppmCloseout.ActualWeight),
		ProGearWeightCustomer: handlers.FmtPoundPtr(ppmCloseout.ProGearWeightCustomer),
		ProGearWeightSpouse:   handlers.FmtPoundPtr(ppmCloseout.ProGearWeightSpouse),
		GrossIncentive:        handlers.FmtCost(ppmCloseout.GrossIncentive),
		Gcc:                   handlers.FmtCost(ppmCloseout.GCC),
		Aoa:                   handlers.FmtCost(ppmCloseout.AOA),
		RemainingIncentive:    handlers.FmtCost(ppmCloseout.RemainingIncentive),
		HaulType:              (*string)(&ppmCloseout.HaulType),
		HaulPrice:             handlers.FmtCost(ppmCloseout.HaulPrice),
		HaulFSC:               handlers.FmtCost(ppmCloseout.HaulFSC),
		Dop:                   handlers.FmtCost(ppmCloseout.DOP),
		Ddp:                   handlers.FmtCost(ppmCloseout.DDP),
		PackPrice:             handlers.FmtCost(ppmCloseout.PackPrice),
		UnpackPrice:           handlers.FmtCost(ppmCloseout.UnpackPrice),
		SITReimbursement:      handlers.FmtCost(ppmCloseout.SITReimbursement),
	}

	return payload
}

// PPMActualWeight payload
func PPMActualWeight(ppmActualWeight *unit.Pound) *ghcmessages.PPMActualWeight {
	if ppmActualWeight == nil {
		return nil
	}
	payload := &ghcmessages.PPMActualWeight{
		ActualWeight: handlers.FmtPoundPtr(ppmActualWeight),
	}

	return payload
}

func PPMSITEstimatedCostParamsFirstDaySIT(ppmSITFirstDayParams models.PPMSITEstimatedCostParams) *ghcmessages.PPMSITEstimatedCostParamsFirstDaySIT {
	payload := &ghcmessages.PPMSITEstimatedCostParamsFirstDaySIT{
		ContractYearName:       ppmSITFirstDayParams.ContractYearName,
		PriceRateOrFactor:      ppmSITFirstDayParams.PriceRateOrFactor,
		IsPeak:                 ppmSITFirstDayParams.IsPeak,
		EscalationCompounded:   ppmSITFirstDayParams.EscalationCompounded,
		ServiceAreaOrigin:      &ppmSITFirstDayParams.ServiceAreaOrigin,
		ServiceAreaDestination: &ppmSITFirstDayParams.ServiceAreaDestination,
	}
	return payload
}

func PPMSITEstimatedCostParamsAdditionalDaySIT(ppmSITAdditionalDayParams models.PPMSITEstimatedCostParams) *ghcmessages.PPMSITEstimatedCostParamsAdditionalDaySIT {
	payload := &ghcmessages.PPMSITEstimatedCostParamsAdditionalDaySIT{
		ContractYearName:       ppmSITAdditionalDayParams.ContractYearName,
		PriceRateOrFactor:      ppmSITAdditionalDayParams.PriceRateOrFactor,
		IsPeak:                 ppmSITAdditionalDayParams.IsPeak,
		EscalationCompounded:   ppmSITAdditionalDayParams.EscalationCompounded,
		ServiceAreaOrigin:      &ppmSITAdditionalDayParams.ServiceAreaOrigin,
		ServiceAreaDestination: &ppmSITAdditionalDayParams.ServiceAreaDestination,
		NumberDaysSIT:          &ppmSITAdditionalDayParams.NumberDaysSIT,
	}
	return payload
}

func PPMSITEstimatedCost(ppmSITEstimatedCost *models.PPMSITEstimatedCostInfo) *ghcmessages.PPMSITEstimatedCost {
	if ppmSITEstimatedCost == nil {
		return nil
	}
	payload := &ghcmessages.PPMSITEstimatedCost{
		SitCost:                handlers.FmtCost(ppmSITEstimatedCost.EstimatedSITCost),
		PriceFirstDaySIT:       handlers.FmtCost(ppmSITEstimatedCost.PriceFirstDaySIT),
		PriceAdditionalDaySIT:  handlers.FmtCost(ppmSITEstimatedCost.PriceAdditionalDaySIT),
		ParamsFirstDaySIT:      PPMSITEstimatedCostParamsFirstDaySIT(ppmSITEstimatedCost.ParamsFirstDaySIT),
		ParamsAdditionalDaySIT: PPMSITEstimatedCostParamsAdditionalDaySIT(ppmSITEstimatedCost.ParamsAdditionalDaySIT),
	}

	return payload
}

// ShipmentAddressUpdate payload
func ShipmentAddressUpdate(shipmentAddressUpdate *models.ShipmentAddressUpdate) *ghcmessages.ShipmentAddressUpdate {
	if shipmentAddressUpdate == nil || shipmentAddressUpdate.ID.IsNil() {
		return nil
	}

	payload := &ghcmessages.ShipmentAddressUpdate{
		ID:                    strfmt.UUID(shipmentAddressUpdate.ID.String()),
		ShipmentID:            strfmt.UUID(shipmentAddressUpdate.ShipmentID.String()),
		NewAddress:            Address(&shipmentAddressUpdate.NewAddress),
		OriginalAddress:       Address(&shipmentAddressUpdate.OriginalAddress),
		SitOriginalAddress:    Address(shipmentAddressUpdate.SitOriginalAddress),
		ContractorRemarks:     shipmentAddressUpdate.ContractorRemarks,
		OfficeRemarks:         shipmentAddressUpdate.OfficeRemarks,
		Status:                ghcmessages.ShipmentAddressUpdateStatus(shipmentAddressUpdate.Status),
		NewSitDistanceBetween: handlers.FmtIntPtrToInt64(shipmentAddressUpdate.NewSitDistanceBetween),
		OldSitDistanceBetween: handlers.FmtIntPtrToInt64(shipmentAddressUpdate.OldSitDistanceBetween),
	}

	return payload
}

// LineOfAccounting payload
func LineOfAccounting(lineOfAccounting *models.LineOfAccounting) *ghcmessages.LineOfAccounting {
	// Nil check
	if lineOfAccounting == nil {
		return nil
	}

	return &ghcmessages.LineOfAccounting{
		ID:                        strfmt.UUID(lineOfAccounting.ID.String()),
		LoaActvtyID:               lineOfAccounting.LoaActvtyID,
		LoaAgncAcntngCd:           lineOfAccounting.LoaAgncAcntngCd,
		LoaAgncDsbrCd:             lineOfAccounting.LoaAgncDsbrCd,
		LoaAlltSnID:               lineOfAccounting.LoaAlltSnID,
		LoaBafID:                  lineOfAccounting.LoaBafID,
		LoaBdgtAcntClsNm:          lineOfAccounting.LoaBdgtAcntClsNm,
		LoaBetCd:                  lineOfAccounting.LoaBetCd,
		LoaBgFyTx:                 handlers.FmtIntPtrToInt64(lineOfAccounting.LoaBgFyTx),
		LoaBgnDt:                  handlers.FmtDatePtr(lineOfAccounting.LoaBgnDt),
		LoaBgtLnItmID:             lineOfAccounting.LoaBgtLnItmID,
		LoaBgtRstrCd:              lineOfAccounting.LoaBgtRstrCd,
		LoaBgtSubActCd:            lineOfAccounting.LoaBgtSubActCd,
		LoaClsRefID:               lineOfAccounting.LoaClsRefID,
		LoaCstCd:                  lineOfAccounting.LoaCstCd,
		LoaCstCntrID:              lineOfAccounting.LoaCstCntrID,
		LoaCustNm:                 lineOfAccounting.LoaCustNm,
		LoaDfAgncyAlctnRcpntID:    lineOfAccounting.LoaDfAgncyAlctnRcpntID,
		LoaDocID:                  lineOfAccounting.LoaDocID,
		LoaDptID:                  lineOfAccounting.LoaDptID,
		LoaDscTx:                  lineOfAccounting.LoaDscTx,
		LoaDtlRmbsmtSrcID:         lineOfAccounting.LoaDtlRmbsmtSrcID,
		LoaEndDt:                  handlers.FmtDatePtr(lineOfAccounting.LoaEndDt),
		LoaEndFyTx:                handlers.FmtIntPtrToInt64(lineOfAccounting.LoaEndFyTx),
		LoaFmsTrnsactnID:          lineOfAccounting.LoaFmsTrnsactnID,
		LoaFnclArID:               lineOfAccounting.LoaFnclArID,
		LoaFnctPrsNm:              lineOfAccounting.LoaFnctPrsNm,
		LoaFndCntrID:              lineOfAccounting.LoaFndCntrID,
		LoaFndTyFgCd:              lineOfAccounting.LoaFndTyFgCd,
		LoaHistStatCd:             lineOfAccounting.LoaHistStatCd,
		LoaHsGdsCd:                lineOfAccounting.LoaHsGdsCd,
		LoaInstlAcntgActID:        lineOfAccounting.LoaInstlAcntgActID,
		LoaJbOrdNm:                lineOfAccounting.LoaJbOrdNm,
		LoaLclInstlID:             lineOfAccounting.LoaLclInstlID,
		LoaMajClmNm:               lineOfAccounting.LoaMajClmNm,
		LoaMajRmbsmtSrcID:         lineOfAccounting.LoaMajRmbsmtSrcID,
		LoaObjClsID:               lineOfAccounting.LoaObjClsID,
		LoaOpAgncyID:              lineOfAccounting.LoaOpAgncyID,
		LoaPgmElmntID:             lineOfAccounting.LoaPgmElmntID,
		LoaPrjID:                  lineOfAccounting.LoaPrjID,
		LoaSbaltmtRcpntID:         lineOfAccounting.LoaSbaltmtRcpntID,
		LoaScrtyCoopCustCd:        lineOfAccounting.LoaScrtyCoopCustCd,
		LoaScrtyCoopDsgntrCd:      lineOfAccounting.LoaScrtyCoopDsgntrCd,
		LoaScrtyCoopImplAgncCd:    lineOfAccounting.LoaScrtyCoopImplAgncCd,
		LoaScrtyCoopLnItmID:       lineOfAccounting.LoaScrtyCoopLnItmID,
		LoaSpclIntrID:             lineOfAccounting.LoaSpclIntrID,
		LoaSrvSrcID:               lineOfAccounting.LoaSrvSrcID,
		LoaStatCd:                 lineOfAccounting.LoaStatCd,
		LoaSubAcntID:              lineOfAccounting.LoaSubAcntID,
		LoaSysID:                  lineOfAccounting.LoaSysID,
		LoaTnsfrDptNm:             lineOfAccounting.LoaTnsfrDptNm,
		LoaTrnsnID:                lineOfAccounting.LoaTrnsnID,
		LoaTrsySfxTx:              lineOfAccounting.LoaTrsySfxTx,
		LoaTskBdgtSblnTx:          lineOfAccounting.LoaTskBdgtSblnTx,
		LoaUic:                    lineOfAccounting.LoaUic,
		LoaWkCntrRcpntNm:          lineOfAccounting.LoaWkCntrRcpntNm,
		LoaWrkOrdID:               lineOfAccounting.LoaWrkOrdID,
		OrgGrpDfasCd:              lineOfAccounting.OrgGrpDfasCd,
		UpdatedAt:                 strfmt.DateTime(lineOfAccounting.UpdatedAt),
		CreatedAt:                 strfmt.DateTime(lineOfAccounting.CreatedAt),
		ValidLoaForTac:            lineOfAccounting.ValidLoaForTac,
		ValidHhgProgramCodeForLoa: lineOfAccounting.ValidHhgProgramCodeForLoa,
	}
}

// MarketCode payload
func MarketCode(marketCode *models.MarketCode) string {
	if marketCode == nil {
		return "" // Or a default string value
	}
	return string(*marketCode)
}

// MTOShipment payload
func MTOShipment(storer storage.FileStorer, mtoShipment *models.MTOShipment, sitStatusPayload *ghcmessages.SITStatus) *ghcmessages.MTOShipment {

	payload := &ghcmessages.MTOShipment{
		ID:                          strfmt.UUID(mtoShipment.ID.String()),
		MoveTaskOrderID:             strfmt.UUID(mtoShipment.MoveTaskOrderID.String()),
		ShipmentType:                ghcmessages.MTOShipmentType(mtoShipment.ShipmentType),
		Status:                      ghcmessages.MTOShipmentStatus(mtoShipment.Status),
		CounselorRemarks:            mtoShipment.CounselorRemarks,
		CustomerRemarks:             mtoShipment.CustomerRemarks,
		RejectionReason:             mtoShipment.RejectionReason,
		PickupAddress:               Address(mtoShipment.PickupAddress),
		SecondaryDeliveryAddress:    Address(mtoShipment.SecondaryDeliveryAddress),
		SecondaryPickupAddress:      Address(mtoShipment.SecondaryPickupAddress),
		DestinationAddress:          Address(mtoShipment.DestinationAddress),
		HasSecondaryDeliveryAddress: mtoShipment.HasSecondaryDeliveryAddress,
		HasSecondaryPickupAddress:   mtoShipment.HasSecondaryPickupAddress,
		TertiaryDeliveryAddress:     Address(mtoShipment.TertiaryDeliveryAddress),
		TertiaryPickupAddress:       Address(mtoShipment.TertiaryPickupAddress),
		HasTertiaryDeliveryAddress:  mtoShipment.HasTertiaryDeliveryAddress,
		HasTertiaryPickupAddress:    mtoShipment.HasTertiaryPickupAddress,
		ActualProGearWeight:         handlers.FmtPoundPtr(mtoShipment.ActualProGearWeight),
		ActualSpouseProGearWeight:   handlers.FmtPoundPtr(mtoShipment.ActualSpouseProGearWeight),
		PrimeEstimatedWeight:        handlers.FmtPoundPtr(mtoShipment.PrimeEstimatedWeight),
		PrimeActualWeight:           handlers.FmtPoundPtr(mtoShipment.PrimeActualWeight),
		NtsRecordedWeight:           handlers.FmtPoundPtr(mtoShipment.NTSRecordedWeight),
		MtoAgents:                   *MTOAgents(&mtoShipment.MTOAgents),
		MtoServiceItems:             MTOServiceItemModels(mtoShipment.MTOServiceItems, storer),
		Diversion:                   mtoShipment.Diversion,
		DiversionReason:             mtoShipment.DiversionReason,
		Reweigh:                     Reweigh(mtoShipment.Reweigh, sitStatusPayload),
		CreatedAt:                   strfmt.DateTime(mtoShipment.CreatedAt),
		UpdatedAt:                   strfmt.DateTime(mtoShipment.UpdatedAt),
		ETag:                        etag.GenerateEtag(mtoShipment.UpdatedAt),
		DeletedAt:                   handlers.FmtDateTimePtr(mtoShipment.DeletedAt),
		ApprovedDate:                handlers.FmtDateTimePtr(mtoShipment.ApprovedDate),
		SitDaysAllowance:            handlers.FmtIntPtrToInt64(mtoShipment.SITDaysAllowance),
		SitExtensions:               *SITDurationUpdates(&mtoShipment.SITDurationUpdates),
		BillableWeightCap:           handlers.FmtPoundPtr(mtoShipment.BillableWeightCap),
		BillableWeightJustification: mtoShipment.BillableWeightJustification,
		UsesExternalVendor:          mtoShipment.UsesExternalVendor,
		ServiceOrderNumber:          mtoShipment.ServiceOrderNumber,
		StorageFacility:             StorageFacility(mtoShipment.StorageFacility),
		PpmShipment:                 PPMShipment(storer, mtoShipment.PPMShipment),
		BoatShipment:                BoatShipment(storer, mtoShipment.BoatShipment),
		MobileHomeShipment:          MobileHomeShipment(storer, mtoShipment.MobileHome),
		DeliveryAddressUpdate:       ShipmentAddressUpdate(mtoShipment.DeliveryAddressUpdate),
		ShipmentLocator:             handlers.FmtStringPtr(mtoShipment.ShipmentLocator),
		MarketCode:                  MarketCode(&mtoShipment.MarketCode),
	}

	if mtoShipment.Distance != nil {
		payload.Distance = handlers.FmtInt64(int64(*mtoShipment.Distance))
	}

	if sitStatusPayload != nil {
		// If we have a sitStatusPayload, overwrite SitDaysAllowance from the shipment model.
		totalSITAllowance := 0
		if sitStatusPayload.TotalDaysRemaining != nil {
			totalSITAllowance += int(*sitStatusPayload.TotalDaysRemaining)
		}
		if sitStatusPayload.TotalSITDaysUsed != nil {
			totalSITAllowance += int(*sitStatusPayload.TotalSITDaysUsed)
		}
		payload.SitDaysAllowance = handlers.FmtIntPtrToInt64(&totalSITAllowance)
	}

	if len(mtoShipment.SITDurationUpdates) > 0 {
		payload.SitExtensions = *SITDurationUpdates(&mtoShipment.SITDurationUpdates)
	}

	if mtoShipment.RequestedPickupDate != nil && !mtoShipment.RequestedPickupDate.IsZero() {
		payload.RequestedPickupDate = handlers.FmtDatePtr(mtoShipment.RequestedPickupDate)
	}

	if mtoShipment.ActualPickupDate != nil && !mtoShipment.ActualPickupDate.IsZero() {
		payload.ActualPickupDate = handlers.FmtDatePtr(mtoShipment.ActualPickupDate)
	}

	if mtoShipment.ActualDeliveryDate != nil && !mtoShipment.ActualDeliveryDate.IsZero() {
		payload.ActualDeliveryDate = handlers.FmtDatePtr(mtoShipment.ActualDeliveryDate)
	}

	if mtoShipment.RequestedDeliveryDate != nil && !mtoShipment.RequestedDeliveryDate.IsZero() {
		payload.RequestedDeliveryDate = handlers.FmtDatePtr(mtoShipment.RequestedDeliveryDate)
	}

	if mtoShipment.RequiredDeliveryDate != nil && !mtoShipment.RequiredDeliveryDate.IsZero() {
		payload.RequiredDeliveryDate = handlers.FmtDatePtr(mtoShipment.RequiredDeliveryDate)
	}

	if mtoShipment.ScheduledPickupDate != nil {
		payload.ScheduledPickupDate = handlers.FmtDatePtr(mtoShipment.ScheduledPickupDate)
	}

	if mtoShipment.ScheduledDeliveryDate != nil {
		payload.ScheduledDeliveryDate = handlers.FmtDatePtr(mtoShipment.ScheduledDeliveryDate)
	}

	if mtoShipment.DestinationType != nil {
		destinationType := ghcmessages.DestinationType(*mtoShipment.DestinationType)
		payload.DestinationType = &destinationType
	}

	if sitStatusPayload != nil {
		payload.SitStatus = sitStatusPayload
	}

	if mtoShipment.TACType != nil {
		tt := ghcmessages.LOAType(*mtoShipment.TACType)
		payload.TacType = &tt
	}

	if mtoShipment.SACType != nil {
		st := ghcmessages.LOAType(*mtoShipment.SACType)
		payload.SacType = &st
	}

	weightsCalculator := mtoshipment.NewShipmentBillableWeightCalculator()
	calculatedWeights := weightsCalculator.CalculateShipmentBillableWeight(mtoShipment)

	// CalculatedBillableWeight is intentionally not a part of the mto_shipments model
	// because we don't want to store a derived value in the database
	payload.CalculatedBillableWeight = handlers.FmtPoundPtr(calculatedWeights.CalculatedBillableWeight)

	return payload
}

// MTOShipments payload
func MTOShipments(storer storage.FileStorer, mtoShipments *models.MTOShipments, sitStatusPayload map[string]*ghcmessages.SITStatus) *ghcmessages.MTOShipments {
	payload := make(ghcmessages.MTOShipments, len(*mtoShipments))

	for i, m := range *mtoShipments {
		copyOfMtoShipment := m // Make copy to avoid implicit memory aliasing of items from a range statement.
		if sitStatus, ok := sitStatusPayload[copyOfMtoShipment.ID.String()]; ok {
			payload[i] = MTOShipment(storer, &copyOfMtoShipment, sitStatus)
		} else {
			payload[i] = MTOShipment(storer, &copyOfMtoShipment, nil)
		}
	}
	return &payload
}

// MTOAgent payload
func MTOAgent(mtoAgent *models.MTOAgent) *ghcmessages.MTOAgent {
	payload := &ghcmessages.MTOAgent{
		ID:            strfmt.UUID(mtoAgent.ID.String()),
		MtoShipmentID: strfmt.UUID(mtoAgent.MTOShipmentID.String()),
		CreatedAt:     strfmt.DateTime(mtoAgent.CreatedAt),
		UpdatedAt:     strfmt.DateTime(mtoAgent.UpdatedAt),
		FirstName:     mtoAgent.FirstName,
		LastName:      mtoAgent.LastName,
		AgentType:     string(mtoAgent.MTOAgentType),
		Email:         mtoAgent.Email,
		Phone:         mtoAgent.Phone,
		ETag:          etag.GenerateEtag(mtoAgent.UpdatedAt),
	}
	return payload
}

// MTOAgents payload
func MTOAgents(mtoAgents *models.MTOAgents) *ghcmessages.MTOAgents {
	payload := make(ghcmessages.MTOAgents, len(*mtoAgents))
	for i, m := range *mtoAgents {
		copyOfMtoAgent := m // Make copy to avoid implicit memory aliasing of items from a range statement.
		payload[i] = MTOAgent(&copyOfMtoAgent)
	}
	return &payload
}

// PaymentRequests payload
func PaymentRequests(appCtx appcontext.AppContext, prs *models.PaymentRequests, storer storage.FileStorer) (*ghcmessages.PaymentRequests, error) {
	payload := make(ghcmessages.PaymentRequests, len(*prs))

	for i, p := range *prs {
		paymentRequest := p
		pr, err := PaymentRequest(appCtx, &paymentRequest, storer)
		if err != nil {
			return nil, err
		}
		payload[i] = pr
	}
	return &payload, nil
}

// PaymentRequest payload
func PaymentRequest(appCtx appcontext.AppContext, pr *models.PaymentRequest, storer storage.FileStorer) (*ghcmessages.PaymentRequest, error) {
	serviceDocs := make(ghcmessages.ProofOfServiceDocs, len(pr.ProofOfServiceDocs))

	if len(pr.ProofOfServiceDocs) > 0 {
		for i, proofOfService := range pr.ProofOfServiceDocs {
			payload, err := ProofOfServiceDoc(proofOfService, storer)
			if err != nil {
				return nil, err
			}
			serviceDocs[i] = payload
		}
	}

	move, err := Move(&pr.MoveTaskOrder, storer)
	if err != nil {
		return nil, err
	}

	ediErrorInfoEDIType := ""
	ediErrorInfoEDICode := ""
	ediErrorInfoEDIDescription := ""
	ediErrorInfo := pr.EdiErrors
	if ediErrorInfo != nil {
		mostRecentEdiError := ediErrorInfo[0]
		if mostRecentEdiError.EDIType != "" {
			ediErrorInfoEDIType = string(mostRecentEdiError.EDIType)
		}
		if mostRecentEdiError.Code != nil {
			ediErrorInfoEDICode = *mostRecentEdiError.Code
		}
		if mostRecentEdiError.Description != nil {
			ediErrorInfoEDIDescription = *mostRecentEdiError.Description
		}
	}

	var totalTPPSPaidInvoicePriceMillicents *int64
	var tppsPaidInvoiceSellerPaidDate *time.Time
	var TPPSPaidInvoiceReportsForPR models.TPPSPaidInvoiceReportEntrys
	if pr.TPPSPaidInvoiceReports != nil {
		TPPSPaidInvoiceReportsForPR = pr.TPPSPaidInvoiceReports
		if len(TPPSPaidInvoiceReportsForPR) > 0 {
			if TPPSPaidInvoiceReportsForPR[0].InvoiceTotalChargesInMillicents >= 0 {
				totalTPPSPaidInvoicePriceMillicents = models.Int64Pointer(int64(TPPSPaidInvoiceReportsForPR[0].InvoiceTotalChargesInMillicents))
				tppsPaidInvoiceSellerPaidDate = &TPPSPaidInvoiceReportsForPR[0].SellerPaidDate
			}
		}
	}

	return &ghcmessages.PaymentRequest{
		ID:                                   *handlers.FmtUUID(pr.ID),
		IsFinal:                              &pr.IsFinal,
		MoveTaskOrderID:                      *handlers.FmtUUID(pr.MoveTaskOrderID),
		MoveTaskOrder:                        move,
		PaymentRequestNumber:                 pr.PaymentRequestNumber,
		RecalculationOfPaymentRequestID:      handlers.FmtUUIDPtr(pr.RecalculationOfPaymentRequestID),
		RejectionReason:                      pr.RejectionReason,
		Status:                               ghcmessages.PaymentRequestStatus(pr.Status),
		ETag:                                 etag.GenerateEtag(pr.UpdatedAt),
		ServiceItems:                         *PaymentServiceItems(&pr.PaymentServiceItems, &TPPSPaidInvoiceReportsForPR),
		ReviewedAt:                           handlers.FmtDateTimePtr(pr.ReviewedAt),
		ProofOfServiceDocs:                   serviceDocs,
		CreatedAt:                            strfmt.DateTime(pr.CreatedAt),
		SentToGexAt:                          (*strfmt.DateTime)(pr.SentToGexAt),
		ReceivedByGexAt:                      (*strfmt.DateTime)(pr.ReceivedByGexAt),
		EdiErrorType:                         &ediErrorInfoEDIType,
		EdiErrorCode:                         &ediErrorInfoEDICode,
		EdiErrorDescription:                  &ediErrorInfoEDIDescription,
		TppsInvoiceAmountPaidTotalMillicents: totalTPPSPaidInvoicePriceMillicents,
		TppsInvoiceSellerPaidDate:            (*strfmt.DateTime)(tppsPaidInvoiceSellerPaidDate),
	}, nil
}

// PaymentServiceItem payload
func PaymentServiceItem(ps *models.PaymentServiceItem) *ghcmessages.PaymentServiceItem {
	if ps == nil {
		return nil
	}
	paymentServiceItemParams := PaymentServiceItemParams(&ps.PaymentServiceItemParams)

	return &ghcmessages.PaymentServiceItem{
		ID:                       *handlers.FmtUUID(ps.ID),
		MtoServiceItemID:         *handlers.FmtUUID(ps.MTOServiceItemID),
		MtoServiceItemCode:       string(ps.MTOServiceItem.ReService.Code),
		MtoServiceItemName:       ps.MTOServiceItem.ReService.Name,
		MtoShipmentType:          ghcmessages.MTOShipmentType(ps.MTOServiceItem.MTOShipment.ShipmentType),
		MtoShipmentID:            handlers.FmtUUIDPtr(ps.MTOServiceItem.MTOShipmentID),
		CreatedAt:                strfmt.DateTime(ps.CreatedAt),
		PriceCents:               handlers.FmtCost(ps.PriceCents),
		RejectionReason:          ps.RejectionReason,
		Status:                   ghcmessages.PaymentServiceItemStatus(ps.Status),
		ReferenceID:              ps.ReferenceID,
		ETag:                     etag.GenerateEtag(ps.UpdatedAt),
		PaymentServiceItemParams: *paymentServiceItemParams,
	}
}

// PaymentServiceItems payload
func PaymentServiceItems(paymentServiceItems *models.PaymentServiceItems, tppsPaidReportData *models.TPPSPaidInvoiceReportEntrys) *ghcmessages.PaymentServiceItems {
	payload := make(ghcmessages.PaymentServiceItems, len(*paymentServiceItems))
	for i, m := range *paymentServiceItems {
		copyOfPaymentServiceItem := m // Make copy to avoid implicit memory aliasing of items from a range statement.
		payload[i] = PaymentServiceItem(&copyOfPaymentServiceItem)

		// We process TPPS Paid Invoice Reports to get payment information for each payment service item
		// This report tells us how much TPPS paid HS for each item, then we store and display it
		if *tppsPaidReportData != nil {
			tppsDataForPaymentRequest := *tppsPaidReportData
			for tppsDataRowIndex := range tppsDataForPaymentRequest {
				if tppsDataForPaymentRequest[tppsDataRowIndex].ProductDescription == payload[i].MtoServiceItemCode {
					payload[i].TppsInvoiceAmountPaidPerServiceItemMillicents = handlers.FmtMilliCentsPtr(&tppsDataForPaymentRequest[tppsDataRowIndex].LineNetCharge)
				}
			}
		}
	}
	return &payload
}

// PaymentServiceItemParam payload
func PaymentServiceItemParam(paymentServiceItemParam models.PaymentServiceItemParam) *ghcmessages.PaymentServiceItemParam {
	return &ghcmessages.PaymentServiceItemParam{
		ID:                   strfmt.UUID(paymentServiceItemParam.ID.String()),
		PaymentServiceItemID: strfmt.UUID(paymentServiceItemParam.PaymentServiceItemID.String()),
		Key:                  ghcmessages.ServiceItemParamName(paymentServiceItemParam.ServiceItemParamKey.Key),
		Value:                paymentServiceItemParam.Value,
		Type:                 ghcmessages.ServiceItemParamType(paymentServiceItemParam.ServiceItemParamKey.Type),
		Origin:               ghcmessages.ServiceItemParamOrigin(paymentServiceItemParam.ServiceItemParamKey.Origin),
		ETag:                 etag.GenerateEtag(paymentServiceItemParam.UpdatedAt),
	}
}

// PaymentServiceItemParams payload
func PaymentServiceItemParams(paymentServiceItemParams *models.PaymentServiceItemParams) *ghcmessages.PaymentServiceItemParams {
	if paymentServiceItemParams == nil {
		return nil
	}

	payload := make(ghcmessages.PaymentServiceItemParams, len(*paymentServiceItemParams))

	for i, p := range *paymentServiceItemParams {
		payload[i] = PaymentServiceItemParam(p)
	}
	return &payload
}

func ServiceRequestDoc(serviceRequest models.ServiceRequestDocument, storer storage.FileStorer) (*ghcmessages.ServiceRequestDocument, error) {

	uploads := make([]*ghcmessages.Upload, len(serviceRequest.ServiceRequestDocumentUploads))

	if len(serviceRequest.ServiceRequestDocumentUploads) > 0 {
		for i, serviceRequestUpload := range serviceRequest.ServiceRequestDocumentUploads {
			url, err := storer.PresignedURL(serviceRequestUpload.Upload.StorageKey, serviceRequestUpload.Upload.ContentType, serviceRequestUpload.Upload.Filename)
			if err != nil {
				return nil, err
			}
			uploads[i] = Upload(storer, serviceRequestUpload.Upload, url)
		}
	}

	return &ghcmessages.ServiceRequestDocument{
		Uploads: uploads,
	}, nil

}

// MTOServiceItemSingleModel payload
func MTOServiceItemSingleModel(s *models.MTOServiceItem) *ghcmessages.MTOServiceItemSingle {
	return &ghcmessages.MTOServiceItemSingle{
		SitPostalCode:            handlers.FmtStringPtr(s.SITPostalCode),
		ApprovedAt:               handlers.FmtDateTimePtr(s.ApprovedAt),
		CreatedAt:                *handlers.FmtDateTime(s.CreatedAt),
		ID:                       *handlers.FmtUUID(s.ID),
		MoveTaskOrderID:          *handlers.FmtUUID(s.MoveTaskOrderID),
		MtoShipmentID:            handlers.FmtUUID(*s.MTOShipmentID),
		PickupPostalCode:         handlers.FmtStringPtr(s.PickupPostalCode),
		ReServiceID:              *handlers.FmtUUID(s.ReServiceID),
		RejectedAt:               handlers.FmtDateTimePtr(s.RejectedAt),
		RejectionReason:          handlers.FmtStringPtr(s.RejectionReason),
		SitCustomerContacted:     handlers.FmtDatePtr(s.SITCustomerContacted),
		SitDepartureDate:         handlers.FmtDateTimePtr(s.SITDepartureDate),
		SitEntryDate:             handlers.FmtDateTimePtr(s.SITEntryDate),
		SitRequestedDelivery:     handlers.FmtDatePtr(s.SITRequestedDelivery),
		Status:                   handlers.FmtString(string(s.Status)),
		UpdatedAt:                *handlers.FmtDateTime(s.UpdatedAt),
		ConvertToCustomerExpense: *handlers.FmtBool(s.CustomerExpense),
		CustomerExpenseReason:    handlers.FmtStringPtr(s.CustomerExpenseReason),
	}
}

// MTOServiceItemModel payload
func MTOServiceItemModel(s *models.MTOServiceItem, storer storage.FileStorer) *ghcmessages.MTOServiceItem {
	if s == nil {
		return nil
	}

	serviceRequestDocs := make(ghcmessages.ServiceRequestDocuments, len(s.ServiceRequestDocuments))

	if len(s.ServiceRequestDocuments) > 0 {
		for i, serviceRequest := range s.ServiceRequestDocuments {
			payload, err := ServiceRequestDoc(serviceRequest, storer)
			if err != nil {
				return nil
			}
			serviceRequestDocs[i] = payload
		}
	}

	return &ghcmessages.MTOServiceItem{
		ID:                            handlers.FmtUUID(s.ID),
		MoveTaskOrderID:               handlers.FmtUUID(s.MoveTaskOrderID),
		MtoShipmentID:                 handlers.FmtUUIDPtr(s.MTOShipmentID),
		ReServiceID:                   handlers.FmtUUID(s.ReServiceID),
		ReServiceCode:                 handlers.FmtString(string(s.ReService.Code)),
		ReServiceName:                 handlers.FmtStringPtr(&s.ReService.Name),
		Reason:                        handlers.FmtStringPtr(s.Reason),
		RejectionReason:               handlers.FmtStringPtr(s.RejectionReason),
		PickupPostalCode:              handlers.FmtStringPtr(s.PickupPostalCode),
		SITPostalCode:                 handlers.FmtStringPtr(s.SITPostalCode),
		SitEntryDate:                  handlers.FmtDateTimePtr(s.SITEntryDate),
		SitDepartureDate:              handlers.FmtDateTimePtr(s.SITDepartureDate),
		SitCustomerContacted:          handlers.FmtDatePtr(s.SITCustomerContacted),
		SitRequestedDelivery:          handlers.FmtDatePtr(s.SITRequestedDelivery),
		Status:                        ghcmessages.MTOServiceItemStatus(s.Status),
		Description:                   handlers.FmtStringPtr(s.Description),
		Dimensions:                    MTOServiceItemDimensions(s.Dimensions),
		CustomerContacts:              MTOServiceItemCustomerContacts(s.CustomerContacts),
		SitOriginHHGOriginalAddress:   Address(s.SITOriginHHGOriginalAddress),
		SitOriginHHGActualAddress:     Address(s.SITOriginHHGActualAddress),
		SitDestinationOriginalAddress: Address(s.SITDestinationOriginalAddress),
		SitDestinationFinalAddress:    Address(s.SITDestinationFinalAddress),
		EstimatedWeight:               handlers.FmtPoundPtr(s.EstimatedWeight),
		CreatedAt:                     strfmt.DateTime(s.CreatedAt),
		ApprovedAt:                    handlers.FmtDateTimePtr(s.ApprovedAt),
		RejectedAt:                    handlers.FmtDateTimePtr(s.RejectedAt),
		ETag:                          etag.GenerateEtag(s.UpdatedAt),
		ServiceRequestDocuments:       serviceRequestDocs,
		ConvertToCustomerExpense:      *handlers.FmtBool(s.CustomerExpense),
		CustomerExpenseReason:         handlers.FmtStringPtr(s.CustomerExpenseReason),
		SitDeliveryMiles:              handlers.FmtIntPtrToInt64(s.SITDeliveryMiles),
		EstimatedPrice:                handlers.FmtCost(s.PricingEstimate),
		StandaloneCrate:               s.StandaloneCrate,
		LockedPriceCents:              handlers.FmtCost(s.LockedPriceCents),
	}
}

// SITServiceItemGrouping payload
func SITServiceItemGrouping(s models.SITServiceItemGrouping, storer storage.FileStorer) *ghcmessages.SITServiceItemGrouping {
	if len(s.ServiceItems) == 0 {
		return nil
	}

	summary := ghcmessages.SITSummary{
		FirstDaySITServiceItemID: strfmt.UUID(s.Summary.FirstDaySITServiceItemID.String()),
		Location:                 s.Summary.Location,
		DaysInSIT:                handlers.FmtIntPtrToInt64(&s.Summary.DaysInSIT),
		SitEntryDate:             *handlers.FmtDateTime(s.Summary.SITEntryDate),
		SitDepartureDate:         handlers.FmtDateTimePtr(s.Summary.SITDepartureDate),
		SitAuthorizedEndDate:     *handlers.FmtDateTime(s.Summary.SITAuthorizedEndDate),
		SitCustomerContacted:     handlers.FmtDateTimePtr(s.Summary.SITCustomerContacted),
		SitRequestedDelivery:     handlers.FmtDateTimePtr(s.Summary.SITRequestedDelivery),
	}

	serviceItems := MTOServiceItemModels(s.ServiceItems, storer)

	return &ghcmessages.SITServiceItemGrouping{
		Summary:      &summary,
		ServiceItems: serviceItems,
	}
}

// SITServiceItemGroupings payload
func SITServiceItemGroupings(s models.SITServiceItemGroupings, storer storage.FileStorer) ghcmessages.SITServiceItemGroupings {
	sitGroupings := ghcmessages.SITServiceItemGroupings{}
	for _, sitGroup := range s {
		if sitPayload := SITServiceItemGrouping(sitGroup, storer); sitPayload != nil {
			sitGroupings = append(sitGroupings, sitPayload)
		}
	}
	return sitGroupings
}

// MTOServiceItemModels payload
func MTOServiceItemModels(s models.MTOServiceItems, storer storage.FileStorer) ghcmessages.MTOServiceItems {
	serviceItems := ghcmessages.MTOServiceItems{}
	for _, item := range s {
		copyOfServiceItem := item // Make copy to avoid implicit memory aliasing of items from a range statement.
		serviceItems = append(serviceItems, MTOServiceItemModel(&copyOfServiceItem, storer))
	}

	return serviceItems
}

// MTOServiceItemDimension payload
func MTOServiceItemDimension(d *models.MTOServiceItemDimension) *ghcmessages.MTOServiceItemDimension {
	return &ghcmessages.MTOServiceItemDimension{
		ID:     *handlers.FmtUUID(d.ID),
		Type:   ghcmessages.DimensionType(d.Type),
		Length: *d.Length.Int32Ptr(),
		Height: *d.Height.Int32Ptr(),
		Width:  *d.Width.Int32Ptr(),
	}
}

// MTOServiceItemDimensions payload
func MTOServiceItemDimensions(d models.MTOServiceItemDimensions) ghcmessages.MTOServiceItemDimensions {
	payload := make(ghcmessages.MTOServiceItemDimensions, len(d))
	for i, item := range d {
		copyOfServiceItem := item // Make copy to avoid implicit memory aliasing of items from a range statement.
		payload[i] = MTOServiceItemDimension(&copyOfServiceItem)
	}
	return payload
}

// MTOServiceItemCustomerContact payload
func MTOServiceItemCustomerContact(c *models.MTOServiceItemCustomerContact) *ghcmessages.MTOServiceItemCustomerContact {
	return &ghcmessages.MTOServiceItemCustomerContact{
		Type:                       ghcmessages.CustomerContactType(c.Type),
		DateOfContact:              *handlers.FmtDate(c.DateOfContact),
		TimeMilitary:               c.TimeMilitary,
		FirstAvailableDeliveryDate: *handlers.FmtDate(c.FirstAvailableDeliveryDate),
	}
}

// MTOServiceItemCustomerContacts payload
func MTOServiceItemCustomerContacts(c models.MTOServiceItemCustomerContacts) ghcmessages.MTOServiceItemCustomerContacts {
	payload := make(ghcmessages.MTOServiceItemCustomerContacts, len(c))
	for i, item := range c {
		copyOfServiceItem := item // Make copy to avoid implicit memory aliasing of items from a range statement.
		payload[i] = MTOServiceItemCustomerContact(&copyOfServiceItem)
	}
	return payload
}

// Upload payload
func Upload(storer storage.FileStorer, upload models.Upload, url string) *ghcmessages.Upload {
	uploadPayload := &ghcmessages.Upload{
		ID:          handlers.FmtUUIDValue(upload.ID),
		Filename:    upload.Filename,
		ContentType: upload.ContentType,
		UploadType:  string(upload.UploadType),
		URL:         strfmt.URI(url),
		Bytes:       upload.Bytes,
		CreatedAt:   strfmt.DateTime(upload.CreatedAt),
		UpdatedAt:   strfmt.DateTime(upload.UpdatedAt),
		DeletedAt:   (*strfmt.DateTime)(upload.DeletedAt),
	}

	if upload.Rotation != nil {
		uploadPayload.Rotation = *upload.Rotation
	}

	tags, err := storer.Tags(upload.StorageKey)
	if err != nil || len(tags) == 0 {
		uploadPayload.Status = "PROCESSING"
	} else {
		uploadPayload.Status = tags["av-status"]
	}
	return uploadPayload
}

// Upload payload for when a Proof of Service doc is designated as a weight ticket
// This adds an isWeightTicket key to the payload for the UI to use
func WeightTicketUpload(storer storage.FileStorer, upload models.Upload, url string, isWeightTicket bool) *ghcmessages.Upload {
	uploadPayload := &ghcmessages.Upload{
		ID:             handlers.FmtUUIDValue(upload.ID),
		Filename:       upload.Filename,
		ContentType:    upload.ContentType,
		URL:            strfmt.URI(url),
		Bytes:          upload.Bytes,
		CreatedAt:      strfmt.DateTime(upload.CreatedAt),
		UpdatedAt:      strfmt.DateTime(upload.UpdatedAt),
		IsWeightTicket: isWeightTicket,
	}
	tags, err := storer.Tags(upload.StorageKey)
	if err != nil || len(tags) == 0 {
		uploadPayload.Status = "PROCESSING"
	} else {
		uploadPayload.Status = tags["av-status"]
	}
	return uploadPayload
}

// ProofOfServiceDoc payload from model
func ProofOfServiceDoc(proofOfService models.ProofOfServiceDoc, storer storage.FileStorer) (*ghcmessages.ProofOfServiceDoc, error) {

	uploads := make([]*ghcmessages.Upload, len(proofOfService.PrimeUploads))
	if len(proofOfService.PrimeUploads) > 0 {
		for i, primeUpload := range proofOfService.PrimeUploads {
			url, err := storer.PresignedURL(primeUpload.Upload.StorageKey, primeUpload.Upload.ContentType, primeUpload.Upload.Filename)
			if err != nil {
				return nil, err
			}
			// if the doc is a weight ticket then we need to return a different payload so the UI can differentiate
			weightTicket := proofOfService.IsWeightTicket
			if weightTicket {
				uploads[i] = WeightTicketUpload(storer, primeUpload.Upload, url, proofOfService.IsWeightTicket)
			} else {
				uploads[i] = Upload(storer, primeUpload.Upload, url)
			}
		}
	}

	return &ghcmessages.ProofOfServiceDoc{
		IsWeightTicket: proofOfService.IsWeightTicket,
		Uploads:        uploads,
	}, nil
}

func PayloadForUploadModel(
	storer storage.FileStorer,
	upload models.Upload,
	url string,
) *ghcmessages.Upload {
	uploadPayload := &ghcmessages.Upload{
		ID:          handlers.FmtUUIDValue(upload.ID),
		Filename:    upload.Filename,
		ContentType: upload.ContentType,
		UploadType:  string(upload.UploadType),
		URL:         strfmt.URI(url),
		Bytes:       upload.Bytes,
		CreatedAt:   strfmt.DateTime(upload.CreatedAt),
		UpdatedAt:   strfmt.DateTime(upload.UpdatedAt),
		DeletedAt:   (*strfmt.DateTime)(upload.DeletedAt),
	}

	if upload.Rotation != nil {
		uploadPayload.Rotation = *upload.Rotation
	}

	tags, err := storer.Tags(upload.StorageKey)
	if err != nil || len(tags) == 0 {
		uploadPayload.Status = "PROCESSING"
	} else {
		uploadPayload.Status = tags["av-status"]
	}
	return uploadPayload
}

func PayloadForDocumentModel(storer storage.FileStorer, document models.Document) (*ghcmessages.Document, error) {
	uploads := make([]*ghcmessages.Upload, len(document.UserUploads))
	for i, userUpload := range document.UserUploads {
		if userUpload.Upload.ID == uuid.Nil {
			return nil, errors.New("no uploads for user")
		}
		url, err := storer.PresignedURL(userUpload.Upload.StorageKey, userUpload.Upload.ContentType, userUpload.Upload.Filename)
		if err != nil {
			return nil, err
		}

		uploadPayload := PayloadForUploadModel(storer, userUpload.Upload, url)
		uploads[i] = uploadPayload
	}

	documentPayload := &ghcmessages.Document{
		ID:              handlers.FmtUUID(document.ID),
		ServiceMemberID: handlers.FmtUUID(document.ServiceMemberID),
		Uploads:         uploads,
	}
	return documentPayload, nil
}

// In the TOO queue response we only want to count shipments in these statuses (excluding draft and cancelled)
// For the Services Counseling queue we will find the earliest move date from shipments in these statuses
func queueIncludeShipmentStatus(status models.MTOShipmentStatus) bool {
	return status == models.MTOShipmentStatusSubmitted ||
		status == models.MTOShipmentStatusApproved ||
		status == models.MTOShipmentStatusDiversionRequested ||
		status == models.MTOShipmentStatusCancellationRequested
}

func QueueAvailableOfficeUsers(officeUsers []models.OfficeUser) *ghcmessages.AvailableOfficeUsers {
	availableOfficeUsers := make(ghcmessages.AvailableOfficeUsers, len(officeUsers))
	for i, officeUser := range officeUsers {

		hasSafety := officeUser.User.Privileges.HasPrivilege(models.PrivilegeTypeSafety)

		availableOfficeUsers[i] = &ghcmessages.AvailableOfficeUser{
			LastName:           officeUser.LastName,
			FirstName:          officeUser.FirstName,
			OfficeUserID:       *handlers.FmtUUID(officeUser.ID),
			HasSafetyPrivilege: swag.BoolValue(&hasSafety),
		}
	}

	return &availableOfficeUsers
}

// QueueMoves payload
func QueueMoves(moves []models.Move, officeUsers []models.OfficeUser, requestedPpmStatus *models.PPMShipmentStatus, role roles.RoleType, officeUser models.OfficeUser, isSupervisor bool, isHQRole bool) *ghcmessages.QueueMoves {
	queueMoves := make(ghcmessages.QueueMoves, len(moves))
	for i, move := range moves {
		customer := move.Orders.ServiceMember

		var transportationOffice string
		var transportationOfficeId uuid.UUID
		if move.CounselingOffice != nil {
			transportationOffice = move.CounselingOffice.Name
			transportationOfficeId = move.CounselingOffice.ID
		}
		var validMTOShipments []models.MTOShipment
		var earliestRequestedPickup *time.Time
		// we can't easily modify our sql query to find the earliest shipment pickup date so we must do it here
		for _, shipment := range move.MTOShipments {
			if queueIncludeShipmentStatus(shipment.Status) && shipment.DeletedAt == nil {
				earliestDateInCurrentShipment := findEarliestDateForRequestedMoveDate(shipment)
				if earliestRequestedPickup == nil || (earliestDateInCurrentShipment != nil && earliestDateInCurrentShipment.Before(*earliestRequestedPickup)) {
					earliestRequestedPickup = earliestDateInCurrentShipment
				}

				validMTOShipments = append(validMTOShipments, shipment)
			}
		}

		var deptIndicator ghcmessages.DeptIndicator
		if move.Orders.DepartmentIndicator != nil {
			deptIndicator = ghcmessages.DeptIndicator(*move.Orders.DepartmentIndicator)
		}

		var gbloc string
		if move.Status == models.MoveStatusNeedsServiceCounseling {
			gbloc = swag.StringValue(move.Orders.OriginDutyLocationGBLOC)
		} else if len(move.ShipmentGBLOC) > 0 && move.ShipmentGBLOC[0].GBLOC != nil {
			// There is a Pop bug that prevents us from using a has_one association for
			// Move.ShipmentGBLOC, so we have to treat move.ShipmentGBLOC as an array, even
			// though there can never be more than one GBLOC for a move.
			gbloc = swag.StringValue(move.ShipmentGBLOC[0].GBLOC)
		} else {
			// If the move's first shipment doesn't have a pickup address (like with an NTS-Release),
			// we need to fall back to the origin duty location GBLOC.  If that's not available for
			// some reason, then we should get the empty string (no GBLOC).
			gbloc = swag.StringValue(move.Orders.OriginDutyLocationGBLOC)
		}
		var closeoutLocation string
		if move.CloseoutOffice != nil {
			closeoutLocation = move.CloseoutOffice.Name
		}
		var closeoutInitiated time.Time
		var ppmStatus models.PPMShipmentStatus
		for _, shipment := range move.MTOShipments {
			if shipment.PPMShipment != nil {
				if requestedPpmStatus != nil {
					if shipment.PPMShipment.Status == *requestedPpmStatus {
						ppmStatus = shipment.PPMShipment.Status
					}
				} else {
					ppmStatus = shipment.PPMShipment.Status
				}
				if shipment.PPMShipment.SubmittedAt != nil {
					if closeoutInitiated.Before(*shipment.PPMShipment.SubmittedAt) {
						closeoutInitiated = *shipment.PPMShipment.SubmittedAt
					}
				}
			}
		}

		queueMoves[i] = &ghcmessages.QueueMove{
			Customer:                Customer(&customer),
			Status:                  ghcmessages.MoveStatus(move.Status),
			ID:                      *handlers.FmtUUID(move.ID),
			Locator:                 move.Locator,
			SubmittedAt:             handlers.FmtDateTimePtr(move.SubmittedAt),
			AppearedInTooAt:         handlers.FmtDateTimePtr(findLastSentToTOO(move)),
			RequestedMoveDate:       handlers.FmtDatePtr(earliestRequestedPickup),
			DepartmentIndicator:     &deptIndicator,
			ShipmentsCount:          int64(len(validMTOShipments)),
			OriginDutyLocation:      DutyLocation(move.Orders.OriginDutyLocation),
			DestinationDutyLocation: DutyLocation(&move.Orders.NewDutyLocation), // #nosec G601 new in 1.22.2
			OriginGBLOC:             ghcmessages.GBLOC(gbloc),
			PpmType:                 move.PPMType,
			CloseoutInitiated:       handlers.FmtDateTimePtr(&closeoutInitiated),
			CloseoutLocation:        &closeoutLocation,
			OrderType:               (*string)(move.Orders.OrdersType.Pointer()),
			LockedByOfficeUserID:    handlers.FmtUUIDPtr(move.LockedByOfficeUserID),
			LockedByOfficeUser:      OfficeUser(move.LockedByOfficeUser),
			LockExpiresAt:           handlers.FmtDateTimePtr(move.LockExpiresAt),
			PpmStatus:               ghcmessages.PPMStatus(ppmStatus),
			CounselingOffice:        &transportationOffice,
			CounselingOfficeID:      handlers.FmtUUID(transportationOfficeId),
		}

		if role == roles.RoleTypeServicesCounselor && move.SCAssignedUser != nil {
			queueMoves[i].AssignedTo = AssignedOfficeUser(move.SCAssignedUser)
		}
		if role == roles.RoleTypeTOO && move.TOOAssignedUser != nil {
			queueMoves[i].AssignedTo = AssignedOfficeUser(move.TOOAssignedUser)
		}

		// scenarios where a move is assinable:

		// if it is unassigned, it is always assignable
		isAssignable := false
		if queueMoves[i].AssignedTo == nil {
			isAssignable = true
		}

		// in TOO queues, all moves are assignable for supervisor users
		if role == roles.RoleTypeTOO && isSupervisor {
			isAssignable = true
		}

		// if it is assigned in the SCs queue
		// it is only assignable if the user is a supervisor
		// and if the move's counseling office is the supervisor's transportation office
		if role == roles.RoleTypeServicesCounselor && isSupervisor && move.CounselingOfficeID != nil && *move.CounselingOfficeID == officeUser.TransportationOfficeID {
			isAssignable = true
		}

		if isHQRole {
			isAssignable = false
		}

		queueMoves[i].Assignable = isAssignable

		// only need to attach available office users if move is assignable
		if queueMoves[i].Assignable {
			availableOfficeUsers := officeUsers
			if role == roles.RoleTypeServicesCounselor {
				// if there is no counseling office
				// OR if our current user doesn't work at the move's counseling office
				// only available user should be themself
				if (move.CounselingOfficeID == nil) || (move.CounselingOfficeID != nil && *move.CounselingOfficeID != officeUser.TransportationOfficeID) {
					availableOfficeUsers = models.OfficeUsers{officeUser}
				}

				// if the office user currently assigned to move works outside of the logged in users counseling office
				// add them to the set
				if move.SCAssignedUser != nil && move.SCAssignedUser.TransportationOfficeID != officeUser.TransportationOfficeID {
					availableOfficeUsers = append(availableOfficeUsers, *move.SCAssignedUser)
				}
			}
			queueMoves[i].AvailableOfficeUsers = *QueueAvailableOfficeUsers(availableOfficeUsers)
		}
	}
	return &queueMoves
}

func findLastSentToTOO(move models.Move) (latestOccurance *time.Time) {
	possibleValues := [3]*time.Time{move.SubmittedAt, move.ServiceCounselingCompletedAt, move.ApprovalsRequestedAt}
	for _, time := range possibleValues {
		if time != nil && (latestOccurance == nil || time.After(*latestOccurance)) {
			latestOccurance = time
		}
	}
	return latestOccurance
}

func findEarliestDateForRequestedMoveDate(shipment models.MTOShipment) (earliestDate *time.Time) {
	var possibleValues []*time.Time

	if shipment.RequestedPickupDate != nil {
		possibleValues = append(possibleValues, shipment.RequestedPickupDate)
	}
	if shipment.RequestedDeliveryDate != nil {
		possibleValues = append(possibleValues, shipment.RequestedDeliveryDate)
	}
	if shipment.PPMShipment != nil {
		possibleValues = append(possibleValues, &shipment.PPMShipment.ExpectedDepartureDate)
	}

	for _, date := range possibleValues {
		if earliestDate == nil || date.Before(*earliestDate) {
			earliestDate = date
		}
	}

	return earliestDate
}

// This is a helper function to calculate the inferred status needed for QueuePaymentRequest payload
func queuePaymentRequestStatus(paymentRequest models.PaymentRequest) string {
	// If a payment request is in the PENDING state, let's use the term 'payment requested'
	if paymentRequest.Status == models.PaymentRequestStatusPending {
		return models.QueuePaymentRequestPaymentRequested
	}

	// If a payment request is either reviewed, sent_to_gex or recieved_by_gex then we'll use 'reviewed'
	if paymentRequest.Status == models.PaymentRequestStatusSentToGex ||
		paymentRequest.Status == models.PaymentRequestStatusTppsReceived ||
		paymentRequest.Status == models.PaymentRequestStatusReviewed {
		return models.QueuePaymentRequestReviewed
	}

	if paymentRequest.Status == models.PaymentRequestStatusReviewedAllRejected {
		return models.QueuePaymentRequestRejected
	}

	if paymentRequest.Status == models.PaymentRequestStatusPaid {
		return models.QueuePaymentRequestPaid
	}

	if paymentRequest.Status == models.PaymentRequestStatusDeprecated {
		return models.QueuePaymentRequestDeprecated
	}

	return models.QueuePaymentRequestError

}

// QueuePaymentRequests payload
func QueuePaymentRequests(paymentRequests *models.PaymentRequests, officeUsers []models.OfficeUser, officeUser models.OfficeUser, isSupervisor bool, isHQRole bool) *ghcmessages.QueuePaymentRequests {

	queuePaymentRequests := make(ghcmessages.QueuePaymentRequests, len(*paymentRequests))

	for i, paymentRequest := range *paymentRequests {
		moveTaskOrder := paymentRequest.MoveTaskOrder
		orders := moveTaskOrder.Orders
		var gbloc ghcmessages.GBLOC
		if moveTaskOrder.ShipmentGBLOC[0].GBLOC != nil {
			gbloc = ghcmessages.GBLOC(*moveTaskOrder.ShipmentGBLOC[0].GBLOC)
		}

		queuePaymentRequests[i] = &ghcmessages.QueuePaymentRequest{
			ID:                   *handlers.FmtUUID(paymentRequest.ID),
			MoveID:               *handlers.FmtUUID(moveTaskOrder.ID),
			Customer:             Customer(&orders.ServiceMember),
			Status:               ghcmessages.QueuePaymentRequestStatus(queuePaymentRequestStatus(paymentRequest)),
			Age:                  math.Ceil(time.Since(paymentRequest.CreatedAt).Hours() / 24.0),
			SubmittedAt:          *handlers.FmtDateTime(paymentRequest.CreatedAt),
			Locator:              moveTaskOrder.Locator,
			OriginGBLOC:          gbloc,
			OriginDutyLocation:   DutyLocation(orders.OriginDutyLocation),
			OrderType:            (*string)(orders.OrdersType.Pointer()),
			LockedByOfficeUserID: handlers.FmtUUIDPtr(moveTaskOrder.LockedByOfficeUserID),
			LockExpiresAt:        handlers.FmtDateTimePtr(moveTaskOrder.LockExpiresAt),
		}

		if paymentRequest.MoveTaskOrder.TIOAssignedUser != nil {
			queuePaymentRequests[i].AssignedTo = AssignedOfficeUser(paymentRequest.MoveTaskOrder.TIOAssignedUser)
		}

		isAssignable := false
		if queuePaymentRequests[i].AssignedTo == nil {
			isAssignable = true
		}

		if isSupervisor {
			isAssignable = true
		}

		if isHQRole {
			isAssignable = false
		}

		queuePaymentRequests[i].Assignable = isAssignable

		// only need to attach available office users if move is assignable
		if queuePaymentRequests[i].Assignable {
			availableOfficeUsers := officeUsers
			if !isSupervisor {
				availableOfficeUsers = models.OfficeUsers{officeUser}
			}

			queuePaymentRequests[i].AvailableOfficeUsers = *QueueAvailableOfficeUsers(availableOfficeUsers)
		}

		if orders.DepartmentIndicator != nil {
			deptIndicator := ghcmessages.DeptIndicator(*orders.DepartmentIndicator)
			queuePaymentRequests[i].DepartmentIndicator = &deptIndicator
		}
	}

	return &queuePaymentRequests
}

// Reweigh payload
func Reweigh(reweigh *models.Reweigh, _ *ghcmessages.SITStatus) *ghcmessages.Reweigh {
	if reweigh == nil || reweigh.ID == uuid.Nil {
		return nil
	}
	payload := &ghcmessages.Reweigh{
		ID:                     strfmt.UUID(reweigh.ID.String()),
		RequestedAt:            strfmt.DateTime(reweigh.RequestedAt),
		RequestedBy:            ghcmessages.ReweighRequester(reweigh.RequestedBy),
		VerificationReason:     reweigh.VerificationReason,
		Weight:                 handlers.FmtPoundPtr(reweigh.Weight),
		VerificationProvidedAt: handlers.FmtDateTimePtr(reweigh.VerificationProvidedAt),
		ShipmentID:             strfmt.UUID(reweigh.ShipmentID.String()),
	}

	return payload
}

// SearchMoves payload
func SearchMoves(appCtx appcontext.AppContext, moves models.Moves) *ghcmessages.SearchMoves {
	searchMoves := make(ghcmessages.SearchMoves, len(moves))
	for i, move := range moves {
		customer := move.Orders.ServiceMember

		numShipments := 0

		for _, shipment := range move.MTOShipments {
			if shipment.Status != models.MTOShipmentStatusDraft {
				numShipments++
			}
		}

		var pickupDate, deliveryDate *strfmt.Date

		if numShipments > 0 && move.MTOShipments[0].ScheduledPickupDate != nil {
			pickupDate = handlers.FmtDatePtr(move.MTOShipments[0].ScheduledPickupDate)
		} else {
			pickupDate = nil
		}

		if numShipments > 0 && move.MTOShipments[0].ScheduledDeliveryDate != nil {
			deliveryDate = handlers.FmtDatePtr(move.MTOShipments[0].ScheduledDeliveryDate)
		} else {
			deliveryDate = nil
		}

		var originGBLOC string
		if move.Status == models.MoveStatusNeedsServiceCounseling {
			originGBLOC = swag.StringValue(move.Orders.OriginDutyLocationGBLOC)
		} else if len(move.ShipmentGBLOC) > 0 && move.ShipmentGBLOC[0].GBLOC != nil {
			// There is a Pop bug that prevents us from using a has_one association for
			// Move.ShipmentGBLOC, so we have to treat move.ShipmentGBLOC as an array, even
			// though there can never be more than one GBLOC for a move.
			originGBLOC = swag.StringValue(move.ShipmentGBLOC[0].GBLOC)
		} else {
			// If the move's first shipment doesn't have a pickup address (like with an NTS-Release),
			// we need to fall back to the origin duty location GBLOC.  If that's not available for
			// some reason, then we should get the empty string (no GBLOC).
			originGBLOC = swag.StringValue(move.Orders.OriginDutyLocationGBLOC)
		}

		var destinationGBLOC ghcmessages.GBLOC
		var PostalCodeToGBLOC models.PostalCodeToGBLOC
		var err error
		if numShipments > 0 && move.MTOShipments[0].DestinationAddress != nil {
			PostalCodeToGBLOC, err = models.FetchGBLOCForPostalCode(appCtx.DB(), move.MTOShipments[0].DestinationAddress.PostalCode)
		} else {
			// If the move has no shipments or the shipment has no destination address fall back to the origin duty location GBLOC
			PostalCodeToGBLOC, err = models.FetchGBLOCForPostalCode(appCtx.DB(), move.Orders.NewDutyLocation.Address.PostalCode)
		}

		if err != nil {
			destinationGBLOC = *ghcmessages.NewGBLOC("")
		} else if customer.Affiliation.String() == "MARINES" {
			destinationGBLOC = ghcmessages.GBLOC("USMC/" + PostalCodeToGBLOC.GBLOC)
		} else {
			destinationGBLOC = ghcmessages.GBLOC(PostalCodeToGBLOC.GBLOC)
		}

		searchMoves[i] = &ghcmessages.SearchMove{
			FirstName:                         customer.FirstName,
			LastName:                          customer.LastName,
			DodID:                             customer.Edipi,
			Emplid:                            customer.Emplid,
			Branch:                            customer.Affiliation.String(),
			Status:                            ghcmessages.MoveStatus(move.Status),
			ID:                                *handlers.FmtUUID(move.ID),
			Locator:                           move.Locator,
			ShipmentsCount:                    int64(numShipments),
			OriginDutyLocationPostalCode:      move.Orders.OriginDutyLocation.Address.PostalCode,
			DestinationDutyLocationPostalCode: move.Orders.NewDutyLocation.Address.PostalCode,
			OrderType:                         string(move.Orders.OrdersType),
			RequestedPickupDate:               pickupDate,
			RequestedDeliveryDate:             deliveryDate,
			OriginGBLOC:                       ghcmessages.GBLOC(originGBLOC),
			DestinationGBLOC:                  destinationGBLOC,
			LockedByOfficeUserID:              handlers.FmtUUIDPtr(move.LockedByOfficeUserID),
			LockExpiresAt:                     handlers.FmtDateTimePtr(move.LockExpiresAt),
		}
	}
	return &searchMoves
}

// ShipmentPaymentSITBalance payload
func ShipmentPaymentSITBalance(shipmentSITBalance *services.ShipmentPaymentSITBalance) *ghcmessages.ShipmentPaymentSITBalance {
	if shipmentSITBalance == nil {
		return nil
	}

	payload := &ghcmessages.ShipmentPaymentSITBalance{
		PendingBilledStartDate:  handlers.FmtDate(shipmentSITBalance.PendingBilledStartDate),
		PendingBilledEndDate:    handlers.FmtDate(shipmentSITBalance.PendingBilledEndDate),
		PendingSITDaysInvoiced:  int64(shipmentSITBalance.PendingSITDaysInvoiced),
		PreviouslyBilledDays:    handlers.FmtIntPtrToInt64(shipmentSITBalance.PreviouslyBilledDays),
		PreviouslyBilledEndDate: handlers.FmtDatePtr(shipmentSITBalance.PreviouslyBilledEndDate),
		ShipmentID:              *handlers.FmtUUID(shipmentSITBalance.ShipmentID),
		TotalSITDaysAuthorized:  int64(shipmentSITBalance.TotalSITDaysAuthorized),
		TotalSITDaysRemaining:   int64(shipmentSITBalance.TotalSITDaysRemaining),
		TotalSITEndDate:         handlers.FmtDate(shipmentSITBalance.TotalSITEndDate),
	}

	return payload
}

// ShipmentsPaymentSITBalance payload
func ShipmentsPaymentSITBalance(shipmentsSITBalance []services.ShipmentPaymentSITBalance) ghcmessages.ShipmentsPaymentSITBalance {
	if len(shipmentsSITBalance) == 0 {
		return nil
	}

	payload := make(ghcmessages.ShipmentsPaymentSITBalance, len(shipmentsSITBalance))
	for i, shipmentSITBalance := range shipmentsSITBalance {
		shipmentSITBalanceCopy := shipmentSITBalance
		payload[i] = ShipmentPaymentSITBalance(&shipmentSITBalanceCopy)
	}

	return payload
}

func SearchCustomers(customers models.ServiceMemberSearchResults) *ghcmessages.SearchCustomers {
	searchCustomers := make(ghcmessages.SearchCustomers, len(customers))
	for i, customer := range customers {
		searchCustomers[i] = &ghcmessages.SearchCustomer{
			FirstName:     customer.FirstName,
			LastName:      customer.LastName,
			DodID:         customer.Edipi,
			Emplid:        customer.Emplid,
			Branch:        customer.Affiliation.String(),
			ID:            *handlers.FmtUUID(customer.ID),
			PersonalEmail: customer.PersonalEmail,
			Telephone:     customer.Telephone,
		}
	}
	return &searchCustomers
}

// VLocation payload
func VLocation(vLocation *models.VLocation) *ghcmessages.VLocation {
	if vLocation == nil {
		return nil
	}
	if *vLocation == (models.VLocation{}) {
		return nil
	}

	return &ghcmessages.VLocation{
		City:                 vLocation.CityName,
		State:                vLocation.StateName,
		PostalCode:           vLocation.UsprZipID,
		County:               &vLocation.UsprcCountyNm,
		UsPostRegionCitiesID: *handlers.FmtUUID(*vLocation.UprcId),
	}
}

// VLocations payload
func VLocations(vLocations models.VLocations) ghcmessages.VLocations {
	payload := make(ghcmessages.VLocations, len(vLocations))
	for i, vLocation := range vLocations {
		copyOfVLocation := vLocation
		payload[i] = VLocation(&copyOfVLocation)
	}
	return payload
}
