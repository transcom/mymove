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

func priceIntlCratingUncrating(appCtx appcontext.AppContext, cratingUncratingCode models.ReServiceCode, contractCode string, referenceDate time.Time, billedCubicFeet unit.CubicFeet, standaloneCrate bool, standaloneCrateCap unit.Cents, externalCrate bool, market models.Market) (unit.Cents, services.PricingDisplayParams, error) {
	if cratingUncratingCode != models.ReServiceCodeICRT && cratingUncratingCode != models.ReServiceCodeIUCRT {
		return 0, nil, fmt.Errorf("unsupported international crating/uncrating code of %s", cratingUncratingCode)
	}

	// Validate parameters
	if len(contractCode) == 0 {
		return 0, nil, errors.New("ContractCode is required")
	}
	if referenceDate.IsZero() {
		return 0, nil, errors.New("ReferenceDate is required")
	}
	if market == "" {
		return 0, nil, errors.New("Market is required")
	}

	if externalCrate && billedCubicFeet < minIntlExternalCrateBilledCubicFeet {
		return 0, nil, fmt.Errorf("external crates must be billed for a minimum of %.2f cubic feet", minIntlExternalCrateBilledCubicFeet)
	}

	internationalAccessorialPrice, err := fetchInternationalAccessorialPrice(appCtx, contractCode, cratingUncratingCode, market)
	if err != nil {
		return 0, nil, fmt.Errorf("could not lookup International Accessorial Area Price: %w", err)
	}

	basePrice := internationalAccessorialPrice.PerUnitCents.Float64()
	escalatedPrice, contractYear, err := escalatePriceForContractYear(appCtx, internationalAccessorialPrice.ContractID, referenceDate, false, basePrice)
	if err != nil {
		return 0, nil, fmt.Errorf("could not calculate escalated price: %w", err)
	}

	escalatedPrice = escalatedPrice * float64(billedCubicFeet)
	totalCost := unit.Cents(math.Round(escalatedPrice))

	displayParams := services.PricingDisplayParams{
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
		{
			Key:   models.ServiceItemParamNameUncappedRequestTotal,
			Value: FormatCents(totalCost),
		},
	}

	if (standaloneCrate) && (totalCost > standaloneCrateCap) {
		totalCost = standaloneCrateCap
	}

	return totalCost, displayParams, nil
}

func priceIntlFuelSurchargeSIT(_ appcontext.AppContext, fuelSurchargeCode models.ReServiceCode, actualPickupDate time.Time, distance unit.Miles, weight unit.Pound, fscWeightBasedDistanceMultiplier float64, eiaFuelPrice unit.Millicents) (unit.Cents, services.PricingDisplayParams, error) {
	if fuelSurchargeCode != models.ReServiceCodeIOSFSC && fuelSurchargeCode != models.ReServiceCodeIDSFSC {
		return 0, nil, fmt.Errorf("unsupported international fuel surcharge code of %s", fuelSurchargeCode)
	}

	// Validate parameters
	if actualPickupDate.IsZero() {
		return 0, nil, errors.New("ActualPickupDate is required")
	}
	// zero represents pricing will not be calculated
	// this to handle when origin/destination addresses are OCONUS
	if distance < 0 {
		return 0, nil, errors.New("Distance cannot be less than 0")
	}
	if weight < minInternationalWeight {
		return 0, nil, fmt.Errorf("Weight must be a minimum of %d", minInternationalWeight)
	}
	if fscWeightBasedDistanceMultiplier == 0 {
		return 0, nil, errors.New("WeightBasedDistanceMultiplier is required")
	}
	if eiaFuelPrice == 0 {
		return 0, nil, errors.New("EIAFuelPrice is required")
	}

	fscPriceDifferenceInCents := (eiaFuelPrice - baseGHCDieselFuelPrice).Float64() / 1000.0
	fscMultiplier := fscWeightBasedDistanceMultiplier * distance.Float64()
	fscPrice := fscMultiplier * fscPriceDifferenceInCents * 100
	totalCost := unit.Cents(math.Round(fscPrice))

	displayParams := services.PricingDisplayParams{
		{Key: models.ServiceItemParamNameFSCPriceDifferenceInCents, Value: FormatFloat(fscPriceDifferenceInCents, 1)},
		{Key: models.ServiceItemParamNameFSCMultiplier, Value: FormatFloat(fscMultiplier, 7)},
	}

	return totalCost, displayParams, nil
}

func priceIntlPickupDeliverySIT(appCtx appcontext.AppContext, pickupDeliverySITCode models.ReServiceCode, contractCode string, referenceDate time.Time, weight unit.Pound, perUnitCents int) (unit.Cents, services.PricingDisplayParams, error) {
	if pickupDeliverySITCode != models.ReServiceCodeIOPSIT && pickupDeliverySITCode != models.ReServiceCodeIDDSIT {
		return 0, nil, fmt.Errorf("unsupported Intl PickupDeliverySIT code of %s", pickupDeliverySITCode)
	}

	// Validate parameters
	if len(contractCode) == 0 {
		return 0, nil, errors.New("ContractCode is required")
	}

	if referenceDate.IsZero() {
		return 0, nil, errors.New("ReferenceDate is required")
	}

	if weight < minInternationalWeight {
		return 0, nil, fmt.Errorf("weight must be a minimum of %d", minInternationalWeight)
	}

	if perUnitCents == 0 {
		return 0, nil, errors.New("perUnitCents is required")
	}

	isPeakPeriod := IsPeakPeriod(referenceDate)

	var reContract models.ReContract
	err := appCtx.DB().Where("re_contracts.code = ?", contractCode).First(&reContract)
	if err != nil {
		return 0, nil, fmt.Errorf("could not retrieve contract by code: %w", err)
	}

	escalatedPrice, contractYear, err := escalatePriceForContractYear(appCtx, reContract.ID, referenceDate, false, float64(perUnitCents))
	if err != nil {
		return 0, nil, fmt.Errorf("could not calculate escalated price: %w", err)
	}
	escalatedPrice = escalatedPrice * weight.ToCWTFloat64()

	totalPriceCents := unit.Cents(math.Round(escalatedPrice))

	displayParams := services.PricingDisplayParams{
		{
			Key:   models.ServiceItemParamNamePriceRateOrFactor,
			Value: FormatCents(unit.Cents(perUnitCents)),
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

	return totalPriceCents, displayParams, nil
}
