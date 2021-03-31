package models

import (
	"fmt"
	"time"

	"github.com/gobuffalo/pop/v5"
	"github.com/gobuffalo/validate/v3"
	"github.com/gobuffalo/validate/v3/validators"
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"

	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
)

// UploadedOrdersDocumentName is the name of an uploaded orders document
const UploadedOrdersDocumentName = "uploaded_orders"

// OrderStatus represents the state of an order record in the UX manual orders flow
type OrderStatus string

const (
	// OrderStatusDRAFT captures enum value "DRAFT"
	OrderStatusDRAFT OrderStatus = "DRAFT"
	// OrderStatusSUBMITTED captures enum value "SUBMITTED"
	OrderStatusSUBMITTED OrderStatus = "SUBMITTED"
	// OrderStatusAPPROVED captures enum value "APPROVED"
	OrderStatusAPPROVED OrderStatus = "APPROVED"
	// OrderStatusCANCELED captures enum value "CANCELED"
	OrderStatusCANCELED OrderStatus = "CANCELED"
)

// Order is a set of orders received by a service member
type Order struct {
	ID                  uuid.UUID                          `json:"id" db:"id"`
	CreatedAt           time.Time                          `json:"created_at" db:"created_at"`
	UpdatedAt           time.Time                          `json:"updated_at" db:"updated_at"`
	ServiceMemberID     uuid.UUID                          `json:"service_member_id" db:"service_member_id"`
	ServiceMember       ServiceMember                      `belongs_to:"service_members"`
	IssueDate           time.Time                          `json:"issue_date" db:"issue_date"`
	ReportByDate        time.Time                          `json:"report_by_date" db:"report_by_date"`
	OrdersType          internalmessages.OrdersType        `json:"orders_type" db:"orders_type"`
	OrdersTypeDetail    *internalmessages.OrdersTypeDetail `json:"orders_type_detail" db:"orders_type_detail"`
	HasDependents       bool                               `json:"has_dependents" db:"has_dependents"`
	SpouseHasProGear    bool                               `json:"spouse_has_pro_gear" db:"spouse_has_pro_gear"`
	OriginDutyStation   *DutyStation                       `belongs_to:"duty_stations" fk_id:"origin_duty_station_id"`
	OriginDutyStationID *uuid.UUID                         `json:"origin_duty_station_id" db:"origin_duty_station_id"`
	NewDutyStationID    uuid.UUID                          `json:"new_duty_station_id" db:"new_duty_station_id"`
	NewDutyStation      DutyStation                        `belongs_to:"duty_stations" fk_id:"new_duty_station_id"`
	UploadedOrders      Document                           `belongs_to:"documents"`
	UploadedOrdersID    uuid.UUID                          `json:"uploaded_orders_id" db:"uploaded_orders_id"`
	OrdersNumber        *string                            `json:"orders_number" db:"orders_number"`
	Moves               Moves                              `has_many:"moves" fk_id:"orders_id" order_by:"created_at desc"`
	Status              OrderStatus                        `json:"status" db:"status"`
	TAC                 *string                            `json:"tac" db:"tac"`
	SAC                 *string                            `json:"sac" db:"sac"`
	DepartmentIndicator *string                            `json:"department_indicator" db:"department_indicator"`
	Grade               *string                            `json:"grade" db:"grade"`
	Entitlement         *Entitlement                       `belongs_to:"entitlements"`
	EntitlementID       *uuid.UUID                         `json:"entitlement_id" db:"entitlement_id"`
}

// Orders is not required by pop and may be deleted
type Orders []Order

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (o *Order) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&OrdersTypeIsPresent{Field: o.OrdersType, Name: "OrdersType"},
		&validators.TimeIsPresent{Field: o.IssueDate, Name: "IssueDate"},
		&validators.TimeIsPresent{Field: o.ReportByDate, Name: "ReportByDate"},
		&validators.UUIDIsPresent{Field: o.ServiceMemberID, Name: "ServiceMemberID"},
		&validators.UUIDIsPresent{Field: o.NewDutyStationID, Name: "NewDutyStationID"},
		&validators.StringIsPresent{Field: string(o.Status), Name: "Status"},
		&StringIsNilOrNotBlank{Field: o.TAC, Name: "TransportationAccountingCode"},
		&StringIsNilOrNotBlank{Field: o.SAC, Name: "SAC"},
		&StringIsNilOrNotBlank{Field: o.DepartmentIndicator, Name: "DepartmentIndicator"},
		&CannotBeTrueIfFalse{Field1: o.SpouseHasProGear, Name1: "SpouseHasProGear", Field2: o.HasDependents, Name2: "HasDependents"},
		&OptionalUUIDIsPresent{Field: o.EntitlementID, Name: "EntitlementID"},
		&OptionalUUIDIsPresent{Field: o.OriginDutyStationID, Name: "OriginDutyStationID"},
		&StringIsPresentAfterSubmission{Name: "TransportationAccountingCode", Field: o.TAC, Order: *o, DB: tx},
		&StringIsPresentAfterSubmission{Name: "DepartmentIndicator", Field: o.DepartmentIndicator, Order: *o, DB: tx},
		&StringIsPresentAfterSubmission{Name: "OrdersNumber", Field: o.OrdersNumber, Order: *o, DB: tx},
		&OrdersTypeDetailIsPresentAfterSubmission{Name: "OrdersTypeDetail", Field: o.OrdersTypeDetail, Order: *o, DB: tx},
		&OptionalRegexMatch{Name: "TransportationAccountingCode", Field: o.TAC, Expr: `\A([A-Za-z0-9]){4}\z`, Message: "TAC must be exactly 4 alphanumeric characters."},
	), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (o *Order) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (o *Order) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// StringIsPresentAfterSubmission checks presence of fields after an order has been submitted
type StringIsPresentAfterSubmission struct {
	Name  string
	Field *string
	Order Order
	DB    *pop.Connection
}

// IsValid adds an error if the field is blank
func (v *StringIsPresentAfterSubmission) IsValid(errors *validate.Errors) {
	order := v.Order

	if len(order.Moves) <= 0 || order.Moves[0].Status == MoveStatusDRAFT {
		return
	}

	if v.Field == nil || *v.Field == "" {
		errors.Add(validators.GenerateKey(v.Name), fmt.Sprintf("%s cannot be blank.", v.Name))
	}
}

// OrdersTypeDetailIsPresentAfterSubmission validates that orders type field is present
type OrdersTypeDetailIsPresentAfterSubmission struct {
	Name  string
	Field *internalmessages.OrdersTypeDetail
	Order Order
	DB    *pop.Connection
}

// IsValid adds an error if the string value is blank.
func (v *OrdersTypeDetailIsPresentAfterSubmission) IsValid(errors *validate.Errors) {
	order := v.Order

	if len(order.Moves) <= 0 || order.Moves[0].Status == MoveStatusDRAFT {
		return
	}

	if v.Field == nil || string(*v.Field) == "" {
		errors.Add(validators.GenerateKey(v.Name), fmt.Sprintf("%s cannot be blank.", v.Name))
	}
}

// SaveOrder saves an order
func SaveOrder(db *pop.Connection, order *Order) (*validate.Errors, error) {
	responseVErrors := validate.NewErrors()
	var responseError error

	transactionErr := db.Transaction(func(dbConnection *pop.Connection) error {
		transactionError := errors.New("Rollback The transaction")

		ppm, err := FetchPersonallyProcuredMoveByOrderID(db, order.ID)
		if err != nil {
			responseError = err
			return transactionError
		}

		ppm.DestinationPostalCode = &order.NewDutyStation.Address.PostalCode
		if verrs, err := dbConnection.ValidateAndSave(ppm); verrs.HasAny() || err != nil {
			responseVErrors.Append(verrs)
			responseError = err
			return transactionError
		}

		if verrs, err := db.ValidateAndSave(order); verrs.HasAny() || err != nil {
			responseVErrors.Append(verrs)
			responseError = err
			return transactionError
		}
		return nil
	})

	if transactionErr != nil {
		return responseVErrors, responseError
	}

	return responseVErrors, responseError
}

// State Machine
// Avoid calling Order.Status = ... ever. Use these methods to change the state.

// Submit submits the Order
func (o *Order) Submit() error {
	if o.Status != OrderStatusDRAFT {
		return errors.Wrap(ErrInvalidTransition, "Submit")
	}

	o.Status = OrderStatusSUBMITTED
	return nil
}

// Cancel cancels the Order
func (o *Order) Cancel() error {
	if o.Status == OrderStatusCANCELED {
		return errors.Wrap(ErrInvalidTransition, "Cancel")
	}

	o.Status = OrderStatusCANCELED
	return nil
}

// FetchOrderForUser returns orders only if it is allowed for the given user to access those orders.
func FetchOrderForUser(db *pop.Connection, session *auth.Session, id uuid.UUID) (Order, error) {
	var order Order
	err := db.Q().EagerPreload("ServiceMember.User",
		"OriginDutyStation.Address",
		"OriginDutyStation.TransportationOffice",
		"NewDutyStation.Address",
		"NewDutyStation.TransportationOffice",
		"UploadedOrders",
		"Moves.PersonallyProcuredMoves",
		"Moves.SignedCertifications",
		"Entitlement",
		"OriginDutyStation").
		Find(&order, id)
	if err != nil {
		if errors.Cause(err).Error() == RecordNotFoundErrorString {
			return Order{}, ErrFetchNotFound
		}
		// Otherwise, it's an unexpected err so we return that.
		return Order{}, err
	}

	// Eager loading of nested has_many associations is broken
	err = db.Load(&order.UploadedOrders, "UserUploads.Upload")
	if err != nil {
		return Order{}, err
	}

	// Only return user uploads that haven't been deleted
	userUploads := order.UploadedOrders.UserUploads
	relevantUploads := make([]UserUpload, 0, len(userUploads))
	for _, userUpload := range userUploads {
		if userUpload.DeletedAt == nil {
			relevantUploads = append(relevantUploads, userUpload)
		}
	}
	order.UploadedOrders.UserUploads = relevantUploads

	// TODO: Handle case where more than one user is authorized to modify orders
	if session.IsMilApp() && order.ServiceMember.ID != session.ServiceMemberID {
		return Order{}, ErrFetchForbidden
	}
	return order, nil
}

// FetchOrder returns orders without REGARDLESS OF USER.
// DO NOT USE IF YOU NEED USER AUTH
func FetchOrder(db *pop.Connection, id uuid.UUID) (Order, error) {
	var order Order
	err := db.Q().Find(&order, id)
	if err != nil {
		if errors.Cause(err).Error() == RecordNotFoundErrorString {
			return Order{}, ErrFetchNotFound
		}
		// Otherwise, it's an unexpected err so we return that.
		return Order{}, err
	}

	return order, nil
}

// FetchOrderForPDFConversion returns orders and any attached uploads
func FetchOrderForPDFConversion(db *pop.Connection, id uuid.UUID) (Order, error) {
	var order Order
	err := db.Q().Eager("UploadedOrders.UserUploads.Upload").Find(&order, id)
	if err != nil {
		if errors.Cause(err).Error() == RecordNotFoundErrorString {
			return Order{}, ErrFetchNotFound
		}
		// Otherwise, it's an unexpected err so we return that.
		return Order{}, err
	}
	return order, nil
}

// CreateNewMove creates a move associated with these Orders
func (o *Order) CreateNewMove(db *pop.Connection, moveOptions MoveOptions) (*Move, *validate.Errors, error) {
	return createNewMove(db, *o, moveOptions)
}

// IsComplete checks if orders have all fields necessary to approve a move
func (o *Order) IsComplete() bool {

	if o.OrdersTypeDetail == nil {
		return false
	}
	if o.OrdersNumber == nil {
		return false
	}
	if o.DepartmentIndicator == nil {
		return false
	}
	if o.TAC == nil {
		return false
	}
	if o.SAC == nil {
		return false
	}
	return true
}

// IsCompleteForGBL checks if orders have all fields necessary to generate a GBL
func (o *Order) IsCompleteForGBL() bool {

	if o.DepartmentIndicator == nil {
		return false
	}
	if o.TAC == nil {
		return false
	}
	if o.SAC == nil {
		return false
	}
	return true
}
