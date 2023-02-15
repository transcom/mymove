package weightticket

import (
	"fmt"
	"time"

	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/mock"

	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services/mocks"
	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/unit"
)

func (suite *WeightTicketSuite) TestDeleteWeightTicket() {

	setupForTest := func(overrides *models.WeightTicket, hasDocumentUploads bool) *models.WeightTicket {
		serviceMember := testdatagen.MakeDefaultServiceMember(suite.DB())
		ppmShipment := testdatagen.MakeMinimalPPMShipment(suite.DB(), testdatagen.Assertions{
			Order: models.Order{
				ServiceMemberID: serviceMember.ID,
				ServiceMember:   serviceMember,
			},
		})

		emptyDocument := testdatagen.MakeDocument(suite.DB(), testdatagen.Assertions{
			Document: models.Document{
				ServiceMemberID: serviceMember.ID,
			},
		})
		fullDocument := testdatagen.MakeDocument(suite.DB(), testdatagen.Assertions{
			Document: models.Document{
				ServiceMemberID: serviceMember.ID,
			},
		})
		trailerDocument := testdatagen.MakeDocument(suite.DB(), testdatagen.Assertions{
			Document: models.Document{
				ServiceMemberID: serviceMember.ID,
			},
		})

		if hasDocumentUploads {
			for i := 0; i < 2; i++ {
				var deletedAt *time.Time
				if i == 1 {
					deletedAt = models.TimePointer(time.Now())
				}
				testdatagen.MakeUserUpload(suite.DB(), testdatagen.Assertions{
					UserUpload: models.UserUpload{
						UploaderID: serviceMember.UserID,
						DocumentID: &emptyDocument.ID,
						Document:   emptyDocument,
						DeletedAt:  deletedAt,
					},
				})
			}
		}

		originalWeightTicket := models.WeightTicket{
			PPMShipmentID:                     ppmShipment.ID,
			EmptyDocument:                     emptyDocument,
			EmptyDocumentID:                   emptyDocument.ID,
			FullDocument:                      fullDocument,
			FullDocumentID:                    fullDocument.ID,
			ProofOfTrailerOwnershipDocument:   trailerDocument,
			ProofOfTrailerOwnershipDocumentID: trailerDocument.ID,
		}

		if overrides != nil {
			testdatagen.MergeModels(&originalWeightTicket, overrides)
		}

		verrs, err := suite.DB().ValidateAndCreate(&originalWeightTicket)

		suite.NoVerrs(verrs)
		suite.Nil(err)
		suite.NotNil(originalWeightTicket.ID)

		return &originalWeightTicket
	}
	suite.Run("Returns an error if the original doesn't exist", func() {
		notFoundWeightTicketID := uuid.Must(uuid.NewV4())
		fetcher := NewWeightTicketFetcher()
		estimator := mocks.PPMEstimator{}
		deleter := NewWeightTicketDeleter(fetcher, &estimator)

		err := deleter.DeleteWeightTicket(suite.AppContextForTest(), notFoundWeightTicketID)

		if suite.Error(err) {
			suite.IsType(apperror.NotFoundError{}, err)

			suite.Equal(
				fmt.Sprintf("ID: %s not found while looking for WeightTicket", notFoundWeightTicketID.String()),
				err.Error(),
			)
		}
	})

	suite.Run("Successfully deletes as a customer's weight ticket", func() {
		originalWeightTicket := setupForTest(nil, true)
		appCtx := suite.AppContextWithSessionForTest(&auth.Session{
			ApplicationName: auth.MilApp,
			ServiceMemberID: originalWeightTicket.EmptyDocument.ServiceMemberID,
		})

		fetcher := NewWeightTicketFetcher()
		estimator := mocks.PPMEstimator{}
		mockIncentive := unit.Cents(10000)
		estimator.On("FinalIncentiveWithDefaultChecks", appCtx, mock.AnythingOfType("models.PPMShipment"), mock.AnythingOfType("*models.PPMShipment")).Return(&mockIncentive, nil)
		deleter := NewWeightTicketDeleter(fetcher, &estimator)

		suite.Nil(originalWeightTicket.DeletedAt)
		err := deleter.DeleteWeightTicket(appCtx, originalWeightTicket.ID)
		suite.NoError(err)

		var weightTicketInDB models.WeightTicket
		err = suite.DB().Find(&weightTicketInDB, originalWeightTicket.ID)
		suite.NoError(err)
		suite.NotNil(weightTicketInDB.DeletedAt)

		// Should not delete associated PPM shipment
		var dbPPMShipment models.PPMShipment
		suite.NotNil(originalWeightTicket.PPMShipmentID)
		err = suite.DB().Find(&dbPPMShipment, originalWeightTicket.PPMShipmentID)
		suite.NoError(err)
		suite.Nil(dbPPMShipment.DeletedAt)

		// Should delete associated documents
		var dbDocument models.Document
		suite.NotNil(originalWeightTicket.EmptyDocumentID)
		err = suite.DB().Find(&dbDocument, originalWeightTicket.EmptyDocumentID)
		suite.NoError(err)
		suite.NotNil(dbDocument.DeletedAt)

		suite.NotNil(originalWeightTicket.FullDocumentID)
		err = suite.DB().Find(&dbDocument, originalWeightTicket.FullDocumentID)
		suite.NoError(err)
		suite.NotNil(dbDocument.DeletedAt)

		suite.NotNil(originalWeightTicket.ProofOfTrailerOwnershipDocumentID)
		err = suite.DB().Find(&dbDocument, originalWeightTicket.ProofOfTrailerOwnershipDocumentID)
		suite.NoError(err)
		suite.NotNil(dbDocument.DeletedAt)
	})

	suite.Run("Updates final incentive estimate", func() {
		originalWeightTicket := setupForTest(nil, true)
		appCtx := suite.AppContextWithSessionForTest(&auth.Session{
			ApplicationName: auth.MilApp,
			ServiceMemberID: originalWeightTicket.EmptyDocument.ServiceMemberID,
		})

		fetcher := NewWeightTicketFetcher()
		estimator := mocks.PPMEstimator{}
		mockIncentive := unit.Cents(10000)
		estimator.On("FinalIncentiveWithDefaultChecks",
			appCtx,
			mock.AnythingOfType("models.PPMShipment"),
			mock.AnythingOfType("*models.PPMShipment")).Return(&mockIncentive, nil).Once()
		deleter := NewWeightTicketDeleter(fetcher, &estimator)
		err := deleter.DeleteWeightTicket(appCtx, originalWeightTicket.ID)
		suite.NoError(err)

		estimator.AssertCalled(suite.T(), "FinalIncentiveWithDefaultChecks",
			appCtx,
			mock.AnythingOfType("models.PPMShipment"),
			mock.AnythingOfType("*models.PPMShipment"))

		var dbPPMShipment models.PPMShipment
		err = suite.DB().Find(&dbPPMShipment, originalWeightTicket.PPMShipmentID)
		suite.NoError(err)
		suite.Equal(mockIncentive, *dbPPMShipment.FinalIncentive)
	})
}
