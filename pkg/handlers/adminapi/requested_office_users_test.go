package adminapi

import (
	"fmt"
	"net/http"
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
		created := time.Date(2005, 03, 05, 1, 1, 1, 1, time.Local)
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
					TransportationOfficeID: transportationOffice2.ID,
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
					CreatedAt:              created,
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
		requestedOnSearch := "2005"
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

		mockRoleAssociator := &mocks.RoleAssociater{}
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

		mockRoleAssociator := &mocks.RoleAssociater{}
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

		mockRoleAssociator := &mocks.RoleAssociater{}
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
			newMockQueryFilterBuilder(&mocks.QueryFilter{}),
		}

		response := handler.Handle(params)

		expectedResponse := &handlers.ErrResponse{
			Code: http.StatusNotFound,
			Err:  expectedError,
		}
		suite.Equal(expectedResponse, response)
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

		officeUserID := requestedOfficeUser.ID
		officeUser := models.OfficeUser{ID: officeUserID, FirstName: "Billy", LastName: "Bob", UserID: requestedOfficeUser.UserID, CreatedAt: time.Now(),
			UpdatedAt: time.Now()}

		mockUserRoleAssociator := &mocks.UserRoleAssociator{}
		mockRoleAssociator := &mocks.RoleAssociater{}
		requestedOfficeUserUpdater := &mocks.RequestedOfficeUserUpdater{}

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

		handler := UpdateRequestedOfficeUserHandler{
			suite.HandlerConfig(),
			requestedOfficeUserUpdater,
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

		mockAndActivateOktaEndpoints(provider, 200)

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
		mockRoleAssociator := &mocks.RoleAssociater{}
		requestedOfficeUserUpdater := &mocks.RequestedOfficeUserUpdater{}

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
				Email:         email,
				Telephone:     &telephone,
				OtherUniqueID: "0000000000",
				Edipi:         "0000000000",
			},
			OfficeUserID: strfmt.UUID(officeUserID.String()),
		}

		defer goth.ClearProviders()
		goth.UseProviders(provider)

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

		handler := UpdateRequestedOfficeUserHandler{
			suite.HandlerConfig(),
			requestedOfficeUserUpdater,
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
		mockRoleAssociator := &mocks.RoleAssociater{}
		requestedOfficeUserUpdater := &mocks.RequestedOfficeUserUpdater{}

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
				Email:         email,
				Telephone:     &telephone,
				OtherUniqueID: "0000000000",
				Edipi:         "0000000000",
			},
			OfficeUserID: strfmt.UUID(officeUserID.String()),
		}

		defer goth.ClearProviders()
		goth.UseProviders(provider)

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

		handler := UpdateRequestedOfficeUserHandler{
			suite.HandlerConfig(),
			requestedOfficeUserUpdater,
			mockUserRoleAssociator,
			mockRoleAssociator,
		}

		response := handler.Handle(params)
		suite.IsType(requestedofficeuserop.NewGetRequestedOfficeUserInternalServerError(), response)
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
