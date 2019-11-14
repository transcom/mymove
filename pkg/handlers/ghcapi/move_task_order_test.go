package ghcapi

import (
	"errors"
	"net/http/httptest"

	"github.com/go-openapi/strfmt"

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
		movetaskorder.NewMoveTaskOrderStatusUpdater(suite.DB()),
	}
	response := handler.Handle(params)

	suite.IsNotErrResponse(response)
	moveTaskOrdersResponse := response.(*movetaskorderops.UpdateMoveTaskOrderStatusOK)
	moveTaskOrdersPayload := moveTaskOrdersResponse.Payload

	suite.Assertions.IsType(&move_task_order.UpdateMoveTaskOrderStatusOK{}, response)
	suite.Equal(moveTaskOrdersPayload.ID, strfmt.UUID(moveTaskOrder.ID.String()))
	suite.Equal(moveTaskOrdersPayload.Status, "DRAFT")
	suite.Regexp("^\\d{4}-\\d{4}$", moveTaskOrdersPayload.ReferenceID)
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
	mtoStatusUpdater := &mocks.MoveTaskOrderStatusUpdater{}
	mtoStatusUpdater.On("UpdateMoveTaskOrderStatus", moveTaskOrderID, models.MoveTaskOrderStatusDraft).
		Return(&models.MoveTaskOrder{}, movetaskorder.ErrNotFound{})
	handler := UpdateMoveTaskOrderStatusHandlerFunc{context, mtoStatusUpdater}
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
	mtoStatusUpdater := &mocks.MoveTaskOrderStatusUpdater{}
	mtoStatusUpdater.On("UpdateMoveTaskOrderStatus", moveTaskOrderID, models.MoveTaskOrderStatusDraft).
		Return(&models.MoveTaskOrder{}, errors.New("something bad happened"))
	handler := UpdateMoveTaskOrderStatusHandlerFunc{context, mtoStatusUpdater}
	response := handler.Handle(params)

	suite.Assertions.IsType(&move_task_order.UpdateMoveTaskOrderStatusInternalServerError{}, response)
}
