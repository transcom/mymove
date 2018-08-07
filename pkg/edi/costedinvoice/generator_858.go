package edicostedinvoice

import (
	"errors"
	"fmt"
	"time"

	"github.com/gobuffalo/pop"

	"github.com/transcom/mymove/pkg/edi/segment"
	"github.com/transcom/mymove/pkg/models"
)

const delimiter = "~"
const dateFormat = "20060102"
const timeFormat = "1504"
const senderCode = "W28GPR-DPS"   // TODO: update with our own after talking to US Bank
const receiverCode = "8004171844" // Syncada

// Generate858C generates an EDI X12 858C transaction set
func Generate858C(shipments []models.Shipment, db *pop.Connection) (string, error) {
	currentTime := time.Now()
	isa := edisegment.ISA{
		AuthorizationInformationQualifier: "00", // No authorization information
		AuthorizationInformation:          fmt.Sprintf("%010d", 0),
		SecurityInformationQualifier:      "00", // No security information
		SecurityInformation:               fmt.Sprintf("%010d", 0),
		InterchangeSenderIDQualifier:      "ZZ",
		InterchangeSenderID:               fmt.Sprintf("%-15v", senderCode), // Must be 15 characters
		InterchangeReceiverIDQualifier:    "12",
		InterchangeReceiverID:             fmt.Sprintf("%-15v", "8004171844"), // Must be 15 characters
		InterchangeDate:                   currentTime.Format("060102"),
		InterchangeTime:                   currentTime.Format(timeFormat),
		InterchangeControlStandards:       "U",
		InterchangeControlVersionNumber:   "00401",
		InterchangeControlNumber:          1,
		AcknowledgementRequested:          1,
		UsageIndicator:                    "T", // T for test, P for production
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

	for index, shipment := range shipments {
		shipment, err := models.FetchShipmentForInvoice(db, shipment.ID)
		if err != nil {
			return transaction, err
		}
		shipment858c, err := generate858CShipment(shipment, index+1)
		if err != nil {
			return transaction, err
		}
		transaction += shipment858c
	}

	ge := edisegment.GE{
		NumberOfTransactionSetsIncluded: len(shipments),
		GroupControlNumber:              1,
	}
	iea := edisegment.IEA{
		NumberOfIncludedFunctionalGroups: 1,
		InterchangeControlNumber:         1,
	}

	transaction += (ge.String(delimiter) + iea.String(delimiter))

	return transaction, nil
}

func generate858CShipment(shipment models.Shipment, sequenceNum int) (string, error) {
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
	if shipment.Move.Orders.ServiceMember.LastName != nil {
		name = *shipment.Move.Orders.ServiceMember.LastName
	}
	if shipment.PickupAddress == nil {
		return "", errors.New("Shipment is missing pick up address")
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
		return "", errors.New("Orders is missing orders number")
	}
	tac := orders.TAC
	if tac == nil {
		return "", errors.New("Orders is missing TAC")
	}
	affiliation := orders.ServiceMember.Affiliation
	if orders.ServiceMember.Affiliation == nil {
		return "", errors.New("Service member is missing affiliation")
	}

	transactionNumber := fmt.Sprintf("%04d", sequenceNum)

	segments := []edisegment.Segment{
		&edisegment.ST{
			TransactionSetIdentifierCode: "858",
			TransactionSetControlNumber:  transactionNumber,
		},
		&edisegment.BX{
			TransactionSetPurposeCode:    "00",        // Original
			TransactionMethodTypeCode:    "J",         // Motor
			ShipmentMethodOfPayment:      "PP",        // Prepaid by seller
			ShipmentIdentificationNumber: "TODO:GBL",  // GBL
			StandardCarrierAlphaCode:     "TODO:SCAC", // tsp.StandardCarrierAlphaCode,
			ShipmentQualifier:            "4",         // HHG Bill of Lading
		},
		&edisegment.N9{
			ReferenceIdentificationQualifier: "DY", // DoD transportation service code #
			ReferenceIdentification:          "SC", // Shipment & cost information
		},
		&edisegment.N9{
			ReferenceIdentificationQualifier: "CN", // Invoice number
			ReferenceIdentification:          "TODO:InvoiceNumber",
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
			LocationQualifier:   "TODO", // CY (county), IP (postal), or RA (rate area)",
			LocationIdentifier:  "TODO",
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
			Weight:          100.11, // TODO: weight
			WeightQualifier: "B",    // Billing weight
			WeightUnitCode:  "L",    // Pounds
		},
	}

	// Add line items and linehaul

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
