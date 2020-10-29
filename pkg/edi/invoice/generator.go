package ediinvoice

import (
	"bytes"
	"fmt"

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
	Header       []edisegment.Segment `validate:"min=1,dive"`
	ServiceItems []edisegment.Segment `validate:"min=1,dive"`
	SE           edisegment.SE
	GE           edisegment.GE
	IEA          edisegment.IEA
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

	for _, line := range invoice.Header {
		records = append(records, line.StringArray())
	}
	for _, line := range invoice.ServiceItems {
		records = append(records, line.StringArray())
	}
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
