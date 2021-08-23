package supportmovetaskorder

import (
	"fmt"
	"time"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/services/support"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/gen/supportmessages"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/office_user/customer"
	"github.com/transcom/mymove/pkg/unit"
)

type moveTaskOrderCreator struct {
}

// InternalCreateMoveTaskOrder creates a move task order for the supportapi (internal use only, not used in production)
func (f moveTaskOrderCreator) InternalCreateMoveTaskOrder(appCtx appcontext.AppContext, payload supportmessages.MoveTaskOrder, logger *zap.Logger) (*models.Move, error) {
	var moveTaskOrder *models.Move
	var refID string
	if payload.Order == nil {
		return nil, services.NewQueryError("MoveTaskOrder", nil, "Order is necessary")
	}

	transactionError := appCtx.NewTransaction(func(txnAppCfg appcontext.AppContext) error {
		// Create or get customer
		customer, err := createOrGetCustomer(txnAppCfg, customer.NewCustomerFetcher(), payload.Order.CustomerID, payload.Order.Customer)
		if err != nil {
			return err
		}

		// Create order and entitlement
		order, err := createOrder(txnAppCfg, customer, payload.Order)
		if err != nil {
			return err
		}

		// Convert payload to model for moveTaskOrder
		moveTaskOrder = MoveTaskOrderModel(&payload)
		// referenceID cannot be set by user so generate it
		refID, err = models.GenerateReferenceID(txnAppCfg.DB())
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

		verrs, err := txnAppCfg.DB().ValidateAndCreate(moveTaskOrder)

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
func NewInternalMoveTaskOrderCreator() support.InternalMoveTaskOrderCreator {
	return &moveTaskOrderCreator{}
}

// createOrder creates a basic order - this is a support function do not use in production
func createOrder(appCtx appcontext.AppContext, customer *models.ServiceMember, orderPayload *supportmessages.Order) (*models.Order, error) {
	if orderPayload == nil {
		returnErr := services.NewInvalidInputError(uuid.Nil, nil, nil, "Order definition is required to create MoveTaskOrder")
		return nil, returnErr
	}

	// Order model will also contain the entitlement
	order := OrderModel(orderPayload)

	// Check that the order destination duty station exists, then hook up to order
	// It's required in the payload
	destinationDutyStation := models.DutyStation{}
	destinationDutyStationID := uuid.FromStringOrNil(orderPayload.DestinationDutyStationID.String())
	err := appCtx.DB().Find(&destinationDutyStation, destinationDutyStationID)
	if err != nil {
		appCtx.Logger().Error("supportapi.createOrder error", zap.Error(err))
		return nil, services.NewNotFoundError(destinationDutyStationID, ". The destinationDutyStation does not exist.")
	}
	order.NewDutyStation = destinationDutyStation
	order.NewDutyStationID = destinationDutyStationID
	// Check that if provided, the origin duty station exists, then hook up to order
	var originDutyStation *models.DutyStation
	if orderPayload.OriginDutyStationID != nil {
		originDutyStation = &models.DutyStation{}
		originDutyStationID := uuid.FromStringOrNil(orderPayload.OriginDutyStationID.String())
		err = appCtx.DB().Find(originDutyStation, originDutyStationID)
		if err != nil {
			appCtx.Logger().Error("supportapi.createOrder error", zap.Error(err))
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
		err = appCtx.DB().Find(uploadedOrders, uploadedOrdersID)
		if err != nil {
			appCtx.Logger().Error("supportapi.createOrder error", zap.Error(err))
			return nil, services.NewNotFoundError(uploadedOrdersID, ". The uploadedOrders does not exist.")
		}
		order.UploadedOrders = *uploadedOrders
		order.UploadedOrdersID = uploadedOrdersID
	}

	// Add customer to mO
	order.ServiceMember = *customer
	order.ServiceMemberID = customer.ID

	// Creates the order and the entitlement at the same time
	verrs, err := appCtx.DB().Eager().ValidateAndCreate(order)
	if verrs.Count() > 0 {
		appCtx.Logger().Error("supportapi.createOrder error", zap.Error(verrs))
		return nil, services.NewInvalidInputError(uuid.Nil, nil, verrs, "")
	} else if err != nil {
		appCtx.Logger().Error("supportapi.createOrder error", zap.Error(err))
		e := services.NewQueryError("Order", err, "Unable to create Order.")
		return nil, e
	}
	return order, nil
}

// createUser creates a user but this is a fake login.gov user
// this is support code only, do not use in a production case
func createUser(appCtx appcontext.AppContext, userEmail *string) (*models.User, error) {
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
	verrs, err := appCtx.DB().ValidateAndCreate(&user)
	if verrs.Count() > 0 {
		appCtx.Logger().Error("supportapi.createUser error", zap.Error(verrs))
		return nil, services.NewInvalidInputError(uuid.Nil, nil, verrs, "")
	} else if err != nil {
		appCtx.Logger().Error("supportapi.createUser error", zap.Error(err))
		e := services.NewQueryError("User", err, "Unable to create User.")
		return nil, e
	}
	return &user, nil
}

// createOrGetCustomer creates a customer or gets one if id was provided
func createOrGetCustomer(appCtx appcontext.AppContext, f services.CustomerFetcher, payloadCustomerID *strfmt.UUID, customerBody *supportmessages.Customer) (*models.ServiceMember, error) {
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
		customer, err := f.FetchCustomer(appCtx, customerID)
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
	user, err := createUser(appCtx, customerBody.Email)
	if err != nil {
		return nil, err
	}

	customer := CustomerModel(customerBody)
	customer.User = *user
	customer.UserID = user.ID

	// Create the new customer in the db
	verrs, err = appCtx.DB().ValidateAndCreate(customer)
	if verrs.Count() > 0 {
		appCtx.Logger().Error("supportapi.createOrGetCustomer error", zap.Error(verrs))
		return nil, services.NewInvalidInputError(uuid.Nil, nil, verrs, "")
	} else if err != nil {
		appCtx.Logger().Error("supportapi.createOrGetCustomer error", zap.Error(err))
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
		Rank:          (*models.ServiceMemberRank)(customer.Rank),
		FirstName:     customer.FirstName,
		LastName:      customer.LastName,
		PersonalEmail: customer.Email,
		Telephone:     customer.Phone,
	}
}

// OrderModel converts payload to model - it does not convert nested
// duty stations but will preserve the ID if provided.
// It will create nested customer and entitlement models
// if those are provided in the payload
func OrderModel(orderPayload *supportmessages.Order) *models.Order {
	if orderPayload == nil {
		return nil
	}

	model := &models.Order{
		ID:           uuid.FromStringOrNil(orderPayload.ID.String()),
		OrdersNumber: orderPayload.OrderNumber,
		Entitlement:  EntitlementModel(orderPayload.Entitlement),
		IssueDate:    (time.Time)(*orderPayload.IssueDate),
		TAC:          orderPayload.Tac,
	}

	if orderPayload.Rank != nil {
		model.Grade = swag.String((string)(*orderPayload.Rank))
	}

	if orderPayload.Status != nil {
		model.Status = (models.OrderStatus)(*orderPayload.Status)
	}

	if orderPayload.OrdersType != nil {
		model.OrdersType = (internalmessages.OrdersType)(*orderPayload.OrdersType)
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

	if orderPayload.OrdersTypeDetail != nil {
		model.OrdersTypeDetail = (*internalmessages.OrdersTypeDetail)(orderPayload.OrdersTypeDetail)
	}

	if orderPayload.DepartmentIndicator != nil {
		model.DepartmentIndicator = (*string)(orderPayload.DepartmentIndicator)
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

	if mtoPayload.SelectedMoveType != nil {
		model.SelectedMoveType = (*models.SelectedMoveType)(mtoPayload.SelectedMoveType)
	}

	return model
}
