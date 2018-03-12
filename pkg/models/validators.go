package models

import (
	"fmt"
	"strings"

	"github.com/markbates/validate"
	"github.com/markbates/validate/validators"
)

type OptionalStringIsPresent struct {
	Name  string
	Field *string
}

func (v *OptionalStringIsPresent) IsValid(errors *validate.Errors) {
	if v.Field == nil {
		errors.Add(validators.GenerateKey(v.Name), fmt.Sprintf("%s can not be blank.", v.Name))
		return
	}
	if strings.TrimSpace(*v.Field) == "" {
		errors.Add(validators.GenerateKey(v.Name), fmt.Sprintf("%s can not be blank.", v.Name))
	}
}
