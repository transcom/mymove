package iws

import "encoding/xml"

func (suite *iwsSuite) TestWkemaSuccessResponseUnmarshal() {
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
	rec := Record{}
	unmarshalErr := xml.Unmarshal([]byte(data), &rec)
	suite.Nil(unmarshalErr)
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

func (suite *iwsSuite) TestGetPersonUsingWorkEmailInvalid() {
	// An empty SSN should get an RbsError from the API
	edipi, person, personnel, err := GetPersonUsingWorkEmail(suite.client, suite.host, suite.custNum, "")
	// error should still be nil - no match is not an error like connection failure
	suite.NotNil(err)
	suite.IsType(&RbsError{}, err)
	suite.Zero(edipi)
	suite.Nil(person)
	suite.Empty(personnel)
}
