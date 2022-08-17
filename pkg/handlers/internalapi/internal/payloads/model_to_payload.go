package payloads

import (
	"errors"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/etag"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
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
		PickupPostalCode:               &ppmShipment.PickupPostalCode,
		SecondaryPickupPostalCode:      ppmShipment.SecondaryPickupPostalCode,
		ActualPickupPostalCode:         ppmShipment.ActualPickupPostalCode,
		DestinationPostalCode:          &ppmShipment.DestinationPostalCode,
		SecondaryDestinationPostalCode: ppmShipment.SecondaryDestinationPostalCode,
		ActualDestinationPostalCode:    ppmShipment.ActualDestinationPostalCode,
		SitExpected:                    ppmShipment.SITExpected,
		EstimatedWeight:                handlers.FmtPoundPtr(ppmShipment.EstimatedWeight),
		EstimatedIncentive:             handlers.FmtCost(ppmShipment.EstimatedIncentive),
		NetWeight:                      handlers.FmtPoundPtr(ppmShipment.NetWeight),
		HasProGear:                     ppmShipment.HasProGear,
		ProGearWeight:                  handlers.FmtPoundPtr(ppmShipment.ProGearWeight),
		SpouseProGearWeight:            handlers.FmtPoundPtr(ppmShipment.SpouseProGearWeight),
		HasRequestedAdvance:            ppmShipment.HasRequestedAdvance,
		AdvanceAmountRequested:         handlers.FmtCost(ppmShipment.AdvanceAmountRequested),
		HasReceivedAdvance:             ppmShipment.HasReceivedAdvance,
		AdvanceAmountReceived:          handlers.FmtCost(ppmShipment.AdvanceAmountReceived),
		WeightTickets:                  WeightTickets(storer, ppmShipment.WeightTickets),
		ETag:                           etag.GenerateEtag(ppmShipment.UpdatedAt),
	}

	return payloadPPMShipment
}

// MTOShipment payload
func MTOShipment(storer storage.FileStorer, mtoShipment *models.MTOShipment) *internalmessages.MTOShipment {
	payload := &internalmessages.MTOShipment{
		ID:                       strfmt.UUID(mtoShipment.ID.String()),
		Agents:                   *MTOAgents(&mtoShipment.MTOAgents),
		MoveTaskOrderID:          strfmt.UUID(mtoShipment.MoveTaskOrderID.String()),
		ShipmentType:             internalmessages.MTOShipmentType(mtoShipment.ShipmentType),
		CustomerRemarks:          mtoShipment.CustomerRemarks,
		PickupAddress:            Address(mtoShipment.PickupAddress),
		SecondaryPickupAddress:   Address(mtoShipment.SecondaryPickupAddress),
		DestinationAddress:       Address(mtoShipment.DestinationAddress),
		SecondaryDeliveryAddress: Address(mtoShipment.SecondaryDeliveryAddress),
		CreatedAt:                strfmt.DateTime(mtoShipment.CreatedAt),
		UpdatedAt:                strfmt.DateTime(mtoShipment.UpdatedAt),
		Status:                   internalmessages.MTOShipmentStatus(mtoShipment.Status),
		PpmShipment:              PPMShipment(storer, mtoShipment.PPMShipment),
		ETag:                     etag.GenerateEtag(mtoShipment.UpdatedAt),
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
		Name:       swag.String(office.Name),
		Gbloc:      office.Gbloc,
		Address:    Address(&office.Address),
		PhoneLines: phoneLines,
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

func PayloadForDocumentModel(storer storage.FileStorer, document models.Document) (*internalmessages.DocumentPayload, error) {
	uploads := make([]*internalmessages.UploadPayload, len(document.UserUploads))
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

	documentPayload := &internalmessages.DocumentPayload{
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
) *internalmessages.UploadPayload {
	uploadPayload := &internalmessages.UploadPayload{
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

// MovingExpense payload
func MovingExpense(storer storage.FileStorer, movingExpense *models.MovingExpense) *internalmessages.MovingExpense {

	document, err := PayloadForDocumentModel(storer, movingExpense.Document)
	if err != nil {
		return nil
	}

	payload := &internalmessages.MovingExpense{
		ID:                *handlers.FmtUUID(movingExpense.ID),
		PpmShipmentID:     handlers.FmtUUID(movingExpense.PPMShipmentID),
		PpmShipment:       &internalmessages.PPMShipment{},
		DocumentID:        handlers.FmtUUID(movingExpense.DocumentID),
		Document:          document,
		CreatedAt:         handlers.FmtDateTime(movingExpense.CreatedAt),
		UpdatedAt:         handlers.FmtDateTime(movingExpense.UpdatedAt),
		MovingExpenseType: internalmessages.MovingExpenseType(*movingExpense.MovingExpenseType),
		Description:       movingExpense.Description,
		PaidWithGtcc:      movingExpense.PaidWithGTCC,
		Amount:            handlers.FmtCost(movingExpense.Amount),
		MissingReceipt:    movingExpense.MissingReceipt,
		Status:            internalmessages.MovingExpense{}.Status,
		Reason:            movingExpense.Reason,
		SitStartDate:      handlers.FmtDate(*movingExpense.SITStartDate),
		SitEndDate:        handlers.FmtDate(*movingExpense.SITEndDate),
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
		ETag:                              etag.GenerateEtag(weightTicket.UpdatedAt),
	}

	return payload
}
