package ghcrateengine

import (
	"fmt"
	"time"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"

	"github.com/pkg/errors"

	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/unit"
)

const (
	dlhPricerMinimumWeight   = unit.Pound(500)
	dlhPricerMinimumDistance = unit.Miles(50)
)

type domesticLinehaulPricer struct {
}

// NewDomesticLinehaulPricer creates a new pricer for domestic linehaul services
func NewDomesticLinehaulPricer() services.DomesticLinehaulPricer {
	return &domesticLinehaulPricer{}
}

// Price determines the price for a domestic linehaul
func (p domesticLinehaulPricer) Price(appCtx appcontext.AppContext, contractCode string, referenceDate time.Time, distance unit.Miles, weight unit.Pound, serviceArea string, isPPM bool) (unit.Cents, services.PricingDisplayParams, error) {
	// Validate parameters
	if len(contractCode) == 0 {
		return 0, nil, errors.New("ContractCode is required")
	}
	if referenceDate.IsZero() {
		return 0, nil, errors.New("ReferenceDate is required")
	}
	if !isPPM && distance < dlhPricerMinimumDistance {
		return 0, nil, fmt.Errorf("Distance must be at least %d", dlhPricerMinimumDistance)
	}
	if !isPPM && weight < dlhPricerMinimumWeight {
		return 0, nil, fmt.Errorf("Weight must be at least %d", dlhPricerMinimumWeight)
	}
	if len(serviceArea) == 0 {
		return 0, nil, errors.New("ServiceArea is required")
	}

	var domesticLinehaulPrice models.ReDomesticLinehaulPrice
	var err error
	isPeakPeriod := IsPeakPeriod(referenceDate)
	if isPPM && weight < dlhPricerMinimumWeight {
		domesticLinehaulPrice, err = fetchDomesticLinehaulPrice(appCtx, contractCode, isPeakPeriod, distance, dlhPricerMinimumWeight, serviceArea)
		if err != nil {
			return unit.Cents(0), nil, fmt.Errorf("could not fetch domestic linehaul rate: %w", err)
		}
	} else {
		domesticLinehaulPrice, err = fetchDomesticLinehaulPrice(appCtx, contractCode, isPeakPeriod, distance, weight, serviceArea)
		if err != nil {
			return unit.Cents(0), nil, fmt.Errorf("could not fetch domestic linehaul rate: %w", err)
		}
	}

	contractYear, err := fetchContractYear(appCtx, domesticLinehaulPrice.ContractID, referenceDate)
	if err != nil {
		return 0, nil, fmt.Errorf("Could not lookup contract year: %w", err)
	}

	baseTotalPrice := weight.ToCWTFloat64() * distance.Float64() * domesticLinehaulPrice.PriceMillicents.Float64()
	escalatedTotalPrice := contractYear.EscalationCompounded * baseTotalPrice

	totalPriceMillicents := unit.Millicents(escalatedTotalPrice)
	totalPriceCents := totalPriceMillicents.ToCents()

	params := services.PricingDisplayParams{
		{Key: models.ServiceItemParamNameContractYearName, Value: contractYear.Name},
		{Key: models.ServiceItemParamNameEscalationCompounded, Value: FormatEscalation(contractYear.EscalationCompounded)},
		{Key: models.ServiceItemParamNameIsPeak, Value: FormatBool(isPeakPeriod)},
		{Key: models.ServiceItemParamNamePriceRateOrFactor, Value: FormatFloat(domesticLinehaulPrice.PriceMillicents.ToDollarFloatNoRound(), 3)},
	}

	// if isPPM && weight < dlhPricerMinimumWeight {
	// 	weightFactor := float64(weight) / float64(dlhPricerMinimumWeight)
	// 	cost := float64(weightFactor) * float64(totalPriceCents)
	// 	cost := weightFactor.Float64()
	// 	fmt.Printf("==================== %v: %v : %v : %v ===============", totalPriceCents, cost, weight, weightFactor)
	// 	return totalPriceCents, params, nil
	// }

	return totalPriceCents, params, nil
}

// PriceUsingParams determines the price for a domestic linehaul given PaymentServiceItemParams
func (p domesticLinehaulPricer) PriceUsingParams(appCtx appcontext.AppContext, params models.PaymentServiceItemParams) (unit.Cents, services.PricingDisplayParams, error) {
	contractCode, err := getParamString(params, models.ServiceItemParamNameContractCode)
	if err != nil {
		return unit.Cents(0), nil, err
	}

	distanceZip3, err := getParamInt(params, models.ServiceItemParamNameDistanceZip3)
	if err != nil {
		return unit.Cents(0), nil, err
	}

	referenceDate, err := getParamTime(params, models.ServiceItemParamNameReferenceDate)
	if err != nil {
		return unit.Cents(0), nil, err
	}

	serviceAreaOrigin, err := getParamString(params, models.ServiceItemParamNameServiceAreaOrigin)
	if err != nil {
		return unit.Cents(0), nil, err
	}

	weightBilled, err := getParamInt(params, models.ServiceItemParamNameWeightBilled)
	if err != nil {
		return unit.Cents(0), nil, err
	}

	var isPPM = false
	if params[0].PaymentServiceItem.MTOServiceItem.MTOShipment.ShipmentType == models.MTOShipmentTypePPM {
		// PPMs do not require minimums for a shipment's weight or distance
		// this flag is passed into the Price function to ensure the weight and distance mins
		// are not enforced for PPMs
		isPPM = true
	}

	return p.Price(appCtx, contractCode, referenceDate, unit.Miles(distanceZip3), unit.Pound(weightBilled), serviceAreaOrigin, isPPM)
}

func fetchDomesticLinehaulPrice(appCtx appcontext.AppContext, contractCode string, isPeakPeriod bool, distance unit.Miles, weight unit.Pound, serviceArea string) (models.ReDomesticLinehaulPrice, error) {
	var domesticLinehaulPrice models.ReDomesticLinehaulPrice
	err := appCtx.DB().Q().
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
