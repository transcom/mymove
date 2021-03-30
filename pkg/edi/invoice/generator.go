package ediinvoice

import (
	"bytes"
	"fmt"
	"reflect"

	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/edi"
	edisegment "github.com/transcom/mymove/pkg/edi/segment"
)

// ICNSequenceName used to query Interchange Control Numbers from DB
const ICNSequenceName = "interchange_control_number"

// ICNRandomMin is the smallest allowed random-number based ICN (we use random ICN numbers in development)
const ICNRandomMin int64 = 100000000

// ICNRandomMax is the largest allowed random-number based ICN (we use random ICN numbers in development)
const ICNRandomMax int64 = 999999999

// Invoice858C holds all the segments that are generated
type Invoice858C struct {
	ISA          edisegment.ISA
	GS           edisegment.GS
	ST           edisegment.ST
	Header       InvoiceHeader
	ServiceItems []ServiceItemSegments `validate:"min=1,dive"`
	L3           edisegment.L3
	SE           edisegment.SE
	GE           edisegment.GE
	IEA          edisegment.IEA
}

// InvoiceHeader holds all of the segments that are part of an Invoice858C's Header
type InvoiceHeader struct {
	ShipmentInformation      edisegment.BX
	PaymentRequestNumber     edisegment.N9
	ContractCode             edisegment.N9
	ServiceMemberName        edisegment.N9
	ServiceMemberRank        edisegment.N9
	ServiceMemberBranch      edisegment.N9
	RequestedPickupDate      *edisegment.G62
	ScheduledPickupDate      *edisegment.G62
	ActualPickupDate         *edisegment.G62
	BuyerOrganizationName    edisegment.N1
	SellerOrganizationName   edisegment.N1
	DestinationName          edisegment.N1
	DestinationStreetAddress *edisegment.N3
	DestinationPostalDetails edisegment.N4
	DestinationPhone         *edisegment.PER
	OriginName               edisegment.N1
	OriginStreetAddress      *edisegment.N3
	OriginPostalDetails      edisegment.N4
	OriginPhone              *edisegment.PER
}

// InvoiceResponseHeader holds all the segments used in the headers of the 997, 824 and 810 response types
type InvoiceResponseHeader struct {
	ISA edisegment.ISA
	GS  edisegment.GS
	ST  edisegment.ST
}

// ServiceItemSegmentsSize is the number of fields in the ServiceItemSegments struct
const ServiceItemSegmentsSize int = 7

// ServiceItemSegments holds segments that are required for every service item
type ServiceItemSegments struct {
	HL  edisegment.HL
	N9  edisegment.N9
	L5  edisegment.L5
	L0  edisegment.L0
	L1  edisegment.L1
	FA1 edisegment.FA1
	FA2 edisegment.FA2
}

// NonEmptySegments produces an array of all of the fields
// in an InvoiceHeader that are not nil
func (ih *InvoiceHeader) NonEmptySegments() []edisegment.Segment {
	var result []edisegment.Segment

	// This array should contain every field of InvoiceHeader
	fields := []edisegment.Segment{
		&ih.ShipmentInformation,
		&ih.PaymentRequestNumber,
		&ih.ContractCode,
		&ih.ServiceMemberName,
		&ih.ServiceMemberRank,
		&ih.ServiceMemberBranch,
		ih.RequestedPickupDate,
		ih.ScheduledPickupDate,
		ih.ActualPickupDate,
		&ih.BuyerOrganizationName,
		&ih.SellerOrganizationName,
		&ih.DestinationName,
		ih.DestinationStreetAddress,
		&ih.DestinationPostalDetails,
		ih.DestinationPhone,
		&ih.OriginName,
		ih.OriginStreetAddress,
		&ih.OriginPostalDetails,
		ih.OriginPhone,
	}

	for _, f := range fields {
		// An interface value holding a nil pointer is not nil, so we have to use
		// reflect here instead of just checking f != nil
		if !(reflect.ValueOf(f).Kind() == reflect.Ptr &&
			reflect.ValueOf(f).IsNil()) {
			result = append(result, f)
		}
	}
	return result
}

// Size returns the number of fields in an InvoiceHeader that are not nil
func (ih *InvoiceHeader) Size() int {
	return len(ih.NonEmptySegments())
}

var validate *validator.Validate

func init() {
	validate = validator.New()
}

// Segments returns the invoice as an array of rows (string arrays),
// each containing a segment, to prepare it for writing
func (invoice Invoice858C) Segments() [][]string {
	records := [][]string{
		invoice.ISA.StringArray(),
		invoice.GS.StringArray(),
		invoice.ST.StringArray(),
	}

	for _, line := range invoice.Header.NonEmptySegments() {
		records = append(records, line.StringArray())
	}

	for _, line := range invoice.ServiceItems {
		records = append(records,
			line.HL.StringArray(),
			line.N9.StringArray(),
			line.L5.StringArray(),
			line.L0.StringArray(),
			line.L1.StringArray(),
			line.FA1.StringArray(),
			line.FA2.StringArray(),
		)
	}
	records = append(records, invoice.L3.StringArray())
	records = append(records, invoice.SE.StringArray())
	records = append(records, invoice.GE.StringArray())
	records = append(records, invoice.IEA.StringArray())
	return records
}

func logValidationErrors(logger Logger, err error) {
	// saftey check err is nil just return
	if err == nil {
		return
	}
	if _, ok := err.(*validator.InvalidValidationError); ok {
		logger.Error("InvalidValidationError", zap.Error(err))
		return
	}

	errs := err.(validator.ValidationErrors)
	strErrs := make([]string, len(errs))
	for i, err := range errs {
		strErrs[i] = fmt.Sprintf("%v (value '%s')", err, err.Value())
	}

	logger.Error("ValidationErrors", zap.Strings("errors", strErrs))
}

// EDIString returns the EDI representation of an 858C
func (invoice Invoice858C) EDIString(logger Logger) (string, error) {
	err := invoice.Validate()
	if err != nil {
		// Log validation details, but do not expose details via API
		logValidationErrors(logger, err)
		return "", fmt.Errorf("EDI failed validation: %w", err)
	}

	var b bytes.Buffer
	ediWriter := edi.NewWriter(&b)
	err = ediWriter.WriteAll(invoice.Segments())
	if err != nil {
		return "", fmt.Errorf("EDI failed write: %w", err)
	}
	return b.String(), err
}

// Validate will validate the invoice struct (and nested structs) to make sure they will produce legal EDI.
// This returns either an InvalidValidationError or a validator.ValidationErrors that allows all validation
// errors to be introspected individually.
func (invoice Invoice858C) Validate() error {
	return validate.Struct(invoice)
}
