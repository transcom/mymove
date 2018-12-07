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
	"github.com/transcom/mymove/pkg/rateengine"
)

const dateFormat = "20060102"
const timeFormat = "1504"
const senderCode = "MYMOVE"

//const senderCode = "W28GPR-DPS"   // TODO: update with ours when US Bank gets it to us
const receiverCode = "8004171844" // Syncada

// ICNSequenceName used to query Interchange Control Numbers from DB
const ICNSequenceName = "interchange_control_number"

var logger *zap.Logger

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
		Date:                     currentTime.Format(dateFormat),
		Time:                     currentTime.Format(timeFormat),
		GroupControlNumber:       interchangeControlNumber,
		ResponsibleAgencyCode:    "X", // Accredited Standards Committee X12
		Version:                  "004010",
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
		return nil, errors.New("Shipment is missing the NetWeight")
	}
	netCentiWeight := float64(*weightLbs) / 100 // convert to CW

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

func getLineItemSegments(shipmentWithCost rateengine.CostByShipment) ([]edisegment.Segment, error) {
	// follows HL loop (p.13) in https://www.ustranscom.mil/cmd/associated/dteb/files/transportationics/dt858c41.pdf
	// HL segment: p. 51
	// L0 segment: p. 77
	// L1 segment: p. 82

	lineItems := shipmentWithCost.Shipment.ShipmentLineItems
	shipment := shipmentWithCost.Shipment
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
		const rateValueQualifier = "RC"    // Rate
		const hierarchicalLevelCode = "SS" // Services
		const weightQualifier = "B"        // Billed Weight
		const weightUnitCode = "L"         // Pounds
		const ladingLineItemNumber = 1
		const billedRatedAsQuantity = 1

		// Place holders that currently exist TODO: Replace this constants with real value
		const freightRate = 4.07

		// Initialize empty edisegment
		var tariffSegment []edisegment.Segment

		// Initialize hierarchicalLevelCode
		var hierarchicalLevelID string

		// Determine HierarchicalLevelCode
		switch lineItem.Location {

		case models.ShipmentLineItemLocationORIGIN:
			hierarchicalLevelID = "304"

		case models.ShipmentLineItemLocationDESTINATION:
			hierarchicalLevelID = "303"

		}

		// Using Maps to group up MeasurementUnit types into categories
		unitBasedMeasurementUnits := map[models.Tariff400ngItemMeasurementUnit]int{
			models.Tariff400ngItemMeasurementUnitFLATRATE: 0,
			models.Tariff400ngItemMeasurementUnitEACH:     0,
		}

		weightBasedMeasurements := map[models.Tariff400ngItemMeasurementUnit]int{
			models.Tariff400ngItemMeasurementUnitWEIGHT: 0,
		}

		// This will check if the Measurement unit is in one of the maps above.
		// Doing this allows us to have two generic paths based on groups of MeasurementUnits
		// This is a way to do something a-kin to OR logic in our comparison for the category.
		_, isUnitBased := unitBasedMeasurementUnits[lineItem.Tariff400ngItem.MeasurementUnit1]
		_, isWeightBased := weightBasedMeasurements[lineItem.Tariff400ngItem.MeasurementUnit1]

		tariffSegment = []edisegment.Segment{
			&edisegment.HL{
				HierarchicalIDNumber:  hierarchicalLevelID,
				HierarchicalLevelCode: hierarchicalLevelCode,
			},
		}

		if isUnitBased {
			unitBasedSegment := &edisegment.L0{
				LadingLineItemNumber:   ladingLineItemNumber,
				BilledRatedAsQuantity:  billedRatedAsQuantity,
				BilledRatedAsQualifier: string(lineItem.Tariff400ngItem.MeasurementUnit1),
			}
			tariffSegment = append(tariffSegment, unitBasedSegment)

		} else if isWeightBased {
			var weight float64

			if lineItem.Tariff400ngItem.RequiresPreApproval {
				weight = lineItem.Quantity1.ToUnitFloat()

			} else {
				weight = netCentiWeight
			}

			weightBasedSegment := &edisegment.L0{
				LadingLineItemNumber: ladingLineItemNumber,
				Weight:               weight,
				WeightQualifier:      weightQualifier,
				WeightUnitCode:       weightUnitCode,
			}
			tariffSegment = append(tariffSegment, weightBasedSegment)

		} else {
			logger.Error(string(lineItem.Tariff400ngItem.MeasurementUnit1) + "Used with " +
				lineItem.ID.String() + " is an EDI meaasurement unit we're not prepared for.")
		}

		segmentL1 := &edisegment.L1{
			FreightRate:              freightRate, //TODO: Replace this with the actual rate. It's a placeholder.
			RateValueQualifier:       rateValueQualifier,
			Charge:                   lineItem.AmountCents.ToDollarFloat(),
			SpecialChargeDescription: lineItem.Tariff400ngItem.Code,
		}

		tariffSegment = append(tariffSegment, segmentL1)

		segments = append(segments, tariffSegment...)

	}

	return segments, nil
}
