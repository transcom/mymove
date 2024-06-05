package payloads

import (
	"errors"

	"github.com/go-openapi/strfmt"
	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/etag"
	internalmessages "github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/storage"
)

// Address payload
func Address(address *models.Address) *internalmessages.Address {
	if address == nil {
		return nil
	}
	if *address == (models.Address{}) {
		return nil
	}
	return &internalmessages.Address{
		ID:             strfmt.UUID(address.ID.String()),
		StreetAddress1: &address.StreetAddress1,
		StreetAddress2: address.StreetAddress2,
		StreetAddress3: address.StreetAddress3,
		City:           &address.City,
		State:          &address.State,
		PostalCode:     &address.PostalCode,
		Country:        address.Country,
		County:         &address.County,
	}
}

// MTOAgent payload
func MTOAgent(mtoAgent *models.MTOAgent) *internalmessages.MTOAgent {
	if mtoAgent == nil {
		return nil
	}

	return &internalmessages.MTOAgent{
		AgentType:     internalmessages.MTOAgentType(mtoAgent.MTOAgentType),
		FirstName:     mtoAgent.FirstName,
		LastName:      mtoAgent.LastName,
		Phone:         mtoAgent.Phone,
		Email:         mtoAgent.Email,
		ID:            strfmt.UUID(mtoAgent.ID.String()),
		MtoShipmentID: strfmt.UUID(mtoAgent.MTOShipmentID.String()),
		CreatedAt:     strfmt.DateTime(mtoAgent.CreatedAt),
		UpdatedAt:     strfmt.DateTime(mtoAgent.UpdatedAt),
	}
}

// MTOAgents payload
func MTOAgents(mtoAgents *models.MTOAgents) *internalmessages.MTOAgents {
	if mtoAgents == nil {
		return nil
	}

	agents := make(internalmessages.MTOAgents, len(*mtoAgents))

	for i, m := range *mtoAgents {
		copyOfAgent := m // Make copy to avoid implicit memory aliasing of items from a range statement.
		agents[i] = MTOAgent(&copyOfAgent)
	}

	return &agents
}

// PPMShipment payload
func PPMShipment(storer storage.FileStorer, ppmShipment *models.PPMShipment) *internalmessages.PPMShipment {
	if ppmShipment == nil || ppmShipment.ID.IsNil() {
		return nil
	}

	payloadPPMShipment := &internalmessages.PPMShipment{
		ID:                             *handlers.FmtUUID(ppmShipment.ID),
		ShipmentID:                     *handlers.FmtUUID(ppmShipment.ShipmentID),
		CreatedAt:                      strfmt.DateTime(ppmShipment.CreatedAt),
		UpdatedAt:                      strfmt.DateTime(ppmShipment.UpdatedAt),
		Status:                         internalmessages.PPMShipmentStatus(ppmShipment.Status),
		ExpectedDepartureDate:          handlers.FmtDate(ppmShipment.ExpectedDepartureDate),
		ActualMoveDate:                 handlers.FmtDatePtr(ppmShipment.ActualMoveDate),
		SubmittedAt:                    handlers.FmtDateTimePtr(ppmShipment.SubmittedAt),
		ReviewedAt:                     handlers.FmtDateTimePtr(ppmShipment.ReviewedAt),
		ApprovedAt:                     handlers.FmtDateTimePtr(ppmShipment.ApprovedAt),
		PickupAddress:                  Address(ppmShipment.PickupAddress),
		SecondaryPickupAddress:         Address(ppmShipment.SecondaryPickupAddress),
		HasSecondaryPickupAddress:      ppmShipment.HasSecondaryPickupAddress,
		ActualPickupPostalCode:         ppmShipment.ActualPickupPostalCode,
		DestinationAddress:             Address(ppmShipment.DestinationAddress),
		SecondaryDestinationAddress:    Address(ppmShipment.SecondaryDestinationAddress),
		HasSecondaryDestinationAddress: ppmShipment.HasSecondaryDestinationAddress,
		ActualDestinationPostalCode:    ppmShipment.ActualDestinationPostalCode,
		W2Address:                      Address(ppmShipment.W2Address),
		SitExpected:                    ppmShipment.SITExpected,
		EstimatedWeight:                handlers.FmtPoundPtr(ppmShipment.EstimatedWeight),
		EstimatedIncentive:             handlers.FmtCost(ppmShipment.EstimatedIncentive),
		FinalIncentive:                 handlers.FmtCost(ppmShipment.FinalIncentive),
		HasProGear:                     ppmShipment.HasProGear,
		ProGearWeight:                  handlers.FmtPoundPtr(ppmShipment.ProGearWeight),
		SpouseProGearWeight:            handlers.FmtPoundPtr(ppmShipment.SpouseProGearWeight),
		HasRequestedAdvance:            ppmShipment.HasRequestedAdvance,
		AdvanceAmountRequested:         handlers.FmtCost(ppmShipment.AdvanceAmountRequested),
		HasReceivedAdvance:             ppmShipment.HasReceivedAdvance,
		AdvanceAmountReceived:          handlers.FmtCost(ppmShipment.AdvanceAmountReceived),
		AdvanceStatus:                  (*internalmessages.PPMAdvanceStatus)(ppmShipment.AdvanceStatus),
		WeightTickets:                  WeightTickets(storer, ppmShipment.WeightTickets),
		MovingExpenses:                 MovingExpenses(storer, ppmShipment.MovingExpenses),
		ProGearWeightTickets:           ProGearWeightTickets(storer, ppmShipment.ProgearWeightTickets),
		SignedCertification:            SignedCertification(ppmShipment.SignedCertification),
		ETag:                           etag.GenerateEtag(ppmShipment.UpdatedAt),
	}

	return payloadPPMShipment
}

// MTOShipment payload
func MTOShipment(storer storage.FileStorer, mtoShipment *models.MTOShipment) *internalmessages.MTOShipment {
	payload := &internalmessages.MTOShipment{
		ID:                          strfmt.UUID(mtoShipment.ID.String()),
		Agents:                      *MTOAgents(&mtoShipment.MTOAgents),
		MoveTaskOrderID:             strfmt.UUID(mtoShipment.MoveTaskOrderID.String()),
		ShipmentType:                internalmessages.MTOShipmentType(mtoShipment.ShipmentType),
		CustomerRemarks:             mtoShipment.CustomerRemarks,
		PickupAddress:               Address(mtoShipment.PickupAddress),
		SecondaryPickupAddress:      Address(mtoShipment.SecondaryPickupAddress),
		HasSecondaryPickupAddress:   mtoShipment.HasSecondaryPickupAddress,
		DestinationAddress:          Address(mtoShipment.DestinationAddress),
		SecondaryDeliveryAddress:    Address(mtoShipment.SecondaryDeliveryAddress),
		HasSecondaryDeliveryAddress: mtoShipment.HasSecondaryDeliveryAddress,
		ActualProGearWeight:         handlers.FmtPoundPtr(mtoShipment.ActualProGearWeight),
		ActualSpouseProGearWeight:   handlers.FmtPoundPtr(mtoShipment.ActualSpouseProGearWeight),
		CreatedAt:                   strfmt.DateTime(mtoShipment.CreatedAt),
		UpdatedAt:                   strfmt.DateTime(mtoShipment.UpdatedAt),
		Status:                      internalmessages.MTOShipmentStatus(mtoShipment.Status),
		PpmShipment:                 PPMShipment(storer, mtoShipment.PPMShipment),
		ETag:                        etag.GenerateEtag(mtoShipment.UpdatedAt),
		ShipmentLocator:             handlers.FmtStringPtr(mtoShipment.ShipmentLocator),
	}
	if mtoShipment.HasSecondaryPickupAddress != nil && !*mtoShipment.HasSecondaryPickupAddress {
		payload.SecondaryPickupAddress = nil
	}
	if mtoShipment.HasSecondaryDeliveryAddress != nil && !*mtoShipment.HasSecondaryDeliveryAddress {
		payload.SecondaryDeliveryAddress = nil
	}

	if mtoShipment.RequestedPickupDate != nil && !mtoShipment.RequestedPickupDate.IsZero() {
		payload.RequestedPickupDate = handlers.FmtDatePtr(mtoShipment.RequestedPickupDate)
	}

	if mtoShipment.RequestedDeliveryDate != nil && !mtoShipment.RequestedDeliveryDate.IsZero() {
		payload.RequestedDeliveryDate = handlers.FmtDatePtr(mtoShipment.RequestedDeliveryDate)
	}

	return payload
}

// TransportationOffice internal payload
func TransportationOffice(office models.TransportationOffice) *internalmessages.TransportationOffice {
	if office.ID == uuid.Nil {
		return nil
	}

	phoneLines := []string{}
	for _, phoneLine := range office.PhoneLines {
		if phoneLine.Type == "voice" {
			phoneLines = append(phoneLines, phoneLine.Number)
		}
	}

	payload := &internalmessages.TransportationOffice{
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

func TransportationOffices(transportationOffices models.TransportationOffices) internalmessages.TransportationOffices {
	payload := make(internalmessages.TransportationOffices, len(transportationOffices))

	for i, to := range transportationOffices {
		transportationOffice := to
		payload[i] = TransportationOffice(transportationOffice)
	}
	return payload
}

// OfficeUser internal payload
func OfficeUser(officeUser *models.OfficeUser) *internalmessages.OfficeUser {
	if officeUser == nil || officeUser.ID == uuid.Nil {
		return nil
	}

	payload := &internalmessages.OfficeUser{
		ID:                   strfmt.UUID(officeUser.ID.String()),
		UserID:               strfmt.UUID(officeUser.UserID.String()),
		Email:                &officeUser.Email,
		FirstName:            &officeUser.FirstName,
		LastName:             &officeUser.LastName,
		MiddleName:           officeUser.MiddleInitials,
		Telephone:            &officeUser.Telephone,
		TransportationOffice: TransportationOffice(officeUser.TransportationOffice),
		CreatedAt:            strfmt.DateTime(officeUser.CreatedAt),
		UpdatedAt:            strfmt.DateTime(officeUser.UpdatedAt),
	}

	return payload
}

// MTOShipments payload
func MTOShipments(storer storage.FileStorer, mtoShipments *models.MTOShipments) *internalmessages.MTOShipments {
	payload := make(internalmessages.MTOShipments, len(*mtoShipments))

	for i, m := range *mtoShipments {
		copyOfMtoShipment := m // Make copy to avoid implicit memory aliasing of items from a range statement.
		payload[i] = MTOShipment(storer, &copyOfMtoShipment)
	}
	return &payload
}

// InternalServerError describes errors in a standard structure to be returned in the payload.
// If detail is nil, string defaults to "An internal server error has occurred."
func InternalServerError(detail *string, traceID uuid.UUID) *internalmessages.Error {
	payload := internalmessages.Error{
		Title:    handlers.FmtString(handlers.InternalServerErrMessage),
		Detail:   handlers.FmtString(handlers.InternalServerErrDetail),
		Instance: strfmt.UUID(traceID.String()),
	}
	if detail != nil {
		payload.Detail = detail
	}
	return &payload
}

// ValidationError describes validation errors from the model or properties
func ValidationError(detail string, instance uuid.UUID, validationErrors *validate.Errors) *internalmessages.ValidationError {
	payload := &internalmessages.ValidationError{
		ClientError: *ClientError(handlers.ValidationErrMessage, detail, instance),
	}
	if validationErrors != nil {
		payload.InvalidFields = handlers.NewValidationErrorListResponse(validationErrors).Errors
	}
	return payload
}

// ClientError describes errors in a standard structure to be returned in the payload
func ClientError(title string, detail string, instance uuid.UUID) *internalmessages.ClientError {
	return &internalmessages.ClientError{
		Title:    handlers.FmtString(title),
		Detail:   handlers.FmtString(detail),
		Instance: handlers.FmtUUID(instance),
	}
}

func PayloadForDocumentModel(storer storage.FileStorer, document models.Document) (*internalmessages.Document, error) {
	uploads := make([]*internalmessages.Upload, len(document.UserUploads))
	for i, userUpload := range document.UserUploads {
		if userUpload.Upload.ID == uuid.Nil {
			return nil, errors.New("no uploads for user")
		}
		url, err := storer.PresignedURL(userUpload.Upload.StorageKey, userUpload.Upload.ContentType)
		if err != nil {
			return nil, err
		}

		uploadPayload := PayloadForUploadModel(storer, userUpload.Upload, url)
		uploads[i] = uploadPayload
	}

	documentPayload := &internalmessages.Document{
		ID:              handlers.FmtUUID(document.ID),
		ServiceMemberID: handlers.FmtUUID(document.ServiceMemberID),
		Uploads:         uploads,
	}
	return documentPayload, nil
}

func PayloadForUploadModel(
	storer storage.FileStorer,
	upload models.Upload,
	url string,
) *internalmessages.Upload {
	uploadPayload := &internalmessages.Upload{
		ID:          handlers.FmtUUIDValue(upload.ID),
		Filename:    upload.Filename,
		ContentType: upload.ContentType,
		URL:         strfmt.URI(url),
		Bytes:       upload.Bytes,
		CreatedAt:   strfmt.DateTime(upload.CreatedAt),
		UpdatedAt:   strfmt.DateTime(upload.UpdatedAt),
	}
	tags, err := storer.Tags(upload.StorageKey)
	if err != nil || len(tags) == 0 {
		uploadPayload.Status = "PROCESSING"
	} else {
		uploadPayload.Status = tags["av-status"]
	}
	return uploadPayload
}

// MovingExpense payload
func MovingExpense(storer storage.FileStorer, movingExpense *models.MovingExpense) *internalmessages.MovingExpense {

	document, err := PayloadForDocumentModel(storer, movingExpense.Document)
	if err != nil {
		return nil
	}

	payload := &internalmessages.MovingExpense{
		ID:             *handlers.FmtUUID(movingExpense.ID),
		PpmShipmentID:  *handlers.FmtUUID(movingExpense.PPMShipmentID),
		DocumentID:     *handlers.FmtUUID(movingExpense.DocumentID),
		Document:       document,
		CreatedAt:      strfmt.DateTime(movingExpense.CreatedAt),
		UpdatedAt:      strfmt.DateTime(movingExpense.UpdatedAt),
		Description:    movingExpense.Description,
		PaidWithGtcc:   movingExpense.PaidWithGTCC,
		Amount:         handlers.FmtCost(movingExpense.Amount),
		MissingReceipt: movingExpense.MissingReceipt,
		ETag:           etag.GenerateEtag(movingExpense.UpdatedAt),
	}
	if movingExpense.MovingExpenseType != nil {
		movingExpenseType := internalmessages.OmittableMovingExpenseType(*movingExpense.MovingExpenseType)
		payload.MovingExpenseType = &movingExpenseType
	}

	if movingExpense.Status != nil {
		status := internalmessages.OmittablePPMDocumentStatus(*movingExpense.Status)
		payload.Status = &status
	}

	if movingExpense.Reason != nil {
		reason := internalmessages.PPMDocumentStatusReason(*movingExpense.Reason)
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
		sitLocation := internalmessages.SITLocationType(*movingExpense.SITLocation)
		payload.SitLocation = &sitLocation
	}

	return payload
}

func MovingExpenses(storer storage.FileStorer, movingExpenses models.MovingExpenses) []*internalmessages.MovingExpense {
	payload := make([]*internalmessages.MovingExpense, len(movingExpenses))
	for i, movingExpense := range movingExpenses {
		copyOfMovingExpense := movingExpense
		payload[i] = MovingExpense(storer, &copyOfMovingExpense)
	}
	return payload
}

func WeightTickets(storer storage.FileStorer, weightTickets models.WeightTickets) []*internalmessages.WeightTicket {
	payload := make([]*internalmessages.WeightTicket, len(weightTickets))
	for i, weightTicket := range weightTickets {
		copyOfWeightTicket := weightTicket
		weightTicketPayload := WeightTicket(storer, &copyOfWeightTicket)
		payload[i] = weightTicketPayload
	}
	return payload
}

// WeightTicket payload
func WeightTicket(storer storage.FileStorer, weightTicket *models.WeightTicket) *internalmessages.WeightTicket {
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

	payload := &internalmessages.WeightTicket{
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
		NetWeightRemarks:                  weightTicket.NetWeightRemarks,
		ETag:                              etag.GenerateEtag(weightTicket.UpdatedAt),
	}

	if weightTicket.Status != nil {
		status := internalmessages.OmittablePPMDocumentStatus(*weightTicket.Status)
		payload.Status = &status
	}

	if weightTicket.Reason != nil {
		reason := internalmessages.PPMDocumentStatusReason(*weightTicket.Reason)
		payload.Reason = &reason
	}

	return payload
}

// ProGearWeightTickets sets up a ProGearWeightTicket slice for the api using model data.
func ProGearWeightTickets(storer storage.FileStorer, proGearWeightTickets models.ProgearWeightTickets) []*internalmessages.ProGearWeightTicket {
	payload := make([]*internalmessages.ProGearWeightTicket, len(proGearWeightTickets))
	for i, proGearWeightTicket := range proGearWeightTickets {
		copyOfProGearWeightTicket := proGearWeightTicket
		proGearWeightTicketPayload := ProGearWeightTicket(storer, &copyOfProGearWeightTicket)
		payload[i] = proGearWeightTicketPayload
	}
	return payload
}

// ProGearWeightTicket payload
func ProGearWeightTicket(storer storage.FileStorer, progear *models.ProgearWeightTicket) *internalmessages.ProGearWeightTicket {
	ppmShipmentID := strfmt.UUID(progear.PPMShipmentID.String())

	document, err := PayloadForDocumentModel(storer, progear.Document)
	if err != nil {
		return nil
	}

	payload := &internalmessages.ProGearWeightTicket{
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
		status := internalmessages.OmittablePPMDocumentStatus(*progear.Status)
		payload.Status = &status
	}

	if progear.Reason != nil {
		reason := internalmessages.PPMDocumentStatusReason(*progear.Reason)
		payload.Reason = &reason
	}

	return payload
}

// SignedCertification converts a model to the api payload type
func SignedCertification(signedCertification *models.SignedCertification) *internalmessages.SignedCertification {
	if signedCertification == nil {
		return nil
	}

	model := &internalmessages.SignedCertification{
		ID:                handlers.FmtUUIDValue(signedCertification.ID),
		SubmittingUserID:  handlers.FmtUUIDValue(signedCertification.SubmittingUserID),
		MoveID:            handlers.FmtUUIDValue(signedCertification.MoveID),
		PpmID:             handlers.FmtUUIDPtr(signedCertification.PpmID),
		CertificationText: &signedCertification.CertificationText,
		Signature:         &signedCertification.Signature,
		Date:              handlers.FmtDate(signedCertification.Date),
		CreatedAt:         strfmt.DateTime(signedCertification.CreatedAt),
		UpdatedAt:         strfmt.DateTime(signedCertification.UpdatedAt),
		ETag:              etag.GenerateEtag(signedCertification.UpdatedAt),
	}

	// CertificationType is required from the api perspective, but at the model and DB level, it's nullable. In
	// practice, it shouldn't ever actually be null though, so we should always be matching the API spec, but
	// regardless, we need to do this nil check. It would be good to go back and make it required in the model/table.
	if signedCertification.CertificationType != nil {
		model.CertificationType = internalmessages.SignedCertificationType(*signedCertification.CertificationType)
	}

	return model
}
