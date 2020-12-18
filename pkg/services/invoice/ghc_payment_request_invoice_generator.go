package invoice

import (
	"fmt"
	"strconv"
	"time"

	"github.com/emirpasic/gods/maps/linkedhashmap"

	"github.com/gobuffalo/pop/v5"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/db/sequence"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"

	ediinvoice "github.com/transcom/mymove/pkg/edi/invoice"
	edisegment "github.com/transcom/mymove/pkg/edi/segment"
)

type ghcPaymentRequestInvoiceGenerator struct {
	db           *pop.Connection
	icnSequencer sequence.Sequencer
}

// NewGHCPaymentRequestInvoiceGenerator returns an implementation of the GHCPaymentRequestInvoiceGenerator interface
func NewGHCPaymentRequestInvoiceGenerator(db *pop.Connection, icnSequencer sequence.Sequencer) services.GHCPaymentRequestInvoiceGenerator {
	return &ghcPaymentRequestInvoiceGenerator{
		db:           db,
		icnSequencer: icnSequencer,
	}
}

const dateFormat = "20060102"
const isaDateFormat = "060102"
const timeFormat = "1504"

// Generate method takes a payment request and returns an Invoice858C
func (g ghcPaymentRequestInvoiceGenerator) Generate(paymentRequest models.PaymentRequest, sendProductionInvoice bool) (ediinvoice.Invoice858C, error) {
	var moveTaskOrder models.Move
	if paymentRequest.MoveTaskOrder.ID == uuid.Nil {
		// load mto
		err := g.db.Q().
			Where("id = ?", paymentRequest.MoveTaskOrderID).
			First(&moveTaskOrder)
		if err != nil {
			if err.Error() == models.RecordNotFoundErrorString {
				return ediinvoice.Invoice858C{}, services.NewNotFoundError(paymentRequest.MoveTaskOrder.ID, "for MoveTaskOrder")
			}
			return ediinvoice.Invoice858C{}, services.NewQueryError("MoveTaskOrder", err, "Unexpected error")
		}
	} else {
		moveTaskOrder = paymentRequest.MoveTaskOrder
	}

	// check or load orders
	if moveTaskOrder.ReferenceID == nil {
		return ediinvoice.Invoice858C{}, services.NewConflictError(moveTaskOrder.ID, "Invalid move taskorder. Must have a ReferenceID value")
	}

	if moveTaskOrder.Orders.ID == uuid.Nil {
		err := g.db.
			Load(&moveTaskOrder, "Orders")
		if err != nil {
			if err.Error() == models.RecordNotFoundErrorString {
				return ediinvoice.Invoice858C{}, services.NewNotFoundError(moveTaskOrder.Orders.ID, "for Orders")
			}
			return ediinvoice.Invoice858C{}, services.NewQueryError("Orders", err, "Unexpected error")
		}
	}

	// check or load service member
	if moveTaskOrder.Orders.ServiceMember.ID == uuid.Nil {
		err := g.db.
			Load(&moveTaskOrder.Orders, "ServiceMember")

		if err != nil {
			if err.Error() == models.RecordNotFoundErrorString {
				return ediinvoice.Invoice858C{}, services.NewNotFoundError(moveTaskOrder.Orders.ServiceMemberID, "for ServiceMember")
			}
			return ediinvoice.Invoice858C{}, services.NewQueryError("ServiceMember", err, fmt.Sprintf("cannot load ServiceMember %s for PaymentRequest %s: %s", moveTaskOrder.Orders.ServiceMemberID, paymentRequest.ID, err))
		}
	}

	currentTime := time.Now()

	interchangeControlNumber, err := g.icnSequencer.NextVal()
	if err != nil {
		return ediinvoice.Invoice858C{}, fmt.Errorf("Failed to get next Interchange Control Number: %w", err)
	}

	var usageIndicator string
	if sendProductionInvoice {
		usageIndicator = "P"
	} else {
		usageIndicator = "T"
	}

	var edi858 ediinvoice.Invoice858C
	edi858.Header = linkedhashmap.New()
	edi858.ISA = edisegment.ISA{
		AuthorizationInformationQualifier: "00", // No authorization information
		AuthorizationInformation:          "0084182369",
		SecurityInformationQualifier:      "00", // No security information
		SecurityInformation:               "0000000000",
		InterchangeSenderIDQualifier:      "ZZ",
		InterchangeSenderID:               fmt.Sprintf("%-15s", "MILMOVE"),
		InterchangeReceiverIDQualifier:    "12",
		InterchangeReceiverID:             fmt.Sprintf("%-15s", "8004171844"),
		InterchangeDate:                   currentTime.Format(isaDateFormat),
		InterchangeTime:                   currentTime.Format(timeFormat),
		InterchangeControlStandards:       "U",
		InterchangeControlVersionNumber:   "00401",
		InterchangeControlNumber:          interchangeControlNumber,
		AcknowledgementRequested:          0,
		UsageIndicator:                    usageIndicator, // T for test, P for production
		ComponentElementSeparator:         "|",
	}

	edi858.GS = edisegment.GS{
		FunctionalIdentifierCode: "SI",
		ApplicationSendersCode:   "MILMOVE",
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

	edi858.Header.Put("BX_ShipmentInformation", &bx)

	paymentRequestNumberSegment := edisegment.N9{
		ReferenceIdentificationQualifier: "CN",
		ReferenceIdentification:          paymentRequest.PaymentRequestNumber,
	}
	edi858.Header.Put("N9_PaymentRequestNumber", &paymentRequestNumberSegment)

	// contract code to header
	var contractCodeServiceItemParam models.PaymentServiceItemParam
	err = g.db.Q().
		Join("service_item_param_keys sipk", "payment_service_item_params.service_item_param_key_id = sipk.id").
		Join("payment_service_items psi", "payment_service_item_params.payment_service_item_id = psi.id").
		Join("payment_requests pr", "psi.payment_request_id = pr.id").
		Where("pr.id = ?", paymentRequest.ID).
		Where("sipk.key = ?", models.ServiceItemParamNameContractCode).
		First(&contractCodeServiceItemParam)
	if err != nil {
		if err.Error() == models.RecordNotFoundErrorString {
			return ediinvoice.Invoice858C{}, services.NewNotFoundError(contractCodeServiceItemParam.ID, "for ContractCode")
		}
		return ediinvoice.Invoice858C{}, services.NewQueryError("ContractCode", err, fmt.Sprintf("Couldn't find contract code: %s", err))
	}

	contractCodeSegment := edisegment.N9{
		ReferenceIdentificationQualifier: "CT",
		ReferenceIdentification:          contractCodeServiceItemParam.Value,
	}
	edi858.Header.Put("N9_ContractCode", &contractCodeSegment)

	// Add service member details to header
	err = g.createServiceMemberDetailSegments(paymentRequest.ID, moveTaskOrder.Orders.ServiceMember, edi858.Header)
	if err != nil {
		return ediinvoice.Invoice858C{}, err
	}

	var paymentServiceItems models.PaymentServiceItems
	err = g.db.Q().
		Eager("MTOServiceItem.ReService").
		Where("payment_request_id = ?", paymentRequest.ID).
		All(&paymentServiceItems)
	if err != nil {
		if err.Error() == models.RecordNotFoundErrorString {
			return ediinvoice.Invoice858C{}, services.NewNotFoundError(paymentRequest.ID, "for paayment service items in PaymentRequest")
		}
		return ediinvoice.Invoice858C{}, services.NewQueryError("PaymentServiceItems", err, fmt.Sprintf("Could not find payment service items: %s", err))
	}

	if !msOrCsOnly(paymentServiceItems) {
		_, err = g.createG62Segments(paymentRequest.ID, edi858.Header)
		if err != nil {
			return ediinvoice.Invoice858C{}, err
		}
	}

	// Add buyer and seller organization names
	err = g.createBuyerAndSellerOrganizationNamesSegments(paymentRequest.ID, moveTaskOrder.Orders, edi858.Header)
	if err != nil {
		return ediinvoice.Invoice858C{}, err
	}

	// Add origin and destination details to header
	err = g.createOriginAndDestinationSegments(paymentRequest.ID, moveTaskOrder.Orders, edi858.Header)
	if err != nil {
		return ediinvoice.Invoice858C{}, err
	}

	paymentServiceItemSegments, err := g.generatePaymentServiceItemSegments(paymentServiceItems, moveTaskOrder.Orders)
	if err != nil {
		return ediinvoice.Invoice858C{}, err
	}
	edi858.ServiceItems = append(edi858.ServiceItems, paymentServiceItemSegments...)

	// the total NumberOfIncludedSegments is ST + SE + all segments other than GS, GE, ISA, and IEA
	edi858.SE = edisegment.SE{
		NumberOfIncludedSegments:    2 + edi858.Header.Size() + len(edi858.ServiceItems),
		TransactionSetControlNumber: "0001",
	}

	edi858.GE = edisegment.GE{
		NumberOfTransactionSetsIncluded: 1,
		GroupControlNumber:              100001251,
	}

	edi858.IEA = edisegment.IEA{
		NumberOfIncludedFunctionalGroups: 1,
		InterchangeControlNumber:         interchangeControlNumber,
	}

	return edi858, nil
}

func (g ghcPaymentRequestInvoiceGenerator) createServiceMemberDetailSegments(paymentRequestID uuid.UUID, serviceMember models.ServiceMember, header *linkedhashmap.Map) error {

	// name
	serviceMemberName := edisegment.N9{
		ReferenceIdentificationQualifier: "1W",
		ReferenceIdentification:          serviceMember.ReverseNameLineFormat(),
	}
	header.Put("N9_ServiceMemberName", &serviceMemberName)

	// rank
	rank := serviceMember.Rank
	if rank == nil {
		return services.NewBadDataError(fmt.Sprintf("no rank found for ServiceMember ID: %s Payment Request ID: %s", serviceMember.ID, paymentRequestID))
	}
	serviceMemberRank := edisegment.N9{
		ReferenceIdentificationQualifier: "ML",
		ReferenceIdentification:          string(*rank),
	}
	header.Put("N9_ServiceMemberRank", &serviceMemberRank)

	// branch
	branch := serviceMember.Affiliation
	if branch == nil {
		return services.NewBadDataError(fmt.Sprintf("no branch found for ServiceMember ID: %s Payment Request ID: %s", serviceMember.ID, paymentRequestID))
	}
	serviceMemberBranch := edisegment.N9{
		ReferenceIdentificationQualifier: "3L",
		ReferenceIdentification:          string(*branch),
	}
	header.Put("N9_ServiceMemberBranch", &serviceMemberBranch)

	return nil
}

func (g ghcPaymentRequestInvoiceGenerator) createG62Segments(paymentRequestID uuid.UUID, header *linkedhashmap.Map) ([]edisegment.Segment, error) {
	var g62Segments []edisegment.Segment

	// Get all the shipments associated with this payment request's service items, ordered by shipment creation date.
	var shipments models.MTOShipments
	err := g.db.Q().
		Join("mto_service_items msi", "mto_shipments.id = msi.mto_shipment_id").
		Join("payment_service_items psi", "msi.id = psi.mto_service_item_id").
		Where("psi.payment_request_id = ?", paymentRequestID).
		Order("msi.created_at").
		All(&shipments)
	if err != nil {
		if err.Error() == models.RecordNotFoundErrorString {
			return nil, services.NewNotFoundError(paymentRequestID, "for mto shipments associated with PaymentRequest")
		}
		return nil, services.NewQueryError("MTOShipments", err, fmt.Sprintf("error querying for shipments to use in G62 segments in PaymentRequest %s: %s", paymentRequestID, err))
	}

	// If no shipments, then just return because we will not have access to the dates.
	if len(shipments) == 0 {
		return g62Segments, nil
	}

	// Use the first (earliest) shipment.
	shipment := shipments[0]

	// Insert request pickup date, if available.
	if shipment.RequestedPickupDate != nil {
		requestedPickupDateSegment := edisegment.G62{
			DateQualifier: 10,
			Date:          shipment.RequestedPickupDate.Format(dateFormat),
		}
		header.Put("G62_RequestedPickupDate", &requestedPickupDateSegment)
	}

	// Insert expected pickup date, if available.
	if shipment.ScheduledPickupDate != nil {
		scheduledPickupDateSegment := edisegment.G62{
			DateQualifier: 76,
			Date:          shipment.ScheduledPickupDate.Format(dateFormat),
		}
		header.Put("G62_ScheduledPickupDate", &scheduledPickupDateSegment)
	}

	// Insert expected pickup date, if available.
	if shipment.ActualPickupDate != nil {
		actualPickupDateSegment := edisegment.G62{
			DateQualifier: 86,
			Date:          shipment.ActualPickupDate.Format(dateFormat),
		}
		header.Put("G62_ActualPickupDate", &actualPickupDateSegment)
	}

	return g62Segments, nil
}

func (g ghcPaymentRequestInvoiceGenerator) createBuyerAndSellerOrganizationNamesSegments(paymentRequestID uuid.UUID, orders models.Order, header *linkedhashmap.Map) error {

	var err error
	var originDutyStation models.DutyStation

	if orders.OriginDutyStationID != nil && *orders.OriginDutyStationID != uuid.Nil {
		originDutyStation, err = models.FetchDutyStation(g.db, *orders.OriginDutyStationID)
		if err != nil {
			return services.NewInvalidInputError(*orders.OriginDutyStationID, err, nil, "unable to find origin duty station")
		}
	} else {
		return services.NewBadDataError("Invalid Order, must have OriginDutyStation")
	}

	originTransportationOffice, err := models.FetchDutyStationTransportationOffice(g.db, originDutyStation.ID)
	if err != nil {
		return services.NewInvalidInputError(originDutyStation.ID, err, nil, "unable to find origin duty station")
	}

	// buyer organization name
	buyerOrganizationName := edisegment.N1{
		EntityIdentifierCode:        "BY",
		Name:                        originTransportationOffice.Name,
		IdentificationCodeQualifier: "92",
		IdentificationCode:          originTransportationOffice.Gbloc,
	}

	// seller organization name
	sellerOrganizationName := edisegment.N1{
		EntityIdentifierCode:        "SE",
		Name:                        "Prime",
		IdentificationCodeQualifier: "2",
		IdentificationCode:          "PRME",
	}

	header.Put("N1_BuyerOrganizationName", &buyerOrganizationName)
	header.Put("N1_SellerOrganizationName", &sellerOrganizationName)

	return nil
}

func (g ghcPaymentRequestInvoiceGenerator) createOriginAndDestinationSegments(paymentRequestID uuid.UUID, orders models.Order, header *linkedhashmap.Map) error {
	var err error
	var destinationDutyStation models.DutyStation
	if orders.NewDutyStationID != uuid.Nil {
		destinationDutyStation, err = models.FetchDutyStation(g.db, orders.NewDutyStationID)
		if err != nil {
			return services.NewInvalidInputError(orders.NewDutyStationID, err, nil, "unable to find new duty station")
		}
	} else {
		return services.NewBadDataError("Invalid Order, must have NewDutyStation")
	}

	destTransportationOffice, err := models.FetchDutyStationTransportationOffice(g.db, destinationDutyStation.ID)
	if err != nil {
		return services.NewInvalidInputError(destinationDutyStation.ID, err, nil, "unable to find destination duty station")
	}

	// destination name
	destinationName := edisegment.N1{
		EntityIdentifierCode:        "ST",
		Name:                        destinationDutyStation.Name,
		IdentificationCodeQualifier: "10",
		IdentificationCode:          destTransportationOffice.Gbloc,
	}
	header.Put("N1_DestinationName", &destinationName)

	// destination address
	if len(destinationDutyStation.Address.StreetAddress1) > 0 {
		destinationStreetAddress := edisegment.N3{
			AddressInformation1: destinationDutyStation.Address.StreetAddress1,
		}
		if destinationDutyStation.Address.StreetAddress2 != nil {
			destinationStreetAddress.AddressInformation2 = *destinationDutyStation.Address.StreetAddress2
		}
		header.Put("N3_DestinationStreetAddress", &destinationStreetAddress)
	}

	// destination city/state/postal
	destinationPostalDetails := edisegment.N4{
		CityName:            destinationDutyStation.Address.City,
		StateOrProvinceCode: destinationDutyStation.Address.State,
		PostalCode:          destinationDutyStation.Address.PostalCode,
	}
	if destinationDutyStation.Address.Country != nil {
		countryCode, ccErr := destinationDutyStation.Address.CountryCode()
		if ccErr != nil {
			return ccErr
		}
		destinationPostalDetails.CountryCode = string(*countryCode)
	}
	header.Put("N4_DestinationPostalDetails", &destinationPostalDetails)

	// Destination PER
	destinationStationPhoneLines := destTransportationOffice.PhoneLines
	var destPhoneLines []string
	for _, phoneLine := range destinationStationPhoneLines {
		if phoneLine.Type == "voice" {
			destPhoneLines = append(destPhoneLines, phoneLine.Number)
		}
	}

	if len(destPhoneLines) > 0 {
		destinationPhone := edisegment.PER{
			ContactFunctionCode:          "CN",
			CommunicationNumberQualifier: "TE",
			CommunicationNumber:          destPhoneLines[0],
		}
		header.Put("PER_DestinationPhone", &destinationPhone)
	}

	// ========  ORIGIN ========= //
	// origin station name
	var originDutyStation models.DutyStation

	if orders.OriginDutyStationID != nil && *orders.OriginDutyStationID != uuid.Nil {
		originDutyStation, err = models.FetchDutyStation(g.db, *orders.OriginDutyStationID)
		if err != nil {
			return services.NewInvalidInputError(*orders.OriginDutyStationID, err, nil, "unable to find origin duty station")
		}
	} else {
		return services.NewBadDataError("Invalid Order, must have OriginDutyStation")
	}

	originTransportationOffice, err := models.FetchDutyStationTransportationOffice(g.db, originDutyStation.ID)
	if err != nil {
		return services.NewInvalidInputError(originDutyStation.ID, err, nil, "unable to find transportation office of origin duty station")
	}

	originName := edisegment.N1{
		EntityIdentifierCode:        "SF",
		Name:                        originDutyStation.Name,
		IdentificationCodeQualifier: "10",
		IdentificationCode:          originTransportationOffice.Gbloc,
	}
	header.Put("N1_OriginName", &originName)

	// origin address
	if len(originDutyStation.Address.StreetAddress1) > 0 {
		originStreetAddress := edisegment.N3{
			AddressInformation1: originDutyStation.Address.StreetAddress1,
		}
		if originDutyStation.Address.StreetAddress2 != nil {
			originStreetAddress.AddressInformation2 = *originDutyStation.Address.StreetAddress2
		}
		header.Put("N3_OriginStreetAddress", &originStreetAddress)
	}

	// origin city/state/postal
	originPostalDetails := edisegment.N4{
		CityName:            originDutyStation.Address.City,
		StateOrProvinceCode: originDutyStation.Address.State,
		PostalCode:          originDutyStation.Address.PostalCode,
	}
	if originDutyStation.Address.Country != nil {
		countryCode, ccErr := originDutyStation.Address.CountryCode()
		if ccErr != nil {
			return ccErr
		}
		originPostalDetails.CountryCode = string(*countryCode)
	}

	header.Put("N4_OriginStreetAddress", &originPostalDetails)

	// Origin Station Phone
	originStationPhoneLines := originTransportationOffice.PhoneLines
	var originPhoneLines []string
	for _, phoneLine := range originStationPhoneLines {
		if phoneLine.Type == "voice" {
			originPhoneLines = append(originPhoneLines, phoneLine.Number)
		}
	}

	if len(originPhoneLines) > 0 {
		originPhone := edisegment.PER{
			ContactFunctionCode:          "CN",
			CommunicationNumberQualifier: "TE",
			CommunicationNumber:          originPhoneLines[0],
		}
		header.Put("PER_OriginPhone", &originPhone)
	}

	return nil
}

func (g ghcPaymentRequestInvoiceGenerator) createLoaSegments(orders models.Order) ([]edisegment.Segment, error) {
	segments := []edisegment.Segment{}
	if orders.TAC == nil || *orders.TAC == "" {
		return segments, services.NewConflictError(orders.ID, "Invalid order. Must have a TAC value")
	}
	affiliation := models.ServiceMemberAffiliation(*orders.DepartmentIndicator)
	agencyQualifierCode, found := edisegment.AffiliationToAgency[affiliation]

	if !found {
		agencyQualifierCode = "DF"
	}

	fa1 := edisegment.FA1{
		AgencyQualifierCode: agencyQualifierCode,
	}

	segments = append(segments, &fa1)

	fa2 := edisegment.FA2{
		BreakdownStructureDetailCode: "TA",
		FinancialInformationCode:     *orders.TAC,
	}

	segments = append(segments, &fa2)

	return segments, nil
}

func (g ghcPaymentRequestInvoiceGenerator) fetchPaymentServiceItemParam(serviceItemID uuid.UUID, key models.ServiceItemParamName) (models.PaymentServiceItemParam, error) {
	var paymentServiceItemParam models.PaymentServiceItemParam

	err := g.db.Q().
		Join("service_item_param_keys sk", "payment_service_item_params.service_item_param_key_id = sk.id").
		Where("payment_service_item_id = ?", serviceItemID).
		Where("sk.key = ?", key).
		First(&paymentServiceItemParam)
	if err != nil {
		if err.Error() == models.RecordNotFoundErrorString {
			return models.PaymentServiceItemParam{}, services.NewNotFoundError(serviceItemID, "for paymentServiceItemParam")
		}
		return models.PaymentServiceItemParam{}, services.NewQueryError("paymentServiceItemParam", err, fmt.Sprintf("Could not lookup PaymentServiceItemParam key (%s) payment service item id (%s): %s", key, serviceItemID, err))
	}
	return paymentServiceItemParam, nil
}

func (g ghcPaymentRequestInvoiceGenerator) getWeightParams(serviceItem models.PaymentServiceItem) (float64, error) {
	weight, err := g.fetchPaymentServiceItemParam(serviceItem.ID, models.ServiceItemParamNameWeightBilledActual)
	if err != nil {
		return 0, err
	}
	weightFloat, err := strconv.ParseFloat(weight.Value, 64)
	if err != nil {
		return 0, fmt.Errorf("Could not parse weight for PaymentServiceItem %s: %w", serviceItem.ID, err)
	}

	return weightFloat, nil
}

func (g ghcPaymentRequestInvoiceGenerator) getWeightAndDistanceParams(serviceItem models.PaymentServiceItem) (float64, float64, error) {
	// TODO: update to have a case statement as different service items may or may not have weight
	// and the distance key can differ (zip3 v zip5, and distances for SIT)
	weightFloat, err := g.getWeightParams(serviceItem)
	if err != nil {
		return 0, 0, err
	}
	distanceModel := models.ServiceItemParamNameDistanceZip3
	if serviceItem.MTOServiceItem.ReService.Code == models.ReServiceCodeDSH {
		distanceModel = models.ServiceItemParamNameDistanceZip5
	}
	distance, err := g.fetchPaymentServiceItemParam(serviceItem.ID, distanceModel)
	if err != nil {
		return 0, 0, err
	}
	distanceFloat, err := strconv.ParseFloat(distance.Value, 64)
	if err != nil {
		return 0, 0, fmt.Errorf("Could not parse Distance Zip3 for PaymentServiceItem %s: %w", serviceItem.ID, err)
	}
	return weightFloat, distanceFloat, nil
}

func (g ghcPaymentRequestInvoiceGenerator) generatePaymentServiceItemSegments(paymentServiceItems models.PaymentServiceItems, orders models.Order) ([]edisegment.Segment, error) {
	//Initialize empty collection of segments
	var segments []edisegment.Segment
	var weightFloat, distanceFloat float64
	// Iterate over payment service items
	for idx, serviceItem := range paymentServiceItems {
		if serviceItem.PriceCents == nil {
			return segments, services.NewConflictError(serviceItem.ID, "Invalid service item. Must have a PriceCents value")
		}
		hierarchicalIDNumber := idx + 1
		// Build and put together the segments
		hlSegment := edisegment.HL{
			HierarchicalIDNumber:  strconv.Itoa(hierarchicalIDNumber), // may need to change if sending multiple payment request in a single edi
			HierarchicalLevelCode: "I",
		}

		n9Segment := edisegment.N9{
			ReferenceIdentificationQualifier: "PO",
			ReferenceIdentification:          serviceItem.ReferenceID,
		}
		// TODO: add another n9 for SIT

		// Determine the correct params to use based off of the particular ReService code
		serviceCode := serviceItem.MTOServiceItem.ReService.Code
		switch serviceCode {
		case models.ReServiceCodeCS, models.ReServiceCodeMS:
			l5Segment := edisegment.L5{
				LadingLineItemNumber:   hierarchicalIDNumber,
				LadingDescription:      string(serviceCode),
				CommodityCode:          "TBD",
				CommodityCodeQualifier: "D",
			}

			l0Segment := edisegment.L0{
				LadingLineItemNumber: hierarchicalIDNumber,
			}

			segments = append(segments, &hlSegment, &n9Segment, &l5Segment, &l0Segment)
		// pack and unpack, dom dest and dom origin have weight no distance
		case models.ReServiceCodeDOP, models.ReServiceCodeDUPK,
			models.ReServiceCodeDPK, models.ReServiceCodeDDP:
			var err error
			weightFloat, err = g.getWeightParams(serviceItem)
			if err != nil {
				return segments, err
			}

			l5Segment := edisegment.L5{
				LadingLineItemNumber:   hierarchicalIDNumber,
				LadingDescription:      string(serviceCode),
				CommodityCode:          "TBD",
				CommodityCodeQualifier: "D",
			}

			l0Segment := edisegment.L0{
				LadingLineItemNumber: hierarchicalIDNumber,
				Weight:               weightFloat,
				WeightQualifier:      "B",
				WeightUnitCode:       "L",
			}

			segments = append(segments, &hlSegment, &n9Segment, &l5Segment, &l0Segment)

		default:
			var err error
			weightFloat, distanceFloat, err = g.getWeightAndDistanceParams(serviceItem)
			if err != nil {
				return segments, err
			}

			l5Segment := edisegment.L5{
				LadingLineItemNumber:   hierarchicalIDNumber,
				LadingDescription:      string(serviceCode),
				CommodityCode:          "TBD",
				CommodityCodeQualifier: "D",
			}

			l0Segment := edisegment.L0{
				LadingLineItemNumber:   hierarchicalIDNumber,
				BilledRatedAsQuantity:  distanceFloat,
				BilledRatedAsQualifier: "DM",
				Weight:                 weightFloat,
				WeightQualifier:        "B",
				WeightUnitCode:         "L",
			}

			segments = append(segments, &hlSegment, &n9Segment, &l5Segment, &l0Segment)
		}

		loaSegments, err := g.createLoaSegments(orders)
		if err != nil {
			return segments, err
		}
		segments = append(segments, loaSegments...)
	}

	l3Segment := edisegment.L3{
		PriceCents: 0, // TODO: hard-coded to zero for now
	}

	segments = append(segments, &l3Segment)

	return segments, nil
}

func msOrCsOnly(paymentServiceItems models.PaymentServiceItems) bool {
	for _, psi := range paymentServiceItems {
		code := psi.MTOServiceItem.ReService.Code
		if code != models.ReServiceCodeMS && code != models.ReServiceCodeCS {
			return false
		}
	}

	return true
}
