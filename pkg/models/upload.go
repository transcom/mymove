package models

import (
	"path"
	"time"

	"github.com/gobuffalo/pop/v6"
	"github.com/gobuffalo/validate/v3"
	"github.com/gobuffalo/validate/v3/validators"
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"

	"github.com/transcom/mymove/pkg/db/utilities"
)

// UploadType represents the type of upload this is, whether is it uploaded for a User or for the Prime
type UploadType string

const (
	// UploadTypeUSER string USER
	UploadTypeUSER UploadType = "USER"
	// UploadTypePRIME string PRIME
	UploadTypePRIME UploadType = "PRIME"
	// UploadTypePRIME string OFFICE
	UploadTypeOFFICE UploadType = "OFFICE"
)

// AVStatusType represents the type of the anti-virus status, whether it is still processing, clean or infected
type AVStatusType string

const (
	// AVStatusTypePROCESSING string PROCESSING
	AVStatusTypePROCESSING AVStatusType = "PROCESSING"
	// AVStatusTypeCLEAN string CLEAN
	AVStatusTypeCLEAN AVStatusType = "CLEAN"
	// AVStatusTypeINFECTED string INFECTED
	AVStatusTypeINFECTED AVStatusType = "INFECTED"
)

// An Upload represents an uploaded file, such as an image or PDF.
type Upload struct {
	ID          uuid.UUID     `db:"id"`
	Filename    string        `db:"filename"`
	Bytes       int64         `db:"bytes"`
	Rotation    *int64        `db:"rotation"`
	ContentType string        `db:"content_type"`
	Checksum    string        `db:"checksum"`
	StorageKey  string        `db:"storage_key"`
	AVStatus    *AVStatusType `db:"av_status"`
	UploadType  UploadType    `db:"upload_type"`
	CreatedAt   time.Time     `db:"created_at"`
	UpdatedAt   time.Time     `db:"updated_at"`
	DeletedAt   *time.Time    `db:"deleted_at"`
}

// TableName overrides the table name used by Pop.
func (u Upload) TableName() string {
	return "uploads"
}

type Uploads []Upload

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
func (u *Upload) Validate(_ *pop.Connection) (*validate.Errors, error) {
	var vs []validate.Validator
	vs = append(vs, &validators.StringInclusion{Field: string(u.UploadType), Name: "UploadType", List: []string{
		string(UploadTypeUSER),
		string(UploadTypePRIME),
		string(UploadTypeOFFICE),
	}})
	vs = append(vs,
		&validators.StringIsPresent{Field: u.Filename, Name: "Filename"},
		&Int64IsPresent{Field: u.Bytes, Name: "Bytes"},
		&validators.StringIsPresent{Field: u.ContentType, Name: "ContentType"},
		&validators.StringIsPresent{Field: u.Checksum, Name: "Checksum"},
	)
	return validate.Validate(vs...), nil
}

// BeforeCreate populates the StorageKey on a newly created UserUpload
func (u *Upload) BeforeCreate(_ *pop.Connection) error {
	// Populate ID if not exists
	if u.ID == uuid.Nil {
		u.ID = uuid.Must(uuid.NewV4())
	}

	if u.StorageKey == "" {
		u.StorageKey = path.Join(string(u.UploadType), "uploads", u.ID.String())
	}

	return nil
}

func FetchUpload(dbConn *pop.Connection, uploadID uuid.UUID) (*Upload, error) {
	var upload Upload
	err := dbConn.Q().Find(&upload, uploadID)
	if err != nil {
		if errors.Cause(err).Error() == RecordNotFoundErrorString {
			return &Upload{}, ErrFetchNotFound
		}
		// Otherwise, it's an unexpected err so we return that.
		return &Upload{}, err
	}

	return &upload, nil
}

// DeleteUpload deletes an upload from the database
func DeleteUpload(dbConn *pop.Connection, upload *Upload) error {
	if dbConn.TX != nil {
		err := utilities.SoftDestroy(dbConn, upload)
		if err != nil {
			return err
		}
	} else {
		return dbConn.Transaction(func(db *pop.Connection) error {
			err := utilities.SoftDestroy(db, upload)
			if err != nil {
				return err
			}
			return nil
		})
	}
	return nil
}

// Valid checks if UploadType is set to a valid value
func (ut UploadType) Valid() bool {
	for _, value := range []string{
		string(UploadTypeUSER),
		string(UploadTypePRIME),
		string(UploadTypeOFFICE),
	} {
		if string(ut) == value {
			return true
		}
	}
	return false
}
