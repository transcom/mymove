package ghcrateengine

import (
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
func (p counselingServicesPricer) Price(appCtx appcontext.AppContext, lockedPriceCents unit.Cents) (unit.Cents, services.PricingDisplayParams, error) {

	params := services.PricingDisplayParams{
		{
			Key:   models.ServiceItemParamNamePriceRateOrFactor,
			Value: FormatCents(lockedPriceCents),
		},
	}

	return lockedPriceCents, params, nil
}

// PriceUsingParams determines the price for a counseling service given PaymentServiceItemParams
func (p counselingServicesPricer) PriceUsingParams(appCtx appcontext.AppContext, params models.PaymentServiceItemParams) (unit.Cents, services.PricingDisplayParams, error) {

	lockedPriceCents, err := getParamInt(params, models.ServiceItemParamNameLockedPriceCents)
	if err != nil {
		return unit.Cents(0), nil, err
	}

	return p.Price(appCtx, unit.Cents(lockedPriceCents))
}
