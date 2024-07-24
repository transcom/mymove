package models_test

import (
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
)

func (suite *ModelSuite) TestMobileHomeShipmentValidation() {
	testCases := map[string]struct {
		mobileHome   models.MobileHome
		expectedErrs map[string][]string
	}{
		"Successful Minimal Validation": {
			mobileHome: models.MobileHome{
				ShipmentID: uuid.Must(uuid.NewV4()),
				Make:       "Mobile Home Make",
				Model:      "Mobile Home Model",
				Year:       1996,
				Length:     200,
				Height:     84,
				Width:      96,
			},
			expectedErrs: nil,
		},
		"Missing Required Fields": {
			mobileHome: models.MobileHome{},
			expectedErrs: map[string][]string{
				"shipment_id": {"ShipmentID can not be blank."},
				"make":        {"Make can not be blank."},
				"model":       {"Model can not be blank."},
				"year":        {"0 is not greater than 0."},
				"length":      {"0 is not greater than 0."},
				"height":      {"0 is not greater than 0."},
				"width":       {"0 is not greater than 0."},
			},
		},
	}

	for name, testCase := range testCases {
		name, testCase := name, testCase

		suite.Run(name, func() {
			suite.verifyValidationErrors(testCase.mobileHome, testCase.expectedErrs)
		})
	}
}
