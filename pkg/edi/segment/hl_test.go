package edisegment

func (suite *SegmentSuite) TestValidateHL() {
	validHL := HL{
		HierarchicalIDNumber:  "303",
		HierarchicalLevelCode: "SS",
	}

	suite.Run("validate success", func() {
		err := suite.validator.Struct(validHL)
		suite.NoError(err)
	})

	suite.Run("validate failure", func() {
		hl := HL{
			HierarchicalIDNumber:       "A-123", // alphanum
			HierarchicalParentIDNumber: "1",     // isdefault
			HierarchicalLevelCode:      "XX",    // eq
		}

		err := suite.validator.Struct(hl)
		suite.ValidateError(err, "HierarchicalIDNumber", "alphanum")
		suite.ValidateError(err, "HierarchicalParentIDNumber", "isdefault")
		suite.ValidateError(err, "HierarchicalLevelCode", "oneof")
		suite.ValidateErrorLen(err, 3)
	})

	suite.Run("validate failure 2", func() {
		hl := validHL
		hl.HierarchicalIDNumber = "" // alphanum takes precidence over min

		err := suite.validator.Struct(hl)
		suite.ValidateError(err, "HierarchicalIDNumber", "alphanum")
		suite.ValidateErrorLen(err, 1)
	})

	suite.Run("validate failure 3", func() {
		hl := validHL
		hl.HierarchicalIDNumber = "0123456789ABCDF" // max

		err := suite.validator.Struct(hl)
		suite.ValidateError(err, "HierarchicalIDNumber", "max")
		suite.ValidateErrorLen(err, 1)
	})
}
