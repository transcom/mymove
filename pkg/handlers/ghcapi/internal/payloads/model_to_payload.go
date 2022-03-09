package payloads

import (
	"fmt"
	"math"
	"time"

	"github.com/transcom/mymove/pkg/services"

	"github.com/gofrs/uuid"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"

	"github.com/transcom/mymove/pkg/etag"
	"github.com/transcom/mymove/pkg/gen/ghcmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	mtoshipment "github.com/transcom/mymove/pkg/services/mto_shipment"
	"github.com/transcom/mymove/pkg/storage"
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

// Move payload
func Move(move *models.Move) *ghcmessages.Move {
	if move == nil {
		return nil
	}

	payload := &ghcmessages.Move{
		ID:                           strfmt.UUID(move.ID.String()),
		AvailableToPrimeAt:           handlers.FmtDateTimePtr(move.AvailableToPrimeAt),
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
		UpdatedAt:                    strfmt.DateTime(move.UpdatedAt),
		ETag:                         etag.GenerateEtag(move.UpdatedAt),
		ServiceCounselingCompletedAt: handlers.FmtDateTimePtr(move.ServiceCounselingCompletedAt),
		ExcessWeightAcknowledgedAt:   handlers.FmtDateTimePtr(move.ExcessWeightAcknowledgedAt),
		TioRemarks:                   handlers.FmtStringPtr(move.TIORemarks),
		FinancialReviewFlag:          move.FinancialReviewFlag,
		FinancialReviewRemarks:       move.FinancialReviewRemarks,
	}

	return payload
}

// MoveHistory payload
func MoveHistory(moveHistory *models.MoveHistory) *ghcmessages.MoveHistory {
	payload := &ghcmessages.MoveHistory{
		HistoryRecords: moveHistoryRecords(moveHistory.AuditHistories),
		ID:             strfmt.UUID(moveHistory.ID.String()),
		Locator:        moveHistory.Locator,
		ReferenceID:    moveHistory.ReferenceID,
	}

	return payload
}

// MoveAuditHistory payload
func MoveAuditHistory(auditHistory models.AuditHistory) *ghcmessages.MoveAuditHistory {

	payload := &ghcmessages.MoveAuditHistory{
		Action:          auditHistory.Action,
		ActionTstampClk: strfmt.DateTime(auditHistory.ActionTstampClk),
		ActionTstampStm: strfmt.DateTime(auditHistory.ActionTstampStm),
		ActionTstampTx:  strfmt.DateTime(auditHistory.ActionTstampTx),
		ChangedValues:   moveHistoryValues(auditHistory.ChangedData, "changed_data"),
		OldValues:       moveHistoryValues(auditHistory.OldData, "old_values"),
		ClientQuery:     auditHistory.ClientQuery,
		EventName:       auditHistory.EventName,
		ID:              strfmt.UUID(auditHistory.ID.String()),
		ObjectID:        handlers.FmtUUIDPtr(auditHistory.ObjectID),
		RelID:           auditHistory.RelID,
		SessionUserID:   handlers.FmtUUIDPtr(auditHistory.SessionUserID),
		StatementOnly:   auditHistory.StatementOnly,
		TableName:       auditHistory.TableName,
		SchemaName:      auditHistory.SchemaName,
		TransactionID:   auditHistory.TransactionID,
	}

	return payload
}

func moveHistoryValues(data *models.JSONMap, fieldName string) ghcmessages.MoveAuditHistoryItems {
	if data == nil {
		return ghcmessages.MoveAuditHistoryItems{}
	}

	payload := ghcmessages.MoveAuditHistoryItems{}

	for k, v := range *data {
		if v != nil {
			item := ghcmessages.MoveAuditHistoryItem{
				ColumnName:  k,
				ColumnValue: fmt.Sprint(v),
			}
			payload = append(payload, &item)
		}
	}

	return payload
}

func moveHistoryRecords(auditHistories models.AuditHistories) ghcmessages.MoveAuditHistories {
	payload := make(ghcmessages.MoveAuditHistories, len(auditHistories))

	for i, a := range auditHistories {
		payload[i] = MoveAuditHistory(a)
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
		Agency:         swag.StringValue((*string)(customer.Affiliation)),
		CurrentAddress: Address(customer.ResidentialAddress),
		DodID:          swag.StringValue(customer.Edipi),
		Email:          customer.PersonalEmail,
		FirstName:      swag.StringValue(customer.FirstName),
		ID:             strfmt.UUID(customer.ID.String()),
		LastName:       swag.StringValue(customer.LastName),
		Phone:          customer.Telephone,
		Suffix:         customer.Suffix,
		MiddleName:     customer.MiddleName,
		UserID:         strfmt.UUID(customer.UserID.String()),
		ETag:           etag.GenerateEtag(customer.UpdatedAt),
		BackupContact:  BackupContact(customer.BackupContacts),
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

	destinationDutyStation := DutyLocation(&order.NewDutyLocation)
	originDutyLocation := DutyLocation(order.OriginDutyLocation)
	if order.Grade != nil && order.Entitlement != nil {
		order.Entitlement.SetWeightAllotment(*order.Grade)
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
	if order.Moves != nil && len(order.Moves) > 0 {
		moveCode = order.Moves[0].Locator
		moveTaskOrderID = strfmt.UUID(order.Moves[0].ID.String())
	}

	payload := ghcmessages.Order{
		DestinationDutyLocation:     destinationDutyStation,
		Entitlement:                 entitlements,
		Grade:                       &grade,
		OrderNumber:                 order.OrdersNumber,
		OrderTypeDetail:             &ordersTypeDetail,
		ID:                          strfmt.UUID(order.ID.String()),
		OriginDutyLocation:          originDutyLocation,
		ETag:                        etag.GenerateEtag(order.UpdatedAt),
		Agency:                      &affiliation,
		CustomerID:                  strfmt.UUID(order.ServiceMemberID.String()),
		Customer:                    Customer(&order.ServiceMember),
		FirstName:                   swag.StringValue(order.ServiceMember.FirstName),
		LastName:                    swag.StringValue(order.ServiceMember.LastName),
		ReportByDate:                strfmt.Date(order.ReportByDate),
		DateIssued:                  strfmt.Date(order.IssueDate),
		OrderType:                   ghcmessages.OrdersType(order.OrdersType),
		DepartmentIndicator:         &deptIndicator,
		Tac:                         handlers.FmtStringPtr(order.TAC),
		Sac:                         handlers.FmtStringPtr(order.SAC),
		NtsTac:                      handlers.FmtStringPtr(order.NtsTAC),
		NtsSac:                      handlers.FmtStringPtr(order.NtsSAC),
		UploadedOrderID:             strfmt.UUID(order.UploadedOrdersID.String()),
		UploadedAmendedOrderID:      handlers.FmtUUIDPtr(order.UploadedAmendedOrdersID),
		AmendedOrdersAcknowledgedAt: handlers.FmtDateTimePtr(order.AmendedOrdersAcknowledgedAt),
		MoveCode:                    moveCode,
		MoveTaskOrderID:             moveTaskOrderID,
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
		ETag: etag.GenerateEtag(entitlement.UpdatedAt),
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
		Country:        address.Country,
		ETag:           etag.GenerateEtag(address.UpdatedAt),
	}
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

// SITExtension payload
func SITExtension(sitExtension *models.SITExtension) *ghcmessages.SITExtension {
	if sitExtension == nil {
		return nil
	}
	payload := &ghcmessages.SITExtension{
		ID:                strfmt.UUID(sitExtension.ID.String()),
		ETag:              etag.GenerateEtag(sitExtension.UpdatedAt),
		MtoShipmentID:     strfmt.UUID(sitExtension.MTOShipmentID.String()),
		RequestReason:     string(sitExtension.RequestReason),
		RequestedDays:     int64(sitExtension.RequestedDays),
		Status:            string(sitExtension.Status),
		CreatedAt:         strfmt.DateTime(sitExtension.CreatedAt),
		UpdatedAt:         strfmt.DateTime(sitExtension.UpdatedAt),
		ApprovedDays:      handlers.FmtIntPtrToInt64(sitExtension.ApprovedDays),
		ContractorRemarks: handlers.FmtStringPtr(sitExtension.ContractorRemarks),
		DecisionDate:      handlers.FmtDateTimePtr(sitExtension.DecisionDate),
		OfficeRemarks:     handlers.FmtStringPtr(sitExtension.OfficeRemarks),
	}

	return payload
}

// SITExtensions payload
func SITExtensions(sitExtensions *models.SITExtensions) *ghcmessages.SITExtensions {
	payload := make(ghcmessages.SITExtensions, len(*sitExtensions))

	if len(*sitExtensions) > 0 {
		for i, m := range *sitExtensions {
			copyOfSITExtension := m // Make copy to avoid implicit memory aliasing of items from a range statement.
			payload[i] = SITExtension(&copyOfSITExtension)
		}
	}
	return &payload
}

// SITStatus payload
func SITStatus(shipmentSITStatuses *services.SITStatus) *ghcmessages.SITStatus {
	if shipmentSITStatuses == nil {
		return nil
	}
	payload := &ghcmessages.SITStatus{
		DaysInSIT:           handlers.FmtIntPtrToInt64(&shipmentSITStatuses.DaysInSIT),
		TotalDaysRemaining:  handlers.FmtIntPtrToInt64(&shipmentSITStatuses.TotalDaysRemaining),
		Location:            shipmentSITStatuses.Location,
		PastSITServiceItems: MTOServiceItemModels(shipmentSITStatuses.PastSITs),
		SitDepartureDate:    handlers.FmtDateTimePtr(shipmentSITStatuses.SITDepartureDate),
		SitEntryDate:        strfmt.DateTime(shipmentSITStatuses.SITEntryDate),
		TotalSITDaysUsed:    handlers.FmtIntPtrToInt64(&shipmentSITStatuses.TotalSITDaysUsed),
	}

	return payload
}

// SITStatuses payload
func SITStatuses(shipmentSITStatuses map[string]services.SITStatus) map[string]*ghcmessages.SITStatus {
	sitStatuses := map[string]*ghcmessages.SITStatus{}
	if len(shipmentSITStatuses) == 0 {
		return sitStatuses
	}

	for _, sitStatus := range shipmentSITStatuses {
		copyOfSITStatus := sitStatus
		sitStatuses[sitStatus.ShipmentID.String()] = SITStatus(&copyOfSITStatus)
	}

	return sitStatuses
}

// MTOShipment payload
func MTOShipment(mtoShipment *models.MTOShipment, sitStatusPayload *ghcmessages.SITStatus) *ghcmessages.MTOShipment {

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
		PrimeEstimatedWeight:        handlers.FmtPoundPtr(mtoShipment.PrimeEstimatedWeight),
		PrimeActualWeight:           handlers.FmtPoundPtr(mtoShipment.PrimeActualWeight),
		NtsRecordedWeight:           handlers.FmtPoundPtr(mtoShipment.NTSRecordedWeight),
		MtoAgents:                   *MTOAgents(&mtoShipment.MTOAgents),
		MtoServiceItems:             MTOServiceItemModels(mtoShipment.MTOServiceItems),
		Diversion:                   mtoShipment.Diversion,
		Reweigh:                     Reweigh(mtoShipment.Reweigh, sitStatusPayload),
		CreatedAt:                   strfmt.DateTime(mtoShipment.CreatedAt),
		UpdatedAt:                   strfmt.DateTime(mtoShipment.UpdatedAt),
		ETag:                        etag.GenerateEtag(mtoShipment.UpdatedAt),
		DeletedAt:                   handlers.FmtDateTimePtr(mtoShipment.DeletedAt),
		ApprovedDate:                handlers.FmtDateTimePtr(mtoShipment.ApprovedDate),
		SitDaysAllowance:            handlers.FmtIntPtrToInt64(mtoShipment.SITDaysAllowance),
		SitExtensions:               *SITExtensions(&mtoShipment.SITExtensions),
		BillableWeightCap:           handlers.FmtPoundPtr(mtoShipment.BillableWeightCap),
		BillableWeightJustification: mtoShipment.BillableWeightJustification,
		UsesExternalVendor:          mtoShipment.UsesExternalVendor,
		ServiceOrderNumber:          mtoShipment.ServiceOrderNumber,
		StorageFacility:             StorageFacility(mtoShipment.StorageFacility),
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

	if mtoShipment.SITExtensions != nil && len(mtoShipment.SITExtensions) > 0 {
		payload.SitExtensions = *SITExtensions(&mtoShipment.SITExtensions)
	}

	if mtoShipment.RequestedPickupDate != nil && !mtoShipment.RequestedPickupDate.IsZero() {
		payload.RequestedPickupDate = *handlers.FmtDatePtr(mtoShipment.RequestedPickupDate)
	}

	if mtoShipment.ActualPickupDate != nil && !mtoShipment.ActualPickupDate.IsZero() {
		payload.ActualPickupDate = handlers.FmtDatePtr(mtoShipment.ActualPickupDate)
	}

	if mtoShipment.RequestedDeliveryDate != nil && !mtoShipment.RequestedDeliveryDate.IsZero() {
		payload.RequestedDeliveryDate = *handlers.FmtDatePtr(mtoShipment.RequestedDeliveryDate)
	}

	if mtoShipment.ScheduledPickupDate != nil {
		payload.ScheduledPickupDate = handlers.FmtDatePtr(mtoShipment.ScheduledPickupDate)
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
	calculatedWeights, _ := weightsCalculator.CalculateShipmentBillableWeight(mtoShipment)

	// CalculatedBillableWeight is intentionally not a part of the mto_shipments model
	// because we don't want to store a derived value in the database
	payload.CalculatedBillableWeight = handlers.FmtPoundPtr(calculatedWeights.CalculatedBillableWeight)

	return payload
}

// MTOShipments payload
func MTOShipments(mtoShipments *models.MTOShipments, sitStatusPayload map[string]*ghcmessages.SITStatus) *ghcmessages.MTOShipments {
	payload := make(ghcmessages.MTOShipments, len(*mtoShipments))

	for i, m := range *mtoShipments {
		copyOfMtoShipment := m // Make copy to avoid implicit memory aliasing of items from a range statement.
		if sitStatus, ok := sitStatusPayload[copyOfMtoShipment.ID.String()]; ok {
			payload[i] = MTOShipment(&copyOfMtoShipment, sitStatus)
		} else {
			payload[i] = MTOShipment(&copyOfMtoShipment, nil)
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
func PaymentRequests(prs *models.PaymentRequests, storer storage.FileStorer) (*ghcmessages.PaymentRequests, error) {
	payload := make(ghcmessages.PaymentRequests, len(*prs))

	for i, p := range *prs {
		paymentRequest := p
		pr, err := PaymentRequest(&paymentRequest, storer)
		if err != nil {
			return nil, err
		}
		payload[i] = pr
	}
	return &payload, nil
}

// PaymentRequest payload
func PaymentRequest(pr *models.PaymentRequest, storer storage.FileStorer) (*ghcmessages.PaymentRequest, error) {
	serviceDocs := make(ghcmessages.ProofOfServiceDocs, len(pr.ProofOfServiceDocs))

	if pr.ProofOfServiceDocs != nil && len(pr.ProofOfServiceDocs) > 0 {
		for i, proofOfService := range pr.ProofOfServiceDocs {
			payload, err := ProofOfServiceDoc(proofOfService, storer)
			if err != nil {
				return nil, err
			}
			serviceDocs[i] = payload
		}
	}

	return &ghcmessages.PaymentRequest{
		ID:                              *handlers.FmtUUID(pr.ID),
		IsFinal:                         &pr.IsFinal,
		MoveTaskOrderID:                 *handlers.FmtUUID(pr.MoveTaskOrderID),
		MoveTaskOrder:                   Move(&pr.MoveTaskOrder),
		PaymentRequestNumber:            pr.PaymentRequestNumber,
		RecalculationOfPaymentRequestID: handlers.FmtUUIDPtr(pr.RecalculationOfPaymentRequestID),
		RejectionReason:                 pr.RejectionReason,
		Status:                          ghcmessages.PaymentRequestStatus(pr.Status),
		ETag:                            etag.GenerateEtag(pr.UpdatedAt),
		ServiceItems:                    *PaymentServiceItems(&pr.PaymentServiceItems),
		ReviewedAt:                      handlers.FmtDateTimePtr(pr.ReviewedAt),
		ProofOfServiceDocs:              serviceDocs,
		CreatedAt:                       strfmt.DateTime(pr.CreatedAt),
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
func PaymentServiceItems(paymentServiceItems *models.PaymentServiceItems) *ghcmessages.PaymentServiceItems {
	payload := make(ghcmessages.PaymentServiceItems, len(*paymentServiceItems))
	for i, m := range *paymentServiceItems {
		copyOfPaymentServiceItem := m // Make copy to avoid implicit memory aliasing of items from a range statement.
		payload[i] = PaymentServiceItem(&copyOfPaymentServiceItem)
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

// MTOServiceItemModel payload
func MTOServiceItemModel(s *models.MTOServiceItem) *ghcmessages.MTOServiceItem {
	if s == nil {
		return nil
	}

	return &ghcmessages.MTOServiceItem{
		ID:               handlers.FmtUUID(s.ID),
		MoveTaskOrderID:  handlers.FmtUUID(s.MoveTaskOrderID),
		MtoShipmentID:    handlers.FmtUUIDPtr(s.MTOShipmentID),
		ReServiceID:      handlers.FmtUUID(s.ReServiceID),
		ReServiceCode:    handlers.FmtString(string(s.ReService.Code)),
		ReServiceName:    handlers.FmtStringPtr(&s.ReService.Name),
		Reason:           handlers.FmtStringPtr(s.Reason),
		RejectionReason:  handlers.FmtStringPtr(s.RejectionReason),
		PickupPostalCode: handlers.FmtStringPtr(s.PickupPostalCode),
		SITPostalCode:    handlers.FmtStringPtr(s.SITPostalCode),
		SitEntryDate:     handlers.FmtDateTimePtr(s.SITEntryDate),
		SitDepartureDate: handlers.FmtDateTimePtr(s.SITDepartureDate),
		Status:           ghcmessages.MTOServiceItemStatus(s.Status),
		Description:      handlers.FmtStringPtr(s.Description),
		Dimensions:       MTOServiceItemDimensions(s.Dimensions),
		CustomerContacts: MTOServiceItemCustomerContacts(s.CustomerContacts),
		EstimatedWeight:  handlers.FmtPoundPtr(s.EstimatedWeight),
		CreatedAt:        strfmt.DateTime(s.CreatedAt),
		ApprovedAt:       handlers.FmtDateTimePtr(s.ApprovedAt),
		RejectedAt:       handlers.FmtDateTimePtr(s.RejectedAt),
		ETag:             etag.GenerateEtag(s.UpdatedAt),
	}
}

// MTOServiceItemModels payload
func MTOServiceItemModels(s models.MTOServiceItems) ghcmessages.MTOServiceItems {
	serviceItems := ghcmessages.MTOServiceItems{}
	for _, item := range s {
		copyOfServiceItem := item // Make copy to avoid implicit memory aliasing of items from a range statement.
		serviceItems = append(serviceItems, MTOServiceItemModel(&copyOfServiceItem))
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
		ID:          handlers.FmtUUID(upload.ID),
		Filename:    swag.String(upload.Filename),
		ContentType: swag.String(upload.ContentType),
		URL:         handlers.FmtURI(url),
		Bytes:       &upload.Bytes,
		CreatedAt:   handlers.FmtDateTime(upload.CreatedAt),
		UpdatedAt:   handlers.FmtDateTime(upload.UpdatedAt),
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
	if proofOfService.PrimeUploads != nil && len(proofOfService.PrimeUploads) > 0 {
		for i, primeUpload := range proofOfService.PrimeUploads {
			url, err := storer.PresignedURL(primeUpload.Upload.StorageKey, primeUpload.Upload.ContentType)
			if err != nil {
				return nil, err
			}
			uploads[i] = Upload(storer, primeUpload.Upload, url)
		}
	}

	return &ghcmessages.ProofOfServiceDoc{
		Uploads: uploads,
	}, nil
}

// In the TOO queue response we only want to count shipments in these statuses (excluding draft and cancelled)
// For the Services Counseling queue we will find the earliest move date from shipments in these statuses
func queueIncludeShipmentStatus(status models.MTOShipmentStatus) bool {
	return status == models.MTOShipmentStatusSubmitted ||
		status == models.MTOShipmentStatusApproved ||
		status == models.MTOShipmentStatusDiversionRequested ||
		status == models.MTOShipmentStatusCancellationRequested
}

// QueueMoves payload
func QueueMoves(moves []models.Move) *ghcmessages.QueueMoves {
	queueMoves := make(ghcmessages.QueueMoves, len(moves))
	for i, move := range moves {
		customer := move.Orders.ServiceMember

		var validMTOShipments []models.MTOShipment
		var earliestRequestedPickup *time.Time
		// we can't easily modify our sql query to find the earliest shipment pickup date so we must do it here
		for _, shipment := range move.MTOShipments {
			if queueIncludeShipmentStatus(shipment.Status) {
				if earliestRequestedPickup == nil {
					earliestRequestedPickup = shipment.RequestedPickupDate
				} else if shipment.RequestedPickupDate != nil && shipment.RequestedPickupDate.Before(*earliestRequestedPickup) {
					earliestRequestedPickup = shipment.RequestedPickupDate
				}
				validMTOShipments = append(validMTOShipments, shipment)
			}
		}

		var deptIndicator ghcmessages.DeptIndicator
		if move.Orders.DepartmentIndicator != nil {
			deptIndicator = ghcmessages.DeptIndicator(*move.Orders.DepartmentIndicator)
		}

		var gbloc ghcmessages.GBLOC
		if move.Status == models.MoveStatusNeedsServiceCounseling {
			gbloc = ghcmessages.GBLOC(move.OriginDutyLocationGBLOC.GBLOC)
		} else if len(move.ShipmentGBLOC) > 0 {
			// There is a Pop bug that prevents us from using a has_one association for
			// Move.ShipmentGBLOC, so we have to treat move.ShipmentGBLOC as an array, even
			// though there can never be more than one GBLOC for a move.
			gbloc = ghcmessages.GBLOC(move.ShipmentGBLOC[0].GBLOC)
		} else {
			// If the move's first shipment doesn't have a pickup address (like with an NTS-Release),
			// we need to fall back to the origin duty location GBLOC.  If that's not available for
			// some reason, then we should get the empty string (no GBLOC).
			gbloc = ghcmessages.GBLOC(move.OriginDutyLocationGBLOC.GBLOC)
		}

		queueMoves[i] = &ghcmessages.QueueMove{
			Customer:                Customer(&customer),
			Status:                  ghcmessages.MoveStatus(move.Status),
			ID:                      *handlers.FmtUUID(move.ID),
			Locator:                 move.Locator,
			SubmittedAt:             handlers.FmtDateTimePtr(move.SubmittedAt),
			RequestedMoveDate:       handlers.FmtDatePtr(earliestRequestedPickup),
			DepartmentIndicator:     &deptIndicator,
			ShipmentsCount:          int64(len(validMTOShipments)),
			OriginDutyLocation:      DutyLocation(move.Orders.OriginDutyLocation),
			DestinationDutyLocation: DutyLocation(&move.Orders.NewDutyLocation),
			OriginGBLOC:             gbloc,
		}
	}
	return &queueMoves
}

var (
	// QueuePaymentRequestPaymentRequested status payment requested
	QueuePaymentRequestPaymentRequested = "Payment requested"
	// QueuePaymentRequestReviewed status Payment request reviewed
	QueuePaymentRequestReviewed = "Reviewed"
	// QueuePaymentRequestRejected status Payment request rejected
	QueuePaymentRequestRejected = "Rejected"
	// QueuePaymentRequestPaid status PaymentRequest paid
	QueuePaymentRequestPaid = "Paid"
)

// This is a helper function to calculate the inferred status needed for QueuePaymentRequest payload
func queuePaymentRequestStatus(paymentRequest models.PaymentRequest) string {
	// If a payment request is in the PENDING state, let's use the term 'payment requested'
	if paymentRequest.Status == models.PaymentRequestStatusPending {
		return QueuePaymentRequestPaymentRequested
	}

	// If a payment request is either reviewed, sent_to_gex or recieved_by_gex then we'll use 'reviewed'
	if paymentRequest.Status == models.PaymentRequestStatusSentToGex ||
		paymentRequest.Status == models.PaymentRequestStatusReceivedByGex ||
		paymentRequest.Status == models.PaymentRequestStatusReviewed {
		return QueuePaymentRequestReviewed
	}

	if paymentRequest.Status == models.PaymentRequestStatusReviewedAllRejected {
		return QueuePaymentRequestRejected
	}

	return QueuePaymentRequestPaid
}

// QueuePaymentRequests payload
func QueuePaymentRequests(paymentRequests *models.PaymentRequests) *ghcmessages.QueuePaymentRequests {
	queuePaymentRequests := make(ghcmessages.QueuePaymentRequests, len(*paymentRequests))

	for i, paymentRequest := range *paymentRequests {
		moveTaskOrder := paymentRequest.MoveTaskOrder
		orders := moveTaskOrder.Orders

		queuePaymentRequests[i] = &ghcmessages.QueuePaymentRequest{
			ID:                 *handlers.FmtUUID(paymentRequest.ID),
			MoveID:             *handlers.FmtUUID(moveTaskOrder.ID),
			Customer:           Customer(&orders.ServiceMember),
			Status:             ghcmessages.PaymentRequestStatus(queuePaymentRequestStatus(paymentRequest)),
			Age:                math.Ceil(time.Since(paymentRequest.CreatedAt).Hours() / 24.0),
			SubmittedAt:        *handlers.FmtDateTime(paymentRequest.CreatedAt), // RequestedAt does not seem to be populated
			Locator:            moveTaskOrder.Locator,
			OriginGBLOC:        ghcmessages.GBLOC(moveTaskOrder.ShipmentGBLOC[0].GBLOC),
			OriginDutyLocation: DutyLocation(orders.OriginDutyLocation),
		}

		if orders.DepartmentIndicator != nil {
			deptIndicator := ghcmessages.DeptIndicator(*orders.DepartmentIndicator)
			queuePaymentRequests[i].DepartmentIndicator = &deptIndicator
		}
	}

	return &queuePaymentRequests
}

// Reweigh payload
func Reweigh(reweigh *models.Reweigh, sitStatusPayload *ghcmessages.SITStatus) *ghcmessages.Reweigh {
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

// ShipmentPaymentSITBalance payload
func ShipmentPaymentSITBalance(shipmentSITBalance *services.ShipmentPaymentSITBalance) *ghcmessages.ShipmentPaymentSITBalance {
	if shipmentSITBalance == nil {
		return nil
	}

	payload := &ghcmessages.ShipmentPaymentSITBalance{
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
