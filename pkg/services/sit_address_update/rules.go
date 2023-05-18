package sitaddressupdate

import (
	"github.com/gobuffalo/validate/v3"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
)

// checkRequiredFields checks that the required fields are included
func checkRequiredFields() sitAddressUpdateValidator {
	return sitAddressUpdateValidatorFunc(func(appCtx appcontext.AppContext, sitAddressUpdate *models.SITAddressUpdate) error {
		verrs := validate.NewErrors()

		// Distance and Status are required fields but aren't validated here
		// Distance should be calcualted
		// Status should be updated with using approve/reject service objects

		if !sitAddressUpdate.NewAddressID.IsNil() {
			verrs.Add("NewAddress:id", "NewAddress:id cannot be set for new addresses")
		}
		if sitAddressUpdate.NewAddress.PostalCode == "" {
			verrs.Add("NewAddress", "NewAddress is required")
		}
		if sitAddressUpdate.MTOServiceItemID.IsNil() {
			verrs.Add("serviceItem", "MTOServiceItem is required")
		}

		var serviceItem models.MTOServiceItem
		err := appCtx.DB().Where("id = ?", sitAddressUpdate.MTOServiceItemID).First(&serviceItem)
		if err != nil {
			verrs.Add("MTOServiceItem", "MTOServiceItem was not found")
		}

		if serviceItem.Status != models.MTOServiceItemStatusApproved {
			verrs.Add("MTOServiceItemID", "MTOServiceItem must be approved")
		}

		if serviceItem.SITDestinationFinalAddressID == nil || serviceItem.SITDestinationFinalAddressID.IsNil() {
			verrs.Add("SITDestinationFinalAddressID", "SITDestinationFinalAddressID is required")
		}

		return verrs
	})
}

func checkTOORequiredFields() sitAddressUpdateValidator {
	return sitAddressUpdateValidatorFunc(func(_ appcontext.AppContext, sitAddressUpdate *models.SITAddressUpdate) error {
		verrs := validate.NewErrors()

		if sitAddressUpdate.OfficeRemarks == nil {
			verrs.Add("OfficeRemarks", "OfficeRemarks are required")
		}

		return verrs
	})
}

func checkPrimeRequiredFields() sitAddressUpdateValidator {
	return sitAddressUpdateValidatorFunc(func(_ appcontext.AppContext, sitAddressUpdate *models.SITAddressUpdate) error {
		verrs := validate.NewErrors()

		if sitAddressUpdate.ContractorRemarks == nil {
			verrs.Add("ContractorRemarks", "ContractorRemarks are required")
		}

		return verrs
	})
}
