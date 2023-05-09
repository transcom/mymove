package models_test

import (
	"github.com/transcom/mymove/pkg/factory"
	. "github.com/transcom/mymove/pkg/models"
)

func (suite *ModelSuite) TestSITAddressUpdateInstantiation() {
	testMTOServiceItem := factory.BuildMTOServiceItem(suite.DB(), nil, nil)
	testOldAddress := factory.BuildAddress(suite.DB(), nil, nil)
	testNewAddress := factory.BuildAddress(suite.DB(), nil, nil)
	testContractorRemarksString := "test contractor remarks"

	type TestCaseType struct {
		name             string
		sitAddressUpdate SITAddressUpdate
		expectedErrs     map[string][]string
	}

	testCases := []TestCaseType{
		{
			name: "Successful create",
			sitAddressUpdate: SITAddressUpdate{
				MTOServiceItemID:  testMTOServiceItem.ID,
				OldAddressID:      testOldAddress.ID,
				NewAddressID:      testNewAddress.ID,
				ContractorRemarks: &testContractorRemarksString,
				Distance:          1323,
				Reason:            "Not a real reason",
				Status:            SITAddressStatusRejected,
			},
			expectedErrs: nil,
		},
		{
			name:             "Missing UUIDs",
			sitAddressUpdate: SITAddressUpdate{},
			expectedErrs: map[string][]string{
				"mtoservice_item_id": {"MTOServiceItemID can not be blank."},
				"old_address_id":     {"OldAddressID can not be blank."},
				"new_address_id":     {"NewAddressID can not be blank."},
				"distance":           {"Distance can not be blank."},
				"reason":             {"Reason can not be blank."},
				"status":             {"Status is not in the list [REQUESTED, REJECTED, APPROVED]."},
			},
		},
		{
			name: "Optional fields are invalid",
			sitAddressUpdate: SITAddressUpdate{
				MTOServiceItemID:  testMTOServiceItem.ID,
				OldAddressID:      testOldAddress.ID,
				NewAddressID:      testNewAddress.ID,
				ContractorRemarks: StringPointer(""),
				Distance:          1323,
				Reason:            "Not a real reason",
				Status:            SITAddressStatusRejected,
				OfficeRemarks:     StringPointer(""),
			},
			expectedErrs: map[string][]string{
				"office_remarks":     {"OfficeRemarks can not be blank."},
				"contractor_remarks": {"ContractorRemarks can not be blank."},
			},
		},
	}

	for _, testCase := range testCases {
		name := testCase.name
		model := testCase.sitAddressUpdate
		expectedErrs := testCase.expectedErrs

		suite.Run(name, func() {
			suite.verifyValidationErrors(&model, expectedErrs)
		})
	}

}
