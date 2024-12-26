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

func priceIntlPackUnpack(appCtx appcontext.AppContext, packUnpackCode models.ReServiceCode, contractCode string, referenceDate time.Time, weight unit.Pound, servicesSchedule int, isPPM bool) (unit.Cents, services.PricingDisplayParams, error) {
	// Validate parameters
	var intlOtherPriceCode models.ReServiceCode
	switch packUnpackCode {
	case models.ReServiceCodeIHPK:
		intlOtherPriceCode = models.ReServiceCodeIHPK
	case models.ReServiceCodeIHUPK:
		intlOtherPriceCode = models.ReServiceCodeIHUPK
	default:
		return 0, nil, fmt.Errorf("unsupported pack/unpack code of %s", packUnpackCode)
	}
	if len(contractCode) == 0 {
		return 0, nil, errors.New("ContractCode is required")
	}
	if referenceDate.IsZero() {
		return 0, nil, errors.New("ReferenceDate is required")
	}
	if !isPPM && weight < minDomesticWeight {
		return 0, nil, fmt.Errorf("Weight must be a minimum of %d", minDomesticWeight)
	}
	if servicesSchedule == 0 {
		return 0, nil, errors.New("Services schedule is required")
	}

	isPeakPeriod := IsPeakPeriod(referenceDate)

	intlOtherPrice, err := fetchIntlOtherPrice(appCtx, contractCode, intlOtherPriceCode, servicesSchedule, isPeakPeriod)
	if err != nil {
		return 0, nil, fmt.Errorf("could not lookup domestic other price: %w", err)
	}

	finalWeight := weight
	if isPPM && weight < minDomesticWeight {
		finalWeight = minDomesticWeight
	}

	basePrice := intlOtherPrice.PerUnitCents.Float64()
	escalatedPrice, contractYear, err := escalatePriceForContractYear(appCtx, intlOtherPrice.ContractID, referenceDate, false, basePrice)
	if err != nil {
		return 0, nil, fmt.Errorf("could not calculate escalated price: %w", err)
	}

	escalatedPrice = escalatedPrice * finalWeight.ToCWTFloat64()

	displayParams := services.PricingDisplayParams{
		{
			Key:   models.ServiceItemParamNameContractYearName,
			Value: contractYear.Name,
		},
		{
			Key:   models.ServiceItemParamNamePriceRateOrFactor,
			Value: FormatCents(intlOtherPrice.PerUnitCents),
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

	totalCost := unit.Cents(math.Round(escalatedPrice))
	if isPPM && weight < minDomesticWeight {
		weightFactor := float64(weight) / float64(minDomesticWeight)
		cost := float64(weightFactor) * float64(totalCost)
		return unit.Cents(cost), displayParams, nil
	}
	return totalCost, displayParams, nil
}
