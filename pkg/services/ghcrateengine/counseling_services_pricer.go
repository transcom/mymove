package ghcrateengine

import (
	"fmt"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/unit"
)

type counselingServicesPricer struct {
}

// NewCounselingServicesPricer creates a new pricer for counseling services
func NewCounselingServicesPricer() services.CounselingServicesPricer {
	return &counselingServicesPricer{}
}

// Price determines the price for a counseling service
func (p managementServicesPricer) Price(appCtx appcontext.AppContext, serviceItem models.MTOServiceItem) (unit.Cents, services.PricingDisplayParams, error) {

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

// PriceUsingParams determines the price for a counseling service given PaymentServiceItemParams
func (p managementServicesPricer) PriceUsingParams(appCtx appcontext.AppContext, params models.PaymentServiceItemParams) (unit.Cents, services.PricingDisplayParams, error) {

	var serviceItem models.MTOServiceItem
	for _, param := range params {
		if param.PaymentServiceItem.MTOServiceItem.LockedPriceCents != nil {
			serviceItem = param.PaymentServiceItem.MTOServiceItem
			break
		}
	}

	if serviceItem.LockedPriceCents == nil {
		return unit.Cents(0), nil, fmt.Errorf("service item did not contain value for locked price cents")
	}

	return p.Price(appCtx, serviceItem)
}
