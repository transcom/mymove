package ppmshipment

import (
	"fmt"
	"time"

	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/mock"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/etag"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/mocks"
	mtoshipment "github.com/transcom/mymove/pkg/services/mto_shipment"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *PPMShipmentSuite) TestSubmitCustomerCloseOut() {
	setUpSignedCertificationUpdaterMock := func(returnValue ...interface{}) services.SignedCertificationUpdater {
		mockUpdater := &mocks.SignedCertificationUpdater{}

		mockUpdater.On(
			"UpdateSignedCertification",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.AnythingOfType("models.SignedCertification"),
			mock.AnythingOfType("string"),
		).Return(returnValue...)

		return mockUpdater
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
		submitter := NewPPMShipmentUpdatedSubmitter(
			setUpSignedCertificationUpdaterMock(nil, nil),
			setUpPPMShipperRouterMock(nil),
		)

		updatedPPMShipment, err := submitter.SubmitUpdatedCustomerCloseOut(
			suite.AppContextForTest(),
			uuid.Nil,
			models.SignedCertification{},
			"",
		)

		if suite.Error(err) {
			suite.Nil(updatedPPMShipment)

			suite.IsType(&apperror.BadDataError{}, err)
			suite.Contains(err.Error(), "PPM ID is required")
		}
	})

	suite.Run("Returns an error if PPM shipment does not exist", func() {
		nonexistentPPMShipmentID := uuid.Must(uuid.NewV4())

		submitter := NewPPMShipmentUpdatedSubmitter(
			setUpSignedCertificationUpdaterMock(nil, nil),
			setUpPPMShipperRouterMock(nil),
		)

		updatedPPMShipment, err := submitter.SubmitUpdatedCustomerCloseOut(
			suite.AppContextForTest(),
			nonexistentPPMShipmentID,
			models.SignedCertification{},
			"",
		)

		if suite.Error(err) {
			suite.Nil(updatedPPMShipment)

			suite.IsType(apperror.NotFoundError{}, err)
			suite.Contains(err.Error(), "not found while looking for PPMShipment")
		}
	})

	suite.Run("Returns an error if updating an existing signed certification fails", func() {
		appCtx := suite.AppContextForTest()

		existingPPMShipment := testdatagen.MakePPMShipmentThatNeedsPaymentApproval(appCtx.DB(), testdatagen.Assertions{})

		appCtx = suite.AppContextWithSessionForTest(&auth.Session{
			UserID: existingPPMShipment.Shipment.MoveTaskOrder.Orders.ServiceMember.User.ID,
		})

		fakeErr := apperror.NewQueryError("SignedCertification", nil, "Unable to update signed certification")
		updater := setUpSignedCertificationUpdaterMock(nil, fakeErr)

		submitter := NewPPMShipmentUpdatedSubmitter(
			updater,
			setUpPPMShipperRouterMock(nil),
		)

		updatedPPMShipment, err := submitter.SubmitUpdatedCustomerCloseOut(
			appCtx,
			existingPPMShipment.ID,
			*existingPPMShipment.SignedCertification,
			etag.GenerateEtag(existingPPMShipment.SignedCertification.UpdatedAt),
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
			UserID: existingPPMShipment.Shipment.MoveTaskOrder.Orders.ServiceMember.User.ID,
		})

		fakeErr := apperror.NewConflictError(
			existingPPMShipment.ID,
			"PPM shipment can't be submitted for close out.",
		)
		router := setUpPPMShipperRouterMock(fakeErr)

		submitter := NewPPMShipmentUpdatedSubmitter(
			setUpSignedCertificationUpdaterMock(nil, nil),
			router,
		)

		updatedPPMShipment, err := submitter.SubmitUpdatedCustomerCloseOut(
			appCtx,
			existingPPMShipment.ID,
			models.SignedCertification{},
			"",
		)

		if suite.Error(err) {
			suite.Nil(updatedPPMShipment)

			suite.IsType(apperror.ConflictError{}, err)
			suite.Equal(fakeErr, err)
		}
	})

	suite.Run("Can update a signed certification and route the PPMShipment properly", func() {
		appCtx := suite.AppContextForTest()

		existingPPMShipment := testdatagen.MakePPMShipmentThatNeedsPaymentApproval(appCtx.DB(), testdatagen.Assertions{})

		userID := existingPPMShipment.Shipment.MoveTaskOrder.Orders.ServiceMember.User.ID
		appCtx = suite.AppContextWithSessionForTest(&auth.Session{
			UserID: userID,
		})

		realRouter := NewPPMShipmentRouter(mtoshipment.NewShipmentRouter())
		err := realRouter.SendToCustomer(appCtx, &existingPPMShipment)
		suite.FatalNoError(err)

		verrs, err := appCtx.DB().ValidateAndUpdate(&existingPPMShipment)

		suite.FatalNoVerrs(verrs)
		suite.FatalNoError(err)

		serviceMember := existingPPMShipment.Shipment.MoveTaskOrder.Orders.ServiceMember
		inputSignedCertification := models.SignedCertification{
			ID:                existingPPMShipment.SignedCertification.ID,
			CertificationText: "I again certify that...",
			Signature:         fmt.Sprintf("%s %s", *serviceMember.FirstName, *serviceMember.LastName),
			Date:              testdatagen.NextValidMoveDate.AddDate(0, 0, 1),
		}

		existingSignedCertification := existingPPMShipment.SignedCertification

		updatedSignedCertification := existingSignedCertification
		updatedSignedCertification.SubmittingUserID = existingSignedCertification.SubmittingUserID
		updatedSignedCertification.MoveID = existingSignedCertification.MoveID
		updatedSignedCertification.PpmID = existingSignedCertification.PpmID
		updatedSignedCertification.CertificationType = existingSignedCertification.CertificationType
		updatedSignedCertification.CreatedAt = existingSignedCertification.CreatedAt
		updatedSignedCertification.UpdatedAt = time.Now()

		updater := setUpSignedCertificationUpdaterMock(updatedSignedCertification, nil)

		mockRouter := setUpPPMShipperRouterMock(
			func(_ appcontext.AppContext, ppmShipment *models.PPMShipment) error {
				ppmShipment.Status = models.PPMShipmentStatusNeedsPaymentApproval

				return nil
			})

		submitter := NewPPMShipmentUpdatedSubmitter(
			updater,
			mockRouter,
		)

		eTag := etag.GenerateEtag(existingPPMShipment.SignedCertification.UpdatedAt)

		// starting a transaction so that the txAppCtx can be used to check the mock call
		txErr := appCtx.NewTransaction(func(txAppCtx appcontext.AppContext) error {
			updatedPPMShipment, err := submitter.SubmitUpdatedCustomerCloseOut(
				txAppCtx,
				existingPPMShipment.ID,
				inputSignedCertification,
				eTag,
			)

			if suite.NoError(err) && suite.NotNil(updatedPPMShipment) {
				suite.Equal(models.PPMShipmentStatusNeedsPaymentApproval, updatedPPMShipment.Status)

				if suite.NotNil(updatedPPMShipment.SignedCertification) {
					suite.Equal(updatedSignedCertification.ID, updatedPPMShipment.SignedCertification.ID)
					suite.Equal(updatedSignedCertification.CertificationText, updatedPPMShipment.SignedCertification.CertificationText)
					suite.Equal(updatedSignedCertification.Signature, updatedPPMShipment.SignedCertification.Signature)
					suite.True(updatedSignedCertification.Date.Equal(updatedPPMShipment.SignedCertification.Date), "SignedCertification dates should be equal")
					suite.True(updatedSignedCertification.UpdatedAt.Equal(updatedPPMShipment.SignedCertification.UpdatedAt), "SignedCertification UpdatedAt times should be equal")
				}

				updater.(*mocks.SignedCertificationUpdater).AssertCalled(
					suite.T(),
					"UpdateSignedCertification",
					txAppCtx,
					inputSignedCertification,
					eTag,
				)

				return nil
			}

			// just fulfilling the return type at this point since we already checked for an error
			return err
		})

		suite.NoError(txErr)
	})
}
