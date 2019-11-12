package ghcapi

import (
	"fmt"
	"net/http/httptest"

	"github.com/gofrs/uuid"

	entitlementscodeop "github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations/entitlements"
	movetaskordercodeop "github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations/move_task_order"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	movetaskorder "github.com/transcom/mymove/pkg/services/move_task_order"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *HandlerSuite) TestUpdateMoveTaskOrderActualWeightHandlerIntegration() {
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

func (suite *HandlerSuite) TestGetPrimeEntitlementsHandlerIntegration() {
	// set up what needs to be passed to handler
	moveTaskOrderID, _ := uuid.NewV4()
	mto := testdatagen.MakeMoveTaskOrder(suite.DB(), testdatagen.Assertions{
		MoveTaskOrder: models.MoveTaskOrder{ID: moveTaskOrderID},
	})
	testdatagen.MakeEntitlement(suite.DB(), testdatagen.Assertions{
		GHCEntitlement: models.GHCEntitlement{MoveTaskOrder: &mto}},
	)
	request := httptest.NewRequest("GET", "/move-task-orders/move_task_order_id/prime-entitlements", nil)
	params := entitlementscodeop.GetPrimeEntitlementsParams{
		HTTPRequest:     request,
		MoveTaskOrderID: mto.ID.String(),
	}
	context := handlers.NewHandlerContext(suite.DB(), suite.TestLogger())

	// make the request
	handler := GetPrimeEntitlementsHandler{context,
		movetaskorder.NewMoveTaskOrderFetcher(suite.DB())}
	response := handler.Handle(params)

	suite.IsNotErrResponse(response)
	suite.Assertions.IsType(&entitlementscodeop.GetPrimeEntitlementsOK{}, response)
	getPrimeEntitlementsResponse := response.(*entitlementscodeop.GetPrimeEntitlementsOK)
	getPrimeEntitlementsPayload := getPrimeEntitlementsResponse.Payload

	suite.NotNil(getPrimeEntitlementsPayload)

	suite.Equal(getPrimeEntitlementsPayload.DependentsAuthorized, true)
	suite.Equal(*getPrimeEntitlementsPayload.NonTemporaryStorage, true)
	suite.Equal(*getPrimeEntitlementsPayload.PrivatelyOwnedVehicle, true)
	suite.Equal(int(getPrimeEntitlementsPayload.ProGearWeight), 100)
	suite.Equal(int(getPrimeEntitlementsPayload.ProGearWeightSpouse), 200)
	suite.Equal(int(getPrimeEntitlementsPayload.StorageInTransit), 2)
	suite.Equal(int(getPrimeEntitlementsPayload.TotalDependents), 1)
	suite.Equal(int(getPrimeEntitlementsPayload.TotalWeightSelf), 0)
}
