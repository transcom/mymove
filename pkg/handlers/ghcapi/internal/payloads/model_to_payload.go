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

func CustomerInfo(CustomerInfo models.CustomerInfo) *ghcmessages.Customer {
	CustomerInfoPayload := ghcmessages.Customer{
		ID:                     strfmt.UUID(CustomerInfo.ID.String()),
		CustomerName:           &CustomerInfo.CustomerName,
		Agency:                 &CustomerInfo.Agency,
		Grade:                  &CustomerInfo.Grade,
		Email:                  &CustomerInfo.Email,
		Telephone:              &CustomerInfo.Telephone,
		OriginDutyStation:      &CustomerInfo.OriginDutyStationName,
		DestinationDutyStation: &CustomerInfo.DestinationDutyStationName,
	}
	return &CustomerInfoPayload
}
