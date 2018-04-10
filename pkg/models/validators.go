package models

import (
	"fmt"
	"regexp"
	"sort"
	"strings"

	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"

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

// AllowedFiletype validates that a content-type is contained in our list of accepted
// types.
type AllowedFiletype struct {
	Name  string
	Field string
}

// AllowedFiletypes are the types of files that are accepted for upload.
var AllowedFiletypes = map[string]string{
	"JPG": "image/jpeg",
	"PNG": "image/png",
	"PDF": "application/pdf",
}

// IsValid adds an error if the value is equal to 0.
func (v *AllowedFiletype) IsValid(errors *validate.Errors) {
	for _, filetype := range AllowedFiletypes {
		if filetype == v.Field {
			return
		}
	}

	filetypes := []string{}
	for name := range AllowedFiletypes {
		filetypes = append(filetypes, name)
	}
	sort.Strings(filetypes)
	list := strings.Join(filetypes, ", ")
	errors.Add(validators.GenerateKey(v.Name), fmt.Sprintf("%s must be one of: %s.", v.Name, list))
}

// BranchIsPresent validates that a branch is present
type BranchIsPresent struct {
	Name  string
	Field internalmessages.MilitaryBranch
}

// IsValid adds an error if the string value is blank.
func (v *BranchIsPresent) IsValid(errors *validate.Errors) {
	if string(v.Field) == "" {
		errors.Add(strings.ToLower(string(v.Field)), fmt.Sprintf("%s must not be blank!", v.Name))
	}
}
