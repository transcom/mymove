package internalapi

import (
	"fmt"
	"net/http/httptest"
	"strings"

	postalcodesops "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/postal_codes"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/mocks"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *HandlerSuite) TestValidatePostalCodeWithRateDataHandler_Valid() {
	// create user
	user := testdatagen.MakeStubbedUser(suite.DB())

	postalCode := "30813"
	postalCodeTypeString := "Destination"
	postalCodeType := services.PostalCodeType(postalCodeTypeString)

	// makes request
	request := httptest.NewRequest("GET", fmt.Sprintf("/postal_codes/%s", postalCode), strings.NewReader("postal_code_type=origin"))
	request = suite.AuthenticateUserRequest(request, user)

	params := postalcodesops.ValidatePostalCodeWithRateDataParams{
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

	handler := ValidatePostalCodeWithRateDataHandler{context, postalCodeValidator}
	response := handler.Handle(params)

	suite.IsNotErrResponse(response)
	validatePostalCodeResponse := response.(*postalcodesops.ValidatePostalCodeWithRateDataOK)
	validatePostalCodePayload := validatePostalCodeResponse.Payload

	suite.NotNil(validatePostalCodePayload.PostalCode)
	suite.NotNil(validatePostalCodePayload.PostalCodeType)
	suite.True(*validatePostalCodePayload.Valid)
	suite.Assertions.IsType(&postalcodesops.ValidatePostalCodeWithRateDataOK{}, response)
}

func (suite *HandlerSuite) TestValidatePostalCodeWithRateDataHandler_Invalid() {
	// create user
	user := testdatagen.MakeStubbedUser(suite.DB())

	postalCode := "00000"
	postalCodeTypeString := "Destination"
	postalCodeType := services.PostalCodeType(postalCodeTypeString)

	// makes request
	request := httptest.NewRequest("GET", fmt.Sprintf("/postal_codes/%s", postalCode), strings.NewReader("postal_code_type=origin"))
	request = suite.AuthenticateUserRequest(request, user)

	params := postalcodesops.ValidatePostalCodeWithRateDataParams{
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

	handler := ValidatePostalCodeWithRateDataHandler{context, postalCodeValidator}
	response := handler.Handle(params)

	suite.IsNotErrResponse(response)
	validatePostalCodeResponse := response.(*postalcodesops.ValidatePostalCodeWithRateDataOK)
	validatePostalCodePayload := validatePostalCodeResponse.Payload

	suite.NotNil(validatePostalCodePayload.PostalCode)
	suite.NotNil(validatePostalCodePayload.PostalCodeType)
	suite.False(*validatePostalCodePayload.Valid)
	suite.Assertions.IsType(&postalcodesops.ValidatePostalCodeWithRateDataOK{}, response)
}
