package ghcrateengine

import (
	"fmt"
	"math"
	"time"

	"github.com/pkg/errors"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/featureflag"
	"github.com/transcom/mymove/pkg/unit"
)

const dlhPricerMinimumWeight = unit.Pound(500)

type domesticLinehaulPricer struct {
}

// NewDomesticLinehaulPricer creates a new pricer for domestic linehaul services
func NewDomesticLinehaulPricer() services.DomesticLinehaulPricer {
	return &domesticLinehaulPricer{}
}

// Price determines the price for a domestic linehaul
func (p domesticLinehaulPricer) Price(appCtx appcontext.AppContext, contractCode string, referenceDate time.Time, distance unit.Miles, weight unit.Pound, serviceArea string, isPPM bool, isMobileHome bool) (unit.Cents, services.PricingDisplayParams, error) {
	// Validate parameters
	if len(contractCode) == 0 {
		return 0, nil, errors.New("ContractCode is required")
	}
	if referenceDate.IsZero() {
		return 0, nil, errors.New("ReferenceDate is required")
	}
	if !isPPM && weight < dlhPricerMinimumWeight {
		return 0, nil, fmt.Errorf("Weight must be at least %d", dlhPricerMinimumWeight)
	}
	if len(serviceArea) == 0 {
		return 0, nil, errors.New("ServiceArea is required")
	}

	isPeakPeriod := IsPeakPeriod(referenceDate)
	finalWeight := weight

	if isPPM && weight < dlhPricerMinimumWeight {
		finalWeight = dlhPricerMinimumWeight
	}

	domesticLinehaulPrice, err := fetchDomesticLinehaulPrice(appCtx, contractCode, isPeakPeriod, distance, finalWeight, serviceArea)
	if err != nil {
		return unit.Cents(0), nil, fmt.Errorf("could not fetch domestic linehaul rate: %w", err)
	}

	basePrice := domesticLinehaulPrice.PriceMillicents.Float64() / 1000
	escalatedPrice, contractYear, err := escalatePriceForContractYear(
		appCtx,
		domesticLinehaulPrice.ContractID,
		referenceDate,
		true,
		basePrice)
	if err != nil {
		return 0, nil, fmt.Errorf("could not calculate escalated price: %w", err)
	}

	var params services.PricingDisplayParams

	if isMobileHome { // Need to apply mobile home factor to calculation
		mobileHomeFactorRow, err := fetchShipmentTypePrice(appCtx, contractCode, models.ReServiceCodeDMHF, models.MarketConus)
		if err != nil {
			return 0, nil, fmt.Errorf("could not fetch mobile home factor from database: %w", err)
		}

		escalatedPrice = roundToPrecision(escalatedPrice*mobileHomeFactorRow.Factor, 3)
		params = services.PricingDisplayParams{
			{Key: models.ServiceItemParamNameContractYearName, Value: contractYear.Name},
			{Key: models.ServiceItemParamNameEscalationCompounded, Value: FormatEscalation(contractYear.EscalationCompounded)},
			{Key: models.ServiceItemParamNameIsPeak, Value: FormatBool(isPeakPeriod)},
			{Key: models.ServiceItemParamNamePriceRateOrFactor, Value: FormatFloat(domesticLinehaulPrice.PriceMillicents.ToDollarFloatNoRound(), 3)},
			{Key: models.ServiceItemParamNameMobileHomeFactor, Value: FormatFloat(mobileHomeFactorRow.Factor, 3)},
		}
	} else { // Return display params without the mobile home factor
		params = services.PricingDisplayParams{
			{Key: models.ServiceItemParamNameContractYearName, Value: contractYear.Name},
			{Key: models.ServiceItemParamNameEscalationCompounded, Value: FormatEscalation(contractYear.EscalationCompounded)},
			{Key: models.ServiceItemParamNameIsPeak, Value: FormatBool(isPeakPeriod)},
			{Key: models.ServiceItemParamNamePriceRateOrFactor, Value: FormatFloat(domesticLinehaulPrice.PriceMillicents.ToDollarFloatNoRound(), 3)},
		}
	}

	totalPrice := finalWeight.ToCWTFloat64() * distance.Float64() * escalatedPrice

	totalPriceCents := unit.Cents(math.Round(totalPrice))

	if isPPM && weight < dlhPricerMinimumWeight {
		weightFactor := float64(weight) / float64(dlhPricerMinimumWeight)
		cost := float64(weightFactor) * float64(totalPriceCents)
		return unit.Cents(cost), params, nil
	}

	return totalPriceCents, params, nil
}

// PriceUsingParams determines the price for a domestic linehaul given PaymentServiceItemParams
func (p domesticLinehaulPricer) PriceUsingParams(appCtx appcontext.AppContext, params models.PaymentServiceItemParams, featureFlagValues map[string]bool) (unit.Cents, services.PricingDisplayParams, error) {
	contractCode, err := getParamString(params, models.ServiceItemParamNameContractCode)
	if err != nil {
		return unit.Cents(0), nil, err
	}

	distance, err := getParamInt(params, models.ServiceItemParamNameDistanceZip)
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

	// Check if mobile home is enabled and check for shipment type
	isMobileHome := false
	if featureFlagValues[featureflag.DomesticMobileHome] && params[0].PaymentServiceItem.MTOServiceItem.MTOShipment.ShipmentType == models.MTOShipmentTypeMobileHome {
		isMobileHome = true
	}

	return p.Price(appCtx, contractCode, referenceDate, unit.Miles(distance), unit.Pound(weightBilled), serviceAreaOrigin, isPPM, isMobileHome)
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
