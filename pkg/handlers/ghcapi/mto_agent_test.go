package ghcapi

import (
	"errors"
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/go-openapi/strfmt"
	"github.com/stretchr/testify/mock"

	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services/mocks"
	"github.com/transcom/mymove/pkg/testdatagen"

	mtoagentop "github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations/mto_agent"
)

func (suite *HandlerSuite) TestListMTOAgentsHandler() {
	testMTOAgent := testdatagen.MakeMTOAgent(suite.DB(), testdatagen.Assertions{})

	requestUser := testdatagen.MakeDefaultUser(suite.DB())
	req := httptest.NewRequest("GET", fmt.Sprintf("/move-task-orders/%s/mto-agents", testMTOAgent.ID.String()), nil)
	req = suite.AuthenticateAdminRequest(req, requestUser)

	suite.T().Run("Successful Response", func(t *testing.T) {
		params := mtoagentop.FetchMTOAgentListParams{
			HTTPRequest: req,
			ShipmentID:  strfmt.UUID(testMTOAgent.MTOShipmentID.String()),
		}
		mtoAgentListFetcher := &mocks.MTOAgentListFetcher{}
		mtoAgentListFetcher.On("FetchMTOAgentList",
			mock.Anything,
		).Return(&models.MTOAgents{testMTOAgent}, nil).Once()

		handler := ListMTOAgentsHandler{
			HandlerContext:      handlers.NewHandlerContext(suite.DB(), suite.TestLogger()),
			MTOAgentListFetcher: mtoAgentListFetcher,
		}

		response := handler.Handle(params)
		suite.IsType(mtoagentop.NewFetchMTOAgentListOK(), response)
		okResponse := response.(*mtoagentop.FetchMTOAgentListOK)
		suite.Len(okResponse.Payload, 1)
		suite.Equal(testMTOAgent.ID.String(), okResponse.Payload[0].ID.String())
	})

	suite.T().Run("Error Response", func(t *testing.T) {
		params := mtoagentop.FetchMTOAgentListParams{
			HTTPRequest: req,
			ShipmentID:  strfmt.UUID(testMTOAgent.MTOShipmentID.String()),
		}
		mtoAgentListFetcher := &mocks.MTOAgentListFetcher{}
		mtoAgentListFetcher.On("FetchMTOAgentList",
			mock.Anything,
		).Return(nil, errors.New("an error happened")).Once()

		handler := ListMTOAgentsHandler{
			HandlerContext:      handlers.NewHandlerContext(suite.DB(), suite.TestLogger()),
			MTOAgentListFetcher: mtoAgentListFetcher,
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
		mtoAgentListFetcher := &mocks.MTOAgentListFetcher{}
		mtoAgentListFetcher.On("FetchMTOAgentList",
			mock.Anything,
		).Return(nil, nil).Once()

		handler := ListMTOAgentsHandler{
			HandlerContext:      handlers.NewHandlerContext(suite.DB(), suite.TestLogger()),
			MTOAgentListFetcher: mtoAgentListFetcher,
		}
		response := handler.Handle(params)
		expectedResponse := mtoagentop.NewFetchMTOAgentListNotFound()
		suite.Equal(expectedResponse, response)
	})
}