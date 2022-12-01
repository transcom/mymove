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

func (suite *WeightTicketSuite) TestListWeightTicketFetcher() {
	weightTicketFetcher := NewWeightTicketFetcher()
	var ppmShipment models.PPMShipment
	var serviceMember models.ServiceMember

	suite.PreloadData(func() {
		existingWeightTicket := testdatagen.MakeDefaultWeightTicket(suite.DB())
		ppmShipment = existingWeightTicket.PPMShipment
		serviceMember = existingWeightTicket.EmptyDocument.ServiceMember

		secondWeightTicket := testdatagen.MakeWeightTicket(suite.DB(), testdatagen.Assertions{
			ServiceMember: serviceMember,
			PPMShipment:   ppmShipment,
		})

		ppmShipment.WeightTickets = models.WeightTickets{existingWeightTicket, secondWeightTicket}
	})

	suite.Run("Can fetch weight tickets with associated uploads for the service member that they belong to", func() {
		appCtx := suite.AppContextWithSessionForTest(&auth.Session{
			ApplicationName: auth.MilApp,
			ServiceMemberID: serviceMember.ID,
		})

		weightTickets, err := weightTicketFetcher.ListWeightTickets(appCtx, ppmShipment.ID)

		if suite.NoError(err) && suite.NotNil(weightTickets) {
			suite.Equal(len(ppmShipment.WeightTickets), len(weightTickets))

			for i, expectedWeightTicket := range ppmShipment.WeightTickets {
				suite.Equal(expectedWeightTicket.ID, weightTickets[i].ID)

				suite.Equal(expectedWeightTicket.EmptyDocument.ID, weightTickets[i].EmptyDocument.ID)
				suite.Equal(expectedWeightTicket.FullDocument.ID, weightTickets[i].FullDocument.ID)
				suite.Equal(expectedWeightTicket.ProofOfTrailerOwnershipDocument.ID, weightTickets[i].ProofOfTrailerOwnershipDocument.ID)

				suite.Equal(len(expectedWeightTicket.EmptyDocument.UserUploads), len(weightTickets[i].EmptyDocument.UserUploads))
				suite.Equal(len(expectedWeightTicket.FullDocument.UserUploads), len(weightTickets[i].FullDocument.UserUploads))
				suite.Equal(len(expectedWeightTicket.ProofOfTrailerOwnershipDocument.UserUploads), len(weightTickets[i].ProofOfTrailerOwnershipDocument.UserUploads))
			}
		}
	})

	suite.Run("Can fetch weight tickets with associated uploads for an office user", func() {
		appCtx := suite.AppContextWithSessionForTest(&auth.Session{
			ApplicationName: auth.OfficeApp,
			OfficeUserID:    uuid.Must(uuid.NewV4()),
		})

		weightTickets, err := weightTicketFetcher.ListWeightTickets(appCtx, ppmShipment.ID)

		if suite.NoError(err) && suite.NotNil(weightTickets) {
			suite.Equal(len(ppmShipment.WeightTickets), len(weightTickets))

			for i, expectedWeightTicket := range ppmShipment.WeightTickets {
				suite.Equal(expectedWeightTicket.ID, weightTickets[i].ID)

				suite.Equal(expectedWeightTicket.EmptyDocument.ID, weightTickets[i].EmptyDocument.ID)
				suite.Equal(expectedWeightTicket.FullDocument.ID, weightTickets[i].FullDocument.ID)
				suite.Equal(expectedWeightTicket.ProofOfTrailerOwnershipDocument.ID, weightTickets[i].ProofOfTrailerOwnershipDocument.ID)

				suite.Equal(len(expectedWeightTicket.EmptyDocument.UserUploads), len(weightTickets[i].EmptyDocument.UserUploads))
				suite.Equal(len(expectedWeightTicket.FullDocument.UserUploads), len(weightTickets[i].FullDocument.UserUploads))
				suite.Equal(len(expectedWeightTicket.ProofOfTrailerOwnershipDocument.UserUploads), len(weightTickets[i].ProofOfTrailerOwnershipDocument.UserUploads))
			}
		}
	})

	suite.Run("Returns a forbidden error when the weight tickets do not belong to the service member", func() {
		appCtx := suite.AppContextWithSessionForTest(&auth.Session{
			ApplicationName: auth.MilApp,
			ServiceMemberID: uuid.Must(uuid.NewV4()),
		})

		weightTickets, err := weightTicketFetcher.ListWeightTickets(appCtx, ppmShipment.ID)

		if suite.Error(err) {
			suite.IsType(apperror.ForbiddenError{}, err)

			suite.Nil(weightTickets)
		}
	})

	suite.Run("Does not return a weight ticket if it has been deleted", func() {
		appCtx := suite.AppContextWithSessionForTest(&auth.Session{
			ApplicationName: auth.MilApp,
			ServiceMemberID: serviceMember.ID,
		})

		err := utilities.SoftDestroy(appCtx.DB(), &ppmShipment.WeightTickets[0])

		suite.FatalNoError(err)

		weightTickets, err := weightTicketFetcher.ListWeightTickets(appCtx, ppmShipment.ID)

		if suite.NoError(err) && suite.NotNil(weightTickets) {
			suite.Equal(len(ppmShipment.WeightTickets)-1, len(weightTickets))

			suite.Equal(ppmShipment.WeightTickets[1].ID, weightTickets[0].ID)
		}
	})

	suite.Run("Does not return any weight tickets if they have all been deleted", func() {
		appCtx := suite.AppContextWithSessionForTest(&auth.Session{
			ApplicationName: auth.MilApp,
			ServiceMemberID: serviceMember.ID,
		})

		for i := range ppmShipment.WeightTickets {
			err := utilities.SoftDestroy(appCtx.DB(), &ppmShipment.WeightTickets[i])

			suite.FatalNoError(err)
		}

		weightTickets, err := weightTicketFetcher.ListWeightTickets(appCtx, ppmShipment.ID)

		if suite.NoError(err) {
			suite.Nil(weightTickets)
		}
	})

	suite.Run("Does not return any weight tickets if the PPMShipment does not exist", func() {
		appCtx := suite.AppContextWithSessionForTest(&auth.Session{
			ApplicationName: auth.MilApp,
			ServiceMemberID: serviceMember.ID,
		})

		nonexistentPPMShipmentID := uuid.Must(uuid.NewV4())

		weightTickets, err := weightTicketFetcher.ListWeightTickets(appCtx, nonexistentPPMShipmentID)

		if suite.NoError(err) {
			suite.Nil(weightTickets)
		}
	})

	suite.Run("Excludes deleted uploads", func() {
		appCtx := suite.AppContextWithSessionForTest(&auth.Session{
			ApplicationName: auth.MilApp,
			ServiceMemberID: serviceMember.ID,
		})

		originalWeightTicket := ppmShipment.WeightTickets[0]
		numValidEmptyUploads := len(originalWeightTicket.EmptyDocument.UserUploads)
		suite.FatalTrue(numValidEmptyUploads > 0)

		deletedUserUpload := testdatagen.MakeUserUpload(appCtx.DB(), testdatagen.Assertions{
			Document: originalWeightTicket.EmptyDocument,
			UserUpload: models.UserUpload{
				DocumentID: &originalWeightTicket.EmptyDocument.ID,
				Document:   originalWeightTicket.EmptyDocument,
				UploaderID: originalWeightTicket.EmptyDocument.ServiceMember.UserID,
			},
		})

		storer := storageTest.NewFakeS3Storage(true)
		userUploader, err := uploader.NewUserUploader(storer, uploader.MaxCustomerUserUploadFileSizeLimit)
		suite.FatalNoError(err)

		err = userUploader.DeleteUserUpload(appCtx, &deletedUserUpload)

		suite.FatalNoError(err)

		suite.FatalNotNil(deletedUserUpload.Upload.DeletedAt)
		suite.FatalNotNil(deletedUserUpload.DeletedAt)

		weightTickets, err := weightTicketFetcher.ListWeightTickets(appCtx, ppmShipment.ID)

		if suite.NoError(err) && suite.NotNil(weightTickets) {
			suite.Equal(len(ppmShipment.WeightTickets), len(weightTickets))

			suite.Equal(originalWeightTicket.ID, weightTickets[0].ID)

			retrievedWeightTicket := weightTickets[0]

			suite.Equal(originalWeightTicket.EmptyDocument.ID, retrievedWeightTicket.EmptyDocument.ID)

			suite.Equal(numValidEmptyUploads, len(retrievedWeightTicket.EmptyDocument.UserUploads))

			for _, userUpload := range retrievedWeightTicket.EmptyDocument.UserUploads {
				suite.NotEqual(deletedUserUpload.ID, userUpload.Upload.ID)
				suite.Nil(userUpload.Upload.DeletedAt)
			}
		}
	})
}
