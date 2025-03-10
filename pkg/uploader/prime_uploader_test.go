package uploader_test

import (
	"fmt"

	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/storage/test"
	"github.com/transcom/mymove/pkg/uploader"
)

func (suite *UploaderSuite) TestPrimeUploadFromLocalFile() {
	document := factory.BuildProofOfServiceDoc(suite.DB(), nil, nil)

	primeUploader, err := uploader.NewPrimeUploader(suite.storer, 25*uploader.MB)
	suite.NoError(err)
	file := suite.fixture("test.pdf")

	contractor := factory.FetchOrBuildDefaultContractor(suite.DB(), nil, nil)

	primeUpload, verrs, err := primeUploader.CreatePrimeUploadForDocument(suite.AppContextForTest(), &document.ID, contractor.ID, uploader.File{File: file}, uploader.AllowedTypesPDF)
	suite.Nil(err, "failed to create upload")
	suite.False(verrs.HasAny(), "failed to validate upload", verrs)
	suite.Equal(primeUpload.Upload.ContentType, uploader.FileTypePDF)
	suite.Equal(primeUpload.Upload.Checksum, "w7rJQqzlaazDW+mxTU9Q40Qchr3DW7FPQD7f8Js2J88=")
}

func (suite *UploaderSuite) TestPrimeUploadFromLocalFileZeroLength() {
	document := factory.BuildProofOfServiceDoc(suite.DB(), nil, nil)

	primeUploader, err := uploader.NewPrimeUploader(suite.storer, 25*uploader.MB)
	suite.NoError(err)
	file, cleanup, err := suite.createFileOfArbitrarySize(uint64(0 * uploader.MB))
	suite.Nil(err, "failed to create upload")
	defer cleanup()

	contractor := factory.FetchOrBuildDefaultContractor(suite.DB(), nil, nil)

	primeUpload, verrs, err := primeUploader.CreatePrimeUploadForDocument(suite.AppContextForTest(), &document.ID, contractor.ID, uploader.File{File: file}, uploader.AllowedTypesAny)
	suite.Equal(uploader.ErrZeroLengthFile, err)
	suite.False(verrs.HasAny(), "failed to validate upload")
	suite.Nil(primeUpload, "returned an upload when erroring")
}

func (suite *UploaderSuite) TestPrimeUploadFromLocalFileWrongContentType() {
	document := factory.BuildProofOfServiceDoc(suite.DB(), nil, nil)

	primeUploader, err := uploader.NewPrimeUploader(suite.storer, 25*uploader.MB)
	suite.NoError(err)
	file, cleanup, err := suite.createFileOfArbitrarySize(uint64(1 * uploader.MB))
	suite.Nil(err, "failed to create upload")
	defer cleanup()

	contractor := factory.FetchOrBuildDefaultContractor(suite.DB(), nil, nil)

	upload, verrs, err := primeUploader.CreatePrimeUploadForDocument(suite.AppContextForTest(), &document.ID, contractor.ID, uploader.File{File: file}, uploader.AllowedTypesPDF)
	suite.Error(err)
	suite.Equal(fmt.Sprintf("content type \"application/octet-stream\" is not one of the supported types [%s]", uploader.FileTypePDF), err.Error())
	suite.True(verrs.HasAny(), "invalid content type for upload")
	suite.Nil(upload, "returned an upload when erroring")
}

func (suite *UploaderSuite) TestTooLargePrimeUploadFromLocalFile() {
	document := factory.BuildProofOfServiceDoc(suite.DB(), nil, nil)

	primeUploader, err := uploader.NewPrimeUploader(suite.storer, 25*uploader.MB)
	suite.NoError(err)
	f, cleanup, err := suite.createFileOfArbitrarySize(uint64(26 * uploader.MB))
	suite.NoError(err)
	defer cleanup()

	contractor := factory.FetchOrBuildDefaultContractor(suite.DB(), nil, nil)

	_, verrs, err := primeUploader.CreatePrimeUploadForDocument(suite.AppContextForTest(), &document.ID, contractor.ID, uploader.File{File: f}, uploader.AllowedTypesAny)
	suite.Error(err)
	suite.IsType(uploader.ErrTooLarge{}, err)
	suite.False(verrs.HasAny(), "failed to validate upload")
}

func (suite *UploaderSuite) TestPrimeUploadStorerCalledWithTags() {
	document := factory.BuildProofOfServiceDoc(suite.DB(), nil, nil)
	fakeS3 := test.NewFakeS3Storage(true)

	primeUploader, err := uploader.NewPrimeUploader(fakeS3, 25*uploader.MB)
	suite.NoError(err)
	f, cleanup, err := suite.createFileOfArbitrarySize(uint64(5 * uploader.MB))
	suite.NoError(err)
	defer cleanup()

	tags := "metaDataTag=value"

	contractor := factory.FetchOrBuildDefaultContractor(suite.DB(), nil, nil)

	// assert tags are passed along to storer
	_, verrs, err := primeUploader.CreatePrimeUploadForDocument(suite.AppContextForTest(), &document.ID, contractor.ID, uploader.File{File: f, Tags: &tags}, uploader.AllowedTypesAny)

	suite.NoError(err)
	suite.False(verrs.HasAny(), "failed to validate upload")
}
