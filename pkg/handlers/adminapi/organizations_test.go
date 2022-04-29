package adminapi

import (
	"net/http"

	"github.com/transcom/mymove/pkg/gen/adminapi/adminoperations/organization"
	organization2 "github.com/transcom/mymove/pkg/services/organization"

	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/mock"

	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services/mocks"
	"github.com/transcom/mymove/pkg/services/pagination"
	"github.com/transcom/mymove/pkg/services/query"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *HandlerSuite) TestIndexOrganizationsHandler() {
	// test that everything is wired up
	suite.Run("integration test ok response", func() {
		org := testdatagen.MakeDefaultOrganization(suite.DB())
		params := organization.IndexOrganizationsParams{
			HTTPRequest: suite.setupAuthenticatedRequest("GET", "/organizations"),
		}
		queryBuilder := query.NewQueryBuilder()
		handler := IndexOrganizationsHandler{
			HandlerContext:          handlers.NewHandlerContext(suite.DB(), suite.Logger()),
			NewQueryFilter:          query.NewQueryFilter,
			OrganizationListFetcher: organization2.NewOrganizationListFetcher(queryBuilder),
			NewPagination:           pagination.NewPagination,
		}

		response := handler.Handle(params)

		suite.IsType(&organization.IndexOrganizationsOK{}, response)
		okResponse := response.(*organization.IndexOrganizationsOK)
		suite.Len(okResponse.Payload, 1)
		suite.Equal(org.ID.String(), okResponse.Payload[0].ID.String())
	})

	queryFilter := mocks.QueryFilter{}
	newQueryFilter := newMockQueryFilterBuilder(&queryFilter)

	suite.Run("successful response", func() {
		id, _ := uuid.FromString("5ce7162a-8d5c-41fc-b0e7-bae726f98fa2")
		org := models.Organization{ID: id}
		params := organization.IndexOrganizationsParams{
			HTTPRequest: suite.setupAuthenticatedRequest("GET", "/organizations"),
		}
		organizationListFetcher := &mocks.OrganizationListFetcher{}
		organizationListFetcher.On("FetchOrganizationList",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything,
			mock.Anything,
			mock.Anything,
			mock.Anything,
		).Return(models.Organizations{org}, nil).Once()
		organizationListFetcher.On("FetchOrganizationCount",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything,
		).Return(1, nil).Once()
		handler := IndexOrganizationsHandler{
			HandlerContext:          handlers.NewHandlerContext(suite.DB(), suite.Logger()),
			NewQueryFilter:          newQueryFilter,
			OrganizationListFetcher: organizationListFetcher,
			NewPagination:           pagination.NewPagination,
		}

		response := handler.Handle(params)

		suite.IsType(&organization.IndexOrganizationsOK{}, response)
		okResponse := response.(*organization.IndexOrganizationsOK)
		suite.Len(okResponse.Payload, 1)
		suite.Equal(id.String(), okResponse.Payload[0].ID.String())
	})

	suite.Run("unsuccesful response when fetch fails", func() {
		params := organization.IndexOrganizationsParams{
			HTTPRequest: suite.setupAuthenticatedRequest("GET", "/organizations"),
		}
		expectedError := models.ErrFetchNotFound
		organizationListFetcher := &mocks.OrganizationListFetcher{}
		organizationListFetcher.On("FetchOrganizationList",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything,
			mock.Anything,
			mock.Anything,
			mock.Anything,
			mock.Anything,
		).Return(nil, expectedError).Once()
		organizationListFetcher.On("FetchOrganizationCount",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything,
		).Return(0, expectedError).Once()
		handler := IndexOrganizationsHandler{
			HandlerContext:          handlers.NewHandlerContext(suite.DB(), suite.Logger()),
			NewQueryFilter:          newQueryFilter,
			OrganizationListFetcher: organizationListFetcher,
			NewPagination:           pagination.NewPagination,
		}

		response := handler.Handle(params)

		expectedResponse := &handlers.ErrResponse{
			Code: http.StatusNotFound,
			Err:  expectedError,
		}
		suite.Equal(expectedResponse, response)
	})
}
