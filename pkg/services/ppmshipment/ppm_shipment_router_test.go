package ppmshipment

import (
	"fmt"
	"time"

	"github.com/stretchr/testify/mock"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/mocks"
)

func setUpPPMShipmentRouter(mtoShipmentRouterMethod string, mtoShipmentRouterReturnValue ...interface{}) services.PPMShipmentRouter {
	mtoShipmentRouter := &mocks.ShipmentRouter{}

	mtoShipmentRouter.
		On(
			mtoShipmentRouterMethod,
			mock.AnythingOfType("*appcontext.appContext"),
			mock.AnythingOfType("*models.MTOShipment"),
		).
		Return(mtoShipmentRouterReturnValue...)

	return NewPPMShipmentRouter(mtoShipmentRouter)
}

func (suite *PPMShipmentSuite) TestSetToDraft() {
	mtoShipmentRouterMethodToMock := ""

	suite.Run(fmt.Sprintf("Can set status to %s", models.PPMShipmentStatusDraft), func() {
		ppmShipment := models.PPMShipment{}

		ppmShipmentRouter := setUpPPMShipmentRouter(mtoShipmentRouterMethodToMock, nil)

		err := ppmShipmentRouter.SetToDraft(suite.AppContextForTest(), &ppmShipment)

		suite.NoError(err)
		suite.Equal(models.PPMShipmentStatusDraft, ppmShipment.Status)
		suite.Equal(models.MTOShipmentStatusDraft, ppmShipment.Shipment.Status)
	})

	suite.Run(fmt.Sprintf("Can't set status to %s if it's not new", models.PPMShipmentStatusDraft), func() {
		ppmShipment := factory.BuildPPMShipment(nil, nil, nil)
		originalPPMShipment := ppmShipment

		ppmShipmentRouter := setUpPPMShipmentRouter(mtoShipmentRouterMethodToMock, nil)

		err := ppmShipmentRouter.SetToDraft(suite.AppContextForTest(), &ppmShipment)

		if suite.Error(err) {
			suite.IsType(apperror.ConflictError{}, err)
			suite.Contains(
				err.Error(),
				fmt.Sprintf("PPM shipment can't be set to %s because it's not new.", models.PPMShipmentStatusDraft),
			)

			suite.Equal(originalPPMShipment.Status, ppmShipment.Status)
			suite.Equal(originalPPMShipment.Shipment.Status, ppmShipment.Shipment.Status)
		}
	})
}

func (suite *PPMShipmentSuite) TestSubmit() {
	mtoShipmentRouterMethodToMock := "Submit"

	successTestCases := map[string]struct {
		mtoShipmentStatus models.MTOShipmentStatus
		ppmShipmentStatus models.PPMShipmentStatus
	}{
		"a new PPM Shipment": {
			mtoShipmentStatus: models.MTOShipmentStatus(""),
			ppmShipmentStatus: models.PPMShipmentStatus(""),
		},
		fmt.Sprintf("a PPM Shipment in %s status", models.PPMShipmentStatusDraft): {
			mtoShipmentStatus: models.MTOShipmentStatusDraft,
			ppmShipmentStatus: models.PPMShipmentStatusDraft,
		},
	}

	for currentStatus, testCase := range successTestCases {
		currentStatus := currentStatus
		testCase := testCase

		suite.Run(fmt.Sprintf("Can set status to %s if it's %s", models.PPMShipmentStatusSubmitted, currentStatus), func() {
			ppmShipment := models.PPMShipment{
				Status:   testCase.ppmShipmentStatus,
				Shipment: models.MTOShipment{Status: testCase.mtoShipmentStatus},
			}

			ppmShipmentRouter := setUpPPMShipmentRouter(
				mtoShipmentRouterMethodToMock,
				func(_ appcontext.AppContext, mtoShipment *models.MTOShipment) error {
					mtoShipment.Status = models.MTOShipmentStatusSubmitted
					return nil
				},
			)

			err := ppmShipmentRouter.Submit(suite.AppContextForTest(), &ppmShipment)

			suite.NoError(err)

			suite.Equal(models.PPMShipmentStatusSubmitted, ppmShipment.Status)
			suite.Equal(models.MTOShipmentStatusSubmitted, ppmShipment.Shipment.Status)
		})
	}

	suite.Run(fmt.Sprintf("Can't set status to %s if the MTOShipment router returns an error", models.PPMShipmentStatusSubmitted), func() {
		ppmShipment := factory.BuildMinimalPPMShipment(nil, nil, nil)

		// Not using the real error that gets returned because it's fields are private and we don't export a constructor
		fakeMTOShipmentRouterErr := apperror.NewConflictError(ppmShipment.Shipment.ID, "can't submit shipment")

		ppmShipmentRouter := setUpPPMShipmentRouter(mtoShipmentRouterMethodToMock, fakeMTOShipmentRouterErr)

		err := ppmShipmentRouter.Submit(suite.AppContextForTest(), &ppmShipment)

		if suite.Error(err) {
			suite.IsType(apperror.ConflictError{}, err)
			suite.Equal(fakeMTOShipmentRouterErr.Error(), err.Error())
		}
	})

	suite.Run(fmt.Sprintf("Can't set status to %s if it's not new or in the %s status", models.PPMShipmentStatusSubmitted, models.PPMShipmentStatusDraft), func() {
		ppmShipment := factory.BuildPPMShipment(nil, nil, []factory.Trait{factory.GetTraitApprovedPPMShipment})
		originalPPMShipmentStatus := ppmShipment.Status
		originalMTOShipmentStatus := ppmShipment.Shipment.Status

		ppmShipmentRouter := setUpPPMShipmentRouter(mtoShipmentRouterMethodToMock, nil)

		err := ppmShipmentRouter.Submit(suite.AppContextForTest(), &ppmShipment)

		if suite.Error(err) {
			suite.IsType(apperror.ConflictError{}, err)
			suite.Contains(
				err.Error(),
				fmt.Sprintf("PPM shipment can't be set to %s because it's not new or in the %s status.", models.PPMShipmentStatusSubmitted, models.PPMShipmentStatusDraft),
			)

			suite.Equal(originalPPMShipmentStatus, ppmShipment.Status)
			suite.Equal(originalMTOShipmentStatus, ppmShipment.Shipment.Status)
		}
	})
}

func (suite *PPMShipmentSuite) TestSendToCustomer() {
	mtoShipmentRouterMethodToMock := "Approve"

	successTestCases := map[models.PPMShipmentStatus]func() models.PPMShipment{
		models.PPMShipmentStatusSubmitted: func() models.PPMShipment {
			return factory.BuildPPMShipment(nil, nil, nil)
		},
		models.PPMShipmentStatusNeedsCloseout: func() models.PPMShipment {
			return factory.BuildPPMShipmentThatNeedsCloseout(nil, nil, nil)
		},
	}

	for currentStatus, makePPMShipment := range successTestCases {
		currentStatus := currentStatus
		makePPMShipment := makePPMShipment

		suite.Run(fmt.Sprintf("Can set status to %s if it is currently %s", models.PPMShipmentStatusWaitingOnCustomer, currentStatus), func() {
			ppmShipment := makePPMShipment()

			ppmShipmentRouter := setUpPPMShipmentRouter(
				mtoShipmentRouterMethodToMock,
				func(_ appcontext.AppContext, mtoShipment *models.MTOShipment) error {
					mtoShipment.Status = models.MTOShipmentStatusApproved
					mtoShipment.ApprovedDate = models.TimePointer(time.Now())
					return nil
				},
			)

			err := ppmShipmentRouter.SendToCustomer(suite.AppContextForTest(), &ppmShipment)

			if suite.NoError(err) {
				suite.Equal(models.PPMShipmentStatusWaitingOnCustomer, ppmShipment.Status)
				suite.Equal(models.MTOShipmentStatusApproved, ppmShipment.Shipment.Status)

				suite.NotNil(ppmShipment.ApprovedAt)
				suite.NotNil(ppmShipment.Shipment.ApprovedDate)
				suite.True(
					ppmShipment.Shipment.ApprovedDate.Equal(*ppmShipment.ApprovedAt),
					"PPMShipment.ApprovedAt and MTOShipment.ApprovedDate should be equal",
				)
			}

		})
	}

	statusFailureTestCases := map[models.PPMShipmentStatus]func() models.PPMShipment{
		models.PPMShipmentStatusDraft: func() models.PPMShipment {
			return factory.BuildMinimalPPMShipment(nil, []factory.Customization{
				{
					Model: models.MTOShipment{
						Status: models.MTOShipmentStatusDraft,
					},
				},
			}, nil)
		},
		models.PPMShipmentStatusWaitingOnCustomer: func() models.PPMShipment {
			return factory.BuildPPMShipment(nil, nil, []factory.Trait{factory.GetTraitApprovedPPMShipment})
		},
	}

	for currentStatus, makePPMShipment := range statusFailureTestCases {
		currentStatus := currentStatus
		makePPMShipment := makePPMShipment

		suite.Run(fmt.Sprintf("Can't set status to %s if it is currently %s", models.PPMShipmentStatusWaitingOnCustomer, currentStatus), func() {
			ppmShipment := makePPMShipment()
			originalMTOShipmentStatus := ppmShipment.Shipment.Status

			ppmShipmentRouter := setUpPPMShipmentRouter(mtoShipmentRouterMethodToMock, nil)

			err := ppmShipmentRouter.SendToCustomer(suite.AppContextForTest(), &ppmShipment)

			if suite.Error(err) {
				suite.IsType(apperror.ConflictError{}, err)
				suite.Contains(
					err.Error(),
					fmt.Sprintf(
						"PPM shipment can't be set to %s because it's not in a %s or %s status.",
						models.PPMShipmentStatusWaitingOnCustomer,
						models.PPMShipmentStatusSubmitted,
						models.PPMShipmentStatusNeedsCloseout,
					),
				)

				suite.Equal(currentStatus, ppmShipment.Status)
				suite.Equal(originalMTOShipmentStatus, ppmShipment.Shipment.Status)
			}
		})
	}

	suite.Run("Can't set status to WaitingOnCustomer if MTOShipment can't be approved", func() {
		ppmShipment := factory.BuildPPMShipment(nil, nil, nil)

		// Not using the real error that gets returned because it's fields are private and we don't export a constructor
		fakeMTOShipmentRouterErr := apperror.NewConflictError(ppmShipment.Shipment.ID, "error approving MTOShipment")

		ppmShipmentRouter := setUpPPMShipmentRouter(
			mtoShipmentRouterMethodToMock,
			func(_ appcontext.AppContext, _ *models.MTOShipment) error {
				return fakeMTOShipmentRouterErr
			},
		)

		err := ppmShipmentRouter.SendToCustomer(suite.AppContextForTest(), &ppmShipment)

		if suite.Error(err) {
			suite.IsType(apperror.ConflictError{}, err)
			suite.Equal(fakeMTOShipmentRouterErr.Error(), err.Error())

			suite.Equal(models.PPMShipmentStatusSubmitted, ppmShipment.Status)
			suite.Equal(models.MTOShipmentStatusSubmitted, ppmShipment.Shipment.Status)
		}
	})

	suite.Run("Skips approving MTOShipment if it is already approved", func() {
		ppmShipment := factory.BuildPPMShipment(nil, nil, []factory.Trait{factory.GetTraitApprovedPPMShipment})
		ppmShipment.Status = models.PPMShipmentStatusSubmitted

		mtoShipmentRouter := &mocks.ShipmentRouter{}

		ppmShipmentRouter := NewPPMShipmentRouter(mtoShipmentRouter)

		err := ppmShipmentRouter.SendToCustomer(suite.AppContextForTest(), &ppmShipment)

		if suite.NoError(err) {
			suite.Equal(models.PPMShipmentStatusWaitingOnCustomer, ppmShipment.Status)
			suite.Equal(models.MTOShipmentStatusApproved, ppmShipment.Shipment.Status)

			mtoShipmentRouter.AssertNotCalled(suite.T(), mtoShipmentRouterMethodToMock)
		}
	})

	suite.Run("Doesn't set a new approval time if there is one already.", func() {
		ppmShipment := factory.BuildPPMShipment(nil, nil, []factory.Trait{factory.GetTraitApprovedPPMShipment})

		differentApprovedAt := time.Now().AddDate(0, 0, 1)
		ppmShipment.ApprovedAt = &differentApprovedAt
		ppmShipment.Status = models.PPMShipmentStatusSubmitted

		ppmShipmentRouter := NewPPMShipmentRouter(&mocks.ShipmentRouter{})

		err := ppmShipmentRouter.SendToCustomer(suite.AppContextForTest(), &ppmShipment)

		if suite.NoError(err) {
			suite.Equal(models.PPMShipmentStatusWaitingOnCustomer, ppmShipment.Status)
			suite.Equal(models.MTOShipmentStatusApproved, ppmShipment.Shipment.Status)

			suite.True(differentApprovedAt.Equal(*ppmShipment.ApprovedAt), "ApprovedAt should not have changed")
		}
	})
}

func (suite *PPMShipmentSuite) TestSubmitCloseOutDocumentation() {
	mtoShipmentRouterMethodToMock := ""

	suite.Run(fmt.Sprintf("Can set status to %s if it is currently %s", models.PPMShipmentStatusNeedsCloseout, models.PPMShipmentStatusWaitingOnCustomer), func() {
		ppmShipment := factory.BuildPPMShipmentReadyForFinalCustomerCloseOut(nil, nil, nil)

		ppmShipmentRouter := setUpPPMShipmentRouter(mtoShipmentRouterMethodToMock, nil)

		err := ppmShipmentRouter.SubmitCloseOutDocumentation(suite.AppContextForTest(), &ppmShipment)

		if suite.NoError(err) {
			suite.Equal(models.PPMShipmentStatusNeedsCloseout, ppmShipment.Status)

			suite.NotNil(ppmShipment.SubmittedAt)
		}
	})

	suite.Run("Does not set the SubmittedAt time if it is already set", func() {
		ppmShipment := factory.BuildPPMShipmentThatNeedsToBeResubmitted(nil, nil, nil)

		suite.FatalNotNil(ppmShipment.SubmittedAt)
		originalSubmittedAt := *ppmShipment.SubmittedAt

		if !suite.Equal(models.PPMShipmentStatusWaitingOnCustomer, ppmShipment.Status) {
			suite.Failf(
				"Test data is in an unexpected state",
				"Expected PPMShipment to be in %s status",
				models.PPMShipmentStatusWaitingOnCustomer,
			)
		}

		ppmShipmentRouter := setUpPPMShipmentRouter(mtoShipmentRouterMethodToMock, nil)

		err := ppmShipmentRouter.SubmitCloseOutDocumentation(suite.AppContextForTest(), &ppmShipment)

		if suite.NoError(err) {
			suite.Equal(models.PPMShipmentStatusNeedsCloseout, ppmShipment.Status)

			suite.True(originalSubmittedAt.Equal(*ppmShipment.SubmittedAt), "SubmittedAt should not have changed")
		}
	})

	statusFailureTestCases := map[models.PPMShipmentStatus]func() models.PPMShipment{
		models.PPMShipmentStatusDraft: func() models.PPMShipment {
			return factory.BuildMinimalPPMShipment(nil, []factory.Customization{
				{
					Model: models.MTOShipment{
						Status: models.MTOShipmentStatusDraft,
					},
				},
			}, nil)
		},
		models.PPMShipmentStatusSubmitted: func() models.PPMShipment {
			return factory.BuildPPMShipment(nil, nil, nil)
		},
		models.PPMShipmentStatusNeedsCloseout: func() models.PPMShipment {
			return factory.BuildPPMShipmentThatNeedsCloseout(nil, nil, nil)
		},
	}

	for currentStatus, makePPMShipment := range statusFailureTestCases {
		currentStatus := currentStatus
		makePPMShipment := makePPMShipment

		suite.Run(fmt.Sprintf("Can't set status to %s if it is currently %s", models.PPMShipmentStatusNeedsCloseout, currentStatus), func() {
			ppmShipment := makePPMShipment()

			ppmShipmentRouter := setUpPPMShipmentRouter(mtoShipmentRouterMethodToMock, nil)

			err := ppmShipmentRouter.SubmitCloseOutDocumentation(suite.AppContextForTest(), &ppmShipment)

			if suite.Error(err) {
				suite.IsType(apperror.ConflictError{}, err)
				suite.Contains(
					err.Error(),
					fmt.Sprintf(
						"PPM shipment can't be set to %s because it's not in the %s status.",
						models.PPMShipmentStatusNeedsCloseout,
						models.PPMShipmentStatusWaitingOnCustomer,
					),
				)

				suite.Equal(currentStatus, ppmShipment.Status)
			}
		})
	}
}

func (suite *PPMShipmentSuite) TestSubmitReviewPPMDocuments() {
	mtoShipmentRouterMethodToMock := ""

	suite.Run("Update PPMShipment Status to CLOSEOUT_COMPLETE when there are rejected weight tickets", func() {
		ppmShipment := factory.BuildPPMShipmentReadyForFinalCustomerCloseOut(nil, nil, nil)
		rejected := models.PPMDocumentStatusRejected
		weightTicket := factory.BuildWeightTicket(suite.DB(), []factory.Customization{
			{
				Model: models.WeightTicket{
					Status: &rejected,
				},
			},
		}, nil)
		ppmShipment.Status = models.PPMShipmentStatusNeedsCloseout
		movingExpense := factory.BuildMovingExpense(suite.DB(), nil, nil)
		progear := factory.BuildProgearWeightTicket(suite.DB(), nil, nil)

		ppmShipmentRouter := setUpPPMShipmentRouter(mtoShipmentRouterMethodToMock, nil)
		ppmShipment.WeightTickets = models.WeightTickets{weightTicket}
		ppmShipment.MovingExpenses = models.MovingExpenses{movingExpense}
		ppmShipment.ProgearWeightTickets = models.ProgearWeightTickets{progear}
		err := ppmShipmentRouter.SubmitReviewedDocuments(suite.AppContextForTest(), &ppmShipment)

		if suite.NoError(err) {
			suite.Equal(models.PPMShipmentStatusCloseoutComplete, ppmShipment.Status)
		}
	})

	suite.Run("Update PPMShipment Status to CLOSEOUT_COMPLETE when there are rejected  progear weight tickets", func() {
		ppmShipment := factory.BuildPPMShipmentReadyForFinalCustomerCloseOut(nil, nil, nil)
		ppmShipment.Status = models.PPMShipmentStatusNeedsCloseout
		rejected := models.PPMDocumentStatusRejected
		progear := factory.BuildProgearWeightTicket(suite.DB(), []factory.Customization{
			{
				Model: models.ProgearWeightTicket{
					Status: &rejected,
				},
			},
		}, nil)
		movingExpense := factory.BuildMovingExpense(suite.DB(), nil, nil)
		weightTicket := factory.BuildWeightTicket(suite.DB(), nil, nil)

		ppmShipmentRouter := setUpPPMShipmentRouter(mtoShipmentRouterMethodToMock, nil)
		ppmShipment.WeightTickets = models.WeightTickets{weightTicket}
		ppmShipment.MovingExpenses = models.MovingExpenses{movingExpense}
		ppmShipment.ProgearWeightTickets = models.ProgearWeightTickets{progear}
		err := ppmShipmentRouter.SubmitReviewedDocuments(suite.AppContextForTest(), &ppmShipment)

		if suite.NoError(err) {
			suite.Equal(models.PPMShipmentStatusCloseoutComplete, ppmShipment.Status)
		}
	})

	suite.Run("Update PPMShipment Status to CLOSEOUT_COMPLETE when there are rejected  moving expenses", func() {
		ppmShipment := factory.BuildPPMShipmentReadyForFinalCustomerCloseOut(nil, nil, nil)
		ppmShipment.Status = models.PPMShipmentStatusNeedsCloseout
		rejected := models.PPMDocumentStatusRejected
		movingExpense := factory.BuildMovingExpense(suite.DB(), []factory.Customization{
			{
				Model: models.MovingExpense{
					Status: &rejected,
				},
			},
		}, nil)
		progear := factory.BuildProgearWeightTicket(suite.DB(), nil, nil)
		weightTicket := factory.BuildWeightTicket(suite.DB(), nil, nil)

		ppmShipmentRouter := setUpPPMShipmentRouter(mtoShipmentRouterMethodToMock, nil)
		ppmShipment.WeightTickets = models.WeightTickets{weightTicket}
		ppmShipment.MovingExpenses = models.MovingExpenses{movingExpense}
		ppmShipment.ProgearWeightTickets = models.ProgearWeightTickets{progear}
		err := ppmShipmentRouter.SubmitReviewedDocuments(suite.AppContextForTest(), &ppmShipment)

		if suite.NoError(err) {
			suite.Equal(models.PPMShipmentStatusCloseoutComplete, ppmShipment.Status)
		}
	})

	suite.Run("Update PPMShipment Status to CLOSEOUT_COMPLETE when there are no rejected PPM Documents", func() {
		ppmShipment := factory.BuildPPMShipmentReadyForFinalCustomerCloseOut(nil, nil, nil)
		ppmShipment.Status = models.PPMShipmentStatusNeedsCloseout
		movingExpense := factory.BuildMovingExpense(suite.DB(), nil, nil)
		progear := factory.BuildProgearWeightTicket(suite.DB(), nil, nil)
		weightTicket := factory.BuildWeightTicket(suite.DB(), nil, nil)

		ppmShipmentRouter := setUpPPMShipmentRouter(mtoShipmentRouterMethodToMock, nil)
		ppmShipment.WeightTickets = models.WeightTickets{weightTicket}
		ppmShipment.MovingExpenses = models.MovingExpenses{movingExpense}
		ppmShipment.ProgearWeightTickets = models.ProgearWeightTickets{progear}
		err := ppmShipmentRouter.SubmitReviewedDocuments(suite.AppContextForTest(), &ppmShipment)

		if suite.NoError(err) {
			suite.Equal(models.PPMShipmentStatusCloseoutComplete, ppmShipment.Status)
		}
	})

	statusFailureTestCases := map[models.PPMShipmentStatus]func() models.PPMShipment{
		models.PPMShipmentStatusDraft: func() models.PPMShipment {
			return factory.BuildMinimalPPMShipment(nil, []factory.Customization{
				{
					Model: models.MTOShipment{
						Status: models.MTOShipmentStatusDraft,
					},
				},
			}, nil)
		},
		models.PPMShipmentStatusSubmitted: func() models.PPMShipment {
			return factory.BuildPPMShipment(nil, nil, nil)
		},
		models.PPMShipmentStatusWaitingOnCustomer: func() models.PPMShipment {
			return factory.BuildPPMShipment(nil, nil, []factory.Trait{factory.GetTraitApprovedPPMWithActualInfo})
		},
		models.PPMShipmentStatusCloseoutComplete: func() models.PPMShipment {
			return factory.BuildPPMShipment(nil, []factory.Customization{
				{
					Model: models.PPMShipment{
						Status: models.PPMShipmentStatusCloseoutComplete,
					},
				},
			}, []factory.Trait{factory.GetTraitApprovedPPMWithActualInfo})
		},
		models.PPMShipmentStatusNeedsAdvanceApproval: func() models.PPMShipment {
			return factory.BuildPPMShipment(nil, []factory.Customization{
				{
					Model: models.PPMShipment{
						Status: models.PPMShipmentStatusNeedsAdvanceApproval,
					},
				},
			}, []factory.Trait{factory.GetTraitApprovedPPMWithActualInfo})
		},
	}

	for currentStatus, makePPMShipment := range statusFailureTestCases {
		currentStatus := currentStatus
		makePPMShipment := makePPMShipment

		suite.Run(fmt.Sprintf("Can't set status to %s if it is currently %s", models.PPMShipmentStatusNeedsCloseout, currentStatus), func() {
			ppmShipment := makePPMShipment()

			ppmShipmentRouter := setUpPPMShipmentRouter(mtoShipmentRouterMethodToMock, nil)

			err := ppmShipmentRouter.SubmitReviewedDocuments(suite.AppContextForTest(), &ppmShipment)

			if suite.Error(err) {
				suite.IsType(apperror.ConflictError{}, err)
				suite.Contains(
					err.Error(),
					fmt.Sprintf(
						"PPM shipment documents cannot be submitted because it's not in the %s status.",
						models.PPMShipmentStatusNeedsCloseout,
					),
				)

				suite.Equal(currentStatus, ppmShipment.Status)
			}
		})
	}
}
