package migrate

func (suite *MigrateSuite) TestIsAfterSpaceZero() {

	in := "hello world"
	buf := NewBuffer()

	_, err := buf.WriteString(in)
	suite.NoError(err)
	buf.Close()

	suite.True(isAfterSpace(buf, 0))
}

func (suite *MigrateSuite) TestIsAfterSpace() {

	in := "hello world"
	buf := NewBuffer()

	_, err := buf.WriteString(in)
	suite.NoError(err)
	buf.Close()

	suite.True(isAfterSpace(buf, 6))
}
