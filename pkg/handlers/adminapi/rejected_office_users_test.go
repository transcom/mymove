package adminapi

import (
	"fmt"
	"net/http"
	"slices"
	"time"

	"github.com/go-openapi/strfmt"
	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/mock"

	"github.com/transcom/mymove/pkg/factory"
	rejectedofficeuserop "github.com/transcom/mymove/pkg/gen/adminapi/adminoperations/rejected_office_users"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/models/roles"
	"github.com/transcom/mymove/pkg/services/mocks"
	"github.com/transcom/mymove/pkg/services/pagination"
	"github.com/transcom/mymove/pkg/services/query"
	rejectedofficeusers "github.com/transcom/mymove/pkg/services/rejected_office_users"
)

func (suite *HandlerSuite) TestIndexRejectedOfficeUsersHandler() {
	// test that everything is wired up
	suite.Run("rejected users result in ok response", func() {
		// building two office user with rejected status
		rejectedOfficeUsers := models.OfficeUsers{
			factory.BuildOfficeUserWithRoles(suite.DB(), factory.GetTraitRejectedOfficeUser(), []roles.RoleType{roles.RoleTypeQae}),
			factory.BuildOfficeUserWithRoles(suite.DB(), factory.GetTraitRejectedOfficeUser(), []roles.RoleType{roles.RoleTypeQae})}

		params := rejectedofficeuserop.IndexRejectedOfficeUsersParams{
			HTTPRequest: suite.setupAuthenticatedRequest("GET", "/rejected_office_users"),
		}

		queryBuilder := query.NewQueryBuilder()
		handler := IndexRejectedOfficeUsersHandler{
			HandlerConfig:                 suite.HandlerConfig(),
			NewQueryFilter:                query.NewQueryFilter,
			RejectedOfficeUserListFetcher: rejectedofficeusers.NewRejectedOfficeUsersListFetcher(queryBuilder),
			NewPagination:                 pagination.NewPagination,
		}

		response := handler.Handle(params)

		// should get an ok response
		suite.IsType(&rejectedofficeuserop.IndexRejectedOfficeUsersOK{}, response)
		okResponse := response.(*rejectedofficeuserop.IndexRejectedOfficeUsersOK)
		suite.Equal(len(okResponse.Payload), len(rejectedOfficeUsers))

		actualID := []string{okResponse.Payload[0].ID.String(), okResponse.Payload[1].ID.String()}
		expected := []string{rejectedOfficeUsers[0].ID.String(), rejectedOfficeUsers[1].ID.String()}
		for _, expectedID := range expected {
			suite.True(slices.Contains(actualID, expectedID))
		}
	})
}

func (suite *HandlerSuite) TestGetRejectedOfficeUserHandler() {
	suite.Run("integration test ok response", func() {
		rejectedOfficeUser := factory.BuildOfficeUserWithRoles(suite.DB(), factory.GetTraitRejectedOfficeUser(), []roles.RoleType{roles.RoleTypeQae})
		params := rejectedofficeuserop.GetRejectedOfficeUserParams{
			HTTPRequest:  suite.setupAuthenticatedRequest("GET", fmt.Sprintf("/rejected_office_users/%s", rejectedOfficeUser.ID)),
			OfficeUserID: strfmt.UUID(rejectedOfficeUser.ID.String()),
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
		handler := GetRejectedOfficeUserHandler{
			suite.HandlerConfig(),
			rejectedofficeusers.NewRejectedOfficeUserFetcher(queryBuilder),
			mockRoleAssociator,
			query.NewQueryFilter,
		}

		response := handler.Handle(params)

		suite.IsType(&rejectedofficeuserop.GetRejectedOfficeUserOK{}, response)
		okResponse := response.(*rejectedofficeuserop.GetRejectedOfficeUserOK)
		suite.Equal(rejectedOfficeUser.ID.String(), okResponse.Payload.ID.String())
	})

	suite.Run("successful response", func() {
		rejectedOfficeUser := factory.BuildOfficeUserWithRoles(suite.DB(), factory.GetTraitRejectedOfficeUser(), []roles.RoleType{roles.RoleTypeQae})
		params := rejectedofficeuserop.GetRejectedOfficeUserParams{
			HTTPRequest:  suite.setupAuthenticatedRequest("GET", fmt.Sprintf("/rejected_office_users/%s", rejectedOfficeUser.ID)),
			OfficeUserID: strfmt.UUID(rejectedOfficeUser.ID.String()),
		}

		rejectedOfficeUserFetcher := &mocks.RejectedOfficeUserFetcher{}
		rejectedOfficeUserFetcher.On("FetchRejectedOfficeUser",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything,
		).Return(rejectedOfficeUser, nil).Once()

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

		handler := GetRejectedOfficeUserHandler{
			suite.HandlerConfig(),
			rejectedOfficeUserFetcher,
			mockRoleAssociator,
			newMockQueryFilterBuilder(&mocks.QueryFilter{}),
		}

		response := handler.Handle(params)

		suite.IsType(&rejectedofficeuserop.GetRejectedOfficeUserOK{}, response)
		okResponse := response.(*rejectedofficeuserop.GetRejectedOfficeUserOK)
		suite.Equal(rejectedOfficeUser.ID.String(), okResponse.Payload.ID.String())
	})

	suite.Run("unsuccessful response when fetch fails", func() {
		rejectedOfficeUser := factory.BuildOfficeUserWithRoles(suite.DB(), factory.GetTraitRejectedOfficeUser(), []roles.RoleType{roles.RoleTypeQae})
		params := rejectedofficeuserop.GetRejectedOfficeUserParams{
			HTTPRequest:  suite.setupAuthenticatedRequest("GET", fmt.Sprintf("/rejected_office_users/%s", rejectedOfficeUser.ID)),
			OfficeUserID: strfmt.UUID(rejectedOfficeUser.ID.String()),
		}

		expectedError := models.ErrFetchNotFound
		rejectedOfficeUserFetcher := &mocks.RejectedOfficeUserFetcher{}
		rejectedOfficeUserFetcher.On("FetchRejectedOfficeUser",
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

		handler := GetRejectedOfficeUserHandler{
			suite.HandlerConfig(),
			rejectedOfficeUserFetcher,
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
