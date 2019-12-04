package payloads

import (
	"github.com/go-openapi/strfmt"

	"github.com/transcom/mymove/pkg/gen/ghcmessages"
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

func CustomerInfo(Customer models.CustomerInfo) *ghcmessages.Customer {
	CustomerInfoPayload := ghcmessages.Customer{
		ID:                     strfmt.UUID(Customer.ID.String()),
		CustomerName:           &Customer.CustomerName,
		Agency:                 &Customer.Agency,
		Grade:                  &Customer.Grade,
		Email:                  &Customer.Email,
		Telephone:              &Customer.Telephone,
		OriginDutyStation:      &Customer.OriginDutyStationName,
		DestinationDutyStation: &Customer.DestinationDutyStationName,
	}
	return &CustomerInfoPayload
}
