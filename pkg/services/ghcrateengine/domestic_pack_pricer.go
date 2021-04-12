package ghcrateengine

import (
	"fmt"
	"math"
	"time"

	"github.com/transcom/mymove/pkg/models"

	"github.com/gobuffalo/pop/v5"
	"github.com/pkg/errors"

	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/unit"
)

type domesticPackPricer struct {
	db *pop.Connection
}

// NewDomesticPackPricer creates a new pricer for domestic pack services
func NewDomesticPackPricer(db *pop.Connection) services.DomesticPackPricer {
	return &domesticPackPricer{
		db: db,
	}
}

// Price determines the price for a domestic pack/unpack service
func (p domesticPackPricer) Price(contractCode string, requestedPickupDate time.Time, weight unit.Pound, servicesScheduleOrigin int) (unit.Cents, services.PricingDisplayParams, error) {
	// Validate parameters
	if len(contractCode) == 0 {
		return 0, nil, errors.New("ContractCode is required")
	}
	if requestedPickupDate.IsZero() {
		return 0, nil, errors.New("RequestedPickupDate is required")
	}
	if weight < minDomesticWeight {
		return 0, nil, fmt.Errorf("Weight must be a minimum of %d", minDomesticWeight)
	}
	if servicesScheduleOrigin == 0 {
		return 0, nil, errors.New("Service schedule is required")
	}

	isPeakPeriod := IsPeakPeriod(requestedPickupDate)
	var contractYear models.ReContractYear
	domOtherPrice, err := fetchDomOtherPrice(p.db, contractCode, models.ReServiceCodeDPK, servicesScheduleOrigin, isPeakPeriod)

	if err != nil {
		return 0, nil, fmt.Errorf("Could not lookup Domestic Other Price: %w", err)
	}

	err = p.db.Where("contract_id = $1", domOtherPrice.ContractID).
		Where("$2 between start_date and end_date", requestedPickupDate).
		First(&contractYear)
	if err != nil {
		return 0, nil, fmt.Errorf("Could not lookup contract year: %w", err)
	}

	basePrice := domOtherPrice.PriceCents.Float64() * weight.ToCWTFloat64()
	escalatedPrice := basePrice * contractYear.EscalationCompounded
	totalCost := unit.Cents(math.Round(escalatedPrice))

	displayParams := services.PricingDisplayParams{
		{
			Key:   models.ServiceItemParamNameContractYearName,
			Value: contractYear.Name,
		},
		{
			Key:   models.ServiceItemParamNamePriceRateOrFactor,
			Value: FormatCents(domOtherPrice.PriceCents),
		},
		{
			Key:   models.ServiceItemParamNameIsPeak,
			Value: FormatBool(isPeakPeriod),
		},
		{
			Key:   models.ServiceItemParamNameEscalationCompounded,
			Value: FormatFloat(contractYear.EscalationCompounded, 5),
		},
	}

	return totalCost, displayParams, nil
}

// PriceUsingParams determines the price for a domestic pack given PaymentServiceItemParams
func (p domesticPackPricer) PriceUsingParams(params models.PaymentServiceItemParams) (unit.Cents, services.PricingDisplayParams, error) {
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

	servicesScheduleOrigin, err := getParamInt(params, models.ServiceItemParamNameServicesScheduleOrigin)
	if err != nil {
		return unit.Cents(0), nil, err
	}

	return p.Price(contractCode, requestedPickupDate, unit.Pound(weightBilledActual), servicesScheduleOrigin)
}
