package cli

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/aws/aws-sdk-go/service/cloudwatchevents"
	"github.com/aws/aws-sdk-go/service/ecr"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/aws/aws-sdk-go/service/rds"
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

func (suite *cliTestSuite) TestCheckAWSRegionForService() {
	region := endpoints.UsWest2RegionID

	err := CheckAWSRegionForService(region, cloudwatchevents.ServiceName)
	suite.Nil(err)

	err = CheckAWSRegionForService(region, ecs.ServiceName)
	suite.Nil(err)

	// This service is not listed in endpoints.AwsPartition().Services()
	// Want this to pass anyway
	err = CheckAWSRegionForService(region, ecr.ServiceName)
	suite.Nil(err)

	err = CheckAWSRegionForService(region, rds.ServiceName)
	suite.Nil(err)
}
