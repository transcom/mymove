package payloads

import (
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/etag"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
)

// Address payload
func Address(address *models.Address) *internalmessages.Address {
	if address == nil {
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
func PPMShipment(ppmShipment *models.PPMShipment) *internalmessages.PPMShipment {
	if ppmShipment == nil || ppmShipment.ID.IsNil() {
		return nil
	}

	payloadPPMShipment := &internalmessages.PPMShipment{
		ID:                    *handlers.FmtUUID(ppmShipment.ID),
		ShipmentID:            *handlers.FmtUUID(ppmShipment.ShipmentID),
		CreatedAt:             strfmt.DateTime(ppmShipment.CreatedAt),
		UpdatedAt:             strfmt.DateTime(ppmShipment.UpdatedAt),
		Status:                internalmessages.PPMShipmentStatus(ppmShipment.Status),
		ExpectedDepartureDate: handlers.FmtDate(ppmShipment.ExpectedDepartureDate),
		PickupPostalCode:      &ppmShipment.PickupPostalCode,
		DestinationPostalCode: &ppmShipment.DestinationPostalCode,
		SitExpected:           &ppmShipment.SitExpected,
		ETag:                  etag.GenerateEtag(ppmShipment.UpdatedAt),
	}

	if ppmShipment.SecondaryPickupPostalCode != nil {
		payloadPPMShipment.SecondaryPickupPostalCode = *ppmShipment.SecondaryPickupPostalCode
	}

	if ppmShipment.SecondaryDestinationPostalCode != nil {
		payloadPPMShipment.SecondaryDestinationPostalCode = *ppmShipment.SecondaryDestinationPostalCode
	}

	if ppmShipment.ActualMoveDate != nil && !ppmShipment.ActualMoveDate.IsZero() {
		payloadPPMShipment.ActualMoveDate = handlers.FmtDate(*ppmShipment.ActualMoveDate)
	}

	if ppmShipment.SubmittedAt != nil && !ppmShipment.SubmittedAt.IsZero() {
		payloadPPMShipment.SubmittedAt = handlers.FmtDateTimePtr(ppmShipment.SubmittedAt)
	}

	if ppmShipment.ReviewedAt != nil && !ppmShipment.ReviewedAt.IsZero() {
		payloadPPMShipment.ReviewedAt = handlers.FmtDateTimePtr(ppmShipment.ReviewedAt)
	}

	if ppmShipment.EstimatedWeight != nil {
		payloadPPMShipment.EstimatedWeight = int64(*ppmShipment.EstimatedWeight)
	}

	if ppmShipment.NetWeight != nil {
		payloadPPMShipment.NetWeight = int64(*ppmShipment.NetWeight)
	}

	if ppmShipment.HasProGear != nil {
		payloadPPMShipment.HasProGear = *ppmShipment.HasProGear
	}

	if ppmShipment.ProGearWeight != nil {
		payloadPPMShipment.ProGearWeight = int64(*ppmShipment.ProGearWeight)
	}

	if ppmShipment.SpouseProGearWeight != nil {
		payloadPPMShipment.SpouseProGearWeight = int64(*ppmShipment.SpouseProGearWeight)
	}

	if ppmShipment.EstimatedIncentive != nil {
		payloadPPMShipment.EstimatedIncentive = int64(*ppmShipment.EstimatedIncentive)
	}

	if ppmShipment.AdvanceRequested != nil {
		payloadPPMShipment.AdvanceRequested = *ppmShipment.AdvanceRequested
	}

	if ppmShipment.AdvanceID != nil {
		payloadPPMShipment.AdvanceID = *handlers.FmtUUIDPtr(ppmShipment.AdvanceID)
	}

	if ppmShipment.AdvanceWorksheetID != nil {
		payloadPPMShipment.AdvanceWorksheetID = *handlers.FmtUUIDPtr(ppmShipment.AdvanceWorksheetID)
	}

	return payloadPPMShipment
}

// MTOShipment payload
func MTOShipment(mtoShipment *models.MTOShipment) *internalmessages.MTOShipment {
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
		PpmShipment:              PPMShipment(mtoShipment.PPMShipment),
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
func MTOShipments(mtoShipments *models.MTOShipments) *internalmessages.MTOShipments {
	payload := make(internalmessages.MTOShipments, len(*mtoShipments))

	for i, m := range *mtoShipments {
		copyOfMtoShipment := m // Make copy to avoid implicit memory aliasing of items from a range statement.
		payload[i] = MTOShipment(&copyOfMtoShipment)
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
