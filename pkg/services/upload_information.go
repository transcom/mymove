package services

import (
	"time"

	"github.com/gofrs/uuid"
)

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

// UploadInformationFetcher is the service object interface for FetchUploadInformation
//go:generate mockery -name UploadInformationFetcher
type UploadInformationFetcher interface {
	FetchUploadInformation(uuid uuid.UUID) (UploadInformation, error)
}
