package ghcapi

import (
	"net/http/httptest"

	"github.com/transcom/mymove/pkg/models"

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

	suite.Equal(getCustomerInfoPayload.FirstName, models.StringPointer("First"))
	suite.Equal(getCustomerInfoPayload.MiddleName, models.StringPointer("Middle"))
	suite.Equal(getCustomerInfoPayload.LastName, models.StringPointer("Last"))
	suite.Equal(getCustomerInfoPayload.Agency, models.StringPointer("Agency"))
	suite.Equal(getCustomerInfoPayload.Grade, models.StringPointer("Grade"))
	suite.Equal(getCustomerInfoPayload.Email, models.StringPointer("Example@example.com"))
	suite.Equal(getCustomerInfoPayload.Telephone, models.StringPointer("213-213-3232"))
	suite.Equal(getCustomerInfoPayload.OriginDutyStation, models.StringPointer("Origin Station"))
	suite.Equal(getCustomerInfoPayload.DestinationDutyStation, models.StringPointer("Destination Station"))
	suite.Equal(getCustomerInfoPayload.DependentsAuthorized, true)
}
