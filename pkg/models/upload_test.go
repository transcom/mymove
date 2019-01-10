package models_test

import (
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"

	"github.com/gofrs/uuid"
)

func (suite *ModelSuite) Test_UploadCreate() {
	t := suite.T()

	document := testdatagen.MakeDefaultDocument(suite.DB())

	upload := models.Upload{
		DocumentID:  &document.ID,
		UploaderID:  document.ServiceMember.UserID,
		Filename:    "test.pdf",
		Bytes:       1048576,
		ContentType: "application/pdf",
		Checksum:    "ImGQ2Ush0bDHsaQthV5BnQ==",
	}

	verrs, err := suite.DB().ValidateAndSave(&upload)

	if err != nil {
		t.Fatalf("could not save Upload: %v", err)
	}

	if verrs.Count() != 0 {
		t.Errorf("did not expect validation errors: %v", verrs)
	}
}

func (suite *ModelSuite) Test_UploadCreateWithID() {
	t := suite.T()

	document := testdatagen.MakeDefaultDocument(suite.DB())

	id := uuid.Must(uuid.NewV4())
	upload := models.Upload{
		ID:          id,
		DocumentID:  &document.ID,
		UploaderID:  document.ServiceMemberID,
		Filename:    "test.pdf",
		Bytes:       1048576,
		ContentType: "application/pdf",
		Checksum:    "ImGQ2Ush0bDHsaQthV5BnQ==",
	}

	verrs, err := suite.DB().ValidateAndSave(&upload)

	if err != nil {
		t.Fatalf("could not save Upload: %v", err)
	}

	if verrs.Count() != 0 {
		t.Errorf("did not expect validation errors: %v", verrs)
	}

	if upload.ID.String() != id.String() {
		t.Errorf("wrong uuid for upload: expected %s, got %s", id.String(), upload.ID.String())
	}
}

func (suite *ModelSuite) Test_UploadValidations() {
	upload := &models.Upload{}

	var expErrors = map[string][]string{
		"uploader_id":  {"UploaderID can not be blank."},
		"checksum":     {"Checksum can not be blank."},
		"bytes":        {"Bytes can not be blank."},
		"filename":     {"Filename can not be blank."},
		"content_type": {"ContentType is not in the list [image/jpeg, image/png, application/pdf, text/plain, text/plain; charset=utf-8]."},
	}

	suite.verifyValidationErrors(upload, expErrors)
}
