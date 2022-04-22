package testdatagen

import (
	"fmt"
	"log"

	"github.com/gobuffalo/pop/v5"
	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/uploader"
)

// MakeUserUpload creates a single UserUpload.
func MakeUserUpload(db *pop.Connection, assertions Assertions) models.UserUpload {
	document := assertions.UserUpload.Document
	if assertions.UserUpload.DocumentID == nil || isZeroUUID(*assertions.UserUpload.DocumentID) {
		document = MakeDocument(db, assertions)
	}

	uploaderID := assertions.UserUpload.UploaderID
	if isZeroUUID(uploaderID) {
		uploaderID = document.ServiceMember.UserID
	}

	// Users can either assert an UserUploader (and a real file is used), or can optionally assert fields
	var userUpload *models.UserUpload
	if assertions.UserUploader != nil {
		// If an UserUploader is passed in, UserUpload assertions are ignored
		var err error
		var verrs *validate.Errors
		file := Fixture("test.pdf")
		if assertions.File != nil {
			file = assertions.File
		}
		// Ugh. Use the global logger. All testdatagen methods should
		// take a logger
		appCtx := appcontext.NewAppContext(db, zap.L(), nil)
		userUpload, verrs, err = assertions.UserUploader.CreateUserUploadForDocument(appCtx, &document.ID, uploaderID, uploader.File{File: file}, uploader.AllowedTypesServiceMember)
		if verrs.HasAny() || err != nil {
			log.Panic(fmt.Errorf("errors encountered saving user upload %v, %v", verrs, err))
		}
		userUpload.Document = document
	} else {
		// If no UserUploader is being stored, use asserted fields

		if assertions.UserUpload.Upload.ID != uuid.Nil {
			assertions.Upload = assertions.UserUpload.Upload
		}
		assertions.Upload.UploadType = models.UploadTypeUSER
		upload := MakeUpload(db, assertions)

		userUpload = &models.UserUpload{
			DocumentID: &document.ID,
			Document:   document,
			UploaderID: uploaderID,
			Upload:     upload,
			UploadID:   upload.ID,
		}

		mergeModels(userUpload, assertions.UserUpload)

		mustCreate(db, userUpload, assertions.Stub)
	}

	return *userUpload
}

// MakeDefaultUserUpload makes an UserUpload with default values
func MakeDefaultUserUpload(db *pop.Connection) models.UserUpload {
	return MakeUserUpload(db, Assertions{})
}
