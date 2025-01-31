package move

import (
	"time"

	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/models/roles"
	"github.com/transcom/mymove/pkg/services"
)

func (suite *MoveServiceSuite) TestMoveFetcher() {
	moveFetcher := NewMoveFetcher()
	defaultSearchParams := services.MoveFetcherParams{}

	suite.Run("successfully returns default draft move", func() {
		expectedMove := factory.BuildMove(suite.DB(), nil, nil)

		actualMove, err := moveFetcher.FetchMove(suite.AppContextForTest(), expectedMove.Locator, &defaultSearchParams)
		suite.FatalNoError(err)

		suite.Equal(expectedMove.ID, actualMove.ID)
		suite.Equal(expectedMove.Locator, actualMove.Locator)
		suite.Equal(expectedMove.CreatedAt.Format(time.RFC3339), actualMove.CreatedAt.Format(time.RFC3339))
		suite.Equal(expectedMove.UpdatedAt.Format(time.RFC3339), actualMove.UpdatedAt.Format(time.RFC3339))
		suite.Equal(expectedMove.SubmittedAt, actualMove.SubmittedAt)
		suite.Equal(expectedMove.OrdersID, actualMove.OrdersID)
		suite.Equal(expectedMove.Status, actualMove.Status)
		suite.Equal(expectedMove.AvailableToPrimeAt, actualMove.AvailableToPrimeAt)
		suite.Equal(expectedMove.ApprovedAt, actualMove.ApprovedAt)
		suite.Equal(expectedMove.ContractorID, actualMove.ContractorID)
		suite.Equal(expectedMove.Contractor.ContractNumber, actualMove.Contractor.ContractNumber)
		suite.Equal(expectedMove.ReferenceID, actualMove.ReferenceID)
	})

	suite.Run("successfully returns submitted move available to prime", func() {
		expectedMove := factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil)

		actualMove, err := moveFetcher.FetchMove(suite.AppContextForTest(), expectedMove.Locator, &defaultSearchParams)
		suite.FatalNoError(err)

		suite.Equal(expectedMove.ID, actualMove.ID)
		suite.Equal(expectedMove.Locator, actualMove.Locator)
		suite.Equal(expectedMove.CreatedAt.Format(time.RFC3339), actualMove.CreatedAt.Format(time.RFC3339))
		suite.Equal(expectedMove.UpdatedAt.Format(time.RFC3339), actualMove.UpdatedAt.Format(time.RFC3339))
		suite.Equal(expectedMove.SubmittedAt, actualMove.SubmittedAt)
		suite.Equal(expectedMove.OrdersID, actualMove.OrdersID)
		suite.Equal(expectedMove.Status, actualMove.Status)
		suite.Equal(expectedMove.AvailableToPrimeAt.Format(time.RFC3339), actualMove.AvailableToPrimeAt.Format(time.RFC3339))
		suite.Equal(expectedMove.ApprovedAt.Format(time.RFC3339), actualMove.ApprovedAt.Format(time.RFC3339))
		suite.Equal(expectedMove.ContractorID, actualMove.ContractorID)
		suite.Equal(expectedMove.Contractor.Name, actualMove.Contractor.Name)
		suite.Equal(expectedMove.ReferenceID, actualMove.ReferenceID)
	})

	suite.Run("returns not found error for unknown locator", func() {
		_ = factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil)

		_, err := moveFetcher.FetchMove(suite.AppContextForTest(), "QX97UY", &defaultSearchParams)
		suite.Error(err)
		suite.IsType(apperror.NotFoundError{}, err)
	})

	suite.Run("Returns not found for a move that is marked hidden in the db", func() {
		hide := false
		hiddenMove := factory.BuildMove(suite.DB(), []factory.Customization{
			{
				Model: models.Move{
					Show: &hide,
				},
			},
		}, nil)
		locator := hiddenMove.Locator
		searchParams := services.MoveFetcherParams{
			IncludeHidden: false,
		}

		_, err := moveFetcher.FetchMove(suite.AppContextForTest(), locator, &searchParams)

		suite.Error(err)
		suite.IsType(apperror.NotFoundError{}, err)
	})

	suite.Run("Returns hidden move if explicit param is passed in", func() {
		hide := false
		actualMove := factory.BuildMove(suite.DB(), []factory.Customization{
			{
				Model: models.Move{
					Show: &hide,
				},
			},
		}, nil)
		locator := actualMove.Locator
		searchParams := services.MoveFetcherParams{
			IncludeHidden: true,
		}

		expectedMove, err := moveFetcher.FetchMove(suite.AppContextForTest(), locator, &searchParams)

		suite.FatalNoError(err)
		suite.Equal(expectedMove.ID, actualMove.ID)
		suite.Equal(expectedMove.Locator, actualMove.Locator)
	})
}

func (suite *MoveServiceSuite) TestMoveFetcherBulkAssignment() {
	setupTestData := func() (services.MoveFetcherBulkAssignment, models.Move, models.TransportationOffice, models.OfficeUser) {
		moveFetcher := NewMoveFetcherBulkAssignment()
		transportationOffice := factory.BuildTransportationOffice(suite.DB(), nil, nil)

		// this move has a transportation office associated with it that matches
		// the SC's transportation office and should be found
		move := factory.BuildMoveWithShipment(suite.DB(), []factory.Customization{
			{
				Model: models.Move{
					Status: models.MoveStatusNeedsServiceCounseling,
				},
			},
			{
				Model:    transportationOffice,
				LinkOnly: true,
				Type:     &factory.TransportationOffices.CounselingOffice,
			},
		}, nil)

		factory.BuildMoveWithShipment(suite.DB(), []factory.Customization{
			{
				Model: models.Move{
					Status: models.MoveStatusNeedsServiceCounseling,
				},
			},
			{
				Model:    transportationOffice,
				LinkOnly: true,
				Type:     &factory.TransportationOffices.CounselingOffice,
			},
		}, nil)

		officeUser := factory.BuildOfficeUserWithRoles(suite.DB(), []factory.Customization{
			{
				Model:    transportationOffice,
				LinkOnly: true,
				Type:     &factory.TransportationOffices.CounselingOffice,
			},
		}, []roles.RoleType{roles.RoleTypeServicesCounselor})

		return moveFetcher, move, transportationOffice, officeUser
	}

	suite.Run("Returns moves that fulfill the query's criteria", func() {
		moveFetcher, _, _, officeUser := setupTestData()
		moves, err := moveFetcher.FetchMovesForBulkAssignmentCounseling(suite.AppContextForTest(), "KKFA", officeUser.TransportationOffice.ID)
		suite.FatalNoError(err)
		suite.Equal(2, len(moves))
	})

	suite.Run("Does not return moves that are counseled by a different counseling office", func() {
		moveFetcher, _, _, officeUser := setupTestData()
		transportationOffice := factory.BuildTransportationOffice(suite.DB(), nil, nil)

		factory.BuildMoveWithShipment(suite.DB(), []factory.Customization{
			{
				Model: models.Move{
					Status: models.MoveStatusNeedsServiceCounseling,
				},
			},
			{
				Model:    transportationOffice,
				LinkOnly: true,
				Type:     &factory.TransportationOffices.CounselingOffice,
			},
		}, nil)

		moves, err := moveFetcher.FetchMovesForBulkAssignmentCounseling(suite.AppContextForTest(), "KKFA", officeUser.TransportationOffice.ID)
		suite.FatalNoError(err)
		suite.Equal(2, len(moves))
	})

	suite.Run("Does not return moves with safety, bluebark, or wounded warrior order types", func() {
		moveFetcher, _, transportationOffice, officeUser := setupTestData()
		factory.BuildMoveWithShipment(suite.DB(), []factory.Customization{
			{
				Model: models.Order{
					OrdersType: internalmessages.OrdersTypeSAFETY,
				},
			},
			{
				Model: models.Move{
					Status: models.MoveStatusNeedsServiceCounseling,
				},
			},
			{
				Model:    transportationOffice,
				LinkOnly: true,
				Type:     &factory.TransportationOffices.CounselingOffice,
			},
		}, nil)
		factory.BuildMoveWithShipment(suite.DB(), []factory.Customization{
			{
				Model: models.Order{
					OrdersType: internalmessages.OrdersTypeBLUEBARK,
				},
			},
			{
				Model: models.Move{
					Status: models.MoveStatusNeedsServiceCounseling,
				},
			},
			{
				Model:    transportationOffice,
				LinkOnly: true,
				Type:     &factory.TransportationOffices.CounselingOffice,
			},
		}, nil)
		factory.BuildMoveWithShipment(suite.DB(), []factory.Customization{
			{
				Model: models.Order{
					OrdersType: internalmessages.OrdersTypeWOUNDEDWARRIOR,
				},
			},
			{
				Model: models.Move{
					Status: models.MoveStatusNeedsServiceCounseling,
				},
			},
			{
				Model:    transportationOffice,
				LinkOnly: true,
				Type:     &factory.TransportationOffices.CounselingOffice,
			},
		}, nil)

		moves, err := moveFetcher.FetchMovesForBulkAssignmentCounseling(suite.AppContextForTest(), "KKFA", officeUser.TransportationOffice.ID)
		suite.FatalNoError(err)
		suite.Equal(2, len(moves))
	})

	// BuildMoveWithPPMShipment apparently builds 3 moves each time its run, so the best way
	// to test is to make sure that the moveWithPPM move is not returned in these 3 separate tests
	suite.Run("Does not return moves with PPMs in waiting on customer status", func() {
		moveFetcher := NewMoveFetcherBulkAssignment()
		transportationOffice := factory.BuildTransportationOffice(suite.DB(), nil, nil)
		moveWithWaitingOnCustomerPPM := factory.BuildMoveWithPPMShipment(suite.DB(), []factory.Customization{
			{
				Model:    transportationOffice,
				LinkOnly: true,
				Type:     &factory.TransportationOffices.CounselingOffice,
			},
			{
				Model: models.PPMShipment{
					Status: models.PPMShipmentStatusWaitingOnCustomer,
				},
			},
		}, []factory.Trait{factory.GetTraitNeedsServiceCounselingMove})

		officeUser := factory.BuildOfficeUserWithRoles(suite.DB(), []factory.Customization{
			{
				Model:    transportationOffice,
				LinkOnly: true,
				Type:     &factory.TransportationOffices.CounselingOffice,
			},
		}, []roles.RoleType{roles.RoleTypeServicesCounselor})

		moves, err := moveFetcher.FetchMovesForBulkAssignmentCounseling(suite.AppContextForTest(), "KKFA", officeUser.TransportationOffice.ID)
		suite.FatalNoError(err)
		// confirm that the there is only one move appearing
		suite.Equal(1, len(moves))
		// confirm that the move appearing iS NOT the moveWithPPM
		suite.NotEqual(moves[0].ID, moveWithWaitingOnCustomerPPM.ID)
		// confirm that the rest of the details are correct
		// and that it SHOULD show up in the queue if it wasn't for PPM status
		// move is NEEDS SERVICE COUNSELING STATUS
		suite.Equal(moveWithWaitingOnCustomerPPM.Status, models.MoveStatusNeedsServiceCounseling)
		// move is not assigned to anyone
		suite.Nil(moveWithWaitingOnCustomerPPM.SCAssignedID)
		// GBLOC is the same
		suite.Equal(*moveWithWaitingOnCustomerPPM.Orders.OriginDutyLocationGBLOC, officeUser.TransportationOffice.Gbloc)
		// Show is true
		suite.Equal(moveWithWaitingOnCustomerPPM.Show, models.BoolPointer(true))
		// Move is counseled by the office user's office
		suite.Equal(*moveWithWaitingOnCustomerPPM.CounselingOfficeID, officeUser.TransportationOfficeID)
		// Orders type isn't WW, BB, or Safety
		suite.Equal(moveWithWaitingOnCustomerPPM.Orders.OrdersType, internalmessages.OrdersTypePERMANENTCHANGEOFSTATION)
	})

	suite.Run("Does not return moves with PPMs in needs closeout status", func() {
		moveFetcher := NewMoveFetcherBulkAssignment()
		transportationOffice := factory.BuildTransportationOffice(suite.DB(), nil, nil)
		moveWithNeedsCloseoutPPM := factory.BuildMoveWithPPMShipment(suite.DB(), []factory.Customization{
			{
				Model:    transportationOffice,
				LinkOnly: true,
				Type:     &factory.TransportationOffices.CounselingOffice,
			},
			{
				Model: models.PPMShipment{
					Status: models.PPMShipmentStatusNeedsCloseout,
				},
			},
		}, []factory.Trait{factory.GetTraitNeedsServiceCounselingMove})

		officeUser := factory.BuildOfficeUserWithRoles(suite.DB(), []factory.Customization{
			{
				Model:    transportationOffice,
				LinkOnly: true,
				Type:     &factory.TransportationOffices.CounselingOffice,
			},
		}, []roles.RoleType{roles.RoleTypeServicesCounselor})

		moves, err := moveFetcher.FetchMovesForBulkAssignmentCounseling(suite.AppContextForTest(), "KKFA", officeUser.TransportationOffice.ID)
		suite.FatalNoError(err)
		// confirm that the there is only one move appearing
		suite.Equal(1, len(moves))
		// confirm that the move appearing iS NOT the moveWithPPM
		suite.NotEqual(moves[0].ID, moveWithNeedsCloseoutPPM.ID)
		// confirm that the rest of the details are correct
		// and that it SHOULD show up in the queue if it wasn't for PPM status
		// move is NEEDS SERVICE COUNSELING STATUS
		suite.Equal(moveWithNeedsCloseoutPPM.Status, models.MoveStatusNeedsServiceCounseling)
		// move is not assigned to anyone
		suite.Nil(moveWithNeedsCloseoutPPM.SCAssignedID)
		// GBLOC is the same
		suite.Equal(*moveWithNeedsCloseoutPPM.Orders.OriginDutyLocationGBLOC, officeUser.TransportationOffice.Gbloc)
		// Show is true
		suite.Equal(moveWithNeedsCloseoutPPM.Show, models.BoolPointer(true))
		// Move is counseled by the office user's office
		suite.Equal(*moveWithNeedsCloseoutPPM.CounselingOfficeID, officeUser.TransportationOfficeID)
		// Orders type isn't WW, BB, or Safety
		suite.Equal(moveWithNeedsCloseoutPPM.Orders.OrdersType, internalmessages.OrdersTypePERMANENTCHANGEOFSTATION)
	})
	suite.Run("Does not return moves with PPMs in closeout complete status", func() {
		moveFetcher := NewMoveFetcherBulkAssignment()
		transportationOffice := factory.BuildTransportationOffice(suite.DB(), []factory.Customization{
			{
				Model: models.TransportationOffice{
					ProvidesCloseout: true,
				},
			},
		}, nil)
		moveWithCloseoutCompletePPM := factory.BuildMoveWithPPMShipment(suite.DB(), []factory.Customization{
			{
				Model:    transportationOffice,
				LinkOnly: true,
				Type:     &factory.TransportationOffices.CounselingOffice,
			},
			{
				Model: models.PPMShipment{
					Status: models.PPMShipmentStatusCloseoutComplete,
				},
			},
		}, []factory.Trait{factory.GetTraitNeedsServiceCounselingMove})

		officeUser := factory.BuildOfficeUserWithRoles(suite.DB(), []factory.Customization{
			{
				Model:    transportationOffice,
				LinkOnly: true,
				Type:     &factory.TransportationOffices.CounselingOffice,
			},
		}, []roles.RoleType{roles.RoleTypeServicesCounselor})

		moves, err := moveFetcher.FetchMovesForBulkAssignmentCounseling(suite.AppContextForTest(), "KKFA", officeUser.TransportationOffice.ID)
		suite.FatalNoError(err)
		// confirm that the there is only one move appearing
		suite.Equal(1, len(moves))
		// confirm that the move appearing iS NOT the moveWithPPM
		suite.NotEqual(moves[0].ID, moveWithCloseoutCompletePPM.ID)
		// confirm that the rest of the details are correct
		// and that it SHOULD show up in the queue if it wasn't for PPM status
		// move is NEEDS SERVICE COUNSELING STATUS
		suite.Equal(moveWithCloseoutCompletePPM.Status, models.MoveStatusNeedsServiceCounseling)
		// move is not assigned to anyone
		suite.Nil(moveWithCloseoutCompletePPM.SCAssignedID)
		// GBLOC is the same
		suite.Equal(*moveWithCloseoutCompletePPM.Orders.OriginDutyLocationGBLOC, officeUser.TransportationOffice.Gbloc)
		// Show is true
		suite.Equal(moveWithCloseoutCompletePPM.Show, models.BoolPointer(true))
		// Move is counseled by the office user's office
		suite.Equal(*moveWithCloseoutCompletePPM.CounselingOfficeID, officeUser.TransportationOfficeID)
		// Orders type isn't WW, BB, or Safety
		suite.Equal(moveWithCloseoutCompletePPM.Orders.OrdersType, internalmessages.OrdersTypePERMANENTCHANGEOFSTATION)
	})

	suite.Run("Does not return moves that are already assigned", func() {
		// moveFetcher, _, transOffice, officeUser := setupTestData()
		moveFetcher := NewMoveFetcherBulkAssignment()
		transportationOffice := factory.BuildTransportationOffice(suite.DB(), nil, nil)

		officeUser := factory.BuildOfficeUserWithRoles(suite.DB(), []factory.Customization{
			{
				Model:    transportationOffice,
				LinkOnly: true,
				Type:     &factory.TransportationOffices.CounselingOffice,
			},
		}, []roles.RoleType{roles.RoleTypeServicesCounselor})

		assignedMove := factory.BuildMoveWithShipment(suite.DB(), []factory.Customization{
			{
				Model: models.Move{
					Status: models.MoveStatusNeedsServiceCounseling,
				},
			},
			{
				Model:    transportationOffice,
				LinkOnly: true,
				Type:     &factory.TransportationOffices.CounselingOffice,
			},
			{
				Model:    officeUser,
				LinkOnly: true,
				Type:     &factory.OfficeUsers.SCAssignedUser,
			},
		}, nil)

		moves, err := moveFetcher.FetchMovesForBulkAssignmentCounseling(suite.AppContextForTest(), "KKFA", officeUser.TransportationOffice.ID)
		suite.FatalNoError(err)

		// confirm that the assigned move isn't returned
		for _, move := range moves {
			suite.NotEqual(move.ID, assignedMove.ID)
		}

		// confirm that the rest of the details are correct
		// move is NEEDS SERVICE COUNSELING STATUS
		suite.Equal(assignedMove.Status, models.MoveStatusNeedsServiceCounseling)
		// GBLOC is the same
		suite.Equal(*assignedMove.Orders.OriginDutyLocationGBLOC, officeUser.TransportationOffice.Gbloc)
		// Show is true
		suite.Equal(assignedMove.Show, models.BoolPointer(true))
		// Move is counseled by the office user's office
		suite.Equal(*assignedMove.CounselingOfficeID, officeUser.TransportationOfficeID)
		// Orders type isn't WW, BB, or Safety
		suite.Equal(assignedMove.Orders.OrdersType, internalmessages.OrdersTypePERMANENTCHANGEOFSTATION)
	})

  suite.Run("Closeout returns non Navy/USCG/USMC ppms in needs closeout status", func() {
		moveFetcher := NewMoveFetcherBulkAssignment()
		transportationOffice := factory.BuildTransportationOffice(suite.DB(), nil, nil)
		officeUser := factory.BuildOfficeUserWithRoles(suite.DB(), []factory.Customization{
			{
				Model:    transportationOffice,
				LinkOnly: true,
				Type:     &factory.TransportationOffices.CloseoutOffice,
			},
		}, []roles.RoleType{roles.RoleTypeServicesCounselor})

		submittedAt := time.Now()

		// create non USMC/USCG/NAVY ppm in need closeout status
		factory.BuildMoveWithPPMShipment(suite.DB(), []factory.Customization{
			{
				Model:    transportationOffice,
				LinkOnly: true,
				Type:     &factory.TransportationOffices.CloseoutOffice,
			},
			{
				Model: models.PPMShipment{
					Status:      models.PPMShipmentStatusNeedsCloseout,
					SubmittedAt: &submittedAt,
				},
			},
			{
				Model: models.Move{
					Status: models.MoveStatusAPPROVED,
				},
			},
		}, nil)

		// create non closeout needed ppm
		factory.BuildMoveWithPPMShipment(suite.DB(), []factory.Customization{
			{
				Model:    transportationOffice,
				LinkOnly: true,
				Type:     &factory.TransportationOffices.CloseoutOffice,
			},
			{
				Model: models.PPMShipment{
					Status:      models.PPMShipmentStatusWaitingOnCustomer,
					SubmittedAt: &submittedAt,
				},
			},
			{
				Model: models.Move{
					Status: models.MoveStatusAPPROVED,
				},
			},
		}, nil)

		marine := models.AffiliationMARINES
		marinePPM := factory.BuildPPMShipment(suite.DB(), []factory.Customization{
			{
				Model: models.Move{
					Status: models.MoveStatusAPPROVED,
				},
			},
			{
				Model: models.MTOShipment{
					ShipmentType: models.MTOShipmentTypePPM,
				},
			},
			{
				Model: models.PPMShipment{
					Status:      models.PPMShipmentStatusNeedsCloseout,
					SubmittedAt: &submittedAt,
				},
			},
			{
				Model: models.ServiceMember{
					Affiliation: &marine,
				},
			},
		}, nil)

		moves, err := moveFetcher.FetchMovesForBulkAssignmentCloseout(suite.AppContextForTest(), "KKFA", officeUser.TransportationOffice.ID)
		suite.FatalNoError(err)
		suite.Equal(1, len(moves))
		suite.NotEqual(marinePPM.ID, moves[0].ID)
	})
  
  suite.Run("TOO: Returns moves that fulfill the query criteria", func() {
		moveFetcher := NewMoveFetcherBulkAssignment()
		transportationOffice := factory.BuildTransportationOffice(suite.DB(), nil, nil)
		officeUser := factory.BuildOfficeUserWithRoles(suite.DB(), []factory.Customization{
			{
				Model:    transportationOffice,
				LinkOnly: true,
				Type:     &factory.TransportationOffices.CounselingOffice,
			},
		}, []roles.RoleType{roles.RoleTypeTOO})
		factory.BuildMoveWithShipment(suite.DB(), []factory.Customization{
			{
				Model: models.Move{
					Status: models.MoveStatusAPPROVALSREQUESTED,
				},
			},
			{
				Model:    transportationOffice,
				LinkOnly: true,
				Type:     &factory.TransportationOffices.CounselingOffice,
			},
		}, nil)
		factory.BuildMoveWithShipment(suite.DB(), []factory.Customization{
			{
				Model: models.Move{
					Status: models.MoveStatusServiceCounselingCompleted,
				},
			},
			{
				Model:    transportationOffice,
				LinkOnly: true,
				Type:     &factory.TransportationOffices.CounselingOffice,
			},
		}, nil)
		marine := models.AffiliationMARINES
		factory.BuildMoveWithShipment(suite.DB(), []factory.Customization{
			{
				Model: models.Move{
					Status: models.MoveStatusServiceCounselingCompleted,
				},
			},
			{
				Model:    transportationOffice,
				LinkOnly: true,
				Type:     &factory.TransportationOffices.CounselingOffice,
			},
			{
				Model: models.ServiceMember{
					Affiliation: &marine,
				},
			},
		}, nil)
		moves, err := moveFetcher.FetchMovesForBulkAssignmentTaskOrder(suite.AppContextForTest(), "KKFA", officeUser.TransportationOffice.ID)
		suite.FatalNoError(err)
		suite.Equal(2, len(moves))
	})
}
