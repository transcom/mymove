package invoice

import (
	"fmt"
	"strconv"

	"github.com/benbjohnson/clock"
	"github.com/gobuffalo/pop/v5"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/db/sequence"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"

	ediinvoice "github.com/transcom/mymove/pkg/edi/invoice"
	edisegment "github.com/transcom/mymove/pkg/edi/segment"
)

/*
	NOTE: The GCN from GS06 and GE02 will match the ICN in ISA13 and IEA02,
	which restricts the 858 to only ever have 1 functional group.
	If multiple functional groups are needed, this will have to change.
*/

type ghcPaymentRequestInvoiceGenerator struct {
	db           *pop.Connection
	icnSequencer sequence.Sequencer
	clock        clock.Clock
}

// NewGHCPaymentRequestInvoiceGenerator returns an implementation of the GHCPaymentRequestInvoiceGenerator interface
func NewGHCPaymentRequestInvoiceGenerator(db *pop.Connection, icnSequencer sequence.Sequencer, clock clock.Clock) services.GHCPaymentRequestInvoiceGenerator {
	return &ghcPaymentRequestInvoiceGenerator{
		db:           db,
		icnSequencer: icnSequencer,
		clock:        clock,
	}
}

const dateFormat = "20060102"
const isaDateFormat = "060102"
const timeFormat = "1504"
const maxCityLength = 30

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

	currentTime := g.clock.Now()

	interchangeControlNumber, err := g.icnSequencer.NextVal()
	if err != nil {
		return ediinvoice.Invoice858C{}, fmt.Errorf("Failed to get next Interchange Control Number: %w", err)
	}

	// save ICN
	pr2icn := models.PaymentRequestToInterchangeControlNumber{
		PaymentRequestID:         paymentRequest.ID,
		InterchangeControlNumber: int(interchangeControlNumber),
	}
	verrs, err := g.db.ValidateAndSave(&pr2icn)
	if err != nil {
		return ediinvoice.Invoice858C{}, fmt.Errorf("Failed to save Interchange Control Number: %w", err)
	} else if verrs != nil && verrs.HasAny() {
		return ediinvoice.Invoice858C{}, fmt.Errorf("Failed to save Interchange Control Number: %s", verrs.String())
	}

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
		GroupControlNumber:       interchangeControlNumber,
		ResponsibleAgencyCode:    "X",
		Version:                  "004010",
	}

	edi858.ST = edisegment.ST{
		TransactionSetIdentifierCode: "858",
		TransactionSetControlNumber:  "0001",
	}

	edi858.Header.ShipmentInformation = edisegment.BX{
		TransactionSetPurposeCode:    "00",
		TransactionMethodTypeCode:    "J",
		ShipmentMethodOfPayment:      "PP",
		ShipmentIdentificationNumber: paymentRequest.PaymentRequestNumber,
		StandardCarrierAlphaCode:     "TRUS",
		ShipmentQualifier:            "4",
	}

	edi858.Header.PaymentRequestNumber = edisegment.N9{
		ReferenceIdentificationQualifier: "CN",
		ReferenceIdentification:          paymentRequest.PaymentRequestNumber,
	}

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

	edi858.Header.ContractCode = edisegment.N9{
		ReferenceIdentificationQualifier: "CT",
		ReferenceIdentification:          contractCodeServiceItemParam.Value,
	}

	// Add service member details to header
	err = g.createServiceMemberDetailSegments(paymentRequest.ID, moveTaskOrder.Orders.ServiceMember, &edi858.Header)
	if err != nil {
		return ediinvoice.Invoice858C{}, err
	}

	var paymentServiceItems models.PaymentServiceItems
	err = g.db.Q().
		Eager("MTOServiceItem.ReService").
		Where("payment_request_id = ?", paymentRequest.ID).
		Where("status = ?", models.PaymentServiceItemStatusApproved).
		All(&paymentServiceItems)
	if err != nil {
		if err.Error() == models.RecordNotFoundErrorString {
			return ediinvoice.Invoice858C{}, services.NewNotFoundError(paymentRequest.ID, "for payment service items in PaymentRequest")
		}
		return ediinvoice.Invoice858C{}, services.NewQueryError("PaymentServiceItems", err, fmt.Sprintf("error while looking for payment service items on payment request: %s", err))
	}

	if len(paymentServiceItems) == 0 {
		return ediinvoice.Invoice858C{}, services.NewConflictError(paymentRequest.ID, "this payment request has no approved PaymentServiceItems")
	}

	if !msOrCsOnly(paymentServiceItems) {
		err = g.createG62Segments(paymentRequest.ID, &edi858.Header)
		if err != nil {
			return ediinvoice.Invoice858C{}, err
		}
	}

	// Add buyer and seller organization names
	err = g.createBuyerAndSellerOrganizationNamesSegments(paymentRequest.ID, moveTaskOrder.Orders, &edi858.Header)
	if err != nil {
		return ediinvoice.Invoice858C{}, err
	}

	// Add origin and destination details to header
	err = g.createOriginAndDestinationSegments(paymentRequest.ID, moveTaskOrder.Orders, &edi858.Header)
	if err != nil {
		return ediinvoice.Invoice858C{}, err
	}

	var l3 edisegment.L3
	paymentServiceItemSegments, l3, err := g.generatePaymentServiceItemSegments(paymentServiceItems, moveTaskOrder.Orders)
	if err != nil {
		return ediinvoice.Invoice858C{}, err
	}
	edi858.ServiceItems = append(edi858.ServiceItems, paymentServiceItemSegments...)
	edi858.L3 = l3

	// the total NumberOfIncludedSegments is ST + SE + all segments other than GS, GE, ISA, and IEA
	stCount := 1
	l3Count := 1
	seCount := 1
	headerSegmentCount := edi858.Header.Size()
	serviceItemSegmentCount := len(edi858.ServiceItems) * ediinvoice.ServiceItemSegmentsSize
	totalNumberOfSegments := stCount + headerSegmentCount + serviceItemSegmentCount + l3Count + seCount

	edi858.SE = edisegment.SE{
		NumberOfIncludedSegments:    totalNumberOfSegments,
		TransactionSetControlNumber: "0001",
	}

	edi858.GE = edisegment.GE{
		NumberOfTransactionSetsIncluded: 1,
		GroupControlNumber:              interchangeControlNumber,
	}

	edi858.IEA = edisegment.IEA{
		NumberOfIncludedFunctionalGroups: 1,
		InterchangeControlNumber:         interchangeControlNumber,
	}

	return edi858, nil
}

func (g ghcPaymentRequestInvoiceGenerator) createServiceMemberDetailSegments(paymentRequestID uuid.UUID, serviceMember models.ServiceMember, header *ediinvoice.InvoiceHeader) error {

	// name
	header.ServiceMemberName = edisegment.N9{
		ReferenceIdentificationQualifier: "1W",
		ReferenceIdentification:          serviceMember.ReverseNameLineFormat(),
	}

	// rank
	rank := serviceMember.Rank
	if rank == nil {
		return services.NewConflictError(serviceMember.ID, fmt.Sprintf("no rank found for ServiceMember ID: %s Payment Request ID: %s", serviceMember.ID, paymentRequestID))
	}
	header.ServiceMemberRank = edisegment.N9{
		ReferenceIdentificationQualifier: "ML",
		ReferenceIdentification:          string(*rank),
	}

	// branch
	branch := serviceMember.Affiliation
	if branch == nil {
		return services.NewConflictError(serviceMember.ID, fmt.Sprintf("no branch found for ServiceMember ID: %s Payment Request ID: %s", serviceMember.ID, paymentRequestID))
	}
	header.ServiceMemberBranch = edisegment.N9{
		ReferenceIdentificationQualifier: "3L",
		ReferenceIdentification:          string(*branch),
	}

	return nil
}

func (g ghcPaymentRequestInvoiceGenerator) createG62Segments(paymentRequestID uuid.UUID, header *ediinvoice.InvoiceHeader) error {
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
			return services.NewNotFoundError(paymentRequestID, "for mto shipments associated with PaymentRequest")
		}
		return services.NewQueryError("MTOShipments", err, fmt.Sprintf("error querying for shipments to use in G62 segments in PaymentRequest %s: %s", paymentRequestID, err))
	}

	// If no shipments, then just return because we will not have access to the dates.
	if len(shipments) == 0 {
		return nil
	}

	// Use the first (earliest) shipment.
	shipment := shipments[0]

	// Insert request pickup date, if available.
	if shipment.RequestedPickupDate != nil {
		requestedPickupDateSegment := edisegment.G62{
			DateQualifier: 10,
			Date:          shipment.RequestedPickupDate.Format(dateFormat),
		}
		header.RequestedPickupDate = &requestedPickupDateSegment
	}

	// Insert expected pickup date, if available.
	if shipment.ScheduledPickupDate != nil {
		scheduledPickupDateSegment := edisegment.G62{
			DateQualifier: 76,
			Date:          shipment.ScheduledPickupDate.Format(dateFormat),
		}
		header.ScheduledPickupDate = &scheduledPickupDateSegment
	}

	// Insert expected pickup date, if available.
	if shipment.ActualPickupDate != nil {
		actualPickupDateSegment := edisegment.G62{
			DateQualifier: 86,
			Date:          shipment.ActualPickupDate.Format(dateFormat),
		}
		header.ActualPickupDate = &actualPickupDateSegment
	}

	return nil
}

func (g ghcPaymentRequestInvoiceGenerator) createBuyerAndSellerOrganizationNamesSegments(paymentRequestID uuid.UUID, orders models.Order, header *ediinvoice.InvoiceHeader) error {

	var err error
	var originDutyStation models.DutyStation

	if orders.OriginDutyStationID != nil && *orders.OriginDutyStationID != uuid.Nil {
		originDutyStation, err = models.FetchDutyStation(g.db, *orders.OriginDutyStationID)
		if err != nil {
			return services.NewInvalidInputError(*orders.OriginDutyStationID, err, nil, "unable to find origin duty station")
		}
	} else {
		return services.NewConflictError(orders.ID, "Invalid Order, must have OriginDutyStation")
	}

	originTransportationOffice, err := models.FetchDutyStationTransportationOffice(g.db, originDutyStation.ID)
	if err != nil {
		return services.NewInvalidInputError(originDutyStation.ID, err, nil, "unable to find origin duty station")
	}

	// buyer organization name
	header.BuyerOrganizationName = edisegment.N1{
		EntityIdentifierCode:        "BY",
		Name:                        originTransportationOffice.Name,
		IdentificationCodeQualifier: "92",
		IdentificationCode:          originTransportationOffice.Gbloc,
	}

	// seller organization name
	header.SellerOrganizationName = edisegment.N1{
		EntityIdentifierCode:        "SE",
		Name:                        "Prime",
		IdentificationCodeQualifier: "2",
		IdentificationCode:          "PRME",
	}

	return nil
}

func (g ghcPaymentRequestInvoiceGenerator) createOriginAndDestinationSegments(paymentRequestID uuid.UUID, orders models.Order, header *ediinvoice.InvoiceHeader) error {
	var err error
	var destinationDutyStation models.DutyStation
	if orders.NewDutyStationID != uuid.Nil {
		destinationDutyStation, err = models.FetchDutyStation(g.db, orders.NewDutyStationID)
		if err != nil {
			return services.NewInvalidInputError(orders.NewDutyStationID, err, nil, "unable to find new duty station")
		}
	} else {
		return services.NewConflictError(orders.ID, "Invalid Order, must have NewDutyStation")
	}

	destTransportationOffice, err := models.FetchDutyStationTransportationOffice(g.db, destinationDutyStation.ID)
	if err != nil {
		return services.NewInvalidInputError(destinationDutyStation.ID, err, nil, "unable to find destination duty station")
	}

	// destination name
	header.DestinationName = edisegment.N1{
		EntityIdentifierCode:        "ST",
		Name:                        destinationDutyStation.Name,
		IdentificationCodeQualifier: "10",
		IdentificationCode:          destTransportationOffice.Gbloc,
	}

	// destination address
	if len(destinationDutyStation.Address.StreetAddress1) > 0 {
		destinationStreetAddress := edisegment.N3{
			AddressInformation1: destinationDutyStation.Address.StreetAddress1,
		}
		if destinationDutyStation.Address.StreetAddress2 != nil {
			destinationStreetAddress.AddressInformation2 = *destinationDutyStation.Address.StreetAddress2
		}
		header.DestinationStreetAddress = &destinationStreetAddress
	}

	// destination city/state/postal
	header.DestinationPostalDetails = edisegment.N4{
		CityName:            truncateStr(destinationDutyStation.Address.City, maxCityLength),
		StateOrProvinceCode: destinationDutyStation.Address.State,
		PostalCode:          destinationDutyStation.Address.PostalCode,
	}
	if destinationDutyStation.Address.Country != nil {
		countryCode, ccErr := destinationDutyStation.Address.CountryCode()
		if ccErr != nil {
			return ccErr
		}
		header.DestinationPostalDetails.CountryCode = string(*countryCode)
	}

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
		header.DestinationPhone = &destinationPhone
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
		return services.NewConflictError(orders.ID, "Invalid Order, must have OriginDutyStation")
	}

	originTransportationOffice, err := models.FetchDutyStationTransportationOffice(g.db, originDutyStation.ID)
	if err != nil {
		return services.NewInvalidInputError(originDutyStation.ID, err, nil, "unable to find transportation office of origin duty station")
	}

	header.OriginName = edisegment.N1{
		EntityIdentifierCode:        "SF",
		Name:                        originDutyStation.Name,
		IdentificationCodeQualifier: "10",
		IdentificationCode:          originTransportationOffice.Gbloc,
	}

	// origin address
	if len(originDutyStation.Address.StreetAddress1) > 0 {
		originStreetAddress := edisegment.N3{
			AddressInformation1: originDutyStation.Address.StreetAddress1,
		}
		if originDutyStation.Address.StreetAddress2 != nil {
			originStreetAddress.AddressInformation2 = *originDutyStation.Address.StreetAddress2
		}
		header.OriginStreetAddress = &originStreetAddress
	}

	// origin city/state/postal
	header.OriginPostalDetails = edisegment.N4{
		CityName:            truncateStr(originDutyStation.Address.City, maxCityLength),
		StateOrProvinceCode: originDutyStation.Address.State,
		PostalCode:          originDutyStation.Address.PostalCode,
	}
	if originDutyStation.Address.Country != nil {
		countryCode, ccErr := originDutyStation.Address.CountryCode()
		if ccErr != nil {
			return ccErr
		}
		header.OriginPostalDetails.CountryCode = string(*countryCode)
	}

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
		header.OriginPhone = &originPhone
	}

	return nil
}

func (g ghcPaymentRequestInvoiceGenerator) createLoaSegments(orders models.Order) (edisegment.FA1, edisegment.FA2, error) {
	if orders.TAC == nil || *orders.TAC == "" {
		return edisegment.FA1{}, edisegment.FA2{}, services.NewConflictError(orders.ID, "Invalid order. Must have a TAC value")
	}
	affiliation := models.ServiceMemberAffiliation(*orders.DepartmentIndicator)
	agencyQualifierCode, found := edisegment.AffiliationToAgency[affiliation]

	if !found {
		agencyQualifierCode = "DF"
	}

	fa1 := edisegment.FA1{
		AgencyQualifierCode: agencyQualifierCode,
	}

	fa2 := edisegment.FA2{
		BreakdownStructureDetailCode: "TA",
		FinancialInformationCode:     *orders.TAC,
	}

	return fa1, fa2, nil
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

func (g ghcPaymentRequestInvoiceGenerator) getWeightParams(serviceItem models.PaymentServiceItem) (int, error) {
	weight, err := g.fetchPaymentServiceItemParam(serviceItem.ID, models.ServiceItemParamNameWeightBilledActual)
	if err != nil {
		return 0, err
	}
	weightInt, err := strconv.Atoi(weight.Value)
	if err != nil {
		return 0, fmt.Errorf("Could not parse weight for PaymentServiceItem %s: %w", serviceItem.ID, err)
	}

	return weightInt, nil
}

func (g ghcPaymentRequestInvoiceGenerator) getWeightAndDistanceParams(serviceItem models.PaymentServiceItem) (int, float64, error) {
	weight, err := g.getWeightParams(serviceItem)
	if err != nil {
		return 0, 0, err
	}

	var distanceModel models.ServiceItemParamName
	switch serviceItem.MTOServiceItem.ReService.Code {
	case models.ReServiceCodeDSH:
		distanceModel = models.ServiceItemParamNameDistanceZip5
	case models.ReServiceCodeDDDSIT:
		distanceModel = models.ServiceItemParamNameDistanceZipSITDest
	case models.ReServiceCodeDOPSIT:
		distanceModel = models.ServiceItemParamNameDistanceZipSITOrigin
	default:
		distanceModel = models.ServiceItemParamNameDistanceZip3
	}

	distance, err := g.fetchPaymentServiceItemParam(serviceItem.ID, distanceModel)
	if err != nil {
		return 0, 0, err
	}
	distanceFloat, err := strconv.ParseFloat(distance.Value, 64)
	if err != nil {
		return 0, 0, fmt.Errorf("Could not parse Distance Zip3 for PaymentServiceItem %s: %w", serviceItem.ID, err)
	}
	return weight, distanceFloat, nil
}

func (g ghcPaymentRequestInvoiceGenerator) generatePaymentServiceItemSegments(paymentServiceItems models.PaymentServiceItems, orders models.Order) ([]ediinvoice.ServiceItemSegments, edisegment.L3, error) {
	//Initialize empty collection of segments
	var segments []ediinvoice.ServiceItemSegments
	l3 := edisegment.L3{
		PriceCents: 0,
	}
	// Iterate over payment service items
	for idx, serviceItem := range paymentServiceItems {
		var newSegment ediinvoice.ServiceItemSegments
		if serviceItem.PriceCents == nil {
			return segments, l3, services.NewConflictError(serviceItem.ID, "Invalid service item. Must have a PriceCents value")
		}
		l3.PriceCents += int64(*serviceItem.PriceCents)
		hierarchicalIDNumber := idx + 1
		// Build and put together the segments
		newSegment.HL = edisegment.HL{
			HierarchicalIDNumber:  strconv.Itoa(hierarchicalIDNumber), // may need to change if sending multiple payment request in a single edi
			HierarchicalLevelCode: "I",
		}

		newSegment.N9 = edisegment.N9{
			ReferenceIdentificationQualifier: "PO",
			ReferenceIdentification:          serviceItem.ReferenceID,
		}
		// TODO: add another n9 for SIT

		// Determine the correct params to use based off of the particular ReService code
		serviceCode := serviceItem.MTOServiceItem.ReService.Code
		switch serviceCode {
		// cs and ms have no weight and no distance
		case models.ReServiceCodeCS, models.ReServiceCodeMS:
			newSegment.L5 = edisegment.L5{
				LadingLineItemNumber:   hierarchicalIDNumber,
				LadingDescription:      string(serviceCode),
				CommodityCode:          "TBD",
				CommodityCodeQualifier: "D",
			}

			newSegment.L0 = edisegment.L0{
				LadingLineItemNumber: hierarchicalIDNumber,
			}

			newSegment.L1 = edisegment.L1{
				LadingLineItemNumber: hierarchicalIDNumber,
				Charge:               serviceItem.PriceCents.Int64(),
			}

		// following service items have weight and no distance
		case models.ReServiceCodeDOP, models.ReServiceCodeDUPK,
			models.ReServiceCodeDPK, models.ReServiceCodeDDP,
			models.ReServiceCodeDDFSIT, models.ReServiceCodeDDASIT,
			models.ReServiceCodeDOFSIT, models.ReServiceCodeDOASIT:
			var err error
			weight, err := g.getWeightParams(serviceItem)
			if err != nil {
				return segments, l3, err
			}

			newSegment.L5 = edisegment.L5{
				LadingLineItemNumber:   hierarchicalIDNumber,
				LadingDescription:      string(serviceCode),
				CommodityCode:          "TBD",
				CommodityCodeQualifier: "D",
			}

			newSegment.L0 = edisegment.L0{
				LadingLineItemNumber: hierarchicalIDNumber,
				Weight:               float64(weight),
				WeightQualifier:      "B",
				WeightUnitCode:       "L",
			}

			newSegment.L1 = edisegment.L1{
				LadingLineItemNumber: hierarchicalIDNumber,
				FreightRate:          &weight,
				RateValueQualifier:   "LB",
				Charge:               serviceItem.PriceCents.Int64(),
			}

		default:
			var err error
			weight, distanceFloat, err := g.getWeightAndDistanceParams(serviceItem)
			if err != nil {
				return segments, l3, err
			}

			newSegment.L5 = edisegment.L5{
				LadingLineItemNumber:   hierarchicalIDNumber,
				LadingDescription:      string(serviceCode),
				CommodityCode:          "TBD",
				CommodityCodeQualifier: "D",
			}

			newSegment.L0 = edisegment.L0{
				LadingLineItemNumber:   hierarchicalIDNumber,
				BilledRatedAsQuantity:  distanceFloat,
				BilledRatedAsQualifier: "DM",
				Weight:                 float64(weight),
				WeightQualifier:        "B",
				WeightUnitCode:         "L",
			}

			newSegment.L1 = edisegment.L1{
				LadingLineItemNumber: hierarchicalIDNumber,
				FreightRate:          &weight,
				RateValueQualifier:   "LB",
				Charge:               serviceItem.PriceCents.Int64(),
			}

		}

		fa1, fa2, err := g.createLoaSegments(orders)
		if err != nil {
			return segments, l3, err
		}
		newSegment.FA1 = fa1
		newSegment.FA2 = fa2
		segments = append(segments, newSegment)
	}

	return segments, l3, nil
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

func truncateStr(str string, cutoff int) string {
	if len(str) >= cutoff {
		if cutoff-3 > 0 {
			return str[:cutoff-3] + "..."
		}
		return str[:cutoff]
	}
	return str
}
