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
				Year:           models.IntPointer(2000),
				Make:           models.StringPointer("Boat Make"),
				Model:          models.StringPointer("Boat Model"),
				LengthInInches: models.IntPointer(300),
				WidthInInches:  models.IntPointer(108),
				HeightInInches: models.IntPointer(72),
				HasTrailer:     models.BoolPointer(false),
			},
			expectedErrs: nil,
		},
		"Missing Required Fields": {
			boatShipment: models.BoatShipment{
				Year:           models.IntPointer(0),
				Make:           models.StringPointer(""),
				Model:          models.StringPointer(""),
				LengthInInches: models.IntPointer(0),
				WidthInInches:  models.IntPointer(0),
				HeightInInches: models.IntPointer(0),
			},
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
			suite.verifyValidationErrors(testCase.boatShipment, testCase.expectedErrs, nil)
		})
	}
}
