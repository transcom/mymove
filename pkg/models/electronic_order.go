package models

import (
	"encoding/json"
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/uuid"
	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
	"github.com/pkg/errors"
	"github.com/transcom/mymove/pkg/gen/ordersmessages"
)

// ElectronicOrder contains the unchanging data of a set of orders across all amendments / revisions
type ElectronicOrder struct {
	ID           uuid.UUID                 `json:"id" db:"id"`
	CreatedAt    time.Time                 `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time                 `json:"updated_at" db:"updated_at"`
	OrdersNumber string                    `json:"orders_number" db:"orders_number"`
	Edipi        string                    `json:"edipi" db:"edipi"`
	Issuer       ordersmessages.Issuer     `json:"issuer" db:"issuer"`
	Revisions    ElectronicOrdersRevisions `has_many:"electronic_orders_revisions" order_by:"seq_num desc"`
}

// String is not required by pop and may be deleted
func (e ElectronicOrder) String() string {
	je, _ := json.Marshal(e)
	return string(je)
}

// ElectronicOrders is not required by pop and may be deleted
type ElectronicOrders []ElectronicOrder

// String is not required by pop and may be deleted
func (e ElectronicOrders) String() string {
	je, _ := json.Marshal(e)
	return string(je)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (e *ElectronicOrder) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.StringIsPresent{Field: e.OrdersNumber, Name: "OrdersNumber"},
		&validators.StringIsPresent{Field: e.Edipi, Name: "Edipi"},
		&validators.StringIsPresent{Field: string(e.Issuer), Name: "Issuer"},
	), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (e *ElectronicOrder) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (e *ElectronicOrder) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// FetchElectronicOrderByID gets all revisions of a set of Orders by their shared UUID,
// sorted in descending order by their sequence number
func FetchElectronicOrderByID(db *pop.Connection, id uuid.UUID) (ElectronicOrder, error) {
	var order ElectronicOrder
	err := db.Q().Eager("Revisions").Find(&order, id)
	if err != nil {
		if errors.Cause(err).Error() == recordNotFoundErrorString {
			return ElectronicOrder{}, ErrFetchNotFound
		}
		// Otherwise, it's an unexpected err so we return that.
		return ElectronicOrder{}, err
	}

	return order, nil
}

// FetchElectronicOrderByUniqueFeatures gets all revisions of a set Orders by the combination of features
// that make Orders unique - the Orders number, the EDIPI of the member, and the issuer. The results are
// returned in descending order by their sequence number
func FetchElectronicOrderByUniqueFeatures(db *pop.Connection, ordersNum string, edipi string, issuer string) ([]ElectronicOrder, error) {
	var orders []ElectronicOrder
	err := db.Q().Eager("ServiceMember").Where("orders_number = $1 AND edipi = $2 AND issuer = $3", ordersNum, edipi, issuer).All(orders)
	if err != nil {
		if errors.Cause(err).Error() == recordNotFoundErrorString {
			return []ElectronicOrder{}, ErrFetchNotFound
		}
		// Otherwise, it's an unexpected err so we return that.
		return []ElectronicOrder{}, err
	}
	return orders, err
}
