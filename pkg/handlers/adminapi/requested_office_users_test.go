package adminapi

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/go-openapi/strfmt"
	"github.com/gofrs/uuid"
	"github.com/jarcoal/httpmock"
	"github.com/markbates/goth"
	"github.com/stretchr/testify/mock"

	"github.com/transcom/mymove/pkg/factory"
	requestedofficeuserop "github.com/transcom/mymove/pkg/gen/adminapi/adminoperations/requested_office_users"
	"github.com/transcom/mymove/pkg/gen/adminmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/handlers/authentication/okta"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/models/roles"
	notificationMocks "github.com/transcom/mymove/pkg/notifications/mocks"
	"github.com/transcom/mymove/pkg/services/mocks"
	"github.com/transcom/mymove/pkg/services/pagination"
	"github.com/transcom/mymove/pkg/services/query"
	requestedofficeusers "github.com/transcom/mymove/pkg/services/requested_office_users"
)

func (suite *HandlerSuite) TestIndexRequestedOfficeUsersHandler() {
	suite.Run("requested users result in ok response", func() {
		requestedOfficeUsers := models.OfficeUsers{
			factory.BuildOfficeUserWithRoles(suite.DB(), factory.GetTraitRequestedOfficeUser(), []roles.RoleType{roles.RoleTypeQae}),
			factory.BuildOfficeUserWithRoles(suite.DB(), factory.GetTraitRequestedOfficeUser(), []roles.RoleType{roles.RoleTypeQae})}
		params := requestedofficeuserop.IndexRequestedOfficeUsersParams{
			HTTPRequest: suite.setupAuthenticatedRequest("GET", "/requested_office_users"),
		}

		queryBuilder := query.NewQueryBuilder()
		handler := IndexRequestedOfficeUsersHandler{
			HandlerConfig:                  suite.HandlerConfig(),
			NewQueryFilter:                 query.NewQueryFilter,
			RequestedOfficeUserListFetcher: requestedofficeusers.NewRequestedOfficeUsersListFetcher(queryBuilder),
			NewPagination:                  pagination.NewPagination,
		}

		response := handler.Handle(params)

		suite.IsType(&requestedofficeuserop.IndexRequestedOfficeUsersOK{}, response)
		okResponse := response.(*requestedofficeuserop.IndexRequestedOfficeUsersOK)
		suite.Len(okResponse.Payload, 2)
		requestedOfficeUser1Id := requestedOfficeUsers[0].ID.String()
		requestedOfficeUser2Id := requestedOfficeUsers[1].ID.String()
		payloadRequestedUser1Id := okResponse.Payload[0].ID.String()
		payloadRequestedUser2Id := okResponse.Payload[1].ID.String()

		// requested office users should exist in response no matter the ordering that has been applied
		user1ExistsInResponse := requestedOfficeUser1Id == payloadRequestedUser1Id || requestedOfficeUser1Id == payloadRequestedUser2Id
		user2ExistsInResponse := requestedOfficeUser2Id == payloadRequestedUser1Id || requestedOfficeUser2Id == payloadRequestedUser2Id
		suite.True(user1ExistsInResponse)
		suite.True(user2ExistsInResponse)
	})

	suite.Run("able to search by name and filter", func() {
		status := models.OfficeUserStatusREQUESTED
		createdAt := time.Date(2007, 03, 05, 1, 1, 1, 1, time.Local)
		createdAt2 := time.Date(2006, 03, 07, 1, 1, 1, 1, time.Local)
		transportationOffice := factory.BuildTransportationOffice(suite.DB(), []factory.Customization{
			{
				Model: models.TransportationOffice{
					Name: "JPPO Test Office",
				},
			},
		}, []factory.Trait{})
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
					TransportationOfficeID: transportationOffice2.ID,
					CreatedAt:              createdAt2,
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
					TransportationOfficeID: transportationOffice2.ID,
					CreatedAt:              createdAt2,
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
					TransportationOfficeID: transportationOffice2.ID,
					CreatedAt:              createdAt2,
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
					TransportationOfficeID: transportationOffice.ID,
					CreatedAt:              createdAt,
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
		params := requestedofficeuserop.IndexRequestedOfficeUsersParams{
			HTTPRequest: suite.setupAuthenticatedRequest("GET", "/requested_office_users"),
			Filter:      &filterJSON,
		}

		queryBuilder := query.NewQueryBuilder()
		handler := IndexRequestedOfficeUsersHandler{
			HandlerConfig:                  suite.HandlerConfig(),
			NewQueryFilter:                 query.NewQueryFilter,
			RequestedOfficeUserListFetcher: requestedofficeusers.NewRequestedOfficeUsersListFetcher(queryBuilder),
			NewPagination:                  pagination.NewPagination,
		}

		response := handler.Handle(params)

		suite.IsType(&requestedofficeuserop.IndexRequestedOfficeUsersOK{}, response)
		okResponse := response.(*requestedofficeuserop.IndexRequestedOfficeUsersOK)
		suite.Len(okResponse.Payload, 2)
		suite.Contains(*okResponse.Payload[0].FirstName, nameSearch)
		suite.Contains(*okResponse.Payload[1].FirstName, nameSearch)

		// email search
		emailSearch := "AirKilluh2"
		filterJSON = fmt.Sprintf("{\"email\":\"%s\"}", emailSearch)
		params = requestedofficeuserop.IndexRequestedOfficeUsersParams{
			HTTPRequest: suite.setupAuthenticatedRequest("GET", "/requested_office_users"),
			Filter:      &filterJSON,
		}
		response = handler.Handle(params)

		suite.IsType(&requestedofficeuserop.IndexRequestedOfficeUsersOK{}, response)
		okResponse = response.(*requestedofficeuserop.IndexRequestedOfficeUsersOK)
		suite.Len(okResponse.Payload, 1)
		suite.Contains(*okResponse.Payload[0].Email, emailSearch)

		// firstName search
		firstSearch := "Angel"
		filterJSON = fmt.Sprintf("{\"firstName\":\"%s\"}", firstSearch)
		params = requestedofficeuserop.IndexRequestedOfficeUsersParams{
			HTTPRequest: suite.setupAuthenticatedRequest("GET", "/requested_office_users"),
			Filter:      &filterJSON,
		}
		response = handler.Handle(params)

		suite.IsType(&requestedofficeuserop.IndexRequestedOfficeUsersOK{}, response)
		okResponse = response.(*requestedofficeuserop.IndexRequestedOfficeUsersOK)
		suite.Len(okResponse.Payload, 1)
		suite.Contains(*okResponse.Payload[0].FirstName, firstSearch)

		// lastName search
		lastSearch := "Jo"
		filterJSON = fmt.Sprintf("{\"lastName\":\"%s\"}", lastSearch)
		params = requestedofficeuserop.IndexRequestedOfficeUsersParams{
			HTTPRequest: suite.setupAuthenticatedRequest("GET", "/requested_office_users"),
			Filter:      &filterJSON,
		}
		response = handler.Handle(params)

		suite.IsType(&requestedofficeuserop.IndexRequestedOfficeUsersOK{}, response)
		okResponse = response.(*requestedofficeuserop.IndexRequestedOfficeUsersOK)
		suite.Len(okResponse.Payload, 1)
		suite.Contains(*okResponse.Payload[0].LastName, lastSearch)

		// transportation office search
		filterJSON = "{\"office\":\"JP\"}"
		params = requestedofficeuserop.IndexRequestedOfficeUsersParams{
			HTTPRequest: suite.setupAuthenticatedRequest("GET", "/requested_office_users"),
			Filter:      &filterJSON,
		}
		response = handler.Handle(params)

		suite.IsType(&requestedofficeuserop.IndexRequestedOfficeUsersOK{}, response)
		okResponse = response.(*requestedofficeuserop.IndexRequestedOfficeUsersOK)
		suite.Len(okResponse.Payload, 1)
		suite.Equal(strfmt.UUID(transportationOffice.ID.String()), *okResponse.Payload[0].TransportationOfficeID)

		// requestedOn search
		requestedOnSearch := "5"
		filterJSON = fmt.Sprintf("{\"requestedOn\":\"%s\"}", requestedOnSearch)
		params = requestedofficeuserop.IndexRequestedOfficeUsersParams{
			HTTPRequest: suite.setupAuthenticatedRequest("GET", "/requested_office_users"),
			Filter:      &filterJSON,
		}
		response = handler.Handle(params)

		suite.IsType(&requestedofficeuserop.IndexRequestedOfficeUsersOK{}, response)
		okResponse = response.(*requestedofficeuserop.IndexRequestedOfficeUsersOK)
		suite.Equal(1, len(okResponse.Payload))
		suite.Contains(okResponse.Payload[0].CreatedAt.String(), requestedOnSearch)

		// roles search
		roleSearch := "Counselor"
		filterJSON = fmt.Sprintf("{\"roles\":\"%s\"}", roleSearch)
		params = requestedofficeuserop.IndexRequestedOfficeUsersParams{
			HTTPRequest: suite.setupAuthenticatedRequest("GET", "/requested_office_users"),
			Filter:      &filterJSON,
		}
		response = handler.Handle(params)

		suite.IsType(&requestedofficeuserop.IndexRequestedOfficeUsersOK{}, response)
		okResponse = response.(*requestedofficeuserop.IndexRequestedOfficeUsersOK)
		suite.Len(okResponse.Payload, 2)
		suite.Contains(*okResponse.Payload[0].Roles[0].RoleName, roleSearch)
		suite.Contains(*okResponse.Payload[1].Roles[0].RoleName, roleSearch)

	})

	suite.Run("test the return of sorted requested office users in asc order", func() {
		requestedStatus := models.OfficeUserStatusREQUESTED
		officeUser1 := factory.BuildOfficeUserWithRoles(suite.DB(), []factory.Customization{
			{
				Model: models.OfficeUser{
					FirstName: "Angelina",
					LastName:  "Jolie",
					Email:     "laraCroft@mail.mil",
					Status:    &requestedStatus,
				},
			},
			{
				Model: models.TransportationOffice{
					Name: "PPPO Kirtland AFB - USAF",
				},
			},
		}, []roles.RoleType{roles.RoleTypeTOO})
		officeUser2 := factory.BuildOfficeUserWithRoles(suite.DB(), []factory.Customization{
			{
				Model: models.OfficeUser{
					FirstName: "Billy",
					LastName:  "Bob",
					Email:     "bigBob@mail.mil",
					Status:    &requestedStatus,
				},
			},
			{
				Model: models.TransportationOffice{
					Name: "PPPO Fort Knox - USA",
				},
			},
		}, []roles.RoleType{roles.RoleTypeTIO})
		officeUser3 := factory.BuildOfficeUserWithRoles(suite.DB(), []factory.Customization{
			{
				Model: models.OfficeUser{
					FirstName: "Nick",
					LastName:  "Cage",
					Email:     "conAirKilluh@mail.mil",
					Status:    &requestedStatus,
				},
			},
			{
				Model: models.TransportationOffice{
					Name: "PPPO Detroit Arsenal - USA",
				},
			},
		}, []roles.RoleType{roles.RoleTypeServicesCounselor})

		sortColumn := "transportation_office_id"
		params := requestedofficeuserop.IndexRequestedOfficeUsersParams{
			HTTPRequest: suite.setupAuthenticatedRequest("GET", "/requested_office_users"),
			Sort:        &sortColumn,
			Order:       models.BoolPointer(true),
		}

		queryBuilder := query.NewQueryBuilder()
		handler := IndexRequestedOfficeUsersHandler{
			HandlerConfig:                  suite.HandlerConfig(),
			NewQueryFilter:                 query.NewQueryFilter,
			RequestedOfficeUserListFetcher: requestedofficeusers.NewRequestedOfficeUsersListFetcher(queryBuilder),
			NewPagination:                  pagination.NewPagination,
		}

		response := handler.Handle(params)

		suite.IsType(&requestedofficeuserop.IndexRequestedOfficeUsersOK{}, response)
		okResponse := response.(*requestedofficeuserop.IndexRequestedOfficeUsersOK)
		suite.Len(okResponse.Payload, 3)
		suite.Equal(officeUser3.ID.String(), okResponse.Payload[0].ID.String())
		suite.Equal(officeUser2.ID.String(), okResponse.Payload[1].ID.String())
		suite.Equal(officeUser1.ID.String(), okResponse.Payload[2].ID.String())

		// sort by transportation office name in desc order
		sortColumn = "transportation_office_id"
		params = requestedofficeuserop.IndexRequestedOfficeUsersParams{
			HTTPRequest: suite.setupAuthenticatedRequest("GET", "/requested_office_users"),
			Sort:        &sortColumn,
			Order:       models.BoolPointer(false),
		}

		queryBuilder = query.NewQueryBuilder()
		handler = IndexRequestedOfficeUsersHandler{
			HandlerConfig:                  suite.HandlerConfig(),
			NewQueryFilter:                 query.NewQueryFilter,
			RequestedOfficeUserListFetcher: requestedofficeusers.NewRequestedOfficeUsersListFetcher(queryBuilder),
			NewPagination:                  pagination.NewPagination,
		}

		response = handler.Handle(params)

		suite.IsType(&requestedofficeuserop.IndexRequestedOfficeUsersOK{}, response)
		okResponse = response.(*requestedofficeuserop.IndexRequestedOfficeUsersOK)
		suite.Len(okResponse.Payload, 3)
		suite.Equal(officeUser1.ID.String(), okResponse.Payload[0].ID.String())
		suite.Equal(officeUser2.ID.String(), okResponse.Payload[1].ID.String())
		suite.Equal(officeUser3.ID.String(), okResponse.Payload[2].ID.String())

		// sort by first name in asc order
		sortColumn = "first_name"
		params = requestedofficeuserop.IndexRequestedOfficeUsersParams{
			HTTPRequest: suite.setupAuthenticatedRequest("GET", "/requested_office_users"),
			Sort:        &sortColumn,
			Order:       models.BoolPointer(true),
		}

		queryBuilder = query.NewQueryBuilder()
		handler = IndexRequestedOfficeUsersHandler{
			HandlerConfig:                  suite.HandlerConfig(),
			NewQueryFilter:                 query.NewQueryFilter,
			RequestedOfficeUserListFetcher: requestedofficeusers.NewRequestedOfficeUsersListFetcher(queryBuilder),
			NewPagination:                  pagination.NewPagination,
		}

		response = handler.Handle(params)

		suite.IsType(&requestedofficeuserop.IndexRequestedOfficeUsersOK{}, response)
		okResponse = response.(*requestedofficeuserop.IndexRequestedOfficeUsersOK)
		suite.Len(okResponse.Payload, 3)
		suite.Equal(officeUser1.ID.String(), okResponse.Payload[0].ID.String())
		suite.Equal(officeUser2.ID.String(), okResponse.Payload[1].ID.String())
		suite.Equal(officeUser3.ID.String(), okResponse.Payload[2].ID.String())
	})

	suite.Run("able to search by transportation office", func() {
		transportationOffice := factory.BuildTransportationOffice(suite.DB(), []factory.Customization{
			{
				Model: models.TransportationOffice{
					Name: "Tinker",
				},
			},
		}, nil)
		requestedStatus := models.OfficeUserStatusREQUESTED
		officeUser1 := factory.BuildOfficeUserWithRoles(suite.DB(), []factory.Customization{
			{
				Model: models.OfficeUser{
					TransportationOfficeID: transportationOffice.ID,
					Status:                 &requestedStatus,
				},
			},
			{
				Model:    transportationOffice,
				LinkOnly: true,
				Type:     &factory.TransportationOffices.CounselingOffice,
			},
		}, []roles.RoleType{roles.RoleTypeServicesCounselor})
		factory.BuildOfficeUserWithRoles(suite.DB(), factory.GetTraitRequestedOfficeUser(), []roles.RoleType{roles.RoleTypeTOO})
		factory.BuildOfficeUserWithRoles(suite.DB(), factory.GetTraitRequestedOfficeUser(), []roles.RoleType{roles.RoleTypeTIO})

		filterJSON := "{\"office\":\"Tinker\"}"
		params := requestedofficeuserop.IndexRequestedOfficeUsersParams{
			HTTPRequest: suite.setupAuthenticatedRequest("GET", "/requested_office_users"),
			Filter:      &filterJSON,
		}

		queryBuilder := query.NewQueryBuilder()
		handler := IndexRequestedOfficeUsersHandler{
			HandlerConfig:                  suite.HandlerConfig(),
			NewQueryFilter:                 query.NewQueryFilter,
			RequestedOfficeUserListFetcher: requestedofficeusers.NewRequestedOfficeUsersListFetcher(queryBuilder),
			NewPagination:                  pagination.NewPagination,
		}

		response := handler.Handle(params)

		suite.IsType(&requestedofficeuserop.IndexRequestedOfficeUsersOK{}, response)
		okResponse := response.(*requestedofficeuserop.IndexRequestedOfficeUsersOK)
		suite.Len(okResponse.Payload, 1)
		suite.Equal(officeUser1.ID.String(), okResponse.Payload[0].ID.String())
	})

	suite.Run("able to search by role", func() {
		requestedStatus := models.OfficeUserStatusREQUESTED
		officeUser1 := factory.BuildOfficeUserWithRoles(suite.DB(), []factory.Customization{
			{
				Model: models.OfficeUser{
					Status: &requestedStatus,
				},
			},
		}, []roles.RoleType{roles.RoleTypeServicesCounselor})

		filterJSON := "{\"rolesSearch\":\"services\"}"
		params := requestedofficeuserop.IndexRequestedOfficeUsersParams{
			HTTPRequest: suite.setupAuthenticatedRequest("GET", "/requested_office_users"),
			Filter:      &filterJSON,
		}

		queryBuilder := query.NewQueryBuilder()
		handler := IndexRequestedOfficeUsersHandler{
			HandlerConfig:                  suite.HandlerConfig(),
			NewQueryFilter:                 query.NewQueryFilter,
			RequestedOfficeUserListFetcher: requestedofficeusers.NewRequestedOfficeUsersListFetcher(queryBuilder),
			NewPagination:                  pagination.NewPagination,
		}

		response := handler.Handle(params)

		suite.IsType(&requestedofficeuserop.IndexRequestedOfficeUsersOK{}, response)
		okResponse := response.(*requestedofficeuserop.IndexRequestedOfficeUsersOK)
		suite.Len(okResponse.Payload, 1)
		suite.Equal(officeUser1.ID.String(), okResponse.Payload[0].ID.String())
	})

	suite.Run("return error when querying for unhandled data", func() {
		requestedStatus := models.OfficeUserStatusREQUESTED
		factory.BuildOfficeUserWithRoles(suite.DB(), []factory.Customization{
			{
				Model: models.OfficeUser{
					Status: &requestedStatus,
				},
			},
		}, []roles.RoleType{roles.RoleTypeServicesCounselor})

		sortColumn := "unknown_column"
		params := requestedofficeuserop.IndexRequestedOfficeUsersParams{
			HTTPRequest: suite.setupAuthenticatedRequest("GET", "/requested_office_users"),
			Sort:        &sortColumn,
			Order:       models.BoolPointer(true),
		}

		queryBuilder := query.NewQueryBuilder()
		handler := IndexRequestedOfficeUsersHandler{
			HandlerConfig:                  suite.HandlerConfig(),
			NewQueryFilter:                 query.NewQueryFilter,
			RequestedOfficeUserListFetcher: requestedofficeusers.NewRequestedOfficeUsersListFetcher(queryBuilder),
			NewPagination:                  pagination.NewPagination,
		}

		response := handler.Handle(params)

		suite.IsType(&handlers.ErrResponse{}, response)
		errResponse := response.(*handlers.ErrResponse)
		suite.Equal(http.StatusInternalServerError, errResponse.Code)
		errMsg := errResponse.Err.Error()
		suite.Equal(errMsg, "Unhandled data error encountered")
	})

	suite.Run("should error when a param filter format is incorrect", func() {
		requestedStatus := models.OfficeUserStatusREQUESTED
		factory.BuildOfficeUserWithRoles(suite.DB(), []factory.Customization{
			{
				Model: models.OfficeUser{
					Status: &requestedStatus,
				},
			},
		}, []roles.RoleType{roles.RoleTypeServicesCounselor})

		// Invalid format for filter params
		filterJSON := "test{\"unknown\":\"value\"}test"
		params := requestedofficeuserop.IndexRequestedOfficeUsersParams{
			HTTPRequest: suite.setupAuthenticatedRequest("GET", "/requested_office_users"),
			Filter:      &filterJSON,
		}

		queryBuilder := query.NewQueryBuilder()
		handler := IndexRequestedOfficeUsersHandler{
			HandlerConfig:                  suite.HandlerConfig(),
			NewQueryFilter:                 query.NewQueryFilter,
			RequestedOfficeUserListFetcher: requestedofficeusers.NewRequestedOfficeUsersListFetcher(queryBuilder),
			NewPagination:                  pagination.NewPagination,
		}

		response := handler.Handle(params)

		expectedError := models.ErrInvalidFilterFormat
		expectedResponse := &handlers.ErrResponse{
			Code: http.StatusInternalServerError,
			Err:  expectedError,
		}

		suite.Equal(expectedResponse, response)
	})
}

func (suite *HandlerSuite) TestGetRequestedOfficeUserHandler() {
	suite.Run("integration test ok response", func() {
		requestedOfficeUser := factory.BuildOfficeUserWithRoles(suite.DB(), factory.GetTraitRequestedOfficeUser(), []roles.RoleType{roles.RoleTypeQae})
		params := requestedofficeuserop.GetRequestedOfficeUserParams{
			HTTPRequest:  suite.setupAuthenticatedRequest("GET", fmt.Sprintf("/requested_office_users/%s", requestedOfficeUser.ID)),
			OfficeUserID: strfmt.UUID(requestedOfficeUser.ID.String()),
		}

		mockRoleAssociator := &mocks.RoleAssociator{}
		mockPrivilegesAssociator := &mocks.UserPrivilegeAssociator{}
		mockRoles := roles.Roles{
			roles.Role{
				ID:        uuid.Must(uuid.NewV4()),
				RoleType:  roles.RoleTypeTOO,
				RoleName:  "Task Ordering Officer",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
		}
		mockRoleAssociator.On(
			"FetchRolesForUser",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything,
		).Return(mockRoles, nil)

		queryBuilder := query.NewQueryBuilder()
		handler := GetRequestedOfficeUserHandler{
			suite.HandlerConfig(),
			requestedofficeusers.NewRequestedOfficeUserFetcher(queryBuilder),
			mockRoleAssociator,
			mockPrivilegesAssociator,
			query.NewQueryFilter,
		}

		response := handler.Handle(params)

		suite.IsType(&requestedofficeuserop.GetRequestedOfficeUserOK{}, response)
		okResponse := response.(*requestedofficeuserop.GetRequestedOfficeUserOK)
		suite.Equal(requestedOfficeUser.ID.String(), okResponse.Payload.ID.String())
	})

	suite.Run("successful response", func() {
		requestedOfficeUser := factory.BuildOfficeUserWithRoles(suite.DB(), factory.GetTraitRequestedOfficeUser(), []roles.RoleType{roles.RoleTypeQae})
		params := requestedofficeuserop.GetRequestedOfficeUserParams{
			HTTPRequest:  suite.setupAuthenticatedRequest("GET", fmt.Sprintf("/requested_office_users/%s", requestedOfficeUser.ID)),
			OfficeUserID: strfmt.UUID(requestedOfficeUser.ID.String()),
		}

		requestedOfficeUserFetcher := &mocks.RequestedOfficeUserFetcher{}
		requestedOfficeUserFetcher.On("FetchRequestedOfficeUser",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything,
		).Return(requestedOfficeUser, nil).Once()

		mockRoleAssociator := &mocks.RoleAssociator{}
		userPrivilegeAssociator := &mocks.UserPrivilegeAssociator{}
		mockRoles := roles.Roles{
			roles.Role{
				ID:        uuid.Must(uuid.NewV4()),
				RoleType:  roles.RoleTypeTOO,
				RoleName:  "Task Ordering Officer",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
		}

		mockRoleAssociator.On(
			"FetchRolesForUser",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything,
		).Return(mockRoles, nil)

		handler := GetRequestedOfficeUserHandler{
			suite.HandlerConfig(),
			requestedOfficeUserFetcher,
			mockRoleAssociator,
			userPrivilegeAssociator,
			newMockQueryFilterBuilder(&mocks.QueryFilter{}),
		}

		response := handler.Handle(params)

		suite.IsType(&requestedofficeuserop.GetRequestedOfficeUserOK{}, response)
		okResponse := response.(*requestedofficeuserop.GetRequestedOfficeUserOK)
		suite.Equal(requestedOfficeUser.ID.String(), okResponse.Payload.ID.String())
	})

	suite.Run("unsuccessful response when fetch fails", func() {
		requestedOfficeUser := factory.BuildOfficeUserWithRoles(suite.DB(), factory.GetTraitRequestedOfficeUser(), []roles.RoleType{roles.RoleTypeQae})
		params := requestedofficeuserop.GetRequestedOfficeUserParams{
			HTTPRequest:  suite.setupAuthenticatedRequest("GET", fmt.Sprintf("/requested_office_users/%s", requestedOfficeUser.ID)),
			OfficeUserID: strfmt.UUID(requestedOfficeUser.ID.String()),
		}

		expectedError := models.ErrFetchNotFound
		requestedOfficeUserFetcher := &mocks.RequestedOfficeUserFetcher{}
		requestedOfficeUserFetcher.On("FetchRequestedOfficeUser",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything,
		).Return(models.OfficeUser{}, expectedError).Once()

		mockRoleAssociator := &mocks.RoleAssociator{}
		userPrivilegeAssociator := &mocks.UserPrivilegeAssociator{}
		mockRoles := roles.Roles{
			roles.Role{
				ID:        uuid.Must(uuid.NewV4()),
				RoleType:  roles.RoleTypeTOO,
				RoleName:  "Task Ordering Officer",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
		}
		mockRoleAssociator.On(
			"FetchRolesForUser",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything,
		).Return(mockRoles, nil)

		mockRoleAssociator.On(
			"FetchRolesForUser",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything,
		).Return(mockRoles, nil)

		handler := GetRequestedOfficeUserHandler{
			suite.HandlerConfig(),
			requestedOfficeUserFetcher,
			mockRoleAssociator,
			userPrivilegeAssociator,
			newMockQueryFilterBuilder(&mocks.QueryFilter{}),
		}

		response := handler.Handle(params)

		expectedResponse := &handlers.ErrResponse{
			Code: http.StatusNotFound,
			Err:  expectedError,
		}
		suite.Equal(expectedResponse, response)
	})

	suite.Run("test - get requested office user handler with privileges", func() {
		requestedOfficeUser := factory.BuildOfficeUserWithPrivileges(suite.DB(), []factory.Customization{
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

		params := requestedofficeuserop.GetRequestedOfficeUserParams{
			HTTPRequest:  suite.setupAuthenticatedRequest("GET", fmt.Sprintf("/requested_office_users/%s", requestedOfficeUser.ID)),
			OfficeUserID: strfmt.UUID(requestedOfficeUser.ID.String()),
		}

		requestedOfficeUserFetcher := &mocks.RequestedOfficeUserFetcher{}
		requestedOfficeUserFetcher.On("FetchRequestedOfficeUser",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything,
		).Return(requestedOfficeUser, nil).Once()

		mockRoleAssociator := &mocks.RoleAssociater{}
		userPrivilegeAssociator := &mocks.UserPrivilegeAssociator{}

		mockPrivileges := []roles.Privilege{
			{
				ID:            uuid.Must(uuid.NewV4()),
				PrivilegeType: roles.PrivilegeTypeSupervisor,
				PrivilegeName: "Supervisor",
				CreatedAt:     time.Now(),
				UpdatedAt:     time.Now(),
			},
		}

		mockRoles := roles.Roles{
			roles.Role{
				ID:        uuid.Must(uuid.NewV4()),
				RoleType:  roles.RoleTypeTOO,
				RoleName:  "Task Ordering Officer",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
		}
		mockRoleAssociator.On(
			"FetchRolesForUser",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything,
		).Return(mockRoles, nil)

		mockRoleAssociator.On(
			"FetchRolesForUser",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything,
		).Return(mockRoles, nil)

		mockRoleAssociator.On(
			"FetchPrivilegesForUser",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything,
		).Return(mockPrivileges, nil)

		userPrivilegeAssociator.On(
			"FetchPrivilegesForUser",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything,
		).Return(mockPrivileges, nil)

		handler := GetRequestedOfficeUserHandler{
			suite.HandlerConfig(),
			requestedOfficeUserFetcher,
			mockRoleAssociator,
			userPrivilegeAssociator,
			newMockQueryFilterBuilder(&mocks.QueryFilter{}),
		}

		response := handler.Handle(params)
		suite.IsType(&requestedofficeuserop.GetRequestedOfficeUserOK{}, response)
		okResponse := response.(*requestedofficeuserop.GetRequestedOfficeUserOK)
		suite.Equal(requestedOfficeUser.ID.String(), okResponse.Payload.ID.String())
		suite.Len(okResponse.Payload.Privileges, 1)
		suite.Equal("Supervisor", okResponse.Payload.Privileges[0].PrivilegeName)
		suite.Equal(string(roles.PrivilegeTypeSupervisor), okResponse.Payload.Privileges[0].PrivilegeType)
	})

}

func (suite *HandlerSuite) TestGetRequestedOfficeUserHandler_WithSupervisorPrivilege() {
	suite.Run("test - get requested office user handler with privileges", func() {
		requestedOfficeUser := factory.BuildOfficeUserWithPrivileges(suite.DB(), []factory.Customization{
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

		params := requestedofficeuserop.GetRequestedOfficeUserParams{
			HTTPRequest:  suite.setupAuthenticatedRequest("GET", fmt.Sprintf("/requested_office_users/%s", requestedOfficeUser.ID)),
			OfficeUserID: strfmt.UUID(requestedOfficeUser.ID.String()),
		}

		requestedOfficeUserFetcher := &mocks.RequestedOfficeUserFetcher{}
		requestedOfficeUserFetcher.On("FetchRequestedOfficeUser",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything,
		).Return(requestedOfficeUser, nil).Once()

		mockRoleAssociator := &mocks.RoleAssociater{}
		userPrivilegeAssociator := &mocks.UserPrivilegeAssociator{}

		mockPrivileges := []roles.Privilege{
			{
				ID:            uuid.Must(uuid.NewV4()),
				PrivilegeType: roles.PrivilegeTypeSupervisor,
				PrivilegeName: "Supervisor",
				CreatedAt:     time.Now(),
				UpdatedAt:     time.Now(),
			},
		}

		mockRoles := roles.Roles{
			roles.Role{
				ID:        uuid.Must(uuid.NewV4()),
				RoleType:  roles.RoleTypeTOO,
				RoleName:  "Task Ordering Officer",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
		}
		mockRoleAssociator.On(
			"FetchRolesForUser",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything,
		).Return(mockRoles, nil)

		mockRoleAssociator.On(
			"FetchRolesForUser",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything,
		).Return(mockRoles, nil)

		mockRoleAssociator.On(
			"FetchPrivilegesForUser",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything,
		).Return(mockPrivileges, nil)

		userPrivilegeAssociator.On(
			"FetchPrivilegesForUser",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything,
		).Return(mockPrivileges, nil)

		handler := GetRequestedOfficeUserHandler{
			suite.HandlerConfig(),
			requestedOfficeUserFetcher,
			mockRoleAssociator,
			userPrivilegeAssociator,
			newMockQueryFilterBuilder(&mocks.QueryFilter{}),
		}

		response := handler.Handle(params)
		suite.IsType(&requestedofficeuserop.GetRequestedOfficeUserOK{}, response)
		okResponse := response.(*requestedofficeuserop.GetRequestedOfficeUserOK)
		suite.Equal(requestedOfficeUser.ID.String(), okResponse.Payload.ID.String())
		suite.Len(okResponse.Payload.Privileges, 1)
		suite.Equal("Supervisor", okResponse.Payload.Privileges[0].PrivilegeName)
		suite.Equal(string(roles.PrivilegeTypeSupervisor), okResponse.Payload.Privileges[0].PrivilegeType)
	})

}

func (suite *HandlerSuite) TestUpdateRequestedOfficeUserHandler_WithSupervisorPrivilege() {
	suite.Run("test - update requested office user handler with privileges", func() {
		tooRoleName := "Task Ordering Officer"
		tooRoleType := string(roles.RoleTypeTOO)
		supervisorPrivilegeName := "Supervisor"
		supervisorPrivilegeType := string(roles.PrivilegeTypeSupervisor)
		requestedOfficeUser := factory.BuildOfficeUserWithPrivileges(suite.DB(), []factory.Customization{
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

		officeUser := models.OfficeUser{ID: requestedOfficeUser.ID, FirstName: "Billy", LastName: "Bob", UserID: requestedOfficeUser.UserID, CreatedAt: time.Now(),
			UpdatedAt: time.Now()}

		mockUserRoleAssociator := &mocks.UserRoleAssociator{}
		mockRoleAssociator := &mocks.RoleAssociater{}
		requestedOfficeUserFetcher := &mocks.RequestedOfficeUserFetcher{}
		requestedOfficeUserUpdater := &mocks.RequestedOfficeUserUpdater{}
		userPrivilegeAssociator := &mocks.UserPrivilegeAssociator{}
		privilegeFetcher := &mocks.PrivilegeAssociator{}

		params := requestedofficeuserop.UpdateRequestedOfficeUserParams{
			HTTPRequest: suite.setupAuthenticatedRequest("PATCH", fmt.Sprintf("/requested_office_users/%s", requestedOfficeUser.ID)),
			Body: &adminmessages.RequestedOfficeUserUpdate{

				Roles: []*adminmessages.OfficeUserRole{
					{
						Name:     &tooRoleName,
						RoleType: &tooRoleType,
					},
				},
				Privileges: []*adminmessages.OfficeUserPrivilege{
					{
						Name:          &supervisorPrivilegeName,
						PrivilegeType: &supervisorPrivilegeType,
					},
				},
			},
			OfficeUserID: strfmt.UUID(requestedOfficeUser.ID.String()),
		}

		requestedOfficeUserFetcher.On("FetchRequestedOfficeUser",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything,
		).Return(requestedOfficeUser, nil, nil).Once()

		requestedOfficeUserUpdater.On("UpdateRequestedOfficeUser",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything,
			mock.Anything,
		).Return(&officeUser, nil, nil).Once()

		mockRoles := roles.Roles{
			roles.Role{
				ID:        uuid.Must(uuid.NewV4()),
				RoleType:  roles.RoleTypeTOO,
				RoleName:  "Task Ordering Officer",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
		}

		mockPrivileges := roles.Privileges{
			{
				ID:            uuid.Must(uuid.NewV4()),
				PrivilegeType: roles.PrivilegeTypeSupervisor,
				PrivilegeName: "Supervisor",
				CreatedAt:     time.Now(),
				UpdatedAt:     time.Now(),
			},
		}

		mockUserRoleAssociator.On(
			"UpdateUserRoles",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything,
			mock.Anything,
		).Return(nil, nil, nil).Once()

		userPrivilegeAssociator.On(
			"UpdateUserPrivileges",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything,
			mock.Anything,
		).Return(nil, nil, nil).Once()

		privilegeFetcher.On(
			"FetchPrivilegesForUser",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything,
		).Return(mockPrivileges, nil)

		mockRoleAssociator.On(
			"FetchRolesForUser",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything,
		).Return(mockRoles, nil)

		handler := UpdateRequestedOfficeUserHandler{
			suite.HandlerConfig(),
			requestedOfficeUserFetcher,
			requestedOfficeUserUpdater,
			userPrivilegeAssociator,
			privilegeFetcher,
			mockUserRoleAssociator,
			mockRoleAssociator,
		}

		response := handler.Handle(params)
		suite.IsType(&requestedofficeuserop.UpdateRequestedOfficeUserOK{}, response)
		okResponse := response.(*requestedofficeuserop.UpdateRequestedOfficeUserOK)
		suite.Equal(officeUser.ID.String(), okResponse.Payload.ID.String())
		suite.Len(okResponse.Payload.Privileges, 1)
		suite.Equal(supervisorPrivilegeName, okResponse.Payload.Privileges[0].PrivilegeName)
		suite.Equal(supervisorPrivilegeType, okResponse.Payload.Privileges[0].PrivilegeType)
	})

	suite.Run("test - update requested office user handler - REJECT supervisor privilege", func() {
		tooRoleName := "Task Ordering Officer"
		tooRoleType := string(roles.RoleTypeTOO)
		requestedOfficeUser := factory.BuildOfficeUserWithPrivileges(suite.DB(), []factory.Customization{
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

		officeUser := models.OfficeUser{ID: requestedOfficeUser.ID, FirstName: "Billy", LastName: "Bob", UserID: requestedOfficeUser.UserID, CreatedAt: time.Now(),
			UpdatedAt: time.Now()}

		mockUserRoleAssociator := &mocks.UserRoleAssociator{}
		mockRoleAssociator := &mocks.RoleAssociater{}
		requestedOfficeUserUpdater := &mocks.RequestedOfficeUserUpdater{}
		requestedOfficeUserFetcher := &mocks.RequestedOfficeUserFetcher{}
		userPrivilegeAssociator := &mocks.UserPrivilegeAssociator{}
		privilegeFetcher := &mocks.PrivilegeAssociator{}

		params := requestedofficeuserop.UpdateRequestedOfficeUserParams{
			HTTPRequest: suite.setupAuthenticatedRequest("PATCH", fmt.Sprintf("/requested_office_users/%s", requestedOfficeUser.ID)),
			Body: &adminmessages.RequestedOfficeUserUpdate{
				Roles: []*adminmessages.OfficeUserRole{
					{
						Name:     &tooRoleName,
						RoleType: &tooRoleType,
					},
				},
				Privileges: []*adminmessages.OfficeUserPrivilege{},
			},
			OfficeUserID: strfmt.UUID(requestedOfficeUser.ID.String()),
		}

		requestedOfficeUserFetcher.On("FetchRequestedOfficeUser",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything,
		).Return(requestedOfficeUser, nil, nil).Once()

		requestedOfficeUserUpdater.On("UpdateRequestedOfficeUser",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything,
			mock.Anything,
		).Return(&officeUser, nil, nil).Once()

		mockRoles := roles.Roles{
			roles.Role{
				ID:        uuid.Must(uuid.NewV4()),
				RoleType:  roles.RoleTypeTOO,
				RoleName:  "Task Ordering Officer",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
		}
		mockPrivileges := roles.Privileges{}

		mockUserRoleAssociator.On(
			"UpdateUserRoles",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything,
			mock.Anything,
		).Return(nil, nil, nil).Once()

		userPrivilegeAssociator.On(
			"UpdateUserPrivileges",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything,
			mock.Anything,
		).Return(nil, nil, nil).Once()

		privilegeFetcher.On(
			"FetchPrivilegesForUser",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything,
		).Return(mockPrivileges, nil)

		mockRoleAssociator.On(
			"FetchRolesForUser",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything,
		).Return(mockRoles, nil)

		notificationSender := &notificationMocks.NotificationSender{}

		notificationSender.On("SendNotification", mock.AnythingOfType("*appcontext.appContext"), mock.Anything).Return(nil)

		suiteHandler := suite.HandlerConfig()
		suiteHandler.SetNotificationSender(notificationSender)

		handler := UpdateRequestedOfficeUserHandler{
			suiteHandler,
			requestedOfficeUserFetcher,
			requestedOfficeUserUpdater,
			userPrivilegeAssociator,
			privilegeFetcher,
			mockUserRoleAssociator,
			mockRoleAssociator,
		}

		response := handler.Handle(params)
		suite.IsType(&requestedofficeuserop.UpdateRequestedOfficeUserOK{}, response)
		okResponse := response.(*requestedofficeuserop.UpdateRequestedOfficeUserOK)
		suite.Equal(officeUser.ID.String(), okResponse.Payload.ID.String())
		suite.Empty(okResponse.Payload.Privileges)
		notificationSender.AssertNumberOfCalls(suite.T(), "SendNotification", 1)
	})
}

func (suite *HandlerSuite) TestUpdateRequestedOfficeUserHandlerWithoutOktaAccountCreation() {
	suite.Run("Successful update", func() {
		user := factory.BuildDefaultUser(suite.DB())
		tooRoleName := "Task Ordering Officer"
		tooRoleType := string(roles.RoleTypeTOO)
		tioRoleName := "Task Invoicing Officer"
		tioRoleType := string(roles.RoleTypeTIO)
		requestedOfficeUser := factory.BuildOfficeUserWithRoles(suite.DB(), []factory.Customization{
			{
				Model: models.OfficeUser{
					Active: true,
					UserID: &user.ID,
				},
			},
			{
				Model: models.User{
					Roles: roles.Roles{
						{RoleName: roles.RoleName(tioRoleName),
							RoleType: roles.RoleType(tioRoleType)},
					},
				},
			},
		}, []roles.RoleType{roles.RoleTypeTOO})

		requestedOfficeUser.User.Privileges = []roles.Privilege{
			{
				ID:            uuid.Must(uuid.NewV4()),
				PrivilegeType: roles.PrivilegeTypeSupervisor,
				PrivilegeName: "Supervisor",
				CreatedAt:     time.Now(),
				UpdatedAt:     time.Now(),
			},
		}

		officeUserID := requestedOfficeUser.ID
		officeUser := models.OfficeUser{ID: officeUserID, FirstName: "Billy", LastName: "Bob", UserID: requestedOfficeUser.UserID, CreatedAt: time.Now(),
			UpdatedAt: time.Now()}

		mockUserRoleAssociator := &mocks.UserRoleAssociator{}
		mockRoleAssociator := &mocks.RoleAssociator{}
		requestedOfficeUserFetcher := &mocks.RequestedOfficeUserFetcher{}
		requestedOfficeUserUpdater := &mocks.RequestedOfficeUserUpdater{}
		userPrivilegeAssociator := &mocks.UserPrivilegeAssociator{}
		privilegeFetcher := &mocks.PrivilegeAssociator{}

		params := requestedofficeuserop.UpdateRequestedOfficeUserParams{
			HTTPRequest: suite.setupAuthenticatedRequest("PATCH", fmt.Sprintf("/requested_office_users/%s", officeUserID)),
			Body: &adminmessages.RequestedOfficeUserUpdate{
				FirstName: &officeUser.FirstName,
				LastName:  &officeUser.LastName,
				Roles: []*adminmessages.OfficeUserRole{
					{
						Name:     &tooRoleName,
						RoleType: &tooRoleType,
					},
				},
			},
			OfficeUserID: strfmt.UUID(officeUserID.String()),
		}

		requestedOfficeUserFetcher.On("FetchRequestedOfficeUser",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything,
		).Return(requestedOfficeUser, nil, nil).Once()

		requestedOfficeUserUpdater.On("UpdateRequestedOfficeUser",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything,
			mock.Anything,
		).Return(&officeUser, nil, nil).Once()

		mockRoles := roles.Roles{
			roles.Role{
				ID:        uuid.Must(uuid.NewV4()),
				RoleType:  roles.RoleTypeTOO,
				RoleName:  "Task Ordering Officer",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
		}

		// Mock roles
		mockUserRoleAssociator.On(
			"UpdateUserRoles",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything,
			mock.Anything,
		).Return(nil, nil, nil).Once()

		mockRoleAssociator.On(
			"FetchRolesForUser",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything,
		).Return(mockRoles, nil)

		privilegeFetcher.On(
			"FetchPrivilegesForUser",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything,
		).Return(roles.Privileges{}, nil)

		handler := UpdateRequestedOfficeUserHandler{
			suite.HandlerConfig(),
			requestedOfficeUserFetcher,
			requestedOfficeUserUpdater,
			userPrivilegeAssociator,
			privilegeFetcher,
			mockUserRoleAssociator,
			mockRoleAssociator,
		}

		response := handler.Handle(params)
		suite.IsType(&requestedofficeuserop.UpdateRequestedOfficeUserOK{}, response)
	})
}

func (suite *HandlerSuite) TestUpdateRequestedOfficeUserHandlerWithOktaAccountCreation() {
	suite.Run("Successful okta account creation and update", func() {

		// Build provider
		provider, err := factory.BuildOktaProvider("adminProvider")
		suite.NoError(err)

		// mocking the okta customer group id env variable
		originalGroupID := os.Getenv("OKTA_OFFICE_GROUP_ID")
		os.Setenv("OKTA_OFFICE_GROUP_ID", "notrealofficegroupId")
		defer os.Setenv("OKTA_OFFICE_GROUP_ID", originalGroupID)

		mockAndActivateOktaGETEndpointNoUserNoError(provider)
		mockAndActivateOktaEndpoints(provider, 200)
		mockAndActivateOktaGroupGETEndpointNoError(provider)
		mockAndActivateOktaGroupAddEndpointNoError(provider)

		user := factory.BuildDefaultUser(suite.DB())
		tooRoleName := "Task Ordering Officer"
		tooRoleType := string(roles.RoleTypeTOO)
		tioRoleName := "Task Invoicing Officer"
		tioRoleType := string(roles.RoleTypeTIO)
		requestedOfficeUser := factory.BuildOfficeUserWithRoles(suite.DB(), []factory.Customization{
			{
				Model: models.OfficeUser{
					Active: true,
					UserID: &user.ID,
				},
			},
			{
				Model: models.User{
					Roles: roles.Roles{
						{RoleName: roles.RoleName(tioRoleName),
							RoleType: roles.RoleType(tioRoleType)},
					},
				},
			},
		}, []roles.RoleType{roles.RoleTypeTOO})

		officeUserID := requestedOfficeUser.ID
		officeUser := models.OfficeUser{ID: officeUserID, FirstName: "Billy", LastName: "Bob", UserID: requestedOfficeUser.UserID, CreatedAt: time.Now(),
			UpdatedAt: time.Now()}

		mockUserRoleAssociator := &mocks.UserRoleAssociator{}
		mockRoleAssociator := &mocks.RoleAssociator{}
		requestedOfficeUserFetcher := &mocks.RequestedOfficeUserFetcher{}
		requestedOfficeUserUpdater := &mocks.RequestedOfficeUserUpdater{}
		userPrivilegeAssociator := &mocks.UserPrivilegeAssociator{}
		privilegeFetcher := &mocks.PrivilegeAssociator{}

		status := "APPROVED"
		email := "example@example.com"
		telephone := "000-000-0000"
		params := requestedofficeuserop.UpdateRequestedOfficeUserParams{
			HTTPRequest: suite.setupAuthenticatedRequest("PATCH", fmt.Sprintf("/requested_office_users/%s", officeUserID)),
			Body: &adminmessages.RequestedOfficeUserUpdate{
				FirstName: &officeUser.FirstName,
				LastName:  &officeUser.LastName,
				Roles: []*adminmessages.OfficeUserRole{
					{
						Name:     &tooRoleName,
						RoleType: &tooRoleType,
					},
				},
				Status:        status,
				Email:         &email,
				Telephone:     &telephone,
				OtherUniqueID: "0000000000",
				Edipi:         "0000000000",
			},
			OfficeUserID: strfmt.UUID(officeUserID.String()),
		}

		defer goth.ClearProviders()
		goth.UseProviders(provider)

		requestedOfficeUserFetcher.On("FetchRequestedOfficeUser",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything,
		).Return(requestedOfficeUser, nil, nil).Once()

		requestedOfficeUserUpdater.On("UpdateRequestedOfficeUser",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything,
			mock.Anything,
		).Return(&officeUser, nil, nil).Once()

		mockRoles := roles.Roles{
			roles.Role{
				ID:        uuid.Must(uuid.NewV4()),
				RoleType:  roles.RoleTypeTOO,
				RoleName:  "Task Ordering Officer",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
		}

		// Mock roles
		mockUserRoleAssociator.On(
			"UpdateUserRoles",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything,
			mock.Anything,
		).Return(nil, nil, nil).Once()

		mockRoleAssociator.On(
			"FetchRolesForUser",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything,
		).Return(mockRoles, nil)

		privilegeFetcher.On(
			"FetchPrivilegesForUser",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything,
		).Return(roles.Privileges{}, nil)

		handler := UpdateRequestedOfficeUserHandler{
			suite.HandlerConfig(),
			requestedOfficeUserFetcher,
			requestedOfficeUserUpdater,
			userPrivilegeAssociator,
			privilegeFetcher,
			mockUserRoleAssociator,
			mockRoleAssociator,
		}

		response := handler.Handle(params)
		suite.IsType(&requestedofficeuserop.UpdateRequestedOfficeUserOK{}, response)
	})
}

func (suite *HandlerSuite) TestUpdateRequestedOfficeUserHandlerWithOktaAccountCreationFail() {
	suite.Run("Should fail if an attempt to create an okta account fails", func() {

		// Build provider
		provider, err := factory.BuildOktaProvider("adminProvider")
		suite.NoError(err)

		mockAndActivateOktaGETEndpointNoUserNoError(provider)
		mockAndActivateOktaEndpoints(provider, 500)

		user := factory.BuildDefaultUser(suite.DB())
		tooRoleName := "Task Ordering Officer"
		tooRoleType := string(roles.RoleTypeTOO)
		tioRoleName := "Task Invoicing Officer"
		tioRoleType := string(roles.RoleTypeTIO)
		requestedOfficeUser := factory.BuildOfficeUserWithRoles(suite.DB(), []factory.Customization{
			{
				Model: models.OfficeUser{
					Active: true,
					UserID: &user.ID,
				},
			},
			{
				Model: models.User{
					Roles: roles.Roles{
						{RoleName: roles.RoleName(tioRoleName),
							RoleType: roles.RoleType(tioRoleType)},
					},
				},
			},
		}, []roles.RoleType{roles.RoleTypeTOO})

		officeUserID := requestedOfficeUser.ID
		officeUser := models.OfficeUser{ID: officeUserID, FirstName: "Billy", LastName: "Bob", UserID: requestedOfficeUser.UserID, CreatedAt: time.Now(),
			UpdatedAt: time.Now()}

		mockUserRoleAssociator := &mocks.UserRoleAssociator{}
		mockRoleAssociator := &mocks.RoleAssociator{}
		requestedOfficeUserFetcher := &mocks.RequestedOfficeUserFetcher{}
		requestedOfficeUserUpdater := &mocks.RequestedOfficeUserUpdater{}
		userPrivilegeAssociator := &mocks.UserPrivilegeAssociator{}
		privilegeFetcher := &mocks.PrivilegeAssociator{}

		status := "APPROVED"
		email := "example@example.com"
		telephone := "000-000-0000"
		params := requestedofficeuserop.UpdateRequestedOfficeUserParams{
			HTTPRequest: suite.setupAuthenticatedRequest("PATCH", fmt.Sprintf("/requested_office_users/%s", officeUserID)),
			Body: &adminmessages.RequestedOfficeUserUpdate{
				FirstName: &officeUser.FirstName,
				LastName:  &officeUser.LastName,
				Roles: []*adminmessages.OfficeUserRole{
					{
						Name:     &tooRoleName,
						RoleType: &tooRoleType,
					},
				},
				Status:        status,
				Email:         &email,
				Telephone:     &telephone,
				OtherUniqueID: "0000000000",
				Edipi:         "0000000000",
			},
			OfficeUserID: strfmt.UUID(officeUserID.String()),
		}

		defer goth.ClearProviders()
		goth.UseProviders(provider)

		requestedOfficeUserFetcher.On("FetchRequestedOfficeUser",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything,
		).Return(requestedOfficeUser, nil, nil).Once()

		requestedOfficeUserUpdater.On("UpdateRequestedOfficeUser",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything,
			mock.Anything,
		).Return(&officeUser, nil, nil).Once()

		mockRoles := roles.Roles{
			roles.Role{
				ID:        uuid.Must(uuid.NewV4()),
				RoleType:  roles.RoleTypeTOO,
				RoleName:  "Task Ordering Officer",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
		}

		// Mock roles
		mockUserRoleAssociator.On(
			"UpdateUserRoles",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything,
			mock.Anything,
		).Return(nil, nil, nil).Once()

		mockRoleAssociator.On(
			"FetchRolesForUser",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything,
		).Return(mockRoles, nil)

		privilegeFetcher.On(
			"FetchPrivilegesForUser",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything,
		).Return(roles.Privileges{}, nil)

		handler := UpdateRequestedOfficeUserHandler{
			suite.HandlerConfig(),
			requestedOfficeUserFetcher,
			requestedOfficeUserUpdater,
			userPrivilegeAssociator,
			privilegeFetcher,
			mockUserRoleAssociator,
			mockRoleAssociator,
		}

		response := handler.Handle(params)
		suite.IsType(requestedofficeuserop.NewUpdateRequestedOfficeUserInternalServerError(), response)
	})
}

// Generate and activate Okta endpoints that will be using during the handler
func mockAndActivateOktaEndpoints(provider *okta.Provider, responseCode int) {
	activate := "true"

	createAccountEndpoint := provider.GetCreateAccountURL(activate)
	oktaID := "fakeSub"

	if responseCode == 200 {
		httpmock.RegisterResponder("POST", createAccountEndpoint,
			httpmock.NewStringResponder(200, fmt.Sprintf(`{
		"id": "%s",
		"profile": {
			"firstName": "First",
			"lastName": "Last",
			"email": "email@email.com",
			"login": "email@email.com"
		}
	}`, oktaID)))
	} else if responseCode == 500 {
		httpmock.RegisterResponder("POST", createAccountEndpoint,
			httpmock.NewStringResponder(500, ""))
	}

	httpmock.Activate()
}

func mockAndActivateOktaGETEndpointNoUserNoError(provider *okta.Provider) {
	getUsersEndpoint := provider.GetUsersURL()
	response := "[]"

	httpmock.RegisterResponder("GET", getUsersEndpoint,
		httpmock.NewStringResponder(200, response))
	httpmock.Activate()
}

func mockAndActivateOktaGroupGETEndpointNoError(provider *okta.Provider) {

	oktaID := "fakeSub"
	getGroupsEndpoint := provider.GetUserGroupsURL(oktaID)

	httpmock.RegisterResponder("GET", getGroupsEndpoint,
		httpmock.NewStringResponder(200, `[]`))

	httpmock.Activate()
}

func mockAndActivateOktaGroupAddEndpointNoError(provider *okta.Provider) {

	oktaID := "fakeSub"
	groupID := "notrealofficegroupId"
	addGroupEndpoint := provider.AddUserToGroupURL(groupID, oktaID)

	httpmock.RegisterResponder("PUT", addGroupEndpoint,
		httpmock.NewStringResponder(204, ""))

	httpmock.Activate()
}
