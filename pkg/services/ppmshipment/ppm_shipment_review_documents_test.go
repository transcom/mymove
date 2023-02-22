package ppmshipment

import (
	"fmt"

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

func (suite *PPMShipmentSuite) TestReviewDocuments() {
	setUpPPMShipperRouterMock := func(returnValue ...interface{}) services.PPMShipmentRouter {
		mockRouter := &mocks.PPMShipmentRouter{}

		mockRouter.On(
			"SubmitReviewedDocuments",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.AnythingOfType("*models.PPMShipment"),
		).Return(returnValue...)

		return mockRouter
	}

	suite.Run("Returns an error if PPM ID is invalid", func() {
		submitter := NewPPMShipmentReviewDocuments(
			setUpPPMShipperRouterMock(nil),
		)

		updatedPPMShipment, err := submitter.SubmitReviewedDocuments(
			suite.AppContextForTest(),
			uuid.Nil,
		)

		if suite.Error(err) {
			suite.Nil(updatedPPMShipment)

			suite.IsType(&apperror.BadDataError{}, err)
			suite.Contains(err.Error(), "PPM ID is required")
		}
	})

	suite.Run("Returns an error if PPM shipment does not exist", func() {
		nonexistentPPMShipmentID := uuid.Must(uuid.NewV4())

		submitter := NewPPMShipmentReviewDocuments(
			setUpPPMShipperRouterMock(nil),
		)

		updatedPPMShipment, err := submitter.SubmitReviewedDocuments(
			suite.AppContextForTest(),
			nonexistentPPMShipmentID,
		)

		if suite.Error(err) {
			suite.Nil(updatedPPMShipment)

			suite.IsType(apperror.NotFoundError{}, err)
			suite.Contains(err.Error(), "not found while looking for PPMShipment")
		}
	})

	suite.Run("Returns an error if submitting the close out documentation fails", func() {
		appCtx := suite.AppContextForTest()

		existingPPMShipment := testdatagen.MakePPMShipmentThatNeedsPaymentApproval(appCtx.DB(), testdatagen.Assertions{})

		appCtx = suite.AppContextWithSessionForTest(&auth.Session{
			UserID: existingPPMShipment.Shipment.MoveTaskOrder.Orders.ServiceMember.UserID,
		})

		fakeErr := apperror.NewConflictError(
			existingPPMShipment.ID,
			fmt.Sprintf(
				"PPM shipment documents cannot be submitted because it's not in the %s status.",
				models.PPMShipmentStatusNeedsPaymentApproval,
			),
		)

		submitter := NewPPMShipmentReviewDocuments(
			setUpPPMShipperRouterMock(fakeErr),
		)

		updatedPPMShipment, err := submitter.SubmitReviewedDocuments(
			appCtx,
			existingPPMShipment.ID,
		)

		if suite.Error(err) {
			suite.Nil(updatedPPMShipment)

			suite.IsType(apperror.ConflictError{}, err)
			suite.Equal(fakeErr, err)
		}
	})

	suite.Run("Can route the PPMShipment properly", func() {
		appCtx := suite.AppContextForTest()

		existingPPMShipment := testdatagen.MakePPMShipmentThatNeedsPaymentApproval(appCtx.DB(), testdatagen.Assertions{})

		appCtx = suite.AppContextWithSessionForTest(&auth.Session{
			UserID: existingPPMShipment.Shipment.MoveTaskOrder.Orders.ServiceMember.UserID,
		})

		router := setUpPPMShipperRouterMock(
			func(_ appcontext.AppContext, ppmShipment *models.PPMShipment) error {
				ppmShipment.Status = models.PPMShipmentStatusPaymentApproved

				return nil
			})

		submitter := NewPPMShipmentReviewDocuments(
			router,
		)

		txErr := appCtx.NewTransaction(func(txAppCtx appcontext.AppContext) error {
			txAppCtx.Session()
			updatedPPMShipment, err := submitter.SubmitReviewedDocuments(
				txAppCtx,
				existingPPMShipment.ID,
			)

			if suite.NoError(err) && suite.NotNil(updatedPPMShipment) {
				suite.Equal(models.PPMShipmentStatusPaymentApproved, updatedPPMShipment.Status)

				router.(*mocks.PPMShipmentRouter).AssertCalled(
					suite.T(),
					"SubmitReviewedDocuments",
					txAppCtx,
					mock.AnythingOfType("*models.PPMShipment"),
				)

				return nil
			}

			return err
		})

		suite.NoError(txErr)
	})
}
