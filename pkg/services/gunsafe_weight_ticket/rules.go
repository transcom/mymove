package gunsafeweightticket

import (
	"github.com/gobuffalo/validate/v3"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
)

func checkID() gunSafeWeightTicketValidator {
	return gunSafeWeightTicketValidatorFunc(func(_ appcontext.AppContext, newGunSafeWeightTicket *models.GunSafeWeightTicket, oldGunSafeWeightTicket *models.GunSafeWeightTicket) error {
		verrs := validate.NewErrors()

		if newGunSafeWeightTicket == nil || oldGunSafeWeightTicket == nil {
			return verrs
		}

		if newGunSafeWeightTicket.ID != oldGunSafeWeightTicket.ID {
			verrs.Add("ID", "new GunSafeWeightTicket ID must match the old GunSafeWeightTicket ID")
		}

		return verrs
	})
}

func checkBaseRequiredFields() gunSafeWeightTicketValidator {
	return gunSafeWeightTicketValidatorFunc(func(_ appcontext.AppContext, newGunSafeWeightTicket *models.GunSafeWeightTicket, _ *models.GunSafeWeightTicket) error {
		verrs := validate.NewErrors()

		if newGunSafeWeightTicket == nil {
			return verrs
		}

		if newGunSafeWeightTicket.PPMShipmentID.IsNil() {
			verrs.Add("PPMShipmentID", "PPMShipmentID must exist")
		}

		if newGunSafeWeightTicket.Document.ServiceMemberID.IsNil() {
			verrs.Add("ServiceMemberID", "Document ServiceMemberID must exist")
		}

		return verrs
	})
}

func checkAdditionalRequiredFields() gunSafeWeightTicketValidator {
	return gunSafeWeightTicketValidatorFunc(func(_ appcontext.AppContext, newGunSafeWeightTicket *models.GunSafeWeightTicket, oldGunSafeWeightTicket *models.GunSafeWeightTicket) error {
		verrs := validate.NewErrors()

		if newGunSafeWeightTicket.Description == nil || *newGunSafeWeightTicket.Description == "" {
			verrs.Add("Description", "Description must have a value")
		}

		if newGunSafeWeightTicket.HasWeightTickets == nil {
			verrs.Add("HasWeightTickets", "HasWeightTickets is required")
		}

		if newGunSafeWeightTicket.Weight == nil || *newGunSafeWeightTicket.Weight < 1 {
			verrs.Add("Weight", "Weight must have a value of at least 1")
		}

		if len(oldGunSafeWeightTicket.Document.UserUploads) < 1 {
			verrs.Add("Document", "At least 1 weight ticket is required")
		}

		return verrs
	})
}

// verifyReasonAndStatusAreConstant ensures that the reason and status fields do not change
func verifyReasonAndStatusAreConstant() gunSafeWeightTicketValidator {
	return gunSafeWeightTicketValidatorFunc(func(_ appcontext.AppContext, newGunSafeWeightTicket *models.GunSafeWeightTicket, originalGunSafeWeightTicket *models.GunSafeWeightTicket) error {
		verrs := validate.NewErrors()

		if (originalGunSafeWeightTicket.Status == nil && newGunSafeWeightTicket.Status != nil) ||
			(originalGunSafeWeightTicket.Status != nil && newGunSafeWeightTicket.Status == nil) ||
			(originalGunSafeWeightTicket.Status != nil && newGunSafeWeightTicket.Status != nil && *originalGunSafeWeightTicket.Status != *newGunSafeWeightTicket.Status) {
			verrs.Add("Status", "status cannot be modified")
		}

		if (originalGunSafeWeightTicket.Reason == nil && newGunSafeWeightTicket.Reason != nil) ||
			(originalGunSafeWeightTicket.Reason != nil && newGunSafeWeightTicket.Reason == nil) ||
			(originalGunSafeWeightTicket.Reason != nil && newGunSafeWeightTicket.Reason != nil && *originalGunSafeWeightTicket.Reason != *newGunSafeWeightTicket.Reason) {
			verrs.Add("Reason", "reason cannot be modified")
		}

		return verrs
	})
}

func verifyReasonAndStatusAreValid() gunSafeWeightTicketValidator {
	return gunSafeWeightTicketValidatorFunc(func(_ appcontext.AppContext, newGunSafeWeightTicket *models.GunSafeWeightTicket, _ *models.GunSafeWeightTicket) error {
		verrs := validate.NewErrors()

		if newGunSafeWeightTicket.Status != nil {
			if *newGunSafeWeightTicket.Status == models.PPMDocumentStatusApproved && newGunSafeWeightTicket.Reason != nil {
				verrs.Add("Reason", "reason must not be set if the status is Approved")
			} else if (*newGunSafeWeightTicket.Status == models.PPMDocumentStatusExcluded || *newGunSafeWeightTicket.Status == models.PPMDocumentStatusRejected) &&
				(newGunSafeWeightTicket.Reason == nil || *newGunSafeWeightTicket.Reason == "") {
				verrs.Add("Reason", "reason is mandatory if the status is Excluded or Rejected")
			}
		} else if newGunSafeWeightTicket.Reason != nil {
			verrs.Add("Reason", "reason should not be set if the status is not set")
		}

		return verrs
	})
}

func basicChecksForCreate() []gunSafeWeightTicketValidator {
	return []gunSafeWeightTicketValidator{
		checkID(),
		checkBaseRequiredFields(),
	}
}

func basicChecksForCustomer() []gunSafeWeightTicketValidator {
	return []gunSafeWeightTicketValidator{
		checkID(),
		checkBaseRequiredFields(),
		checkAdditionalRequiredFields(),
		verifyReasonAndStatusAreConstant(),
	}
}

func basicChecksForOffice() []gunSafeWeightTicketValidator {
	return []gunSafeWeightTicketValidator{
		checkID(),
		checkBaseRequiredFields(),
		checkAdditionalRequiredFields(),
		verifyReasonAndStatusAreValid(),
	}
}
