package models_test

import (
	"context"

	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"

	"github.com/gofrs/uuid"
)

func (suite *ModelSuite) Test_UploadCreate() {
	t := suite.T()

	upload := models.Upload{
		Filename:    "test.pdf",
		Bytes:       1048576,
		ContentType: "application/pdf",
		Checksum:    "ImGQ2Ush0bDHsaQthV5BnQ==",
		UploadType:  models.UploadTypeUSER,
	}

	verrs, err := suite.DB().ValidateAndSave(&upload)

	if err != nil {
		t.Fatalf("could not save UserUpload: %v", err)
	}

	if verrs.Count() != 0 {
		t.Errorf("did not expect UserUpload validation errors: %v", verrs)
	}
}

func (suite *ModelSuite) Test_UploadCreateWithID() {
	t := suite.T()

	id := uuid.Must(uuid.NewV4())
	upload := models.Upload{
		ID:          id,
		Filename:    "test.pdf",
		Bytes:       1048576,
		ContentType: "application/pdf",
		Checksum:    "ImGQ2Ush0bDHsaQthV5BnQ==",
		UploadType:  models.UploadTypeUSER,
	}

	verrs, err := suite.DB().ValidateAndSave(&upload)

	if err != nil {
		t.Fatalf("could not save UserUpload: %v", err)
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
		"checksum":     {"Checksum can not be blank."},
		"bytes":        {"Bytes can not be blank."},
		"filename":     {"Filename can not be blank."},
		"content_type": {"ContentType can not be blank."},
		"upload_type":  {"UploadType is not in the list [USER, PRIME]."},
	}

	suite.verifyValidationErrors(upload, expErrors)
}

func (suite *ModelSuite) TestFetchUpload() {
	t := suite.T()

	ctx := context.Background()
	document := testdatagen.MakeDefaultDocument(suite.DB())

	session := auth.Session{
		UserID:          document.ServiceMember.UserID,
		ApplicationName: auth.MilApp,
		ServiceMemberID: document.ServiceMember.ID,
	}
	upload := models.Upload{
		Filename:    "test.pdf",
		Bytes:       1048576,
		ContentType: "application/pdf",
		Checksum:    "ImGQ2Ush0bDHsaQthV5BnQ==",
		UploadType:  models.UploadTypeUSER,
	}

	verrs, err := suite.DB().ValidateAndSave(&upload)
	if err != nil {
		t.Fatalf("could not save UserUpload: %v", err)
	}
	if verrs.Count() != 0 {
		t.Errorf("did not expect UserUpload validation errors: %v", verrs)
	}

	uploadUser := models.UserUpload{
		DocumentID: &document.ID,
		UploaderID: document.ServiceMember.UserID,
		Upload:     upload,
		UploadID:   upload.ID,
	}

	verrs, err = suite.DB().ValidateAndSave(&uploadUser)
	if err != nil {
		t.Fatalf("could not save UserUpload: %v", err)
	}
	if verrs.Count() != 0 {
		t.Errorf("did not expect UserUpload validation errors: %v", verrs)
	}

	upUser, _ := models.FetchUserUpload(ctx, suite.DB(), &session, uploadUser.ID)
	suite.Equal(upUser.UploadID, upload.ID)
	suite.Equal(upUser.Upload.ID, upload.ID)
	suite.Equal(upUser.ID, uploadUser.ID)

	upUser, _ = models.FetchUserUploadFromUploadID(ctx, suite.DB(), &session, upload.ID)
	suite.Equal(upUser.UploadID, upload.ID)
	suite.Equal(upUser.Upload.ID, upload.ID)
	suite.Equal(upUser.ID, uploadUser.ID)
}

func (suite *ModelSuite) TestFetchDeletedUpload() {
	t := suite.T()

	ctx := context.Background()
	document := testdatagen.MakeDefaultDocument(suite.DB())

	session := auth.Session{
		UserID:          document.ServiceMember.UserID,
		ApplicationName: auth.MilApp,
		ServiceMemberID: document.ServiceMember.ID,
	}
	upload := models.Upload{
		Filename:    "test.pdf",
		Bytes:       1048576,
		ContentType: "application/pdf",
		Checksum:    "ImGQ2Ush0bDHsaQthV5BnQ==",
		UploadType:  models.UploadTypeUSER,
	}

	verrs, err := suite.DB().ValidateAndSave(&upload)
	if err != nil {
		t.Fatalf("could not save UserUpload: %v", err)
	}

	if verrs.Count() != 0 {
		t.Errorf("did not expect validation errors: %v", verrs)
	}

	//RA Summary: gosec - errcheck - Unchecked return value
	//RA: Linter flags errcheck error: Ignoring a method's return value can cause the program to overlook unexpected states and conditions.
	//RA: Functions with unchecked return values in the file are used fetch data and assign data to a variable that is checked later on
	//RA: Given the return value is being checked in a different line and the functions that are flagged by the linter are being used to assign variables
	//RA: in a unit test, then there is no risk
	//RA Developer Status: False Positive
	//RA Validator Status: {RA Accepted, Return to Developer, Known Issue, Mitigated, False Positive, Bad Practice}
	//RA Validator: jneuner@mitre.org
	//RA Modified Severity:
	models.DeleteUpload(suite.DB(), &upload) // nolint:errcheck
	up, _ := models.FetchUserUpload(ctx, suite.DB(), &session, upload.ID)

	// fetches a nil upload
	suite.Equal(up.ID, uuid.Nil)
}
