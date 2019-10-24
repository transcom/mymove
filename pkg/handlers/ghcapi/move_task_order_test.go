package ghcapi

import (
	"errors"
	"log"
	"net/http/httptest"

	"github.com/transcom/mymove/pkg/models"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/services/mocks"

	"github.com/transcom/mymove/pkg/gen/ghcmessages"
	"github.com/transcom/mymove/pkg/testdatagen"

	movetaskorder "github.com/transcom/mymove/pkg/services/move_task_order"

	"github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations/move_task_order"
	movetaskorderops "github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations/move_task_order"
	"github.com/transcom/mymove/pkg/handlers"
)

func (suite *HandlerSuite) TestUpdateMoveTaskOrderHandlerIntegration() {
	serviceItem := testdatagen.MakeServiceItem(suite.DB(), testdatagen.Assertions{})
	moveTaskOrder := serviceItem.MoveTaskOrder
	request := httptest.NewRequest("PATCH", "/move-task-orders/{moveTaskOrderID}/status", nil)
	params := move_task_order.UpdateMoveTaskOrderStatusParams{
		HTTPRequest:     request,
		Body:            &ghcmessages.MoveTaskOrderStatus{Status: "DRAFT"},
		MoveTaskOrderID: moveTaskOrder.ID.String(),
	}
	context := handlers.NewHandlerContext(suite.DB(), suite.TestLogger())

	// make the request
	handler := UpdateMoveTaskOrderStatusHandlerFunc{context,
		movetaskorder.NewMoveTaskOrderFetcher(suite.DB()),
	}
	response := handler.Handle(params)

	suite.IsNotErrResponse(response)
	log.Println(response)
	moveTaskOrdersResponse := response.(*movetaskorderops.UpdateMoveTaskOrderStatusOK)
	moveTaskOrdersPayload := moveTaskOrdersResponse.Payload

	suite.NotNil(moveTaskOrdersPayload)

	suite.NotNil(moveTaskOrdersPayload.ID, false)
}

func (suite *HandlerSuite) TestUpdateMoveTaskOrderHandlerNotFoundError() {
	moveTaskOrderID, _ := uuid.NewV4()
	request := httptest.NewRequest("PATCH", "/move-task-orders/{moveTaskOrderID}/status", nil)
	params := move_task_order.UpdateMoveTaskOrderStatusParams{
		HTTPRequest:     request,
		Body:            &ghcmessages.MoveTaskOrderStatus{Status: "DRAFT"},
		MoveTaskOrderID: moveTaskOrderID.String(),
	}
	context := handlers.NewHandlerContext(suite.DB(), suite.TestLogger())

	// make the request
	mtoFetcher := &mocks.MoveTaskOrderFetcher{}
	mtoFetcher.On("FetchMoveTaskOrder", moveTaskOrderID).
		Return(&models.MoveTaskOrder{}, movetaskorder.ErrNotFound{})
	handler := UpdateMoveTaskOrderStatusHandlerFunc{context, mtoFetcher}
	response := handler.Handle(params)

	suite.Assertions.IsType(&move_task_order.UpdateMoveTaskOrderStatusNotFound{}, response)
}

func (suite *HandlerSuite) TestUpdateMoveTaskOrderHandlerServerError() {
	moveTaskOrderID, _ := uuid.NewV4()
	request := httptest.NewRequest("PATCH", "/move-task-orders/{moveTaskOrderID}/status", nil)
	params := move_task_order.UpdateMoveTaskOrderStatusParams{
		HTTPRequest:     request,
		Body:            &ghcmessages.MoveTaskOrderStatus{Status: "DRAFT"},
		MoveTaskOrderID: moveTaskOrderID.String(),
	}
	context := handlers.NewHandlerContext(suite.DB(), suite.TestLogger())

	// make the request
	mtoFetcher := &mocks.MoveTaskOrderFetcher{}
	mtoFetcher.On("FetchMoveTaskOrder", moveTaskOrderID).
		Return(&models.MoveTaskOrder{}, errors.New("something bad happened"))
	handler := UpdateMoveTaskOrderStatusHandlerFunc{context, mtoFetcher}
	response := handler.Handle(params)

	suite.Assertions.IsType(&move_task_order.UpdateMoveTaskOrderStatusInternalServerError{}, response)
}
