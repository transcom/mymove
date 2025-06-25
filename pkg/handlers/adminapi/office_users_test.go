package adminapi

import (
	"database/sql"
	"fmt"
	"net/http"
	"net/http/httptest"
	"slices"

	"github.com/go-openapi/strfmt"
	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/mock"

	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/factory"
	officeuserop "github.com/transcom/mymove/pkg/gen/adminapi/adminoperations/office_users"
	"github.com/transcom/mymove/pkg/gen/adminmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/handlers/adminapi/payloads"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/models/roles"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/mocks"
	officeuser "github.com/transcom/mymove/pkg/services/office_user"
	"github.com/transcom/mymove/pkg/services/pagination"
	"github.com/transcom/mymove/pkg/services/query"
	rolesservice "github.com/transcom/mymove/pkg/services/roles"
	transportationofficeassignments "github.com/transcom/mymove/pkg/services/transportation_office_assignments"
	usersprivileges "github.com/transcom/mymove/pkg/services/users_privileges"
	usersroles "github.com/transcom/mymove/pkg/services/users_roles"
)

func (suite *HandlerSuite) TestIndexOfficeUsersHandler() {
	// test that everything is wired up
	suite.Run("integration test ok response", func() {
		officeUsers := models.OfficeUsers{
			factory.BuildOfficeUserWithRoles(suite.DB(), factory.GetTraitApprovedOfficeUser(), []roles.RoleType{roles.RoleTypeQae}),
			factory.BuildOfficeUserWithRoles(suite.DB(), factory.GetTraitApprovedOfficeUser(), []roles.RoleType{roles.RoleTypeQae}),
			factory.BuildOfficeUserWithRoles(suite.DB(), factory.GetTraitApprovedOfficeUser(), []roles.RoleType{roles.RoleTypeQae, roles.RoleTypeQae, roles.RoleTypeCustomer, roles.RoleTypeContractingOfficer, roles.RoleTypeContractingOfficer}),
		}
		params := officeuserop.IndexOfficeUsersParams{
			HTTPRequest: suite.setupAuthenticatedRequest("GET", "/office_users"),
		}

		queryBuilder := query.NewQueryBuilder()
		handler := IndexOfficeUsersHandler{
			HandlerConfig:         suite.NewHandlerConfig(),
			NewQueryFilter:        query.NewQueryFilter,
			OfficeUserListFetcher: officeuser.NewOfficeUsersListFetcher(queryBuilder),
			NewPagination:         pagination.NewPagination,
		}

		response := handler.Handle(params)

		suite.IsType(&officeuserop.IndexOfficeUsersOK{}, response)
		okResponse := response.(*officeuserop.IndexOfficeUsersOK)

		actualOfficeUsers := okResponse.Payload
		suite.Equal(len(officeUsers), len(actualOfficeUsers))

		expectedOfficeUser1Id := officeUsers[0].ID.String()
		expectedOfficeUser2Id := officeUsers[1].ID.String()
		expectedOfficeUser3Id := officeUsers[2].ID.String()
		expectedOfficeUserIDs := []string{expectedOfficeUser1Id, expectedOfficeUser2Id, expectedOfficeUser3Id}

		for i := 0; i < len(actualOfficeUsers); i++ {
			suite.True(slices.Contains(expectedOfficeUserIDs, actualOfficeUsers[i].ID.String()))
		}
	})

	// Test that user roles list is not returning duplicate roles
	suite.Run("roles list has no duplicate roles", func() {
		officeUsers := models.OfficeUsers{
			factory.BuildOfficeUserWithRoles(suite.DB(), factory.GetTraitApprovedOfficeUser(), []roles.RoleType{roles.RoleTypeQae, roles.RoleTypeQae, roles.RoleTypeCustomer, roles.RoleTypeContractingOfficer, roles.RoleTypeContractingOfficer}),
		}
		params := officeuserop.IndexOfficeUsersParams{
			HTTPRequest: suite.setupAuthenticatedRequest("GET", "/office_users"),
		}

		queryBuilder := query.NewQueryBuilder()
		handler := IndexOfficeUsersHandler{
			HandlerConfig:         suite.NewHandlerConfig(),
			NewQueryFilter:        query.NewQueryFilter,
			OfficeUserListFetcher: officeuser.NewOfficeUsersListFetcher(queryBuilder),
			NewPagination:         pagination.NewPagination,
		}

		response := handler.Handle(params)

		suite.IsType(&officeuserop.IndexOfficeUsersOK{}, response)
		okResponse := response.(*officeuserop.IndexOfficeUsersOK)

		// Check payload user roles list for duplicate roles
		for _, r := range officeUsers[0].User.Roles {
			var dup = false
			var count = 0
			for _, r2 := range officeUsers[0].User.Roles {
				if r.RoleName == r2.RoleName {
					count++
				}
			}

			if count > 1 {
				dup = true
			}
			suite.False(dup)
		}

		suite.Len(okResponse.Payload, 1)
		suite.Len(officeUsers[0].User.Roles, 3)
	})

	suite.Run("invalid search returns no results", func() {
		fakeFilter := "{\"search\":\"invalidSearch\"}"

		params := officeuserop.IndexOfficeUsersParams{
			HTTPRequest: suite.setupAuthenticatedRequest("GET", "/office_users"),
			Filter:      &fakeFilter,
		}

		queryBuilder := query.NewQueryBuilder()
		handler := IndexOfficeUsersHandler{
			HandlerConfig:         suite.NewHandlerConfig(),
			OfficeUserListFetcher: officeuser.NewOfficeUsersListFetcher(queryBuilder),
			NewQueryFilter:        query.NewQueryFilter,
			NewPagination:         pagination.NewPagination,
		}

		response := handler.Handle(params)
		okResponse := response.(*officeuserop.IndexOfficeUsersOK)

		suite.Len(okResponse.Payload, 0)
	})

	suite.Run("able to search by name and filter", func() {
		status := models.OfficeUserStatusAPPROVED
		transportationOffice := factory.BuildTransportationOffice(suite.DB(), []factory.Customization{
			{
				Model: models.TransportationOffice{
					Name: "JPPO Test Office",
				},
			},
		}, nil)
		transportationOffice2 := factory.BuildTransportationOffice(suite.DB(), []factory.Customization{
			{
				Model: models.TransportationOffice{
					Name: "PPO Rome Test Office",
				},
			},
		}, nil)
		factory.BuildOfficeUserWithRoles(suite.DB(), []factory.Customization{
			{
				Model: models.OfficeUser{
					FirstName:              "Angelina",
					LastName:               "Jolie",
					Email:                  "laraCroft@mail.mil",
					Status:                 &status,
					Telephone:              "555-555-5555",
					TransportationOfficeID: transportationOffice.ID,
				},
			},
		}, []roles.RoleType{roles.RoleTypeTOO})
		factory.BuildOfficeUserWithRoles(suite.DB(), []factory.Customization{
			{
				Model: models.OfficeUser{
					FirstName:              "Billy",
					LastName:               "Bob",
					Email:                  "bigBob@mail.mil",
					Status:                 &status,
					Telephone:              "555-555-5555",
					TransportationOfficeID: transportationOffice.ID,
				},
			},
		}, []roles.RoleType{roles.RoleTypeTIO})
		factory.BuildOfficeUserWithRoles(suite.DB(), []factory.Customization{
			{
				Model: models.OfficeUser{
					FirstName:              "Nick",
					LastName:               "Cage",
					Email:                  "conAirKilluh@mail.mil",
					Status:                 &status,
					Telephone:              "555-555-5555",
					TransportationOfficeID: transportationOffice.ID,
				},
			},
		}, []roles.RoleType{roles.RoleTypeServicesCounselor})
		factory.BuildOfficeUserWithRoles(suite.DB(), []factory.Customization{
			{
				Model: models.OfficeUser{
					FirstName:              "Nick",
					LastName:               "Cage",
					Email:                  "conAirKilluh2@mail.mil",
					Status:                 &status,
					TransportationOfficeID: transportationOffice2.ID,
					Telephone:              "415-555-5555",
				},
			},
			{
				Model:    transportationOffice,
				LinkOnly: true,
				Type:     &factory.TransportationOffices.CounselingOffice,
			},
		}, []roles.RoleType{roles.RoleTypeServicesCounselor})

		// partial search
		nameSearch := "ic"
		filterJSON := fmt.Sprintf("{\"search\":\"%s\"}", nameSearch)
		params := officeuserop.IndexOfficeUsersParams{
			HTTPRequest: suite.setupAuthenticatedRequest("GET", "/office_users"),
			Filter:      &filterJSON,
		}

		queryBuilder := query.NewQueryBuilder()
		handler := IndexOfficeUsersHandler{
			HandlerConfig:         suite.NewHandlerConfig(),
			NewQueryFilter:        query.NewQueryFilter,
			OfficeUserListFetcher: officeuser.NewOfficeUsersListFetcher(queryBuilder),
			NewPagination:         pagination.NewPagination,
		}

		response := handler.Handle(params)

		suite.IsType(&officeuserop.IndexOfficeUsersOK{}, response)
		okResponse := response.(*officeuserop.IndexOfficeUsersOK)
		suite.Len(okResponse.Payload, 2)
		suite.Contains(*okResponse.Payload[0].FirstName, nameSearch)
		suite.Contains(*okResponse.Payload[1].FirstName, nameSearch)

		// email search
		emailSearch := "AirKilluh2"
		filterJSON = fmt.Sprintf("{\"email\":\"%s\"}", emailSearch)
		params = officeuserop.IndexOfficeUsersParams{
			HTTPRequest: suite.setupAuthenticatedRequest("GET", "/office_users"),
			Filter:      &filterJSON,
		}
		response = handler.Handle(params)

		suite.IsType(&officeuserop.IndexOfficeUsersOK{}, response)
		okResponse = response.(*officeuserop.IndexOfficeUsersOK)
		suite.Len(okResponse.Payload, 1)
		suite.Contains(*okResponse.Payload[0].Email, emailSearch)

		// telephone search
		phoneSearch := "415-"
		filterJSON = fmt.Sprintf("{\"phone\":\"%s\"}", phoneSearch)
		params = officeuserop.IndexOfficeUsersParams{
			HTTPRequest: suite.setupAuthenticatedRequest("GET", "/office_users"),
			Filter:      &filterJSON,
		}
		response = handler.Handle(params)

		suite.IsType(&officeuserop.IndexOfficeUsersOK{}, response)
		okResponse = response.(*officeuserop.IndexOfficeUsersOK)
		suite.Len(okResponse.Payload, 1)
		suite.Contains(*okResponse.Payload[0].Telephone, phoneSearch)

		// firstName search
		firstSearch := "Angel"
		filterJSON = fmt.Sprintf("{\"firstName\":\"%s\"}", firstSearch)
		params = officeuserop.IndexOfficeUsersParams{
			HTTPRequest: suite.setupAuthenticatedRequest("GET", "/office_users"),
			Filter:      &filterJSON,
		}
		response = handler.Handle(params)

		suite.IsType(&officeuserop.IndexOfficeUsersOK{}, response)
		okResponse = response.(*officeuserop.IndexOfficeUsersOK)
		suite.Len(okResponse.Payload, 1)
		suite.Contains(*okResponse.Payload[0].FirstName, firstSearch)

		// lastName search
		lastSearch := "Jo"
		filterJSON = fmt.Sprintf("{\"lastName\":\"%s\"}", lastSearch)
		params = officeuserop.IndexOfficeUsersParams{
			HTTPRequest: suite.setupAuthenticatedRequest("GET", "/office_users"),
			Filter:      &filterJSON,
		}
		response = handler.Handle(params)

		suite.IsType(&officeuserop.IndexOfficeUsersOK{}, response)
		okResponse = response.(*officeuserop.IndexOfficeUsersOK)
		suite.Len(okResponse.Payload, 1)
		suite.Contains(*okResponse.Payload[0].LastName, lastSearch)

		// transportation office search
		filterJSON = "{\"office\":\"Ro\"}"
		params = officeuserop.IndexOfficeUsersParams{
			HTTPRequest: suite.setupAuthenticatedRequest("GET", "/office_users"),
			Filter:      &filterJSON,
		}
		response = handler.Handle(params)

		suite.IsType(&officeuserop.IndexOfficeUsersOK{}, response)
		okResponse = response.(*officeuserop.IndexOfficeUsersOK)
		suite.Len(okResponse.Payload, 1)
		suite.Equal(strfmt.UUID(transportationOffice2.ID.String()), *okResponse.Payload[0].TransportationOfficeID)

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
			suite.NewHandlerConfig(),
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
		// Expected Outcome:	The office user is not returned and we get a 500 server error.
		fakeID := "3b9c2975-4e54-40ea-a781-bab7d6e4a502"
		params := officeuserop.GetOfficeUserParams{
			HTTPRequest:  suite.setupAuthenticatedRequest("GET", fmt.Sprintf("/office_users/%s", fakeID)),
			OfficeUserID: strfmt.UUID(fakeID),
		}

		// queryBuilder := query.NewQueryBuilder()
		handler := GetOfficeUserHandler{
			suite.NewHandlerConfig(),
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

	scRoleName := "Services Counselor"
	scRoleType := string(roles.RoleTypeServicesCounselor)

	supervisorPrivilegeName := "Supervisor"
	supervisorPrivilegeType := string(roles.PrivilegeTypeSupervisor)

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
						Name:     &scRoleName,
						RoleType: &scRoleType,
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
			suite.NewHandlerConfig(),
			officeuser.NewOfficeUserCreator(queryBuilder, suite.TestNotificationSender()),
			query.NewQueryFilter,
			usersroles.NewUsersRolesCreator(),
			rolesservice.NewRolesFetcher(),
			usersprivileges.NewUsersPrivilegesCreator(),
			transportationofficeassignments.NewTransportationOfficeAssignmentUpdater(),
		}
		suite.NoError(params.OfficeUser.Validate(strfmt.Default))
		response := handler.Handle(params)
		suite.IsType(&officeuserop.CreateOfficeUserCreated{}, response)
	})

	suite.Run("200 - Successfully create Office User with two transportation offices", func() {
		// Test:				CreateOfficeUserHandler, Fetcher
		// Set up:				Create a new Office User with two transportation offices, save new user to the DB
		// Expected Outcome:	The office user is created and we get a 200 OK.
		transportationOfficeID := factory.BuildDefaultTransportationOffice(suite.DB()).ID
		secondTransportationOfficeID := factory.BuildTransportationOffice(suite.DB(), nil, nil).ID

		params := officeuserop.CreateOfficeUserParams{
			HTTPRequest: suite.setupAuthenticatedRequest("POST", "/office_users"),
			OfficeUser: &adminmessages.OfficeUserCreate{
				FirstName:      "J",
				MiddleInitials: models.StringPointer("Jonah"),
				LastName:       "Jameson",
				Telephone:      "212-555-5555",
				Email:          "fakeemail2@dailybugle.com",
				Roles: []*adminmessages.OfficeUserRole{
					{
						Name:     &scRoleName,
						RoleType: &scRoleType,
					},
				},
				TransportationOfficeAssignments: []*adminmessages.OfficeUserTransportationOfficeAssignment{
					{
						TransportationOfficeID: strfmt.UUID(transportationOfficeID.String()),
						PrimaryOffice:          models.BoolPointer(true),
					},
					{
						TransportationOfficeID: strfmt.UUID(secondTransportationOfficeID.String()),
						PrimaryOffice:          models.BoolPointer(false),
					},
				},
			},
		}
		queryBuilder := query.NewQueryBuilder()
		handler := CreateOfficeUserHandler{
			suite.NewHandlerConfig(),
			officeuser.NewOfficeUserCreator(queryBuilder, suite.TestNotificationSender()),
			query.NewQueryFilter,
			usersroles.NewUsersRolesCreator(),
			rolesservice.NewRolesFetcher(),
			usersprivileges.NewUsersPrivilegesCreator(),
			transportationofficeassignments.NewTransportationOfficeAssignmentUpdater(),
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
						Name:     &scRoleName,
						RoleType: &scRoleName,
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
			suite.NewHandlerConfig(),
			officeuser.NewOfficeUserCreator(queryBuilder, suite.TestNotificationSender()),
			query.NewQueryFilter,
			usersroles.NewUsersRolesCreator(),
			rolesservice.NewRolesFetcher(),
			usersprivileges.NewUsersPrivilegesCreator(),
			transportationofficeassignments.NewTransportationOfficeAssignmentUpdater(),
		}

		response := handler.Handle(params)
		suite.IsType(&officeuserop.CreateOfficeUserUnprocessableEntity{}, response)
	})

	suite.Run("Failed create due to missing transportation office assignments", func() {
		// Test:				CreateOfficeUserHandler, Fetcher
		// Set up:				Submit an office user with no transportation offices
		// Expected Outcome:	The office user is not created and we get an unprocessible entity error.
		params := officeuserop.CreateOfficeUserParams{
			HTTPRequest: suite.setupAuthenticatedRequest("POST", "/office_users"),
			OfficeUser: &adminmessages.OfficeUserCreate{
				FirstName: "Sam",
				LastName:  "Cook",
				Telephone: "555-555-5555",
				Email:     "fakeemail3@gmail.com",
				Roles: []*adminmessages.OfficeUserRole{
					{
						Name:     &tooRoleName,
						RoleType: &tooRoleType,
					},
				},
			},
		}

		queryBuilder := query.NewQueryBuilder()
		handler := CreateOfficeUserHandler{
			suite.NewHandlerConfig(),
			officeuser.NewOfficeUserCreator(queryBuilder, suite.TestNotificationSender()),
			query.NewQueryFilter,
			usersroles.NewUsersRolesCreator(),
			rolesservice.NewRolesFetcher(),
			usersprivileges.NewUsersPrivilegesCreator(),
			transportationofficeassignments.NewTransportationOfficeAssignmentUpdater(),
		}

		response := handler.Handle(params)
		suite.IsType(&officeuserop.CreateOfficeUserUnprocessableEntity{}, response)
	})

	suite.Run("Fail create Office User due to no primary transportation office", func() {
		// Test:				CreateOfficeUserHandler, Fetcher
		// Set up:				Submit an office user with two non-primary transportation offices
		// Expected Outcome:	The office user is not created and we get an unprocessible entity error.
		transportationOfficeID := factory.BuildDefaultTransportationOffice(suite.DB()).ID
		secondTransportationOfficeID := factory.BuildTransportationOffice(suite.DB(), nil, nil).ID

		params := officeuserop.CreateOfficeUserParams{
			HTTPRequest: suite.setupAuthenticatedRequest("POST", "/office_users"),
			OfficeUser: &adminmessages.OfficeUserCreate{
				FirstName: "Bob",
				LastName:  "Loblaw",
				Telephone: "949-555-5555",
				Email:     "bob.loblaw@attornyatlaw.org",
				Roles: []*adminmessages.OfficeUserRole{
					{
						Name:     &scRoleName,
						RoleType: &scRoleType,
					},
				},
				TransportationOfficeAssignments: []*adminmessages.OfficeUserTransportationOfficeAssignment{
					{
						TransportationOfficeID: strfmt.UUID(transportationOfficeID.String()),
						PrimaryOffice:          models.BoolPointer(false),
					},
					{
						TransportationOfficeID: strfmt.UUID(secondTransportationOfficeID.String()),
						PrimaryOffice:          models.BoolPointer(false),
					},
				},
			},
		}
		queryBuilder := query.NewQueryBuilder()
		handler := CreateOfficeUserHandler{
			suite.NewHandlerConfig(),
			officeuser.NewOfficeUserCreator(queryBuilder, suite.TestNotificationSender()),
			query.NewQueryFilter,
			usersroles.NewUsersRolesCreator(),
			rolesservice.NewRolesFetcher(),
			usersprivileges.NewUsersPrivilegesCreator(),
			transportationofficeassignments.NewTransportationOfficeAssignmentUpdater(),
		}

		response := handler.Handle(params)
		suite.IsType(&officeuserop.CreateOfficeUserUnprocessableEntity{}, response)
	})

	suite.Run("Failed create due to missing roles", func() {
		// Test:				CreateOfficeUserHandler, Fetcher
		// Set up:				Submit an office user with no roles
		// Expected Outcome:	The office user is not created and we get an unprocessible entity error.
		transportationOfficeID := factory.BuildDefaultTransportationOffice(suite.DB()).ID

		params := officeuserop.CreateOfficeUserParams{
			HTTPRequest: suite.setupAuthenticatedRequest("POST", "/office_users"),
			OfficeUser: &adminmessages.OfficeUserCreate{
				FirstName: "Bob",
				LastName:  "Loblaw",
				Telephone: "949-555-5555",
				Email:     "bob.loblaw2@attornyatlaw.org",
				TransportationOfficeAssignments: []*adminmessages.OfficeUserTransportationOfficeAssignment{
					{
						TransportationOfficeID: strfmt.UUID(transportationOfficeID.String()),
						PrimaryOffice:          models.BoolPointer(true),
					},
				},
			},
		}

		queryBuilder := query.NewQueryBuilder()
		handler := CreateOfficeUserHandler{
			suite.NewHandlerConfig(),
			officeuser.NewOfficeUserCreator(queryBuilder, suite.TestNotificationSender()),
			query.NewQueryFilter,
			usersroles.NewUsersRolesCreator(),
			rolesservice.NewRolesFetcher(),
			usersprivileges.NewUsersPrivilegesCreator(),
			transportationofficeassignments.NewTransportationOfficeAssignmentUpdater(),
		}

		response := handler.Handle(params)
		suite.IsType(&officeuserop.CreateOfficeUserUnprocessableEntity{}, response)
	})

	suite.Run("Failed create due validation errors creating office user", func() {
		// Test:				CreateOfficeUserHandler, Fetcher
		// Set up:				Submit a valid office user and mock failed office user validation
		// Expected Outcome:	The office user is not created and we get an unprocessible entity error.
		transportationOfficeID := factory.BuildDefaultTransportationOffice(suite.DB()).ID

		params := officeuserop.CreateOfficeUserParams{
			HTTPRequest: suite.setupAuthenticatedRequest("POST", "/office_users"),
			OfficeUser: &adminmessages.OfficeUserCreate{
				FirstName: "Sam",
				LastName:  "Cook",
				Telephone: "555-555-5555",
				Email:     "fakeemail5@gmail.com",
				Roles: []*adminmessages.OfficeUserRole{
					{
						Name:     &tooRoleName,
						RoleType: &tooRoleType,
					},
				},
				TransportationOfficeAssignments: []*adminmessages.OfficeUserTransportationOfficeAssignment{
					{
						TransportationOfficeID: strfmt.UUID(transportationOfficeID.String()),
						PrimaryOffice:          models.BoolPointer(true),
					},
				},
			},
		}

		expectedError := &validate.Errors{}
		officeUserCreator := &mocks.OfficeUserCreator{}
		officeUserCreator.On("CreateOfficeUser",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything,
			mock.Anything,
		).Return(nil, expectedError, nil).Once()
		handler := CreateOfficeUserHandler{
			suite.NewHandlerConfig(),
			officeUserCreator,
			query.NewQueryFilter,
			usersroles.NewUsersRolesCreator(),
			rolesservice.NewRolesFetcher(),
			usersprivileges.NewUsersPrivilegesCreator(),
			transportationofficeassignments.NewTransportationOfficeAssignmentUpdater(),
		}

		response := handler.Handle(params)
		suite.IsType(&officeuserop.CreateOfficeUserUnprocessableEntity{}, response)
	})

	suite.Run("Failed create due validation errors creating office user", func() {
		// Test:				CreateOfficeUserHandler, Fetcher
		// Set up:				Submit a valid office user and mock failed roles validation
		// Expected Outcome:	The office user is not created and we get an unprocessible entity error.
		transportationOfficeID := factory.BuildDefaultTransportationOffice(suite.DB()).ID

		params := officeuserop.CreateOfficeUserParams{
			HTTPRequest: suite.setupAuthenticatedRequest("POST", "/office_users"),
			OfficeUser: &adminmessages.OfficeUserCreate{
				FirstName: "Sam",
				LastName:  "Cook",
				Telephone: "555-555-5555",
				Email:     "fakeemail5@gmail.com",
				Roles: []*adminmessages.OfficeUserRole{
					{
						Name:     &tooRoleName,
						RoleType: &tooRoleType,
					},
				},
				TransportationOfficeAssignments: []*adminmessages.OfficeUserTransportationOfficeAssignment{
					{
						TransportationOfficeID: strfmt.UUID(transportationOfficeID.String()),
						PrimaryOffice:          models.BoolPointer(true),
					},
				},
			},
		}

		expectedError := &validate.Errors{Errors: map[string][]string{
			"role_name": {"What's a TXO? I've only heard of TIO and TOO."},
		}}
		userRoleAssociator := &mocks.UserRoleAssociator{}
		userRoleAssociator.On("UpdateUserRoles",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything,
			mock.Anything,
		).Return(nil, expectedError, nil).Once()
		queryBuilder := query.NewQueryBuilder()
		handler := CreateOfficeUserHandler{
			suite.NewHandlerConfig(),
			officeuser.NewOfficeUserCreator(queryBuilder, suite.TestNotificationSender()),
			query.NewQueryFilter,
			userRoleAssociator,
			rolesservice.NewRolesFetcher(),
			usersprivileges.NewUsersPrivilegesCreator(),
			transportationofficeassignments.NewTransportationOfficeAssignmentUpdater(),
		}

		response := handler.Handle(params)
		suite.IsType(&officeuserop.CreateOfficeUserUnprocessableEntity{}, response)
	})

	suite.Run("Failed to create due to Supervisor privileges not authorized", func() {
		transportationOfficeID := factory.BuildDefaultTransportationOffice(suite.DB()).ID
		supervisorPrivilegeName := "Supervisor"
		supervisorPrivilegeType := string(roles.PrivilegeTypeSupervisor)
		primeRoleName := "Prime"
		primeRoleType := string(roles.RoleTypePrime)
		params := officeuserop.CreateOfficeUserParams{
			HTTPRequest: suite.setupAuthenticatedRequest("POST", "/office_users"),
			OfficeUser: &adminmessages.OfficeUserCreate{
				FirstName: "Sam",
				LastName:  "Cook",
				Telephone: "555-555-5555",
				Email:     "fakeemail5@gmail.com",
				Roles: []*adminmessages.OfficeUserRole{
					{
						Name:     &primeRoleName,
						RoleType: &primeRoleType,
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
						PrimaryOffice:          models.BoolPointer(true),
					},
				},
			},
		}

		queryBuilder := query.NewQueryBuilder()
		handler := CreateOfficeUserHandler{
			suite.NewHandlerConfig(),
			officeuser.NewOfficeUserCreator(queryBuilder, suite.TestNotificationSender()),
			query.NewQueryFilter,
			usersroles.NewUsersRolesCreator(),
			rolesservice.NewRolesFetcher(),
			usersprivileges.NewUsersPrivilegesCreator(),
			transportationofficeassignments.NewTransportationOfficeAssignmentUpdater(),
		}

		response := handler.Handle(params)
		suite.IsType(&officeuserop.CreateOfficeUserUnprocessableEntity{}, response)
	})

	suite.Run("Update fails due to Safety privileges not authorized", func() {
		transportationOfficeID := factory.BuildDefaultTransportationOffice(suite.DB()).ID
		safetyPrivilegeName := "Safety"
		safetyPrivilegeType := string(roles.PrivilegeSearchTypeSafety)
		contractingOfficerRoleName := "Contracting Officer"
		contractingOfficerRoleType := string(roles.RoleTypeContractingOfficer)
		params := officeuserop.CreateOfficeUserParams{
			HTTPRequest: suite.setupAuthenticatedRequest("POST", "/office_users"),
			OfficeUser: &adminmessages.OfficeUserCreate{
				FirstName: "Sam",
				LastName:  "Cook",
				Telephone: "555-555-5555",
				Email:     "fakeemail5@gmail.com",
				Roles: []*adminmessages.OfficeUserRole{
					{
						Name:     &contractingOfficerRoleName,
						RoleType: &contractingOfficerRoleType,
					},
				},
				Privileges: []*adminmessages.OfficeUserPrivilege{
					{
						Name:          &safetyPrivilegeName,
						PrivilegeType: &safetyPrivilegeType,
					},
				},
				TransportationOfficeAssignments: []*adminmessages.OfficeUserTransportationOfficeAssignment{
					{
						TransportationOfficeID: strfmt.UUID(transportationOfficeID.String()),
						PrimaryOffice:          models.BoolPointer(true),
					},
				},
			},
		}

		queryBuilder := query.NewQueryBuilder()
		handler := CreateOfficeUserHandler{
			suite.NewHandlerConfig(),
			officeuser.NewOfficeUserCreator(queryBuilder, suite.TestNotificationSender()),
			query.NewQueryFilter,
			usersroles.NewUsersRolesCreator(),
			rolesservice.NewRolesFetcher(),
			usersprivileges.NewUsersPrivilegesCreator(),
			transportationofficeassignments.NewTransportationOfficeAssignmentUpdater(),
		}

		response := handler.Handle(params)
		suite.IsType(&officeuserop.CreateOfficeUserUnprocessableEntity{}, response)
	})
}

func (suite *HandlerSuite) TestUpdateOfficeUserHandler() {
	setupHandler := func(updater services.OfficeUserUpdater, revoker services.UserSessionRevocation) UpdateOfficeUserHandler {
		handlerConfig := suite.NewHandlerConfig()
		return UpdateOfficeUserHandler{
			handlerConfig,
			updater,
			query.NewQueryFilter,
			usersroles.NewUsersRolesCreator(), // a special can of worms, TODO mocked tests
			usersprivileges.NewUsersPrivilegesCreator(),
			revoker,
			transportationofficeassignments.NewTransportationOfficeAssignmentUpdater(),
			rolesservice.NewRolesFetcher(),
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
		transportationOffice := factory.BuildTransportationOffice(suite.DB(), nil, nil)
		primaryOffice := true
		firstName := "Riley"
		middleInitials := "RB"
		telephone := "865-555-5309"
		supervisorPrivilegeName := "Supervisor"
		supervisorPrivilegeType := string(roles.PrivilegeTypeSupervisor)
		tooRoleName := "Task Ordering Officer"
		tooRoleType := string(roles.RoleTypeTOO)

		officeUserUpdates := &adminmessages.OfficeUserUpdate{
			FirstName:      &firstName,
			MiddleInitials: &middleInitials,
			Telephone:      &telephone,
			Privileges: []*adminmessages.OfficeUserPrivilege{
				{
					Name:          &supervisorPrivilegeName,
					PrivilegeType: &supervisorPrivilegeType,
				},
			},
			Roles: []*adminmessages.OfficeUserRole{
				{
					Name:     &tooRoleName,
					RoleType: &tooRoleType,
				},
			},
			TransportationOfficeAssignments: []*adminmessages.OfficeUserTransportationOfficeAssignment{
				{
					TransportationOfficeID: strfmt.UUID(transportationOffice.ID.String()),
					PrimaryOffice:          &primaryOffice,
				},
			},
		}

		officeUserUpdatesModel := payloads.OfficeUserModelFromUpdate(officeUserUpdates, &officeUser)

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
		expectedOfficeUser.User.Roles = roles.Roles{roles.Role{RoleType: roles.RoleTypeTOO}}
		expectedOfficeUser.User.Privileges = roles.Privileges{roles.Privilege{PrivilegeType: roles.PrivilegeSearchTypeSupervisor}}

		mockUpdater := mocks.OfficeUserUpdater{}
		mockUpdater.On("UpdateOfficeUser", mock.AnythingOfType("*appcontext.appContext"), officeUser.ID, officeUserUpdatesModel, transportationOffice.ID).Return(&expectedOfficeUser, nil, nil)
		queryBuilder := query.NewQueryBuilder()
		officeUserUpdater := officeuser.NewOfficeUserUpdater(queryBuilder)

		expectedSessionUpdate := &adminmessages.UserUpdate{
			RevokeOfficeSession: models.BoolPointer(true),
		}
		mockRevoker := mocks.UserSessionRevocation{}
		mockRevoker.
			On("RevokeUserSession", mock.AnythingOfType("*appcontext.appContext"), *officeUser.UserID, expectedSessionUpdate, mock.Anything).
			Return(nil, nil, nil).
			Times(3)

		response := setupHandler(officeUserUpdater, &mockRevoker).Handle(params)
		suite.IsType(&officeuserop.UpdateOfficeUserOK{}, response)

		okResponse := response.(*officeuserop.UpdateOfficeUserOK)
		// Check updates:
		suite.Equal(firstName, *okResponse.Payload.FirstName)
		suite.Equal(middleInitials, *okResponse.Payload.MiddleInitials)
		suite.Equal(telephone, *okResponse.Payload.Telephone)
		suite.Equal(transportationOffice.ID.String(), okResponse.Payload.TransportationOfficeID.String())
		suite.Equal(transportationOffice.ID.String(), okResponse.Payload.TransportationOfficeAssignments[0].TransportationOfficeID.String())
		suite.Equal(tooRoleName, *okResponse.Payload.Roles[0].RoleName)
		suite.Equal(supervisorPrivilegeName, okResponse.Payload.Privileges[0].PrivilegeName)
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

		officeUserDB, _ := models.FetchOfficeUserByID(suite.DB(), officeUser.ID)

		officeUserUpdatesModel := payloads.OfficeUserModelFromUpdate(officeUserUpdates, officeUserDB)

		mockUpdater := mocks.OfficeUserUpdater{}
		mockUpdater.On("UpdateOfficeUser",
			mock.AnythingOfType("*appcontext.appContext"),
			officeUser.ID,
			officeUserUpdatesModel,
			uuid.FromStringOrNil(officeUserUpdates.TransportationOfficeAssignments[0].TransportationOfficeID.String()),
		).Return(nil, nil, sql.ErrNoRows)

		expectedSessionUpdate := &adminmessages.UserUpdate{
			RevokeOfficeSession: models.BoolPointer(true),
		}
		mockRevoker := mocks.UserSessionRevocation{}
		mockRevoker.
			On("RevokeUserSession", mock.AnythingOfType("*appcontext.appContext"), *officeUser.UserID, expectedSessionUpdate, mock.Anything).
			Return(nil, nil, nil).
			Once()

		response := setupHandler(&mockUpdater, &mockRevoker).Handle(params)
		suite.IsType(&officeuserop.UpdateOfficeUserInternalServerError{}, response)
	})

	suite.Run("Returns not found when office user does not exist in DB", func() {
		officeUser := setupTestData()
		transportationOffice := factory.BuildTransportationOffice(suite.DB(), nil, nil)
		primaryOffice := true
		firstName := "Riley"
		middleInitials := "RB"
		telephone := "865-555-5309"
		supervisorPrivilegeName := "Supervisor"
		supervisorPrivilegeType := string(roles.PrivilegeTypeSupervisor)
		tooRoleName := "Task Ordering Officer"
		tooRoleType := string(roles.RoleTypeTOO)

		officeUserUpdates := &adminmessages.OfficeUserUpdate{
			FirstName:      &firstName,
			MiddleInitials: &middleInitials,
			Telephone:      &telephone,
			Privileges: []*adminmessages.OfficeUserPrivilege{
				{
					Name:          &supervisorPrivilegeName,
					PrivilegeType: &supervisorPrivilegeType,
				},
			},
			Roles: []*adminmessages.OfficeUserRole{
				{
					Name:     &tooRoleName,
					RoleType: &tooRoleType,
				},
			},
			TransportationOfficeAssignments: []*adminmessages.OfficeUserTransportationOfficeAssignment{
				{
					TransportationOfficeID: strfmt.UUID(transportationOffice.ID.String()),
					PrimaryOffice:          &primaryOffice,
				},
			},
		}

		fakeID := uuid.Must(uuid.NewV4())

		officeUserDB, err := models.FetchOfficeUserByID(suite.DB(), fakeID)
		suite.Error(err)
		suite.Equal(uuid.Nil, officeUserDB.ID)

		params := officeuserop.UpdateOfficeUserParams{
			HTTPRequest:  suite.setupAuthenticatedRequest("PUT", fmt.Sprintf("/office_users/%s", officeUser.ID)),
			OfficeUserID: strfmt.UUID(fakeID.String()),
			OfficeUser:   officeUserUpdates,
		}
		suite.NoError(params.OfficeUser.Validate(strfmt.Default))

		mockUpdater := mocks.OfficeUserUpdater{}
		mockRevoker := mocks.UserSessionRevocation{}

		response := setupHandler(&mockUpdater, &mockRevoker).Handle(params)
		suite.IsType(&officeuserop.UpdateOfficeUserNotFound{}, response)
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

		officeUserDB, _ := models.FetchOfficeUserByID(suite.DB(), officeUser.ID)

		officeUserUpdatesModel := payloads.OfficeUserModelFromUpdate(officeUserUpdates, officeUserDB)

		mockUpdater := mocks.OfficeUserUpdater{}
		mockUpdater.
			On("UpdateOfficeUser", mock.AnythingOfType("*appcontext.appContext"), officeUser.ID, officeUserUpdatesModel, uuid.Nil).
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

	suite.Run("Update fails due to Safety privileges not authorized", func() {
		officeUser := setupTestData()
		transportationOffice := factory.BuildTransportationOffice(suite.DB(), nil, nil)
		primaryOffice := true
		safetyPrivilegeName := "Safety"
		safetyPrivilegeType := string(roles.PrivilegeSearchTypeSafety)
		contractingOfficerRoleName := "Contracting Officer"
		contractingOfficerRoleType := string(roles.RoleTypeContractingOfficer)

		officeUserUpdates := &adminmessages.OfficeUserUpdate{
			Privileges: []*adminmessages.OfficeUserPrivilege{
				{
					Name:          &safetyPrivilegeName,
					PrivilegeType: &safetyPrivilegeType,
				},
			},
			Roles: []*adminmessages.OfficeUserRole{
				{
					Name:     &contractingOfficerRoleName,
					RoleType: &contractingOfficerRoleType,
				},
			},
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

		expectedInput := *officeUserUpdates
		mockUpdater := mocks.OfficeUserUpdater{}
		mockUpdater.On("UpdateOfficeUser",
			mock.AnythingOfType("*appcontext.appContext"),
			officeUser.ID,
			&expectedInput,
			uuid.FromStringOrNil(officeUserUpdates.TransportationOfficeAssignments[0].TransportationOfficeID.String()),
		).Return(nil, nil, sql.ErrNoRows)

		expectedSessionUpdate := &adminmessages.UserUpdate{
			RevokeOfficeSession: models.BoolPointer(true),
		}
		mockRevoker := mocks.UserSessionRevocation{}
		mockRevoker.
			On("RevokeUserSession", mock.AnythingOfType("*appcontext.appContext"), *officeUser.UserID, expectedSessionUpdate, mock.Anything).
			Return(nil, nil, nil).
			Once()

		response := setupHandler(&mockUpdater, &mockRevoker).Handle(params)
		suite.IsType(&officeuserop.CreateOfficeUserUnprocessableEntity{}, response)
	})
}

func (suite *HandlerSuite) TestDeleteOfficeUsersHandler() {
	suite.Run("deleted requested users results in no content (successful) response", func() {
		user := factory.BuildDefaultUser(suite.DB())
		status := models.OfficeUserStatusREQUESTED
		officeUser := factory.BuildOfficeUserWithRoles(suite.DB(), []factory.Customization{
			{
				Model: models.OfficeUser{
					Active: true,
					UserID: &user.ID,
					Email:  user.OktaEmail,
					Status: &status,
				},
			},
			{
				Model:    user,
				LinkOnly: true,
			},
		}, []roles.RoleType{roles.RoleTypeTOO})
		officeUserID := officeUser.ID

		params := officeuserop.DeleteOfficeUserParams{
			HTTPRequest:  suite.setupAuthenticatedRequest("DELETE", fmt.Sprintf("/office_users/%s", officeUserID)),
			OfficeUserID: *handlers.FmtUUID(officeUserID),
		}

		queryBuilder := query.NewQueryBuilder()
		handler := DeleteOfficeUserHandler{
			HandlerConfig:     suite.NewHandlerConfig(),
			OfficeUserDeleter: officeuser.NewOfficeUserDeleter(queryBuilder),
		}

		response := handler.Handle(params)

		suite.IsType(&officeuserop.DeleteOfficeUserNoContent{}, response)

		var dbUser models.User
		err := suite.DB().Where("id = ?", user.ID).First(&dbUser)
		suite.Error(err)
		suite.Equal(sql.ErrNoRows, err, "sql: no rows in result set")

		var dbOfficeUser models.OfficeUser
		err = suite.DB().Where("user_id = ?", user.ID).First(&dbOfficeUser)
		suite.Error(err)
		suite.Equal(sql.ErrNoRows, err, "sql: no rows in result set")

		// .All does not return a sql no rows error, so we will verify that the struct is empty
		var userRoles []models.UsersRoles
		err = suite.DB().Where("user_id = ?", user.ID).All(&userRoles)
		suite.NoError(err)
		suite.Empty(userRoles, "Expected no roles to remain for the user")

		var userPrivileges []models.UsersPrivileges
		err = suite.DB().Where("user_id = ?", user.ID).All(&userPrivileges)
		suite.NoError(err)
		suite.Empty(userPrivileges, "Expected no privileges to remain for the user")
	})

	suite.Run("get an error when the office user does not exist", func() {
		officeUserID := uuid.Must(uuid.NewV4())

		params := officeuserop.DeleteOfficeUserParams{
			HTTPRequest:  suite.setupAuthenticatedRequest("DELETE", fmt.Sprintf("/office_users/%s", officeUserID)),
			OfficeUserID: *handlers.FmtUUID(officeUserID),
		}

		queryBuilder := query.NewQueryBuilder()
		handler := DeleteOfficeUserHandler{
			HandlerConfig:     suite.NewHandlerConfig(),
			OfficeUserDeleter: officeuser.NewOfficeUserDeleter(queryBuilder),
		}

		response := handler.Handle(params)

		suite.IsType(&officeuserop.DeleteOfficeUserNotFound{}, response)
	})

	suite.Run("error response when a user is not in the admin application", func() {
		officeUser := factory.BuildOfficeUser(suite.DB(), nil, nil)
		officeUserID := officeUser.ID
		req := httptest.NewRequest("DELETE", fmt.Sprintf("/office_users/%s", officeUserID), nil)

		session := &auth.Session{
			ApplicationName: auth.OfficeApp,
			OfficeUserID:    officeUserID,
		}
		ctx := auth.SetSessionInRequestContext(req, session)

		params := officeuserop.DeleteOfficeUserParams{
			HTTPRequest:  req.WithContext(ctx),
			OfficeUserID: *handlers.FmtUUID(officeUserID),
		}

		queryBuilder := query.NewQueryBuilder()
		handler := DeleteOfficeUserHandler{
			HandlerConfig:     suite.NewHandlerConfig(),
			OfficeUserDeleter: officeuser.NewOfficeUserDeleter(queryBuilder),
		}

		response := handler.Handle(params)

		suite.IsType(&officeuserop.DeleteOfficeUserUnauthorized{}, response)
	})
}

func (suite *HandlerSuite) TestGetRolesPrivilegesHandler() {
	suite.Run("200 OK - successfully retrieve unique role privilege mappings", func() {
		// Test:				GetOfficeUserHandler, Fetcher
		// Set up:				Login as admin user
		// Expected Outcome:	The list of unique role privlege mappings
		params := officeuserop.GetRolesPrivilegesParams{
			HTTPRequest: suite.setupAuthenticatedRequest("GET", "/office_users/roles-privileges"),
		}

		handler := GetRolesPrivilegesHandler{
			suite.NewHandlerConfig(),
			rolesservice.NewRolesFetcher(),
		}

		rolePrivs, err := handler.RoleAssociater.FetchRolesPrivileges(suite.AppContextForTest())

		suite.NoError(err)

		response := handler.Handle(params)

		suite.IsType(&officeuserop.GetRolesPrivilegesOK{}, response)
		okResponse := response.(*officeuserop.GetRolesPrivilegesOK)
		suite.Len(okResponse.Payload, len(rolePrivs))

		type privValidation struct {
			PrivilegeType string
			PrivilegeName string
		}

		type rolePrivValidation struct {
			RoleType   string
			RoleName   string
			Privileges []privValidation
		}

		rolePrivEntries := make(map[uuid.UUID]*rolePrivValidation)

		for _, rp := range rolePrivs {
			rid := rp.ID

			if _, ok := rolePrivEntries[rid]; !ok {
				rolePrivEntries[rid] = &rolePrivValidation{
					RoleType:   string(rp.RoleType),
					RoleName:   string(rp.RoleName),
					Privileges: []privValidation{},
				}
			}
			for _, resPriv := range rp.RolePrivileges {
				rolePrivEntries[rid].Privileges = append(rolePrivEntries[rid].Privileges, privValidation{
					PrivilegeType: string(resPriv.Privilege.PrivilegeType),
					PrivilegeName: string(resPriv.Privilege.PrivilegeName),
				})
			}
		}

		for _, resRolePriv := range okResponse.Payload {
			entryKey, err := uuid.FromString(resRolePriv.ID.String())
			suite.NoError(err)
			rolePriv, ok := rolePrivEntries[entryKey]
			suite.NotNil(ok)
			suite.Equal(rolePriv.RoleType, *resRolePriv.RoleType)
			suite.Equal(rolePriv.RoleName, *resRolePriv.RoleName)
			for i, priv := range rolePriv.Privileges {
				suite.Equal(priv.PrivilegeType, resRolePriv.Privileges[i].PrivilegeType)
				suite.Equal(priv.PrivilegeName, resRolePriv.Privileges[i].PrivilegeName)
			}
			delete(rolePrivEntries, entryKey) // remove to ensure unique values
		}
		suite.Len(rolePrivEntries, 0, "all entries should have been deleted")
	})

	suite.Run("401 ERROR - Unauthorized ", func() {
		// Test:				GetOfficeUserHandler, Fetcher - Unauthorized
		// Set up:				Run request when NOT logged in as admin user
		// Expected Outcome:	Unauthorized response returned, no data
		requestUser := factory.BuildOfficeUser(nil, []factory.Customization{
			{
				Model: models.User{
					Roles: roles.Roles{
						{
							RoleType: roles.RoleTypeTOO,
						},
					},
				},
			},
		}, nil)

		req := httptest.NewRequest("GET", "/office_users/roles-privileges", nil) // We never need to set a body this endpoint

		params := officeuserop.GetRolesPrivilegesParams{
			HTTPRequest: suite.AuthenticateOfficeRequest(req, requestUser),
		}

		handler := GetRolesPrivilegesHandler{
			suite.NewHandlerConfig(),
			rolesservice.NewRolesFetcher(),
		}

		response := handler.Handle(params)

		suite.IsType(&officeuserop.GetRolesPrivilegesUnauthorized{}, response)
	})

	suite.Run("404 ERROR - Not Found ", func() {
		// Test:				GetOfficeUserHandler, Fetcher - Not Found
		// Set up:				Run request when logged in as admin user
		// Expected Outcome:	Not found response returned, no data
		params := officeuserop.GetRolesPrivilegesParams{
			HTTPRequest: suite.setupAuthenticatedRequest("GET", "/office_users/roles-privileges"),
		}

		mockFetcher := mocks.RoleAssociater{}
		mockFetcher.On("FetchRolesPrivileges", mock.AnythingOfType("*appcontext.appContext")).Return(nil, sql.ErrNoRows)

		handler := GetRolesPrivilegesHandler{
			suite.NewHandlerConfig(),
			&mockFetcher,
		}

		response := handler.Handle(params)

		suite.IsType(&officeuserop.GetRolesPrivilegesNotFound{}, response)
	})

	suite.Run("500 ERROR - Internal Server Error ", func() {
		// Test:				GetOfficeUserHandler, Fetcher - Internal Server Error
		// Set up:				Run request when logged in as admin user
		// Expected Outcome:	Internal Server Error response returned, no data
		params := officeuserop.GetRolesPrivilegesParams{
			HTTPRequest: suite.setupAuthenticatedRequest("GET", "/office_users/roles-privileges"),
		}

		mockFetcher := mocks.RoleAssociater{}
		mockFetcher.On("FetchRolesPrivileges", mock.AnythingOfType("*appcontext.appContext")).Return(nil, apperror.InternalServerError{})

		handler := GetRolesPrivilegesHandler{
			suite.NewHandlerConfig(),
			&mockFetcher,
		}

		response := handler.Handle(params)

		suite.IsType(&officeuserop.GetRolesPrivilegesInternalServerError{}, response)
	})
}
