package ghcapi

import (
	"net/http/httptest"

	customerservice "github.com/transcom/mymove/pkg/services/office_user/customer"

	customerops "github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations/customer"

	"github.com/go-openapi/strfmt"

	"github.com/transcom/mymove/pkg/testdatagen"

	"github.com/transcom/mymove/pkg/handlers"
)

func (suite *HandlerSuite) TestGetCustomerHandlerIntegration() {
	customer := testdatagen.MakeDefaultServiceMember(suite.DB())

	request := httptest.NewRequest("GET", "/customer/{customerID}", nil)
	params := customerops.GetCustomerParams{
		HTTPRequest: request,
		CustomerID:  strfmt.UUID(customer.ID.String()),
	}
	context := handlers.NewHandlerContext(suite.DB(), suite.TestLogger())
	handler := GetCustomerHandler{
		context,
		customerservice.NewCustomerFetcher(suite.DB()),
	}

	response := handler.Handle(params)
	suite.IsNotErrResponse(response)

	getCustomerResponse := response.(*customerops.GetCustomerOK)
	getCustomerPayload := getCustomerResponse.Payload
	suite.Assertions.IsType(&customerops.GetCustomerOK{}, response)
	suite.Equal(strfmt.UUID(customer.ID.String()), getCustomerPayload.ID)
	suite.Equal(*customer.Edipi, getCustomerPayload.DodID)
	suite.Equal(strfmt.UUID(customer.UserID.String()), getCustomerPayload.UserID)
	suite.Equal(customer.Affiliation.String(), getCustomerPayload.Agency)
	suite.Equal(customer.PersonalEmail, getCustomerPayload.Email)
	suite.Equal(customer.Telephone, getCustomerPayload.Phone)
	suite.NotZero(getCustomerPayload.CurrentAddress)
}
