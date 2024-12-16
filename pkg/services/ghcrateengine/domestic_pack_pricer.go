package ghcrateengine

import (
	"time"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/unit"
)

type domesticPackPricer struct {
	services.FeatureFlagFetcher
}

// NewDomesticPackPricer creates a new pricer for the domestic pack service
func NewDomesticPackPricer(featureFlagFetcher services.FeatureFlagFetcher) services.DomesticPackPricer {
	return &domesticPackPricer{featureFlagFetcher}
}

// Price determines the price for a domestic pack service
func (p domesticPackPricer) Price(appCtx appcontext.AppContext, contractCode string, referenceDate time.Time, weight unit.Pound, servicesScheduleOrigin int, isPPM bool, isMobileHome bool) (unit.Cents, services.PricingDisplayParams, error) {
	return priceDomesticPackUnpack(appCtx, models.ReServiceCodeDPK, contractCode, referenceDate, weight, servicesScheduleOrigin, isPPM, isMobileHome, p.FeatureFlagFetcher)
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

	var isPPM = false
	if params[0].PaymentServiceItem.MTOServiceItem.MTOShipment.ShipmentType == models.MTOShipmentTypePPM {
		// PPMs do not require minimums for a shipment's weight
		// this flag is passed into the Price function to ensure the weight min
		// are not enforced for PPMs
		isPPM = true
	}

	// Check if packing service items have been enabled for Mobile Home shipments
	isMobileHomePackingItemOn, err := getFeatureFlagValue(appCtx, p.FeatureFlagFetcher, services.DomesticMobileHomePackingEnabled)
	if err != nil {
		return unit.Cents(0), nil, err
	}

	var isMobileHome = false
	if isMobileHomePackingItemOn && params[0].PaymentServiceItem.MTOServiceItem.MTOShipment.ShipmentType == models.MTOShipmentTypeMobileHome {
		isMobileHome = true
	}

	return p.Price(appCtx, contractCode, referenceDate, unit.Pound(weightBilled), servicesScheduleOrigin, isPPM, isMobileHome)
}
