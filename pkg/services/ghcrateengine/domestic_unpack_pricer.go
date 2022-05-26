package ghcrateengine

import (
	"time"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"

	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/unit"
)

type domesticUnpackPricer struct {
}

// NewDomesticUnpackPricer creates a new pricer for the domestic unpack service
func NewDomesticUnpackPricer() services.DomesticUnpackPricer {
	return &domesticUnpackPricer{}
}

// Price determines the price for a domestic unpack service
func (p domesticUnpackPricer) Price(appCtx appcontext.AppContext, contractCode string, referenceDate time.Time, weight unit.Pound, servicesScheduleDest int, isPPM bool) (unit.Cents, services.PricingDisplayParams, error) {
	return priceDomesticPackUnpack(appCtx, models.ReServiceCodeDUPK, contractCode, referenceDate, weight, servicesScheduleDest, isPPM)
}

// PriceUsingParams determines the price for a domestic unpack service given PaymentServiceItemParams
func (p domesticUnpackPricer) PriceUsingParams(appCtx appcontext.AppContext, params models.PaymentServiceItemParams) (unit.Cents, services.PricingDisplayParams, error) {
	contractCode, err := getParamString(params, models.ServiceItemParamNameContractCode)
	if err != nil {
		return unit.Cents(0), nil, err
	}

	referenceDate, err := getParamTime(params, models.ServiceItemParamNameReferenceDate)
	if err != nil {
		return unit.Cents(0), nil, err
	}

	servicesScheduleDest, err := getParamInt(params, models.ServiceItemParamNameServicesScheduleDest)
	if err != nil {
		return unit.Cents(0), nil, err
	}

	weightBilled, err := getParamInt(params, models.ServiceItemParamNameWeightBilled)
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

	return p.Price(appCtx, contractCode, referenceDate, unit.Pound(weightBilled), servicesScheduleDest, isPPM)
}
