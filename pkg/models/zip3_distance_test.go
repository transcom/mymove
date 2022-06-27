package models_test

import (
	"github.com/transcom/mymove/pkg/models"
)

func (suite *ModelSuite) TestZip3DistanceValidations() {
	suite.Run("test valid Zip3Distance", func() {
		validZip3Distance := models.Zip3Distance{
			FromZip3:      "010",
			ToZip3:        "011",
			DistanceMiles: 24,
		}
		expErrors := map[string][]string{}
		suite.verifyValidationErrors(&validZip3Distance, expErrors)
	})

	suite.Run("test invalid Zip3Distance", func() {
		emptyZip3Distance := models.Zip3Distance{}
		expErrors := map[string][]string{
			"from_zip3":      {"FromZip3 not in range(3, 3)"},
			"to_zip3":        {"ToZip3 not in range(3, 3)"},
			"distance_miles": {"DistanceMiles can not be blank."},
		}
		suite.verifyValidationErrors(&emptyZip3Distance, expErrors)
	})

	suite.Run("test when from_zip3 is not a length of 3", func() {
		invalidFromZip3Distance := models.Zip3Distance{
			FromZip3:      "01",
			ToZip3:        "011",
			DistanceMiles: 24,
		}
		expErrors := map[string][]string{
			"from_zip3": {"FromZip3 not in range(3, 3)"},
		}
		suite.verifyValidationErrors(&invalidFromZip3Distance, expErrors)
	})

	suite.Run("test when to_zip3 is not a length of 3", func() {
		invalidToZip3Distance := models.Zip3Distance{
			FromZip3:      "010",
			ToZip3:        "0115",
			DistanceMiles: 24,
		}
		expErrors := map[string][]string{
			"to_zip3": {"ToZip3 not in range(3, 3)"},
		}
		suite.verifyValidationErrors(&invalidToZip3Distance, expErrors)
	})

	suite.Run("test when distance_miles is not provided", func() {
		invalidDistance := models.Zip3Distance{
			FromZip3: "010",
			ToZip3:   "011",
		}
		expErrors := map[string][]string{
			"distance_miles": {"DistanceMiles can not be blank."},
		}
		suite.verifyValidationErrors(&invalidDistance, expErrors)
	})
}
