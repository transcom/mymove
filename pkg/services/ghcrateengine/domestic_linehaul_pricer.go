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
func (p domesticLinehaulPricer) Price(isShortHaul bool, contractCode string, requestedPickupDate time.Time, distance unit.Miles, weight unit.Pound, serviceArea string) (unit.Cents, services.PricingDisplayParams, error) {
	// Validate parameters
	if len(contractCode) == 0 {
		return 0, nil, errors.New("ContractCode is required")
	}
	if requestedPickupDate.IsZero() {
		return 0, nil, errors.New("RequestedPickupDate is required")
	}
	if distance <= 0 {
		return 0, nil, errors.New("Distance must be greater than 0")
	}
	if weight < minDomesticWeight {
		return 0, nil, fmt.Errorf("Weight must be at least %d", minDomesticWeight)
	}
	if len(serviceArea) == 0 {
		return 0, nil, errors.New("ServiceArea is required")
	}

	var contractYear models.ReContractYear
	var totalPriceCents unit.Cents
	var ratePriceDollarFloat string
	var err error

	isPeakPeriod := IsPeakPeriod(requestedPickupDate)

	if isShortHaul {
		ratePriceDollarFloat, totalPriceCents, contractYear, err = p.priceShorthaul(contractCode, requestedPickupDate, isPeakPeriod, distance, weight, serviceArea)
		if err != nil {
			return 0, nil, err
		}
	} else {
		ratePriceDollarFloat, totalPriceCents, contractYear, err = p.priceLinehaul(contractCode, requestedPickupDate, isPeakPeriod, distance, weight, serviceArea)
		if err != nil {
			return 0, nil, err
		}
	}

	params := services.PricingDisplayParams{
		{Key: models.ServiceItemParamNameContractYearName, Value: contractYear.Name},
		{Key: models.ServiceItemParamNameEscalationCompounded, Value: FormatEscalation(contractYear.EscalationCompounded)},
		{Key: models.ServiceItemParamNameIsPeak, Value: FormatBool(isPeakPeriod)},
		{Key: models.ServiceItemParamNamePriceRateOrFactor, Value: ratePriceDollarFloat},
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

	weightBilledActual, err := getParamInt(params, models.ServiceItemParamNameWeightBilledActual)
	if err != nil {
		return unit.Cents(0), nil, err
	}

	serviceAreaOrigin, err := getParamString(params, models.ServiceItemParamNameServiceAreaOrigin)
	if err != nil {
		return unit.Cents(0), nil, err
	}

	zipPickup, err := getParamString(params, models.ServiceItemParamNameZipPickupAddress)
	if err != nil {
		return unit.Cents(0), nil, err
	}

	zipDestination, err := getParamString(params, models.ServiceItemParamNameZipDestAddress)
	if err != nil {
		return unit.Cents(0), nil, err
	}

	isShortHaul := isSameZip3(zipPickup, zipDestination)
	var distanceZip int
	if isShortHaul {
		distanceZip, err = getParamInt(params, models.ServiceItemParamNameDistanceZip5)
		if err != nil {
			return unit.Cents(0), nil, err
		}
	} else {
		distanceZip, err = getParamInt(params, models.ServiceItemParamNameDistanceZip3)
		if err != nil {
			return unit.Cents(0), nil, err
		}
	}

	return p.Price(isShortHaul, contractCode, requestedPickupDate, unit.Miles(distanceZip), unit.Pound(weightBilledActual), serviceAreaOrigin)
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

func (p domesticLinehaulPricer) priceLinehaul(contractCode string, requestedPickupDate time.Time, isPeakPeriod bool, distance unit.Miles, weight unit.Pound, serviceArea string) (string, unit.Cents, models.ReContractYear, error) {
	domesticLinehaulPrice, err := fetchDomesticLinehaulPrice(p.db, contractCode, isPeakPeriod, distance, weight, serviceArea)
	if err != nil {
		return "", unit.Cents(0), models.ReContractYear{}, fmt.Errorf("could not fetch domestic linehaul rate: %w", err)
	}

	ratePriceDollarFloat := FormatFloat(domesticLinehaulPrice.PriceMillicents.ToDollarFloatNoRound(), 3)

	contractYear, err := fetchContractYear(p.db, domesticLinehaulPrice.ContractID, requestedPickupDate)
	if err != nil {
		return "", 0, models.ReContractYear{}, fmt.Errorf("Could not lookup contract year: %w", err)
	}

	// linehaul price is in millicents
	baseTotalPrice := weight.ToCWTFloat64() * distance.Float64() * domesticLinehaulPrice.PriceMillicents.Float64()
	escalatedTotalPrice := contractYear.EscalationCompounded * baseTotalPrice
	totalPriceMillicents := unit.Millicents(escalatedTotalPrice)
	totalPriceCents := totalPriceMillicents.ToCents()

	return ratePriceDollarFloat, totalPriceCents, contractYear, nil
}

func (p domesticLinehaulPricer) priceShorthaul(contractCode string, requestedPickupDate time.Time, isPeakPeriod bool, distance unit.Miles, weight unit.Pound, serviceArea string) (string, unit.Cents, models.ReContractYear, error) {
	domServiceAreaPrice, err := fetchDomServiceAreaPrice(p.db, contractCode, models.ReServiceCodeDSH, serviceArea, isPeakPeriod)
	if err != nil {
		return "", 0, models.ReContractYear{}, fmt.Errorf("Could not lookup Domestic Service Area Price: %w", err)
	}

	ratePriceDollarFloat := FormatFloat(domServiceAreaPrice.PriceCents.ToDollarFloatNoRound(), 3)

	contractYear, err := fetchContractYear(p.db, domServiceAreaPrice.ContractID, requestedPickupDate)
	if err != nil {
		return "", 0, models.ReContractYear{}, fmt.Errorf("Could not lookup contract year: %w", err)
	}

	// shorthaul price is in cents
	baseTotalPrice := domServiceAreaPrice.PriceCents.Float64() * distance.Float64() * weight.ToCWTFloat64()
	escalatedTotalPrice := contractYear.EscalationCompounded * baseTotalPrice
	totalPriceCents := unit.Cents(math.Round(escalatedTotalPrice))

	return ratePriceDollarFloat, totalPriceCents, contractYear, nil
}
