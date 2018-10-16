package rateengine

import (
	"github.com/gobuffalo/uuid"
	"github.com/pkg/errors"
	"github.com/transcom/mymove/pkg/models"
	"time"

	"github.com/gobuffalo/pop"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/transcom/mymove/pkg/route"
	"github.com/transcom/mymove/pkg/unit"
)

// MaxSITDays is the maximum number of days of SIT that will be reimbursed.
const MaxSITDays = 90

// RateEngine encapsulates the TSP rate engine process
type RateEngine struct {
	db      *pop.Connection
	logger  *zap.Logger
	planner route.Planner
}

// CostComputation represents the results of a computation.
type CostComputation struct {
	LinehaulCostComputation
	NonLinehaulCostComputation
	SITFee unit.Cents
	SITMax unit.Cents
	GCC    unit.Cents
}

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
	encoder.AddInt("BaseLinehaul", c.BaseLinehaul.Int())
	encoder.AddInt("OriginLinehaulFactor", c.OriginLinehaulFactor.Int())
	encoder.AddInt("DestinationLinehaulFactor", c.DestinationLinehaulFactor.Int())
	encoder.AddInt("ShorthaulCharge", c.ShorthaulCharge.Int())
	encoder.AddInt("LinehaulChargeTotal", c.LinehaulChargeTotal.Int())

	encoder.AddInt("OriginServiceFee", c.OriginServiceFee.Int())
	encoder.AddInt("DestinationServiceFee", c.DestinationServiceFee.Int())
	encoder.AddInt("PackFee", c.PackFee.Int())
	encoder.AddInt("UnpackFee", c.UnpackFee.Int())
	encoder.AddInt("SITMax", c.SITMax.Int())
	encoder.AddInt("SITFee", c.SITFee.Int())

	encoder.AddInt("GCC", c.GCC.Int())

	return nil
}

// Zip5ToZip3 takes a ZIP5 string and returns the ZIP3 representation of it.
func Zip5ToZip3(zip5 string) string {
	return zip5[0:3]
}

// ComputePPM Calculates the cost of a PPM move.
func (re *RateEngine) ComputePPM(
	weight unit.Pound,
	originZip5 string,
	destinationZip5 string,
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
	linehaulCostComputation, err := re.linehaulChargeComputation(weight, originZip5, destinationZip5, date)
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
	nonLinehaulCostComputation.OriginServiceFee = lhDiscount.Apply(nonLinehaulCostComputation.OriginServiceFee)
	nonLinehaulCostComputation.DestinationServiceFee = lhDiscount.Apply(nonLinehaulCostComputation.DestinationServiceFee)
	nonLinehaulCostComputation.PackFee = lhDiscount.Apply(nonLinehaulCostComputation.PackFee)
	nonLinehaulCostComputation.UnpackFee = lhDiscount.Apply(nonLinehaulCostComputation.UnpackFee)

	// SIT
	// Note that SIT has a different discount rate than [non]linehaul charges
	destinationZip3 := Zip5ToZip3(destinationZip5)
	sit, err := re.SitCharge(weight.ToCWT(), daysInSIT, destinationZip3, date, true)
	if err != nil {
		re.logger.Info("Can't calculate sit")
		return
	}
	sitFee := sitDiscount.Apply(sit)

	/// Max SIT
	maxSIT, err := re.SitCharge(weight.ToCWT(), MaxSITDays, destinationZip3, date, true)
	if err != nil {
		re.logger.Info("Can't calculate max sit")
		return
	}
	// Note that SIT has a different discount rate than [non]linehaul charges
	maxSITFee := sitDiscount.Apply(maxSIT)

	// Totals
	gcc := linehaulCostComputation.LinehaulChargeTotal +
		nonLinehaulCostComputation.OriginServiceFee +
		nonLinehaulCostComputation.DestinationServiceFee +
		nonLinehaulCostComputation.PackFee +
		nonLinehaulCostComputation.UnpackFee

	cost = CostComputation{
		LinehaulCostComputation:    linehaulCostComputation,
		NonLinehaulCostComputation: nonLinehaulCostComputation,
		SITFee: sitFee,
		SITMax: maxSITFee,
		GCC:    gcc,
	}

	// Finally, scale by prorate factor
	cost.Scale(prorateFactor)

	re.logger.Info("PPM cost computation", zap.Object("cost", cost))

	return cost, nil
}

// ComputeShipment Calculates the cost of an HHG move.
func (re *RateEngine) ComputeShipment(
	weight unit.Pound,
	originZip5 string,
	destinationZip5 string,
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
	linehaulCostComputation, err := re.linehaulChargeComputation(weight, originZip5, destinationZip5, date)
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
	nonLinehaulCostComputation.OriginServiceFee = lhDiscount.Apply(nonLinehaulCostComputation.OriginServiceFee)
	nonLinehaulCostComputation.DestinationServiceFee = lhDiscount.Apply(nonLinehaulCostComputation.DestinationServiceFee)
	nonLinehaulCostComputation.PackFee = lhDiscount.Apply(nonLinehaulCostComputation.PackFee)
	nonLinehaulCostComputation.UnpackFee = lhDiscount.Apply(nonLinehaulCostComputation.UnpackFee)

	// SIT
	// Note that SIT has a different discount rate than [non]linehaul charges
	destinationZip3 := Zip5ToZip3(destinationZip5)
	sit, err := re.SitCharge(weight.ToCWT(), daysInSIT, destinationZip3, date, true)
	if err != nil {
		re.logger.Info("Can't calculate sit")
		return
	}
	sitFee := sitDiscount.Apply(sit)

	/// Max SIT
	maxSIT, err := re.SitCharge(weight.ToCWT(), MaxSITDays, destinationZip3, date, true)
	if err != nil {
		re.logger.Info("Can't calculate max sit")
		return
	}
	// Note that SIT has a different discount rate than [non]linehaul charges
	maxSITFee := sitDiscount.Apply(maxSIT)

	// Totals
	gcc := linehaulCostComputation.LinehaulChargeTotal +
		nonLinehaulCostComputation.OriginServiceFee +
		nonLinehaulCostComputation.DestinationServiceFee +
		nonLinehaulCostComputation.PackFee +
		nonLinehaulCostComputation.UnpackFee

	cost = CostComputation{
		LinehaulCostComputation:    linehaulCostComputation,
		NonLinehaulCostComputation: nonLinehaulCostComputation,
		SITFee: sitFee,
		SITMax: maxSITFee,
		GCC:    gcc,
	}

	// Finally, scale by prorate factor
	cost.Scale(prorateFactor)

	re.logger.Info("PPM cost computation", zap.Object("cost", cost))

	return cost, nil
}

// CostByShipment struct containing shipment and cost
type CostByShipment struct {
	Shipment models.Shipment
	Cost     CostComputation
}

// HandleRunOnShipment runs the rate engine on a shipment and returns the shipment and cost.
// Assumptions: Shipment model passed in has eagerly fetched PickupAddress,
// Move.Orders.NewDutyStation.Address, and ShipmentOffers.TransportationServiceProviderPerformance.
func (re *RateEngine) HandleRunOnShipment(shipment models.Shipment) (CostByShipment, error) {
	// Validate expected model relationships are available.
	if shipment.PickupAddress == nil {
		return CostByShipment{}, errors.New("PickupAddress is nil")
	}

	// NewDutyStation's address/postal code is required per model/schema, so no nil check needed.

	if shipment.ShipmentOffers == nil {
		return CostByShipment{}, errors.New("ShipmentOffers is nil")
	} else if len(shipment.ShipmentOffers) == 0 {
		return CostByShipment{}, errors.New("ShipmentOffers fetched, but none found")
	}

	if shipment.ShipmentOffers[0].TransportationServiceProviderPerformance.ID == uuid.Nil {
		return CostByShipment{}, errors.New("TransportationServiceProviderPerformance is nil")
	}

	if shipment.NetWeight == nil {
		return CostByShipment{}, errors.New("NetWeight is nil")
	}

	// All required relationships should exist at this point.
	daysInSIT := 0
	var sitDiscount unit.DiscountRate
	sitDiscount = 0.0

	// Assume the most recent matching shipment offer is the right one.
	lhDiscount := shipment.ShipmentOffers[0].TransportationServiceProviderPerformance.LinehaulRate

	// Apply rate engine to shipment
	var shipmentCost CostByShipment
	cost, err := re.ComputeShipment(*shipment.NetWeight,
		shipment.PickupAddress.PostalCode,
		shipment.Move.Orders.NewDutyStation.Address.PostalCode,
		time.Time(*shipment.ActualPickupDate),
		daysInSIT, // We don't want any SIT charges
		lhDiscount,
		sitDiscount,
	)
	if err != nil {
		return CostByShipment{}, err
	}

	shipmentCost = CostByShipment{
		Shipment: shipment,
		Cost:     cost,
	}
	return shipmentCost, err
}

// NewRateEngine creates a new RateEngine
func NewRateEngine(db *pop.Connection, logger *zap.Logger, planner route.Planner) *RateEngine {
	return &RateEngine{db: db, logger: logger, planner: planner}
}
