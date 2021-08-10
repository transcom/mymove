package ghcrateengine

import (
	"time"

	"github.com/transcom/mymove/pkg/appconfig"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/unit"
)

type domesticUncratingPricer struct {
}

// NewDomesticUncratingPricer creates a new pricer for domestic destination first day SIT
func NewDomesticUncratingPricer() services.DomesticUncratingPricer {
	return &domesticUncratingPricer{}
}

// Price determines the price for domestic destination first day SIT
func (p domesticUncratingPricer) Price(appCfg appconfig.AppConfig, contractCode string, requestedPickupDate time.Time, billedCubicFeet unit.CubicFeet, serviceSchedule int) (unit.Cents, services.PricingDisplayParams, error) {
	return priceDomesticCrating(appCfg, models.ReServiceCodeDUCRT, contractCode, requestedPickupDate, billedCubicFeet, serviceSchedule)
}

// PriceUsingParams determines the price for domestic destination first day SIT given PaymentServiceItemParams
func (p domesticUncratingPricer) PriceUsingParams(appCfg appconfig.AppConfig, params models.PaymentServiceItemParams) (unit.Cents, services.PricingDisplayParams, error) {
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

	serviceScheduleDestination, err := getParamInt(params, models.ServiceItemParamNameServicesScheduleDest)
	if err != nil {
		return unit.Cents(0), nil, err
	}

	return p.Price(appCfg, contractCode, requestedPickupDate, cubicFeetBilled, serviceScheduleDestination)
}
