package serviceparamvaluelookups

import (
	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
)

// ZipAddressLookup does lookup on the postal code for the pickup address
type ZipAddressLookup struct {
	Address models.Address
}

func (r ZipAddressLookup) lookup(appCtx appcontext.AppContext, keyData *ServiceItemParamKeyData) (string, error) {
	value := r.Address.PostalCode
	return value, nil
}
