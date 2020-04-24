package supportapi

import (
	"encoding/json"
	"fmt"

	"github.com/gofrs/uuid"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/gen/primemessages"
	"github.com/transcom/mymove/pkg/gen/supportmessages"
	"github.com/transcom/mymove/pkg/models"

	"github.com/transcom/mymove/pkg/handlers/supportapi/internal/payloads"

	"github.com/go-openapi/runtime/middleware"

	movetaskorderops "github.com/transcom/mymove/pkg/gen/supportapi/supportoperations/move_task_order"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/services"
)

// UpdateMoveTaskOrderStatusHandlerFunc updates the status of a Move Task Order
type UpdateMoveTaskOrderStatusHandlerFunc struct {
	handlers.HandlerContext
	moveTaskOrderStatusUpdater services.MoveTaskOrderUpdater
}

// Handle updates the status of a MoveTaskOrder
func (h UpdateMoveTaskOrderStatusHandlerFunc) Handle(params movetaskorderops.UpdateMoveTaskOrderStatusParams) middleware.Responder {
	_, logger := h.SessionAndLoggerFromRequest(params.HTTPRequest)
	eTag := params.IfMatch

	moveTaskOrderID := uuid.FromStringOrNil(params.MoveTaskOrderID)

	mto, err := h.moveTaskOrderStatusUpdater.MakeAvailableToPrime(moveTaskOrderID, eTag)

	if err != nil {
		logger.Error("supportapi.MoveTaskOrderHandler error", zap.Error(err))
		switch err.(type) {
		case services.NotFoundError:
			return movetaskorderops.NewUpdateMoveTaskOrderStatusNotFound().WithPayload(&supportmessages.Error{Message: handlers.FmtString(err.Error())})
		case services.InvalidInputError:
			return movetaskorderops.NewUpdateMoveTaskOrderStatusBadRequest().WithPayload(&supportmessages.Error{Message: handlers.FmtString(err.Error())})
		case services.PreconditionFailedError:
			return movetaskorderops.NewUpdateMoveTaskOrderStatusPreconditionFailed().WithPayload(&supportmessages.Error{Message: handlers.FmtString(err.Error())})
		default:
			return movetaskorderops.NewUpdateMoveTaskOrderStatusInternalServerError().WithPayload(&supportmessages.Error{Message: handlers.FmtString(err.Error())})
		}
	}

	moveTaskOrderPayload := payloads.MoveTaskOrder(mto)

	return movetaskorderops.NewUpdateMoveTaskOrderStatusOK().WithPayload(moveTaskOrderPayload)
}

// GetMoveTaskOrderHandlerFunc updates the status of a Move Task Order
type GetMoveTaskOrderHandlerFunc struct {
	handlers.HandlerContext
	moveTaskOrderFetcher services.MoveTaskOrderFetcher
}

// Handle updates the status of a MoveTaskOrder
func (h GetMoveTaskOrderHandlerFunc) Handle(params movetaskorderops.GetMoveTaskOrderParams) middleware.Responder {
	logger := h.LoggerFromRequest(params.HTTPRequest)

	moveTaskOrderID := uuid.FromStringOrNil(params.MoveTaskOrderID)
	mto, err := h.moveTaskOrderFetcher.FetchMoveTaskOrder(moveTaskOrderID)
	if err != nil {
		logger.Error("primeapi.support.GetMoveTaskOrderHandler error", zap.Error(err))
		switch err.(type) {
		case services.NotFoundError:
			return movetaskorderops.NewGetMoveTaskOrderNotFound().WithPayload(&primemessages.Error{Message: handlers.FmtString(err.Error())})
		case services.InvalidInputError:
			return movetaskorderops.NewGetMoveTaskOrderBadRequest().WithPayload(&primemessages.Error{Message: handlers.FmtString(err.Error())})
		default:
			return movetaskorderops.NewGetMoveTaskOrderInternalServerError().WithPayload(&primemessages.Error{Message: handlers.FmtString(err.Error())})
		}
	}
	moveTaskOrderPayload := payloads.MoveTaskOrder(mto)
	return movetaskorderops.NewGetMoveTaskOrderOK().WithPayload(moveTaskOrderPayload)
}

// CreateMoveTaskOrderHandler creates a move task order
type CreateMoveTaskOrderHandler struct {
	handlers.HandlerContext
	services.CustomerFetcher
	services.MoveTaskOrderCreator
}

// Handle updates to move task order post-counseling
func (h CreateMoveTaskOrderHandler) Handle(params movetaskorderops.CreateMoveTaskOrderParams) middleware.Responder {
	logger := h.LoggerFromRequest(params.HTTPRequest)
	payload := params.Body

	// Create or get customer
	customer, err := createOrGetCustomer(h, payload.MoveOrder.CustomerID.String(), payload.MoveOrder.Customer, logger)
	if err == nil {
		fmt.Println("\n\n >>", *customer.FirstName, *customer.LastName)
		fmt.Println("\n\n --")
	}
	if err != nil {
		logger.Error("supportapi.CreateMoveTaskOrderHandler error", zap.Error(err))
		errMsg := primemessages.Error{Message: handlers.FmtString(err.Error())}
		switch err.(type) {
		case services.NotFoundError:
			return movetaskorderops.NewCreateMoveTaskOrderNotFound().WithPayload(errMsg)
		case services.InvalidInputError:
			return movetaskorderops.NewCreateMoveTaskOrderBadRequest().WithPayload(errMsg)
		default:
			return movetaskorderops.NewCreateMoveTaskOrderInternalServerError().WithPayload(errMsg)
		}
	}

	return movetaskorderops.NewCreateMoveTaskOrderCreated()

}

// createUser creates a user
func createUser(h CreateMoveTaskOrderHandler, userEmail *string, logger handlers.Logger) (*models.User, error) {
	if userEmail == nil {
		defaultEmail := "generatedMTOuser@example.com"
		userEmail = &defaultEmail
	}
	id := uuid.Must(uuid.NewV4())
	user := models.User{
		LoginGovUUID:  id,
		LoginGovEmail: *userEmail,
		Active:        true,
	}
	verrs, err := h.DB().ValidateAndCreate(&user)
	if err != nil || verrs.Count() != 0 {
		returnErr := services.NewCreateObjectError("User", err, nil, "Error creating a user")
		return nil, returnErr
	}
	return &user, nil
}

// createOrGetCustomer creates a customer or gets one if id was provided
func createOrGetCustomer(h CreateMoveTaskOrderHandler, customerIDString string, customerBody *supportmessages.Customer, logger handlers.Logger) (*models.Customer, error) {
	// If customer ID string is provided, we should find this customer
	if customerIDString != "" {
		customerID, err := uuid.FromString(customerIDString)
		// Error on bad customer id string
		if err != nil {
			returnErr := services.NewInvalidInputError(uuid.Nil, err, nil, "Invalid customerID: params CustomerID cannot be converted to a UUID")
			return nil, returnErr
		}
		// Find customer and return
		customer, err := h.FetchCustomer(customerID)
		if err != nil {
			returnErr := services.NewNotFoundError(customerID, "Customer with that ID not found")
			return nil, returnErr
		}
		return customer, nil

	}
	// Else customerIDString is empty and we need to create a customer
	// Since each customer has a unique userid we need to create a user
	user, err := createUser(h, customerBody.Email, logger)
	if err != nil {
		return nil, err
	}

	// Create the customer model and populate the new user
	customer := payloads.CustomerModel(customerBody)
	customer.User = *user
	customer.UserID = user.ID

	data, _ := json.Marshal(customer)
	fmt.Printf("\n\n >> %s\n", data)

	// Create the new customer in the db
	verrs, err := h.DB().ValidateAndCreate(customer)
	if err != nil || verrs.Count() > 0 {
		logger.Error("createOrGetCustomer", zap.String("Error", err.Error()))
		returnErr := services.NewCreateObjectError("Customer", err, nil, "Error creating a customer")
		return nil, returnErr
	}
	return customer, nil

}
