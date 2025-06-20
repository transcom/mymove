package internalapi

import (
	"database/sql"
	"fmt"
	"net/http/httptest"

	"github.com/go-openapi/strfmt"
	"github.com/stretchr/testify/mock"

	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/factory"
	postalcodesops "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/postal_codes"
	"github.com/transcom/mymove/pkg/services/mocks"
)

func (suite *HandlerSuite) TestValidatePostalCodeWithRateDataHandler() {
	suite.Run("Valid postal code", func() {
		user := factory.BuildUser(nil, nil, nil)

		postalCode := "30813"
		postalCodeTypeString := "origin"

		request := httptest.NewRequest("GET", fmt.Sprintf("/rate_engine_postal_codes/%s?postal_code_type=%s", postalCode, postalCodeTypeString), nil)
		request = suite.AuthenticateUserRequest(request, user)

		params := postalcodesops.ValidatePostalCodeWithRateDataParams{
			HTTPRequest:    request,
			PostalCode:     postalCode,
			PostalCodeType: postalCodeTypeString,
		}

		handlerConfig := suite.HandlerConfig()
		postalCodeValidator := &mocks.PostalCodeValidator{}
		postalCodeValidator.On("ValidatePostalCode",
			mock.AnythingOfType("*appcontext.appContext"),
			postalCode,
		).Return(true, nil)
		handler := ValidatePostalCodeWithRateDataHandler{handlerConfig, postalCodeValidator}

		// Validate incoming payload: no body to validate

		response := handler.Handle(params)

		suite.IsType(&postalcodesops.ValidatePostalCodeWithRateDataOK{}, response)
		validatePostalCodeResponse := response.(*postalcodesops.ValidatePostalCodeWithRateDataOK)
		validatePostalCodePayload := validatePostalCodeResponse.Payload

		// Validate outgoing payload
		suite.NoError(validatePostalCodePayload.Validate(strfmt.Default))

		suite.NotNil(validatePostalCodePayload.PostalCode)
		suite.NotNil(validatePostalCodePayload.PostalCodeType)
		suite.True(*validatePostalCodePayload.Valid)
	})

	suite.Run("Invalid postal code", func() {
		user := factory.BuildUser(nil, nil, nil)

		postalCode := "00988"
		postalCodeTypeString := "destination"

		request := httptest.NewRequest("GET", fmt.Sprintf("/rate_engine_postal_codes/%s?postal_code_type=%s", postalCode, postalCodeTypeString), nil)
		request = suite.AuthenticateUserRequest(request, user)

		params := postalcodesops.ValidatePostalCodeWithRateDataParams{
			HTTPRequest:    request,
			PostalCode:     postalCode,
			PostalCodeType: postalCodeTypeString,
		}

		handlerConfig := suite.HandlerConfig()
		postalCodeValidator := &mocks.PostalCodeValidator{}
		postalCodeValidator.On("ValidatePostalCode",
			mock.AnythingOfType("*appcontext.appContext"),
			postalCode,
		).Return(false, apperror.NewUnsupportedPostalCodeError(postalCode, "bad postal code"))
		handler := ValidatePostalCodeWithRateDataHandler{handlerConfig, postalCodeValidator}

		// Validate incoming payload: no body to validate

		response := handler.Handle(params)

		suite.IsType(&postalcodesops.ValidatePostalCodeWithRateDataOK{}, response)
		validatePostalCodeResponse := response.(*postalcodesops.ValidatePostalCodeWithRateDataOK)
		validatePostalCodePayload := validatePostalCodeResponse.Payload

		// Validate outgoing payload
		suite.NoError(validatePostalCodePayload.Validate(strfmt.Default))

		suite.NotNil(validatePostalCodePayload.PostalCode)
		suite.NotNil(validatePostalCodePayload.PostalCodeType)
		suite.False(*validatePostalCodePayload.Valid)
	})

	suite.Run("Database error", func() {
		user := factory.BuildUser(nil, nil, nil)

		postalCode := "30813"
		postalCodeTypeString := "destination"

		request := httptest.NewRequest("GET", fmt.Sprintf("/rate_engine_postal_codes/%s?postal_code_type=%s", postalCode, postalCodeTypeString), nil)
		request = suite.AuthenticateUserRequest(request, user)

		params := postalcodesops.ValidatePostalCodeWithRateDataParams{
			HTTPRequest:    request,
			PostalCode:     postalCode,
			PostalCodeType: postalCodeTypeString,
		}

		handlerConfig := suite.HandlerConfig()
		postalCodeValidator := &mocks.PostalCodeValidator{}
		postalCodeValidator.On("ValidatePostalCode",
			mock.AnythingOfType("*appcontext.appContext"),
			postalCode,
		).Return(false, sql.ErrNoRows)
		handler := ValidatePostalCodeWithRateDataHandler{handlerConfig, postalCodeValidator}

		// Validate incoming payload: no body to validate

		response := handler.Handle(params)
		suite.IsType(&postalcodesops.ValidatePostalCodeWithRateDataBadRequest{}, response)

		// Validate outgoing payload: no payload
	})
}
