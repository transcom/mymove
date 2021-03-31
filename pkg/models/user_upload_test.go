package models_test

import (
	"github.com/jackc/pgerrcode"

	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/db/dberr"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"

	"github.com/gofrs/uuid"
)

func (suite *ModelSuite) Test_UserUploadCreate() {
	t := suite.T()

	document := testdatagen.MakeDefaultDocument(suite.DB())

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
	}

	verrs, err = suite.DB().ValidateAndSave(&uploadUser)

	if err != nil {
		t.Fatalf("could not save UserUpload: %v", err)
	}

	if verrs.Count() != 0 {
		t.Errorf("did not expect UserUpload validation errors: %v", verrs)
	}
}

func (suite *ModelSuite) Test_UserUploadCreateWithID() {
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

	document := testdatagen.MakeDefaultDocument(suite.DB())

	id := uuid.Must(uuid.NewV4())
	uploadUser := models.UserUpload{
		ID:         id,
		DocumentID: &document.ID,
		UploaderID: document.ServiceMemberID,
		Upload:     upload,
	}

	verrs, err = suite.DB().ValidateAndSave(&uploadUser)

	if err != nil {
		t.Fatalf("could not save UserUpload: %v", err)
	}

	if verrs.Count() != 0 {
		t.Errorf("did not expect UserUpload validation errors: %v", verrs)
	}

	if uploadUser.ID.String() != id.String() {
		t.Errorf("wrong uuid for UserUpload: expected %s, got %s", id.String(), uploadUser.ID.String())
	}
}

func (suite *ModelSuite) Test_UserUploadValidations() {
	uploadUser := &models.UserUpload{}

	var expErrors = map[string][]string{
		"uploader_id": {"UploaderID can not be blank."},
	}

	suite.verifyValidationErrors(uploadUser, expErrors)
}

func (suite *ModelSuite) TestFetchUserUploadWithNoUpload() {
	document := testdatagen.MakeDefaultDocument(suite.DB())

	uploadUser := models.UserUpload{
		DocumentID: &document.ID,
		UploaderID: document.ServiceMember.UserID,
	}

	_, err := suite.DB().ValidateAndSave(&uploadUser)

	suite.True(dberr.IsDBErrorForConstraint(err, pgerrcode.ForeignKeyViolation, "user_uploads_uploads_id_fkey"), "expected userupload error")
}

func (suite *ModelSuite) TestFetchUserUpload() {
	t := suite.T()

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
	}

	verrs, err = suite.DB().ValidateAndSave(&uploadUser)
	if err != nil {
		t.Fatalf("could not save UserUpload: %v", err)
	}

	if verrs.Count() != 0 {
		t.Errorf("did not expect UserUpload validation errors: %v", verrs)
	}

	upUser, _ := models.FetchUserUpload(suite.DB(), &session, uploadUser.ID)
	suite.Equal(upUser.ID, uploadUser.ID)
	suite.Equal(upload.ID, uploadUser.Upload.ID)
	suite.Equal(upload.ID, uploadUser.UploadID)
}

func (suite *ModelSuite) TestFetchDeletedUserUpload() {
	t := suite.T()

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
		t.Fatalf("could not save Upload: %v", err)
	}

	if verrs.Count() != 0 {
		t.Errorf("did not expect rUpload validation errors: %v", verrs)
	}

	uploadUser := models.UserUpload{
		DocumentID: &document.ID,
		UploaderID: document.ServiceMember.UserID,
		UploadID:   upload.ID,
		Upload:     upload,
	}

	verrs, err = suite.DB().ValidateAndSave(&uploadUser)
	if err != nil {
		t.Fatalf("could not save UserUpload: %v", err)
	}

	if verrs.Count() != 0 {
		t.Errorf("did not expect validation errors: %v", verrs)
	}

	err = models.DeleteUserUpload(suite.DB(), &uploadUser)
	suite.Nil(err)
	userUp, err := models.FetchUserUpload(suite.DB(), &session, uploadUser.ID)
	suite.Equal("error fetching user_uploads: FETCH_NOT_FOUND", err.Error())

	// fetches a nil userupload
	suite.Equal(userUp.ID, uuid.Nil)
}
