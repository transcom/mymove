package ghcapi

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/gofrs/uuid"
	"go.uber.org/zap"

	customercodeop "github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations/customer"
	"github.com/transcom/mymove/pkg/gen/ghcmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
)

// GetCustomerInfoHandler fetches the information of a specific customer
type GetCustomerInfoHandler struct {
	handlers.HandlerContext
}

func payloadForCustomerInfo(Customer models.Customer) *ghcmessages.Customer {
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

// Handle getting the information of a specific customer
func (h GetCustomerInfoHandler) Handle(params customercodeop.GetCustomerInfoParams) middleware.Responder {
	session, logger := h.SessionAndLoggerFromRequest(params.HTTPRequest)

	if !session.IsOfficeUser() {
		return customercodeop.NewGetCustomerInfoForbidden()
	}
	customerID, _ := uuid.FromString(params.CustomerID.String())

	customer, err := models.GetCustomerInfo(h.DB(), customerID)

	if err != nil {
		logger.Error("Loading Customer Info", zap.Error(err))
		return handlers.ResponseForError(logger, err)
	}
	customerInfoPayload := payloadForCustomerInfo(customer)
	return customercodeop.NewGetCustomerInfoOK().WithPayload(customerInfoPayload)
}

func payloadForCustomerMoveItem(CustomerMoveItem models.CustomerMoveItem) *ghcmessages.CustomerMoveItem {
	CustomerMoveItemPayload := ghcmessages.CustomerMoveItem{
		ID:                    *handlers.FmtUUID(CustomerMoveItem.ID),
		CustomerID:            *handlers.FmtUUID(CustomerMoveItem.CustomerID),
		CreatedAt:             strfmt.DateTime(CustomerMoveItem.CreatedAt),
		CustomerName:          swag.String(CustomerMoveItem.CustomerName),
		ConfirmationNumber:    CustomerMoveItem.ConfirmationNumber,
		BranchOfService:       CustomerMoveItem.BranchOfService,
		OriginDutyStationName: models.StringPointer(CustomerMoveItem.OriginDutyStationName),
	}
	return &CustomerMoveItemPayload
}

// GetAllCustomerMovesHandler fetches the information of a specific customer
type GetAllCustomerMovesHandler struct {
	handlers.HandlerContext
}

// Handle getting the information of all customers
func (h GetAllCustomerMovesHandler) Handle(params customercodeop.GetAllCustomerMovesParams) middleware.Responder {
	session, logger := h.SessionAndLoggerFromRequest(params.HTTPRequest)

	if !session.IsOfficeUser() {
		return customercodeop.NewGetAllCustomerMovesForbidden()
	}

	CustomerMoveItems, err := models.GetCustomerMoveItems(h.DB())

	if err != nil {
		logger.Error("Loading Customer Moves", zap.Error(err))
		return handlers.ResponseForError(logger, err)
	}

	CustomerMoveItemPayloads := make([]*ghcmessages.CustomerMoveItem, len(CustomerMoveItems))
	for i, MoveQueueItem := range CustomerMoveItems {
		MoveQueueItemPayload := payloadForCustomerMoveItem(MoveQueueItem)
		CustomerMoveItemPayloads[i] = MoveQueueItemPayload
	}

	return customercodeop.NewGetAllCustomerMovesOK().WithPayload(CustomerMoveItemPayloads)
}
