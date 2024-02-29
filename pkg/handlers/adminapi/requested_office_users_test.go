package adminapi

import (
	"fmt"
	"net/http"

	"github.com/go-openapi/strfmt"
	"github.com/stretchr/testify/mock"
	"github.com/transcom/mymove/pkg/factory"
	requestedofficeuserop "github.com/transcom/mymove/pkg/gen/adminapi/adminoperations/requested_office_users"
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

		queryBuilder := query.NewQueryBuilder()
		handler := GetRequestedOfficeUserHandler{
			suite.HandlerConfig(),
			requestedofficeusers.NewRequestedOfficeUserFetcher(queryBuilder),
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
		handler := GetRequestedOfficeUserHandler{
			suite.HandlerConfig(),
			requestedOfficeUserFetcher,
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
		handler := GetRequestedOfficeUserHandler{
			suite.HandlerConfig(),
			requestedOfficeUserFetcher,
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
