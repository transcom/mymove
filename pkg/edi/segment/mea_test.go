package edisegment

import (
	"testing"
)

func (suite *SegmentSuite) TestValidateMEA() {
	validMEADefault := MEA{
		MeasurementValue: 100.0,
	}

	validMEAWithValues := MEA{
		MeasurementReferenceIDCode: "XX",
		MeasurementQualifier:       "HT",
		MeasurementValue:           100.0,
	}

	suite.T().Run("validate success default", func(t *testing.T) {
		err := suite.validator.Struct(validMEADefault)
		suite.NoError(err)
	})

	suite.T().Run("validate success with values", func(t *testing.T) {
		err := suite.validator.Struct(validMEAWithValues)
		suite.NoError(err)
	})

	suite.T().Run("validate failure", func(t *testing.T) {
		mea := MEA{
			MeasurementReferenceIDCode: "ABC", // len
			MeasurementQualifier:       "XX",  // oneof
			// MeasurementValue is required
		}

		err := suite.validator.Struct(mea)
		suite.ValidateError(err, "MeasurementReferenceIDCode", "len")
		suite.ValidateError(err, "MeasurementQualifier", "oneof")
		suite.ValidateError(err, "MeasurementValue", "required")
		suite.ValidateErrorLen(err, 3)
	})
}
