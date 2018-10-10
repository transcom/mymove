package operations

import (
	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/uuid"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/notifications"
)

// CancelMove is a struct on the service object layer to handle move cancelations
type CancelMove struct {
	Operation
	Notifier notifications.NotificationSender
}

// * Make transaction occur on the service object level rather than save level
// * Check attempted transition is valid (potentially pass in accepted states) (maybe abstract this out)
// * Return a specific error type from the operations level
// * Write a simpler example maybe (e.g. Cancel PPM, which is called here)

// Run runs CancelMove
func (cm *CancelMove) Run(moveID uuid.UUID, cancelReason string) (move *models.Move) {
	move, err := models.FetchMove(cm.DB, cm.Session, moveID)
	if err != nil {
		cm.Err = err
		return nil
	}

	// We can cancel any move that isn't already complete.
	if move.Status == models.MoveStatusCOMPLETED || move.Status == models.MoveStatusCANCELED {
		cm.Err = errors.Wrap(models.ErrInvalidTransition, "Cancel")
		return nil
	}

	move.Status = models.MoveStatusCANCELED

	// If a reason was submitted, add it to the move record.
	if cancelReason != "" {
		move.CancelReason = &cancelReason
	}

	// TODO: cancel any shipments

	cm.DB.Transaction(func(db *pop.Connection) error {
		transactionError := errors.New("Rollback The transaction")

		for _, ppm := range move.PersonallyProcuredMoves {
			cancelPPM := CancelPPM{Operation: Operation{DB: db, Logger: cm.Logger, Session: cm.Session}}
			cancelPPM.Run(ppm.ID)

			if cancelPPM.Err != nil {
				cm.Err = errors.Wrap(cancelPPM.Err, "Failed to cancel PPM")
				return nil
			}
			if cancelPPM.Verrs != nil {
				cm.Verrs = cancelPPM.Verrs
				return nil
			}
		}

		if cm.hadErrors(db.ValidateAndSave(move)) {
			return transactionError
		}
		return nil
	})

	if cm.Err != nil {
		return nil
	}

	err = cm.Notifier.SendNotification(
		notifications.NewMoveCanceled(cm.DB, cm.Logger, cm.Session, moveID),
	)

	if err != nil {
		cm.Logger.Error("problem sending email to user", zap.Error(err))
		cm.Err = err
		return nil
	}

	return move
}
