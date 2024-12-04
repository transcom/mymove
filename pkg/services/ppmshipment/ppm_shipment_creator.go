package ppmshipment

import (
	"fmt"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

// ppmShipmentCreator sets up the service object, and passes in
type ppmShipmentCreator struct {
	estimator      services.PPMEstimator
	checks         []ppmShipmentValidator
	addressCreator services.AddressCreator
}

// NewPPMShipmentCreator creates a new struct with the service dependencies
func NewPPMShipmentCreator(estimator services.PPMEstimator, addressCreator services.AddressCreator) services.PPMShipmentCreator {
	return &ppmShipmentCreator{
		estimator:      estimator,
		addressCreator: addressCreator,
		checks: []ppmShipmentValidator{
			checkShipmentID(),
			checkPPMShipmentID(),
			checkRequiredFields(),
			checkPPMShipmentSequenceValidForCreate(),
		},
	}
}

// CreatePPMShipmentWithDefaultCheck passes a validator key to CreatePPMShipment
func (f *ppmShipmentCreator) CreatePPMShipmentWithDefaultCheck(appCtx appcontext.AppContext, ppmShipment *models.PPMShipment) (*models.PPMShipment, error) {
	return f.createPPMShipment(appCtx, ppmShipment, f.checks...)
}

func (f *ppmShipmentCreator) createPPMShipment(appCtx appcontext.AppContext, ppmShipment *models.PPMShipment, checks ...ppmShipmentValidator) (*models.PPMShipment, error) {
	var address *models.Address
	var err error

	transactionError := appCtx.NewTransaction(func(txnAppCtx appcontext.AppContext) error {
		if ppmShipment.Shipment.ShipmentType != models.MTOShipmentTypePPM {
			return apperror.NewInvalidInputError(uuid.Nil, nil, nil, "MTO shipment type must be PPM shipment")
		}

		if ppmShipment.Shipment.Status != models.MTOShipmentStatusDraft && ppmShipment.Shipment.Status != models.MTOShipmentStatusSubmitted {
			return apperror.NewInvalidInputError(uuid.Nil, nil, nil, "Must have a DRAFT or SUBMITTED status associated with MTO shipment")
		}

		if ppmShipment.Status == "" {
			ppmShipment.Status = models.PPMShipmentStatusDraft
		} else if ppmShipment.Status != models.PPMShipmentStatusDraft && ppmShipment.Status != models.PPMShipmentStatusSubmitted {
			return apperror.NewInvalidInputError(uuid.Nil, nil, nil, "Must have a DRAFT or SUBMITTED status associated with PPM shipment")
		}

		// create pickup and destination addresses
		if ppmShipment.PickupAddress != nil {
			address, err = f.addressCreator.CreateAddress(txnAppCtx, ppmShipment.PickupAddress)
			if err != nil {
				switch err := err.(type) {
				case apperror.EventError:
					return err
				default:
					return fmt.Errorf("failed to create pickup address %e", err)
				}
			}
			ppmShipment.PickupAddressID = &address.ID
			ppmShipment.PickupAddress = address
		}

		if ppmShipment.SecondaryPickupAddress != nil {
			address, err = f.addressCreator.CreateAddress(txnAppCtx, ppmShipment.SecondaryPickupAddress)
			if err != nil {
				switch err := err.(type) {
				case apperror.EventError:
					return err
				default:
					return fmt.Errorf("failed to create secondary pickup address %e", err)
				}
			}
			ppmShipment.SecondaryPickupAddressID = &address.ID
			// ensure HasSecondaryPickupAddress property is set true on create
			ppmShipment.HasSecondaryPickupAddress = models.BoolPointer(true)
		}

		if ppmShipment.TertiaryPickupAddress != nil {
			address, err = f.addressCreator.CreateAddress(txnAppCtx, ppmShipment.TertiaryPickupAddress)
			if err != nil {
				return fmt.Errorf("failed to create secondary pickup address %e", err)
			}
			ppmShipment.TertiaryPickupAddressID = &address.ID
			// ensure HasTertiaryPickupAddress property is set true on create
			ppmShipment.HasTertiaryPickupAddress = models.BoolPointer(true)
		}

		if ppmShipment.DestinationAddress != nil {
			address, err = f.addressCreator.CreateAddress(txnAppCtx, ppmShipment.DestinationAddress)
			if err != nil {
				switch err := err.(type) {
				case apperror.EventError:
					return err
				default:
					return fmt.Errorf("failed to create destination address %e", err)
				}
			}
			ppmShipment.DestinationAddressID = &address.ID
			ppmShipment.DestinationAddress = address
		}

		if ppmShipment.SecondaryDestinationAddress != nil {
			address, err = f.addressCreator.CreateAddress(txnAppCtx, ppmShipment.SecondaryDestinationAddress)
			if err != nil {
				switch err := err.(type) {
				case apperror.EventError:
					return err
				default:
					return fmt.Errorf("failed to create secondary destination address %e", err)
				}
			}
			ppmShipment.SecondaryDestinationAddressID = &address.ID
			// ensure HasSecondaryDestinationAddress property is set true on create
			ppmShipment.HasSecondaryDestinationAddress = models.BoolPointer(true)
		}

		if ppmShipment.TertiaryDestinationAddress != nil {
			address, err = f.addressCreator.CreateAddress(txnAppCtx, ppmShipment.TertiaryDestinationAddress)
			if err != nil {
				return fmt.Errorf("failed to create tertiary delivery address %e", err)
			}
			ppmShipment.TertiaryDestinationAddressID = &address.ID
			// ensure HasTertiaryDestinationAddress property is set true on create
			ppmShipment.HasTertiaryDestinationAddress = models.BoolPointer(true)
		}

		// Validate the ppmShipment, and return an error
		if err := validatePPMShipment(txnAppCtx, *ppmShipment, nil, &ppmShipment.Shipment, checks...); err != nil {
			return err
		}

		estimatedIncentive, estimatedSITCost, err := f.estimator.EstimateIncentiveWithDefaultChecks(appCtx, models.PPMShipment{}, ppmShipment)
		if err != nil {
			return err
		}
		ppmShipment.EstimatedIncentive = estimatedIncentive
		ppmShipment.SITEstimatedCost = estimatedSITCost

		maxIncentive, err := f.estimator.MaxIncentive(appCtx, models.PPMShipment{}, ppmShipment)
		if err != nil {
			return err
		}
		ppmShipment.MaxIncentive = maxIncentive

		// Validate ppm shipment model object and save it to DB
		verrs, err := txnAppCtx.DB().ValidateAndCreate(ppmShipment)
		// Check validation errors
		if verrs != nil && verrs.HasAny() {
			return apperror.NewInvalidInputError(uuid.Nil, err, verrs, "Invalid input found while creating the PPM shipment.")
		} else if err != nil {
			// If the error is something else (this is unexpected), we create a QueryError
			return apperror.NewQueryError("PPM Shipment", err, "")
		}

		// updating the shipment after PPM creation due to addresses not being created until PPM shipment is created
		// when populating the market_code column, it is considered domestic if both pickup & dest on the PPM are CONUS addresses
		var mtoShipment models.MTOShipment
		if err := txnAppCtx.DB().Find(&mtoShipment, ppmShipment.ShipmentID); err != nil {
			return err
		}
		if ppmShipment.PickupAddress != nil && ppmShipment.DestinationAddress != nil &&
			ppmShipment.PickupAddress.IsOconus != nil && ppmShipment.DestinationAddress.IsOconus != nil {
			pickupAddress := ppmShipment.PickupAddress
			destAddress := ppmShipment.DestinationAddress
			marketCode, err := models.DetermineMarketCode(pickupAddress, destAddress)
			if err != nil {
				return err
			}
			mtoShipment.MarketCode = marketCode
			if err := txnAppCtx.DB().Update(&mtoShipment); err != nil {
				return err
			}
			ppmShipment.Shipment = mtoShipment
		}

		return err
	})

	if transactionError != nil {
		return nil, transactionError
	}

	return ppmShipment, nil
}
