package ghcapi

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
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

// Handle getting the information of a specific customer
func (h GetCustomerInfoHandler) Handle(params customercodeop.GetCustomerInfoParams) middleware.Responder {
	// for now just return static data
	customer := &ghcmessages.Customer{
		FirstName:              models.StringPointer("First"),
		MiddleName:             models.StringPointer("Middle"),
		LastName:               models.StringPointer("Last"),
		Agency:                 models.StringPointer("Agency"),
		Grade:                  models.StringPointer("Grade"),
		Email:                  models.StringPointer("Example@example.com"),
		Telephone:              models.StringPointer("213-213-3232"),
		OriginDutyStation:      models.StringPointer("Origin Station"),
		DestinationDutyStation: models.StringPointer("Destination Station"),
		DependentsAuthorized:   true,
	}
	return customercodeop.NewGetCustomerInfoOK().WithPayload(customer)
}

func payloadForCustomerMoveItem(CustomerMoveItem models.CustomerMoveItem) *ghcmessages.CustomerMoveItem {
	CustomerMoveItemPayload := ghcmessages.CustomerMoveItem{
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
