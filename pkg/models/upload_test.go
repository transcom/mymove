package models_test

import (
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *ModelSuite) Test_UploadCreate() {
	t := suite.T()

	move, err := testdatagen.MakeMove(suite.db)
	if err != nil {
		t.Fatalf("could not create move: %v", err)
	}

	document := models.Document{
		UploaderID: move.UserID,
		MoveID:     move.ID,
	}
	suite.mustSave(&document)

	upload := models.Upload{
		DocumentID:  document.ID,
		UploaderID:  move.UserID,
		Filename:    "test.pdf",
		Bytes:       1048576,
		ContentType: "application/pdf",
		Checksum:    "ImGQ2Ush0bDHsaQthV5BnQ==",
	}

	verrs, err := suite.db.ValidateAndSave(&upload)

	if err != nil {
		t.Fatalf("could not save Upload: %v", err)
	}

	if verrs.Count() != 0 {
		t.Errorf("did not expect validation errors: %v", verrs)
	}
}

func (suite *ModelSuite) Test_UploadValidations() {
	upload := &models.Upload{}

	var expErrors = map[string][]string{
		"document_id":  []string{"DocumentID can not be blank."},
		"uploader_id":  []string{"UploaderID can not be blank."},
		"checksum":     []string{"Checksum can not be blank."},
		"bytes":        []string{"Bytes can not be blank."},
		"filename":     []string{"Filename can not be blank."},
		"content_type": []string{"ContentType can not be blank."},
	}

	suite.verifyValidationErrors(upload, expErrors)
}
