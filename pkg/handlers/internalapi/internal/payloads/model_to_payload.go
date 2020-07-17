package payloads

import (
	"github.com/go-openapi/strfmt"
	"github.com/gobuffalo/validate"
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
		agents[i] = MTOAgent(&m)
	}

	return &agents
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
		Status:                   string(mtoShipment.Status),
		DestinationAddress:       Address(mtoShipment.DestinationAddress),
		SecondaryPickupAddress:   Address(mtoShipment.SecondaryPickupAddress),
		SecondaryDeliveryAddress: Address(mtoShipment.SecondaryDeliveryAddress),
		CreatedAt:                strfmt.DateTime(mtoShipment.CreatedAt),
		UpdatedAt:                strfmt.DateTime(mtoShipment.UpdatedAt),
		ETag:                     etag.GenerateEtag(mtoShipment.UpdatedAt),
	}

	if mtoShipment.ApprovedDate != nil {
		payload.ApprovedDate = strfmt.Date(*mtoShipment.ApprovedDate)
	}

	if mtoShipment.ScheduledPickupDate != nil {
		payload.ScheduledPickupDate = strfmt.Date(*mtoShipment.ScheduledPickupDate)
	}

	if mtoShipment.RequestedPickupDate != nil && !mtoShipment.RequestedPickupDate.IsZero() {
		payload.RequestedPickupDate = strfmt.Date(*mtoShipment.RequestedPickupDate)
	}

	if mtoShipment.ActualPickupDate != nil && !mtoShipment.ActualPickupDate.IsZero() {
		payload.ActualPickupDate = strfmt.Date(*mtoShipment.ActualPickupDate)
	}

	if mtoShipment.FirstAvailableDeliveryDate != nil && !mtoShipment.FirstAvailableDeliveryDate.IsZero() {
		payload.FirstAvailableDeliveryDate = strfmt.Date(*mtoShipment.FirstAvailableDeliveryDate)
	}

	if mtoShipment.RequiredDeliveryDate != nil && !mtoShipment.RequiredDeliveryDate.IsZero() {
		payload.RequiredDeliveryDate = strfmt.Date(*mtoShipment.RequiredDeliveryDate)
	}

	if mtoShipment.PrimeEstimatedWeight != nil && mtoShipment.PrimeEstimatedWeightRecordedDate != nil {
		payload.PrimeEstimatedWeight = int64(*mtoShipment.PrimeEstimatedWeight)
		payload.PrimeEstimatedWeightRecordedDate = strfmt.Date(*mtoShipment.PrimeEstimatedWeightRecordedDate)
	}

	if mtoShipment.PrimeActualWeight != nil {
		payload.PrimeActualWeight = int64(*mtoShipment.PrimeActualWeight)
	}

	return payload
}

// MTOShipments payload
func MTOShipments(mtoShipments *models.MTOShipments) *internalmessages.MTOShipments {
	payload := make(internalmessages.MTOShipments, len(*mtoShipments))

	for i, m := range *mtoShipments {
		payload[i] = MTOShipment(&m)
	}
	return &payload
}

// MTOServiceItem payload
func MTOServiceItem(mtoServiceItem *models.MTOServiceItem) internalmessages.MTOServiceItem {
	var payload internalmessages.MTOServiceItem

	// here we determine which payload model to use based on the re service code
	switch mtoServiceItem.ReService.Code {
	case models.ReServiceCodeDOFSIT:
		payload = &internalmessages.MTOServiceItemDOFSIT{
			ReServiceCode:    string(mtoServiceItem.ReService.Code),
			PickupPostalCode: mtoServiceItem.PickupPostalCode,
			Reason:           mtoServiceItem.Reason,
		}
	case models.ReServiceCodeDDFSIT:
		firstContact := getCustomerContact(mtoServiceItem.CustomerContacts, models.CustomerContactTypeFirst)
		secondContact := getCustomerContact(mtoServiceItem.CustomerContacts, models.CustomerContactTypeSecond)
		payload = &internalmessages.MTOServiceItemDDFSIT{
			ReServiceCode:               string(mtoServiceItem.ReService.Code),
			TimeMilitary1:               handlers.FmtString(firstContact.TimeMilitary),
			FirstAvailableDeliveryDate1: handlers.FmtDate(firstContact.FirstAvailableDeliveryDate),
			TimeMilitary2:               handlers.FmtString(secondContact.TimeMilitary),
			FirstAvailableDeliveryDate2: handlers.FmtDate(secondContact.FirstAvailableDeliveryDate),
		}
	case models.ReServiceCodeDCRT, models.ReServiceCodeDUCRT, models.ReServiceCodeDCRTSA:
		item := getDimension(mtoServiceItem.Dimensions, models.DimensionTypeItem)
		crate := getDimension(mtoServiceItem.Dimensions, models.DimensionTypeCrate)
		payload = &internalmessages.MTOServiceItemDomesticCrating{
			ReServiceCode: handlers.FmtString(string(mtoServiceItem.ReService.Code)),
			Item: &internalmessages.MTOServiceItemDimension{
				ID:     strfmt.UUID(item.ID.String()),
				Type:   internalmessages.DimensionType(item.Type),
				Height: item.Height.Int32Ptr(),
				Length: item.Length.Int32Ptr(),
				Width:  item.Width.Int32Ptr(),
			},
			Crate: &internalmessages.MTOServiceItemDimension{
				ID:     strfmt.UUID(crate.ID.String()),
				Type:   internalmessages.DimensionType(crate.Type),
				Height: crate.Height.Int32Ptr(),
				Length: crate.Length.Int32Ptr(),
				Width:  crate.Width.Int32Ptr(),
			},
			Description: mtoServiceItem.Description,
		}
	case models.ReServiceCodeDDSHUT, models.ReServiceCodeDOSHUT:
		payload = &internalmessages.MTOServiceItemShuttle{
			Description:   mtoServiceItem.Description,
			ReServiceCode: handlers.FmtString(string(mtoServiceItem.ReService.Code)),
			Reason:        mtoServiceItem.Reason,
		}
	default:
		// otherwise, basic service item
		payload = &internalmessages.MTOServiceItemBasic{
			ReServiceCode: internalmessages.ReServiceCode(mtoServiceItem.ReService.Code),
		}
	}

	// set all relevant fields that apply to all service items
	var shipmentIDStr string
	if mtoServiceItem.MTOShipmentID != nil {
		shipmentIDStr = mtoServiceItem.MTOShipmentID.String()
	}

	payload.SetID(strfmt.UUID(mtoServiceItem.ID.String()))
	payload.SetMoveTaskOrderID(handlers.FmtUUID(mtoServiceItem.MoveTaskOrderID))
	payload.SetMtoShipmentID(strfmt.UUID(shipmentIDStr))
	payload.SetReServiceName(mtoServiceItem.ReService.Name)
	payload.SetStatus(internalmessages.MTOServiceItemStatus(mtoServiceItem.Status))
	payload.SetRejectionReason(mtoServiceItem.Reason)
	payload.SetETag(etag.GenerateEtag(mtoServiceItem.UpdatedAt))
	return payload
}

// MTOServiceItems payload
func MTOServiceItems(mtoServiceItems *models.MTOServiceItems) *[]internalmessages.MTOServiceItem {
	var payload []internalmessages.MTOServiceItem

	for _, p := range *mtoServiceItems {
		payload = append(payload, MTOServiceItem(&p))
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

// getDimension will get the first dimension of the passed in type.
func getDimension(dimensions models.MTOServiceItemDimensions, dimensionType models.DimensionType) models.MTOServiceItemDimension {
	if len(dimensions) == 0 {
		return models.MTOServiceItemDimension{}
	}

	for _, dimension := range dimensions {
		if dimension.Type == dimensionType {
			return dimension
		}
	}

	return models.MTOServiceItemDimension{}
}

// getFirstCustomerContact will get the first customer contact for destination 1st day SIT based on type.
func getCustomerContact(customerContacts models.MTOServiceItemCustomerContacts, customerContactType models.CustomerContactType) models.MTOServiceItemCustomerContact {
	if len(customerContacts) == 0 {
		return models.MTOServiceItemCustomerContact{}
	}

	for _, customerContact := range customerContacts {
		if customerContact.Type == customerContactType {
			return customerContact
		}
	}

	return models.MTOServiceItemCustomerContact{}
}
