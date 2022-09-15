package services

import (
	"io"
	"time"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
)

// UploadInformation contains information for uploads
type UploadInformation struct {
	UploadID               uuid.UUID `db:"upload_id"`
	ContentType            string    `db:"content_type"`
	CreatedAt              time.Time `db:"created_at"`
	Filename               string
	Bytes                  int64
	MoveLocator            *string    `db:"locator"`
	ServiceMemberID        *uuid.UUID `db:"service_member_id"`
	ServiceMemberFirstName *string    `db:"service_member_first_name"`
	ServiceMemberLastName  *string    `db:"service_member_last_name"`
	ServiceMemberPhone     *string    `db:"service_member_telephone"`
	ServiceMemberEmail     *string    `db:"service_member_email"`
	OfficeUserID           *uuid.UUID `db:"office_user_id"`
	OfficeUserFirstName    *string    `db:"office_user_first_name"`
	OfficeUserLastName     *string    `db:"office_user_last_name"`
	OfficeUserPhone        *string    `db:"office_user_telephone"`
	OfficeUserEmail        *string    `db:"office_user_email"`
}

// UploadInformationFetcher is the service object interface for FetchUploadInformation
//go:generate mockery --name UploadInformationFetcher --disable-version-string
type UploadInformationFetcher interface {
	FetchUploadInformation(appCtx appcontext.AppContext, uuid uuid.UUID) (UploadInformation, error)
}

// UploadCreator is the service object interface for CreateUpload
//go:generate mockery --name UploadCreator --disable-version-string
type UploadCreator interface {
	CreateUpload(appCtx appcontext.AppContext, file io.ReadCloser, uploadFilename string, uploadType models.UploadType) (*models.Upload, error)
}
