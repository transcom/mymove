package ghcrateengine

import (
	"fmt"
	"math"
	"time"

	"github.com/pkg/errors"

	"github.com/gobuffalo/pop/v5"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/unit"
)

const baseGHCDieselFuelPrice = unit.Millicents(250000)

// FuelSurchargePricer is a service object to price domestic shorthaul
type fuelSurchargePricer struct {
	db *pop.Connection
}

// NewFuelSurchargePricer is the public constructor for a domesticFuelSurchargePricer using Pop
func NewFuelSurchargePricer(db *pop.Connection) services.FuelSurchargePricer {
	return &fuelSurchargePricer{
		db: db,
	}
}

// Price determines the price for a counseling service
func (p fuelSurchargePricer) Price(contractCode string, actualPickupDate time.Time, distance unit.Miles, weight unit.Pound, weightBasedDistanceMultiplier float64, fuelPrice unit.Millicents) (unit.Cents, services.PricingDisplayParams, error) {
	// Validate parameters
	if len(contractCode) == 0 {
		return 0, nil, errors.New("ContractCode is required")
	}
	if actualPickupDate.IsZero() {
		return 0, nil, errors.New("RequestedPickupDate is required")
	}
	if weight < minDomesticWeight {
		return 0, nil, fmt.Errorf("Weight must be a minimum of %d", minDomesticWeight)
	}
	if distance <= 0 {
		return 0, nil, errors.New("Distance must be greater than 0")
	}
	if weightBasedDistanceMultiplier == 0 {
		return 0, nil, errors.New("WeightBasedDistanceMultiplier is required")
	}

	priceDifference := (fuelPrice - baseGHCDieselFuelPrice).Float64() / 1000.00
	surchargeMultiplier := weightBasedDistanceMultiplier * distance.Float64()
	fscPrice := surchargeMultiplier * priceDifference * 100
	totalCost := unit.Cents(math.Round(fscPrice))

	return totalCost, nil, nil
}

func (p fuelSurchargePricer) PriceUsingParams(params models.PaymentServiceItemParams) (unit.Cents, services.PricingDisplayParams, error) {
	contractCode, err := getParamString(params, models.ServiceItemParamNameContractCode)
	if err != nil {
		return unit.Cents(0), nil, err
	}

	actualPickupDate, err := getParamTime(params, models.ServiceItemParamNameActualPickupDate)
	if err != nil {
		return unit.Cents(0), nil, err
	}

	var paymentServiceItem models.PaymentServiceItem
	err = p.db.Eager("MTOServiceItem", "MTOServiceItem.MTOShipment").Find(&paymentServiceItem, params[0].PaymentServiceItemID)
	if err != nil {
		return unit.Cents(0), nil, err
	}

	mtoShipment := paymentServiceItem.MTOServiceItem.MTOShipment
	distance := *mtoShipment.Distance

	weightBilledActual, err := getParamInt(params, models.ServiceItemParamNameWeightBilledActual)
	if err != nil {
		return unit.Cents(0), nil, err
	}

	weightBasedDistanceMultiplier, err := getParamFloat(params, models.ServiceItemParamNameFSCWeightBasedDistanceMultiplier)
	if err != nil {
		return unit.Cents(0), nil, err
	}

	fuelPrice, err := getParamInt(params, models.ServiceItemParamNameEIAFuelPrice)
	if err != nil {
		return unit.Cents(0), nil, err
	}

	return p.Price(contractCode, actualPickupDate, distance, unit.Pound(weightBilledActual), weightBasedDistanceMultiplier, unit.Millicents(fuelPrice))
}
