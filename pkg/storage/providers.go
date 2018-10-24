package storage

import (
	"github.com/transcom/mymove/pkg/di"
	"go.uber.org/zap"
)

// NewStorageProvider constructs either a local FileStorer or an S3 based one depending on the presence of S3 config
func NewStorageProvider(cfg *S3StorerConfig, l *zap.Logger) (FileStorer, error) {
	if cfg == nil {
		l.Info("Using filesystem storage backend")
		fsParams := DefaultFilesystemParams(l)
		return NewFilesystem(fsParams), nil
	}
	l.Info("Using s3 storage backend")
	return NewS3(cfg, l)
}

// AddProviders adds the DI providers to generate storage objectss
func AddProviders(c *di.Container) {
	c.MustProvide(NewStorageProvider)
}
