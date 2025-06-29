// RA Summary: gosec - errcheck - Unchecked return value
// RA: Linter flags errcheck error: Ignoring a method's return value can cause the program to overlook unexpected states and conditions.
// RA: Functions with unchecked return values in the file are used to clean up file created for unit test
// RA: Given the functions causing the lint errors are used to clean up local storage space after a unit test, it does not present a risk
// RA Developer Status: Mitigated
// RA Validator Status: Mitigated
// RA Modified Severity: N/A
// nolint:errcheck
package uploader_test

import (
	"fmt"
	"io"

	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/storage/test"
	"github.com/transcom/mymove/pkg/uploader"
)

func (suite *UploaderSuite) TestUserUploadFromLocalFile() {
	document := factory.BuildDocument(suite.DB(), nil, nil)

	userUploader, err := uploader.NewUserUploader(suite.storer, 25*uploader.MB)
	suite.NoError(err)
	file := suite.fixture("test.pdf")

	userUpload, verrs, err := userUploader.CreateUserUploadForDocument(suite.AppContextForTest(), &document.ID, document.ServiceMember.UserID, uploader.File{File: file}, uploader.AllowedTypesPDF)
	suite.Nil(err, "failed to create upload")
	suite.False(verrs.HasAny(), "failed to validate upload", verrs)
	suite.Equal(userUpload.Upload.ContentType, uploader.FileTypePDF)
	suite.Equal(userUpload.Upload.Checksum, "w7rJQqzlaazDW+mxTU9Q40Qchr3DW7FPQD7f8Js2J88=")
}

func (suite *UploaderSuite) TestUserUploadFromLocalFileZeroLength() {
	document := factory.BuildDocument(suite.DB(), nil, nil)

	userUploader, err := uploader.NewUserUploader(suite.storer, 25*uploader.MB)
	suite.NoError(err)
	file, cleanup, err := suite.createFileOfArbitrarySize(uint64(0 * uploader.MB))
	suite.Nil(err, "failed to create upload")
	defer cleanup()

	userUpload, verrs, err := userUploader.CreateUserUploadForDocument(suite.AppContextForTest(), &document.ID, document.ServiceMember.UserID, uploader.File{File: file}, uploader.AllowedTypesAny)
	suite.Equal(uploader.ErrZeroLengthFile, err)
	suite.False(verrs.HasAny(), "failed to validate upload")
	suite.Nil(userUpload, "returned an upload when erroring")
}

func (suite *UploaderSuite) TestUserUploadFromLocalFileWrongContentType() {
	document := factory.BuildDocument(suite.DB(), nil, nil)

	userUploader, err := uploader.NewUserUploader(suite.storer, 25*uploader.MB)
	suite.NoError(err)
	file, cleanup, err := suite.createFileOfArbitrarySize(uint64(1 * uploader.MB))
	suite.Nil(err, "failed to create upload")
	defer cleanup()

	upload, verrs, err := userUploader.CreateUserUploadForDocument(suite.AppContextForTest(), &document.ID, document.ServiceMember.UserID, uploader.File{File: file}, uploader.AllowedTypesPDF)
	suite.Error(err)
	suite.Equal(fmt.Sprintf("content type \"application/octet-stream\" is not one of the supported types [%s]", uploader.FileTypePDF), err.Error())
	suite.True(verrs.HasAny(), "invalid content type for upload")
	suite.Nil(upload, "returned an upload when erroring")
}

func (suite *UploaderSuite) TestTooLargeUserUploadFromLocalFile() {
	document := factory.BuildDocument(suite.DB(), nil, nil)

	userUploader, err := uploader.NewUserUploader(suite.storer, 25*uploader.MB)
	suite.NoError(err)
	f, cleanup, err := suite.createFileOfArbitrarySize(uint64(26 * uploader.MB))
	suite.NoError(err)
	defer cleanup()

	_, verrs, err := userUploader.CreateUserUploadForDocument(suite.AppContextForTest(), &document.ID, document.ServiceMember.UserID, uploader.File{File: f}, uploader.AllowedTypesAny)
	suite.Error(err)
	suite.IsType(uploader.ErrTooLarge{}, err)
	suite.False(verrs.HasAny(), "failed to validate upload")
}

func (suite *UploaderSuite) TestUserUploadStorerCalledWithTags() {
	document := factory.BuildDocument(suite.DB(), nil, nil)
	fakeS3 := test.NewFakeS3Storage(true)

	userUploader, err := uploader.NewUserUploader(fakeS3, 25*uploader.MB)
	suite.NoError(err)
	f, cleanup, err := suite.createFileOfArbitrarySize(uint64(5 * uploader.MB))
	suite.NoError(err)
	defer cleanup()

	tags := "metaDataTag=value"

	// assert tags are passed along to storer
	_, verrs, err := userUploader.CreateUserUploadForDocument(suite.AppContextForTest(), &document.ID, document.ServiceMember.UserID, uploader.File{File: f, Tags: &tags}, uploader.AllowedTypesAny)

	suite.NoError(err)
	suite.False(verrs.HasAny(), "failed to validate upload")
}

func (suite *UploaderSuite) TestCreateUserUploadNoDocument() {
	document := factory.BuildDocument(suite.DB(), nil, nil)
	userID := document.ServiceMember.UserID

	userUploader, err := uploader.NewUserUploader(suite.storer, 25*uploader.MB)
	suite.NoError(err)
	file := suite.fixture("test.pdf")
	fixtureFileInfo, err := file.Stat()
	suite.NoError(err)

	// Create file and upload
	userUpload, verrs, err := userUploader.CreateUserUpload(suite.AppContextForTest(), userID, uploader.File{File: file}, uploader.AllowedTypesPDF)
	suite.Nil(err, "failed to create upload")
	suite.Empty(verrs.Error(), "verrs returned error")
	suite.NotNil(userUpload, "failed to create upload structure")
	file.Close()

	// Download file and test size
	download, err := userUploader.Download(suite.AppContextForTest(), userUpload)
	suite.NoError(err)
	defer download.Close()

	outputFile, err := suite.helperNewTempFile()
	suite.NoError(err)
	defer outputFile.Close()

	written, err := io.Copy(outputFile, download)
	suite.NoError(err)
	suite.NotEqual(0, written)

	info, err := outputFile.Stat()
	suite.Equal(fixtureFileInfo.Size(), info.Size())
	suite.NoError(err)

	// Delete file previously uploaded
	err = userUploader.DeleteUserUpload(suite.AppContextForTest(), userUpload)
	suite.NoError(err)
}

func (suite *UploaderSuite) TestUpdateUserUploadFilename_Prefix() {
	document := factory.BuildDocument(suite.DB(), nil, nil)

	userUploader, err := uploader.NewUserUploader(suite.storer, 25*uploader.MB)
	suite.NoError(err)
	file := suite.fixture("weightEstimatorExpectSuccessfulUpload.xlsx")

	filename := "weightEstimatorExpectSuccessfulUpload.xlsx"
	userUpload, verrs, err := userUploader.CreateUserUploadForDocument(suite.AppContextForTest(), &document.ID, document.ServiceMember.UserID, uploader.File{File: file}, uploader.AllowedTypesPPMDocuments)
	suite.Nil(err, "failed to create upload")
	suite.False(verrs.HasAny(), "failed to validate upload")

	updated, verrs, err := userUploader.UpdateUserXlsxUploadFilename(suite.AppContextForTest(), userUpload, filename)
	suite.NoError(err, "failed to update filename")
	suite.False(verrs.HasAny(), "validation errors on update")

	suite.Equal(
		filename,
		updated.Upload.Filename,
		"expected updated Upload.Filename to equal %q, got %q",
		filename,
		updated.Upload.Filename,
	)
}
