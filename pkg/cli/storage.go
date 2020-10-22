package cli

import (
	"fmt"
	"path/filepath"

	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	// StorageBackendFlag is the Storage Backend Flag
	StorageBackendFlag string = "storage-backend"
	// LocalStorageRootFlag is the Local Storage Root Flag
	LocalStorageRootFlag string = "local-storage-root"
	// LocalStorageWebRootFlag is the Local Storage WebRoot Flag
	LocalStorageWebRootFlag string = "local-storage-web-root"
	// AWSS3BucketNameFlag is the AWS S3 Bucket Name Flag
	AWSS3BucketNameFlag string = "aws-s3-bucket-name"
	// AWSS3RegionFlag is the AWS S3 Region Flag
	AWSS3RegionFlag string = "aws-s3-region"
	// AWSS3KeyNamespaceFlag is the AWS S3 Key Namespace Flag
	AWSS3KeyNamespaceFlag string = "aws-s3-key-namespace"
	// AWSCfDomain is domain assets use
	AWSCfDomain string = "aws-cf-domain" //  so gosec doesn't claim its a hard coded cred
	// CFPrivateKeyFlag is cloudfront private key flag
	CFPrivateKeyFlag string = "cloud-front-private-key"
	// CFKeyIDFlag is cloudfront key id flag
	CFKeyIDFlag string = "cloud-front-key-id"
)

// InitStorageFlags initializes Storage command line flags
func InitStorageFlags(flag *pflag.FlagSet) {
	flag.String(StorageBackendFlag, "local", "Storage backend to use, either local, memory or s3.")
	flag.String(LocalStorageRootFlag, "tmp", "Local storage root directory. Default is tmp.")
	flag.String(LocalStorageWebRootFlag, "storage", "Local storage web root directory. Default is storage.")
	flag.String(AWSS3BucketNameFlag, "", "S3 bucket used for file storage")
	flag.String(AWSS3RegionFlag, "", "AWS region used for S3 file storage")
	flag.String(AWSS3KeyNamespaceFlag, "", "Key prefix for all objects written to S3")
	flag.String(AWSCfDomain, "assets.devlocal.move.mil", "Hostname according to environment.")
	flag.String(CFPrivateKeyFlag, "", "Cloudfront private key")
	flag.String(CFKeyIDFlag, "", "Cloudfront private key id")
}

// CheckStorage validates Storage command line flags
func CheckStorage(v *viper.Viper) error {

	storageBackend := v.GetString(StorageBackendFlag)
	if !stringSliceContains([]string{"local", "memory", "s3", "cdn"}, storageBackend) {
		return fmt.Errorf("invalid storage-backend %s, expecting local, memory, s3 or cdn", storageBackend)
	}

	if storageBackend == "s3" {
		r := v.GetString(AWSS3RegionFlag)
		if err := CheckAWSRegionForService(r, s3.ServiceName); err != nil {
			//return errors.Wrap(err, fmt.Sprintf("%s is invalid, value for region: %s", AWSS3RegionFlag, r))
			return nil
		}
	} else if storageBackend == "cdn" {
		privateKey := v.GetString(CFPrivateKeyFlag)
		privateKeyID := v.GetString(CFKeyIDFlag)
		cfDomain := v.GetString(AWSCfDomain)

		if len(privateKeyID) == 0 {
			return fmt.Errorf("cloudfront key id flag %q cannot be empty when using CDN for %q flag, exiting", CFKeyIDFlag, StorageBackendFlag)
		}
		if len(privateKey) == 0 {
			return fmt.Errorf("cloudfront private key flag %q cannot be empty when using CDN for %q flag, exiting", CFKeyIDFlag, StorageBackendFlag)
		}
		if len(cfDomain) == 0 {
			return fmt.Errorf("cloudfront domain flag %q cannot be empty when using CDN for %q flag, exiting", AWSCfDomain, StorageBackendFlag)
		}

	} else if storageBackend == "local" {
		localStorageRoot := v.GetString(LocalStorageRootFlag)
		if _, err := filepath.Abs(localStorageRoot); err != nil {
			return fmt.Errorf("could not get absolute path for %s", localStorageRoot)
		}
	}

	return nil
}
