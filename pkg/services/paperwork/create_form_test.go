package paperwork

func (suite *PaperworkServiceSuite) TestCreateFormServiceCreateAssetByteReaderFailure() {
	badAssetPath := "pkg/paperwork/formtemplates/someUndefinedTemplatePath.png"
	templateBuffer, err := createAssetByteReader(badAssetPath)
	suite.Nil(templateBuffer)
	suite.NotNil(err)
	suite.Equal("Error creating asset from path. Check image path.: Asset pkg/paperwork/formtemplates/someUndefinedTemplatePath.png not found", err.Error())
}
