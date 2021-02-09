package cli

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	// AWSRegionFlag is the generic AWS Region Flag
	AWSRegionFlag string = "aws-region"
)

type errInvalidRegion struct {
	Region string
}

type errUnknownPartition struct {
	Region string
}

func (e *errInvalidRegion) Error() string {
	return fmt.Sprintf("invalid region %s", e.Region)
}

func (e *errUnknownPartition) Error() string {
	return fmt.Sprintf("unknown partition for region %s", e.Region)
}

// InitAWSFlags initializes AWS command line flags
func InitAWSFlags(flag *pflag.FlagSet) {
	flag.String(AWSRegionFlag, endpoints.UsGovWest1RegionID, "The AWS Region")
}

// CheckAWSRegion validates the AWS Region command line flags
func CheckAWSRegion(v *viper.Viper) (string, error) {
	region := v.GetString(AWSRegionFlag)
	if len(region) == 0 {
		return "", &errInvalidRegion{Region: region}
	}
	return region, nil
}

// CheckAWSRegionForService validates AWS command line flags against a region
func CheckAWSRegionForService(region, awsServiceName string) error {
	partition, ok := endpoints.PartitionForRegion(endpoints.DefaultPartitions(), region)
	if !ok {
		return &errUnknownPartition{Region: region}
	}

	if service, ok := partition.Services()[awsServiceName]; ok {
		regions := service.Regions()
		if _, ok := regions[region]; !ok {
			return &errInvalidRegion{Region: region}
		}
	}
	return nil
}
