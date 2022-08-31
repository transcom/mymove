package models_test

import (
	"github.com/gofrs/uuid"
	"github.com/jackc/pgerrcode"

	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/db/dberr"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *ModelSuite) Test_UserUploadCreate() {
	document := testdatagen.MakeDefaultDocument(suite.DB())

	upload := models.Upload{
		Filename:    "test.pdf",
		Bytes:       1048576,
		ContentType: "application/pdf",
		Checksum:    "ImGQ2Ush0bDHsaQthV5BnQ==",
		UploadType:  models.UploadTypeUSER,
	}

	uploadUser := models.UserUpload{
		ID:         uuid.Must(uuid.NewV4()),
		DocumentID: &document.ID,
		UploaderID: document.ServiceMember.UserID,
		Upload:     upload,
	}

	var expErrors = map[string][]string{}

	suite.verifyValidationErrors(&uploadUser, expErrors)
}

func (suite *ModelSuite) Test_UserUploadValidations() {
	uploadUser := &models.UserUpload{}

	var expErrors = map[string][]string{
		"uploader_id": {"UploaderID can not be blank."},
	}

	suite.verifyValidationErrors(uploadUser, expErrors)
}

func (suite *ModelSuite) TestCreateUserUploadWithNoUpload() {
	document := testdatagen.MakeDefaultDocument(suite.DB())

	uploadUser := models.UserUpload{
		DocumentID: &document.ID,
		UploaderID: document.ServiceMember.UserID,
	}

	_, err := suite.DB().ValidateAndSave(&uploadUser)

	suite.True(dberr.IsDBErrorForConstraint(err, pgerrcode.ForeignKeyViolation, "user_uploads_uploads_id_fkey"), "expected userupload error")
}

func (suite *ModelSuite) TestFetchUserUpload() {
	userUpload := testdatagen.MakeDefaultUserUpload(suite.DB())

	session := auth.Session{
		UserID:          userUpload.Document.ServiceMember.UserID,
		ApplicationName: auth.MilApp,
		ServiceMemberID: userUpload.Document.ServiceMember.ID,
	}

	fetchedUserUpload, _ := models.FetchUserUpload(suite.DB(), &session, userUpload.ID)

	suite.Equal(fetchedUserUpload.ID, userUpload.ID)

	savedUpload := userUpload.Upload
	fetchedUserUpload, _ = models.FetchUserUploadFromUploadID(suite.DB(), &session, savedUpload.ID)

	suite.Equal(fetchedUserUpload.UploadID, savedUpload.ID)
	suite.Equal(fetchedUserUpload.Upload.ID, savedUpload.ID)
	suite.Equal(fetchedUserUpload.ID, userUpload.ID)
}

func (suite *ModelSuite) TestFetchDeletedUserUpload() {
	userUpload := testdatagen.MakeDefaultUserUpload(suite.DB())
	session := auth.Session{
		UserID:          userUpload.Document.ServiceMember.UserID,
		ApplicationName: auth.MilApp,
		ServiceMemberID: userUpload.Document.ServiceMember.ID,
	}
	err := models.DeleteUserUpload(suite.DB(), &userUpload)

	suite.Nil(err)

	userUp, err := models.FetchUserUpload(suite.DB(), &session, userUpload.ID)

	suite.Equal("error fetching user_uploads: FETCH_NOT_FOUND", err.Error())
	// fetches a nil userupload
	suite.Equal(userUp.ID, uuid.Nil)
}
