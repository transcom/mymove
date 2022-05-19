package shipment

import (
	"github.com/gobuffalo/validate/v3"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/models"
)

// shipmentValidator defines the interface for checking business rules for a shipment
type shipmentValidator interface {
	Validate(appCtx appcontext.AppContext, newShipment models.MTOShipment) error
}

// shipmentValidatorFunc is an adapter that will convert a function into an implementation of shipmentValidator
type shipmentValidatorFunc func(appcontext.AppContext, models.MTOShipment) error

// Validate fulfills the shipmentValidator interface
func (fn shipmentValidatorFunc) Validate(appCtx appcontext.AppContext, newShipment models.MTOShipment) error {
	return fn(appCtx, newShipment)
}

// validateShipment runs a shipment through the checks that are passed in.
func validateShipment(appCtx appcontext.AppContext, newShipment models.MTOShipment, checks ...shipmentValidator) error {
	verrs := validate.NewErrors()

	for _, check := range checks {
		if err := check.Validate(appCtx, newShipment); err != nil {
			switch e := err.(type) {
			case *validate.Errors:
				verrs.Append(e)
			default:
				return err
			}
		}
	}

	if verrs.HasAny() {
		return apperror.NewInvalidInputError(newShipment.ID, nil, verrs, "Invalid input found while validating the shipment.")
	}

	return nil
}
