package operations

import (
	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/uuid"
	"github.com/gobuffalo/validate"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/models"
)

// StatusOperation is the interface for models' status operations
type StatusOperation interface {
	DB() *pop.Connection
	Logger() *zap.Logger
	Session() *auth.Session
	Verrs() *validate.Errors
	err() error
}

type cancelMove struct {
	db      *pop.Connection
	logger  *zap.Logger
	session *auth.Session
	verrs   *validate.Errors
	err     error
}

// DB returns a POP db connection for the operation
func (cancelMove *cancelMove) DB() *pop.Connection {
	return cancelMove.db
}

// Logger returns the logger to use in the operation
func (cancelMove *cancelMove) Logger() *zap.Logger {
	return cancelMove.logger
}

func (cancelMove *cancelMove) Run(moveID uuid.UUID, cancelReason string) (move *models.Move, err error) {
	move, err := models.FetchMove(db, session, moveID)
	if err != nil {
		return nil, err
	}

	// We can cancel any move that isn't already complete.
	if move.Status == models.MoveStatusCOMPLETED || move.Status == models.MoveStatusCANCELED {
		return nil, errors.Wrap(models.ErrInvalidTransition, "Cancel")
	}

	move.Status = models.MoveStatusCANCELED

	// If a reason was submitted, add it to the move record.
	if reason != "" {
		move.CancelReason = &reason
	}

	// This will work only if you use the PPM in question rather than a var representing it
	// i.e. you can't use _, ppm := range PPMs, has to be PPMS[i] as below
	for i := range move.PersonallyProcuredMoves {
		err := move.PersonallyProcuredMoves[i].Cancel()
		if err != nil {
			return nil, err
		}
	}

	// Save move, orders, and PPMs statuses
	verrs, err := models.SaveMoveDependencies(db, move)
	if err != nil || verrs.HasAny() {
		return handlers.ResponseForVErrors(logger, verrs, err)
	}

	err = h.NotificationSender().SendNotification(
		notifications.NewMoveCanceled(h.DB(), h.Logger(), session, moveID),
	)

	if err != nil {
		h.Logger().Error("problem sending email to user", zap.Error(err))
		return handlers.ResponseForError(h.Logger(), err)
	}

	return &move, nil
}
