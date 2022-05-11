package invoice

import (
	"database/sql"
	"fmt"
	"regexp"
	"strconv"

	"github.com/benbjohnson/clock"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/db/sequence"
	ediinvoice "github.com/transcom/mymove/pkg/edi/invoice"
	edisegment "github.com/transcom/mymove/pkg/edi/segment"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

/*
	NOTE: The GCN from GS06 and GE02 will match the ICN in ISA13 and IEA02,
	which restricts the 858 to only ever have 1 functional group.
	If multiple functional groups are needed, this will have to change.
*/

type ghcPaymentRequestInvoiceGenerator struct {
	icnSequencer sequence.Sequencer
	clock        clock.Clock
}

// NewGHCPaymentRequestInvoiceGenerator returns an implementation of the GHCPaymentRequestInvoiceGenerator interface
func NewGHCPaymentRequestInvoiceGenerator(icnSequencer sequence.Sequencer, clock clock.Clock) services.GHCPaymentRequestInvoiceGenerator {
	return &ghcPaymentRequestInvoiceGenerator{
		icnSequencer: icnSequencer,
		clock:        clock,
	}
}

const dateFormat = "20060102"
const isaDateFormat = "060102"
const timeFormat = "1504"
const maxCityLength = 30

// Generate method takes a payment request and returns an Invoice858C
func (g ghcPaymentRequestInvoiceGenerator) Generate(appCtx appcontext.AppContext, paymentRequest models.PaymentRequest, sendProductionInvoice bool) (ediinvoice.Invoice858C, error) {
	var moveTaskOrder models.Move
	if paymentRequest.MoveTaskOrder.ID == uuid.Nil {
		// load mto
		err := appCtx.DB().Q().
			Where("id = ?", paymentRequest.MoveTaskOrderID).
			First(&moveTaskOrder)
		if err != nil {
			switch err {
			case sql.ErrNoRows:
				return ediinvoice.Invoice858C{}, apperror.NewNotFoundError(paymentRequest.MoveTaskOrder.ID, "for MoveTaskOrder")
			default:
				return ediinvoice.Invoice858C{}, apperror.NewQueryError("MoveTaskOrder", err, "Unexpected error")
			}
		}
	} else {
		moveTaskOrder = paymentRequest.MoveTaskOrder
	}

	// check or load orders
	if moveTaskOrder.ReferenceID == nil {
		return ediinvoice.Invoice858C{}, apperror.NewConflictError(moveTaskOrder.ID, "Invalid move taskorder. Must have a ReferenceID value")
	}

	if moveTaskOrder.Orders.ID == uuid.Nil {
		err := appCtx.DB().
			Load(&moveTaskOrder, "Orders")
		if err != nil {
			switch err {
			case sql.ErrNoRows:
				return ediinvoice.Invoice858C{}, apperror.NewNotFoundError(moveTaskOrder.Orders.ID, "for Orders")
			default:
				return ediinvoice.Invoice858C{}, apperror.NewQueryError("Orders", err, "Unexpected error")
			}
		}
	}

	// check or load service member
	if moveTaskOrder.Orders.ServiceMember.ID == uuid.Nil {
		err := appCtx.DB().
			Load(&moveTaskOrder.Orders, "ServiceMember")

		if err != nil {
			switch err {
			case sql.ErrNoRows:
				return ediinvoice.Invoice858C{}, apperror.NewNotFoundError(moveTaskOrder.Orders.ServiceMemberID, "for ServiceMember")
			default:
				return ediinvoice.Invoice858C{}, apperror.NewQueryError("ServiceMember", err, fmt.Sprintf("cannot load ServiceMember %s for PaymentRequest %s: %s", moveTaskOrder.Orders.ServiceMemberID, paymentRequest.ID, err))
			}
		}
	}

	currentTime := g.clock.Now()

	interchangeControlNumber, err := g.icnSequencer.NextVal(appCtx)
	if err != nil {
		return ediinvoice.Invoice858C{}, fmt.Errorf("Failed to get next Interchange Control Number: %w", err)
	}

	// save ICN
	pr2icn := models.PaymentRequestToInterchangeControlNumber{
		PaymentRequestID:         paymentRequest.ID,
		InterchangeControlNumber: int(interchangeControlNumber),
		EDIType:                  models.EDIType858,
	}
	verrs, err := appCtx.DB().ValidateAndSave(&pr2icn)
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
		StandardCarrierAlphaCode:     "BLKW",
		ShipmentQualifier:            "4",
	}

	edi858.Header.PaymentRequestNumber = edisegment.N9{
		ReferenceIdentificationQualifier: "CN",
		ReferenceIdentification:          paymentRequest.PaymentRequestNumber,
	}

	// contract code to header
	var contractCodeServiceItemParam models.PaymentServiceItemParam
	err = appCtx.DB().Q().
		Join("service_item_param_keys sipk", "payment_service_item_params.service_item_param_key_id = sipk.id").
		Join("payment_service_items psi", "payment_service_item_params.payment_service_item_id = psi.id").
		Join("payment_requests pr", "psi.payment_request_id = pr.id").
		Where("pr.id = ?", paymentRequest.ID).
		Where("sipk.key = ?", models.ServiceItemParamNameContractCode).
		First(&contractCodeServiceItemParam)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return ediinvoice.Invoice858C{}, apperror.NewNotFoundError(contractCodeServiceItemParam.ID, "for ContractCode")
		default:
			return ediinvoice.Invoice858C{}, apperror.NewQueryError("ContractCode", err, fmt.Sprintf("Couldn't find contract code: %s", err))
		}
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
	err = appCtx.DB().Q().
		Eager("MTOServiceItem.ReService", "MTOServiceItem.MTOShipment").
		Where("payment_request_id = ?", paymentRequest.ID).
		Where("status = ?", models.PaymentServiceItemStatusApproved).
		All(&paymentServiceItems)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return ediinvoice.Invoice858C{}, apperror.NewNotFoundError(paymentRequest.ID, "for payment service items in PaymentRequest")
		default:
			return ediinvoice.Invoice858C{}, apperror.NewQueryError("PaymentServiceItems", err, fmt.Sprintf("error while looking for payment service items on payment request: %s", err))
		}
	}

	// Add C3 segment here
	err = g.createC3Segment(&edi858.Header)
	if err != nil {
		return ediinvoice.Invoice858C{}, err
	}

	if len(paymentServiceItems) == 0 {
		return ediinvoice.Invoice858C{}, apperror.NewConflictError(paymentRequest.ID, "this payment request has no approved PaymentServiceItems")
	}

	if !msOrCsOnly(paymentServiceItems) {
		err = g.createG62Segments(appCtx, paymentRequest.ID, &edi858.Header)
		if err != nil {
			return ediinvoice.Invoice858C{}, err
		}
	}

	// Add buyer and seller organization names
	err = g.createBuyerAndSellerOrganizationNamesSegments(appCtx, paymentRequest.ID, moveTaskOrder.Orders, &edi858.Header)
	if err != nil {
		return ediinvoice.Invoice858C{}, err
	}

	// Add origin and destination details to header
	err = g.createOriginAndDestinationSegments(appCtx, paymentRequest.ID, moveTaskOrder.Orders, &edi858.Header)
	if err != nil {
		return ediinvoice.Invoice858C{}, err
	}

	var l3 edisegment.L3
	paymentServiceItemSegments, l3, err := g.generatePaymentServiceItemSegments(appCtx, paymentServiceItems, moveTaskOrder.Orders)
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
		return apperror.NewConflictError(serviceMember.ID, fmt.Sprintf("no rank found for ServiceMember ID: %s Payment Request ID: %s", serviceMember.ID, paymentRequestID))
	}
	header.ServiceMemberRank = edisegment.N9{
		ReferenceIdentificationQualifier: "ML",
		ReferenceIdentification:          string(*rank),
	}

	// branch
	branch := serviceMember.Affiliation
	if branch == nil {
		return apperror.NewConflictError(serviceMember.ID, fmt.Sprintf("no branch found for ServiceMember ID: %s Payment Request ID: %s", serviceMember.ID, paymentRequestID))
	}
	header.ServiceMemberBranch = edisegment.N9{
		ReferenceIdentificationQualifier: "3L",
		ReferenceIdentification:          string(*branch),
	}

	return nil
}

func (g ghcPaymentRequestInvoiceGenerator) createC3Segment(header *ediinvoice.InvoiceHeader) error {
	header.Currency = edisegment.C3{
		CurrencyCodeC301: "USD",
	}
	return nil
}

func (g ghcPaymentRequestInvoiceGenerator) createG62Segments(appCtx appcontext.AppContext, paymentRequestID uuid.UUID, header *ediinvoice.InvoiceHeader) error {
	// Get all the shipments associated with this payment request's service items, ordered by shipment creation date.
	var shipments models.MTOShipments
	err := appCtx.DB().Q().
		Join("mto_service_items msi", "mto_shipments.id = msi.mto_shipment_id").
		Join("payment_service_items psi", "msi.id = psi.mto_service_item_id").
		Where("psi.payment_request_id = ?", paymentRequestID).
		Order("msi.created_at").
		All(&shipments)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return apperror.NewNotFoundError(paymentRequestID, "for mto shipments associated with PaymentRequest")
		default:
			return apperror.NewQueryError("MTOShipments", err, fmt.Sprintf("error querying for shipments to use in G62 segments in PaymentRequest %s: %s", paymentRequestID, err))
		}
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

func (g ghcPaymentRequestInvoiceGenerator) createBuyerAndSellerOrganizationNamesSegments(appCtx appcontext.AppContext, paymentRequestID uuid.UUID, orders models.Order, header *ediinvoice.InvoiceHeader) error {

	var err error
	var originDutyLocation models.DutyLocation

	if orders.OriginDutyLocationID != nil && *orders.OriginDutyLocationID != uuid.Nil {
		originDutyLocation, err = models.FetchDutyLocation(appCtx.DB(), *orders.OriginDutyLocationID)
		if err != nil {
			return apperror.NewInvalidInputError(*orders.OriginDutyLocationID, err, nil, "unable to find origin duty location")
		}
	} else {
		return apperror.NewConflictError(orders.ID, "Invalid Order, must have OriginDutyLocation")
	}

	originTransportationOffice, err := models.FetchDutyLocationTransportationOffice(appCtx.DB(), originDutyLocation.ID)
	if err != nil {
		return apperror.NewInvalidInputError(originDutyLocation.ID, err, nil, "unable to find origin duty location")
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
		IdentificationCode:          "BLKW",
	}

	return nil
}

func (g ghcPaymentRequestInvoiceGenerator) createOriginAndDestinationSegments(appCtx appcontext.AppContext, paymentRequestID uuid.UUID, orders models.Order, header *ediinvoice.InvoiceHeader) error {
	var err error
	var destinationDutyLocation models.DutyLocation
	if orders.NewDutyLocationID != uuid.Nil {
		destinationDutyLocation, err = models.FetchDutyLocation(appCtx.DB(), orders.NewDutyLocationID)
		if err != nil {
			return apperror.NewInvalidInputError(orders.NewDutyLocationID, err, nil, "unable to find new duty location")
		}
	} else {
		return apperror.NewConflictError(orders.ID, "Invalid Order, must have NewDutyLocation")
	}

	destTransportationOffice, err := models.FetchDutyLocationTransportationOffice(appCtx.DB(), destinationDutyLocation.ID)
	if err != nil {
		return apperror.NewInvalidInputError(destinationDutyLocation.ID, err, nil, "unable to find destination duty location")
	}

	// destination name
	header.DestinationName = edisegment.N1{
		EntityIdentifierCode:        "ST",
		Name:                        destinationDutyLocation.Name,
		IdentificationCodeQualifier: "10",
		IdentificationCode:          destTransportationOffice.Gbloc,
	}

	// destination address
	if len(destinationDutyLocation.Address.StreetAddress1) > 0 {
		destinationStreetAddress := edisegment.N3{
			AddressInformation1: destinationDutyLocation.Address.StreetAddress1,
		}
		if destinationDutyLocation.Address.StreetAddress2 != nil {
			destinationStreetAddress.AddressInformation2 = *destinationDutyLocation.Address.StreetAddress2
		}
		header.DestinationStreetAddress = &destinationStreetAddress
	}

	// destination city/state/postal
	header.DestinationPostalDetails = edisegment.N4{
		CityName:            truncateStr(destinationDutyLocation.Address.City, maxCityLength),
		StateOrProvinceCode: destinationDutyLocation.Address.State,
		PostalCode:          destinationDutyLocation.Address.PostalCode,
	}
	if destinationDutyLocation.Address.Country != nil {
		countryCode, ccErr := destinationDutyLocation.Address.CountryCode()
		if ccErr != nil {
			return ccErr
		}
		header.DestinationPostalDetails.CountryCode = string(*countryCode)
	}

	// Destination PER
	destinationDutyLocationPhoneLines := destTransportationOffice.PhoneLines
	var destPhoneLines []string
	for _, phoneLine := range destinationDutyLocationPhoneLines {
		if phoneLine.Type == "voice" {
			destPhoneLines = append(destPhoneLines, phoneLine.Number)
		}
	}

	if len(destPhoneLines) > 0 {
		digits, digitsErr := g.getPhoneNumberDigitsOnly(destPhoneLines[0])
		if digitsErr != nil {
			return apperror.NewInvalidInputError(destinationDutyLocation.ID, digitsErr, nil, "unable to get destination duty location phone number")
		}
		destinationPhone := edisegment.PER{
			ContactFunctionCode:          "CN",
			CommunicationNumberQualifier: "TE",
			CommunicationNumber:          digits,
		}
		header.DestinationPhone = &destinationPhone
	}

	// ========  ORIGIN ========= //
	// origin duty location name
	var originDutyLocation models.DutyLocation

	if orders.OriginDutyLocationID != nil && *orders.OriginDutyLocationID != uuid.Nil {
		originDutyLocation, err = models.FetchDutyLocation(appCtx.DB(), *orders.OriginDutyLocationID)
		if err != nil {
			return apperror.NewInvalidInputError(*orders.OriginDutyLocationID, err, nil, "unable to find origin duty location")
		}
	} else {
		return apperror.NewConflictError(orders.ID, "Invalid Order, must have OriginDutyLocation")
	}

	originTransportationOffice, err := models.FetchDutyLocationTransportationOffice(appCtx.DB(), originDutyLocation.ID)
	if err != nil {
		return apperror.NewInvalidInputError(originDutyLocation.ID, err, nil, "unable to find transportation office of origin duty location")
	}

	header.OriginName = edisegment.N1{
		EntityIdentifierCode:        "SF",
		Name:                        originDutyLocation.Name,
		IdentificationCodeQualifier: "10",
		IdentificationCode:          originTransportationOffice.Gbloc,
	}

	// origin address
	if len(originDutyLocation.Address.StreetAddress1) > 0 {
		originStreetAddress := edisegment.N3{
			AddressInformation1: originDutyLocation.Address.StreetAddress1,
		}
		if originDutyLocation.Address.StreetAddress2 != nil {
			originStreetAddress.AddressInformation2 = *originDutyLocation.Address.StreetAddress2
		}
		header.OriginStreetAddress = &originStreetAddress
	}

	// origin city/state/postal
	header.OriginPostalDetails = edisegment.N4{
		CityName:            truncateStr(originDutyLocation.Address.City, maxCityLength),
		StateOrProvinceCode: originDutyLocation.Address.State,
		PostalCode:          originDutyLocation.Address.PostalCode,
	}
	if originDutyLocation.Address.Country != nil {
		countryCode, ccErr := originDutyLocation.Address.CountryCode()
		if ccErr != nil {
			return ccErr
		}
		header.OriginPostalDetails.CountryCode = string(*countryCode)
	}

	// Origin Duty Location Phone
	originDutyLocationPhoneLines := originTransportationOffice.PhoneLines
	var originPhoneLines []string
	for _, phoneLine := range originDutyLocationPhoneLines {
		if phoneLine.Type == "voice" {
			originPhoneLines = append(originPhoneLines, phoneLine.Number)
		}
	}

	if len(originPhoneLines) > 0 {
		digits, digitsErr := g.getPhoneNumberDigitsOnly(originPhoneLines[0])
		if digitsErr != nil {
			return apperror.NewInvalidInputError(originDutyLocation.ID, digitsErr, nil, "unable to get origin duty location phone number")
		}
		originPhone := edisegment.PER{
			ContactFunctionCode:          "CN",
			CommunicationNumberQualifier: "TE",
			CommunicationNumber:          digits,
		}
		header.OriginPhone = &originPhone
	}

	return nil
}

func (g ghcPaymentRequestInvoiceGenerator) createLoaSegments(orders models.Order, shipment models.MTOShipment) (edisegment.FA1, edisegment.FA2, error) {
	// We need to determine which TAC to use. We'll default to using the HHG TAC as that's what we've been doing
	// up to this point. But now we need to look at the service item's MTOShipment (if there is one -- some
	// service items like MS/CS aren't associated with a shipment) and see if it prefers the NTS TAC instead.
	useHHGTac := true
	if shipment.ID != uuid.Nil {
		// We do have a shipment, so see if the shipment prefers the NTS TAC.
		if shipment.TACType != nil && *shipment.TACType == models.LOATypeNTS {
			useHHGTac = false
		}
	}

	// Now grab the preferred TAC, making sure that it actually exists on the orders record.
	var tac string
	if useHHGTac {
		if orders.TAC == nil || *orders.TAC == "" {
			return edisegment.FA1{}, edisegment.FA2{}, apperror.NewConflictError(orders.ID, "Invalid order. Must have an HHG TAC value")
		}
		tac = *orders.TAC
	} else {
		if orders.NtsTAC == nil || *orders.NtsTAC == "" {
			return edisegment.FA1{}, edisegment.FA2{}, apperror.NewConflictError(orders.ID, "Invalid order. Must have an NTS TAC value")
		}
		tac = *orders.NtsTAC
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
		FinancialInformationCode:     tac,
	}

	return fa1, fa2, nil
}

func (g ghcPaymentRequestInvoiceGenerator) fetchPaymentServiceItemParam(appCtx appcontext.AppContext, serviceItemID uuid.UUID, key models.ServiceItemParamName) (models.PaymentServiceItemParam, error) {
	var paymentServiceItemParam models.PaymentServiceItemParam

	err := appCtx.DB().Q().
		Join("service_item_param_keys sk", "payment_service_item_params.service_item_param_key_id = sk.id").
		Where("payment_service_item_id = ?", serviceItemID).
		Where("sk.key = ?", key).
		First(&paymentServiceItemParam)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return models.PaymentServiceItemParam{}, apperror.NewNotFoundError(serviceItemID, "for paymentServiceItemParam")
		default:
			return models.PaymentServiceItemParam{}, apperror.NewQueryError("paymentServiceItemParam", err, fmt.Sprintf("Could not lookup PaymentServiceItemParam key (%s) payment service item id (%s): %s", key, serviceItemID, err))
		}
	}
	return paymentServiceItemParam, nil
}

func (g ghcPaymentRequestInvoiceGenerator) getPhoneNumberDigitsOnly(phoneString string) (string, error) {
	reg, err := regexp.Compile("[^0-9]+")
	if err != nil {
		return "", err
	}
	digitsOnly := reg.ReplaceAllString(phoneString, "")
	return digitsOnly, nil
}

func (g ghcPaymentRequestInvoiceGenerator) getWeightParams(appCtx appcontext.AppContext, serviceItem models.PaymentServiceItem) (int, error) {
	weight, err := g.fetchPaymentServiceItemParam(appCtx, serviceItem.ID, models.ServiceItemParamNameWeightBilled)
	if err != nil {
		return 0, err
	}
	weightInt, err := strconv.Atoi(weight.Value)
	if err != nil {
		return 0, fmt.Errorf("Could not parse weight for PaymentServiceItem %s: %w", serviceItem.ID, err)
	}

	return weightInt, nil
}

func (g ghcPaymentRequestInvoiceGenerator) getServiceItemDimensionRateParams(appCtx appcontext.AppContext, serviceItem models.PaymentServiceItem) (float64, float64, error) {
	cubicFeet, err := g.fetchPaymentServiceItemParam(appCtx, serviceItem.ID, models.ServiceItemParamNameCubicFeetBilled)
	if err != nil {
		return 0, 0, err
	}

	cubicFeetFloat, err := strconv.ParseFloat(cubicFeet.Value, 64)
	if err != nil {
		return 0, 0, fmt.Errorf("Could not parse cubic feet as a float for PaymentServiceItem %s: %w", serviceItem.ID, err)
	}

	rate, err := g.fetchPaymentServiceItemParam(appCtx, serviceItem.ID, models.ServiceItemParamNamePriceRateOrFactor)
	if err != nil {
		return 0, 0, err
	}
	rateFloat, err := strconv.ParseFloat(rate.Value, 64)
	if err != nil {
		return 0, 0, fmt.Errorf("Could not parse rate as a float for PaymentServiceItem %s: %w", serviceItem.ID, err)
	}

	return cubicFeetFloat, rateFloat, nil
}

func (g ghcPaymentRequestInvoiceGenerator) getWeightAndDistanceParams(appCtx appcontext.AppContext, serviceItem models.PaymentServiceItem) (int, float64, error) {
	weight, err := g.getWeightParams(appCtx, serviceItem)
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

	distance, err := g.fetchPaymentServiceItemParam(appCtx, serviceItem.ID, distanceModel)
	if err != nil {
		return 0, 0, err
	}
	distanceFloat, err := strconv.ParseFloat(distance.Value, 64)
	if err != nil {
		return 0, 0, fmt.Errorf("Could not parse Distance Zip3 for PaymentServiceItem %s: %w", serviceItem.ID, err)
	}
	return weight, distanceFloat, nil
}

func (g ghcPaymentRequestInvoiceGenerator) generatePaymentServiceItemSegments(appCtx appcontext.AppContext, paymentServiceItems models.PaymentServiceItems, orders models.Order) ([]ediinvoice.ServiceItemSegments, edisegment.L3, error) {
	//Initialize empty collection of segments
	var segments []ediinvoice.ServiceItemSegments
	l3 := edisegment.L3{
		PriceCents: 0,
	}
	// Iterate over payment service items
	for idx, serviceItem := range paymentServiceItems {
		var newSegment ediinvoice.ServiceItemSegments
		if serviceItem.PriceCents == nil {
			return segments, l3, apperror.NewConflictError(serviceItem.ID, "Invalid service item. Must have a PriceCents value")
		}
		l3.PriceCents += int64(*serviceItem.PriceCents)
		hierarchicalIDNumber := idx + 1
		// Build and put together the segments
		newSegment.HL = edisegment.HL{
			HierarchicalIDNumber:  strconv.Itoa(hierarchicalIDNumber), // may need to change if sending multiple payment request in a single edi
			HierarchicalLevelCode: "9",
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
			models.ReServiceCodeDOFSIT, models.ReServiceCodeDOASIT,
			models.ReServiceCodeDOSHUT, models.ReServiceCodeDDSHUT,
			models.ReServiceCodeDNPK:
			var err error
			weight, err := g.getWeightParams(appCtx, serviceItem)
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

			weightFloat := float64(weight)
			newSegment.L1 = edisegment.L1{
				LadingLineItemNumber: hierarchicalIDNumber,
				FreightRate:          &weightFloat,
				RateValueQualifier:   "LB",
				Charge:               serviceItem.PriceCents.Int64(),
			}

		// following service items have service item dimensions and rate but no distance
		case models.ReServiceCodeDCRT, models.ReServiceCodeDUCRT:
			var err error
			dimensions, rate, err := g.getServiceItemDimensionRateParams(appCtx, serviceItem)
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
				Volume:               dimensions,
				VolumeUnitQualifier:  "E",
				LadingQuantity:       1,
				PackagingFormCode:    "CRT",
			}

			newSegment.L1 = edisegment.L1{
				LadingLineItemNumber: hierarchicalIDNumber,
				FreightRate:          &rate,
				RateValueQualifier:   "PF", // Per Cubic Foot
				Charge:               serviceItem.PriceCents.Int64(),
			}

		default:
			var err error
			weight, distanceFloat, err := g.getWeightAndDistanceParams(appCtx, serviceItem)
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

			weightFloat := float64(weight)
			newSegment.L1 = edisegment.L1{
				LadingLineItemNumber: hierarchicalIDNumber,
				FreightRate:          &weightFloat,
				RateValueQualifier:   "LB",
				Charge:               serviceItem.PriceCents.Int64(),
			}

		}

		fa1, fa2, err := g.createLoaSegments(orders, serviceItem.MTOServiceItem.MTOShipment)
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
