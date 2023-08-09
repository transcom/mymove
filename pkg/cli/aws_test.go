package cli

import (
	"fmt"
)

func (suite *cliTestSuite) TestConfigAWS() {
	flagSet := []string{
		fmt.Sprintf("--%s=%s", AWSRegionFlag, "us-west-2"),
	}
	suite.Setup(InitAWSFlags, flagSet)
	region, err := CheckAWSRegion(suite.viper)
	suite.NoError(err)
	suite.Equal("us-west-2", region)
}
