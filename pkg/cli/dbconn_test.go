package cli

func (suite *cliTestSuite) TestConfigDatabase() {
	suite.Setup(InitDatabaseFlags, []string{})
	suite.NoError(CheckDatabase(suite.viper, suite.logger))
}

func (suite *cliTestSuite) TestInitDatabase() {
	suite.Setup(InitDatabaseFlags, []string{})
	conn, err := InitDatabase(suite.viper, nil, suite.logger)
	suite.NoError(err)
	suite.NotNil(conn)
	defer conn.Close()
}

func (suite *cliTestSuite) TestConfigDatabaseRetry() {
	suite.Setup(InitDatabaseRetryFlags, []string{})
	suite.NoError(CheckDatabaseRetry(suite.viper))
}

func (suite *cliTestSuite) TestPingPopConnectionOk() {
	suite.Setup(InitDatabaseFlags, []string{})
	conn, err := InitDatabase(suite.viper, nil, suite.logger)
	suite.NoError(err)
	suite.NotNil(conn)
	defer conn.Close()
	suite.NoError(PingPopConnection(conn, suite.logger))
}

func (suite *cliTestSuite) TestPingPopConnectionFail() {
	// intentionally misconfigure the db so the ping will fail
	suite.Setup(InitDatabaseFlags, []string{"--" + DbNameFlag, "missingdb"})
	conn, err := InitDatabase(suite.viper, nil, suite.logger)
	suite.NoError(err)
	suite.NotNil(conn)
	defer conn.Close()
	suite.Error(PingPopConnection(conn, suite.logger))
}
