package operations

import (
	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/uuid"
	"github.com/gobuffalo/validate"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/notifications"
)

// CancelMove is a struct on the service object layer to handle move cancelations
type CancelMove struct {
	DB       *pop.Connection
	Logger   *zap.Logger
	Session  *auth.Session
	Notifier notifications.NotificationSender
	Verrs    *validate.Errors
	Err      error
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

	// This will work only if you use the PPM in question rather than a var representing it
	// i.e. you can't use _, ppm := range PPMs, has to be PPMS[i] as below
	for i := range move.PersonallyProcuredMoves {
		err := move.PersonallyProcuredMoves[i].Cancel()
		if err != nil {
			cm.Err = err
			return nil
		}
	}

	cm.DB.Transaction(func(db *pop.Connection) error {
		transactionError := errors.New("Rollback The transaction")

		if cm.hadErrors(cm.savePPMsAndDependencies(move)) {
			return transactionError
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

func (cm *CancelMove) hadErrors(verrs *validate.Errors, saveErr error) bool {
	if saveErr != nil {
		cm.Err = errors.Wrap(saveErr, "error saving model")
		return true
	}
	if verrs.HasAny() {
		cm.Verrs = verrs
		cm.Err = errors.New("Model validation failure")
		return true
	}
	return false
}

// SaveMoveDependencies safely saves a Move status, ppms' advances' statuses, orders statuses,
// and shipment GBLOCs.
func (cm *CancelMove) savePPMsAndDependencies(move *models.Move) (*validate.Errors, error) {
	validationErrors := validate.NewErrors()

	for _, ppm := range move.PersonallyProcuredMoves {
		if ppm.Advance != nil {
			if verrs, err := cm.DB.ValidateAndSave(ppm.Advance); verrs.HasAny() || err != nil {
				validationErrors.Append(verrs)
				return validationErrors, errors.Wrap(err, "Error Saving Advance")
			}
		}

		if verrs, err := cm.DB.ValidateAndSave(&ppm); verrs.HasAny() || err != nil {
			validationErrors.Append(verrs)
			return validationErrors, errors.Wrap(err, "Error Saving PPM")
		}
	}

	return validationErrors, nil
}

func getGbloc(db *pop.Connection, dutyStationID uuid.UUID) (gbloc string, err error) {
	transportationOffice, err := models.FetchDutyStationTransportationOffice(db, dutyStationID)
	if err != nil {
		return "", errors.Wrap(err, "could not load transportation office for duty station")
	}
	return transportationOffice.Gbloc, nil
}
