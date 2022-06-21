package edisegment

func (suite *SegmentSuite) TestValidatePER() {
	validPERDefault := PER{
		ContactFunctionCode: "IC",
	}

	validPERAll := PER{
		ContactFunctionCode:          "IC",
		Name:                         "Cross Dock",
		CommunicationNumberQualifier: "TE",
		CommunicationNumber:          "5551234567",
	}

	suite.Run("validate success default", func() {
		err := suite.validator.Struct(validPERDefault)
		suite.NoError(err)
	})

	suite.Run("validate success all", func() {
		err := suite.validator.Struct(validPERAll)
		suite.NoError(err)
	})

	suite.Run("validate failure 1", func() {
		per := PER{
			Name:                         "Cross Dock Testing With Too Long a Name Cross Dock Testing With Too Long a Name",
			CommunicationNumberQualifier: "BX",
			CommunicationNumber:          "111111111111111111111111111111111111111115551234567111111111111111111111111111111111111111115551234567",
		}

		err := suite.validator.Struct(per)
		suite.ValidateError(err, "ContactFunctionCode", "required")
		suite.ValidateError(err, "Name", "max")
		suite.ValidateError(err, "CommunicationNumberQualifier", "eq")
		suite.ValidateError(err, "CommunicationNumber", "max")
		suite.ValidateErrorLen(err, 4)
	})

	suite.Run("validate segment is parsed correctly", func() {
		values := []string{"IC", "Cross Dock", "TE", "5551234567"}
		per := PER{
			ContactFunctionCode:          "IC",
			Name:                         "Cross Dock",
			CommunicationNumberQualifier: "TE",
			CommunicationNumber:          "5551234567",
		}
		err := (*PER).Parse(&per, values)
		suite.NoError(err)
	})
}
