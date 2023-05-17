package sitaddressupdate

import (
	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/route"
	"github.com/transcom/mymove/pkg/services"
)

type sitAddressUpdateRequestCreator struct {
	planner        route.Planner
	addressCreator services.AddressCreator
	checks         []sitAddressUpdateValidator
}

func NewsitAddressUpdateRequestCreator(planner route.Planner, addressCreator services.AddressCreator) services.SITAddressUpdateRequestCreator {
	return &sitAddressUpdateRequestCreator{
		planner:        planner,
		addressCreator: addressCreator,
		checks: []sitAddressUpdateValidator{
			checkRequiredFields(),
			checkPrimeRequiredFields(),
		},
	}
}

// CreateSITAddressUpdateRequest creates a SIT address update for requests with a distance greater than 50 miles
func (f *sitAddressUpdateRequestCreator) CreateSITAddressUpdateRequest(appCtx appcontext.AppContext, sitAddressUpdateRequest *models.SITAddressUpdate) (*models.SITAddressUpdate, error) {
	var err error
	if err = validateSITAddressUpdate(appCtx, sitAddressUpdateRequest, f.checks...); err != nil {
		return nil, err
	}

	sitAddressUpdateRequest.Status = models.SITAddressUpdateStatusRequested

	txErr := appCtx.NewTransaction(func(txnAppCtx appcontext.AppContext) (err error) {
		// Grabbing the service item in question - must be an approved service item
		var serviceItem models.MTOServiceItem
		err = txnAppCtx.DB().Eager("SITDestinationFinalAddress").Where("id = ?", sitAddressUpdateRequest.MTOServiceItemID).Where("status = ?", models.MTOServiceItemStatusApproved).First(&serviceItem)
		if err != nil {
			return err
		}

		//The SITDestinationFinalAddress is the most up to date address on the service item, so that is the one we set as "OldAddress" since we wish to update it
		sitAddressUpdateRequest.OldAddressID = *serviceItem.SITDestinationFinalAddressID
		sitAddressUpdateRequest.OldAddress = *serviceItem.SITDestinationFinalAddress

		//We create an address from the new address being requested by the prime
		newAddress, err := f.addressCreator.CreateAddress(txnAppCtx, &sitAddressUpdateRequest.NewAddress)
		if err != nil {
			return err
		}

		//Set that new created address in our update request
		sitAddressUpdateRequest.NewAddressID = newAddress.ID
		sitAddressUpdateRequest.NewAddress = *newAddress

		//We calculate and set the distance between the old and new address
		sitAddressUpdateRequest.Distance, err = f.planner.TransitDistance(appCtx, &sitAddressUpdateRequest.OldAddress, &sitAddressUpdateRequest.NewAddress)
		if err != nil {
			return err
		} else if sitAddressUpdateRequest.Distance > 50 {
			return apperror.NewInvalidInputError(sitAddressUpdateRequest.ID, nil, nil, "Distance exceeds 50 miles from final address, should be 50 miles or less.")
		}

		verrs, err := appCtx.DB().ValidateAndCreate(sitAddressUpdateRequest)

		if verrs != nil && verrs.HasAny() {
			return apperror.NewInvalidInputError(sitAddressUpdateRequest.ID, nil, verrs, "Invalid input found while creating sit address update request.")
		} else if err != nil {
			return apperror.NewQueryError("SITAddressUpdate", err, "Unable to create SIT address update request.")
		}

		return nil
	})

	if txErr != nil {
		return nil, txErr
	}

	return sitAddressUpdateRequest, nil
}
