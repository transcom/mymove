package entitlements

import (
	"fmt"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/services"
)

type weightRestrictor struct {
}

func NewWeightRestrictor() services.WeightRestrictor {
	return &weightRestrictor{}
}

func (wr *weightRestrictor) ApplyWeightRestrictionToEntitlement(appCtx appcontext.AppContext, entitlementID uuid.UUID, weightRestriction int) error {

	maxHhgAllowance, err := wr.fetchMaxHhgAllowance(appCtx)
	if err != nil {
		return err
	}

	// Don't allow applying a weight restriction above teh max allowance, that's silly
	if weightRestriction > maxHhgAllowance {
		return apperror.NewInvalidInputError(entitlementID, fmt.Errorf("weight restriction %d exceeds max HHG allowance %d", weightRestriction, maxHhgAllowance), nil, "error applying weight restriction")
	}

	// If we reached this spot we're good to apply the restriction to the entitlement
	err = appCtx.DB().
		RawQuery(`
            UPDATE entitlements
            SET weight_restriction = $1, is_weight_restricted = true
            WHERE id = $2
        `, weightRestriction, entitlementID).
		Exec()
	if err != nil {
		return apperror.NewQueryError("Entitlements", err, "error updating weight restriction for entitlement")
	}

	return nil
}

func (wr *weightRestrictor) RemoveWeightRestrictionFromEntitlement(appCtx appcontext.AppContext, entitlementID uuid.UUID) error {
	// Remove the restriction by setting weight_restriction = NULL and is_weight_restricted = false
	err := appCtx.DB().
		RawQuery(`
            UPDATE entitlements
            SET weight_restriction = NULL, is_weight_restricted = false
            WHERE id = $1
        `, entitlementID).
		Exec()
	if err != nil {
		return apperror.NewQueryError("Entitlements", err, "error removing weight restriction for entitlement")
	}

	return nil
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
