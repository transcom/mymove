package ghcrateengine

import (
	"fmt"
	"math"
	"time"

	"github.com/pkg/errors"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/unit"
)

type domesticDestinationPricer struct {
}

// NewDomesticDestinationPricer instantiates a new pricer
func NewDomesticDestinationPricer() services.DomesticDestinationPricer {
	return &domesticDestinationPricer{}
}

// Price determines the price for the destination service area
func (p domesticDestinationPricer) Price(appCtx appcontext.AppContext, contractCode string, requestedPickupDate time.Time, weight unit.Pound, serviceArea string) (unit.Cents, services.PricingDisplayParams, error) {
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
	if len(serviceArea) == 0 {
		return 0, nil, errors.New("ServiceArea is required")
	}

	isPeakPeriod := IsPeakPeriod(requestedPickupDate)

	// look up rate for domestic destination price
	domServiceAreaPrice, err := fetchDomServiceAreaPrice(appCtx, contractCode, models.ReServiceCodeDDP, serviceArea, isPeakPeriod)
	if err != nil {
		return 0, nil, fmt.Errorf("Could not lookup Domestic Service Area Price: %w", err)
	}

	contractYear, err := fetchContractYear(appCtx, domServiceAreaPrice.ContractID, requestedPickupDate)
	if err != nil {
		return 0, nil, fmt.Errorf("Could not lookup contract year: %w", err)
	}

	basePrice := domServiceAreaPrice.PriceCents.Float64() * weight.ToCWTFloat64()
	escalatedPrice := basePrice * contractYear.EscalationCompounded
	totalCost := unit.Cents(math.Round(escalatedPrice))

	pricingParams := services.PricingDisplayParams{
		{Key: models.ServiceItemParamNamePriceRateOrFactor, Value: FormatCents(domServiceAreaPrice.PriceCents)},
		{Key: models.ServiceItemParamNameContractYearName, Value: contractYear.Name},
		{Key: models.ServiceItemParamNameIsPeak, Value: FormatBool(isPeakPeriod)},
		{Key: models.ServiceItemParamNameEscalationCompounded, Value: FormatEscalation(contractYear.EscalationCompounded)},
	}

	return totalCost, pricingParams, nil
}

func (p domesticDestinationPricer) PriceUsingParams(appCtx appcontext.AppContext, params models.PaymentServiceItemParams) (unit.Cents, services.PricingDisplayParams, error) {
	contractCode, err := getParamString(params, models.ServiceItemParamNameContractCode)
	if err != nil {
		return unit.Cents(0), nil, err
	}

	requestedPickupDate, err := getParamTime(params, models.ServiceItemParamNameRequestedPickupDate)
	if err != nil {
		return unit.Cents(0), nil, err
	}

	weightBilled, err := getParamInt(params, models.ServiceItemParamNameWeightBilled)
	if err != nil {
		return unit.Cents(0), nil, err
	}

	serviceAreaDest, err := getParamString(params, models.ServiceItemParamNameServiceAreaDest)
	if err != nil {
		return unit.Cents(0), nil, err
	}

	return p.Price(appCtx, contractCode, requestedPickupDate, unit.Pound(weightBilled), serviceAreaDest)
}
