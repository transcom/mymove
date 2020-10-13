package ghcrateengine

import (
	"fmt"
	"math"
	"time"

	"github.com/gobuffalo/pop/v5"
	"github.com/pkg/errors"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/unit"
)

type domesticDestinationPricer struct {
	db *pop.Connection
}

// NewDomesticDestinationPricer instantiates a new pricer
func NewDomesticDestinationPricer(db *pop.Connection) services.DomesticDestinationPricer {
	return &domesticDestinationPricer{
		db: db,
	}
}

// Price determines the price for the destination service area
func (p domesticDestinationPricer) Price(contractCode string, requestedPickupDate time.Time, weight unit.Pound, serviceArea string) (totalCost unit.Cents, err error) {
	// Validate parameters
	if len(contractCode) == 0 {
		return 0, errors.New("ContractCode is required")
	}
	if requestedPickupDate.IsZero() {
		return 0, errors.New("RequestedPickupDate is required")
	}
	if weight < minDomesticWeight {
		return 0, fmt.Errorf("Weight must be a minimum of %d", minDomesticWeight)
	}
	if len(serviceArea) == 0 {
		return 0, errors.New("ServiceArea is required")
	}

	isPeakPeriod := IsPeakPeriod(requestedPickupDate)
	// look up rate for domestic destination price
	var contractYear models.ReContractYear
	var domServiceAreaPrice models.ReDomesticServiceAreaPrice
	err = p.db.Q().
		Join("re_domestic_service_areas sa", "domestic_service_area_id = sa.id").
		Join("re_services", "service_id = re_services.id").
		Join("re_contracts", "re_contracts.id = re_domestic_service_area_prices.contract_id").
		Where("sa.service_area = $1", serviceArea).
		Where("re_services.code = $2", models.ReServiceCodeDDP).
		Where("re_contracts.code = $3", contractCode).
		Where("is_peak_period = $4", isPeakPeriod).
		First(&domServiceAreaPrice)
	if err != nil {
		return 0, fmt.Errorf("Could not lookup Domestic Service Area Price: %w", err)
	}
	err = p.db.Where("contract_id = $1", domServiceAreaPrice.ContractID).
		Where("$2 between start_date and end_date", requestedPickupDate).
		First(&contractYear)
	if err != nil {
		return 0, fmt.Errorf("Could not lookup contract year: %w", err)
	}

	basePrice := domServiceAreaPrice.PriceCents.Float64() * weight.ToCWTFloat64()
	escalatedPrice := basePrice * contractYear.EscalationCompounded
	totalCost = unit.Cents(math.Round(escalatedPrice))

	return totalCost, err
}

func (p domesticDestinationPricer) PriceUsingParams(params models.PaymentServiceItemParams) (unit.Cents, error) {
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

	total, err := p.Price(contractCode, requestedPickupDate, unit.Pound(weightBilledActual), serviceAreaDest)
	return total, err
}
