package services

import (
	"time"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
)

// UploadFetcher is the service object interface for FetchUploads
//go:generate mockery -name UploadFetcher
type UploadFetcher interface {
	FetchUploads(filters []QueryFilter, associations QueryAssociations, pagination Pagination) (models.Uploads, error)
}

type UploadInformation struct {
	UploadID        uuid.UUID `db:"upload_id"`
	ContentType     string    `db:"content_type"`
	CreatedAt       time.Time `db:"created_at"`
	Filename        string
	Bytes           int64
	MoveLocator     *string    `db:"locator"`
	ServiceMemberID *uuid.UUID `db:"service_member_id"`
	OfficeUserID    *uuid.UUID `db:"office_user_id"`
	OfficeUserEmail *string    `db:"office_user_email"`
}

//go:generate mockery -name UploadInformationFetcher
type UploadInformationFetcher interface {
	FetchUploadInformation(uuid uuid.UUID) (UploadInformation, error)
}
