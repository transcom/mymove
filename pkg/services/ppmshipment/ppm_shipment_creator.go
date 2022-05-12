package ppmshipment

import (
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

// ppmShipmentCreator sets up the service object, and passes in
type ppmShipmentCreator struct {
	checks []ppmShipmentValidator
}

// NewPPMShipmentCreator creates a new struct with the service dependencies
func NewPPMShipmentCreator() services.PPMShipmentCreator {
	return &ppmShipmentCreator{
		checks: []ppmShipmentValidator{
			checkShipmentID(),
			checkPPMShipmentID(),
			checkRequiredFields(),
		},
	}
}

// CreatePPMShipmentWithDefaultCheck passes a validator key to CreatePPMShipment
func (f *ppmShipmentCreator) CreatePPMShipmentWithDefaultCheck(appCtx appcontext.AppContext, ppmShipment *models.PPMShipment) (*models.PPMShipment, error) {
	return f.createPPMShipment(appCtx, ppmShipment, f.checks...)
}

func (f *ppmShipmentCreator) createPPMShipment(appCtx appcontext.AppContext, ppmShipment *models.PPMShipment, checks ...ppmShipmentValidator) (*models.PPMShipment, error) {
	transactionError := appCtx.NewTransaction(func(txnAppCtx appcontext.AppContext) error {
		if ppmShipment.Shipment.ShipmentType != models.MTOShipmentTypePPM {
			return apperror.NewInvalidInputError(uuid.Nil, nil, nil, "MTO shipment type must be PPM shipment")
		}

		if ppmShipment.Shipment.Status != models.MTOShipmentStatusDraft {
			return apperror.NewInvalidInputError(uuid.Nil, nil, nil, "Must have a DRAFT status associated with MTO shipment")
		}

		if ppmShipment.Status == "" {
			ppmShipment.Status = models.PPMShipmentStatusDraft
		} else if ppmShipment.Status != models.PPMShipmentStatusDraft {
			return apperror.NewInvalidInputError(uuid.Nil, nil, nil, "Must have a DRAFT status associated with PPM shipment")
		}

		// Validate the ppmShipment, and return an error
		if err := validatePPMShipment(txnAppCtx, *ppmShipment, nil, &ppmShipment.Shipment, checks...); err != nil {
			return err
		}

		// Validate ppm shipment model object and save it to DB
		verrs, err := txnAppCtx.DB().ValidateAndCreate(ppmShipment)

		// Check validation errors
		if verrs != nil && verrs.HasAny() {
			return apperror.NewInvalidInputError(uuid.Nil, err, verrs, "Invalid input found while creating the PPM shipment.")
		} else if err != nil {
			// If the error is something else (this is unexpected), we create a QueryError
			return apperror.NewQueryError("PPM Shipment", err, "")
		}

		return err
	})

	if transactionError != nil {
		return nil, transactionError
	}

	return ppmShipment, nil
}
