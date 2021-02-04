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
