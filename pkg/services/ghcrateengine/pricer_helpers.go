package ghcrateengine

import (
	"fmt"
	"math"
	"time"

	"github.com/gofrs/uuid"
	"github.com/pkg/errors"

	"github.com/transcom/mymove/pkg/apperror"

	"github.com/transcom/mymove/pkg/appcontext"
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
		return 0, nil, fmt.Errorf("Could not lookup domestic other price: %w", err)
	}

	var contractYear models.ReContractYear
	err = appCtx.DB().Where("contract_id = $1", domOtherPrice.ContractID).
		Where("$2 between start_date and end_date", referenceDate).
		First(&contractYear)
	if err != nil {
		return 0, nil, fmt.Errorf("Could not lookup contract year: %w", err)
	}

	finalWeight := weight
	if isPPM && weight < minDomesticWeight {
		finalWeight = minDomesticWeight
	}

	basePrice := domOtherPrice.PriceCents.Float64() * finalWeight.ToCWTFloat64()
	escalatedPrice := basePrice * contractYear.EscalationCompounded

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
			return 0, nil, fmt.Errorf("Could not lookup shipment type price: %w", err)
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

	contractYear, err := fetchContractYear(appCtx, serviceAreaPrice.ContractID, referenceDate)
	if err != nil {
		return unit.Cents(0), nil, fmt.Errorf("could not fetch contract year: %w", err)
	}

	baseTotalPrice := serviceAreaPrice.PriceCents.Float64() * weight.ToCWTFloat64()
	escalatedTotalPrice := baseTotalPrice * contractYear.EscalationCompounded

	totalPriceCents := unit.Cents(math.Round(escalatedTotalPrice))

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

	contractYear, err := fetchContractYear(appCtx, serviceAreaPrice.ContractID, referenceDate)
	if err != nil {
		return unit.Cents(0), nil, fmt.Errorf("could not fetch contract year: %w", err)
	}

	baseTotalPrice := serviceAreaPrice.PriceCents.Float64() * weight.ToCWTFloat64()
	escalatedTotalPrice := baseTotalPrice * contractYear.EscalationCompounded
	totalForNumberOfDaysPrice := escalatedTotalPrice * float64(numberOfDaysInSIT)

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
		zipOriginalName = "destination"
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
	contractYear, err := fetchContractYear(appCtx, domOtherPrice.ContractID, referenceDate)
	if err != nil {
		return unit.Cents(0), nil, fmt.Errorf("could not fetch contract year: %w", err)
	}

	baseTotalPrice := domOtherPrice.PriceCents.Float64() * weight.ToCWTFloat64()
	escalatedTotalPrice := baseTotalPrice * contractYear.EscalationCompounded
	totalPriceCents := unit.Cents(math.Round(escalatedTotalPrice))

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
		return 0, nil, fmt.Errorf("Could not lookup Domestic Accessorial Area Price: %w", err)
	}

	contractYear, err := fetchContractYear(appCtx, domAccessorialPrice.ContractID, referenceDate)
	if err != nil {
		return 0, nil, fmt.Errorf("Could not lookup contract year: %w", err)
	}

	basePrice := domAccessorialPrice.PerUnitCents.Float64() * weight.ToCWTFloat64()
	escalatedPrice := basePrice * contractYear.EscalationCompounded
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

func priceDomesticCrating(appCtx appcontext.AppContext, code models.ReServiceCode, contractCode string, referenceDate time.Time, billedCubicFeet unit.CubicFeet, serviceSchedule int) (unit.Cents, services.PricingDisplayParams, error) {
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

	basePrice := domAccessorialPrice.PerUnitCents.Float64() * float64(billedCubicFeet)
	contractYear, err := fetchContractYear(appCtx, domAccessorialPrice.ContractID, referenceDate)
	if err != nil {
		return 0, nil, fmt.Errorf("could not lookup contract year: %w", err)
	}
	escalatedPrice := basePrice * contractYear.EscalationCompounded
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
		} else {
			// Append it to a slice of PaymentServiceItemParams to return
			paymentServiceItemParams = append(paymentServiceItemParams, newParam)
		}
	}
	return paymentServiceItemParams, nil
}
