package ghcapi

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/go-openapi/strfmt"
	"github.com/stretchr/testify/mock"

	mtoagentop "github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations/mto_agent"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services/mocks"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *HandlerSuite) TestListMTOAgentsHandler() {
	var requestUser models.User
	var testMTOAgent models.MTOAgent
	setupTestData := func() *http.Request {
		requestUser = testdatagen.MakeStubbedUser(suite.DB())
		testMTOAgent = testdatagen.MakeDefaultMTOAgent(suite.DB())
		req := httptest.NewRequest("GET", fmt.Sprintf("/move-task-orders/%s/mto-agents", testMTOAgent.ID.String()), nil)
		req = suite.AuthenticateAdminRequest(req, requestUser)
		return req
	}

	suite.Run("Successful Response", func() {
		req := setupTestData()
		params := mtoagentop.FetchMTOAgentListParams{
			HTTPRequest: req,
			ShipmentID:  strfmt.UUID(testMTOAgent.MTOShipmentID.String()),
		}
		listFetcher := &mocks.ListFetcher{}
		listFetcher.On("FetchRecordList",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything,
			mock.Anything,
			mock.Anything,
			mock.Anything,
			mock.Anything,
		).Return(nil).Once()

		handler := ListMTOAgentsHandler{
			HandlerConfig: suite.HandlerConfig(),
			ListFetcher:   listFetcher,
		}

		// Validate incoming payload: no body to validate

		response := handler.Handle(params)
		suite.IsType(mtoagentop.NewFetchMTOAgentListOK(), response)
		payload := response.(*mtoagentop.FetchMTOAgentListOK).Payload

		// Validate outgoing payload
		suite.NoError(payload.Validate(strfmt.Default))
	})

	suite.Run("Error Response", func() {
		req := setupTestData()
		params := mtoagentop.FetchMTOAgentListParams{
			HTTPRequest: req,
			ShipmentID:  strfmt.UUID(testMTOAgent.MTOShipmentID.String()),
		}
		listFetcher := &mocks.ListFetcher{}
		listFetcher.On("FetchRecordList",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything,
			mock.Anything,
			mock.Anything,
			mock.Anything,
			mock.Anything,
		).Return(errors.New("an error happened")).Once()

		handler := ListMTOAgentsHandler{
			HandlerConfig: suite.HandlerConfig(),
			ListFetcher:   listFetcher,
		}

		// Validate incoming payload: no body to validate

		response := handler.Handle(params)
		expectedResponse := mtoagentop.NewFetchMTOAgentListInternalServerError()
		suite.Equal(expectedResponse, response)
		payload := response.(*mtoagentop.FetchMTOAgentListInternalServerError).Payload

		// Validate outgoing payload: nil payload
		suite.Nil(payload)
	})

	suite.Run("404 Response", func() {
		req := setupTestData()
		params := mtoagentop.FetchMTOAgentListParams{
			HTTPRequest: req,
			ShipmentID:  strfmt.UUID(testMTOAgent.MTOShipmentID.String()),
		}
		listFetcher := &mocks.ListFetcher{}
		listFetcher.On("FetchRecordList",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything,
			mock.Anything,
			mock.Anything,
			mock.Anything,
			mock.Anything,
		).Return(models.ErrFetchNotFound).Once()

		handler := ListMTOAgentsHandler{
			HandlerConfig: suite.HandlerConfig(),
			ListFetcher:   listFetcher,
		}

		// Validate incoming payload: no body to validate

		response := handler.Handle(params)
		expectedResponse := mtoagentop.NewFetchMTOAgentListNotFound()
		suite.Equal(expectedResponse, response)
		payload := response.(*mtoagentop.FetchMTOAgentListNotFound).Payload

		// Validate outgoing payload: nil payload
		suite.Nil(payload)
	})
}
