package ghcrateengine

import (
	"fmt"
	"time"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/unit"
)

const PeakRateCycleStartMonth = time.May
const PeakRateCycleStartDay = 15
const PeakRateCycleEndMonth = time.September
const PeakRateCycleEndDay = 30

type DomesticServicePricingData struct {
	MoveDate      time.Time
	ServiceAreaID uuid.UUID
	Distance      unit.Miles
	Weight        unit.Pound
}

func getIsPeakPeriod(moveDate time.Time) bool {
	isPeakPeriod := false
	peakRateCycleStart := time.Date(moveDate.Year(), PeakRateCycleStartMonth, PeakRateCycleStartDay, 0, 0, 0, 0, time.UTC)
	peakRateCycleEnd := time.Date(moveDate.Year(), PeakRateCycleEndMonth, PeakRateCycleEndDay, 0, 0, 0, 0, time.UTC)
	if moveDate.After(peakRateCycleStart.AddDate(0, 0, -1)) && moveDate.Before(peakRateCycleEnd.AddDate(0, 0, 1)) {
		isPeakPeriod = true
	}
	return isPeakPeriod
}

func getDomesticLinehaulRate(moveDate time.Time, serviceAreaID uuid.UUID, weight unit.Pound, distance unit.Miles) unit.Cents {
	isPeakPeriod := getIsPeakPeriod(moveDate)
	fmt.Printf("is Peak Period %v", isPeakPeriod)
	//PSEUDOCODE
	// var rate = look up rate in domestic linehaul prices where
	// is_peak_period or not
	// the weight is between weight_lower and weight_upper (inclusive)
	// it equals serviceArea.serviceAreaNumber
	// the distance is between miles_upper and miles_lower (inclusive)
	// is correct contract

	// To be implemented when models are created
	//query := gre.db.Where(
	//	"isPeakPeriod = isPeakPeriod").Where(
	//		"weight gte weight_lower").Where(
	//			"weight lte weight_upper").Where(
	//				"distance gte miles_lower").Where(
	//					"distance lte miles_upper").LeftJoin(
	//						"serviceAreas sa", "sa.id=re_domestic_linehaul_prices.service_area_id").LeftJoin(
	//							"contract", "contract.id=re_domestic_linehaul_prices.contract_id")
	//domesticLinehaulPricing := models.DomesticLinehaulPricing{}
	//err := gre.db.query.All(&domesticLinehaulPricing)
	var stubPrice unit.Cents = 4045
	return stubPrice
}

// Calculation Functions
// PriceDemesticLinehaul calculates the cost domestic linehaul
func (gre *GHCRateEngine) PriceDomesticLinehaul(d DomesticServicePricingData) unit.Cents {
	rate := getDomesticLinehaulRate(d.MoveDate, d.ServiceAreaID, d.Weight, d.Distance)
	// calculate total
	return rate.Multiply(d.Weight.ToCWT().Int())
}
