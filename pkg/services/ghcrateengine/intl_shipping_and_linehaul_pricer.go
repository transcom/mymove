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

const islhPricerMinimumWeight = unit.Pound(500)

type intlShippingAndLinehaulPricer struct {
}

func NewIntlShippingAndLinehaulPricer() services.IntlShippingAndLinehaulPricer {
	return &intlShippingAndLinehaulPricer{}
}

func (p intlShippingAndLinehaulPricer) Price(appCtx appcontext.AppContext, contractCode string, referenceDate time.Time, distance unit.Miles, weight unit.Pound, perUnitCents int) (unit.Cents, services.PricingDisplayParams, error) {
	if len(contractCode) == 0 {
		return 0, nil, errors.New("ContractCode is required")
	}
	if referenceDate.IsZero() {
		return 0, nil, errors.New("referenceDate is required")
	}
	if weight < islhPricerMinimumWeight {
		return 0, nil, fmt.Errorf("weight must be at least %d", islhPricerMinimumWeight)
	}
	if perUnitCents == 0 {
		return 0, nil, errors.New("PerUnitCents is required")
	}

	isPeakPeriod := IsPeakPeriod(referenceDate)

	contract, err := fetchContractByContractCode(appCtx, contractCode)
	if err != nil {
		return 0, nil, fmt.Errorf("could not find contract with code: %s: %w", contractCode, err)
	}

	basePrice := float64(perUnitCents)
	escalatedPrice, contractYear, err := escalatePriceForContractYear(
		appCtx,
		contract.ID,
		referenceDate,
		false,
		basePrice)
	if err != nil {
		return 0, nil, fmt.Errorf("could not calculate escalated price: %w", err)
	}

	escalatedPrice = escalatedPrice * weight.ToCWTFloat64()
	totalPriceCents := unit.Cents(math.Round(escalatedPrice))

	params := services.PricingDisplayParams{
		{
			Key:   models.ServiceItemParamNameContractYearName,
			Value: contractYear.Name,
		},
		{
			Key:   models.ServiceItemParamNameEscalationCompounded,
			Value: FormatEscalation(contractYear.EscalationCompounded),
		},
		{
			Key:   models.ServiceItemParamNameIsPeak,
			Value: FormatBool(isPeakPeriod),
		},
		{
			Key:   models.ServiceItemParamNamePriceRateOrFactor,
			Value: FormatCents(unit.Cents(perUnitCents)),
		}}

	return totalPriceCents, params, nil
}

func (p intlShippingAndLinehaulPricer) PriceUsingParams(appCtx appcontext.AppContext, params models.PaymentServiceItemParams) (unit.Cents, services.PricingDisplayParams, error) {
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

	weightBilled, err := getParamInt(params, models.ServiceItemParamNameWeightBilled)
	if err != nil {
		return unit.Cents(0), nil, err
	}

	perUnitCents, err := getParamInt(params, models.ServiceItemParamNamePerUnitCents)
	if err != nil {
		return unit.Cents(0), nil, err
	}

	return p.Price(appCtx, contractCode, referenceDate, unit.Miles(distance), unit.Pound(weightBilled), perUnitCents)
}
