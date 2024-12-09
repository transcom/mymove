package iws

func (suite *iwsSuite) TestRbsError() {
	data := `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
	<RbsError>
	  <faultCode>14030</faultCode>
	  <faultMessage> Problem with this argument: EMA_TX</faultMessage>
	</RbsError>`
	_, _, _, err := parseWkEmaResponse([]byte(data))
	suite.NotNil(err)
	rbsError, ok := err.(*RbsError)
	suite.True(ok)
	suite.Equal(uint64(14030), rbsError.FaultCode)
	suite.NotEmpty(rbsError.FaultMessage)
}
