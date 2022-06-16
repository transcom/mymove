package edisegment

func (suite *SegmentSuite) TestValidateMEA() {
	validMEADefault := MEA{
		MeasurementValue: 100.0,
	}

	validMEAWithValues := MEA{
		MeasurementReferenceIDCode: "XX",
		MeasurementQualifier:       "HT",
		MeasurementValue:           100.0,
	}

	suite.Run("validate success default", func() {
		err := suite.validator.Struct(validMEADefault)
		suite.NoError(err)
	})

	suite.Run("validate success with values", func() {
		err := suite.validator.Struct(validMEAWithValues)
		suite.NoError(err)
	})

	suite.Run("validate failure", func() {
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
