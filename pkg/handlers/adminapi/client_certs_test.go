package adminapi

import (
	clientcertop "github.com/transcom/mymove/pkg/gen/adminapi/adminoperations/client_certs"
	"github.com/transcom/mymove/pkg/models"
	fetch "github.com/transcom/mymove/pkg/services/fetch"
	"github.com/transcom/mymove/pkg/services/pagination"
	"github.com/transcom/mymove/pkg/services/query"
)

func (suite *HandlerSuite) TestIndexClientCertsHandler() {
	// test that everything is wired up
	suite.Run("integration test ok response", func() {
		users := models.ClientCerts{
			models.ClientCert{
				Subject:      "CN=example-user,OU=DoD+OU=PKI+OU=CONTRACTOR,O=U.S. Government,C=US",
				Sha256Digest: "01ba4719c80b6fe911b091a7c05124b64eeece964e09c058ef8f9805daca546b",
			},
		}
		params := clientcertop.IndexClientCertsParams{
			HTTPRequest: suite.setupAuthenticatedRequest("GET", "/client_certs"),
		}

		queryBuilder := query.NewQueryBuilder()
		handler := IndexClientCertsHandler{
			HandlerConfig:         suite.HandlerConfig(),
			NewQueryFilter:        query.NewQueryFilter,
			ClientCertListFetcher: fetch.NewListFetcher(queryBuilder),
			NewPagination:         pagination.NewPagination,
		}

		response := handler.Handle(params)

		suite.IsType(&userop.IndexUsersOK{}, response)
		okResponse := response.(*userop.IndexUsersOK)
		suite.Len(okResponse.Payload, 2)
		suite.Equal(users[0].ID.String(), okResponse.Payload[0].ID.String())
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
