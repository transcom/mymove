package ghcrateengine

import (
	"time"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/unit"
)

type intlCratingPricer struct {
}

// NewIntlCratingPricer creates a new pricer for international crating
func NewIntlCratingPricer() services.IntlCratingPricer {
	return &intlCratingPricer{}
}

// Price determines the price for international crating
func (p intlCratingPricer) Price(appCtx appcontext.AppContext, contractCode string, referenceDate time.Time, billedCubicFeet unit.CubicFeet, standaloneCrate bool, standaloneCrateCap unit.Cents, externalCrate bool, market models.Market) (unit.Cents, services.PricingDisplayParams, error) {
	return priceIntlCratingUncrating(appCtx, models.ReServiceCodeICRT, contractCode, referenceDate, billedCubicFeet, standaloneCrate, standaloneCrateCap, externalCrate, market)
}

// PriceUsingParams determines the price for international crating given PaymentServiceItemParams
func (p intlCratingPricer) PriceUsingParams(appCtx appcontext.AppContext, params models.PaymentServiceItemParams) (unit.Cents, services.PricingDisplayParams, error) {
	contractCode, err := getParamString(params, models.ServiceItemParamNameContractCode)
	if err != nil {
		return unit.Cents(0), nil, err
	}

	cubicFeetFloat, err := getParamFloat(params, models.ServiceItemParamNameCubicFeetBilled)
	if err != nil {
		return unit.Cents(0), nil, err
	}

	cubicFeetBilled := unit.CubicFeet(cubicFeetFloat)

	referenceDate, err := getParamTime(params, models.ServiceItemParamNameReferenceDate)
	if err != nil {
		return unit.Cents(0), nil, err
	}

	market, err := getParamMarket(params, models.ServiceItemParamNameMarketOrigin)
	if err != nil {
		return unit.Cents(0), nil, err
	}

	standaloneCrate, err := getParamBool(params, models.ServiceItemParamNameStandaloneCrate)
	if err != nil {
		return unit.Cents(0), nil, err
	}

	externalCrate, err := getParamBool(params, models.ServiceItemParamNameExternalCrate)
	if err != nil {
		return unit.Cents(0), nil, err
	}

	standaloneCrateCapParam, err := getParamInt(params, models.ServiceItemParamNameStandaloneCrateCap)
	if err != nil {
		return unit.Cents(0), nil, err
	}
	standaloneCrateCap := unit.Cents(float64(standaloneCrateCapParam))

	return p.Price(appCtx, contractCode, referenceDate, cubicFeetBilled, standaloneCrate, standaloneCrateCap, externalCrate, market)
}
