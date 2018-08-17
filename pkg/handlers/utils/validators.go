package utils

import (
	"strings"

	"github.com/gobuffalo/validate"

	"github.com/transcom/mymove/pkg/gen/internalmessages"
)

// CreateFailedValidationPayload Converts the value returned by Pop's ValidateAnd* methods into a payload that can
// be returned to clients. This payload contains an object with a key,  `errors`, the
// value of which is a name -> validation error object.
func CreateFailedValidationPayload(verrs *validate.Errors) *internalmessages.InvalidRequestResponsePayload {
	errs := make(map[string]string)
	for _, key := range verrs.Keys() {
		errs[key] = strings.Join(verrs.Get(key), " ")
	}
	return &internalmessages.InvalidRequestResponsePayload{
		Errors: errs,
	}
}
