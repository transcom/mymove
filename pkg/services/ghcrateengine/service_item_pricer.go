package ghcrateengine

import (
	"fmt"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/unit"
)

// serviceItemPricer is a service object to price service items
type serviceItemPricer struct {
}

// NewServiceItemPricer constructs a pricer for service items
func NewServiceItemPricer() services.ServiceItemPricer {
	return &serviceItemPricer{}
}

// PriceServiceItem returns a price for any PaymentServiceItem
func (p serviceItemPricer) PriceServiceItem(appCtx appcontext.AppContext, item models.PaymentServiceItem) (unit.Cents, models.PaymentServiceItemParams, error) {
	pricer, err := p.getPricer(item.MTOServiceItem.ReService.Code)
	if err != nil {
		return unit.Cents(0), nil, err
	}

	// pricingParams are rate engine params that were queried from the pricing tables such as
	// price, rate, escalation etc.
	priceCents, pricingParams, err := pricer.PriceUsingParams(appCtx, item.PaymentServiceItemParams)
	if err != nil {
		return unit.Cents(0), nil, err
	}

	// createPricerGeneratedParams will throw an error if pricingParams is an empty slice
	// currently our pricers are returning empty slices for pricingParams
	// once all pricers have been updated to return pricingParams
	// TODO: this conditional logic should be removed
	var displayParams models.PaymentServiceItemParams
	if len(pricingParams) > 0 {
		displayParams, err = createPricerGeneratedParams(appCtx, item.ID, pricingParams)
	}
	return priceCents, displayParams, err
}

func (p serviceItemPricer) getPricer(serviceCode models.ReServiceCode) (services.ParamsPricer, error) {
	return PricerForServiceItem(serviceCode)
}

func PricerForServiceItem(serviceCode models.ReServiceCode) (services.ParamsPricer, error) {
	switch serviceCode {
	case models.ReServiceCodeMS:
		return NewManagementServicesPricer(), nil
	case models.ReServiceCodeCS:
		return NewCounselingServicesPricer(), nil
	case models.ReServiceCodeDLH:
		return NewDomesticLinehaulPricer(), nil
	case models.ReServiceCodeDSH:
		return NewDomesticShorthaulPricer(), nil
	case models.ReServiceCodeDOP:
		return NewDomesticOriginPricer(), nil
	case models.ReServiceCodeDDP:
		return NewDomesticDestinationPricer(), nil
	case models.ReServiceCodeDDSHUT:
		return NewDomesticDestinationShuttlingPricer(), nil
	case models.ReServiceCodeDOSHUT:
		return NewDomesticOriginShuttlingPricer(), nil
	case models.ReServiceCodeDCRT:
		return NewDomesticCratingPricer(), nil
	case models.ReServiceCodeDUCRT:
		return NewDomesticUncratingPricer(), nil
	case models.ReServiceCodeDPK:
		return NewDomesticPackPricer(), nil
	case models.ReServiceCodeDNPK:
		return NewDomesticNTSPackPricer(), nil
	case models.ReServiceCodeDUPK:
		return NewDomesticUnpackPricer(), nil
	case models.ReServiceCodeFSC:
		return NewFuelSurchargePricer(), nil
	case models.ReServiceCodeDOFSIT:
		return NewDomesticOriginFirstDaySITPricer(), nil
	case models.ReServiceCodeDDFSIT:
		return NewDomesticDestinationFirstDaySITPricer(), nil
	case models.ReServiceCodeDOASIT:
		return NewDomesticOriginAdditionalDaysSITPricer(), nil
	case models.ReServiceCodeDDASIT:
		return NewDomesticDestinationAdditionalDaysSITPricer(), nil
	case models.ReServiceCodeDOPSIT:
		return NewDomesticOriginSITPickupPricer(), nil
	case models.ReServiceCodeDDDSIT:
		return NewDomesticDestinationSITDeliveryPricer(), nil
	default:
		// TODO: We may want a different error type here after all pricers have been implemented
		return nil, apperror.NewNotImplementedError(fmt.Sprintf("pricer not found for code %s", serviceCode))
	}
}
