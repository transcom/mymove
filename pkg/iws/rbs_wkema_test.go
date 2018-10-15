package iws

import "net/url"

func (suite *iwsSuite) TestParseWkEmaResponse() {
	data := `<?xml version="1.0" encoding="UTF-8" standalone="no"?>
	<record>
	  <rule>
		<customer>1234</customer>
		<schemaName>schema_name</schemaName>
		<schemaVersion>1.0</schemaVersion>
	  </rule>
	  <identifier>
		<wkEma>
		  <EMA_TX>nobody_here@mail.mil</EMA_TX>
		</wkEma>
	  </identifier>
	  <adrRecord>
		<WKEMARecord>
		  <DOD_EDI_PN_ID>9995006001</DOD_EDI_PN_ID>
		  <EMA_TX>nobody_here@mail.mil</EMA_TX>
		</WKEMARecord>
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
	edipi, person, personnel, err := parseWkEmaResponse([]byte(data))
	suite.Nil(err)
	suite.Equal(uint64(9995006001), edipi)
	suite.NotNil(person)
	suite.NotEmpty(personnel)
}

func (suite *iwsSuite) TestParseWkEmaResponseError() {
	data := `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
	<RbsError>
	  <faultCode>14030</faultCode>
	  <faultMessage> Problem with this argument: EMA_TX</faultMessage>
	</RbsError>`
	edipi, person, personnel, err := parseWkEmaResponse([]byte(data))
	suite.NotNil(err)
	rbsError, ok := err.(*RbsError)
	suite.True(ok)
	suite.Equal(uint64(14030), rbsError.FaultCode)
	suite.NotEmpty(rbsError.FaultMessage)
	suite.Zero(edipi)
	suite.Nil(person)
	suite.Empty(personnel)
}

func (suite *iwsSuite) TestBuildWkEmaURL() {
	urlString, err := buildWkEmaURL("example.com", "1234", "test@example.com")
	suite.NotEmpty(urlString)
	suite.Nil(err)
	parsedURL, parseErr := url.Parse(urlString)
	suite.Nil(parseErr)
	suite.Equal("https", parsedURL.Scheme)
	suite.Equal("example.com", parsedURL.Host)
	suite.Equal("/appj/rbs/rest/op=wkEma/customer=1234/schemaName=get_cac_data/schemaVersion=1.0/EMA_TX=test@example.com", parsedURL.Path)
}

func (suite *iwsSuite) TestBuildWkEmaURLEmailInvalid() {
	u, err := buildWkEmaURL("example.com", "1234", "invalid@")
	suite.NotNil(err)
	suite.Empty(u)
}

func (suite *iwsSuite) TestBuildWkEmaURLLongEmail() {
	urlString, err := buildWkEmaURL("example.com", "1234", "pneumonoultramicroscopicsilicovolcanoconiosis_is_a_terrible_way_to_expire@unpronounceablediseases.org")
	suite.NotEmpty(urlString)
	suite.Nil(err)
	parsedURL, parseErr := url.Parse(urlString)
	suite.Nil(parseErr)
	suite.Equal("https", parsedURL.Scheme)
	suite.Equal("example.com", parsedURL.Host)
	suite.Equal("/appj/rbs/rest/op=wkEma/customer=1234/schemaName=get_cac_data/schemaVersion=1.0/EMA_TX=pneumonoultramicroscopicsilicovolcanoconiosis_is_a_terrible_way_to_expire@unpron", parsedURL.Path)
}

func (suite *iwsSuite) TestGetPersonUsingWorkEmail() {
	edipi, person, personnel, err := GetPersonUsingWorkEmail(suite.client, suite.host, suite.custNum, "matthew.m.heitner@ctwork.com")
	suite.Nil(err)
	suite.Equal(uint64(1920203960), edipi)
	suite.NotNil(person)
	suite.NotEmpty(personnel)
}

func (suite *iwsSuite) TestGetPersonUsingWorkEmailNotFound() {
	edipi, person, personnel, err := GetPersonUsingWorkEmail(suite.client, suite.host, suite.custNum, "nobody@example.com")
	// error should still be nil - no match is not an error like connection failure
	suite.Nil(err)
	suite.Zero(edipi)
	suite.Nil(person)
	suite.Empty(personnel)
}
