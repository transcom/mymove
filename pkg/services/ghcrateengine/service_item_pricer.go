package ghcrateengine

import (
	"fmt"

	"github.com/gobuffalo/pop/v5"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/unit"
)

// serviceItemPricer is a service object to price service items
type serviceItemPricer struct {
	db *pop.Connection
}

// NewServiceItemPricer constructs a pricer for service items
func NewServiceItemPricer(db *pop.Connection) services.ServiceItemPricer {
	return &serviceItemPricer{
		db: db,
	}
}

// PriceServiceItem returns a price for any PaymentServiceItem
func (p serviceItemPricer) PriceServiceItem(item models.PaymentServiceItem) (unit.Cents, models.PaymentServiceItemParams, error) {
	pricer, err := p.getPricer(item.MTOServiceItem.ReService.Code)
	if err != nil {
		return unit.Cents(0), nil, err
	}

	priceCents, pricingParams, err := pricer.PriceUsingParams(item.PaymentServiceItemParams)
	if err != nil {
		return unit.Cents(0), nil, err
	}

	// createPricerGeneratedParams will throw an error if pricingParams is an empty slice
	// currently our pricers are returning empty slices for pricingParams
	// once all pricers have been updated to return pricingParams
	// TODO: this conditional logic should be removed
	var displayParams models.PaymentServiceItemParams
	if len(pricingParams) > 0 {
		displayParams, err = createPricerGeneratedParams(p.db, item.ID, pricingParams)
	}
	return priceCents, displayParams, err
}

func (p serviceItemPricer) UsingConnection(db *pop.Connection) services.ServiceItemPricer {
	p.db = db
	return p
}

func (p serviceItemPricer) getPricer(serviceCode models.ReServiceCode) (services.ParamsPricer, error) {
	switch serviceCode {
	case models.ReServiceCodeMS:
		return NewManagementServicesPricer(p.db), nil
	case models.ReServiceCodeCS:
		return NewCounselingServicesPricer(p.db), nil
	case models.ReServiceCodeDLH:
		return NewDomesticLinehaulPricer(p.db), nil
	case models.ReServiceCodeDSH:
		return NewDomesticShorthaulPricer(p.db), nil
	case models.ReServiceCodeDOP:
		return NewDomesticOriginPricer(p.db), nil
	case models.ReServiceCodeDDP:
		return NewDomesticDestinationPricer(p.db), nil
	case models.ReServiceCodeDPK:
		return NewDomesticPackPricer(p.db), nil
	case models.ReServiceCodeDUPK:
		return NewDomesticUnpackPricer(p.db), nil
	case models.ReServiceCodeFSC:
		return NewFuelSurchargePricer(p.db), nil
	case models.ReServiceCodeDOFSIT:
		return NewDomesticOriginFirstDaySITPricer(p.db), nil
	case models.ReServiceCodeDDFSIT:
		return NewDomesticDestinationFirstDaySITPricer(p.db), nil
	case models.ReServiceCodeDOASIT:
		return NewDomesticOriginAdditionalDaysSITPricer(p.db), nil
	case models.ReServiceCodeDDASIT:
		return NewDomesticDestinationAdditionalDaysSITPricer(p.db), nil
	case models.ReServiceCodeDOPSIT:
		return NewDomesticOriginSITPickupPricer(p.db), nil
	case models.ReServiceCodeDDDSIT:
		return NewDomesticDestinationSITDeliveryPricer(p.db), nil
	default:
		// TODO: We may want a different error type here after all pricers have been implemented
		return nil, services.NewNotImplementedError(fmt.Sprintf("pricer not found for code %s", serviceCode))
	}
}
