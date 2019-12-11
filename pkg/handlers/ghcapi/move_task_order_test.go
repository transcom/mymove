package ghcapi

import (
	"net/http/httptest"

	"github.com/go-openapi/strfmt"

	"github.com/transcom/mymove/pkg/testdatagen"

	movetaskorder "github.com/transcom/mymove/pkg/services/move_task_order"

	"github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations/move_task_order"
	movetaskorderops "github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations/move_task_order"
	"github.com/transcom/mymove/pkg/handlers"
)

func (suite *HandlerSuite) TestGetMoveTaskOrderHandlerIntegration() {
	moveOrder := testdatagen.MakeMoveOrder(suite.DB(), testdatagen.Assertions{})
	moveTaskOrder := testdatagen.MakeMoveTaskOrder(suite.DB(), testdatagen.Assertions{
		MoveOrder: moveOrder,
	})
	request := httptest.NewRequest("GET", "/move-task-orders/{moveTaskOrderID}", nil)
	params := move_task_order.GetMoveTaskOrderParams{
		HTTPRequest:     request,
		MoveTaskOrderID: moveTaskOrder.ID.String(),
	}
	context := handlers.NewHandlerContext(suite.DB(), suite.TestLogger())

	// make the request
	handler := GetMoveTaskOrderHandler{
		context,
		movetaskorder.NewMoveTaskOrderFetcher(suite.DB()),
	}
	response := handler.Handle(params)

	suite.IsNotErrResponse(response)
	moveTaskOrdersResponse := response.(*movetaskorderops.GetMoveTaskOrderOK)
	moveTaskOrdersPayload := moveTaskOrdersResponse.Payload

	suite.Assertions.IsType(&move_task_order.GetMoveTaskOrderOK{}, response)
	suite.Equal(strfmt.UUID(moveTaskOrder.ID.String()), moveTaskOrdersPayload.ID)
	suite.False(*moveTaskOrdersPayload.IsAvailableToPrime)
	suite.False(*moveTaskOrdersPayload.IsCanceled)
	suite.Equal(strfmt.UUID(moveTaskOrder.MoveOrderID.String()), moveTaskOrdersPayload.MoveOrdersID)
	suite.Nil(moveTaskOrdersPayload.ReferenceID)
}
