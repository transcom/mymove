package move

import (
	"os"

	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/services/upload"
	"github.com/transcom/mymove/pkg/storage/test"
)

func (suite *MoveServiceSuite) TestAdditionalDocumentUploader() {
	fakeFileStorer := test.NewFakeS3Storage(true)
	uploadCreator := upload.NewUploadCreator(fakeFileStorer)

	testFileName := "upload-test.pdf"
	additionalDocumentUploader := NewMoveAdditionalDocumentsUploader(uploadCreator)

	suite.Run("Success - Additional Document upload is created", func() {
		testFile, fileErr := os.Open("../../testdatagen/testdata/test.pdf")
		suite.Require().NoError(fileErr)

		defer func() {
			closeErr := testFile.Close()
			suite.NoError(closeErr, "Error occurred while closing the test file.")
		}()

		serviceMember := factory.BuildServiceMember(suite.DB(), nil, nil)
		move := factory.BuildMove(suite.DB(), nil, nil)

		session := auth.Session{
			ApplicationName: auth.MilApp,
			ServiceMemberID: serviceMember.ID,
			IDToken:         "fake_token",
			AccessToken:     "fakeAccessToken",
		}

		updatedMove, _, _, err := additionalDocumentUploader.CreateAdditionalDocumentsUpload(
			suite.AppContextWithSessionForTest(&session),
			serviceMember.ID,
			move.ID,
			testFile,
			testFileName,
			fakeFileStorer)
		suite.NoError(err)
		suite.NotNil(updatedMove.UpdatedAt.Date())
		suite.NotNil(move.AdditionalDocuments.UserUploads)
	})

	suite.Run("Fail - Cannot create upload for move", func() {
		move := factory.BuildMove(suite.DB(), nil, nil)
		session := auth.Session{
			ApplicationName: auth.MilApp,
			UserID:          move.Orders.ServiceMemberID,
			IDToken:         "fake_token",
			AccessToken:     "fakeAccessToken",
		}

		testFile, fileErr := os.Open("../../testdatagen/testdata/test.pdf")
		suite.Require().NoError(fileErr)

		defer func() {
			closeErr := testFile.Close()
			suite.NoError(closeErr, "Error occurred while closing the test file.")
		}()

		updatedMove, _, _, err := additionalDocumentUploader.CreateAdditionalDocumentsUpload(
			suite.AppContextWithSessionForTest(&session),
			move.Orders.ServiceMemberID,
			move.ID,
			testFile,
			testFileName,
			fakeFileStorer)
		suite.Nil(updatedMove)
		suite.Require().Error(err)

		suite.IsType(apperror.NotFoundError{}, err)
		suite.Contains(err.Error(), move.ID.String())
	})
}
