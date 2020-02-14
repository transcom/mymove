package movetaskorder

import (
	"database/sql"
	"fmt"
	"github.com/go-openapi/strfmt"
	movetaskorderops "github.com/transcom/mymove/pkg/gen/primeapi/primeoperations/move_task_order"
	"github.com/transcom/mymove/pkg/gen/primemessages"
	"strings"
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	ppmService "github.com/transcom/mymove/pkg/services/personally_procured_move"
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

//ErrPreconditionFailed is returned when a given mto shipment if attempting to update after the if-unmodified-since date
type ErrPreconditionFailed struct {
	id              uuid.UUID
	unmodifiedSince time.Time
	message         string
}

func (e ErrPreconditionFailed) Error() string {
	return fmt.Sprintf("%s %s can not be updated after date %s", e.message, e.id.String(), strfmt.Date(e.unmodifiedSince))
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
}

// NewMoveTaskOrderUpdater creates a new struct with the service dependencies
func NewMoveTaskOrderUpdater(db *pop.Connection) services.MoveTaskOrderUpdater {
	return &moveTaskOrderUpdater{db, moveTaskOrderFetcher{db}}
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

// UpdateMTOWithPersonallyProcuredMove updates the PPM's estimated weight and type and associates it to the MTO if it is not already
func (f moveTaskOrderUpdater) UpdateMTOWithPersonallyProcuredMove(params movetaskorderops.UpdateMTOPostCounselingInformationParams) (*models.MoveTaskOrder, error) {
	mtoID := uuid.FromStringOrNil(params.MoveTaskOrderID)
	unmodifiedSince := time.Time(params.IfUnmodifiedSince)
	newPPMInfo := params.Body.PersonallyProcuredMove

	_, ppmErr := ppmService.NewPersonallyProcuredMoveFetcher(f.db).FetchPersonallyProcuredMove(uuid.FromStringOrNil(newPPMInfo.ID.String()))

	if ppmErr != nil {
		return &models.MoveTaskOrder{}, ppmErr
	}

	mto, err := f.FetchMoveTaskOrder(mtoID)
	if err != nil {
		return &models.MoveTaskOrder{}, err
	}
	err = updateMTOWithPersonallyProcuredMove(f.db, mtoID, unmodifiedSince, newPPMInfo)
	return mto, nil
}

// updateMTOWithPersonallyProcuredMove updates the PPM with the info in the body and then associates it with the MTO
func updateMTOWithPersonallyProcuredMove(db *pop.Connection, mtoID uuid.UUID, unmodifiedSince time.Time, newPPMInfo *primemessages.PersonallyProcuredMove) error {
	ppmQuery := `UPDATE personally_procured_moves,
		SET type = ?,
			weight_estimate = ?,
			updated_at = NOW()
		WHERE
			id = ?
		AND
			updated_at = ?
		;`

	affectedRows, err := db.RawQuery(ppmQuery, newPPMInfo.Type, newPPMInfo.EstimatedWeight, newPPMInfo.ID, unmodifiedSince).ExecWithCount()

	if err != nil {
		return err
	}

	if affectedRows != 1 {
		return ErrPreconditionFailed{message: "PPM", id: uuid.FromStringOrNil(newPPMInfo.ID.String()), unmodifiedSince: unmodifiedSince}
	}

	mtoQuery := `UPDATE move_task_orders,
		SET personally_procured_move_id = ?,
			updated_at = NOW()
		WHERE
			id = ?
		AND
			updated_at = ?
		;
	`

	affectedRows2, err2 := db.RawQuery(mtoQuery, newPPMInfo.ID, mtoID, unmodifiedSince).ExecWithCount()

	if err2 != nil {
		return err
	}

	if affectedRows2 != 1 {
		return ErrPreconditionFailed{message: "MTO", id: uuid.FromStringOrNil(mtoID.String()), unmodifiedSince: unmodifiedSince}
	}

	return nil
}
