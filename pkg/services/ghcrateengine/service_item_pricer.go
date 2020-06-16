package ghcrateengine

import (
	"fmt"

	"github.com/gobuffalo/pop"

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
func (p serviceItemPricer) PriceServiceItem(item models.PaymentServiceItem) (unit.Cents, error) {
	pricer, err := p.getPricer(item.MTOServiceItem.ReService.Code)
	if err != nil {
		return unit.Cents(0), err
	}

	return pricer.PriceUsingParams(item.PaymentServiceItemParams)
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
	default:
		// TODO: We may want a different error type here after all pricers have been implemented
		return nil, services.NewNotImplementedError(fmt.Sprintf("pricer not found for code %s", serviceCode))
	}
}
