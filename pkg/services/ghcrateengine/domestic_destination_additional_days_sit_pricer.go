package ghcrateengine

import (
	"time"

	"github.com/transcom/mymove/pkg/appconfig"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/unit"
)

type domesticDestinationAdditionalDaysSITPricer struct {
}

// NewDomesticDestinationAdditionalDaysSITPricer creates a new pricer for domestic destination additional days SIT
func NewDomesticDestinationAdditionalDaysSITPricer() services.DomesticDestinationAdditionalDaysSITPricer {
	return &domesticDestinationAdditionalDaysSITPricer{}
}

// Price determines the price for domestic destination additional days SIT
func (p domesticDestinationAdditionalDaysSITPricer) Price(appCfg appconfig.AppConfig, contractCode string, requestedPickupDate time.Time, weight unit.Pound, serviceArea string, numberOfDaysInSIT int) (unit.Cents, services.PricingDisplayParams, error) {
	return priceDomesticAdditionalDaysSIT(appCfg, models.ReServiceCodeDDASIT, contractCode, requestedPickupDate, weight, serviceArea, numberOfDaysInSIT)
}

// PriceUsingParams determines the price for domestic destination first day SIT given PaymentServiceItemParams
func (p domesticDestinationAdditionalDaysSITPricer) PriceUsingParams(appCfg appconfig.AppConfig, params models.PaymentServiceItemParams) (unit.Cents, services.PricingDisplayParams, error) {
	contractCode, err := getParamString(params, models.ServiceItemParamNameContractCode)
	if err != nil {
		return unit.Cents(0), nil, err
	}

	requestedPickupDate, err := getParamTime(params, models.ServiceItemParamNameRequestedPickupDate)
	if err != nil {
		return unit.Cents(0), nil, err
	}

	weightBilledActual, err := getParamInt(params, models.ServiceItemParamNameWeightBilledActual)
	if err != nil {
		return unit.Cents(0), nil, err
	}

	serviceAreaDest, err := getParamString(params, models.ServiceItemParamNameServiceAreaDest)
	if err != nil {
		return unit.Cents(0), nil, err
	}

	numberOfDaysInSIT, err := getParamInt(params, models.ServiceItemParamNameNumberDaysSIT)
	if err != nil {
		return unit.Cents(0), nil, err
	}

	return p.Price(appCfg, contractCode, requestedPickupDate, unit.Pound(weightBilledActual), serviceAreaDest, numberOfDaysInSIT)
}
