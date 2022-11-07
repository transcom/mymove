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

func checkRequiredFields() progearWeightTicketValidator {
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

		// ARE WE ALLOWING 0 values now?
		if newProgearWeightTicket.Weight == nil || *newProgearWeightTicket.Weight < 1 {
			verrs.Add("Weight", "Weight must have a value of at least 1")
		}

		if len(oldProgearWeightTicket.Document.UserUploads) < 1 {
			verrs.Add("Document", "At least 1 weight ticket is required")
		}

		if newProgearWeightTicket.Status != nil {
			if (*newProgearWeightTicket.Status == models.PPMDocumentStatusExcluded || *newProgearWeightTicket.Status == models.PPMDocumentStatusRejected) && (newProgearWeightTicket.Reason == nil || *newProgearWeightTicket.Reason == "") {
				verrs.Add("Reason", "A reason must be provided when the status is EXCLUDED or REJECTED")
			}
		}

		return verrs
	})
}

func checkCreateRequiredFields() progearWeightTicketValidator {
	return progearWeightTicketValidatorFunc(func(_ appcontext.AppContext, newProgearWeightTicket *models.ProgearWeightTicket, oldProgearWeightTicket *models.ProgearWeightTicket) error {
		verrs := validate.NewErrors()
		if newProgearWeightTicket == nil || oldProgearWeightTicket == nil {
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

// verifyReasonAndStatusAreConstant ensures that the reason and status fields do not change
func verifyReasonAndStatusAreConstant() progearWeightTicketValidator {
	return progearWeightTicketValidatorFunc(func(_ appcontext.AppContext, newProgearWeightTicket *models.ProgearWeightTicket, originalProgearWeightTicket *models.ProgearWeightTicket) error {
		verrs := validate.NewErrors()

		if newProgearWeightTicket == nil || originalProgearWeightTicket == nil {
			return verrs
		}

		if originalProgearWeightTicket.Status == nil && newProgearWeightTicket.Status != nil {
			verrs.Add("Status", "status cannot be modified")
		} else if originalProgearWeightTicket.Status != nil && newProgearWeightTicket.Status != nil && *originalProgearWeightTicket.Status != *newProgearWeightTicket.Status {
			verrs.Add("Status", "status cannot be modified")
		}

		if originalProgearWeightTicket.Reason != newProgearWeightTicket.Reason {
			verrs.Add("Reason", "reason cannot be modified")
		}
		return verrs
	})
}

func verifyReasonAndStatusAreValid() progearWeightTicketValidator {
	return progearWeightTicketValidatorFunc(func(_ appcontext.AppContext, newProgearWeightTicket *models.ProgearWeightTicket, originalProgearWeightTicket *models.ProgearWeightTicket) error {
		verrs := validate.NewErrors()

		if newProgearWeightTicket == nil || originalProgearWeightTicket == nil {
			return verrs
		}

		if newProgearWeightTicket.Status != nil {
			if *newProgearWeightTicket.Status == models.PPMDocumentStatusApproved && newProgearWeightTicket.Reason != nil {
				verrs.Add("Reason", "reason must be blank if the status is Approved")
			}

			if (*newProgearWeightTicket.Status == models.PPMDocumentStatusExcluded || *newProgearWeightTicket.Status == models.PPMDocumentStatusRejected) && (newProgearWeightTicket.Reason == nil || len(*newProgearWeightTicket.Reason) <= 0) {
				verrs.Add("Reason", "reason is mandatory if the status is Excluded or Rejected")
			}
		} else {
			if newProgearWeightTicket.Reason != nil && len(*newProgearWeightTicket.Reason) > 0 {
				verrs.Add("Reason", "reason should be empty")
			}
		}

		return verrs
	})
}

func basicChecksForCreate() []progearWeightTicketValidator {
	return []progearWeightTicketValidator{
		checkID(),
		checkCreateRequiredFields(),
	}
}

func basicChecksForCustomer() []progearWeightTicketValidator {
	return []progearWeightTicketValidator{
		checkID(),
		checkRequiredFields(),
		verifyReasonAndStatusAreConstant(),
	}
}

func basicChecksForOffice() []progearWeightTicketValidator {
	return []progearWeightTicketValidator{
		checkID(),
		checkRequiredFields(),
		verifyReasonAndStatusAreValid(),
	}
}
