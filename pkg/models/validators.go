package models

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
	"github.com/rickar/cal"

	"github.com/gobuffalo/pop"

	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/unit"
)

// StringIsNilOrNotBlank validates OptionalString fields, which we represent as *string.
type StringIsNilOrNotBlank struct {
	Name  string
	Field *string
}

// IsValid adds an error if the pointer is not nil and also an empty string.
func (v *StringIsNilOrNotBlank) IsValid(errors *validate.Errors) {
	if v.Field == nil {
		return
	}
	if strings.TrimSpace(*v.Field) == "" {
		errors.Add(validators.GenerateKey(v.Name), fmt.Sprintf("%s can not be blank.", v.Name))
	}
}

// StringDoesNotContainSSN adds an error if the Field contains an SSN.
type StringDoesNotContainSSN struct {
	Name  string
	Field string
}

var ignoredCharactersRegex = regexp.MustCompile(`(\s|-|\.|_)`)
var nineDigitsRegex = regexp.MustCompile(`^\d{9}$`)

// IsValid adds an error if the Field contains an SSN.
func (v *StringDoesNotContainSSN) IsValid(errors *validate.Errors) {
	cleanSSN := ignoredCharactersRegex.ReplaceAllString(v.Field, "")
	if nineDigitsRegex.MatchString(cleanSSN) {
		errors.Add(validators.GenerateKey(v.Name), fmt.Sprintf("%s Cannot store a raw SSN in this field.", v.Name))
	}
}

// OptionalTimeIsPresent adds an error if the Field is not nil and also not a valid time
type OptionalTimeIsPresent struct {
	Name    string
	Field   *time.Time
	Message string
}

// IsValid adds an error if the Field is not nil and also not a valid time
func (v *OptionalTimeIsPresent) IsValid(errors *validate.Errors) {
	if v.Field != nil {
		timeIsPresent := validators.TimeIsPresent{Name: v.Name, Field: *v.Field, Message: v.Message}
		timeIsPresent.IsValid(errors)
	}
}

// OptionalInt64IsPositive adds an error if the Field is less than or equal to zero
type OptionalInt64IsPositive struct {
	Name  string
	Field *int64
}

// IsValid adds an error if the Field is less than or equal to zero
func (v *OptionalInt64IsPositive) IsValid(errors *validate.Errors) {
	if v.Field != nil {
		if *v.Field <= 0 {
			errors.Add(validators.GenerateKey(v.Name), fmt.Sprintf("%d is less than or equal to zero.", *v.Field))
		}
	}
}

// OptionalIntIsPositive adds an error if the Field is less than or equal to zero
type OptionalIntIsPositive struct {
	Name  string
	Field *int
}

// IsValid adds an error if the Field is less than or equal to zero
func (v *OptionalIntIsPositive) IsValid(errors *validate.Errors) {
	if v.Field != nil {
		if *v.Field <= 0 {
			errors.Add(validators.GenerateKey(v.Name), fmt.Sprintf("%d is less than or equal to zero.", *v.Field))
		}
	}
}

// OptionalPoundIsNonNegative adds an error if the Field is less than zero
type OptionalPoundIsNonNegative struct {
	Name  string
	Field *unit.Pound
}

// IsValid adds an error if the Field is less than zero
func (v *OptionalPoundIsNonNegative) IsValid(errors *validate.Errors) {
	if v.Field != nil {
		if *v.Field < 0 {
			errors.Add(validators.GenerateKey(v.Name), fmt.Sprintf("%d is less than zero.", *v.Field))
		}
	}
}

// Int64IsPresent validates that an int64 is greater than 0.
type Int64IsPresent struct {
	Name  string
	Field int64
}

// IsValid adds an error if the value is equal to 0.
func (v *Int64IsPresent) IsValid(errors *validate.Errors) {
	if v.Field == 0 {
		errors.Add(validators.GenerateKey(v.Name), fmt.Sprintf("%s can not be blank.", v.Name))
	}
}

// DiscountRateIsValid validates that a DiscountRate contains a value between 0 and 1.
type DiscountRateIsValid struct {
	Name  string
	Field unit.DiscountRate
}

// IsValid adds an error if the value is not between 0 and 1.
func (v *DiscountRateIsValid) IsValid(errors *validate.Errors) {
	if v.Field.Float64() < 0 || v.Field.Float64() > 1 {
		errors.Add(validators.GenerateKey(v.Name), fmt.Sprintf("%s must be between 0.0 and 1.0, got %f", v.Name, v.Field))
	}
}

// AllowedFileType validates that a content-type is contained in our list of accepted types.
type AllowedFileType struct {
	validators.StringInclusion
}

// NewAllowedFileTypeValidator constructs as StringInclusion Validator which checks for allowed file upload types
func NewAllowedFileTypeValidator(field string, name string) *AllowedFileType {
	return &AllowedFileType{
		validators.StringInclusion{Name: name,
			Field: field,
			List:  []string{"image/jpeg", "image/png", "application/pdf", "text/plain", "text/plain; charset=utf-8"}}}
}

// AffiliationIsPresent validates that a branch is present
type AffiliationIsPresent struct {
	Name  string
	Field internalmessages.Affiliation
}

// IsValid adds an error if the string value is blank.
func (v *AffiliationIsPresent) IsValid(errors *validate.Errors) {
	if string(v.Field) == "" {
		errors.Add(validators.GenerateKey(v.Name), fmt.Sprintf("%s can not be blank.", v.Name))
	}
}

// OrdersTypeIsPresent validates that orders type field is present
type OrdersTypeIsPresent struct {
	Name  string
	Field internalmessages.OrdersType
}

// IsValid adds an error if the string value is blank.
func (v *OrdersTypeIsPresent) IsValid(errors *validate.Errors) {
	if string(v.Field) == "" {
		errors.Add(validators.GenerateKey(v.Name), fmt.Sprintf("%s can not be blank.", v.Name))
	}
}

// CannotBeTrueIfFalse validates that field1 cannot be true while field2 is false
type CannotBeTrueIfFalse struct {
	Name1  string
	Field1 bool
	Name2  string
	Field2 bool
}

// IsValid adds an error if field1 is true while field2 is false
func (v *CannotBeTrueIfFalse) IsValid(errors *validate.Errors) {
	if v.Field1 == true && v.Field2 == false {
		errors.Add(validators.GenerateKey(v.Name1), fmt.Sprintf("%s can not be true if %s is false", v.Name1, v.Name2))
	}
}

// DateIsWorkday validates that field is on a workday
type DateIsWorkday struct {
	Name     string
	Field    time.Time
	Calendar *cal.Calendar
}

// IsValid adds error if field is not on valid workday
func (v *DateIsWorkday) IsValid(errors *validate.Errors) {
	if !v.Calendar.IsWorkday(v.Field) {
		errors.Add(validators.GenerateKey(v.Name),
			fmt.Sprintf("cannot be on a weekend or holiday, is %v", v.Field))
	}
}

// OptionalDateIsWorkday validates that a field is on a workday if it exists
type OptionalDateIsWorkday struct {
	Name     string
	Field    *time.Time
	Calendar *cal.Calendar
}

// IsValid adds error if field is not on valid workday
// ignores nil field
func (v *OptionalDateIsWorkday) IsValid(errors *validate.Errors) {
	if v.Field == nil {
		return
	}
	dateIsWorkday := DateIsWorkday{v.Name, *v.Field, v.Calendar}
	dateIsWorkday.IsValid(errors)
}

// ValidateableModel is here simply because `validateable` is private to `pop`
type ValidateableModel interface {
	Validate(*pop.Connection) (*validate.Errors, error)
}
