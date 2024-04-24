// RA Summary: gosec - errcheck - Unchecked return value
// RA: Linter flags errcheck error: Ignoring a method's return value can cause the program to overlook unexpected states and conditions.
// RA: Functions with unchecked return values in the file are used to generate stub data for a localized version of the application.
// RA: Given the data is being generated for local use and does not contain any sensitive information, there are no unexpected states and conditions
// RA: in which this would be considered a risk
// RA Developer Status: Mitigated
// RA Validator Status: Mitigated
// RA Modified Severity: N/A
// nolint:errcheck
package paperwork

func (suite *PaperworkServiceSuite) TestCreateFormServiceCreateAssetByteReaderFailure() {
	badAssetPath := "paperwork/formtemplates/someUndefinedTemplatePath.png"
	templateBuffer, err := createAssetByteReader(badAssetPath)
	suite.Nil(templateBuffer)
	suite.NotNil(err)
	suite.Equal("error creating asset from path; check image path: open paperwork/formtemplates/someUndefinedTemplatePath.png: file does not exist", err.Error())
}
