package models

import (
	"fmt"
	"strings"

	"github.com/markbates/validate"
	"github.com/markbates/validate/validators"
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
