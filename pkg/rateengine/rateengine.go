package rateengine

import (
	"time"

	"github.com/transcom/mymove/pkg/models"

	"github.com/gobuffalo/pop/v5"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/transcom/mymove/pkg/unit"
)

// MaxSITDays is the maximum number of days of SIT that will be reimbursed.
const MaxSITDays = 90

// RateEngine encapsulates the TSP rate engine process
type RateEngine struct {
	db     *pop.Connection
	logger Logger
	move   models.Move
}

// CostComputation represents the results of a computation.
type CostComputation struct {
	LinehaulCostComputation
	NonLinehaulCostComputation

	SITFee      unit.Cents
	SITMax      unit.Cents
	GCC         unit.Cents
	LHDiscount  unit.DiscountRate
	SITDiscount unit.DiscountRate
	Weight      unit.Pound
}

// CostDetail holds the costComputation and a bool that signifies if the calculation is the winning (lowest cost) computation
type CostDetail struct {
	Cost      CostComputation
	IsWinning bool
}

// CostDetails is a map of CostDetail
type CostDetails map[string]*CostDetail

// Scale scales a cost computation by a multiplicative factor
func (c *CostComputation) Scale(factor float64) {
	c.LinehaulCostComputation.Scale(factor)
	c.NonLinehaulCostComputation.Scale(factor)

	c.SITFee = c.SITFee.MultiplyFloat64(factor)
	c.SITMax = c.SITMax.MultiplyFloat64(factor)
	c.GCC = c.GCC.MultiplyFloat64(factor)
}

// MarshalLogObject allows CostComputation to be logged by Zap.
func (c CostComputation) MarshalLogObject(encoder zapcore.ObjectEncoder) error {
	if err := encoder.AddObject("Linehaul Components", c.LinehaulCostComputation); err != nil {
		return err
	}

	encoder.AddInt("LinehaulChargeTotal", c.LinehaulChargeTotal.Int())
	encoder.AddInt("FuelSurcharge", c.FuelSurcharge.Fee.Int())

	encoder.AddInt("OriginServiceFee", c.OriginService.Fee.Int())
	encoder.AddInt("DestinationServiceFee", c.DestinationService.Fee.Int())
	encoder.AddInt("PackFee", c.Pack.Fee.Int())
	encoder.AddInt("UnpackFee", c.Unpack.Fee.Int())
	encoder.AddInt("SITMax", c.SITMax.Int())
	encoder.AddInt("SITFee", c.SITFee.Int())

	encoder.AddInt("GCC", c.GCC.Int())

	encoder.AddFloat64("LHDiscount", float64(c.LHDiscount))
	encoder.AddFloat64("SITDiscount", float64(c.SITDiscount))
	encoder.AddInt("Miles", c.Mileage)
	encoder.AddInt("Weight", c.Weight.Int())

	return nil
}

// Zip5ToZip3 takes a ZIP5 string and returns the ZIP3 representation of it.
func Zip5ToZip3(zip5 string) string {
	return zip5[0:3]
}

// computePPM Calculates the cost of a PPM move.
func (re *RateEngine) computePPM(
	weight unit.Pound,
	originZip5 string,
	destinationZip5 string,
	distanceMiles int,
	date time.Time,
	daysInSIT int,
	lhDiscount unit.DiscountRate,
	sitDiscount unit.DiscountRate) (cost CostComputation, err error) {

	// Weights below 1000lbs are prorated to the 1000lb rate
	prorateFactor := 1.0
	if weight.Int() < 1000 {
		prorateFactor = weight.Float64() / 1000.0
		weight = unit.Pound(1000)
	}

	// Linehaul charges
	linehaulCostComputation, err := re.linehaulChargeComputation(weight, originZip5, destinationZip5, distanceMiles, date)
	if err != nil {
		re.logger.Error("Failed to compute linehaul cost", zap.Error(err))
		return
	}

	// Non linehaul charges
	nonLinehaulCostComputation, err := re.nonLinehaulChargeComputation(weight, originZip5, destinationZip5, date)
	if err != nil {
		re.logger.Error("Failed to compute non-linehaul cost", zap.Error(err))
		return
	}

	// Apply linehaul discounts
	linehaulCostComputation.LinehaulChargeTotal = lhDiscount.Apply(linehaulCostComputation.LinehaulChargeTotal)
	nonLinehaulCostComputation.OriginService.Fee = lhDiscount.Apply(nonLinehaulCostComputation.OriginService.Fee)
	nonLinehaulCostComputation.DestinationService.Fee = lhDiscount.Apply(nonLinehaulCostComputation.DestinationService.Fee)
	nonLinehaulCostComputation.Pack.Fee = lhDiscount.Apply(nonLinehaulCostComputation.Pack.Fee)
	nonLinehaulCostComputation.Unpack.Fee = lhDiscount.Apply(nonLinehaulCostComputation.Unpack.Fee)

	// SIT
	// Note that SIT has a different discount rate than [non]linehaul charges
	destinationZip3 := Zip5ToZip3(destinationZip5)
	sitComputation, err := re.SitCharge(weight.ToCWT(), daysInSIT, destinationZip3, date, true)
	if err != nil {
		re.logger.Info("Can't calculate sit",
			zap.String("moveLocator", re.move.Locator),
		)
		return
	}
	sitFee := sitComputation.ApplyDiscount(lhDiscount, sitDiscount)

	/// Max SIT
	maxSITComputation, err := re.SitCharge(weight.ToCWT(), MaxSITDays, destinationZip3, date, true)
	if err != nil {
		re.logger.Info("Can't calculate max sit",
			zap.String("moveLocator", re.move.Locator),
		)
		return
	}
	// Note that SIT has a different discount rate than [non]linehaul charges
	maxSITFee := maxSITComputation.ApplyDiscount(lhDiscount, sitDiscount)

	// Totals
	gcc := linehaulCostComputation.LinehaulChargeTotal +
		nonLinehaulCostComputation.OriginService.Fee +
		nonLinehaulCostComputation.DestinationService.Fee +
		nonLinehaulCostComputation.Pack.Fee +
		nonLinehaulCostComputation.Unpack.Fee

	cost = CostComputation{
		LinehaulCostComputation:    linehaulCostComputation,
		NonLinehaulCostComputation: nonLinehaulCostComputation,
		SITFee:                     sitFee,
		SITMax:                     maxSITFee,
		GCC:                        gcc,
		LHDiscount:                 lhDiscount,
		SITDiscount:                sitDiscount,
		Weight:                     weight,
	}

	// Finally, scale by prorate factor
	cost.Scale(prorateFactor)

	re.logger.Info("PPM cost computation",
		zap.String("moveLocator", re.move.Locator),
		zap.Object("cost", cost),
	)

	return cost, nil
}

//computePPMIncludingLHDiscount Calculates the cost of a PPM move using zip + date derived linehaul discount
func (re *RateEngine) computePPMIncludingLHDiscount(weight unit.Pound, originZip5 string, destinationZip5 string, distanceMiles int, date time.Time, daysInSIT int) (cost CostComputation, err error) {

	lhDiscount, sitDiscount, err := models.PPMDiscountFetch(re.db,
		re.logger,
		re.move,
		originZip5,
		destinationZip5,
		date,
	)
	if err != nil {
		re.logger.Error("Failed to compute linehaul cost", zap.Error(err))
		return
	}

	cost, err = re.computePPM(weight,
		originZip5,
		destinationZip5,
		distanceMiles,
		date,
		daysInSIT,
		lhDiscount,
		sitDiscount,
	)

	if err != nil {
		re.logger.Error("Failed to compute PPM cost", zap.Error(err))
		return
	}
	return cost, nil
}

// ComputePPMMoveCosts uses zip codes to make two calculations for the price of a PPM move - once with the pickup zip and once with the current duty station zip - and returns both calcs.
func (re *RateEngine) ComputePPMMoveCosts(weight unit.Pound, originPickupZip5 string, originDutyStationZip5 string, destinationZip5 string, distanceMilesFromOriginPickupZip int, distanceMilesFromOriginDutyStationZip int, date time.Time, daysInSit int) (costDetails CostDetails, err error) {
	costFromOriginPickupZip, err := re.computePPMIncludingLHDiscount(
		weight,
		originPickupZip5,
		destinationZip5,
		distanceMilesFromOriginPickupZip,
		date,
		daysInSit,
	)
	if err != nil {
		re.logger.Error("Failed to compute PPM cost", zap.Error(err))
		return
	}

	costDetails = make(CostDetails)
	costDetails["pickupLocation"] = &CostDetail{
		costFromOriginPickupZip,
		false,
	}

	costFromOriginDutyStationZip, err := re.computePPMIncludingLHDiscount(
		weight,
		originDutyStationZip5,
		destinationZip5,
		distanceMilesFromOriginDutyStationZip,
		date,
		daysInSit,
	)
	if err != nil {
		re.logger.Error("Failed to compute PPM cost", zap.Error(err))
		return
	}
	costDetails["originDutyStation"] = &CostDetail{
		costFromOriginDutyStationZip,
		false,
	}

	originZipCode := originPickupZip5
	originZipLocation := "Pickup location"
	if costFromOriginPickupZip.GCC > costFromOriginDutyStationZip.GCC {
		costDetails["originDutyStation"].IsWinning = true
		originZipCode = originDutyStationZip5
		originZipLocation = "Origin duty station"
	} else {
		costDetails["pickupLocation"].IsWinning = true
	}

	re.logger.Info("Origin zip code information",
		zap.String("moveLocator", re.move.Locator),
		zap.String("originZipLocation", originZipLocation),
		zap.String("originZipCode", originZipCode),
	)
	return costDetails, nil
}

// GetWinningCostMove returns a costComputation of the winning calculation
func GetWinningCostMove(costDetails CostDetails) CostComputation {
	if costDetails["pickupLocation"].IsWinning {
		return costDetails["pickupLocation"].Cost
	}
	return costDetails["originDutyStation"].Cost
}

// GetNonWinningCostMove returns a costComputation of the non-winning calculation
func GetNonWinningCostMove(costDetails CostDetails) CostComputation {
	if costDetails["pickupLocation"].IsWinning {
		return costDetails["originDutyStation"].Cost
	}
	return costDetails["pickupLocation"].Cost
}

// NewRateEngine creates a new RateEngine
func NewRateEngine(db *pop.Connection, logger Logger, move models.Move) *RateEngine {
	return &RateEngine{db: db, logger: logger, move: move}
}
