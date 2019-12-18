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

//ErrNotFound is returned when a given move task order is not found
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

//ErrInvalidInput is returned when an update to a move task order fails a validation rule
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

type fetchMoveOrder struct {
	db *pop.Connection
}

//Last, First Name | Confirmation # | Branch of Service | Origin Duty Station
func (f fetchMoveOrder) ListMoveOrders() ([]models.MoveOrder, error) {
	var moveOrders []models.MoveOrder
	err := f.db.Eager("Customer", "ConfirmationNumber","DestinationDutyStation.Address", "OriginDutyStation.Address", "Entitlement").All(&moveOrders)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return []models.MoveOrder{}, ErrNotFound{}
		default:
			return []models.MoveOrder{}, err
		}
	}
	return moveOrders, nil
}

// NewMoveOrderFetcher creates a new struct with the service dependencies
func NewMoveOrderFetcher(db *pop.Connection) services.MoveOrderFetcher {
	return &fetchMoveOrder{db}
}

//FetchMoveOrder retrieves a MoveOrder for a given UUID
func (f fetchMoveOrder) FetchMoveOrder(moveOrderID uuid.UUID) (*models.MoveOrder, error) {
	moveOrder := &models.MoveOrder{}
	err := f.db.Eager("DestinationDutyStation.Address", "OriginDutyStation.Address", "Entitlement").Find(moveOrder, moveOrderID)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return &models.MoveOrder{}, ErrNotFound{moveOrderID}
		default:
			return &models.MoveOrder{}, err
		}
	}
	return moveOrder, nil
}
