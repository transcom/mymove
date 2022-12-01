package address

import (
	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/models"
)

type addressValidator interface {
	Validate(appCtx appcontext.AppContext, newAddress *models.Address, originalAddress *models.Address) error
}

type addressValidatorFunc func(appcontext.AppContext, *models.Address, *models.Address) error

func (fn addressValidatorFunc) Validate(appCtx appcontext.AppContext, newAddress *models.Address, originalAddress *models.Address) error {
	return fn(appCtx, newAddress, originalAddress)
}

func validateAddress(appCtx appcontext.AppContext, newAddress *models.Address, originalAddress *models.Address, checks ...addressValidator) error {
	verrs := validate.NewErrors()

	for _, check := range checks {
		if err := check.Validate(appCtx, newAddress, originalAddress); err != nil {
			switch e := err.(type) {
			case *validate.Errors:
				// Accumulate all validation errors
				verrs.Append(e)
			default:
				// Non-validation errors have priority and short-circuit doing any further checks
				return err
			}
		}
	}

	if verrs.HasAny() {
		currentID := uuid.Nil
		if newAddress != nil {
			currentID = newAddress.ID
		}
		return apperror.NewInvalidInputError(currentID, nil, verrs, "Invalid input found while validating the address.")
	}

	return nil
}
