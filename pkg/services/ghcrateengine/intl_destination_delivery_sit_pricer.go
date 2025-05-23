package ghcrateengine

import (
	"time"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/unit"
)

type internationalDestinationSITDeliveryPricer struct {
}

// NewInternationalDestinationSITDeliveryPricer creates a new pricer for international destination SIT delivery
func NewInternationalDestinationSITDeliveryPricer() services.InternationalDestinationSITDeliveryPricer {
	return &internationalDestinationSITDeliveryPricer{}
}

// Price determines the price for international destination SIT delivery
func (p internationalDestinationSITDeliveryPricer) Price(appCtx appcontext.AppContext, contractCode string, referenceDate time.Time, weight unit.Pound, perUnitCents int, distance int) (unit.Cents, services.PricingDisplayParams, error) {
	return priceIntlPickupDeliverySIT(appCtx, models.ReServiceCodeIDDSIT, contractCode, referenceDate, weight, perUnitCents, distance)
}

// PriceUsingParams determines the price for international destination SIT delivery given PaymentServiceItemParams
func (p internationalDestinationSITDeliveryPricer) PriceUsingParams(appCtx appcontext.AppContext, params models.PaymentServiceItemParams) (unit.Cents, services.PricingDisplayParams, error) {
	contractCode, err := getParamString(params, models.ServiceItemParamNameContractCode)
	if err != nil {
		return unit.Cents(0), nil, err
	}

	referenceDate, err := getParamTime(params, models.ServiceItemParamNameReferenceDate)
	if err != nil {
		return unit.Cents(0), nil, err
	}

	perUnitCents, err := getParamInt(params, models.ServiceItemParamNamePerUnitCents)
	if err != nil {
		return unit.Cents(0), nil, err
	}

	weightBilled, err := getParamInt(params, models.ServiceItemParamNameWeightBilled)
	if err != nil {
		return unit.Cents(0), nil, err
	}

	distance, err := getParamInt(params, models.ServiceItemParamNameDistanceZipSITDest)
	if err != nil {
		return unit.Cents(0), nil, err
	}

	return p.Price(appCtx, contractCode, referenceDate, unit.Pound(weightBilled), perUnitCents, distance)
}
