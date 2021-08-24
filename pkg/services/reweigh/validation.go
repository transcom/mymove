package reweigh

import (
	"context"

	"github.com/transcom/mymove/pkg/models"
)

// checks business requirements
// reweighValidator describes a method for checking business requirements
type reweighValidator interface {
	// Validate checks the newReweigh for adherence to business rules. The
	// oldReweigh parameter is expected to be nil in creation use cases.
	// It is safe to return a *validate.Errors with zero added errors as
	// a success case.
	Validate(c context.Context, newReweigh models.Reweigh, oldReweigh *models.Reweigh, shipment *models.MTOShipment) error
}

// reweighValidatorFunc is an adapter type for converting a function into an implementation of reweighValidator
type reweighValidatorFunc func(context.Context, models.Reweigh, *models.Reweigh, *models.MTOShipment) error

// Validate fulfills the reweighValidator interface
func (fn reweighValidatorFunc) Validate(ctx context.Context, newer models.Reweigh, older *models.Reweigh, ship *models.MTOShipment) error {
	return fn(ctx, newer, older, ship)
}
