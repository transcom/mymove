package supportapi

import (
	"net/http/httptest"

	"github.com/transcom/mymove/pkg/etag"
	"github.com/transcom/mymove/pkg/services/query"

	"github.com/go-openapi/strfmt"

	"github.com/transcom/mymove/pkg/gen/supportapi/supportoperations/move_task_order"
	movetaskorderops "github.com/transcom/mymove/pkg/gen/supportapi/supportoperations/move_task_order"
	"github.com/transcom/mymove/pkg/handlers"
	movetaskorder "github.com/transcom/mymove/pkg/services/move_task_order"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *HandlerSuite) TestUpdateMoveTaskOrderHandlerIntegrationSuccess() {
	moveTaskOrder := testdatagen.MakeMoveTaskOrder(suite.DB(), testdatagen.Assertions{})

	request := httptest.NewRequest("PATCH", "/move-task-orders/{moveTaskOrderID}/status", nil)
	params := move_task_order.UpdateMoveTaskOrderStatusParams{
		HTTPRequest:     request,
		MoveTaskOrderID: moveTaskOrder.ID.String(),
		IfMatch:         etag.GenerateEtag(moveTaskOrder.UpdatedAt),
	}
	context := handlers.NewHandlerContext(suite.DB(), suite.TestLogger())
	queryBuilder := query.NewQueryBuilder(suite.DB())

	// make the request
	handler := UpdateMoveTaskOrderStatusHandlerFunc{context,
		movetaskorder.NewMoveTaskOrderUpdater(suite.DB(), queryBuilder),
	}
	response := handler.Handle(params)

	suite.IsNotErrResponse(response)
	moveTaskOrdersResponse := response.(*movetaskorderops.UpdateMoveTaskOrderStatusOK)
	moveTaskOrdersPayload := moveTaskOrdersResponse.Payload

	suite.Assertions.IsType(&move_task_order.UpdateMoveTaskOrderStatusOK{}, response)
	suite.Equal(moveTaskOrdersPayload.ID, strfmt.UUID(moveTaskOrder.ID.String()))
	suite.Equal(*moveTaskOrdersPayload.IsAvailableToPrime, true)
}

func (suite *HandlerSuite) TestGetMoveTaskOrder() {
	moveTaskOrder := testdatagen.MakeMoveTaskOrder(suite.DB(), testdatagen.Assertions{})

	request := httptest.NewRequest("GET", "/move-task-orders/{moveTaskOrderID}", nil)
	params := move_task_order.GetMoveTaskOrderParams{
		HTTPRequest:     request,
		MoveTaskOrderID: moveTaskOrder.ID.String(),
	}

	context := handlers.NewHandlerContext(suite.DB(), suite.TestLogger())
	handler := GetMoveTaskOrderHandlerFunc{context,
		movetaskorder.NewMoveTaskOrderFetcher(suite.DB()),
	}
	response := handler.Handle(params)
	suite.IsNotErrResponse(response)
	moveTaskOrdersResponse := response.(*movetaskorderops.GetMoveTaskOrderOK)
	moveTaskOrdersPayload := moveTaskOrdersResponse.Payload

	suite.Assertions.IsType(&move_task_order.GetMoveTaskOrderOK{}, response)
	suite.Equal(moveTaskOrdersPayload.ID, strfmt.UUID(moveTaskOrder.ID.String()))
	suite.Equal(*moveTaskOrdersPayload.IsAvailableToPrime, false)
}
