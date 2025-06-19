package paymentrequest

import (
	"errors"
	"fmt"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
)

// resolveReServiceForLookup ensures that the correct ReService is used for parameter lookup
// This is because some service items don't have parameters that can be looked up because they inherit the logic from existing items.
// For example, INPK. INPK is for iHHG shipments going into non-temporary storage.
// This means we are packing an iHHG shipment, so we price by IHPK, but with a special
// pricer for INPK. INPK is iHHG -> iNTS. Prices by IHPK multiplied by NTS market factor
func resolveReServiceForLookup(appCtx appcontext.AppContext, mtoServiceItem models.MTOServiceItem) (models.ReService, error) {
	var reService models.ReService

	// Map of swap services - aka codes that are priced using a different code's parameters.
	// It maps the actual MTOServiceItem's ReService.Code to the ReService.Code that we should look up parameters for
	serviceCodeSwaps := map[models.ReServiceCode]models.ReServiceCode{
		models.ReServiceCodeINPK: models.ReServiceCodeIHPK,
	}

	requestedCode := mtoServiceItem.ReService.Code
	if requestedCode == "" {
		return reService, errors.New("Error when resolving ReServiceItemForLookup: mtoServiceItem does not have a joined ReService with a code")
	}

	// If there’s a swap, lookup and return the alternate code. Otherwise, return as is
	if alternateCode, ok := serviceCodeSwaps[requestedCode]; ok {
		reServicePtr, err := models.FetchReServiceByCode(appCtx.DB(), alternateCode)
		if err != nil {
			return models.ReService{}, fmt.Errorf("failed to fetch ReService by code %s: %w", alternateCode, err)
		}
		return *reServicePtr, nil
	}

	// No swap needed
	return mtoServiceItem.ReService, nil
}
