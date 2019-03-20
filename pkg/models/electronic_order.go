package models

import (
	"context"
	"encoding/json"
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
	"github.com/gofrs/uuid"
	beeline "github.com/honeycombio/beeline-go"
	"github.com/pkg/errors"
)

// Issuer is the organization that issues orders.
type Issuer string

const (
	// IssuerArmy captures enum value "army"
	IssuerArmy Issuer = "army"
	// IssuerNavy captures enum value "navy"
	IssuerNavy Issuer = "navy"
	// IssuerAirForce captures enum value "air-force"
	IssuerAirForce Issuer = "air-force"
	// IssuerMarineCorps captures enum value "marine-corps"
	IssuerMarineCorps Issuer = "marine-corps"
	// IssuerCoastGuard captures enum value "coast-guard"
	IssuerCoastGuard Issuer = "coast-guard"
)

// ElectronicOrder contains the unchanging data of a set of orders across all amendments / revisions
type ElectronicOrder struct {
	ID           uuid.UUID                 `json:"id" db:"id"`
	CreatedAt    time.Time                 `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time                 `json:"updated_at" db:"updated_at"`
	OrdersNumber string                    `json:"orders_number" db:"orders_number"`
	Edipi        string                    `json:"edipi" db:"edipi"`
	Issuer       Issuer                    `json:"issuer" db:"issuer"`
	Revisions    ElectronicOrdersRevisions `has_many:"electronic_orders_revisions" order_by:"seq_num asc"`
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
		&validators.RegexMatch{Field: e.Edipi, Name: "Edipi", Expr: "\\d{10}"},
		&validators.StringInclusion{Field: string(e.Issuer), Name: "Issuer", List: []string{
			string(IssuerAirForce),
			string(IssuerArmy),
			string(IssuerCoastGuard),
			string(IssuerMarineCorps),
			string(IssuerNavy),
		}},
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

// CreateElectronicOrder inserts an empty set of electronic Orders into the database
func CreateElectronicOrder(ctx context.Context, dbConnection *pop.Connection, order *ElectronicOrder) (*validate.Errors, error) {
	ctx, span := beeline.StartSpan(ctx, "CreateElectronicOrder")
	defer span.Send()

	responseVErrors := validate.NewErrors()
	verrs, responseError := dbConnection.ValidateAndCreate(order)
	if verrs.HasAny() {
		responseVErrors.Append(verrs)
	}

	return responseVErrors, responseError
}

// CreateElectronicOrderWithRevision inserts a new set of electronic Orders into the database with its first Revision
func CreateElectronicOrderWithRevision(ctx context.Context, dbConnection *pop.Connection, order *ElectronicOrder, firstRevision *ElectronicOrdersRevision) (*validate.Errors, error) {
	ctx, span := beeline.StartSpan(ctx, "CreateElectronicOrder")
	defer span.Send()

	responseVErrors := validate.NewErrors()
	var responseError error

	// If the passed in function returns an error, the transaction is rolled back
	dbConnection.Transaction(func(dbConnection *pop.Connection) error {
		transactionError := errors.New("Rollback The transaction")
		if verrs, err := CreateElectronicOrder(ctx, dbConnection, order); verrs.HasAny() || err != nil {
			responseVErrors.Append(verrs)
			responseError = err
			return transactionError
		}
		firstRevision.ElectronicOrderID = order.ID
		firstRevision.ElectronicOrder = *order
		if verrs, err := CreateElectronicOrdersRevision(ctx, dbConnection, firstRevision); verrs.HasAny() || err != nil {
			responseVErrors.Append(verrs)
			responseError = err
			return transactionError
		}

		return nil
	})

	return responseVErrors, responseError
}

// FetchElectronicOrderByID gets all revisions of a set of Orders by their shared UUID,
// sorted in ascending order by their sequence number
func FetchElectronicOrderByID(db *pop.Connection, id uuid.UUID) (*ElectronicOrder, error) {
	var order ElectronicOrder
	err := db.Q().Eager("Revisions").Find(&order, id)
	if err != nil {
		if errors.Cause(err).Error() == recordNotFoundErrorString {
			return &order, ErrFetchNotFound
		}
		// Otherwise, it's an unexpected err so we return that.
	}

	return &order, nil
}

// FetchElectronicOrderByIssuerAndOrdersNum gets all revisions of a set of Orders by the unique combination of the Orders number and the issuer.
func FetchElectronicOrderByIssuerAndOrdersNum(db *pop.Connection, issuer string, ordersNum string) (*ElectronicOrder, error) {
	var order ElectronicOrder
	err := db.Q().Eager("Revisions").Where("orders_number = $1 AND issuer = $2", ordersNum, issuer).First(&order)
	if err != nil {
		if errors.Cause(err).Error() == recordNotFoundErrorString {
			return &order, ErrFetchNotFound
		}
		// Otherwise, it's an unexpected err so we return that.
	}
	return &order, err
}

// FetchElectronicOrdersByEdipiAndIssuers gets all Orders issued to a member by EDIPI from the specified issuers
func FetchElectronicOrdersByEdipiAndIssuers(db *pop.Connection, edipi string, issuers []string) ([]*ElectronicOrder, error) {
	var orders []ElectronicOrder
	err := db.Q().Eager("Revisions").Where("edipi = ?", edipi).Where("issuer IN (?)", issuers).All(&orders)
	ordersPtrs := make([]*ElectronicOrder, len(orders))
	for i := range orders {
		ordersPtrs[i] = &orders[i]
	}
	if err != nil {
		if errors.Cause(err).Error() == recordNotFoundErrorString {
			return ordersPtrs, ErrFetchNotFound
		}
		// Otherwise, it's an unexpected err so we return that.
	}
	return ordersPtrs, err
}
