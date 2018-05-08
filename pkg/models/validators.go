package models

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"

	"github.com/gobuffalo/pop"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
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

// AllowedFileType validates that a content-type is contained in our list of accepted types.
type AllowedFileType struct {
	validators.StringInclusion
}

// NewAllowedFileTypeValidator constructs as StringInclusion Validator which checks for allowed file upload types
func NewAllowedFileTypeValidator(field string, name string) *AllowedFileType {
	return &AllowedFileType{
		validators.StringInclusion{Name: name,
			Field: field,
			List:  []string{"image/jpeg", "image/png", "application/pdf"}}}
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

// BackupContactPermissionIsPresent validates that permission field is present
type BackupContactPermissionIsPresent struct {
	Name  string
	Field internalmessages.BackupContactPermission
}

// IsValid adds an error if the string value is blank.
func (v *BackupContactPermissionIsPresent) IsValid(errors *validate.Errors) {
	if string(v.Field) == "" {
		errors.Add(validators.GenerateKey(v.Name), fmt.Sprintf("%s can not be blank.", v.Name))
	}
}

// OrdersTypeIsPresent validates that permission field is present
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

// ValidateableModel is here simply because `validateable` is private to `pop`
type ValidateableModel interface {
	Validate(*pop.Connection) (*validate.Errors, error)
}

// FieldValidator is used to chain validations when a field of a is, itself, a model
type FieldValidator struct {
	connection *pop.Connection
	field      ValidateableModel
	name       string
	Error      error
}

// NewFieldValidator constructs and
func NewFieldValidator(c *pop.Connection, f ValidateableModel, name string) *FieldValidator {
	return &FieldValidator{c, f, name, nil}
}

// IsValid adds the field(model)'s validation errors to the Errors for the parent. Also sets v.Error appropriately
func (v *FieldValidator) IsValid(errors *validate.Errors) {
	var localErrors *validate.Errors
	key := validators.GenerateKey(v.name)
	localErrors, v.Error = v.field.Validate(v.connection)
	for _, verr := range localErrors.Errors {
		for _, msg := range verr {
			errors.Add(key, fmt.Sprintf("%s.%s", v.name, msg))
		}
	}
}
