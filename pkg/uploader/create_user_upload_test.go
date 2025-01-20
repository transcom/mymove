package uploader_test

import (
	"os"
	"time"

	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/storage/test"
	"github.com/transcom/mymove/pkg/uploader"
)

func (suite *UploaderSuite) TestCreateUserUpload() {

	suite.Run("Successful S3 upload clean AV", func() {
		document := factory.BuildDocument(suite.DB(), nil, nil)
		user := factory.BuildDefaultUser(suite.DB())
		fakeS3 := test.NewFakeS3Storage(true, map[string]string{"av-status": "CLEAN"})

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
	})

	suite.Run("Infected S3 upload", func() {
		document := factory.BuildDocument(suite.DB(), nil, nil)
		user := factory.BuildDefaultUser(suite.DB())
		fakeS3 := test.NewFakeS3Storage(true, map[string]string{"av-status": "INFECTED"})

		file, err := os.Open("../testdatagen/testdata/test.pdf")
		suite.FatalNoError(err)

		_, _, _, err = uploader.CreateUserUploadForDocumentWrapper(
			suite.AppContextForTest(),
			user.ID,
			fakeS3,
			file,
			file.Name(),
			25*uploader.MB,
			uploader.AllowedTypesPDF,
			&document.ID,
			models.UploadTypeUSER)

		suite.Error(err, "S3 upload should have returned an infected error")
	})

	suite.Run("S3 upload timeout", func() {
		document := factory.BuildDocument(suite.DB(), nil, nil)
		user := factory.BuildDefaultUser(suite.DB())
		fakeS3 := test.NewFakeS3Storage(true, map[string]string{}) // No av-status tag

		file, err := os.Open("../testdatagen/testdata/test.pdf")
		suite.FatalNoError(err)

		// Channel to capture the result of your upload call
		// As this is a timeout test,
		// we use a goroutine to not hold up the processor thread from running other tests.
		// This will place us in the processing queue (go routine) to come back later when ready
		done := make(chan error, 1)

		go func() {
			// This function call WILL timeout
			_, _, verrs, err := uploader.CreateUserUploadForDocumentWrapper(
				suite.AppContextForTest(),
				user.ID,
				fakeS3,
				file,
				file.Name(),
				25*uploader.MB,
				uploader.AllowedTypesPDF,
				&document.ID,
				models.UploadTypeUSER,
			)

			// Return the errs to the channel if any
			if err != nil || verrs.HasAny() {
				done <- err
				return
			}

			done <- nil
		}()

		// Receive channel result
		select {
		case result := <-done:
			// Case err
			suite.NoError(result, "Failed to create upload successfully or encountered validations")
		case <-time.After(10 * time.Second):
			// Case routine still running after 10 seconds
			// This only occurs if it's timing out, otherwise the test runs much faster than
			// 10 seconds
			suite.T().Log("S3 upload exceeded 10 seconds without returning any errors, passing test")
		}
	})

}
