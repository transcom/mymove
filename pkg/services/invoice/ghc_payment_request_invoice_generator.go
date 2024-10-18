package invoice

import (
	"database/sql"
	"fmt"
	"regexp"
	"strconv"
	"strings"

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
	services.LineOfAccountingFetcher
}

// NewGHCPaymentRequestInvoiceGenerator returns an implementation of the GHCPaymentRequestInvoiceGenerator interface
func NewGHCPaymentRequestInvoiceGenerator(icnSequencer sequence.Sequencer, clock clock.Clock, linesOfAccountingFetcher services.LineOfAccountingFetcher) services.GHCPaymentRequestInvoiceGenerator {
	return &ghcPaymentRequestInvoiceGenerator{
		icnSequencer:            icnSequencer,
		clock:                   clock,
		LineOfAccountingFetcher: linesOfAccountingFetcher,
	}
}

const dateFormat = "20060102"
const isaDateFormat = "060102"
const timeFormat = "1504"
const maxCityLength = 30
const maxLocationlength = 60
const maxServiceMemberNameLengthN9 = 30

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
		return ediinvoice.Invoice858C{}, fmt.Errorf("failed to get next Interchange Control Number: %w", err)
	}

	// save ICN
	pr2icn := models.PaymentRequestToInterchangeControlNumber{
		PaymentRequestID:         paymentRequest.ID,
		InterchangeControlNumber: int(interchangeControlNumber),
		EDIType:                  models.EDIType858,
	}
	verrs, err := appCtx.DB().ValidateAndSave(&pr2icn)
	if err != nil {
		return ediinvoice.Invoice858C{}, fmt.Errorf("failed to save Interchange Control Number: %w", err)
	} else if verrs != nil && verrs.HasAny() {
		return ediinvoice.Invoice858C{}, fmt.Errorf("failed to save Interchange Control Number: %s", verrs.String())
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
		Date:                     paymentRequest.RequestedAt.Format(dateFormat),
		Time:                     paymentRequest.RequestedAt.Format(timeFormat),
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
		ShipmentIdentificationNumber: *moveTaskOrder.ReferenceID,
		StandardCarrierAlphaCode:     "HSFR",
		ShipmentQualifier:            "4",
	}

	edi858.Header.PaymentRequestNumber = edisegment.N9{
		ReferenceIdentificationQualifier: "CN",
		ReferenceIdentification:          paymentRequest.PaymentRequestNumber,
	}

	// Add moveCode to header
	edi858.Header.MoveCode = edisegment.N9{
		ReferenceIdentificationQualifier: "CMN",
		ReferenceIdentification:          moveTaskOrder.Locator,
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

	// Add order pay grade detail to header
	if moveTaskOrder.Orders.Grade == nil {
		// Nil check
		return ediinvoice.Invoice858C{}, apperror.NewNotFoundError(moveTaskOrder.Orders.ID, "order pay grade not found")
	}

	edi858.Header.OrderPayGrade = edisegment.N9{
		ReferenceIdentificationQualifier: "ML",
		ReferenceIdentification:          string(*moveTaskOrder.Orders.Grade),
	}

	var paymentServiceItems models.PaymentServiceItems
	err = appCtx.DB().Q().
		Eager("MTOServiceItem.ReService", "MTOServiceItem.MTOShipment").
		Where("payment_request_id = ?", paymentRequest.ID).
		Where("status = ?", models.PaymentServiceItemStatusApproved).
		Order("created_at").
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

	if moveTaskOrder.Orders.OriginDutyLocationGBLOC == nil {
		return ediinvoice.Invoice858C{}, apperror.NewInvalidInputError(moveTaskOrder.OrdersID, fmt.Errorf("origin duty location GBLOC value is missing"), nil, "origin duty location GBLOC is required")
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
	var fa2segments []edisegment.FA2
	for _, serviceItem := range edi858.ServiceItems {
		fa2segments = append(fa2segments, serviceItem.FA2s...)
	}

	serviceItemSegmentCount := len(edi858.ServiceItems)*ediinvoice.ServiceItemSegmentsSizeWithoutFA2s + len(fa2segments)
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
		ReferenceIdentification:          truncateStr(serviceMember.ReverseNameLineFormat(), maxServiceMemberNameLengthN9),
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

	// dod id or emplid
	if branch.String() == models.AffiliationCOASTGUARD.String() {
		emplid := serviceMember.Emplid
		if emplid == nil {
			return apperror.NewConflictError(serviceMember.ID, fmt.Sprintf("no employee id found for ServiceMember ID: %s Payment Request ID: %s", serviceMember.ID, paymentRequestID))
		}
		header.ServiceMemberID = edisegment.N9{
			ReferenceIdentificationQualifier: "4A",
			ReferenceIdentification:          string(*emplid),
		}
	} else {
		dodID := serviceMember.Edipi
		if dodID == nil {
			return apperror.NewConflictError(serviceMember.ID, fmt.Sprintf("no dod id found for ServiceMember ID: %s Payment Request ID: %s", serviceMember.ID, paymentRequestID))
		}
		header.ServiceMemberID = edisegment.N9{
			ReferenceIdentificationQualifier: "4A",
			ReferenceIdentification:          string(*dodID),
		}
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
		originDutyLocation, err = models.FetchDutyLocationWithTransportationOffice(appCtx.DB(), *orders.OriginDutyLocationID)
		if err != nil {
			return apperror.NewInvalidInputError(*orders.OriginDutyLocationID, err, nil, "unable to find origin duty location")
		}
	} else {
		return apperror.NewConflictError(orders.ID, "Invalid Order, must have OriginDutyLocation")
	}

	var address models.Address
	err = appCtx.DB().Q().
		Select("addresses.*").
		Join("mto_shipments", "addresses.id = mto_shipments.pickup_address_id").
		Join("moves", "mto_shipments.move_id = moves.id").
		Join("mto_service_items", "mto_service_items.move_id = moves.id").
		Join("payment_service_items", "payment_service_items.mto_service_item_id = mto_service_items.id").
		Where("payment_service_items.payment_request_id = ?", paymentRequestID).
		Order("mto_shipments.created_at").
		First(&address)

	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return apperror.NewNotFoundError(paymentRequestID, "for mto shipments associated with PaymentRequest")
		default:
			return apperror.NewQueryError("MTOShipments", err, fmt.Sprintf("error querying for shipments pickup address gbloc to use in N1*BY segments in PaymentRequest %s: %s", paymentRequestID, err))
		}
	}
	pickupPostalCodeToGbloc, gblocErr := models.FetchGBLOCForPostalCode(appCtx.DB(), address.PostalCode)
	if gblocErr != nil {
		return apperror.NewInvalidInputError(pickupPostalCodeToGbloc.ID, gblocErr, nil, "unable to determine GBLOC for pickup postal code")
	}

	header.BuyerOrganizationName = edisegment.N1{
		EntityIdentifierCode:        "BY",
		Name:                        truncateStr(originDutyLocation.Name, maxLocationlength),
		IdentificationCodeQualifier: "92",
		IdentificationCode:          modifyGblocIfMarines(*orders.ServiceMember.Affiliation, pickupPostalCodeToGbloc.GBLOC),
	}

	// seller organization name
	header.SellerOrganizationName = edisegment.N1{
		EntityIdentifierCode:        "SE",
		Name:                        "Prime",
		IdentificationCodeQualifier: "2",
		IdentificationCode:          "HSFR",
	}

	return nil
}

func (g ghcPaymentRequestInvoiceGenerator) createOriginAndDestinationSegments(appCtx appcontext.AppContext, _ uuid.UUID, orders models.Order, header *ediinvoice.InvoiceHeader) error {
	var err error
	var destinationDutyLocation models.DutyLocation
	if orders.NewDutyLocationID != uuid.Nil {
		destinationDutyLocation, err = models.FetchDutyLocationWithTransportationOffice(appCtx.DB(), orders.NewDutyLocationID)
		if err != nil {
			return apperror.NewInvalidInputError(orders.NewDutyLocationID, err, nil, "unable to find new duty location")
		}
	} else {
		return apperror.NewConflictError(orders.ID, "Invalid Order, must have NewDutyLocation")
	}

	destPostalCodeToGbloc, gblocErr := models.FetchGBLOCForPostalCode(appCtx.DB(), destinationDutyLocation.Address.PostalCode)
	if gblocErr != nil {
		return apperror.NewInvalidInputError(destinationDutyLocation.ID, gblocErr, nil, "unable to determine GBLOC for duty location postal code")
	}

	// destination name
	header.DestinationName = edisegment.N1{
		EntityIdentifierCode:        "ST",
		Name:                        truncateStr(destinationDutyLocation.Name, maxLocationlength),
		IdentificationCodeQualifier: "10",
		IdentificationCode:          modifyGblocIfMarines(*orders.ServiceMember.Affiliation, destPostalCodeToGbloc.GBLOC),
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
	destPhoneLines := determineDutyLocationPhoneLines(destinationDutyLocation)

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
		originDutyLocation, err = models.FetchDutyLocationWithTransportationOffice(appCtx.DB(), *orders.OriginDutyLocationID)
		if err != nil {
			return apperror.NewInvalidInputError(*orders.OriginDutyLocationID, err, nil, "unable to find origin duty location")
		}
	} else {
		return apperror.NewConflictError(orders.ID, "Invalid Order, must have OriginDutyLocation")
	}

	header.OriginName = edisegment.N1{
		EntityIdentifierCode:        "SF",
		Name:                        truncateStr(originDutyLocation.Name, maxLocationlength),
		IdentificationCodeQualifier: "10",
		IdentificationCode:          modifyGblocIfMarines(*orders.ServiceMember.Affiliation, *orders.OriginDutyLocationGBLOC),
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
	originPhoneLines := determineDutyLocationPhoneLines(originDutyLocation)

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

func (g ghcPaymentRequestInvoiceGenerator) createLoaSegments(appCtx appcontext.AppContext, orders models.Order, shipment models.MTOShipment) (edisegment.FA1, []edisegment.FA2, error) {
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
			return edisegment.FA1{}, nil, apperror.NewConflictError(orders.ID, "Invalid order. Must have an HHG TAC value")
		}
		tac = *orders.TAC
	} else {
		if orders.NtsTAC == nil || *orders.NtsTAC == "" {
			return edisegment.FA1{}, nil, apperror.NewConflictError(orders.ID, "Invalid order. Must have an NTS TAC value")
		}
		tac = *orders.NtsTAC
	}

	// Get SAC or SDN from orders. Use the HHG one by default (blank or SACType HHG), use the NTS if SACType is NTS.
	useHHGSac := true
	if shipment.ID != uuid.Nil {
		// We do have a shipment, so see if the shipment prefers the NTS SAC or SDN.
		if shipment.SACType != nil && *shipment.SACType == models.LOATypeNTS {
			useHHGSac = false
		}
	}

	// Now grab the preferred SAC or SDN, making sure that it actually exists on the orders record.
	var sac string
	if useHHGSac {
		if orders.SAC != nil && *orders.SAC != "" {
			sac = *orders.SAC
			if len(sac) > 80 {
				return edisegment.FA1{}, nil, apperror.NewConflictError(orders.ID, "Invalid order. HHG SAC/SDN must be 80 characters or less")
			}
		}

	} else {
		if orders.NtsSAC != nil && *orders.NtsSAC != "" {
			sac = *orders.NtsSAC
			if len(sac) > 80 {
				return edisegment.FA1{}, nil, apperror.NewConflictError(orders.ID, "Invalid order. NTS SAC/SDN must be 80 characters or less")
			}
		}
	}

	affiliation := models.ServiceMemberAffiliation(*orders.ServiceMember.Affiliation)
	agencyQualifierCode, found := edisegment.AffiliationToAgency[affiliation]

	if !found {
		agencyQualifierCode = "DF"
	}

	fa1 := edisegment.FA1{
		AgencyQualifierCode: agencyQualifierCode,
	}

	// May have multiple FA2 segments: TAC, SAC (optional), and many long LOA values
	var fa2s []edisegment.FA2

	// TAC
	fa2TAC := edisegment.FA2{
		BreakdownStructureDetailCode: edisegment.FA2DetailCodeTA,
		FinancialInformationCode:     tac,
	}
	fa2s = append(fa2s, fa2TAC)

	// SAC (optional)
	if sac != "" {
		fa2SAC := edisegment.FA2{
			BreakdownStructureDetailCode: edisegment.FA2DetailCodeZZ,
			FinancialInformationCode:     sac,
		}
		fa2s = append(fa2s, fa2SAC)
	}

	fa2LongLoaSegments, err := g.createLongLoaSegments(appCtx, orders, tac)
	if err != nil {
		return edisegment.FA1{}, nil, err
	}
	fa2s = append(fa2s, fa2LongLoaSegments...)

	return fa1, fa2s, nil
}

func (g ghcPaymentRequestInvoiceGenerator) createLongLoaSegments(appCtx appcontext.AppContext, orders models.Order, tac string) ([]edisegment.FA2, error) {
	var loas []models.LineOfAccounting
	var loa models.LineOfAccounting

	// Nil check on orders department indicator
	if orders.DepartmentIndicator == nil {
		return nil, apperror.NewQueryError("orders", fmt.Errorf("could not identify department indicator for Order ID %s", orders.ID), "Unexpected error")
	}

	// Fetch the long lines of accounting for an invoice based off an order's department indicator, tacCode, and the orders issue date.
	// There is special logic for whether or not the department indicator is for the US Coast Guard.
	loas, err := g.LineOfAccountingFetcher.FetchLongLinesOfAccounting(models.DepartmentIndicator(*orders.DepartmentIndicator), orders.IssueDate, tac, appCtx)
	if err != nil {
		return nil, apperror.NewQueryError("lineOfAccounting", err, "Unexpected error")
	}
	if len(loas) == 0 {
		return nil, nil
	}
	// pick first one (sorted by FBMC, loa_bgn_dt, tac_fy_txt) inside the service object
	loa = loas[0]

	if models.DepartmentIndicator(*orders.DepartmentIndicator) != models.DepartmentIndicatorCOASTGUARD {

		//"HE" - E-1 through E-9 and Special Enlisted
		//"HO" - O-1 Academy graduate through O-10, W1 - W5, Aviation Cadet, Academy Cadet, and Midshipman
		//"HC" - Civilian employee

		if orders.Grade == nil {
			return nil, apperror.NewConflictError(orders.ServiceMember.ID, "this service member has no pay grade for the specified order")
		}
		grade := *orders.Grade

		hhgCode := ""
		if grade[:2] == "E_" {
			hhgCode = "HE"
		} else if grade[:2] == "O_" || grade[:2] == "W_" || grade == models.ServiceMemberGradeACADEMYCADET || grade == models.ServiceMemberGradeAVIATIONCADET || grade == models.ServiceMemberGradeMIDSHIPMAN {
			hhgCode = "HO"
		} else if grade == models.ServiceMemberGradeCIVILIANEMPLOYEE {
			hhgCode = "HC"
		} else {
			return nil, apperror.NotImplementedError{}
		}
		// if just one, pick it
		// if multiple,lowest FBMC
		var loaWithMatchingCode []models.LineOfAccounting

		for _, line := range loas {
			if line.LoaHsGdsCd != nil && *line.LoaHsGdsCd == hhgCode {
				loaWithMatchingCode = append(loaWithMatchingCode, line)
			}
		}
		if len(loaWithMatchingCode) == 0 {
			// fall back to the whole set and then sort by fbmc
			// take first thing from whole set
			loa = loas[0]
		}
		if len(loaWithMatchingCode) >= 1 {
			// take first of loaWithMatchingCode
			loa = loaWithMatchingCode[0]
		}
	}
	var fa2LongLoaSegments []edisegment.FA2

	var concatDate *string
	if loa.LoaBgFyTx != nil && loa.LoaEndFyTx != nil {
		fiscalYearStr := fmt.Sprintf("%d%d", *loa.LoaBgFyTx, *loa.LoaEndFyTx)
		concatDate = &fiscalYearStr
	} else {
		blankValue := "XXXXXXXX"
		concatDate = &blankValue
	}

	// The FA2 L1 segment must be exactly six characters in length. Our imported database values from TRDM are numeric
	// strings and so we need to left pad with zeros to meet the threshold.  This may not be needed when the real TGET
	// integration is introduced.
	var accountingInstallationNumber *string
	if loa.LoaInstlAcntgActID != nil {
		zeroPaddedInstlAcntgActID := fmt.Sprintf("%06s", *loa.LoaInstlAcntgActID)
		accountingInstallationNumber = &zeroPaddedInstlAcntgActID
	}

	// Create long LOA FA2 segments
	segmentInputs := []struct {
		detailCode edisegment.FA2DetailCode
		infoCode   *string
	}{
		// If order of these changes, tests will also need to be adjusted. Using alpha order by detailCode.
		{edisegment.FA2DetailCodeA1, loa.LoaDptID},
		{edisegment.FA2DetailCodeA2, loa.LoaTnsfrDptNm},
		{edisegment.FA2DetailCodeA3, concatDate},
		{edisegment.FA2DetailCodeA4, loa.LoaBafID},
		{edisegment.FA2DetailCodeA5, loa.LoaTrsySfxTx},
		{edisegment.FA2DetailCodeA6, loa.LoaMajClmNm},
		{edisegment.FA2DetailCodeB1, loa.LoaOpAgncyID},
		{edisegment.FA2DetailCodeB2, loa.LoaAlltSnID},
		{edisegment.FA2DetailCodeB3, loa.LoaUic},
		{edisegment.FA2DetailCodeC1, loa.LoaPgmElmntID},
		{edisegment.FA2DetailCodeC2, loa.LoaTskBdgtSblnTx},
		{edisegment.FA2DetailCodeD1, loa.LoaDfAgncyAlctnRcpntID},
		{edisegment.FA2DetailCodeD4, loa.LoaJbOrdNm},
		{edisegment.FA2DetailCodeD6, loa.LoaSbaltmtRcpntID},
		{edisegment.FA2DetailCodeD7, loa.LoaWkCntrRcpntNm},
		{edisegment.FA2DetailCodeE1, loa.LoaMajRmbsmtSrcID},
		{edisegment.FA2DetailCodeE2, loa.LoaDtlRmbsmtSrcID},
		{edisegment.FA2DetailCodeE3, loa.LoaCustNm},
		{edisegment.FA2DetailCodeF1, loa.LoaObjClsID},
		{edisegment.FA2DetailCodeF3, loa.LoaSrvSrcID},
		{edisegment.FA2DetailCodeG2, loa.LoaSpclIntrID},
		{edisegment.FA2DetailCodeI1, loa.LoaBdgtAcntClsNm},
		{edisegment.FA2DetailCodeJ1, loa.LoaDocID},
		{edisegment.FA2DetailCodeK6, loa.LoaClsRefID},
		{edisegment.FA2DetailCodeL1, accountingInstallationNumber},
		{edisegment.FA2DetailCodeM1, loa.LoaLclInstlID},
		{edisegment.FA2DetailCodeN1, loa.LoaTrnsnID},
		{edisegment.FA2DetailCodeP5, loa.LoaFmsTrnsactnID},
	}

	for _, input := range segmentInputs {
		fa2, loaErr := createLongLoaSegment(input.detailCode, input.infoCode)
		if loaErr != nil {
			return nil, loaErr
		}
		if fa2 != nil {
			fa2LongLoaSegments = append(fa2LongLoaSegments, *fa2)
		}
	}

	return fa2LongLoaSegments, nil
}

func createLongLoaSegment(detailCode edisegment.FA2DetailCode, infoCode *string) (*edisegment.FA2, error) {
	// If we don't have an infoCode value, then just ignore this segment
	if infoCode == nil || strings.TrimSpace(*infoCode) == "" {
		return nil, nil
	}
	value := *infoCode

	// Make sure we have a detailCode
	if len(detailCode) != 2 {
		return nil, apperror.NewImplementationError("Detail code should have length 2")
	}

	// The FinancialInformationCode field is limited to 80 characters, so make sure the value doesn't exceed
	// that (given our LOA field schema types, it shouldn't unless we've made a mistake somewhere).
	if len(value) > 80 {
		return nil, apperror.NewImplementationError(fmt.Sprintf("Value for FA2 code %s exceeds 80 character limit", detailCode))
	}

	return &edisegment.FA2{
		BreakdownStructureDetailCode: detailCode,
		FinancialInformationCode:     value,
	}, nil
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
		return 0, fmt.Errorf("could not parse weight for PaymentServiceItem %s: %w", serviceItem.ID, err)
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
		return 0, 0, fmt.Errorf("could not parse cubic feet as a float for PaymentServiceItem %s: %w", serviceItem.ID, err)
	}

	rate, err := g.fetchPaymentServiceItemParam(appCtx, serviceItem.ID, models.ServiceItemParamNamePriceRateOrFactor)
	if err != nil {
		return 0, 0, err
	}
	rateFloat, err := strconv.ParseFloat(rate.Value, 64)
	if err != nil {
		return 0, 0, fmt.Errorf("could not parse rate as a float for PaymentServiceItem %s: %w", serviceItem.ID, err)
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
	case models.ReServiceCodeDDDSIT, models.ReServiceCodeDDSFSC:
		distanceModel = models.ServiceItemParamNameDistanceZipSITDest
	case models.ReServiceCodeDOPSIT, models.ReServiceCodeDOSFSC:
		distanceModel = models.ServiceItemParamNameDistanceZipSITOrigin
	default:
		distanceModel = models.ServiceItemParamNameDistanceZip
	}

	distance, err := g.fetchPaymentServiceItemParam(appCtx, serviceItem.ID, distanceModel)
	if err != nil {
		return 0, 0, err
	}
	distanceFloat, err := strconv.ParseFloat(distance.Value, 64)
	if err != nil {
		return 0, 0, fmt.Errorf("could not parse Distance Zip3 for PaymentServiceItem %s: %w", serviceItem.ID, err)
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
				FreightRate:          nil,
				RateValueQualifier:   "",
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

		fa1, fa2s, err := g.createLoaSegments(appCtx, orders, serviceItem.MTOServiceItem.MTOShipment)
		if err != nil {
			return segments, l3, err
		}
		newSegment.FA1 = fa1
		newSegment.FA2s = fa2s
		segments = append(segments, newSegment)
	}

	return segments, l3, nil
}

// This business logic should likely live in the transportation_office.go file,
// however, since the change would likely impact other parts of the application it is here so that it only
// updates the Gbloc sent to Syncada
func modifyGblocIfMarines(affiliation models.ServiceMemberAffiliation, gbloc string) string {
	if affiliation == models.AffiliationMARINES {
		gbloc = "USMC"
	}
	return gbloc
}

// determineDutyLocationPhoneLines returns a slice of strings of the phone numbers of all voice type phones lines for
// the associated Transportation Office
func determineDutyLocationPhoneLines(dutyLocation models.DutyLocation) (phoneLines []string) {
	if dutyLocation.TransportationOfficeID == nil {
		return phoneLines
	}

	dutyLocationPhoneLines := dutyLocation.TransportationOffice.PhoneLines
	for _, phoneLine := range dutyLocationPhoneLines {
		if phoneLine.Type == "voice" {
			phoneLines = append(phoneLines, phoneLine.Number)
		}
	}
	return phoneLines
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
