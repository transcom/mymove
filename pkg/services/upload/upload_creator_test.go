package upload

import (
	"os"
	"testing"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/storage/test"
)

// TestCreateUpload tests uploading a new document
func (suite *UploadServiceSuite) TestCreateUpload() {
	fakeFileStorer := test.NewFakeS3Storage(true)
	uploadCreator := NewUploadCreator(fakeFileStorer)

	testFileName := "upload-test.pdf"
	testFile, fileErr := os.Open("../../testdatagen/testdata/test.pdf")
	suite.Require().NoError(fileErr)

	suite.T().Run("Success - Upload is created", func(t *testing.T) {
		upload, err := uploadCreator.CreateUpload(suite.TestAppContext(), testFile, testFileName, models.UploadTypePRIME)
		suite.NoError(err)
		suite.Require().NotNil(upload)

		suite.Equal(models.UploadTypePRIME, upload.UploadType)
		suite.Contains(upload.Filename, testFileName)
		suite.Contains(upload.StorageKey, testFileName)
		suite.Equal(upload.Filename, upload.StorageKey)
	})

	suite.T().Run("Fail - Upload with invalid type causes an error", func(t *testing.T) {
		upload, err := uploadCreator.CreateUpload(suite.TestAppContext(), testFile, testFileName, "INVALID")
		suite.Nil(upload)
		suite.Require().Error(err)
	})

	err := testFile.Close()
	suite.NoError(err, "Error occurred while closing the test file.")
}
