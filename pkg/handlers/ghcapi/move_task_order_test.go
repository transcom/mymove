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

	entitlementscodeop "github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations/entitlements"
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

func (suite *HandlerSuite) TestGetEntitlementsHandlerIntegration() {
	// set up what needs to be passed to handler
	moveTaskOrderID, _ := uuid.NewV4()
	mto := testdatagen.MakeMoveTaskOrder(suite.DB(), testdatagen.Assertions{
		MoveTaskOrder: models.MoveTaskOrder{ID: moveTaskOrderID},
	})
	testdatagen.MakeEntitlement(suite.DB(), testdatagen.Assertions{
		GHCEntitlement: models.GHCEntitlement{MoveTaskOrder: &mto}},
	)
	request := httptest.NewRequest("GET", "/move-task-orders/move_task_order_id/entitlements", nil)
	params := entitlementscodeop.GetEntitlementsParams{
		HTTPRequest:     request,
		MoveTaskOrderID: mto.ID.String(),
	}
	context := handlers.NewHandlerContext(suite.DB(), suite.TestLogger())

	// make the request
	handler := GetEntitlementsHandler{context,
		movetaskorder.NewMoveTaskOrderFetcher(suite.DB())}
	response := handler.Handle(params)

	suite.IsNotErrResponse(response)
	suite.Assertions.IsType(&entitlementscodeop.GetEntitlementsOK{}, response)
	getEntitlementsResponse := response.(*entitlementscodeop.GetEntitlementsOK)
	getEntitlementsPayload := getEntitlementsResponse.Payload

	suite.NotNil(getEntitlementsPayload)

	suite.Equal(getEntitlementsPayload.DependentsAuthorized, true)
	suite.Equal(*getEntitlementsPayload.NonTemporaryStorage, true)
	suite.Equal(*getEntitlementsPayload.PrivatelyOwnedVehicle, true)
	suite.Equal(int(getEntitlementsPayload.ProGearWeight), 100)
	suite.Equal(int(getEntitlementsPayload.ProGearWeightSpouse), 200)
	suite.Equal(int(getEntitlementsPayload.StorageInTransit), 2)
	suite.Equal(int(getEntitlementsPayload.TotalDependents), 1)
	suite.Equal(int(getEntitlementsPayload.TotalWeightSelf), 0)
}
