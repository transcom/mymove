package ppmshipment

import (
	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/transcom/mymove/pkg/apperror"

	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/mocks"
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

	setUpPPMDocumentFetcherMock := func(returnValue ...interface{}) services.PPMDocumentFetcher {
		mockFetcher := &mocks.PPMDocumentFetcher{}

		mockFetcher.On(
			"GetPPMDocuments",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.AnythingOfType("uuid.UUID"),
		).Return(returnValue...)

		return mockFetcher
	}

	suite.Run("Returns an error if PPM ID is invalid", func() {
		submitter := NewPPMShipmentReviewDocuments(
			setUpPPMShipperRouterMock(nil),
			setUpPPMDocumentFetcherMock(nil),
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
			setUpPPMDocumentFetcherMock(nil),
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
}
