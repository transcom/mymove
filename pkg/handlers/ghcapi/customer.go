package ghcapi

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/gofrs/uuid"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/handlers/ghcapi/internal/payloads"

	customercodeop "github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations/customer"
	"github.com/transcom/mymove/pkg/gen/ghcmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
)

// GetCustomerInfoHandler fetches the information of a specific customer
type GetCustomerInfoHandler struct {
	handlers.HandlerContext
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
		logger.Error("Loading CustomerInfo Info", zap.Error(err))
		return handlers.ResponseForError(logger, err)
	}
	customerInfoPayload := payloads.CustomerInfo(customer)
	return customercodeop.NewGetCustomerInfoOK().WithPayload(customerInfoPayload)
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
		logger.Error("Loading CustomerInfo Moves", zap.Error(err))
		return handlers.ResponseForError(logger, err)
	}

	CustomerMoveItemPayloads := make([]*ghcmessages.CustomerMoveItem, len(CustomerMoveItems))
	for i, MoveQueueItem := range CustomerMoveItems {
		MoveQueueItemPayload := payloads.CustomerMoveItem(MoveQueueItem)
		CustomerMoveItemPayloads[i] = MoveQueueItemPayload
	}

	return customercodeop.NewGetAllCustomerMovesOK().WithPayload(CustomerMoveItemPayloads)
}
