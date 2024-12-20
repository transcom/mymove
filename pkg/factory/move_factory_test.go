package factory

import (
	"time"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/models/roles"
)

func (suite *FactorySuite) TestBuildMove() {
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
		locator := "ABC123"
		closeoutOfficeName := "Closeout office"

		customMove := models.Move{
			ReferenceID: &referenceID,
			Show:        &show,
			PPMType:     &ppmType,
			Locator:     locator,
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

		suite.Equal(ppmType, *move.PPMType)
		suite.False(*move.Show)
		suite.Equal(locator, move.Locator)
		suite.Equal(closeoutOfficeName, move.CloseoutOffice.Name)
		suite.Equal(referenceID, *move.ReferenceID)
		suite.NotNil(move.Contractor)
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
		suite.NotNil(move.ApprovedAt)
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
		suite.NotNil(move.ApprovedAt)
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
		suite.NotNil(move.AvailableToPrimeAt)
		suite.NotNil(move.ApprovedAt)
	})
	suite.Run("Successful creation of customized available to prime move", func() {
		// Under test:      BuildAvailableToPrimeMove
		// Set up:          Create a move with status needs SC
		// Expected outcome:Move with needs SC status
		availableToPrimeAt := time.Now().Add(-3 * 24 * time.Hour)
		customMove := models.Move{
			Locator:            "999111",
			AvailableToPrimeAt: &availableToPrimeAt,
			ApprovedAt:         &availableToPrimeAt,
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
		suite.Equal(availableToPrimeAt, *move.ApprovedAt) // Inheriting available to prime at date
		suite.NotNil(move.AvailableToPrimeAt)
		suite.NotNil(move.ApprovedAt)
	})
	suite.Run("Successful creation of a move with an assigned SC", func() {
		officeUser := BuildOfficeUserWithRoles(suite.DB(), nil, []roles.RoleType{roles.RoleTypeServicesCounselor})

		move := BuildMoveWithShipment(suite.DB(), []Customization{
			{
				Model: models.Move{
					Status: models.MoveStatusNeedsServiceCounseling,
				},
			},
			{
				Model:    officeUser,
				LinkOnly: true,
				Type:     &OfficeUsers.SCAssignedUser,
			},
		}, nil)
		suite.Equal(officeUser.ID, *move.SCAssignedID)
	})
	suite.Run("Successful creation of move with shipment", func() {
		// Under test:      BuildMoveWithShipment
		// Set up:          Create a move using BuildMoveWithShipment
		// Expected outcome:Move with shipment

		move := BuildMoveWithShipment(suite.DB(), nil, nil)
		suite.NotEmpty(move.MTOShipments)
		suite.Equal(models.MTOShipmentStatusSubmitted, move.MTOShipments[0].Status)
	})
	suite.Run("Successful creation of a move with an assigned TIO", func() {
		officeUser := BuildOfficeUserWithRoles(suite.DB(), nil, []roles.RoleType{roles.RoleTypeTIO})

		move := BuildMoveWithShipment(suite.DB(), []Customization{
			{
				Model: models.Move{
					Status: models.MoveStatusAPPROVED,
				},
			},
			{
				Model:    officeUser,
				LinkOnly: true,
				Type:     &OfficeUsers.TIOAssignedUser,
			},
		}, nil)
		suite.Equal(officeUser.ID, *move.TIOAssignedID)
	})

	suite.Run("Successful creation of customized move with shipment", func() {
		// Under test:      BuildMoveWithShipment
		// Set up:          Create a custom move using BuildMoveWithShipment
		// Expected outcome:Customized move with shipment
		customMove := models.Move{
			Locator: "999111",
			Status:  models.MoveStatusAPPROVALSREQUESTED,
		}
		customServiceMember := models.ServiceMember{
			FirstName: models.StringPointer("Riley"),
		}
		customOrders := models.Order{
			HasDependents: true,
		}
		customShipment := models.MTOShipment{
			Status: models.MTOShipmentStatusCanceled,
		}
		// Create move
		move := BuildMoveWithShipment(suite.DB(), []Customization{
			{Model: customMove},
			{Model: customServiceMember},
			{Model: customOrders},
			{Model: customShipment},
		}, nil)
		suite.Equal(customMove.Locator, move.Locator)
		suite.Equal(customMove.Status, move.Status)
		suite.Equal(customServiceMember.FirstName, move.Orders.ServiceMember.FirstName)
		suite.Equal(customOrders.HasDependents, move.Orders.HasDependents)
		suite.Equal(customShipment.Status, move.MTOShipments[0].Status)
	})

}
