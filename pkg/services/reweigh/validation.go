package reweigh

import (
	"github.com/transcom/mymove/pkg/appcontext"

	"github.com/transcom/mymove/pkg/models"
)

// checks business requirements
// reweighValidator describes a method for checking business requirements
type reweighValidator interface {
	// Validate checks the newReweigh for adherence to business rules. The
	// oldReweigh parameter is expected to be nil in creation use cases.
	// It is safe to return a *validate.Errors with zero added errors as
	// a success case.
	Validate(a appcontext.AppContext, newReweigh models.Reweigh, oldReweigh *models.Reweigh, shipment *models.MTOShipment) error
}

// reweighValidatorFunc is an adapter type for converting a function into an implementation of reweighValidator
type reweighValidatorFunc func(appcontext.AppContext, models.Reweigh, *models.Reweigh, *models.MTOShipment) error

// Validate fulfills the reweighValidator interface
func (fn reweighValidatorFunc) Validate(appCtx appcontext.AppContext, newer models.Reweigh, older *models.Reweigh, ship *models.MTOShipment) error {
	return fn(appCtx, newer, older, ship)
}
