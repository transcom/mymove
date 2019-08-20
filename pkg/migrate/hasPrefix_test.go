package migrate

func (suite *MigrateSuite) TestHasPrefix() {
	suite.True(hasPrefix("file://some.file.to.load", "file://"))
	suite.False(hasPrefix("file://some.file.to.load", "s3://"))
}
