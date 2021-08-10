package ghcrateengine

import (
	"fmt"
	"math"
	"strconv"
	"time"

	"github.com/pkg/errors"

	"github.com/transcom/mymove/pkg/appconfig"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/unit"
)

// DomesticShorthaulPricer is a service object to price domestic shorthaul
type domesticShorthaulPricer struct {
}

// NewDomesticShorthaulPricer is the public constructor for a DomesticRateAreaPricer using Pop
func NewDomesticShorthaulPricer() services.DomesticShorthaulPricer {
	return &domesticShorthaulPricer{}
}

// Price determines the price for a counseling service
func (p domesticShorthaulPricer) Price(appCfg appconfig.AppConfig, contractCode string,
	requestedPickupDate time.Time,
	distance unit.Miles,
	weight unit.Pound,
	serviceArea string) (totalCost unit.Cents, params services.PricingDisplayParams, err error) {
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
	if distance <= 0 {
		return 0, nil, errors.New("Distance must be greater than 0")
	}
	if len(serviceArea) == 0 {
		return 0, nil, errors.New("ServiceArea is required")
	}

	isPeakPeriod := IsPeakPeriod(requestedPickupDate)

	// look up rate for shorthaul
	domServiceAreaPrice, err := fetchDomServiceAreaPrice(appCfg, contractCode, models.ReServiceCodeDSH, serviceArea, isPeakPeriod)
	if err != nil {
		return 0, nil, fmt.Errorf("Could not lookup Domestic Service Area Price: %w", err)
	}

	contractYear, err := fetchContractYear(appCfg, domServiceAreaPrice.ContractID, requestedPickupDate)
	if err != nil {
		return 0, nil, fmt.Errorf("Could not lookup contract year: %w", err)
	}

	basePrice := domServiceAreaPrice.PriceCents.Float64() * distance.Float64() * weight.ToCWTFloat64()
	escalatedPrice := basePrice * contractYear.EscalationCompounded
	totalCost = unit.Cents(math.Round(escalatedPrice))

	var pricingRateEngineParams = services.PricingDisplayParams{
		{
			Key:   models.ServiceItemParamNameContractYearName,
			Value: contractYear.Name,
		},
		{
			Key:   models.ServiceItemParamNamePriceRateOrFactor,
			Value: FormatCents(domServiceAreaPrice.PriceCents),
		},
		{
			Key:   models.ServiceItemParamNameIsPeak,
			Value: strconv.FormatBool(isPeakPeriod),
		},
		{
			Key:   models.ServiceItemParamNameEscalationCompounded,
			Value: FormatEscalation(contractYear.EscalationCompounded),
		},
	}
	return totalCost, pricingRateEngineParams, nil
}

func (p domesticShorthaulPricer) PriceUsingParams(appCfg appconfig.AppConfig, params models.PaymentServiceItemParams) (unit.Cents, services.PricingDisplayParams, error) {
	contractCode, err := getParamString(params, models.ServiceItemParamNameContractCode)
	if err != nil {
		return unit.Cents(0), nil, err
	}

	requestedPickupDate, err := getParamTime(params, models.ServiceItemParamNameRequestedPickupDate)
	if err != nil {
		return unit.Cents(0), nil, err
	}

	distanceZip5, err := getParamInt(params, models.ServiceItemParamNameDistanceZip5)
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

	return p.Price(appCfg, contractCode, requestedPickupDate, unit.Miles(distanceZip5), unit.Pound(weightBilledActual), serviceAreaOrigin)
}
