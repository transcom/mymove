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

type moveTaskOrderFetcher struct {
	db *pop.Connection
}

func (f moveTaskOrderFetcher) ListMoveTaskOrders(moveOrderID uuid.UUID) ([]models.MoveTaskOrder, error) {
	var moveTaskOrders []models.MoveTaskOrder
	err := f.db.Where("move_order_id = $1", moveOrderID).Eager().All(&moveTaskOrders)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return []models.MoveTaskOrder{}, ErrNotFound{}
		default:
			return []models.MoveTaskOrder{}, err
		}
	}
	return moveTaskOrders, nil
}

// NewMoveTaskOrderFetcher creates a new struct with the service dependencies
func NewMoveTaskOrderFetcher(db *pop.Connection) services.MoveTaskOrderFetcher {
	return &moveTaskOrderFetcher{db}
}

//FetchMoveTaskOrder retrieves a MoveTaskOrder for a given UUID
func (f moveTaskOrderFetcher) FetchMoveTaskOrder(moveTaskOrderID uuid.UUID) (*models.MoveTaskOrder, error) {
	mto := &models.MoveTaskOrder{}
	if err := f.db.Eager().Find(mto, moveTaskOrderID); err != nil {
		switch err {
		case sql.ErrNoRows:
			return &models.MoveTaskOrder{}, ErrNotFound{moveTaskOrderID}
		default:
			return &models.MoveTaskOrder{}, err
		}
	}

	f.createDefaultServiceItems(mto)

	return mto, nil
}

func (f moveTaskOrderFetcher) createDefaultServiceItems(mto *models.MoveTaskOrder) error {
	var reServices []models.ReService
	err := f.db.Where("code in (?)", []string{"MS", "CS"}).All(&reServices)

	if err != nil {
		return err
	}

	defaultServiceItems := make(map[uuid.UUID]models.MTOServiceItem)
	for _, reService := range reServices {
		defaultServiceItems[reService.ID] = models.MTOServiceItem{
			ReServiceID:     reService.ID,
			MoveTaskOrderID: mto.ID,
		}
	}

	// Remove the ones that exist on the mto
	for _, item := range mto.MTOServiceItems {
		for _, reService := range reServices {
			if item.ReServiceID == reService.ID {
				delete(defaultServiceItems, reService.ID)
			}
		}
	}

	for _, serviceItem := range defaultServiceItems {
		_, err := f.db.ValidateAndCreate(&serviceItem)

		if err != nil {
			return err
		}
	}

	return nil
}

type moveTaskOrderUpdater struct {
	db *pop.Connection
	moveTaskOrderFetcher
	//services.PersonallyProcuredMoveFetcher
}

// NewMoveTaskOrderFetcher creates a new struct with the service dependencies
func NewMoveTaskOrderUpdater(db *pop.Connection) services.MoveTaskOrderUpdater {
	return &moveTaskOrderUpdater{db, moveTaskOrderFetcher{db},}
}

//MakeAvailableToPrime updates the status of a MoveTaskOrder for a given UUID to make it available to prime
func (f moveTaskOrderUpdater) MakeAvailableToPrime(moveTaskOrderID uuid.UUID) (*models.MoveTaskOrder, error) {
	mto, err := f.FetchMoveTaskOrder(moveTaskOrderID)
	if err != nil {
		return &models.MoveTaskOrder{}, err
	}
	mto.IsAvailableToPrime = true
	vErrors, err := f.db.ValidateAndUpdate(mto)
	if vErrors.HasAny() {
		return &models.MoveTaskOrder{}, ErrInvalidInput{}
	}
	if err != nil {
		return &models.MoveTaskOrder{}, err
	}
	return mto, nil
}

func (f moveTaskOrderUpdater) UpdateMTOPersonallyProcuredMove(moveTaskOrderID uuid.UUID, personallyProcuredMoveID uuid.UUID) (*models.MoveTaskOrder, error) {
	//ppm, err := f.PersonallyProcuredMoveFetcher.FetchPersonallyProcuredMove(personallyProcuredMoveID)
	mto, err := f.FetchMoveTaskOrder(moveTaskOrderID)
	if err != nil {
		return &models.MoveTaskOrder{}, err
	}

	return mto, nil
}