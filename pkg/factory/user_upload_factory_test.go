package factory

import (
	"time"

	"github.com/transcom/mymove/pkg/models"
	storageTest "github.com/transcom/mymove/pkg/storage/test"
	"github.com/transcom/mymove/pkg/uploader"
)

func (suite *FactorySuite) TestBuildUserUpload() {
	suite.Run("Successful creation of default user upload", func() {
		// Under test:      BuildUserUpload
		// Set up:          Create a default user upload
		// Expected outcome:Create a document and upload
		// This file doesn't actually exist

		// Create user upload
		userUpload := BuildUserUpload(suite.DB(), nil, nil)

		suite.False(userUpload.DocumentID.IsNil())
		suite.False(userUpload.Document.ID.IsNil())
		suite.Equal(userUpload.Document.ServiceMember.UserID,
			userUpload.UploaderID)
		suite.False(userUpload.UploadID.IsNil())
		suite.False(userUpload.Upload.ID.IsNil())
	})

	suite.Run("Successful creation of customized user upload with linked document", func() {
		// Under test:       BuildUserUpload
		// Set up:           Create a customized upload (no uploader)
		// Expected outcome: All fields should match

		// For custom UploaderID, need a valid UserID
		user := BuildUser(suite.DB(), nil, nil)

		// Create document
		document := BuildDocument(suite.DB(), nil, nil)

		// make sure these are not equal to test customization
		suite.NotEqual(user.ID, document.ServiceMember.UserID)

		// Create upload
		upload := BuildUpload(suite.DB(), nil, nil)

		customUserUpload := models.UserUpload{
			UploaderID: user.ID,
		}

		// Create user upload
		userUpload := BuildUserUpload(suite.DB(), []Customization{
			{
				Model: customUserUpload,
			},
			{
				Model:    upload,
				LinkOnly: true,
			},
			{
				Model:    document,
				LinkOnly: true,
			},
		}, nil)
		suite.Equal(document.ID, *userUpload.DocumentID)
		suite.Equal(document, userUpload.Document)
		suite.Equal(customUserUpload.UploaderID, userUpload.UploaderID)
		suite.Equal(upload.ID, userUpload.UploadID)
		suite.Equal(upload, userUpload.Upload)
	})

	suite.Run("Successful creation of customized user upload with linked service member", func() {
		// Under test:       BuildUserUpload
		// Set up:           Create a customized service member
		// Expected outcome: All fields should match

		serviceMember := BuildServiceMember(suite.DB(), nil, nil)

		// Create user upload
		userUpload := BuildUserUpload(suite.DB(), []Customization{
			{
				Model:    serviceMember,
				LinkOnly: true,
			},
		}, nil)
		suite.Equal(serviceMember.ID, userUpload.Document.ServiceMemberID)
		suite.Equal(serviceMember, userUpload.Document.ServiceMember)
		suite.Equal(serviceMember.UserID, userUpload.UploaderID)
	})

	suite.Run("Successful creation of customized user upload with customized upload", func() {
		// Under test:       BuildUserUpload
		// Set up:           Create a customized upload with no uploader
		// Expected outcome: All fields should match

		customUpload := models.Upload{
			Filename:    "BaisWinery.jpg",
			Bytes:       int64(6081979),
			ContentType: "application/jpg",
			Checksum:    "GauMarJosbDHsaQthV5BnQ==",
			CreatedAt:   time.Now(),
		}

		// Create user upload
		userUpload := BuildUserUpload(suite.DB(), []Customization{
			{
				Model: customUpload,
			},
		}, nil)

		suite.Equal(customUpload.Filename, userUpload.Upload.Filename)
		suite.Equal(customUpload.Bytes, userUpload.Upload.Bytes)
		suite.Equal(customUpload.ContentType, userUpload.Upload.ContentType)
		suite.Equal(customUpload.Checksum, userUpload.Upload.Checksum)
		suite.Equal(customUpload.CreatedAt, userUpload.Upload.CreatedAt)
		suite.False(userUpload.DocumentID.IsNil())
		suite.False(userUpload.Document.ID.IsNil())
		suite.False(userUpload.UploadID.IsNil())
		suite.False(userUpload.Upload.ID.IsNil())
	})

	suite.Run("Successful creation of user upload with basic uploader", func() {
		// Under test:      BuildUserUpload
		// Mocked:          None
		// Set up:          Create an upload with an uploader and default file
		// Expected outcome:Upload filename should be the default file
		storer := storageTest.NewFakeS3Storage(true)
		userUploader, err := uploader.NewUserUploader(storer, 100*uploader.MB)
		suite.NoError(err)

		defaultFileName := "testdata/test.pdf"
		userUpload := BuildUserUpload(suite.DB(), []Customization{
			{
				Model: models.UserUpload{},
				ExtendedParams: &UserUploadExtendedParams{
					UserUploader: userUploader,
					AppContext:   suite.AppContextForTest(),
				},
			},
		}, nil)

		upload := userUpload.Upload

		// no need to test every bit of how the UserUploader works
		suite.Contains(upload.Filename, defaultFileName)
		suite.Equal(models.UploadTypeUSER, upload.UploadType)

		// Ensure the associated models are created
		suite.False(userUpload.DocumentID.IsNil())
		suite.False(userUpload.Document.ID.IsNil())
		suite.False(userUpload.UploaderID.IsNil())
		suite.False(userUpload.UploadID.IsNil())
		suite.False(userUpload.Upload.ID.IsNil())
	})

	suite.Run("Failed creation of upload - no appcontext", func() {
		// Under test:      BuildUserUpload
		// Mocked:          None
		// Set up:          Create a user upload with a user uploader
		//                  but no appcontext
		// Expected outcome:Should cause a panic
		storer := storageTest.NewFakeS3Storage(true)
		userUploader, err := uploader.NewUserUploader(storer, 100*uploader.MB)
		suite.NoError(err)

		suite.Panics(func() {
			BuildUserUpload(suite.DB(), []Customization{
				{
					Model: models.UserUpload{},
					ExtendedParams: &UserUploadExtendedParams{
						UserUploader: userUploader,
					},
				},
			}, nil)
		})

	})

	suite.Run("Successful creation of user upload with uploader and custom file", func() {
		// Under test:      BuildUserUploader
		// Mocked:          None
		// Set up:          Create a user upload with a specific file
		// Expected outcome:UserUpload should be created with default values
		storer := storageTest.NewFakeS3Storage(true)
		userUploader, err := uploader.NewUserUploader(storer, 100*uploader.MB)
		suite.NoError(err)

		uploadFile := "testdata/test.jpg"
		userUpload := BuildUserUpload(suite.DB(), []Customization{
			{
				Model: models.UserUpload{},
				ExtendedParams: &UserUploadExtendedParams{
					File:         FixtureOpen("test.jpg"),
					UserUploader: userUploader,
					AppContext:   suite.AppContextForTest(),
				},
			},
		}, nil)
		suite.False(userUpload.Upload.ID.IsNil())
		suite.Contains(userUpload.Upload.Filename, uploadFile)
		suite.Equal(models.UploadTypeUSER, userUpload.Upload.UploadType)
	})
}
