package ghcrateengine

import (
	"time"

	"github.com/transcom/mymove/pkg/appconfig"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/unit"
)

type domesticDestinationFirstDaySITPricer struct {
}

// NewDomesticDestinationFirstDaySITPricer creates a new pricer for domestic destination first day SIT
func NewDomesticDestinationFirstDaySITPricer() services.DomesticDestinationFirstDaySITPricer {
	return &domesticDestinationFirstDaySITPricer{}
}

// Price determines the price for domestic destination first day SIT
func (p domesticDestinationFirstDaySITPricer) Price(appCfg appconfig.AppConfig, contractCode string, requestedPickupDate time.Time, weight unit.Pound, serviceArea string) (unit.Cents, services.PricingDisplayParams, error) {
	return priceDomesticFirstDaySIT(appCfg, models.ReServiceCodeDDFSIT, contractCode, requestedPickupDate, weight, serviceArea)
}

// PriceUsingParams determines the price for domestic destination first day SIT given PaymentServiceItemParams
func (p domesticDestinationFirstDaySITPricer) PriceUsingParams(appCfg appconfig.AppConfig, params models.PaymentServiceItemParams) (unit.Cents, services.PricingDisplayParams, error) {
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

	return p.Price(appCfg, contractCode, requestedPickupDate, unit.Pound(weightBilledActual), serviceAreaDest)
}
