package serviceparamvaluelookups

import (
	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
)

type MarketDestinationLookup struct {
	Address models.Address
}

func (r MarketDestinationLookup) lookup(_ appcontext.AppContext, _ *ServiceItemParamKeyData) (string, error) {
	international := r.Address.IsOconus
	value := handlers.FmtString(models.MarketOconus.String())
	if *international {
		value = handlers.FmtString(models.MarketConus.String())
	}
	return *value, nil
}
