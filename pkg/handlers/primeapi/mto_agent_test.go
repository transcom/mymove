package primeapi

import (
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/go-openapi/strfmt"

	"github.com/transcom/mymove/pkg/etag"
	mtoshipmentops "github.com/transcom/mymove/pkg/gen/primeapi/primeoperations/mto_shipment"
	"github.com/transcom/mymove/pkg/handlers/primeapi/payloads"
	mtoagent "github.com/transcom/mymove/pkg/services/mto_agent"

	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
)

func (suite *HandlerSuite) TestUpdateMTOAgentHandler() {
	agent := models.MTOAgent{
		MTOAgentType: models.MTOAgentReleasing,
	}

	// Create handler
	handler := UpdateMTOAgentHandler{
		handlers.NewHandlerContext(suite.DB(), suite.TestLogger()),
		mtoagent.NewMTOAgentUpdater(suite.DB()),
	}

	// Test a successful request + update

	// Test invalid IDs in the body vs. path values

	// Test stale eTag

	// Test not found response

	// Test not Prime-available (not found response)

	// Test invalid input

	suite.T().Run("NotImplemented response", func(t *testing.T) {
		payload := payloads.MTOAgent(&agent)
		req := httptest.NewRequest("PUT", fmt.Sprintf("/mto-shipments/%s/agents/%s", agent.MTOShipmentID.String(), agent.ID.String()), nil)
		params := mtoshipmentops.UpdateMTOAgentParams{
			HTTPRequest:   req,
			AgentID:       *handlers.FmtUUID(agent.ID),
			MtoShipmentID: *handlers.FmtUUID(agent.MTOShipmentID),
			Body:          payload,
			IfMatch:       etag.GenerateEtag(agent.UpdatedAt),
		}
		// Run swagger validations
		suite.NoError(params.Body.Validate(strfmt.Default))

		// Run handler and check response
		response := handler.Handle(params)
		suite.NotNil(response)
	})
}
