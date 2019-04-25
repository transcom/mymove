package publicapi

import (
	"fmt"
	"net/http/httptest"

	accesscodeops "github.com/transcom/mymove/pkg/gen/restapi/apioperations/accesscode"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	accesscodeservice "github.com/transcom/mymove/pkg/services/accesscode"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *HandlerSuite) TestValidateAccessCodeHandler_Valid() {
	// create user
	user := testdatagen.MakeDefaultUser(suite.DB())
	selectedMoveType := models.SelectedMoveTypePPM

	// creates access code
	code := "TEST1"
	accessCode := models.AccessCode{
		Code:     code,
		MoveType: &selectedMoveType,
	}
	fullCode := fmt.Sprintf("%s-%s", selectedMoveType, code)
	suite.MustSave(&accessCode)
	// makes request
	request := httptest.NewRequest("GET", fmt.Sprintf("/accesscode/valid?code=%s", fullCode), nil)
	request = suite.AuthenticateUserRequest(request, user)

	params := accesscodeops.ValidateAccessCodeParams{
		HTTPRequest: request,
		Code:        &fullCode,
	}

	context := handlers.NewHandlerContext(suite.DB(), suite.TestLogger())
	accessCodeValidator := accesscodeservice.NewAccessCodeValidator(suite.DB())

	handler := ValidateAccessCodeHandler{context, accessCodeValidator}
	response := handler.Handle(params)

	suite.IsNotErrResponse(response)
	validateAccessCodeResponse := response.(*accesscodeops.ValidateAccessCodeOK)
	validateAccessCodePayload := validateAccessCodeResponse.Payload

	suite.NotNil(validateAccessCodePayload.AccessCode)
	suite.True(*validateAccessCodePayload.Valid)
	suite.Assertions.IsType(&accesscodeops.ValidateAccessCodeOK{}, response)
}

func (suite *HandlerSuite) TestValidateAccessCodeHandler_Invalid() {
	// create user
	user := testdatagen.MakeDefaultUser(suite.DB())
	selectedMoveType := models.SelectedMoveTypeHHG

	// creates access code
	code := "TEST2"
	accessCode := models.AccessCode{
		Code:     code,
		MoveType: &selectedMoveType,
		UserID:   &user.ID,
	}
	fullCode := fmt.Sprintf("%s-%s", selectedMoveType, code)
	suite.MustSave(&accessCode)
	// makes request
	request := httptest.NewRequest("GET", fmt.Sprintf("/accesscode/valid?code=%s", fullCode), nil)
	request = suite.AuthenticateUserRequest(request, user)

	params := accesscodeops.ValidateAccessCodeParams{
		HTTPRequest: request,
		Code:        &fullCode,
	}

	context := handlers.NewHandlerContext(suite.DB(), suite.TestLogger())
	accessCodeValidator := accesscodeservice.NewAccessCodeValidator(suite.DB())

	handler := ValidateAccessCodeHandler{context, accessCodeValidator}
	response := handler.Handle(params)

	suite.IsNotErrResponse(response)
	validateAccessCodeResponse := response.(*accesscodeops.ValidateAccessCodeOK)
	validateAccessCodePayload := validateAccessCodeResponse.Payload

	suite.False(*validateAccessCodePayload.Valid)
	suite.Assertions.IsType(&accesscodeops.ValidateAccessCodeOK{}, response)
}
