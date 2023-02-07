package adminapi

import (
	"github.com/transcom/mymove/pkg/factory"
	clientcertop "github.com/transcom/mymove/pkg/gen/adminapi/adminoperations/client_certs"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services/clientcert"
	"github.com/transcom/mymove/pkg/services/pagination"
	"github.com/transcom/mymove/pkg/services/query"
)

func (suite *HandlerSuite) TestIndexClientCertsHandler() {
	// test that everything is wired up
	suite.Run("integration test ok response", func() {
		clientCerts := models.ClientCerts{
			factory.BuildDefaultClientCert(suite.DB()),
		}

		params := clientcertop.IndexClientCertsParams{
			HTTPRequest: suite.setupAuthenticatedRequest("GET", "/client_certs"),
		}

		queryBuilder := query.NewQueryBuilder()
		handler := IndexClientCertsHandler{
			HandlerConfig:         suite.HandlerConfig(),
			NewQueryFilter:        query.NewQueryFilter,
			ClientCertListFetcher: clientcert.NewClientCertListFetcher(queryBuilder),
			NewPagination:         pagination.NewPagination,
		}

		response := handler.Handle(params)

		suite.IsType(&clientcertop.IndexClientCertsOK{}, response)
		okResponse := response.(*clientcertop.IndexClientCertsOK)
		suite.Len(okResponse.Payload, 2)
		suite.Equal(clientCerts[0].ID.String(), okResponse.Payload[0].ID.String())
	})

	// suite.Run("unsuccesful response when fetch fails", func() {
	// 	queryFilter := mocks.QueryFilter{}
	// 	newQueryFilter := newMockQueryFilterBuilder(&queryFilter)
	//
	// 	params := userop.IndexUsersParams{
	// 		HTTPRequest: suite.setupAuthenticatedRequest("GET", "/users"),
	// 	}
	// 	expectedError := models.ErrFetchNotFound
	// 	userListFetcher := &mocks.ListFetcher{}
	// 	userListFetcher.On("FetchRecordList",
	// 		mock.AnythingOfType("*appcontext.appContext"),
	// 		mock.Anything,
	// 		mock.Anything,
	// 		mock.Anything,
	// 		mock.Anything,
	// 		mock.Anything,
	// 	).Return(nil, expectedError).Once()
	// 	userListFetcher.On("FetchRecordCount",
	// 		mock.AnythingOfType("*appcontext.appContext"),
	// 		mock.Anything,
	// 		mock.Anything,
	// 	).Return(0, expectedError).Once()
	// 	handler := IndexUsersHandler{
	// 		HandlerConfig:  suite.HandlerConfig(),
	// 		NewQueryFilter: newQueryFilter,
	// 		ListFetcher:    userListFetcher,
	// 		NewPagination:  pagination.NewPagination,
	// 	}
	//
	// 	response := handler.Handle(params)
	//
	// 	expectedResponse := &handlers.ErrResponse{
	// 		Code: http.StatusNotFound,
	// 		Err:  expectedError,
	// 	}
	// 	suite.Equal(expectedResponse, response)
	// })
}
