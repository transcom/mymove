package ghcrateengine

import (
	"fmt"
	"math"
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

// Price determines the price for a domestic linehaul
func (p domesticDestinationFirstDaySITPricer) Price(contractCode string, requestedPickupDate time.Time, isPeakPeriod bool, weight unit.Pound, serviceArea string) (unit.Cents, error) {
	if weight < minDomesticWeight {
		return 0, fmt.Errorf("weight of %d less than the minimum of %d", weight, minDomesticWeight)
	}

	serviceAreaPrice, err := fetchDomServiceAreaPrice(p.db, contractCode, models.ReServiceCodeDDFSIT, serviceArea, isPeakPeriod)
	if err != nil {
		return unit.Cents(0), fmt.Errorf("could not fetch domestic destination first day SIT rate: %w", err)
	}

	contractYear, err := fetchContractYear(p.db, serviceAreaPrice.ContractID, requestedPickupDate)
	if err != nil {
		return unit.Cents(0), fmt.Errorf("could not fetch contract year: %w", err)
	}

	baseTotalPrice := serviceAreaPrice.PriceCents.Float64() * weight.ToCWTFloat64()
	escalatedTotalPrice := baseTotalPrice * contractYear.EscalationCompounded

	totalPriceCents := unit.Cents(math.Round(escalatedTotalPrice))

	return totalPriceCents, nil
}

// PriceUsingParams determines the price for a domestic linehaul given PaymentServiceItemParams
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
