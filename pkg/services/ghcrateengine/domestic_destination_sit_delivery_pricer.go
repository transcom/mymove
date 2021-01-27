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

type domesticDestinationSITDeliveryPricer struct {
	db *pop.Connection
}

// NewDomesticDestinationSITDeliveryPricer creates a new pricer for domestic destination SIT delivery
func NewDomesticDestinationSITDeliveryPricer(db *pop.Connection) services.DomesticDestinationSITDeliveryPricer {
	return &domesticDestinationSITDeliveryPricer{
		db: db,
	}
}

// Price determines the price for domestic destination SIT delivery
func (p domesticDestinationSITDeliveryPricer) Price(contractCode string, requestedPickupDate time.Time, isPeakPeriod bool,
	weight unit.Pound, serviceArea string, sitSchedule int, zipDest string, zipSITDest string, distance unit.Miles) (unit.Cents, error) {

	if weight < minDomesticWeight {
		return 0, fmt.Errorf("weight of %d less than the minimum of %d", weight, minDomesticWeight)
	}

	// Three different pricing scenarios below.

	// 1) Less than or equal to 50 miles (if same or different zip3s)
	if distance <= 50 {
		// Rate comes from the domestic other price table based on SIT schedule
		domOtherPrice, err := fetchDomOtherPrice(p.db, contractCode, models.ReServiceCodeDDDSIT, sitSchedule, isPeakPeriod)
		if err != nil {
			return unit.Cents(0), fmt.Errorf("could not fetch domestic destination SIT delivery rate: %w", err)
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

	// Distance must be greater than 50 miles at this point.  Now examine zip3s.

	if len(zipDest) < 5 {
		return unit.Cents(0), fmt.Errorf("invalid destination postal code of %s", zipDest)
	}
	zip3Dest := zipDest[:3]

	if len(zipSITDest) < 5 {
		return unit.Cents(0), fmt.Errorf("invalid SIT destination postal code of %s", zipSITDest)
	}
	zip3SITDest := zipSITDest[:3]

	// 2) Greater than 50 miles and different zip3s
	if zip3Dest != zip3SITDest {
		// Do a normal linehaul calculation
		linehaulPricer := NewDomesticLinehaulPricer(p.db)
		totalPriceCents, err := linehaulPricer.Price(contractCode, requestedPickupDate, isPeakPeriod, distance, weight, serviceArea)
		if err != nil {
			return unit.Cents(0), fmt.Errorf("could not price linehaul: %w", err)
		}

		return totalPriceCents, nil
	}

	// Distance must be greater than 50 miles and the zip3s are the same at this point.

	// 3) Greater than 50 miles and same zip3s
	shorthaulPricer := NewDomesticShorthaulPricer(p.db)
	totalPriceCents, err := shorthaulPricer.Price(contractCode, requestedPickupDate, distance, weight, serviceArea)
	if err != nil {
		return unit.Cents(0), fmt.Errorf("could not price shorthaul: %w", err)
	}

	return totalPriceCents, nil
}

// PriceUsingParams determines the price for domestic destination SIT delivery given PaymentServiceItemParams
func (p domesticDestinationSITDeliveryPricer) PriceUsingParams(params models.PaymentServiceItemParams) (unit.Cents, error) {
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

	serviceAreaDest, err := getParamString(params, models.ServiceItemParamNameServiceAreaDest)
	if err != nil {
		return unit.Cents(0), err
	}

	sitScheduleDest, err := getParamInt(params, models.ServiceItemParamNameSITScheduleDest)
	if err != nil {
		return unit.Cents(0), err
	}

	zipDestAddress, err := getParamString(params, models.ServiceItemParamNameZipDestAddress)
	if err != nil {
		return unit.Cents(0), err
	}

	zipSITDestHHGFinalAddress, err := getParamString(params, models.ServiceItemParamNameZipSITDestHHGFinalAddress)
	if err != nil {
		return unit.Cents(0), err
	}

	distanceZipSITDest, err := getParamInt(params, models.ServiceItemParamNameDistanceZipSITDest)
	if err != nil {
		return unit.Cents(0), err
	}

	isPeakPeriod := IsPeakPeriod(requestedPickupDate)

	return p.Price(contractCode, requestedPickupDate, isPeakPeriod, unit.Pound(weightBilledActual), serviceAreaDest,
		sitScheduleDest, zipDestAddress, zipSITDestHHGFinalAddress, unit.Miles(distanceZipSITDest))
}
