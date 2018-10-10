package operations

import (
	"github.com/gobuffalo/uuid"
	"github.com/pkg/errors"

	"github.com/transcom/mymove/pkg/models"
)

// CancelPPM is a struct on the service object layer to handle PPM cancelations
type CancelPPM struct {
	Operation
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

	if cp.hadErrors(cp.DB.ValidateAndSave(ppm)) {
		return nil
	}

	return ppm
}
