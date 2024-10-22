package adminapi

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/go-openapi/strfmt"
	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/mock"

	"github.com/transcom/mymove/pkg/factory"
	officeuserop "github.com/transcom/mymove/pkg/gen/adminapi/adminoperations/office_users"
	"github.com/transcom/mymove/pkg/gen/adminmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/models/roles"
	"github.com/transcom/mymove/pkg/services"
	fetch "github.com/transcom/mymove/pkg/services/fetch"
	"github.com/transcom/mymove/pkg/services/mocks"
	officeuser "github.com/transcom/mymove/pkg/services/office_user"
	"github.com/transcom/mymove/pkg/services/pagination"
	"github.com/transcom/mymove/pkg/services/query"
	rolesservice "github.com/transcom/mymove/pkg/services/roles"
	transportaionofficeassignments "github.com/transcom/mymove/pkg/services/transportation_office_assignments"
	usersprivileges "github.com/transcom/mymove/pkg/services/users_privileges"
	usersroles "github.com/transcom/mymove/pkg/services/users_roles"
)

func (suite *HandlerSuite) TestIndexOfficeUsersHandler() {
	setupTestData := func() models.OfficeUsers {
		return models.OfficeUsers{
			factory.BuildOfficeUserWithRoles(suite.DB(), factory.GetTraitApprovedOfficeUser(), []roles.RoleType{roles.RoleTypeQae}),
			factory.BuildOfficeUserWithRoles(suite.DB(), factory.GetTraitApprovedOfficeUser(), []roles.RoleType{roles.RoleTypeQae}),
		}
	}

	// test that everything is wired up
	suite.Run("integration test ok response", func() {
		officeUsers := setupTestData()
		params := officeuserop.IndexOfficeUsersParams{
			HTTPRequest: suite.setupAuthenticatedRequest("GET", "/office_users"),
		}

		queryBuilder := query.NewQueryBuilder()
		handler := IndexOfficeUsersHandler{
			HandlerConfig:  suite.HandlerConfig(),
			NewQueryFilter: query.NewQueryFilter,
			ListFetcher:    fetch.NewListFetcher(queryBuilder),
			NewPagination:  pagination.NewPagination,
		}

		response := handler.Handle(params)

		suite.IsType(&officeuserop.IndexOfficeUsersOK{}, response)
		okResponse := response.(*officeuserop.IndexOfficeUsersOK)
		suite.Len(okResponse.Payload, 2)
		suite.Equal(officeUsers[0].ID.String(), okResponse.Payload[0].ID.String())
		suite.Equal(string(officeUsers[0].User.Roles[0].RoleType), *okResponse.Payload[0].Roles[0].RoleType)
	})

	suite.Run("fetch return an empty list", func() {
		setupTestData()
		// TEST:				IndexOfficeUserHandler, Fetcher
		// Set up:				Provide an invalid search that won't be found
		// Expected Outcome:	An empty list is returned and we get a 200 OK.
		fakeFilter := "{\"search\":\"something\"}"

		params := officeuserop.IndexOfficeUsersParams{
			HTTPRequest: suite.setupAuthenticatedRequest("GET", "/office_users"),
			Filter:      &fakeFilter,
		}

		queryBuilder := query.NewQueryBuilder()
		handler := IndexOfficeUsersHandler{
			HandlerConfig:  suite.HandlerConfig(),
			ListFetcher:    fetch.NewListFetcher(queryBuilder),
			NewQueryFilter: query.NewQueryFilter,
			NewPagination:  pagination.NewPagination,
		}

		response := handler.Handle(params)
		okResponse := response.(*officeuserop.IndexOfficeUsersOK)

		suite.Len(okResponse.Payload, 0)
	})
}

func (suite *HandlerSuite) TestGetOfficeUserHandler() {
	// test that everything is wired up
	suite.Run("integration test ok response", func() {
		// Test:				GetOfficeUserHandler, Fetcher
		// Set up:				Provide a valid req with the office user ID to the endpoint
		// Expected Outcome:	The office user is returned and we get a 200 OK.
		officeUser := factory.BuildOfficeUserWithRoles(suite.DB(), nil, []roles.RoleType{roles.RoleTypeTOO})
		params := officeuserop.GetOfficeUserParams{
			HTTPRequest:  suite.setupAuthenticatedRequest("GET", fmt.Sprintf("/office_users/%s", officeUser.ID)),
			OfficeUserID: strfmt.UUID(officeUser.ID.String()),
		}

		// queryBuilder := query.NewQueryBuilder()
		handler := GetOfficeUserHandler{
			suite.HandlerConfig(),
			officeuser.NewOfficeUserFetcherPop(),
			query.NewQueryFilter,
		}

		response := handler.Handle(params)

		suite.IsType(&officeuserop.GetOfficeUserOK{}, response)
		okResponse := response.(*officeuserop.GetOfficeUserOK)
		suite.Equal(officeUser.ID.String(), okResponse.Payload.ID.String())
	})

	suite.Run("500 error - Internal Server error. Unsuccessful fetch ", func() {
		// Test:				GetOfficeUserHandler, Fetcher
		// Set up:				Provide a valid req with the fake office user ID to the endpoint
		// Expected Outcome:	The office user is returned and we get a 404 NotFound.
		fakeID := "3b9c2975-4e54-40ea-a781-bab7d6e4a502"
		params := officeuserop.GetOfficeUserParams{
			HTTPRequest:  suite.setupAuthenticatedRequest("GET", fmt.Sprintf("/office_users/%s", fakeID)),
			OfficeUserID: strfmt.UUID(fakeID),
		}

		// queryBuilder := query.NewQueryBuilder()
		handler := GetOfficeUserHandler{
			suite.HandlerConfig(),
			officeuser.NewOfficeUserFetcherPop(),
			query.NewQueryFilter,
		}

		response := handler.Handle(params)

		suite.IsType(&handlers.ErrResponse{}, response)
		errResponse := response.(*handlers.ErrResponse)
		suite.Equal(http.StatusInternalServerError, errResponse.Code)
	})
}

func (suite *HandlerSuite) TestCreateOfficeUserHandler() {
	tooRoleName := "Task Ordering Officer"
	tooRoleType := string(roles.RoleTypeTOO)

	tioRoleName := "Task Invoicing Officer"
	tioRoleType := string(roles.RoleTypeTIO)

	supervisorPrivilegeName := "Supervisor"
	supervisorPrivilegeType := string(models.PrivilegeTypeSupervisor)

	suite.Run("200 - Successfully create Office User", func() {
		// Test:				CreateOfficeUserHandler, Fetcher
		// Set up:				Create a new Office User, save new user to the DB
		// Expected Outcome:	The office user is created and we get a 200 OK.
		transportationOfficeID := factory.BuildDefaultTransportationOffice(suite.DB()).ID
		primaryOffice := true

		params := officeuserop.CreateOfficeUserParams{
			HTTPRequest: suite.setupAuthenticatedRequest("POST", "/office_users"),
			OfficeUser: &adminmessages.OfficeUserCreate{
				FirstName: "Sam",
				LastName:  "Cook",
				Telephone: "555-555-5555",
				Email:     "fakeemail@gmail.com",
				Roles: []*adminmessages.OfficeUserRole{
					{
						Name:     &tooRoleName,
						RoleType: &tooRoleType,
					},
					{
						Name:     &tioRoleName,
						RoleType: &tioRoleType,
					},
				},
				Privileges: []*adminmessages.OfficeUserPrivilege{
					{
						Name:          &supervisorPrivilegeName,
						PrivilegeType: &supervisorPrivilegeType,
					},
				},
				TransportationOfficeAssignments: []*adminmessages.OfficeUserTransportationOfficeAssignment{
					{
						TransportationOfficeID: strfmt.UUID(transportationOfficeID.String()),
						PrimaryOffice:          &primaryOffice,
					},
				},
			},
		}
		queryBuilder := query.NewQueryBuilder()
		handler := CreateOfficeUserHandler{
			suite.HandlerConfig(),
			officeuser.NewOfficeUserCreator(queryBuilder, suite.TestNotificationSender()),
			query.NewQueryFilter,
			usersroles.NewUsersRolesCreator(),
			rolesservice.NewRolesFetcher(),
			usersprivileges.NewUsersPrivilegesCreator(),
			transportaionofficeassignments.NewTransportaionOfficeAssignmentUpdater(),
		}
		suite.NoError(params.OfficeUser.Validate(strfmt.Default))
		response := handler.Handle(params)
		suite.IsType(&officeuserop.CreateOfficeUserCreated{}, response)
	})

	suite.Run("Failed create", func() {
		// Test:				CreateOfficeUserHandler, Fetcher
		// Set up:				Add new Office User to the DB
		// Expected Outcome:	The office user is not created and we get a 500 internal server error.
		fakeTransportationOfficeID := "3b9c2975-4e54-40ea-a781-bab7d6e4a502"
		primaryOffice := true
		officeUser := factory.BuildOfficeUser(suite.DB(), nil, []factory.Trait{
			factory.GetTraitOfficeUserWithID,
		})

		params := officeuserop.CreateOfficeUserParams{
			HTTPRequest: suite.setupAuthenticatedRequest("POST", "/office_users"),
			OfficeUser: &adminmessages.OfficeUserCreate{
				FirstName: officeUser.FirstName,
				LastName:  officeUser.LastName,
				Telephone: officeUser.Telephone,
				Roles: []*adminmessages.OfficeUserRole{
					{
						Name:     &tooRoleName,
						RoleType: &tooRoleType,
					},
					{
						Name:     &tioRoleName,
						RoleType: &tioRoleType,
					},
				},
				Privileges: []*adminmessages.OfficeUserPrivilege{
					{
						Name:          &supervisorPrivilegeName,
						PrivilegeType: &supervisorPrivilegeType,
					},
				},
				TransportationOfficeAssignments: []*adminmessages.OfficeUserTransportationOfficeAssignment{
					{
						TransportationOfficeID: strfmt.UUID(fakeTransportationOfficeID),
						PrimaryOffice:          &primaryOffice,
					},
				},
			},
		}

		queryBuilder := query.NewQueryBuilder()
		handler := CreateOfficeUserHandler{
			suite.HandlerConfig(),
			officeuser.NewOfficeUserCreator(queryBuilder, suite.TestNotificationSender()),
			query.NewQueryFilter,
			usersroles.NewUsersRolesCreator(),
			rolesservice.NewRolesFetcher(),
			usersprivileges.NewUsersPrivilegesCreator(),
			transportaionofficeassignments.NewTransportaionOfficeAssignmentUpdater(),
		}

		response := handler.Handle(params)
		suite.IsType(&officeuserop.CreateOfficeUserInternalServerError{}, response)
	})
}

func (suite *HandlerSuite) TestUpdateOfficeUserHandler() {
	setupHandler := func(updater services.OfficeUserUpdater, revoker services.UserSessionRevocation) UpdateOfficeUserHandler {
		handlerConfig := suite.HandlerConfig()
		return UpdateOfficeUserHandler{
			handlerConfig,
			updater,
			query.NewQueryFilter,
			usersroles.NewUsersRolesCreator(), // a special can of worms, TODO mocked tests
			usersprivileges.NewUsersPrivilegesCreator(),
			revoker,
			transportaionofficeassignments.NewTransportaionOfficeAssignmentUpdater(),
		}
	}

	setupTestData := func() models.OfficeUser {
		return factory.BuildOfficeUser(suite.DB(), []factory.Customization{
			{
				Model: models.TransportationOffice{
					Name: "Random Office",
				},
			},
		}, nil)
	}

	suite.Run("Office user is successfully updated", func() {
		officeUser := setupTestData()
		transportationOffice := factory.BuildTransportationOffice(nil, nil, nil)
		primaryOffice := true
		firstName := "Riley"
		middleInitials := "RB"
		telephone := "865-555-5309"

		officeUserUpdates := &adminmessages.OfficeUserUpdate{
			FirstName:      &firstName,
			MiddleInitials: &middleInitials,
			Telephone:      &telephone,
			TransportationOfficeAssignments: []*adminmessages.OfficeUserTransportationOfficeAssignment{
				{
					TransportationOfficeID: strfmt.UUID(transportationOffice.ID.String()),
					PrimaryOffice:          &primaryOffice,
				},
			},
		}

		params := officeuserop.UpdateOfficeUserParams{
			HTTPRequest:  suite.setupAuthenticatedRequest("PUT", fmt.Sprintf("/office_users/%s", officeUser.ID)),
			OfficeUserID: strfmt.UUID(officeUser.ID.String()),
			OfficeUser:   officeUserUpdates,
		}
		suite.NoError(params.OfficeUser.Validate(strfmt.Default))

		// Mock DB update:
		expectedInput := *officeUserUpdates // make a copy so we can ensure our expected values don't change
		expectedOfficeUser := officeUser
		expectedOfficeUser.FirstName = *expectedInput.FirstName
		expectedOfficeUser.MiddleInitials = expectedInput.MiddleInitials
		expectedOfficeUser.Telephone = *expectedInput.Telephone
		expectedOfficeUser.TransportationOfficeID = transportationOffice.ID

		mockUpdater := mocks.OfficeUserUpdater{}
		mockUpdater.On("UpdateOfficeUser", mock.AnythingOfType("*appcontext.appContext"), officeUser.ID, &expectedInput).Return(&expectedOfficeUser, nil, nil)

		response := setupHandler(&mockUpdater, &mocks.UserSessionRevocation{}).Handle(params)
		suite.IsType(&officeuserop.UpdateOfficeUserOK{}, response)

		okResponse := response.(*officeuserop.UpdateOfficeUserOK)
		// Check updates:
		suite.Equal(firstName, *okResponse.Payload.FirstName)
		suite.Equal(middleInitials, *okResponse.Payload.MiddleInitials)
		suite.Equal(telephone, *okResponse.Payload.Telephone)
		suite.Equal(transportationOffice.ID.String(), okResponse.Payload.TransportationOfficeID.String())
		suite.Equal(officeUser.LastName, *okResponse.Payload.LastName) // should not have been updated
		suite.Equal(officeUser.Email, *okResponse.Payload.Email)       // should not have been updated
	})

	suite.Run("Update fails due to bad Transportation Office", func() {
		officeUser := setupTestData()
		primaryOffice := true
		officeUserUpdates := &adminmessages.OfficeUserUpdate{
			TransportationOfficeAssignments: []*adminmessages.OfficeUserTransportationOfficeAssignment{
				{
					TransportationOfficeID: strfmt.UUID(uuid.Must(uuid.NewV4()).String()),
					PrimaryOffice:          &primaryOffice,
				},
			},
		}

		params := officeuserop.UpdateOfficeUserParams{
			HTTPRequest:  suite.setupAuthenticatedRequest("PUT", fmt.Sprintf("/office_users/%s", officeUser.ID)),
			OfficeUserID: strfmt.UUID(officeUser.ID.String()),
			OfficeUser:   officeUserUpdates,
		}
		suite.NoError(params.OfficeUser.Validate(strfmt.Default))

		expectedInput := *officeUserUpdates
		mockUpdater := mocks.OfficeUserUpdater{}
		mockUpdater.On("UpdateOfficeUser", mock.AnythingOfType("*appcontext.appContext"), officeUser.ID, &expectedInput).Return(nil, nil, sql.ErrNoRows)

		response := setupHandler(&mockUpdater, &mocks.UserSessionRevocation{}).Handle(params)
		suite.IsType(&officeuserop.UpdateOfficeUserInternalServerError{}, response)
	})

	suite.Run("Office user session is revoked when roles are changed", func() {
		officeUser := setupTestData()
		// Setup payload to remove all roles for office user
		newRoles := []*adminmessages.OfficeUserRole{}
		officeUserUpdates := &adminmessages.OfficeUserUpdate{
			Roles: newRoles,
		}
		params := officeuserop.UpdateOfficeUserParams{
			HTTPRequest:  suite.setupAuthenticatedRequest("PUT", fmt.Sprintf("/office_users/%s", officeUser.ID)),
			OfficeUserID: strfmt.UUID(officeUser.ID.String()),
			OfficeUser:   officeUserUpdates,
		}

		suite.NoError(params.OfficeUser.Validate(strfmt.Default))

		mockUpdater := mocks.OfficeUserUpdater{}
		mockUpdater.
			On("UpdateOfficeUser", mock.AnythingOfType("*appcontext.appContext"), officeUser.ID, officeUserUpdates).
			Return(&officeUser, nil, nil)

		expectedSessionUpdate := &adminmessages.UserUpdate{
			RevokeOfficeSession: models.BoolPointer(true),
		}
		mockRevoker := mocks.UserSessionRevocation{}
		mockRevoker.
			On("RevokeUserSession", mock.AnythingOfType("*appcontext.appContext"), *officeUser.UserID, expectedSessionUpdate, mock.Anything).
			Return(nil, nil, nil).
			Once()

		response := setupHandler(&mockUpdater, &mockRevoker).Handle(params)
		suite.IsType(&officeuserop.UpdateOfficeUserOK{}, response)
		mockRevoker.AssertNumberOfCalls(suite.T(), "RevokeUserSession", 1)
	})
}
