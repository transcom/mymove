package cli

import (
	"fmt"

	"github.com/pkg/errors"
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

func (e *errInvalidRegion) Error() string {
	return fmt.Sprintf("invalid region %s", e.Region)
}

// InitAWSFlags initializes AWS command line flags
func InitAWSFlags(flag *pflag.FlagSet) {
	flag.String(AWSRegionFlag, "us-west-2", "The AWS Region")
}

// CheckAWS validates AWS command line flags
func CheckAWS(v *viper.Viper) error {
	region := v.GetString(AWSRegionFlag)
	if len(region) == 0 {
		return errors.Wrap(&errInvalidRegion{Region: region}, fmt.Sprintf("%q is invalid", AWSRegionFlag))
	}
	return nil
}
