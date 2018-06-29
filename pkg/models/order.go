package models

import (
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/uuid"
	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
	"github.com/pkg/errors"

	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
)

// UploadedOrdersDocumentName is the name of an uploaded orders document
const UploadedOrdersDocumentName = "uploaded_orders"

// OrderStatus represents the status of an order record's lifecycle
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
	NewDutyStationID    uuid.UUID                          `json:"new_duty_station_id" db:"new_duty_station_id"`
	NewDutyStation      DutyStation                        `belongs_to:"duty_stations"`
	UploadedOrders      Document                           `belongs_to:"documents"`
	UploadedOrdersID    uuid.UUID                          `json:"uploaded_orders_id" db:"uploaded_orders_id"`
	OrdersNumber        *string                            `json:"orders_number" db:"orders_number"`
	Moves               Moves                              `has_many:"moves" fk_id:"orders_id" order_by:"created_at desc"`
	Status              OrderStatus                        `json:"status" db:"status"`
	TAC                 *string                            `json:"tac" db:"tac"`
	DepartmentIndicator *string                            `json:"department_indicator" db:"department_indicator"`
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
		&StringIsNilOrNotBlank{Field: o.TAC, Name: "Transportation Accounting Code"},
		&StringIsNilOrNotBlank{Field: o.DepartmentIndicator, Name: "Department Indicator"},
		&CannotBeTrueIfFalse{Field1: o.SpouseHasProGear, Name1: "SpouseHasProGear", Field2: o.HasDependents, Name2: "HasDependents"},
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

// SaveOrder saves an order
func SaveOrder(db *pop.Connection, order *Order) (*validate.Errors, error) {
	return db.ValidateAndSave(order)
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
	if o.Status != OrderStatusSUBMITTED {
		return errors.Wrap(ErrInvalidTransition, "Cancel")
	}

	o.Status = OrderStatusCANCELED
	return nil
}

// FetchOrder returns orders only if it is allowed for the given user to access those orders.
func FetchOrder(db *pop.Connection, session *auth.Session, id uuid.UUID) (Order, error) {
	var order Order
	err := db.Q().Eager("ServiceMember.User",
		"NewDutyStation.Address",
		"UploadedOrders.Uploads",
		"Moves.PersonallyProcuredMoves",
		"Moves.SignedCertifications").Find(&order, id)
	if err != nil {
		if errors.Cause(err).Error() == recordNotFoundErrorString {
			return Order{}, ErrFetchNotFound
		}
		// Otherwise, it's an unexpected err so we return that.
		return Order{}, err
	}
	// TODO: Handle case where more than one user is authorized to modify orders
	if session.IsMyApp() && order.ServiceMember.ID != session.ServiceMemberID {
		return Order{}, ErrFetchForbidden
	}

	return order, nil
}

// FetchOrderForPDFConversion returns orders and any attached uploads
func FetchOrderForPDFConversion(db *pop.Connection, id uuid.UUID) (Order, error) {
	var order Order
	err := db.Q().Eager("UploadedOrders.Uploads").Find(&order, id)
	if err != nil {
		if errors.Cause(err).Error() == recordNotFoundErrorString {
			return Order{}, ErrFetchNotFound
		}
		// Otherwise, it's an unexpected err so we return that.
		return Order{}, err
	}
	return order, nil
}

// CreateNewMove creates a move associated with these Orders
func (o *Order) CreateNewMove(db *pop.Connection, moveType *internalmessages.SelectedMoveType) (*Move, *validate.Errors, error) {
	return createNewMove(db, o.ID, moveType)
}
