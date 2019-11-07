package primeapi

import (
	"fmt"
	"net/http/httptest"
	"time"

	"github.com/transcom/mymove/pkg/testdatagen"

	movetaskorderops "github.com/transcom/mymove/pkg/gen/primeapi/primeoperations/move_task_order"
	"github.com/transcom/mymove/pkg/handlers"
)

func (suite *HandlerSuite) TestListMoveTaskOrdersHandler() {
	moveTaskOrder := testdatagen.MakeMoveTaskOrder(suite.DB(), testdatagen.Assertions{})

	request := httptest.NewRequest("GET", "/move-task-orders", nil)

	params := movetaskorderops.ListMoveTaskOrdersParams{HTTPRequest: request}
	context := handlers.NewHandlerContext(suite.DB(), suite.TestLogger())

	// make the request
	handler := ListMoveTaskOrdersHandler{HandlerContext: context}
	response := handler.Handle(params)

	suite.IsNotErrResponse(response)
	moveTaskOrdersResponse := response.(*movetaskorderops.ListMoveTaskOrdersOK)
	moveTaskOrdersPayload := moveTaskOrdersResponse.Payload

	suite.Equal(1, len(moveTaskOrdersPayload))
	suite.Equal(moveTaskOrder.ID.String(), moveTaskOrdersPayload[0].ID.String())
}

func (suite *HandlerSuite) TestListMoveTaskOrdersHandlerReturnsUpdated() {
	now := time.Now()
	lastFetch := now.Add(-time.Second)

	moveTaskOrder := testdatagen.MakeMoveTaskOrder(suite.DB(), testdatagen.Assertions{})

	// this MTO should not be returned
	olderMoveTaskOrder := testdatagen.MakeMoveTaskOrder(suite.DB(), testdatagen.Assertions{})

	// Pop will overwrite UpdatedAt when saving a model, so use SQL to set it in the past
	suite.NoError(suite.DB().RawQuery("UPDATE move_task_orders SET updated_at=? WHERE id=?",
		now.Add(-2*time.Second), olderMoveTaskOrder.ID).Exec())

	since := lastFetch.Unix()
	request := httptest.NewRequest("GET", fmt.Sprintf("/move-task-orders?since=%d", lastFetch.Unix()), nil)

	params := movetaskorderops.ListMoveTaskOrdersParams{HTTPRequest: request, Since: &since}
	context := handlers.NewHandlerContext(suite.DB(), suite.TestLogger())

	// make the request
	handler := ListMoveTaskOrdersHandler{HandlerContext: context}
	response := handler.Handle(params)

	suite.IsNotErrResponse(response)
	moveTaskOrdersResponse := response.(*movetaskorderops.ListMoveTaskOrdersOK)
	moveTaskOrdersPayload := moveTaskOrdersResponse.Payload

	suite.Equal(1, len(moveTaskOrdersPayload))
	suite.Equal(moveTaskOrder.ID.String(), moveTaskOrdersPayload[0].ID.String())
}
