package migrate

func (suite *MigrateSuite) TestByteIsSpace() {
	suite.True(byteIsSpace(byte('\t')))
	suite.True(byteIsSpace(byte('\n')))
	suite.True(byteIsSpace(byte('\v')))
	suite.True(byteIsSpace(byte('\f')))
	suite.True(byteIsSpace(byte('\r')))
	suite.True(byteIsSpace(byte(' ')))
}
