package primeapi

import (
	"fmt"
	"time"

	"github.com/transcom/mymove/pkg/services"

	"github.com/go-openapi/runtime/middleware"
	"github.com/gofrs/uuid"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/handlers/primeapi/internal/payloads"
	"github.com/transcom/mymove/pkg/models"

	movetaskorderops "github.com/transcom/mymove/pkg/gen/primeapi/primeoperations/move_task_order"
	"github.com/transcom/mymove/pkg/gen/primemessages"
	"github.com/transcom/mymove/pkg/handlers"
)

const (
	// MilMoveUserType is the type of user for a Service Member
	MilMoveUserType string = "milmove"
	// OfficeUserType is the type of user for an Office user
	OfficeUserType string = "office"
	// DpsUserType is the type of user for a DPS user
	DpsUserType string = "dps"
	// AdminUserType is the type of user for an admin user
	AdminUserType string = "admin"
)

// FetchMTOUpdatesHandler lists move task orders with the option to filter since a particular date
type FetchMTOUpdatesHandler struct {
	handlers.HandlerContext
}

// Handle fetches all move task orders with the option to filter since a particular date
func (h FetchMTOUpdatesHandler) Handle(params movetaskorderops.FetchMTOUpdatesParams) middleware.Responder {
	logger := h.LoggerFromRequest(params.HTTPRequest)

	var mtos models.MoveTaskOrders

	query := h.DB().Where("is_available_to_prime = ?", true).Eager(
		"PaymentRequests",
		"MTOServiceItems",
		"MTOServiceItems.ReService",
		"MTOServiceItems.Dimensions",
		"MTOServiceItems.CustomerContacts",
		"MTOShipments",
		"MTOShipments.DestinationAddress",
		"MTOShipments.PickupAddress",
		"MTOShipments.SecondaryDeliveryAddress",
		"MTOShipments.SecondaryPickupAddress",
		"MTOShipments.MTOAgents",
		"MoveOrder",
		"MoveOrder.Customer",
		"MoveOrder.Entitlement")
	if params.Since != nil {
		since := time.Unix(*params.Since, 0)
		query = query.Where("updated_at > ?", since)
	}

	err := query.All(&mtos)

	if err != nil {
		logger.Error("Unable to fetch records:", zap.Error(err))
		return movetaskorderops.NewFetchMTOUpdatesInternalServerError()
	}

	payload := payloads.MoveTaskOrders(&mtos)

	return movetaskorderops.NewFetchMTOUpdatesOK().WithPayload(payload)
}

// UpdateMTOPostCounselingInformationHandler updates the move task order with post-counseling information
type UpdateMTOPostCounselingInformationHandler struct {
	handlers.HandlerContext
	services.Fetcher
	services.MoveTaskOrderUpdater
}

// Handle updates to move task order post-counseling
func (h UpdateMTOPostCounselingInformationHandler) Handle(params movetaskorderops.UpdateMTOPostCounselingInformationParams) middleware.Responder {
	logger := h.LoggerFromRequest(params.HTTPRequest)
	mtoID := uuid.FromStringOrNil(params.MoveTaskOrderID)
	eTag := params.IfMatch
	logger.Info("primeapi.UpdateMTOShipmentHandler info", zap.String("pointOfContact", params.Body.PointOfContact))

	mto, err := h.MoveTaskOrderUpdater.UpdatePostCounselingInfo(mtoID, params.Body, eTag)
	if err != nil {
		logger.Error("primeapi.UpdateMTOPostCounselingInformation error", zap.Error(err))
		switch err.(type) {
		case services.NotFoundError:
			return movetaskorderops.NewUpdateMTOPostCounselingInformationNotFound()
		case services.PreconditionFailedError:
			return movetaskorderops.NewUpdateMTOPostCounselingInformationPreconditionFailed().WithPayload(&primemessages.Error{Message: handlers.FmtString(err.Error())})
		case services.InvalidInputError:
			return movetaskorderops.NewUpdateMTOPostCounselingInformationUnprocessableEntity()
		default:
			return movetaskorderops.NewUpdateMTOPostCounselingInformationInternalServerError()
		}
	}
	mtoPayload := payloads.MoveTaskOrder(mto)
	return movetaskorderops.NewUpdateMTOPostCounselingInformationOK().WithPayload(mtoPayload)
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
		logger.Error("primeapi.UpdateMTOShipmentHandler error", zap.Error(err))
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
func createUser(h CreateMoveTaskOrderHandler, userEmail string, logger handlers.Logger) (*models.User, error) {
	if userEmail == "" {
		userEmail = "generatedMTOuser@example.com"
	}
	id := uuid.Must(uuid.NewV4())
	user := models.User{
		LoginGovUUID:  id,
		LoginGovEmail: userEmail,
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
func createOrGetCustomer(h CreateMoveTaskOrderHandler, customerIDString string, customerBody *primemessages.Customer, logger handlers.Logger) (*models.Customer, error) {
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

	// Create the new customer in the db
	verrs, err := h.DB().ValidateAndCreate(customer)
	if err != nil || verrs.Count() > 0 {
		logger.Error("createOrGetCustomer", zap.String("Error", err.Error()))
		returnErr := services.NewCreateObjectError("Customer", err, nil, "Error creating a customer")
		return nil, returnErr
	}
	return customer, nil

}
