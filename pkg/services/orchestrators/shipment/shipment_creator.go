package shipment

import (
	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

// shipmentCreator is the concrete struct implementing the services.ShipmentCreator interface
type shipmentCreator struct {
	checks             []shipmentValidator
	mtoShipmentCreator services.MTOShipmentCreator
	ppmShipmentCreator services.PPMShipmentCreator
	shipmentRouter     services.ShipmentRouter
}

// NewShipmentCreator creates a new shipmentCreator struct with the basic checks and service dependencies.
func NewShipmentCreator(mtoShipmentCreator services.MTOShipmentCreator, ppmShipmentCreator services.PPMShipmentCreator, shipmentRouter services.ShipmentRouter) services.ShipmentCreator {
	return &shipmentCreator{
		checks:             basicShipmentChecks(),
		mtoShipmentCreator: mtoShipmentCreator,
		ppmShipmentCreator: ppmShipmentCreator,
		shipmentRouter:     shipmentRouter,
	}
}

// CreateShipment creates a shipment, taking into account different shipment types and their needs.
func (s *shipmentCreator) CreateShipment(appCtx appcontext.AppContext, shipment *models.MTOShipment) (*models.MTOShipment, error) {
	dataToValidate := models.MTOShipmentUpdate{ShipmentType: shipment.ShipmentType}
	if err := validateShipment(appCtx, dataToValidate, s.checks...); err != nil {
		return nil, err
	}

	isPPMShipment := shipment.ShipmentType == models.MTOShipmentTypePPM

	if shipment.Status == "" {
		if isPPMShipment {
			shipment.Status = models.MTOShipmentStatusDraft
		} else {
			// TODO: remove this status change once MB-3428 is implemented and can update to Submitted on second page
			err := s.shipmentRouter.Submit(appCtx, shipment)
			if err != nil {
				return nil, err
			}
		}
	}

	var mtoShipment *models.MTOShipment

	txErr := appCtx.NewTransaction(func(txnAppCtx appcontext.AppContext) (err error) {
		mtoShipment, err = s.mtoShipmentCreator.CreateMTOShipment(txnAppCtx, shipment)

		if err != nil {
			return err
		}

		if !isPPMShipment {
			return nil
		}

		mtoShipment.PPMShipment.ShipmentID = mtoShipment.ID
		mtoShipment.PPMShipment.Shipment = *mtoShipment

		_, err = s.ppmShipmentCreator.CreatePPMShipmentWithDefaultCheck(txnAppCtx, mtoShipment.PPMShipment)

		if err != nil {
			return err
		}

		return nil
	})

	if txErr != nil {
		return nil, txErr
	}

	return mtoShipment, nil
}
