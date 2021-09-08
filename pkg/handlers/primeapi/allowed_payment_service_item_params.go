package primeapi

import "github.com/transcom/mymove/pkg/models"

// AllowedParamKeys contains a list of param keys that are permitted in a particular context
type AllowedParamKeys map[models.ReServiceCode][]models.ServiceItemParamName

var (
	// AllowedParamKeysPaymentRequest includes the param keys that we allow to be set by the prime while
	// creating a payment request
	AllowedParamKeysPaymentRequest AllowedParamKeys = map[models.ReServiceCode][]models.ServiceItemParamName{
		models.ReServiceCodeDOASIT: {
			models.ServiceItemParamNameSITPaymentRequestStart,
			models.ServiceItemParamNameSITPaymentRequestEnd,
		},
		models.ReServiceCodeDDASIT: {
			models.ServiceItemParamNameSITPaymentRequestStart,
			models.ServiceItemParamNameSITPaymentRequestEnd,
		},
	}
)

// Contains checks to see if the provided file type is acceptable
func (apk AllowedParamKeys) Contains(serviceCode models.ReServiceCode, paramKeyName string) bool {
	allowedKeys, ok := apk[serviceCode]
	if !ok {
		return false
	}
	for _, key := range allowedKeys {
		if string(key) == paramKeyName {
			return true
		}
	}
	return false
}
