package factory

import (
	"time"

	"github.com/transcom/mymove/pkg/models"
)

func (suite *FactorySuite) TestBuildMove() {
	defaultMoveType := models.SelectedMoveTypePPM
	partialType := "PARTIAL"
	defaultPpmType := &partialType
	defaultShow := true

	suite.Run("Successful creation of default move", func() {
		// Under test:      BuildMove
		// Set up:          Create a default move
		// Expected outcome:Create a contractor, order and move

		// SETUP
		// Create a default move info to compare values

		// Create move
		move := BuildMove(suite.DB(), nil, nil)

		suite.Equal(defaultMoveType, *move.SelectedMoveType)
		suite.Equal(*defaultPpmType, *move.PPMType)
		suite.Equal(defaultShow, *move.Show)
		suite.NotNil(move.Contractor)
		suite.False(move.ContractorID.IsNil())
		suite.NotNil(move.ReferenceID)
		suite.NotEmpty(*move.ReferenceID)
	})
	suite.Run("Success creation of stubbed move ", func() {
		// Under test:      BuildMove
		// Set up:          Create a move, but don't pass in a db
		// Expected outcome:Move should be created
		//                  No move should be created in database
		precount, err := suite.DB().Count(&models.Move{})
		suite.NoError(err)

		move := BuildMove(nil, nil, nil)
		suite.Equal(defaultMoveType, *move.SelectedMoveType)
		suite.Equal(*defaultPpmType, *move.PPMType)
		suite.Equal(defaultShow, *move.Show)
		suite.NotNil(move.Contractor)
		suite.Empty(*move.ReferenceID)

		count, err := suite.DB().Count(&models.Move{})
		suite.NoError(err)
		suite.Equal(precount, count)
	})
	suite.Run("Successful creation of custom move", func() {
		// Under test:      BuildMove
		// Set up:          Create a custom move
		// Expected outcome:Create a contractor, order and move

		// SETUP
		// Create default move info to compare values

		// custom move
		referenceID := "refID"
		show := false
		ppmType := "FULL"
		moveType := models.SelectedMoveTypeHHG
		locator := "ABC123"
		closeoutOfficeName := "Closeout office"

		customMove := models.Move{
			ReferenceID:      &referenceID,
			Show:             &show,
			SelectedMoveType: &moveType,
			PPMType:          &ppmType,
			Locator:          locator,
		}
		customs := []Customization{
			{
				Model: customMove,
			},
			{
				Model: models.TransportationOffice{
					Name: closeoutOfficeName,
				},
				Type: &TransportationOffices.CloseoutOffice,
			},
		}
		move := BuildMove(suite.DB(), customs, nil)

		suite.Equal(moveType, *move.SelectedMoveType)
		suite.Equal(ppmType, *move.PPMType)
		suite.False(*move.Show)
		suite.Equal(locator, move.Locator)
		suite.Equal(closeoutOfficeName, move.CloseoutOffice.Name)
		suite.Equal(referenceID, *move.ReferenceID)
		suite.NotNil(move.Contractor)
	})
	suite.Run("Successful creation of move without move type", func() {
		// Under test:      BuildMoveWithoutMoveType
		// Set up:          Create a move without move type set
		// Expected outcome:Create a contractor, order and move

		// Create move
		move := BuildMoveWithoutMoveType(suite.DB(), nil, nil)

		suite.Nil(move.SelectedMoveType)
		suite.Nil(move.PPMType)
		suite.NotNil(move.Contractor)
		suite.False(move.ContractorID.IsNil())
		suite.NotNil(move.ReferenceID)
		suite.NotEmpty(*move.ReferenceID)
	})
	suite.Run("Successful creation of stubbed move with status", func() {
		// Under test:      BuildStubbedMoveWithStatus
		// Set up:          Create a stubbed move with given status
		// Expected outcome:Create a stubbed entitlement, duty
		// location, contractor, order and move

		// Create move
		status := models.MoveStatusCANCELED
		move := BuildStubbedMoveWithStatus(status)

		suite.False(move.ID.IsNil())
		suite.NotEmpty(move.Locator)
		suite.Equal(status, move.Status)

		suite.False(move.OrdersID.IsNil())
		suite.False(move.Orders.ID.IsNil())
		suite.NotNil(move.Orders.Grade)

		suite.False(move.Orders.OriginDutyLocationID.IsNil())
		suite.False(move.Orders.OriginDutyLocation.ID.IsNil())
		suite.NotEmpty(move.Orders.OriginDutyLocation.Name)
		suite.False(move.Orders.EntitlementID.IsNil())
		suite.False(move.Orders.Entitlement.ID.IsNil())

		suite.NotNil(move.Contractor)
		suite.NotEmpty(move.Contractor.Name)
	})
	suite.Run("Successful creation of move with approval status", func() {
		// Under test:      BuildApprovalsRequestedMove
		// Set up:          Create a move with approvals requested
		// Expected outcome:Move with available to prime and approvals
		// requested status

		// Create move
		move := BuildApprovalsRequestedMove(suite.DB(), nil, nil)
		suite.Equal(models.MoveStatusAPPROVALSREQUESTED, move.Status)
		suite.NotNil(move.AvailableToPrimeAt)
	})
	suite.Run("Successful creation of customized move with approval status", func() {
		// Under test:      BuildApprovalsRequestedMove
		// Set up:          Create a move with approvals requested
		// Expected outcome:Move with available to prime and approvals
		// requested status

		customMove := models.Move{
			Locator: "999111",
		}
		// Create move
		move := BuildApprovalsRequestedMove(suite.DB(), []Customization{
			{
				Model: customMove,
			},
		}, nil)
		suite.Equal(customMove.Locator, move.Locator)
		suite.Equal(models.MoveStatusAPPROVALSREQUESTED, move.Status)
		suite.NotNil(move.AvailableToPrimeAt)
	})
	suite.Run("Successful creation of submitted move", func() {
		// Under test:      BuildSubmittedMove
		// Set up:          Create a move with status submitted
		// Expected outcome:Move with submitted status

		move := BuildSubmittedMove(suite.DB(), nil, nil)
		suite.Equal(models.MoveStatusSUBMITTED, move.Status)
		suite.NotNil(move.SubmittedAt)
	})
	suite.Run("Successful creation of customized move with submitted status", func() {
		// Under test:      BuildSubmittedMove
		// Set up:          Create a move with status submitted
		// Expected outcome:Move with submitted status

		customMove := models.Move{
			Locator: "999111",
		}
		// Create move
		move := BuildSubmittedMove(suite.DB(), []Customization{
			{
				Model: customMove,
			},
		}, nil)
		suite.Equal(customMove.Locator, move.Locator)
		suite.Equal(models.MoveStatusSUBMITTED, move.Status)
		suite.NotNil(move.SubmittedAt)
	})
	suite.Run("Successful creation of needs SC move", func() {
		// Under test:      BuildNeedsServiceCounselingMove
		// Set up:          Create a move with status needs SC
		// Expected outcome:Move with needs SC status

		move := BuildNeedsServiceCounselingMove(suite.DB(), nil, nil)
		suite.Equal(models.MoveStatusNeedsServiceCounseling, move.Status)
	})
	suite.Run("Successful creation of customized move with needs SC status", func() {
		// Under test:      BuildNeedsServiceCounselingMove
		// Set up:          Create a move with status needs SC
		// Expected outcome:Move with needs SC status

		customMove := models.Move{
			Locator: "999111",
		}
		// Create move
		move := BuildNeedsServiceCounselingMove(suite.DB(), []Customization{
			{
				Model: customMove,
			},
		}, nil)
		suite.Equal(customMove.Locator, move.Locator)
		suite.Equal(models.MoveStatusNeedsServiceCounseling, move.Status)
		suite.NotNil(move.SubmittedAt)
	})
	suite.Run("Successful creation of SC completed move", func() {
		// Under test:      BuildServiceCounselingCompletedMove
		// Set up:          Create a move with status SC completed
		// Expected outcome:Move with SC completd status

		move := BuildServiceCounselingCompletedMove(suite.DB(), nil, nil)
		suite.Equal(models.MoveStatusServiceCounselingCompleted, move.Status)
		suite.NotNil(move.ServiceCounselingCompletedAt)
	})
	suite.Run("Successful creation of customized SC completed move", func() {
		// Under test:      BuildServiceCounselingCompletedMove
		// Set up:          Create a move with status SC completed
		// Expected outcome:Move with SC completd status

		customMove := models.Move{
			Locator: "999111",
		}
		move := BuildServiceCounselingCompletedMove(suite.DB(), []Customization{
			{
				Model: customMove,
			},
		}, nil)
		suite.Equal(customMove.Locator, move.Locator)
		suite.Equal(models.MoveStatusServiceCounselingCompleted, move.Status)
		suite.NotNil(move.ServiceCounselingCompletedAt)
	})
	suite.Run("Successful creation of available to prime move", func() {
		// Under test:      BuildAvailableToPrimeMove
		// Set up:          Create a move with status SC completed
		// Expected outcome:Move with SC completd status

		move := BuildAvailableToPrimeMove(suite.DB(), nil, nil)
		suite.Equal(models.MoveStatusAPPROVED, move.Status)
		suite.NotNil(move.SubmittedAt)
		suite.NotNil(move.AvailableToPrimeAt)
	})
	suite.Run("Successful creation of customized available to prime move", func() {
		// Under test:      BuildAvailableToPrimeMove
		// Set up:          Create a move with status needs SC
		// Expected outcome:Move with needs SC status
		availableToPrimeAt := time.Now().Add(-3 * 24 * time.Hour)
		customMove := models.Move{
			Locator:            "999111",
			AvailableToPrimeAt: &availableToPrimeAt,
		}
		// Create move
		move := BuildAvailableToPrimeMove(suite.DB(), []Customization{
			{
				Model: customMove,
			},
		}, nil)
		suite.Equal(customMove.Locator, move.Locator)
		suite.Equal(models.MoveStatusAPPROVED, move.Status)
		suite.Equal(availableToPrimeAt, *move.AvailableToPrimeAt)
		suite.NotNil(move.AvailableToPrimeAt)
	})

}
