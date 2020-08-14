package serviceparamvaluelookups

import (
	"fmt"

	"github.com/transcom/mymove/pkg/models"
)

// ZipAddressLookup does lookup on the postal code for the pickup address
type ZipAddressLookup struct {
	Address models.Address
}

func (r ZipAddressLookup) lookup(keyData *ServiceItemParamKeyData) (string, error) {
	value := fmt.Sprintf("%s", r.Address.PostalCode)
	return value, nil
}
