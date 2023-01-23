package factory

import (
	"encoding/base64"
	"encoding/hex"

	"github.com/transcom/mymove/pkg/models"
	storageTest "github.com/transcom/mymove/pkg/storage/test"
	"github.com/transcom/mymove/pkg/uploader"
)

func (suite *FactorySuite) TestBuildUpload() {
	suite.Run("Successful creation of default upload", func() {
		// Under test:      BuildUpload
		// Set up:          Create a default upload
		// Expected outcome:Upload filename should be the default file
		// This file doesn't actually exist
		defaults := models.Upload{
			Filename:    "testFile.pdf",
			Bytes:       int64(2202009),
			ContentType: "application/pdf",
			Checksum:    "ImGQ2Ush0bDHsaQthV5BnQ==",
			UploadType:  models.UploadTypeUSER,
		}

		//Create upload
		upload := BuildUpload(suite.DB(), nil, nil)
		suite.Equal(defaults.Filename, upload.Filename)
		suite.Equal(defaults.Bytes, upload.Bytes)
		suite.Equal(defaults.ContentType, upload.ContentType)
		suite.Equal(defaults.Checksum, upload.Checksum)
		suite.Equal(defaults.UploadType, upload.UploadType)
	})
	suite.Run("Successful creation of customized upload", func() {
		// Under test:       BuildUpload
		// Set up:           Create a customized upload (no uploader)
		// Expected outcome: All fields should match
		// This file doesn't actually exist
		customUpload := models.Upload{
			Filename:    "BaisWinery.jpg",
			Bytes:       int64(6081979),
			ContentType: "application/jpg",
			Checksum:    "GauMarJosbDHsaQthV5BnQ==",
			UploadType:  models.UploadTypePRIME,
		}

		//Create upload
		upload := BuildUpload(suite.DB(), []Customization{
			{
				Model: customUpload,
			},
		}, nil)
		suite.Equal(customUpload.Filename, upload.Filename)
		suite.Equal(customUpload.Bytes, upload.Bytes)
		suite.Equal(customUpload.ContentType, upload.ContentType)
		suite.Equal(customUpload.Checksum, upload.Checksum)
		suite.Equal(customUpload.UploadType, upload.UploadType)
	})
	suite.Run("Successful creation of upload with uploader", func() {
		// Under test:      BuildUser
		// Mocked:          None
		// Set up:          Create an upload with an uploader and default file
		// Expected outcome:Upload filename should be the default file
		storer := storageTest.NewFakeS3Storage(true)
		uploader, err := uploader.NewUploader(storer, 100*uploader.MB, "USER")
		suite.NoError(err)

		defaultFileName := "testdata/test.pdf"
		upload := BuildUpload(suite.DB(), []Customization{
			{
				Model: models.Upload{},
				ExtendedParams: &UploadExtendedParams{
					Uploader:   uploader,
					AppContext: suite.AppContextForTest(),
				},
			},
		}, nil)

		suite.Contains(upload.Filename, defaultFileName)
		suite.Equal(encodeMD5("9ce13a1f0cf21385440d79faed42de78"), upload.Checksum)
		suite.Equal(int64(10596), upload.Bytes)
		suite.Equal("application/pdf", upload.ContentType)
		suite.Equal(models.UploadTypeUSER, upload.UploadType)
	})
	suite.Run("Failed creation of upload - no appcontext", func() {
		// Under test:      BuildUser
		// Mocked:          None
		// Set up:          Create an upload with an uploader but no appcontext
		// Expected outcome:Should cause a panic
		storer := storageTest.NewFakeS3Storage(true)
		uploader, err := uploader.NewUploader(storer, 100*uploader.MB, "USER")
		suite.NoError(err)

		suite.Panics(func() {
			BuildUpload(suite.DB(), []Customization{
				{
					Model: models.Upload{},
					ExtendedParams: &UploadExtendedParams{
						Uploader: uploader,
					},
				},
			}, nil)
		})

	})
	suite.Run("Successful creation of uploader with custom file", func() {
		// Under test:      BuildUser
		// Mocked:          None
		// Set up:          Create an upload with a specific file
		// Expected outcome:User should be created with default values
		storer := storageTest.NewFakeS3Storage(true)
		uploader, err := uploader.NewUploader(storer, 100*uploader.MB, "USER")
		suite.NoError(err)

		uploadFile := "testdata/test.jpg"
		upload := BuildUpload(suite.DB(), []Customization{
			{
				Model: models.Upload{
					Filename: "yoyoyo", // should not be used
				},
				ExtendedParams: &UploadExtendedParams{
					File:       FixtureOpen("test.jpg"),
					Uploader:   uploader,
					AppContext: suite.AppContextForTest(),
				},
			},
		}, nil)
		suite.Contains(upload.Filename, uploadFile)
		suite.Equal(encodeMD5("b1e74a6bc8e52bdf45075927168c4bb0"), upload.Checksum)
		suite.Equal(int64(37986), upload.Bytes)
		suite.Equal("image/jpeg", upload.ContentType)
		suite.Equal(models.UploadTypeUSER, upload.UploadType)
	})
}

func encodeMD5(md5 string) string {
	result, _ := hex.DecodeString(md5)
	return base64.StdEncoding.EncodeToString(result)
}
