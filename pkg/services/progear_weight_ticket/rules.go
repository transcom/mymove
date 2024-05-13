package progearweightticket

import (
	"github.com/gobuffalo/validate/v3"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
)

func checkID() progearWeightTicketValidator {
	return progearWeightTicketValidatorFunc(func(_ appcontext.AppContext, newProgearWeightTicket *models.ProgearWeightTicket, oldProgearWeightTicket *models.ProgearWeightTicket) error {
		verrs := validate.NewErrors()

		if newProgearWeightTicket == nil || oldProgearWeightTicket == nil {
			return verrs
		}

		if newProgearWeightTicket.ID != oldProgearWeightTicket.ID {
			verrs.Add("ID", "new ProgearWeightTicket ID must match the old ProgearWeightTicket ID")
		}

		return verrs
	})
}

func checkBaseRequiredFields() progearWeightTicketValidator {
	return progearWeightTicketValidatorFunc(func(_ appcontext.AppContext, newProgearWeightTicket *models.ProgearWeightTicket, _ *models.ProgearWeightTicket) error {
		verrs := validate.NewErrors()

		if newProgearWeightTicket == nil {
			return verrs
		}

		if newProgearWeightTicket.PPMShipmentID.IsNil() {
			verrs.Add("PPMShipmentID", "PPMShipmentID must exist")
		}

		if newProgearWeightTicket.Document.ServiceMemberID.IsNil() {
			verrs.Add("ServiceMemberID", "Document ServiceMemberID must exist")
		}

		return verrs
	})
}

func checkAdditionalRequiredFields() progearWeightTicketValidator {
	return progearWeightTicketValidatorFunc(func(_ appcontext.AppContext, newProgearWeightTicket *models.ProgearWeightTicket, oldProgearWeightTicket *models.ProgearWeightTicket) error {
		verrs := validate.NewErrors()

		if newProgearWeightTicket.BelongsToSelf == nil {
			verrs.Add("BelongsToSelf", "BelongsToSelf must be a boolean value")
		}

		if newProgearWeightTicket.Description == nil || *newProgearWeightTicket.Description == "" {
			verrs.Add("Description", "Description must have a value of at least 0")
		}

		if newProgearWeightTicket.HasWeightTickets == nil {
			verrs.Add("HasWeightTickets", "HasWeightTickets is required")
		}

		if newProgearWeightTicket.Weight == nil || *newProgearWeightTicket.Weight < 1 {
			verrs.Add("Weight", "Weight must have a value of at least 1")
		}

		if len(oldProgearWeightTicket.Document.UserUploads) < 1 {
			verrs.Add("Document", "At least 1 weight ticket is required")
		}

		return verrs
	})
}

// verifyReasonAndStatusAreConstant ensures that the reason and status fields do not change
func verifyReasonAndStatusAreConstant() progearWeightTicketValidator {
	return progearWeightTicketValidatorFunc(func(_ appcontext.AppContext, newProgearWeightTicket *models.ProgearWeightTicket, originalProgearWeightTicket *models.ProgearWeightTicket) error {
		verrs := validate.NewErrors()

		if (originalProgearWeightTicket.Status == nil && newProgearWeightTicket.Status != nil) ||
			(originalProgearWeightTicket.Status != nil && newProgearWeightTicket.Status == nil) ||
			(originalProgearWeightTicket.Status != nil && newProgearWeightTicket.Status != nil && *originalProgearWeightTicket.Status != *newProgearWeightTicket.Status) {
			verrs.Add("Status", "status cannot be modified")
		}

		if (originalProgearWeightTicket.Reason == nil && newProgearWeightTicket.Reason != nil) ||
			(originalProgearWeightTicket.Reason != nil && newProgearWeightTicket.Reason == nil) ||
			(originalProgearWeightTicket.Reason != nil && newProgearWeightTicket.Reason != nil && *originalProgearWeightTicket.Reason != *newProgearWeightTicket.Reason) {
			verrs.Add("Reason", "reason cannot be modified")
		}

		return verrs
	})
}

func verifyReasonAndStatusAreValid() progearWeightTicketValidator {
	return progearWeightTicketValidatorFunc(func(_ appcontext.AppContext, newProgearWeightTicket *models.ProgearWeightTicket, _ *models.ProgearWeightTicket) error {
		verrs := validate.NewErrors()

		if newProgearWeightTicket.Status != nil {
			if *newProgearWeightTicket.Status == models.PPMDocumentStatusApproved && newProgearWeightTicket.Reason != nil {
				verrs.Add("Reason", "reason must not be set if the status is Approved")
			} else if (*newProgearWeightTicket.Status == models.PPMDocumentStatusExcluded || *newProgearWeightTicket.Status == models.PPMDocumentStatusRejected) &&
				(newProgearWeightTicket.Reason == nil || *newProgearWeightTicket.Reason == "") {
				verrs.Add("Reason", "reason is mandatory if the status is Excluded or Rejected")
			}
		} else if newProgearWeightTicket.Reason != nil {
			verrs.Add("Reason", "reason should not be set if the status is not set")
		}

		return verrs
	})
}

func basicChecksForCreate() []progearWeightTicketValidator {
	return []progearWeightTicketValidator{
		checkID(),
		checkBaseRequiredFields(),
	}
}

func basicChecksForCustomer() []progearWeightTicketValidator {
	return []progearWeightTicketValidator{
		checkID(),
		checkBaseRequiredFields(),
		checkAdditionalRequiredFields(),
		verifyReasonAndStatusAreConstant(),
	}
}

func basicChecksForOffice() []progearWeightTicketValidator {
	return []progearWeightTicketValidator{
		checkID(),
		checkBaseRequiredFields(),
		checkAdditionalRequiredFields(),
		verifyReasonAndStatusAreValid(),
	}
}
