package testdatagen

import (
	"fmt"
	"log"
	"time"

	"github.com/gobuffalo/pop/v6"
	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/uploader"
)

// MakeUpload creates a single Upload.
func MakeUpload(db *pop.Connection, assertions Assertions) models.Upload {
	// Users can either assert an Uploader (and a real file is used), or can optionally assert fields
	var upload *models.Upload
	// if uploader (not a struct) is passed in
	if assertions.Uploader != nil {
		// If an Uploader is passed in, models.Upload assertions are ignored
		// because we actually upload a file using the uploader
		fmt.Printf("🔥🔥🔥")
		var verrs *validate.Errors
		var err error
		// Get file from assertions if available, or a default file if not
		file := Fixture("test.pdf")
		if assertions.File != nil {
			file = assertions.File
		}
		// Ugh. Use the global logger. All testdatagen methods should
		// take a logger
		// Save file to the database
		appCtx := appcontext.NewAppContext(db, zap.L(), nil)
		upload, verrs, err = assertions.Uploader.CreateUpload(appCtx, uploader.File{File: file}, uploader.AllowedTypesServiceMember)
		if verrs.HasAny() || err != nil {
			log.Panic(fmt.Errorf("errors encountered saving upload %v, %v", verrs, err))
		}
	} else {
		fmt.Printf("😡😡😡")
		upload = &models.Upload{
			Filename:    "testFile.pdf",
			Bytes:       int64(2202009),
			ContentType: "application/pdf",
			Checksum:    "ImGQ2Ush0bDHsaQthV5BnQ==",
			UploadType:  models.UploadTypeUSER,
		}

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
