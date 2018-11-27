package ediinvoice

import (
	"bytes"
	"fmt"
	"time"

	"github.com/facebookgo/clock"
	"github.com/gobuffalo/pop"
	"github.com/pkg/errors"

	"github.com/transcom/mymove/pkg/db/sequence"
	"github.com/transcom/mymove/pkg/edi"
	"github.com/transcom/mymove/pkg/edi/segment"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/rateengine"
)

const dateFormat = "20060102"
const timeFormat = "1504"
const senderCode = "MYMOVE"

//const senderCode = "W28GPR-DPS"   // TODO: update with ours when US Bank gets it to us
const receiverCode = "8004171844" // Syncada

// ICNSequenceName used to query Interchange Control Numbers from DB
const ICNSequenceName = "interchange_control_number"

// Invoice858C holds all the segments that are generated
type Invoice858C struct {
	ISA       edisegment.ISA
	GS        edisegment.GS
	Shipments [][]edisegment.Segment
	GE        edisegment.GE
	IEA       edisegment.IEA
}

// Segments returns the invoice as an array of rows (string arrays),
// each containing a segment, to prepare it for writing
func (invoice Invoice858C) Segments() [][]string {
	records := [][]string{
		invoice.ISA.StringArray(),
		invoice.GS.StringArray(),
	}
	for _, shipment := range invoice.Shipments {
		for _, line := range shipment {
			records = append(records, line.StringArray())
		}
	}
	records = append(records, invoice.GE.StringArray())
	records = append(records, invoice.IEA.StringArray())
	return records
}

// EDIString returns the EDI representation of an 858C
func (invoice Invoice858C) EDIString() (string, error) {
	var b bytes.Buffer
	ediWriter := edi.NewWriter(&b)
	err := ediWriter.WriteAll(invoice.Segments())
	if err != nil {
		return "", err
	}
	return b.String(), err
}

// Generate858C generates an EDI X12 858C transaction set
func Generate858C(shipmentsAndCosts []rateengine.CostByShipment, db *pop.Connection, sendProductionInvoice bool, clock clock.Clock) (Invoice858C, error) {
	loc, err := time.LoadLocation("America/Los_Angeles")
	if err != nil {
		return Invoice858C{}, err
	}
	currentTime := clock.Now().In(loc)

	interchangeControlNumber, err := sequence.NextVal(db, ICNSequenceName)
	if err != nil {
		return Invoice858C{}, errors.Wrap(err, fmt.Sprintf("Failed to get next Interchange Control Number"))
	}

	var usageIndicator string
	if sendProductionInvoice {
		usageIndicator = "P"
	} else {
		usageIndicator = "T"
	}

	invoice := Invoice858C{}
	invoice.ISA = edisegment.ISA{
		AuthorizationInformationQualifier: "00", // No authorization information
		AuthorizationInformation:          fmt.Sprintf("%010d", 0),
		SecurityInformationQualifier:      "00", // No security information
		SecurityInformation:               fmt.Sprintf("%010d", 0),
		InterchangeSenderIDQualifier:      "ZZ",
		InterchangeSenderID:               fmt.Sprintf("%-15v", senderCode), // Must be 15 characters
		InterchangeReceiverIDQualifier:    "12",
		InterchangeReceiverID:             fmt.Sprintf("%-15s", receiverCode), // Must be 15 characters
		InterchangeDate:                   currentTime.Format("060102"),
		InterchangeTime:                   currentTime.Format(timeFormat),
		InterchangeControlStandards:       "U",
		InterchangeControlVersionNumber:   "00401",
		InterchangeControlNumber:          interchangeControlNumber,
		AcknowledgementRequested:          1,
		UsageIndicator:                    usageIndicator, // T for test, P for production
		ComponentElementSeparator:         "|",
	}
	invoice.GS = edisegment.GS{
		FunctionalIdentifierCode: "SI", // Shipment Information (858)
		ApplicationSendersCode:   senderCode,
		ApplicationReceiversCode: receiverCode,
		Date:                  currentTime.Format(dateFormat),
		Time:                  currentTime.Format(timeFormat),
		GroupControlNumber:    interchangeControlNumber,
		ResponsibleAgencyCode: "X", // Accredited Standards Committee X12
		Version:               "004010",
	}

	var shipments []models.Shipment

	invoice.Shipments = make([][]edisegment.Segment, 0)
	for index, shipmentWithCost := range shipmentsAndCosts {
		shipment := shipmentWithCost.Shipment

		shipmentSegments, err := generate858CShipment(shipmentWithCost, index+1)
		if err != nil {
			return invoice, err
		}
		invoice.Shipments = append(invoice.Shipments, shipmentSegments)
		shipments = append(shipments, shipment)
	}

	invoice.GE = edisegment.GE{
		NumberOfTransactionSetsIncluded: len(shipments),
		GroupControlNumber:              interchangeControlNumber,
	}
	invoice.IEA = edisegment.IEA{
		NumberOfIncludedFunctionalGroups: 1,
		InterchangeControlNumber:         interchangeControlNumber,
	}

	return invoice, nil
}

func generate858CShipment(shipmentWithCost rateengine.CostByShipment, sequenceNum int) ([]edisegment.Segment, error) {
	transactionNumber := fmt.Sprintf("%04d", sequenceNum)
	segments := []edisegment.Segment{
		&edisegment.ST{
			TransactionSetIdentifierCode: "858",
			TransactionSetControlNumber:  transactionNumber,
		},
	}

	headingSegments, err := getHeadingSegments(shipmentWithCost, sequenceNum)
	if err != nil {
		return segments, err
	}
	segments = append(segments, headingSegments...)

	lineItemSegments, err := getLineItemSegments(shipmentWithCost)
	if err != nil {
		return segments, err
	}
	segments = append(segments, lineItemSegments...)

	segments = append(segments, &edisegment.SE{
		NumberOfIncludedSegments:    len(segments) + 1, // Include SE in count
		TransactionSetControlNumber: transactionNumber,
	})

	return segments, nil
}

func getHeadingSegments(shipmentWithCost rateengine.CostByShipment, sequenceNum int) ([]edisegment.Segment, error) {
	shipment := shipmentWithCost.Shipment
	segments := []edisegment.Segment{}
	/* for bx
	if shipment.TransportationServiceProviderID == nil {
		return "", errors.New("Shipment is missing TSP ID")
	}
	var tsp models.TransportationServiceProvider
	err := db.Find(&tsp, shipment.TransportationServiceProviderID)
	if err != nil {
		return "", err
	}
	*/

	name := ""
	if shipment.ServiceMember.LastName != nil {
		name = *shipment.ServiceMember.LastName
	}
	if shipment.PickupAddress == nil {
		return segments, errors.New("Shipment is missing pick up address")
	}
	street2 := ""
	if shipment.PickupAddress.StreetAddress2 != nil {
		street2 = *shipment.PickupAddress.StreetAddress2
	}
	country := "US"
	if shipment.PickupAddress.Country != nil {
		country = *shipment.PickupAddress.Country
	}

	orders := shipment.Move.Orders
	ordersNumber := orders.OrdersNumber
	if ordersNumber == nil {
		return segments, errors.New("Orders is missing orders number")
	}
	tac := orders.TAC
	if tac == nil {
		return segments, errors.New("Orders is missing TAC")
	}
	affiliation := shipment.ServiceMember.Affiliation
	if shipment.ServiceMember.Affiliation == nil {
		return segments, errors.New("Service member is missing affiliation")
	}
	GBL := shipment.GBLNumber
	if GBL == nil {
		return segments, errors.New("GBL Number is missing for Shipment Identification Number (BX04)")
	}

	return []edisegment.Segment{
		&edisegment.BX{
			TransactionSetPurposeCode:    "00", // Original
			TransactionMethodTypeCode:    "J",  // Motor
			ShipmentMethodOfPayment:      "PP", // Prepaid by seller
			ShipmentIdentificationNumber: *GBL,
			StandardCarrierAlphaCode:     "MCCG", // TODO: real SCAC
			ShipmentQualifier:            "4",    // HHG Government Bill of Lading
		},
		&edisegment.N9{
			ReferenceIdentificationQualifier: "DY", // DoD transportation service code #
			ReferenceIdentification:          "SC", // Shipment & cost information
		},
		&edisegment.N9{
			ReferenceIdentificationQualifier: "CN",          // Invoice number
			ReferenceIdentification:          "ABCD00001-1", // TODO: real invoice number
		},
		&edisegment.N9{
			ReferenceIdentificationQualifier: "PQ",       // Payee code
			ReferenceIdentification:          "ABBV2708", // TODO: add real supplier ID
		},
		&edisegment.N9{
			ReferenceIdentificationQualifier: "OQ", // Order number
			ReferenceIdentification:          *ordersNumber,
			FreeFormDescription:              string(*affiliation),
			Date:                             orders.IssueDate.Format(dateFormat),
		},
		// Ship from address
		&edisegment.N1{
			EntityIdentifierCode: "SF", // Ship From
			Name:                 name,
		},
		&edisegment.N3{
			AddressInformation1: shipment.PickupAddress.StreetAddress1,
			AddressInformation2: street2,
		},
		&edisegment.N4{
			CityName:            shipment.PickupAddress.City,
			StateOrProvinceCode: shipment.PickupAddress.State,
			PostalCode:          shipment.PickupAddress.PostalCode,
			CountryCode:         country,
		},
		// Origin installation information
		&edisegment.N1{
			EntityIdentifierCode: "RG",   // Issuing office name qualifier
			Name:                 "LKNQ", // TODO: pull from TransportationOffice
			IdentificationCodeQualifier: "27", // GBLOC
			IdentificationCode:          *shipment.SourceGBLOC,
		},
		// Destination installation information
		&edisegment.N1{
			EntityIdentifierCode: "RH",   // Destination name qualifier
			Name:                 "MLNQ", // TODO: pull from TransportationOffice
			IdentificationCodeQualifier: "27", // GBLOC
			IdentificationCode:          *shipment.DestinationGBLOC,
		},
		// Accounting info
		&edisegment.FA1{
			AgencyQualifierCode: edisegment.AffiliationToAgency[*affiliation],
		},
		&edisegment.FA2{
			BreakdownStructureDetailCode: "TA", // TAC
			FinancialInformationCode:     *tac,
		},
		&edisegment.L10{
			Weight:          108.2, // TODO: real weight
			WeightQualifier: "B",   // Billing weight
			WeightUnitCode:  "L",   // Pounds
		},
	}, nil
}

func getLineItemSegments(shipmentWithCost rateengine.CostByShipment) ([]edisegment.Segment, error) {
	// follows HL loop (p.13) in https://www.ustranscom.mil/cmd/associated/dteb/files/transportationics/dt858c41.pdf
	// HL segment: p. 51
	// L0 segment: p. 77
	// L1 segment: p. 82

	lineItems := shipmentWithCost.Shipment.ShipmentLineItems

	// TODO: For the moment, we are explicitly grabbing the line items for linehaul, pack, etc.
	// TODO: We ultimately need to process all line items and hopefully abstract out their processing.
	// TODO: See https://www.pivotaltracker.com/story/show/162065870

	var segments []edisegment.Segment

	linehaulSegments, err := generateLinehaulSegments(lineItems)
	if err != nil {
		return nil, err
	}
	segments = append(segments, linehaulSegments...)

	fullPackSegments, err := generateFullPackSegments(lineItems)
	if err != nil {
		return nil, err
	}
	segments = append(segments, fullPackSegments...)

	// TODO: We are missing full unpack (no "105C" currently in our tariff400ng_items table)
	// TODO: Currently, the pack shipment line item covers the charge for both pack/unpack.
	// fullUnpackSegments, err := generateFullUnpackSegments(lineItems)
	// if err != nil {
	//     return nil, err
	// }
	// segments = append(segments, fullUnpackSegments...)

	originServiceSegments, err := generateOriginServiceSegments(lineItems)
	if err != nil {
		return nil, err
	}
	segments = append(segments, originServiceSegments...)

	destinationServiceSegments, err := generateDestinationServiceSegments(lineItems)
	if err != nil {
		return nil, err
	}
	segments = append(segments, destinationServiceSegments...)

	// TODO: We haven't migrated fuel surcharge yet ("16A") to use shipment line items.
	fuelLinehaulSegments, err := generateFuelLinehaulSegments(lineItems)
	if err != nil {
		return nil, err
	}
	segments = append(segments, fuelLinehaulSegments...)

	return segments, nil
}

func generateLinehaulSegments(lineItems []models.ShipmentLineItem) ([]edisegment.Segment, error) {
	lineItem, err := findLineItemByCode(lineItems, "LHS")
	if err != nil {
		return nil, err
	}

	return []edisegment.Segment{
		// Linehaul. Not sure why this uses the 303 code, but that's what I saw from DPS
		&edisegment.HL{
			HierarchicalIDNumber:  "303", // Accessorial services performed at origin
			HierarchicalLevelCode: "SS",  // Services
		},
		&edisegment.L0{
			LadingLineItemNumber:   1,
			BilledRatedAsQuantity:  1,
			BilledRatedAsQualifier: "FR", // Flat rate
		},
		&edisegment.L1{
			FreightRate:        0,    // TODO: placeholder for now
			RateValueQualifier: "RC", // Rate
			Charge:             lineItem.AmountCents.ToDollarFloat(),
			SpecialChargeDescription: "LHS", // Linehaul
		},
	}, nil
}

func generateFullPackSegments(lineItems []models.ShipmentLineItem) ([]edisegment.Segment, error) {
	lineItem, err := findLineItemByCode(lineItems, "105A")
	if err != nil {
		return nil, err
	}

	return []edisegment.Segment{
		// Full pack
		&edisegment.HL{
			HierarchicalIDNumber:  "303", // Accessorial services performed at origin
			HierarchicalLevelCode: "SS",  // Services
		},
		&edisegment.L0{
			LadingLineItemNumber: 1,
			Weight:               lineItem.Quantity1.ToUnitFloat(),
			WeightQualifier:      "B", // Billed weight
			WeightUnitCode:       "L", // Pounds
		},
		&edisegment.L1{
			FreightRate:        65.77, // TODO: placeholder for now
			RateValueQualifier: "RC",  // Rate
			Charge:             lineItem.AmountCents.ToDollarFloat(),
			SpecialChargeDescription: "105A", // Full pack
		},
	}, nil
}

func generateFullUnpackSegments(lineItems []models.ShipmentLineItem) ([]edisegment.Segment, error) {
	lineItem, err := findLineItemByCode(lineItems, "105C")
	if err != nil {
		return nil, err
	}

	return []edisegment.Segment{
		// Full unpack
		&edisegment.HL{
			HierarchicalIDNumber:  "304", // Accessorial services performed at destination
			HierarchicalLevelCode: "SS",  // Services
		},
		&edisegment.L0{
			LadingLineItemNumber: 1,
			Weight:               lineItem.Quantity1.ToUnitFloat(),
			WeightQualifier:      "B", // Billed weight
			WeightUnitCode:       "L", // Pounds
		},
		&edisegment.L1{
			FreightRate:        65.77, // TODO: placeholder for now
			RateValueQualifier: "RC",  // Rate
			Charge:             lineItem.AmountCents.ToDollarFloat(),
			SpecialChargeDescription: "105C", // unpack TODO: verify that GEX can recognize 105C (unpack used to be included with pack above)
		},
	}, nil
}

func generateOriginServiceSegments(lineItems []models.ShipmentLineItem) ([]edisegment.Segment, error) {
	lineItem, err := findLineItemByCode(lineItems, "135A")
	if err != nil {
		return nil, err
	}

	return []edisegment.Segment{
		// Origin service charge
		&edisegment.HL{
			HierarchicalIDNumber:  "303", // Accessorial services performed at origin
			HierarchicalLevelCode: "SS",  // Services
		},
		&edisegment.L0{
			LadingLineItemNumber: 1,
			Weight:               lineItem.Quantity1.ToUnitFloat(),
			WeightQualifier:      "B", // Billed weight
			WeightUnitCode:       "L", // Pounds
		},
		&edisegment.L1{
			FreightRate:        4.07, // TODO: placeholder for now
			RateValueQualifier: "RC", // Rate
			Charge:             lineItem.AmountCents.ToDollarFloat(),
			SpecialChargeDescription: "135A", // Origin service charge
		},
	}, nil
}

func generateDestinationServiceSegments(lineItems []models.ShipmentLineItem) ([]edisegment.Segment, error) {
	lineItem, err := findLineItemByCode(lineItems, "135B")
	if err != nil {
		return nil, err
	}

	return []edisegment.Segment{
		// Destination service charge
		&edisegment.HL{
			HierarchicalIDNumber:  "304", // Accessorial services performed at destination
			HierarchicalLevelCode: "SS",  // Services
		},
		&edisegment.L0{
			LadingLineItemNumber: 1,
			Weight:               lineItem.Quantity1.ToUnitFloat(),
			WeightQualifier:      "B", // Billed weight
			WeightUnitCode:       "L", // Pounds
		},
		&edisegment.L1{
			FreightRate:        4.07, // TODO: placeholder for now
			RateValueQualifier: "RC", // Rate
			Charge:             lineItem.AmountCents.ToDollarFloat(),
			SpecialChargeDescription: "135B", // TODO: check if correct for Destination service charge
		},
	}, nil
}

func generateFuelLinehaulSegments(lineItems []models.ShipmentLineItem) ([]edisegment.Segment, error) {
	// TODO: We haven't migrated fuel surcharge yet ("16A") to use shipment line items.
	// lineItem, err := findLineItemByCode(lineItems, "16A")
	// if err != nil {
	//     return nil, err
	// }

	return []edisegment.Segment{
		// Fuel surcharge - linehaul
		&edisegment.HL{
			HierarchicalIDNumber:  "303", // Accessorial services performed at origin
			HierarchicalLevelCode: "SS",  // Services
		},
		&edisegment.L0{
			LadingLineItemNumber:   1,
			BilledRatedAsQuantity:  1,
			BilledRatedAsQualifier: "FR", // Flat rate
		},
		&edisegment.L1{
			FreightRate:        0.03,   // TODO: placeholder for now
			RateValueQualifier: "RC",   // Rate
			Charge:             227.42, // TODO: add a calculation of this value to rate engine
			SpecialChargeDescription: "16A", // Fuel surchage - linehaul
		},
	}, nil
}

func findLineItemByCode(lineItems []models.ShipmentLineItem, code string) (models.ShipmentLineItem, error) {
	for i := range lineItems {
		if lineItems[i].Tariff400ngItem.Code == code {
			return lineItems[i], nil
		}
	}

	return models.ShipmentLineItem{}, errors.Errorf("Could not find shipment line item with code %s", code)
}
