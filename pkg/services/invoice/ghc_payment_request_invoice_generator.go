package invoice

import (
	"fmt"
	"strconv"
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"

	ediinvoice "github.com/transcom/mymove/pkg/edi/invoice"
	edisegment "github.com/transcom/mymove/pkg/edi/segment"
)

// GHCPaymentRequestInvoiceGenerator is a service object to turn payment requests into 858s
type GHCPaymentRequestInvoiceGenerator struct {
	DB *pop.Connection
}

const dateFormat = "060102"
const timeFormat = "1504"

// Generate method takes a payment request and returns an Invoice858C
func (g GHCPaymentRequestInvoiceGenerator) Generate(paymentRequest models.PaymentRequest, sendProductionInvoice bool) (ediinvoice.Invoice858C, error) {
	// TODO: probably need to check if the MTO is loaded on the paymentRequest that is passed in, not sure what is more in line with go standards to error out if it's not there or look it up.
	// TODO: seems ReferenceID is a *string but cannot be saved as nil, do we need to validate it's not nil here

	currentTime := time.Now()

	// TODO: interchangeControlNumber, err := icnSequencer.NextVal()
	// if err != nil {
	// 	return Invoice858C{}, errors.Wrap(err, fmt.Sprintf("Failed to get next Interchange Control Number"))
	// }
	var usageIndicator string
	if sendProductionInvoice {
		usageIndicator = "P"
	} else {
		usageIndicator = "T"
	}

	var edi858 ediinvoice.Invoice858C
	edi858.ISA = edisegment.ISA{
		AuthorizationInformationQualifier: "00", // No authorization information
		AuthorizationInformation:          "0084182369",
		SecurityInformationQualifier:      "00", // No security information
		SecurityInformation:               "_   _",
		InterchangeSenderIDQualifier:      "ZZ",
		InterchangeSenderID:               "GOVDPIBS",
		InterchangeReceiverIDQualifier:    "12",
		InterchangeReceiverID:             "8004171844",
		InterchangeDate:                   currentTime.Format(dateFormat),
		InterchangeTime:                   currentTime.Format(timeFormat),
		InterchangeControlStandards:       "U",
		InterchangeControlVersionNumber:   "00401",
		InterchangeControlNumber:          100001272,
		AcknowledgementRequested:          0,
		UsageIndicator:                    usageIndicator, // T for test, P for production
		ComponentElementSeparator:         "|",
	}
	edi858.IEA = edisegment.IEA{
		NumberOfIncludedFunctionalGroups: 1,
		InterchangeControlNumber:         100001272,
	}

	edi858.GS = edisegment.GS{
		FunctionalIdentifierCode: "SI",
		ApplicationSendersCode:   "MYMOVE",
		ApplicationReceiversCode: "8004171844",
		Date:                     currentTime.Format(dateFormat),
		Time:                     currentTime.Format(timeFormat),
		GroupControlNumber:       100001251,
		ResponsibleAgencyCode:    "X",
		Version:                  "004010",
	}

	edi858.GE = edisegment.GE{
		NumberOfTransactionSetsIncluded: 1,
		GroupControlNumber:              100001251,
	}

	bx := edisegment.BX{
		TransactionSetPurposeCode:    "00",
		TransactionMethodTypeCode:    "J",
		ShipmentMethodOfPayment:      "PP",
		ShipmentIdentificationNumber: *paymentRequest.MoveTaskOrder.ReferenceID,
		StandardCarrierAlphaCode:     "TRUS",
		ShipmentQualifier:            "4",
	}

	edi858.Header = append(edi858.Header, &bx)

	paymentRequestNumberSegment := edisegment.N9{
		ReferenceIdentificationQualifier: "CN",
		ReferenceIdentification:          paymentRequest.PaymentRequestNumber,
	}
	edi858.Header = append(edi858.Header, &paymentRequestNumberSegment)

	// Add service member details to header
	serviceMemberSegments, err := g.createServiceMemberDetailSegments(paymentRequest)
	if err != nil {
		return ediinvoice.Invoice858C{}, err
	}
	edi858.Header = append(edi858.Header, serviceMemberSegments...)

	// TODO: Determine correct values to fill in the l7 segment
	// l7 := edisegment.L7{
	// 	LadingLineItemNumber:    812,
	// 	TariffNumber:    "T",
	// 	TariffItemNumber:      "",
	// 	TariffDistance: ,
	// }
	// edi858.Header = append(edi858.Header, &l7)

	// Add NTE lines to header
	nteSegment, err := g.createNteLines()
	if err != nil {
		return ediinvoice.Invoice858C{}, err
	}
	edi858.Header = append(edi858.Header, nteSegment...)

	// Add origin and destination details to header
	originDestinationSegments, err := g.createOriginAndDestinationSegments(paymentRequest)
	if err != nil {
		return ediinvoice.Invoice858C{}, err
	}
	edi858.Header = append(edi858.Header, originDestinationSegments...)

	var paymentServiceItems models.PaymentServiceItems
	error := g.DB.Q().
		Eager("MTOServiceItem.ReService").
		Where("payment_request_id = ?", paymentRequest.ID).
		All(&paymentServiceItems)
	if error != nil {
		return ediinvoice.Invoice858C{}, fmt.Errorf("Could not find payment service items: %w", error)
	}

	paymentServiceItemSegments, err := g.generatePaymentServiceItemSegments(paymentServiceItems)
	if err != nil {
		return ediinvoice.Invoice858C{}, fmt.Errorf("Could not generate payment service item segments: %w", err)
	}
	edi858.ServiceItems = append(edi858.ServiceItems, paymentServiceItemSegments...)

	return edi858, nil
}

func (g GHCPaymentRequestInvoiceGenerator) createServiceMemberDetailSegments(paymentRequest models.PaymentRequest) ([]edisegment.Segment, error) {
	serviceMemberDetails := []edisegment.Segment{}

	serviceMember := paymentRequest.MoveTaskOrder.Orders.ServiceMember
	if serviceMember.ID == uuid.Nil {
		return []edisegment.Segment{}, fmt.Errorf("no ServiceMember found for Payment Request ID: %s", paymentRequest.ID)
	}

	// name
	serviceMemberName := edisegment.N9{
		ReferenceIdentificationQualifier: "1W",
		ReferenceIdentification:          serviceMember.ReverseNameLineFormat(),
	}
	serviceMemberDetails = append(serviceMemberDetails, &serviceMemberName)

	// rank
	rank := serviceMember.Rank
	if rank == nil {
		return []edisegment.Segment{}, fmt.Errorf("no rank found for ServiceMember ID: %s Payment Request ID: %s", serviceMember.ID, paymentRequest.ID)
	}
	serviceMemberRank := edisegment.N9{
		ReferenceIdentificationQualifier: "ML",
		ReferenceIdentification:          string(*rank),
	}
	serviceMemberDetails = append(serviceMemberDetails, &serviceMemberRank)

	// branch
	branch := serviceMember.Affiliation
	if branch == nil {
		return []edisegment.Segment{}, fmt.Errorf("no branch found for ServiceMember ID: %s Payment Request ID: %s", serviceMember.ID, paymentRequest.ID)
	}
	serviceMemberBranch := edisegment.N9{
		ReferenceIdentificationQualifier: "3L",
		ReferenceIdentification:          string(*branch),
	}
	serviceMemberDetails = append(serviceMemberDetails, &serviceMemberBranch)

	return serviceMemberDetails, nil
}

func (g GHCPaymentRequestInvoiceGenerator) createNteLines() ([]edisegment.Segment, error) {
	nteSegments := []edisegment.Segment{}

	nteZip3 := edisegment.NTE{
		NoteReferenceCode: "ADD",
		Description:       "DistanceZip3",
	}

	nteZip5 := edisegment.NTE{
		NoteReferenceCode: "ADD",
		Description:       "DistanceZip5",
	}

	nteSitOrigin := edisegment.NTE{
		NoteReferenceCode: "ADD",
		Description:       "DistanceZip5SITOrigin",
	}

	nteSitDest := edisegment.NTE{
		NoteReferenceCode: "ADD",
		Description:       "DistanceZip5SITDest",
	}

	nteSegments = append(nteSegments, &nteZip3, &nteZip5, &nteSitOrigin, &nteSitDest)

	return nteSegments, nil

}

func (g GHCPaymentRequestInvoiceGenerator) createOriginAndDestinationSegments(paymentRequest models.PaymentRequest) ([]edisegment.Segment, error) {
	originAndDestinationSegments := []edisegment.Segment{}

	order := paymentRequest.MoveTaskOrder.Orders
	mtoShipment := paymentRequest.MoveTaskOrder.MTOShipments[0]

	if order.ID == uuid.Nil {
		return []edisegment.Segment{}, fmt.Errorf("no order found for Payment Request ID: %s", paymentRequest.ID)
	}

	if mtoShipment.ID == uuid.Nil {
		return []edisegment.Segment{}, fmt.Errorf("no MTO shipment found for Payment Request ID: %s", paymentRequest.ID)
	}

	// destination name
	destinationStationName := order.NewDutyStation.Name
	if len(destinationStationName) == 0 {
		return []edisegment.Segment{}, fmt.Errorf("no destination duty station name found for Order ID: %s Payment Request ID: %s", order.ID, paymentRequest.ID)
	}
	destinationName := edisegment.N1{
		EntityIdentifierCode:        "ST",
		Name:                        destinationStationName,
		IdentificationCodeQualifier: "10",
		IdentificationCode:          "GBLOC/DODAAC",
	}
	originAndDestinationSegments = append(originAndDestinationSegments, &destinationName)

	// destination address
	destinationAddress := mtoShipment.DestinationAddress.StreetAddress1

	if len(destinationAddress) == 0 {
		return []edisegment.Segment{}, fmt.Errorf("no destination street address found for MTO shipment ID: %s Payment Request ID: %s", mtoShipment.ID, paymentRequest.ID)
	}

	destinationStreetAddress := edisegment.N3{
		AddressInformation1: destinationAddress,
	}
	originAndDestinationSegments = append(originAndDestinationSegments, &destinationStreetAddress)

	// destination city/state/postal
	destinationAddressDetails := mtoShipment.DestinationAddress

	if destinationAddressDetails.ID == uuid.Nil {
		return []edisegment.Segment{}, fmt.Errorf("no destination address found for MTO shipment ID: %s Payment Request ID: %s", mtoShipment.ID, paymentRequest.ID)
	}

	destinationPostalDetails := edisegment.N4{
		CityName:            destinationAddressDetails.City,
		StateOrProvinceCode: destinationAddressDetails.State,
		PostalCode:          destinationAddressDetails.PostalCode,
		CountryCode:         string(*destinationAddressDetails.Country),
		LocationQualifier:   "SL",
		LocationIdentifier:  "237740290",
	}
	originAndDestinationSegments = append(originAndDestinationSegments, &destinationPostalDetails)

	// TODO: Create PER segment and implement Destination POC Phone

	// ========  ORIGIN ========= //
	// origin station name
	originStationName := order.OriginDutyStation.Name
	if len(originStationName) == 0 {
		return []edisegment.Segment{}, fmt.Errorf("no origin duty station name found for Order ID: %s Payment Request ID: %s", order.ID, paymentRequest.ID)
	}
	originName := edisegment.N1{
		EntityIdentifierCode:        "ST",
		Name:                        originStationName,
		IdentificationCodeQualifier: "10",
		IdentificationCode:          "GBLOC/DODAAC",
	}
	originAndDestinationSegments = append(originAndDestinationSegments, &originName)

	// origin address
	originAddress := mtoShipment.PickupAddress.StreetAddress1

	if len(originAddress) == 0 {
		return []edisegment.Segment{}, fmt.Errorf("no origin street address found for MTO shipment ID: %s Payment Request ID: %s", mtoShipment.ID, paymentRequest.ID)
	}

	originStreetAddress := edisegment.N3{
		AddressInformation1: originAddress,
	}
	originAndDestinationSegments = append(originAndDestinationSegments, &originStreetAddress)

	// origin city/state/postal
	originAddressDetails := mtoShipment.PickupAddress

	if originAddressDetails.ID == uuid.Nil {
		return []edisegment.Segment{}, fmt.Errorf("no origin address found for MTO shipment ID: %s Payment Request ID: %s", mtoShipment.ID, paymentRequest.ID)
	}

	originPostalDetails := edisegment.N4{
		CityName:            originAddressDetails.City,
		StateOrProvinceCode: originAddressDetails.State,
		PostalCode:          originAddressDetails.PostalCode,
		CountryCode:         string(*originAddressDetails.Country),
		LocationQualifier:   "SL",
		LocationIdentifier:  "237740290",
	}
	originAndDestinationSegments = append(originAndDestinationSegments, &originPostalDetails)

	return originAndDestinationSegments, nil
}

func (g GHCPaymentRequestInvoiceGenerator) fetchPaymentServiceItemParam(serviceItemID uuid.UUID, key models.ServiceItemParamName) (models.PaymentServiceItemParam, error) {
	var paymentServiceItemParam models.PaymentServiceItemParam

	err := g.DB.Q().
		Join("service_item_param_keys sk", "payment_service_item_params.service_item_param_key_id = sk.id").
		Where("payment_service_item_id = ?", serviceItemID).
		Where("sk.key = ?", key).
		First(&paymentServiceItemParam)
	if err != nil {
		return models.PaymentServiceItemParam{}, fmt.Errorf("Could not lookup PaymentServiceItemParam key (%s) payment service item id (%s): %w", key, serviceItemID, err)
	}
	return paymentServiceItemParam, nil
}

func (g GHCPaymentRequestInvoiceGenerator) generatePaymentServiceItemSegments(paymentServiceItems models.PaymentServiceItems) ([]edisegment.Segment, error) {
	//Initialize empty collection of segments
	var segments []edisegment.Segment

	// Iterate over payment service items
	for idx, serviceItem := range paymentServiceItems {
		hierarchicalIDNumber := idx + 1
		// Build and put together the segments
		hlSegment := edisegment.HL{
			HierarchicalIDNumber:  strconv.Itoa(hierarchicalIDNumber), // may need to change if sending multiple payment request in a single edi
			HierarchicalLevelCode: "|",
		}

		n9Segment := edisegment.N9{
			ReferenceIdentificationQualifier: "PO",
			// pending creation of shorter identifier for payment service item
			// https://dp3.atlassian.net/browse/MB-3718
			ReferenceIdentification: serviceItem.ID.String(),
		}

		// TODO: add another n9 for SIT
		// TODO: add a L5 segment/definition

		var weight models.PaymentServiceItemParam
		// TODO: update to have a case statement as different service items may or may not have weight
		// and the distance key can differ (zip3 v zip5, and distances for SIT)
		weight, err := g.fetchPaymentServiceItemParam(serviceItem.ID, models.ServiceItemParamNameWeightBilledActual)
		if err != nil {
			return nil, err
		}
		weightFloat, err := strconv.ParseFloat(weight.Value, 64)
		if err != nil {
			return nil, fmt.Errorf("Could not parse weight for PaymentServiceItem %s: %w", serviceItem.ID, err)
		}
		distance, err := g.fetchPaymentServiceItemParam(serviceItem.ID, models.ServiceItemParamNameDistanceZip3)
		if err != nil {
			return nil, err
		}
		distanceFloat, err := strconv.ParseFloat(distance.Value, 64)
		if err != nil {
			return nil, fmt.Errorf("Could not parse Distance Zip3 for PaymentServiceItem %s: %w", serviceItem.ID, err)
		}

		l0Segment := edisegment.L0{
			LadingLineItemNumber:   hierarchicalIDNumber,
			BilledRatedAsQuantity:  distanceFloat,
			BilledRatedAsQualifier: "DM",
			Weight:                 weightFloat,
			WeightQualifier:        "B",
			WeightUnitCode:         "L",
		}

		segments = append(segments, &hlSegment, &n9Segment, &l0Segment)
	}

	return segments, nil
}
