package ghcrateengine

import (
	"time"

	"github.com/gobuffalo/pop/v5"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/unit"
)

type domesticOriginShuttlingPricer struct {
	db *pop.Connection
}

// NewDomesticOriginShuttlingPricer creates a new pricer for domestic origin first day SIT
func NewDomesticOriginShuttlingPricer(db *pop.Connection) services.DomesticOriginShuttlingPricer {
	return &domesticOriginShuttlingPricer{
		db: db,
	}
}

// Price determines the price for domestic origin first day SIT
func (p domesticOriginShuttlingPricer) Price(contractCode string, requestedPickupDate time.Time, weight unit.Pound, serviceSchedule int) (unit.Cents, services.PricingDisplayParams, error) {
	return priceDomesticShuttling(p.db, models.ReServiceCodeDOSHUT, contractCode, requestedPickupDate, weight, serviceSchedule)
}

// PriceUsingParams determines the price for domestic origin first day SIT given PaymentServiceItemParams
func (p domesticOriginShuttlingPricer) PriceUsingParams(params models.PaymentServiceItemParams) (unit.Cents, services.PricingDisplayParams, error) {
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

	return p.Price(contractCode, requestedPickupDate, unit.Pound(weightBilledActual), serviceScheduleOrigin)
}
