package testdatagen

import (
	"github.com/gobuffalo/pop"

	"github.com/transcom/mymove/pkg/models"
)

// MakeUpload creates a single Upload.
func MakeUpload(db *pop.Connection, assertions Assertions) models.Upload {
	document := assertions.Upload.Document
	if isZeroUUID(assertions.Upload.DocumentID) {
		document = MakeDocument(db, assertions)
	}

	upload := models.Upload{
		DocumentID:  document.ID,
		Document:    document,
		UploaderID:  document.ServiceMember.UserID,
		Filename:    "testFile.pdf",
		Bytes:       2202009,
		ContentType: "application/pdf",
		Checksum:    "ImGQ2Ush0bDHsaQthV5BnQ==",
	}

	mergeModels(&upload, assertions.Upload)

	mustSave(db, &upload)

	return upload
}

// MakeDefaultUpload makes an Upload with default values
func MakeDefaultUpload(db *pop.Connection) models.Upload {
	return MakeUpload(db, Assertions{})
}
