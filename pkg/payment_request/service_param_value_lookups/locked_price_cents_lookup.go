package serviceparamvaluelookups

import (
	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/models"
)

// LockedPriceCents does lookup on serviceItem
type LockedPriceCentsLookup struct {
	ServiceItem models.MTOServiceItem
}

func (r LockedPriceCentsLookup) lookup(appCtx appcontext.AppContext, _ *ServiceItemParamKeyData) (string, error) {
	lockedPriceCents := r.ServiceItem.LockedPriceCents
	if lockedPriceCents == nil {
		return "0", apperror.NewConflictError(r.ServiceItem.ID, "unable to find locked price cents")
	}

	return lockedPriceCents.ToMillicents().ToCents().String(), nil
}
