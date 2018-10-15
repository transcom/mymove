package iws

import "net/url"

func (suite *iwsSuite) TestBuildEdiURL() {
	urlString, err := buildEdiURL("example.com", "1234", 1234567890)
	suite.Nil(err)
	parsedURL, parseErr := url.Parse(urlString)
	suite.Nil(parseErr)
	suite.Equal("https", parsedURL.Scheme)
	suite.Equal("example.com", parsedURL.Host)
	suite.Equal("/appj/rbs/rest/op=edi/customer=1234/schemaName=get_cac_data/schemaVersion=1.0/DOD_EDI_PN_ID=1234567890", parsedURL.Path)
}

func (suite *iwsSuite) TestBuildEdiURLInvalidEDIPI() {
	urlString, err := buildEdiURL("example.com", "1234", 10000000000)
	suite.NotNil(err)
	suite.Empty(urlString)
}

func (suite *iwsSuite) TestParseEdiResponse() {
	data := `<?xml version="1.0" encoding="UTF-8" standalone="no" ?>
	<record>
		<rule>
			<customer>1234</customer>
			<schemaName>schema_name</schemaName>
			<schemaVersion>1.0</schemaVersion>
		</rule>
		<identifier>
			<DOD_EDI_PN_ID>9995006001</DOD_EDI_PN_ID>
		</identifier>
		<adrRecord>
			<DOD_EDI_PN_ID>9995006001</DOD_EDI_PN_ID>
			<person>
				<PN_ID>xxxx12345</PN_ID>
				<PN_ID_TYP_CD>S</PN_ID_TYP_CD>
				<PN_1ST_NM>Mickey</PN_1ST_NM>
				<PN_MID_NM>Middle</PN_MID_NM>
				<PN_LST_NM>Mantle</PN_LST_NM>
				<PN_CDNCY_NM>III</PN_CDNCY_NM>
				<PN_BRTH_DT>19311020</PN_BRTH_DT>
				<PN_DTH_CD>N</PN_DTH_CD>
			</person>
			<personnel>
				<PNL_CAT_CD>A</PNL_CAT_CD>
				<PNL_PE_DT>20201101</PNL_PE_DT>
				<PNL_TERM_DT>20201101</PNL_TERM_DT>
				<RANK_CD>CPL</RANK_CD>
				<SVC_CD>A</SVC_CD>
				<UNIT_ID_CD>00000</UNIT_ID_CD>
			</personnel>
		</adrRecord>
	</record>`
	person, personnel, err := parseEdiResponse([]byte(data))
	suite.Nil(err)
	suite.NotNil(person)
	suite.Equal("xxxx12345", person.ID)
	suite.NotEmpty(personnel)
}

func (suite *iwsSuite) TestParseEdiResponseError() {
	data := `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
	<RbsError>
		<faultCode>14030</faultCode>
		<faultMessage>DOD_EDI_PN_ID should be in the range between 1000000000 and 9999999999</faultMessage>
	</RbsError>`
	person, personnel, err := parseEdiResponse([]byte(data))
	suite.Nil(person)
	suite.Empty(personnel)
	suite.NotNil(err)
	rbsError, ok := err.(*RbsError)
	suite.True(ok)
	suite.Equal(uint64(14030), rbsError.FaultCode)
}

func (suite *iwsSuite) TestGetPersonUsingEDIPI() {
	person, personnel, err := GetPersonUsingEDIPI(suite.client, suite.host, suite.custNum, 1920203960)
	suite.Nil(err)
	suite.NotNil(person)
	suite.NotEmpty(personnel)
}

func (suite *iwsSuite) TestGetPersonUsingEDIPINotFound() {
	person, personnel, err := GetPersonUsingEDIPI(suite.client, suite.host, suite.custNum, 9999999999)
	// error should still be nil - no match is not an error like connection failure
	suite.Nil(err)
	suite.Nil(person)
	suite.Empty(personnel)
}

func (suite *iwsSuite) TestGetPersonUsingEDIPIInvalid() {
	// Lowest valid EDIPI is 1000000000, so this should get an RbsError from the API
	person, personnel, err := GetPersonUsingEDIPI(suite.client, suite.host, suite.custNum, 0)
	suite.NotNil(err)
	suite.IsType(&RbsError{}, err)
	suite.Nil(person)
	suite.Empty(personnel)
}
