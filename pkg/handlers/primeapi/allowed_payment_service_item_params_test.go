package primeapi

import (
	"fmt"

	"github.com/transcom/mymove/pkg/models"
)

type paramTestCase struct {
	reServiceCode models.ReServiceCode
	paramKeyName  string
}

var allowedParamsTestCases = []paramTestCase{
	{models.ReServiceCodeDDASIT, "SITPaymentRequestStart"},
	{models.ReServiceCodeDDASIT, "SITPaymentRequestEnd"},
	{models.ReServiceCodeDOASIT, "SITPaymentRequestStart"},
	{models.ReServiceCodeDOASIT, "SITPaymentRequestEnd"},
}

var invalidParamsTestCases = []paramTestCase{
	// Invalid params for service items that do have allowed params
	{models.ReServiceCodeDDASIT, "NotARealParamKey"},
	{models.ReServiceCodeDOASIT, "NotARealParamKey"},
	// Real params used for the wrong service items
	{models.ReServiceCodeDOASIT, "ZipSITOriginHHGActualAddress"},
	{models.ReServiceCodeDOFSIT, "SITPaymentRequestStart"},
	{models.ReServiceCodeDDFSIT, "SITPaymentRequestEnd"},
}

func (suite *HandlerSuite) TestAllowedParams() {
	for _, tc := range allowedParamsTestCases {
		suite.Run(fmt.Sprintf("param %s should be allowed for service code %s", tc.paramKeyName, string(tc.reServiceCode)), func() {
			suite.True(AllowedParamKeysPaymentRequest.Contains(tc.reServiceCode, tc.paramKeyName))
		})
	}
}

func (suite *HandlerSuite) TestNotAllowedParams() {
	for _, tc := range invalidParamsTestCases {
		suite.Run(fmt.Sprintf("param %s should not be allowed for service code %s", tc.paramKeyName, string(tc.reServiceCode)), func() {
			suite.False(AllowedParamKeysPaymentRequest.Contains(tc.reServiceCode, tc.paramKeyName))
		})
	}
}
