package sitaddressupdate

import (
	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/etag"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/route"
	"github.com/transcom/mymove/pkg/services"
)

// approvedSITAddressUpdateRequestCreator is the concrete struct implementing the services.ApprovedSITAddressUpdateRequestCreator interface
type approvedSITAddressUpdateRequestCreator struct {
	planner            route.Planner
	addressCreator     services.AddressCreator
	serviceItemUpdater services.MTOServiceItemUpdater
	checks             []sitAddressUpdateValidator
}

// NewApprovedOfficeSITAddressUpdateCreator creates a new struct with the service dependencies
func NewApprovedOfficeSITAddressUpdateCreator(planner route.Planner, addressCreator services.AddressCreator, serviceItemUpdater services.MTOServiceItemUpdater) services.ApprovedSITAddressUpdateRequestCreator {
	return &approvedSITAddressUpdateRequestCreator{
		planner:            planner,
		addressCreator:     addressCreator,
		serviceItemUpdater: serviceItemUpdater,
		checks: []sitAddressUpdateValidator{
			checkAndValidateRequiredFields(),
			checkTOORequiredFields(),
			checkServiceItem(),
		},
	}
}

// CreateSITAddressUpdate creates a SIT Address Update
func (f *approvedSITAddressUpdateRequestCreator) CreateApprovedSITAddressUpdate(appCtx appcontext.AppContext, sitAddressUpdate *models.SITAddressUpdate) (*models.SITAddressUpdate, error) {
	var err error
	if err = validateSITAddressUpdate(appCtx, sitAddressUpdate, f.checks...); err != nil {
		return nil, err
	}

	sitAddressUpdate.Status = models.SITAddressUpdateStatusApproved

	txErr := appCtx.NewTransaction(func(txnAppCtx appcontext.AppContext) (err error) {
		var serviceItem models.MTOServiceItem
		err = txnAppCtx.DB().Eager("SITDestinationFinalAddress", "SITDestinationOriginalAddress").Where("id = ?", sitAddressUpdate.MTOServiceItemID).First(&serviceItem)
		if err != nil {
			return err
		}

		sitAddressUpdate.OldAddressID = *serviceItem.SITDestinationFinalAddressID
		sitAddressUpdate.OldAddress = *serviceItem.SITDestinationFinalAddress

		newAddress, err := f.addressCreator.CreateAddress(txnAppCtx, &sitAddressUpdate.NewAddress)
		if err != nil {
			return err
		}
		sitAddressUpdate.NewAddressID = newAddress.ID
		sitAddressUpdate.NewAddress = *newAddress

		sitAddressUpdate.Distance, err = f.planner.ZipTransitDistance(appCtx, serviceItem.SITDestinationOriginalAddress.PostalCode, sitAddressUpdate.NewAddress.PostalCode)
		if err != nil {
			return err
		}

		verrs, err := txnAppCtx.DB().ValidateAndCreate(sitAddressUpdate)

		if verrs.HasAny() {
			return apperror.NewInvalidInputError(sitAddressUpdate.ID, nil, verrs, "Invalid input found while creating the SIT Address Update.")
		} else if err != nil {
			return apperror.NewQueryError("SITAddressUpdate", err, "Unable to create SIT Address Update")
		}

		serviceItem.SITDestinationFinalAddressID = &newAddress.ID
		serviceItem.SITDestinationFinalAddress = newAddress

		updatedServiceItem, err := f.serviceItemUpdater.UpdateMTOServiceItemBasic(txnAppCtx, &serviceItem, etag.GenerateEtag(serviceItem.UpdatedAt))
		if err != nil {
			return err
		}

		sitAddressUpdate.MTOServiceItem = *updatedServiceItem

		return nil
	})

	if txErr != nil {
		return nil, txErr
	}

	return sitAddressUpdate, nil
}
