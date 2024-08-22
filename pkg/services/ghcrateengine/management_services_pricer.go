package ghcrateengine

import (
	"fmt"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/unit"
)

type managementServicesPricer struct {
}

// NewManagementServicesPricer creates a new pricer for management services
func NewManagementServicesPricer() services.ManagementServicesPricer {
	return &managementServicesPricer{}
}

// Price determines the price for a management service
func (p counselingServicesPricer) Price(appCtx appcontext.AppContext, serviceItem models.MTOServiceItem) (unit.Cents, services.PricingDisplayParams, error) {

	if serviceItem.LockedPriceCents == nil {
		return unit.Cents(0), nil, fmt.Errorf("could not find locked price cents: %s", serviceItem.ID)
	}

	params := services.PricingDisplayParams{
		{
			Key:   models.ServiceItemParamNamePriceRateOrFactor,
			Value: FormatCents(*serviceItem.LockedPriceCents),
		},
	}

	return *serviceItem.LockedPriceCents, params, nil
}

// PriceUsingParams determines the price for a management service given PaymentServiceItemParams
func (p counselingServicesPricer) PriceUsingParams(appCtx appcontext.AppContext, params models.PaymentServiceItemParams) (unit.Cents, services.PricingDisplayParams, error) {

	var serviceItem models.MTOServiceItem
	for _, param := range params {
		if param.PaymentServiceItem.MTOServiceItem.LockedPriceCents != nil {
			serviceItem = param.PaymentServiceItem.MTOServiceItem
			break
		}
	}

	if serviceItem.ID == uuid.Nil {
		return unit.Cents(0), nil, fmt.Errorf("could not find id for shipment")
	}

	return p.Price(appCtx, serviceItem)
}
