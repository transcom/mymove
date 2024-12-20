package ghcrateengine

import (
	"time"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/unit"
)

type domesticCratingPricer struct {
}

// NewDomesticCratingPricer creates a new pricer for domestic destination first day SIT
func NewDomesticCratingPricer() services.DomesticCratingPricer {
	return &domesticCratingPricer{}
}

// Price determines the price for domestic destination first day SIT
func (p domesticCratingPricer) Price(appCtx appcontext.AppContext, contractCode string, referenceDate time.Time, billedCubicFeet unit.CubicFeet, serviceSchedule int, standaloneCrate bool, standaloneCrateCap unit.Cents) (unit.Cents, services.PricingDisplayParams, error) {
	return priceDomesticCrating(appCtx, models.ReServiceCodeDCRT, contractCode, referenceDate, billedCubicFeet, serviceSchedule, standaloneCrate, standaloneCrateCap)
}

// PriceUsingParams determines the price for domestic destination first day SIT given PaymentServiceItemParams
func (p domesticCratingPricer) PriceUsingParams(appCtx appcontext.AppContext, params models.PaymentServiceItemParams, featureFlagValues map[string]bool) (unit.Cents, services.PricingDisplayParams, error) {
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

	serviceScheduleDestination, err := getParamInt(params, models.ServiceItemParamNameServicesScheduleOrigin)
	if err != nil {
		return unit.Cents(0), nil, err
	}

	standaloneCrate, err := getParamBool(params, models.ServiceItemParamNameStandaloneCrate)
	if err != nil {
		return unit.Cents(0), nil, err
	}

	standaloneCrateCapParam, err := getParamInt(params, models.ServiceItemParamNameStandaloneCrateCap)
	if err != nil {
		return unit.Cents(0), nil, err
	}
	standaloneCrateCap := unit.Cents(float64(standaloneCrateCapParam))

	return p.Price(appCtx, contractCode, referenceDate, cubicFeetBilled, serviceScheduleDestination, standaloneCrate, standaloneCrateCap)
}
