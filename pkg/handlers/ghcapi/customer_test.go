package ghcapi

import (
	"net/http/httptest"

	customercodeop "github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations/customer"
	"github.com/transcom/mymove/pkg/handlers"
)

func (suite *HandlerSuite) TestGetCustomerInfoHandler_Success() {
	// set up what needs to be passed to handler
	request := httptest.NewRequest("GET", "/move-task-orders/move_task_order_id/entitlements", nil)
	params := customercodeop.GetCustomerInfoParams{
		HTTPRequest: request,
	}
	context := handlers.NewHandlerContext(suite.DB(), suite.TestLogger())

	// make the request
	handler := GetCustomerInfoHandler{context}
	response := handler.Handle(params)

	suite.IsNotErrResponse(response)
	getCustomerInfo := response.(*customercodeop.GetCustomerInfoOK)
	getCustomerInfoPayload := getCustomerInfo.Payload

	suite.NotNil(getCustomerInfoPayload)

	suite.Equal(getCustomerInfoPayload.FirstName, "First")
	suite.Equal(getCustomerInfoPayload.MiddleName, "Middle")
	suite.Equal(getCustomerInfoPayload.LastName, "Last")
	suite.Equal(getCustomerInfoPayload.Agency, "Agency")
	suite.Equal(getCustomerInfoPayload.Grade, "Grade")
	suite.Equal(getCustomerInfoPayload.Email, "Example@example.com")
	suite.Equal(getCustomerInfoPayload.Telephone, "213-213-3232")
	suite.Equal(getCustomerInfoPayload.OriginDutyStation, "Origin Station")
	suite.Equal(getCustomerInfoPayload.DestinationDutyStation, "Destination Station")
	suite.Equal(getCustomerInfoPayload.DependentsAuthorized, true)
}
