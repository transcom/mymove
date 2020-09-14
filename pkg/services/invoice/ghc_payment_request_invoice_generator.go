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

const dateFormat = "20060102"
const timeFormat = "1504"

// Generate method takes a payment request and returns an Invoice858C
func (g GHCPaymentRequestInvoiceGenerator) Generate(paymentRequest models.PaymentRequest, sendProductionInvoice bool) (ediinvoice.Invoice858C, error) {
	// TODO: seems ReferenceID is a *string but cannot be saved as nil, do we need to validate it's not nil here
	var moveTaskOrder models.Move
	if paymentRequest.MoveTaskOrder.ID == uuid.Nil {
		// load mto
		err := g.DB.Q().
			Where("id = ?", paymentRequest.MoveTaskOrderID).
			First(&moveTaskOrder)
		if err != nil {
			return ediinvoice.Invoice858C{}, fmt.Errorf("cannot load MTO %s for PaymentRequest %s: %w", paymentRequest.MoveTaskOrderID, paymentRequest.ID, err)
		}
	} else {
		moveTaskOrder = paymentRequest.MoveTaskOrder
	}

	// check or load orders
	if moveTaskOrder.Orders.ID == uuid.Nil {
		err := g.DB.
			Load(&moveTaskOrder, "Orders")
		if err != nil {
			return ediinvoice.Invoice858C{}, fmt.Errorf("cannot load Orders %s for PaymentRequest %s: %w", moveTaskOrder.OrdersID, paymentRequest.ID, err)
		}
	}

	// check or load service member
	if moveTaskOrder.Orders.ServiceMember.ID == uuid.Nil {
		err := g.DB.
			Load(&moveTaskOrder.Orders, "ServiceMember")
		if err != nil {
			return ediinvoice.Invoice858C{}, fmt.Errorf("cannot load ServiceMember %s for PaymentRequest %s: %w", moveTaskOrder.Orders.ServiceMemberID, paymentRequest.ID, err)
		}
	}

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

	edi858.ST = edisegment.ST{
		TransactionSetIdentifierCode: "858",
		TransactionSetControlNumber:  "0001",
	}

	bx := edisegment.BX{
		TransactionSetPurposeCode:    "00",
		TransactionMethodTypeCode:    "J",
		ShipmentMethodOfPayment:      "PP",
		ShipmentIdentificationNumber: *moveTaskOrder.ReferenceID,
		StandardCarrierAlphaCode:     "TRUS",
		ShipmentQualifier:            "4",
	}

	edi858.Header = append(edi858.Header, &bx)

	paymentRequestNumberSegment := edisegment.N9{
		ReferenceIdentificationQualifier: "CN",
		ReferenceIdentification:          paymentRequest.PaymentRequestNumber,
	}
	edi858.Header = append(edi858.Header, &paymentRequestNumberSegment)

	// contract code to header
	var contractCodeServiceItemParam models.PaymentServiceItemParam
	err := g.DB.Q().
		Join("service_item_param_keys sipk", "payment_service_item_params.service_item_param_key_id = sipk.id").
		Join("payment_service_items psi", "payment_service_item_params.payment_service_item_id = psi.id").
		Join("payment_requests pr", "psi.payment_request_id = pr.id").
		Where("pr.id = ?", paymentRequest.ID).
		Where("sipk.key = ?", models.ServiceItemParamNameContractCode).
		First(&contractCodeServiceItemParam)
	if err != nil {
		return ediinvoice.Invoice858C{}, err
	}

	contractCodeSegment := edisegment.N9{
		ReferenceIdentificationQualifier: "CT",
		ReferenceIdentification:          contractCodeServiceItemParam.Value,
	}
	edi858.Header = append(edi858.Header, &contractCodeSegment)

	// Add service member details to header
	serviceMemberSegments, err := g.createServiceMemberDetailSegments(paymentRequest.ID, moveTaskOrder.Orders.ServiceMember)
	if err != nil {
		return ediinvoice.Invoice858C{}, err
	}
	edi858.Header = append(edi858.Header, serviceMemberSegments...)

	// Add requested pickup date
	var requestedPickupDateParam models.PaymentServiceItemParam
	err = g.DB.Q().
		Join("service_item_param_keys sipk", "payment_service_item_params.service_item_param_key_id = sipk.id").
		Join("payment_service_items psi", "payment_service_item_params.payment_service_item_id = psi.id").
		Join("payment_requests pr", "psi.payment_request_id = pr.id").
		Where("pr.id = ?", paymentRequest.ID).
		Where("sipk.key = ?", models.ServiceItemParamNameRequestedPickupDate).
		First(&requestedPickupDateParam)
	if err != nil {
		return ediinvoice.Invoice858C{}, err
	}

	requestedPickupDateSegment := edisegment.G62{
		DateQualifier: 86,
		Date:          requestedPickupDateParam.Value,
	}
	edi858.Header = append(edi858.Header, &requestedPickupDateSegment)

	// Add origin and destination details to header
	originDestinationSegments, err := g.createOriginAndDestinationSegments(paymentRequest.ID, moveTaskOrder.Orders)
	if err != nil {
		return ediinvoice.Invoice858C{}, err
	}
	edi858.Header = append(edi858.Header, originDestinationSegments...)

	var paymentServiceItems models.PaymentServiceItems
	err = g.DB.Q().
		Eager("MTOServiceItem.ReService").
		Where("payment_request_id = ?", paymentRequest.ID).
		All(&paymentServiceItems)
	if err != nil {
		return ediinvoice.Invoice858C{}, fmt.Errorf("Could not find payment service items: %w", err)
	}

	paymentServiceItemSegments, err := g.generatePaymentServiceItemSegments(paymentServiceItems)
	if err != nil {
		return ediinvoice.Invoice858C{}, fmt.Errorf("Could not generate payment service item segments: %w", err)
	}
	edi858.ServiceItems = append(edi858.ServiceItems, paymentServiceItemSegments...)

	// the total NumberOfIncludedSegments is ST + SE + all segments other than GS, GE, ISA, and IEA
	edi858.SE = edisegment.SE{
		NumberOfIncludedSegments:    2 + len(edi858.Header) + len(edi858.ServiceItems),
		TransactionSetControlNumber: "0001",
	}

	edi858.GE = edisegment.GE{
		NumberOfTransactionSetsIncluded: 1,
		GroupControlNumber:              100001251,
	}

	edi858.IEA = edisegment.IEA{
		NumberOfIncludedFunctionalGroups: 1,
		InterchangeControlNumber:         100001272,
	}

	return edi858, nil
}

func (g GHCPaymentRequestInvoiceGenerator) createServiceMemberDetailSegments(paymentRequestID uuid.UUID, serviceMember models.ServiceMember) ([]edisegment.Segment, error) {
	serviceMemberDetails := []edisegment.Segment{}

	// name
	serviceMemberName := edisegment.N9{
		ReferenceIdentificationQualifier: "1W",
		ReferenceIdentification:          serviceMember.ReverseNameLineFormat(),
	}
	serviceMemberDetails = append(serviceMemberDetails, &serviceMemberName)

	// rank
	rank := serviceMember.Rank
	if rank == nil {
		return []edisegment.Segment{}, fmt.Errorf("no rank found for ServiceMember ID: %s Payment Request ID: %s", serviceMember.ID, paymentRequestID)
	}
	serviceMemberRank := edisegment.N9{
		ReferenceIdentificationQualifier: "ML",
		ReferenceIdentification:          string(*rank),
	}
	serviceMemberDetails = append(serviceMemberDetails, &serviceMemberRank)

	// branch
	branch := serviceMember.Affiliation
	if branch == nil {
		return []edisegment.Segment{}, fmt.Errorf("no branch found for ServiceMember ID: %s Payment Request ID: %s", serviceMember.ID, paymentRequestID)
	}
	serviceMemberBranch := edisegment.N9{
		ReferenceIdentificationQualifier: "3L",
		ReferenceIdentification:          string(*branch),
	}
	serviceMemberDetails = append(serviceMemberDetails, &serviceMemberBranch)

	return serviceMemberDetails, nil
}

func (g GHCPaymentRequestInvoiceGenerator) createOriginAndDestinationSegments(paymentRequestID uuid.UUID, orders models.Order) ([]edisegment.Segment, error) {
	originAndDestinationSegments := []edisegment.Segment{}

	var err error
	var destinationDutyStation models.DutyStation
	if orders.NewDutyStation.ID == uuid.Nil {
		destinationDutyStation, err = models.FetchDutyStation(g.DB, orders.NewDutyStationID)
		if err != nil {
			return []edisegment.Segment{}, fmt.Errorf("cannot load NewDutyStation %s for PaymentRequest %s: %w", orders.NewDutyStationID, paymentRequestID, err)
		}
	} else {
		destinationDutyStation = orders.NewDutyStation
	}

	destTransportationOffice, err := models.FetchDutyStationTransportationOffice(g.DB, destinationDutyStation.ID)
	if err != nil {
		return []edisegment.Segment{}, fmt.Errorf("cannot load TransportationOffice for DutyStation %s for PaymentRequest %s: %w", orders.NewDutyStationID, paymentRequestID, err)
	}

	// destination name
	destinationStationName := orders.NewDutyStation.Name
	destinationName := edisegment.N1{
		EntityIdentifierCode:        "ST",
		Name:                        destinationStationName,
		IdentificationCodeQualifier: "10",
		IdentificationCode:          destTransportationOffice.Gbloc,
	}
	originAndDestinationSegments = append(originAndDestinationSegments, &destinationName)

	// destination address
	var destinationStreetAddress edisegment.N3
	if destinationDutyStation.Address.StreetAddress2 == nil {
		destinationStreetAddress = edisegment.N3{
			AddressInformation1: destinationDutyStation.Address.StreetAddress1,
		}
	} else {
		destinationStreetAddress = edisegment.N3{
			AddressInformation1: destinationDutyStation.Address.StreetAddress1,
			AddressInformation2: *destinationDutyStation.Address.StreetAddress2,
		}
	}
	originAndDestinationSegments = append(originAndDestinationSegments, &destinationStreetAddress)

	// destination city/state/postal
	destinationPostalDetails := edisegment.N4{
		CityName:            destinationDutyStation.Address.City,
		StateOrProvinceCode: destinationDutyStation.Address.State,
		PostalCode:          destinationDutyStation.Address.PostalCode,
		CountryCode:         string(*destinationDutyStation.Address.Country),
	}
	originAndDestinationSegments = append(originAndDestinationSegments, &destinationPostalDetails)

	// TODO: Create PER segment and implement Destination POC Phone

	// ========  ORIGIN ========= //
	// origin station name
	var originDutyStation models.DutyStation
	if orders.OriginDutyStationID != nil {
		originDutyStation, err = models.FetchDutyStation(g.DB, *orders.OriginDutyStationID)
		if err != nil {
			return []edisegment.Segment{}, fmt.Errorf("cannot load OriginDutyStation %s for PaymentRequest %s: %w", orders.OriginDutyStationID, paymentRequestID, err)
		}
	} else {
		originDutyStation = *orders.OriginDutyStation
	}

	originTransportationOffice, err := models.FetchDutyStationTransportationOffice(g.DB, originDutyStation.ID)
	if err != nil {
		return []edisegment.Segment{}, fmt.Errorf("cannot load TransportationOffice for DutyStation %s for PaymentRequest %s: %w", orders.OriginDutyStationID, paymentRequestID, err)
	}

	originName := edisegment.N1{
		EntityIdentifierCode:        "SF",
		Name:                        originDutyStation.Name,
		IdentificationCodeQualifier: "10",
		IdentificationCode:          originTransportationOffice.Gbloc,
	}
	originAndDestinationSegments = append(originAndDestinationSegments, &originName)

	// origin address
	var originStreetAddress edisegment.N3
	if originDutyStation.Address.StreetAddress2 == nil {
		originStreetAddress = edisegment.N3{
			AddressInformation1: originDutyStation.Address.StreetAddress1,
		}
	} else {
		originStreetAddress = edisegment.N3{
			AddressInformation1: originDutyStation.Address.StreetAddress1,
			AddressInformation2: *originDutyStation.Address.StreetAddress2,
		}
	}
	originAndDestinationSegments = append(originAndDestinationSegments, &originStreetAddress)

	// origin city/state/postal
	originPostalDetails := edisegment.N4{
		CityName:            originDutyStation.Address.City,
		StateOrProvinceCode: originDutyStation.Address.State,
		PostalCode:          originDutyStation.Address.PostalCode,
		CountryCode:         string(*originDutyStation.Address.Country),
	}
	originAndDestinationSegments = append(originAndDestinationSegments, &originPostalDetails)

	// TODO: Create PER segment and implement Origin POC Phone

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
