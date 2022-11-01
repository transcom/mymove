package address

import (
	"github.com/gobuffalo/validate/v3"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
)

// checkID ensures that an ID is not passed if there is no original address and if there is an original address it ensures the id doesn't change
func checkID() addressValidator {
	return addressValidatorFunc(func(_ appcontext.AppContext, newAddress, originalAddress *models.Address) error {
		verrs := validate.NewErrors()
		if originalAddress == nil {
			if !newAddress.ID.IsNil() {
				verrs.Add("ID", "an ID should not be specified")
			}
		} else if originalAddress.ID != newAddress.ID {
			verrs.Add("ID", "an ID should not change")
		}
		return verrs
	})
}
