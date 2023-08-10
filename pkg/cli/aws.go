package cli

import (
	"fmt"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	// AWSRegionFlag is the generic AWS Region Flag
	AWSRegionFlag string = "aws-region"

	// aws-sdk-go-v2 does not expose a constant for the regions anymore
	AWSRegionUSGovWest1 = "us-gov-west-1"
)

type errInvalidRegion struct {
	Region string
}

func (e *errInvalidRegion) Error() string {
	return fmt.Sprintf("invalid region %s", e.Region)
}

// InitAWSFlags initializes AWS command line flags
func InitAWSFlags(flag *pflag.FlagSet) {
	flag.String(AWSRegionFlag, AWSRegionUSGovWest1, "The AWS Region")
}

// CheckAWSRegion validates the AWS Region command line flags
func CheckAWSRegion(v *viper.Viper) (string, error) {
	region := v.GetString(AWSRegionFlag)
	if len(region) == 0 {
		return "", &errInvalidRegion{Region: region}
	}
	return region, nil
}
