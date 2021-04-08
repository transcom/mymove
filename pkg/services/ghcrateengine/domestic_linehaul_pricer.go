package ghcrateengine

import (
	"fmt"
	"time"

	"github.com/transcom/mymove/pkg/models"

	"github.com/gobuffalo/pop/v5"
	"github.com/pkg/errors"

	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/unit"
)

const (
	dlhPricerMinimumWeight   = unit.Pound(500)
	dlhPricerMinimumDistance = unit.Miles(50)
)

type domesticLinehaulPricer struct {
	db *pop.Connection
}

// NewDomesticLinehaulPricer creates a new pricer for domestic linehaul services
func NewDomesticLinehaulPricer(db *pop.Connection) services.DomesticLinehaulPricer {
	return &domesticLinehaulPricer{
		db: db,
	}
}

// Price determines the price for a domestic linehaul
func (p domesticLinehaulPricer) Price(contractCode string, requestedPickupDate time.Time, distance unit.Miles, weight unit.Pound, serviceArea string) (unit.Cents, services.PricingDisplayParams, error) {
	// Validate parameters
	if len(contractCode) == 0 {
		return 0, nil, errors.New("ContractCode is required")
	}
	if requestedPickupDate.IsZero() {
		return 0, nil, errors.New("RequestedPickupDate is required")
	}
	if distance < dlhPricerMinimumDistance {
		return 0, nil, fmt.Errorf("Distance must be at least %d", dlhPricerMinimumDistance)
	}
	if weight < dlhPricerMinimumWeight {
		return 0, nil, fmt.Errorf("Weight must be at least %d", dlhPricerMinimumWeight)
	}
	if len(serviceArea) == 0 {
		return 0, nil, errors.New("ServiceArea is required")
	}

	isPeakPeriod := IsPeakPeriod(requestedPickupDate)
	domesticLinehaulPrice, err := fetchDomesticLinehaulPrice(p.db, contractCode, isPeakPeriod, distance, weight, serviceArea)
	if err != nil {
		return unit.Cents(0), nil, fmt.Errorf("could not fetch domestic linehaul rate: %w", err)
	}

	contractYear, err := fetchContractYear(p.db, domesticLinehaulPrice.ContractID, requestedPickupDate)
	if err != nil {
		return 0, nil, fmt.Errorf("Could not lookup contract year: %w", err)
	}

	baseTotalPrice := weight.ToCWTFloat64() * distance.Float64() * domesticLinehaulPrice.PriceMillicents.Float64()
	escalatedTotalPrice := contractYear.EscalationCompounded * baseTotalPrice

	totalPriceMillicents := unit.Millicents(escalatedTotalPrice)
	totalPriceCents := totalPriceMillicents.ToCents()

	params := services.PricingDisplayParams{
		{Key: models.ServiceItemParamNameContractYearName, Value: contractYear.Name},
		{Key: models.ServiceItemParamNameEscalationCompounded, Value: FormatFloat(contractYear.EscalationCompounded, 5)},
		{Key: models.ServiceItemParamNameIsPeak, Value: FormatBool(isPeakPeriod)},
		{Key: models.ServiceItemParamNamePriceRateOrFactor, Value: FormatFloat(domesticLinehaulPrice.PriceMillicents.ToDollarFloatNoRound(), 3)},
	}

	return totalPriceCents, params, nil
}

// PriceUsingParams determines the price for a domestic linehaul given PaymentServiceItemParams
func (p domesticLinehaulPricer) PriceUsingParams(params models.PaymentServiceItemParams) (unit.Cents, services.PricingDisplayParams, error) {
	contractCode, err := getParamString(params, models.ServiceItemParamNameContractCode)
	if err != nil {
		return unit.Cents(0), nil, err
	}

	requestedPickupDate, err := getParamTime(params, models.ServiceItemParamNameRequestedPickupDate)
	if err != nil {
		return unit.Cents(0), nil, err
	}

	distanceZip3, err := getParamInt(params, models.ServiceItemParamNameDistanceZip3)
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

	return p.Price(contractCode, requestedPickupDate, unit.Miles(distanceZip3), unit.Pound(weightBilledActual), serviceAreaOrigin)
}

func fetchDomesticLinehaulPrice(db *pop.Connection, contractCode string, isPeakPeriod bool, distance unit.Miles, weight unit.Pound, serviceArea string) (models.ReDomesticLinehaulPrice, error) {
	var domesticLinehaulPrice models.ReDomesticLinehaulPrice
	err := db.Q().
		Join("re_domestic_service_areas sa", "domestic_service_area_id = sa.id").
		Join("re_contracts c", "re_domestic_linehaul_prices.contract_id = c.id").
		Where("c.code = $1", contractCode).
		Where("re_domestic_linehaul_prices.is_peak_period = $2", isPeakPeriod).
		Where("$3 between weight_lower and weight_upper", weight).
		Where("$4 between miles_lower and miles_upper", distance).
		Where("sa.service_area = $5", serviceArea).
		First(&domesticLinehaulPrice)

	if err != nil {
		return models.ReDomesticLinehaulPrice{}, err
	}

	return domesticLinehaulPrice, nil
}
