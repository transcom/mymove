package ghcrateengine

import (
	"fmt"
	"math"
	"time"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"

	"github.com/pkg/errors"

	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/unit"
)

type domesticOriginPricer struct {
}

// NewDomesticOriginPricer creates a new pricer for domestic origin services
func NewDomesticOriginPricer() services.DomesticOriginPricer {
	return &domesticOriginPricer{}
}

// Price determines the price for a domestic origin
func (p domesticOriginPricer) Price(appCtx appcontext.AppContext, contractCode string, referenceDate time.Time, weight unit.Pound, serviceArea string, isPPM bool) (unit.Cents, services.PricingDisplayParams, error) {
	// Validate parameters
	if len(contractCode) == 0 {
		return 0, nil, errors.New("ContractCode is required")
	}
	if referenceDate.IsZero() {
		return 0, nil, errors.New("ReferenceDate is required")
	}
	if !isPPM && weight < minDomesticWeight {
		return 0, nil, fmt.Errorf("Weight must be a minimum of %d", minDomesticWeight)
	}
	if len(serviceArea) == 0 {
		return 0, nil, errors.New("ServiceArea is required")
	}

	isPeakPeriod := IsPeakPeriod(referenceDate)

	// look up rate for domestic origin price
	domServiceAreaPrice, err := fetchDomServiceAreaPrice(appCtx, contractCode, models.ReServiceCodeDOP, serviceArea, isPeakPeriod)
	if err != nil {
		return 0, nil, fmt.Errorf("Could not lookup Domestic Service Area Price: %w", err)
	}

	contractYear, err := fetchContractYear(appCtx, domServiceAreaPrice.ContractID, referenceDate)
	if err != nil {
		return 0, nil, fmt.Errorf("Could not lookup contract year: %w", err)
	}

	finalWeight := weight
	if isPPM && weight < dlhPricerMinimumWeight {
		finalWeight = dlhPricerMinimumWeight
	}

	basePrice := domServiceAreaPrice.PriceCents.Float64() * finalWeight.ToCWTFloat64()
	escalatedPrice := basePrice * contractYear.EscalationCompounded
	totalCost := unit.Cents(math.Round(escalatedPrice))

	params := services.PricingDisplayParams{
		{
			Key:   models.ServiceItemParamNamePriceRateOrFactor,
			Value: FormatCents(domServiceAreaPrice.PriceCents),
		},
		{
			Key:   models.ServiceItemParamNameContractYearName,
			Value: contractYear.Name,
		},
		{
			Key:   models.ServiceItemParamNameIsPeak,
			Value: FormatBool(isPeakPeriod),
		},
		{
			Key:   models.ServiceItemParamNameEscalationCompounded,
			Value: FormatEscalation(contractYear.EscalationCompounded),
		},
	}

	if isPPM && weight < dlhPricerMinimumWeight {
		weightFactor := float64(weight) / float64(dlhPricerMinimumWeight)
		cost := float64(weightFactor) * float64(totalCost)
		return unit.Cents(cost), params, nil
	}

	return totalCost, params, nil
}

// PriceUsingParams determines the price for a domestic origin given PaymentServiceItemParams
func (p domesticOriginPricer) PriceUsingParams(appCtx appcontext.AppContext, params models.PaymentServiceItemParams) (unit.Cents, services.PricingDisplayParams, error) {
	contractCode, err := getParamString(params, models.ServiceItemParamNameContractCode)
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
		// PPMs do not require minimums for a shipment's weight
		// this flag is passed into the Price function to ensure the weight min
		// are not enforced for PPMs
		isPPM = true
	}

	return p.Price(appCtx, contractCode, referenceDate, unit.Pound(weightBilled), serviceAreaOrigin, isPPM)
}
