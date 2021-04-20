package models

import (
	"fmt"
	"strings"
	"time"

	"github.com/transcom/mymove/pkg/gen/primemessages"

	"github.com/gofrs/uuid"

	"github.com/gobuffalo/validate/v3"
	"github.com/gobuffalo/validate/v3/validators"
	"github.com/rickar/cal"

	"github.com/gobuffalo/pop/v5"

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

// OptionalDateNotBefore validates that a date is not before the earliest allowable date
type OptionalDateNotBefore struct {
	Name    string
	Field   *time.Time
	MinDate *time.Time
}

// IsValid adds an error if the field has a value and there is not a not-before date or the date is before the not-before date
func (v *OptionalDateNotBefore) IsValid(errors *validate.Errors) {
	if v.Field != nil {
		if v.MinDate == nil {
			errors.Add(validators.GenerateKey(v.Name), fmt.Sprintf("cannot create this date without a no-earlier-than date"))
		} else if (*v.Field).Before(*v.MinDate) {
			errors.Add(validators.GenerateKey(v.Name), fmt.Sprintf("%s must be on or after %s", *v.Field, *v.MinDate))
		}
	}
}

type container interface {
	Contains(string) bool
	Contents() []string
}

// StringInList is an improved validators.StringInclusion validator with better error messages.
type StringInList struct {
	Value     string
	FieldName string
	List      container
}

// NewStringInList returns a new StringInList validator.
func NewStringInList(value string, fieldName string, list container) *StringInList {
	return &StringInList{
		Value:     value,
		FieldName: fieldName,
		List:      list,
	}
}

// IsValid adds an error if the string value is blank.
func (v *StringInList) IsValid(errors *validate.Errors) {
	if !v.List.Contains(v.Value) {
		errors.Add(validators.GenerateKey(v.FieldName), fmt.Sprintf("'%s' is not in the list [%s].", v.Value, strings.Join(v.List.Contents(), ", ")))
	}
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
	if v.Field1 && !v.Field2 {
		errors.Add(validators.GenerateKey(v.Name1), fmt.Sprintf("%s can not be true if %s is false", v.Name1, v.Name2))
	}
}

// MustBeBothNilOrBothHaveValue validates that two fields are either both nil or both have values
type MustBeBothNilOrBothHaveValue struct {
	FieldName1  string
	FieldValue1 *string
	FieldName2  string
	FieldValue2 *string
}

// IsValid adds an error if fieldValue1 or fieldValue2 are neither both empty nor both non-empty
func (v *MustBeBothNilOrBothHaveValue) IsValid(errors *validate.Errors) {
	if (v.FieldValue1 == nil && v.FieldValue2 != nil) || (v.FieldValue1 != nil && v.FieldValue2 == nil) {
		errors.Add(validators.GenerateKey(v.FieldName1), fmt.Sprintf("%s can not be nil if %s has a value and vice versa", v.FieldName1, v.FieldName2))
	}
}

// AtLeastOneNotNil validates that at least one of two fields are not nil
type AtLeastOneNotNil struct {
	FieldName1  string
	FieldValue1 *string
	FieldName2  string
	FieldValue2 *string
}

// IsValid adds an error if fieldValue1 and fieldValue2 are nil
func (v *AtLeastOneNotNil) IsValid(errors *validate.Errors) {
	if v.FieldValue1 == nil && v.FieldValue2 == nil {
		errors.Add(validators.GenerateKey(v.FieldName1), fmt.Sprintf("Both %s and %s cannot be nil, one must be valid", v.FieldName1, v.FieldName2))
		errors.Add(validators.GenerateKey(v.FieldName2), fmt.Sprintf("Both %s and %s cannot be nil, one must be valid", v.FieldName2, v.FieldName1))
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

// OptionalStringInclusion validates that a field is in a list of strings if the field exists
type OptionalStringInclusion struct {
	Name    string
	Field   *string
	List    []string
	Message string
}

// IsValid adds error if field is non-nil and not in the list of strings
func (v *OptionalStringInclusion) IsValid(errors *validate.Errors) {
	if v.Field == nil {
		return
	}
	stringInclusion := validators.StringInclusion{
		Name:    v.Name,
		Field:   *v.Field,
		List:    v.List,
		Message: v.Message,
	}
	stringInclusion.IsValid(errors)
}

// OptionalRegexMatch validates that a field matches the regexp match
type OptionalRegexMatch struct {
	Name    string
	Field   *string
	Expr    string
	Message string
}

// IsValid performs the validation based on the regexp match
func (v *OptionalRegexMatch) IsValid(errors *validate.Errors) {
	if v.Field == nil {
		return
	}
	regexMatch := validators.RegexMatch{
		Name:    v.Name,
		Field:   *v.Field,
		Expr:    v.Expr,
		Message: v.Message,
	}
	regexMatch.IsValid(errors)
}

// Float64IsPresent validates that a float64 is non-zero.
type Float64IsPresent struct {
	Name    string
	Field   float64
	Message string
}

// IsValid adds an error if the field equals 0.
func (v *Float64IsPresent) IsValid(errors *validate.Errors) {
	if v.Field != 0 {
		return
	}

	if len(v.Message) > 0 {
		errors.Add(validators.GenerateKey(v.Name), v.Message)
		return
	}

	errors.Add(validators.GenerateKey(v.Name), fmt.Sprintf("%s can not be blank.", v.Name))
}

// Float64IsGreaterThan validates that a float64 is greater than a given value
type Float64IsGreaterThan struct {
	Name     string
	Field    float64
	Compared float64
	Message  string
}

// IsValid adds an error if the field is not greater than the compared value.
func (v *Float64IsGreaterThan) IsValid(errors *validate.Errors) {
	if v.Field > v.Compared {
		return
	}

	if len(v.Message) > 0 {
		errors.Add(validators.GenerateKey(v.Name), v.Message)
		return
	}

	errors.Add(validators.GenerateKey(v.Name), fmt.Sprintf("%f is not greater than %f.", v.Field, v.Compared))
}

// OptionalUUIDIsPresent is a structure for determining if an Optional UUID is valid
// If it is a nil pointer, it passes validation.
// If it is a pointer to a valid UUID, it passes validation.
// If it is a pointer to a non-valid UUID, it fails validation.
type OptionalUUIDIsPresent struct {
	Name    string
	Field   *uuid.UUID
	Message string
}

// IsValid adds an error if an optional UUID is valid.
// If it is a nil pointer, it passes validation.
// If it is a pointer to a valid UUID, it passes validation.
// If it is a pointer to a non-valid UUID, it fails validation.
func (v *OptionalUUIDIsPresent) IsValid(errors *validate.Errors) {
	if v.Field == nil {
		return
	}

	s := v.Field.String()
	if strings.TrimSpace(s) != "" && *v.Field != uuid.Nil {
		return
	}

	if len(v.Message) > 0 {
		errors.Add(validators.GenerateKey(v.Name), v.Message)
		return
	}

	errors.Add(validators.GenerateKey(v.Name), fmt.Sprintf("%s can not be blank.", v.Name))
}

// ItemCanFitInsideCrate is a structure for determining if an Item Dimension can fit inside a Crate Dimension
type ItemCanFitInsideCrate struct {
	Name         string
	NameCompared string
	Item         *primemessages.MTOServiceItemDimension
	Crate        *primemessages.MTOServiceItemDimension
	Message      string
}

// IsValid adds an error if the Item can not fit inside a Crate
func (v ItemCanFitInsideCrate) IsValid(errors *validate.Errors) {
	if v.Item == nil || v.Crate == nil {
		errors.Add(validators.GenerateKey(v.Name), fmt.Sprintf("%s or %s can not be nil.", v.Name, v.NameCompared))
		return
	}

	if v.Item.Height == nil || v.Item.Length == nil || v.Item.Width == nil ||
		v.Crate.Height == nil || v.Crate.Length == nil || v.Crate.Width == nil {
		errors.Add(validators.GenerateKey(v.Name), fmt.Sprintf("%s or %s has missing dimensions.", v.Name, v.NameCompared))
		return
	}

	if *v.Item.Length < *v.Crate.Length && *v.Item.Width < *v.Crate.Width && *v.Item.Height < *v.Crate.Height {
		return
	}

	if len(v.Message) > 0 {
		errors.Add(validators.GenerateKey(v.Name), v.Message)
		return
	}

	errors.Add(validators.GenerateKey(v.Name), fmt.Sprintf("%s dimensions can not be greater than or equal to %s dimensions.", v.Name, v.NameCompared))
}

// ValidateableModel is here simply because `validateable` is private to `pop`
type ValidateableModel interface {
	Validate(*pop.Connection) (*validate.Errors, error)
}
