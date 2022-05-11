package shipment

import (
	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

// shipmentUpdater is the concrete struct implementing the services.ShipmentUpdater interface
type shipmentUpdater struct {
	checks             []shipmentValidator
	mtoShipmentUpdater services.MTOShipmentUpdater
	ppmShipmentUpdater services.PPMShipmentUpdater
}

// NewShipmentUpdater creates a new shipmentUpdater struct with the basic checks and service dependencies.
func NewShipmentUpdater(mtoShipmentUpdater services.MTOShipmentUpdater, ppmShipmentUpdater services.PPMShipmentUpdater) services.ShipmentUpdater {
	return &shipmentUpdater{
		checks:             basicShipmentChecks(),
		mtoShipmentUpdater: mtoShipmentUpdater,
		ppmShipmentUpdater: ppmShipmentUpdater,
	}
}

// UpdateShipment updates a shipment, taking into account different shipment types and their needs.
func (s *shipmentUpdater) UpdateShipment(appCtx appcontext.AppContext, shipment *models.MTOShipment, eTag string) (*models.MTOShipment, error) {
	if err := validateShipment(appCtx, *shipment, s.checks...); err != nil {
		return nil, err
	}

	var mtoShipment *models.MTOShipment

	txErr := appCtx.NewTransaction(func(txnAppCtx appcontext.AppContext) (err error) {
		mtoShipment, err = s.mtoShipmentUpdater.UpdateMTOShipmentCustomer(txnAppCtx, shipment, eTag)

		if err != nil {
			return err
		}

		if shipment.ShipmentType != models.MTOShipmentTypePPM {
			return nil
		}

		shipment.PPMShipment.ShipmentID = mtoShipment.ID
		shipment.PPMShipment.Shipment = *mtoShipment

		ppmShipment, err := s.ppmShipmentUpdater.UpdatePPMShipmentWithDefaultCheck(txnAppCtx, shipment.PPMShipment, mtoShipment.ID)

		if err != nil {
			return err
		}

		// Update variables with latest versions
		mtoShipment = &ppmShipment.Shipment
		mtoShipment.PPMShipment = ppmShipment

		return nil
	})

	if txErr != nil {
		return nil, txErr
	}

	return mtoShipment, nil
}
