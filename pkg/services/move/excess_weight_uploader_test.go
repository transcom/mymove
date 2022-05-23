package move

import (
	"os"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/apperror"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services/upload"
	"github.com/transcom/mymove/pkg/storage/test"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *MoveServiceSuite) TestCreateExcessWeightUpload() {
	fakeFileStorer := test.NewFakeS3Storage(true)
	uploadCreator := upload.NewUploadCreator(fakeFileStorer)

	testFileName := "upload-test.pdf"
	defaultUploader := NewMoveExcessWeightUploader(uploadCreator)

	suite.Run("Success - Excess weight upload is created and move is updated", func() {
		move := testdatagen.MakeDefaultMove(suite.DB())

		testFile, fileErr := os.Open("../../testdatagen/testdata/test.pdf")
		suite.Require().NoError(fileErr)

		defer func() {
			closeErr := testFile.Close()
			suite.NoError(closeErr, "Error occurred while closing the test file.")
		}()

		updatedMove, err := defaultUploader.CreateExcessWeightUpload(
			suite.AppContextForTest(), move.ID, testFile, testFileName, models.UploadTypeUSER)
		suite.NoError(err)
		suite.Require().NotNil(updatedMove)

		suite.NotNil(updatedMove.ExcessWeightUploadID)
		suite.Require().NotNil(updatedMove.ExcessWeightUpload)
		suite.Equal(updatedMove.ExcessWeightUpload.ID, *updatedMove.ExcessWeightUploadID)

		suite.Equal(models.UploadTypeUSER, updatedMove.ExcessWeightUpload.UploadType)
		suite.Contains(updatedMove.ExcessWeightUpload.Filename, testFileName)
		suite.Contains(updatedMove.ExcessWeightUpload.Filename, move.ID.String())
		suite.Contains(updatedMove.ExcessWeightUpload.StorageKey, testFileName)
	})

	suite.Run("Fail - Move not found", func() {
		testFile, fileErr := os.Open("../../testdatagen/testdata/test.pdf")
		suite.Require().NoError(fileErr)

		defer func() {
			closeErr := testFile.Close()
			suite.NoError(closeErr, "Error occurred while closing the test file.")
		}()

		notFoundUUID := uuid.FromStringOrNil("00000000-0000-0000-0000-000000000001")

		updatedMove, err := defaultUploader.CreateExcessWeightUpload(
			suite.AppContextForTest(), notFoundUUID, testFile, testFileName, models.UploadTypeUSER)
		suite.Nil(updatedMove)
		suite.Require().Error(err)

		suite.IsType(apperror.NotFoundError{}, err)
		suite.Contains(err.Error(), notFoundUUID.String())
	})

	suite.Run("Fail - Move validation error rolls back transaction", func() {
		move := testdatagen.MakeDefaultMove(suite.DB())
		testFile, fileErr := os.Open("../../testdatagen/testdata/test.pdf")
		suite.Require().NoError(fileErr)

		defer func() {
			closeErr := testFile.Close()
			suite.NoError(closeErr, "Error occurred while closing the test file.")
		}()

		// A move cannot have a blank locator, so this will cause an error with the Validate function.
		// This validation happens during the DB update on the move, which happens AFTER the file has been
		// successfully uploaded using the UploadCreator.
		suite.Require().NoError(suite.DB().RawQuery("UPDATE moves SET locator='' WHERE id=$1;", move.ID).Exec())

		// Testing the number of uploads on DB prior to failure so we can make sure the DB rolls back the upload
		numUploadsBefore, countErr := suite.DB().Count(models.Upload{})
		suite.NoError(countErr)
		suite.Greater(numUploadsBefore, 0) // should have at least 1, likely 2 from the test data

		updatedMove, err := defaultUploader.CreateExcessWeightUpload(
			suite.AppContextForTest(), move.ID, testFile, testFileName, models.UploadTypeUSER)
		suite.Nil(updatedMove)
		suite.Require().Error(err)

		// Check the DB rollback
		numUploadsAfter, countErr := suite.DB().Count(models.Upload{})
		suite.NoError(countErr)
		suite.Equal(numUploadsBefore, numUploadsAfter)
	})
}

func (suite *MoveServiceSuite) TestCreateExcessWeightUploadPrime() {
	fakeFileStorer := test.NewFakeS3Storage(true)
	uploadCreator := upload.NewUploadCreator(fakeFileStorer)

	testFileName := "upload-test.pdf"
	primeUploader := NewPrimeMoveExcessWeightUploader(uploadCreator)

	suite.Run("Success - Excess weight upload is created for a Prime-available move", func() {
		testFile, fileErr := os.Open("../../testdatagen/testdata/test.pdf")
		suite.Require().NoError(fileErr)

		defer func() {
			closeErr := testFile.Close()
			suite.NoError(closeErr, "Error occurred while closing the test file.")
		}()

		primeMove := testdatagen.MakeAvailableMove(suite.DB())

		updatedMove, err := primeUploader.CreateExcessWeightUpload(
			suite.AppContextForTest(), primeMove.ID, testFile, testFileName, models.UploadTypePRIME)
		suite.NoError(err)
		suite.Require().NotNil(updatedMove)

		suite.NotNil(updatedMove.ExcessWeightUploadID)
		suite.Require().NotNil(updatedMove.ExcessWeightUpload)
		suite.Equal(updatedMove.ExcessWeightUpload.ID, *updatedMove.ExcessWeightUploadID)

		suite.Equal(models.UploadTypePRIME, updatedMove.ExcessWeightUpload.UploadType)
		suite.Contains(updatedMove.ExcessWeightUpload.Filename, testFileName)
		suite.Contains(updatedMove.ExcessWeightUpload.StorageKey, testFileName)
	})

	suite.Run("Fail - Cannot create upload for non-Prime move", func() {
		move := testdatagen.MakeDefaultMove(suite.DB())

		testFile, fileErr := os.Open("../../testdatagen/testdata/test.pdf")
		suite.Require().NoError(fileErr)

		defer func() {
			closeErr := testFile.Close()
			suite.NoError(closeErr, "Error occurred while closing the test file.")
		}()

		updatedMove, err := primeUploader.CreateExcessWeightUpload(
			suite.AppContextForTest(), move.ID, testFile, testFileName, models.UploadTypePRIME)
		suite.Nil(updatedMove)
		suite.Require().Error(err)

		suite.IsType(apperror.NotFoundError{}, err)
		suite.Contains(err.Error(), move.ID.String())
	})
}
