package ghcrateengine

import (
	"time"

	"github.com/gobuffalo/pop/v5"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/unit"
)

type domesticUncratingPricer struct {
	db *pop.Connection
}

// NewDomesticUncratingPricer creates a new pricer for domestic destination first day SIT
func NewDomesticUncratingPricer(db *pop.Connection) services.DomesticUncratingPricer {
	return &domesticUncratingPricer{
		db: db,
	}
}

// Price determines the price for domestic destination first day SIT
func (p domesticUncratingPricer) Price(contractCode string, requestedPickupDate time.Time, billedCubicFeet unit.CubicFeet, serviceSchedule int) (unit.Cents, services.PricingDisplayParams, error) {
	return priceDomesticCrating(p.db, models.ReServiceCodeDUCRT, contractCode, requestedPickupDate, billedCubicFeet, serviceSchedule)
}

// PriceUsingParams determines the price for domestic destination first day SIT given PaymentServiceItemParams
func (p domesticUncratingPricer) PriceUsingParams(params models.PaymentServiceItemParams) (unit.Cents, services.PricingDisplayParams, error) {
	contractCode, err := getParamString(params, models.ServiceItemParamNameContractCode)
	if err != nil {
		return unit.Cents(0), nil, err
	}

	cubicFeetFloat, err := getParamFloat(params, models.ServiceItemParamNameCubicFeetBilled)
	if err != nil {
		return unit.Cents(0), nil, err
	}

	cubicFeetBilled := unit.CubicFeet(cubicFeetFloat)

	requestedPickupDate, err := getParamTime(params, models.ServiceItemParamNameRequestedPickupDate)
	if err != nil {
		return unit.Cents(0), nil, err
	}

	serviceScheduleDestination, err := getParamInt(params, models.ServiceItemParamNameServicesScheduleDest)
	if err != nil {
		return unit.Cents(0), nil, err
	}

	return p.Price(contractCode, requestedPickupDate, cubicFeetBilled, serviceScheduleDestination)
}
