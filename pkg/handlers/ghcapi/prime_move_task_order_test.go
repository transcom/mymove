package ghcapi

import (
	"fmt"
	"net/http/httptest"

	movetaskordercodeop "github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations/move_task_order"
	"github.com/transcom/mymove/pkg/gen/ghcmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *HandlerSuite) TestUpdateMoveTaskOrderActualWeightHandler_Success() {
	// create the move task order to test with
	mto := testdatagen.MakeDefaultMoveTaskOrder(suite.DB())

	// set up what needs to be passed to handler
	request := httptest.NewRequest("PATCH", fmt.Sprintf("/move-task-orders/%s/prime-actual-weight", "dlafksd"), nil)
	actualWeightPayload := ghcmessages.PatchActualWeight{
		ActualWeight: 2819,
	}
	params := movetaskordercodeop.UpdateMoveTaskOrderActualWeightParams{
		HTTPRequest:       request,
		PatchActualWeight: &actualWeightPayload,
		MoveTaskOrderID:   mto.ID.String(),
	}
	context := handlers.NewHandlerContext(suite.DB(), suite.TestLogger())

	// make the request
	handler := UpdateMoveTaskOrderActualWeightHandler{context}
	response := handler.Handle(params)

	suite.IsNotErrResponse(response)
	updateMoveTaskOrderActualWeightResponse := response.(*movetaskordercodeop.UpdateMoveTaskOrderActualWeightOK)
	updateMoveTaskOrderActualWeightPayload := updateMoveTaskOrderActualWeightResponse.Payload

	suite.NotNil(updateMoveTaskOrderActualWeightPayload)
	suite.Equal(int(updateMoveTaskOrderActualWeightPayload.ActualWeight), 2819)
}
