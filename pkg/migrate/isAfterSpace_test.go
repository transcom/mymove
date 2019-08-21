package migrate

func (suite *MigrateSuite) TestIsAfterSpaceZero() {

	in := "hello world"
	buf := NewBuffer()

	buf.WriteString(in)
	buf.Close()

	suite.True(isAfterSpace(buf, 0))
}

func (suite *MigrateSuite) TestIsAfterSpace() {

	in := "hello world"
	buf := NewBuffer()

	buf.WriteString(in)
	buf.Close()

	suite.True(isAfterSpace(buf, 6))
}
