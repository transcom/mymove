package ppmshipment

import (
	"fmt"
	"time"

	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/mock"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/mocks"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *PPMShipmentSuite) TestSubmitNewCustomerCloseOut() {
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

	suite.Run("Returns an error if PPM shipment does not exist", func() {
		nonexistentPPMShipmentID := uuid.Must(uuid.NewV4())

		submitter := NewPPMShipmentNewSubmitter(
			setUpSignedCertificationCreatorMock(nil, nil),
			setUpPPMShipperRouterMock(nil),
		)

		updatedPPMShipment, err := submitter.SubmitNewCustomerCloseOut(
			suite.AppContextForTest(),
			nonexistentPPMShipmentID,
			models.SignedCertification{},
		)

		if suite.Error(err) {
			suite.Nil(updatedPPMShipment)

			suite.IsType(apperror.NotFoundError{}, err)
			suite.Contains(err.Error(), "not found while looking for PPMShipment")
		}
	})

	suite.Run("Returns an error if creating a new signed certification fails", func() {
		appCtx := suite.AppContextForTest()

		existingPPMShipment := testdatagen.MakePPMShipmentReadyForFinalCustomerCloseOut(appCtx.DB(), testdatagen.Assertions{})

		appCtx = suite.AppContextWithSessionForTest(&auth.Session{
			UserID: existingPPMShipment.Shipment.MoveTaskOrder.Orders.ServiceMember.UserID,
		})

		fakeErr := apperror.NewQueryError("SignedCertification", nil, "Unable to create signed certification")
		creator := setUpSignedCertificationCreatorMock(nil, fakeErr)

		submitter := NewPPMShipmentNewSubmitter(
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
		appCtx := suite.AppContextForTest()

		existingPPMShipment := testdatagen.MakePPMShipmentReadyForFinalCustomerCloseOut(appCtx.DB(), testdatagen.Assertions{})

		appCtx = suite.AppContextWithSessionForTest(&auth.Session{
			UserID: existingPPMShipment.Shipment.MoveTaskOrder.Orders.ServiceMember.UserID,
		})

		fakeErr := apperror.NewConflictError(
			existingPPMShipment.ID,
			"PPM shipment can't be submitted for close out.",
		)
		router := setUpPPMShipperRouterMock(fakeErr)

		submitter := NewPPMShipmentNewSubmitter(
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

	suite.Run("Can create a signed certification and route the PPMShipment properly", func() {
		appCtx := suite.AppContextForTest()

		existingPPMShipment := testdatagen.MakePPMShipmentReadyForFinalCustomerCloseOut(appCtx.DB(), testdatagen.Assertions{})

		appCtx = suite.AppContextWithSessionForTest(&auth.Session{
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
				ppmShipment.Status = models.PPMShipmentStatusNeedsPaymentApproval

				return nil
			})

		submitter := NewPPMShipmentNewSubmitter(
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
				suite.Equal(models.PPMShipmentStatusNeedsPaymentApproval, updatedPPMShipment.Status)

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

				return nil
			}

			// just fulfilling the return type at this point since we already checked for an error
			return err
		})

		suite.NoError(txErr)
	})
}
