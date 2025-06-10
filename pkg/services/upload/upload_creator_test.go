package upload

import (
	"os"
	"regexp"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/storage/test"
	"github.com/transcom/mymove/pkg/utils"
)

// TestCreateUpload tests uploading a new document
func (suite *UploadServiceSuite) TestCreateUpload() {
	fakeFileStorer := test.NewFakeS3Storage(true)
	uploadCreator := NewUploadCreator(fakeFileStorer)

	testFileName := "upload-test.pdf"
	testFileNameNoExtension := "upload-test"
	testFile, fileErr := os.Open("../../testdatagen/testdata/test.pdf")
	suite.Require().NoError(fileErr)

	suite.Run("Success - Upload is created", func() {
		upload, err := uploadCreator.CreateUpload(suite.AppContextForTest(), testFile, testFileName, models.UploadTypePRIME)
		suite.NoError(err)
		suite.Require().NotNil(upload)

		suite.Equal(models.UploadTypePRIME, upload.UploadType)
		suite.Contains(upload.Filename, testFileNameNoExtension)
		suite.Contains(upload.StorageKey, testFileNameNoExtension)
		suite.Equal(upload.Filename, upload.StorageKey)
	})

	suite.Run("Fail - Upload with invalid type causes an error", func() {
		upload, err := uploadCreator.CreateUpload(suite.AppContextForTest(), testFile, testFileName, "INVALID")
		suite.Nil(upload)
		suite.Require().Error(err)
	})

	err := testFile.Close()
	suite.NoError(err, "Error occurred while closing the test file.")
}

// Test_assembleUploadFilePathName tests assembling the file path for saving in storage
func (suite *UploadServiceSuite) Test_assembleUploadFilePathName() {
	filePathName := "move/4b7b7c9b-8023-4843-9c4f-a8185bfb7b11/proof.pdf"
	resultPattern, err := regexp.Compile(
		`move/4b7b7c9b-8023-4843-9c4f-a8185bfb7b11/(proof-\d{14})\.pdf`)
	suite.Require().NoError(err, "Error compiling regex for test")

	result := utils.AppendTimestampToFilename(filePathName)
	suite.True(resultPattern.MatchString(result), "Regex should match filename: %s", result)
}
