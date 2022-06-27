package weightticket

import (
	"github.com/gobuffalo/validate/v3"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
)

func checkID() weightTicketValidator {
	return weightTicketValidatorFunc(func(_ appcontext.AppContext, newWeightTicket *models.WeightTicket, originalWeightTicket *models.WeightTicket) error {
		verrs := validate.NewErrors()

		if newWeightTicket == nil && originalWeightTicket == nil {
			return verrs
		}

		if newWeightTicket != nil && originalWeightTicket != nil {
			if newWeightTicket.ID != originalWeightTicket.ID {
				verrs.Add("ID", "new WeightTicket ID must match original WeightTicket ID")
			}
		}

		return verrs
	})
}

func basicChecks() []weightTicketValidator {
	return []weightTicketValidator{
		checkID(),
	}
}
