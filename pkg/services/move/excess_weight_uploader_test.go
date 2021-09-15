package move

import (
	"os"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/upload"
	"github.com/transcom/mymove/pkg/storage/test"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *MoveServiceSuite) TestCreateExcessWeightUpload() {
	fakeFileStorer := test.NewFakeS3Storage(true)
	uploadCreator := upload.NewUploadCreator(fakeFileStorer)

	move := testdatagen.MakeDefaultMove(suite.DB())

	suite.Run("Default", func() {
		excessWeightUploader := NewMoveExcessWeightUploader(uploadCreator)

		testFileName := "upload-test.pdf"
		testFile, fileErr := os.Open("../../testdatagen/testdata/test.pdf")
		suite.Require().NoError(fileErr)

		suite.Run("Success - Excess weight upload is created and move is updated", func() {
			updatedMove, err := excessWeightUploader.CreateExcessWeightUpload(
				suite.TestAppContext(), move.ID, testFile, testFileName, models.UploadTypeUSER)
			suite.NoError(err)
			suite.Require().NotNil(updatedMove)

			suite.NotNil(updatedMove.ExcessWeightUploadID)
			suite.NotNil(updatedMove.ExcessWeightQualifiedAt)
			suite.False(updatedMove.ExcessWeightQualifiedAt.IsZero())
			suite.Require().NotNil(updatedMove.ExcessWeightUpload)

			suite.Equal(models.UploadTypeUSER, updatedMove.ExcessWeightUpload.UploadType)
			suite.Contains(updatedMove.ExcessWeightUpload.Filename, testFileName)
			suite.Contains(updatedMove.ExcessWeightUpload.StorageKey, testFileName)
		})

		suite.Run("Fail - Move not found", func() {
			notFoundUUID := uuid.FromStringOrNil("00000000-0000-0000-0000-000000000001")

			updatedMove, err := excessWeightUploader.CreateExcessWeightUpload(
				suite.TestAppContext(), notFoundUUID, testFile, testFileName, models.UploadTypeUSER)
			suite.Nil(updatedMove)
			suite.Require().Error(err)

			suite.IsType(services.NotFoundError{}, err)
			suite.Contains(err.Error(), notFoundUUID.String())
		})

		suite.Run("Fail - Invalid upload type causes error and rolls back transaction", func() {
			// Testing the number of uploads on DB prior to failure so we can make sure the DB rolls back the upload
			numUploadsBefore, countErr := suite.DB().Count(models.Upload{})
			suite.NoError(countErr)
			suite.Greater(numUploadsBefore, 0) // should have at least 1, likely 2 from the test data

			updatedMove, err := excessWeightUploader.CreateExcessWeightUpload(
				suite.TestAppContext(), move.ID, testFile, testFileName, "INVALID")
			suite.Nil(updatedMove)
			suite.Require().Error(err)

			// Check the DB rollback
			numUploadsAfter, countErr := suite.DB().Count(models.Upload{})
			suite.NoError(countErr)
			suite.Equal(numUploadsBefore, numUploadsAfter)
		})

		err := testFile.Close()
		suite.NoError(err, "Error occurred while closing the test file for basic uploader.")
	})

	suite.Run("Prime", func() {
		primeExcessWeightUploader := NewPrimeMoveExcessWeightUploader(uploadCreator)

		testFileName := "upload-test.pdf"
		testFile, fileErr := os.Open("../../testdatagen/testdata/test.pdf")
		suite.Require().NoError(fileErr)

		suite.Run("Success - Excess weight upload is created for a Prime-available move", func() {
			primeMove := testdatagen.MakeAvailableMove(suite.DB())

			updatedMove, err := primeExcessWeightUploader.CreateExcessWeightUpload(
				suite.TestAppContext(), primeMove.ID, testFile, testFileName, models.UploadTypePRIME)
			suite.NoError(err)
			suite.Require().NotNil(updatedMove)

			suite.NotNil(updatedMove.ExcessWeightUploadID)
			suite.NotNil(updatedMove.ExcessWeightQualifiedAt)
			suite.False(updatedMove.ExcessWeightQualifiedAt.IsZero())
			suite.Require().NotNil(updatedMove.ExcessWeightUpload)

			suite.Equal(models.UploadTypePRIME, updatedMove.ExcessWeightUpload.UploadType)
			suite.Contains(updatedMove.ExcessWeightUpload.Filename, testFileName)
			suite.Contains(updatedMove.ExcessWeightUpload.StorageKey, testFileName)
		})

		suite.Run("Fail - Cannot create upload for non-Prime move", func() {
			updatedMove, err := primeExcessWeightUploader.CreateExcessWeightUpload(
				suite.TestAppContext(), move.ID, testFile, testFileName, models.UploadTypePRIME)
			suite.Nil(updatedMove)
			suite.Require().Error(err)

			suite.IsType(services.NotFoundError{}, err)
			suite.Contains(err.Error(), move.ID.String())
		})

		err := testFile.Close()
		suite.NoError(err, "Error occurred while closing the test file for Prime uploader.")
	})
}
