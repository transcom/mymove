package ghcapi

import (
	"net/http/httptest"

	entitlementscodeop "github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations/entitlements"
	"github.com/transcom/mymove/pkg/handlers"
)

func (suite *HandlerSuite) TestGetEntitlementsHandler_Success() {
	// set up what needs to be passed to handler
	request := httptest.NewRequest("GET", "/move-task-orders/move_task_order_id/entitlements", nil)
	params := entitlementscodeop.GetEntitlementsParams{
		HTTPRequest: request,
	}
	context := handlers.NewHandlerContext(suite.DB(), suite.TestLogger())

	// make the request
	handler := GetEntitlementsHandler{context}
	response := handler.Handle(params)

	suite.IsNotErrResponse(response)
	getEntitlementsResponse := response.(*entitlementscodeop.GetEntitlementsOK)
	getEntitlementsPayload := getEntitlementsResponse.Payload

	suite.NotNil(getEntitlementsPayload)

	suite.Equal(getEntitlementsPayload.DependentsAuthorized, false)
	suite.Equal(getEntitlementsPayload.NonTemporaryStorage, false)
	suite.Equal(getEntitlementsPayload.PrivatelyOwnedVehicle, true)
	suite.Equal(int(getEntitlementsPayload.ProGearWeight), 200)
	suite.Equal(int(getEntitlementsPayload.ProGearWeightSpouse), 100)
	suite.Equal(int(getEntitlementsPayload.StorageInTransit), 1000)
	suite.Equal(int(getEntitlementsPayload.TotalDependents), 3)
	suite.Equal(int(getEntitlementsPayload.TotalWeightSelf), 1300)
}
