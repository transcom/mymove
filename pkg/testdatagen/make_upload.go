package testdatagen

import (
	"fmt"
	"log"

	"github.com/gobuffalo/pop/v6"
	"github.com/gobuffalo/validate/v3"
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
		appCtx := appcontext.NewAppContext(db, zap.L(), nil, nil)
		upload, verrs, err = assertions.Uploader.CreateUpload(appCtx, uploader.File{File: file}, uploader.AllowedTypesServiceMember)
		if verrs.HasAny() || err != nil {
			log.Panic(fmt.Errorf("errors encountered saving upload %v, %v", verrs, err))
		}
	} else {
		upload = &models.Upload{
			Filename:    "testFile.pdf",
			Bytes:       int64(2202009),
			ContentType: uploader.FileTypePDF,
			Checksum:    "ImGQ2Ush0bDHsaQthV5BnQ==",
			UploadType:  models.UploadTypeUSER,
		}

		mergeModels(upload, assertions.Upload)

		mustCreate(db, upload, assertions.Stub)
	}

	return *upload
}
