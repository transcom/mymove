package adminapi

import (
	"net/http"
	"net/http/httptest"
	"testing"

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
	// replace this with generated UUID when filter param is built out
	uuidString := "5ce7162a-8d5c-41fc-b0e7-bae726f98fa2"
	id, _ := uuid.FromString(uuidString)
	assertions := testdatagen.Assertions{
		Organization: models.Organization{
			ID: id,
		},
	}
	testdatagen.MakeOrganization(suite.DB(), assertions)

	requestUser := testdatagen.MakeStubbedUser(suite.DB())
	req := httptest.NewRequest("GET", "/organizations", nil)
	req = suite.AuthenticateUserRequest(req, requestUser)

	// test that everything is wired up
	suite.T().Run("integration test ok response", func(t *testing.T) {
		params := organization.IndexOrganizationsParams{
			HTTPRequest: req,
		}
		queryBuilder := query.NewQueryBuilder(suite.DB())
		handler := IndexOrganizationsHandler{
			HandlerContext:          handlers.NewHandlerContext(suite.DB(), suite.TestLogger()),
			NewQueryFilter:          query.NewQueryFilter,
			OrganizationListFetcher: organization2.NewOrganizationListFetcher(queryBuilder),
			NewPagination:           pagination.NewPagination,
		}

		response := handler.Handle(params)

		suite.IsType(&organization.IndexOrganizationsOK{}, response)
		okResponse := response.(*organization.IndexOrganizationsOK)
		suite.Len(okResponse.Payload, 1)
		suite.Equal(uuidString, okResponse.Payload[0].ID.String())
	})

	queryFilter := mocks.QueryFilter{}
	newQueryFilter := newMockQueryFilterBuilder(&queryFilter)

	suite.T().Run("successful response", func(t *testing.T) {
		org := models.Organization{ID: id}
		params := organization.IndexOrganizationsParams{
			HTTPRequest: req,
		}
		organizationListFetcher := &mocks.OrganizationListFetcher{}
		organizationListFetcher.On("FetchOrganizationList",
			mock.Anything,
			mock.Anything,
			mock.Anything,
			mock.Anything,
		).Return(models.Organizations{org}, nil).Once()
		organizationListFetcher.On("FetchOrganizationCount",
			mock.Anything,
		).Return(1, nil).Once()
		handler := IndexOrganizationsHandler{
			HandlerContext:          handlers.NewHandlerContext(suite.DB(), suite.TestLogger()),
			NewQueryFilter:          newQueryFilter,
			OrganizationListFetcher: organizationListFetcher,
			NewPagination:           pagination.NewPagination,
		}

		response := handler.Handle(params)

		suite.IsType(&organization.IndexOrganizationsOK{}, response)
		okResponse := response.(*organization.IndexOrganizationsOK)
		suite.Len(okResponse.Payload, 1)
		suite.Equal(uuidString, okResponse.Payload[0].ID.String())
	})

	suite.T().Run("unsuccesful response when fetch fails", func(t *testing.T) {
		params := organization.IndexOrganizationsParams{
			HTTPRequest: req,
		}
		expectedError := models.ErrFetchNotFound
		organizationListFetcher := &mocks.OrganizationListFetcher{}
		organizationListFetcher.On("FetchOrganizationList",
			mock.Anything,
			mock.Anything,
			mock.Anything,
			mock.Anything,
		).Return(nil, expectedError).Once()
		organizationListFetcher.On("FetchOrganizationCount",
			mock.Anything,
		).Return(0, expectedError).Once()
		handler := IndexOrganizationsHandler{
			HandlerContext:          handlers.NewHandlerContext(suite.DB(), suite.TestLogger()),
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
