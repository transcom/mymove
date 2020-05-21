package supportapi

import (
	"github.com/gobuffalo/pop"
	"github.com/gofrs/uuid"
	"github.com/transcom/mymove/pkg/services/support"
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
		switch typedErr := err.(type) {
		case services.NotFoundError:
			return movetaskorderops.NewUpdateMoveTaskOrderStatusNotFound().WithPayload(
				payloads.ClientError(handlers.NotFoundMessage, *handlers.FmtString(err.Error()), h.GetTraceID()))
		case services.InvalidInputError:
			return movetaskorderops.NewCreateMoveTaskOrderUnprocessableEntity().WithPayload(
				payloads.ValidationError(err.Error(), h.GetTraceID(), typedErr.ValidationErrors))
		case services.PreconditionFailedError:
			return movetaskorderops.NewUpdateMoveTaskOrderStatusPreconditionFailed().WithPayload(
				payloads.ClientError(handlers.PreconditionErrMessage, *handlers.FmtString(err.Error()), h.GetTraceID()))
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
		switch typedErr := err.(type) {
		case services.NotFoundError:
			return movetaskorderops.NewGetMoveTaskOrderNotFound().WithPayload(
				payloads.ClientError(handlers.NotFoundMessage, *handlers.FmtString(err.Error()), h.GetTraceID()))
		case services.InvalidInputError:
			return movetaskorderops.NewCreateMoveTaskOrderUnprocessableEntity().WithPayload(
				payloads.ValidationError(err.Error(), h.GetTraceID(), typedErr.ValidationErrors))
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
	support.InternalMoveTaskOrderCreator
}

// Handle updates to move task order post-counseling
func (h CreateMoveTaskOrderHandler) Handle(params movetaskorderops.CreateMoveTaskOrderParams) middleware.Responder {
	logger := h.LoggerFromRequest(params.HTTPRequest)

	moveTaskOrder, err := createMoveTaskOrderSupport(h, params, logger)

	if err != nil {
		switch typedErr := err.(type) {
		case services.NotFoundError:
			return movetaskorderops.NewCreateMoveTaskOrderNotFound().WithPayload(
				payloads.ClientError(handlers.NotFoundMessage, *handlers.FmtString(err.Error()), h.GetTraceID()))
		case services.InvalidInputError:
			errPayload := payloads.ValidationError(err.Error(), h.GetTraceID(), typedErr.ValidationErrors)
			return movetaskorderops.NewCreateMoveTaskOrderUnprocessableEntity().WithPayload(errPayload)
		case services.QueryError:
			// This error is generated when the validation passed but there was an error in creation
			// Usually this is due to a more complex dependency like a foreign key constraint
			return movetaskorderops.NewCreateMoveTaskOrderBadRequest().WithPayload(
				payloads.ClientError(handlers.SQLErrMessage, *handlers.FmtString(err.Error()), h.GetTraceID()))
		default:
			return movetaskorderops.NewCreateMoveTaskOrderInternalServerError().WithPayload(&supportmessages.Error{Message: handlers.FmtString(err.Error())})
		}
	}
	moveTaskOrderPayload := payloads.MoveTaskOrder(moveTaskOrder)
	return movetaskorderops.NewCreateMoveTaskOrderCreated().WithPayload(moveTaskOrderPayload)

}

// createMoveTaskOrderSupport creates a move task order - this is a support function do not use in production
// It creates customers, users, move orders as well.
func createMoveTaskOrderSupport(h CreateMoveTaskOrderHandler, params movetaskorderops.CreateMoveTaskOrderParams, logger handlers.Logger) (*models.MoveTaskOrder, error) {

	var moveTaskOrder *models.MoveTaskOrder
	payload := params.Body

	if payload.MoveOrder == nil {
		return nil, services.NewQueryError("MoveTaskOrder", nil, "MoveOrder is necessary")
	}

	transactionError := h.DB().Transaction(func(tx *pop.Connection) error {
		// Create or get customer
		customer, err := createOrGetCustomer(tx, h.CustomerFetcher, payload.MoveOrder.CustomerID.String(), payload.MoveOrder.Customer, logger)
		if err != nil {
			return err
		}

		// Create move order and entitlement
		moveOrder, err := createMoveOrder(tx, customer, payload.MoveOrder, logger)
		if err != nil {
			return err
		}

		// Create move task order
		moveTaskOrder = payloads.MoveTaskOrderModel(payload)
		moveTaskOrder.MoveOrder = *moveOrder
		moveTaskOrder.MoveOrderID = moveOrder.ID

		// Creates the moveOrder and the entitlement at the same time
		verrs, err := tx.ValidateAndCreate(moveTaskOrder)
		if verrs.Count() > 0 {
			logger.Error("supportapi.createMoveTaskOrderSupport error", zap.Error(verrs))
			return services.NewInvalidInputError(uuid.Nil, nil, verrs, "")
		} else if err != nil {
			logger.Error("supportapi.createMoveTaskOrderSupport error", zap.Error(err))
			return services.NewQueryError("MoveTaskOrder", err, "Unable to create MoveTaskOrder.")
		}
		return nil
	})

	if transactionError != nil {
		return nil, transactionError
	}
	return moveTaskOrder, nil

}

// createMoveOrder creates a basic move order - this is a support function do not use in production
func createMoveOrder(tx *pop.Connection, customer *models.Customer, moveOrderPayload *supportmessages.MoveOrder, logger handlers.Logger) (*models.MoveOrder, error) {
	if moveOrderPayload == nil {
		returnErr := services.NewInvalidInputError(uuid.Nil, nil, nil, "MoveOrder definition is required to create MoveTaskOrder")
		return nil, returnErr
	}

	// Move order model will also contain the entitlement
	moveOrder := payloads.MoveOrderModel(moveOrderPayload)

	// Add customer to mO
	moveOrder.Customer = customer
	moveOrder.CustomerID = &customer.ID

	// Creates the moveOrder and the entitlement at the same time
	verrs, err := tx.Eager().ValidateAndCreate(moveOrder)
	if verrs.Count() > 0 {
		logger.Error("supportapi.createMoveOrder error", zap.Error(verrs))
		return nil, services.NewInvalidInputError(uuid.Nil, nil, verrs, "")
	} else if err != nil {
		logger.Error("supportapi.createMoveOrder error", zap.Error(err))
		e := services.NewQueryError("MoveOrder", err, "Unable to create MoveOrder.")
		return nil, e
	}
	return moveOrder, nil
}

// createUser creates a user but this is a fake login.gov user
// this is support code only, do not use in a production case
func createUser(tx *pop.Connection, userEmail *string, logger handlers.Logger) (*models.User, error) {
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
	verrs, err := tx.ValidateAndCreate(&user)
	if verrs.Count() > 0 {
		logger.Error("supportapi.createUser error", zap.Error(verrs))
		return nil, services.NewInvalidInputError(uuid.Nil, nil, verrs, "")
	} else if err != nil {
		logger.Error("supportapi.createUser error", zap.Error(err))
		e := services.NewQueryError("User", err, "Unable to create User.")
		return nil, e
	}
	return &user, nil
}

// createOrGetCustomer creates a customer or gets one if id was provided
func createOrGetCustomer(tx *pop.Connection, f services.CustomerFetcher, customerIDString string, customerBody *supportmessages.Customer, logger handlers.Logger) (*models.Customer, error) {
	// If customer ID string is provided, we should find this customer
	if customerIDString != "" {
		customerID, err := uuid.FromString(customerIDString)
		// Error on bad customer id string
		if err != nil {
			returnErr := services.NewInvalidInputError(uuid.Nil, err, nil, "Invalid customerID: params CustomerID cannot be converted to a UUID")
			return nil, returnErr
		}
		// Find customer and return
		customer, err := f.FetchCustomer(customerID)
		if err != nil {
			returnErr := services.NewNotFoundError(customerID, "Customer with that ID not found")
			return nil, returnErr
		}
		return customer, nil

	}
	// Else customerIDString is empty and we need to create a customer
	// Since each customer has a unique userid we need to create a user
	if customerBody == nil {
		returnErr := services.NewInvalidInputError(uuid.Nil, nil, nil, "If CustomerID is not provided, customer object is required to create Customer")
		return nil, returnErr
	}
	user, err := createUser(tx, customerBody.Email, logger)
	if err != nil {
		return nil, err
	}

	customer := payloads.CustomerModel(customerBody)
	customer.User = *user
	customer.UserID = user.ID

	// Create the new customer in the db
	verrs, err := tx.ValidateAndCreate(customer)
	if verrs.Count() > 0 {
		logger.Error("supportapi.createOrGetCustomer error", zap.Error(verrs))
		return nil, services.NewInvalidInputError(uuid.Nil, nil, verrs, "")
	} else if err != nil {
		logger.Error("supportapi.createOrGetCustomer error", zap.Error(err))
		e := services.NewQueryError("Customer", err, "Unable to create Customer")
		return nil, e
	}
	return customer, nil

}
