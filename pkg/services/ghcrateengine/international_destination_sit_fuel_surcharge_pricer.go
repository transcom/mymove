package ghcrateengine

import (
	"database/sql"
	"time"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/unit"
)

type internationalDestinationSITFuelSurchargePricer struct {
}

func NewInternationalDestinationSITFuelSurchargePricer() services.InternationalDestinationSITFuelSurchargePricer {
	return &internationalDestinationSITFuelSurchargePricer{}
}

func (p internationalDestinationSITFuelSurchargePricer) Price(appCtx appcontext.AppContext, actualPickupDate time.Time, distance unit.Miles, weight unit.Pound, fscWeightBasedDistanceMultiplier float64, eiaFuelPrice unit.Millicents) (unit.Cents, services.PricingDisplayParams, error) {
	return priceIntlFuelSurchargeSIT(appCtx, models.ReServiceCodeIDSFSC, actualPickupDate, distance, weight, fscWeightBasedDistanceMultiplier, eiaFuelPrice)
}

func (p internationalDestinationSITFuelSurchargePricer) PriceUsingParams(appCtx appcontext.AppContext, params models.PaymentServiceItemParams) (unit.Cents, services.PricingDisplayParams, error) {
	actualPickupDate, err := getParamTime(params, models.ServiceItemParamNameActualPickupDate)
	if err != nil {
		return unit.Cents(0), nil, err
	}

	var paymentServiceItem models.PaymentServiceItem

	err = appCtx.DB().Eager("MTOServiceItem", "MTOServiceItem.SITDestinationFinalAddress").Find(&paymentServiceItem, params[0].PaymentServiceItemID)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return unit.Cents(0), nil, apperror.NewNotFoundError(params[0].PaymentServiceItemID, "looking for PaymentServiceItem")
		default:
			return unit.Cents(0), nil, apperror.NewQueryError("PaymentServiceItem", err, "")
		}
	}

	// do not calculate mileage if destination address is OCONUS. this is to prevent pricing to be calculated.
	distance := 0
	if paymentServiceItem.MTOServiceItem.SITDestinationFinalAddress != nil &&
		!*paymentServiceItem.MTOServiceItem.SITDestinationFinalAddress.IsOconus {
		distance, err = getParamInt(params, models.ServiceItemParamNameDistanceZipSITDest)
		if err != nil {
			return unit.Cents(0), nil, err
		}
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

	return p.Price(appCtx, actualPickupDate, unit.Miles(distance), unit.Pound(weightBilled), fscWeightBasedDistanceMultiplier, unit.Millicents(eiaFuelPrice))
}
