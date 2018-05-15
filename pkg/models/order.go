package models

import (
	"encoding/json"
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/uuid"
	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
	"github.com/pkg/errors"

	"github.com/transcom/mymove/pkg/app"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
)

// UploadedOrdersDocumentName is the name of an uploaded orders document
const UploadedOrdersDocumentName = "uploaded_orders"

// Order is a set of orders received by a service member
type Order struct {
	ID                   uuid.UUID                          `json:"id" db:"id"`
	CreatedAt            time.Time                          `json:"created_at" db:"created_at"`
	UpdatedAt            time.Time                          `json:"updated_at" db:"updated_at"`
	ServiceMemberID      uuid.UUID                          `json:"service_member_id" db:"service_member_id"`
	ServiceMember        ServiceMember                      `belongs_to:"service_members"`
	IssueDate            time.Time                          `json:"issue_date" db:"issue_date"`
	ReportByDate         time.Time                          `json:"report_by_date" db:"report_by_date"`
	OrdersType           internalmessages.OrdersType        `json:"orders_type" db:"orders_type"`
	OrdersTypeDetail     *internalmessages.OrdersTypeDetail `json:"orders_type_detail" db:"orders_type_detail"`
	HasDependents        bool                               `json:"has_dependents" db:"has_dependents"`
	NewDutyStationID     uuid.UUID                          `json:"new_duty_station_id" db:"new_duty_station_id"`
	NewDutyStation       DutyStation                        `belongs_to:"duty_stations"`
	CurrentDutyStationID *uuid.UUID                         `json:"current_duty_station_id" db:"current_duty_station_id"`
	CurrentDutyStation   *DutyStation                       `belongs_to:"duty_stations"`
	UploadedOrders       Document                           `belongs_to:"documents"`
	UploadedOrdersID     uuid.UUID                          `json:"uploaded_orders_id" db:"uploaded_orders_id"`
	OrdersNumber         *string                            `json:"orders_number" db:"orders_number"`
}

// String is not required by pop and may be deleted
func (o Order) String() string {
	jo, _ := json.Marshal(o)
	return string(jo)
}

// Orders is not required by pop and may be deleted
type Orders []Order

// String is not required by pop and may be deleted
func (o Orders) String() string {
	jo, _ := json.Marshal(o)
	return string(jo)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (o *Order) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&OrdersTypeIsPresent{Field: o.OrdersType, Name: "OrdersType"},
		&validators.TimeIsPresent{Field: o.IssueDate, Name: "IssueDate"},
		&validators.TimeIsPresent{Field: o.ReportByDate, Name: "ReportByDate"},
		&validators.UUIDIsPresent{Field: o.ServiceMemberID, Name: "ServiceMemberID"},
		&validators.UUIDIsPresent{Field: o.NewDutyStationID, Name: "NewDutyStationID"},
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

// FetchOrder returns orders only if it is allowed for the given user to access those orders.
func FetchOrder(db *pop.Connection, user User, reqApp string, id uuid.UUID) (Order, error) {
	var order Order
	err := db.Q().Eager("ServiceMember.User", "NewDutyStation.Address", "UploadedOrders.Uploads").Find(&order, id)
	if err != nil {
		if errors.Cause(err).Error() == recordNotFoundErrorString {
			return Order{}, ErrFetchNotFound
		}
		// Otherwise, it's an unexpected err so we return that.
		return Order{}, err
	}
	// TODO: Handle case where more than one user is authorized to modify orders
	if reqApp == app.MyApp && order.ServiceMember.UserID != user.ID {
		return Order{}, ErrFetchForbidden
	}

	return order, nil
}

// CreateNewMove creates a move associated with these Orders
func (o *Order) CreateNewMove(db *pop.Connection, moveType *internalmessages.SelectedMoveType) (*Move, *validate.Errors, error) {
	return createNewMove(db, o.ID, moveType)
}
