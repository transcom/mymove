package storage

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"io"
	"path"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/gabriel-vasile/mimetype"
	"github.com/pkg/errors"
	"github.com/spf13/afero"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/cli"
)

// StoreResult represents the result of a call to Store().
type StoreResult struct{}

// FileStorer is the set of methods needed to store and retrieve objects.
//
//go:generate mockery --name FileStorer
type FileStorer interface {
	Store(string, io.ReadSeeker, string, *string) (*StoreResult, error)
	Fetch(string) (io.ReadCloser, error)
	Delete(string) error
	PresignedURL(string, string, string) (string, error)
	FileSystem() *afero.Afero
	TempFileSystem() *afero.Afero
	Tags(string) (map[string]string, error)
}

// ComputeChecksum calculates the SHA256 checksum for the provided data. It expects that
// the passed io object will be seeked to its beginning and will seek back to the
// beginning after reading its content.
func ComputeChecksum(data io.ReadSeeker) (string, error) {
	hash := sha256.New()
	if _, err := io.Copy(hash, data); err != nil {
		return "", errors.Wrap(err, "could not read file")
	}

	if _, err := data.Seek(0, io.SeekStart); err != nil { // seek back to beginning of file
		return "", errors.Wrap(err, "could not seek to beginning of file")
	}
	return base64.StdEncoding.EncodeToString(hash.Sum(nil)), nil
}

// DetectContentType leverages http.DetectContentType to identify the content type
// of the provided data. It expects that the passed io object will be seeked to its
// beginning and will seek back to the beginning after reading its content.
func DetectContentType(data io.ReadSeeker) (string, error) {
	if _, err := data.Seek(0, io.SeekStart); err != nil { // seek back to beginning of file
		return "", errors.Wrap(err, "could not seek to beginning of file")
	}

	// the default return value will default to application/octet-stream if unable to detect the MIME type
	contentType, readErr := mimetype.DetectReader(data)
	if readErr != nil {
		return "", errors.Wrap(readErr, "encountered error reading file content type")
	}

	if _, err := data.Seek(0, io.SeekStart); err != nil { // seek back to beginning of file
		return "", errors.Wrap(err, "could not seek to beginning of file")
	}
	return contentType.String(), nil
}

// InitStorage initializes the storage backend
func InitStorage(v *viper.Viper, logger *zap.Logger) FileStorer {
	storageBackend := v.GetString(cli.StorageBackendFlag)
	localStorageRoot := v.GetString(cli.LocalStorageRootFlag)
	localStorageWebRoot := v.GetString(cli.LocalStorageWebRootFlag)

	var storer FileStorer
	if storageBackend == "s3" {
		awsS3Bucket := v.GetString(cli.AWSS3BucketNameFlag)
		awsS3Region := v.GetString(cli.AWSS3RegionFlag)
		awsS3KeyNamespace := v.GetString(cli.AWSS3KeyNamespaceFlag)

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
		cfg, err := config.LoadDefaultConfig(context.Background(),
			config.WithRegion(awsS3Region),
		)
		if err != nil {
			logger.Fatal("error loading S3 aws config", zap.Error(err))
		}

		storer = NewS3(awsS3Bucket, awsS3KeyNamespace, cfg)
	} else if storageBackend == "memory" {
		logger.Info("Using memory storage backend",
			zap.String(cli.LocalStorageRootFlag, path.Join(localStorageRoot, localStorageWebRoot)),
			zap.String(cli.LocalStorageWebRootFlag, localStorageWebRoot))
		fsParams := NewMemoryParams(localStorageRoot, localStorageWebRoot)
		storer = NewMemory(fsParams)
	} else {
		logger.Info("Using local storage backend",
			zap.String(cli.LocalStorageRootFlag, path.Join(localStorageRoot, localStorageWebRoot)),
			zap.String(cli.LocalStorageWebRootFlag, localStorageWebRoot))
		fsParams := NewFilesystemParams(localStorageRoot, localStorageWebRoot)
		storer = NewFilesystem(fsParams)
	}
	return storer
}
