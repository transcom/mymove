package ghcrateengine

import (
	"fmt"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/unit"
)

type counselingServicesPricer struct {
}

// NewCounselingServicesPricer creates a new pricer for counseling services
func NewCounselingServicesPricer() services.CounselingServicesPricer {
	return &counselingServicesPricer{}
}

// Price determines the price for a counseling service
func (p counselingServicesPricer) Price(appCtx appcontext.AppContext, lockedPriceCents *unit.Cents) (unit.Cents, services.PricingDisplayParams, error) {

	if lockedPriceCents == nil {
		return 0, nil, fmt.Errorf("invalid value for locked_price_cents")
	}

	params := services.PricingDisplayParams{
		{
			Key:   models.ServiceItemParamNamePriceRateOrFactor,
			Value: FormatCents(*lockedPriceCents),
		},
	}

	return *lockedPriceCents, params, nil
}

// PriceUsingParams determines the price for a counseling service given PaymentServiceItemParams
func (p counselingServicesPricer) PriceUsingParams(appCtx appcontext.AppContext, params models.PaymentServiceItemParams, featureFlagValues map[string]bool) (unit.Cents, services.PricingDisplayParams, error) {

	lockedPriceCents, err := getParamInt(params, models.ServiceItemParamNameLockedPriceCents)
	if err != nil {
		return unit.Cents(0), nil, err
	}

	lockedPrice := unit.Cents(lockedPriceCents)
	return p.Price(appCtx, &lockedPrice)
}
