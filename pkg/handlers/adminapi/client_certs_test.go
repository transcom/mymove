package adminapi

import (
	"fmt"
	"net/http"

	"github.com/go-openapi/strfmt"
	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/mock"

	"github.com/transcom/mymove/pkg/factory"
	clientcertop "github.com/transcom/mymove/pkg/gen/adminapi/adminoperations/client_certificates"
	userop "github.com/transcom/mymove/pkg/gen/adminapi/adminoperations/users"
	"github.com/transcom/mymove/pkg/gen/adminmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services/clientcert"
	"github.com/transcom/mymove/pkg/services/mocks"
	"github.com/transcom/mymove/pkg/services/pagination"
	"github.com/transcom/mymove/pkg/services/query"
)

func (suite *HandlerSuite) TestIndexClientCertsHandler() {
	// test that everything is wired up
	suite.Run("integration test ok response", func() {
		clientCerts := models.ClientCerts{
			factory.BuildClientCert(suite.DB(), nil, nil),
		}

		params := clientcertop.IndexClientCertificatesParams{
			HTTPRequest: suite.setupAuthenticatedRequest("GET", "/client-certificates"),
		}

		queryBuilder := query.NewQueryBuilder()
		handler := IndexClientCertsHandler{
			HandlerConfig:         suite.NewHandlerConfig(),
			NewQueryFilter:        query.NewQueryFilter,
			ClientCertListFetcher: clientcert.NewClientCertListFetcher(queryBuilder),
			NewPagination:         pagination.NewPagination,
		}

		response := handler.Handle(params)

		suite.IsType(&clientcertop.IndexClientCertificatesOK{}, response)
		okResponse := response.(*clientcertop.IndexClientCertificatesOK)
		suite.Len(okResponse.Payload, 1)
		suite.Equal(clientCerts[0].ID.String(), okResponse.Payload[0].ID.String())
	})

	suite.Run("unsuccesful response when fetch fails", func() {
		queryFilter := mocks.QueryFilter{}
		newQueryFilter := newMockQueryFilterBuilder(&queryFilter)

		params := userop.IndexUsersParams{
			HTTPRequest: suite.setupAuthenticatedRequest("GET", "/client-certificates"),
		}
		expectedError := models.ErrFetchNotFound
		clientCertListFetcher := &mocks.ListFetcher{}
		clientCertListFetcher.On("FetchRecordList",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything,
			mock.Anything,
			mock.Anything,
			mock.Anything,
			mock.Anything,
		).Return(nil, expectedError).Once()
		clientCertListFetcher.On("FetchRecordCount",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything,
			mock.Anything,
		).Return(0, expectedError).Once()
		handler := IndexUsersHandler{
			HandlerConfig:  suite.NewHandlerConfig(),
			NewQueryFilter: newQueryFilter,
			ListFetcher:    clientCertListFetcher,
			NewPagination:  pagination.NewPagination,
		}

		response := handler.Handle(params)

		expectedResponse := &handlers.ErrResponse{
			Code: http.StatusNotFound,
			Err:  expectedError,
		}
		suite.Equal(expectedResponse, response)
	})
}

func (suite *HandlerSuite) TestGetClientCertHandler() {
	// test that everything is wired up
	suite.Run("integration test ok response", func() {
		clientCert := factory.BuildClientCert(suite.DB(), nil, nil)
		params := clientcertop.GetClientCertificateParams{
			HTTPRequest:         suite.setupAuthenticatedRequest("GET", fmt.Sprintf("/client-certificates/%s", clientCert.ID)),
			ClientCertificateID: strfmt.UUID(clientCert.ID.String()),
		}

		queryBuilder := query.NewQueryBuilder()
		handler := GetClientCertHandler{
			suite.NewHandlerConfig(),
			clientcert.NewClientCertFetcher(queryBuilder),
			query.NewQueryFilter,
		}

		response := handler.Handle(params)

		suite.IsType(&clientcertop.GetClientCertificateOK{}, response)
		okResponse := response.(*clientcertop.GetClientCertificateOK)
		suite.Equal(clientCert.ID.String(), okResponse.Payload.ID.String())
	})

	suite.Run("successful response", func() {
		clientCert := factory.BuildClientCert(suite.DB(), nil, nil)
		params := clientcertop.GetClientCertificateParams{
			HTTPRequest:         suite.setupAuthenticatedRequest("GET", fmt.Sprintf("/client-certificates/%s", clientCert.ID)),
			ClientCertificateID: strfmt.UUID(clientCert.ID.String()),
		}
		clientCertFetcher := &mocks.ClientCertFetcher{}
		clientCertFetcher.On("FetchClientCert",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything,
		).Return(clientCert, nil).Once()
		handler := GetClientCertHandler{
			suite.NewHandlerConfig(),
			clientCertFetcher,
			newMockQueryFilterBuilder(&mocks.QueryFilter{}),
		}

		response := handler.Handle(params)

		suite.IsType(&clientcertop.GetClientCertificateOK{}, response)
		okResponse := response.(*clientcertop.GetClientCertificateOK)
		suite.Equal(clientCert.ID.String(), okResponse.Payload.ID.String())
	})

	suite.Run("unsuccessful response when fetch fails", func() {
		clientCert := factory.BuildClientCert(suite.DB(), nil, nil)
		params := clientcertop.GetClientCertificateParams{
			HTTPRequest:         suite.setupAuthenticatedRequest("GET", fmt.Sprintf("/client-certificates/%s", clientCert.ID)),
			ClientCertificateID: strfmt.UUID(clientCert.ID.String()),
		}
		expectedError := models.ErrFetchNotFound
		clientCertFetcher := &mocks.ClientCertFetcher{}
		clientCertFetcher.On("FetchClientCert",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything,
		).Return(models.ClientCert{}, expectedError).Once()
		handler := GetClientCertHandler{
			suite.NewHandlerConfig(),
			clientCertFetcher,
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

func (suite *HandlerSuite) TestCreateClientCertificateHandler() {
	clientCert := models.ClientCert{
		ID:           uuid.Nil,
		Subject:      "CN=fake-cn",
		Sha256Digest: "fakesha",
		AllowPrime:   true,
		UserID:       uuid.Nil,
	}
	email := "fakecncert@example.com"

	suite.Run("Successful create", func() {
		params := clientcertop.CreateClientCertificateParams{
			HTTPRequest: suite.setupAuthenticatedRequest("POST", "/client-certificates"),
			ClientCertificate: &adminmessages.ClientCertificateCreate{
				Subject:          &clientCert.Subject,
				Sha256Digest:     &clientCert.Sha256Digest,
				Email:            &email,
				AllowPrime:       true,
				AllowPPTAS:       true,
				PptasAffiliation: (*adminmessages.Affiliation)(models.StringPointer("MARINES")),
			},
		}

		clientCertCreator := &mocks.ClientCertCreator{}
		clientCertCreator.On("CreateClientCert",
			mock.AnythingOfType("*appcontext.appContext"),
			email,
			&clientCert).Return(&clientCert, nil, nil).Once()

		handler := CreateClientCertHandler{
			suite.NewHandlerConfig(),
			clientCertCreator,
		}

		response := handler.Handle(params)
		suite.IsType(&clientcertop.CreateClientCertificateCreated{}, response)
	})

	suite.Run("Failed create", func() {
		params := clientcertop.CreateClientCertificateParams{
			HTTPRequest: suite.setupAuthenticatedRequest("POST", "/client-certificates"),
			ClientCertificate: &adminmessages.ClientCertificateCreate{
				Subject:          &clientCert.Subject,
				Sha256Digest:     &clientCert.Sha256Digest,
				Email:            &email,
				AllowPrime:       true,
				AllowPPTAS:       true,
				PptasAffiliation: (*adminmessages.Affiliation)(models.StringPointer("MARINES")),
			},
		}

		expectedError := models.ErrWriteConflict
		clientCertCreator := &mocks.ClientCertCreator{}
		clientCertCreator.On("CreateClientCert",
			mock.AnythingOfType("*appcontext.appContext"),
			email,
			&clientCert).Return(&models.ClientCert{}, nil, expectedError).Once()

		handler := CreateClientCertHandler{
			suite.NewHandlerConfig(),
			clientCertCreator,
		}

		response := handler.Handle(params)
		suite.IsType(&handlers.ErrResponse{}, response)
	})
}

func (suite *HandlerSuite) TestUpdateClientCertificateHandler() {
	clientCert := models.ClientCert{
		ID:               uuid.Nil,
		AllowPrime:       false,
		UserID:           uuid.Nil,
		AllowPPTAS:       true,
		PPTASAffiliation: (*models.ServiceMemberAffiliation)(models.StringPointer("MARINES")),
	}
	queryFilter := mocks.QueryFilter{}
	newQueryFilter := newMockQueryFilterBuilder(&queryFilter)

	suite.Run("Successful update", func() {
		params := clientcertop.UpdateClientCertificateParams{
			HTTPRequest: suite.setupAuthenticatedRequest("PUT", fmt.Sprintf("/client-certificates/%s", clientCert.ID)),
			ClientCertificate: &adminmessages.ClientCertificateUpdate{
				AllowPrime:       &clientCert.AllowPrime,
				AllowPPTAS:       &clientCert.AllowPPTAS,
				PptasAffiliation: (*adminmessages.Affiliation)(clientCert.PPTASAffiliation),
			},
		}

		clientCertUpdater := &mocks.ClientCertUpdater{}
		clientCertUpdater.On("UpdateClientCert",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything,
			params.ClientCertificate,
		).Return(&clientCert, nil, nil).Once()

		handler := UpdateClientCertHandler{
			suite.NewHandlerConfig(),
			clientCertUpdater,
			newQueryFilter,
		}

		response := handler.Handle(params)
		suite.IsType(&clientcertop.UpdateClientCertificateOK{}, response)
	})

	suite.Run("Failed update", func() {
		// TESTCASE SCENARIO
		// Under test: UpdateClientCertificateHandler
		// Mocked: UpdateClientCertificate
		// Set up: UpdateClientCertificate is mocked to return validation errors as if an error was encountered
		// Expected outcome: The handler should see the validation errors and also return an error
		params := clientcertop.UpdateClientCertificateParams{
			HTTPRequest: suite.setupAuthenticatedRequest("PUT", fmt.Sprintf("/client-certificates/%s", clientCert.ID)),
			ClientCertificate: &adminmessages.ClientCertificateUpdate{
				AllowPrime:       &clientCert.AllowPrime,
				AllowPPTAS:       &clientCert.AllowPPTAS,
				PptasAffiliation: (*adminmessages.Affiliation)(clientCert.PPTASAffiliation),
			},
		}

		err := validate.NewErrors()
		clientCertUpdater := &mocks.ClientCertUpdater{}
		clientCertUpdater.On("UpdateClientCert",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything,
			params.ClientCertificate,
		).Return(nil, err, nil).Once()

		handler := UpdateClientCertHandler{
			suite.NewHandlerConfig(),
			clientCertUpdater,
			newQueryFilter,
		}

		response := handler.Handle(params)
		suite.IsType(&clientcertop.UpdateClientCertificateBadRequest{}, response)
	})
}

func (suite *HandlerSuite) TestDeleteClientCertificateHandler() {
	clientCert := models.ClientCert{
		ID:         uuid.Must(uuid.NewV4()),
		AllowPrime: false,
		UserID:     uuid.Nil,
	}
	queryFilter := mocks.QueryFilter{}
	newQueryFilter := newMockQueryFilterBuilder(&queryFilter)
	suite.Run("Successful delete", func() {
		params := clientcertop.RemoveClientCertificateParams{
			HTTPRequest:         suite.setupAuthenticatedRequest("DELETE", fmt.Sprintf("/client-certificates/%s", clientCert.ID)),
			ClientCertificateID: strfmt.UUID(clientCert.ID.String()),
		}

		clientCertRemover := &mocks.ClientCertRemover{}
		clientCertRemover.On("RemoveClientCert",
			mock.AnythingOfType("*appcontext.appContext"),
			clientCert.ID,
		).Return(&clientCert, nil, nil).Once()
		handler := RemoveClientCertHandler{
			suite.NewHandlerConfig(),
			clientCertRemover,
			newQueryFilter,
		}
		response := handler.Handle(params)
		suite.IsType(&clientcertop.UpdateClientCertificateOK{}, response)

	})

	suite.Run("Failed update", func() {
		// TESTCASE SCENARIO
		// Under test: DeleteClientCertificateHandler
		// Mocked: DeleteClientCertificate
		// Set up: DeleteClientCertificate is mocked to return validation errors as if an error was encountered
		// Expected outcome: The handler should see the validation errors and also return an error

		params := clientcertop.RemoveClientCertificateParams{
			HTTPRequest:         suite.setupAuthenticatedRequest("DELETE", fmt.Sprintf("/client-certificates/%s", clientCert.ID)),
			ClientCertificateID: strfmt.UUID(clientCert.ID.String()),
		}

		err := validate.NewErrors()
		clientCertRemover := &mocks.ClientCertRemover{}
		clientCertRemover.On("RemoveClientCert",
			mock.AnythingOfType("*appcontext.appContext"),
			clientCert.ID,
		).Return(nil, err, nil).Once()
		handler := RemoveClientCertHandler{
			suite.NewHandlerConfig(),
			clientCertRemover,
			newQueryFilter,
		}

		response := handler.Handle(params)
		suite.IsType(&clientcertop.RemoveClientCertificateInternalServerError{}, response)
	})
}
