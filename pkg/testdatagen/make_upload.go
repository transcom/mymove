package testdatagen

import (
	"log"

	"github.com/markbates/pop"
	"github.com/satori/go.uuid"

	"github.com/transcom/mymove/pkg/models"
)

// MakeUpload creates a single User.
func MakeUpload(db *pop.Connection, document *models.Document) (models.Upload, error) {
	id, err := uuid.NewV4()
	if err != nil {
		log.Panic(err)
	}

	if document == nil {
		newDocument, err := MakeDocument(db, nil)
		if err != nil {
			log.Panic(err)
		}
		document = &newDocument
	}

	upload := models.Upload{
		DocumentID:  document.ID,
		UploaderID:  document.UploaderID,
		Filename:    "testFile.pdf",
		Bytes:       2202009,
		ContentType: "application/pdf",
		Checksum:    "ImGQ2Ush0bDHsaQthV5BnQ==",
		S3ID:        id,
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
