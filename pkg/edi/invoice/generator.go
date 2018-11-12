package ediinvoice

import (
	"fmt"
	"time"

	"github.com/gobuffalo/pop"
	"github.com/pkg/errors"

	"github.com/transcom/mymove/pkg/db/sequence"
	"github.com/transcom/mymove/pkg/edi/segment"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/rateengine"
)

const delimiter = "*"
const dateFormat = "20060102"
const timeFormat = "1504"
const senderCode = "MYMOVE"

//const senderCode = "W28GPR-DPS"   // TODO: update with ours when US Bank gets it to us
const receiverCode = "8004171844" // Syncada

// ICNSequenceName used to query Interchange Control Numbers from DB
const ICNSequenceName = "interchange_control_number"

// Generate858C generates an EDI X12 858C transaction set
func Generate858C(shipmentsAndCosts []rateengine.CostByShipment, db *pop.Connection, sendProductionInvoice bool) (string, error) {
	interchangeControlNumber, err := getNextICN(db)
	if err != nil {
		return "", errors.Wrap(err, fmt.Sprintf("Failed to get next Interchange Control Number"))
	}

	currentTime := time.Now()
	var usageIndicator string

	if sendProductionInvoice {
		usageIndicator = "P"
	} else {
		usageIndicator = "T"
	}

	isa := edisegment.ISA{
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
	gs := edisegment.GS{
		FunctionalIdentifierCode: "SI", // Shipment Information (858)
		ApplicationSendersCode:   senderCode,
		ApplicationReceiversCode: receiverCode,
		Date:                  currentTime.Format(dateFormat),
		Time:                  currentTime.Format(timeFormat),
		GroupControlNumber:    1,
		ResponsibleAgencyCode: "X", // Accredited Standards Committee X12
		Version:               "004010",
	}
	transaction := isa.String(delimiter) + gs.String(delimiter)

	var shipments []models.Shipment

	for index, shipmentWithCost := range shipmentsAndCosts {
		shipment := shipmentWithCost.Shipment

		shipment858c, err := generate858CShipment(shipmentWithCost, index+1)
		if err != nil {
			return "", err
		}
		transaction += shipment858c
		shipments = append(shipments, shipment)
	}

	ge := edisegment.GE{
		NumberOfTransactionSetsIncluded: len(shipments),
		GroupControlNumber:              1,
	}
	iea := edisegment.IEA{
		NumberOfIncludedFunctionalGroups: 1,
		InterchangeControlNumber:         interchangeControlNumber,
	}

	transaction += (ge.String(delimiter) + iea.String(delimiter))

	return transaction, nil
}

func generate858CShipment(shipmentWithCost rateengine.CostByShipment, sequenceNum int) (string, error) {
	transactionNumber := fmt.Sprintf("%04d", sequenceNum)
	segments := []edisegment.Segment{
		&edisegment.ST{
			TransactionSetIdentifierCode: "858",
			TransactionSetControlNumber:  transactionNumber,
		},
	}

	headingSegments, err := getHeadingSegments(shipmentWithCost, sequenceNum)
	if err != nil {
		return "", err
	}
	segments = append(segments, headingSegments...)

	lineItemSegments, err := getLineItemSegments(shipmentWithCost)
	if err != nil {
		return "", err
	}
	segments = append(segments, lineItemSegments...)

	segments = append(
		segments,
		&edisegment.SE{
			NumberOfIncludedSegments:    len(segments) + 1, // Include SE in count
			TransactionSetControlNumber: transactionNumber,
		},
	)

	transaction := ""
	for _, seg := range segments {
		transaction += seg.String(delimiter)
	}

	return transaction, nil
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
			IdentificationCode:          "LKNQ",
		},
		// Destination installation information
		&edisegment.N1{
			EntityIdentifierCode: "RH",   // Destination name qualifier
			Name:                 "MLNQ", // TODO: pull from TransportationOffice
			IdentificationCodeQualifier: "27", // GBLOC
			IdentificationCode:          "MLNQ",
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
	// TODO: These are sample line items, need to pull actual line items from shipment
	// that are ready to be invoiced
	cost := shipmentWithCost.Cost

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
			FreightRate:        0,
			RateValueQualifier: "RC", // Rate
			Charge:             cost.LinehaulCostComputation.LinehaulChargeTotal.ToDollarFloat(),
			SpecialChargeDescription: "LHS", // Linehaul
		},
		// Full pack
		&edisegment.HL{
			HierarchicalIDNumber:  "303", // Accessorial services performed at origin
			HierarchicalLevelCode: "SS",  // Services
		},
		&edisegment.L0{
			LadingLineItemNumber: 1,
			Weight:               108.2,
			WeightQualifier:      "B", // Billed weight
			WeightUnitCode:       "L", // Pounds
		},
		&edisegment.L1{
			FreightRate:        65.77,
			RateValueQualifier: "RC", // Rate
			Charge:             cost.NonLinehaulCostComputation.PackFee.ToDollarFloat(),
			SpecialChargeDescription: "105A", // Full pack
		},
		// Full unpack
		&edisegment.HL{
			HierarchicalIDNumber:  "304", // Accessorial services performed at destination
			HierarchicalLevelCode: "SS",  // Services
		},
		&edisegment.L0{
			LadingLineItemNumber: 1,
			Weight:               108.2,
			WeightQualifier:      "B", // Billed weight
			WeightUnitCode:       "L", // Pounds
		},
		&edisegment.L1{
			FreightRate:        65.77,
			RateValueQualifier: "RC", // Rate
			Charge:             cost.NonLinehaulCostComputation.UnpackFee.ToDollarFloat(),
			SpecialChargeDescription: "105C", // unpack TODO: verify that GEX can recognize 105C (unpack used to be included with pack above)
		},
		// Origin service charge
		&edisegment.HL{
			HierarchicalIDNumber:  "303", // Accessorial services performed at origin
			HierarchicalLevelCode: "SS",  // Services
		},
		&edisegment.L0{
			LadingLineItemNumber: 1,
			Weight:               108.2,
			WeightQualifier:      "B", // Billed weight
			WeightUnitCode:       "L", // Pounds
		},
		&edisegment.L1{
			FreightRate:        4.07,
			RateValueQualifier: "RC", // Rate
			Charge:             cost.NonLinehaulCostComputation.OriginServiceFee.ToDollarFloat(),
			SpecialChargeDescription: "135A", // Origin service charge
		},
		// Destination service charge
		&edisegment.HL{
			HierarchicalIDNumber:  "304", // Accessorial services performed at destination
			HierarchicalLevelCode: "SS",  // Services
		},
		&edisegment.L0{
			LadingLineItemNumber: 1,
			Weight:               108.2,
			WeightQualifier:      "B", // Billed weight
			WeightUnitCode:       "L", // Pounds
		},
		&edisegment.L1{
			FreightRate:        4.07,
			RateValueQualifier: "RC", // Rate
			Charge:             cost.NonLinehaulCostComputation.DestinationServiceFee.ToDollarFloat(),
			SpecialChargeDescription: "135B", // TODO: check if correct for Destination service charge
		},
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
			FreightRate:        0.03,
			RateValueQualifier: "RC",   // Rate
			Charge:             227.42, // TODO: add a calculation of this value to rate engine
			SpecialChargeDescription: "16A", // Fuel surchage - linehaul
		},
	}, nil
}

// GetNextICN is a public wrapper around getNextICN for testing
// See: link issue when created
var GetNextICN = getNextICN

// getNextICN returns the next Interchange Control Number in a PostgreSQL sequence
func getNextICN(db *pop.Connection) (int64, error) {
	return sequence.NextVal(db, ICNSequenceName)
}
