package supportapi

import (
	"encoding/json"
	"fmt"

	"github.com/gofrs/uuid"
	"go.uber.org/zap"

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
			return movetaskorderops.NewGetMoveTaskOrderNotFound().WithPayload(&supportmessages.Error{Message: handlers.FmtString(err.Error())})
		case services.InvalidInputError:
			return movetaskorderops.NewGetMoveTaskOrderBadRequest().WithPayload(&supportmessages.Error{Message: handlers.FmtString(err.Error())})
		default:
			return movetaskorderops.NewGetMoveTaskOrderInternalServerError().WithPayload(&supportmessages.Error{Message: handlers.FmtString(err.Error())})
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

	moveTaskOrder, err := createMoveTaskOrderAndChildren(h, params, logger)

	if err != nil {
		errorForPayload := supportmessages.Error{Message: handlers.FmtString(err.Error())}

		var returnOp middleware.Responder
		var detailedErrMsg string
		switch typedErr := err.(type) {
		case services.NotFoundError:
			return movetaskorderops.NewCreateMoveTaskOrderNotFound().WithPayload(errorForPayload)
		case services.InvalidInputError:
			return movetaskorderops.NewCreateMoveTaskOrderBadRequest().WithPayload(errorForPayload)
		case services.CreateObjectError:
			detailedErrMsg = typedErr.DetailedMsg()
			returnOp = movetaskorderops.NewCreateMoveTaskOrderInternalServerError().WithPayload(errorForPayload)
			break
		default:
			return movetaskorderops.NewCreateMoveTaskOrderInternalServerError().WithPayload(errorForPayload)
		}
		logger.Error("supportapi.CreateMoveTaskOrderHandler", zap.String("Error", detailedErrMsg))

		return returnOp
	}
	moveTaskOrderPayload := payloads.MoveTaskOrder(moveTaskOrder)
	return movetaskorderops.NewCreateMoveTaskOrderCreated().WithPayload(moveTaskOrderPayload)

}

// createMoveTaskOrderAndChildren creates a move task order - this is a support function do not use in production
func createMoveTaskOrderAndChildren(h CreateMoveTaskOrderHandler, params movetaskorderops.CreateMoveTaskOrderParams, logger handlers.Logger) (*models.MoveTaskOrder, error) {

	payload := params.Body

	// Create or get customer
	customer, err := createOrGetCustomer(h, payload.MoveOrder.CustomerID.String(), payload.MoveOrder.Customer, logger)
	if err != nil {
		fmt.Println("returning here")
		data, _ := json.Marshal(err)
		fmt.Printf("\n\n >>1 %s\n", data)

		return nil, err
	}
	fmt.Println("\n\n >> Customer created! ", *customer.FirstName, *customer.LastName)
	fmt.Println("\n\n --")

	moveOrder, err := createMoveOrder(h, customer, payload.MoveOrder, logger)
	if err != nil {
		return nil, err
	}
	fmt.Println("\n\n >> moveOrder created", *moveOrder.Grade)
	fmt.Println("\n\n --")

	moveTaskOrder := payloads.MoveTaskOrderModel(payload)
	moveTaskOrder.MoveOrder = *moveOrder
	moveTaskOrder.MoveOrderID = moveOrder.ID

	// Creates the moveOrder and the entitlement at the same time
	verrs, err := h.DB().ValidateAndCreate(moveTaskOrder)
	if err != nil || verrs.Count() > 0 {
		return nil, services.NewCreateObjectError("MoveTaskOrder", err, verrs, "")
	}
	fmt.Println("\n\n >> moveTaskOrder created", moveTaskOrder.ID.String())
	fmt.Println("\n\n --")
	return moveTaskOrder, nil

}

// createMoveOrder creates a basic move order - this is a support function do not use in production
func createMoveOrder(h CreateMoveTaskOrderHandler, customer *models.Customer, moveOrderPayload *supportmessages.MoveOrder, logger handlers.Logger) (*models.MoveOrder, error) {
	if moveOrderPayload == nil {
		returnErr := services.NewInvalidInputError(uuid.Nil, nil, nil, "MoveOrder definition is required to create MoveTaskOrder")
		return nil, returnErr
	}

	// We need to create an entitlement
	// let's try to create both with eager.
	// assume true
	moveOrder := payloads.MoveOrderModel(moveOrderPayload)
	moveOrder.Entitlement = payloads.EntitlementModel(moveOrderPayload.Entitlement)

	// Check if dutystations are valid uuids
	// Doesn't check if in database, this will be automatically checked on creation of mO
	if *moveOrder.DestinationDutyStationID == uuid.Nil ||
		*moveOrder.OriginDutyStationID == uuid.Nil {
		returnErr := services.NewInvalidInputError(uuid.Nil, nil, nil, "MoveOrder must contain valid destination and origin duty station UUIDs")
		return nil, returnErr
	}

	// Add customer to mO
	moveOrder.Customer = customer
	moveOrder.CustomerID = &customer.ID

	data, _ := json.Marshal(moveOrder)
	fmt.Printf("\n\n >> moveOrder model : \n%s\n", data)

	// Creates the moveOrder and the entitlement at the same time
	verrs, err := h.DB().Eager().ValidateAndCreate(moveOrder)
	if err != nil || verrs.Count() > 0 {
		return nil, services.NewCreateObjectError("MoveOrder", err, verrs, "")
	}
	return moveOrder, nil
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
		return nil, services.NewCreateObjectError("User", err, nil, "")
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

	customer := payloads.CustomerModel(customerBody)
	customer.User = *user
	customer.UserID = user.ID

	data, _ := json.Marshal(customer) //TODO remove
	fmt.Printf("\n\n >> %s\n", data)

	// Create the new customer in the db
	verrs, err := h.DB().ValidateAndCreate(customer)
	if err != nil || verrs.Count() > 0 {
		e := services.NewCreateObjectError("Customer", err, verrs, "")

		return nil, e
	}
	return customer, nil

}
