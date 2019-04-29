package cli

import (
	"fmt"
	"path"
	"path/filepath"

	"github.com/aws/aws-sdk-go/service/s3"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/endpoints"
	awssession "github.com/aws/aws-sdk-go/aws/session"
	"github.com/pkg/errors"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/storage"
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
		return fmt.Errorf("invalid storage-backend %s, expecting local, memory or s3", storageBackend)
	}

	if storageBackend == "s3" {
		r := v.GetString(AWSS3RegionFlag)
		if len(r) == 0 {
			return errors.Wrap(&errInvalidRegion{Region: r}, fmt.Sprintf("%s is invalid", "aws-s3-region"))
		}

		regions := endpoints.AwsPartition().Services()[s3.ServiceName].Regions()
		if _, ok := regions[r]; !ok {
			return errors.Wrap(&errInvalidRegion{Region: r}, fmt.Sprintf("%s is invalid", "aws-s3-region"))
		}
	} else if storageBackend == "local" {
		localStorageRoot := v.GetString(LocalStorageRootFlag)
		if _, err := filepath.Abs(localStorageRoot); err != nil {
			return fmt.Errorf("could not get absolute path for %s", localStorageRoot)
		}
	}

	return nil
}

// InitStorage initializes the storage backend
func InitStorage(v *viper.Viper, logger Logger) storage.FileStorer {
	storageBackend := v.GetString(StorageBackendFlag)
	localStorageRoot := v.GetString(LocalStorageRootFlag)
	localStorageWebRoot := v.GetString(LocalStorageWebRootFlag)

	var storer storage.FileStorer
	if storageBackend == "s3" {
		awsS3Bucket := v.GetString(AWSS3BucketNameFlag)
		awsS3Region := v.GetString(AWSS3RegionFlag)
		awsS3KeyNamespace := v.GetString(AWSS3KeyNamespaceFlag)
		logger.Info("Using s3 storage backend",
			zap.String("bucket", awsS3Bucket),
			zap.String("region", awsS3Region),
			zap.String("key", awsS3KeyNamespace))
		if len(awsS3Bucket) == 0 {
			logger.Fatal("must provide aws-s3-bucket-name parameter, exiting")
		}
		if len(awsS3Region) == 0 {
			logger.Fatal("Must provide aws-s3-region parameter, exiting")
		}
		if len(awsS3KeyNamespace) == 0 {
			logger.Fatal("Must provide aws_s3_key_namespace parameter, exiting")
		}
		aws := awssession.Must(awssession.NewSession(&aws.Config{
			Region: aws.String(awsS3Region),
		}))
		storer = storage.NewS3(awsS3Bucket, awsS3KeyNamespace, logger, aws)
	} else if storageBackend == "memory" {
		logger.Info("Using memory storage backend",
			zap.String(LocalStorageRootFlag, path.Join(localStorageRoot, localStorageWebRoot)),
			zap.String(LocalStorageWebRootFlag, localStorageWebRoot))
		fsParams := storage.NewMemoryParams(localStorageRoot, localStorageWebRoot, logger)
		storer = storage.NewMemory(fsParams)
	} else {
		logger.Info("Using local storage backend",
			zap.String(LocalStorageRootFlag, path.Join(localStorageRoot, localStorageWebRoot)),
			zap.String(LocalStorageWebRootFlag, localStorageWebRoot))
		fsParams := storage.NewFilesystemParams(localStorageRoot, localStorageWebRoot, logger)
		storer = storage.NewFilesystem(fsParams)
	}
	return storer
}
