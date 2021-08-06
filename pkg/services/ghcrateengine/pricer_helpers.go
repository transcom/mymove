package ghcrateengine

import (
	"fmt"
	"math"
	"time"

	"github.com/gobuffalo/pop/v5"
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/unit"
)

func priceDomesticFirstDaySIT(db *pop.Connection, firstDaySITCode models.ReServiceCode, contractCode string, requestedPickupDate time.Time, weight unit.Pound, serviceArea string) (unit.Cents, services.PricingDisplayParams, error) {
	var sitType string
	if firstDaySITCode == models.ReServiceCodeDDFSIT {
		sitType = "destination"
	} else if firstDaySITCode == models.ReServiceCodeDOFSIT {
		sitType = "origin"
	} else {
		return 0, nil, fmt.Errorf("unsupported first day sit code of %s", firstDaySITCode)
	}

	if weight < minDomesticWeight {
		return 0, nil, fmt.Errorf("weight of %d less than the minimum of %d", weight, minDomesticWeight)
	}

	isPeakPeriod := IsPeakPeriod(requestedPickupDate)
	serviceAreaPrice, err := fetchDomServiceAreaPrice(db, contractCode, firstDaySITCode, serviceArea, isPeakPeriod)
	if err != nil {
		return unit.Cents(0), nil, fmt.Errorf("could not fetch domestic %s first day SIT rate: %w", sitType, err)
	}

	contractYear, err := fetchContractYear(db, serviceAreaPrice.ContractID, requestedPickupDate)
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

func priceDomesticAdditionalDaysSIT(db *pop.Connection, additionalDaySITCode models.ReServiceCode, contractCode string, requestedPickupDate time.Time, weight unit.Pound, serviceArea string, numberOfDaysInSIT int) (unit.Cents, services.PricingDisplayParams, error) {
	var sitType string
	if additionalDaySITCode == models.ReServiceCodeDDASIT {
		sitType = "destination"
	} else if additionalDaySITCode == models.ReServiceCodeDOASIT {
		sitType = "origin"
	} else {
		return 0, nil, fmt.Errorf("unsupported additional day sit code of %s", additionalDaySITCode)
	}

	if weight < minDomesticWeight {
		return 0, nil, fmt.Errorf("weight of %d less than the minimum of %d", weight, minDomesticWeight)
	}

	isPeakPeriod := IsPeakPeriod(requestedPickupDate)
	serviceAreaPrice, err := fetchDomServiceAreaPrice(db, contractCode, additionalDaySITCode, serviceArea, isPeakPeriod)
	if err != nil {
		return unit.Cents(0), nil, fmt.Errorf("could not fetch domestic %s additional days SIT rate: %w", sitType, err)
	}

	contractYear, err := fetchContractYear(db, serviceAreaPrice.ContractID, requestedPickupDate)
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

func priceDomesticPickupDeliverySIT(db *pop.Connection, pickupDeliverySITCode models.ReServiceCode, contractCode string, requestedPickupDate time.Time, weight unit.Pound, serviceArea string, sitSchedule int, zipOriginal string, zipActual string, distance unit.Miles) (unit.Cents, services.PricingDisplayParams, error) {
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
		shorthaulPricer := NewDomesticShorthaulPricer(db)
		totalPriceCents, displayParams, err := shorthaulPricer.Price(contractCode, requestedPickupDate, distance, weight, serviceArea)
		if err != nil {
			return unit.Cents(0), nil, fmt.Errorf("could not price shorthaul: %w", err)
		}

		return totalPriceCents, displayParams, nil
	}

	// Zip3s must be different at this point.  Now examine distance.

	// 2) Zip3 to different zip3 and > 50 miles
	if distance > 50 {
		// Do a normal linehaul calculation
		linehaulPricer := NewDomesticLinehaulPricer(db)
		totalPriceCents, displayParams, err := linehaulPricer.Price(contractCode, requestedPickupDate, distance, weight, serviceArea)
		if err != nil {
			return unit.Cents(0), nil, fmt.Errorf("could not price linehaul: %w", err)
		}

		return totalPriceCents, displayParams, nil
	}

	// Zip3s must be different at this point and distance is <= 50.

	// 3) Zip3 to different zip3 and <= 50 miles

	// Rate comes from the domestic other price table based on SIT schedule
	isPeakPeriod := IsPeakPeriod(requestedPickupDate)
	domOtherPrice, err := fetchDomOtherPrice(db, contractCode, pickupDeliverySITCode, sitSchedule, isPeakPeriod)
	if err != nil {
		return unit.Cents(0), nil, fmt.Errorf("could not fetch domestic %s SIT %s rate: %w", sitType, sitModifier, err)
	}
	contractYear, err := fetchContractYear(db, domOtherPrice.ContractID, requestedPickupDate)
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

func priceDomesticShuttling(db *pop.Connection, shuttlingCode models.ReServiceCode, contractCode string, requestedPickupDate time.Time, weight unit.Pound, serviceSchedule int) (unit.Cents, services.PricingDisplayParams, error) {
	if shuttlingCode != models.ReServiceCodeDOSHUT && shuttlingCode != models.ReServiceCodeDDSHUT {
		return 0, nil, fmt.Errorf("unsupported domestic shuttling code of %s", shuttlingCode)
	}
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
	if serviceSchedule == 0 {
		return 0, nil, errors.New("Service schedule is required")
	}

	// look up rate for domestic accessorial price
	domAccessorialPrice, err := fetchAccessorialPrice(db, contractCode, shuttlingCode, serviceSchedule)
	if err != nil {
		return 0, nil, fmt.Errorf("Could not lookup Domestic Accessorial Area Price: %w", err)
	}

	contractYear, err := fetchContractYear(db, domAccessorialPrice.ContractID, requestedPickupDate)
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

func priceDomesticCrating(db *pop.Connection, code models.ReServiceCode, contractCode string, requestedPickupDate time.Time, billedCubicFeet unit.CubicFeet, serviceSchedule int) (unit.Cents, services.PricingDisplayParams, error) {
	if code != models.ReServiceCodeDCRT && code != models.ReServiceCodeDUCRT {
		return 0, nil, fmt.Errorf("unsupported domestic crating code of %s", code)
	}

	if billedCubicFeet < 4.0 {
		return 0, nil, fmt.Errorf("crate must be billed for a minimum of 4 cubic feet")
	}
	domAccessorialPrice, err := fetchAccessorialPrice(db, contractCode, code, serviceSchedule)
	if err != nil {
		return 0, nil, fmt.Errorf("could not lookup Domestic Accessorial Area Price: %w", err)
	}

	basePrice := domAccessorialPrice.PerUnitCents.Float64() * float64(billedCubicFeet)
	contractYear, err := fetchContractYear(db, domAccessorialPrice.ContractID, requestedPickupDate)
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
func createPricerGeneratedParams(db *pop.Connection, paymentServiceItemID uuid.UUID, params services.PricingDisplayParams) (models.PaymentServiceItemParams, error) {
	var paymentServiceItemParams models.PaymentServiceItemParams

	if len(params) == 0 {
		return paymentServiceItemParams, fmt.Errorf("PricingDisplayParams must not be empty")
	}

	for _, param := range params {

		// Find the paymentServiceItemParam associated with this PricingDisplayParam
		var serviceItemParamKey models.ServiceItemParamKey
		err := db.Q().
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

		verrs, err := db.ValidateAndCreate(&newParam)
		if err != nil {
			return paymentServiceItemParams, fmt.Errorf("failure creating payment service item param: %w", err)
		} else if verrs.HasAny() {
			return paymentServiceItemParams, services.NewInvalidCreateInputError(verrs, "validation error with creating payment service item param")
		} else {
			// Append it to a slice of PaymentServiceItemParams to return
			paymentServiceItemParams = append(paymentServiceItemParams, newParam)
		}
	}
	return paymentServiceItemParams, nil
}
