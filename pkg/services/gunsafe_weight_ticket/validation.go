package gunsafeweightticket

import (
	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/models"
)

// gunSafeWeightTicketValidator defines the interface for checking business rules for a gunSafeWeightTicket
type gunSafeWeightTicketValidator interface {
	// The newGunSafeWeightTicket is optional, as create requires no payload.
	// The originalGunSafeWeightTicket is optional, so it's a pointer type.
	Validate(appCtx appcontext.AppContext, newGunSafeWeightTicket *models.GunSafeWeightTicket, originalGunSafeWeightTicket *models.GunSafeWeightTicket) error
}

// gunSafeWeightTicketValidatorFunc is an adapter that will convert a function into an implementation of gunSafeWeightTicketValidator
type gunSafeWeightTicketValidatorFunc func(appcontext.AppContext, *models.GunSafeWeightTicket, *models.GunSafeWeightTicket) error

// Validate fulfills the gunSafeWeightTicketValidator interface
func (fn gunSafeWeightTicketValidatorFunc) Validate(appCtx appcontext.AppContext, newGunSafeWeightTicket *models.GunSafeWeightTicket, originalGunSafeWeightTicket *models.GunSafeWeightTicket) error {
	return fn(appCtx, newGunSafeWeightTicket, originalGunSafeWeightTicket)
}

func validateGunSafeWeightTicket(appCtx appcontext.AppContext, newGunSafeWeightTicket *models.GunSafeWeightTicket, originalGunSafeWeightTicket *models.GunSafeWeightTicket, checks ...gunSafeWeightTicketValidator) error {
	verrs := validate.NewErrors()

	for _, check := range checks {
		if err := check.Validate(appCtx, newGunSafeWeightTicket, originalGunSafeWeightTicket); err != nil {
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
		if newGunSafeWeightTicket != nil {
			currentID = newGunSafeWeightTicket.ID
		}
		return apperror.NewInvalidInputError(currentID, nil, verrs, "Invalid input found while validating the gunSafe weight ticket.")
	}

	return nil
}
