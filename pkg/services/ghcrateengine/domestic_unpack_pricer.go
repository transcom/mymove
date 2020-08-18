package ghcrateengine

import (
	"fmt"
	"math"
	"time"

	"github.com/transcom/mymove/pkg/models"

	"github.com/gobuffalo/pop"
	"github.com/pkg/errors"

	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/unit"
)

type domesticUnpackPricer struct {
	db *pop.Connection
}

// NewDomesticUnpackPricer creates a new pricer for domestic pack services
func NewDomesticUnpackPricer(db *pop.Connection) services.DomesticUnpackPricer {
	return &domesticUnpackPricer{
		db: db,
	}
}

// Price determines the price for a domestic pack/unpack service
func (p domesticUnpackPricer) Price(contractCode string, requestedPickupDate time.Time, weight unit.Pound, servicesScheduleDest int) (totalCost unit.Cents, err error) {
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
	if servicesScheduleDest == 0 {
		return 0, errors.New("Service schedule is required")
	}

	isPeakPeriod := IsPeakPeriod(requestedPickupDate)
	// look up rate for domestic pack/unpack price
	var contractYear models.ReContractYear
	// domOtherPrice, err := fetchDomOtherPrice(p.db, contractCode, models.ReServiceCodeDUPK, servicesScheduleDest, isPeakPeriod)
	var domOtherPrice models.ReDomesticOtherPrice
	err = p.db.Q().
		Join("re_services", "service_id = re_services.id").
		Join("re_contracts", "re_contracts.id = re_domestic_other_prices.contract_id").
		Where("re_contracts.code = $1", contractCode).
		Where("re_services.code = $2", models.ReServiceCodeDUPK).
		Where("schedule = $3", servicesScheduleDest).
		Where("is_peak_period = $4", isPeakPeriod).
		First(&domOtherPrice)

	if err != nil {
		return 0, fmt.Errorf("Could not lookup Domestic Other Price: %w", err)
	}

	err = p.db.Where("contract_id = $1", domOtherPrice.ContractID).
		Where("$2 between start_date and end_date", requestedPickupDate).
		First(&contractYear)
	if err != nil {
		return 0, fmt.Errorf("Could not lookup contract year: %w", err)
	}

	basePrice := domOtherPrice.PriceCents.Float64() * weight.ToCWTFloat64()
	escalatedPrice := basePrice * contractYear.EscalationCompounded
	totalCost = unit.Cents(math.Round(escalatedPrice))

	return totalCost, err
}

// PriceUsingParams determines the price for a domestic pack given PaymentServiceItemParams
func (p domesticUnpackPricer) PriceUsingParams(params models.PaymentServiceItemParams) (unit.Cents, error) {
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

	servicesScheduleDest, err := getParamInt(params, models.ServiceItemParamNameServicesScheduleDest)
	if err != nil {
		return unit.Cents(0), err
	}

	return p.Price(contractCode, requestedPickupDate, unit.Pound(weightBilledActual), servicesScheduleDest)
}
