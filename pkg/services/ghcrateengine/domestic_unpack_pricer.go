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

type domesticUnpackPricer struct {
}

// NewDomesticUnpackPricer creates a new pricer for domestic pack services
func NewDomesticUnpackPricer() services.DomesticUnpackPricer {
	return &domesticUnpackPricer{}
}

// Price determines the price for a domestic pack/unpack service
func (p domesticUnpackPricer) Price(appCtx appcontext.AppContext, contractCode string, requestedPickupDate time.Time, weight unit.Pound, servicesScheduleDest int) (unit.Cents, services.PricingDisplayParams, error) {
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
	if servicesScheduleDest == 0 {
		return 0, nil, errors.New("Service schedule is required")
	}

	isPeakPeriod := IsPeakPeriod(requestedPickupDate)
	var contractYear models.ReContractYear
	domOtherPrice, err := fetchDomOtherPrice(appCtx, contractCode, models.ReServiceCodeDUPK, servicesScheduleDest, isPeakPeriod)

	if err != nil {
		return 0, nil, fmt.Errorf("Could not lookup Domestic Other Price: %w", err)
	}

	err = appCtx.DB().Where("contract_id = $1", domOtherPrice.ContractID).
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
			Value: FormatEscalation(contractYear.EscalationCompounded),
		},
	}

	return totalCost, displayParams, nil
}

// PriceUsingParams determines the price for a domestic pack given PaymentServiceItemParams
func (p domesticUnpackPricer) PriceUsingParams(appCtx appcontext.AppContext, params models.PaymentServiceItemParams) (unit.Cents, services.PricingDisplayParams, error) {
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

	servicesScheduleDest, err := getParamInt(params, models.ServiceItemParamNameServicesScheduleDest)
	if err != nil {
		return unit.Cents(0), nil, err
	}

	return p.Price(appCtx, contractCode, requestedPickupDate, unit.Pound(weightBilled), servicesScheduleDest)
}
