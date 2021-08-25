package ghcrateengine

import (
	"time"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/unit"
)

type domesticOriginSITPickupPricer struct {
}

// NewDomesticOriginSITPickupPricer creates a new pricer for domestic origin SIT pickup
func NewDomesticOriginSITPickupPricer() services.DomesticOriginSITPickupPricer {
	return &domesticOriginSITPickupPricer{}
}

// Price determines the price for domestic origin SIT pickup
func (p domesticOriginSITPickupPricer) Price(appCtx appcontext.AppContext, contractCode string, requestedPickupDate time.Time, weight unit.Pound, serviceArea string, sitSchedule int, zipSITOriginOriginal string, zipSITOriginActual string, distance unit.Miles) (unit.Cents, services.PricingDisplayParams, error) {
	return priceDomesticPickupDeliverySIT(appCtx, models.ReServiceCodeDOPSIT, contractCode, requestedPickupDate, weight, serviceArea, sitSchedule, zipSITOriginOriginal, zipSITOriginActual, distance)
}

// PriceUsingParams determines the price for domestic origin SIT pickup given PaymentServiceItemParams
func (p domesticOriginSITPickupPricer) PriceUsingParams(appCtx appcontext.AppContext, params models.PaymentServiceItemParams) (unit.Cents, services.PricingDisplayParams, error) {
	contractCode, err := getParamString(params, models.ServiceItemParamNameContractCode)
	if err != nil {
		return unit.Cents(0), nil, err
	}

	requestedPickupDate, err := getParamTime(params, models.ServiceItemParamNameRequestedPickupDate)
	if err != nil {
		return unit.Cents(0), nil, err
	}

	weightBilledActual, err := getParamInt(params, models.ServiceItemParamNameWeightBilledActual)
	if err != nil {
		return unit.Cents(0), nil, err
	}

	serviceAreaOrigin, err := getParamString(params, models.ServiceItemParamNameServiceAreaOrigin)
	if err != nil {
		return unit.Cents(0), nil, err
	}

	sitScheduleOrigin, err := getParamInt(params, models.ServiceItemParamNameSITScheduleOrigin)
	if err != nil {
		return unit.Cents(0), nil, err
	}

	zipSITOriginOriginalAddress, err := getParamString(params, models.ServiceItemParamNameZipSITOriginHHGOriginalAddress)
	if err != nil {
		return unit.Cents(0), nil, err
	}

	zipSITOriginActualAddress, err := getParamString(params, models.ServiceItemParamNameZipSITOriginHHGActualAddress)
	if err != nil {
		return unit.Cents(0), nil, err
	}

	distanceZipSITOrigin, err := getParamInt(params, models.ServiceItemParamNameDistanceZipSITOrigin)
	if err != nil {
		return unit.Cents(0), nil, err
	}

	return p.Price(appCtx, contractCode, requestedPickupDate, unit.Pound(weightBilledActual), serviceAreaOrigin,
		sitScheduleOrigin, zipSITOriginOriginalAddress, zipSITOriginActualAddress, unit.Miles(distanceZipSITOrigin))
}
