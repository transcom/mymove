package moveorder

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/gobuffalo/pop"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

// ErrNotFound is returned when a given move task order is not found
type ErrNotFound struct {
	id uuid.UUID
}

func (e ErrNotFound) Error() string {
	return fmt.Sprintf("move task order id: %s not found", e.id.String())
}

type errInvalidInput struct {
	id uuid.UUID
	error
	validationErrors map[string][]string
}

// ErrInvalidInput is returned when an update to a move task order fails a validation rule
type ErrInvalidInput struct {
	errInvalidInput
}

func NewErrInvalidInput(id uuid.UUID, err error, validationErrors map[string][]string) ErrInvalidInput {
	return ErrInvalidInput{
		errInvalidInput{
			id:               id,
			error:            err,
			validationErrors: validationErrors,
		},
	}
}

func (e ErrInvalidInput) Error() string {
	return fmt.Sprintf("invalid input for move task order id: %s. %s", e.id.String(), e.InvalidFields())
}

func (e ErrInvalidInput) InvalidFields() map[string]string {
	es := make(map[string]string)
	if e.validationErrors == nil {
		return es
	}
	for k, v := range e.validationErrors {
		es[k] = strings.Join(v, " ")
	}
	return es
}

type moveOrderFetcher struct {
	db *pop.Connection
}

func (f moveOrderFetcher) ListMoveOrders() ([]models.MoveOrder, error) {
	var moveOrders []models.MoveOrder
	err := f.db.All(&moveOrders)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return []models.MoveOrder{}, ErrNotFound{}
		default:
			return []models.MoveOrder{}, err
		}
	}

	// Attempting to load these associations using Eager() returns an error, so this loop
	// loads them one at a time. This is creating a N + 1 query for each association, which is
	// bad. But that's also what the current implementation of Eager does, so this is no worse
	// that what we had.
	for i := range moveOrders {
		f.db.Load(&moveOrders[i], "Customer")
		f.db.Load(&moveOrders[i], "ConfirmationNumber")
		f.db.Load(&moveOrders[i], "DestinationDutyStation.Address")
		f.db.Load(&moveOrders[i], "OriginDutyStation.Address")
		f.db.Load(&moveOrders[i], "Entitlement")
	}

	return moveOrders, nil
}

// NewMoveOrderFetcher creates a new struct with the service dependencies
func NewMoveOrderFetcher(db *pop.Connection) services.MoveOrderFetcher {
	return &moveOrderFetcher{db}
}

// FetchMoveOrder retrieves a MoveOrder for a given UUID
func (f moveOrderFetcher) FetchMoveOrder(moveOrderID uuid.UUID) (*models.MoveOrder, error) {
	moveOrder := &models.MoveOrder{}
	err := f.db.Find(moveOrder, moveOrderID)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return &models.MoveOrder{}, ErrNotFound{moveOrderID}
		default:
			return &models.MoveOrder{}, err
		}
	}

	f.db.Load(moveOrder, "Customer")
	f.db.Load(moveOrder, "ConfirmationNumber")
	f.db.Load(moveOrder, "DestinationDutyStation.Address")
	f.db.Load(moveOrder, "OriginDutyStation.Address")
	f.db.Load(moveOrder, "Entitlement")

	return moveOrder, nil
}
