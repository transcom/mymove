package supportmovetaskorder

import (
	"fmt"
	"time"

	"github.com/transcom/mymove/pkg/services/support"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/gobuffalo/pop/v5"
	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/gen/supportmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/office_user/customer"
	"github.com/transcom/mymove/pkg/unit"
)

type moveTaskOrderCreator struct {
	db *pop.Connection
}

// InternalCreateMoveTaskOrder creates a move task order for the supportapi (internal use only, not used in production)
func (f moveTaskOrderCreator) InternalCreateMoveTaskOrder(payload supportmessages.MoveTaskOrder, logger handlers.Logger) (*models.Move, error) {
	var moveTaskOrder *models.Move
	var refID string
	if payload.MoveOrder == nil {
		return nil, services.NewQueryError("MoveTaskOrder", nil, "MoveOrder is necessary")
	}

	transactionError := f.db.Transaction(func(tx *pop.Connection) error {
		// Create or get customer
		customer, err := createOrGetCustomer(tx, customer.NewCustomerFetcher(tx), payload.MoveOrder.CustomerID, payload.MoveOrder.Customer, logger)
		if err != nil {
			return err
		}

		// Create move order and entitlement
		order, err := createMoveOrder(tx, customer, payload.MoveOrder, logger)
		if err != nil {
			return err
		}

		// Convert payload to model for moveTaskOrder
		moveTaskOrder = MoveTaskOrderModel(&payload)
		// referenceID cannot be set by user so generate it
		refID, err = models.GenerateReferenceID(tx)
		if err != nil {
			return err
		}
		moveTaskOrder.ReferenceID = &refID

		if moveTaskOrder.Locator == "" {
			moveTaskOrder.Locator = models.GenerateLocator()
		}
		if moveTaskOrder.Status == "" {
			moveTaskOrder.Status = models.MoveStatusDRAFT
		}
		moveTaskOrder.Show = swag.Bool(true)
		moveTaskOrder.Orders = *order
		moveTaskOrder.OrdersID = order.ID

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

// NewInternalMoveTaskOrderCreator creates a new struct with the service dependencies
func NewInternalMoveTaskOrderCreator(db *pop.Connection) support.InternalMoveTaskOrderCreator {
	return &moveTaskOrderCreator{db}
}

// createMoveOrder creates a basic move order - this is a support function do not use in production
func createMoveOrder(tx *pop.Connection, customer *models.ServiceMember, orderPayload *supportmessages.MoveOrder, logger handlers.Logger) (*models.Order, error) {
	if orderPayload == nil {
		returnErr := services.NewInvalidInputError(uuid.Nil, nil, nil, "MoveOrder definition is required to create MoveTaskOrder")
		return nil, returnErr
	}

	// Move order model will also contain the entitlement
	order := MoveOrderModel(orderPayload)

	// Check that the order destination duty station exists, then hook up to order
	// It's required in the payload
	destinationDutyStation := models.DutyStation{}
	destinationDutyStationID := uuid.FromStringOrNil(orderPayload.DestinationDutyStationID.String())
	err := tx.Find(&destinationDutyStation, destinationDutyStationID)
	if err != nil {
		logger.Error("supportapi.createMoveOrder error", zap.Error(err))
		return nil, services.NewNotFoundError(destinationDutyStationID, ". The destinationDutyStation does not exist.")
	}
	order.NewDutyStation = destinationDutyStation
	order.NewDutyStationID = destinationDutyStationID
	// Check that if provided, the origin duty station exists, then hook up to order
	var originDutyStation *models.DutyStation
	if orderPayload.OriginDutyStationID != nil {
		originDutyStation = &models.DutyStation{}
		originDutyStationID := uuid.FromStringOrNil(orderPayload.OriginDutyStationID.String())
		err = tx.Find(originDutyStation, originDutyStationID)
		if err != nil {
			logger.Error("supportapi.createMoveOrder error", zap.Error(err))
			return nil, services.NewNotFoundError(originDutyStationID, ". The originDutyStation does not exist.")
		}
		order.OriginDutyStation = originDutyStation
		order.OriginDutyStationID = &originDutyStationID
	}
	// Check that the uploaded orders document exists
	var uploadedOrders *models.Document
	if orderPayload.UploadedOrdersID != nil {
		uploadedOrders = &models.Document{}
		uploadedOrdersID := uuid.FromStringOrNil(orderPayload.UploadedOrdersID.String())
		fmt.Println("\n\nUploaded orders id is ", uploadedOrdersID)
		err = tx.Find(uploadedOrders, uploadedOrdersID)
		if err != nil {
			logger.Error("supportapi.createMoveOrder error", zap.Error(err))
			return nil, services.NewNotFoundError(uploadedOrdersID, ". The uploadedOrders does not exist.")
		}
		order.UploadedOrders = *uploadedOrders
		order.UploadedOrdersID = uploadedOrdersID
	}

	// Add customer to mO
	order.ServiceMember = *customer
	order.ServiceMemberID = customer.ID

	// Creates the order and the entitlement at the same time
	verrs, err := tx.Eager().ValidateAndCreate(order)
	if verrs.Count() > 0 {
		logger.Error("supportapi.createMoveOrder error", zap.Error(verrs))
		return nil, services.NewInvalidInputError(uuid.Nil, nil, verrs, "")
	} else if err != nil {
		logger.Error("supportapi.createMoveOrder error", zap.Error(err))
		e := services.NewQueryError("MoveOrder", err, "Unable to create MoveOrder.")
		return nil, e
	}
	return order, nil
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
		LoginGovUUID:  &id,
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
func createOrGetCustomer(tx *pop.Connection, f services.CustomerFetcher, payloadCustomerID *strfmt.UUID, customerBody *supportmessages.Customer, logger handlers.Logger) (*models.ServiceMember, error) {
	verrs := validate.NewErrors()

	// If customer ID string is provided, we should find this customer
	if payloadCustomerID != nil {
		customerID, err := uuid.FromString(payloadCustomerID.String())
		if err != nil {
			verrs.Add("customerID", "UUID is invalid")
			returnErr := services.NewInvalidInputError(uuid.Nil, err, verrs, "ID provided for Customer was invalid")
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
		verrs.Add("customer", "If no customerID is provided, nested Customer object is required")
		returnErr := services.NewInvalidInputError(uuid.Nil, nil, verrs, "If CustomerID is not provided, customer object is required to create Customer")
		return nil, returnErr
	}
	user, err := createUser(tx, customerBody.Email, logger)
	if err != nil {
		return nil, err
	}

	customer := CustomerModel(customerBody)
	customer.User = *user
	customer.UserID = user.ID

	// Create the new customer in the db
	verrs, err = tx.ValidateAndCreate(customer)
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

// CustomerModel converts payload to model - currently does not tackle addresses
func CustomerModel(customer *supportmessages.Customer) *models.ServiceMember {
	if customer == nil {
		return nil
	}
	return &models.ServiceMember{
		ID:            uuid.FromStringOrNil(customer.ID.String()),
		Affiliation:   (*models.ServiceMemberAffiliation)(customer.Agency),
		Edipi:         customer.DodID,
		Rank:          (*models.ServiceMemberRank)(&customer.Rank),
		FirstName:     customer.FirstName,
		LastName:      customer.LastName,
		PersonalEmail: customer.Email,
		Telephone:     customer.Phone,
	}
}

// MoveOrderModel converts payload to model - it does not convert nested
// duty stations but will preserve the ID if provided.
// It will create nested customer and entitlement models
// if those are provided in the payload
func MoveOrderModel(orderPayload *supportmessages.MoveOrder) *models.Order {
	if orderPayload == nil {
		return nil
	}
	model := &models.Order{
		ID:           uuid.FromStringOrNil(orderPayload.ID.String()),
		Grade:        swag.String((string)(orderPayload.Rank)),
		OrdersNumber: orderPayload.OrderNumber,
		Entitlement:  EntitlementModel(orderPayload.Entitlement),
		Status:       (models.OrderStatus)(orderPayload.Status),
		IssueDate:    (time.Time)(*orderPayload.IssueDate),
		OrdersType:   (internalmessages.OrdersType)(orderPayload.OrdersType),
		TAC:          orderPayload.Tac,
	}

	if orderPayload.CustomerID != nil {
		customerID := uuid.FromStringOrNil(orderPayload.CustomerID.String())
		model.ServiceMemberID = customerID
	}

	if orderPayload.DestinationDutyStationID != nil {
		model.NewDutyStationID = uuid.FromStringOrNil(orderPayload.DestinationDutyStationID.String())
	}

	if orderPayload.OriginDutyStationID != nil {
		originDutyStationID := uuid.FromStringOrNil(orderPayload.OriginDutyStationID.String())
		model.OriginDutyStationID = &originDutyStationID
	}

	if orderPayload.Customer != nil {
		model.ServiceMember = *CustomerModel(orderPayload.Customer)
	}

	if orderPayload.UploadedOrdersID != nil {
		uploadedOrdersID := uuid.FromStringOrNil(orderPayload.UploadedOrdersID.String())
		model.UploadedOrdersID = uploadedOrdersID
	}

	reportByDate := time.Time(*orderPayload.ReportByDate)
	if !reportByDate.IsZero() {
		model.ReportByDate = reportByDate
	}
	return model
}

// EntitlementModel converts the payload to model
func EntitlementModel(entitlementPayload *supportmessages.Entitlement) *models.Entitlement {
	if entitlementPayload == nil {
		return nil
	}

	// proGearWeight and ProGearWeightSpouse currently not handled as
	// they are not in the entitlement record in the db.
	model := &models.Entitlement{
		ID:                    uuid.FromStringOrNil(entitlementPayload.ID.String()),
		DependentsAuthorized:  entitlementPayload.DependentsAuthorized,
		NonTemporaryStorage:   entitlementPayload.NonTemporaryStorage,
		PrivatelyOwnedVehicle: entitlementPayload.PrivatelyOwnedVehicle,
	}

	if entitlementPayload.AuthorizedWeight != nil {
		model.DBAuthorizedWeight = swag.Int(int(*entitlementPayload.AuthorizedWeight))
	}

	totalDependents := int(entitlementPayload.TotalDependents)
	model.TotalDependents = &totalDependents

	storageInTransit := int(entitlementPayload.StorageInTransit)
	model.StorageInTransit = &storageInTransit

	return model
}

// MoveTaskOrderModel return an MTO model constructed from the payload.
// Does not create nested mtoServiceItems, mtoShipments, or paymentRequests
func MoveTaskOrderModel(mtoPayload *supportmessages.MoveTaskOrder) *models.Move {
	if mtoPayload == nil {
		return nil
	}
	ppmEstimatedWeight := unit.Pound(mtoPayload.PpmEstimatedWeight)
	contractorID := uuid.FromStringOrNil(mtoPayload.ContractorID.String())
	model := &models.Move{
		ReferenceID:        &mtoPayload.ReferenceID,
		Locator:            mtoPayload.MoveCode,
		PPMEstimatedWeight: &ppmEstimatedWeight,
		PPMType:            &mtoPayload.PpmType,
		ContractorID:       &contractorID,
		Status:             (models.MoveStatus)(mtoPayload.Status),
	}

	if mtoPayload.AvailableToPrimeAt != nil {
		availableToPrimeAt := time.Time(*mtoPayload.AvailableToPrimeAt)
		model.AvailableToPrimeAt = &availableToPrimeAt
	}

	return model
}
