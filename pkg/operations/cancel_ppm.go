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

// CancelPPM is a struct on the service object layer to handle PPM cancelations
type CancelPPM struct {
	DB      *pop.Connection
	Logger  *zap.Logger
	Session *auth.Session
	Verrs   *validate.Errors
	Err     error
}

// Run runs CancelPPM
func (cp *CancelPPM) Run(ppmID uuid.UUID) (ppm *models.PersonallyProcuredMove) {
	ppm, err := models.FetchPersonallyProcuredMove(cp.DB, cp.Session, ppmID)
	if err != nil {
		cp.Err = err
		return nil
	}
	if ppm.Status == models.PPMStatusCOMPLETED || ppm.Status == models.PPMStatusCANCELED {
		cp.Err = errors.Wrap(models.ErrInvalidTransition, "Cancel")
		return nil
	}

	ppm.Status = models.PPMStatusCANCELED

	if verrs, err := cp.DB.ValidateAndSave(&ppm); verrs.HasAny() || err != nil {
		cp.Verrs = verrs
		cp.Err = err
		return nil
	}

	return ppm
}
