package testdatagen

import (
	"fmt"
	"log"

	"github.com/transcom/mymove/pkg/uploader"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/validate"

	"github.com/transcom/mymove/pkg/models"
)

// MakeUpload creates a single UserUpload.
func MakeUpload(db *pop.Connection, assertions Assertions) models.Upload {
	// Users can either assert an UserUploader (and a real file is used), or can optionally assert fields
	var upload *models.Upload
	if assertions.Uploader != nil {
		// If an UserUploader is passed in, UserUpload assertions are ignored
		var verrs *validate.Errors
		var err error
		file := fixture("test.pdf")
		upload, verrs, err = assertions.Uploader.CreateUploadForDocument(uploader.File{File: file}, uploader.AllowedTypesServiceMember)
		if verrs.HasAny() || err != nil {
			log.Panic(fmt.Errorf("errors encountered saving upload %v, %v", verrs, err))
		}
	} else {
		// If no file is being stored, use asserted fields
		upload = &models.Upload{
			Filename:    "testFile.pdf",
			Bytes:       2202009,
			ContentType: "application/pdf",
			Checksum:    "ImGQ2Ush0bDHsaQthV5BnQ==",
		}

		mergeModels(upload, assertions.Upload)

		mustCreate(db, upload)
	}

	return *upload
}

// MakeDefaultUpload makes an UserUpload with default values
func MakeDefaultUpload(db *pop.Connection) models.Upload {
	return MakeUpload(db, Assertions{})
}
