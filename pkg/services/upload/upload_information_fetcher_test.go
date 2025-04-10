package upload

import (
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/models"
)

func (suite *UploadServiceSuite) TestFetchUploadInformation() {
	suite.Run("fetch office user upload", func() {
		email := "officeuser1@example.com"
		ou := factory.BuildOfficeUser(suite.DB(), []factory.Customization{
			{
				Model: models.User{
					ID:        uuid.Must(uuid.FromString("9bfa91d2-7a0c-4de0-ae02-b8cf8b4b858b")),
					OktaEmail: email,
				},
			},
			{
				Model: models.OfficeUser{
					ID:        uuid.FromStringOrNil("9c5911a7-5885-4cf4-abec-021a40692403"),
					Email:     email,
					FirstName: "Office",
					LastName:  "User",
					Telephone: "212-312-1234",
				},
			},
		}, nil)
		uu := factory.BuildUserUpload(suite.DB(), []factory.Customization{
			{
				Model: models.UserUpload{
					UploaderID: *ou.UserID,
				},
			},
		}, nil)
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
		uu := factory.BuildUserUpload(suite.DB(), nil, nil)
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
