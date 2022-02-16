package ghcrateengine

import (
	"time"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/unit"
)

type domesticUncratingPricer struct {
}

// NewDomesticUncratingPricer creates a new pricer for domestic destination first day SIT
func NewDomesticUncratingPricer() services.DomesticUncratingPricer {
	return &domesticUncratingPricer{}
}

// Price determines the price for domestic destination first day SIT
func (p domesticUncratingPricer) Price(appCtx appcontext.AppContext, contractCode string, referenceDate time.Time, billedCubicFeet unit.CubicFeet, serviceSchedule int) (unit.Cents, services.PricingDisplayParams, error) {
	return priceDomesticCrating(appCtx, models.ReServiceCodeDUCRT, contractCode, referenceDate, billedCubicFeet, serviceSchedule)
}

// PriceUsingParams determines the price for domestic destination first day SIT given PaymentServiceItemParams
func (p domesticUncratingPricer) PriceUsingParams(appCtx appcontext.AppContext, params models.PaymentServiceItemParams) (unit.Cents, services.PricingDisplayParams, error) {
	contractCode, err := getParamString(params, models.ServiceItemParamNameContractCode)
	if err != nil {
		return unit.Cents(0), nil, err
	}

	cubicFeetFloat, err := getParamFloat(params, models.ServiceItemParamNameCubicFeetBilled)
	if err != nil {
		return unit.Cents(0), nil, err
	}

	cubicFeetBilled := unit.CubicFeet(cubicFeetFloat)

	referenceDate, err := getParamTime(params, models.ServiceItemParamNameReferenceDate)
	if err != nil {
		return unit.Cents(0), nil, err
	}

	serviceScheduleDestination, err := getParamInt(params, models.ServiceItemParamNameServicesScheduleDest)
	if err != nil {
		return unit.Cents(0), nil, err
	}

	return p.Price(appCtx, contractCode, referenceDate, cubicFeetBilled, serviceScheduleDestination)
}
