package ghcrateengine

import (
	"time"

	"github.com/transcom/mymove/pkg/appconfig"
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
func (p domesticCratingPricer) Price(appCfg appconfig.AppConfig, contractCode string, requestedPickupDate time.Time, billedCubicFeet unit.CubicFeet, serviceSchedule int) (unit.Cents, services.PricingDisplayParams, error) {
	return priceDomesticCrating(appCfg, models.ReServiceCodeDCRT, contractCode, requestedPickupDate, billedCubicFeet, serviceSchedule)
}

// PriceUsingParams determines the price for domestic destination first day SIT given PaymentServiceItemParams
func (p domesticCratingPricer) PriceUsingParams(appCfg appconfig.AppConfig, params models.PaymentServiceItemParams) (unit.Cents, services.PricingDisplayParams, error) {
	contractCode, err := getParamString(params, models.ServiceItemParamNameContractCode)
	if err != nil {
		return unit.Cents(0), nil, err
	}

	cubicFeetFloat, err := getParamFloat(params, models.ServiceItemParamNameCubicFeetBilled)
	if err != nil {
		return unit.Cents(0), nil, err
	}

	cubicFeetBilled := unit.CubicFeet(cubicFeetFloat)

	requestedPickupDate, err := getParamTime(params, models.ServiceItemParamNameRequestedPickupDate)
	if err != nil {
		return unit.Cents(0), nil, err
	}

	serviceScheduleDestination, err := getParamInt(params, models.ServiceItemParamNameServicesScheduleOrigin)
	if err != nil {
		return unit.Cents(0), nil, err
	}

	return p.Price(appCfg, contractCode, requestedPickupDate, cubicFeetBilled, serviceScheduleDestination)
}
