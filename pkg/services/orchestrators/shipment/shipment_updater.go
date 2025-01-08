package shipment

import (
	"database/sql"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/route"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/unit"
)

// shipmentUpdater is the concrete struct implementing the services.ShipmentUpdater interface
type shipmentUpdater struct {
	checks                    []shipmentValidator
	mtoShipmentUpdater        services.MTOShipmentUpdater
	ppmShipmentUpdater        services.PPMShipmentUpdater
	boatShipmentUpdater       services.BoatShipmentUpdater
	mobileHomeShipmentUpdater services.MobileHomeShipmentUpdater
	mtoServiceItemCreator     services.MTOServiceItemCreator
}

// NewShipmentUpdater creates a new shipmentUpdater struct with the basic checks and service dependencies.
func NewShipmentUpdater(mtoShipmentUpdater services.MTOShipmentUpdater, ppmShipmentUpdater services.PPMShipmentUpdater, boatShipmentUpdater services.BoatShipmentUpdater, mobileHomeShipmentUpdater services.MobileHomeShipmentUpdater, mtoServiceItemCreator services.MTOServiceItemCreator) services.ShipmentUpdater {
	return &shipmentUpdater{
		checks:                    basicShipmentChecks(),
		mtoShipmentUpdater:        mtoShipmentUpdater,
		ppmShipmentUpdater:        ppmShipmentUpdater,
		boatShipmentUpdater:       boatShipmentUpdater,
		mobileHomeShipmentUpdater: mobileHomeShipmentUpdater,
		mtoServiceItemCreator:     mtoServiceItemCreator,
	}
}

// UpdateShipment updates a shipment, taking into account different shipment types and their needs.
func (s *shipmentUpdater) UpdateShipment(appCtx appcontext.AppContext, shipment *models.MTOShipment, eTag string, api string, planner route.Planner) (*models.MTOShipment, error) {
	if err := validateShipment(appCtx, *shipment, s.checks...); err != nil {
		return nil, err
	}

	var mtoShipment *models.MTOShipment

	txErr := appCtx.NewTransaction(func(txnAppCtx appcontext.AppContext) (err error) {
		mtoShipment, err = s.mtoShipmentUpdater.UpdateMTOShipment(txnAppCtx, shipment, eTag, api)

		if err != nil {
			return err
		}

		if mtoShipment != nil && planner != nil {
			if mtoShipment.ShipmentType != models.MTOShipmentTypePPM && (shipment.PrimeEstimatedWeight != nil || mtoShipment.PrimeEstimatedWeight != nil) && mtoShipment.Status == models.MTOShipmentStatusApproved {
				for index, serviceItem := range mtoShipment.MTOServiceItems {
					var estimatedWeightToUse unit.Pound
					if shipment.PrimeEstimatedWeight != nil {
						estimatedWeightToUse = *shipment.PrimeEstimatedWeight
					} else {
						estimatedWeightToUse = *mtoShipment.PrimeEstimatedWeight
					}
					mtoShipment.MTOServiceItems[index].EstimatedWeight = &estimatedWeightToUse
					serviceItemEstimatedPrice, err := s.mtoServiceItemCreator.FindEstimatedPrice(appCtx, &serviceItem, *mtoShipment)
					if serviceItemEstimatedPrice != 0 && err == nil {
						mtoShipment.MTOServiceItems[index].PricingEstimate = &serviceItemEstimatedPrice
					}
					if err != nil {
						return err
					}
				}
			}
		}

		if mtoShipment.MTOServiceItems != nil {
			_, mtoErr := appCtx.DB().ValidateAndUpdate(&mtoShipment.MTOServiceItems)
			if mtoErr != nil {
				return mtoErr
			}
		}

		isBoatShipment := shipment.ShipmentType == models.MTOShipmentTypeBoatHaulAway || shipment.ShipmentType == models.MTOShipmentTypeBoatTowAway

		if shipment.ShipmentType == models.MTOShipmentTypePPM {
			shipment.PPMShipment.ShipmentID = mtoShipment.ID
			shipment.PPMShipment.Shipment = *mtoShipment

			ppmShipment, err := s.ppmShipmentUpdater.UpdatePPMShipmentWithDefaultCheck(txnAppCtx, shipment.PPMShipment, mtoShipment.ID)

			if err != nil {
				return err
			}

			// getting updated shipment since market code value was updated after PPM creation
			var updatedShipment models.MTOShipment
			err = txnAppCtx.DB().Find(&updatedShipment, mtoShipment.ID)
			if err != nil && err != sql.ErrNoRows {
				return err
			}
			if mtoShipment.MarketCode != updatedShipment.MarketCode {
				mtoShipment.MarketCode = updatedShipment.MarketCode
			}
			// since the shipment was updated, we need to ensure we have the most recent eTag
			if mtoShipment.UpdatedAt != updatedShipment.UpdatedAt {
				mtoShipment.UpdatedAt = updatedShipment.UpdatedAt
			}
			// Update variables with latest versions
			mtoShipment = &updatedShipment
			mtoShipment.PPMShipment = ppmShipment

			return nil
		} else if isBoatShipment && shipment.BoatShipment != nil {
			shipment.BoatShipment.ShipmentID = mtoShipment.ID
			shipment.BoatShipment.Shipment = *mtoShipment

			// Match boatShipment.Type with shipmentType incase they are different
			if shipment.ShipmentType == models.MTOShipmentTypeBoatHaulAway && shipment.BoatShipment.Type != models.BoatShipmentTypeHaulAway {
				shipment.BoatShipment.Type = models.BoatShipmentTypeHaulAway
			} else if shipment.ShipmentType == models.MTOShipmentTypeBoatTowAway && shipment.BoatShipment.Type != models.BoatShipmentTypeTowAway {
				shipment.BoatShipment.Type = models.BoatShipmentTypeTowAway
			}

			boatShipment, err := s.boatShipmentUpdater.UpdateBoatShipmentWithDefaultCheck(txnAppCtx, shipment.BoatShipment, mtoShipment.ID)

			if err != nil {
				return err
			}

			// Update variables with latest versions
			mtoShipment = &boatShipment.Shipment
			mtoShipment.BoatShipment = boatShipment

			return nil
		} else if shipment.ShipmentType == models.MTOShipmentTypeMobileHome && shipment.MobileHome != nil {
			shipment.MobileHome.ShipmentID = mtoShipment.ID
			shipment.MobileHome.Shipment = *mtoShipment

			mobileHomeShipment, err := s.mobileHomeShipmentUpdater.UpdateMobileHomeShipmentWithDefaultCheck(txnAppCtx, shipment.MobileHome, mtoShipment.ID)

			if err != nil {
				return err
			}

			// Update variables with latest versions
			mtoShipment = &mobileHomeShipment.Shipment
			mtoShipment.MobileHome = mobileHomeShipment

			return nil
		}

		return nil

	})

	if txErr != nil {
		return nil, txErr
	}

	return mtoShipment, nil
}
