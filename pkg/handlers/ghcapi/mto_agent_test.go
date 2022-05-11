package ghcapi

import (
	"errors"
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/go-openapi/strfmt"
	"github.com/stretchr/testify/mock"

	mtoagentop "github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations/mto_agent"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services/mocks"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *HandlerSuite) TestListMTOAgentsHandler() {
	testMTOAgent := testdatagen.MakeDefaultMTOAgent(suite.DB())

	requestUser := testdatagen.MakeStubbedUser(suite.DB())
	req := httptest.NewRequest("GET", fmt.Sprintf("/move-task-orders/%s/mto-agents", testMTOAgent.ID.String()), nil)
	req = suite.AuthenticateAdminRequest(req, requestUser)

	suite.T().Run("Successful Response", func(t *testing.T) {
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
			HandlerContext: handlers.NewHandlerContext(suite.DB(), suite.Logger()),
			ListFetcher:    listFetcher,
		}

		response := handler.Handle(params)
		suite.IsType(mtoagentop.NewFetchMTOAgentListOK(), response)
	})

	suite.T().Run("Error Response", func(t *testing.T) {
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
			HandlerContext: handlers.NewHandlerContext(suite.DB(), suite.Logger()),
			ListFetcher:    listFetcher,
		}
		response := handler.Handle(params)
		expectedResponse := mtoagentop.NewFetchMTOAgentListInternalServerError()
		suite.Equal(expectedResponse, response)
	})

	suite.T().Run("404 Response", func(t *testing.T) {
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
			HandlerContext: handlers.NewHandlerContext(suite.DB(), suite.Logger()),
			ListFetcher:    listFetcher,
		}
		response := handler.Handle(params)
		expectedResponse := mtoagentop.NewFetchMTOAgentListNotFound()
		suite.Equal(expectedResponse, response)
	})
}
