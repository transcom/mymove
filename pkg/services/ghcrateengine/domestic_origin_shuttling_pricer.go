package ghcrateengine

import (
	"time"

	"github.com/transcom/mymove/pkg/appconfig"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/unit"
)

type domesticOriginShuttlingPricer struct {
}

// NewDomesticOriginShuttlingPricer creates a new pricer for domestic origin first day SIT
func NewDomesticOriginShuttlingPricer() services.DomesticOriginShuttlingPricer {
	return &domesticOriginShuttlingPricer{}
}

// Price determines the price for domestic origin first day SIT
func (p domesticOriginShuttlingPricer) Price(appCfg appconfig.AppConfig, contractCode string, requestedPickupDate time.Time, weight unit.Pound, serviceSchedule int) (unit.Cents, services.PricingDisplayParams, error) {
	return priceDomesticShuttling(appCfg, models.ReServiceCodeDOSHUT, contractCode, requestedPickupDate, weight, serviceSchedule)
}

// PriceUsingParams determines the price for domestic origin first day SIT given PaymentServiceItemParams
func (p domesticOriginShuttlingPricer) PriceUsingParams(appCfg appconfig.AppConfig, params models.PaymentServiceItemParams) (unit.Cents, services.PricingDisplayParams, error) {
	contractCode, err := getParamString(params, models.ServiceItemParamNameContractCode)
	if err != nil {
		return unit.Cents(0), nil, err
	}

	weightBilledActual, err := getParamInt(params, models.ServiceItemParamNameWeightBilledActual)
	if err != nil {
		return unit.Cents(0), nil, err
	}

	requestedPickupDate, err := getParamTime(params, models.ServiceItemParamNameRequestedPickupDate)
	if err != nil {
		return unit.Cents(0), nil, err
	}

	serviceScheduleOrigin, err := getParamInt(params, models.ServiceItemParamNameServicesScheduleOrigin)
	if err != nil {
		return unit.Cents(0), nil, err
	}

	return p.Price(appCfg, contractCode, requestedPickupDate, unit.Pound(weightBilledActual), serviceScheduleOrigin)
}
