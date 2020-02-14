package personally_procured_move


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

type personallyProcuredMoveFetcher struct {
	db *pop.Connection
}

// NewMoveTaskOrderFetcher creates a new struct with the service dependencies
func NewMoveTaskOrderFetcher(db *pop.Connection) services.PersonallyProcuredMoveFetcher {
	return &personallyProcuredMoveFetcher{db}
}

//FetchPersonallyProcuredMove retrieves a PersonallyProcuredMove for a given UUID
func (f personallyProcuredMoveFetcher) FetchPersonallyProcuredMove(personallyProcuredMoveID uuid.UUID) (*models.PersonallyProcuredMove, error) {
	ppm := &models.PersonallyProcuredMove{}
	if err := f.db.Eager().Find(ppm, personallyProcuredMoveID); err != nil {
		switch err {
		case sql.ErrNoRows:
			return &models.PersonallyProcuredMove{}, ErrNotFound{personallyProcuredMoveID}
		default:
			return &models.PersonallyProcuredMove{}, err
		}
	}

	return ppm, nil
}
