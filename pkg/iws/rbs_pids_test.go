package iws

import "net/url"

func (suite *iwsSuite) TestParsePidsResponse() {
	data := `<?xml version="1.0" encoding="UTF-8" standalone="no" ?>
	<record>
		<rule>
			<customer>1234</customer>
			<schemaName>schema_name</schemaName>
			<schemaVersion>1.0</schemaVersion>
		</rule>
		<identifier>
			<pids>
				<PN_ID>xxxx12345</PN_ID>
				<PN_ID_TYP_CD>S</PN_ID_TYP_CD>
				<PN_LST_NM>Mantle</PN_LST_NM>
				<PN_1ST_NM>Mickey</PN_1ST_NM>
				<PN_BRTH_DT>19311020</PN_BRTH_DT>
			</pids>
		</identifier>
		<adrRecord>
			<PIDSRecord>
				<DOD_EDI_PN_ID>9995006001</DOD_EDI_PN_ID>
				<MTCH_RSN_CD>PMC</MTCH_RSN_CD>
			</PIDSRecord>
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

	reason, edipi, person, personnel, err := parsePidsResponse([]byte(data))
	suite.Nil(err)
	suite.Equal(MatchReasonCodeFull, reason)
	suite.Equal(uint64(9995006001), edipi)
	suite.NotNil(person)
	suite.NotEmpty(personnel)
}

func (suite *iwsSuite) TestParsePidsResponseError() {
	data := `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
	<RbsError>
	  <faultCode>14030</faultCode>
	  <faultMessage> Problem with this argument: PN_ID</faultMessage>
	</RbsError>`
	reason, edipi, person, personnel, err := parsePidsResponse([]byte(data))
	suite.NotNil(err)
	rbsErr, typeErr := err.(*RbsError)
	suite.True(typeErr)
	suite.Equal(uint64(14030), rbsErr.FaultCode)
	suite.Equal(MatchReasonCodeNone, reason)
	suite.Zero(edipi)
	suite.Nil(person)
	suite.Empty(personnel)
}

func (suite *iwsSuite) TestBuildPidsUrl() {
	urlString, err := buildPidsURL("example.com", "1234", "000000000", "Last", "First")
	suite.NotEmpty(urlString)
	suite.Nil(err)
	parsedURL, parseErr := url.Parse(urlString)
	suite.Nil(parseErr)
	suite.Equal("https", parsedURL.Scheme)
	suite.Equal("example.com", parsedURL.Host)
	suite.Equal("/appj/rbs/rest/op=pids-P/customer=1234/schemaName=get_cac_data/schemaVersion=1.0/PN_ID=000000000/PN_ID_TYP_CD=S/PN_LST_NM=Last/PN_1ST_NM=First", parsedURL.Path)
}

func (suite *iwsSuite) TestBuildPidsUrlLongNames() {
	urlString, err := buildPidsURL("example.com", "1234", "000000000", "abcdefghijklmnopqrstuvwxyzyxwvutsrqponmlkjihgfedcba", "abcdefghijklmnopqrstuvwxyz")
	suite.NotEmpty(urlString)
	suite.Nil(err)
	parsedURL, parseErr := url.Parse(urlString)
	suite.Nil(parseErr)
	suite.Equal("https", parsedURL.Scheme)
	suite.Equal("example.com", parsedURL.Host)
	suite.Equal("/appj/rbs/rest/op=pids-P/customer=1234/schemaName=get_cac_data/schemaVersion=1.0/PN_ID=000000000/PN_ID_TYP_CD=S/PN_LST_NM=abcdefghijklmnopqrstuvwxyz/PN_1ST_NM=abcdefghijklmnopqrst", parsedURL.Path)
}

func (suite *iwsSuite) TestBuildPidsUrlNoFirstName() {
	urlString, err := buildPidsURL("example.com", "1234", "000000000", "Last", "")
	suite.NotEmpty(urlString)
	suite.Nil(err)
	parsedURL, parseErr := url.Parse(urlString)
	suite.Nil(parseErr)
	suite.Equal("https", parsedURL.Scheme)
	suite.Equal("example.com", parsedURL.Host)
	suite.Equal("/appj/rbs/rest/op=pids-P/customer=1234/schemaName=get_cac_data/schemaVersion=1.0/PN_ID=000000000/PN_ID_TYP_CD=S/PN_LST_NM=Last", parsedURL.Path)
}

func (suite *iwsSuite) TestBuildPidsUrlBadSSN() {
	urlString, err := buildPidsURL("example.com", "1234", "12345678", "Last", "First")
	suite.Empty(urlString)
	suite.NotNil(err)
}
