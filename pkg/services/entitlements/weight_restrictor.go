package entitlements

import (
	"fmt"

	"github.com/pkg/errors"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/etag"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

type weightRestrictor struct {
}

func NewWeightRestrictor() services.WeightRestrictor {
	return &weightRestrictor{}
}

func (wr *weightRestrictor) ApplyWeightRestrictionToEntitlement(appCtx appcontext.AppContext, entitlement models.Entitlement, weightRestriction int, eTag string) (*models.Entitlement, error) {
	// First, fetch the latest version of the entitlement for etag check
	var originalEntitlement models.Entitlement
	err := appCtx.DB().Find(&originalEntitlement, entitlement.ID)
	if err != nil {
		return nil, apperror.NewQueryError("Entitlements", err, "error fetching entitlement")
	}

	// verify ETag
	if etag.GenerateEtag(originalEntitlement.UpdatedAt) != eTag {
		return nil, apperror.NewPreconditionFailedError(originalEntitlement.ID, nil)
	}

	maxHhgAllowance, err := wr.fetchMaxHhgAllowance(appCtx)
	if err != nil {
		return nil, err
	}

	// Don't allow applying a weight restriction above teh max allowance, that's silly
	if weightRestriction > maxHhgAllowance {
		return nil, apperror.NewInvalidInputError(entitlement.ID,
			fmt.Errorf("weight restriction %d exceeds max HHG allowance %d", weightRestriction, maxHhgAllowance),
			nil, "error applying weight restriction")
	}

	// Update the restriction fields
	originalEntitlement.IsWeightRestricted = true
	originalEntitlement.WeightRestriction = &weightRestriction

	verrs, err := appCtx.DB().ValidateAndUpdate(&originalEntitlement)
	if err != nil {
		return nil, apperror.NewQueryError("Entitlements", err, "error updating weight restriction for entitlement")
	}
	if verrs != nil && verrs.HasAny() {
		return nil, apperror.NewInvalidInputError(originalEntitlement.ID, err, verrs, "invalid input while updating entitlement")
	}

	return &originalEntitlement, nil
}

func (wr *weightRestrictor) RemoveWeightRestrictionFromEntitlement(appCtx appcontext.AppContext, entitlement models.Entitlement, eTag string) (*models.Entitlement, error) {
	// Fetch the latest version of the entitlement for etag check
	var originalEntitlement models.Entitlement
	err := appCtx.DB().Find(&originalEntitlement, entitlement.ID)
	if err != nil {
		if errors.Cause(err).Error() == models.RecordNotFoundErrorString {
			return nil, apperror.NewNotFoundError(entitlement.ID, "entitlement not found")
		}
		return nil, apperror.NewQueryError("Entitlements", err, "error fetching entitlement")
	}

	// verify ETag
	if etag.GenerateEtag(originalEntitlement.UpdatedAt) != eTag {
		return nil, apperror.NewPreconditionFailedError(originalEntitlement.ID, nil)
	}

	// Update the restriction fields
	originalEntitlement.IsWeightRestricted = false
	originalEntitlement.WeightRestriction = nil

	verrs, err := appCtx.DB().ValidateAndUpdate(&originalEntitlement)
	if err != nil {
		return nil, apperror.NewQueryError("Entitlements", err, "error removing weight restriction for entitlement")
	}
	if verrs != nil && verrs.HasAny() {
		return nil, apperror.NewInvalidInputError(originalEntitlement.ID, err, verrs, "invalid input while updating entitlement")
	}

	return &originalEntitlement, nil
}

func (wr *weightRestrictor) fetchMaxHhgAllowance(appCtx appcontext.AppContext) (int, error) {
	var maxHhgAllowance int
	err := appCtx.DB().
		RawQuery(`
            SELECT parameter_value::int
            FROM application_parameters
            WHERE parameter_name = $1
            LIMIT 1
        `, "maxHhgAllowance").
		First(&maxHhgAllowance)

	if err != nil {
		return maxHhgAllowance, apperror.NewQueryError("ApplicationParameters", err, "error fetching max HHG allowance")
	}
	return maxHhgAllowance, nil
}
