package ghcrateengine

import (
	"time"

	"github.com/gobuffalo/pop/v5"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/unit"
)

type domesticDestinationShuttlingPricer struct {
	db *pop.Connection
}

// NewDomesticDestinationShuttlingPricer creates a new pricer for domestic destination first day SIT
func NewDomesticDestinationShuttlingPricer(db *pop.Connection) services.DomesticDestinationShuttlingPricer {
	return &domesticDestinationShuttlingPricer{
		db: db,
	}
}

// Price determines the price for domestic destination first day SIT
func (p domesticDestinationShuttlingPricer) Price(contractCode string, requestedPickupDate time.Time, weight unit.Pound, serviceSchedule int) (unit.Cents, services.PricingDisplayParams, error) {
	return priceDomesticShuttling(p.db, models.ReServiceCodeDDSHUT, contractCode, requestedPickupDate, weight, serviceSchedule)
}

// PriceUsingParams determines the price for domestic destination first day SIT given PaymentServiceItemParams
func (p domesticDestinationShuttlingPricer) PriceUsingParams(params models.PaymentServiceItemParams) (unit.Cents, services.PricingDisplayParams, error) {
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

	serviceScheduleDestination, err := getParamInt(params, models.ServiceItemParamNameServicesScheduleDest)
	if err != nil {
		return unit.Cents(0), nil, err
	}

	return p.Price(contractCode, requestedPickupDate, unit.Pound(weightBilledActual), serviceScheduleDestination)
}
