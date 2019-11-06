package primeapi

import (
	"fmt"
	"net/http/httptest"
	"time"

	"github.com/go-openapi/strfmt"

	movetaskorder "github.com/transcom/mymove/pkg/services/move_task_order"

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

func (suite *HandlerSuite) TestUpdateMoveTaskOrderEstimatedWeightHandlerSuccess() {
	serviceItem := testdatagen.MakeServiceItem(suite.DB(), testdatagen.Assertions{})
	moveTaskOrder := serviceItem.MoveTaskOrder

	// set up what needs to be passed to handler
	request := httptest.NewRequest("PATCH", fmt.Sprintf("/move-task-orders/%s/prime-estimated-weight", moveTaskOrder.ID), nil)
	params := movetaskorderops.UpdateMoveTaskOrderEstimatedWeightParams{
		HTTPRequest: request,
		Body: movetaskorderops.UpdateMoveTaskOrderEstimatedWeightBody{
			PrimeEstimatedWeight: 1220,
		},
		MoveTaskOrderID: moveTaskOrder.ID.String(),
	}
	context := handlers.NewHandlerContext(suite.DB(), suite.TestLogger())

	// make the request
	handler := UpdateMoveTaskOrderEstimatedWeightHandler{context,
		movetaskorder.NewMoveTaskOrderEstimatedWeightUpdater(suite.DB())}
	response := handler.Handle(params)

	suite.IsNotErrResponse(response)
	updateMoveTaskOrderEstimatedWeightResponse := response.(*movetaskorderops.UpdateMoveTaskOrderEstimatedWeightOK)
	updateMoveTaskOrderEstimatedWeightPayload := updateMoveTaskOrderEstimatedWeightResponse.Payload

	suite.NotNil(updateMoveTaskOrderEstimatedWeightPayload)
	suite.NotNil(updateMoveTaskOrderEstimatedWeightPayload.PrimeEstimatedWeight)
	suite.Equal(params.Body.PrimeEstimatedWeight, *updateMoveTaskOrderEstimatedWeightPayload.PrimeEstimatedWeight)
	now := strfmt.Date(time.Now())
	suite.NotNil(updateMoveTaskOrderEstimatedWeightPayload.PrimeEstimatedWeightRecordedDate)
	suite.Equal(now.String(), updateMoveTaskOrderEstimatedWeightPayload.PrimeEstimatedWeightRecordedDate.String())
}
