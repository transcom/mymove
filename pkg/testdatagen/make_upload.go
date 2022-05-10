package testdatagen

import (
	"fmt"
	"log"
	"time"

	"github.com/gofrs/uuid"

	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/uploader"

	"github.com/gobuffalo/pop/v6"
	"github.com/gobuffalo/validate/v3"

	"github.com/transcom/mymove/pkg/models"
)

// MakeUpload creates a single Upload.
func MakeUpload(db *pop.Connection, assertions Assertions) models.Upload {
	// Users can either assert an Uploader (and a real file is used), or can optionally assert fields
	var upload *models.Upload
	if assertions.Uploader != nil {
		// If an Uploader is passed in, Upload assertions are ignored
		var verrs *validate.Errors
		var err error
		file := Fixture("test.pdf")
		if assertions.File != nil {
			file = assertions.File
		}
		// Ugh. Use the global logger. All testdatagen methods should
		// take a logger
		appCtx := appcontext.NewAppContext(db, zap.L(), nil)
		upload, verrs, err = assertions.Uploader.CreateUpload(appCtx, uploader.File{File: file}, uploader.AllowedTypesServiceMember)
		if verrs.HasAny() || err != nil {
			log.Panic(fmt.Errorf("errors encountered saving upload %v, %v", verrs, err))
		}
	} else {
		// If no file is being stored, use asserted fields
		upload = &models.Upload{}

		filename := "testFile.pdf"
		if assertions.Upload.Filename != "" {
			filename = assertions.Upload.Filename
		}
		upload.Filename = filename

		bytes := int64(2202009)
		if assertions.UploadUseZeroBytes {
			bytes = 0
		} else if assertions.Upload.Bytes > 0 {
			bytes = assertions.Upload.Bytes
		}
		upload.Bytes = bytes

		contentType := "application/pdf"
		if assertions.Upload.ContentType != "" {
			contentType = assertions.Upload.ContentType
		}
		upload.ContentType = contentType

		checksum := "ImGQ2Ush0bDHsaQthV5BnQ=="
		if assertions.Upload.Checksum != "" {
			checksum = assertions.Upload.Checksum
		}
		upload.Checksum = checksum

		uploadType := models.UploadTypeUSER
		if assertions.Upload.UploadType.Valid() {
			uploadType = assertions.Upload.UploadType
		}
		upload.UploadType = uploadType

		mergeModels(upload, assertions.Upload)

		mustCreate(db, upload, assertions.Stub)
	}

	return *upload
}

// MakeDefaultUpload makes an Upload with default values
func MakeDefaultUpload(db *pop.Connection) models.Upload {
	return MakeUpload(db, Assertions{})
}

// MakeStubbedUpload makes a fake Upload that is not saved to the DB
func MakeStubbedUpload(db *pop.Connection, assertions Assertions) models.Upload {
	assertions.Stub = true
	assertions.Upload.ID = uuid.Must(uuid.NewV4())
	assertions.Upload.CreatedAt = time.Now()
	assertions.Upload.UpdatedAt = time.Now()
	return MakeUpload(db, assertions)
}
