package ghcrateengine

import (
	"time"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/unit"
)

type internationalOriginSITPickupPricer struct {
}

// NewInternationalOriginSITPickupPricer creates a new pricer for international origin SIT pickup
func NewInternationalOriginSITPickupPricer() services.InternationalOriginSITPickupPricer {
	return &internationalOriginSITPickupPricer{}
}

// Price determines the price for international origin SIT pickup
func (p internationalOriginSITPickupPricer) Price(appCtx appcontext.AppContext, contractCode string, referenceDate time.Time, weight unit.Pound, perUnitCents int, distance int) (unit.Cents, services.PricingDisplayParams, error) {
	return priceIntlPickupDeliverySIT(appCtx, models.ReServiceCodeIOPSIT, contractCode, referenceDate, weight, perUnitCents, distance)
}

// PriceUsingParams determines the price for international origin SIT pickup given PaymentServiceItemParams
func (p internationalOriginSITPickupPricer) PriceUsingParams(appCtx appcontext.AppContext, params models.PaymentServiceItemParams) (unit.Cents, services.PricingDisplayParams, error) {
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

	distance, err := getParamInt(params, models.ServiceItemParamNameDistanceZipSITOrigin)
	if err != nil {
		return unit.Cents(0), nil, err
	}

	return p.Price(appCtx, contractCode, referenceDate, unit.Pound(weightBilled), perUnitCents, distance)
}
