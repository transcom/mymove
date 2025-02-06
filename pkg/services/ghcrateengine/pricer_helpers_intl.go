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

func priceInternationalShuttling(appCtx appcontext.AppContext, shuttlingCode models.ReServiceCode, contractCode string, referenceDate time.Time, weight unit.Pound, market models.Market) (unit.Cents, services.PricingDisplayParams, error) {
	if shuttlingCode != models.ReServiceCodeIOSHUT && shuttlingCode != models.ReServiceCodeIDSHUT {
		return 0, nil, fmt.Errorf("unsupported international shuttling code of %s", shuttlingCode)
	}
	// Validate parameters
	if len(contractCode) == 0 {
		return 0, nil, errors.New("ContractCode is required")
	}
	if referenceDate.IsZero() {
		return 0, nil, errors.New("ReferenceDate is required")
	}
	if weight < minInternationalWeight {
		return 0, nil, fmt.Errorf("Weight must be a minimum of %d", minInternationalWeight)
	}
	if market == "" {
		return 0, nil, errors.New("Market is required")
	}

	// look up rate for international accessorial price
	internationalAccessorialPrice, err := fetchInternationalAccessorialPrice(appCtx, contractCode, shuttlingCode, market)
	if err != nil {
		return 0, nil, fmt.Errorf("could not lookup International Accessorial Area Price: %w", err)
	}

	basePrice := internationalAccessorialPrice.PerUnitCents.Float64()
	escalatedPrice, contractYear, err := escalatePriceForContractYear(appCtx, internationalAccessorialPrice.ContractID, referenceDate, false, basePrice)
	if err != nil {
		return 0, nil, fmt.Errorf("could not calculate escalated price: %w", err)
	}

	escalatedPrice = escalatedPrice * weight.ToCWTFloat64()
	totalCost := unit.Cents(math.Round(escalatedPrice))

	params := services.PricingDisplayParams{
		{
			Key:   models.ServiceItemParamNamePriceRateOrFactor,
			Value: FormatCents(internationalAccessorialPrice.PerUnitCents),
		},
		{
			Key:   models.ServiceItemParamNameContractYearName,
			Value: contractYear.Name,
		},
		{
			Key:   models.ServiceItemParamNameEscalationCompounded,
			Value: FormatEscalation(contractYear.EscalationCompounded),
		},
	}
	return totalCost, params, nil
}

func priceIntlPackUnpack(appCtx appcontext.AppContext, packUnpackCode models.ReServiceCode, contractCode string, referenceDate time.Time, weight unit.Pound, perUnitCents int) (unit.Cents, services.PricingDisplayParams, error) {
	if packUnpackCode != models.ReServiceCodeIHPK && packUnpackCode != models.ReServiceCodeIHUPK {
		return 0, nil, fmt.Errorf("unsupported pack/unpack code of %s", packUnpackCode)
	}
	if len(contractCode) == 0 {
		return 0, nil, errors.New("ContractCode is required")
	}
	if referenceDate.IsZero() {
		return 0, nil, errors.New("ReferenceDate is required")
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
	escalatedPrice, contractYear, err := escalatePriceForContractYear(appCtx, contract.ID, referenceDate, false, basePrice)
	if err != nil {
		return 0, nil, fmt.Errorf("could not calculate escalated price: %w", err)
	}

	escalatedPrice = escalatedPrice * weight.ToCWTFloat64()
	totalCost := unit.Cents(math.Round(escalatedPrice))

	displayParams := services.PricingDisplayParams{
		{
			Key:   models.ServiceItemParamNameContractYearName,
			Value: contractYear.Name,
		},
		{
			Key:   models.ServiceItemParamNamePriceRateOrFactor,
			Value: FormatCents(unit.Cents(perUnitCents)),
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

func priceIntlFirstDaySIT(appCtx appcontext.AppContext, firstDaySITCode models.ReServiceCode, contractCode string, referenceDate time.Time, weight unit.Pound, perUnitCents int) (unit.Cents, services.PricingDisplayParams, error) {
	if firstDaySITCode != models.ReServiceCodeIOFSIT && firstDaySITCode != models.ReServiceCodeIDFSIT {
		return 0, nil, fmt.Errorf("unsupported first day SIT code of %s", firstDaySITCode)
	}
	if len(contractCode) == 0 {
		return 0, nil, errors.New("ContractCode is required")
	}
	if referenceDate.IsZero() {
		return 0, nil, errors.New("ReferenceDate is required")
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
	escalatedPrice, contractYear, err := escalatePriceForContractYear(appCtx, contract.ID, referenceDate, false, basePrice)
	if err != nil {
		return 0, nil, fmt.Errorf("could not calculate escalated price: %w", err)
	}

	escalatedPrice = escalatedPrice * weight.ToCWTFloat64()
	totalCost := unit.Cents(math.Round(escalatedPrice))

	displayParams := services.PricingDisplayParams{
		{
			Key:   models.ServiceItemParamNameContractYearName,
			Value: contractYear.Name,
		},
		{
			Key:   models.ServiceItemParamNamePriceRateOrFactor,
			Value: FormatCents(unit.Cents(perUnitCents)),
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

func priceIntlAdditionalDaySIT(appCtx appcontext.AppContext, additionalDaySITCode models.ReServiceCode, contractCode string, referenceDate time.Time, numberOfDaysInSIT int, weight unit.Pound, perUnitCents int) (unit.Cents, services.PricingDisplayParams, error) {
	if additionalDaySITCode != models.ReServiceCodeIOASIT && additionalDaySITCode != models.ReServiceCodeIDASIT {
		return 0, nil, fmt.Errorf("unsupported additional day SIT code of %s", additionalDaySITCode)
	}
	if len(contractCode) == 0 {
		return 0, nil, errors.New("ContractCode is required")
	}
	if referenceDate.IsZero() {
		return 0, nil, errors.New("ReferenceDate is required")
	}
	if numberOfDaysInSIT == 0 {
		return 0, nil, errors.New("NumberDaysSIT is required")
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
	escalatedPrice, contractYear, err := escalatePriceForContractYear(appCtx, contract.ID, referenceDate, false, basePrice)
	if err != nil {
		return 0, nil, fmt.Errorf("could not calculate escalated price: %w", err)
	}

	escalatedPrice = escalatedPrice * weight.ToCWTFloat64()
	totalForNumberOfDaysPrice := escalatedPrice * float64(numberOfDaysInSIT)
	totalCost := unit.Cents(math.Round(totalForNumberOfDaysPrice))

	displayParams := services.PricingDisplayParams{
		{
			Key:   models.ServiceItemParamNameContractYearName,
			Value: contractYear.Name,
		},
		{
			Key:   models.ServiceItemParamNamePriceRateOrFactor,
			Value: FormatCents(unit.Cents(perUnitCents)),
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
