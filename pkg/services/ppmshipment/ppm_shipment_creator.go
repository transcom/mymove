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
	mtoShipmentCreator services.MTOShipmentCreator
	estimator          services.PPMEstimator
	checks             []ppmShipmentValidator
}

// NewPPMShipmentCreator creates a new struct with the service dependencies
func NewPPMShipmentCreator(mtoShipmentCreator services.MTOShipmentCreator, estimator services.PPMEstimator) services.PPMShipmentCreator {
	return &ppmShipmentCreator{
		mtoShipmentCreator: mtoShipmentCreator,
		estimator:          estimator,
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
	// Start a transaction that will create a Shipment, then create a PPM
	transactionError := appCtx.NewTransaction(func(txnAppCtx appcontext.AppContext) error {
		var err error
		if ppmShipment.Shipment.ShipmentType == "" {
			ppmShipment.Shipment.ShipmentType = models.MTOShipmentTypePPM
		} else if ppmShipment.Shipment.ShipmentType != models.MTOShipmentTypePPM {
			return apperror.NewInvalidInputError(uuid.Nil, nil, nil, "MTO shipment type must be PPM shipment")
		}

		if ppmShipment.Shipment.Status == "" {
			ppmShipment.Shipment.Status = models.MTOShipmentStatusDraft
		} else if ppmShipment.Shipment.Status != models.MTOShipmentStatusDraft && ppmShipment.Shipment.Status != models.MTOShipmentStatusSubmitted {
			return apperror.NewInvalidInputError(uuid.Nil, nil, nil, "Must have a DRAFT or SUBMITTED status associated with MTO shipment")
		}

		if ppmShipment.Status == "" {
			ppmShipment.Status = models.PPMShipmentStatusDraft
		} else if ppmShipment.Status != models.PPMShipmentStatusDraft && ppmShipment.Status != models.PPMShipmentStatusSubmitted {
			return apperror.NewInvalidInputError(uuid.Nil, nil, nil, "Must have a DRAFT or SUBMITTED status associated with PPM shipment")
		}

		// NOTE: The ppm may require a service item for pricing.Passing an HHG service item may be sufficient for the pricer.
		createShipment, err := f.mtoShipmentCreator.CreateMTOShipment(txnAppCtx, &ppmShipment.Shipment, nil)
		// Check that mtoshipment is created. If not, bail out.
		if err != nil {
			return err
		}

		ppmShipment.ShipmentID = createShipment.ID
		// Update the model with ppm shipment data:
		ppmShipment.Shipment = *createShipment
		// Validate the ppmShipment, and return an error
		err = validatePPMShipment(txnAppCtx, *ppmShipment, nil, &ppmShipment.Shipment, checks...)
		if err != nil {
			return err
		}

		estimatedIncentive, err := f.estimator.EstimateIncentiveWithDefaultChecks(appCtx, models.PPMShipment{}, ppmShipment)
		if err != nil {
			return err
		}
		ppmShipment.EstimatedIncentive = estimatedIncentive

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
