package ghcrateengine

import (
	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/unit"
)

type managementServicesPricer struct {
}

// NewManagementServicesPricer creates a new pricer for management services
func NewManagementServicesPricer() services.ManagementServicesPricer {
	return &managementServicesPricer{}
}

// Price determines the price for a management service
func (p managementServicesPricer) Price(appCtx appcontext.AppContext, lockedPriceCents unit.Cents) (unit.Cents, services.PricingDisplayParams, error) {

	params := services.PricingDisplayParams{
		{
			Key:   models.ServiceItemParamNamePriceRateOrFactor,
			Value: FormatCents(lockedPriceCents),
		},
	}

	return lockedPriceCents, params, nil
}

// PriceUsingParams determines the price for a management service given PaymentServiceItemParams
func (p managementServicesPricer) PriceUsingParams(appCtx appcontext.AppContext, params models.PaymentServiceItemParams) (unit.Cents, services.PricingDisplayParams, error) {

	lockedPriceCents, err := getParamInt(params, models.ServiceItemParamNameLockedPriceCents)
	if err != nil {
		return unit.Cents(0), nil, err
	}

	return p.Price(appCtx, unit.Cents(lockedPriceCents))
}
