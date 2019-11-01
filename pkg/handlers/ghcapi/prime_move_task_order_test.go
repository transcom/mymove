package ghcapi

import (
	"fmt"
	"net/http/httptest"

	movetaskordercodeop "github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations/move_task_order"
	"github.com/transcom/mymove/pkg/handlers"
	movetaskorder "github.com/transcom/mymove/pkg/services/move_task_order"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *HandlerSuite) TestUpdateMoveTaskOrderActualWeightHandler_Success() {
	serviceItem := testdatagen.MakeServiceItem(suite.DB(), testdatagen.Assertions{})
	moveTaskOrder := serviceItem.MoveTaskOrder

	// set up what needs to be passed to handler
	request := httptest.NewRequest("PATCH", fmt.Sprintf("/move-task-orders/%s/prime-actual-weight", moveTaskOrder.ID), nil)
	params := movetaskordercodeop.UpdateMoveTaskOrderActualWeightParams{
		HTTPRequest:     request,
		Body:            movetaskordercodeop.UpdateMoveTaskOrderActualWeightBody{ActualWeight: 2819},
		MoveTaskOrderID: moveTaskOrder.ID.String(),
	}
	context := handlers.NewHandlerContext(suite.DB(), suite.TestLogger())

	// make the request
	handler := UpdateMoveTaskOrderActualWeightHandler{context,
		movetaskorder.NewMoveTaskOrderActualWeightUpdater(suite.DB())}
	response := handler.Handle(params)

	suite.IsNotErrResponse(response)
	updateMoveTaskOrderActualWeightResponse := response.(*movetaskordercodeop.UpdateMoveTaskOrderActualWeightOK)
	updateMoveTaskOrderActualWeightPayload := updateMoveTaskOrderActualWeightResponse.Payload

	suite.NotNil(updateMoveTaskOrderActualWeightPayload)
	suite.Equal(int(updateMoveTaskOrderActualWeightPayload.ActualWeight), 2819)
}
