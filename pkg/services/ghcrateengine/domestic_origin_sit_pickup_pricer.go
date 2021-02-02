package ghcrateengine

import (
	"fmt"
	"math"
	"time"

	"github.com/gobuffalo/pop/v5"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/unit"
)

type domesticOriginSITPickupPricer struct {
	db *pop.Connection
}

// NewDomesticOriginSITPickupPricer creates a new pricer for domestic origin SIT pickup
func NewDomesticOriginSITPickupPricer(db *pop.Connection) services.DomesticOriginSITPickupPricer {
	return &domesticOriginSITPickupPricer{
		db: db,
	}
}

// Price determines the price for domestic origin SIT pickup
func (p domesticOriginSITPickupPricer) Price(contractCode string, requestedPickupDate time.Time, isPeakPeriod bool,
	weight unit.Pound, serviceArea string, sitSchedule int, zipSITOriginOriginal string, zipSITOriginActual string, distance unit.Miles) (unit.Cents, error) {

	if weight < minDomesticWeight {
		return 0, fmt.Errorf("weight of %d less than the minimum of %d", weight, minDomesticWeight)
	}

	if len(zipSITOriginOriginal) < 5 {
		return unit.Cents(0), fmt.Errorf("invalid SIT origin original postal code of %s", zipSITOriginOriginal)
	}
	zip3Original := zipSITOriginOriginal[:3]

	if len(zipSITOriginActual) < 5 {
		return unit.Cents(0), fmt.Errorf("invalid SIT origin actual postal code of %s", zipSITOriginActual)
	}
	zip3Actual := zipSITOriginActual[:3]

	// Three different pricing scenarios below.

	// 1) Zip3 to same zip3
	if zip3Original == zip3Actual {
		// Do a normal shorthaul calculation
		shorthaulPricer := NewDomesticShorthaulPricer(p.db)
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
		linehaulPricer := NewDomesticLinehaulPricer(p.db)
		totalPriceCents, err := linehaulPricer.Price(contractCode, requestedPickupDate, isPeakPeriod, distance, weight, serviceArea)
		if err != nil {
			return unit.Cents(0), fmt.Errorf("could not price linehaul: %w", err)
		}

		return totalPriceCents, nil
	}

	// Zip3s must be different at this point and distance is <= 50.

	// 3) Zip3 to different zip3 and <= 50 miles

	// Rate comes from the domestic other price table based on SIT schedule
	domOtherPrice, err := fetchDomOtherPrice(p.db, contractCode, models.ReServiceCodeDOPSIT, sitSchedule, isPeakPeriod)
	if err != nil {
		return unit.Cents(0), fmt.Errorf("could not fetch domestic origin SIT pickup rate: %w", err)
	}
	contractYear, err := fetchContractYear(p.db, domOtherPrice.ContractID, requestedPickupDate)
	if err != nil {
		return unit.Cents(0), fmt.Errorf("could not fetch contract year: %w", err)
	}

	baseTotalPrice := domOtherPrice.PriceCents.Float64() * weight.ToCWTFloat64()
	escalatedTotalPrice := baseTotalPrice * contractYear.EscalationCompounded
	totalPriceCents := unit.Cents(math.Round(escalatedTotalPrice))

	return totalPriceCents, nil
}

// PriceUsingParams determines the price for domestic origin SIT pickup given PaymentServiceItemParams
func (p domesticOriginSITPickupPricer) PriceUsingParams(params models.PaymentServiceItemParams) (unit.Cents, error) {
	contractCode, err := getParamString(params, models.ServiceItemParamNameContractCode)
	if err != nil {
		return unit.Cents(0), err
	}

	requestedPickupDate, err := getParamTime(params, models.ServiceItemParamNameRequestedPickupDate)
	if err != nil {
		return unit.Cents(0), err
	}

	weightBilledActual, err := getParamInt(params, models.ServiceItemParamNameWeightBilledActual)
	if err != nil {
		return unit.Cents(0), err
	}

	serviceAreaOrigin, err := getParamString(params, models.ServiceItemParamNameServiceAreaOrigin)
	if err != nil {
		return unit.Cents(0), err
	}

	sitScheduleOrigin, err := getParamInt(params, models.ServiceItemParamNameSITScheduleOrigin)
	if err != nil {
		return unit.Cents(0), err
	}

	zipSITOriginOriginalAddress, err := getParamString(params, models.ServiceItemParamNameZipSITOriginHHGOriginalAddress)
	if err != nil {
		return unit.Cents(0), err
	}

	zipSITOriginActualAddress, err := getParamString(params, models.ServiceItemParamNameZipSITOriginHHGActualAddress)
	if err != nil {
		return unit.Cents(0), err
	}

	distanceZipSITOrigin, err := getParamInt(params, models.ServiceItemParamNameDistanceZipSITOrigin)
	if err != nil {
		return unit.Cents(0), err
	}

	isPeakPeriod := IsPeakPeriod(requestedPickupDate)

	return p.Price(contractCode, requestedPickupDate, isPeakPeriod, unit.Pound(weightBilledActual), serviceAreaOrigin,
		sitScheduleOrigin, zipSITOriginOriginalAddress, zipSITOriginActualAddress, unit.Miles(distanceZipSITOrigin))
}
