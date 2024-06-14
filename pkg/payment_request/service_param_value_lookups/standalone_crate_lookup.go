package serviceparamvaluelookups

import (
	"strconv"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
)

// StandaloneCrateLookup does lookup on standaloneCrate
type StandaloneCrateLookup struct {
	ServiceItem models.MTOServiceItem
}

func (r StandaloneCrateLookup) lookup(_ appcontext.AppContext, _ *ServiceItemParamKeyData) (string, error) {
	standaloneCrate := r.ServiceItem.StandaloneCrate
	if standaloneCrate == nil {
		return "false", nil
	}

	return strconv.FormatBool(*standaloneCrate), nil
}
