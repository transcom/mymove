package upload

import (
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *UploadServiceSuite) TestFetchUploadInformation() {
	suite.Run("fetch office user upload", func() {
		email := "officeuser1@example.com"
		ou := testdatagen.MakeOfficeUser(suite.DB(), testdatagen.Assertions{
			User: models.User{
				ID:            uuid.Must(uuid.FromString("9bfa91d2-7a0c-4de0-ae02-b8cf8b4b858b")),
				LoginGovEmail: email,
			},
			OfficeUser: models.OfficeUser{
				ID:        uuid.FromStringOrNil("9c5911a7-5885-4cf4-abec-021a40692403"),
				Email:     email,
				FirstName: "Office",
				LastName:  "User",
				Telephone: "212-312-1234",
			},
		})

		assertions := testdatagen.Assertions{UserUpload: models.UserUpload{UploaderID: *ou.UserID}}
		uu := testdatagen.MakeUserUpload(suite.DB(), assertions)
		uif := NewUploadInformationFetcher()
		suite.NotNil(uu.UploadID)
		suite.NotNil(uu.Upload)
		u := uu.Upload
		ui, err := uif.FetchUploadInformation(suite.AppContextForTest(), u.ID)

		suite.NoError(err)
		suite.Nil(ui.ServiceMemberID)
		suite.Equal(ou.ID, *ui.OfficeUserID)
		suite.Equal(ou.Email, *ui.OfficeUserEmail)
		suite.Equal(ou.FirstName, *ui.OfficeUserFirstName)
		suite.Equal(ou.LastName, *ui.OfficeUserLastName)
		suite.Equal(ou.Telephone, *ui.OfficeUserPhone)
	})

	suite.Run("fetch service member upload", func() {
		uu := testdatagen.MakeDefaultUserUpload(suite.DB())
		uif := NewUploadInformationFetcher()
		suite.NotNil(uu.UploadID)
		suite.NotNil(uu.Upload)
		u := uu.Upload
		ui, err := uif.FetchUploadInformation(suite.AppContextForTest(), u.ID)

		suite.NoError(err)
		suite.Nil(ui.OfficeUserID)
		sm := uu.Document.ServiceMember
		suite.Equal(sm.ID, *ui.ServiceMemberID)
		suite.Equal(*sm.PersonalEmail, *ui.ServiceMemberEmail)
		suite.Equal(*sm.FirstName, *ui.ServiceMemberFirstName)
		suite.Equal(*sm.LastName, *ui.ServiceMemberLastName)
		suite.Equal(*sm.Telephone, *ui.ServiceMemberPhone)
		suite.Equal(u.ID, ui.UploadID)
		suite.Equal(u.ContentType, ui.ContentType)
		suite.Equal(u.Bytes, ui.Bytes)
		suite.Equal(u.Filename, ui.Filename)
	})
}
