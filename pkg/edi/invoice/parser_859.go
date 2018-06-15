package invoice

import (
	"edi"
	"fmt"
	"time"
)

const dateFormat = "20060102"

var locationType = map[string]string{
	"77": "SERVICE_LOCATION",    // TODO: What is this?
	"7C": "PLACE_OF_OCCURRENCE", // TODO: What is this?
	"AB": "ADDITIONAL_PICKUP_ADDRESS",
	"AE": "ADDITIONAL_DELIVERY_ADDRESS",
	"DA": "DELIVERY_ADDRESS",
	"ND": "NEXT_DESTINATION",
	"OT": "ORIGIN_PORT",
	"PW": "PICKUP_ADDRESS",
	"SF": "SHIP_FROM",
	"ST": "SHIP_TO",
	"TR": "DESTINATION_PORT",
	"WD": "DESTINATION_STORAGE_FACILITY",
	"WO": "ORIGIN_STORAGE_FACILITY",
}

// Parser859 represents an EDI 859 parser
type Parser859 struct {
	invoice  Invoice
	segments []edi.Segment
}

// NewParser859 creates a new Parser859
func NewParser859(segments []edi.Segment) *Parser859 {
	return &Parser859{segments: segments}
}

// Invoice returns the parsed Invoice struct from the EDI 859
func (p *Parser859) Invoice() Invoice {
	return p.invoice
}

// Parse parses an array of EDI segments in X12 format into an Invoice
func (p *Parser859) Parse() error {
	for i := 0; i < len(p.segments); i++ {
		var err error
		segment := p.segments[i]

		switch segment.(type) {
		case *edi.B3:
			err = p.parseB3(*segment.(*edi.B3))
		case *edi.B3A:
			err = p.parseB3A(*segment.(*edi.B3A))
		case *edi.G62:
			err = p.parseG62(*segment.(*edi.G62))
		case *edi.LX:
			i, err = p.parseLXLoop(i)
		case *edi.N1:
			i, err = p.parseN1Loop(i, nil)
		case *edi.N9:
			err = p.parseN9(*segment.(*edi.N9), nil)
		}

		if err != nil {
			return err
		}
	}

	return nil
}

func (p *Parser859) parseB3(b3 edi.B3) error {
	p.invoice.InvoiceID = b3.InvoiceNumber
	p.invoice.Shipment.ShipmentID = b3.ShipmentIdentificationNumber
	if b3.ShipmentMethodOfPayment != "PP" {
		return fmt.Errorf("B3: Expected B304 to be PP (prepaid by seller), got: %s", b3.ShipmentMethodOfPayment)
	}

	var err error
	p.invoice.InvoiceDate, err = time.Parse(dateFormat, b3.Date) // TODO: timezone?
	if err != nil {
		return err
	}
	p.invoice.NetAmountDue = b3.NetAmountDue

	if b3.DateTimeQualifier == "035" {
		p.invoice.Shipment.DeliveryDate, err = time.Parse(dateFormat, b3.DeliveryDate) // TODO: timezone?
		if err != nil {
			return err
		}
	} else if b3.DateTimeQualifier != "" {
		return fmt.Errorf("B3: Expected B310 to be 035 (delivered), got: %s", b3.DateTimeQualifier)
	}

	p.invoice.Shipment.SCAC = b3.StandardCarrierAlphaCode
	return nil
}

func (p *Parser859) parseB3A(b3a edi.B3A) error {
	if b3a.TransactionTypeCode != "DI" {
		return fmt.Errorf("B3A: Expected B3A01 to be DI (debit invoice), got: %s", b3a.TransactionTypeCode)
	}
	return nil
}

func (p *Parser859) parseG62(g62 edi.G62) error {
	date, err := time.Parse(dateFormat, g62.Date)
	if err != nil {
		return err
	}

	switch g62.DateQualifier {
	case "86": // Actual pickup date
		p.invoice.Shipment.PickupDate = date
	case "35": // Actual delivery date
		// TODO: This can already be set by the B3 segment. Should we verify that they match?
		p.invoice.Shipment.DeliveryDate = date
	default:
		return fmt.Errorf("G62: Expected G6201 to be 86 or 35, got: %s", g62.DateQualifier)
	}

	return nil
}

// Parses the loop starting with the LX segment. The int argment is the index of
// the beginning LX segment. The int in the return value is the index of the
// last segment in the LX loop.
func (p *Parser859) parseLXLoop(i int) (int, error) {
	lineItem := LineItem{}
	var err error

F:
	for i++; i < len(p.segments); i++ {
		segment := p.segments[i]
		switch segment.(type) {
		case *edi.L0:
			i, err = p.parseL0Loop(i, &lineItem)
		case *edi.L5:
			lineItem.Descriptions = append(lineItem.Descriptions, segment.(*edi.L5).LadingDescription)
		case *edi.L7:
			lineItem.ItemCode = segment.(*edi.L7).TariffItemNumber
		case *edi.N1:
			i, err = p.parseN1Loop(i, &lineItem)
		case *edi.N9:
			err = p.parseN9(*segment.(*edi.N9), &lineItem)
		default:
			break F
		}

		if err != nil {
			return i, err
		}
	}

	p.invoice.LineItems = append(p.invoice.LineItems, lineItem)
	return nil, i - 1
}

func (p *Parser859) parseL0Loop(i int, lineItem *LineItem) (int, error) {
	charge := LineItemCharge{}
	segment := p.segments[i] // L0 segment
	err := p.parseL0(*segment.(*edi.L0), &charge)
	if err != nil {
		return i, err
	}

F:
	for i++; i < len(p.segments); i++ {
		segment = p.segments[i]
		switch segment.(type) {
		case *edi.L1:
			err = p.parseL1(*segment.(*edi.L1), &charge)
		case *edi.MEA:
			err = p.parseMEA(*segment.(*edi.MEA), &charge)
		default:
			break F
		}

		if err != nil {
			return i, err
		}
	}

	lineItem.Charge = charge
	return i - 1, nil
}

func (p *Parser859) parseL0(l0 edi.L0, charge *LineItemCharge) error {
	switch l0.BilledRatedAsQualifier {
	case "CF": // Cubic foot
		charge.VolumeCubicFeet = l0.BilledRatedAsQuantity
	case "EA": // Each
		charge.Quantity = l0.BilledRatedAsQuantity
	case "FR": // Flat rate
	case "MV": // Monetary value
	case "NR": // Container
	case "OR": // Other (weight)
	case "RV": // Release value
	case "TD": // Days
		charge.Days = l0.BilledRatedAsQuantity
	case "TH": // Hours
		charge.Hours = l0.BilledRatedAsQuantity
	default:
		return fmt.Errorf("L0: Unexpected L003: %s", l0.BilledRatedAsQualifier)
	}

	if l0.WeightQualifier == "B" { // Billed weight
		if l0.WeightUnitCode != "L" { // Pounds
			return fmt.Errorf("L0: Expected L011 to be L, got %s", l0.WeightUnitCode)
		}
		charge.BilledWeightPounds = l0.Weight * 100 // L004 is in CWT
	}
	return nil
}

func (p *Parser859) parseL1(l1 edi.L1, charge *LineItemCharge) error {
	if l1.RateValueQualifier != "RC" {
		return fmt.Errorf("L1: expected L103 to be RC, got %s", l1.RateValueQualifier)
	}

	charge.Rate = l1.FreightRate
	charge.TotalCharge = l1.Charge
	charge.RateMultiplier = l1.Percent
	return nil
}

func (p *Parser859) parseMEA(mea edi.MEA, charge *LineItemCharge) error {
	switch mea.MeasurementReferenceIDCode {
	case "BC":
		if mea.MeasurementQualifier == "B" {
			charge.BilledWeightPounds = mea.MeasurementValue
			return nil
		}
	case "EF":
		if mea.MeasurementQualifier == "DS" {
			// TODO: What is this used for?? What are the units?
			charge.Distance = mea.MeasurementValue
			return nil
		} else if mea.MeasurementQualifier == "VOL" {
			// TODO: What is this used for?? What are the units?
			charge.VolumeCubicFeet = mea.MeasurementValue
			return nil
		}
	case "WT":
		if mea.MeasurementQualifier == "N" || mea.MeasurementQualifier == "G" {
			charge.WeightPounds = mea.MeasurementValue
			return nil
		}
	}
	return fmt.Errorf("MEA: Unparseable segment: %s", mea.String("*"))
}

func (p *Parser859) parseN1Loop(i int, lineItem *LineItem) (int, error) {
	n1 := *p.segments[i].(*edi.N1)
	switch n1.EntityIdentifierCode {
	case "RG":
		if n1.IdentificationCodeQualifier == "27" { // 27 = GBLOC
			p.invoice.Shipment.OriginGBLOC = n1.IdentificationCode
		} else if n1.IdentificationCodeQualifier != "" {
			return i, fmt.Errorf("N1: Expected N103 to be 27, got: %s", n1.IdentificationCodeQualifier)
		}
	case "RH":
		if n1.IdentificationCodeQualifier == "27" { // 27 = GBLOC
			p.invoice.Shipment.DestinationGBLOC = n1.IdentificationCode
		} else if n1.IdentificationCodeQualifier != "" {
			return i, fmt.Errorf("N1: Expected N103 to be 27, got: %s", n1.IdentificationCodeQualifier)
		}
	default:
		t := locationType[n1.EntityIdentifierCode]
		if t == "" {
			return i, fmt.Errorf("N1: Unexpected N101: %s", n1.EntityIdentifierCode)
		}
		i++
		location := Location{Name: n1.Name, LocationType: t}
		n4, ok := p.segments[i].(*edi.N4)
		if !ok {
			return i, fmt.Errorf("Missing expected N4 in N1 loop")
		}
		err := p.parseN4(*n4, &location)
		if err != nil {
			return i, err
		}
		if lineItem != nil {
			lineItem.Locations = append(lineItem.Locations, location)
		}
	}

	return i, nil
}

func (p *Parser859) parseN4(n4 edi.N4, location *Location) error {
	location.City = n4.CityName
	location.State = n4.StateOrProvinceCode
	location.PostalCode = n4.PostalCode
	location.CountryCode = n4.CountryCode
	switch n4.LocationQualifier {
	case "CC":
		location.County = n4.LocationIdentifier
	case "RA":
		location.RateAreaCode = n4.LocationIdentifier
	default:
		return fmt.Errorf("N4: Expected N405 to be CC or RA, got: %s", n4.LocationQualifier)
	}
	return nil
}

func (p *Parser859) parseN9(n9 edi.N9, lineItem *LineItem) error {
	switch n9.ReferenceIdentificationQualifier {
	case "0L":
		if lineItem == nil {
			return fmt.Errorf("Erroneous N9:0L segment")
		}
		lineItem.Note += n9.FreeFormDescription
	case "2I":
		if lineItem == nil || lineItem.SyncadaMatchNumber != "" {
			return fmt.Errorf("Erroneous N9:2I segment")
		}
		lineItem.SyncadaMatchNumber = n9.ReferenceIdentification
	case "8M": // Originating Company Identifier
		// TODO: What is this used for? I only see "DAYCOSXRSBIL"
	case "CN": // Carrier's Reference Number (PRO/Invoice)
		// TODO: What is this used for?
	case "KK": // Delivery Reference
		if n9.ReferenceIdentification == "DV" {
			p.invoice.Shipment.Status = "DELIVERED"
		} else if n9.ReferenceIdentification == "RC" {
			p.invoice.Shipment.Status = "RECONSIGNED"
		} else if n9.ReferenceIdentification == "ST" {
			p.invoice.Shipment.Status = "IN_SIT"
		}
	}
	return nil
}
