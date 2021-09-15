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
	excessWeightUploader := NewMoveExcessWeightUploader(uploadCreator)

	move := testdatagen.MakeDefaultMove(suite.DB())

	testFileName := "upload-test.pdf"
	testFile, fileErr := os.Open("../../testdatagen/testdata/test.pdf")
	suite.Require().NoError(fileErr)

	suite.Run("Success - Excess weight upload is created and move is updated", func() {
		updatedMove, err := excessWeightUploader.CreateExcessWeightUpload(suite.TestAppContext(), move.ID, testFile, testFileName, models.UploadTypePRIME)
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

	suite.Run("Fail - Move not found", func() {
		notFoundUUID := uuid.FromStringOrNil("00000000-0000-0000-0000-000000000001")

		updatedMove, err := excessWeightUploader.CreateExcessWeightUpload(suite.TestAppContext(), notFoundUUID, testFile, testFileName, models.UploadTypePRIME)
		suite.Nil(updatedMove)
		suite.Require().Error(err)

		suite.IsType(services.NotFoundError{}, err)
		suite.Contains(err.Error(), notFoundUUID.String())
	})

	err := testFile.Close()
	suite.NoError(err, "Error occurred while closing the test file.")
}
