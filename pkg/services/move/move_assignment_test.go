package move

import (
	"github.com/go-openapi/strfmt"

	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/gen/ghcmessages"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/models/roles"
)

func (suite *MoveServiceSuite) TestBulkMoveAssignment() {
	moveAssigner := NewMoveAssignerBulkAssignment()

	setupTestData := func() (models.TransportationOffice, models.Move, models.Move, models.Move) {
		transportationOffice := factory.BuildTransportationOffice(suite.DB(), nil, nil)
		move1 := factory.BuildMoveWithShipment(suite.DB(), []factory.Customization{
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

		move2 := factory.BuildMoveWithShipment(suite.DB(), []factory.Customization{
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

		move3 := factory.BuildMoveWithShipment(suite.DB(), []factory.Customization{
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

		return transportationOffice, move1, move2, move3
	}

	suite.Run("successfully assigns multiple counseling moves to a SC user", func() {
		transportationOffice, move1, move2, move3 := setupTestData()

		officeUser := factory.BuildOfficeUserWithPrivileges(suite.DB(), []factory.Customization{
			{
				Model: models.OfficeUser{
					Email:  "officeuser1@example.com",
					Active: true,
				},
			},
			{
				Model:    transportationOffice,
				LinkOnly: true,
				Type:     &factory.TransportationOffices.CounselingOffice,
			},
			{
				Model: models.User{
					Privileges: []models.Privilege{
						{
							PrivilegeType: models.PrivilegeTypeSupervisor,
						},
					},
					Roles: []roles.Role{
						{
							RoleType: roles.RoleTypeServicesCounselor,
						},
					},
				},
			},
		}, nil)

		moves := []models.Move{move1, move2, move3}
		userData := []*ghcmessages.BulkAssignmentForUser{
			{ID: strfmt.UUID(officeUser.ID.String()), MoveAssignments: 2},
		}

		_, err := moveAssigner.BulkMoveAssignment(suite.AppContextForTest(), string(models.QueueTypeCounseling), userData, moves)
		suite.NoError(err)

		// reload move data to check assigned
		suite.NoError(suite.DB().Reload(&move1))
		suite.NoError(suite.DB().Reload(&move2))
		suite.NoError(suite.DB().Reload(&move3))

		suite.Equal(officeUser.ID, *move1.SCAssignedID)
		suite.Equal(officeUser.ID, *move2.SCAssignedID)
		suite.Nil(move3.SCAssignedID)
	})

	suite.Run("successfully assigns multiple closeout moves to a SC user", func() {
		transportationOffice, move1, move2, move3 := setupTestData()

		officeUser := factory.BuildOfficeUserWithPrivileges(suite.DB(), []factory.Customization{
			{
				Model: models.OfficeUser{
					Email:  "officeuser1@example.com",
					Active: true,
				},
			},
			{
				Model:    transportationOffice,
				LinkOnly: true,
				Type:     &factory.TransportationOffices.CounselingOffice,
			},
			{
				Model: models.User{
					Privileges: []models.Privilege{
						{
							PrivilegeType: models.PrivilegeTypeSupervisor,
						},
					},
					Roles: []roles.Role{
						{
							RoleType: roles.RoleTypeServicesCounselor,
						},
					},
				},
			},
		}, nil)

		moves := []models.Move{move1, move2, move3}
		userData := []*ghcmessages.BulkAssignmentForUser{
			{ID: strfmt.UUID(officeUser.ID.String()), MoveAssignments: 2},
		}

		_, err := moveAssigner.BulkMoveAssignment(suite.AppContextForTest(), string(models.QueueTypeCloseout), userData, moves)
		suite.NoError(err)

		// reload move data to check assigned
		suite.NoError(suite.DB().Reload(&move1))
		suite.NoError(suite.DB().Reload(&move2))
		suite.NoError(suite.DB().Reload(&move3))

		suite.Equal(officeUser.ID, *move1.SCAssignedID)
		suite.Equal(officeUser.ID, *move2.SCAssignedID)
		suite.Nil(move3.SCAssignedID)
	})

	suite.Run("successfully assigns multiple task order moves to a TOO user", func() {
		transportationOffice, move1, move2, move3 := setupTestData()

		officeUser := factory.BuildOfficeUserWithPrivileges(suite.DB(), []factory.Customization{
			{
				Model: models.OfficeUser{
					Email:  "officeuser1@example.com",
					Active: true,
				},
			},
			{
				Model:    transportationOffice,
				LinkOnly: true,
				Type:     &factory.TransportationOffices.CounselingOffice,
			},
			{
				Model: models.User{
					Privileges: []models.Privilege{
						{
							PrivilegeType: models.PrivilegeTypeSupervisor,
						},
					},
					Roles: []roles.Role{
						{
							RoleType: roles.RoleTypeTOO,
						},
					},
				},
			},
		}, nil)

		moves := []models.Move{move1, move2, move3}
		userData := []*ghcmessages.BulkAssignmentForUser{
			{ID: strfmt.UUID(officeUser.ID.String()), MoveAssignments: 2},
		}

		_, err := moveAssigner.BulkMoveAssignment(suite.AppContextForTest(), string(models.QueueTypeTaskOrder), userData, moves)
		suite.NoError(err)

		// reload move data to check assigned
		suite.NoError(suite.DB().Reload(&move1))
		suite.NoError(suite.DB().Reload(&move2))
		suite.NoError(suite.DB().Reload(&move3))

		suite.Equal(officeUser.ID, *move1.TOOAssignedID)
		suite.Equal(officeUser.ID, *move2.TOOAssignedID)
		suite.Nil(move3.TOOAssignedID)
	})

	suite.Run("successfully assigns payment requests to a TIO user", func() {
		transportationOffice, move1, move2, move3 := setupTestData()
		officeUser := factory.BuildOfficeUserWithPrivileges(suite.DB(), []factory.Customization{
			{
				Model: models.OfficeUser{
					Email:  "officeuser1@example.com",
					Active: true,
				},
			},
			{
				Model:    transportationOffice,
				LinkOnly: true,
				Type:     &factory.TransportationOffices.CounselingOffice,
			},
			{
				Model: models.User{
					Privileges: []models.Privilege{
						{
							PrivilegeType: models.PrivilegeTypeSupervisor,
						},
					},
					Roles: []roles.Role{
						{
							RoleType: roles.RoleTypeTIO,
						},
					},
				},
			},
		}, nil)

		moves := []models.Move{move1, move2, move3}
		userData := []*ghcmessages.BulkAssignmentForUser{
			{ID: strfmt.UUID(officeUser.ID.String()), MoveAssignments: 2},
		}

		_, err := moveAssigner.BulkMoveAssignment(suite.AppContextForTest(), string(models.QueueTypePaymentRequest), userData, moves)
		suite.NoError(err)

		// reload move data to check assigned
		suite.NoError(suite.DB().Reload(&move1))
		suite.NoError(suite.DB().Reload(&move2))

		suite.Equal(officeUser.ID, *move1.TIOAssignedID)
		suite.Equal(officeUser.ID, *move2.TIOAssignedID)
		suite.Nil(move3.TIOAssignedID)
	})
}
