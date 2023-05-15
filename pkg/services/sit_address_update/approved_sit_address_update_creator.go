package sitaddressupdate

import (
	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/etag"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/route"
	"github.com/transcom/mymove/pkg/services"
)

// approvedSITAddressUpdateCreator is the concrete struct implementing the services.ApprovedSITAddressUpdateCreator interface
type approvedSITAddressUpdateCreator struct {
	planner            route.Planner
	addressCreator     services.AddressCreator
	serviceItemUpdater services.MTOServiceItemUpdater
	checks             []sitAddressUpdateValidator
}

// NewApprovedOfficeSITAddressUpdateCreator creates a new struct with the service dependencies
func NewApprovedOfficeSITAddressUpdateCreator(planner route.Planner, addressCreator services.AddressCreator, serviceItemUpdater services.MTOServiceItemUpdater) services.ApprovedSITAddressUpdateCreator {
	return &approvedSITAddressUpdateCreator{
		planner:            planner,
		addressCreator:     addressCreator,
		serviceItemUpdater: serviceItemUpdater,
		checks: []sitAddressUpdateValidator{
			checkRequiredFields(),
			checkTOORequiredFields(),
		},
	}
}

// CreateSITAddressUpdate creates a SIT Address Update
func (f *approvedSITAddressUpdateCreator) CreateApprovedSITAddressUpdate(appCtx appcontext.AppContext, sitAddressUpdate *models.SITAddressUpdate) (*models.SITAddressUpdate, error) {
	var err error
	if err = validateSITAddressUpdate(appCtx, sitAddressUpdate, f.checks...); err != nil {
		return nil, err
	}

	sitAddressUpdate.Status = models.SITAddressUpdateStatusApproved
	sitAddressUpdate.Distance, err = f.planner.TransitDistance(appCtx, &sitAddressUpdate.OldAddress, &sitAddressUpdate.NewAddress)
	if err != nil {
		return nil, err
	}

	txErr := appCtx.NewTransaction(func(txnAppCtx appcontext.AppContext) (err error) {
		newAddress, err := f.addressCreator.CreateAddress(txnAppCtx, &sitAddressUpdate.NewAddress)
		if err != nil {
			return err
		}
		sitAddressUpdate.NewAddressID = newAddress.ID
		sitAddressUpdate.NewAddress = *newAddress

		verrs, err := txnAppCtx.DB().ValidateAndCreate(sitAddressUpdate)

		if verrs.HasAny() {
			return apperror.NewInvalidInputError(sitAddressUpdate.ID, nil, verrs, "Invalid input found while creating the SIT Address Update.")
		} else if err != nil {
			return apperror.NewQueryError("SITAddressUpdate", err, "Unable to create SIT Address Update")
		}

		var oldServiceItem models.MTOServiceItem
		err = txnAppCtx.DB().Where("id = ?", sitAddressUpdate.MTOServiceItemID).First(&oldServiceItem)
		if err != nil {
			return err
		}

		oldServiceItem.SITDestinationFinalAddressID = &newAddress.ID
		oldServiceItem.SITDestinationFinalAddress = newAddress
		updatedServiceItem, err := f.serviceItemUpdater.UpdateMTOServiceItemBasic(txnAppCtx, &oldServiceItem, etag.GenerateEtag(oldServiceItem.UpdatedAt))
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
