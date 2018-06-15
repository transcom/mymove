package invoice

import (
	"time"
)

// Invoice represents an invoice from a TSP for a shipment
type Invoice struct {
	InvoiceDate  time.Time
	InvoiceID    string
	Shipment     Shipment
	LineItems    []InvoiceLineItem
	NetAmountDue float64
}

// Shipment includes information about the shipment from the Invoice.
// This info should either be used to identify the shipment to associate this
// invoice to or to update data in the shipment.
type Shipment struct {
	DeliveryDate     time.Time
	DestinationGBLOC string
	OriginGBLOC      string
	PickupDate       time.Time
	SCAC             string
	ShipmentID       string
	Status           string // TODO: enum
}

// LineItem is a line item in the invoice
type LineItem struct {
	Descriptions       []string
	Charge             LineItemCharge
	ItemCode           string
	Locations          []InvoiceLocation
	Note               string
	SITControlNumber   string
	SyncadaMatchNumber string
}

// Location represents a physical location in the invoice
type Location struct {
	City         string
	CountryCode  string
	County       string
	Name         string
	PostalCode   string
	RateAreaCode string
	State        string
	LocationType string
}

// LineItemCharge includes the charge for a particular line item in an invoice,
// as well as the values used to calculate that charge
type LineItemCharge struct {
	BilledWeightPounds float64
	Days               float64
	Distance           float64
	Hours              float64
	Quantity           float64
	Rate               float64
	RateMultiplier     float64
	TotalCharge        float64
	VolumeCubicFeet    float64
	WeightPounds       float64
}
