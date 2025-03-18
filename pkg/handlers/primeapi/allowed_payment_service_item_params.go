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
			models.ServiceItemParamNameWeightBilled,
		},
		models.ReServiceCodeDDASIT: {
			models.ServiceItemParamNameSITPaymentRequestStart,
			models.ServiceItemParamNameSITPaymentRequestEnd,
			models.ServiceItemParamNameWeightBilled,
		},
		models.ReServiceCodeDOFSIT: {
			models.ServiceItemParamNameWeightBilled,
		},
		models.ReServiceCodeDLH: {
			models.ServiceItemParamNameWeightBilled,
		},
		models.ReServiceCodeFSC: {
			models.ServiceItemParamNameWeightBilled,
		},
		models.ReServiceCodeDSH: {
			models.ServiceItemParamNameWeightBilled,
		},
		models.ReServiceCodeDUPK: {
			models.ServiceItemParamNameWeightBilled,
		},
		models.ReServiceCodeDNPK: {
			models.ServiceItemParamNameWeightBilled,
		},
		models.ReServiceCodeDOPSIT: {
			models.ServiceItemParamNameWeightBilled,
		},
		models.ReServiceCodeDDDSIT: {
			models.ServiceItemParamNameWeightBilled,
		},
		models.ReServiceCodeDDSHUT: {
			models.ServiceItemParamNameWeightBilled,
		},
		models.ReServiceCodeDOSHUT: {
			models.ServiceItemParamNameWeightBilled,
		},
		models.ReServiceCodeDDFSIT: {
			models.ServiceItemParamNameWeightBilled,
		},
		models.ReServiceCodeDOP: {
			models.ServiceItemParamNameWeightBilled,
		},
		models.ReServiceCodeDDP: {
			models.ServiceItemParamNameWeightBilled,
		},
		models.ReServiceCodeDPK: {
			models.ServiceItemParamNameWeightBilled,
		},
		models.ReServiceCodeDDSFSC: {
			models.ServiceItemParamNameWeightBilled,
		},
		models.ReServiceCodeDOSFSC: {
			models.ServiceItemParamNameWeightBilled,
		},
		models.ReServiceCodeISLH: {
			models.ServiceItemParamNameWeightBilled,
		},
		models.ReServiceCodeIHPK: {
			models.ServiceItemParamNameWeightBilled,
		},
		models.ReServiceCodeIHUPK: {
			models.ServiceItemParamNameWeightBilled,
		},
		models.ReServiceCodePOEFSC: {
			models.ServiceItemParamNameWeightBilled,
		},
		models.ReServiceCodePODFSC: {
			models.ServiceItemParamNameWeightBilled,
		},
	}
)

// Contains checks to see if the provided param key is valid for the given service code
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
