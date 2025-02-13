package serviceparamvaluelookups

import (
	"strconv"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
)

// ExternalCrateLookup does lookup on externalCrate
type ExternalCrateLookup struct {
	ServiceItem models.MTOServiceItem
}

func (r ExternalCrateLookup) lookup(_ appcontext.AppContext, _ *ServiceItemParamKeyData) (string, error) {
	externalCrate := r.ServiceItem.ExternalCrate
	if externalCrate == nil {
		return "false", nil
	}

	return strconv.FormatBool(*externalCrate), nil
}
