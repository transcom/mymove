package cli

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws/endpoints"
)

func (suite *cliTestSuite) TestConfigAWS() {
	flagSet := []string{
		fmt.Sprintf("--%s=%s", AWSRegionFlag, endpoints.UsWest2RegionID),
	}
	suite.Setup(InitAWSFlags, flagSet)
	region, err := CheckAWSRegion(suite.viper)
	suite.Nil(err)
	suite.Equal(endpoints.UsWest2RegionID, region)
}
