package edisegment

import (
	"time"

	"gopkg.in/go-playground/validator.v9"
)

// HasTimeFormat is a custom validator to verify time format matches expectations.
// Example usage: timeformat=20060102 or timeformat=1504
// See https://golang.org/pkg/time/#Parse for how to interpret formats.
func HasTimeFormat(fl validator.FieldLevel) bool {
	field := fl.Field()
	param := fl.Param()

	_, err := time.Parse(param, field.String())

	return err == nil
}
