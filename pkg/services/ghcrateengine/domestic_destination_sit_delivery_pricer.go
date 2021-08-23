package ghcrateengine

import (
	"time"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/unit"
)

type domesticDestinationSITDeliveryPricer struct {
}

// NewDomesticDestinationSITDeliveryPricer creates a new pricer for domestic destination SIT delivery
func NewDomesticDestinationSITDeliveryPricer() services.DomesticDestinationSITDeliveryPricer {
	return &domesticDestinationSITDeliveryPricer{}
}

// Price determines the price for domestic destination SIT delivery
func (p domesticDestinationSITDeliveryPricer) Price(appCtx appcontext.AppContext, contractCode string, requestedPickupDate time.Time, weight unit.Pound, serviceArea string, sitSchedule int, zipDest string, zipSITDest string, distance unit.Miles) (unit.Cents, services.PricingDisplayParams, error) {
	return priceDomesticPickupDeliverySIT(appCtx, models.ReServiceCodeDDDSIT, contractCode, requestedPickupDate, weight, serviceArea, sitSchedule, zipDest, zipSITDest, distance)
}

// PriceUsingParams determines the price for domestic destination SIT delivery given PaymentServiceItemParams
func (p domesticDestinationSITDeliveryPricer) PriceUsingParams(appCtx appcontext.AppContext, params models.PaymentServiceItemParams) (unit.Cents, services.PricingDisplayParams, error) {
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

	serviceAreaDest, err := getParamString(params, models.ServiceItemParamNameServiceAreaDest)
	if err != nil {
		return unit.Cents(0), nil, err
	}

	sitScheduleDest, err := getParamInt(params, models.ServiceItemParamNameSITScheduleDest)
	if err != nil {
		return unit.Cents(0), nil, err
	}

	zipDestAddress, err := getParamString(params, models.ServiceItemParamNameZipDestAddress)
	if err != nil {
		return unit.Cents(0), nil, err
	}

	zipSITDestHHGFinalAddress, err := getParamString(params, models.ServiceItemParamNameZipSITDestHHGFinalAddress)
	if err != nil {
		return unit.Cents(0), nil, err
	}

	distanceZipSITDest, err := getParamInt(params, models.ServiceItemParamNameDistanceZipSITDest)
	if err != nil {
		return unit.Cents(0), nil, err
	}

	return p.Price(appCtx, contractCode, requestedPickupDate, unit.Pound(weightBilledActual), serviceAreaDest,
		sitScheduleDest, zipDestAddress, zipSITDestHHGFinalAddress, unit.Miles(distanceZipSITDest))
}
