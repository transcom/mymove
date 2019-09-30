package services

import (
	"github.com/transcom/mymove/pkg/models"
)

// UploadFetcher is the service object interface for FetchUploads
//go:generate mockery -name UploadFetcher
type UploadFetcher interface {
	FetchUploads(filters []QueryFilter, associations QueryAssociations) (models.Uploads, error)
}