package cli

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	awssession "github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/notifications"

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
		// SES is only available in 3 regions: us-east-1, us-west-2, and eu-west-1
		// - see https://docs.aws.amazon.com/ses/latest/DeveloperGuide/regions.html#region-endpoints
		if r := v.GetString("aws-ses-region"); len(r) == 0 || !stringSliceContains([]string{"us-east-1", "us-west-2", "eu-west-1"}, r) {
			return errors.Wrap(&errInvalidRegion{Region: r}, fmt.Sprintf("%s is invalid", "aws-ses-region"))
		}
		if h := v.GetString("aws-ses-domain"); len(h) == 0 {
			return errors.Wrap(&errInvalidHost{Host: h}, fmt.Sprintf("%s is invalid", "aws-ses-domain"))
		}
	}

	return nil
}

// InitEmail initializes the email backend
func InitEmail(v *viper.Viper, logger Logger) notifications.NotificationSender {
	if v.GetString(EmailBackendFlag) == "ses" {
		// Setup Amazon SES (email) service
		// TODO: This might be able to be combined with the AWS Session that we're using for S3 down
		// below.
		awsSESRegion := v.GetString(AWSSESRegionFlag)
		awsSESDomain := v.GetString(AWSSESDomainFlag)
		logger.Info("Using ses email backend",
			zap.String("region", awsSESRegion),
			zap.String("domain", awsSESDomain))
		sesSession, err := awssession.NewSession(&aws.Config{
			Region: aws.String(awsSESRegion),
		})
		if err != nil {
			logger.Fatal("Failed to create a new AWS client config provider", zap.Error(err))
		}
		sesService := ses.New(sesSession)
		return notifications.NewNotificationSender(sesService, awsSESDomain, logger)
	}

	domain := "milmovelocal"
	logger.Info("Using local email backend", zap.String("domain", domain))
	return notifications.NewStubNotificationSender(domain, logger)
}
