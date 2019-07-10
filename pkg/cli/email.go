package cli

import (
	"fmt"

	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/pkg/errors"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	// EmailBackendFlag is the Email Backend Flag
	EmailBackendFlag string = "email-backend"
	// AWSSESRegionFlag is the AWS SES Region Flag
	AWSSESRegionFlag string = "aws-ses-region"
	// AWSSESDomainFlag is the AWS SES Domain Flag
	AWSSESDomainFlag string = "aws-ses-domain"
)

// InitEmailFlags initializes Email command line flags
func InitEmailFlags(flag *pflag.FlagSet) {
	flag.String(EmailBackendFlag, "local", "Email backend to use, either 'ses' or 'local'")
	flag.String(AWSSESRegionFlag, "", "AWS region used for SES")
	flag.String(AWSSESDomainFlag, "", "Domain used for SES")
}

// CheckEmail validates Email command line flags
func CheckEmail(v *viper.Viper) error {
	emailBackend := v.GetString(EmailBackendFlag)
	if !stringSliceContains([]string{"local", "ses"}, emailBackend) {
		return fmt.Errorf("invalid email backend %s, expecting local or ses", emailBackend)
	}

	if emailBackend == "ses" {
		r := v.GetString(AWSSESRegionFlag)
		if err := CheckAWSRegionForService(r, ses.ServiceName); err != nil {
			return errors.Wrap(err, fmt.Sprintf("%s is invalid", AWSSESRegionFlag))
		}
		if h := v.GetString(AWSSESDomainFlag); len(h) == 0 {
			return errors.Wrap(&errInvalidHost{Host: h}, fmt.Sprintf("%s is invalid", AWSSESDomainFlag))
		}
	}

	return nil
}
