package models_test

import (
	"fmt"
	"strings"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
)

func (suite *ModelSuite) TestBoatShipmentValidation() {
	validBoatShipmentTypes := strings.Join(models.AllowedBoatShipmentTypes, ", ")

	testCases := map[string]struct {
		boatShipment models.BoatShipment
		expectedErrs map[string][]string
	}{
		"Successful Minimal Validation": {
			boatShipment: models.BoatShipment{
				ShipmentID:     uuid.Must(uuid.NewV4()),
				Type:           models.BoatShipmentTypeHaulAway,
				Year:           1991,
				Make:           "Boat Make",
				Model:          "Boat Model",
				LengthInInches: 200,
				WidthInInches:  96,
				HeightInInches: 84,
				HasTrailer:     false,
			},
			expectedErrs: nil,
		},
		"Missing Required Fields": {
			boatShipment: models.BoatShipment{},
			expectedErrs: map[string][]string{
				"shipment_id":      {"ShipmentID can not be blank."},
				"type":             {fmt.Sprintf("Type is not in the list [%s].", validBoatShipmentTypes)},
				"year":             {"0 is not greater than 0."},
				"make":             {"Make can not be blank."},
				"model":            {"Model can not be blank."},
				"length_in_inches": {"0 is not greater than 0."},
				"width_in_inches":  {"0 is not greater than 0."},
				"height_in_inches": {"0 is not greater than 0."},
			},
		},
	}

	for name, testCase := range testCases {
		name, testCase := name, testCase

		suite.Run(name, func() {
			suite.verifyValidationErrors(testCase.boatShipment, testCase.expectedErrs)
		})
	}
}
