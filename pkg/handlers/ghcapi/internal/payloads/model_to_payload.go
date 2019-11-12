package payloads

import (
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"

	"github.com/transcom/mymove/pkg/gen/ghcmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
)

func CustomerMoveItem(CustomerMoveItem models.CustomerMoveItem) *ghcmessages.CustomerMoveItem {
	CustomerMoveItemPayload := ghcmessages.CustomerMoveItem{
		ID:                    strfmt.UUID(CustomerMoveItem.ID.String()),
		CustomerID:            strfmt.UUID(CustomerMoveItem.CustomerID.String()),
		CreatedAt:             strfmt.DateTime(CustomerMoveItem.CreatedAt),
		CustomerName:          &CustomerMoveItem.CustomerName,
		ConfirmationNumber:    CustomerMoveItem.ConfirmationNumber,
		BranchOfService:       CustomerMoveItem.BranchOfService,
		OriginDutyStationName: &CustomerMoveItem.OriginDutyStationName,
		ReferenceID:           CustomerMoveItem.ReferenceID,
	}
	return &CustomerMoveItemPayload
}

func CustomerInfo(Customer models.Customer) *ghcmessages.Customer {
	CustomerInfoPayload := ghcmessages.Customer{
		ID:                     *handlers.FmtUUID(Customer.ID),
		CustomerName:           swag.String(Customer.CustomerName),
		Agency:                 swag.String(Customer.Agency),
		Grade:                  swag.String(Customer.Grade),
		Email:                  swag.String(Customer.Email),
		Telephone:              swag.String(Customer.Telephone),
		OriginDutyStation:      swag.String(Customer.OriginDutyStationName),
		DestinationDutyStation: swag.String(Customer.DestinationDutyStationName),
	}
	return &CustomerInfoPayload
}
