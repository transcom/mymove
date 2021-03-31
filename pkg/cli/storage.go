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
)

// InitStorageFlags initializes Storage command line flags
func InitStorageFlags(flag *pflag.FlagSet) {
	flag.String(StorageBackendFlag, "local", "Storage backend to use, either local, memory or s3.")
	flag.String(LocalStorageRootFlag, "tmp", "Local storage root directory. Default is tmp.")
	flag.String(LocalStorageWebRootFlag, "storage", "Local storage web root directory. Default is storage.")
	flag.String(AWSS3BucketNameFlag, "", "S3 bucket used for file storage")
	flag.String(AWSS3RegionFlag, "", "AWS region used for S3 file storage")
	flag.String(AWSS3KeyNamespaceFlag, "", "Key prefix for all objects written to S3")
}

// CheckStorage validates Storage command line flags
func CheckStorage(v *viper.Viper) error {

	storageBackend := v.GetString(StorageBackendFlag)
	if !stringSliceContains([]string{"local", "memory", "s3"}, storageBackend) {
		return fmt.Errorf("invalid storage-backend %s, expecting local, memory, or s3", storageBackend)
	}

	if storageBackend == "s3" {
		r := v.GetString(AWSS3RegionFlag)
		if err := CheckAWSRegionForService(r, s3.ServiceName); err != nil {
			//return errors.Wrap(err, fmt.Sprintf("%s is invalid, value for region: %s", AWSS3RegionFlag, r))
			return nil
		}
	} else if storageBackend == "local" {
		localStorageRoot := v.GetString(LocalStorageRootFlag)
		if _, err := filepath.Abs(localStorageRoot); err != nil {
			return fmt.Errorf("could not get absolute path for %s", localStorageRoot)
		}
	}

	return nil
}
