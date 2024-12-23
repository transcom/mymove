package ghcrateengine

import (
	"fmt"
	"math"
	"slices"
	"strconv"
	"time"

	"github.com/gofrs/uuid"
	"github.com/pkg/errors"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/unit"
)

func priceDomesticPackUnpack(appCtx appcontext.AppContext, packUnpackCode models.ReServiceCode, contractCode string, referenceDate time.Time, weight unit.Pound, servicesSchedule int, isPPM bool) (unit.Cents, services.PricingDisplayParams, error) {
	// Validate parameters
	var domOtherPriceCode models.ReServiceCode
	switch packUnpackCode {
	case models.ReServiceCodeDPK, models.ReServiceCodeDNPK:
		domOtherPriceCode = models.ReServiceCodeDPK
	case models.ReServiceCodeDUPK:
		domOtherPriceCode = models.ReServiceCodeDUPK
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

	domOtherPrice, err := fetchDomOtherPrice(appCtx, contractCode, domOtherPriceCode, servicesSchedule, isPeakPeriod)
	if err != nil {
		return 0, nil, fmt.Errorf("could not lookup domestic other price: %w", err)
	}

	finalWeight := weight
	if isPPM && weight < minDomesticWeight {
		finalWeight = minDomesticWeight
	}

	basePrice := domOtherPrice.PriceCents.Float64()
	escalatedPrice, contractYear, err := escalatePriceForContractYear(appCtx, domOtherPrice.ContractID, referenceDate, false, basePrice)
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

	// Adjust for NTS packing factor if needed.
	if packUnpackCode == models.ReServiceCodeDNPK {
		shipmentTypePrice, err := fetchShipmentTypePrice(appCtx, contractCode, models.ReServiceCodeDNPK, models.MarketConus)
		if err != nil {
			return 0, nil, fmt.Errorf("could not lookup shipment type price: %w", err)
		}
		escalatedPrice = escalatedPrice * shipmentTypePrice.Factor

		displayParams = append(displayParams, services.PricingDisplayParam{
			Key:   models.ServiceItemParamNameNTSPackingFactor,
			Value: FormatFloat(shipmentTypePrice.Factor, 2),
		})
	}

	totalCost := unit.Cents(math.Round(escalatedPrice))
	if isPPM && weight < minDomesticWeight {
		weightFactor := float64(weight) / float64(minDomesticWeight)
		cost := float64(weightFactor) * float64(totalCost)
		return unit.Cents(cost), displayParams, nil
	}
	return totalCost, displayParams, nil
}

func priceDomesticFirstDaySIT(appCtx appcontext.AppContext, firstDaySITCode models.ReServiceCode, contractCode string, referenceDate time.Time, weight unit.Pound, serviceArea string, disableWeightMinimum bool) (unit.Cents, services.PricingDisplayParams, error) {
	var sitType string
	if firstDaySITCode == models.ReServiceCodeDDFSIT {
		sitType = "destination"
	} else if firstDaySITCode == models.ReServiceCodeDOFSIT {
		sitType = "origin"
	} else {
		return 0, nil, fmt.Errorf("unsupported first day sit code of %s", firstDaySITCode)
	}

	if !disableWeightMinimum && weight < minDomesticWeight {
		return 0, nil, fmt.Errorf("weight of %d less than the minimum of %d", weight, minDomesticWeight)
	}

	isPeakPeriod := IsPeakPeriod(referenceDate)
	serviceAreaPrice, err := fetchDomServiceAreaPrice(appCtx, contractCode, firstDaySITCode, serviceArea, isPeakPeriod)
	if err != nil {
		return unit.Cents(0), nil, fmt.Errorf("could not fetch domestic %s first day SIT rate: %w", sitType, err)
	}

	basePrice := serviceAreaPrice.PriceCents.Float64()
	escalatedPrice, contractYear, err := escalatePriceForContractYear(appCtx, serviceAreaPrice.ContractID, referenceDate, false, basePrice)
	if err != nil {
		return 0, nil, fmt.Errorf("could not calculate escalated price: %w", err)
	}

	escalatedPrice = escalatedPrice * weight.ToCWTFloat64()
	totalPriceCents := unit.Cents(math.Round(escalatedPrice))

	params := services.PricingDisplayParams{
		{Key: models.ServiceItemParamNameContractYearName, Value: contractYear.Name},
		{Key: models.ServiceItemParamNameEscalationCompounded, Value: FormatEscalation(contractYear.EscalationCompounded)},
		{Key: models.ServiceItemParamNameIsPeak, Value: FormatBool(isPeakPeriod)},
		{Key: models.ServiceItemParamNamePriceRateOrFactor, Value: FormatCents(serviceAreaPrice.PriceCents)},
	}

	return totalPriceCents, params, nil
}

func priceDomesticAdditionalDaysSIT(appCtx appcontext.AppContext, additionalDaySITCode models.ReServiceCode, contractCode string, referenceDate time.Time, weight unit.Pound, serviceArea string, numberOfDaysInSIT int, disableWeightMinimum bool) (unit.Cents, services.PricingDisplayParams, error) {
	var sitType string
	if additionalDaySITCode == models.ReServiceCodeDDASIT {
		sitType = "destination"
	} else if additionalDaySITCode == models.ReServiceCodeDOASIT {
		sitType = "origin"
	} else {
		return 0, nil, fmt.Errorf("unsupported additional day sit code of %s", additionalDaySITCode)
	}

	if !disableWeightMinimum && weight < minDomesticWeight {
		return 0, nil, fmt.Errorf("weight of %d less than the minimum of %d", weight, minDomesticWeight)
	}

	isPeakPeriod := IsPeakPeriod(referenceDate)
	serviceAreaPrice, err := fetchDomServiceAreaPrice(appCtx, contractCode, additionalDaySITCode, serviceArea, isPeakPeriod)
	if err != nil {
		return unit.Cents(0), nil, fmt.Errorf("could not fetch domestic %s additional days SIT rate: %w", sitType, err)
	}

	basePrice := serviceAreaPrice.PriceCents.Float64()
	escalatedPrice, contractYear, err := escalatePriceForContractYear(appCtx, serviceAreaPrice.ContractID, referenceDate, false, basePrice)
	if err != nil {
		return 0, nil, fmt.Errorf("could not calculate escalated price: %w", err)
	}

	escalatedPrice = escalatedPrice * weight.ToCWTFloat64()
	totalForNumberOfDaysPrice := escalatedPrice * float64(numberOfDaysInSIT)
	totalPriceCents := unit.Cents(math.Round(totalForNumberOfDaysPrice))

	displayParams := services.PricingDisplayParams{
		{
			Key:   models.ServiceItemParamNameContractYearName,
			Value: contractYear.Name,
		},
		{
			Key:   models.ServiceItemParamNamePriceRateOrFactor,
			Value: FormatCents(serviceAreaPrice.PriceCents),
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

func priceDomesticPickupDeliverySIT(appCtx appcontext.AppContext, pickupDeliverySITCode models.ReServiceCode, contractCode string, referenceDate time.Time, weight unit.Pound, serviceArea string, sitSchedule int, zipOriginal string, zipActual string, distance unit.Miles) (unit.Cents, services.PricingDisplayParams, error) {
	var sitType, sitModifier, zipOriginalName, zipActualName string
	if pickupDeliverySITCode == models.ReServiceCodeDDDSIT {
		sitType = "destination"
		sitModifier = "delivery"
		zipOriginalName = "SIT original destination"
		zipActualName = "SIT final destination"
	} else if pickupDeliverySITCode == models.ReServiceCodeDOPSIT {
		sitType = "origin"
		sitModifier = "pickup"
		zipOriginalName = "SIT origin original"
		zipActualName = "SIT origin actual"
	} else {
		return 0, nil, fmt.Errorf("unsupported pickup/delivery SIT code of %s", pickupDeliverySITCode)
	}

	if weight < minDomesticWeight {
		return 0, nil, fmt.Errorf("weight of %d less than the minimum of %d", weight, minDomesticWeight)
	}

	if len(zipOriginal) < 5 {
		return unit.Cents(0), nil, fmt.Errorf("invalid %s postal code of %s", zipOriginalName, zipOriginal)
	}
	zip3Original := zipOriginal[:3]

	if len(zipActual) < 5 {
		return unit.Cents(0), nil, fmt.Errorf("invalid %s postal code of %s", zipActualName, zipActual)
	}
	zip3Actual := zipActual[:3]

	// Three different pricing scenarios below.

	if distance <= 50 {
		isPeakPeriod := IsPeakPeriod(referenceDate)
		domOtherPrice, err := fetchDomOtherPrice(appCtx, contractCode, pickupDeliverySITCode, sitSchedule, isPeakPeriod)
		if err != nil {
			return unit.Cents(0), nil, fmt.Errorf("could not fetch domestic %s SIT %s rate: %w", sitType, sitModifier, err)
		}
		basePrice := domOtherPrice.PriceCents.Float64()

		escalatedPrice, contractYear, err := escalatePriceForContractYear(appCtx, domOtherPrice.ContractID, referenceDate, false, basePrice)
		if err != nil {
			return 0, nil, fmt.Errorf("could not calculate escalated price: %w", err)
		}
		escalatedPrice = escalatedPrice * weight.ToCWTFloat64()

		totalPriceCents := unit.Cents(math.Round(escalatedPrice))

		displayParams := services.PricingDisplayParams{
			{
				Key:   models.ServiceItemParamNamePriceRateOrFactor,
				Value: FormatCents(domOtherPrice.PriceCents),
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

	// 1) Zip3 to same zip3
	if zip3Original == zip3Actual {
		// Do a normal shorthaul calculation
		shorthaulPricer := NewDomesticShorthaulPricer()
		totalPriceCents, displayParams, err := shorthaulPricer.Price(appCtx, contractCode, referenceDate, distance, weight, serviceArea)
		if err != nil {
			return unit.Cents(0), nil, fmt.Errorf("could not price shorthaul: %w", err)
		}

		return totalPriceCents, displayParams, nil
	}

	// Zip3s must be different at this point.  Now examine distance.

	// 2) Zip3 to different zip3 and > 50 miles
	if distance > 50 {
		// Do a normal linehaul calculation
		linehaulPricer := NewDomesticLinehaulPricer()
		// TODO: This will need adjusting once SIT is implemented for PPMs
		isPPM := false
		totalPriceCents, displayParams, err := linehaulPricer.Price(appCtx, contractCode, referenceDate, distance, weight, serviceArea, isPPM)
		if err != nil {
			return unit.Cents(0), nil, fmt.Errorf("could not price linehaul: %w", err)
		}

		return totalPriceCents, displayParams, nil
	}

	// Zip3s must be different at this point and distance is <= 50.

	// 3) Zip3 to different zip3 and <= 50 miles

	// Rate comes from the domestic other price table based on SIT schedule
	isPeakPeriod := IsPeakPeriod(referenceDate)
	domOtherPrice, err := fetchDomOtherPrice(appCtx, contractCode, pickupDeliverySITCode, sitSchedule, isPeakPeriod)
	if err != nil {
		return unit.Cents(0), nil, fmt.Errorf("could not fetch domestic %s SIT %s rate: %w", sitType, sitModifier, err)
	}
	basePrice := domOtherPrice.PriceCents.Float64()

	escalatedPrice, contractYear, err := escalatePriceForContractYear(appCtx, domOtherPrice.ContractID, referenceDate, false, basePrice)
	if err != nil {
		return 0, nil, fmt.Errorf("could not calculate escalated price: %w", err)
	}
	escalatedPrice = escalatedPrice * weight.ToCWTFloat64()

	totalPriceCents := unit.Cents(math.Round(escalatedPrice))

	displayParams := services.PricingDisplayParams{
		{
			Key:   models.ServiceItemParamNamePriceRateOrFactor,
			Value: FormatCents(domOtherPrice.PriceCents),
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

func priceDomesticShuttling(appCtx appcontext.AppContext, shuttlingCode models.ReServiceCode, contractCode string, referenceDate time.Time, weight unit.Pound, serviceSchedule int) (unit.Cents, services.PricingDisplayParams, error) {
	if shuttlingCode != models.ReServiceCodeDOSHUT && shuttlingCode != models.ReServiceCodeDDSHUT {
		return 0, nil, fmt.Errorf("unsupported domestic shuttling code of %s", shuttlingCode)
	}
	// Validate parameters
	if len(contractCode) == 0 {
		return 0, nil, errors.New("ContractCode is required")
	}
	if referenceDate.IsZero() {
		return 0, nil, errors.New("ReferenceDate is required")
	}
	if weight < minDomesticWeight {
		return 0, nil, fmt.Errorf("Weight must be a minimum of %d", minDomesticWeight)
	}
	if serviceSchedule == 0 {
		return 0, nil, errors.New("Service schedule is required")
	}

	// look up rate for domestic accessorial price
	domAccessorialPrice, err := fetchAccessorialPrice(appCtx, contractCode, shuttlingCode, serviceSchedule)
	if err != nil {
		return 0, nil, fmt.Errorf("could not lookup Domestic Accessorial Area Price: %w", err)
	}

	basePrice := domAccessorialPrice.PerUnitCents.Float64()
	escalatedPrice, contractYear, err := escalatePriceForContractYear(appCtx, domAccessorialPrice.ContractID, referenceDate, false, basePrice)
	if err != nil {
		return 0, nil, fmt.Errorf("could not calculate escalated price: %w", err)
	}

	escalatedPrice = escalatedPrice * weight.ToCWTFloat64()
	totalCost := unit.Cents(math.Round(escalatedPrice))

	params := services.PricingDisplayParams{
		{
			Key:   models.ServiceItemParamNamePriceRateOrFactor,
			Value: FormatCents(domAccessorialPrice.PerUnitCents),
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
	if weight < minDomesticWeight {
		return 0, nil, fmt.Errorf("Weight must be a minimum of %d", minInternationalWeight)
	}
	if market == "" {
		return 0, nil, errors.New("Market is required")
	}

	// look up rate for international accessorial price
	internationalAccessorialPrice, err := fetchInternationalAccessorialPrice(appCtx, contractCode, shuttlingCode, market)
	if err != nil {
		return 0, nil, fmt.Errorf("could not lookup Interntional Accessorial Area Price: %w", err)
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

func priceDomesticCrating(appCtx appcontext.AppContext, code models.ReServiceCode, contractCode string, referenceDate time.Time, billedCubicFeet unit.CubicFeet, serviceSchedule int, standaloneCrate bool, standaloneCrateCap unit.Cents) (unit.Cents, services.PricingDisplayParams, error) {
	if code != models.ReServiceCodeDCRT && code != models.ReServiceCodeDUCRT {
		return 0, nil, fmt.Errorf("unsupported domestic crating code of %s", code)
	}

	if billedCubicFeet < 4.0 {
		return 0, nil, fmt.Errorf("crate must be billed for a minimum of 4 cubic feet")
	}
	domAccessorialPrice, err := fetchAccessorialPrice(appCtx, contractCode, code, serviceSchedule)
	if err != nil {
		return 0, nil, fmt.Errorf("could not lookup Domestic Accessorial Area Price: %w", err)
	}

	basePrice := domAccessorialPrice.PerUnitCents.Float64()
	escalatedPrice, contractYear, err := escalatePriceForContractYear(appCtx, domAccessorialPrice.ContractID, referenceDate, false, basePrice)
	if err != nil {
		return 0, nil, fmt.Errorf("could not calculate escalated price: %w", err)
	}

	escalatedPrice = escalatedPrice * float64(billedCubicFeet)
	totalCost := unit.Cents(math.Round(escalatedPrice))

	params := services.PricingDisplayParams{
		{
			Key:   models.ServiceItemParamNamePriceRateOrFactor,
			Value: FormatCents(domAccessorialPrice.PerUnitCents),
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

	return totalCost, params, nil
}

// createPricerGeneratedParams stores PaymentServiceItemParams, whose origin is the PRICER, into the database
// It also returns the newly created PaymentServiceItemParams.
func createPricerGeneratedParams(appCtx appcontext.AppContext, paymentServiceItemID uuid.UUID, params services.PricingDisplayParams) (models.PaymentServiceItemParams, error) {
	var paymentServiceItemParams models.PaymentServiceItemParams

	if len(params) == 0 {
		return paymentServiceItemParams, fmt.Errorf("PricingDisplayParams must not be empty")
	}

	for _, param := range params {

		// Find the paymentServiceItemParam associated with this PricingDisplayParam
		var serviceItemParamKey models.ServiceItemParamKey
		err := appCtx.DB().Q().
			Where("key = ?", param.Key).
			First(&serviceItemParamKey)
		if err != nil {
			return paymentServiceItemParams, fmt.Errorf("Unable to find service item param key for %v", param.Key)
		}
		if serviceItemParamKey.Origin != models.ServiceItemParamOriginPricer {
			return paymentServiceItemParams, fmt.Errorf("Service item param key is not a pricer param. Param key: %v", serviceItemParamKey.Key)
		}

		// Create the PaymentServiceItemParam from the PricingDisplayParam and store it in the DB
		newParam := models.PaymentServiceItemParam{
			PaymentServiceItemID:  paymentServiceItemID,
			ServiceItemParamKeyID: serviceItemParamKey.ID,
			ServiceItemParamKey:   serviceItemParamKey,
			Value:                 param.Value,
		}

		verrs, err := appCtx.DB().ValidateAndCreate(&newParam)
		if err != nil {
			return paymentServiceItemParams, fmt.Errorf("failure creating payment service item param: %w", err)
		} else if verrs.HasAny() {
			return paymentServiceItemParams, apperror.NewInvalidCreateInputError(verrs, "validation error with creating payment service item param")
			//nolint:revive //
		} else {
			// Append it to a slice of PaymentServiceItemParams to return
			paymentServiceItemParams = append(paymentServiceItemParams, newParam)
		}
	}
	return paymentServiceItemParams, nil
}

// escalatePriceForContractYear calculates the escalated price from the base price in cents, which is provided by the caller/pricer,
// and the escalation factor, which is provided by the contract year. The product is rounded to the nearest cent, or to the
// nearest tenth-cent for linehaul prices, before and after multiplication. The resulting price is returned in cents along
// with the contract year.
func escalatePriceForContractYear(appCtx appcontext.AppContext, contractID uuid.UUID, referenceDate time.Time, isLinehaul bool, basePriceCents float64) (float64, models.ReContractYear, error) {
	contractYear, err := fetchContractYear(appCtx, contractID, referenceDate)
	if err != nil {
		return 0, contractYear, fmt.Errorf("could not lookup contract year: %w", err)
	}

	escalatedPrice := basePriceCents

	// round escalated price to the nearest cent, or the nearest tenth-of-a-cent if linehaul
	precision := 0
	if isLinehaul {
		precision = 1
	}

	escalatedPrice = roundToPrecision(escalatedPrice, precision)

	if slices.Contains(models.ContractYears, contractYear.Name) {
		escalatedPrice, err = compoundEscalationFactors(appCtx, contractID, contractYear, escalatedPrice)
		if err != nil {
			return 0, contractYear, err
		}
	} else {
		escalatedPrice = escalatedPrice * contractYear.EscalationCompounded
	}

	escalatedPrice = roundToPrecision(escalatedPrice, precision)
	return escalatedPrice, contractYear, nil
}

func compoundEscalationFactors(appCtx appcontext.AppContext, contractID uuid.UUID, contractYear models.ReContractYear, escalatedPrice float64) (float64, error) {
	// Get all contracts based on contract Id
	contractYearsFromDB, err := fetchContractsByContractId(appCtx, contractID)
	if err != nil {
		return escalatedPrice, fmt.Errorf("could not lookup contracts by Id: %w", err)
	}

	// A contract may have Option Year 3 but it is not guaranteed. Need to know if it does or not
	contractsYearsFromDBMap := make(map[string]models.ReContractYear)
	for _, contract := range contractYearsFromDB {
		// Add re_contract_years record to map
		contractsYearsFromDBMap[contract.Name] = contract
	}

	// Get expectations for price escalations calculations
	expectations, err := models.GetExpectedEscalationPriceContractsCount(contractYear.Name)
	if err != nil {
		return escalatedPrice, err
	}

	// Adding contracts that are expected to be in the calculations based on the contract year to a map
	contractYearsForCalculation := make(map[string]models.ReContractYear)
	if expectations.ExpectedAmountOfAwardTermsForCalculation > 0 {
		contractYearsForCalculation, err = addContractsForEscalationCalculation(contractYearsForCalculation, contractsYearsFromDBMap, expectations.ExpectedAmountOfAwardTermsForCalculation, models.AwardTerm)
		if err != nil {
			return escalatedPrice, err
		}
	}
	if expectations.ExpectedAmountOfOptionPeriodYearsForCalculation > 0 {
		contractYearsForCalculation, err = addContractsForEscalationCalculation(contractYearsForCalculation, contractsYearsFromDBMap, expectations.ExpectedAmountOfOptionPeriodYearsForCalculation, models.OptionPeriod)
		if err != nil {
			return escalatedPrice, err
		}
	}
	if expectations.ExpectedAmountOfBasePeriodYearsForCalculation > 0 {
		contractYearsForCalculation, err = addContractsForEscalationCalculation(contractYearsForCalculation, contractsYearsFromDBMap, expectations.ExpectedAmountOfBasePeriodYearsForCalculation, models.BasePeriodYear)
		if err != nil {
			return escalatedPrice, err
		}
	}

	// Make sure the expected amount of contracts are being used in the escalated Price calculation
	if expectations.ExpectedAmountOfContractYearsForCalculation > 0 && len(contractYearsForCalculation) != expectations.ExpectedAmountOfContractYearsForCalculation {
		err := apperror.NewInternalServerError("Unexpected amount of contract years being used in escalated price calculation")
		return escalatedPrice, err
	}

	// Multiply the escalated price by each re_contract_years record escalation factor. EscalatedPrice = EscalatedPrice * ContractEscalationFactor
	var compoundedEscalatedPrice = escalatedPrice

	if expectations.ExpectedAmountOfContractYearsForCalculation > 0 {
		for _, contract := range contractYearsForCalculation {
			compoundedEscalatedPrice = compoundedEscalatedPrice * contract.Escalation
		}
	}

	return compoundedEscalatedPrice, nil
}

// roundToPrecision rounds a float64 value to the number of decimal points indicated by the precision.
// TODO: Future cleanup could involve moving this function to a math/utility package with some simple tests
func roundToPrecision(value float64, precision int) float64 {
	ratio := math.Pow(10, float64(precision))
	return math.Round(value*ratio) / ratio
}

func addContractsForEscalationCalculation(contractsMap map[string]models.ReContractYear, contractsMapDB map[string]models.ReContractYear, contractsAmount int, contractName string) (map[string]models.ReContractYear, error) {
	if contractsAmount > 0 {
		for i := contractsAmount; i != 0; i-- {
			name := fmt.Sprintf("%s %s", contractName, strconv.FormatInt(int64(i), 10))
			// If a contract that is expected to be used in the calculations is not found then return error
			if _, exist := contractsMapDB[name]; exist {
				contractsMap[contractsMapDB[name].Name] = contractsMapDB[name]
			} else {
				err := fmt.Errorf("expected contract %s not found", name)
				return contractsMap, err
			}
		}
	}
	return contractsMap, nil
}
