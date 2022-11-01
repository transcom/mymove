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

func checkCreateRequiredFields() progearWeightTicketValidator {
	return progearWeightTicketValidatorFunc(func(_ appcontext.AppContext, newProgearWeightTicket *models.ProgearWeightTicket, oldProgearWeightTicket *models.ProgearWeightTicket) error {
		verrs := validate.NewErrors()

		if newProgearWeightTicket.PPMShipmentID.IsNil() {
			verrs.Add("PPMShipmentID", "PPMShipmentID must exist")
		}

		if newProgearWeightTicket.FullDocument.ServiceMemberID.IsNil() {
			verrs.Add("ServiceMemberID", "Full Document ServiceMemberID must exist")
		}

		return verrs
	})
}

func checkUpdateRequiredFields() progearWeightTicketValidator {
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
		if newProgearWeightTicket.EmptyWeight == nil || *newProgearWeightTicket.EmptyWeight < 0 {
			verrs.Add("EmptyWeight", "EmptyWeight must have a value of at least 0")
		}

		if newProgearWeightTicket.FullWeight == nil || *newProgearWeightTicket.EmptyWeight < 1 {
			verrs.Add("FullWeight", "FullWeight must have a value of at least 1")
		}

		if newProgearWeightTicket.ConstructedWeight == nil || *newProgearWeightTicket.ConstructedWeight < 1 {
			verrs.Add("ConstructedWeight", "ConstructedWeight must have a value of at least 1")
		}

		if len(oldProgearWeightTicket.FullDocument.UserUploads) < 1 {
			verrs.Add("FullDocument", "At least 1 weight ticket is required")
		}

		if newProgearWeightTicket.Status != nil {
			if (*newProgearWeightTicket.Status == models.PPMDocumentStatusExcluded || *newProgearWeightTicket.Status == models.PPMDocumentStatusRejected) && (newProgearWeightTicket.Reason == nil || *newProgearWeightTicket.Reason == "") {
				verrs.Add("Reason", "A reason must be provided when the status is EXCLUDED or REJECTED")
			}
		}

		return verrs
	})
}

func createChecks() []progearWeightTicketValidator {
	return []progearWeightTicketValidator{
		checkID(),
		checkCreateRequiredFields(),
	}
}

func updateChecks() []progearWeightTicketValidator {
	return []progearWeightTicketValidator{
		checkID(),
		checkCreateRequiredFields(),
		checkUpdateRequiredFields(),
	}
}
