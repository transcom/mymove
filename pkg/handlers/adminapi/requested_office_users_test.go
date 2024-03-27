package adminapi

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-openapi/strfmt"
	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/mock"

	"github.com/transcom/mymove/pkg/factory"
	requestedofficeuserop "github.com/transcom/mymove/pkg/gen/adminapi/adminoperations/requested_office_users"
	"github.com/transcom/mymove/pkg/gen/adminmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/models/roles"
	"github.com/transcom/mymove/pkg/services/mocks"
	"github.com/transcom/mymove/pkg/services/pagination"
	"github.com/transcom/mymove/pkg/services/query"
	requestedofficeusers "github.com/transcom/mymove/pkg/services/requested_office_users"
)

func (suite *HandlerSuite) TestIndexRequestedOfficeUsersHandler() {
	// test that everything is wired up
	suite.Run("requested users result in ok response", func() {
		// building two office user with requested status
		requestedOfficeUsers := models.OfficeUsers{
			factory.BuildOfficeUserWithRoles(suite.DB(), factory.GetTraitRequestedOfficeUser(), []roles.RoleType{roles.RoleTypeQaeCsr}),
			factory.BuildOfficeUserWithRoles(suite.DB(), factory.GetTraitRequestedOfficeUser(), []roles.RoleType{roles.RoleTypeQaeCsr})}
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

		// should get an ok response
		suite.IsType(&requestedofficeuserop.IndexRequestedOfficeUsersOK{}, response)
		okResponse := response.(*requestedofficeuserop.IndexRequestedOfficeUsersOK)
		suite.Len(okResponse.Payload, 2)
		suite.Equal(requestedOfficeUsers[0].ID.String(), okResponse.Payload[0].ID.String())
	})
}

func (suite *HandlerSuite) TestGetRequestedOfficeUserHandler() {
	// test that everything is wired up
	suite.Run("integration test ok response", func() {
		requestedOfficeUser := factory.BuildOfficeUserWithRoles(suite.DB(), factory.GetTraitRequestedOfficeUser(), []roles.RoleType{roles.RoleTypeQaeCsr})
		params := requestedofficeuserop.GetRequestedOfficeUserParams{
			HTTPRequest:  suite.setupAuthenticatedRequest("GET", fmt.Sprintf("/requested_office_users/%s", requestedOfficeUser.ID)),
			OfficeUserID: strfmt.UUID(requestedOfficeUser.ID.String()),
		}

		mockRoleAssociator := &mocks.RoleAssociater{}
		mockRoles := roles.Roles{
			roles.Role{
				ID:        uuid.Must(uuid.NewV4()),
				RoleType:  roles.RoleTypeTOO,
				RoleName:  "Transportation Ordering Officer",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
		}
		mockRoleAssociator.On(
			"FetchRoles",
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
		requestedOfficeUser := factory.BuildOfficeUserWithRoles(suite.DB(), factory.GetTraitRequestedOfficeUser(), []roles.RoleType{roles.RoleTypeQaeCsr})
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
				RoleName:  "Transportation Ordering Officer",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
		}
		mockRoleAssociator.On(
			"FetchRoles",
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
		requestedOfficeUser := factory.BuildOfficeUserWithRoles(suite.DB(), factory.GetTraitRequestedOfficeUser(), []roles.RoleType{roles.RoleTypeQaeCsr})
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
				RoleName:  "Transportation Ordering Officer",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
		}
		mockRoleAssociator.On(
			"FetchRoles",
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

func (suite *HandlerSuite) TestUpdateRequestedOfficeUserHandler() {
	suite.Run("Successful update", func() {
		user := factory.BuildDefaultUser(suite.DB())
		tooRoleName := "Transportation Ordering Officer"
		tooRoleType := string(roles.RoleTypeTOO)
		tioRoleName := "Transportation Invoicing Officer"
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
				RoleName:  "Transportation Ordering Officer",
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
		).Return(nil, nil).Once()

		mockRoleAssociator.On(
			"FetchRoles",
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
