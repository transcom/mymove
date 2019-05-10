package cli

import (
	"fmt"
)

func (suite *cliTestSuite) TestConfigAWS() {
	expectedRegion := "us-west-2"
	flagSet := []string{
		fmt.Sprintf("--%s=%s", AWSRegionFlag, expectedRegion),
	}
	suite.Setup(InitAWSFlags, flagSet)
	region, err := CheckAWSRegion(suite.viper)
	suite.Nil(err)
	suite.Equal(region, expectedRegion)
}
