package rateengine

import (
	"github.com/pkg/errors"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/unit"
)

var tariff400ngItemPricing = map[string]pricer{
	// Re-weigh
	"4A": newBasicQuantityPricer(),
	"4B": newBasicQuantityPricer(),

	// Attempted delivery - 31-50 miles
	"17B": newFlatRatePricer(),
	"17F": newFlatRatePricer(),

	// Extra pickups, diversions
	"28A": newBasicQuantityPricer(),
	"28B": newBasicQuantityPricer(),
	"28C": newBasicQuantityPricer(),

	// Key west service charge
	"35A": newFlatRatePricer(),

	// Crating - subject to a minimum charge of 4 cu. ft.
	"105B": newMinimumQuantityPricer(4),
	"105E": newMinimumQuantityPricer(4),

	// Debris removal
	"105D": newBasicQuantityPricer(),

	// Extra labor, waiting time
	"120A": newBasicQuantityPricer(),
	"120B": newBasicQuantityPricer(),
	"120C": newBasicQuantityPricer(),
	"120D": newBasicQuantityPricer(),
	"120E": newBasicQuantityPricer(),
	"120F": newBasicQuantityPricer(),

	// Bulky items
	"130A": newBasicQuantityPricer(),
	"130B": newBasicQuantityPricer(),
	"130C": newBasicQuantityPricer(),
	"130D": newBasicQuantityPricer(),
	"130E": newBasicQuantityPricer(),
	"130F": newBasicQuantityPricer(),
	"130G": newBasicQuantityPricer(),
	"130H": newBasicQuantityPricer(),
	"130I": newBasicQuantityPricer(),
	"130J": newBasicQuantityPricer(),

	// Pickup/delivery at third-party and self-storage warehouses
	"225A": newMinimumQuantityPricer(1000),
	"225B": newScaledRateMinimumQuantityPricer(10, 1000),

	// Misc. charge
	"226A": newBasicQuantityPricer(),
}

// ComputeShipmentLineItemCharge calculates the total charge for a supplied shipment line item
func (re *RateEngine) ComputeShipmentLineItemCharge(shipmentLineItem models.ShipmentLineItem, shipment models.Shipment) (unit.Cents, error) {
	zip := Zip5ToZip3(shipment.PickupAddress.PostalCode)
	if shipmentLineItem.Location == models.ShipmentLineItemLocationDESTINATION {
		zip = Zip5ToZip3(shipment.Move.Orders.NewDutyStation.Address.PostalCode)
	}
	shipDate := shipment.BookDate

	serviceArea, err := models.FetchTariff400ngServiceAreaForZip3(re.db, zip, *shipDate)
	if err != nil {
		return unit.Cents(0), errors.Wrap(err, "Fetching 400ng service area from db")
	}

	rate, err := models.FetchTariff400ngItemRate(re.db,
		shipmentLineItem.Tariff400ngItem.Code,
		serviceArea.ServicesSchedule,
		*shipment.NetWeight,
		*shipDate,
	)
	if err != nil {
		return unit.Cents(0), errors.Wrap(err, "Fetching 400ng item rate from db")
	}

	var discountRate *unit.DiscountRate
	if shipmentLineItem.Tariff400ngItem.DiscountType == models.Tariff400ngItemDiscountTypeHHG {
		discountRate = &shipment.ShipmentOffers[0].TransportationServiceProviderPerformance.LinehaulRate
	} else if shipmentLineItem.Tariff400ngItem.DiscountType == models.Tariff400ngItemDiscountTypeSIT {
		discountRate = &shipment.ShipmentOffers[0].TransportationServiceProviderPerformance.SITRate
	}

	if itemPricer, ok := tariff400ngItemPricing[shipmentLineItem.Tariff400ngItem.Code]; ok {
		return itemPricer.price(rate.RateCents, shipmentLineItem.Quantity1, discountRate), nil
	}

	return unit.Cents(0), errors.New("Could not find pricing function for given code")
}
