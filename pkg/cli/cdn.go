package cli

import (
	"fmt"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	// AwsCfDomain is domain assets use
	AwsCfDomain string = "aws-cf-domain" // #nosec so gosec doesn't claim its a hard coded cred
	// CFPrivateKeyFlag is cloudfront private key flag
	CFPrivateKeyFlag string = "cloud-front-private-key"
	// CFKeyIDFlag is cloudfront key id flag
	CFKeyIDFlag string = "cloud-front-key-id"
	// CDNBackendFlag is the CDN Backend Flag
	CDNBackendFlag string = "cdn-backend"
)

// InitCDNFlags initializes the Hosts command line flags
func InitCDNFlags(flag *pflag.FlagSet) {
	flag.String(AwsCfDomain, "", "Hostname according to environment.")
	flag.String(CFPrivateKeyFlag, "", "Cloudfront private key")
	flag.String(CFKeyIDFlag, "", "Cloudfront private key id")
	flag.String(CDNBackendFlag, "s3", "CDN backend for serving files")
}

func CheckCDNValues(v *viper.Viper) error {

	flags := []string{
		AwsCfDomain,
	}

	for _, c := range flags {
		val := v.GetString(c)
		if len(val) < 1 {
			return fmt.Errorf("%q value cannot be empty, expected length to be more than 1 found: %q", c, val)
		}
	}
	return nil
}