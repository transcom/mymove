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

type domesticDestinationPricer struct {
}

// NewDomesticDestinationPricer instantiates a new pricer
func NewDomesticDestinationPricer() services.DomesticDestinationPricer {
	return &domesticDestinationPricer{}
}

// Price determines the price for the destination service area
func (p domesticDestinationPricer) Price(appCtx appcontext.AppContext, contractCode string, referenceDate time.Time, weight unit.Pound, serviceArea string, isPPM bool, isMobileHome bool, featureFlagValues map[string]bool) (unit.Cents, services.PricingDisplayParams, error) {
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

	isFactorToggleOn := false // Track whether DMHF Factor FF toggle is on for this Pack or Unpack item
	if isMobileHome {         // Only check for mobile home factor FF if this is a mobile home shipment.
		if featureFlagValues == nil || len(featureFlagValues) <= 0 {
			return 0, nil, fmt.Errorf("Expected a map of feature flag values when checking pricing for DPK item, received nil or empty map instead.")
		}
		if featureFlagValues[featureflag.DomesticMobileHomeDDPFactor] {
			isFactorToggleOn = true
		}
	}

	isPeakPeriod := IsPeakPeriod(referenceDate)

	// look up rate for domestic destination price
	domServiceAreaPrice, err := fetchDomServiceAreaPrice(appCtx, contractCode, models.ReServiceCodeDDP, serviceArea, isPeakPeriod)
	if err != nil {
		return 0, nil, fmt.Errorf("could not lookup Domestic Service Area Price: %w", err)
	}

	finalWeight := weight
	if isPPM && weight < minDomesticWeight {
		finalWeight = minDomesticWeight
	}

	basePrice := domServiceAreaPrice.PriceCents.Float64()
	escalatedPrice, contractYear, err := escalatePriceForContractYear(appCtx, domServiceAreaPrice.ContractID, referenceDate, false, basePrice)
	if err != nil {
		return 0, nil, fmt.Errorf("could not calculate escalated price: %w", err)
	}

	totalCost := unit.Cents(0)
	var pricingParams services.PricingDisplayParams
	if isFactorToggleOn {
		mobileHomeFactorRow, err := fetchShipmentTypePrice(appCtx, contractCode, models.ReServiceCodeDMHF, models.MarketConus)
		if err != nil {
			return 0, nil, fmt.Errorf("could not fetch mobile home factor from database: %w", err)
		}
		escalatedPrice = roundToPrecision(escalatedPrice*mobileHomeFactorRow.Factor, 2)
		totalCost = unit.Cents(escalatedPrice * finalWeight.ToCWTFloat64())

		pricingParams = services.PricingDisplayParams{
			{Key: models.ServiceItemParamNameContractYearName, Value: contractYear.Name},
			{Key: models.ServiceItemParamNamePriceRateOrFactor, Value: FormatCents(domServiceAreaPrice.PriceCents)},
			{Key: models.ServiceItemParamNameIsPeak, Value: FormatBool(isPeakPeriod)},
			{Key: models.ServiceItemParamNameEscalationCompounded, Value: FormatEscalation(contractYear.EscalationCompounded)},
			{Key: models.ServiceItemParamNameMobileHomeFactor, Value: FormatFloat(mobileHomeFactorRow.Factor, 2)},
		}
	} else {
		escalatedPrice = escalatedPrice * finalWeight.ToCWTFloat64()
		totalCost = unit.Cents(math.Round(escalatedPrice))

		pricingParams = services.PricingDisplayParams{
			{Key: models.ServiceItemParamNamePriceRateOrFactor, Value: FormatCents(domServiceAreaPrice.PriceCents)},
			{Key: models.ServiceItemParamNameContractYearName, Value: contractYear.Name},
			{Key: models.ServiceItemParamNameIsPeak, Value: FormatBool(isPeakPeriod)},
			{Key: models.ServiceItemParamNameEscalationCompounded, Value: FormatEscalation(contractYear.EscalationCompounded)},
		}
	}
	if isPPM && weight < minDomesticWeight {
		weightFactor := float64(weight) / float64(minDomesticWeight)
		cost := float64(weightFactor) * float64(totalCost)
		return unit.Cents(cost), pricingParams, nil
	}

	return totalCost, pricingParams, nil
}

func (p domesticDestinationPricer) PriceUsingParams(appCtx appcontext.AppContext, params models.PaymentServiceItemParams, featureFlagValues map[string]bool) (unit.Cents, services.PricingDisplayParams, error) {
	contractCode, err := getParamString(params, models.ServiceItemParamNameContractCode)
	if err != nil {
		return unit.Cents(0), nil, err
	}

	referenceDate, err := getParamTime(params, models.ServiceItemParamNameReferenceDate)
	if err != nil {
		return unit.Cents(0), nil, err
	}

	serviceAreaDest, err := getParamString(params, models.ServiceItemParamNameServiceAreaDest)
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

	// Check if DDP service item has been enabled for Mobile Home shipments
	isMobileHome := false
	if featureFlagValues[featureflag.DomesticMobileHome] &&
		featureFlagValues[featureflag.DomesticMobileHomeDDPEnabled] &&
		params[0].PaymentServiceItem.MTOServiceItem.MTOShipment.ShipmentType == models.MTOShipmentTypeMobileHome {
		isMobileHome = true
	}

	return p.Price(appCtx, contractCode, referenceDate, unit.Pound(weightBilled), serviceAreaDest, isPPM, isMobileHome, featureFlagValues)
}
