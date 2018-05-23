package testdatagen

import (
	"log"

	"github.com/gobuffalo/pop"

	"github.com/transcom/mymove/pkg/models"
)

// MakeUpload creates a single Upload.
func MakeUpload(db *pop.Connection, document *models.Document) (models.Upload, error) {
	if document == nil {
		newDocument, err := MakeDocument(db, nil, "")
		if err != nil {
			log.Panic(err)
		}
		document = &newDocument
	}

	upload := models.Upload{
		DocumentID:  document.ID,
		Document:    *document,
		UploaderID:  document.ServiceMember.UserID,
		Filename:    "testFile.pdf",
		Bytes:       2202009,
		ContentType: "application/pdf",
		Checksum:    "ImGQ2Ush0bDHsaQthV5BnQ==",
	}

	verrs, err := db.ValidateAndSave(&upload)
	if err != nil {
		log.Panic(err)
	}
	if verrs.Count() != 0 {
		log.Panic(verrs.Error())
	}

	return upload, err
}
