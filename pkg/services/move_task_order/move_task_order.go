package movetaskorder

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

type fetchMoveTaskOrder struct {
	db *pop.Connection
}

// NewMoveTaskOrderFetcher creates a new struct with the service dependencies
func NewMoveTaskOrderFetcher(db *pop.Connection) services.MoveTaskOrderFetcher {
	return &fetchMoveTaskOrder{db}
}

//FetchMoveTaskOrder retrieves a MoveTaskOrder for a given UUID
func (f fetchMoveTaskOrder) FetchMoveTaskOrder(moveTaskOrderID uuid.UUID) (*models.MoveTaskOrder, error) {
	mto := &models.MoveTaskOrder{}
	if err := f.db.Eager().Find(mto, moveTaskOrderID); err != nil {
		switch err {
		case sql.ErrNoRows:
			return &models.MoveTaskOrder{}, ErrNotFound{moveTaskOrderID}
		default:
			return &models.MoveTaskOrder{}, err
		}
	}

	var defaultServiceItems models.MTOServiceItems
	for _, item := range mto.MTOServiceItems {
		defaultServiceItemIDs := []uuid.UUID{
			uuid.FromStringOrNil("1130e612-94eb-49a7-973d-72f33685e551"),
			uuid.FromStringOrNil("9dc919da-9b66-407b-9f17-05c0f03fcb50"),
		}

		for _, id := range defaultServiceItemIDs {
			if item.ReServiceID == id {
				defaultServiceItems = append(defaultServiceItems, item)
			}
		}
	}

	if len(defaultServiceItems) == 0 {
		serviceItems := []models.MTOServiceItem{
			{
				ReServiceID:     uuid.FromStringOrNil("1130e612-94eb-49a7-973d-72f33685e551"), // Shipment Management Services
				MoveTaskOrderID: mto.ID,
			},
			{
				ReServiceID: uuid.FromStringOrNil("9dc919da-9b66-407b-9f17-05c0f03fcb50"), // Counseling Services

				MoveTaskOrderID: mto.ID,
			},
		}

		_, err := f.db.ValidateAndCreate(&serviceItems[0])

		if err != nil {
			return nil, err
		}

		_, err = f.db.ValidateAndCreate(&serviceItems[1])

		if err != nil {
			return nil, err
		}

		return mto, nil
	}

	return mto, nil
}
