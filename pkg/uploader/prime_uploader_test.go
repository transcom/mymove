package uploader_test

import (
	"github.com/transcom/mymove/pkg/storage/test"
	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/uploader"
)

func (suite *UploaderSuite) TestPrimeUploadFromLocalFile() {
	document := testdatagen.MakeDefaultProofOfServiceDoc(suite.DB())

	primeUploader, err := uploader.NewPrimeUploader(suite.storer, 25*uploader.MB)
	suite.NoError(err)
	file := suite.fixture("test.pdf")

	contractor := testdatagen.MakeDefaultContractor(suite.DB())

	primeUpload, verrs, err := primeUploader.CreatePrimeUploadForDocument(suite.AppContextForTest(), &document.ID, contractor.ID, uploader.File{File: file}, uploader.AllowedTypesPDF)
	suite.Nil(err, "failed to create upload")
	suite.False(verrs.HasAny(), "failed to validate upload", verrs)
	suite.Equal(primeUpload.Upload.ContentType, "application/pdf")
	suite.Equal(primeUpload.Upload.Checksum, "nOE6HwzyE4VEDXn67ULeeA==")
}

func (suite *UploaderSuite) TestPrimeUploadFromLocalFileZeroLength() {
	document := testdatagen.MakeDefaultProofOfServiceDoc(suite.DB())

	primeUploader, err := uploader.NewPrimeUploader(suite.storer, 25*uploader.MB)
	suite.NoError(err)
	file, cleanup, err := suite.createFileOfArbitrarySize(uint64(0 * uploader.MB))
	suite.Nil(err, "failed to create upload")
	defer cleanup()

	contractor := testdatagen.MakeDefaultContractor(suite.DB())

	primeUpload, verrs, err := primeUploader.CreatePrimeUploadForDocument(suite.AppContextForTest(), &document.ID, contractor.ID, uploader.File{File: file}, uploader.AllowedTypesAny)
	suite.Equal(uploader.ErrZeroLengthFile, err)
	suite.False(verrs.HasAny(), "failed to validate upload")
	suite.Nil(primeUpload, "returned an upload when erroring")
}

func (suite *UploaderSuite) TestPrimeUploadFromLocalFileWrongContentType() {
	document := testdatagen.MakeDefaultProofOfServiceDoc(suite.DB())

	primeUploader, err := uploader.NewPrimeUploader(suite.storer, 25*uploader.MB)
	suite.NoError(err)
	file, cleanup, err := suite.createFileOfArbitrarySize(uint64(1 * uploader.MB))
	suite.Nil(err, "failed to create upload")
	defer cleanup()

	contractor := testdatagen.MakeDefaultContractor(suite.DB())

	upload, verrs, err := primeUploader.CreatePrimeUploadForDocument(suite.AppContextForTest(), &document.ID, contractor.ID, uploader.File{File: file}, uploader.AllowedTypesPDF)
	suite.NoError(err)
	suite.True(verrs.HasAny(), "invalid content type for upload")
	suite.Nil(upload, "returned an upload when erroring")
}

func (suite *UploaderSuite) TestTooLargePrimeUploadFromLocalFile() {
	document := testdatagen.MakeDefaultProofOfServiceDoc(suite.DB())

	primeUploader, err := uploader.NewPrimeUploader(suite.storer, 25*uploader.MB)
	suite.NoError(err)
	f, cleanup, err := suite.createFileOfArbitrarySize(uint64(26 * uploader.MB))
	suite.NoError(err)
	defer cleanup()

	contractor := testdatagen.MakeDefaultContractor(suite.DB())

	_, verrs, err := primeUploader.CreatePrimeUploadForDocument(suite.AppContextForTest(), &document.ID, contractor.ID, uploader.File{File: f}, uploader.AllowedTypesAny)
	suite.Error(err)
	suite.IsType(uploader.ErrTooLarge{}, err)
	suite.False(verrs.HasAny(), "failed to validate upload")
}

func (suite *UploaderSuite) TestPrimeUploadStorerCalledWithTags() {
	document := testdatagen.MakeDefaultProofOfServiceDoc(suite.DB())
	fakeS3 := test.NewFakeS3Storage(true)

	primeUploader, err := uploader.NewPrimeUploader(fakeS3, 25*uploader.MB)
	suite.NoError(err)
	f, cleanup, err := suite.createFileOfArbitrarySize(uint64(5 * uploader.MB))
	suite.NoError(err)
	defer cleanup()

	tags := "metaDataTag=value"

	contractor := testdatagen.MakeDefaultContractor(suite.DB())

	// assert tags are passed along to storer
	_, verrs, err := primeUploader.CreatePrimeUploadForDocument(suite.AppContextForTest(), &document.ID, contractor.ID, uploader.File{File: f, Tags: &tags}, uploader.AllowedTypesAny)

	suite.NoError(err)
	suite.False(verrs.HasAny(), "failed to validate upload")
}
