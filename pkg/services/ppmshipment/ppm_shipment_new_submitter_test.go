package ppmshipment

import (
	"fmt"
	"time"

	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/mock"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/mocks"
	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/unit"
)

func (suite *PPMShipmentSuite) TestSubmitNewCustomerCloseOut() {
	refectchPPMShipment := func(ppmShipmentID uuid.UUID) *models.PPMShipment {
		// The submitter uses a copier which runs into an issue because of all of the extra references our test data
		// will have filled out because of how our factories work, including some circular references. In practice,
		// we wouldn't have all of those relationships loaded at once, so the copier works fine during regular usage.
		// Here we'll only retrieve the bare minimum.

		var ppmShipment models.PPMShipment

		err := suite.DB().EagerPreload(EagerPreloadAssociationShipment, EagerPreloadAssociationWeightTickets).Find(&ppmShipment, ppmShipmentID)

		suite.FatalNoError(err)

		return &ppmShipment
	}

	setUpPPMShipmentFetcherMock := func(returnValue interface{}, err error) services.PPMShipmentFetcher {
		mockFetcher := &mocks.PPMShipmentFetcher{}

		mockFetcher.On(
			"GetPPMShipment",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.AnythingOfType("uuid.UUID"),
			mock.AnythingOfType("[]string"),
			mock.AnythingOfType("[]string"),
		).Return(returnValue, err)

		return mockFetcher
	}

	setUpSignedCertificationCreatorMock := func(returnValue ...interface{}) services.SignedCertificationCreator {
		mockCreator := &mocks.SignedCertificationCreator{}

		mockCreator.On(
			"CreateSignedCertification",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.AnythingOfType("models.SignedCertification"),
		).Return(returnValue...)

		return mockCreator
	}

	setUpPPMShipperRouterMock := func(returnValue ...interface{}) services.PPMShipmentRouter {
		mockRouter := &mocks.PPMShipmentRouter{}

		mockRouter.On(
			"SubmitCloseOutDocumentation",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.AnythingOfType("*models.PPMShipment"),
		).Return(returnValue...)

		return mockRouter
	}

	suite.Run("Returns an error if PPM ID is invalid", func() {
		submitter := NewPPMShipmentNewSubmitter(
			setUpPPMShipmentFetcherMock(nil, nil),
			setUpSignedCertificationCreatorMock(nil, nil),
			setUpPPMShipperRouterMock(nil),
		)

		updatedPPMShipment, err := submitter.SubmitNewCustomerCloseOut(
			suite.AppContextForTest(),
			uuid.Nil,
			models.SignedCertification{},
		)

		if suite.Error(err) {
			suite.Nil(updatedPPMShipment)

			suite.IsType(&apperror.BadDataError{}, err)
			suite.Contains(err.Error(), "PPM ID is required")
		}
	})

	suite.Run("Returns an error if there is a failure fetching the PPM shipment", func() {
		nonexistentPPMShipmentID := uuid.Must(uuid.NewV4())

		fakeErr := apperror.NewNotFoundError(nonexistentPPMShipmentID, "while looking for PPMShipment")

		submitter := NewPPMShipmentNewSubmitter(
			setUpPPMShipmentFetcherMock(nil, fakeErr),
			setUpSignedCertificationCreatorMock(nil, nil),
			setUpPPMShipperRouterMock(nil),
		)

		updatedPPMShipment, err := submitter.SubmitNewCustomerCloseOut(
			suite.AppContextForTest(),
			nonexistentPPMShipmentID,
			models.SignedCertification{},
		)

		suite.ErrorIs(err, fakeErr)
		suite.Nil(updatedPPMShipment)
	})

	suite.Run("Returns an error if creating a new signed certification fails", func() {
		existingPPMShipment := factory.BuildPPMShipmentReadyForFinalCustomerCloseOut(suite.DB(), nil, nil)

		appCtx := suite.AppContextWithSessionForTest(&auth.Session{
			UserID: existingPPMShipment.Shipment.MoveTaskOrder.Orders.ServiceMember.UserID,
		})

		fakeErr := apperror.NewQueryError("SignedCertification", nil, "Unable to create signed certification")
		creator := setUpSignedCertificationCreatorMock(nil, fakeErr)

		expectedShipment := refectchPPMShipment(existingPPMShipment.ID)
		mockFetcher := setUpPPMShipmentFetcherMock(expectedShipment, nil)

		submitter := NewPPMShipmentNewSubmitter(
			mockFetcher,
			creator,
			setUpPPMShipperRouterMock(nil),
		)

		updatedPPMShipment, err := submitter.SubmitNewCustomerCloseOut(
			appCtx,
			existingPPMShipment.ID,
			models.SignedCertification{},
		)

		if suite.Error(err) {
			suite.Nil(updatedPPMShipment)

			suite.IsType(apperror.QueryError{}, err)
			suite.Equal(fakeErr, err)
		}
	})

	suite.Run("Returns an error if submitting the close out documentation fails", func() {
		existingPPMShipment := factory.BuildPPMShipmentReadyForFinalCustomerCloseOut(suite.DB(), nil, nil)

		appCtx := suite.AppContextWithSessionForTest(&auth.Session{
			UserID: existingPPMShipment.Shipment.MoveTaskOrder.Orders.ServiceMember.UserID,
		})

		fakeErr := apperror.NewConflictError(
			existingPPMShipment.ID,
			"PPM shipment can't be submitted for close out.",
		)
		router := setUpPPMShipperRouterMock(fakeErr)

		expectedShipment := refectchPPMShipment(existingPPMShipment.ID)
		mockFetcher := setUpPPMShipmentFetcherMock(expectedShipment, nil)

		submitter := NewPPMShipmentNewSubmitter(
			mockFetcher,
			setUpSignedCertificationCreatorMock(nil, nil),
			router,
		)

		updatedPPMShipment, err := submitter.SubmitNewCustomerCloseOut(
			appCtx,
			existingPPMShipment.ID,
			models.SignedCertification{},
		)

		if suite.Error(err) {
			suite.Nil(updatedPPMShipment)

			suite.IsType(apperror.ConflictError{}, err)
			suite.Equal(fakeErr, err)
		}
	})

	suite.Run("Can create a signed certification, route the PPMShipment, and calculate allowable weight properly", func() {
		existingPPMShipment := factory.BuildPPMShipmentReadyForFinalCustomerCloseOut(suite.DB(), nil, nil)

		appCtx := suite.AppContextWithSessionForTest(&auth.Session{
			UserID: existingPPMShipment.Shipment.MoveTaskOrder.Orders.ServiceMember.UserID,
		})

		serviceMember := existingPPMShipment.Shipment.MoveTaskOrder.Orders.ServiceMember
		inputSignedCertification := models.SignedCertification{
			CertificationText: "I certify that...",
			Signature:         fmt.Sprintf("%s %s", *serviceMember.FirstName, *serviceMember.LastName),
			Date:              testdatagen.NextValidMoveDate,
		}

		move := existingPPMShipment.Shipment.MoveTaskOrder
		certType := models.SignedCertificationTypePPMPAYMENT

		filledOutSignedCertification := inputSignedCertification
		filledOutSignedCertification.SubmittingUserID = move.Orders.ServiceMember.User.ID
		filledOutSignedCertification.MoveID = move.ID
		filledOutSignedCertification.PpmID = &existingPPMShipment.ID
		filledOutSignedCertification.CertificationType = &certType

		newSignedCertification := filledOutSignedCertification
		now := time.Now()
		newSignedCertification.ID = uuid.Must(uuid.NewV4())
		newSignedCertification.CreatedAt = now
		newSignedCertification.UpdatedAt = now

		creator := setUpSignedCertificationCreatorMock(&newSignedCertification, nil)

		router := setUpPPMShipperRouterMock(
			func(_ appcontext.AppContext, ppmShipment *models.PPMShipment) error {
				ppmShipment.Status = models.PPMShipmentStatusNeedsCloseout

				return nil
			})

		expectedShipment := refectchPPMShipment(existingPPMShipment.ID)
		mockFetcher := setUpPPMShipmentFetcherMock(expectedShipment, nil)

		submitter := NewPPMShipmentNewSubmitter(
			mockFetcher,
			creator,
			router,
		)

		// starting a transaction so that the txAppCtx can be used to check the mock call
		txErr := appCtx.NewTransaction(func(txAppCtx appcontext.AppContext) error {
			txAppCtx.Session()
			updatedPPMShipment, err := submitter.SubmitNewCustomerCloseOut(
				txAppCtx,
				existingPPMShipment.ID,
				inputSignedCertification,
			)

			if suite.NoError(err) && suite.NotNil(updatedPPMShipment) {
				suite.Equal(models.PPMShipmentStatusNeedsCloseout, updatedPPMShipment.Status)

				if suite.NotNil(updatedPPMShipment.SignedCertification) {
					suite.Equal(newSignedCertification.ID, updatedPPMShipment.SignedCertification.ID)
				}

				creator.(*mocks.SignedCertificationCreator).AssertCalled(
					suite.T(),
					"CreateSignedCertification",
					txAppCtx,
					filledOutSignedCertification,
				)

				router.(*mocks.PPMShipmentRouter).AssertCalled(
					suite.T(),
					"SubmitCloseOutDocumentation",
					txAppCtx,
					mock.AnythingOfType("*models.PPMShipment"),
				)

				var expectedAllowableWeight = unit.Pound(0)
				if len(existingPPMShipment.WeightTickets) >= 1 {
					for _, weightTicket := range existingPPMShipment.WeightTickets {
						expectedAllowableWeight += *weightTicket.FullWeight - *weightTicket.EmptyWeight
					}
				}
				if suite.NotNil(updatedPPMShipment.AllowableWeight) {
					suite.Equal(*updatedPPMShipment.AllowableWeight, expectedAllowableWeight)
				}

				return nil
			}

			// just fulfilling the return type at this point since we already checked for an error
			return err
		})

		suite.NoError(txErr)
	})
}
