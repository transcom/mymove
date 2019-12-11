package ghcapi

import (
	"net/http/httptest"

	customerops "github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations/customer"

	"github.com/go-openapi/strfmt"

	"github.com/transcom/mymove/pkg/testdatagen"

	"github.com/transcom/mymove/pkg/handlers"
)

func (suite *HandlerSuite) TestGetCustomerHandlerIntegration() {
	customer := testdatagen.MakeDefaultCustomer(suite.DB())

	request := httptest.NewRequest("GET", "/customer/{customerID}", nil)
	params := customerops.GetCustomerParams{
		HTTPRequest: request,
		CustomerID:  strfmt.UUID(customer.ID.String()),
	}
	context := handlers.NewHandlerContext(suite.DB(), suite.TestLogger())
	handler := GetCustomerHandler{
		context,
	}
	response := handler.Handle(params)

	suite.IsNotErrResponse(response)
	moveTaskOrdersResponse := response.(*customerops.GetCustomerOK)
	moveTaskOrdersPayload := moveTaskOrdersResponse.Payload

	suite.Assertions.IsType(&customerops.GetCustomerOK{}, response)
	suite.Equal(strfmt.UUID(customer.ID.String()), moveTaskOrdersPayload.ID)
	suite.Equal(customer.DODID, moveTaskOrdersPayload.DodID)
	suite.Equal(strfmt.UUID(customer.UserID.String()), moveTaskOrdersPayload.UserID)
}
