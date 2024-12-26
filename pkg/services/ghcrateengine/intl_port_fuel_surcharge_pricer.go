package ghcrateengine

import (
	"database/sql"
	"fmt"
	"math"
	"time"

	"github.com/gofrs/uuid"
	"github.com/pkg/errors"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/unit"
)

// FuelSurchargePricer is a service object to price domestic shorthaul
type portFuelSurchargePricer struct {
}

// NewFuelSurchargePricer is the public constructor for a domesticFuelSurchargePricer using Pop
func NewPortFuelSurchargePricer() services.IntlPortFuelSurchargePricer {
	return &portFuelSurchargePricer{}
}

// Price determines the price for fuel surcharge
func (p portFuelSurchargePricer) Price(_ appcontext.AppContext, actualPickupDate time.Time, distance unit.Miles, weight unit.Pound, fscWeightBasedDistanceMultiplier float64, eiaFuelPrice unit.Millicents, isPPM bool) (unit.Cents, services.PricingDisplayParams, error) {
	// Validate parameters
	if actualPickupDate.IsZero() {
		return 0, nil, errors.New("ActualPickupDate is required")
	}
	if distance <= 0 {
		return 0, nil, errors.New("Distance must be greater than 0")
	}
	if !isPPM && weight < minIntlWeightHHG {
		return 0, nil, fmt.Errorf("weight must be a minimum of %d", minIntlWeightHHG)
	}
	if fscWeightBasedDistanceMultiplier == 0 {
		return 0, nil, errors.New("WeightBasedDistanceMultiplier is required")
	}
	if eiaFuelPrice == 0 {
		return 0, nil, errors.New("EIAFuelPrice is required")
	}

	fscPriceDifferenceInCents := (eiaFuelPrice - baseGHCDieselFuelPrice).Float64() / 1000.0
	fscMultiplier := fscWeightBasedDistanceMultiplier * distance.Float64()
	fscPrice := fscMultiplier * fscPriceDifferenceInCents * 100
	totalCost := unit.Cents(math.Round(fscPrice))

	displayParams := services.PricingDisplayParams{
		{Key: models.ServiceItemParamNameFSCPriceDifferenceInCents, Value: FormatFloat(fscPriceDifferenceInCents, 1)},
		{Key: models.ServiceItemParamNameFSCMultiplier, Value: FormatFloat(fscMultiplier, 7)},
	}

	return totalCost, displayParams, nil
}

func (p portFuelSurchargePricer) PriceUsingParams(appCtx appcontext.AppContext, params models.PaymentServiceItemParams) (unit.Cents, services.PricingDisplayParams, error) {
	actualPickupDate, err := getParamTime(params, models.ServiceItemParamNameActualPickupDate)
	if err != nil {
		return unit.Cents(0), nil, err
	}

	var paymentServiceItem models.PaymentServiceItem
	mtoShipment := params[0].PaymentServiceItem.MTOServiceItem.MTOShipment

	if mtoShipment.ID == uuid.Nil {
		err = appCtx.DB().Eager("MTOServiceItem", "MTOServiceItem.MTOShipment").Find(&paymentServiceItem, params[0].PaymentServiceItemID)
		if err != nil {
			switch err {
			case sql.ErrNoRows:
				return unit.Cents(0), nil, apperror.NewNotFoundError(params[0].PaymentServiceItemID, "looking for PaymentServiceItem")
			default:
				return unit.Cents(0), nil, apperror.NewQueryError("PaymentServiceItem", err, "")
			}
		}
		mtoShipment = paymentServiceItem.MTOServiceItem.MTOShipment
	}

	distance, err := getParamInt(params, models.ServiceItemParamNameDistanceZip)
	if err != nil {
		return unit.Cents(0), nil, err
	}

	weightBilled, err := getParamInt(params, models.ServiceItemParamNameWeightBilled)
	if err != nil {
		return unit.Cents(0), nil, err
	}

	fscWeightBasedDistanceMultiplier, err := getParamFloat(params, models.ServiceItemParamNameFSCWeightBasedDistanceMultiplier)
	if err != nil {
		return unit.Cents(0), nil, err
	}

	eiaFuelPrice, err := getParamInt(params, models.ServiceItemParamNameEIAFuelPrice)
	if err != nil {
		return unit.Cents(0), nil, err
	}

	var isPPM = false
	if params[0].PaymentServiceItem.MTOServiceItem.MTOShipment.ShipmentType == models.MTOShipmentTypePPM {
		// PPMs do not require minimums for a shipment's weight
		// this flag is passed into the Price function to ensure the weight min
		// are not enforced for PPMs
		isPPM = true
	}

	return p.Price(appCtx, actualPickupDate, unit.Miles(distance), unit.Pound(weightBilled), fscWeightBasedDistanceMultiplier, unit.Millicents(eiaFuelPrice), isPPM)
}
