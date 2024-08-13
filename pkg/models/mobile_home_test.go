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
				ShipmentID:     uuid.Must(uuid.NewV4()),
				Make:           "Mobile Home Make",
				Model:          "Mobile Home Model",
				Year:           1996,
				LengthInInches: 200,
				HeightInInches: 84,
				WidthInInches:  96,
			},
			expectedErrs: nil,
		},
		"Missing Required Fields": {
			mobileHome: models.MobileHome{},
			expectedErrs: map[string][]string{
				"shipment_id":      {"ShipmentID can not be blank."},
				"make":             {"Make can not be blank."},
				"model":            {"Model can not be blank."},
				"year":             {"0 is not greater than 0."},
				"length_in_inches": {"0 is not greater than 0."},
				"height_in_inches": {"0 is not greater than 0."},
				"width_in_inches":  {"0 is not greater than 0."},
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
