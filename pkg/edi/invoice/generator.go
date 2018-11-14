package ediinvoice

import (
	"errors"
	"fmt"

	"github.com/facebookgo/clock"
	"github.com/gobuffalo/pop"

	"github.com/transcom/mymove/pkg/edi/segment"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/rateengine"
)

const dateFormat = "20060102"
const timeFormat = "1504"
const senderCode = "MYMOVE"

//const senderCode = "W28GPR-DPS"   // TODO: update with ours when US Bank gets it to us
const receiverCode = "8004171844" // Syncada

// Invoice858C holds all the segments that are generated
type Invoice858C struct {
	ISA       edisegment.ISA
	GS        edisegment.GS
	Shipments []InvoiceShipment
	SE        edisegment.SE
	GE        edisegment.GE
	IEA       edisegment.IEA
}

// InvoiceShipment holds the 858C items per shipment
type InvoiceShipment struct {
	BX                   *edisegment.BX
	N9DY                 *edisegment.N9
	N9CN                 *edisegment.N9
	N9PQ                 *edisegment.N9
	N9OQ                 *edisegment.N9
	N1SF                 *edisegment.N1
	N3                   *edisegment.N3
	N4                   *edisegment.N4
	N1RG                 *edisegment.N1
	N1RH                 *edisegment.N1
	FA1                  *edisegment.FA1
	FA2                  *edisegment.FA2
	L10                  *edisegment.L10
	HLLinehaul           *edisegment.HL
	L0Linehaul           *edisegment.L0
	L1Linehaul           *edisegment.L1
	HLFullPack           *edisegment.HL
	L0FullPack           *edisegment.L0
	L1FullPack           *edisegment.L1
	HLFullUnpack         *edisegment.HL
	L0FullUnpack         *edisegment.L0
	L1FullUnpack         *edisegment.L1
	HLOriginService      *edisegment.HL
	L0OriginService      *edisegment.L0
	L1OriginService      *edisegment.L1
	HLDestinationService *edisegment.HL
	L0DestinationService *edisegment.L0
	L1DestinationService *edisegment.L1
	HLFuel               *edisegment.HL
	L0Fuel               *edisegment.L0
	L1Fuel               *edisegment.L1
	SE                   *edisegment.SE
}

// Generate858C generates an EDI X12 858C transaction set
func Generate858C(shipmentsAndCosts []rateengine.CostByShipment, db *pop.Connection, sendProductionInvoice bool, clock clock.Clock) (Invoice858C, error) {
	interchangeControlNumber := 1 //TODO: increment this
	currentTime := clock.Now()
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
		GroupControlNumber:    1,
		ResponsibleAgencyCode: "X", // Accredited Standards Committee X12
		Version:               "004010",
	}

	var shipments []models.Shipment

	invoice.Shipments = make([]InvoiceShipment, 0)
	for index, shipmentWithCost := range shipmentsAndCosts {
		shipment := shipmentWithCost.Shipment

		invoiceShipment, err := generate858CShipment(shipmentWithCost, index+1)
		if err != nil {
			return Invoice858C{}, err
		}
		invoice.Shipments = append(invoice.Shipments, invoiceShipment)
		shipments = append(shipments, shipment)
	}

	invoice.GE = edisegment.GE{
		NumberOfTransactionSetsIncluded: len(shipments),
		GroupControlNumber:              1,
	}
	invoice.IEA = edisegment.IEA{
		NumberOfIncludedFunctionalGroups: 1,
		InterchangeControlNumber:         interchangeControlNumber,
	}

	return invoice, nil
}

func generate858CShipment(shipmentWithCost rateengine.CostByShipment, sequenceNum int) (InvoiceShipment, error) {
	invoiceShipment := InvoiceShipment{}
	transactionNumber := fmt.Sprintf("%04d", sequenceNum)
	segments := []edisegment.Segment{
		&edisegment.ST{
			TransactionSetIdentifierCode: "858",
			TransactionSetControlNumber:  transactionNumber,
		},
	}

	err := setHeadingSegments(&invoiceShipment, shipmentWithCost, sequenceNum)
	if err != nil {
		return invoiceShipment, err
	}

	err = setLineItemSegments(&invoiceShipment, shipmentWithCost)
	if err != nil {
		return invoiceShipment, err
	}

	invoiceShipment.SE = &edisegment.SE{
		NumberOfIncludedSegments:    len(segments) + 1, // Include SE in count
		TransactionSetControlNumber: transactionNumber,
	}

	return invoiceShipment, nil
}

func setHeadingSegments(invoiceShipment *InvoiceShipment, shipmentWithCost rateengine.CostByShipment, sequenceNum int) error {
	shipment := shipmentWithCost.Shipment
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
		return errors.New("Shipment is missing pick up address")
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
		return errors.New("Orders is missing orders number")
	}
	tac := orders.TAC
	if tac == nil {
		return errors.New("Orders is missing TAC")
	}
	affiliation := shipment.ServiceMember.Affiliation
	if shipment.ServiceMember.Affiliation == nil {
		return errors.New("Service member is missing affiliation")
	}
	GBL := shipment.GBLNumber
	if GBL == nil {
		return errors.New("GBL Number is missing for Shipment Identification Number (BX04)")
	}

	invoiceShipment.BX = &edisegment.BX{
		TransactionSetPurposeCode:    "00", // Original
		TransactionMethodTypeCode:    "J",  // Motor
		ShipmentMethodOfPayment:      "PP", // Prepaid by seller
		ShipmentIdentificationNumber: *GBL,
		StandardCarrierAlphaCode:     "MCCG", // TODO: real SCAC
		ShipmentQualifier:            "4",    // HHG Government Bill of Lading
	}
	invoiceShipment.N9DY = &edisegment.N9{
		ReferenceIdentificationQualifier: "DY", // DoD transportation service code #
		ReferenceIdentification:          "SC", // Shipment & cost information
	}
	invoiceShipment.N9CN = &edisegment.N9{
		ReferenceIdentificationQualifier: "CN",          // Invoice number
		ReferenceIdentification:          "ABCD00001-1", // TODO: real invoice number
	}
	invoiceShipment.N9PQ = &edisegment.N9{
		ReferenceIdentificationQualifier: "PQ",       // Payee code
		ReferenceIdentification:          "ABBV2708", // TODO: add real supplier ID
	}
	invoiceShipment.N9OQ = &edisegment.N9{
		ReferenceIdentificationQualifier: "OQ", // Order number
		ReferenceIdentification:          *ordersNumber,
		FreeFormDescription:              string(*affiliation),
		Date:                             orders.IssueDate.Format(dateFormat),
	}
	// Ship from address
	invoiceShipment.N1SF = &edisegment.N1{
		EntityIdentifierCode: "SF", // Ship From
		Name:                 name,
	}
	invoiceShipment.N3 = &edisegment.N3{
		AddressInformation1: shipment.PickupAddress.StreetAddress1,
		AddressInformation2: street2,
	}
	invoiceShipment.N4 = &edisegment.N4{
		CityName:            shipment.PickupAddress.City,
		StateOrProvinceCode: shipment.PickupAddress.State,
		PostalCode:          shipment.PickupAddress.PostalCode,
		CountryCode:         country,
	}
	// Origin installation information
	invoiceShipment.N1RG = &edisegment.N1{
		EntityIdentifierCode: "RG",   // Issuing office name qualifier
		Name:                 "LKNQ", // TODO: pull from TransportationOffice
		IdentificationCodeQualifier: "27", // GBLOC
		IdentificationCode:          "LKNQ",
	}
	// Destination installation information
	invoiceShipment.N1RH = &edisegment.N1{
		EntityIdentifierCode: "RH",   // Destination name qualifier
		Name:                 "MLNQ", // TODO: pull from TransportationOffice
		IdentificationCodeQualifier: "27", // GBLOC
		IdentificationCode:          "MLNQ",
	}
	// Accounting info
	invoiceShipment.FA1 = &edisegment.FA1{
		AgencyQualifierCode: edisegment.AffiliationToAgency[*affiliation],
	}
	invoiceShipment.FA2 = &edisegment.FA2{
		BreakdownStructureDetailCode: "TA", // TAC
		FinancialInformationCode:     *tac,
	}
	invoiceShipment.L10 = &edisegment.L10{
		Weight:          108.2, // TODO: real weight
		WeightQualifier: "B",   // Billing weight
		WeightUnitCode:  "L",   // Pounds
	}
	return nil
}

func setLineItemSegments(invoiceShipment *InvoiceShipment, shipmentWithCost rateengine.CostByShipment) error {
	// follows HL loop (p.13) in https://www.ustranscom.mil/cmd/associated/dteb/files/transportationics/dt858c41.pdf
	// HL segment: p. 51
	// L0 segment: p. 77
	// L1 segment: p. 82
	// TODO: These are sample line items, need to pull actual line items from shipment
	// that are ready to be invoiced
	cost := shipmentWithCost.Cost

	// Linehaul. Not sure why this uses the 303 code, but that's what I saw from DPS
	invoiceShipment.HLLinehaul = &edisegment.HL{
		HierarchicalIDNumber:  "303", // Accessorial services performed at origin
		HierarchicalLevelCode: "SS",  // Services
	}
	invoiceShipment.L0Linehaul = &edisegment.L0{
		LadingLineItemNumber:   1,
		BilledRatedAsQuantity:  1,
		BilledRatedAsQualifier: "FR", // Flat rate
	}
	invoiceShipment.L1Linehaul = &edisegment.L1{
		FreightRate:        0,
		RateValueQualifier: "RC", // Rate
		Charge:             cost.LinehaulCostComputation.LinehaulChargeTotal.ToDollarFloat(),
		SpecialChargeDescription: "LHS", // Linehaul
	}
	// Full pack
	invoiceShipment.HLFullPack = &edisegment.HL{
		HierarchicalIDNumber:  "303", // Accessorial services performed at origin
		HierarchicalLevelCode: "SS",  // Services
	}
	invoiceShipment.L0FullPack = &edisegment.L0{
		LadingLineItemNumber: 1,
		Weight:               108.2,
		WeightQualifier:      "B", // Billed weight
		WeightUnitCode:       "L", // Pounds
	}
	invoiceShipment.L1FullPack = &edisegment.L1{
		FreightRate:        65.77,
		RateValueQualifier: "RC", // Rate
		Charge:             cost.NonLinehaulCostComputation.PackFee.ToDollarFloat(),
		SpecialChargeDescription: "105A", // Full pack
	}
	// Full unpack
	invoiceShipment.HLFullUnpack = &edisegment.HL{
		HierarchicalIDNumber:  "304", // Accessorial services performed at destination
		HierarchicalLevelCode: "SS",  // Services
	}
	invoiceShipment.L0FullUnpack = &edisegment.L0{
		LadingLineItemNumber: 1,
		Weight:               108.2,
		WeightQualifier:      "B", // Billed weight
		WeightUnitCode:       "L", // Pounds
	}
	invoiceShipment.L1FullUnpack = &edisegment.L1{
		FreightRate:        65.77,
		RateValueQualifier: "RC", // Rate
		Charge:             cost.NonLinehaulCostComputation.UnpackFee.ToDollarFloat(),
		SpecialChargeDescription: "105C", // unpack TODO: verify that GEX can recognize 105C (unpack used to be included with pack above)
	}
	// Origin service charge
	invoiceShipment.HLOriginService = &edisegment.HL{
		HierarchicalIDNumber:  "303", // Accessorial services performed at origin
		HierarchicalLevelCode: "SS",  // Services
	}
	invoiceShipment.L0OriginService = &edisegment.L0{
		LadingLineItemNumber: 1,
		Weight:               108.2,
		WeightQualifier:      "B", // Billed weight
		WeightUnitCode:       "L", // Pounds
	}
	invoiceShipment.L1OriginService = &edisegment.L1{
		FreightRate:        4.07,
		RateValueQualifier: "RC", // Rate
		Charge:             cost.NonLinehaulCostComputation.OriginServiceFee.ToDollarFloat(),
		SpecialChargeDescription: "135A", // Origin service charge
	}
	// Destination service charge
	invoiceShipment.HLDestinationService = &edisegment.HL{
		HierarchicalIDNumber:  "304", // Accessorial services performed at destination
		HierarchicalLevelCode: "SS",  // Services
	}
	invoiceShipment.L0DestinationService = &edisegment.L0{
		LadingLineItemNumber: 1,
		Weight:               108.2,
		WeightQualifier:      "B", // Billed weight
		WeightUnitCode:       "L", // Pounds
	}
	invoiceShipment.L1DestinationService = &edisegment.L1{
		FreightRate:        4.07,
		RateValueQualifier: "RC", // Rate
		Charge:             cost.NonLinehaulCostComputation.DestinationServiceFee.ToDollarFloat(),
		SpecialChargeDescription: "135B", // TODO: check if correct for Destination service charge
	}
	// Fuel surcharge - linehaul
	invoiceShipment.HLFuel = &edisegment.HL{
		HierarchicalIDNumber:  "303", // Accessorial services performed at origin
		HierarchicalLevelCode: "SS",  // Services
	}
	invoiceShipment.L0Fuel = &edisegment.L0{
		LadingLineItemNumber:   1,
		BilledRatedAsQuantity:  1,
		BilledRatedAsQualifier: "FR", // Flat rate
	}
	invoiceShipment.L1Fuel = &edisegment.L1{
		FreightRate:        0.03,
		RateValueQualifier: "RC",   // Rate
		Charge:             227.42, // TODO: add a calculation of this value to rate engine
		SpecialChargeDescription: "16A", // Fuel surchage - linehaul
	}
	return nil
}
