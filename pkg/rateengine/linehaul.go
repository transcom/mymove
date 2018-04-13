package rateengine

import (
	"os"
	"time"

	"github.com/go-openapi/swag"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/route"
	"github.com/transcom/mymove/pkg/unit"
)

func (re *RateEngine) determineMileage(originZip string, destinationZip string) (mileage int, err error) {
	// JAMES ATHEY AT DDS SAYS NOT TO USE BASEPOINT CITY/STATE COMBO
	// originZip3, err := models.FetchCityAndStateForZip3(re.db, originZip)
	// destinationZip3, err := models.FetchCityAndStateForZip3(re.db, destinationZip)
	bingEndpoint := os.Getenv("BING_MAPS_ENDPOINT")
	bingKey := os.Getenv("BING_MAPS_KEY")
	planner := route.NewBingPlanner(re.logger, &bingEndpoint, &bingKey)
	sourceAddress := models.Address{
		StreetAddress1: "",
		StreetAddress2: swag.String(""),
		StreetAddress3: swag.String(""),
		City:           "",
		State:          originZip,
		PostalCode:     "",
	}
	destinationAddress := models.Address{
		StreetAddress1: "",
		StreetAddress2: swag.String(""),
		StreetAddress3: swag.String(""),
		City:           "",
		State:          "",
		PostalCode:     destinationZip,
	}
	// sourceAddress := models.Address{
	// 	StreetAddress1: "",
	// 	StreetAddress2: swag.String(""),
	// 	StreetAddress3: swag.String(""),
	// 	City:           originZip3.BasepointCity,
	// 	State:          originZip3.State,
	// 	PostalCode:     "",
	// }
	// destinationAddress := models.Address{
	// 	StreetAddress1: "",
	// 	StreetAddress2: swag.String(""),
	// 	StreetAddress3: swag.String(""),
	// 	City:           destinationZip3.BasepointCity,
	// 	State:          destinationZip3.State,
	// 	PostalCode:     "",
	// }
	distance, err := planner.TransitDistance(&sourceAddress, &destinationAddress)
	if err != nil {
		re.logger.Error("Failed to get distance from Bing - %v", zap.Error(err))
	}
	return distance, err
}

// Determine the Base Linehaul (BLH)
func (re *RateEngine) baseLinehaul(mileage int, cwt int, date time.Time) (baseLinehaulChargeCents unit.Cents, err error) {
	baseLinehaulChargeCents, err = models.FetchBaseLinehaulRate(re.db, mileage, cwt, date)
	if err != nil {
		re.logger.Error("Base Linehaul query didn't complete: ", zap.Error(err))
	}

	return baseLinehaulChargeCents, err
}

// Determine the Linehaul Factors (OLF and DLF)
func (re *RateEngine) linehaulFactors(cwt int, zip3 string, date time.Time) (linehaulFactorCents unit.Cents, err error) {
	serviceArea, err := models.FetchTariff400ngServiceAreaForZip3(re.db, zip3)
	if err != nil {
		return 0, err
	}
	linehaulFactorCents, err = models.FetchTariff400ngLinehaulFactor(re.db, serviceArea.ServiceArea, date)
	if err != nil {
		return 0, err
	}
	return linehaulFactorCents.Multiply(cwt), nil
}

// Determine Shorthaul (SH) Charge (ONLY applies if shipment moves 800 miles and less)
func (re *RateEngine) shorthaulCharge(mileage int, cwt int, date time.Time) (shorthaulChargeCents unit.Cents, err error) {
	if mileage >= 800 {
		return 0, nil
	}
	re.logger.Debug("Shipment qualifies for shorthaul fee",
		zap.Int("miles", mileage))

	cwtMiles := mileage * cwt
	shorthaulChargeCents, err = models.FetchShorthaulRateCents(re.db, cwtMiles, date)

	return shorthaulChargeCents, err
}

// Determine Linehaul Charge (LC) TOTAL
// Formula: LC= [BLH + OLF + DLF + [SH]
func (re *RateEngine) linehaulChargeTotal(weight int, originZip string, destinationZip string, date time.Time) (linehaulChargeCents unit.Cents, err error) {
	mileage, err := re.determineMileage(originZip, destinationZip)
	cwt := re.determineCWT(weight)
	baseLinehaulChargeCents, err := re.baseLinehaul(mileage, cwt, date)
	if err != nil {
		return 0, err
	}
	originLinehaulFactorCents, err := re.linehaulFactors(cwt, originZip, date)
	if err != nil {
		return 0, err
	}
	destinationLinehaulFactorCents, err := re.linehaulFactors(cwt, destinationZip, date)
	if err != nil {
		return 0, err
	}
	shorthaulChargeCents, err := re.shorthaulCharge(mileage, cwt, date)
	if err != nil {
		return 0, err
	}

	linehaulChargeCents = baseLinehaulChargeCents + originLinehaulFactorCents + destinationLinehaulFactorCents + shorthaulChargeCents
	re.logger.Info("Linehaul charge total calculated",
		zap.Int("linehaul total", linehaulChargeCents.Int()),
		zap.Int("linehaul", baseLinehaulChargeCents.Int()),
		zap.Int("origin lh factor", originLinehaulFactorCents.Int()),
		zap.Int("destination lh factor", destinationLinehaulFactorCents.Int()),
		zap.Int("shorthaul", shorthaulChargeCents.Int()))

	return linehaulChargeCents, err
}
