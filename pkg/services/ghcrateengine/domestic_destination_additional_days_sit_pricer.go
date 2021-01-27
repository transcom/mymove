package ghcrateengine

import (
	"time"

	"github.com/gobuffalo/pop/v5"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/unit"
)

type domesticDestinationAdditionalDaysSITPricer struct {
	db *pop.Connection
}

// NewDomesticDestinationAdditionalDaysSITPricer creates a new pricer for domestic destination additional days SIT
func NewDomesticDestinationAdditionalDaysSITPricer(db *pop.Connection) services.DomesticDestinationAdditionalDaysSITPricer {
	return &domesticDestinationAdditionalDaysSITPricer{
		db: db,
	}
}

// Price determines the price for domestic destination additional days SIT
func (p domesticDestinationAdditionalDaysSITPricer) Price(contractCode string, requestedPickupDate time.Time, isPeakPeriod bool, weight unit.Pound, serviceArea string, numberOfDaysInSIT int) (unit.Cents, error) {
	return priceDomesticAdditionalDaysSit(p.db, models.ReServiceCodeDDASIT, contractCode, requestedPickupDate, isPeakPeriod, weight, serviceArea, numberOfDaysInSIT)
}

// PriceUsingParams determines the price for domestic destination first day SIT given PaymentServiceItemParams
func (p domesticDestinationAdditionalDaysSITPricer) PriceUsingParams(params models.PaymentServiceItemParams) (unit.Cents, error) {
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

	numberOfDaysInSIT, err := getParamInt(params, models.ServiceItemParamNameNumberDaysSIT)
	if err != nil {
		return unit.Cents(0), err
	}

	isPeakPeriod := IsPeakPeriod(requestedPickupDate)

	return p.Price(contractCode, requestedPickupDate, isPeakPeriod, unit.Pound(weightBilledActual), serviceAreaDest, numberOfDaysInSIT)
}
