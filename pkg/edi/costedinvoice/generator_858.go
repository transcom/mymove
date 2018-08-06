package edicostedinvoice

import (
	"errors"
	"fmt"
	"time"

	"github.com/gobuffalo/pop"

	"github.com/transcom/mymove/pkg/edi/segment"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/models"
)

const delimiter = "~"
const dateFormat = "20060102"
const timeFormat = "1504"

// Generate858C generates an EDI X12 858C transaction set
func Generate858C(shipments []models.Shipment, db *pop.Connection) (string, error) {
	currentTime := time.Now()
	isa := edisegment.ISA{
		AuthorizationInformationQualifier: "00",
		AuthorizationInformation:          "          ",
		SecurityInformationQualifier:      "00",
		SecurityInformation:               "          ",
		InterchangeSenderIDQualifier:      "ZZ",
		InterchangeSenderID:               fmt.Sprintf("%-15v", "W28GPR-DPS"), // TODO: update with our own after talking to US Bank. Must be 15 characters
		InterchangeReceiverIDQualifier:    "12",
		InterchangeReceiverID:             fmt.Sprintf("%-15v", "8004171844"), // Syncada - extra blank spaces are intentional
		InterchangeDate:                   currentTime.Format(dateFormat),
		InterchangeTime:                   currentTime.Format(timeFormat),
		InterchangeControlStandards:       "U",
		InterchangeControlVersionNumber:   "00401",
		InterchangeControlNumber:          "000000001",
		AcknowledgementRequested:          "1",
		UsageIndicator:                    "T", // T for test, P for production
		ComponentElementSeparator:         "|",
	}
	transaction := ""
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
			ReferenceIdentificationQualifier: "PQ", // Payee code
			ReferenceIdentification:          "TODO:SupplierID",
		},
		&edisegment.N9{
			ReferenceIdentificationQualifier: "OQ", // Order number
			ReferenceIdentification:          "TODO:OrdersNumber",
			FreeFormDescription:              "TODO:BranchOfService",
			Date:                             "TODO:OrdersDate",
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
		// Origin installation
		&edisegment.N1{
			EntityIdentifierCode: "RG", // Ship From
			Name:                 "TODO:GBLOC",
			IdentificationCodeQualifier: "27", // GBLOC
			IdentificationCode:          "TODO:GBLOC",
		},
		&edisegment.N4{
			LocationQualifier:  "RA", // Rate area
			LocationIdentifier: "TODO:RA",
		},
		// Destination installation
		&edisegment.N1{
			EntityIdentifierCode: "RH", // Ship From
			Name:                 "TODO:GBLOC",
			IdentificationCodeQualifier: "27", // GBLOC
			IdentificationCode:          "TODO:GBLOC",
		},
		&edisegment.N4{
			LocationQualifier:  "RA", // Rate area
			LocationIdentifier: "TODO:RA",
		},
		// Accounting info
		&edisegment.FA1{
			AgencyQualifierCode: edisegment.AffiliationToAgency[internalmessages.AffiliationAIRFORCE], // TODO: get correct agency
		},
		&edisegment.FA2{
			BreakdownStructureDetailCode: "TA",   // TAC
			FinancialInformationCode:     "NAL8", // TODO: sample TAC
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
