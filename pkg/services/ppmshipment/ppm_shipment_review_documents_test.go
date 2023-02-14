package ppmshipment

import (
	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/auth"

	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/mocks"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *PPMShipmentSuite) TestReviewDocuments() {
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

		existingPPMShipment := testdatagen.MakePPMShipmentReadyForFinalCustomerCloseOut(appCtx.DB(), testdatagen.Assertions{})

		appCtx = suite.AppContextWithSessionForTest(&auth.Session{
			UserID: existingPPMShipment.Shipment.MoveTaskOrder.Orders.ServiceMember.UserID,
		})

		fakeErr := apperror.NewConflictError(
			existingPPMShipment.ID,
			"PPM shipment can't be submitted for close out.",
		)
		router := setUpPPMShipperRouterMock(fakeErr)

		submitter := NewPPMShipmentReviewDocuments(
			router,
		)

		updatedPPMShipment, err := submitter.SubmitNewCustomerCloseOut(
			appCtx,
			existingPPMShipment.ID,
		)

		if suite.Error(err) {
			suite.Nil(updatedPPMShipment)

			suite.IsType(apperror.ConflictError{}, err)
			suite.Equal(fakeErr, err)
		}
	})
}
