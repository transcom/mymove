package ghcrateengine

import (
	"time"

	"github.com/gobuffalo/pop/v5"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/unit"
)

type domesticOriginFirstDaySITPricer struct {
	db *pop.Connection
}

// NewDomesticOriginFirstDaySITPricer creates a new pricer for domestic origin first day SIT
func NewDomesticOriginFirstDaySITPricer(db *pop.Connection) services.DomesticOriginFirstDaySITPricer {
	return &domesticOriginFirstDaySITPricer{
		db: db,
	}
}

// Price determines the price for domestic origin first day SIT
func (p domesticOriginFirstDaySITPricer) Price(contractCode string, requestedPickupDate time.Time, weight unit.Pound, serviceArea string) (unit.Cents, services.PricingDisplayParams, error) {
	return priceDomesticFirstDaySIT(p.db, models.ReServiceCodeDOFSIT, contractCode, requestedPickupDate, weight, serviceArea)
}

// PriceUsingParams determines the price for domestic origin first day SIT given PaymentServiceItemParams
func (p domesticOriginFirstDaySITPricer) PriceUsingParams(params models.PaymentServiceItemParams) (unit.Cents, services.PricingDisplayParams, error) {
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

	serviceAreaOrigin, err := getParamString(params, models.ServiceItemParamNameServiceAreaOrigin)
	if err != nil {
		return unit.Cents(0), nil, err
	}

	return p.Price(contractCode, requestedPickupDate, unit.Pound(weightBilledActual), serviceAreaOrigin)
}
