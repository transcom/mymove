package ediinvoice

import (
	"bytes"
	"fmt"
	"go.uber.org/zap"
	"time"

	"github.com/facebookgo/clock"
	"github.com/gobuffalo/pop"
	"github.com/pkg/errors"

	"github.com/transcom/mymove/pkg/db/sequence"
	"github.com/transcom/mymove/pkg/edi"
	"github.com/transcom/mymove/pkg/edi/segment"
	"github.com/transcom/mymove/pkg/models"
)

const dateFormat = "20060102"
const timeFormat = "1504"
const senderCode = "MYMOVE"

const receiverCode = "8004171844" // Syncada

// ICNSequenceName used to query Interchange Control Numbers from DB
const ICNSequenceName = "interchange_control_number"

const rateValueQualifier = "RC"    // Rate
const hierarchicalLevelCode = "SS" // Services
const weightQualifier = "B"        // Billed Weight
const weightUnitCode = "L"         // Pounds
const ladingLineItemNumber = 1
const billedRatedAsQuantity = 1

// Place holders that currently exist TODO: Replace this constants with real value
const freightRate = 4.07

var logger *zap.Logger

// Invoice858C holds all the segments that are generated
type Invoice858C struct {
	ISA      edisegment.ISA
	GS       edisegment.GS
	Shipment []edisegment.Segment
	GE       edisegment.GE
	IEA      edisegment.IEA
}

// Segments returns the invoice as an array of rows (string arrays),
// each containing a segment, to prepare it for writing
func (invoice Invoice858C) Segments() [][]string {
	records := [][]string{
		invoice.ISA.StringArray(),
		invoice.GS.StringArray(),
	}

	for _, line := range invoice.Shipment {
		records = append(records, line.StringArray())
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
func Generate858C(shipment models.Shipment, db *pop.Connection, sendProductionInvoice bool, clock clock.Clock) (Invoice858C, error) {
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
		Date:                     currentTime.Format(dateFormat),
		Time:                     currentTime.Format(timeFormat),
		GroupControlNumber:       interchangeControlNumber,
		ResponsibleAgencyCode:    "X", // Accredited Standards Committee X12
		Version:                  "004010",
	}

	shipmentSegments, err := generate858CShipment(shipment, 1)
	if err != nil {
		return invoice, err
	}
	invoice.Shipment = shipmentSegments

	invoice.GE = edisegment.GE{
		NumberOfTransactionSetsIncluded: 1,
		GroupControlNumber:              interchangeControlNumber,
	}
	invoice.IEA = edisegment.IEA{
		NumberOfIncludedFunctionalGroups: 1,
		InterchangeControlNumber:         interchangeControlNumber,
	}

	return invoice, nil
}

func generate858CShipment(shipment models.Shipment, sequenceNum int) ([]edisegment.Segment, error) {
	transactionNumber := fmt.Sprintf("%04d", sequenceNum)
	segments := []edisegment.Segment{
		&edisegment.ST{
			TransactionSetIdentifierCode: "858",
			TransactionSetControlNumber:  transactionNumber,
		},
	}

	headingSegments, err := getHeadingSegments(shipment, sequenceNum)
	if err != nil {
		return segments, err
	}
	segments = append(segments, headingSegments...)

	lineItemSegments, err := getLineItemSegments(shipment)
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

func getHeadingSegments(shipment models.Shipment, sequenceNum int) ([]edisegment.Segment, error) {
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
	originTransportationOfficeName := shipment.ServiceMember.DutyStation.TransportationOffice.Name
	if originTransportationOfficeName == "" {
		return segments, errors.New("Transportation Office Name is missing (for N102)")
	}
	destinationTransportationOfficeName := shipment.Move.Orders.NewDutyStation.TransportationOffice.Name
	if destinationTransportationOfficeName == "" {
		return segments, errors.New("Transportation Office Name is missing (for N102)")
	}
	weightLbs := shipment.NetWeight
	if weightLbs == nil {
		return segments, errors.New("Shipment is missing the NetWeight")
	}
	netCentiWeight := float64(*weightLbs) / 100 // convert to CW

	acceptedOffer, err := shipment.AcceptedShipmentOffer()
	if err != nil || acceptedOffer == nil {
		return segments, errors.Wrap(err, "Error retrieving ACCEPTED ShipmentOffer for EDI generator")
	}

	scac, err := acceptedOffer.SCAC()
	if err != nil {
		return segments, err
	}

	supplierID, err := acceptedOffer.SupplierID()
	if err != nil {
		return segments, err
	}

	return []edisegment.Segment{
		&edisegment.BX{
			TransactionSetPurposeCode:    "00", // Original
			TransactionMethodTypeCode:    "J",  // Motor
			ShipmentMethodOfPayment:      "PP", // Prepaid by seller
			ShipmentIdentificationNumber: *GBL,
			StandardCarrierAlphaCode:     scac,
			ShipmentQualifier:            "4", // HHG Government Bill of Lading
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
			ReferenceIdentificationQualifier: "PQ", // Payee code
			ReferenceIdentification:          *supplierID,
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
			EntityIdentifierCode:        "RG", // Issuing office name qualifier
			Name:                        originTransportationOfficeName,
			IdentificationCodeQualifier: "27", // GBLOC
			IdentificationCode:          *shipment.SourceGBLOC,
		},
		// Destination installation information
		&edisegment.N1{
			EntityIdentifierCode:        "RH", // Destination name qualifier
			Name:                        destinationTransportationOfficeName,
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
			Weight:          netCentiWeight,
			WeightQualifier: "B", // Billing weight
			WeightUnitCode:  "L", // Pounds
		},
	}, nil
}

func getLineItemSegments(shipment models.Shipment) ([]edisegment.Segment, error) {
	// follows HL loop (p.13) in https://www.ustranscom.mil/cmd/associated/dteb/files/transportationics/dt858c41.pdf
	// HL segment: p. 51
	// L0 segment: p. 77
	// L1 segment: p. 82

	lineItems := shipment.ShipmentLineItems
	weightLbs := shipment.NetWeight
	if weightLbs == nil {
		return nil, errors.New("Shipment is missing the NetWeight")
	}
	netCentiWeight := float64(*weightLbs) / 100 // convert to CW

	//Initialize empty collection of segments
	var segments []edisegment.Segment

	// Iterate over lineitems
	for _, lineItem := range lineItems {
		// Some hardcoded values that are being used

		// Initialize empty edisegment
		var tariffSegments []edisegment.Segment

		// Build and put together the segments
		hlSegment := MakeHLSegment(lineItem)
		l0Segment := MakeL0Segment(lineItem, netCentiWeight)
		l1Segment := MakeL1Segment(lineItem)
		tariffSegments = append(tariffSegments, hlSegment, l0Segment, l1Segment)

		segments = append(segments, tariffSegments...)

	}

	return segments, nil
}

// MakeHLSegment builds HL segment based on shipment line item input.
func MakeHLSegment(lineItem models.ShipmentLineItem) *edisegment.HL {
	// Initialize hierarchicalLevelCode
	var hierarchicalLevelID string

	// Determine HierarchicalLevelCode
	switch lineItem.Location {

	case models.ShipmentLineItemLocationORIGIN:
		hierarchicalLevelID = "303"

	case models.ShipmentLineItemLocationDESTINATION:
		hierarchicalLevelID = "304"

	case models.ShipmentLineItemLocationNEITHER:
		hierarchicalLevelID = "303"
	}
	return &edisegment.HL{
		HierarchicalIDNumber:  hierarchicalLevelID,
		HierarchicalLevelCode: hierarchicalLevelCode,
	}

}

// MakeL0Segment builds L0 segment based on shipment line item input and shipment centiweight input.
func MakeL0Segment(lineItem models.ShipmentLineItem, netCentiWeight float64) *edisegment.L0 {
	// Using Maps to group up MeasurementUnit types into categories
	unitBasedMeasurementUnits := map[models.Tariff400ngItemMeasurementUnit]int{
		models.Tariff400ngItemMeasurementUnitFLATRATE:       0,
		models.Tariff400ngItemMeasurementUnitEACH:           0,
		models.Tariff400ngItemMeasurementUnitHOURS:          0,
		models.Tariff400ngItemMeasurementUnitDAYS:           0,
		models.Tariff400ngItemMeasurementUnitCUBICFOOT:      0,
		models.Tariff400ngItemMeasurementUnitFUELPERCENTAGE: 0,
		models.Tariff400ngItemMeasurementUnitCONTAINER:      0,
		models.Tariff400ngItemMeasurementUnitMONETARYVALUE:  0,
		models.Tariff400ngItemMeasurementUnitNONE:           0,
	}

	weightBasedMeasurements := map[models.Tariff400ngItemMeasurementUnit]int{
		models.Tariff400ngItemMeasurementUnitWEIGHT: 0,
	}

	measurementUnit := lineItem.Tariff400ngItem.MeasurementUnit1

	// This will check if the Measurement unit is in one of the maps above.
	// Doing this allows us to have two generic paths based on groups of MeasurementUnits
	// This is a way to do something a-kin to OR logic in our comparison for the category.
	_, isUnitBased := unitBasedMeasurementUnits[measurementUnit]
	_, isWeightBased := weightBasedMeasurements[measurementUnit]

	if isUnitBased {

		actualBilledRatedAsQuantity := float64(billedRatedAsQuantity)
		if lineItem.Tariff400ngItem.MeasurementUnit1 != models.Tariff400ngItemMeasurementUnitFLATRATE {
			actualBilledRatedAsQuantity = lineItem.Quantity1.ToUnitFloat()
		}

		return &edisegment.L0{
			LadingLineItemNumber:   ladingLineItemNumber,
			BilledRatedAsQuantity:  actualBilledRatedAsQuantity,
			BilledRatedAsQualifier: string(measurementUnit),
		}

	} else if isWeightBased {
		var weight float64

		if lineItem.Tariff400ngItem.RequiresPreApproval {
			weight = lineItem.Quantity1.ToUnitFloat()

		} else {
			weight = netCentiWeight
		}

		return &edisegment.L0{
			LadingLineItemNumber: ladingLineItemNumber,
			Weight:               weight,
			WeightQualifier:      weightQualifier,
			WeightUnitCode:       weightUnitCode,
		}

	} else {
		logger.Error(string(measurementUnit) + "Used with " +
			lineItem.ID.String() + " is an EDI measurement unit we're not prepared for.")
		return nil
	}

}

// MakeL1Segment builds L1 segment based on shipment lineitem input.
func MakeL1Segment(lineItem models.ShipmentLineItem) *edisegment.L1 {

	return &edisegment.L1{
		FreightRate:              lineItem.AppliedRate.ToDollarFloat(),
		RateValueQualifier:       rateValueQualifier,
		Charge:                   lineItem.AmountCents.ToDollarFloat(),
		SpecialChargeDescription: lineItem.Tariff400ngItem.Code,
	}

}
