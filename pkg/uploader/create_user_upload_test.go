package uploader_test

import (
	"os"

	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/storage/test"
	"github.com/transcom/mymove/pkg/uploader"
)

// TODO: Simulating S3 is giving me a struggle, this test doesn't
// trigger the av scan poller yet
func (suite *UploaderSuite) TestCreateUserUpload() {
	document := factory.BuildDocument(suite.DB(), nil, nil)
	user := factory.BuildDefaultUser(suite.DB())
	fakeS3 := test.NewFakeS3Storage(true)

	file, err := os.Open("../testdatagen/testdata/test.pdf")
	suite.FatalNoError(err)

	_, _, verrs, err := uploader.CreateUserUploadForDocumentWrapper(
		suite.AppContextForTest(),
		user.ID,
		fakeS3,
		file,
		file.Name(),
		25*uploader.MB,
		uploader.AllowedTypesPDF,
		&document.ID,
		models.UploadTypeUSER)

	suite.Nil(err, "failed to create upload")
	suite.False(verrs.HasAny(), "failed to validate upload", verrs)
}
