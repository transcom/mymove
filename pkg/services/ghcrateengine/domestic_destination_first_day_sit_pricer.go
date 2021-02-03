package ghcrateengine

import (
	"time"

	"github.com/gobuffalo/pop/v5"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/unit"
)

type domesticDestinationFirstDaySITPricer struct {
	db *pop.Connection
}

// NewDomesticDestinationFirstDaySITPricer creates a new pricer for domestic destination first day SIT
func NewDomesticDestinationFirstDaySITPricer(db *pop.Connection) services.DomesticDestinationFirstDaySITPricer {
	return &domesticDestinationFirstDaySITPricer{
		db: db,
	}
}

// Price determines the price for domestic destination first day SIT
func (p domesticDestinationFirstDaySITPricer) Price(contractCode string, requestedPickupDate time.Time, isPeakPeriod bool, weight unit.Pound, serviceArea string) (unit.Cents, error) {
	return priceDomesticFirstDaySIT(p.db, models.ReServiceCodeDDFSIT, contractCode, requestedPickupDate, isPeakPeriod, weight, serviceArea)
}

// PriceUsingParams determines the price for domestic destination first day SIT given PaymentServiceItemParams
func (p domesticDestinationFirstDaySITPricer) PriceUsingParams(params models.PaymentServiceItemParams) (unit.Cents, error) {
	contractCode, err := getParamString(params, models.ServiceItemParamNameContractCode)
	if err != nil {
		return unit.Cents(0), err
	}

	requestedPickupDate, err := getParamTime(params, models.ServiceItemParamNameRequestedPickupDate)
	if err != nil {
		return unit.Cents(0), err
	}

	weightBilledActual, err := getParamInt(params, models.ServiceItemParamNameWeightBilledActual)
	if err != nil {
		return unit.Cents(0), err
	}

	serviceAreaDest, err := getParamString(params, models.ServiceItemParamNameServiceAreaDest)
	if err != nil {
		return unit.Cents(0), err
	}

	isPeakPeriod := IsPeakPeriod(requestedPickupDate)

	return p.Price(contractCode, requestedPickupDate, isPeakPeriod, unit.Pound(weightBilledActual), serviceAreaDest)
}
