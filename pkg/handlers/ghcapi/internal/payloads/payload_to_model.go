package payloads

import (
	"errors"
	"time"

	"github.com/go-openapi/strfmt"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/gen/ghcmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/unit"
)

const (
	timeHHMMFormat = "15:04"
)

// MTOAgentModel model
func MTOAgentModel(mtoAgent *ghcmessages.MTOAgent) *models.MTOAgent {
	if mtoAgent == nil {
		return nil
	}

	return &models.MTOAgent{
		ID:            uuid.FromStringOrNil(mtoAgent.ID.String()),
		MTOShipmentID: uuid.FromStringOrNil(mtoAgent.MtoShipmentID.String()),
		FirstName:     mtoAgent.FirstName,
		LastName:      mtoAgent.LastName,
		Email:         mtoAgent.Email,
		Phone:         mtoAgent.Phone,
		MTOAgentType:  models.MTOAgentType(mtoAgent.AgentType),
	}
}

// MTOAgentsModel model
func MTOAgentsModel(mtoAgents *ghcmessages.MTOAgents) *models.MTOAgents {
	if mtoAgents == nil {
		return nil
	}

	agents := make(models.MTOAgents, len(*mtoAgents))

	for i, m := range *mtoAgents {
		agents[i] = *MTOAgentModel(m)
	}

	return &agents
}

// CustomerToServiceMember transforms UpdateCustomerPayload to ServiceMember model
func CustomerToServiceMember(payload ghcmessages.UpdateCustomerPayload) models.ServiceMember {
	address := AddressModel(&payload.CurrentAddress.Address)
	backupAddress := AddressModel(&payload.BackupAddress.Address)

	var backupContacts []models.BackupContact
	if payload.BackupContact != nil {
		backupContacts = []models.BackupContact{{
			Email: *payload.BackupContact.Email,
			Name:  *payload.BackupContact.Name,
			Phone: payload.BackupContact.Phone,
		}}
	}

	return models.ServiceMember{
		ResidentialAddress:   address,
		BackupContacts:       backupContacts,
		FirstName:            &payload.FirstName,
		LastName:             &payload.LastName,
		Suffix:               payload.Suffix,
		MiddleName:           payload.MiddleName,
		PersonalEmail:        payload.Email,
		Telephone:            payload.Phone,
		SecondaryTelephone:   payload.SecondaryTelephone,
		PhoneIsPreferred:     &payload.PhoneIsPreferred,
		EmailIsPreferred:     &payload.EmailIsPreferred,
		BackupMailingAddress: backupAddress,
	}
}

// AddressModel model
func AddressModel(address *ghcmessages.Address) *models.Address {
	// To check if the model is intended to be blank, we'll look at both ID and StreetAddress1
	// We should always have ID if the user intends to update an Address,
	// and StreetAddress1 is a required field on creation. If both are blank, it should be treated as nil.
	var blankSwaggerID strfmt.UUID
	if address == nil || (address.ID == blankSwaggerID && address.StreetAddress1 == nil) {
		return nil
	}

	modelAddress := &models.Address{
		ID:             uuid.FromStringOrNil(address.ID.String()),
		StreetAddress2: address.StreetAddress2,
		StreetAddress3: address.StreetAddress3,
		Country:        address.Country,
	}
	if address.StreetAddress1 != nil {
		modelAddress.StreetAddress1 = *address.StreetAddress1
	}
	if address.City != nil {
		modelAddress.City = *address.City
	}
	if address.State != nil {
		modelAddress.State = *address.State
	}
	if address.PostalCode != nil {
		modelAddress.PostalCode = *address.PostalCode
	}
	return modelAddress
}

// StorageFacilityModel model
func StorageFacilityModel(storageFacility *ghcmessages.StorageFacility) *models.StorageFacility {
	// To check if the model is intended to be blank, we'll look at both ID and FacilityName
	// We should always have ID if the user intends to update a Storage Facility,
	// and FacilityName is a required field on creation. If both are blank, it should be treated as nil.
	var blankSwaggerID strfmt.UUID
	if storageFacility == nil || (storageFacility.ID == blankSwaggerID && storageFacility.FacilityName == "") {
		return nil
	}

	modelStorageFacility := &models.StorageFacility{
		ID:           uuid.FromStringOrNil(storageFacility.ID.String()),
		FacilityName: storageFacility.FacilityName,
		LotNumber:    storageFacility.LotNumber,
		Phone:        storageFacility.Phone,
		Email:        storageFacility.Email,
	}

	addressModel := AddressModel(storageFacility.Address)
	if addressModel != nil {
		modelStorageFacility.Address = *addressModel
	}

	return modelStorageFacility
}

// ApprovedSITExtensionFromCreate model
func ApprovedSITExtensionFromCreate(sitExtension *ghcmessages.CreateApprovedSITDurationUpdate, shipmentID strfmt.UUID) *models.SITDurationUpdate {
	if sitExtension == nil {
		return nil
	}
	now := time.Now()
	ad := int(*sitExtension.ApprovedDays)
	model := &models.SITDurationUpdate{
		MTOShipmentID: uuid.FromStringOrNil(shipmentID.String()),
		RequestReason: models.SITDurationUpdateRequestReason(*sitExtension.RequestReason),
		RequestedDays: int(*sitExtension.ApprovedDays),
		Status:        models.SITExtensionStatusApproved,
		ApprovedDays:  &ad,
		OfficeRemarks: sitExtension.OfficeRemarks,
		DecisionDate:  &now,
	}

	return model
}

// MTOShipmentModelFromCreate model
func MTOShipmentModelFromCreate(mtoShipment *ghcmessages.CreateMTOShipment) *models.MTOShipment {
	if mtoShipment == nil {
		return nil
	}

	var tacType *models.LOAType
	if mtoShipment.TacType != nil {
		tt := models.LOAType(*mtoShipment.TacType)
		tacType = &tt
	}

	var sacType *models.LOAType
	if mtoShipment.SacType != nil {
		st := models.LOAType(*mtoShipment.SacType)
		sacType = &st
	}

	var usesExternalVendor bool
	if mtoShipment.UsesExternalVendor != nil {
		usesExternalVendor = *mtoShipment.UsesExternalVendor
	}

	model := &models.MTOShipment{
		MoveTaskOrderID:    uuid.FromStringOrNil(mtoShipment.MoveTaskOrderID.String()),
		Status:             models.MTOShipmentStatusSubmitted,
		CustomerRemarks:    mtoShipment.CustomerRemarks,
		CounselorRemarks:   mtoShipment.CounselorRemarks,
		TACType:            tacType,
		SACType:            sacType,
		UsesExternalVendor: usesExternalVendor,
		ServiceOrderNumber: mtoShipment.ServiceOrderNumber,
	}

	if mtoShipment.ShipmentType != nil {
		model.ShipmentType = models.MTOShipmentType(*mtoShipment.ShipmentType)
	}

	if mtoShipment.RequestedPickupDate != nil {
		model.RequestedPickupDate = models.TimePointer(time.Time(*mtoShipment.RequestedPickupDate))
	}

	if mtoShipment.RequestedDeliveryDate != nil {
		model.RequestedDeliveryDate = models.TimePointer(time.Time(*mtoShipment.RequestedDeliveryDate))
	}

	// Set up address models
	var addressModel *models.Address

	addressModel = AddressModel(&mtoShipment.PickupAddress.Address)
	if addressModel != nil {
		model.PickupAddress = addressModel
	}

	addressModel = AddressModel(&mtoShipment.DestinationAddress.Address)
	if addressModel != nil {
		model.DestinationAddress = addressModel
	}

	if mtoShipment.DestinationType != nil {
		valDestinationType := models.DestinationType(*mtoShipment.DestinationType)
		model.DestinationType = &valDestinationType
	}

	if mtoShipment.Agents != nil {
		model.MTOAgents = *MTOAgentsModel(&mtoShipment.Agents)
	}

	if mtoShipment.NtsRecordedWeight != nil {
		ntsRecordedWeight := unit.Pound(*mtoShipment.NtsRecordedWeight)
		model.NTSRecordedWeight = &ntsRecordedWeight
	}

	storageFacilityModel := StorageFacilityModel(mtoShipment.StorageFacility)
	if storageFacilityModel != nil {
		model.StorageFacility = storageFacilityModel
	}

	if mtoShipment.PpmShipment != nil {
		model.PPMShipment = PPMShipmentModelFromCreate(mtoShipment.PpmShipment)
		model.PPMShipment.Shipment = *model
	}

	return model
}

// PPMShipmentModelFromCreate model
func PPMShipmentModelFromCreate(ppmShipment *ghcmessages.CreatePPMShipment) *models.PPMShipment {
	if ppmShipment == nil {
		return nil
	}

	model := &models.PPMShipment{
		Status:                         models.PPMShipmentStatusSubmitted,
		SITExpected:                    ppmShipment.SitExpected,
		SecondaryPickupPostalCode:      ppmShipment.SecondaryPickupPostalCode,
		SecondaryDestinationPostalCode: ppmShipment.SecondaryDestinationPostalCode,
		EstimatedWeight:                handlers.PoundPtrFromInt64Ptr(ppmShipment.EstimatedWeight),
		HasProGear:                     ppmShipment.HasProGear,
	}

	expectedDepartureDate := handlers.FmtDatePtrToPopPtr(ppmShipment.ExpectedDepartureDate)
	if expectedDepartureDate != nil && !expectedDepartureDate.IsZero() {
		model.ExpectedDepartureDate = *expectedDepartureDate
	}

	if ppmShipment.PickupPostalCode != nil {
		model.PickupPostalCode = *ppmShipment.PickupPostalCode
	}
	if ppmShipment.DestinationPostalCode != nil {
		model.DestinationPostalCode = *ppmShipment.DestinationPostalCode
	}

	if model.SITExpected != nil && *model.SITExpected {
		if ppmShipment.SitLocation != nil {
			sitLocation := models.SITLocationType(*ppmShipment.SitLocation)
			model.SITLocation = &sitLocation
		}

		model.SITEstimatedWeight = handlers.PoundPtrFromInt64Ptr(ppmShipment.SitEstimatedWeight)

		sitEstimatedEntryDate := handlers.FmtDatePtrToPopPtr(ppmShipment.SitEstimatedEntryDate)
		if sitEstimatedEntryDate != nil && !sitEstimatedEntryDate.IsZero() {
			model.SITEstimatedEntryDate = sitEstimatedEntryDate
		}
		sitEstimatedDepartureDate := handlers.FmtDatePtrToPopPtr(ppmShipment.SitEstimatedDepartureDate)
		if sitEstimatedDepartureDate != nil && !sitEstimatedDepartureDate.IsZero() {
			model.SITEstimatedDepartureDate = sitEstimatedDepartureDate
		}
	}

	if model.HasProGear != nil && *model.HasProGear {
		model.ProGearWeight = handlers.PoundPtrFromInt64Ptr(ppmShipment.ProGearWeight)
		model.SpouseProGearWeight = handlers.PoundPtrFromInt64Ptr(ppmShipment.SpouseProGearWeight)
	}

	return model
}

func CustomerSupportRemarkModelFromCreate(remark *ghcmessages.CreateCustomerSupportRemark) *models.CustomerSupportRemark {
	if remark == nil {
		return nil
	}

	model := &models.CustomerSupportRemark{
		Content:      *remark.Content,
		OfficeUserID: uuid.FromStringOrNil(remark.OfficeUserID.String()),
	}

	return model
}

// MTOShipmentModelFromUpdate model
func MTOShipmentModelFromUpdate(mtoShipment *ghcmessages.UpdateShipment) *models.MTOShipment {
	if mtoShipment == nil {
		return nil
	}

	var requestedPickupDate *time.Time
	if mtoShipment.RequestedPickupDate != nil {
		rpd := time.Time(*mtoShipment.RequestedPickupDate)
		requestedPickupDate = &rpd
	}
	var requestedDeliveryDate *time.Time
	if mtoShipment.RequestedDeliveryDate != nil {
		rdd := time.Time(*mtoShipment.RequestedDeliveryDate)
		requestedDeliveryDate = &rdd
	}
	var billableWeightCap *unit.Pound
	if mtoShipment.BillableWeightCap != nil {
		bwc := unit.Pound(*mtoShipment.BillableWeightCap)
		billableWeightCap = &bwc
	}

	var tacType *models.LOAType
	if mtoShipment.TacType.Present {
		tt := models.LOAType(*mtoShipment.TacType.Value)
		tacType = &tt
	}

	var sacType *models.LOAType
	if mtoShipment.SacType.Present {
		tt := models.LOAType(*mtoShipment.SacType.Value)
		sacType = &tt
	}

	var usesExternalVendor bool
	if mtoShipment.UsesExternalVendor != nil {
		usesExternalVendor = *mtoShipment.UsesExternalVendor
	}

	model := &models.MTOShipment{
		BillableWeightCap:           billableWeightCap,
		BillableWeightJustification: mtoShipment.BillableWeightJustification,
		ShipmentType:                models.MTOShipmentType(mtoShipment.ShipmentType),
		RequestedPickupDate:         requestedPickupDate,
		RequestedDeliveryDate:       requestedDeliveryDate,
		CustomerRemarks:             mtoShipment.CustomerRemarks,
		CounselorRemarks:            mtoShipment.CounselorRemarks,
		TACType:                     tacType,
		SACType:                     sacType,
		UsesExternalVendor:          usesExternalVendor,
		ServiceOrderNumber:          mtoShipment.ServiceOrderNumber,
		HasSecondaryPickupAddress:   mtoShipment.HasSecondaryPickupAddress,
		HasSecondaryDeliveryAddress: mtoShipment.HasSecondaryDeliveryAddress,
	}

	model.PickupAddress = AddressModel(&mtoShipment.PickupAddress.Address)
	model.DestinationAddress = AddressModel(&mtoShipment.DestinationAddress.Address)
	if mtoShipment.HasSecondaryPickupAddress != nil {
		if *mtoShipment.HasSecondaryPickupAddress {
			model.SecondaryPickupAddress = AddressModel(&mtoShipment.SecondaryPickupAddress.Address)
		}
	}
	if mtoShipment.HasSecondaryDeliveryAddress != nil {
		if *mtoShipment.HasSecondaryDeliveryAddress {
			model.SecondaryDeliveryAddress = AddressModel(&mtoShipment.SecondaryDeliveryAddress.Address)
		}
	}

	if mtoShipment.DestinationType != nil {
		valDestinationType := models.DestinationType(*mtoShipment.DestinationType)
		model.DestinationType = &valDestinationType
	}

	if mtoShipment.Agents != nil {
		model.MTOAgents = *MTOAgentsModel(&mtoShipment.Agents)
	}

	if mtoShipment.NtsRecordedWeight != nil {
		ntsRecordedWeight := handlers.PoundPtrFromInt64Ptr(mtoShipment.NtsRecordedWeight)
		model.NTSRecordedWeight = ntsRecordedWeight
	}

	storageFacilityModel := StorageFacilityModel(mtoShipment.StorageFacility)
	if storageFacilityModel != nil {
		model.StorageFacility = storageFacilityModel
	}

	if mtoShipment.PpmShipment != nil {
		model.PPMShipment = PPMShipmentModelFromUpdate(mtoShipment.PpmShipment)
		model.PPMShipment.Shipment = *model
	}

	return model
}

// PPMShipmentModelFromUpdate model
func PPMShipmentModelFromUpdate(ppmShipment *ghcmessages.UpdatePPMShipment) *models.PPMShipment {
	if ppmShipment == nil {
		return nil
	}
	model := &models.PPMShipment{
		ActualMoveDate:                 (*time.Time)(ppmShipment.ActualMoveDate),
		SecondaryPickupPostalCode:      ppmShipment.SecondaryPickupPostalCode,
		SecondaryDestinationPostalCode: ppmShipment.SecondaryDestinationPostalCode,
		SITExpected:                    ppmShipment.SitExpected,
		EstimatedWeight:                handlers.PoundPtrFromInt64Ptr(ppmShipment.EstimatedWeight),
		HasProGear:                     ppmShipment.HasProGear,
		ProGearWeight:                  handlers.PoundPtrFromInt64Ptr(ppmShipment.ProGearWeight),
		SpouseProGearWeight:            handlers.PoundPtrFromInt64Ptr(ppmShipment.SpouseProGearWeight),
		HasRequestedAdvance:            ppmShipment.HasRequestedAdvance,
		AdvanceAmountRequested:         handlers.FmtInt64PtrToPopPtr(ppmShipment.AdvanceAmountRequested),
	}

	expectedDepartureDate := handlers.FmtDatePtrToPopPtr(ppmShipment.ExpectedDepartureDate)
	if expectedDepartureDate != nil && !expectedDepartureDate.IsZero() {
		model.ExpectedDepartureDate = *expectedDepartureDate
	}

	if ppmShipment.PickupPostalCode != nil {
		model.PickupPostalCode = *ppmShipment.PickupPostalCode
	}
	if ppmShipment.DestinationPostalCode != nil {
		model.DestinationPostalCode = *ppmShipment.DestinationPostalCode
	}

	var addressModel *models.Address

	if ppmShipment.W2Address != nil {
		addressModel = AddressModel(ppmShipment.W2Address)
		model.W2Address = addressModel
	}

	if ppmShipment.SitLocation != nil {
		sitLocation := models.SITLocationType(*ppmShipment.SitLocation)
		model.SITLocation = &sitLocation
	}

	if ppmShipment.AdvanceStatus != nil {
		advanceStatus := models.PPMAdvanceStatus(*ppmShipment.AdvanceStatus)
		model.AdvanceStatus = &advanceStatus
	}

	model.SITEstimatedWeight = handlers.PoundPtrFromInt64Ptr(ppmShipment.SitEstimatedWeight)

	sitEstimatedEntryDate := handlers.FmtDatePtrToPopPtr(ppmShipment.SitEstimatedEntryDate)
	if sitEstimatedEntryDate != nil && !sitEstimatedEntryDate.IsZero() {
		model.SITEstimatedEntryDate = sitEstimatedEntryDate
	}
	sitEstimatedDepartureDate := handlers.FmtDatePtrToPopPtr(ppmShipment.SitEstimatedDepartureDate)
	if sitEstimatedDepartureDate != nil && !sitEstimatedDepartureDate.IsZero() {
		model.SITEstimatedDepartureDate = sitEstimatedDepartureDate
	}

	return model
}

// ProgearWeightTicketModelFromUpdate model
func ProgearWeightTicketModelFromUpdate(progearWeightTicket *ghcmessages.UpdateProGearWeightTicket) *models.ProgearWeightTicket {
	if progearWeightTicket == nil {
		return nil
	}

	model := &models.ProgearWeightTicket{
		Weight:           handlers.PoundPtrFromInt64Ptr(progearWeightTicket.Weight),
		HasWeightTickets: handlers.FmtBool(progearWeightTicket.HasWeightTickets),
		BelongsToSelf:    handlers.FmtBool(progearWeightTicket.BelongsToSelf),
		Status:           (*models.PPMDocumentStatus)(handlers.FmtString(string(progearWeightTicket.Status))),
		Reason:           handlers.FmtString(progearWeightTicket.Reason),
	}
	return model
}

// WeightTicketModelFromUpdate
func WeightTicketModelFromUpdate(weightTicket *ghcmessages.UpdateWeightTicket) *models.WeightTicket {
	if weightTicket == nil {
		return nil
	}
	model := &models.WeightTicket{
		EmptyWeight:          handlers.PoundPtrFromInt64Ptr(weightTicket.EmptyWeight),
		FullWeight:           handlers.PoundPtrFromInt64Ptr(weightTicket.FullWeight),
		OwnsTrailer:          handlers.FmtBool(weightTicket.OwnsTrailer),
		TrailerMeetsCriteria: handlers.FmtBool(weightTicket.TrailerMeetsCriteria),
		Status:               (*models.PPMDocumentStatus)(handlers.FmtString(string(weightTicket.Status))),
		Reason:               handlers.FmtString(weightTicket.Reason),
		AdjustedNetWeight:    handlers.PoundPtrFromInt64Ptr(weightTicket.AdjustedNetWeight),
		NetWeightRemarks:     handlers.FmtString(weightTicket.NetWeightRemarks),
		AllowableWeight:      handlers.PoundPtrFromInt64Ptr(weightTicket.AllowableWeight),
	}
	return model
}

// MovingExpenseModelFromUpdate
func MovingExpenseModelFromUpdate(movingExpense *ghcmessages.UpdateMovingExpense) *models.MovingExpense {
	if movingExpense == nil {
		return nil
	}
	model := &models.MovingExpense{
		Amount:       handlers.FmtInt64PtrToPopPtr(&movingExpense.Amount),
		SITStartDate: handlers.FmtDatePtrToPopPtr(&movingExpense.SitStartDate),
		SITEndDate:   handlers.FmtDatePtrToPopPtr(&movingExpense.SitEndDate),
		Status:       (*models.PPMDocumentStatus)(handlers.FmtString(string(movingExpense.Status))),
		Reason:       handlers.FmtString(movingExpense.Reason),
	}

	return model
}

func EvaluationReportFromUpdate(evaluationReport *ghcmessages.EvaluationReport) (*models.EvaluationReport, error) {
	if evaluationReport == nil {
		err := apperror.NewPreconditionFailedError(uuid.UUID{}, errors.New("Cannot update empty report"))
		return nil, err
	}

	var inspectionType *models.EvaluationReportInspectionType
	if evaluationReport.InspectionType != nil {
		tempInspectionType := models.EvaluationReportInspectionType(*evaluationReport.InspectionType)
		inspectionType = &tempInspectionType
	}

	var location *models.EvaluationReportLocationType
	if evaluationReport.Location != nil {
		tempLocation := models.EvaluationReportLocationType(*evaluationReport.Location)
		location = &tempLocation
	}

	var timeDepart *time.Time
	if evaluationReport.TimeDepart != nil {
		td, err := time.Parse(timeHHMMFormat, *evaluationReport.TimeDepart)

		if err != nil {
			return nil, apperror.NewPreconditionFailedError(uuid.UUID{}, err)
		}

		timeDepart = &td
	}

	var evalStart *time.Time
	if evaluationReport.EvalStart != nil {
		es, err := time.Parse(timeHHMMFormat, *evaluationReport.EvalStart)

		if err != nil {
			return nil, apperror.NewPreconditionFailedError(uuid.UUID{}, err)
		}

		evalStart = &es
	}

	var evalEnd *time.Time
	if evaluationReport.EvalEnd != nil {
		ee, err := time.Parse(timeHHMMFormat, *evaluationReport.EvalEnd)

		if err != nil {
			return nil, apperror.NewPreconditionFailedError(uuid.UUID{}, err)
		}

		evalEnd = &ee
	}

	model := models.EvaluationReport{
		ID:                                 uuid.FromStringOrNil(evaluationReport.ID.String()),
		OfficeUser:                         models.OfficeUser{},
		OfficeUserID:                       uuid.Nil,
		Move:                               models.Move{},
		MoveID:                             uuid.Nil,
		Shipment:                           nil,
		ShipmentID:                         nil,
		Type:                               models.EvaluationReportType(evaluationReport.Type),
		InspectionDate:                     (*time.Time)(evaluationReport.InspectionDate),
		InspectionType:                     inspectionType,
		TimeDepart:                         timeDepart,
		EvalStart:                          evalStart,
		EvalEnd:                            evalEnd,
		Location:                           location,
		LocationDescription:                evaluationReport.LocationDescription,
		ObservedShipmentDeliveryDate:       (*time.Time)(evaluationReport.ObservedShipmentDeliveryDate),
		ObservedShipmentPhysicalPickupDate: (*time.Time)(evaluationReport.ObservedShipmentPhysicalPickupDate),
		ViolationsObserved:                 evaluationReport.ViolationsObserved,
		Remarks:                            evaluationReport.Remarks,
		SeriousIncident:                    evaluationReport.SeriousIncident,
		SeriousIncidentDesc:                evaluationReport.SeriousIncidentDesc,
		ObservedClaimsResponseDate:         (*time.Time)(evaluationReport.ObservedClaimsResponseDate),
		ObservedPickupDate:                 (*time.Time)(evaluationReport.ObservedPickupDate),
		ObservedPickupSpreadStartDate:      (*time.Time)(evaluationReport.ObservedPickupSpreadStartDate),
		ObservedPickupSpreadEndDate:        (*time.Time)(evaluationReport.ObservedPickupSpreadEndDate),
		ObservedDeliveryDate:               (*time.Time)(evaluationReport.ObservedDeliveryDate),
		SubmittedAt:                        handlers.FmtDateTimePtrToPopPtr(evaluationReport.SubmittedAt),
	}
	return &model, nil
}
