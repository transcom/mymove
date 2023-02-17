package weightticket

import (
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/db/utilities"
	"github.com/transcom/mymove/pkg/models"
	storageTest "github.com/transcom/mymove/pkg/storage/test"
	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/uploader"
)

func (suite *WeightTicketSuite) TestGetWeightTicketFetcher() {
	weightTicketFetcher := NewWeightTicketFetcher()
	var existingWeightTicket models.WeightTicket

	suite.PreloadData(func() {
		existingWeightTicket = testdatagen.MakeDefaultWeightTicket(suite.DB())
	})

	suite.Run("Can fetch a weight ticket with associated uploads for the service member that they belong to", func() {
		appCtx := suite.AppContextWithSessionForTest(&auth.Session{
			ApplicationName: auth.MilApp,
			ServiceMemberID: existingWeightTicket.EmptyDocument.ServiceMember.ID,
		})

		weightTicket, err := weightTicketFetcher.GetWeightTicket(appCtx, existingWeightTicket.ID)

		if suite.NoError(err) && suite.NotNil(weightTicket) {
			suite.Equal(existingWeightTicket.ID, weightTicket.ID)

			suite.Equal(existingWeightTicket.EmptyDocument.ID, weightTicket.EmptyDocument.ID)
			suite.Equal(existingWeightTicket.FullDocument.ID, weightTicket.FullDocument.ID)
			suite.Equal(existingWeightTicket.ProofOfTrailerOwnershipDocument.ID, weightTicket.ProofOfTrailerOwnershipDocument.ID)
		}
	})

	suite.Run("Can fetch a weight ticket with associated uploads for an office user", func() {
		appCtx := suite.AppContextWithSessionForTest(&auth.Session{
			ApplicationName: auth.OfficeApp,
			OfficeUserID:    uuid.Must(uuid.NewV4()),
		})

		weightTicket, err := weightTicketFetcher.GetWeightTicket(appCtx, existingWeightTicket.ID)

		if suite.NoError(err) && suite.NotNil(weightTicket) {
			suite.Equal(existingWeightTicket.ID, weightTicket.ID)

			suite.Equal(existingWeightTicket.EmptyDocument.ID, weightTicket.EmptyDocument.ID)
			suite.Equal(existingWeightTicket.FullDocument.ID, weightTicket.FullDocument.ID)
			suite.Equal(existingWeightTicket.ProofOfTrailerOwnershipDocument.ID, weightTicket.ProofOfTrailerOwnershipDocument.ID)
		}
	})

	suite.Run("Returns a forbidden error when the weight ticket does not belong to the service member", func() {
		appCtx := suite.AppContextWithSessionForTest(&auth.Session{
			ApplicationName: auth.MilApp,
			ServiceMemberID: uuid.Must(uuid.NewV4()),
		})

		weightTicket, err := weightTicketFetcher.GetWeightTicket(appCtx, existingWeightTicket.ID)

		if suite.Error(err) {
			suite.IsType(apperror.ForbiddenError{}, err)

			suite.Nil(weightTicket)
		}
	})

	suite.Run("Does not return weight ticket if it has been deleted", func() {
		appCtx := suite.AppContextWithSessionForTest(&auth.Session{
			ApplicationName: auth.MilApp,
			ServiceMemberID: existingWeightTicket.EmptyDocument.ServiceMember.ID,
		})

		err := utilities.SoftDestroy(appCtx.DB(), &existingWeightTicket)

		suite.FatalNoError(err)

		weightTicket, err := weightTicketFetcher.GetWeightTicket(appCtx, existingWeightTicket.ID)

		if suite.Error(err) {
			suite.IsType(apperror.NotFoundError{}, err)

			suite.Nil(weightTicket)
		}
	})

	suite.Run("Excludes deleted uploads", func() {
		appCtx := suite.AppContextWithSessionForTest(&auth.Session{
			ApplicationName: auth.MilApp,
			ServiceMemberID: existingWeightTicket.EmptyDocument.ServiceMember.ID,
		})

		numValidUploads := len(existingWeightTicket.EmptyDocument.UserUploads)
		suite.FatalTrue(numValidUploads > 0)

		deletedUserUpload := testdatagen.MakeUserUpload(appCtx.DB(), testdatagen.Assertions{
			Document: existingWeightTicket.EmptyDocument,
			UserUpload: models.UserUpload{
				DocumentID: &existingWeightTicket.EmptyDocument.ID,
				Document:   existingWeightTicket.EmptyDocument,
				UploaderID: existingWeightTicket.EmptyDocument.ServiceMember.UserID,
			},
		})

		storer := storageTest.NewFakeS3Storage(true)
		userUploader, err := uploader.NewUserUploader(storer, uploader.MaxCustomerUserUploadFileSizeLimit)

		suite.FatalNoError(err)

		err = userUploader.DeleteUserUpload(appCtx, &deletedUserUpload)

		suite.FatalNoError(err)

		suite.FatalNotNil(deletedUserUpload.Upload.DeletedAt)
		suite.FatalNotNil(deletedUserUpload.DeletedAt)

		weightTicket, err := weightTicketFetcher.GetWeightTicket(appCtx, existingWeightTicket.ID)

		if suite.NoError(err) && suite.NotNil(weightTicket) {
			suite.Equal(existingWeightTicket.ID, weightTicket.ID)

			suite.Equal(existingWeightTicket.EmptyDocument.ID, weightTicket.EmptyDocument.ID)
			suite.Equal(existingWeightTicket.FullDocument.ID, weightTicket.FullDocument.ID)
			suite.Equal(existingWeightTicket.ProofOfTrailerOwnershipDocument.ID, weightTicket.ProofOfTrailerOwnershipDocument.ID)

			suite.Equal(numValidUploads, len(weightTicket.EmptyDocument.UserUploads))

			for _, userUpload := range weightTicket.EmptyDocument.UserUploads {
				suite.NotEqual(deletedUserUpload.ID, userUpload.Upload.ID)
				suite.Nil(userUpload.Upload.DeletedAt)
			}
		}
	})

	suite.Run("Returns a not found error when the weight ticket does not exist", func() {
		appCtx := suite.AppContextWithSessionForTest(&auth.Session{
			ApplicationName: auth.MilApp,
			ServiceMemberID: existingWeightTicket.EmptyDocument.ServiceMember.ID,
		})

		nonexistentWeightTicketID := uuid.Must(uuid.NewV4())

		weightTicket, err := weightTicketFetcher.GetWeightTicket(appCtx, nonexistentWeightTicketID)

		if suite.Error(err) {
			suite.IsType(apperror.NotFoundError{}, err)

			suite.Nil(weightTicket)
		}
	})
}
