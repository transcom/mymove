package shipment

import (
	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

// shipmentCreator is the concrete struct implementing the services.ShipmentCreator interface
type shipmentCreator struct {
	checks                    []shipmentValidator
	mtoShipmentCreator        services.MTOShipmentCreator
	ppmShipmentCreator        services.PPMShipmentCreator
	boatShipmentCreator       services.BoatShipmentCreator
	mobileHomeShipmentCreator services.MobileHomeShipmentCreator
	shipmentRouter            services.ShipmentRouter
	moveTaskOrderUpdater      services.MoveTaskOrderUpdater
}

// NewShipmentCreator creates a new shipmentCreator struct with the basic checks and service dependencies.
func NewShipmentCreator(mtoShipmentCreator services.MTOShipmentCreator, ppmShipmentCreator services.PPMShipmentCreator, boatShipmentCreator services.BoatShipmentCreator, mobileHomeShipmentCreator services.MobileHomeShipmentCreator, shipmentRouter services.ShipmentRouter, moveTaskOrderUpdater services.MoveTaskOrderUpdater) services.ShipmentCreator {
	return &shipmentCreator{
		checks:                    basicShipmentChecks(),
		mtoShipmentCreator:        mtoShipmentCreator,
		ppmShipmentCreator:        ppmShipmentCreator,
		boatShipmentCreator:       boatShipmentCreator,
		mobileHomeShipmentCreator: mobileHomeShipmentCreator,
		shipmentRouter:            shipmentRouter,
		moveTaskOrderUpdater:      moveTaskOrderUpdater,
	}
}

// CreateShipment creates a shipment, taking into account different shipment types and their needs.
func (s *shipmentCreator) CreateShipment(appCtx appcontext.AppContext, shipment *models.MTOShipment) (*models.MTOShipment, error) {
	if err := validateShipment(appCtx, *shipment, s.checks...); err != nil {
		return nil, err
	}

	isPPMShipment := shipment.ShipmentType == models.MTOShipmentTypePPM
	isBoatShipment := (shipment.ShipmentType == models.MTOShipmentTypeBoatHaulAway || shipment.ShipmentType == models.MTOShipmentTypeBoatTowAway)
	isMobileHomeShipment := shipment.ShipmentType == models.MTOShipmentTypeMobileHome

	if isBoatShipment {
		// Match boatShipment.Type with shipmentType incase they are different
		if shipment.ShipmentType == models.MTOShipmentTypeBoatHaulAway && shipment.BoatShipment.Type != models.BoatShipmentTypeHaulAway {
			shipment.BoatShipment.Type = models.BoatShipmentTypeHaulAway
		} else if shipment.ShipmentType == models.MTOShipmentTypeBoatTowAway && shipment.BoatShipment.Type != models.BoatShipmentTypeTowAway {
			shipment.BoatShipment.Type = models.BoatShipmentTypeTowAway
		}
	}

	if shipment.Status == "" {
		if isPPMShipment || isBoatShipment || isMobileHomeShipment {
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
			// Update PPMType once shipment gets created.
			_, err = s.moveTaskOrderUpdater.UpdatePPMType(txnAppCtx, mtoShipment.MoveTaskOrderID)
			if err != nil {
				return err
			}
		}

		if isPPMShipment {
			mtoShipment.PPMShipment.ShipmentID = mtoShipment.ID
			mtoShipment.PPMShipment.Shipment = *mtoShipment

			_, err = s.ppmShipmentCreator.CreatePPMShipmentWithDefaultCheck(txnAppCtx, mtoShipment.PPMShipment)
			if err != nil {
				return err
			}

			// Update PPMType once shipment gets created.
			_, err = s.moveTaskOrderUpdater.UpdatePPMType(txnAppCtx, mtoShipment.MoveTaskOrderID)
			if err != nil {
				return err
			}
			return nil
		} else if isBoatShipment {
			mtoShipment.BoatShipment.ShipmentID = mtoShipment.ID
			mtoShipment.BoatShipment.Shipment = *mtoShipment

			_, err = s.boatShipmentCreator.CreateBoatShipmentWithDefaultCheck(txnAppCtx, mtoShipment.BoatShipment)

			if err != nil {
				return err
			}
			return nil
		} else if isMobileHomeShipment {
			mtoShipment.MobileHome.ShipmentID = mtoShipment.ID
			mtoShipment.MobileHome.Shipment = *mtoShipment

			_, err = s.mobileHomeShipmentCreator.CreateMobileHomeShipmentWithDefaultCheck(txnAppCtx, mtoShipment.MobileHome)

			if err != nil {
				return err
			}
			return nil
		}
		return nil
	})
	if txErr != nil {
		return nil, txErr
	}

	return mtoShipment, nil
}
