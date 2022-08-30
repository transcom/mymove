package rateengine

import (
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/unit"
)

// MaxSITDays is the maximum number of days of SIT that will be reimbursed.
const MaxSITDays = 90

// RateEngine encapsulates the TSP rate engine process
type RateEngine struct {
	move models.Move
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

// computePPM this is returning a hardcoded value because we aren't loading tariff400ng data
// hardcoded value to prevent errors when scheduling a ppm beyond 2021-05-15
// PPMs will be addressed in outcome 7
func (re *RateEngine) computePPM(
	appCtx appcontext.AppContext,
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

	linehaulCost := LinehaulCostComputation{
		BaseLinehaul:              310300,
		OriginLinehaulFactor:      1995,
		DestinationLinehaulFactor: 1770,
		ShorthaulCharge:           1,
	}

	nonLinehaulCost := NonLinehaulCostComputation{
		OriginService:      FeeAndRate{11025, 735000},
		DestinationService: FeeAndRate{4113, 554000},
		Pack:               FeeAndRate{49228, 6630000},
		Unpack:             FeeAndRate{5169, 696150},
	}

	cost = CostComputation{
		LinehaulCostComputation:    linehaulCost,
		NonLinehaulCostComputation: nonLinehaulCost,
		SITFee:                     0,
		SITMax:                     106166,
		GCC:                        219429,
		LHDiscount:                 lhDiscount,
		SITDiscount:                sitDiscount,
		Weight:                     weight,
	}

	// this formula means nothing - it's only so the estimate changes when the slider moves
	weightValue := weight.Float64()
	cost.Scale(weightValue / 1000 * .4)
	cost.Scale(prorateFactor)

	appCtx.Logger().Info("PPM cost computation",
		zap.String("moveLocator", re.move.Locator),
		zap.Object("cost", cost),
	)

	return cost, nil
}

// computePPM commented out in favor of returning a hardcoded struct above.
// left this function in rather than deleting because we will use this function when we get to outcome 7
// // computePPM Calculates the cost of a PPM move.
// func (re *RateEngine) computePPM(
// 	weight unit.Pound,
// 	originZip5 string,
// 	destinationZip5 string,
// 	distanceMiles int,
// 	date time.Time,
// 	daysInSIT int,
// 	lhDiscount unit.DiscountRate,
// 	sitDiscount unit.DiscountRate) (cost CostComputation, err error) {

// 	// Weights below 1000lbs are prorated to the 1000lb rate
// 	prorateFactor := 1.0
// 	if weight.Int() < 1000 {
// 		prorateFactor = weight.Float64() / 1000.0
// 		weight = unit.Pound(1000)
// 	}

// 	// Linehaul charges
// 	linehaulCostComputation, err := re.linehaulChargeComputation(weight, originZip5, destinationZip5, distanceMiles, date)
// 	if err != nil {
// 		re.logger.Error("Failed to compute linehaul cost", zap.Error(err))
// 		return
// 	}

// 	// Non linehaul charges
// 	nonLinehaulCostComputation, err := re.nonLinehaulChargeComputation(weight, originZip5, destinationZip5, date)
// 	if err != nil {
// 		re.logger.Error("Failed to compute non-linehaul cost", zap.Error(err))
// 		return
// 	}

// 	// Apply linehaul discounts
// 	linehaulCostComputation.LinehaulChargeTotal = lhDiscount.Apply(linehaulCostComputation.LinehaulChargeTotal)
// 	nonLinehaulCostComputation.OriginService.Fee = lhDiscount.Apply(nonLinehaulCostComputation.OriginService.Fee)
// 	nonLinehaulCostComputation.DestinationService.Fee = lhDiscount.Apply(nonLinehaulCostComputation.DestinationService.Fee)
// 	nonLinehaulCostComputation.Pack.Fee = lhDiscount.Apply(nonLinehaulCostComputation.Pack.Fee)
// 	nonLinehaulCostComputation.Unpack.Fee = lhDiscount.Apply(nonLinehaulCostComputation.Unpack.Fee)

// 	// SIT
// 	// Note that SIT has a different discount rate than [non]linehaul charges
// 	destinationZip3 := Zip5ToZip3(destinationZip5)
// 	sitComputation, err := re.SitCharge(weight.ToCWT(), daysInSIT, destinationZip3, date, true)
// 	if err != nil {
// 		re.logger.Info("Can't calculate sit",
// 			zap.String("moveLocator", re.move.Locator),
// 		)
// 		return
// 	}
// 	sitFee := sitComputation.ApplyDiscount(lhDiscount, sitDiscount)

// 	/// Max SIT
// 	maxSITComputation, err := re.SitCharge(weight.ToCWT(), MaxSITDays, destinationZip3, date, true)
// 	if err != nil {
// 		re.logger.Info("Can't calculate max sit",
// 			zap.String("moveLocator", re.move.Locator),
// 		)
// 		return
// 	}
// 	// Note that SIT has a different discount rate than [non]linehaul charges
// 	maxSITFee := maxSITComputation.ApplyDiscount(lhDiscount, sitDiscount)

// 	// Totals
// 	gcc := linehaulCostComputation.LinehaulChargeTotal +
// 		nonLinehaulCostComputation.OriginService.Fee +
// 		nonLinehaulCostComputation.DestinationService.Fee +
// 		nonLinehaulCostComputation.Pack.Fee +
// 		nonLinehaulCostComputation.Unpack.Fee

// 	cost = CostComputation{
// 		LinehaulCostComputation:    linehaulCostComputation,
// 		NonLinehaulCostComputation: nonLinehaulCostComputation,
// 		SITFee:                     sitFee,
// 		SITMax:                     maxSITFee,
// 		GCC:                        gcc,
// 		LHDiscount:                 lhDiscount,
// 		SITDiscount:                sitDiscount,
// 		Weight:                     weight,
// 	}

// 	// Finally, scale by prorate factor
// 	cost.Scale(prorateFactor)

// 	re.logger.Info("PPM cost computation",
// 		zap.String("moveLocator", re.move.Locator),
// 		zap.Object("cost", cost),
// 	)

// 	return cost, nil
// }

//computePPMIncludingLHDiscount Calculates the cost of a PPM move using zip + date derived linehaul discount
func (re *RateEngine) computePPMIncludingLHDiscount(appCtx appcontext.AppContext, weight unit.Pound, originZip5 string, destinationZip5 string, distanceMiles int, date time.Time, daysInSIT int) (cost CostComputation, err error) {

	lhDiscount, sitDiscount, err := models.PPMDiscountFetch(appCtx.DB(),
		appCtx.Logger(),
		re.move,
		originZip5,
		destinationZip5,
		date,
	)
	if err != nil {
		appCtx.Logger().Error("Failed to compute linehaul cost", zap.Error(err))
		return
	}

	cost, err = re.computePPM(appCtx,
		weight,
		originZip5,
		destinationZip5,
		distanceMiles,
		date,
		daysInSIT,
		lhDiscount,
		sitDiscount,
	)

	if err != nil {
		appCtx.Logger().Error("Failed to compute PPM cost", zap.Error(err))
		return
	}
	return cost, nil
}

// ComputePPMMoveCosts uses zip codes to make two calculations for the price of a PPM move - once with the pickup zip and once with the current duty location zip - and returns both calcs.
func (re *RateEngine) ComputePPMMoveCosts(appCtx appcontext.AppContext, weight unit.Pound, originPickupZip5 string, originDutyLocationZip5 string, destinationZip5 string, distanceMilesFromOriginPickupZip int, distanceMilesFromOriginDutyLocationZip int, date time.Time, daysInSit int) (costDetails CostDetails, err error) {
	costFromOriginPickupZip, err := re.computePPMIncludingLHDiscount(
		appCtx,
		weight,
		originPickupZip5,
		destinationZip5,
		distanceMilesFromOriginPickupZip,
		date,
		daysInSit,
	)
	if err != nil {
		appCtx.Logger().Error("Failed to compute PPM cost", zap.Error(err))
		return
	}

	costDetails = make(CostDetails)
	costDetails["pickupLocation"] = &CostDetail{
		costFromOriginPickupZip,
		false,
	}

	costFromOriginDutyLocationZip, err := re.computePPMIncludingLHDiscount(
		appCtx,
		weight,
		originDutyLocationZip5,
		destinationZip5,
		distanceMilesFromOriginDutyLocationZip,
		date,
		daysInSit,
	)
	if err != nil {
		appCtx.Logger().Error("Failed to compute PPM cost", zap.Error(err))
		return
	}
	costDetails["originDutyLocation"] = &CostDetail{
		costFromOriginDutyLocationZip,
		false,
	}

	originZipCode := originPickupZip5
	originZipLocation := "Pickup location"
	if costFromOriginPickupZip.GCC > costFromOriginDutyLocationZip.GCC {
		costDetails["originDutyLocation"].IsWinning = true
		originZipCode = originDutyLocationZip5
		originZipLocation = "Origin duty location"
	} else {
		costDetails["pickupLocation"].IsWinning = true
	}

	appCtx.Logger().Info("Origin zip code information",
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
	return costDetails["originDutyLocation"].Cost
}

// GetNonWinningCostMove returns a costComputation of the non-winning calculation
func GetNonWinningCostMove(costDetails CostDetails) CostComputation {
	if costDetails["pickupLocation"].IsWinning {
		return costDetails["originDutyLocation"].Cost
	}
	return costDetails["pickupLocation"].Cost
}

// NewRateEngine creates a new RateEngine
func NewRateEngine(move models.Move) *RateEngine {
	return &RateEngine{move: move}
}
