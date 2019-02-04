package rateengine

import (
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/unit"
)

var tariff400ngItemPricing = map[string]pricer{
	// Re-weigh
	"4A": newBasicQuantityPricer(),
	"4B": newBasicQuantityPricer(),

	// Below are SIT-related codes, uncomment after SIT is implemented
	// Attempted delivery from SIT
	// "17A": newFlatRatePricer(),
	// "17B": newFlatRatePricer(),
	// Priced using base linehaul rate, needs more work
	// "17C": newFlatRatePricer(),
	// "17D": newMinimumQuantityHundredweightPricer(1000),
	// "17E": newFlatRatePricer(),
	// "17F": newFlatRatePricer(),
	// Priced using base linehaul rate, needs more work
	// "17G": newFlatRatePricer(),

	// Extra pickups, diversions
	"28A": newBasicQuantityPricer(),
	"28B": newBasicQuantityPricer(),
	"28C": newBasicQuantityPricer(),

	// Third Party Service (TPS) charge
	"35A": newBasicQuantityPricer(),
	// Key West service charge
	"35B": newFlatRatePricer(),

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

	// Shuttle service
	"125A": newFlatRatePricer(),
	"125B": newFlatRatePricer(),
	"125C": newFlatRatePricer(),
	"125D": newFlatRatePricer(),

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

	// Overtime loading/unloading
	// Note: this pricer doesn't allow for weights under 1,000, which the below excerpt would imply is possible
	// "If only a portion of a shipment is loaded/unloaded a separate weight ticket MUST be provided,
	// otherwise TSP is limited to billing 1,000 lbs."
	"175A": newMinimumQuantityPricer(1000),

	// Below are SIT-related codes, uncomment after SIT is implemented
	// SIT
	// "185A": newMinimumQuantityHundredweightPricer(1000),
	// Priced using two quantities (weight and days), needs more work
	// "185B": newMinimumQuantityHundredweightPricer(1000),

	// Below are SIT-related codes, uncomment after SIT is implemented
	// SIT P/D OT
	// "210D": newFlatRatePricer(),
	// "210E": newFlatRatePricer(),

	// Pickup/delivery at third-party and self-storage warehouses
	"225A": newFlatRatePricer(),
	"225B": newFlatRatePricer(),

	// Misc. charge
	"226A": newBasicQuantityPricer(),
}

// Some codes (17, mainly) are explicitly priced using rates corresponding to a different item code
var tariff400ngItemRateMap = map[string]string{
	"17A": "210A",
	"17D": "185A",
	"17E": "210D",
}

// These codes have charges based on weight, which will use the final measured shipment weight
var tariff400ngWeightBasedItems = map[string]bool{
	"17D":  true,
	"175A": true,
	"185A": true,
}

// ComputeShipmentLineItemCharge calculates the total charge for a supplied shipment line item and returns it and the DISCOUNTED rate
func (re *RateEngine) ComputeShipmentLineItemCharge(shipmentLineItem models.ShipmentLineItem) (FeeAndRate, error) {
	itemCode := shipmentLineItem.Tariff400ngItem.Code
	shipment := shipmentLineItem.Shipment

	if shipment.NetWeight == nil {
		return FeeAndRate{}, errors.New("Can't price a shipment line item for a shipment without NetWeight")
	}

	// Defaults to origin postal code, but if location is NEITHER than this doesn't matter
	zip := Zip5ToZip3(shipment.PickupAddress.PostalCode)
	if shipmentLineItem.Location == models.ShipmentLineItemLocationDESTINATION {
		zip = Zip5ToZip3(shipment.Move.Orders.NewDutyStation.Address.PostalCode)
	}
	shipDate := shipment.BookDate

	serviceArea, err := models.FetchTariff400ngServiceAreaForZip3(re.db, zip, *shipDate)
	if err != nil {
		return FeeAndRate{}, errors.Wrap(err, "Fetching 400ng service area from db")
	}

	var rateCents unit.Cents
	if itemCode == "185A" {
		// Rates for SIT are stored  on the service area
		rateCents = serviceArea.SIT185ARateCents
	} else if itemCode == "185B" {
		rateCents = serviceArea.SIT185BRateCents
	} else if itemCode == "226A" {
		// 226A is Misc charge, allow user to enter dollar amount as quantity
		rateCents = unit.Cents(100)
	} else if itemCode == "35A" {
		// 35A is a Third Party Service (TPS) charge, allow user to enter dollar amount as quantity
		rateCents = unit.Cents(100)
	} else {
		// Most rates should be in the tariff400ngItemRates table though

		// If code is priced using rate from separate code, use that
		effectiveItemCode := itemCode
		if mappedCode, ok := tariff400ngItemRateMap[effectiveItemCode]; ok {
			effectiveItemCode = mappedCode
		}

		rate, err := models.FetchTariff400ngItemRate(re.db,
			effectiveItemCode,
			serviceArea.ServicesSchedule,
			*shipment.NetWeight,
			*shipDate,
		)
		if err != nil {
			return FeeAndRate{}, errors.Wrap(err, "Fetching 400ng item rate from db")
		}
		rateCents = rate.RateCents
	}

	// Make sure we have a ShipmentOffer and TSPP if we need to apply a discount
	hasTSPP := len(shipment.ShipmentOffers) == 0 || shipment.ShipmentOffers[0].TransportationServiceProviderPerformance.ID == uuid.Nil
	if shipmentLineItem.Tariff400ngItem.DiscountType != models.Tariff400ngItemDiscountTypeNONE && hasTSPP {
		return FeeAndRate{}, errors.New("No TSPP provided for Shipment, something is very wrong")
	}

	var discountRate *unit.DiscountRate
	if shipmentLineItem.Tariff400ngItem.DiscountType == models.Tariff400ngItemDiscountTypeHHG || shipmentLineItem.Tariff400ngItem.DiscountType == models.Tariff400ngItemDiscountTypeHHGLINEHAUL50 {
		discountRate = &shipment.ShipmentOffers[0].TransportationServiceProviderPerformance.LinehaulRate
	} else if shipmentLineItem.Tariff400ngItem.DiscountType == models.Tariff400ngItemDiscountTypeSIT {
		discountRate = &shipment.ShipmentOffers[0].TransportationServiceProviderPerformance.SITRate
	}

	// Weight-based items will pull final weight values from the shipment when available
	appliedQuantity := shipmentLineItem.Quantity1
	if _, ok := tariff400ngWeightBasedItems[itemCode]; ok {
		if shipment.NetWeight == nil {
			return FeeAndRate{}, errors.New("Can't price a weight-based accessorial without shipment net weight")
		}
		appliedQuantity = unit.BaseQuantityFromInt(shipment.NetWeight.Int())
	}

	appliedRate := rateCents
	if discountRate != nil {
		appliedRate = discountRate.Apply(rateCents)
	}

	if itemPricer, ok := tariff400ngItemPricing[itemCode]; ok {
		return FeeAndRate{Fee: itemPricer.price(rateCents, appliedQuantity, discountRate), Rate: appliedRate.ToMillicents()}, nil
	}

	return FeeAndRate{}, errors.New("Could not find pricing function for given code")
}

// PricePreapprovalRequestsForShipment for a shipment, computes prices for all approved pre-approval requests and populates amount_cents field and applied_rate on those models
func (re *RateEngine) PricePreapprovalRequestsForShipment(shipment models.Shipment) ([]models.ShipmentLineItem, error) {
	items, err := models.FetchApprovedPreapprovalRequestsByShipment(re.db, shipment)
	if err != nil {
		return []models.ShipmentLineItem{}, err
	}

	for i := 0; i < len(items); i++ {
		err := re.PricePreapprovalRequest(&items[i])
		if err != nil {
			return []models.ShipmentLineItem{}, err
		}
	}

	return items, nil
}

// PricePreapprovalRequest computes price for given pre-approval requests and populates amount_cents field and applied_rate on those models
func (re *RateEngine) PricePreapprovalRequest(shipmentLineItem *models.ShipmentLineItem) error {

	feeAndRate, err := re.ComputeShipmentLineItemCharge(*shipmentLineItem)
	if err != nil {
		return err
	}
	shipmentLineItem.AmountCents = &feeAndRate.Fee
	shipmentLineItem.AppliedRate = &feeAndRate.Rate

	return nil
}
