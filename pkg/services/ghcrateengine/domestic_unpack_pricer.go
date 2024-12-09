package ghcrateengine

import (
	"fmt"
	"time"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/unit"
)

type domesticUnpackPricer struct {
	services.FeatureFlagFetcher
}

// NewDomesticUnpackPricer creates a new pricer for the domestic unpack service
func NewDomesticUnpackPricer(featureFlagFetcher services.FeatureFlagFetcher) services.DomesticUnpackPricer {
	return &domesticUnpackPricer{featureFlagFetcher}
}

// Price determines the price for a domestic unpack service
func (p domesticUnpackPricer) Price(appCtx appcontext.AppContext, contractCode string, referenceDate time.Time, weight unit.Pound, servicesScheduleDest int, isPPM bool, isMobileHome bool) (unit.Cents, services.PricingDisplayParams, error) {
	return priceDomesticPackUnpack(appCtx, models.ReServiceCodeDUPK, contractCode, referenceDate, weight, servicesScheduleDest, isPPM, isMobileHome, p.FeatureFlagFetcher)
}

// Determines if this DUPK item should actually be added to the payment request by checking for relevant feature flags
func (p domesticUnpackPricer) ShouldPrice(appCtx appcontext.AppContext) (bool, error) {
	isOn, err := GetFeatureFlagValue(appCtx, p.FeatureFlagFetcher, services.DomesticMobileHomeUnpackingEnabled) // This should be edited later to also include the Boat Shipment FFs
	if err != nil {
		return false, fmt.Errorf("could not fetch feature flag to determine unpack pricing formula: %w", err)
	}
	return isOn, nil
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

	// Check if unpacking service items have been enabled for Mobile Home shipments
	isMobileHomePackingItemOn, err := GetFeatureFlagValue(appCtx, p.FeatureFlagFetcher, "domestic_mobile_home_unpacking_enabled")
	if err != nil {
		return unit.Cents(0), nil, err
	}

	var isMobileHome = false
	if isMobileHomePackingItemOn && params[0].PaymentServiceItem.MTOServiceItem.MTOShipment.ShipmentType == models.MTOShipmentTypeMobileHome {
		isMobileHome = true
	}

	return p.Price(appCtx, contractCode, referenceDate, unit.Pound(weightBilled), servicesScheduleDest, isPPM, isMobileHome)
}
