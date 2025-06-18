package ppmshipment

import (
	"fmt"

	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/mock"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/models/roles"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/mocks"
	shipmentsummaryworksheet "github.com/transcom/mymove/pkg/services/shipment_summary_worksheet"
	signedcertification "github.com/transcom/mymove/pkg/services/signed_certification"
	storageTest "github.com/transcom/mymove/pkg/storage/test"
	"github.com/transcom/mymove/pkg/unit"
	"github.com/transcom/mymove/pkg/uploader"
)

func (suite *PPMShipmentSuite) TestReviewDocuments() {
	mockSSWPPMComputer := mocks.SSWPPMComputer{}

	setUpPPMShipperRouterMock := func(returnValue ...interface{}) services.PPMShipmentRouter {
		mockRouter := &mocks.PPMShipmentRouter{}

		mockRouter.On(
			"SubmitReviewedDocuments",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.AnythingOfType("*models.PPMShipment"),
		).Return(returnValue...)

		return mockRouter
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

	suite.Run("Returns an error if PPM ID is invalid", func() {
		submitter := NewPPMShipmentReviewDocuments(
			setUpPPMShipperRouterMock(nil), setUpSignedCertificationCreatorMock(nil, nil),
			setUpSignedCertificationUpdaterMock(nil, nil), &mockSSWPPMComputer,
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
			setUpPPMShipperRouterMock(nil), setUpSignedCertificationCreatorMock(nil, nil),
			setUpSignedCertificationUpdaterMock(nil, nil), &mockSSWPPMComputer,
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
		existingPPMShipment := factory.BuildPPMShipmentThatNeedsCloseout(suite.DB(), nil, nil)

		appCtx := suite.AppContextWithSessionForTest(&auth.Session{
			UserID: existingPPMShipment.Shipment.MoveTaskOrder.Orders.ServiceMember.UserID,
		})

		fakeErr := apperror.NewConflictError(
			existingPPMShipment.ID,
			fmt.Sprintf(
				"PPM shipment documents cannot be submitted because it's not in the %s status.",
				models.PPMShipmentStatusNeedsCloseout,
			),
		)

		submitter := NewPPMShipmentReviewDocuments(
			setUpPPMShipperRouterMock(fakeErr), setUpSignedCertificationCreatorMock(nil, nil),
			setUpSignedCertificationUpdaterMock(nil, nil), &mockSSWPPMComputer,
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

	suite.Run("Incentive-based PPM - Can route the PPMShipment properly", func() {
		existingPPMShipment := factory.BuildPPMShipmentThatNeedsCloseout(suite.DB(), nil, nil)
		sm := factory.BuildServiceMember(suite.DB(), nil, nil)
		session := suite.AppContextWithSessionForTest(&auth.Session{
			ApplicationName: auth.OfficeApp,
			OfficeUserID:    sm.User.ID,
			UserID:          sm.User.ID,
			FirstName:       "Nelson",
			LastName:        "Muntz",
		})
		setupTestData := func() models.OfficeUser {

			transportationOffice := factory.BuildTransportationOffice(suite.DB(), []factory.Customization{
				{
					Model: models.TransportationOffice{
						ProvidesCloseout: true,
					},
				},
			}, nil)

			officeUser := factory.BuildOfficeUserWithRoles(suite.DB(), []factory.Customization{
				{
					Model:    transportationOffice,
					LinkOnly: true,
					Type:     &factory.TransportationOffices.CloseoutOffice,
				},
			}, []roles.RoleType{roles.RoleTypeServicesCounselor})

			return officeUser
		}

		officeUser := setupTestData()
		move := factory.BuildMoveWithShipment(suite.DB(), []factory.Customization{
			{
				Model: models.Move{
					Status: models.MoveStatusNeedsServiceCounseling,
				},
			},
			{
				Model:    officeUser,
				LinkOnly: true,
				Type:     &factory.OfficeUsers.SCAssignedUser,
			},
		}, nil)

		existingPPMShipment.Shipment.MoveTaskOrder = move
		suite.NotNil(existingPPMShipment.Shipment.MoveTaskOrder.SCAssignedID)

		router := setUpPPMShipperRouterMock(
			func(_ appcontext.AppContext, ppmShipment *models.PPMShipment) error {
				ppmShipment.Status = models.PPMShipmentStatusCloseoutComplete

				return nil
			})

		mockPPMCloseoutFetcher := &mocks.PPMCloseoutFetcher{}
		ppmEstimator := &mocks.PPMEstimator{}
		maxIncentive := 987654
		ppmEstimator.On("MaxIncentive",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.AnythingOfType("models.PPMShipment"),
			mock.AnythingOfType("*models.PPMShipment")).
			Return(models.CentPointer(unit.Cents(maxIncentive)), nil)
		SSWPPMComputer := shipmentsummaryworksheet.NewSSWPPMComputer(mockPPMCloseoutFetcher, ppmEstimator)
		mockPPMCloseoutFetcher.On("GetActualWeight", mock.AnythingOfType("*models.PPMShipment")).Return(unit.Pound(1000))
		submitter := NewPPMShipmentReviewDocuments(
			router, signedcertification.NewSignedCertificationCreator(), signedcertification.NewSignedCertificationUpdater(), SSWPPMComputer,
		)

		txErr := session.NewTransaction(func(txAppCtx appcontext.AppContext) error {
			txAppCtx.Session()
			updatedPPMShipment, err := submitter.SubmitReviewedDocuments(
				txAppCtx,
				existingPPMShipment.ID,
			)

			//check removal of the SC Assigned User
			suite.Nil(updatedPPMShipment.Shipment.MoveTaskOrder.SCAssignedID)

			if suite.NoError(err) && suite.NotNil(updatedPPMShipment) {
				suite.Equal(models.PPMShipmentStatusCloseoutComplete, updatedPPMShipment.Status)

				router.(*mocks.PPMShipmentRouter).AssertCalled(
					suite.T(),
					"SubmitReviewedDocuments",
					txAppCtx,
					mock.AnythingOfType("*models.PPMShipment"),
				)

				closeout, err := models.FetchPPMCloseoutByPPMID(suite.DB(), updatedPPMShipment.ID)
				suite.NoError(err)
				suite.NotNil(closeout)
				suite.Equal(unit.Cents(401300), *closeout.MemberDisbursement)

				return nil
			}
			return err
		})

		suite.NoError(txErr)

		certs, err := models.FetchSignedCertificationPPMByType(suite.DB(), session.Session(), existingPPMShipment.Shipment.MoveTaskOrderID, existingPPMShipment.ID, models.SignedCertificationTypeCloseoutReviewedPPMPAYMENT)
		suite.NotNil(certs)
		suite.Nil(err)
		suite.True(len(certs) == 1)

		// run SubmitReviewedDocuments again to test certification is only added once. will update current if exist
		// this is to simulate N times customer and reviewer can resubmit/reconfirm
		txErr = session.NewTransaction(func(txAppCtx appcontext.AppContext) error {
			txAppCtx.Session()
			updatedPPMShipment, err := submitter.SubmitReviewedDocuments(
				txAppCtx,
				existingPPMShipment.ID,
			)

			if suite.NoError(err) && suite.NotNil(updatedPPMShipment) {
				suite.Equal(models.PPMShipmentStatusCloseoutComplete, updatedPPMShipment.Status)

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

		// verfify only one certification record is present
		certs, err = models.FetchSignedCertificationPPMByType(suite.DB(), session.Session(), existingPPMShipment.Shipment.MoveTaskOrderID, existingPPMShipment.ID, models.SignedCertificationTypeCloseoutReviewedPPMPAYMENT)
		suite.NotNil(certs)
		suite.Nil(err)
		suite.True(len(certs) == 1)
	})

	suite.Run("Small package PPM - Can route the PPMShipment properly", func() {
		fakeS3 := storageTest.NewFakeS3Storage(true)
		userUploader, uploaderErr := uploader.NewUserUploader(fakeS3, uploader.MaxCustomerUserUploadFileSizeLimit)
		suite.FatalNoError(uploaderErr)

		approvedStatus := models.PPMDocumentStatusApproved
		trackingNumber := "TRK1234"
		isProGear := true
		proGearBelongsToSelf := true
		proGearDescription := "Pro gear updated description"
		weightShipped := 2000
		spr := models.MovingExpenseReceiptTypeSmallPackage
		// this factory has two small package member-paid moving expenses that total 4000 pounds
		existingPPMShipment := factory.BuildPPMSPRShipmentWithoutPaymentPacketTwoExpenses(suite.DB(), userUploader)
		// adding one more expense that will be GTCC paid
		factory.BuildMovingExpense(suite.DB(), []factory.Customization{
			{
				Model:    existingPPMShipment,
				LinkOnly: true,
			},
			{
				Model: models.MovingExpense{
					MovingExpenseType:    &spr,
					Status:               &approvedStatus,
					PaidWithGTCC:         models.BoolPointer(true),
					MissingReceipt:       models.BoolPointer(false),
					Amount:               models.CentPointer(unit.Cents(8675309)),
					TrackingNumber:       &trackingNumber,
					IsProGear:            &isProGear,
					ProGearBelongsToSelf: &proGearBelongsToSelf,
					ProGearDescription:   &proGearDescription,
					WeightShipped:        (*unit.Pound)(&weightShipped),
				},
			},
		}, nil)

		sm := factory.BuildServiceMember(suite.DB(), nil, nil)
		session := suite.AppContextWithSessionForTest(&auth.Session{
			ApplicationName: auth.OfficeApp,
			OfficeUserID:    sm.User.ID,
			UserID:          sm.User.ID,
			FirstName:       "Nelson",
			LastName:        "Muntz",
		})
		setupTestData := func() models.OfficeUser {

			transportationOffice := factory.BuildTransportationOffice(suite.DB(), []factory.Customization{
				{
					Model: models.TransportationOffice{
						ProvidesCloseout: true,
					},
				},
			}, nil)

			officeUser := factory.BuildOfficeUserWithRoles(suite.DB(), []factory.Customization{
				{
					Model:    transportationOffice,
					LinkOnly: true,
					Type:     &factory.TransportationOffices.CloseoutOffice,
				},
			}, []roles.RoleType{roles.RoleTypeServicesCounselor})

			return officeUser
		}

		officeUser := setupTestData()
		move := factory.BuildMoveWithShipment(suite.DB(), []factory.Customization{
			{
				Model: models.Move{
					Status: models.MoveStatusNeedsServiceCounseling,
				},
			},
			{
				Model:    officeUser,
				LinkOnly: true,
				Type:     &factory.OfficeUsers.SCAssignedUser,
			},
		}, nil)

		existingPPMShipment.Shipment.MoveTaskOrder = move
		suite.NotNil(existingPPMShipment.Shipment.MoveTaskOrder.SCAssignedID)

		router := setUpPPMShipperRouterMock(
			func(_ appcontext.AppContext, ppmShipment *models.PPMShipment) error {
				ppmShipment.Status = models.PPMShipmentStatusCloseoutComplete

				return nil
			})

		mockPPMCloseoutFetcher := &mocks.PPMCloseoutFetcher{}
		ppmEstimator := &mocks.PPMEstimator{}
		maxIncentive := 987654
		ppmEstimator.On("MaxIncentive",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.AnythingOfType("models.PPMShipment"),
			mock.AnythingOfType("*models.PPMShipment")).
			Return(models.CentPointer(unit.Cents(maxIncentive)), nil)
		SSWPPMComputer := shipmentsummaryworksheet.NewSSWPPMComputer(mockPPMCloseoutFetcher, ppmEstimator)
		mockPPMCloseoutFetcher.On("GetActualWeight", mock.AnythingOfType("*models.PPMShipment")).Return(unit.Pound(1000))
		submitter := NewPPMShipmentReviewDocuments(
			router, signedcertification.NewSignedCertificationCreator(), signedcertification.NewSignedCertificationUpdater(), SSWPPMComputer,
		)

		txErr := session.NewTransaction(func(txAppCtx appcontext.AppContext) error {
			txAppCtx.Session()
			updatedPPMShipment, err := submitter.SubmitReviewedDocuments(
				txAppCtx,
				existingPPMShipment.ID,
			)

			//check removal of the SC Assigned User
			suite.Nil(updatedPPMShipment.Shipment.MoveTaskOrder.SCAssignedID)

			if suite.NoError(err) && suite.NotNil(updatedPPMShipment) {
				suite.Equal(models.PPMShipmentStatusCloseoutComplete, updatedPPMShipment.Status)

				router.(*mocks.PPMShipmentRouter).AssertCalled(
					suite.T(),
					"SubmitReviewedDocuments",
					txAppCtx,
					mock.AnythingOfType("*models.PPMShipment"),
				)

				closeout, err := models.FetchPPMCloseoutByPPMID(suite.DB(), updatedPPMShipment.ID)
				suite.NoError(err)
				suite.NotNil(closeout)
				suite.Equal(unit.Cents(8675309), *closeout.GTCCPaidSmallPackage)
				suite.Equal(unit.Cents(17350618), *closeout.MemberPaidSmallPackage)

				suite.Equal(unit.Cents(8675309), *closeout.TotalGTCCPaidExpenses)
				suite.Equal(unit.Cents(17350618), *closeout.TotalMemberPaidExpenses)

				return nil
			}
			return err
		})

		suite.NoError(txErr)

		certs, err := models.FetchSignedCertificationPPMByType(suite.DB(), session.Session(), existingPPMShipment.Shipment.MoveTaskOrderID, existingPPMShipment.ID, models.SignedCertificationTypeCloseoutReviewedPPMPAYMENT)
		suite.NotNil(certs)
		suite.Nil(err)
		suite.True(len(certs) == 1)

		// run SubmitReviewedDocuments again to test certification is only added once. will update current if exist
		// this is to simulate N times customer and reviewer can resubmit/reconfirm
		txErr = session.NewTransaction(func(txAppCtx appcontext.AppContext) error {
			txAppCtx.Session()
			updatedPPMShipment, err := submitter.SubmitReviewedDocuments(
				txAppCtx,
				existingPPMShipment.ID,
			)

			if suite.NoError(err) && suite.NotNil(updatedPPMShipment) {
				suite.Equal(models.PPMShipmentStatusCloseoutComplete, updatedPPMShipment.Status)

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

		// verfify only one certification record is present
		certs, err = models.FetchSignedCertificationPPMByType(suite.DB(), session.Session(), existingPPMShipment.Shipment.MoveTaskOrderID, existingPPMShipment.ID, models.SignedCertificationTypeCloseoutReviewedPPMPAYMENT)
		suite.NotNil(certs)
		suite.Nil(err)
		suite.True(len(certs) == 1)
	})

	suite.Run("Dollar input handled properly", func() {
		dollarInput := "$245.78"
		expected := 24578

		result, err := priceToCents(dollarInput)
		suite.Nil(err)
		suite.Equal(expected, result)
	})

	suite.Run("Non-dollar input handled properly, returns 0", func() {
		nonDollarInput := ""
		expected := 0

		result, err := priceToCents(nonDollarInput)
		suite.Nil(err)
		suite.Equal(expected, result)
	})
}
