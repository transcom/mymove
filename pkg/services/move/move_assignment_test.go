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

	suite.Run("properly distributes moves", func() {
		transportationOffice, move1, move2, move3 := setupTestData()
		move4 := factory.BuildMoveWithShipment(suite.DB(), []factory.Customization{
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
		move5 := factory.BuildMoveWithShipment(suite.DB(), []factory.Customization{
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
		move6 := factory.BuildMoveWithShipment(suite.DB(), []factory.Customization{
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

		officeUser1 := factory.BuildOfficeUserWithPrivileges(suite.DB(), []factory.Customization{
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
					Privileges: []roles.Privilege{
						{
							PrivilegeType: roles.PrivilegeTypeSupervisor,
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
		officeUser2 := factory.BuildOfficeUserWithPrivileges(suite.DB(), []factory.Customization{
			{
				Model: models.OfficeUser{
					Email:  "officeuser2@example.com",
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
					Roles: []roles.Role{
						{
							RoleType: roles.RoleTypeServicesCounselor,
						},
					},
				},
			},
		}, nil)
		officeUser3 := factory.BuildOfficeUserWithPrivileges(suite.DB(), []factory.Customization{
			{
				Model: models.OfficeUser{
					Email:  "officeuser3@example.com",
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
					Roles: []roles.Role{
						{
							RoleType: roles.RoleTypeServicesCounselor,
						},
					},
				},
			},
		}, nil)

		moves := []models.Move{move1, move2, move3, move4, move5, move6}
		userData := []*ghcmessages.BulkAssignmentForUser{
			{ID: strfmt.UUID(officeUser1.ID.String()), MoveAssignments: 1},
			{ID: strfmt.UUID(officeUser2.ID.String()), MoveAssignments: 2},
			{ID: strfmt.UUID(officeUser3.ID.String()), MoveAssignments: 3},
		}

		_, err := moveAssigner.BulkMoveAssignment(suite.AppContextForTest(), string(models.QueueTypeCounseling), userData, moves)
		suite.NoError(err)

		// reload move data to check assigned
		suite.NoError(suite.DB().Reload(&move1))
		suite.NoError(suite.DB().Reload(&move2))
		suite.NoError(suite.DB().Reload(&move3))
		suite.NoError(suite.DB().Reload(&move4))
		suite.NoError(suite.DB().Reload(&move5))
		suite.NoError(suite.DB().Reload(&move6))

		suite.Equal(officeUser1.ID, *move1.SCCounselingAssignedID)
		suite.Equal(officeUser2.ID, *move2.SCCounselingAssignedID)
		suite.Equal(officeUser3.ID, *move3.SCCounselingAssignedID)
		suite.Equal(officeUser2.ID, *move4.SCCounselingAssignedID)
		suite.Equal(officeUser3.ID, *move5.SCCounselingAssignedID)
		suite.Equal(officeUser3.ID, *move6.SCCounselingAssignedID)
	})

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
					Privileges: []roles.Privilege{
						{
							PrivilegeType: roles.PrivilegeTypeSupervisor,
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

		suite.Equal(officeUser.ID, *move1.SCCounselingAssignedID)
		suite.Equal(officeUser.ID, *move2.SCCounselingAssignedID)
		suite.Nil(move3.SCCounselingAssignedID)
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
					Privileges: []roles.Privilege{
						{
							PrivilegeType: roles.PrivilegeTypeSupervisor,
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

		suite.Equal(officeUser.ID, *move1.SCCloseoutAssignedID)
		suite.Equal(officeUser.ID, *move2.SCCloseoutAssignedID)
		suite.Nil(move3.SCCloseoutAssignedID)
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
					Privileges: []roles.Privilege{
						{
							PrivilegeType: roles.PrivilegeTypeSupervisor,
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

		suite.Equal(officeUser.ID, *move1.TOOTaskOrderAssignedID)
		suite.Equal(officeUser.ID, *move2.TOOTaskOrderAssignedID)
		suite.Nil(move3.TOOTaskOrderAssignedID)
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
					Privileges: []roles.Privilege{
						{
							PrivilegeType: roles.PrivilegeTypeSupervisor,
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

	suite.Run("successfully assigns multiple destination requests moves to a TOO destination user", func() {
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
					Privileges: []roles.Privilege{
						{
							PrivilegeType: roles.PrivilegeTypeSupervisor,
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

		_, err := moveAssigner.BulkMoveAssignment(suite.AppContextForTest(), string(models.QueueTypeDestinationRequest), userData, moves)
		suite.NoError(err)

		// reload move data to check assigned
		suite.NoError(suite.DB().Reload(&move1))
		suite.NoError(suite.DB().Reload(&move2))
		suite.NoError(suite.DB().Reload(&move3))

		suite.Equal(officeUser.ID, *move1.TOODestinationAssignedID)
		suite.Equal(officeUser.ID, *move2.TOODestinationAssignedID)
		suite.Nil(move3.TOODestinationAssignedID)
	})
}
