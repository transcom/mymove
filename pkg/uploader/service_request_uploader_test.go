package uploader_test

import (
	"fmt"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/storage/test"
	"github.com/transcom/mymove/pkg/uploader"
)

func (suite *UploaderSuite) TestServiceRequestUploadFromLocalFile() {
	document := factory.BuildServiceRequestDocument(suite.DB(), nil, nil)

	serviceRequestUploader, err := uploader.NewServiceRequestUploader(suite.storer, 25*uploader.MB)
	suite.NoError(err)
	file := suite.fixture("test.pdf")

	contractor := factory.FetchOrBuildDefaultContractor(suite.DB(), nil, nil)

	serviceRequestUpload, verrs, err := serviceRequestUploader.CreateServiceRequestUploadForDocument(suite.AppContextForTest(), &document.ID, contractor.ID, uploader.File{File: file}, uploader.AllowedTypesPDF)
	suite.Nil(err, "failed to create upload")
	suite.False(verrs.HasAny(), "failed to validate upload", verrs)
	suite.Equal(serviceRequestUpload.Upload.ContentType, uploader.FileTypePDF)
	suite.Equal(serviceRequestUpload.Upload.Checksum, "w7rJQqzlaazDW+mxTU9Q40Qchr3DW7FPQD7f8Js2J88=")
}

func (suite *UploaderSuite) TestServiceRequestUploadFromLocalFileZeroLength() {
	document := factory.BuildServiceRequestDocument(suite.DB(), nil, nil)

	serviceRequestUploader, err := uploader.NewServiceRequestUploader(suite.storer, 25*uploader.MB)
	suite.NoError(err)
	file, cleanup, err := suite.createFileOfArbitrarySize(uint64(0 * uploader.MB))
	suite.Nil(err, "failed to create upload")
	defer cleanup()

	contractor := factory.FetchOrBuildDefaultContractor(suite.DB(), nil, nil)

	serviceRequestUpload, verrs, err := serviceRequestUploader.CreateServiceRequestUploadForDocument(suite.AppContextForTest(), &document.ID, contractor.ID, uploader.File{File: file}, uploader.AllowedTypesAny)
	suite.Equal(uploader.ErrZeroLengthFile, err)
	suite.False(verrs.HasAny(), "failed to validate upload")
	suite.Nil(serviceRequestUpload, "returned an upload when erroring")
}

func (suite *UploaderSuite) TestServiceRequestUploadFromLocalFileWrongContentType() {
	document := factory.BuildServiceRequestDocument(suite.DB(), nil, nil)

	serviceRequestUploader, err := uploader.NewServiceRequestUploader(suite.storer, 25*uploader.MB)
	suite.NoError(err)
	file, cleanup, err := suite.createFileOfArbitrarySize(uint64(1 * uploader.MB))
	suite.Nil(err, "failed to create upload")
	defer cleanup()

	contractor := factory.FetchOrBuildDefaultContractor(suite.DB(), nil, nil)

	upload, verrs, err := serviceRequestUploader.CreateServiceRequestUploadForDocument(suite.AppContextForTest(), &document.ID, contractor.ID, uploader.File{File: file}, uploader.AllowedTypesPDF)
	suite.Error(err)
	suite.Equal(fmt.Sprintf("content type \"application/octet-stream\" is not one of the supported types [%s]", uploader.FileTypePDF), err.Error())
	suite.True(verrs.HasAny(), "invalid content type for upload")
	suite.Nil(upload, "returned an upload when erroring")
}

func (suite *UploaderSuite) TestTooLargeServiceRequestUploadFromLocalFile() {
	document := factory.BuildProofOfServiceDoc(suite.DB(), nil, nil)

	serviceRequestUploader, err := uploader.NewServiceRequestUploader(suite.storer, 25*uploader.MB)
	suite.NoError(err)
	f, cleanup, err := suite.createFileOfArbitrarySize(uint64(26 * uploader.MB))
	suite.NoError(err)
	defer cleanup()

	contractor := factory.FetchOrBuildDefaultContractor(suite.DB(), nil, nil)

	_, verrs, err := serviceRequestUploader.CreateServiceRequestUploadForDocument(suite.AppContextForTest(), &document.ID, contractor.ID, uploader.File{File: f}, uploader.AllowedTypesAny)
	suite.Error(err)
	suite.IsType(uploader.ErrTooLarge{}, err)
	suite.False(verrs.HasAny(), "failed to validate upload")
}

func (suite *UploaderSuite) TestFailureCreatingServiceRequestUpload() {
	document := factory.BuildServiceRequestDocument(suite.DB(), nil, nil)

	serviceRequestUploader, err := uploader.NewServiceRequestUploader(suite.storer, 25*uploader.MB)
	suite.NoError(err)
	file := suite.fixture("test.pdf")

	contractor := factory.FetchOrBuildDefaultContractor(suite.DB(), nil, nil)

	serviceRequestUpload, verrs, err := serviceRequestUploader.CreateServiceRequestUploadForDocument(suite.AppContextForTest(), &document.ID, contractor.ID, uploader.File{File: file}, uploader.AllowedTypesPDF)
	suite.Nil(err, "failed to create upload")
	suite.False(verrs.HasAny(), "failed to validate upload", verrs)
	suite.Equal(serviceRequestUpload.Upload.ContentType, uploader.FileTypePDF)
	suite.Equal(serviceRequestUpload.Upload.Checksum, "w7rJQqzlaazDW+mxTU9Q40Qchr3DW7FPQD7f8Js2J88=")
	suite.False(verrs.HasAny(), "failed to validate upload")
}

func (suite *UploaderSuite) TestServiceRequestUploadStorerCalledWithTags() {
	document := factory.BuildServiceRequestDocument(suite.DB(), nil, nil)
	fakeS3 := test.NewFakeS3Storage(true)

	serviceRequestUploader, err := uploader.NewServiceRequestUploader(fakeS3, 25*uploader.MB)
	suite.NoError(err)
	f, cleanup, err := suite.createFileOfArbitrarySize(uint64(5 * uploader.MB))
	suite.NoError(err)
	defer cleanup()

	tags := "metaDataTag=value"

	// assert tags are passed along to storer
	_, verrs, err := serviceRequestUploader.CreateServiceRequestUploadForDocument(suite.AppContextForTest(), &document.ID, uuid.Nil, uploader.File{File: f, Tags: &tags}, uploader.AllowedTypesAny)

	suite.NoError(err)
	suite.True(verrs.HasAny(), "error creating new prime upload")
}
