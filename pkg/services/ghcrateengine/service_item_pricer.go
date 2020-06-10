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
	pricer, err := p.getPricer(item)
	if err != nil {
		return unit.Cents(0), err
	}

	return pricer.Price()
}

func (p serviceItemPricer) getPricer(item models.PaymentServiceItem) (Pricer, error) {
	serviceCode := item.MTOServiceItem.ReService.Code

	switch serviceCode {
	case models.ReServiceCodeMS, models.ReServiceCodeCS:
		return NewTaskOrderServicesPricerFromParams(p.db, serviceCode, item.PaymentServiceItemParams)
	default:
		return nil, services.NewNotImplementedError(fmt.Sprintf("pricer not found for code %s", serviceCode))
	}
}
