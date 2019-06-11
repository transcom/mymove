package publicapi

import (
	"fmt"
	"net/http/httptest"
	"strings"

	"github.com/transcom/mymove/mocks"
	postalcodesops "github.com/transcom/mymove/pkg/gen/restapi/apioperations/postal_codes"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *HandlerSuite) TestValidatePostalCodeHandler_Valid() {
	// create user
	user := testdatagen.MakeDefaultUser(suite.DB())

	postalCode := "30813"
	postalCodeTypeString := "Destination"
	postalCodeType := services.PostalCodeType(postalCodeTypeString)

	// makes request
	request := httptest.NewRequest("GET", fmt.Sprintf("/postal_codes/%s", postalCode), strings.NewReader("postal_code_type=origin"))
	request = suite.AuthenticateUserRequest(request, user)

	params := postalcodesops.ValidatePostalCodeParams{
		HTTPRequest:    request,
		PostalCode:     postalCode,
		PostalCodeType: postalCodeTypeString,
	}

	context := handlers.NewHandlerContext(suite.DB(), suite.TestLogger())
	postalCodeValidator := &mocks.PostalCodeValidator{}
	postalCodeValidator.On("ValidatePostalCode",
		postalCode,
		postalCodeType,
	).Return(true, nil)

	handler := ValidatePostalCodeHandler{context, postalCodeValidator}
	response := handler.Handle(params)

	suite.IsNotErrResponse(response)
	validatePostalCodeResponse := response.(*postalcodesops.ValidatePostalCodeOK)
	validatePostalCodePayload := validatePostalCodeResponse.Payload

	suite.NotNil(validatePostalCodePayload.PostalCode)
	suite.NotNil(validatePostalCodePayload.PostalCodeType)
	suite.True(*validatePostalCodePayload.Valid)
	suite.Assertions.IsType(&postalcodesops.ValidatePostalCodeOK{}, response)
}

func (suite *HandlerSuite) TestValidatePostalCodeHandler_Invalid() {
	// create user
	user := testdatagen.MakeDefaultUser(suite.DB())

	postalCode := "00000"
	postalCodeTypeString := "Destination"
	postalCodeType := services.PostalCodeType(postalCodeTypeString)

	// makes request
	request := httptest.NewRequest("GET", fmt.Sprintf("/postal_codes/%s", postalCode), strings.NewReader("postal_code_type=origin"))
	request = suite.AuthenticateUserRequest(request, user)

	params := postalcodesops.ValidatePostalCodeParams{
		HTTPRequest:    request,
		PostalCode:     postalCode,
		PostalCodeType: postalCodeTypeString,
	}

	context := handlers.NewHandlerContext(suite.DB(), suite.TestLogger())
	postalCodeValidator := &mocks.PostalCodeValidator{}
	postalCodeValidator.On("ValidatePostalCode",
		postalCode,
		postalCodeType,
	).Return(false, nil)

	handler := ValidatePostalCodeHandler{context, postalCodeValidator}
	response := handler.Handle(params)

	suite.IsNotErrResponse(response)
	validatePostalCodeResponse := response.(*postalcodesops.ValidatePostalCodeOK)
	validatePostalCodePayload := validatePostalCodeResponse.Payload

	suite.NotNil(validatePostalCodePayload.PostalCode)
	suite.NotNil(validatePostalCodePayload.PostalCodeType)
	suite.False(*validatePostalCodePayload.Valid)
	suite.Assertions.IsType(&postalcodesops.ValidatePostalCodeOK{}, response)
}
