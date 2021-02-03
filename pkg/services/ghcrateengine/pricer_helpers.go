package ghcrateengine

import (
	"fmt"
	"math"
	"time"

	"github.com/gobuffalo/pop/v5"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/unit"
)

func priceDomesticFirstDaySit(db *pop.Connection, firstDaySitCode models.ReServiceCode, contractCode string, requestedPickupDate time.Time, isPeakPeriod bool, weight unit.Pound, serviceArea string) (unit.Cents, error) {
	var sitType string
	if firstDaySitCode == models.ReServiceCodeDDFSIT {
		sitType = "destination"
	} else if firstDaySitCode == models.ReServiceCodeDOFSIT {
		sitType = "origin"
	} else {
		return 0, fmt.Errorf("unsupported first day sit code of %s", firstDaySitCode)
	}

	if weight < minDomesticWeight {
		return 0, fmt.Errorf("weight of %d less than the minimum of %d", weight, minDomesticWeight)
	}

	serviceAreaPrice, err := fetchDomServiceAreaPrice(db, contractCode, firstDaySitCode, serviceArea, isPeakPeriod)
	if err != nil {
		return unit.Cents(0), fmt.Errorf("could not fetch domestic %s first day SIT rate: %w", sitType, err)
	}

	contractYear, err := fetchContractYear(db, serviceAreaPrice.ContractID, requestedPickupDate)
	if err != nil {
		return unit.Cents(0), fmt.Errorf("could not fetch contract year: %w", err)
	}

	baseTotalPrice := serviceAreaPrice.PriceCents.Float64() * weight.ToCWTFloat64()
	escalatedTotalPrice := baseTotalPrice * contractYear.EscalationCompounded

	totalPriceCents := unit.Cents(math.Round(escalatedTotalPrice))

	return totalPriceCents, nil
}

func priceDomesticAdditionalDaysSit(db *pop.Connection, additionalDaySitCode models.ReServiceCode, contractCode string, requestedPickupDate time.Time, isPeakPeriod bool, weight unit.Pound, serviceArea string, numberOfDaysInSIT int) (unit.Cents, error) {
	var sitType string
	if additionalDaySitCode == models.ReServiceCodeDDASIT {
		sitType = "destination"
	} else if additionalDaySitCode == models.ReServiceCodeDOASIT {
		sitType = "origin"
	} else {
		return 0, fmt.Errorf("unsupported additional day sit code of %s", additionalDaySitCode)
	}

	if weight < minDomesticWeight {
		return 0, fmt.Errorf("weight of %d less than the minimum of %d", weight, minDomesticWeight)
	}

	serviceAreaPrice, err := fetchDomServiceAreaPrice(db, contractCode, additionalDaySitCode, serviceArea, isPeakPeriod)
	if err != nil {
		return unit.Cents(0), fmt.Errorf("could not fetch domestic %s additional days SIT rate: %w", sitType, err)
	}

	contractYear, err := fetchContractYear(db, serviceAreaPrice.ContractID, requestedPickupDate)
	if err != nil {
		return unit.Cents(0), fmt.Errorf("could not fetch contract year: %w", err)
	}

	baseTotalPrice := serviceAreaPrice.PriceCents.Float64() * weight.ToCWTFloat64()
	escalatedTotalPrice := baseTotalPrice * contractYear.EscalationCompounded
	totalForNumberOfDaysPrice := escalatedTotalPrice * float64(numberOfDaysInSIT)

	totalPriceCents := unit.Cents(math.Round(totalForNumberOfDaysPrice))

	return totalPriceCents, nil
}

func priceDomesticPickupDeliverySIT(db *pop.Connection, pickupDeliverySITCode models.ReServiceCode, contractCode string, requestedPickupDate time.Time, isPeakPeriod bool,
	weight unit.Pound, serviceArea string, sitSchedule int, zipOriginal string, zipActual string, distance unit.Miles) (unit.Cents, error) {

	var sitType, sitModifier, zipOriginalName, zipActualName string
	if pickupDeliverySITCode == models.ReServiceCodeDDDSIT {
		sitType = "destination"
		sitModifier = "delivery"
		zipOriginalName = "destination"
		zipActualName = "SIT destination"
	} else if pickupDeliverySITCode == models.ReServiceCodeDOPSIT {
		sitType = "origin"
		sitModifier = "pickup"
		zipOriginalName = "SIT origin original"
		zipActualName = "SIT origin actual"
	} else {
		return 0, fmt.Errorf("unsupported pickup/delivery SIT code of %s", pickupDeliverySITCode)
	}

	if weight < minDomesticWeight {
		return 0, fmt.Errorf("weight of %d less than the minimum of %d", weight, minDomesticWeight)
	}

	if len(zipOriginal) < 5 {
		return unit.Cents(0), fmt.Errorf("invalid %s postal code of %s", zipOriginalName, zipOriginal)
	}
	zip3Original := zipOriginal[:3]

	if len(zipActual) < 5 {
		return unit.Cents(0), fmt.Errorf("invalid %s postal code of %s", zipActualName, zipActual)
	}
	zip3Actual := zipActual[:3]

	// Three different pricing scenarios below.

	// 1) Zip3 to same zip3
	if zip3Original == zip3Actual {
		// Do a normal shorthaul calculation
		shorthaulPricer := NewDomesticShorthaulPricer(db)
		totalPriceCents, err := shorthaulPricer.Price(contractCode, requestedPickupDate, distance, weight, serviceArea)
		if err != nil {
			return unit.Cents(0), fmt.Errorf("could not price shorthaul: %w", err)
		}

		return totalPriceCents, nil
	}

	// Zip3s must be different at this point.  Now examine distance.

	// 2) Zip3 to different zip3 and > 50 miles
	if distance > 50 {
		// Do a normal linehaul calculation
		linehaulPricer := NewDomesticLinehaulPricer(db)
		totalPriceCents, err := linehaulPricer.Price(contractCode, requestedPickupDate, isPeakPeriod, distance, weight, serviceArea)
		if err != nil {
			return unit.Cents(0), fmt.Errorf("could not price linehaul: %w", err)
		}

		return totalPriceCents, nil
	}

	// Zip3s must be different at this point and distance is <= 50.

	// 3) Zip3 to different zip3 and <= 50 miles

	// Rate comes from the domestic other price table based on SIT schedule
	domOtherPrice, err := fetchDomOtherPrice(db, contractCode, pickupDeliverySITCode, sitSchedule, isPeakPeriod)
	if err != nil {
		return unit.Cents(0), fmt.Errorf("could not fetch domestic %s SIT %s rate: %w", sitType, sitModifier, err)
	}
	contractYear, err := fetchContractYear(db, domOtherPrice.ContractID, requestedPickupDate)
	if err != nil {
		return unit.Cents(0), fmt.Errorf("could not fetch contract year: %w", err)
	}

	baseTotalPrice := domOtherPrice.PriceCents.Float64() * weight.ToCWTFloat64()
	escalatedTotalPrice := baseTotalPrice * contractYear.EscalationCompounded
	totalPriceCents := unit.Cents(math.Round(escalatedTotalPrice))

	return totalPriceCents, nil
}
