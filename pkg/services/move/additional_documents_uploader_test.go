package move

import (
	"github.com/go-openapi/runtime"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services/upload"
	storageTest "github.com/transcom/mymove/pkg/storage/test"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *MoveServiceSuite) TestAdditionalDocumentUploader() {
	fakeFileStorer := storageTest.NewFakeS3Storage(true)
	uploadCreator := upload.NewUploadCreator(fakeFileStorer)

	additionalDocumentUploader := NewMoveAdditionalDocumentsUploader(uploadCreator)

	setUpOrders := func(setUpPreExistingAdditionalDocuments bool) *models.Order {
		var moves models.Moves
		var customs []factory.Customization

		if setUpPreExistingAdditionalDocuments {
			customs = []factory.Customization{
				{
					Model: models.Document{},
					Type:  &factory.AdditionalDocuments,
				},
			}
		}

		mto := factory.BuildServiceCounselingCompletedMove(suite.DB(), customs, nil)

		order := mto.Orders
		order.Moves = append(moves, mto)

		return &order
	}

	setUpFileToUpload := func() (*runtime.File, func()) {
		file := testdatagen.FixtureRuntimeFile("filled-out-orders.pdf")

		cleanUpFunc := func() {
			fileCloseErr := file.Close()
			suite.NoError(fileCloseErr)
		}

		return file, cleanUpFunc
	}

	suite.Run("Creates and saves new additionalDocument doc when the move.AdditionalDocuments is nil", func() {
		order := setUpOrders(false)
		appCtx := suite.AppContextWithSessionForTest(&auth.Session{
			ApplicationName: auth.MilApp,
			ServiceMemberID: order.ServiceMemberID,
		})

		file, cleanUpFunc := setUpFileToUpload()
		defer cleanUpFunc()

		fakeS3 := storageTest.NewFakeS3Storage(true)

		suite.NotEqual(uuid.Nil, order.ServiceMemberID, "ServiceMember has ID that is not 0/empty")
		suite.NotEqual(uuid.Nil, order.ServiceMember.UserID, "ServiceMember.UserID has ID that is not 0/empty")

		upload, url, verrs, err := additionalDocumentUploader.CreateAdditionalDocumentsUpload(
			appCtx,
			order.ServiceMember.UserID,
			order.Moves[0].ID,
			file.Data,
			file.Header.Filename,
			fakeS3,
			models.UploadTypeUSER)
		suite.NoError(err)
		suite.NoVerrs(verrs)

		expectedChecksum := "+XM59C3+hSg3Qrs0dPRuUhng5IQTWdYZtmcXhEH0SYU="
		if upload.Checksum != expectedChecksum {
			suite.Fail("Did not calculate the correct MD5: expected %s, got %s", expectedChecksum, upload.Checksum)
		}

		var moveInDB models.Move
		err = suite.DB().
			EagerPreload("AdditionalDocuments").
			Find(&moveInDB, order.Moves[0].ID)

		suite.NoError(err)
		suite.Equal(moveInDB.ID.String(), order.Moves[0].ID.String())
		suite.NotNil(moveInDB.AdditionalDocuments)

		findUpload := models.Upload{}
		err = suite.DB().Find(&findUpload, upload.ID)
		if err != nil {
			suite.Fail("Couldn't find expected upload.")
		}
		suite.Equal(upload.ID.String(), findUpload.ID.String(), "found upload in db")
		suite.NotEmpty(url, "URL is populated")
	})

	suite.Run("Saves userUpload if the document already exists", func() {
		order := setUpOrders(true)

		appCtx := suite.AppContextWithSessionForTest(&auth.Session{
			ApplicationName: auth.MilApp,
			ServiceMemberID: order.ServiceMemberID,
		})

		file, cleanUpFunc := setUpFileToUpload()
		defer cleanUpFunc()

		fakeS3 := storageTest.NewFakeS3Storage(true)

		suite.NotEqual(uuid.Nil, order.ServiceMemberID, "ServiceMember has ID that is not 0/empty")
		suite.NotEqual(uuid.Nil, order.ServiceMember.UserID, "ServiceMember.UserID has ID that is not 0/empty")

		_, _, verrs, err := additionalDocumentUploader.CreateAdditionalDocumentsUpload(
			appCtx,
			order.ServiceMember.UserID,
			order.Moves[0].ID,
			file.Data,
			file.Header.Filename,
			fakeS3,
			models.UploadTypeUSER)
		suite.NoError(err)
		suite.NoVerrs(verrs)
		suite.NoError(err)

		var moveInDB models.Move
		err = suite.DB().
			EagerPreload("AdditionalDocuments").
			Find(&moveInDB, order.Moves[0].ID)

		suite.NoError(err)
		suite.NotNil(moveInDB.ID)
		suite.NotNil(moveInDB.AdditionalDocuments)

		reloadErr := suite.DB().Reload(order.Moves)
		suite.NoError(reloadErr, "error reloading orders")

		suite.Equal(order.Moves[0].AdditionalDocumentsID, moveInDB.AdditionalDocumentsID)
	})
}
