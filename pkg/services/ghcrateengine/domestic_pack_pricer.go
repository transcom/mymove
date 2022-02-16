package ghcrateengine

import (
	"time"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"

	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/unit"
)

type domesticPackPricer struct {
}

// NewDomesticPackPricer creates a new pricer for the domestic pack service
func NewDomesticPackPricer() services.DomesticPackPricer {
	return &domesticPackPricer{}
}

// Price determines the price for a domestic pack service
func (p domesticPackPricer) Price(appCtx appcontext.AppContext, contractCode string, referenceDate time.Time, weight unit.Pound, servicesScheduleOrigin int) (unit.Cents, services.PricingDisplayParams, error) {
	return priceDomesticPackUnpack(appCtx, models.ReServiceCodeDPK, contractCode, referenceDate, weight, servicesScheduleOrigin)
}

// PriceUsingParams determines the price for a domestic pack service given PaymentServiceItemParams
func (p domesticPackPricer) PriceUsingParams(appCtx appcontext.AppContext, params models.PaymentServiceItemParams) (unit.Cents, services.PricingDisplayParams, error) {
	contractCode, err := getParamString(params, models.ServiceItemParamNameContractCode)
	if err != nil {
		return unit.Cents(0), nil, err
	}

	referenceDate, err := getParamTime(params, models.ServiceItemParamNameReferenceDate)
	if err != nil {
		return unit.Cents(0), nil, err
	}

	servicesScheduleOrigin, err := getParamInt(params, models.ServiceItemParamNameServicesScheduleOrigin)
	if err != nil {
		return unit.Cents(0), nil, err
	}

	weightBilled, err := getParamInt(params, models.ServiceItemParamNameWeightBilled)
	if err != nil {
		return unit.Cents(0), nil, err
	}

	return p.Price(appCtx, contractCode, referenceDate, unit.Pound(weightBilled), servicesScheduleOrigin)
}
