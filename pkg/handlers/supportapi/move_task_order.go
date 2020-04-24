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
		logger.Error("supportapi.CreateMoveTaskOrderHandler error", zap.Error(err))
		errorForPayload := supportmessages.Error{Message: handlers.FmtString(err.Error())}
		// if errnew, ok := err.(services.CreateObjectError); ok {
		// 	fmt.Println("new code", errnew.Unwrap().Error())
		// 	// query failed because of a permission problem
		// } else {
		// 	fmt.Println("is isn't!")
		// }

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
	fmt.Println(moveTaskOrder)
	return movetaskorderops.NewCreateMoveTaskOrderBadRequest()

}

// createMoveTaskOrderAndChildren creates a move task order - this is a support function do not use in production
func createMoveTaskOrderAndChildren(h CreateMoveTaskOrderHandler, params movetaskorderops.CreateMoveTaskOrderParams, logger handlers.Logger) (*models.MoveTaskOrder, error) {

	payload := params.Body

	// Create or get customer
	customer, err := createOrGetCustomer(h, payload.MoveOrder.CustomerID.String(), payload.MoveOrder.Customer, logger)
	if err != nil {
		return nil, err
	}
	fmt.Println("\n\n >>", *customer.FirstName, *customer.LastName)
	fmt.Println("\n\n --")

	moveOrder, err := createMoveOrder(h, customer, payload.MoveOrder, logger)
	if err != nil {
		return nil, err
	}
	fmt.Println("\n\n >>", *moveOrder.Grade)
	fmt.Println("\n\n --")

	return nil, err

}

// createMoveOrder creates a basic move order - this is a support function do not use in production
func createMoveOrder(h CreateMoveTaskOrderHandler, customer *models.Customer, moveOrderPayload *supportmessages.MoveOrder, logger handlers.Logger) (*models.MoveOrder, error) {
	if moveOrderPayload == nil {
	}
	// create entitlement and move order at same time
	// move order --> customer (done)
	// move order --> entitlement
	// move order --> dutystations (id only)
	moveOrder := payloads.MoveOrderModel(moveOrderPayload)
	moveOrder.Customer = customer
	moveOrder.CustomerID = &customer.ID
	ddsUUID, _ := uuid.FromString("025efeb9-bd3e-47ea-be4f-b41365ee123c")
	moveOrder.DestinationDutyStationID = &ddsUUID
	odsUUID, _ := uuid.FromString("db061ae3-5da1-40f2-8317-631318752247")
	moveOrder.OriginDutyStationID = &odsUUID
	eID, _ := uuid.FromString("3dd5b72a-75d2-492a-a260-fc7d97a57b13")
	moveOrder.EntitlementID = &eID

	data, _ := json.Marshal(moveOrder)
	fmt.Printf("\n\n >> moveOrder model : \n%s\n", data)

	verrs, err := h.DB().ValidateAndCreate(moveOrder)
	//err := h.DB().Create(moveOrder)
	if err != nil || verrs.Count() > 0 {
		returnErr := services.NewCreateObjectError("MoveOrder", err, verrs, "")
		fmt.Println("\n\n :: ", verrs.Count())
		for _, key := range verrs.Keys() {
			fmt.Println(" ::", verrs.Get(key))
		}
		return nil, returnErr
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
