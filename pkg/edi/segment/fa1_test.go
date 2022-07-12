package edisegment

func (suite *SegmentSuite) TestValidateFA1() {
	validFA1 := FA1{
		AgencyQualifierCode: "DF",
	}

	suite.Run("validate success", func() {
		err := suite.validator.Struct(validFA1)
		suite.NoError(err)
	})

	suite.Run("validate failure", func() {
		fa1 := FA1{
			AgencyQualifierCode: "XX", // oneof
		}

		err := suite.validator.Struct(fa1)
		suite.ValidateError(err, "AgencyQualifierCode", "oneof")
		suite.ValidateErrorLen(err, 1)
	})
}
